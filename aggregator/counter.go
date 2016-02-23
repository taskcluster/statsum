package aggregator

import "sync"

type counterMetric struct {
	m           sync.Mutex
	count       float64
	hourlyCount float64
}

func newCounterMetric() *counterMetric {
	return &counterMetric{}
}

func (c *counterMetric) process(count float64) {
	c.m.Lock()
	defer c.m.Unlock()
	c.count += count
}

func (c *counterMetric) sendMetrics(prefix string, name string, handler MetricHandler) {
	c.m.Lock()
	count := c.count
	c.hourlyCount += count
	c.count = 0
	c.m.Unlock()

	handler(prefix+"."+name+".5m.count", count)
}

func (c *counterMetric) sendHourlyMetrics(prefix string, name string, handler MetricHandler) {
	c.m.Lock()
	count := c.hourlyCount
	c.hourlyCount = 0
	c.m.Unlock()

	handler(prefix+"."+name+".1h.count", count)
}
