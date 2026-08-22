package main

import (
	"fmt"
	"sync"
)

// opsTransitionTable declares every legal status transition. Each status may
// move to itself as a no-op (the resume button re-confirming the current state
// must not be rejected), and paused plans can resume back to active so a paused
// observation can continue. Closed is terminal and has no outgoing moves.
var opsTransitionTable = map[OpsStatus]map[OpsStatus]bool{
	OpsStatusQueued: {OpsStatusQueued: true, OpsStatusActive: true, OpsStatusClosed: true},
	OpsStatusActive: {OpsStatusActive: true, OpsStatusPaused: true, OpsStatusClosed: true},
	OpsStatusPaused: {OpsStatusPaused: true, OpsStatusActive: true, OpsStatusClosed: true},
	OpsStatusClosed: {},
}

type OpsTransition struct {
	From   OpsStatus `json:"from"`
	To     OpsStatus `json:"to"`
	Reason string    `json:"reason"`
}
type OpsStateMachine struct {
	mu      sync.RWMutex
	history []OpsTransition
}

func newOpsStateMachine() *OpsStateMachine { return &OpsStateMachine{history: []OpsTransition{}} }
func (m *OpsStateMachine) CanMove(from, to OpsStatus) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return opsTransitionTable[from][to]
}
func (m *OpsStateMachine) Move(from, to OpsStatus, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// A same-status move is a no-op confirmation: it is always permitted but must
	// never be recorded, otherwise the status history accumulates spurious entries
	// for actions that changed nothing.
	if from == to {
		return nil
	}
	if !opsTransitionTable[from][to] {
		return fmt.Errorf("%w: %s to %s", ErrOpsTransition, from, to)
	}
	m.history = append(m.history, OpsTransition{From: from, To: to, Reason: reason})
	return nil
}
func (m *OpsStateMachine) History() []OpsTransition {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]OpsTransition(nil), m.history...)
}
func (m *OpsStateMachine) Last() (OpsTransition, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.history) == 0 {
		return OpsTransition{}, false
	}
	return m.history[len(m.history)-1], true
}
func (m *OpsStateMachine) Reset() { m.mu.Lock(); defer m.mu.Unlock(); m.history = m.history[:0] }
func opsStatusValid(value OpsStatus) bool {
	return value == OpsStatusQueued || value == OpsStatusActive || value == OpsStatusPaused || value == OpsStatusClosed
}
func opsStatusTerminal(value OpsStatus) bool { return value == OpsStatusClosed }
