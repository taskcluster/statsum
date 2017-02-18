package statsum

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	got "github.com/taskcluster/go-got"
	"github.com/taskcluster/statsum/payload"

	"github.com/pborman/uuid"
	jwt "gopkg.in/dgrijalva/jwt-go.v2"
)

// A Configurer is a function that given a project name returns Config options
// required for submitting metrics.
type Configurer func(project string) (Config, error)

// Config required for statsum client to submit metrics.
type Config struct {
	Project string
	BaseURL string
	Token   string
	expires time.Time
}

// StaticConfigurer computes a JWT that is valid for 25 minutes, and rotated
// every 10 minutes. This allows for 15 minutes clock drift, while ensuring that
// all credentials being transmitted are temporary.
func StaticConfigurer(baseURL string, secret []byte) Configurer {
	return func(project string) (Config, error) {
		now := time.Now()
		token := jwt.New(jwt.SigningMethodHS256)
		token.Claims["project"] = project
		token.Claims["exp"] = now.Add(25 * time.Minute).Unix()
		signedToken, err := token.SignedString(secret)
		if err != nil {
			panic(fmt.Errorf("Unable to sign JWT token, error: %s", err))
		}
		return Config{
			Project: project,
			BaseURL: baseURL,
			Token:   signedToken,
			expires: now.Add(10 * time.Minute),
		}, nil
	}
}

type cache struct {
	project    string
	options    Options
	configurer Configurer
	mConfig    sync.Mutex
	config     *Config
	g          *got.Got
	mData      sync.Mutex
	counters   map[string]float64
	measures   map[string][]float64
	dataPoints int
	timer      *time.Timer
}

// A ServerError represents a server-side error from statsum.
type ServerError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (s *ServerError) Error() string {
	return fmt.Sprintf("%s: %s", s.Code, s.Message)
}

// Statsum client for collecting metrics
type Statsum struct {
	prefix string
	cache  *cache
}

// Options for Statsum client.
type Options struct {
	MaxDataPoints int           // Maximum data points before sending, default 10k
	MaxDelay      time.Duration // Maximum delay before sending, default 90s
	MinDelay      time.Duration // Minimum delay before sending, default 30s
	OnError       func(error)   // Error callback, ignored by default
}

// New returns a new Statum client.
func New(project string, configurer Configurer, options Options) *Statsum {
	if options.MaxDataPoints == 0 {
		options.MaxDataPoints = 10000
	}
	if options.MaxDelay == 0 {
		options.MaxDelay = 90 * time.Second
	}
	if options.MinDelay == 0 {
		options.MinDelay = 30 * time.Second
	}
	g := got.New()
	g.Client = &http.Client{Timeout: 45 * time.Second}
	g.Retries = 7
	g.MaxSize = 3 * 1024 * 1024 // 30 MiB
	g.BackOff = &got.BackOff{
		DelayFactor:         500 * time.Millisecond,
		RandomizationFactor: 0.25,
		MaxDelay:            120 * time.Second,
	}
	return &Statsum{
		prefix: "",
		cache: &cache{
			project:    project,
			options:    options,
			configurer: configurer,
			g:          g,
		},
	}
}

// Measure one or more values for the given name
func (s *Statsum) Measure(name string, values ...float64) {
	// Lock cache data
	s.cache.mData.Lock()
	defer s.cache.mData.Unlock()

	// Ensure map exists, and insert data
	if s.cache.measures == nil {
		s.cache.measures = make(map[string][]float64)
	}
	key := s.prefix + name
	s.cache.measures[key] = append(s.cache.measures[key], values...)

	// Increment data points, check we don't have too many
	s.cache.dataPoints += len(values)
	if s.cache.dataPoints >= s.cache.options.MaxDataPoints {
		s.taskAndFlush()
	} else {
		s.ensureTimeout()
	}
}

// Count increments counter by name with given value
func (s *Statsum) Count(name string, value float64) {
	// Lock cache data
	s.cache.mData.Lock()
	defer s.cache.mData.Unlock()

	// Ensure map exists, and insert data
	if s.cache.counters == nil {
		s.cache.counters = make(map[string]float64)
	}
	key := s.prefix + name
	s.cache.counters[key] += value

	// Increment data points, check we don't have too many
	s.cache.dataPoints++
	if s.cache.dataPoints >= s.cache.options.MaxDataPoints {
		s.taskAndFlush()
	} else {
		s.ensureTimeout()
	}
}

