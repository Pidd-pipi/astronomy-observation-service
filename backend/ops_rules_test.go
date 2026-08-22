package main

import (
	"testing"
)

func TestRulesTotalCount112(t *testing.T) {
	rules := opsRules()
	if len(rules) != 112 {
		t.Fatalf("rules total = %d, want 112", len(rules))
	}
}

func TestRules01Has0108(t *testing.T) {
	found := false
	for _, rule := range opsRules01() {
		if rule.Code == "OPS-0108" {
			found = true
		}
	}
	if !found {
		t.Fatalf("opsRules01 missing OPS-0108")
	}
}

func TestRules01CodesUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, rule := range opsRules01() {
		if seen[rule.Code] {
			t.Fatalf("duplicate rule code %s", rule.Code)
		}
		seen[rule.Code] = true
	}
}

func TestRules02Has0205(t *testing.T) {
	found := false
	for _, rule := range opsRules02() {
		if rule.Code == "OPS-0205" {
			found = true
		}
	}
	if !found {
		t.Fatalf("opsRules02 missing OPS-0205")
	}
}

func TestRules02CodesUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, rule := range opsRules02() {
		if seen[rule.Code] {
			t.Fatalf("duplicate rule code %s", rule.Code)
		}
		seen[rule.Code] = true
	}
}

func TestRule0303Advisory(t *testing.T) {
	rule := opsRule0303()
	if rule.Terminal {
		t.Fatalf("rule 0303 must not be terminal")
	}
}

func TestRule0704RequiresLabels(t *testing.T) {
	rule := opsRule0704()
	want := len(rule.RequiredLabels)
	if want != 4 {
		t.Fatalf("rule 0704 required labels = %d, want 4", want)
	}
	record := OpsRecord{ID: "r", Subject: "s", Owner: "o", Priority: OpsPriorityNormal, Labels: map[string]string{}}
	if missing := opsRuleMissingLabels(record, rule); len(missing) != want {
		t.Fatalf("missing labels = %d, want %d", len(missing), want)
	}
}

func TestRuleCountsNoPanic(t *testing.T) {
	rules := opsRules()
	counts := opsRuleCounts(rules)
	if len(counts) == 0 {
		t.Fatalf("rule counts empty")
	}
	total := 0
	for _, n := range counts {
		total += n
	}
	if total != len(rules) {
		t.Fatalf("rule counts total %d != %d", total, len(rules))
	}
}

func TestTerminalCountMatches(t *testing.T) {
	rules := opsRules()
	manual := 0
	for _, rule := range rules {
		if rule.Terminal {
			manual++
		}
	}
	if got := opsRuleTerminalCount(rules); got != manual {
		t.Fatalf("terminal count = %d, want %d", got, manual)
	}
}

func TestRules04Has0405(t *testing.T) {
	found := false
	for _, rule := range opsRules04() {
		if rule.Code == "OPS-0405" {
			found = true
		}
	}
	if !found {
		t.Fatalf("opsRules04 missing OPS-0405")
	}
}

func TestRules04Has0408(t *testing.T) {
	found := false
	for _, rule := range opsRules04() {
		if rule.Code == "OPS-0408" {
			found = true
		}
	}
	if !found {
		t.Fatalf("opsRules04 missing OPS-0408")
	}
}

func TestRules04CodesUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, rule := range opsRules04() {
		if seen[rule.Code] {
			t.Fatalf("duplicate rule code %s", rule.Code)
		}
		seen[rule.Code] = true
	}
}

func TestRules05Has0507(t *testing.T) {
	found := false
	for _, rule := range opsRules05() {
		if rule.Code == "OPS-0507" {
			found = true
		}
	}
	if !found {
		t.Fatalf("opsRules05 missing OPS-0507")
	}
}

func TestRules05CodesUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, rule := range opsRules05() {
		if seen[rule.Code] {
			t.Fatalf("duplicate rule code %s", rule.Code)
		}
		seen[rule.Code] = true
	}
}

func TestRuleMissingLabelsSurfaced(t *testing.T) {
	rule := opsRule0101()
	record := OpsRecord{ID: "r", Labels: map[string]string{"site": "s", "operator": "o", "evidence": "e", "reviewed": "y"}}
	if missing := opsRuleMissingLabels(record, rule); len(missing) != 0 {
		t.Fatalf("clean record reported missing %v", missing)
	}
}

func TestRules01Has0103(t *testing.T) {
	found := false
	for _, rule := range opsRules01() {
		if rule.Code == "OPS-0103" {
			found = true
		}
	}
	if !found {
		t.Fatalf("opsRules01 missing OPS-0103")
	}
}
