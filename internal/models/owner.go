package models

type SetOwnerRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type OwnerResponse struct {
	Code  int32 `json:"code"`
	Owner Owner `json:"owner"`
}

type Owner struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SetOwnerResponse struct {
	Code int32 `json:"code"`
	Set  bool  `json:"set"`
}
