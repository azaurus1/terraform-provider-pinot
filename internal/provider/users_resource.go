package provider

import (
	"context"
	"encoding/json"
	"log"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	goPinotModel "github.com/azaurus1/go-pinot-api/model"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	_ resource.Resource              = &userResource{}
	_ resource.ResourceWithConfigure = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client *goPinotAPI.PinotAPIClient
}

type userResourceModel struct {
	Username  string `tfsdk:"username"`
	Password  string `tfsdk:"password"`
	Component string `tfsdk:"component"`
	Role      string `tfsdk:"role"`
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
		},
	}
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
		Username:  plan.Username,
		Password:  plan.Password,
		Component: plan.Component,
		Role:      plan.Role,
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

}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// get user from state
	// compare for password change

	var state userResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var passwordChanged bool
	if state.Password != plan.Password {
		passwordChanged = true
	}

	// Generate API request from plan
	// Convert plan into []byte

	user := goPinotModel.User{
		Username:  plan.Username,
		Password:  plan.Password,
		Component: plan.Component,
		Role:      plan.Role,
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Panic(err)
	}

	_, err = r.client.UpdateUser(plan.Username, plan.Component, passwordChanged, userBytes)
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

func (r *userResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
