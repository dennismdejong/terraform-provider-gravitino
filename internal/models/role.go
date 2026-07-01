package models

type Role struct {
	Name            string   `json:"name"`
	Privileges      []string `json:"privileges,omitempty"`
	SecurableObject string   `json:"securableObject,omitempty"`
}

type RoleListResponse struct {
	Code  int    `json:"code"`
	Roles []Role `json:"roles"`
}
