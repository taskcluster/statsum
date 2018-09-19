package server

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/pkg/errors"
	got "github.com/taskcluster/go-got"
)

const signalfxEndPoint = "https://ingest.signalfx.com/v2/datapoint"

type datapoint struct {
	Metric    string  `json:"metric"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

type signalfxDrain struct {
	token  string
	g      *got.Got
	points []datapoint
}

func newSignalfxDrain(token string) *signalfxDrain {
	return &signalfxDrain{
		token: token,
		g:     got.New(),
	}
}

func (s *signalfxDrain) Name() string {
	return "signalfx"
}

func (s *signalfxDrain) Send(name string, value float64, time time.Time) {
	if math.IsNaN(value) {
		fmt.Printf("Ignoring NaN value for metric %s\n", name)
		return
	}
	s.points = append(s.points, datapoint{
		Metric:    name,
		Value:     value,
		Timestamp: time.Unix() * 1000,
	})
}

func (s *signalfxDrain) Flush() error {
	req := s.g.Post(signalfxEndPoint, nil)
	err := req.JSON(map[string][]datapoint{
		"gauge": s.points,
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal payload to JSON"))
	}

	req.Header.Set("X-SF-Token", s.token)
	res, err := req.Send()
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Wrong status code from signalfx: %d", res.StatusCode)
	}

	s.points = []datapoint{} // release accumulated memory

	return nil
}
