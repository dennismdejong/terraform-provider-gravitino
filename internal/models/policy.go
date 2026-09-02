package models

type PolicyContent struct {
	SupportedObjectTypes []string          `json:"supportedObjectTypes"`
	Properties           map[string]string `json:"properties,omitempty"`
	CustomRules          map[string]string `json:"customRules,omitempty"`
}

type Policy struct {
	Name       string         `json:"name"`
	Comment    string         `json:"comment,omitempty"`
	PolicyType string         `json:"policyType"`
	Enabled    bool           `json:"enabled"`
	Content    *PolicyContent `json:"content,omitempty"`
	Inherited  *bool          `json:"inherited,omitempty"`
	Audit      *Audit         `json:"audit,omitempty"`
}

type PolicyListResponse struct {
	Code     int32    `json:"code"`
	Policies []Policy `json:"policies"`
}

type PolicyResponse struct {
	Code   int32  `json:"code"`
	Policy Policy `json:"policy"`
}

type PolicyCreateRequest struct {
	Name       string         `json:"name"`
	Comment    string         `json:"comment,omitempty"`
	PolicyType string         `json:"policyType"`
	Enabled    bool           `json:"enabled"`
	Content    *PolicyContent `json:"content,omitempty"`
}

type PolicyAssociationRequest struct {
	PoliciesToAdd    []string `json:"policiesToAdd,omitempty"`
	PoliciesToRemove []string `json:"policiesToRemove,omitempty"`
}
