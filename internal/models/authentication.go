package models

type Principal struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles,omitempty"`
}

type PrincipalResponse struct {
	Code      int       `json:"code"`
	Principal Principal `json:"principal"`
}
