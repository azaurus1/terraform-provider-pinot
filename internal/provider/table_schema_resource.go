package provider

import (
	"context"

	pinot "github.com/azaurus1/go-pinot-api"
	"github.com/azaurus1/go-pinot-api/model"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &tableSchemaResource{}
	_ resource.ResourceWithConfigure = &tableSchemaResource{}
)

type tableSchemaResource struct {
	client *pinot.PinotAPIClient
}

func NewTableSchemaResource() resource.Resource {
	return &tableSchemaResource{}

}

type fieldSpec struct {
	Name     string              `tfsdk:"name"`
	DataType string              `tfsdk:"data_type"`
	NotNull  basetypes.BoolValue `tfsdk:"not_null"`
}

type dateTimeFieldSpec struct {
	Name        string              `tfsdk:"name"`
	DataType    string              `tfsdk:"data_type"`
	NotNull     basetypes.BoolValue `tfsdk:"not_null"`
	Format      string              `tfsdk:"format"`
	Granularity string              `tfsdk:"granularity"`
}

type tableSchemaResourceModel struct {
	SchemaName                    types.String        `tfsdk:"schema_name"`
	EnableColumnBasedNullHandling basetypes.BoolValue `tfsdk:"enable_column_based_null_handling"`
	DimensionFieldSpecs           []fieldSpec         `tfsdk:"dimension_field_specs"`
	MetricFieldSpecs              []fieldSpec         `tfsdk:"metric_field_specs"`
	DateTimeFieldSpecs            []dateTimeFieldSpec `tfsdk:"date_time_field_specs"`
	PrimaryKeyColumns             []string            `tfsdk:"primary_key_columns"`
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

func (t *tableSchemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan tableSchemaResourceModel

	diagnostics := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	pinotSchema := model.Schema{
		SchemaName:          plan.SchemaName.ValueString(),
		DimensionFieldSpecs: toPinotModelFieldSpec(plan.DimensionFieldSpecs),
		MetricFieldSpecs:    toPinotModelFieldSpec(plan.MetricFieldSpecs),
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
		DimensionFieldSpecs: toPinotModelFieldSpec(plan.DimensionFieldSpecs),
		MetricFieldSpecs:    toPinotModelFieldSpec(plan.MetricFieldSpecs),
		DateTimeFieldSpecs:  toDateTimeFieldSpecs(plan.DateTimeFieldSpecs),
		PrimaryKeyColumns:   plan.PrimaryKeyColumns,
	}

	_, err := t.client.UpdateSchema(pinotSchema)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update schema", err.Error())
		return
	}

	// Update succeeded, reload the table segments

	// First check if the table exists, if it doesnt, dont reload

	tableResp, err := t.client.GetTable(plan.SchemaName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Update Failed: Unable to get table", err.Error())
		return
	}

	// iterate over the table segments and reload them
	// tableResp.REALTIME or tableResp.OFFLINE

	// if tableResp.OFFLINE doesnt exist and tableResp.REALTIME doesnt exist, return

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

func toPinotModelFieldSpec(fieldSpecs []fieldSpec) []model.FieldSpec {

	var pinotFieldSpecs []model.FieldSpec
	for _, fs := range fieldSpecs {
		pinotFieldSpecs = append(pinotFieldSpecs, model.FieldSpec{
			Name:     fs.Name,
			DataType: fs.DataType,
			NotNull:  fs.NotNull.ValueBool(),
		})
	}
	return pinotFieldSpecs
}

func toDateTimeFieldSpecs(fieldSpecs []dateTimeFieldSpec) []model.FieldSpec {
	var pinotFieldSpecs []model.FieldSpec
	for _, fs := range fieldSpecs {
		pinotFieldSpecs = append(pinotFieldSpecs, model.FieldSpec{
			Name:        fs.Name,
			DataType:    fs.DataType,
			Format:      fs.Format,
			Granularity: fs.Granularity,
			NotNull:     fs.NotNull.ValueBool(),
		})
	}
	return pinotFieldSpecs
}

func setState(state *tableSchemaResourceModel, schema *model.Schema) {

	dimensionFieldSpecs := make([]fieldSpec, len(schema.DimensionFieldSpecs))
	for i, fs := range schema.DimensionFieldSpecs {
		dimensionFieldSpecs[i] = fieldSpec{
			Name:     fs.Name,
			DataType: fs.DataType,
			NotNull:  basetypes.NewBoolValue(fs.NotNull),
		}
	}

	metricFieldSpecs := make([]fieldSpec, len(schema.MetricFieldSpecs))
	for i, fs := range schema.MetricFieldSpecs {
		metricFieldSpecs[i] = fieldSpec{
			Name:     fs.Name,
			DataType: fs.DataType,
			NotNull:  basetypes.NewBoolValue(fs.NotNull),
		}
	}

	dateTimeFieldSpecs := make([]dateTimeFieldSpec, len(schema.DateTimeFieldSpecs))
	for i, fs := range schema.DateTimeFieldSpecs {
		dateTimeFieldSpecs[i] = dateTimeFieldSpec{
			Name:        fs.Name,
			DataType:    fs.DataType,
			Format:      fs.Format,
			Granularity: fs.Granularity,
			NotNull:     basetypes.NewBoolValue(fs.NotNull),
		}
	}

	state.SchemaName = types.StringValue(schema.SchemaName)
	state.EnableColumnBasedNullHandling = types.BoolValue(false)
	state.DimensionFieldSpecs = dimensionFieldSpecs
	state.MetricFieldSpecs = metricFieldSpecs
	state.DateTimeFieldSpecs = dateTimeFieldSpecs
	state.PrimaryKeyColumns = schema.PrimaryKeyColumns

}
