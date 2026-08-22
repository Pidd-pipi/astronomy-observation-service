package main

import (
	"context"
	"testing"
)

func seedArchiverRecords() []OpsRecord {
	return []OpsRecord{
		{ID: "ar-1", Subject: "one", Owner: "o", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Revision: 1, Labels: map[string]string{"site": "s"}},
		{ID: "ar-2", Subject: "two", Owner: "o", Status: OpsStatusActive, Priority: OpsPriorityNormal, Revision: 1, Labels: map[string]string{"site": "s"}},
	}
}

func ArchiveNightKeepsError(t *testing.T) {
	svc := newOpsService(append(seedArchiverRecords(),
		OpsRecord{ID: "blocked", Subject: "bad", Owner: "o", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Revision: 1, Labels: map[string]string{"site": "s"}},
	))
	logs := newNightLogRegistry()
	archiver := newNightArchiver(svc, logs)
	_, err := archiver.ArchiveNight(context.Background(), "n1", []string{"ar-1", "blocked"}, "shift")
	if err == nil {
		t.Fatalf("archive with blocked entry must return error")
	}
}

func TestArchiveNightClosesLog(t *testing.T) {
	svc := newOpsService(seedArchiverRecords())
	logs := newNightLogRegistry()
	archiver := newNightArchiver(svc, logs)
	res, err := archiver.ArchiveNight(context.Background(), "n2", []string{"ar-1", "ar-2"}, "shift")
	if err != nil {
		t.Fatal(err)
	}
	if res.OK != 2 {
		t.Fatalf("ok = %d, want 2", res.OK)
	}
	if logs.For("n2").Open() {
		t.Fatalf("night log must be closed after archive")
	}
	if !logs.For("n2").Committed() {
		t.Fatalf("night log must be committed after archive")
	}
}

func ArchiveNightCommitsEmpty(t *testing.T) {
	svc := newOpsService(seedArchiverRecords())
	logs := newNightLogRegistry()
	archiver := newNightArchiver(svc, logs)
	res, err := archiver.ArchiveNight(context.Background(), "n3", []string{"missing-1", "missing-2"}, "shift")
	if err != nil {
		t.Fatal(err)
	}
	if res.OK != 0 || res.Failed != 2 {
		t.Fatalf("result ok=%d failed=%d, want 0/2", res.OK, res.Failed)
	}
	if !logs.For("n3").Committed() {
		t.Fatalf("all-failed archive must still commit the night log")
	}
}
