package models

type RoleCreateRequest struct {
	Name             string            `json:"name"`
	Properties       map[string]string `json:"properties,omitempty"`
	SecurableObjects []SecurableObject `json:"securableObjects,omitempty"`
}

type RoleResponse struct {
	Code int32      `json:"code"`
	Role RoleDetail `json:"role"`
}

type RoleDetail struct {
	Name             string            `json:"name"`
	Properties       map[string]string `json:"properties,omitempty"`
	SecurableObjects []SecurableObject `json:"securableObjects,omitempty"`
	Audit            *Audit            `json:"audit,omitempty"`
}

type SecurableObject struct {
	FullName   string      `json:"fullName"`
	Type       string      `json:"type"`
	Privileges []Privilege `json:"privileges"`
}

type Privilege struct {
	Name      string `json:"name"`
	Condition string `json:"condition"`
}

type RoleListResponse struct {
	Code  int32  `json:"code"`
	Roles []Role `json:"roles"`
}

type Role struct {
	Name            string   `json:"name"`
	SecurableObject string   `json:"securableObject"`
	Privileges      []string `json:"privileges"`
	Audit           *Audit   `json:"audit,omitempty"`
}

type PrivilegesRequest struct {
	Privileges []Privilege `json:"privileges"`
}

type PrivilegeOverrideRequest struct {
	Overrides []SecurableObject `json:"overrides"`
}
