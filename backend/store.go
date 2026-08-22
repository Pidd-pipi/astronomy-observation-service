package main

import (
	"errors"
	"sync"
)

var errRunNotFound = errors.New("observation run not found")

type RunStore struct {
	mu    sync.RWMutex
	items map[string]ObservationRun
}

func newRunStore() *RunStore {
	return &RunStore{items: map[string]ObservationRun{"run-501": {ID: "run-501", Target: "M31", Instrument: "0.8m reflector", Quality: "photometric", Status: "planned", StartedAt: "2026-08-22T20:00:00Z"}, "run-502": {ID: "run-502", Target: "NGC 7000", Instrument: "wide-field camera", Quality: "spectroscopic", Status: "running", StartedAt: "2026-08-21T19:30:00Z"}}}
}
func (s *RunStore) list() []ObservationRun {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o := make([]ObservationRun, 0, len(s.items))
	for _, v := range s.items {
		o = append(o, v)
	}
	return o
}
func (s *RunStore) changeStatus(id, status string) (ObservationRun, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.items[id]
	if !ok {
		return ObservationRun{}, errRunNotFound
	}
	v.Status = status
	s.items[id] = v
	return v, nil
}
