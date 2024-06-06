package provider

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	goPinotModel "github.com/azaurus1/go-pinot-api/model"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client *goPinotAPI.PinotAPIClient
}

type userResourceModel struct {
	Username    basetypes.StringValue `tfsdk:"username"`
	Password    basetypes.StringValue `tfsdk:"password"`
	Component   basetypes.StringValue `tfsdk:"component"`
	Role        basetypes.StringValue `tfsdk:"role"`
	Permissions *[]string             `tfsdk:"permissions"`
	Tables      *[]string             `tfsdk:"tables"`
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Description: "The username of the user.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password of the user.",
				Required:    true,
				Sensitive:   true,
			},
			"component": schema.StringAttribute{
				Description: "The component of the user.",
				Required:    true,
			},
			"role": schema.StringAttribute{
				Description: "The role of the user.",
				Required:    true,
			},
			"permissions": schema.ListAttribute{
				Description: "A list of permissions for the user.",
				Optional:    true,
				ElementType: basetypes.StringType{},
			},
			"tables": schema.ListAttribute{
				Description: "A list of tables which the permissions are applied to.",
				Optional:    true,
				ElementType: basetypes.StringType{},
			},
		},
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("Invalid ID format", "Expected ID to be in the format <username>,<component>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("username"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("component"), idParts[1])...)

}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request from plan
	// Convert plan into []byte

	user := goPinotModel.User{
		Username:  plan.Username.ValueString(),
		Password:  plan.Password.ValueString(),
		Component: plan.Component.ValueString(),
		Role:      plan.Role.ValueString(),
	}

	if plan.Permissions != nil {
		user.Permissions = plan.Permissions
	}

	if plan.Tables != nil {
		user.Tables = plan.Tables
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Panic(err)
	}

	_, err = r.client.CreateUser(userBytes)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create user", err.Error())
		return
	}

	// set state to populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUser(state.Username.ValueString(), state.Component.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get user", err.Error())
		return
	}

	state.Username = basetypes.NewStringValue(user.Username)
	state.Password = basetypes.NewStringValue(user.Password)
	state.Component = basetypes.NewStringValue(user.Component)
	state.Role = basetypes.NewStringValue(user.Role)

	if user.Permissions != nil {
		state.Permissions = user.Permissions
	}

	if user.Tables != nil {
		state.Tables = user.Tables
	}

	// set state to populated data
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// compare for password change

	var state userResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var passwordChanged = false
	var user goPinotModel.User

	if state.Password != plan.Password {
		passwordChanged = true
		user = goPinotModel.User{
			Username:  plan.Username.ValueString(),
			Password:  plan.Password.ValueString(),
			Component: plan.Component.ValueString(),
			Role:      plan.Role.ValueString(),
		}
		if plan.Permissions != nil {
			user.Permissions = plan.Permissions
		}

		if plan.Tables != nil {
			user.Tables = plan.Tables
		}

	} else {
		user = goPinotModel.User{
			Username:  plan.Username.ValueString(),
			Password:  state.Password.ValueString(),
			Component: plan.Component.ValueString(),
			Role:      plan.Role.ValueString(),
		}
		if plan.Permissions != nil {
			user.Permissions = plan.Permissions
		}

		if plan.Tables != nil {
			user.Tables = plan.Tables
		}
	}

	// Generate API request from plan
	// Convert plan into []byte

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Panic(err)
	}

	_, err = r.client.UpdateUser(plan.Username.ValueString(), plan.Component.ValueString(), passwordChanged, userBytes)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update user", err.Error())
		return
	}

	// set state to populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteUser(state.Username.ValueString(), state.Component.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete user", err.Error())
		return
	}

	// set state to populated data
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
