package aggregator

import (
	"sync"

	"github.com/taskcluster/statsum/payload"
)

// A Project aggregates metrics for a entire project.
type Project struct {
	sync.RWMutex
	counters map[string]*counter
	measures map[string]*measure
}

func newProject() *Project {
	return &Project{
		counters: make(map[string]*counter),
		measures: make(map[string]*measure),
	}
}

// processCounters aggregates counters
func (p *Project) processCounters(counters []payload.Counter) {
	// Lock the project for reading, we can process an entry if it already exists
	// so this will be the default mode when entering the loop below
	p.RLock()
	defer p.RUnlock()
	for _, c := range counters {
		counter := p.counters[c.Key]

		if counter != nil {
			// If the metric exists this is as simple as aggregating it
			counter.process(c.Value)
		} else {
			// Get a lock for writing
			p.RUnlock()
			p.Lock()

			// Check that metric wasn't added while we got the write lock
			counter := p.counters[c.Key]
			if counter == nil {
				counter = newCounter()
				p.counters[c.Key] = counter
			}
			counter.process(c.Value)

			p.Unlock()
			p.RLock() // Ensure that we lock for reading before we continue
		}
	}
}

// processMeasures aggregates values
func (p *Project) processMeasures(measures []payload.Measure) {
	// Lock the project for reading, we can process an entry if it already exists
	// so this will be the default mode when entering the loop below
	p.RLock()
	defer p.RUnlock()
	for _, m := range measures {
		measure := p.measures[m.Key]

		if measure != nil {
			// If the metric exists this is as simple as aggregating it
			measure.process(m.Value)
		} else {
			// Get a lock for writing
			p.RUnlock()
			p.Lock()

			// Check that metric wasn't added while we got the write lock
			measure := p.measures[m.Key]
			if measure == nil {
				measure = newMeasure()
				p.measures[m.Key] = measure
			}
			measure.process(m.Value)

			p.Unlock()
			p.RLock() // Ensure that we lock for reading before we continue
		}
	}
}

// Process will aggregate all metrics from a payload
func (p *Project) process(payload *payload.Payload) {
	p.processCounters(payload.Counters)
	p.processMeasures(payload.Measures)
}

// SendMetrics metrics to hander under prefix
func (p *Project) sendMetrics(prefix string, handler MetricHandler) {
	p.RLock()
	defer p.RUnlock()
	for name, measure := range p.measures {
		measure.sendMetrics(prefix, name, handler)
	}
	for name, counter := range p.counters {
		counter.sendMetrics(prefix, name, handler)
	}
}

// SendMetrics metrics to hander under prefix
func (p *Project) sendHourlyMetrics(prefix string, handler MetricHandler) {
	p.RLock()
	defer p.RUnlock()
	for name, measure := range p.measures {
		measure.sendHourlyMetrics(prefix, name, handler)
	}
	for name, counter := range p.counters {
		counter.sendHourlyMetrics(prefix, name, handler)
	}
}
