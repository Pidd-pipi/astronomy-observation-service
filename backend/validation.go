package main

import "fmt"

func validateRunStatus(s string) error {
	switch s {
	case "planned", "running", "archived":
		return nil
	default:
		return fmt.Errorf("status must be planned, running, or archived")
	}
}
