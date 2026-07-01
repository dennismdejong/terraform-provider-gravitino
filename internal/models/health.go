package models

type HealthResponse struct {
	Code   int           `json:"code"`
	Status string        `json:"status"`
	Checks []HealthCheck `json:"checks"`
}

type HealthCheck struct {
	Name    string            `json:"name"`
	Status  string            `json:"status"`
	Details map[string]string `json:"details,omitempty"`
}
