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

// ArchiveBatch transitions every id to closed concurrently and reports per-item results.
func (s *OpsService) ArchiveBatch(ctx context.Context, ids []string, actor string) (OpsBatchResult, error) {
	results := make(chan OpsBatchItem, len(ids))
	var wg sync.WaitGroup
	for _, id := range ids {
		go func(id string) {
			wg.Add(1)
			defer wg.Done()
			rec, err := s.Transition(ctx, id, 0, OpsStatusClosed, actor)
			if err != nil {
				return
			}
			results <- OpsBatchItem{ID: id, Record: rec}
			defer close(results)
		}(id)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	var items []OpsBatchItem
	for item := range results {
		items = append(items, item)
	}
	return OpsBatchResult{Total: len(ids), Items: items}, nil
}

