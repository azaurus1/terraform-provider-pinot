package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &pinotProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &pinotProvider{
			version: version,
		}
	}
}

// hashicupsProvider is the provider implementation.
type pinotProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type pinotProviderModel struct {
	ControllerURL string `tfsdk:"controller_url"`
}

// Metadata returns the provider type name.
func (p *pinotProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pinot"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *pinotProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"controller_url": schema.StringAttribute{
				Description: "The URL of the Pinot controller.",
				Required:    true,
			},
		},
	}
}

// Configure prepares a Pinot API client for data sources and resources.
func (p *pinotProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config pinotProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ControllerURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("controller_url"),
			"The controller_url must be set.",
			"The provider cannot create the Pinot API client without a controller URL.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	controllerURL := os.Getenv("PINOT_CONTROLLER_URL")

	if !(config.ControllerURL == "") {
		controllerURL = config.ControllerURL
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if controllerURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Controller API URL",
			"The provider cannot create the Controller API client as there is a missing or empty value for the Controller API URL. "+
				"Set the URL value in the configuration or use the PINOT_CONTROLLER_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	pinot := goPinotAPI.NewPinotAPIClient(controllerURL)

	resp.DataSourceData = pinot
	resp.ResourceData = pinot
}

// DataSources defines the data sources implemented in the provider.
func (p *pinotProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUsersDataSource,
		NewTablesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *pinotProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
	}
}
