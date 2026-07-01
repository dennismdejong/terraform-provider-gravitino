package table

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gravitino/terraform-provider-gravitino/internal/client"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*tableDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*tableDataSource)(nil)

type tableDataSource struct {
	client *client.Client
}

func NewTableDataSource() datasource.DataSource {
	return &tableDataSource{}
}

func (d *tableDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	d.client = c
}

func (d *tableDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_table"
}

func (d *tableDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"metalake": schema.StringAttribute{
				Required: true,
			},
			"catalog": schema.StringAttribute{
				Required: true,
			},
			"schema": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"comment": schema.StringAttribute{
				Computed: true,
			},
			"properties": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
		Blocks: map[string]schema.Block{
			"audit": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"creator":            schema.StringAttribute{Computed: true},
					"create_time":        schema.StringAttribute{Computed: true},
					"last_modifier":      schema.StringAttribute{Computed: true},
					"last_modified_time": schema.StringAttribute{Computed: true},
				},
			},
			"column": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":           schema.StringAttribute{Computed: true},
						"type":           schema.StringAttribute{Computed: true},
						"length":         schema.Int64Attribute{Computed: true},
						"precision":      schema.Int64Attribute{Computed: true},
						"scale":          schema.Int64Attribute{Computed: true},
						"comment":        schema.StringAttribute{Computed: true},
						"nullable":       schema.BoolAttribute{Computed: true},
						"auto_increment": schema.BoolAttribute{Computed: true},
						"default_value":  schema.StringAttribute{Computed: true},
					},
				},
			},
			"sort_order": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"field_name": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"direction":     schema.StringAttribute{Computed: true},
						"null_ordering": schema.StringAttribute{Computed: true},
					},
				},
			},
			"distribution": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"strategy": schema.StringAttribute{Computed: true},
					"number":   schema.Int64Attribute{Computed: true},
					"func_args": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"partitioning": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"strategy": schema.StringAttribute{Computed: true},
						"field_name": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"field_names": schema.ListAttribute{
							Computed:    true,
							ElementType: types.ListType{ElemType: types.StringType},
						},
						"num_buckets": schema.Int64Attribute{Computed: true},
						"width":       schema.Int64Attribute{Computed: true},
						"func_name":   schema.StringAttribute{Computed: true},
						"func_args": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"index": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"index_type": schema.StringAttribute{Computed: true},
						"name":       schema.StringAttribute{Computed: true},
						"field_names": schema.ListAttribute{
							Computed:    true,
							ElementType: types.ListType{ElemType: types.StringType},
						},
					},
				},
			},
		},
	}
}

