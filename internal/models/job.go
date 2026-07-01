package models

type Job struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name"`
	Template   string                 `json:"template,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Schedule   string                 `json:"schedule,omitempty"`
	Status     string                 `json:"status,omitempty"`
	Audit      *Audit                 `json:"audit,omitempty"`
}

type JobResponse struct {
	Code int `json:"code"`
	Job  Job `json:"job"`
}

type JobListResponse struct {
	Code int   `json:"code"`
	Jobs []Job `json:"jobs"`
}

type JobCreateRequest struct {
	Name       string                 `json:"name"`
	Template   string                 `json:"template,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Schedule   string                 `json:"schedule,omitempty"`
}

type JobStatusResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

type JobTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type JobTemplateListResponse struct {
	Code      int           `json:"code"`
	Templates []JobTemplate `json:"templates"`
}
