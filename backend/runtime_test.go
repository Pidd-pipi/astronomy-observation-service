package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func wrappedHealthHandler() http.Handler {
	return newEnterpriseServer("127.0.0.1:0", newRouter(newRunStore())).Handler
}

func HealthDeadlineTrue(t *testing.T) {
	handler := wrappedHealthHandler()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("health code=%d", rr.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["deadline"] != "true" {
		t.Fatalf("health deadline = %q, want true", body["deadline"])
	}
}

func HealthEchoesRequestID(t *testing.T) {
	handler := wrappedHealthHandler()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("X-Request-ID", "req-test-123")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["requestId"] != "req-test-123" {
		t.Fatalf("health requestId = %q, want req-test-123", body["requestId"])
	}
}
