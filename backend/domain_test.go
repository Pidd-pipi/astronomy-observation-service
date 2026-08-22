package main

import "testing"

func TestOpsCloneInitializesLabels(t *testing.T) {
	original := OpsRecord{ID: "op-x", Subject: "s", Owner: "o", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "north-dome"}}
	cloned := original.Clone()
	cloned.Labels["site"] = "mutated"
	if original.Labels["site"] != "north-dome" {
		t.Fatalf("clone shares labels with original: %v", original.Labels)
	}
	if cloned.Labels["site"] != "mutated" {
		t.Fatalf("clone labels not writable: %v", cloned.Labels)
	}
}

func TestNormalizeInitializesNilLabels(t *testing.T) {
	record := normalizeOpsRecord(OpsRecord{ID: "op-y", Subject: "s", Owner: "o", Priority: OpsPriorityNormal})
	if record.Labels == nil {
		t.Fatalf("normalized record must have non-nil labels")
	}
}

func TestRunDefaultsApplied(t *testing.T) {
	run := ObservationRun{ID: "r-1"}.EnsureDefaults()
	if run.Target != "unscheduled" || run.Instrument != "unspecified" || run.Quality != "unknown" || run.Status != "planned" {
		t.Fatalf("defaults not applied: %+v", run)
	}
}
