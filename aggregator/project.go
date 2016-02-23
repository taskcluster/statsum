package aggregator

import (
	"sync"

	"github.com/jonasfj/statsum/payload"
)

// A Project aggregates metrics for a entire project.
type Project struct {
	sync.RWMutex
	counterMetrics map[string]*counterMetric
	valueMetrics   map[string]*valueMetric
}

func newProject() *Project {
	return &Project{
		counterMetrics: make(map[string]*counterMetric),
		valueMetrics:   make(map[string]*valueMetric),
	}
}

// ProcessCounts aggregates counts
func (p *Project) processCounts(counts []payload.CountMetric) {
	// Lock the project for reading, we can process an entry if it already exists
	// so this will be the default mode when entering the loop below
	p.RLock()
	defer p.RUnlock()
	for _, count := range counts {
		counterMetric := p.counterMetrics[count.Key]

		if counterMetric != nil {
			// If the metric exists this is as simple as aggregating it
			counterMetric.process(count.Value)
		} else {
			// Get a lock for writing
			p.RUnlock()
			p.Lock()

			// Check that metric wasn't added while we got the write lock
			counterMetric := p.counterMetrics[count.Key]
			if counterMetric == nil {
				counterMetric = newCounterMetric()
				p.counterMetrics[count.Key] = counterMetric
			}
			counterMetric.process(count.Value)

			p.Unlock()
			p.RLock() // Ensure that we lock for reading before we continue
		}
	}
}

// ProcessValues aggregates values
func (p *Project) processValues(values []payload.ValueMetric) {
	// Lock the project for reading, we can process an entry if it already exists
	// so this will be the default mode when entering the loop below
	p.RLock()
	defer p.RUnlock()
	for _, value := range values {
		valueMetric := p.valueMetrics[value.Key]

		if valueMetric != nil {
			// If the metric exists this is as simple as aggregating it
			valueMetric.process(value.Value)
		} else {
			// Get a lock for writing
			p.RUnlock()
			p.Lock()

			// Check that metric wasn't added while we got the write lock
			valueMetric := p.valueMetrics[value.Key]
			if valueMetric == nil {
				valueMetric = newValueMetric()
				p.valueMetrics[value.Key] = valueMetric
			}
			valueMetric.process(value.Value)

			p.Unlock()
			p.RLock() // Ensure that we lock for reading before we continue
		}
	}
}

// Process will aggregate all metrics from a payload
func (p *Project) process(payload *payload.Payload) {
	p.processCounts(payload.CountMetrics)
	p.processValues(payload.ValueMetrics)
}

// SendMetrics metrics to hander under prefix
func (p *Project) sendMetrics(prefix string, handler MetricHandler) {
	p.RLock()
	defer p.RUnlock()
	for name, valueMetric := range p.valueMetrics {
		valueMetric.sendMetrics(prefix, name, handler)
	}
	for name, counterMetric := range p.counterMetrics {
		counterMetric.sendMetrics(prefix, name, handler)
	}
}

// SendMetrics metrics to hander under prefix
func (p *Project) sendHourlyMetrics(prefix string, handler MetricHandler) {
	p.RLock()
	defer p.RUnlock()
	for name, valueMetric := range p.valueMetrics {
		valueMetric.sendHourlyMetrics(prefix, name, handler)
	}
	for name, counterMetric := range p.counterMetrics {
		counterMetric.sendHourlyMetrics(prefix, name, handler)
	}
}
