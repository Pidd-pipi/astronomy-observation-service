package main

import (
	"testing"
)

func LogAppendClosedRejected(t *testing.T) {
	log := newNightLog()
	if err := log.Append("entry-1"); err != nil {
		t.Fatal(err)
	}
	if err := log.Close(); err != nil {
		t.Fatal(err)
	}
	if err := log.Append("entry-2"); err == nil {
		t.Fatalf("append after close must fail")
	}
}

func LogCommitClosedRejected(t *testing.T) {
	log := newNightLog()
	if err := log.Close(); err != nil {
		t.Fatal(err)
	}
	if err := log.Commit(); err == nil {
		t.Fatalf("commit after close must fail")
	}
	if log.Committed() {
		t.Fatalf("closed log must not be committed")
	}
}

func LogEntriesIsolated(t *testing.T) {
	log := newNightLog()
	if err := log.Append("entry-1"); err != nil {
		t.Fatal(err)
	}
	entries := log.Entries()
	entries[0] = "mutated"
	if got := log.Entries(); got[0] != "entry-1" {
		t.Fatalf("entries leaked internal slice: %v", got)
	}
}

func TestNightLogCommitMarksCommitted(t *testing.T) {
	log := newNightLog()
	if err := log.Append("entry-1"); err != nil {
		t.Fatal(err)
	}
	if err := log.Commit(); err != nil {
		t.Fatal(err)
	}
	if !log.Committed() {
		t.Fatalf("committed log must be marked committed")
	}
}
