package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunHTTP(t *testing.T) {
	h := newRouter(newRunStore())
	cases := []struct {
		name, path, body string
		code             int
	}{{"collection", "/api/runs", "", 200}, {"archive", "/api/runs/run-501/status", `{"status":"archived"}`, 200}, {"invalid", "/api/runs/run-501/status", `{"status":"lost"}`, 400}, {"missing", "/api/runs/nope/status", `{"status":"running"}`, 404}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRecorder()
			method := http.MethodGet
			if tc.body != "" {
				method = http.MethodPost
			}
			h.ServeHTTP(r, httptest.NewRequest(method, tc.path, bytes.NewBufferString(tc.body)))
			if r.Code != tc.code {
				t.Fatalf("got %d", r.Code)
			}
		})
	}
}
