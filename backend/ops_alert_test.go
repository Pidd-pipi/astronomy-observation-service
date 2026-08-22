package main

import (
	"context"
	"errors"
	"testing"
)

func TestAlertNotifyKeepsChain(t *testing.T) {
	sink := newAlertSink()
	dispatcher := newAlertDispatcher(sink)
	if err := sink.Close(); err != nil {
		t.Fatal(err)
	}
	err := dispatcher.Notify(context.Background(), Alert{ID: "a1", Severity: OpsPriorityHigh, Message: "seeing"})
	if err == nil {
		t.Fatalf("notify on closed sink must fail")
	}
	if !errors.Is(err, ErrAlertSinkClosed) {
		t.Fatalf("notify error chain broken: %v", err)
	}
}

func TestAlertPolicyChainIntact(t *testing.T) {
	sink := newAlertSink()
	dispatcher := newAlertDispatcher(sink)
	err := dispatcher.Notify(context.Background(), Alert{ID: "a2", Severity: OpsPriority("urgent"), Message: "moon"})
	if err == nil {
		t.Fatalf("unknown severity must be rejected")
	}
	if !errors.Is(err, ErrAlertSinkBlocked) {
		t.Fatalf("policy error chain broken: %v", err)
	}
}

func TestAlertIDChainIntact(t *testing.T) {
	sink := newAlertSink()
	dispatcher := newAlertDispatcher(sink)
	err := dispatcher.Notify(context.Background(), Alert{Severity: OpsPriorityHigh, Message: "x"})
	if err == nil {
		t.Fatalf("empty alert id must be rejected")
	}
	if !errors.Is(err, ErrAlertSinkBlocked) {
		t.Fatalf("alert id error chain broken: %v", err)
	}
}

func TestAlertWriteClosedRejected(t *testing.T) {
	sink := newAlertSink()
	if err := sink.Close(); err != nil {
		t.Fatal(err)
	}
	if err := sink.WriteLine("late"); err == nil {
		t.Fatalf("write to closed sink must fail")
	}
}

func TestAlertCommitDeniedState(t *testing.T) {
	sink := newAlertSink()
	if err := sink.Close(); err != nil {
		t.Fatal(err)
	}
	if err := sink.Commit(); err == nil {
		t.Fatalf("commit on closed sink must fail")
	}
}

func TestAlertNotifyHonorsCancel(t *testing.T) {
	sink := newAlertSink()
	dispatcher := newAlertDispatcher(sink)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := dispatcher.Notify(ctx, Alert{ID: "a3", Severity: OpsPriorityHigh, Message: "moon"}); err == nil {
		t.Fatalf("notify must honor a cancelled context")
	}
}

func TestAlertLinesIsolated(t *testing.T) {
	sink := newAlertSink()
	if err := sink.WriteLine("one"); err != nil {
		t.Fatal(err)
	}
	lines := sink.Lines()
	lines[0] = "mutated"
	if got := sink.Lines(); got[0] != "one" {
		t.Fatalf("lines leaked internal slice: %v", got)
	}
}
