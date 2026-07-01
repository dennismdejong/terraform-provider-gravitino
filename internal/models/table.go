package models

type Column struct {
	Name          string      `json:"name"`
	Type          DataType    `json:"type"`
	Comment       string      `json:"comment,omitempty"`
	Nullable      bool        `json:"nullable,omitempty"`
	AutoIncrement bool        `json:"autoIncrement,omitempty"`
	DefaultValue  interface{} `json:"defaultValue,omitempty"`
}

type DataType struct {
	Type         string        `json:"type"`
	Length       *int64        `json:"length,omitempty"`
	Precision    *int64        `json:"precision,omitempty"`
	Scale        *int64        `json:"scale,omitempty"`
	StructType   *StructType   `json:"struct,omitempty"`
	ListType     *ListType     `json:"list,omitempty"`
	MapType      *MapType      `json:"map,omitempty"`
	UnionType    *UnionType    `json:"union,omitempty"`
	UnparsedType *UnparsedType `json:"unparsed,omitempty"`
}

type StructType struct {
	Fields []StructField `json:"fields"`
}

type StructField struct {
	Name     string   `json:"name"`
	Type     DataType `json:"type"`
	Nullable bool     `json:"nullable,omitempty"`
	Comment  string   `json:"comment,omitempty"`
}

type ListType struct {
	ContainsNull bool     `json:"containsNull,omitempty"`
	ElementType  DataType `json:"elementType"`
}

type MapType struct {
	KeyType           DataType `json:"keyType"`
	ValueType         DataType `json:"valueType"`
	ValueContainsNull bool     `json:"valueContainsNull,omitempty"`
}

type UnionType struct {
	Types []DataType `json:"types"`
}

type UnparsedType struct {
	UnparsedType string `json:"unparsedType"`
}

type SortOrder struct {
	SortTerm     Expression `json:"sortTerm"`
	Direction    string     `json:"direction"`
	NullOrdering string     `json:"nullOrdering,omitempty"`
}

type Expression struct {
	Type      string        `json:"type"`
	DataType  *DataType     `json:"dataType,omitempty"`
	Value     string        `json:"value,omitempty"`
	FieldName []string      `json:"fieldName,omitempty"`
	FuncName  string        `json:"funcName,omitempty"`
	FuncArgs  []interface{} `json:"funcArgs,omitempty"`
}

type Distribution struct {
	Strategy string       `json:"strategy"`
	Number   int32        `json:"number"`
	FuncArgs []Expression `json:"funcArgs,omitempty"`
}

type Partitioning struct {
	Strategy    string                `json:"strategy"`
	FieldName   []string              `json:"fieldName,omitempty"`
	FieldNames  [][]string            `json:"fieldNames,omitempty"`
	NumBuckets  int                   `json:"numBuckets,omitempty"`
	Width       int                   `json:"width,omitempty"`
	Assignments []PartitionAssignment `json:"assignments,omitempty"`
	FuncName    string                `json:"funcName,omitempty"`
	FuncArgs    []Expression          `json:"funcArgs,omitempty"`
}

type PartitionAssignment struct {
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	Lists      [][]interface{}   `json:"lists,omitempty"`
	Upper      *Literal          `json:"upper,omitempty"`
	Lower      *Literal          `json:"lower,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type Literal struct {
	Type     string    `json:"type"`
	DataType *DataType `json:"dataType"`
	Value    string    `json:"value"`
}

type Index struct {
	IndexType  string     `json:"indexType"`
	Name       string     `json:"name,omitempty"`
	FieldNames [][]string `json:"fieldNames"`
}

type Table struct {
	Name         string            `json:"name"`
	Columns      []Column          `json:"columns"`
	Comment      string            `json:"comment,omitempty"`
	Properties   map[string]string `json:"properties,omitempty"`
	Audit        *Audit            `json:"audit,omitempty"`
	SortOrders   []SortOrder       `json:"sortOrders,omitempty"`
	Distribution *Distribution     `json:"distribution,omitempty"`
	Partitioning []Partitioning    `json:"partitioning,omitempty"`
	Indexes      []Index           `json:"indexes,omitempty"`
}

type TableResponse struct {
	Code  int   `json:"code"`
	Table Table `json:"table"`
}

type TableCreateRequest struct {
	Name         string            `json:"name"`
	Columns      []Column          `json:"columns"`
	Comment      string            `json:"comment,omitempty"`
	Properties   map[string]string `json:"properties,omitempty"`
	SortOrders   []SortOrder       `json:"sortOrders,omitempty"`
	Distribution *Distribution     `json:"distribution,omitempty"`
	Partitioning []Partitioning    `json:"partitioning,omitempty"`
	Indexes      []Index           `json:"indexes,omitempty"`
}

type TableUpdate struct {
	Type string `json:"@type"`
}

type TableUpdateRequest struct {
	Updates []interface{} `json:"updates"`
}

func NewRenameTableRequest(newName string) interface{} {
	return struct {
		Type    string `json:"@type"`
		NewName string `json:"newName"`
	}{
		Type:    "rename",
		NewName: newName,
	}
}

func NewUpdateTableCommentRequest(newComment string) interface{} {
	return struct {
		Type       string `json:"@type"`
		NewComment string `json:"newComment"`
	}{
		Type:       "updateComment",
		NewComment: newComment,
	}
}

func NewSetTablePropertyRequest(property, value string) interface{} {
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

func NewRemoveTablePropertyRequest(property string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
	}{
		Type:     "removeProperty",
		Property: property,
	}
}
