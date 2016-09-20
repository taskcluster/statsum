package aggregator

import (
	"fmt"
	"sync"

	"github.com/taskcluster/statsum/payload"
)

// Aggregator holds all Projects
type Aggregator struct {
	m        sync.RWMutex
	projects map[string]*Project
}

// NewAggregator creates a new aggregator
func NewAggregator() *Aggregator {
	return &Aggregator{
		projects: make(map[string]*Project),
	}
}

// Process aggreates a payload
func (a *Aggregator) Process(project string, payload *payload.Payload) {
	// Lock the aggregator for reading, this works if project already exists
	a.m.RLock()
	p := a.projects[project]
	if p != nil {
		p.process(payload)
		a.m.RUnlock()
		return // Return once processed
	}
	a.m.RUnlock()

	// Lock aggregator for write-access
	a.m.Lock()
	p = a.projects[project]
	if p == nil {
		p = newProject()
		a.projects[project] = p
	}
	p.process(payload)
	a.m.Unlock()
}

// A MetricHandler takes care of sending metrics to somewhere
type MetricHandler func(name string, value float64)

// SendMetrics passes all metrics to MetricHandler
func (a *Aggregator) SendMetrics(handler MetricHandler) {
	projects := a.ProjectNames()

	for _, project := range projects {
		a.m.RLock()
		p := a.projects[project]
		p.sendMetrics(project, handler)
		a.m.RUnlock()
	}
}

// SendHourlyMetrics passes all metrics to MetricHandler
func (a *Aggregator) SendHourlyMetrics(handler MetricHandler) {
	projects := a.ProjectNames()

	for _, project := range projects {
		// At end of hour we swap out the project and replace it
		a.m.Lock()
		p := a.projects[project]
		if p != nil {
			a.projects[project] = newProject()
		}
		a.m.Unlock()
		p.sendHourlyMetrics(project, handler)
	}
}

// ProjectNames returns a list of projects that currently exist
func (a *Aggregator) ProjectNames() []string {
	a.m.RLock()
	defer a.m.RUnlock()

	projects := make([]string, 0, len(a.projects))
	for project := range a.projects {
		projects = append(projects, project)
	}
	return projects
}

// PrintHealthMetrics renders some simple health metrics
func (a *Aggregator) PrintHealthMetrics() {
	a.m.RLock()
	defer a.m.RUnlock()
	for project, p := range a.projects {
		p.RLock()
		defer p.RUnlock()
		counters := len(p.counters)
		measures := len(p.measures)
		fmt.Println("Project: ", project,
			" counter metrics: ", counters,
			" value metrics: ", measures)
	}
}
