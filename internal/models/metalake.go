package models

type Metalake struct {
	Name       string            `json:"name"`
	Comment    string            `json:"comment,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Audit      *Audit            `json:"audit,omitempty"`
}

type MetalakeListResponse struct {
	Code      int        `json:"code"`
	Metalakes []Metalake `json:"metalakes"`
}

type MetalakeResponse struct {
	Code     int      `json:"code"`
	Metalake Metalake `json:"metalake"`
}

type MetalakeCreateRequest struct {
	Name       string            `json:"name"`
	Comment    string            `json:"comment,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type MetalakeUpdate struct {
	Type string `json:"@type"`
}

type RenameMetalakeRequest struct {
	MetalakeUpdate
	NewName string `json:"newName"`
}

type UpdateMetalakeCommentRequest struct {
	MetalakeUpdate
	NewComment string `json:"newComment"`
}

type SetMetalakePropertyRequest struct {
	MetalakeUpdate
	Property string `json:"property"`
	Value    string `json:"value"`
}

type RemoveMetalakePropertyRequest struct {
	MetalakeUpdate
	Property string `json:"property"`
}

type MetalakeUpdateRequest struct {
	Updates []interface{} `json:"updates"`
}

func NewRenameMetalakeRequest(newName string) interface{} {
	return struct {
		Type    string `json:"@type"`
		NewName string `json:"newName"`
	}{
		Type:    "rename",
		NewName: newName,
	}
}

func NewUpdateMetalakeCommentRequest(newComment string) interface{} {
	return struct {
		Type       string `json:"@type"`
		NewComment string `json:"newComment"`
	}{
		Type:       "updateComment",
		NewComment: newComment,
	}
}

func NewSetMetalakePropertyRequest(property, value string) interface{} {
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

func NewRemoveMetalakePropertyRequest(property string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
	}{
		Type:     "removeProperty",
		Property: property,
	}
}