func (d *tableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.TableDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tableResp, err := d.client.GetTable(
		config.Metalake.ValueString(),
		config.Catalog.ValueString(),
		config.Schema.ValueString(),
		config.Name.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read table", err.Error())
		return
	}

	state := models.TableDataSourceModel{
		Metalake: config.Metalake,
		Catalog:  config.Catalog,
		Schema:   config.Schema,
		Name:     config.Name,
	}

	diags := mapTableResponseToDataSourceState(ctx, tableResp, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func mapTableResponseToDataSourceState(ctx context.Context, resp *models.TableResponse, state *models.TableDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	t := resp.Table

	state.Name = types.StringValue(t.Name)
	state.Comment = types.StringValue(t.Comment)

	if t.Properties != nil {
		props, d := types.MapValueFrom(ctx, types.StringType, t.Properties)
		diags.Append(d...)
		state.Properties = props
	} else {
		state.Properties = types.MapNull(types.StringType)
	}

	if t.Audit != nil {
		audit := &models.AuditTFSDK{
			Creator:      types.StringValue(t.Audit.Creator),
			LastModifier: types.StringValue(t.Audit.LastModifier),
		}
		if t.Audit.CreateTime != nil {
			audit.CreateTime = types.StringValue(t.Audit.CreateTime.String())
		} else {
			audit.CreateTime = types.StringNull()
		}
		if t.Audit.LastModifiedTime != nil {
			audit.LastModifiedTime = types.StringValue(t.Audit.LastModifiedTime.String())
		} else {
			audit.LastModifiedTime = types.StringNull()
		}
		state.Audit = audit
	} else {
		state.Audit = nil
	}

	state.Columns = mapDataSourceColumnsToState(ctx, t.Columns, &diags)
	state.SortOrders = mapDataSourceSortOrdersToState(ctx, t.SortOrders, &diags)
	state.Distribution = mapDataSourceDistributionToState(ctx, t.Distribution, &diags)
	state.Partitioning = mapDataSourcePartitioningToState(ctx, t.Partitioning, &diags)
	state.Indexes = mapDataSourceIndexesToState(ctx, t.Indexes, &diags)

	return diags
}

func mapDataSourceColumnsToState(ctx context.Context, cols []models.Column, diags *diag.Diagnostics) []models.ColumnTFSDK {
	result := make([]models.ColumnTFSDK, 0, len(cols))
	for _, col := range cols {
		m := models.ColumnTFSDK{
			Name:          types.StringValue(col.Name),
			Type:          types.StringValue(col.Type.Type),
			Comment:       types.StringValue(col.Comment),
			Nullable:      types.BoolValue(col.Nullable),
			AutoIncrement: types.BoolValue(col.AutoIncrement),
		}

		if col.Type.Length != nil {
			m.Length = types.Int64Value(*col.Type.Length)
		} else {
			m.Length = types.Int64Null()
		}
		if col.Type.Precision != nil {
			m.Precision = types.Int64Value(*col.Type.Precision)
		} else {
			m.Precision = types.Int64Null()
		}
		if col.Type.Scale != nil {
			m.Scale = types.Int64Value(*col.Type.Scale)
		} else {
			m.Scale = types.Int64Null()
		}

		if col.DefaultValue != nil {
			b, err := json.Marshal(col.DefaultValue)
			if err == nil {
				m.DefaultValue = types.StringValue(string(b))
			} else {
				m.DefaultValue = types.StringNull()
			}
		} else {
			m.DefaultValue = types.StringNull()
		}

		result = append(result, m)
	}
	return result
}

func mapDataSourceSortOrdersToState(ctx context.Context, orders []models.SortOrder, diags *diag.Diagnostics) []models.SortOrderTFSDK {
	result := make([]models.SortOrderTFSDK, 0, len(orders))
	for _, o := range orders {
		m := models.SortOrderTFSDK{
			Direction: types.StringValue(o.Direction),
		}

		if o.NullOrdering != "" {
			m.NullOrdering = types.StringValue(o.NullOrdering)
		} else {
			m.NullOrdering = types.StringNull()
		}

		fieldVals := make([]attr.Value, 0, len(o.SortTerm.FieldName))
		for _, fn := range o.SortTerm.FieldName {
			fieldVals = append(fieldVals, types.StringValue(fn))
		}
		listVal, d := types.ListValue(types.StringType, fieldVals)
		diags.Append(d...)
		m.FieldName = listVal

		result = append(result, m)
	}
	return result
}

func mapDataSourceDistributionToState(ctx context.Context, dist *models.Distribution, diags *diag.Diagnostics) *models.DistributionTFSDK {
	if dist == nil {
		return nil
	}

	m := &models.DistributionTFSDK{
		Strategy: types.StringValue(dist.Strategy),
		Number:   types.Int64Value(int64(dist.Number)),
	}

	funcArgVals := make([]attr.Value, 0, len(dist.FuncArgs))
	for _, expr := range dist.FuncArgs {
		funcArgVals = append(funcArgVals, types.StringValue(strings.Join(expr.FieldName, ".")))
	}
	listVal, d := types.ListValue(types.StringType, funcArgVals)
	diags.Append(d...)
	m.FuncArgs = listVal

	return m
}

func mapDataSourcePartitioningToState(ctx context.Context, parts []models.Partitioning, diags *diag.Diagnostics) []models.PartitioningTFSDK {
	result := make([]models.PartitioningTFSDK, 0, len(parts))
	for _, p := range parts {
		m := models.PartitioningTFSDK{
			Strategy: types.StringValue(p.Strategy),
		}

		if len(p.FieldName) > 0 {
			fieldVals := make([]attr.Value, 0, len(p.FieldName))
			for _, fn := range p.FieldName {
				fieldVals = append(fieldVals, types.StringValue(fn))
			}
			listVal, d := types.ListValue(types.StringType, fieldVals)
			diags.Append(d...)
			m.FieldName = listVal
		} else {
			m.FieldName = types.ListNull(types.StringType)
		}

		if len(p.FieldNames) > 0 {
			fieldNamesVals := make([]attr.Value, 0, len(p.FieldNames))
			for _, fns := range p.FieldNames {
				innerVals := make([]attr.Value, 0, len(fns))
				for _, fn := range fns {
					innerVals = append(innerVals, types.StringValue(fn))
				}
				innerList, d := types.ListValue(types.StringType, innerVals)
				diags.Append(d...)
				fieldNamesVals = append(fieldNamesVals, innerList)
			}
			listVal, d := types.ListValue(types.ListType{ElemType: types.StringType}, fieldNamesVals)
			diags.Append(d...)
			m.FieldNames = listVal
		} else {
			m.FieldNames = types.ListNull(types.ListType{ElemType: types.StringType})
		}

		if p.NumBuckets > 0 {
			m.NumBuckets = types.Int64Value(int64(p.NumBuckets))
		} else {
			m.NumBuckets = types.Int64Null()
		}

		if p.Width > 0 {
			m.Width = types.Int64Value(int64(p.Width))
		} else {
			m.Width = types.Int64Null()
		}

		if p.FuncName != "" {
			m.FuncName = types.StringValue(p.FuncName)
		} else {
			m.FuncName = types.StringNull()
		}

		if len(p.FuncArgs) > 0 {
			funcArgVals := make([]attr.Value, 0, len(p.FuncArgs))
			for _, expr := range p.FuncArgs {
				funcArgVals = append(funcArgVals, types.StringValue(strings.Join(expr.FieldName, ".")))
			}
			listVal, d := types.ListValue(types.StringType, funcArgVals)
			diags.Append(d...)
			m.FuncArgs = listVal
		} else {
			m.FuncArgs = types.ListNull(types.StringType)
		}

		result = append(result, m)
	}
	return result
}

func mapDataSourceIndexesToState(ctx context.Context, indexes []models.Index, diags *diag.Diagnostics) []models.IndexTFSDK {
	result := make([]models.IndexTFSDK, 0, len(indexes))
	for _, idx := range indexes {
		m := models.IndexTFSDK{
			IndexType: types.StringValue(idx.IndexType),
		}

		if idx.Name != "" {
			m.Name = types.StringValue(idx.Name)
		} else {
			m.Name = types.StringNull()
		}

		fieldNamesVals := make([]attr.Value, 0, len(idx.FieldNames))
		for _, fns := range idx.FieldNames {
			innerVals := make([]attr.Value, 0, len(fns))
			for _, fn := range fns {
				innerVals = append(innerVals, types.StringValue(fn))
			}
			innerList, d := types.ListValue(types.StringType, innerVals)
			diags.Append(d...)
			fieldNamesVals = append(fieldNamesVals, innerList)
		}
		listVal, d := types.ListValue(types.ListType{ElemType: types.StringType}, fieldNamesVals)
		diags.Append(d...)
		m.FieldNames = listVal

		result = append(result, m)
	}
	return result
}
