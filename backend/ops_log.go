package main

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrNightLogClosed  = errors.New("night log is closed")
	ErrNightLogBlocked = errors.New("night log rejects this entry")
)

// NightLog is an append-only in-memory log for one observing night with
// commit/close semantics: entries may only be committed once, before close.
type NightLog struct {
	mu        sync.Mutex
	entries   []string
	open      bool
	committed bool
}

func newNightLog() *NightLog { return &NightLog{open: true} }

func (l *NightLog) Append(entry string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.open {
		return ErrNightLogClosed
	}
	if entry == "" || strings.Contains(entry, "blocked") {
		return ErrNightLogBlocked
	}
	l.entries = append(l.entries, entry)
	return nil
}

func (l *NightLog) Commit() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.open {
		return ErrNightLogClosed
	}
	l.committed = true
	return nil
}

func (l *NightLog) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.open = false
	return nil
}

func (l *NightLog) Open() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.open
}

func (l *NightLog) Committed() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.committed
}

func (l *NightLog) Entries() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.entries...)
}

// NightLogRegistry hands out one log per night.
type NightLogRegistry struct {
	mu   sync.Mutex
	logs map[string]*NightLog
}

func newNightLogRegistry() *NightLogRegistry {
	return &NightLogRegistry{logs: map[string]*NightLog{}}
}

func (r *NightLogRegistry) For(night string) *NightLog {
	r.mu.Lock()
	defer r.mu.Unlock()
	if log, ok := r.logs[night]; ok {
		return log
	}
	log := newNightLog()
	r.logs[night] = log
	return log
}

func (r *NightLogRegistry) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.logs)
}
