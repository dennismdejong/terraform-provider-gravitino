package models

type ModelVersionLinkRequest struct {
	Version    string            `json:"version"`
	URI        string            `json:"uri,omitempty"`
	Aliases    []string          `json:"aliases,omitempty"`
	Comment    string            `json:"comment,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type ModelVersionResponse struct {
	Code         int32        `json:"code"`
	ModelVersion ModelVersion `json:"modelVersion"`
}

type ModelVersion struct {
	Version    string            `json:"version"`
	URI        string            `json:"uri,omitempty"`
	Aliases    []string          `json:"aliases,omitempty"`
	Comment    string            `json:"comment,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Audit      *Audit            `json:"audit,omitempty"`
}

type ModelVersionListResponse struct {
	Code          int32          `json:"code"`
	ModelVersions []ModelVersion `json:"modelVersions"`
}

type ModelVersionURIResponse struct {
	Code int32  `json:"code"`
	URI  string `json:"uri"`
}
