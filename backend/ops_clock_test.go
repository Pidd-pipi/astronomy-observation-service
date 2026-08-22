package main

import (
	"context"
	"testing"
	"time"
)

func StampParseErrorReported(t *testing.T) {
	if _, err := opsParseStamp("not-a-time"); err == nil {
		t.Fatalf("opsParseStamp must report parse error")
	}
}

func AgeInvalidStampZero(t *testing.T) {
	now := time.Now().UTC()
	if age := opsAge(now, "not-a-time"); age != 0 {
		t.Fatalf("opsAge invalid stamp = %v, want 0", age)
	}
}

func BackoffCappedAtSix(t *testing.T) {
	if d := opsBackoff(10); d > time.Second {
		t.Fatalf("opsBackoff not capped: %v", d)
	}
}

func TestOpsDelayHonorsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := opsDelay(ctx, time.Hour); err == nil {
		t.Fatalf("opsDelay must return error on cancelled ctx")
	}
}
