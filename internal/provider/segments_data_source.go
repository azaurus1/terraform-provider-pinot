package provider

import (
	"context"
	"fmt"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &segmentsDataSource{}
	_ datasource.DataSourceWithConfigure = &segmentsDataSource{}
)

// NewSegmentsDataSource is a helper function to simplify the provider implementation.
func NewSegmentsDataSource() datasource.DataSource {
	return &segmentsDataSource{}
}

// segmentsDataSource is the data source implementation.
type segmentsDataSource struct {
	client *goPinotAPI.PinotAPIClient
}

type segmentsDataSourceModel struct {
	TableName        types.String    `tfsdk:"table_name"`
	OfflineSegments  []segmentsModel `tfsdk:"offline_segments"`
	RealtimeSegments []segmentsModel `tfsdk:"realtime_segments"`
}

type segmentsModel struct {
	// TableName   string `tfsdk:"table_name"`
	SegmentName string `tfsdk:"segment_name"`
}

// Configure adds the provider configured client to the data source.
func (d *segmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *segmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segments"
}

// Schema returns the data source schema.
func (d *segmentsDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"table_name": schema.StringAttribute{
				Description: "The name of the table to get segments for",
				Required:    true,
			},
			"offline_segments": schema.ListNestedAttribute{
				Description: "The list of offline segments.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"segment_name": schema.StringAttribute{
							Description: "The name of the segment.",
							Computed:    true,
						},
					},
				},
			},
			"realtime_segments": schema.ListNestedAttribute{
				Description: "The list of realtime segments.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"segment_name": schema.StringAttribute{
							Description: "The name of the segment.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read the data source.
func (d *segmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state segmentsDataSourceModel

	diags := req.Config.Get(ctx, &state)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	segments, err := d.client.GetSegments(state.TableName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get segments", err.Error())
		return
	}

	for _, segment := range segments {

		// create slices for offline and realtime segments per segment
		offlineSegments := make([]segmentsModel, 0)
		realtimeSegments := make([]segmentsModel, 0)

		if len(segment.Offline) > 0 {
			for _, segment := range segment.Offline {
				offlineSegments = append(offlineSegments, segmentsModel{
					SegmentName: segment,
				})
			}
		}

		if len(segment.Realtime) > 0 {
			for _, segment := range segment.Realtime {
				realtimeSegments = append(realtimeSegments, segmentsModel{
					SegmentName: segment,
				})
			}
		}

		state.OfflineSegments = offlineSegments
		state.RealtimeSegments = realtimeSegments
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}
