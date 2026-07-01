package models

type Partition struct {
	Name       string            `json:"name"`
	Properties map[string]string `json:"properties,omitempty"`
	Audit      *Audit            `json:"audit,omitempty"`
}

type PartitionResponse struct {
	Code      int       `json:"code"`
	Partition Partition `json:"partition"`
}

type PartitionCreateRequest struct {
	Name       string            `json:"name"`
	Properties map[string]string `json:"properties,omitempty"`
}

type PartitionUpdateRequest struct {
	Updates []interface{} `json:"updates"`
}

func NewRenamePartitionRequest(newName string) interface{} {
	return struct {
		Type    string `json:"@type"`
		NewName string `json:"newName"`
	}{
		Type:    "rename",
		NewName: newName,
	}
}

func NewSetPartitionPropertyRequest(property, value string) interface{} {
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

func NewRemovePartitionPropertyRequest(property string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
	}{
		Type:     "removeProperty",
		Property: property,
	}
}
