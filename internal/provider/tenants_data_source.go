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
	_ datasource.DataSource              = &tenantsDataSource{}
	_ datasource.DataSourceWithConfigure = &tenantsDataSource{}
)

// NewTenantsDataSource is a helper function to simplify the provider implementation.
func NewTenantsDataSource() datasource.DataSource {
	return &tenantsDataSource{}
}

// tenantsDataSource is the data source implementation.
type tenantsDataSource struct {
	client *goPinotAPI.PinotAPIClient
}

type tenantsDataSourceModel struct {
	ServerTenants []tenantsModel `tfsdk:"server_tenants"`
	BrokerTenants []tenantsModel `tfsdk:"broker_tenants"`
}

type tenantsModel struct {
	TenantName string `tfsdk:"tenant_name"`
}

// Configure adds the provider configured client to the data source.
func (d *tenantsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *tenantsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tenants"
}

// Schema defines the schema for the data source.
func (d *tenantsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_tenants": schema.ListNestedAttribute{
				Description: "Server tenants",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tenant_name": schema.StringAttribute{
							Description: "The name of the tenant",
							Computed:    true,
						},
					},
				},
			},
			"broker_tenants": schema.ListNestedAttribute{
				Description: "Broker tenants",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tenant_name": schema.StringAttribute{
							Description: "The name of the tenant",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *tenantsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tenantsDataSourceModel

	diags := req.Config.Get(ctx, &state)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	tenants, err := d.client.GetTenants()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get tenants",
			fmt.Sprintf("Failed to get tenants: %s", err),
		)

		return
	}

	for _, brokerTenant := range tenants.BrokerTenants {
		state.BrokerTenants = append(state.BrokerTenants, tenantsModel{
			TenantName: brokerTenant,
		})
	}

	for _, serverTenant := range tenants.ServerTenants {
		state.ServerTenants = append(state.ServerTenants, tenantsModel{
			TenantName: serverTenant,
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}
