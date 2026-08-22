package main

import (
	"context"
	"sync"
)

// OpsBatchItem is one unit of a batch archive: either the archived record or the failure.
type OpsBatchItem struct {
	ID     string    `json:"id"`
	Record OpsRecord `json:"record,omitempty"`
	Err    string    `json:"error,omitempty"`
}

// OpsBatchResult summarizes a batch archive run.
type OpsBatchResult struct {
	Total  int            `json:"total"`
	OK     int            `json:"ok"`
	Failed int            `json:"failed"`
	Items  []OpsBatchItem `json:"items"`
}

// ArchiveBatch transitions every id to closed concurrently and reports per-item
// results. A missing or invalid id produces an item carrying its error instead of
// blocking the rest of the batch; the channel is closed exactly once after every
// worker finishes so callers never hang or panic on a partial run.
func (s *OpsService) ArchiveBatch(ctx context.Context, ids []string, actor string) (OpsBatchResult, error) {
	results := make(chan OpsBatchItem, len(ids))
	var wg sync.WaitGroup
	wg.Add(len(ids))
	for _, id := range ids {
		go func(id string) {
			defer wg.Done()
			rec, err := s.Transition(ctx, id, 0, OpsStatusClosed, actor)
			if err != nil {
				results <- OpsBatchItem{ID: id, Err: err.Error()}
				return
			}
			results <- OpsBatchItem{ID: id, Record: rec}
		}(id)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	items := make([]OpsBatchItem, 0, len(ids))
	for item := range results {
		items = append(items, item)
	}
	ok, failed := batchCounts(items)
	return OpsBatchResult{Total: len(ids), OK: ok, Failed: failed, Items: items}, nil
}

