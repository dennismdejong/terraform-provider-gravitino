package models

import (
	"fmt"
	"time"
)

type ErrorResponse struct {
	Code    int      `json:"code"`
	Type    string   `json:"type"`
	Message string   `json:"message"`
	Stack   []string `json:"stack,omitempty"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s (%s)", e.Message, e.Type)
}

type Audit struct {
	Creator          string     `json:"creator,omitempty"`
	CreateTime       *time.Time `json:"createTime,omitempty"`
	LastModifier     string     `json:"lastModifier,omitempty"`
	LastModifiedTime *time.Time `json:"lastModifiedTime,omitempty"`
}

type NameIdentifier struct {
	Namespace []string `json:"namespace"`
	Name      string   `json:"name"`
}

type DropResponse struct {
	Code    int  `json:"code"`
	Dropped bool `json:"dropped"`
}

type IdentifiersResponse struct {
	Code        int              `json:"code"`
	Identifiers []NameIdentifier `json:"identifiers"`
}

type BaseResponse struct {
	Code int `json:"code"`
}
