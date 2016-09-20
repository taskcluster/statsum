package server

import (
	"time"

	"golang.org/x/net/context"

	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/golib/sfxclient"
)

type signalfxDrain struct {
	token  string
	client *sfxclient.HTTPDatapointSink
	points []*datapoint.Datapoint
}

func newSignalfxDrain(token string) *signalfxDrain {
	client := sfxclient.NewHTTPDatapointSink()
	client.AuthToken = token
	return &signalfxDrain{
		token:  token,
		client: client,
	}
}

func (s *signalfxDrain) Name() string {
	return "signalfx"
}

func (s *signalfxDrain) Send(name string, value float64, time time.Time) {
	dataPoint := datapoint.New(name, nil, datapoint.NewFloatValue(value), datapoint.Gauge, time)
	s.points = append(s.points, dataPoint)
}

func (s *signalfxDrain) Flush() error {
	ctx := context.Background()
	points := s.points
	s.points = nil
	return s.client.AddDatapoints(ctx, points)
}
