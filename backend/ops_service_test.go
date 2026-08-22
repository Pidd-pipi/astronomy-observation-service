package main

import (
	"context"
	"testing"
)

func TestSearchTotalFiltered(t *testing.T) {
	svc := newOpsService(seedOpsRecords())
	page, err := svc.Search(context.Background(), OpsQuery{Status: OpsStatusActive})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 {
		t.Fatalf("search total = %d, want 1 (active)", page.Total)
	}
	if len(page.Items) != 1 {
		t.Fatalf("search items = %d, want 1", len(page.Items))
	}
}

func TestSearchPriorityOnly(t *testing.T) {
	svc := newOpsService(seedOpsRecords())
	page, err := svc.Search(context.Background(), OpsQuery{Priority: OpsPriorityCritical})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 {
		t.Fatalf("search total = %d, want 1 (critical)", page.Total)
	}
	for _, item := range page.Items {
		if item.Priority != OpsPriorityCritical {
			t.Fatalf("search returned non-critical record: %s", item.ID)
		}
	}
}

func TestSearchPagedAllStatuses(t *testing.T) {
	svc := newOpsService(seedOpsRecords())
	page, err := svc.Search(context.Background(), OpsQuery{Page: 1, PageSize: 100})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 4 {
		t.Fatalf("search total = %d, want 4", page.Total)
	}
}
