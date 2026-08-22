package main

import (
	"context"
	"strings"
	"time"
)

// PlanIngest ingests observation plans with a per-request deadline.
type PlanIngest struct {
	store *IngestStore
}

func newPlanIngest(store *IngestStore) *PlanIngest { return &PlanIngest{store: store} }

// Ingest writes each plan using the request context; cancellation stops the batch.
func (p *PlanIngest) Ingest(ctx context.Context, plans map[string]string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	written := 0
	for id, plan := range plans {
		select {
		case <-ctx.Done():
			return written, nil
		default:
		}
		if strings.TrimSpace(plan) == "" {
			continue
		}
		if err := p.store.Put(ctx, id, plan); err != nil {
			return written, err
		}
		written++
	}
	return written, nil
}
