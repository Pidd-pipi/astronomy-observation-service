package main

import "testing"

func TestAuditHistoryIntact(t *testing.T) {
	audit := newOpsAudit()
	audit.Add("r1", "created", "a")
	audit.Add("r2", "created", "b")
	if got := audit.For("r1"); len(got) != 1 {
		t.Fatalf("r1 history = %d, want 1", len(got))
	}
	if got := audit.For("r2"); len(got) != 1 {
		t.Fatalf("r2 history = %d, want 1", len(got))
	}
	if got := audit.For("r1"); len(got) != 1 {
		t.Fatalf("r1 history corrupted by reading r2: %d", len(got))
	}
}
