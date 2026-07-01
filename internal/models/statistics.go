package models

type Statistics struct {
	Name       string            `json:"name,omitempty"`
	Type       string            `json:"type,omitempty"`
	Value      string            `json:"value,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type StatisticsResponse struct {
	Code       int          `json:"code"`
	Statistics []Statistics `json:"statistics"`
}

type PartitionStatistics struct {
	PartitionName string       `json:"partitionName,omitempty"`
	Statistics    []Statistics `json:"statistics,omitempty"`
}

type PartitionStatisticsResponse struct {
	Code       int                   `json:"code"`
	Statistics []PartitionStatistics `json:"partitionStatistics,omitempty"`
}
