// Package main contains the main function which loads configuration,
// initializes and starts the server.
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"gopkg.in/dgrijalva/jwt-go.v2"

	"github.com/docopt/docopt-go"
	"github.com/jonasfj/statsum/payload"
	"github.com/jonasfj/statsum/server"
)

const usage = `statsum

Usage:
  statsum server
  statsum test <duration> -t <threads> -c <counters> -v <values> -p <points> -d <dimensionality> [-j]
  statsum -h | --help
  statsum --version

Options:
  -h, --help                Show this screen.
  --version                 Show version.
  -t, --threads <threads>   Number of threads sending requests.
  -c, --counters <counters> Number of counters per request.
  -v, --values <values>     Number of values per request.
  -p, --points <points>     Number of points per value.
  -d, --dim <dim>           Dimensionality of metrics.
  -j, --json                Send requests in JSON`

func main() {
	secret := []byte(os.Getenv("JWT_SECRET_KEY"))
	args, _ := docopt.Parse(usage, nil, true, "statsum 0.1.0", false)

	switch {
	case args["server"].(bool):
		s, err := server.New(server.Config{
			Port:           os.Getenv("PORT"),
			TLSCertificate: os.Getenv("TLS_CERTIFICATE"),
			TLSKey:         os.Getenv("TLS_KEY"),
			JwtSecret:      secret,
			SignalFxToken:  os.Getenv("SIGNALFX_TOKEN"),
			DatadogAPIKey:  os.Getenv("DATADOG_API_KEY"),
			DatadogAppKey:  os.Getenv("DATADOG_APP_KEY"),
			SentryDSN:      os.Getenv("SENTRY_DSN"),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting server")
		err = s.Start()
		if err != nil {
			panic(err)
		}

	case args["test"].(bool):
		duration, err := strconv.Atoi(args["<duration>"].(string))
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		threads, err := strconv.Atoi(args["--threads"].(string))
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		counters, err := strconv.Atoi(args["--counters"].(string))
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		values, err := strconv.Atoi(args["--values"].(string))
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		points, err := strconv.Atoi(args["--points"].(string))
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		dim, err := strconv.Atoi(args["--dim"].(string))
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		json := args["--json"].(bool)
		format := "msgpack"
		if json {
			format = "json"
		}

		baseURL := os.Getenv("PUBLIC_URL")
		if baseURL == "" {
			fmt.Println("Missing PUBLIC_URL")
			os.Exit(1)
		}

		fmt.Println("Running tests with:")
		fmt.Println("  duration:       ", duration, " seconds")
		fmt.Println("  threads:        ", threads)
		fmt.Println("  counters:       ", counters)
		fmt.Println("  values:         ", values)
		fmt.Println("  points:         ", points)
		fmt.Println("  dimensionality: ", dim)
		fmt.Println("  format:         ", format)

		errCount := int32(0)
		failCount := int32(0)
		readFailCount := int32(0)
		okCount := int32(0)

		for t := 0; t < threads; t++ {
			go func(t int) {
				// Setup token
				token := jwt.New(jwt.SigningMethodHS256)
				token.Claims["project"] = "my-super-awesome-test-project"
				token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
				jwt, err := token.SignedString(secret)
				if err != nil {
					panic(err)
				}
				auth := "Bearer " + jwt

				// Create url
				url := baseURL + "/v1/project/my-super-awesome-test-project"

				// Create payload
				i := int((counters / threads) * t)
				j := int((values / threads) * t)
				newPayload := func() *payload.Payload {
					p := payload.Payload{}
					for c := 0; c < counters; c++ {
						p.Counters = append(p.Counters, payload.Counter{
							Key:   fmt.Sprint(i, "-my-test-counter-", i),
							Value: math.Floor(rand.Float64() * 100),
						})
						i = (i + 1) % dim
					}
					for v := 0; v < values; v++ {
						data := make([]float64, points)
						for x := 0; x < points; x++ {
							data[x] = rand.Float64() * 30000
						}
						p.Measures = append(p.Measures, payload.Measure{
							Key:   fmt.Sprint(j, "-my-test-value-", j),
							Value: data,
						})
						j = (j + 1) % dim
					}
					return &p
				}

				// Encode payload
				contentType := "application/x-msgpack"
				newBody := func() io.Reader {
					payload := newPayload()
					b, err := payload.MarshalMsg([]byte{})
					if err != nil {
						panic(err)
					}
					return bytes.NewReader(b)
				}
				if json {
					contentType = "application/json"
					newBody = func() io.Reader {
						payload := newPayload()
						b, err := payload.MarshalJSON()
						if err != nil {
							panic(err)
						}
						return bytes.NewReader(b)
					}
				}

				// Create http client
				client := http.Client{
					Timeout: 25 * time.Second,
				}

				for {
					r, err := http.NewRequest("POST", url, newBody())
					if err != nil {
						panic(err)
					}
					r.Header.Set("Content-Type", contentType)
					r.Header.Set("Authorization", auth)
					res, err := client.Do(r)
					if err != nil {
						atomic.AddInt32(&errCount, 1)
						continue
					}
					if res.StatusCode != http.StatusOK {
						atomic.AddInt32(&failCount, 1)
					} else {
						atomic.AddInt32(&okCount, 1)
					}
					_, err = ioutil.ReadAll(res.Body)
					res.Body.Close()
					if err != nil {
						atomic.AddInt32(&readFailCount, 1)
					}
				}
			}(t)
		}

		time.Sleep(time.Duration(duration) * time.Second)

		// Print summary
		ok := atomic.LoadInt32(&okCount)
		fmt.Println("errCount:       ", atomic.LoadInt32(&errCount))
		fmt.Println("failCount:      ", atomic.LoadInt32(&failCount))
		fmt.Println("readFailCount:  ", atomic.LoadInt32(&readFailCount))
		fmt.Println("okCount:        ", ok)
		fmt.Println("req/s:          ", float64(ok)/float64(duration))
		os.Exit(0)

	default:
		fmt.Println(usage)
		os.Exit(1)
	}
}