// WithPrefix returns a Statsum client object that prefixes all values with
// the given prefix.
//
// The new Statsum client object uses the same underlying cache as its parent,
// making this a cheap operation that does not number of increase submission
// requests, or require additional Flush calls.
func (s *Statsum) WithPrefix(prefix string) *Statsum {
	return &Statsum{
		prefix: s.prefix + prefix + ".",
		cache:  s.cache,
	}
}

// Flush will send all metrics immediately,
func (s *Statsum) Flush() error {
	// Lock cache data
	s.cache.mData.Lock()

	// Take data
	counters := s.cache.counters
	measures := s.cache.measures
	s.cache.counters = nil
	s.cache.measures = nil
	s.cache.dataPoints = 0

	// Cancel timeouts
	if s.cache.timer != nil {
		s.cache.timer.Stop()
		s.cache.timer = nil
	}

	// Unlock so others can continue
	s.cache.mData.Unlock()

	// If we have no metrics we're done
	if counters == nil && measures == nil {
		return nil
	}

	return s.cache.submit(counters, measures)
}

// takeAndFlush takes the metrics and flushes later, this requires that caller
// has the mData lock!
func (s *Statsum) taskAndFlush() {
	// Take data
	counters := s.cache.counters
	measures := s.cache.measures
	s.cache.counters = nil
	s.cache.measures = nil
	s.cache.dataPoints = 0

	// Cancel timeouts
	if s.cache.timer != nil {
		s.cache.timer.Stop()
		s.cache.timer = nil
	}

	go func() {
		err := s.cache.submit(counters, measures)
		if err != nil && s.cache.options.OnError != nil {
			s.cache.options.OnError(err)
		}
	}()
}

// ensureTimeout sets a timeout for submit, if there isn't already one there.
// This requires that caller has the mData lock!
func (s *Statsum) ensureTimeout() {
	if s.cache.timer != nil {
		return
	}

	interval := (s.cache.options.MaxDelay - s.cache.options.MinDelay).Seconds()
	duration := s.cache.options.MinDelay + time.Duration(interval*rand.Float64())*time.Second
	s.cache.timer = time.AfterFunc(duration, func() {
		s.cache.mData.Lock()
		s.taskAndFlush()
		s.cache.mData.Unlock()
	})
}

func (c *cache) fetchConfig() (*Config, error) {
	c.mConfig.Lock()
	defer c.mConfig.Unlock()

	// Update config if necessary
	if c.config == nil || c.config.expires.Before(time.Now()) {
		cfg, err := c.configurer(c.project)
		if err != nil {
			return nil, err
		}
		// Strip trailing slash if any
		if strings.HasSuffix(cfg.BaseURL, "/") {
			cfg.BaseURL = cfg.BaseURL[:len(cfg.BaseURL)-1]
		}
		// Store config for next time
		c.config = &cfg
	}

	return c.config, nil
}

func (c *cache) submit(counters map[string]float64, measures map[string][]float64) error {
	// Get a config object
	config, err := c.fetchConfig()
	if err != nil {
		return err
	}

	// Construct payload
	p := payload.Payload{
		Counters: make([]payload.Counter, 0, len(counters)),
		Measures: make([]payload.Measure, 0, len(measures)),
	}
	for name, count := range counters {
		p.Counters = append(p.Counters, payload.Counter{
			Key:   name,
			Value: count,
		})
	}
	for name, values := range measures {
		p.Measures = append(p.Measures, payload.Measure{
			Key:   name,
			Value: values,
		})
	}

	// Construct request
	request := c.g.Post(config.BaseURL+"/v1/project/"+config.Project, nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+config.Token)
	request.Header.Set("X-Statsum-Request-Id", uuid.NewRandom().String())
	if err = request.JSON(p); err != nil {
		panic(fmt.Sprintf("Internal error serializing JSON, error: %s", err))
	}

	// Send request
	response, err := request.Send()
	if err != nil {
		return fmt.Errorf("Failed to submit statsum metrics, error: %s", err)
	}

	// Handle server errors
	if response.StatusCode != http.StatusOK {
		// Parse error message, if any
		var serr ServerError
		if json.Unmarshal(response.Body, &serr) == nil {
			return &serr
		}
		// Some default error codes and messages
		if 400 <= response.StatusCode && response.StatusCode < 500 {
			serr.Code = "RequestError"
		} else if 500 <= response.StatusCode && response.StatusCode < 600 {
			serr.Code = "InternalServerError"
		} else {
			serr.Code = "UnknownError"
		}
		serr.Message = fmt.Sprintf("Received status code %d", response.StatusCode)
		return &serr
	}

	return nil
}
