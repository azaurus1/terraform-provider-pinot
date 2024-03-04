package provider

import (
	"context"
	"encoding/json"
	pinot "github.com/azaurus1/go-pinot-api"
	"github.com/azaurus1/go-pinot-api/model"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/appengine/log"
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

type tableSchemaResourceModel struct {
	Name                          types.String         `json:"schemaName" tfsdk:"schema_name"`
	EnableColumnBasedNullHandling types.Bool           `json:"enableColumnBasedNullHandling" tfsdk:"enable_column_based_null_handling"`
	DimensionFieldSpec            []dimensionFieldSpec `json:"dimensionFieldSpecs" tfsdk:"dimension_field_specs"`
	MetricFieldSpec               []metricFieldSpec    `json:"metricFieldSpecs" tfsdk:"metric_field_specs"`
	DateTimeFieldSpec             []dateTimeFieldSpec  `json:"dateTimeFieldSpecs" tfsdk:"date_time_field_specs"`
}

type dimensionFieldSpec struct {
	Name     types.String `json:"name" tfsdk:"name"`
	DataType types.String `json:"dataType" tfsdk:"data_type"`
	NotNull  types.Bool   `json:"notNull" tfsdk:"not_null"`
}

type metricFieldSpec struct {
	Name     types.String `json:"name" tfsdk:"name"`
	DataType types.String `json:"dataType" tfsdk:"data_type"`
	NotNull  types.Bool   `json:"notNull" tfsdk:"not_null"`
}

type dateTimeFieldSpec struct {
	Name        types.String `json:"name" tfsdk:"name"`
	DataType    types.String `json:"dataType" tfsdk:"data_type"`
	NotNull     types.Bool   `json:"notNull" tfsdk:"not_null"`
	Format      types.String `json:"format" tfsdk:"format"`
	Granularity types.String `json:"granularity" tfsdk:"granularity"`
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
				Description: "Enable column based null handling.",
				Optional:    true,
			},
			"dimension_field_specs": schema.ListNestedAttribute{
				Description: "The dimension field specs.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the dimension field.",
							Required:    true,
							Computed:    true,
						},
						"data_type": schema.StringAttribute{
							Description: "The data type of the dimension field.",
							Required:    true,
							Computed:    true,
						},
						"not_null": schema.BoolAttribute{
							Description: "The not null of the dimension field.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"metric_field_specs": schema.ListNestedAttribute{
				Description: "The metric field specs.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the metric field.",
							Required:    true,
							Computed:    true,
						},
						"data_type": schema.StringAttribute{
							Description: "The data type of the metric field.",
							Required:    true,
							Computed:    true,
						},
						"not_null": schema.BoolAttribute{
							Description: "The not null of the metric field.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"date_time_field_specs": schema.ListNestedAttribute{
				Description: "The date time field specs.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the date time field.",
							Required:    true,
							Computed:    true,
						},
						"data_type": schema.StringAttribute{
							Description: "The data type of the date time field.",
							Required:    true,
							Computed:    true,
						},
						"not_null": schema.BoolAttribute{
							Description: "The not null of the date time field.",
							Optional:    true,
							Computed:    true,
						},
						"format": schema.StringAttribute{
							Description: "The format of the date time field.",
							Optional:    true,
							Computed:    true,
						},
						"granularity": schema.StringAttribute{
							Description: "The granularity of the date time field.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
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

	pinotModelSchema, err := toPinotSchemaModel(*t)
	if err != nil {
		resp.Diagnostics.AddError("Failed to convert to pinot schema model", err.Error())
		return
	}

	_, err = t.client.CreateSchema(pinotModelSchema)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create schema", err.Error())
		return
	}

	diagnostics = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		log.Errorf(ctx, "Failed to set state: %v", diagnostics)
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

	tableSchema, err := t.client.GetSchema(state.Name.String())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get schema", err.Error())
		return
	}

	resourceModel, err := fromPinotSchemaModel(*tableSchema)
	if err != nil {
		resp.Diagnostics.AddError("Failed to convert to resource model", err.Error())
		return

	}

	state.Name = resourceModel.Name
	state.EnableColumnBasedNullHandling = resourceModel.EnableColumnBasedNullHandling
	state.DimensionFieldSpec = resourceModel.DimensionFieldSpec
	state.MetricFieldSpec = resourceModel.MetricFieldSpec
	state.DateTimeFieldSpec = resourceModel.DateTimeFieldSpec

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

	schemaToUpdate, err := toPinotSchemaModel(*t)
	if err != nil {
		resp.Diagnostics.AddError("Failed to convert to pinot schema model", err.Error())
		return
	}

	_, err = t.client.UpdateSchema(schemaToUpdate)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update schema", err.Error())
		return
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

	schemaToDelete, err := toPinotSchemaModel(*t)
	if err != nil {
		resp.Diagnostics.AddError("Failed to convert to pinot schema model", err.Error())
		return
	}

	_, err = t.client.DeleteSchema(schemaToDelete.SchemaName)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete schema", err.Error())
		return
	}

	diagnostics = resp.State.Set(ctx, nil)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func toPinotSchemaModel(tfModel tableSchemaResource) (model.Schema, error) {

	var pinotSchemaModel model.Schema

	outBytes, err := json.Marshal(tfModel)
	if err != nil {
		return pinotSchemaModel, err
	}

	err = json.Unmarshal(outBytes, &pinotSchemaModel)
	if err != nil {
		return pinotSchemaModel, err
	}

	return pinotSchemaModel, nil
}

func fromPinotSchemaModel(pinotModel model.Schema) (tableSchemaResourceModel, error) {

	var tfModel tableSchemaResourceModel

	outBytes, err := json.Marshal(pinotModel)
	if err != nil {
		return tfModel, err
	}

	err = json.Unmarshal(outBytes, &tfModel)
	if err != nil {
		return tfModel, err
	}

	return tfModel, nil
}
