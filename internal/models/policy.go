package models

type Policy struct {
	Name       string            `json:"name"`
	Condition  string            `json:"condition,omitempty"`
	Effect     string            `json:"effect"`
	Actions    []string          `json:"actions,omitempty"`
	Subjects   []string          `json:"subjects,omitempty"`
	Object     string            `json:"object,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Audit      *Audit            `json:"audit,omitempty"`
}

type PolicyListResponse struct {
	Code     int      `json:"code"`
	Policies []Policy `json:"policies"`
}

type PolicyResponse struct {
	Code   int    `json:"code"`
	Policy Policy `json:"policy"`
}

type PolicyCreateRequest struct {
	Name       string            `json:"name"`
	Condition  string            `json:"condition,omitempty"`
	Effect     string            `json:"effect"`
	Actions    []string          `json:"actions,omitempty"`
	Subjects   []string          `json:"subjects,omitempty"`
	Object     string            `json:"object,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type PolicyUpdateRequest struct {
	Updates []interface{} `json:"updates"`
}

func NewRenamePolicyRequest(newName string) interface{} {
	return struct {
		Type    string `json:"@type"`
		NewName string `json:"newName"`
	}{
		Type:    "rename",
		NewName: newName,
	}
}

func NewSetPolicyPropertyRequest(property, value string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
		Value    string `json:"value"`
	}{
		Type:     "setProperty",
		Property: property,
		Value:    value,
	}
}

func NewRemovePolicyPropertyRequest(property string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
	}{
		Type:     "removeProperty",
		Property: property,
	}
}
