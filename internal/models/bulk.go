package models

type BulkUserAddRequest struct {
	Users []UserCreateRequest `json:"users"`
}

type BulkGroupAddRequest struct {
	Groups []GroupCreateRequest `json:"groups"`
}

type BulkRemoveRequest struct {
	Names []string `json:"names"`
}

type BulkError struct {
	Index   int32  `json:"index"`
	Name    string `json:"name,omitempty"`
	Code    int32  `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type BulkSummary struct {
	Total     int32 `json:"total"`
	Succeeded int32 `json:"succeeded"`
	Failed    int32 `json:"failed"`
}

type BulkUserResponse struct {
	Code    int32        `json:"code"`
	Users   []User       `json:"users"`
	Errors  []BulkError  `json:"errors,omitempty"`
	Summary *BulkSummary `json:"summary,omitempty"`
}

type BulkGroupResponse struct {
	Code    int32        `json:"code"`
	Groups  []Group      `json:"groups"`
	Errors  []BulkError  `json:"errors,omitempty"`
	Summary *BulkSummary `json:"summary,omitempty"`
}

type BulkRemoveResponse struct {
	Code    int32        `json:"code"`
	Names   []string     `json:"names"`
	Errors  []BulkError  `json:"errors,omitempty"`
	Summary *BulkSummary `json:"summary,omitempty"`
}
