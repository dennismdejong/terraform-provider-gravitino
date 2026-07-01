package models

type GroupCreateRequest struct {
	Name string `json:"name"`
}

type GroupResponse struct {
	Code  int32 `json:"code"`
	Group Group `json:"group"`
}

type Group struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles,omitempty"`
	Audit *Audit   `json:"audit,omitempty"`
}
