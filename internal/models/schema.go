package models

type Schema struct {
	Name       string            `json:"name"`
	Comment    string            `json:"comment,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Audit      *Audit            `json:"audit,omitempty"`
}

type SchemaResponse struct {
	Code   int    `json:"code"`
	Schema Schema `json:"schema"`
}

type SchemaCreateRequest struct {
	Name       string            `json:"name"`
	Comment    string            `json:"comment,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type SchemaUpdate struct {
	Type string `json:"@type"`
}

type SchemaUpdateRequest struct {
	Updates []interface{} `json:"updates"`
}

func NewRenameSchemaRequest(newName string) interface{} {
	return struct {
		Type    string `json:"@type"`
		NewName string `json:"newName"`
	}{
		Type:    "rename",
		NewName: newName,
	}
}

func NewUpdateSchemaCommentRequest(newComment string) interface{} {
	return struct {
		Type       string `json:"@type"`
		NewComment string `json:"newComment"`
	}{
		Type:       "updateComment",
		NewComment: newComment,
	}
}

func NewSetSchemaPropertyRequest(property, value string) interface{} {
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

func NewRemoveSchemaPropertyRequest(property string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
	}{
		Type:     "removeProperty",
		Property: property,
	}
}
