package main

import (
	"errors"
	"sync"
)

var (
	ErrAlertSinkClosed  = errors.New("alert sink is closed")
	ErrAlertSinkBlocked = errors.New("alert sink rejects this alert")
)

// AlertSink is an in-memory alert output with commit/close semantics.
type AlertSink struct {
	mu        sync.Mutex
	lines     []string
	open      bool
	committed bool
}

func newAlertSink() *AlertSink { return &AlertSink{open: true} }

func (s *AlertSink) WriteLine(line string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.open {
		return ErrAlertSinkClosed
	}
	if line == "" {
		return ErrAlertSinkBlocked
	}
	s.lines = append(s.lines, line)
	return nil
}

func (s *AlertSink) Commit() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.open {
		return ErrAlertSinkClosed
	}
	s.committed = true
	return nil
}

func (s *AlertSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.open = false
	return nil
}

func (s *AlertSink) Open() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.open
}

func (s *AlertSink) Committed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.committed
}

func (s *AlertSink) Lines() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, len(s.lines))
	copy(out, s.lines)
	return out
}
