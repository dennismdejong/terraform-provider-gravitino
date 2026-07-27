package table

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

func TestTypeNameToDataType(t *testing.T) {
	tests := []struct {
		name      string
		typeName  string
		length    int64
		precision int64
		scale     int64
		wantType  string
		wantLen   *int64
		wantPrec  *int64
		wantScale *int64
	}{
		{"integer", "integer", 0, 0, 0, "integer", nil, nil, nil},
		{"string", "string", 0, 0, 0, "string", nil, nil, nil},
		{"long", "long", 0, 0, 0, "long", nil, nil, nil},
		{"varchar_with_length", "varchar", 255, 0, 0, "varchar", i64ptr(255), nil, nil},
		{"decimal", "decimal", 0, 10, 2, "decimal", nil, i64ptr(10), i64ptr(2)},
		{"fixed", "fixed", 16, 0, 0, "fixed", i64ptr(16), nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := typeNameToDataType(tt.typeName, tt.length, tt.precision, tt.scale)
			if result.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", result.Type, tt.wantType)
			}
			if tt.wantLen != nil {
				if result.Length == nil || *result.Length != *tt.wantLen {
					t.Errorf("Length = %v, want %v", result.Length, tt.wantLen)
				}
			}
			if tt.wantPrec != nil {
				if result.Precision == nil || *result.Precision != *tt.wantPrec {
					t.Errorf("Precision = %v, want %v", result.Precision, tt.wantPrec)
				}
			}
			if tt.wantScale != nil {
				if result.Scale == nil || *result.Scale != *tt.wantScale {
					t.Errorf("Scale = %v, want %v", result.Scale, tt.wantScale)
				}
			}
		})
	}
}

func TestColumnTypeToString(t *testing.T) {
	tests := []struct {
		name     string
		dataType models.DataType
		want     string
	}{
		{"integer", models.DataType{Type: "integer"}, "integer"},
		{"varchar", models.DataType{Type: "varchar", Length: i64ptr(255)}, "varchar"},
		{"decimal", models.DataType{Type: "decimal", Precision: i64ptr(10), Scale: i64ptr(2)}, "decimal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := columnTypeToString(tt.dataType)
			if result != tt.want {
				t.Errorf("columnTypeToString() = %q, want %q", result, tt.want)
			}
		})
	}
}

func TestClientCreateTable(t *testing.T) {
	now := time.Now()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/api/metalakes/ml/catalogs/cat/schemas/sch/tables" {
			var req models.TableCreateRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			resp := models.TableResponse{
				Code: 0,
				Table: models.Table{
					Name:       req.Name,
					Columns:    req.Columns,
					Comment:    req.Comment,
					Properties: req.Properties,
					Audit: &models.Audit{
						Creator:          "test-user",
						CreateTime:       &now,
						LastModifier:     "test-user",
						LastModifiedTime: &now,
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	c, err := client.New(srv.URL, nil)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	length := int64(255)
	req := &models.TableCreateRequest{
		Name:    "test_table",
		Comment: "test comment",
		Columns: []models.Column{
			{
				Name:     "id",
				Type:     models.DataType{Type: "integer"},
				Nullable: false,
			},
			{
				Name:     "name",
				Type:     models.DataType{Type: "varchar", Length: &length},
				Nullable: true,
			},
		},
	}

	got, err := c.CreateTable("ml", "cat", "sch", req)
	if err != nil {
		t.Fatalf("CreateTable() error = %v", err)
	}
	if got.Table.Name != "test_table" {
		t.Errorf("Name = %q, want 'test_table'", got.Table.Name)
	}
	if len(got.Table.Columns) != 2 {
		t.Errorf("Columns len = %d, want 2", len(got.Table.Columns))
	}
}

func TestClientGetTable(t *testing.T) {
	now := time.Now()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/api/metalakes/ml/catalogs/cat/schemas/sch/tables/test_table" {
			resp := models.TableResponse{
				Code: 0,
				Table: models.Table{
					Name:    "test_table",
					Comment: "a test table",
					Audit: &models.Audit{
						Creator:          "admin",
						CreateTime:       &now,
						LastModifier:     "admin",
						LastModifiedTime: &now,
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	c, err := client.New(srv.URL, nil)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got, err := c.GetTable("ml", "cat", "sch", "test_table")
	if err != nil {
		t.Fatalf("GetTable() error = %v", err)
	}
	if got.Table.Name != "test_table" {
		t.Errorf("Name = %q, want 'test_table'", got.Table.Name)
	}
	if got.Table.Audit == nil || got.Table.Audit.Creator != "admin" {
		t.Error("Audit not correctly populated")
	}
}

func TestClientListTables(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/api/metalakes/ml/catalogs/cat/schemas/sch/tables" {
			resp := models.IdentifiersResponse{
				Code: 0,
				Identifiers: []models.NameIdentifier{
					{Namespace: []string{"ml", "cat", "sch"}, Name: "table1"},
					{Namespace: []string{"ml", "cat", "sch"}, Name: "table2"},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	c, err := client.New(srv.URL, nil)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got, err := c.ListTables("ml", "cat", "sch")
	if err != nil {
		t.Fatalf("ListTables() error = %v", err)
	}
	if len(got.Identifiers) != 2 {
		t.Errorf("Identifiers len = %d, want 2", len(got.Identifiers))
	}
	if got.Identifiers[0].Name != "table1" {
		t.Errorf("First table = %q, want 'table1'", got.Identifiers[0].Name)
	}
}

func TestClientDropTable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			resp := models.DropResponse{Code: 0, Dropped: true}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	c, err := client.New(srv.URL, nil)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got, err := c.DropTable("ml", "cat", "sch", "test_table", true, false)
	if err != nil {
		t.Fatalf("DropTable() error = %v", err)
	}
	if !got.Dropped {
		t.Error("expected Dropped = true")
	}
}

func i64ptr(v int64) *int64 { return &v }
