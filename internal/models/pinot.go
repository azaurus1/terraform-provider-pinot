package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type TenantsConfig struct {
	Broker            types.String   `tfsdk:"broker"`
	Server            types.String   `tfsdk:"server"`
	TagOverrideConfig *types.MapType `tfsdk:"tag_override_config"`
}

type SegmentsConfig struct {
	TimeType                      types.String `tfsdk:"time_type"`
	Replication                   types.String `tfsdk:"replication"`
	TimeColumnName                types.String `tfsdk:"time_column_name"`
	SegmentAssignmentStrategy     types.String `tfsdk:"segment_assignment_strategy"`
	SegmentPushType               types.String `tfsdk:"segment_push_type"`
	MinimizeDataMovement          types.Bool   `tfsdk:"minimize_data_movement"`
	RetentionTimeUnit             types.String `tfsdk:"retention_time_unit"`
	RetentionTimeValue            types.String `tfsdk:"retention_time_value"`
	DeletedSegmentRetentionPeriod types.String `tfsdk:"deleted_segment_retention_period"`
}

type UpsertConfig struct {
	Mode                  types.String  `tfsdk:"mode"`
	PartialUpsertStrategy types.MapType `tfsdk:"partial_upsert_strategy"`
}

type TableIndexConfig struct {
	InvertedIndexColumns                       []string                `tfsdk:"inverted_index_columns"`
	SortedColumn                               []string                `tfsdk:"sorted_column"`
	NoDictionaryColumns                        []string                `tfsdk:"no_dictionary_columns"`
	VarLengthDictionaryColumns                 []string                `tfsdk:"var_length_dictionary_columns"`
	RangeIndexColumns                          []string                `tfsdk:"range_index_columns"`
	LoadMode                                   types.String            `tfsdk:"load_mode"`
	NullHandlingEnabled                        types.Bool              `tfsdk:"null_handling_enabled"`
	CreateInvertedIndexDuringSegmentGeneration types.Bool              `tfsdk:"create_inverted_index_during_segment_generation"`
	EnableDynamicStarTree                      types.Bool              `tfsdk:"enable_dynamic_star_tree"`
	EnableDefaultStarTree                      types.Bool              `tfsdk:"enable_default_star_tree"`
	OptimizeDictionary                         types.Bool              `tfsdk:"optimize_dictionary"`
	OptimizeDictionaryForMetrics               types.Bool              `tfsdk:"optimize_dictionary_for_metrics"`
	NoDictionarySizeRatioThreshold             types.Float64           `tfsdk:"no_dictionary_size_ratio_threshold"`
	ColumnMinMaxValueGeneratorMode             types.String            `tfsdk:"column_min_max_value_generator_mode"`
	SegmentNameGeneratorType                   types.String            `tfsdk:"segment_name_generator_type"`
	StarTreeIndexConfigs                       []*StarTreeIndexConfigs `tfsdk:"star_tree_index_configs"`
}

type AggregationConfig struct {
	AggregationFunction types.String `tfsdk:"aggregation_function"`
	ColumnName          types.String `tfsdk:"column_name"`
	CompressionCodec    types.String `tfsdk:"compression_codec"`
}

type StarTreeIndexConfigs struct {
	MaxLeafRecords                  types.Int64          `tfsdk:"max_leaf_records"`
	SkipStarNodeCreationForDimNames []string             `tfsdk:"skip_star_node_creation_for_dim_names"`
	DimensionsSplitOrder            []string             `tfsdk:"dimensions_split_order"`
	FunctionColumnPairs             []string             `tfsdk:"function_column_pairs"`
	AggregationConfigs              []*AggregationConfig `tfsdk:"aggregation_configs"`
}

type IngestionConfig struct {
	SegmentTimeValueCheck types.Bool             `tfsdk:"segment_time_value_check"`
	RowTimeValueCheck     types.Bool             `tfsdk:"row_time_value_check"`
	ContinueOnError       types.Bool             `tfsdk:"continue_on_error"`
	StreamIngestionConfig *StreamIngestionConfig `tfsdk:"stream_ingestion_config"`
}

type StreamIngestionConfig struct {
	StreamConfigMaps []map[string]string `tfsdk:"stream_config_maps"`
}

type TierConfig struct {
	Name                types.String `tfsdk:"name"`
	StorageType         types.String `tfsdk:"storage_type"`
	SegmentSelectorType types.String `tfsdk:"segment_selector_type"`
	SegmentAge          types.String `tfsdk:"segment_age"`
	ServerTag           types.String `tfsdk:"server_tag"`
}

type Metadata struct {
	CustomConfigs map[string]string `tfsdk:"custom_configs"`
}
