package client

import (
	"fmt"
	"net/url"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func (c *Client) ListJobs(metalake string) (*models.JobListResponse, error) {
	var result models.JobListResponse
	err := c.Get("/metalakes/"+url.PathEscape(metalake)+"/jobs", &result)
	return &result, err
}

func (c *Client) GetJob(metalake, name string) (*models.JobResponse, error) {
	var result models.JobResponse
	err := c.Get("/metalakes/"+url.PathEscape(metalake)+"/jobs/"+url.PathEscape(name), &result)
	return &result, err
}

func (c *Client) CreateJob(metalake string, req *models.JobCreateRequest) (*models.JobResponse, error) {
	var result models.JobResponse
	err := c.Post("/metalakes/"+url.PathEscape(metalake)+"/jobs", req, &result)
	return &result, err
}

func (c *Client) DeleteJob(metalake, name string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/jobs/%s", url.PathEscape(metalake), url.PathEscape(name))
	var result models.DropResponse
	err := c.Delete(path, &result)
	return &result, err
}

func (c *Client) ListJobTemplates(metalake string) (*models.JobTemplateListResponse, error) {
	var result models.JobTemplateListResponse
	err := c.Get("/metalakes/"+url.PathEscape(metalake)+"/jobs/templates", &result)
	return &result, err
}

func (c *Client) RegisterJobTemplate(metalake string, req *models.JobTemplateCreateRequest) (*models.JobTemplateResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/jobs/templates", metalake)
	var result models.JobTemplateResponse
	if err := c.Post(path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetJobTemplate(metalake, name string) (*models.JobTemplateResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/jobs/templates/%s", metalake, name)
	var result models.JobTemplateResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateJobTemplate(metalake, name string, req *models.JobTemplateCreateRequest) (*models.JobTemplateResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/jobs/templates/%s", metalake, name)
	var result models.JobTemplateResponse
	if err := c.Put(path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteJobTemplate(metalake, name string) (*models.DropResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/jobs/templates/%s", metalake, name)
	var result models.DropResponse
	if err := c.Delete(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) PauseJob(metalake, name string) (*models.JobResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/jobs/%s/pause", url.PathEscape(metalake), url.PathEscape(name))
	var result models.JobResponse
	err := c.Post(path, nil, &result)
	return &result, err
}

func (c *Client) ResumeJob(metalake, name string) (*models.JobResponse, error) {
	path := fmt.Sprintf("/metalakes/%s/jobs/%s/resume", url.PathEscape(metalake), url.PathEscape(name))
	var result models.JobResponse
	err := c.Post(path, nil, &result)
	return &result, err
}
