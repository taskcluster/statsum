package aggregator

import (
	"sync"

	"github.com/snaury/tdigest-go"
)

const tdigestCompressionLevel = 10

type valueMetric struct {
	m             sync.Mutex
	digest        *tdigest.MergingDigest
	hourlySummary []tdigest.Centroid
}

func newValueMetric() *valueMetric {
	return &valueMetric{
		digest: tdigest.New(tdigestCompressionLevel),
	}
}

func (v *valueMetric) process(values []float64) {
	v.m.Lock()
	defer v.m.Unlock()
	for _, value := range values {
		v.digest.Add(value, 1)
	}
}

func (v *valueMetric) sendMetrics(prefix string, name string, handler MetricHandler) {
	v.m.Lock()
	min := v.digest.Quantile(0)
	p25 := v.digest.Quantile(0.25)
	p50 := v.digest.Quantile(0.50)
	p75 := v.digest.Quantile(0.75)
	p95 := v.digest.Quantile(0.95)
	p99 := v.digest.Quantile(0.99)
	max := v.digest.Quantile(1)
	for _, c := range v.hourlySummary {
		v.digest.Merge(c)
	}
	v.hourlySummary = v.digest.Summary()
	v.digest = tdigest.New(tdigestCompressionLevel)
	v.m.Unlock()

	handler(prefix+"."+name+".5m.min", min)
	handler(prefix+"."+name+".5m.p25", p25)
	handler(prefix+"."+name+".5m.p50", p50)
	handler(prefix+"."+name+".5m.p75", p75)
	handler(prefix+"."+name+".5m.p95", p95)
	handler(prefix+"."+name+".5m.p99", p99)
	handler(prefix+"."+name+".5m.max", max)
}

func (v *valueMetric) sendHourlyMetrics(prefix string, name string, handler MetricHandler) {
	v.m.Lock()
	summary := v.hourlySummary
	v.hourlySummary = nil
	v.m.Unlock()

	digest := tdigest.New(tdigestCompressionLevel)
	for _, c := range summary {
		digest.Merge(c)
	}
	min := digest.Quantile(0)
	p25 := digest.Quantile(0.25)
	p50 := digest.Quantile(0.50)
	p75 := digest.Quantile(0.75)
	p95 := digest.Quantile(0.95)
	p99 := digest.Quantile(0.99)
	max := digest.Quantile(1)

	handler(prefix+"."+name+".1h.min", min)
	handler(prefix+"."+name+".1h.p25", p25)
	handler(prefix+"."+name+".1h.p50", p50)
	handler(prefix+"."+name+".1h.p75", p75)
	handler(prefix+"."+name+".1h.p95", p95)
	handler(prefix+"."+name+".1h.p99", p99)
	handler(prefix+"."+name+".1h.max", max)
}
