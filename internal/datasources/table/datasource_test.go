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

func TestDataSourceListTables(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/api/metalakes/ml/catalogs/cat/schemas/sch/tables" {
			resp := models.IdentifiersResponse{
				Code: 0,
				Identifiers: []models.NameIdentifier{
					{Namespace: []string{"ml", "cat", "sch"}, Name: "t1"},
					{Namespace: []string{"ml", "cat", "sch"}, Name: "t2"},
					{Namespace: []string{"ml", "cat", "sch"}, Name: "t3"},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	c, err := client.New(srv.URL, "", "", "", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got, err := c.ListTables("ml", "cat", "sch")
	if err != nil {
		t.Fatalf("ListTables() error = %v", err)
	}
	if len(got.Identifiers) != 3 {
		t.Fatalf("expected 3 tables, got %d", len(got.Identifiers))
	}

	names := make([]string, len(got.Identifiers))
	for i, id := range got.Identifiers {
		names[i] = id.Name
	}
	expected := []string{"t1", "t2", "t3"}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("table[%d] = %q, want %q", i, n, expected[i])
		}
	}
}

func TestDataSourceGetTableFullDetail(t *testing.T) {
	now := time.Now()
	length := int64(100)
	precision := int64(10)
	scale := int64(2)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/api/metalakes/ml/catalogs/cat/schemas/sch/tables/full_table" {
			resp := models.TableResponse{
				Code: 0,
				Table: models.Table{
					Name:       "full_table",
					Comment:    "a detailed table",
					Properties: map[string]string{"owner": "admin", "format": "parquet"},
					Audit: &models.Audit{
						Creator:          "admin",
						CreateTime:       &now,
						LastModifier:     "admin",
						LastModifiedTime: &now,
					},
					Columns: []models.Column{
						{
							Name: "id",
							Type: models.DataType{Type: "integer"},
						},
						{
							Name: "email",
							Type: models.DataType{Type: "varchar", Length: &length},
						},
						{
							Name: "price",
							Type: models.DataType{Type: "decimal", Precision: &precision, Scale: &scale},
						},
					},
					SortOrders: []models.SortOrder{
						{
							SortTerm:  models.Expression{Type: "field", FieldName: []string{"id"}},
							Direction: "asc",
						},
					},
					Distribution: &models.Distribution{
						Strategy: "hash",
						Number:   4,
						FuncArgs: []models.Expression{
							{Type: "field", FieldName: []string{"id"}},
						},
					},
					Partitioning: []models.Partitioning{
						{
							Strategy:  "identity",
							FieldName: []string{"created_at"},
						},
						{
							Strategy:   "bucket",
							FieldNames: [][]string{{"id"}, {"name"}},
							NumBuckets: 16,
						},
					},
					Indexes: []models.Index{
						{
							IndexType:  "primary_key",
							Name:       "pk_id",
							FieldNames: [][]string{{"id"}},
						},
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

	c, err := client.New(srv.URL, "", "", "", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got, err := c.GetTable("ml", "cat", "sch", "full_table")
	if err != nil {
		t.Fatalf("GetTable() error = %v", err)
	}

	table := got.Table
	if table.Name != "full_table" {
		t.Errorf("Name = %q", table.Name)
	}
	if table.Comment != "a detailed table" {
		t.Errorf("Comment = %q", table.Comment)
	}
	if len(table.Columns) != 3 {
		t.Errorf("expected 3 columns, got %d", len(table.Columns))
	}
	if table.Columns[1].Type.Type != "varchar" {
		t.Errorf("Column[1] type = %q", table.Columns[1].Type.Type)
	}
	if table.Columns[1].Type.Length == nil || *table.Columns[1].Type.Length != 100 {
		t.Error("Column[1] length not correctly set")
	}
	if len(table.SortOrders) != 1 {
		t.Errorf("expected 1 sort order, got %d", len(table.SortOrders))
	}
	if table.Distribution == nil {
		t.Fatal("expected distribution")
	}
	if table.Distribution.Number != 4 {
		t.Errorf("distribution number = %d", table.Distribution.Number)
	}
	if len(table.Partitioning) != 2 {
		t.Errorf("expected 2 partitionings, got %d", len(table.Partitioning))
	}
	if table.Partitioning[1].NumBuckets != 16 {
		t.Errorf("partitioning[1].NumBuckets = %d", table.Partitioning[1].NumBuckets)
	}
	if len(table.Indexes) != 1 {
		t.Errorf("expected 1 index, got %d", len(table.Indexes))
	}
	if table.Indexes[0].IndexType != "primary_key" {
		t.Errorf("index type = %q", table.Indexes[0].IndexType)
	}
}
