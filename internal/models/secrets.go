package models

type SecretsResponse struct {
	Code    int32             `json:"code"`
	Secrets map[string]string `json:"secrets"`
}
