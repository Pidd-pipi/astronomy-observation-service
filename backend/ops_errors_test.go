package main

import (
	"errors"
	"strings"
	"testing"
)

func TestErrChainPreserved(t *testing.T) {
	err := wrapOps("conflict", "store.update", ErrOpsConflict)
	if !errors.Is(err, ErrOpsConflict) {
		t.Fatalf("wrapped conflict error chain broken: %v", err)
	}
}

func TestErrCodeKeepsTypedCode(t *testing.T) {
	err := wrapOps("create", "store.put", ErrOpsConflict)
	if code := opsCode(err); code != "create" {
		t.Fatalf("opsCode = %q, want create", code)
	}
}

func TestErrMessageKeepsCause(t *testing.T) {
	err := wrapOps("conflict", "store.update", ErrOpsConflict)
	text := err.Error()
	if !strings.Contains(text, "operations revision conflict") {
		t.Fatalf("error message lost cause: %q", text)
	}
}

func TestErrHTTPMapsConflict409(t *testing.T) {
	rr := newTestRecorder()
	opsHTTPError(rr, wrapOps("conflict", "store.update", ErrOpsConflict))
	if rr.Code != 409 {
		t.Fatalf("conflict mapped to %d, want 409", rr.Code)
	}
}

func TestErrConflictClassified(t *testing.T) {
	if !opsIsConflict(wrapOps("conflict", "store.update", ErrOpsConflict)) {
		t.Fatalf("wrapped conflict must classify as conflict")
	}
}
