//go:generate msgp.v0
//go:generate ffjson $GOFILE

package payload

// A CountMetric is a metric that is to be aggregated by simple summation.
type CountMetric struct {
	Key   string  `json:"k" msg:"k"`
	Value float64 `json:"v" msg:"v"`
}

// A ValueMetric is metric that is to be aggregated as a series of discrete
// values, by estimation of percentiles.
type ValueMetric struct {
	Key   string    `json:"k" msg:"k"`
	Value []float64 `json:"v" msg:"v"`
}

// A Payload is the payload of a single request.
type Payload struct {
	CountMetrics []CountMetric `json:"countMetrics" msg:"countMetrics"`
	ValueMetrics []ValueMetric `json:"valueMetrics" msg:"valueMetrics"`
}
