package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func newOpsTestHandler(seed []OpsRecord) http.Handler {
	svc := newOpsService(seed)
	api := newOpsAPI(svc)
	m := http.NewServeMux()
	api.register(m)
	return m
}

func TestHandlerGetNoPollution(t *testing.T) {
	handler := newOpsTestHandler(seedOpsRecords())
	req := httptest.NewRequest(http.MethodGet, "/api/ops/records/op-1001", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET record code=%d", rec.Code)
	}
	req2 := httptest.NewRequest(http.MethodGet, "/api/ops/records", nil)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("GET list code=%d", rec2.Code)
	}
	var page OpsPage
	if err := json.Unmarshal(rec2.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	for _, item := range page.Items {
		if item.ID == "op-1001" {
			if _, ok := item.Labels["fetchedAt"]; ok {
				t.Fatalf("read path polluted stored record: %v", item.Labels)
			}
		}
	}
}

func TestHandlerListNoPollution(t *testing.T) {
	handler := newOpsTestHandler(seedOpsRecords())
	req := httptest.NewRequest(http.MethodGet, "/api/ops/records", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET list code=%d", rec.Code)
	}
	req2 := httptest.NewRequest(http.MethodGet, "/api/ops/records", nil)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	var page OpsPage
	if err := json.Unmarshal(rec2.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	for _, item := range page.Items {
		if _, ok := item.Labels["listedAt"]; ok {
			t.Fatalf("list path polluted stored record: %v", item.Labels)
		}
	}
}

func TestReadPathConcurrentSafe(t *testing.T) {
	handler := newOpsTestHandler(seedOpsRecords())
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 20; j++ {
				req := httptest.NewRequest(http.MethodGet, "/api/ops/records/op-1001", nil)
				handler.ServeHTTP(httptest.NewRecorder(), req)
			}
		}()
	}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 20; j++ {
				req := httptest.NewRequest(http.MethodGet, "/api/ops/records?status=active", nil)
				handler.ServeHTTP(httptest.NewRecorder(), req)
			}
		}()
	}
	close(start)
	wg.Wait()
}

func TestCreateDuplicate409(t *testing.T) {
	handler := newOpsTestHandler(nil)
	body := `{"id":"op-dup","subject":"dup","owner":"o","priority":"normal","status":"queued","labels":{"site":"s"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/ops/records", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("first create code=%d", rec.Code)
	}
	req2 := httptest.NewRequest(http.MethodPost, "/api/ops/records", bytes.NewBufferString(body))
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusConflict {
		t.Fatalf("duplicate create code=%d, want 409", rec2.Code)
	}
}

func TestOpsTransitionIllegalReturns422(t *testing.T) {
	handler := newOpsTestHandler(seedOpsRecords())
	req := httptest.NewRequest(http.MethodPost, "/api/ops/records/op-1002/transition", bytes.NewBufferString(`{"status":"queued","actor":"t"}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("illegal transition code=%d, want 422", rec.Code)
	}
}

type testRecorder struct {
	Code int
	body bytes.Buffer
}

func (r *testRecorder) Header() http.Header         { return http.Header{} }
func (r *testRecorder) Write(b []byte) (int, error) { return r.body.Write(b) }
func (r *testRecorder) WriteHeader(code int)        { r.Code = code }

func newTestRecorder() *testRecorder { return &testRecorder{Code: 200} }
