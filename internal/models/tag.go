package models

type Tag struct {
	Name       string            `json:"name"`
	Comment    string            `json:"comment,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Audit      *Audit            `json:"audit,omitempty"`
	Inherited  *bool             `json:"inherited,omitempty"`
}

type TagListResponse struct {
	Code int   `json:"code"`
	Tags []Tag `json:"tags"`
}

type TagNameListResponse struct {
	Code  int      `json:"code"`
	Names []string `json:"names"`
}

type TagResponse struct {
	Code int `json:"code"`
	Tag  Tag `json:"tag"`
}

type TagCreateRequest struct {
	Name       string            `json:"name"`
	Comment    string            `json:"comment,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type TagUpdateRequest struct {
	Updates []interface{} `json:"updates"`
}

func NewRenameTagRequest(newName string) interface{} {
	return struct {
		Type    string `json:"@type"`
		NewName string `json:"newName"`
	}{
		Type:    "rename",
		NewName: newName,
	}
}

func NewUpdateTagCommentRequest(newComment string) interface{} {
	return struct {
		Type       string `json:"@type"`
		NewComment string `json:"newComment"`
	}{
		Type:       "updateComment",
		NewComment: newComment,
	}
}

func NewSetTagPropertyRequest(property, value string) interface{} {
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

func NewRemoveTagPropertyRequest(property string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
	}{
		Type:     "removeProperty",
		Property: property,
	}
}

type TagAssociationRequest struct {
	Tags []string `json:"tags"`
}
