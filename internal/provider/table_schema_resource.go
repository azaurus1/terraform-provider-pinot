package provider

import (
	"context"

	pinot "github.com/azaurus1/go-pinot-api"
	"github.com/azaurus1/go-pinot-api/model"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &tableSchemaResource{}
	_ resource.ResourceWithConfigure   = &tableSchemaResource{}
	_ resource.ResourceWithImportState = &tableSchemaResource{}
)

type tableSchemaResource struct {
	client *pinot.PinotAPIClient
}

func NewTableSchemaResource() resource.Resource {
	return &tableSchemaResource{}

}

type metricFieldSpec struct {
	Name              string                `tfsdk:"name"`
	DataType          string                `tfsdk:"data_type"`
	NotNull           basetypes.BoolValue   `tfsdk:"not_null"`
	TransformFunction basetypes.StringValue `tfsdk:"transform_function"`
}

type dimensionFieldSpec struct {
	Name              string                `tfsdk:"name"`
	DataType          string                `tfsdk:"data_type"`
	NotNull           basetypes.BoolValue   `tfsdk:"not_null"`
	SingleValueField  basetypes.BoolValue   `tfsdk:"single_value_field"`
	TransformFunction basetypes.StringValue `tfsdk:"transform_function"`
}

type dateTimeFieldSpec struct {
	Name              string                `tfsdk:"name"`
	DataType          string                `tfsdk:"data_type"`
	NotNull           basetypes.BoolValue   `tfsdk:"not_null"`
	Format            string                `tfsdk:"format"`
	Granularity       string                `tfsdk:"granularity"`
	TransformFunction basetypes.StringValue `tfsdk:"transform_function"`
}

type tableSchemaResourceModel struct {
	SchemaName                    types.String         `tfsdk:"schema_name"`
	EnableColumnBasedNullHandling basetypes.BoolValue  `tfsdk:"enable_column_based_null_handling"`
	DimensionFieldSpecs           []dimensionFieldSpec `tfsdk:"dimension_field_specs"`
	MetricFieldSpecs              []metricFieldSpec    `tfsdk:"metric_field_specs"`
	DateTimeFieldSpecs            []dateTimeFieldSpec  `tfsdk:"date_time_field_specs"`
	PrimaryKeyColumns             []string             `tfsdk:"primary_key_columns"`
}

func (t *tableSchemaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*pinot.PinotAPIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *pinot.PinotAPIClient, got something else. Please report this issue to the provider developers.",
		)

		return
	}

	t.client = client
}

func (t *tableSchemaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (t *tableSchemaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"schema_name": schema.StringAttribute{
				Description: "The name of the schema.",
				Required:    true,
			},
			"enable_column_based_null_handling": schema.BoolAttribute{
				Description: "Whether to enable column based null handling.",
				Optional:    true,
			},
			"dimension_field_specs": schema.ListNestedAttribute{
				Description: "The dimension field specs.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the dimension.",
							Required:    true,
						},
						"data_type": schema.StringAttribute{
							Description: "The data type of the dimension.",
							Required:    true,
						},
						"not_null": schema.BoolAttribute{
							Description: "Whether the dimension is not null.",
							Optional:    true,
						},
						"single_value_field": schema.BoolAttribute{
							Description: "Whether the dimension is a single value field.",
							Optional:    true,
						},
						"transform_function": schema.StringAttribute{
							Description: "Transform function for specific field.",
							Optional:    true,
						},
					},
				},
			},
			"metric_field_specs": schema.ListNestedAttribute{
				Description: "The dimension field specs.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the dimension.",
							Required:    true,
						},
						"data_type": schema.StringAttribute{
							Description: "The data type of the dimension.",
							Required:    true,
						},
						"not_null": schema.BoolAttribute{
							Description: "Whether the dimension is not null.",
							Optional:    true,
						},
						"transform_function": schema.StringAttribute{
							Description: "Transform function for specific field.",
							Optional:    true,
						},
					},
				},
			},
			"date_time_field_specs": schema.ListNestedAttribute{
				Description: "The dimension field specs.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the dimension.",
							Required:    true,
						},
						"data_type": schema.StringAttribute{
							Description: "The data type of the dimension.",
							Required:    true,
						},
						"not_null": schema.BoolAttribute{
							Description: "Whether the dimension is not null.",
							Optional:    true,
						},
						"format": schema.StringAttribute{
							Description: "The format of the date time.",
							Optional:    true,
						},
						"granularity": schema.StringAttribute{
							Description: "The granularity of the date time.",
							Optional:    true,
						},
						"transform_function": schema.StringAttribute{
							Description: "Transform function for specific field.",
							Optional:    true,
						},
					},
				},
			},
			"primary_key_columns": schema.ListAttribute{
				Description: "The primary key columns.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (t *tableSchemaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("schema_name"), req, resp)
}

