package provider

import (
	"context"
	"fmt"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &instancesDataSource{}
	_ datasource.DataSourceWithConfigure = &instancesDataSource{}
)

// NewUsersDataSource is a helper function to simplify the provider implementation.
func NewInstancesDataSource() datasource.DataSource {
	return &instancesDataSource{}
}

// usersDataSource is the data source implementation.
type instancesDataSource struct {
	client *goPinotAPI.PinotAPIClient
}

type instancesDataSourceModel struct {
	Instances []instancesModel `tfsdk:"instances"`
}

type systemResourceInfoModel struct {
	NumCores      string `tfsdk:"num_cores"`
	TotalMemoryMB string `tfsdk:"total_memory_mb"`
	MaxHeapSizeMB string `tfsdk:"max_heap_size_mb"`
}

type instancesModel struct {
	InstanceName       string                  `tfsdk:"instance_name"`
	HostName           string                  `tfsdk:"host_name"`
	Enabled            bool                    `tfsdk:"enabled"`
	Port               string                  `tfsdk:"port"`
	Tags               []string                `tfsdk:"tags"`
	Pools              []string                `tfsdk:"pools"`
	GRPCPort           int                     `tfsdk:"grpc_port"`
	AdminPort          int                     `tfsdk:"admin_port"`
	QueryServicePort   int                     `tfsdk:"query_service_port"`
	QueryMailboxPort   int                     `tfsdk:"query_mailbox_port"`
	SystemResourceInfo systemResourceInfoModel `tfsdk:"system_resource_info"`
}

// Configure adds the provider configured client to the data source.
func (d *instancesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *instancesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instances"
}

// Schema defines the schema for the data source.
func (d *instancesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"instances": schema.ListNestedAttribute{
				Description: "The list of instances.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_name": schema.StringAttribute{
							Description: "The name of the instance.",
							Computed:    true,
						},
						"host_name": schema.StringAttribute{
							Description: "The hostname of the instance.",
							Computed:    true,
							// Sensitive:   true,
						},
						"enabled": schema.BoolAttribute{
							Description: "If the instance is enabled.",
							Computed:    true,
						},
						"port": schema.StringAttribute{
							Description: "The port of the instance.",
							Computed:    true,
						},
						"tags": schema.ListAttribute{
							Description: "The list of tags.",
							Computed:    true,
							ElementType: basetypes.StringType{},
						},
						"pools": schema.ListAttribute{
							Description: "The list of pools.",
							Computed:    true,
							ElementType: basetypes.StringType{},
						},
						"grpc_port": schema.NumberAttribute{
							Description: "The GRPC port of the instance.",
							Computed:    true,
						},
						"admin_port": schema.NumberAttribute{
							Description: "The admin port of the instance.",
							Computed:    true,
						},
						"query_service_port": schema.NumberAttribute{
							Description: "The query server port of the instance.",
							Computed:    true,
						},
						"query_mailbox_port": schema.NumberAttribute{
							Description: "The query mailbox port of the instance.",
							Computed:    true,
						},
						"system_resource_info": schema.SingleNestedAttribute{
							Description: "The role of the user.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"num_cores": schema.StringAttribute{
									Description: "The number of cores.",
									Computed:    true,
								},
								"total_memory_mb": schema.StringAttribute{
									Description: "The total memory in MB.",
									Computed:    true,
								},
								"max_heap_size_mb": schema.StringAttribute{
									Description: "The max heap size in MB.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *instancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state instancesDataSourceModel

	instancesResp, err := d.client.GetInstances()
	if err != nil {
		resp.Diagnostics.AddError("Failed to get instances", fmt.Sprintf("Failed to get instances: %s", err))
		return
	}

	for _, instance := range instancesResp.Instances {
		instanceResp, err := d.client.GetInstance(instance)
		if err != nil {
			resp.Diagnostics.AddError("Failed to get instance", fmt.Sprintf("Failed to get instance: %s", err))
			return
		}

		systemInfo := systemResourceInfoModel{
			NumCores:      instanceResp.SystemResourceInfo.NumCores,
			TotalMemoryMB: instanceResp.SystemResourceInfo.TotalMemoryMB,
			MaxHeapSizeMB: instanceResp.SystemResourceInfo.MaxHeapSizeMB,
		}

		state.Instances = append(state.Instances, instancesModel{
			InstanceName:       instanceResp.InstanceName,
			HostName:           instanceResp.Hostname,
			Enabled:            instanceResp.Enabled,
			Port:               instanceResp.Port,
			Tags:               instanceResp.Tags,
			Pools:              instanceResp.Pools,
			GRPCPort:           instanceResp.GRPCPort,
			AdminPort:          instanceResp.AdminPort,
			QueryServicePort:   instanceResp.QueryServicePort,
			QueryMailboxPort:   instanceResp.QueryMailboxPort,
			SystemResourceInfo: systemInfo,
		})

	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
