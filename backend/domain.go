package main

type ObservationRun struct {
	ID         string `json:"id"`
	Target     string `json:"target"`
	Instrument string `json:"instrument"`
	Quality    string `json:"quality"`
	Status     string `json:"status"`
	StartedAt  string `json:"startedAt"`
}
type StatusChange struct {
	Status string `json:"status"`
}

// EnsureDefaults fills empty fields with stable fallbacks so consumers never see
// zero-value gaps in a run payload.
func (r ObservationRun) EnsureDefaults() ObservationRun {
	if r.Target == "" {
		r.Target = "unscheduled"
	}
	if r.Instrument == "" {
		r.Instrument = "unspecified"
	}
	if r.Quality == "" {
		r.Quality = "unknown"
	}
	if r.Status == "" {
		r.Status = "planned"
	}
	return r
}
