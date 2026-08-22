package main

import (
	"context"
	"fmt"
)

// Alert is one dispatchable observation alert.
type Alert struct {
	ID       string
	Severity OpsPriority
	Message  string
}

// AlertDispatcher validates alerts against the rule set and writes them to the sink.
type AlertDispatcher struct {
	sink  *AlertSink
	rules []OpsRule
}

func newAlertDispatcher(sink *AlertSink) *AlertDispatcher {
	return &AlertDispatcher{sink: sink, rules: opsRules()}
}

// allowedSeverity reports whether a severity is permitted by the rule set.
func (d *AlertDispatcher) allowedSeverity(severity OpsPriority) error {
	for _, rule := range d.rules {
		if rule.Severity == severity {
			return nil
		}
	}
	return fmt.Errorf("%w: severity %s not covered by rules", ErrAlertSinkBlocked, severity)
}

// Notify validates and writes one alert, then commits the sink.
func (d *AlertDispatcher) Notify(ctx context.Context, alert Alert) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("notify alert %s: %w", alert.ID, err)
	}
	if alert.ID == "" {
		return fmt.Errorf("%w: alert id required", ErrAlertSinkBlocked)
	}
	if err := d.allowedSeverity(alert.Severity); err != nil {
		return err
	}
	line := fmt.Sprintf("%s: %s", alert.ID, alert.Message)
	if err := d.sink.WriteLine(line); err != nil {
		return fmt.Errorf("write alert %s: %w", alert.ID, err)
	}
	if err := d.sink.Commit(); err != nil {
		return fmt.Errorf("%w: commit alert %s", err, alert.ID)
	}
	return nil
}
