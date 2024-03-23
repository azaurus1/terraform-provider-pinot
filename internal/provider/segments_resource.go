package provider

import (
	"context"

	pinot "github.com/azaurus1/go-pinot-api"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource              = &segmentResource{}
	_ resource.ResourceWithConfigure = &segmentResource{}
)

type segmentResource struct {
	client *pinot.PinotAPIClient
}

func NewSegmentResource() resource.Resource {
	return &segmentResource{}
}

func (t *segmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*pinot.PinotAPIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *pinot.PinotAPIClient, got something else. Please report this issue to the provider developers.",
		)

		return
	}

	t.client = client
}

func (t *segmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment"
}

func (t *segmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"content_disposition": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{},
					"parameters": schema.MapAttribute{
						ElementType: basetypes.StringType{},
					},
					"file_name":         schema.StringAttribute{},
					"creation_date":     schema.StringAttribute{},
					"modification_date": schema.StringAttribute{},
					"read_date":         schema.StringAttribute{},
					"size":              schema.NumberAttribute{},
				},
			},
			"entity": schema.MapAttribute{
				ElementType: basetypes.StringType{},
			},
			"headers": schema.MapAttribute{
				ElementType: basetypes.ListType{},
			},
			"media_type": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"type":    schema.StringAttribute{},
					"subtype": schema.StringAttribute{},
					"parameters": schema.MapAttribute{
						ElementType: basetypes.StringType{},
					},
					"wildcard_type":    schema.BoolAttribute{},
					"wildcard_subtype": schema.BoolAttribute{},
				},
			},
			"message_body_workers": schema.MapAttribute{
				ElementType: basetypes.StringType{},
			},
			"parent": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"content_disposition": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{},
							"parameters": schema.MapAttribute{
								ElementType: basetypes.StringType{},
							},
							"file_name":         schema.StringAttribute{},
							"creation_date":     schema.StringAttribute{},
							"modification_date": schema.StringAttribute{},
							"read_date":         schema.StringAttribute{},
							"size":              schema.NumberAttribute{},
						},
					},
					"entity": schema.MapAttribute{
						ElementType: basetypes.StringType{},
					},
					"headers": schema.MapAttribute{
						ElementType: basetypes.ListType{},
					},
					"media_type": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"type":    schema.StringAttribute{},
							"subtype": schema.StringAttribute{},
							"parameters": schema.MapAttribute{
								ElementType: basetypes.StringType{},
							},
							"wildcard_type":    schema.BoolAttribute{},
							"wildcard_subtype": schema.BoolAttribute{},
						},
					},
					"message_body_workers": schema.MapAttribute{
						ElementType: basetypes.StringType{},
					},
					"parent": schema.StringAttribute{},
					"providers": schema.MapAttribute{
						ElementType: basetypes.StringType{},
					},
					"body_parts": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"content_disposition": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{},
										"parameters": schema.MapAttribute{
											ElementType: basetypes.StringType{},
										},
										"file_name":         schema.StringAttribute{},
										"creation_date":     schema.StringAttribute{},
										"modification_date": schema.StringAttribute{},
										"read_date":         schema.StringAttribute{},
										"size":              schema.NumberAttribute{},
									},
								},
								"entity": schema.MapAttribute{
									ElementType: basetypes.StringType{},
								},
								"headers": schema.MapAttribute{
									ElementType: basetypes.ListType{},
								},
								"media_type": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"type":    schema.StringAttribute{},
										"subtype": schema.StringAttribute{},
										"parameters": schema.MapAttribute{
											ElementType: basetypes.StringType{},
										},
										"wildcard_type":    schema.BoolAttribute{},
										"wildcard_subtype": schema.BoolAttribute{},
									},
								},
								"message_body_workers": schema.MapAttribute{
									ElementType: basetypes.StringType{},
								},
								"parent": schema.StringAttribute{},
								"providers": schema.MapAttribute{
									ElementType: basetypes.StringType{},
								},
								"parameterized_headers": schema.MapAttribute{
									ElementType: basetypes.StringType{},
								},
							},
						},
					},
					"parameterized_headers": schema.MapAttribute{
						ElementType: basetypes.StringType{},
					},
				},
			},
		},
	}
}

func (t *segmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

func (t *segmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (t *segmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (t *segmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