func (t *tableSchemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan tableSchemaResourceModel

	diagnostics := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	pinotSchema := model.Schema{
		SchemaName:          plan.SchemaName.ValueString(),
		DimensionFieldSpecs: toDimensionFieldSpecs(plan.DimensionFieldSpecs),
		MetricFieldSpecs:    toMetricFieldSpecs(plan.MetricFieldSpecs),
		DateTimeFieldSpecs:  toDateTimeFieldSpecs(plan.DateTimeFieldSpecs),
		PrimaryKeyColumns:   plan.PrimaryKeyColumns,
	}

	_, err := t.client.CreateSchema(pinotSchema)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create schema", err.Error())
		return
	}

	diagnostics = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (t *tableSchemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state tableSchemaResourceModel
	diagnostics := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	tableSchema, err := t.client.GetSchema(state.SchemaName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get schema", err.Error())
		return
	}

	setState(&state, tableSchema)

	diagnostics = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (t *tableSchemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan tableSchemaResourceModel
	diagnostics := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	pinotSchema := model.Schema{
		SchemaName:          plan.SchemaName.ValueString(),
		DimensionFieldSpecs: toDimensionFieldSpecs(plan.DimensionFieldSpecs),
		MetricFieldSpecs:    toMetricFieldSpecs(plan.MetricFieldSpecs),
		DateTimeFieldSpecs:  toDateTimeFieldSpecs(plan.DateTimeFieldSpecs),
		PrimaryKeyColumns:   plan.PrimaryKeyColumns,
	}

	_, err := t.client.UpdateSchema(pinotSchema)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update schema", err.Error())
		return
	}

	// Update succeeded, reload the table segments
	// First check if the table exists, if it doesn't, don't reload

	tableResp, err := t.client.GetTable(plan.SchemaName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Update Failed: Unable to get table", err.Error())
		return
	}

	// iterate over the table segments and reload them
	// tableResp.REALTIME or tableResp.OFFLINE
	// if tableResp.OFFLINE doesn't exist and tableResp.REALTIME doesn't exist, return

	// TODO: Once go-pinot-api is updated to have IsEmpty() method, use that instead

	tflog.Debug(ctx, "reloading matching table segments")
	if tableResp.OFFLINE.TableName == "" && tableResp.REALTIME.TableName == "" {
		// No tables matching this schema, do not error out
		tflog.Info(ctx, "no tables matching this schema, skipping reload")
	} else {
		_, err := t.client.ReloadTableSegments(plan.SchemaName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Update Failed: Unable to reload table segments", err.Error())
			return
		}
		tflog.Info(ctx, "successfully reloaded matching table segments")
	}

	diagnostics = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (t *tableSchemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state tableSchemaResourceModel
	diagnostics := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := t.client.DeleteSchema(state.SchemaName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete schema", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func toDimensionFieldSpecs(fieldSpecs []dimensionFieldSpec) []model.FieldSpec {

	var pinotFieldSpecs []model.FieldSpec
	for _, fs := range fieldSpecs {
		pinotFieldSpecs = append(pinotFieldSpecs, model.FieldSpec{
			Name:              fs.Name,
			DataType:          fs.DataType,
			NotNull:           fs.NotNull.ValueBoolPointer(),
			SingleValueField:  fs.SingleValueField.ValueBoolPointer(),
			TransformFunction: fs.TransformFunction.ValueString(),
		})
	}
	return pinotFieldSpecs
}

func toMetricFieldSpecs(fieldSpecs []metricFieldSpec) []model.FieldSpec {

	var pinotFieldSpecs []model.FieldSpec
	for _, fs := range fieldSpecs {
		pinotFieldSpecs = append(pinotFieldSpecs, model.FieldSpec{
			Name:              fs.Name,
			DataType:          fs.DataType,
			NotNull:           fs.NotNull.ValueBoolPointer(),
			TransformFunction: fs.TransformFunction.ValueString(),
		})
	}
	return pinotFieldSpecs
}

func toDateTimeFieldSpecs(fieldSpecs []dateTimeFieldSpec) []model.FieldSpec {
	var pinotFieldSpecs []model.FieldSpec
	for _, fs := range fieldSpecs {
		pinotFieldSpecs = append(pinotFieldSpecs, model.FieldSpec{
			Name:              fs.Name,
			DataType:          fs.DataType,
			Format:            fs.Format,
			Granularity:       fs.Granularity,
			NotNull:           fs.NotNull.ValueBoolPointer(),
			TransformFunction: fs.TransformFunction.ValueString(),
		})
	}
	return pinotFieldSpecs
}

func setState(state *tableSchemaResourceModel, schema *model.Schema) {

	dimensionFieldSpecs := make([]dimensionFieldSpec, len(schema.DimensionFieldSpecs))
	for i, fs := range schema.DimensionFieldSpecs {
		dimensionFieldSpec := dimensionFieldSpec{
			Name:             fs.Name,
			DataType:         fs.DataType,
			NotNull:          basetypes.NewBoolPointerValue(fs.NotNull),
			SingleValueField: basetypes.NewBoolPointerValue(fs.SingleValueField),
		}
		if fs.TransformFunction != "" {
			dimensionFieldSpec.TransformFunction = basetypes.NewStringValue(fs.TransformFunction)
		}
		dimensionFieldSpecs[i] = dimensionFieldSpec
	}

	metricFieldSpecs := make([]metricFieldSpec, len(schema.MetricFieldSpecs))
	for i, fs := range schema.MetricFieldSpecs {
		metricFieldSpec := metricFieldSpec{
			Name:     fs.Name,
			DataType: fs.DataType,
			NotNull:  basetypes.NewBoolPointerValue(fs.NotNull),
		}
		if fs.TransformFunction != "" {
			metricFieldSpec.TransformFunction = basetypes.NewStringValue(fs.TransformFunction)
		}
		metricFieldSpecs[i] = metricFieldSpec
	}

	dateTimeFieldSpecs := make([]dateTimeFieldSpec, len(schema.DateTimeFieldSpecs))
	for i, fs := range schema.DateTimeFieldSpecs {
		dateTimeFieldSpec := dateTimeFieldSpec{
			Name:        fs.Name,
			DataType:    fs.DataType,
			Format:      fs.Format,
			Granularity: fs.Granularity,
			NotNull:     basetypes.NewBoolPointerValue(fs.NotNull),
		}
		if fs.TransformFunction != "" {
			dateTimeFieldSpec.TransformFunction = basetypes.NewStringValue(fs.TransformFunction)
		}
		dateTimeFieldSpecs[i] = dateTimeFieldSpec
	}

	state.SchemaName = types.StringValue(schema.SchemaName)
	state.EnableColumnBasedNullHandling = types.BoolValue(false)
	state.DimensionFieldSpecs = dimensionFieldSpecs
	state.MetricFieldSpecs = metricFieldSpecs
	state.DateTimeFieldSpecs = dateTimeFieldSpecs
	state.PrimaryKeyColumns = schema.PrimaryKeyColumns

}
