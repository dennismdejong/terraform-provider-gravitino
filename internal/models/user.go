package models

type UserCreateRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	Code int32 `json:"code"`
	User User  `json:"user"`
}

type User struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles,omitempty"`
	Audit *Audit   `json:"audit,omitempty"`
}

type NameListResponse struct {
	Code  int32    `json:"code"`
	Names []string `json:"names"`
}

type GrantRevokeRequest struct {
	RoleNames []string `json:"roleNames"`
}
