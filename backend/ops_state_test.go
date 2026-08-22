package main

import "testing"

func TestStatePausedResumes(t *testing.T) {
	m := newOpsStateMachine()
	if !m.CanMove(OpsStatusPaused, OpsStatusActive) {
		t.Fatalf("paused must be able to resume to active")
	}
}

func TestStateNoopLeavesEmpty(t *testing.T) {
	m := newOpsStateMachine()
	if err := m.Move(OpsStatusActive, OpsStatusActive, "noop"); err != nil {
		t.Fatal(err)
	}
	if history := m.History(); len(history) != 0 {
		t.Fatalf("no-op move recorded history: %v", history)
	}
}

func TestStatePausedIsValid(t *testing.T) {
	if !opsStatusValid(OpsStatusPaused) {
		t.Fatalf("paused must be a valid status")
	}
}

func TestStateRejectsIllegalMove(t *testing.T) {
	m := newOpsStateMachine()
	if err := m.Move(OpsStatusClosed, OpsStatusActive, "reopen"); err == nil {
		t.Fatalf("closed -> active must be rejected")
	}
}

func TestStateSameMoveAllowed(t *testing.T) {
	m := newOpsStateMachine()
	if !m.CanMove(OpsStatusActive, OpsStatusActive) {
		t.Fatalf("same-status move must be allowed")
	}
}

func TestStateLastEmptySafe(t *testing.T) {
	m := newOpsStateMachine()
	if _, ok := m.Last(); ok {
		t.Fatalf("empty machine must not report a last transition")
	}
}
