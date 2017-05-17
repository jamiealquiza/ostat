package main

import (
	"encoding/json"
	"sync"
)

// Metrics reads from inputs and writes
// output to data.
type Metrics struct {
	sync.Mutex
	data   map[string]map[string]interface{}
	inputs map[string]func() interface{}
}

func (m *Metrics) fetchMetrics() {
	m.Lock()
	defer m.Unlock()

	for c, f := range m.inputs {

		m.data[hostname][c] = f()
	}
}

func (m *Metrics) getMetrics() ([]byte, error) {
	m.Lock()
	defer m.Unlock()

	r, err := json.Marshal(m.data)
	r = append(r, 10)
	return r, err
}

func (m *Metrics) registerInput(c string, f func() interface{}) {
	m.Lock()
	defer m.Unlock()

	m.inputs[c] = f
}
