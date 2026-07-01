package models

type Fileset struct {
	Name            string            `json:"name"`
	Comment         string            `json:"comment,omitempty"`
	Type            string            `json:"type,omitempty"`
	StorageLocation string            `json:"storageLocation,omitempty"`
	Properties      map[string]string `json:"properties,omitempty"`
	Audit           *Audit            `json:"audit,omitempty"`
}

type FilesetResponse struct {
	Code    int     `json:"code"`
	Fileset Fileset `json:"fileset"`
}

type FilesetCreateRequest struct {
	Name            string            `json:"name"`
	Comment         string            `json:"comment,omitempty"`
	Type            string            `json:"type,omitempty"`
	StorageLocation string            `json:"storageLocation,omitempty"`
	Properties      map[string]string `json:"properties,omitempty"`
}

type FilesetUpdateRequest struct {
	Updates []interface{} `json:"updates"`
}

func NewRenameFilesetRequest(newName string) interface{} {
	return struct {
		Type    string `json:"@type"`
		NewName string `json:"newName"`
	}{Type: "rename", NewName: newName}
}

func NewUpdateFilesetCommentRequest(newComment string) interface{} {
	return struct {
		Type       string `json:"@type"`
		NewComment string `json:"newComment"`
	}{Type: "updateComment", NewComment: newComment}
}

func NewSetFilesetPropertyRequest(property, value string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
		Value    string `json:"value"`
	}{Type: "setProperty", Property: property, Value: value}
}

func NewRemoveFilesetPropertyRequest(property string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
	}{Type: "removeProperty", Property: property}
}
