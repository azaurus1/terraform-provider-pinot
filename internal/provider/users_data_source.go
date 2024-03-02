package provider

import (
	"context"
	"fmt"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

// NewUsersDataSource is a helper function to simplify the provider implementation.
func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

// usersDataSource is the data source implementation.
type usersDataSource struct {
	client *goPinotAPI.PinotAPIClient
}

type usersDataSourceModel struct {
	Users []usersModel `tfsdk:"users"`
}

type usersModel struct {
	Username  string `tfsdk:"username"`
	Password  string `tfsdk:"password"`
	Component string `tfsdk:"component"`
	Role      string `tfsdk:"role"`
}

// Configure adds the provider configured client to the data source.
func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*goPinotAPI.PinotAPIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *goPinotAPI.PinotAPIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the data source.
func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				Description: "The list of users.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"username": schema.StringAttribute{
							Description: "The username of the user.",
							Computed:    true,
						},
						"password": schema.StringAttribute{
							Description: "The password of the user.",
							Computed:    true,
							// Sensitive:   true,
						},
						"component": schema.StringAttribute{
							Description: "The component of the user.",
							Computed:    true,
						},
						"role": schema.StringAttribute{
							Description: "The role of the user.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state usersDataSourceModel

	usersResp, err := d.client.GetUsers()
	if err != nil {
		resp.Diagnostics.AddError("Failed to get users", fmt.Sprintf("Failed to get users: %s", err))
		return
	}

	for _, user := range usersResp.Users {
		state.Users = append(state.Users, usersModel{
			Username:  user.Username,
			Password:  user.Password,
			Component: user.Component,
			Role:      user.Role,
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
