package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"log"
	"strconv"
	"terraform-provider-pinot/internal/converter"
	"terraform-provider-pinot/internal/models"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	"github.com/azaurus1/go-pinot-api/model"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &tableResource{}
	_ resource.ResourceWithConfigure = &tableResource{}
)

func NewTableResource() resource.Resource {
	return &tableResource{}
}

type tableResource struct {
	client *goPinotAPI.PinotAPIClient
}

func (r *tableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *tableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_table"
}

func (r *tableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"table_name": schema.StringAttribute{
				Description: "The name of the table.",
				Required:    true,
			},
			"table": schema.StringAttribute{
				Description: "The table definition.",
				Sensitive:   true,
				Required:    true,
			},
			"table_type": schema.StringAttribute{
				Description: "The table type.",
				Required:    true,
			},
			"segments_config": schema.SingleNestedAttribute{
				Description: "The segments configuration for the table.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"replication": schema.StringAttribute{
						Description: "The replication count for the segments.",
						Required:    true,
					},
					"replicas_per_partition": schema.StringAttribute{
						Description: "The replicas per partition for the segments.",
						Optional:    true,
					},
					"time_type": schema.StringAttribute{
						Description: "The time type for the segments.",
						Required:    true,
					},
					"time_column_name": schema.StringAttribute{
						Description: "The time column name for the segments.",
						Required:    true,
					},
					"retention_time_unit": schema.StringAttribute{
						Description: "The retention time unit for the segments.",
						Optional:    true,
					},
					"retention_time_value": schema.StringAttribute{
						Description: "The retention time value for the segments.",
						Optional:    true,
					},
					"deleted_segments_retention_period": schema.StringAttribute{
						Description: "The deleted segments retention period for the segments.",
						Optional:    true,
					},
				},
			},
			"tenants": schema.SingleNestedAttribute{
				Description: "The tenants configuration for the table.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"broker": schema.StringAttribute{
						Description: "The broker for the tenants.",
						Optional:    true,
					},
					"server": schema.StringAttribute{
						Description: "The server for the tenants.",
						Optional:    true,
					},
					"tag_override_config": schema.MapAttribute{
						Description: "The tag override config for the tenants.",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			"table_index_config": schema.SingleNestedAttribute{
				Description: "The table index configuration for the table.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"sorted_column": schema.ListAttribute{
						Description: "The sorted column for the table.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"load_mode": schema.StringAttribute{
						Description: "The load mode for the table.",
						Optional:    true,
					},
					"null_handling_enabled": schema.BoolAttribute{
						Description: "The null handling enabled for the table.",
						Optional:    true,
					},
					"create_inverted_index_during_segment_generation": schema.BoolAttribute{
						Description: "The create inverted index during segment generation for the table.",
						Optional:    true,
					},
					"star_tree_index_configs": schema.ListNestedAttribute{
						Description: "The star tree index configurations for the table.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"dimensions_split_order": schema.ListAttribute{
									Description: "The dimensions split order for the star tree index.",
									Optional:    true,
									ElementType: types.StringType,
								},
								"skip_star_node_creation_for_dim_names": schema.ListAttribute{
									Description: "The skip star node creation for dim names for the star tree index.",
									Optional:    true,
									ElementType: types.StringType,
								},
								"max_leaf_records": schema.Int64Attribute{
									Description: "The max leaf records for the star tree index.",
									Required:    true,
								},
								"function_column_pairs": schema.ListAttribute{
									Description: "The function column pairs for the star tree index.",
									Optional:    true,
									ElementType: types.StringType,
								},
								"aggregation_configs": schema.ListNestedAttribute{
									Description: "The aggregation configurations for the star tree index.",
									Optional:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"column_name": schema.StringAttribute{
												Description: "The column name for the star tree index.",
												Required:    true,
											},
											"aggregate_function": schema.StringAttribute{
												Description: "The aggregate function for the star tree index.",
												Required:    true,
											},
											"compression_codec": schema.StringAttribute{
												Description: "The compression codec for the star tree index.",
												Required:    true,
											},
										},
									},
								},
							},
						},
					},
					"enable_dynamic_star_tree": schema.BoolAttribute{
						Description: "The enable dynamic star tree for the table.",
						Optional:    true,
					},
					"enable_default_star_tree": schema.BoolAttribute{
						Description: "The enable default star tree for the table.",
						Optional:    true,
					},
					"optimize_dictionary": schema.BoolAttribute{
						Description: "The optimize dictionary for the table.",
						Optional:    true,
					},
					"optimize_dictionary_for_metrics": schema.BoolAttribute{
						Description: "The optimize dictionary for metrics for the table.",
						Optional:    true,
					},
					"no_dictionary_size_ratio_threshold": schema.Float64Attribute{
						Description: "The no dictionary size ration threshold for the table.",
						Optional:    true,
					},
					"column_min_max_value_generator_mode": schema.StringAttribute{
						Description: "The column min max value generator mode for the table.",
						Optional:    true,
					},
					"segment_name_generator_type": schema.StringAttribute{
						Description: "The segment name generator type for the table.",
						Optional:    true,
					},
					"aggregate_metrics": schema.BoolAttribute{
						Description: "The aggregate metrics for the table.",
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"segment_partition_config": schema.SingleNestedAttribute{
						Description: "The segment partition configuration for the table.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"column_partition_map": schema.MapAttribute{
								Description: "The column partition map for the segment partition config.",
								Optional:    true,
								ElementType: types.MapType{
									ElemType: types.StringType,
								},
							},
						},
					},
					"range_index_columns": schema.ListAttribute{
						Description: "The range index columns for the table.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"no_dictionary_columns": schema.ListAttribute{
						Description: "The no dictionary columns for the table.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"on_heap_dictionary_columns": schema.ListAttribute{
						Description: "The on heap dictionary columns for the table.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"var_length_dictionary_columns": schema.ListAttribute{
						Description: "The var length dictionary columns for the table.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"bloom_filter_columns": schema.ListAttribute{
						Description: "The bloom filter columns for the table.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"range_index_version": schema.Int64Attribute{
						Description: "The range index version for the table.",
						Optional:    true,
					},
				},
			},
			"upsert_config": schema.SingleNestedAttribute{
				Description: "The upsert configuration for the table.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"mode": schema.StringAttribute{
						Description: "The upsert mode for the table.",
						Required:    true,
					},
					"partial_upsert_strategies": schema.MapAttribute{
						Description: "The partial upsert strategies for the table.",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			"ingestion_config": schema.SingleNestedAttribute{
				Description: "ingestion configuration for the table i.e kafka",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"segment_time_value_check": schema.BoolAttribute{
						Description: "segment time value check.",
						Optional:    true,
					},
					"row_time_value_check": schema.BoolAttribute{
						Description: "row time value check.",
						Optional:    true,
					},
					"continue_on_error": schema.BoolAttribute{
						Description: "continue after error ingesting.",
						Optional:    true,
					},
					"stream_ingestion_config": schema.SingleNestedAttribute{
						Description: "stream ingestion configurations",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"stream_config_maps": schema.ListAttribute{
								Description: "stream configuration",
								Optional:    true,
								ElementType: types.MapType{ElemType: types.StringType},
							},
						},
					},
					"transform_configs": schema.ListNestedAttribute{
						Description: "transform configurations",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"column_name": schema.StringAttribute{
									Description: "column name",
									Required:    true,
								},
								"transform_function": schema.StringAttribute{
									Description: "transform function",
									Required:    true,
								},
							},
						},
					},
				},
			},
			"tier_configs": schema.ListNestedAttribute{
				Description: "tier configurations for the table",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "name of the tier",
							Required:    true,
						},
						"segment_selector_type": schema.StringAttribute{
							Description: "segment selector type",
							Required:    true,
						},
						"segment_age": schema.StringAttribute{
							Description: "segment age",
							Required:    true,
						},
						"storage_type": schema.StringAttribute{
							Description: "storage type",
							Required:    true,
						},
						"server_tag": schema.StringAttribute{
							Description: "server tag",
							Required:    true,
						},
					},
				},
			},
			"is_dim_table": schema.BoolAttribute{
				Description: "is dimension table",
				Optional:    true,
			},
			"metadata": schema.SingleNestedAttribute{
				Description: "metadata for the table",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"custom_configs": schema.MapAttribute{
						Description: "custom configs",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			"field_config_list": schema.ListNestedAttribute{
				Description: "field configurations for the table",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "name of the field",
							Required:    true,
						},
						"encoding_type": schema.StringAttribute{
							Description: "encoding type",
							Required:    true,
						},
						"index_type": schema.StringAttribute{
							Description: "index type",
							Required:    true,
						},
						"index_types": schema.ListAttribute{
							Description: "index types",
							Optional:    true,
							ElementType: types.StringType,
						},
						"timestamp_config": schema.SingleNestedAttribute{
							Description: "timestamp configuration",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"granularities": schema.ListAttribute{
									Description: "granularities",
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
						"indexes": schema.SingleNestedAttribute{
							Description: "indexes",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"inverted": schema.SingleNestedAttribute{
									Description: "inverted",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										"enabled": schema.StringAttribute{
											Description: "enabled",
											Required:    true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *tableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan models.TableResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var table model.Table
	err := json.Unmarshal([]byte(plan.Table.ValueString()), &table)
	if err != nil {
		resp.Diagnostics.AddError("Create Failed: Unable to unmarshal table from config", err.Error())
		return
	}

	overriddenTableBytes, err := json.Marshal(override(&plan))
	if err != nil {
		resp.Diagnostics.AddError("Create Failed: Unable to marshal table", err.Error())
		return
	}

	_, err = r.client.CreateTable(overriddenTableBytes)
	if err != nil {
		resp.Diagnostics.AddError("Create Failed: Unable to create table", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state models.TableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tableResponse, err := r.client.GetTable(state.TableName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get table", err.Error())
		return
	}

	var table model.Table

	// if table.OFFLINE is not nil, set the state to populated data
	if tableResponse.OFFLINE.TableName != "" {
		table = tableResponse.OFFLINE
	} else {
		table = tableResponse.REALTIME
	}

	converter.SetStateFromTable(&state, &table)

	// set state to populated data
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *tableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan models.TableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var table model.Table
	err := json.Unmarshal([]byte(plan.Table.ValueString()), &table)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed: Unable to unmarshal table", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Overriding table config: %s", plan.TableName))

	overriddenTableBytes, err := json.Marshal(override(&plan))
	if err != nil {
		resp.Diagnostics.AddError("Update Failed: Unable to marshal table", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating table: %s", plan.TableName))

	_, err = r.client.UpdateTable(plan.TableName.ValueString(), overriddenTableBytes)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed: Unable to update table", err.Error())
		return
	}

	// set state to populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *tableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.TableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting table: %s", state.TableName))

	_, err := r.client.DeleteTable(state.TableName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Failed: Unable to delete table", err.Error())
		return
	}

	// set state to populated data
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func override(plan *models.TableResourceModel) *model.Table {

	table := model.Table{
		TableName:        plan.TableName.ValueString(),
		TableType:        plan.TableType.ValueString(),
		Tenants:          overrideTenantsConfig(plan),
		SegmentsConfig:   overrideSegmentsConfig(plan),
		TableIndexConfig: overrideTableConfigs(plan),
		IsDimTable:       plan.IsDimTable.ValueBool(),
		IngestionConfig:  overrideIngestionConfig(plan),
	}

	if plan.Metadata != nil {
		table.Metadata = overrideMetadata(plan)
	}

	if plan.TierConfigs != nil {
		table.TierConfigs = overrideTierConfigs(plan)
	}

	if plan.FieldConfigList != nil {
		table.FieldConfigList = overrideFieldConfigList(plan)

	}

	return &table
}

func overrideTableConfigs(plan *models.TableResourceModel) model.TableIndexConfig {
	return model.TableIndexConfig{
		CreateInvertedIndexDuringSegmentGeneration: plan.TableIndexConfig.CreateInvertedIndexDuringSegmentGeneration.ValueBool(),
		SortedColumn:                   plan.TableIndexConfig.SortedColumn,
		StarTreeIndexConfigs:           overrideStarTreeConfigs(plan),
		EnableDefaultStarTree:          plan.TableIndexConfig.EnableDefaultStarTree.ValueBool(),
		EnableDynamicStarTreeCreation:  plan.TableIndexConfig.EnableDynamicStarTree.ValueBool(),
		LoadMode:                       plan.TableIndexConfig.LoadMode.ValueString(),
		ColumnMinMaxValueGeneratorMode: plan.TableIndexConfig.ColumnMinMaxValueGeneratorMode.ValueString(),
		NullHandlingEnabled:            plan.TableIndexConfig.NullHandlingEnabled.ValueBool(),
		AggregateMetrics:               plan.TableIndexConfig.AggregateMetrics.ValueBool(),
		OptimizeDictionary:             plan.TableIndexConfig.OptimizeDictionary.ValueBool(),
		OptimizeDictionaryForMetrics:   plan.TableIndexConfig.OptimizeDictionaryForMetrics.ValueBool(),
		NoDictionarySizeRatioThreshold: plan.TableIndexConfig.NoDictionarySizeRatioThreshold.ValueFloat64(),
		SegmentNameGeneratorType:       plan.TableIndexConfig.SegmentNameGeneratorType.ValueString(),
		SegmentPartitionConfig:         overrideSegmentPartitionConfig(plan),
		RangeIndexColumns:              plan.TableIndexConfig.RangeIndexColumns,
		NoDictionaryColumns:            plan.TableIndexConfig.NoDictionaryColumns,
		RangeIndexVersion:              int(plan.TableIndexConfig.RangeIndexVersion.ValueInt64()),
		OnHeapDictionaryColumns:        plan.TableIndexConfig.OnHeapDictionaryColumns,
		VarLengthDictionaryColumns:     plan.TableIndexConfig.VarLengthDictionaryColumns,
		BloomFilterColumns:             plan.TableIndexConfig.BloomFilterColumns,
	}

}

func overrideSegmentPartitionConfig(plan *models.TableResourceModel) *model.SegmentPartitionConfig {

	if plan.TableIndexConfig.SegmentPartitionConfig == nil {
		return nil
	}

	columnPartitionMap := make(map[string]model.ColumnPartitionMapConfig, 1)
	for key, value := range plan.TableIndexConfig.SegmentPartitionConfig.ColumnPartitionMap {

		numPartitions, err := strconv.Atoi(value["numPartitions"])
		if err != nil {
			log.Panic(err)
		}

		columnPartitionMap[key] = model.ColumnPartitionMapConfig{
			FunctionName: value["functionName"],
			// convert to int
			NumPartitions: numPartitions,
		}
	}
	return &model.SegmentPartitionConfig{ColumnPartitionMap: columnPartitionMap}
}

func overrideStarTreeConfigs(plan *models.TableResourceModel) []*model.StarTreeIndexConfig {

	if plan.TableIndexConfig.StarTreeIndexConfigs == nil {
		return nil
	}

	var starTreeConfigs []*model.StarTreeIndexConfig
	for _, starConfig := range plan.TableIndexConfig.StarTreeIndexConfigs {
		starTreeConfigs = append(starTreeConfigs, &model.StarTreeIndexConfig{
			MaxLeafRecords:                    int(starConfig.MaxLeafRecords.ValueInt64()),
			DimensionsSplitOrder:              starConfig.DimensionsSplitOrder,
			FunctionColumnPairs:               starConfig.FunctionColumnPairs,
			SkipStarNodeCreationForDimensions: starConfig.SkipStarNodeCreationForDimNames,
		})
	}
	return starTreeConfigs
}

func overrideSegmentsConfig(plan *models.TableResourceModel) model.TableSegmentsConfig {

	return model.TableSegmentsConfig{
		TimeType:           plan.SegmentsConfig.TimeType.ValueString(),
		Replication:        plan.SegmentsConfig.Replication.ValueString(),
		TimeColumnName:     plan.SegmentsConfig.TimeColumnName.ValueString(),
		RetentionTimeUnit:  plan.SegmentsConfig.RetentionTimeUnit.ValueString(),
		RetentionTimeValue: plan.SegmentsConfig.RetentionTimeValue.ValueString(),
	}

}

func overrideTenantsConfig(plan *models.TableResourceModel) model.TableTenant {
	return model.TableTenant{
		Broker: plan.TenantsConfig.Broker.ValueString(),
		Server: plan.TenantsConfig.Server.ValueString(),
	}
}

func overrideTierConfigs(plan *models.TableResourceModel) []*model.TierConfig {

	if plan.TierConfigs == nil {
		return nil
	}

	var tierConfigs []*model.TierConfig
	for _, tierConfig := range plan.TierConfigs {
		tierConfigs = append(tierConfigs, &model.TierConfig{
			Name:                tierConfig.Name.ValueString(),
			SegmentSelectorType: tierConfig.SegmentSelectorType.ValueString(),
			SegmentAge:          tierConfig.SegmentAge.ValueString(),
			StorageType:         tierConfig.StorageType.ValueString(),
			ServerTag:           tierConfig.ServerTag.ValueString(),
		})
	}
	return tierConfigs
}

func overrideMetadata(plan *models.TableResourceModel) *model.TableMetadata {

	if plan.Metadata == nil {
		return nil
	}

	return &model.TableMetadata{
		CustomConfigs: plan.Metadata.CustomConfigs,
	}
}

func overrideIngestionConfig(plan *models.TableResourceModel) *model.TableIngestionConfig {

	if plan.IngestionConfig == nil {
		return nil
	}

	return &model.TableIngestionConfig{
		SegmentTimeValueCheck: plan.IngestionConfig.SegmentTimeValueCheck.ValueBool(),
		RowTimeValueCheck:     plan.IngestionConfig.RowTimeValueCheck.ValueBool(),
		ContinueOnError:       plan.IngestionConfig.ContinueOnError.ValueBool(),
		StreamIngestionConfig: &model.StreamIngestionConfig{
			StreamConfigMaps: plan.IngestionConfig.StreamIngestionConfig.StreamConfigMaps,
		},
	}
}

func overrideFieldConfigList(plan *models.TableResourceModel) []model.FieldConfig {

	if plan.FieldConfigList == nil {
		return nil
	}

	var fieldConfigs []model.FieldConfig
	for _, fieldConfig := range plan.FieldConfigList {

		fc := model.FieldConfig{
			Name:         fieldConfig.Name.ValueString(),
			EncodingType: fieldConfig.EncodingType.ValueString(),
			IndexType:    fieldConfig.IndexType.ValueString(),
			IndexTypes:   fieldConfig.IndexTypes,
		}

		if fieldConfig.TimestampConfig != nil {
			fc.TimestampConfig = &model.TimestampConfig{
				Granularities: fieldConfig.TimestampConfig.Granularities,
			}
		}

		if fieldConfig.Indexes != nil {
			fc.Indexes = &model.FieldIndexes{
				Inverted: &model.FiendIndexInverted{
					Enabled: fieldConfig.Indexes.Inverted.Enabled.ValueString(),
				},
			}
		}

		fieldConfigs = append(fieldConfigs, fc)
	}
	return fieldConfigs
}
