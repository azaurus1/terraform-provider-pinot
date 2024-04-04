package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"terraform-provider-pinot/internal/converter"
	"terraform-provider-pinot/internal/models"
	"terraform-provider-pinot/internal/tf_schema"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	"github.com/azaurus1/go-pinot-api/model"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

func (r *tableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"table_name": schema.StringAttribute{
				Description: "The name of the table.",
				Required:    true,
			},
			"table": schema.StringAttribute{
				Description: "The table definition.",
				Required:    true,
				Sensitive:   true,
			},
			"table_type": schema.StringAttribute{
				Description: "The table type.",
				Required:    true,
			},
			"is_dim_table": schema.BoolAttribute{
				Description: "is dimension table",
				Optional:    true,
			},
			"segments_config":    tf_schema.SegmentsConfig(),
			"tenants":            tf_schema.Tenants(),
			"table_index_config": tf_schema.TableIndexConfig(),
			"upsert_config":      tf_schema.UpsertConfig(),
			"ingestion_config":   tf_schema.IngestionConfig(),
			"tier_configs":       tf_schema.TierConfigs(),
			"field_config_list":  tf_schema.FieldConfigList(),
			"routing":            tf_schema.Routing(),
			"metadata":           tf_schema.Metadata(),
		},
	}
}

func (r *tableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan models.TableResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var table model.Table
	err := json.Unmarshal([]byte(plan.Table.ValueString()), &table)
	if err != nil {
		resp.Diagnostics.AddError("Create Failed: Unable to unmarshal table from config", err.Error())
		return
	}

	tableWithPlanOverrides, resultDiags := converter.ToTable(&plan)
	if resultDiags.HasError() {
		resp.Diagnostics.Append(resultDiags...)
		return
	}

	overriddenTableBytes, err := json.Marshal(tableWithPlanOverrides)
	if err != nil {
		resp.Diagnostics.AddError("Create Failed: Unable to marshal table", err.Error())
		return
	}

	_, err = r.client.CreateTable(overriddenTableBytes)
	if err != nil {
		resp.Diagnostics.AddError("Create Failed: Unable to create table", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state models.TableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tableResponse, err := r.client.GetTable(state.TableName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get table", err.Error())
		return
	}

	tflog.Info(ctx, "here\n")

	var table model.Table

	// if table.OFFLINE is not nil, set the state to populated data
	if tableResponse.OFFLINE.TableName != "" {
		table = tableResponse.OFFLINE
	} else {
		table = tableResponse.REALTIME
	}

	tflog.Info(ctx, "setting state\n")

	resultDiags := converter.SetStateFromTable(ctx, &state, &table)
	if resultDiags.HasError() {
		resp.Diagnostics.Append(resultDiags...)
		return
	}

	// set state to populated data
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *tableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan models.TableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Overriding table config: %s", plan.TableName))

	tableWithPlanOverrides, resultDiags := converter.ToTable(&plan)
	if resultDiags.HasError() {
		resp.Diagnostics.Append(resultDiags...)
		return
	}

	overriddenTableBytes, err := json.Marshal(tableWithPlanOverrides)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed: Unable to marshal table", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating table: %s", plan.TableName))

	_, err = r.client.UpdateTable(plan.TableName.ValueString(), overriddenTableBytes)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed: Unable to update table", err.Error())
		return
	}

	// Update succeeded, reload the table segments

	_, err = r.client.ReloadTableSegments(plan.TableName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Update Failed: Unable to reload table segments", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Table segments reloaded: %s", plan.TableName))

	// set state to populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *tableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.TableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting table: %s", state.TableName))

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
