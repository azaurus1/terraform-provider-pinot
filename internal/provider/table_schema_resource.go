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
	SchemaName types.String `tfsdk:"schema_name"`
	Schema     types.String `tfsdk:"schema"`
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
			"schema": schema.StringAttribute{
				Description: "The schema definition.",
				Required:    true,
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

	var pinotSchema model.Schema
	err := json.Unmarshal([]byte(plan.Schema.ValueString()), &pinotSchema)
	if err != nil {
		resp.Diagnostics.AddError("Create Failed: Unable to unmarshal schema", err.Error())
		return
	}

	_, err = t.client.CreateSchema(pinotSchema)
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

	tableSchema, err := t.client.GetSchema(state.SchemaName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get schema", err.Error())
		return
	}

	state.SchemaName = types.StringValue(tableSchema.SchemaName)

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

	var pinotSchema model.Schema
	err := json.Unmarshal([]byte(plan.Schema.ValueString()), &pinotSchema)
	if err != nil {
		resp.Diagnostics.AddWarning(plan.Schema.String(), "")
		resp.Diagnostics.AddError("Update Failed: Unable to unmarshal schema", plan.Schema.String())
		return
	}

	_, err = t.client.UpdateSchema(pinotSchema)
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
