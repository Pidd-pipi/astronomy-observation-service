package main

import "net/http"

func main() {
	c := loadConfig()
	r := newRouter(newRunStore())
	svc := newOpsService(seedOpsRecords())
	api := newOpsAPI(svc)
	m := http.NewServeMux()
	m.Handle("/healthz", r)
	m.Handle("/api/runs", r)
	m.Handle("/api/runs/", r)
	api.register(m)
	m.Handle("/", staticHandler())
	if e := serveAddress(":"+c.Port, m); e != nil {
		panic(e)
	}
}

// seedOpsRecords returns the observatory's operations records for the current season.
func seedOpsRecords() []OpsRecord {
	stamp := "2026-08-22T00:00:00Z"
	return []OpsRecord{
		{ID: "op-1001", Subject: "M31 photometry session", Owner: "yu-obs", Status: OpsStatusQueued, Priority: OpsPriorityHigh, Revision: 1, Labels: map[string]string{"site": "north-dome", "operator": "yu-obs", "evidence": "camera-3", "reviewed": "yes"}, CreatedAt: stamp, UpdatedAt: stamp},
		{ID: "op-1002", Subject: "NGC 7000 wide-field scan", Owner: "lin-obs", Status: OpsStatusActive, Priority: OpsPriorityNormal, Revision: 1, Labels: map[string]string{"site": "south-dome", "operator": "lin-obs", "evidence": "wide-cam"}, CreatedAt: stamp, UpdatedAt: stamp},
		{ID: "op-1003", Subject: "Target handoff review", Owner: "wei-obs", Status: OpsStatusPaused, Priority: OpsPriorityCritical, Revision: 2, Labels: map[string]string{"site": "control-room", "operator": "wei-obs", "evidence": "handoff-form", "reviewed": "yes"}, CreatedAt: stamp, UpdatedAt: stamp},
		{ID: "op-1004", Subject: "Archive night 20260821", Owner: "yu-obs", Status: OpsStatusClosed, Priority: OpsPriorityLow, Revision: 3, Labels: map[string]string{"site": "archive", "operator": "yu-obs", "evidence": "night-log"}, CreatedAt: stamp, UpdatedAt: stamp},
	}
}
