package client

import (
	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

type RunEventRequest struct {
	SourceIdentifier   string `json:"sourceIdentifier"`
	TargetIdentifier   string `json:"targetIdentifier"`
	Operation          string `json:"operation"`
	OperationParameter string `json:"operationParameter,omitempty"`
	RunTime            string `json:"runTime,omitempty"`
}

func (c *Client) PostRunEvent(req *RunEventRequest) (*models.BaseResponse, error) {
	var result models.BaseResponse
	if err := c.Post("/lineage", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
