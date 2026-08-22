package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// opsAPI exposes the operations-management domain over HTTP.
type opsAPI struct {
	svc *OpsService
}

func newOpsAPI(svc *OpsService) *opsAPI { return &opsAPI{svc: svc} }

func (a *opsAPI) register(mux *http.ServeMux) {
	mux.HandleFunc("/api/ops/records", a.handleRecords)
	mux.HandleFunc("/api/ops/records/", a.handleRecord)
	mux.HandleFunc("/api/ops/snapshot", a.handleSnapshot)
	mux.HandleFunc("/api/ops/rules", a.handleRules)
	mux.HandleFunc("/api/ops/batch/archive", a.handleBatchArchive)
}

func (a *opsAPI) handleRecords(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		q := OpsQuery{
			Subject:  r.URL.Query().Get("subject"),
			Status:   OpsStatus(r.URL.Query().Get("status")),
			Priority: OpsPriority(r.URL.Query().Get("priority")),
			Owner:    r.URL.Query().Get("owner"),
		}
		if p := r.URL.Query().Get("page"); p != "" {
			if n, err := strconv.Atoi(p); err == nil {
				q.Page = n
			}
		}
		if ps := r.URL.Query().Get("pageSize"); ps != "" {
			if n, err := strconv.Atoi(ps); err == nil {
				q.PageSize = n
			}
		}
		page, err := a.svc.Search(r.Context(), q)
		if err != nil {
			opsHTTPError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, page)
	case http.MethodPost:
		var record OpsRecord
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			opsJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid record body"})
			return
		}
		created, err := a.svc.Create(r.Context(), record)
		if err != nil {
			opsHTTPError(w, err)
			return
		}
		opsJSON(w, http.StatusCreated, created)
	default:
		opsJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (a *opsAPI) handleRecord(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/ops/records/"), "/")
	if rest == "" {
		opsJSON(w, http.StatusNotFound, map[string]string{"error": "record not found"})
		return
	}
	id, suffix, hasSuffix := strings.Cut(rest, "/")
	if !hasSuffix {
		if r.Method != http.MethodGet {
			opsJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		rec, err := a.svc.Get(r.Context(), id)
		if err != nil {
			opsHTTPError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, rec)
		return
	}
	switch suffix {
	case "audit":
		if r.Method != http.MethodGet {
			opsJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		opsJSON(w, http.StatusOK, a.svc.Audit(id))
	case "transition":
		if r.Method != http.MethodPost {
			opsJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		var req struct {
			Status   OpsStatus `json:"status"`
			Expected int       `json:"expected"`
			Actor    string    `json:"actor"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			opsJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid transition body"})
			return
		}
		actor := req.Actor
		if actor == "" {
			actor = opsActorFromRequest(r)
		}
		rec, err := a.svc.Transition(r.Context(), id, req.Expected, req.Status, actor)
		if err != nil {
			opsHTTPError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, rec)
	default:
		opsJSON(w, http.StatusNotFound, map[string]string{"error": "record path not found"})
	}
}

func (a *opsAPI) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		opsJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	opsJSON(w, http.StatusOK, a.svc.Snapshot())
}

func (a *opsAPI) handleRules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		opsJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	opsJSON(w, http.StatusOK, a.svc.Rules())
}

func (a *opsAPI) handleBatchArchive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		opsJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var req struct {
		IDs   []string `json:"ids"`
		Actor string   `json:"actor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		opsJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid batch body"})
		return
	}
	actor := req.Actor
	if actor == "" {
		actor = opsActorFromRequest(r)
	}
	res, err := a.svc.ArchiveBatch(r.Context(), req.IDs, actor)
	if err != nil {
		opsHTTPError(w, err)
		return
	}
	opsJSON(w, http.StatusOK, res)
}
