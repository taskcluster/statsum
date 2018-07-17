//go:generate msgp

package payload

// A Counter is a metric that is to be aggregated by simple summation.
type Counter struct {
	Key   string  `json:"k" msg:"k"`
	Value float64 `json:"v" msg:"v"`
}

// A Measure is metric that is to be aggregated as a series of discrete
// values, by estimation of percentiles.
type Measure struct {
	Key   string    `json:"k" msg:"k"`
	Value []float64 `json:"v" msg:"v"`
}

// A Payload is the payload of a single request.
type Payload struct {
	Counters []Counter `json:"counters" msg:"counters"`
	Measures []Measure `json:"measures" msg:"measures"`
}
