package main

import (
	"context"
	"fmt"
)

// NightArchiver closes out an observing night: every run in the batch is
// archived to closed and one log entry is written per archived run.
type NightArchiver struct {
	svc  *OpsService
	logs *NightLogRegistry
}

func newNightArchiver(svc *OpsService, logs *NightLogRegistry) *NightArchiver {
	return &NightArchiver{svc: svc, logs: logs}
}

// ArchiveNight archives all ids and writes the night log; the log is committed
// only when every archived run has been recorded.
func (a *NightArchiver) ArchiveNight(ctx context.Context, night string, ids []string, actor string) (OpsBatchResult, error) {
	log := a.logs.For(night)
	defer log.Close()
	res := OpsBatchResult{Total: len(ids)}
	for _, id := range ids {
		rec, err := a.svc.Transition(ctx, id, 0, OpsStatusClosed, actor)
		if err != nil {
			res.Failed++
			continue
		}
		res.OK++
		if err := log.Append(fmt.Sprintf("%s archived by %s", rec.ID, actor)); err != nil {
			return res, err
		}
	}
	if err := log.Commit(); err != nil {
		return res, err
	}
	return res, nil
}
