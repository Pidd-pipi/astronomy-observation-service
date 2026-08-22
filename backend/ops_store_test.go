package main

import (
	"context"
	"testing"
)

func TestStoreGetIsolatedCopy(t *testing.T) {
	store := newOpsStore([]OpsRecord{
		{ID: "op-1", Subject: "session", Owner: "o", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "north-dome"}},
	})
	rec, err := store.Get(context.Background(), "op-1")
	if err != nil {
		t.Fatal(err)
	}
	rec.Labels["site"] = "mutated"
	again, err := store.Get(context.Background(), "op-1")
	if err != nil {
		t.Fatal(err)
	}
	if again.Labels["site"] != "north-dome" {
		t.Fatalf("store mutated by caller: %v", again.Labels)
	}
}

func TestStoreListIsolatedCopies(t *testing.T) {
	store := newOpsStore([]OpsRecord{
		{ID: "op-1", Subject: "session", Owner: "o", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "north-dome"}},
	})
	items, err := store.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("want 1, got %d", len(items))
	}
	items[0].Labels["site"] = "mutated"
	again, _ := store.List(context.Background())
	if again[0].Labels["site"] != "north-dome" {
		t.Fatalf("store mutated by caller: %v", again[0].Labels)
	}
}

func TestStoreUpdateStoresCopy(t *testing.T) {
	store := newOpsStore([]OpsRecord{
		{ID: "op-1", Subject: "session", Owner: "o", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Revision: 1, Labels: map[string]string{"site": "north-dome"}},
	})
	rec, _ := store.Get(context.Background(), "op-1")
	rec.Status = OpsStatusActive
	if err := store.Update(context.Background(), rec, 1); err != nil {
		t.Fatal(err)
	}
	rec.Labels["site"] = "mutated"
	again, _ := store.Get(context.Background(), "op-1")
	if again.Labels["site"] != "north-dome" {
		t.Fatalf("update stored caller reference: %v", again.Labels)
	}
}

func TestStorePutNormalizesId(t *testing.T) {
	store := newOpsStore(nil)
	item := OpsRecord{ID: "OP-X", Subject: "  session  ", Owner: " o ", Status: OpsStatusQueued, Priority: OpsPriorityNormal}
	if err := store.Put(context.Background(), item); err != nil {
		t.Fatal(err)
	}
	got, err := store.Get(context.Background(), "op-x")
	if err != nil {
		t.Fatalf("normalized id not findable: %v", err)
	}
	if got.Owner != "o" || got.Subject != "session" {
		t.Fatalf("record not normalized: %+v", got)
	}
}
