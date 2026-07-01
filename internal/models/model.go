package models

type Model struct {
	Name       string            `json:"name"`
	Comment    string            `json:"comment,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	ModelURI   string            `json:"modelUri,omitempty"`
	Audit      *Audit            `json:"audit,omitempty"`
}

type ModelResponse struct {
	Code  int   `json:"code"`
	Model Model `json:"model"`
}

type ModelCreateRequest struct {
	Name       string            `json:"name"`
	Comment    string            `json:"comment,omitempty"`
	ModelURI   string            `json:"modelUri,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type ModelUpdateRequest struct {
	Updates []interface{} `json:"updates"`
}

func NewRenameModelRequest(newName string) interface{} {
	return struct {
		Type    string `json:"@type"`
		NewName string `json:"newName"`
	}{
		Type:    "rename",
		NewName: newName,
	}
}

func NewUpdateModelCommentRequest(newComment string) interface{} {
	return struct {
		Type       string `json:"@type"`
		NewComment string `json:"newComment"`
	}{
		Type:       "updateComment",
		NewComment: newComment,
	}
}

func NewSetModelPropertyRequest(property, value string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
		Value    string `json:"value"`
	}{
		Type:     "setProperty",
		Property: property,
		Value:    value,
	}
}

func NewRemoveModelPropertyRequest(property string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
	}{
		Type:     "removeProperty",
		Property: property,
	}
}
