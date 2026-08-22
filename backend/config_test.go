package main

import "testing"

func ConfigDefaultsPort(t *testing.T) {
	t.Setenv("PORT", "")
	cfg := loadConfig()
	if cfg.Port != "8080" {
		t.Fatalf("default port = %q, want 8080", cfg.Port)
	}
}

func TestLoadConfigHonorsPort(t *testing.T) {
	t.Setenv("PORT", "9090")
	cfg := loadConfig()
	if cfg.Port != "9090" {
		t.Fatalf("port = %q, want 9090", cfg.Port)
	}
}

func RunStatusEmptyRejected(t *testing.T) {
	if err := validateRunStatus(""); err == nil {
		t.Fatalf("empty status must be rejected")
	}
}
