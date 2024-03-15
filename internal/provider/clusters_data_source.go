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
	_ datasource.DataSource              = &clustersDataSource{}
	_ datasource.DataSourceWithConfigure = &clustersDataSource{}
)

// NewClustersDataSource is a helper function to simplify the provider implementation.
func NewClustersDataSource() datasource.DataSource {
	return &clustersDataSource{}
}

// clustersDataSource is the data source implementation.
type clustersDataSource struct {
	client *goPinotAPI.PinotAPIClient
}

type clustersDataSourceModel struct {
	ClusterName   string             `tfsdk:"cluster_name"`
	ClusterConfig clusterConfigModel `tfsdk:"cluster_config"`
}

type clusterConfigModel struct {
	AllowParticipantAutoJoin            string `tfsdk:"allow_participant_auto_join"`
	EnableCaseInsensitive               string `tfsdk:"enable_case_insensitive"`
	DefaultHyperlogLogLog2m             string `tfsdk:"default_hyperlog_log_log2m"`
	PinotBrokerEnableQueryLimitOverride string `tfsdk:"pinot_broker_enable_query_limit_override"`
}

// Configure adds the provider configured client to the data source.
func (d *clustersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *clustersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clusters"
}

// Schema defines the schema for the data source.
func (d *clustersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_name": schema.StringAttribute{
				Description: "The name of the Pinot cluster.",
				Computed:    true,
			},
			"cluster_config": schema.SingleNestedAttribute{
				Description: "The configuration of the Pinot cluster.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"allow_participant_auto_join": schema.StringAttribute{
						Computed: true,
					},
					"enable_case_insensitive": schema.StringAttribute{
						Computed: true,
					},
					"default_hyperlog_log_log2m": schema.StringAttribute{
						Computed: true,
					},
					"pinot_broker_enable_query_limit_override": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *clustersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state clustersDataSourceModel

	// Get Cluster Name
	clusterNameResp, err := d.client.GetClusterInfo()
	if err != nil {
		resp.Diagnostics.AddError("Failed to get cluster name", err.Error())
		return
	}

	state.ClusterName = clusterNameResp.ClusterName

	// Get Cluster Config
	clusterConfigResp, err := d.client.GetClusterConfigs()
	if err != nil {
		resp.Diagnostics.AddError("Failed to get cluster config", err.Error())
		return
	}

	state.ClusterConfig.AllowParticipantAutoJoin = clusterConfigResp.AllowParticipantAutoJoin
	state.ClusterConfig.EnableCaseInsensitive = clusterConfigResp.EnableCaseInsensitive
	state.ClusterConfig.DefaultHyperlogLogLog2m = clusterConfigResp.DefaultHyperlogLogLog2m
	state.ClusterConfig.PinotBrokerEnableQueryLimitOverride = clusterConfigResp.PinotBrokerEnableQueryLimitOverride

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
