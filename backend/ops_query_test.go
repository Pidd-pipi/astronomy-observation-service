package main

import (
	"testing"
)

func TestFilterRecordsNoAlias(t *testing.T) {
	items := []OpsRecord{
		{ID: "a", Subject: "alpha", Owner: "o", Status: OpsStatusActive, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}},
		{ID: "b", Subject: "beta", Owner: "o", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s2"}},
		{ID: "c", Subject: "gamma", Owner: "o", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s3"}},
	}
	filtered := filterOpsRecords(items, OpsQuery{Status: OpsStatusQueued})
	if len(filtered) != 2 {
		t.Fatalf("filtered = %d, want 2", len(filtered))
	}
	if items[0].ID != "a" || items[1].ID != "b" || items[2].ID != "c" {
		t.Fatalf("filter rewrote the input slice: %v", []string{items[0].ID, items[1].ID, items[2].ID})
	}
	filtered = append(filtered, OpsRecord{ID: "x"})
	if len(items) != 3 {
		t.Fatalf("filter shares backing array with input: len=%d", len(items))
	}
}

func TestClonePageBackingIsolated(t *testing.T) {
	original := OpsPage{Items: make([]OpsRecord, 1, 4)}
	original.Items[0] = OpsRecord{ID: "a", Labels: map[string]string{"site": "s1"}}
	cloned := opsClonePage(original)
	cloned.Items = append(cloned.Items, OpsRecord{ID: "z"})
	if original.Items[0].ID != "a" || len(original.Items) != 1 {
		t.Fatalf("clone mutated original: len=%d", len(original.Items))
	}
	if original.Items[:2][1].ID != "" {
		t.Fatalf("cloned page shares backing array with original: %q", original.Items[:2][1].ID)
	}
}

func TestQueryDefaultsCapsSize(t *testing.T) {
	q := opsQueryDefaults(OpsQuery{PageSize: 999})
	if q.PageSize != 200 {
		t.Fatalf("page size = %d, want 200", q.PageSize)
	}
}

func TestBoundsClamped(t *testing.T) {
	start, end := opsBounds(5, 3, 25)
	if start != 5 || end != 5 {
		t.Fatalf("bounds = (%d,%d), want (5,5)", start, end)
	}
}
