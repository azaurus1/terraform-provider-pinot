package provider

import (
	"context"
	"encoding/json"
	"fmt"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	"github.com/azaurus1/go-pinot-api/model"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &tableResource{}
	_ resource.ResourceWithConfigure = &tableResource{}
)

func NewTableResource() resource.Resource {
	return &tableResource{}
}

type tableResource struct {
	client *goPinotAPI.PinotAPIClient
}

type tableResourceModel struct {
	TableName types.String `tfsdk:"table_name"`
	Table     types.String `tfsdk:"table"`
	TableType types.String `tfsdk:"table_type"`
}

func (r *tableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*goPinotAPI.PinotAPIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *goPinotAPI.PinotAPIClient, got something else. Please report this issue to the provider developers.",
		)

		return
	}

	r.client = client
}

func (r *tableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_table"
}

func (r *tableResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"table_name": schema.StringAttribute{
				Description: "The name of the table.",
				Required:    true,
			},
			"table": schema.StringAttribute{
				Description: "The table definition.",
				Required:    true,
			},
			"table_type": schema.StringAttribute{
				Description: "The table type.",
				Required:    true,
			},
		},
	}
}

func (r *tableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var table model.Table
	err := json.Unmarshal([]byte(plan.Table.ValueString()), &table)
	if err != nil {
		resp.Diagnostics.AddError("Create Failed: Unable to unmarshal table", err.Error())
		return
	}

	_, err = r.client.CreateTable([]byte(plan.Table.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Create Failed: Unable to create table", err.Error())
		return
	}

	// set state to populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *tableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	table, err := r.client.GetTable(state.TableName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get table", err.Error())
		return
	}

	// if table.OFFLINE is not nil, set the state to populated data
	if table.OFFLINE.TableName != "" {
		state.TableName = types.StringValue(table.OFFLINE.TableName)
	}

	// if table.REALTIME is not nil, set the state to populated data
	if table.REALTIME.TableName != "" {
		state.TableName = types.StringValue(table.REALTIME.TableName)
	}

	// set state to populated data
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *tableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var table model.Table
	err := json.Unmarshal([]byte(plan.Table.ValueString()), &table)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed: Unable to unmarshal table", err.Error())
		return
	}

	_, err = r.client.UpdateTable(plan.TableName.String(), []byte(plan.Table.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Update Failed: Unable to update table", err.Error())
		return
	}

	// set state to populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *tableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	log := fmt.Sprintf("Deleting table: %s", state.TableName)

	tflog.Info(ctx, log)

	_, err := r.client.DeleteTable(state.TableName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Failed: Unable to delete table", err.Error())
		return
	}

	// set state to populated data
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
