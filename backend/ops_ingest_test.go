package main

import (
	"context"
	"testing"
	"time"
)

func TestIngestDerivesFromRequest(t *testing.T) {
	ingest := newPlanIngest(newIngestStore())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := ingest.Ingest(ctx, map[string]string{"p1": "plan-1"}); err == nil {
		t.Fatalf("ingest must honor an already-cancelled request context")
	}
}

func TestIngestReportsCancelError(t *testing.T) {
	ingest := newPlanIngest(newIngestStore())
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	if _, err := ingest.Ingest(ctx, map[string]string{"p1": "plan-1"}); err == nil {
		t.Fatalf("ingest must surface the deadline error instead of swallowing it")
	}
}

func TestIngestStoreCancelRejectedPut(t *testing.T) {
	store := newIngestStore()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := store.Put(ctx, "p1", "plan-1"); err == nil {
		t.Fatalf("cancelled write must fail")
	}
}

func TestIngestStoreCancelRejected(t *testing.T) {
	store := newIngestStore()
	ctx, cancel := context.WithCancel(context.Background())
	_ = store.Put(context.Background(), "p1", "plan-1")
	cancel()
	if _, err := store.Get(ctx, "p1"); err == nil {
		t.Fatalf("cancelled read must fail")
	}
}
