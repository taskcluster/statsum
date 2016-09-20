package server

import "time"

const maxDataPoints = 1000

type drain interface {
	Name() string
	Send(name string, value float64, time time.Time)
	Flush() error
}
