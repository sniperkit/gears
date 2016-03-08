package fennel

import (
	"sync"
	"time"
)

// Datapoint is an interface representing one measurement or sample
type Datapoint interface {
	MetricName() string
	Tags() map[string]interface{}
	Value() interface{}
	Unit() string
	Time() time.Time
}

// NewDatapoint returns a datapoint for a metric with a given value, time, and tags.
func NewDatapoint(m *Metric, value interface{}, t time.Time, tags map[string]interface{}) Datapoint {
	return &datapoint{
		Metric: m,
		value:  value,
		t:      t,
		tags:   tags,
	}
}

// Metric is a struct with fields that
// have the same value in every Datapoint
// within a measurement/event stream
type Metric struct {
	name string
	unit string
}

// NewMetric returns a metric which can be used for creating datapoints.
func NewMetric(name, unit string) *Metric {
	return &Metric{
		name: name,
		unit: unit,
	}
}

type datapoint struct {
	*Metric
	l     sync.RWMutex
	value interface{}
	t     time.Time
	tags  map[string]interface{}
}

func (dp *datapoint) MetricName() string {
	return dp.name
}

func (dp *datapoint) Tags() map[string]interface{} {
	c := map[string]interface{}{}
	dp.l.RLock()

	// returns a copy of the tags
	// instead the map itself
	for k, v := range dp.tags {
		c[k] = v
	}

	dp.l.RUnlock()
	return c
}

func (dp *datapoint) Value() interface{} {
	return dp.value
}

func (dp *datapoint) Unit() string {
	return dp.unit
}

func (dp *datapoint) Time() time.Time {
	return dp.t
}
