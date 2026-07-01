package models

type View struct {
	Name       string            `json:"name"`
	Comment    string            `json:"comment,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	ViewDef    string            `json:"viewDef,omitempty"`
	Audit      *Audit            `json:"audit,omitempty"`
}

type ViewResponse struct {
	Code int  `json:"code"`
	View View `json:"view"`
}

type ViewCreateRequest struct {
	Name       string            `json:"name"`
	Comment    string            `json:"comment,omitempty"`
	ViewDef    string            `json:"viewDef,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type ViewUpdateRequest struct {
	Updates []interface{} `json:"updates"`
}

func NewRenameViewRequest(newName string) interface{} {
	return struct {
		Type    string `json:"@type"`
		NewName string `json:"newName"`
	}{Type: "rename", NewName: newName}
}

func NewUpdateViewCommentRequest(newComment string) interface{} {
	return struct {
		Type       string `json:"@type"`
		NewComment string `json:"newComment"`
	}{Type: "updateComment", NewComment: newComment}
}

func NewSetViewPropertyRequest(property, value string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
		Value    string `json:"value"`
	}{Type: "setProperty", Property: property, Value: value}
}

func NewRemoveViewPropertyRequest(property string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
	}{Type: "removeProperty", Property: property}
}
