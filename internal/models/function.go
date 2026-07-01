package models

type Function struct {
	Name         string            `json:"name"`
	Comment      string            `json:"comment,omitempty"`
	Properties   map[string]string `json:"properties,omitempty"`
	FunctionBody string            `json:"functionBody,omitempty"`
	Audit        *Audit            `json:"audit,omitempty"`
}

type FunctionResponse struct {
	Code     int      `json:"code"`
	Function Function `json:"function"`
}

type FunctionCreateRequest struct {
	Name         string            `json:"name"`
	Comment      string            `json:"comment,omitempty"`
	FunctionBody string            `json:"functionBody,omitempty"`
	Properties   map[string]string `json:"properties,omitempty"`
}

type FunctionUpdateRequest struct {
	Updates []interface{} `json:"updates"`
}

func NewRenameFunctionRequest(newName string) interface{} {
	return struct {
		Type    string `json:"@type"`
		NewName string `json:"newName"`
	}{Type: "rename", NewName: newName}
}

func NewUpdateFunctionCommentRequest(newComment string) interface{} {
	return struct {
		Type       string `json:"@type"`
		NewComment string `json:"newComment"`
	}{Type: "updateComment", NewComment: newComment}
}

func NewSetFunctionPropertyRequest(property, value string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
		Value    string `json:"value"`
	}{Type: "setProperty", Property: property, Value: value}
}

func NewRemoveFunctionPropertyRequest(property string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
	}{Type: "removeProperty", Property: property}
}
