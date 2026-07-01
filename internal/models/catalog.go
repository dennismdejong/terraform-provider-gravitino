package models

type Catalog struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Provider   string            `json:"provider"`
	Comment    string            `json:"comment,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Audit      *Audit            `json:"audit,omitempty"`
}

type CatalogListResponse struct {
	Code        int              `json:"code"`
	Identifiers []NameIdentifier `json:"identifiers"`
}

type CatalogInfoListResponse struct {
	Code     int       `json:"code"`
	Catalogs []Catalog `json:"catalogs"`
}

type CatalogResponse struct {
	Code    int     `json:"code"`
	Catalog Catalog `json:"catalog"`
}

type CatalogCreateRequest struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Provider   string            `json:"provider,omitempty"`
	Comment    string            `json:"comment,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type CatalogUpdate struct {
	Type string `json:"@type"`
}

type RenameCatalogRequest struct {
	CatalogUpdate
	NewName string `json:"newName"`
}

type UpdateCatalogCommentRequest struct {
	CatalogUpdate
	NewComment string `json:"newComment"`
}

type SetCatalogPropertyRequest struct {
	CatalogUpdate
	Property string `json:"property"`
	Value    string `json:"value"`
}

type RemoveCatalogPropertyRequest struct {
	CatalogUpdate
	Property string `json:"property"`
}

type CatalogUpdateRequest struct {
	Updates []interface{} `json:"updates"`
}

func NewRenameCatalogRequest(newName string) interface{} {
	return struct {
		Type    string `json:"@type"`
		NewName string `json:"newName"`
	}{
		Type:    "rename",
		NewName: newName,
	}
}

func NewUpdateCatalogCommentRequest(newComment string) interface{} {
	return struct {
		Type       string `json:"@type"`
		NewComment string `json:"newComment"`
	}{
		Type:       "updateComment",
		NewComment: newComment,
	}
}

func NewSetCatalogPropertyRequest(property, value string) interface{} {
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

func NewRemoveCatalogPropertyRequest(property string) interface{} {
	return struct {
		Type     string `json:"@type"`
		Property string `json:"property"`
	}{
		Type:     "removeProperty",
		Property: property,
	}
}
