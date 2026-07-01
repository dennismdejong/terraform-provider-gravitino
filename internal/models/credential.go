package models

type Credential struct {
	Type       string `json:"type,omitempty"`
	Value      string `json:"value,omitempty"`
	ExpireTime string `json:"expireTime,omitempty"`
}

type CredentialResponse struct {
	Code       int        `json:"code"`
	Credential Credential `json:"credential"`
}
