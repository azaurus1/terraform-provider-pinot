package provider

import (
	"context"
	"fmt"
	"log"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	"github.com/azaurus1/go-pinot-api/model"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &schemasDataSource{}
	_ datasource.DataSourceWithConfigure = &schemasDataSource{}
)

func NewSchemasDataSource() datasource.DataSource {
	return &schemasDataSource{}
}

func DimensionFieldSpecsFromFieldSpecs(fieldSpecs []model.FieldSpec) []dimensionFieldSpec {
	var dimensionFieldSpecs []dimensionFieldSpec

	for _, fieldSpec := range fieldSpecs {
		dimensionFieldSpecs = append(dimensionFieldSpecs, dimensionFieldSpec{
			Name:             fieldSpec.Name,
			DataType:         fieldSpec.DataType,
			NotNull:          basetypes.NewBoolPointerValue(fieldSpec.NotNull),
			SingleValueField: basetypes.NewBoolPointerValue(fieldSpec.SingleValueField),
		})
	}

	return dimensionFieldSpecs
}

func MetricFieldSpecsFromFieldSpecs(fieldSpecs []model.FieldSpec) []metricFieldSpec {
	var metricFieldSpecs []metricFieldSpec

	for _, fieldSpec := range fieldSpecs {
		metricFieldSpecs = append(metricFieldSpecs, metricFieldSpec{
			Name:     fieldSpec.Name,
			DataType: fieldSpec.DataType,
			NotNull:  basetypes.NewBoolPointerValue(fieldSpec.NotNull),
		})
	}

	return metricFieldSpecs
}

func DateTimeFieldSpecsFromFieldSpecs(fieldSpecs []model.FieldSpec) []dateTimeFieldSpec {
	var dateTimeFieldSpecs []dateTimeFieldSpec

	for _, fieldSpec := range fieldSpecs {
		dateTimeFieldSpecs = append(dateTimeFieldSpecs, dateTimeFieldSpec{
			Name:        fieldSpec.Name,
			DataType:    fieldSpec.DataType,
			NotNull:     basetypes.NewBoolPointerValue(fieldSpec.NotNull),
			Format:      fieldSpec.Format,
			Granularity: fieldSpec.Granularity,
		})
	}

	return dateTimeFieldSpecs
}

type schemaModel struct {
	SchemaName                    types.String         `tfsdk:"schema_name"`
	EnableColumnBasedNullHandling basetypes.BoolValue  `tfsdk:"enable_column_based_null_handling"`
	DimensionFieldSpecs           []dimensionFieldSpec `tfsdk:"dimension_field_specs"`
	MetricFieldSpecs              []metricFieldSpec    `tfsdk:"metric_field_specs"`
	DateTimeFieldSpecs            []dateTimeFieldSpec  `tfsdk:"date_time_field_specs"`
	PrimaryKeyColumns             []string             `tfsdk:"primary_key_columns"`
}

type schemasDataSource struct {
	client     *goPinotAPI.PinotAPIClient
	SchemaName types.String `tfsdk:"schema_name"`
}

type schemasDataSourceModel struct {
	Schemas []schemaModel `tfsdk:"schemas"`
}

// Configure adds the provider configured client to the data source.
func (d *schemasDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *schemasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schemas"
}

func (d *schemasDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"schemas": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"schema_name": schema.StringAttribute{
							Description: "The name of the schema.",
							Required:    true,
						},
						"enable_column_based_null_handling": schema.BoolAttribute{
							Description: "Whether to enable column based null handling.",
							Optional:    true,
						},
						"dimension_field_specs": schema.ListNestedAttribute{
							Description: "The dimension field specs.",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "The name of the dimension.",
										Required:    true,
									},
									"data_type": schema.StringAttribute{
										Description: "The data type of the dimension.",
										Required:    true,
									},
									"not_null": schema.BoolAttribute{
										Description: "Whether the dimension is not null.",
										Optional:    true,
									},
									"single_value_field": schema.BoolAttribute{
										Description: "Whether the dimension is a single value field.",
										Optional:    true,
									},
								},
							},
						},
						"metric_field_specs": schema.ListNestedAttribute{
							Description: "The dimension field specs.",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "The name of the dimension.",
										Required:    true,
									},
									"data_type": schema.StringAttribute{
										Description: "The data type of the dimension.",
										Required:    true,
									},
									"not_null": schema.BoolAttribute{
										Description: "Whether the dimension is not null.",
										Optional:    true,
									},
								},
							},
						},
						"date_time_field_specs": schema.ListNestedAttribute{
							Description: "The dimension field specs.",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "The name of the dimension.",
										Required:    true,
									},
									"data_type": schema.StringAttribute{
										Description: "The data type of the dimension.",
										Required:    true,
									},
									"not_null": schema.BoolAttribute{
										Description: "Whether the dimension is not null.",
										Optional:    true,
									},
									"format": schema.StringAttribute{
										Description: "The format of the date time.",
										Optional:    true,
									},
									"granularity": schema.StringAttribute{
										Description: "The granularity of the date time.",
										Optional:    true,
									},
								},
							},
						},
						"primary_key_columns": schema.ListAttribute{
							Description: "The primary key columns.",
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *schemasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state schemasDataSourceModel

	currentSchemas, err := d.client.GetSchemas()
	if err != nil {
		log.Panic(err)
	}

	currentSchemas.ForEachSchema(func(schemaName string) {

		schemaResp, err := d.client.GetSchema(schemaName)
		if err != nil {
			log.Panic(err)
		}

		// fmt.Println("Reading Schema:")
		state.Schemas = append(state.Schemas, schemaModel{
			SchemaName:          basetypes.NewStringValue(schemaResp.SchemaName),
			DimensionFieldSpecs: DimensionFieldSpecsFromFieldSpecs(schemaResp.DimensionFieldSpecs),
			MetricFieldSpecs:    MetricFieldSpecsFromFieldSpecs(schemaResp.MetricFieldSpecs),
			DateTimeFieldSpecs:  DateTimeFieldSpecsFromFieldSpecs(schemaResp.DateTimeFieldSpecs),
			PrimaryKeyColumns:   schemaResp.PrimaryKeyColumns,
		})

	})

	diagnostics := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

}
