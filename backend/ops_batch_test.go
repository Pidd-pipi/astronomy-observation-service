package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func seedBatchRecords() []OpsRecord {
	return []OpsRecord{
		{ID: "b1", Subject: "one", Owner: "o", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Revision: 1, Labels: map[string]string{"site": "s"}},
		{ID: "b2", Subject: "two", Owner: "o", Status: OpsStatusActive, Priority: OpsPriorityNormal, Revision: 1, Labels: map[string]string{"site": "s"}},
		{ID: "b3", Subject: "three", Owner: "o", Status: OpsStatusPaused, Priority: OpsPriorityNormal, Revision: 1, Labels: map[string]string{"site": "s"}},
	}
}

func BatchConcurrentComplete(t *testing.T) {
	svc := newOpsService(seedBatchRecords())
	ids := make([]string, 0, 30)
	for i := 0; i < 30; i++ {
		ids = append(ids, fmt.Sprintf("b%d", i%3+1))
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			res, err := svc.ArchiveBatch(context.Background(), ids, "night-shift")
			if err != nil {
				t.Errorf("batch error: %v", err)
				return
			}
			if len(res.Items) != len(ids) {
				t.Errorf("want %d items, got %d", len(ids), len(res.Items))
			}
			if res.OK+res.Failed != len(ids) {
				t.Errorf("counts %d+%d != %d", res.OK, res.Failed, len(ids))
			}
		}()
	}
	close(start)
	wg.Wait()
}

func BatchCollectsFailures(t *testing.T) {
	svc := newOpsService(seedBatchRecords())
	res, err := svc.ArchiveBatch(context.Background(), []string{"b1", "b2", "missing-1"}, "night-shift")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Items) != 3 {
		t.Fatalf("want 3 items, got %d", len(res.Items))
	}
	if res.Failed != 1 {
		t.Fatalf("want 1 failed, got %d", res.Failed)
	}
	found := false
	for _, item := range res.Items {
		if item.ID == "missing-1" {
			found = true
			if item.Err == "" {
				t.Fatalf("missing-1 must carry error")
			}
		}
	}
	if !found {
		t.Fatalf("missing-1 not present in batch items")
	}
}

func BatchCountsTallies(t *testing.T) {
	items := []OpsBatchItem{
		{ID: "a", Err: ""},
		{ID: "b", Err: "boom"},
		{ID: "c", Err: ""},
	}
	ok, failed := batchCounts(items)
	if ok != 2 || failed != 1 {
		t.Fatalf("batchCounts = (%d,%d), want (2,1)", ok, failed)
	}
}

func BatchNoPanicMany(t *testing.T) {
	svc := newOpsService(seedBatchRecords())
	ids := make([]string, 0, 12)
	for i := 0; i < 12; i++ {
		ids = append(ids, fmt.Sprintf("b%d", i%3+1))
	}
	res, err := svc.ArchiveBatch(context.Background(), ids, "night-shift")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Items) != 12 {
		t.Fatalf("want 12 items, got %d", len(res.Items))
	}
}
