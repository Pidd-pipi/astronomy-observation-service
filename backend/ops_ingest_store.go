package main

import (
	"context"
	"sync"
)

// IngestStore stages ingested observation plans.
type IngestStore struct {
	mu    sync.Mutex
	items map[string]string
}

func newIngestStore() *IngestStore { return &IngestStore{items: map[string]string{}} }

func (s *IngestStore) Put(ctx context.Context, id, plan string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[id] = plan
	return nil
}

func (s *IngestStore) Get(ctx context.Context, id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.items[id]
	if !ok {
		return "", ErrOpsNotFound
	}
	return v, nil
}

func (s *IngestStore) Count() int { s.mu.Lock(); defer s.mu.Unlock(); return len(s.items) }
