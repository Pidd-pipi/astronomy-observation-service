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
		wg.Add(1)
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
	var items []OpsBatchItem
	for item := range results {
		items = append(items, item)
	}
	return buildBatchResult(len(ids), items), nil
}

// buildBatchResult assembles a batch result with consistent totals.
func buildBatchResult(total int, items []OpsBatchItem) OpsBatchResult {
	res := OpsBatchResult{Total: total, Items: items}
	res.OK, res.Failed = batchCounts(items)
	return res
}
