package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type TableResourceModel struct {
	TableName        types.String      `tfsdk:"table_name"`
	Table            types.String      `tfsdk:"table"`
	TableType        types.String      `tfsdk:"table_type"`
	SegmentsConfig   *SegmentsConfig   `tfsdk:"segments_config"`
	TenantsConfig    *TenantsConfig    `tfsdk:"tenants"`
	TableIndexConfig *TableIndexConfig `tfsdk:"table_index_config"`
	UpsertConfig     *UpsertConfig     `tfsdk:"upsert_config"`
	IngestionConfig  *IngestionConfig  `tfsdk:"ingestion_config"`
	TierConfigs      []*TierConfig     `tfsdk:"tier_configs"`
	IsDimTable       types.Bool        `tfsdk:"is_dim_table"`
	Metadata         *Metadata         `tfsdk:"metadata"`
	FieldConfigList  []*FieldConfig    `tfsdk:"field_config_list"`
	Routing          *RoutingConfig    `tfsdk:"routing"`
}

type TenantsConfig struct {
	Broker            types.String   `tfsdk:"broker"`
	Server            types.String   `tfsdk:"server"`
	TagOverrideConfig *types.MapType `tfsdk:"tag_override_config"`
}

type SegmentsConfig struct {
	TimeType                       types.String `tfsdk:"time_type"`
	Replication                    types.String `tfsdk:"replication"`
	ReplicasPerPartition           types.String `tfsdk:"replicas_per_partition"`
	TimeColumnName                 types.String `tfsdk:"time_column_name"`
	RetentionTimeUnit              types.String `tfsdk:"retention_time_unit"`
	RetentionTimeValue             types.String `tfsdk:"retention_time_value"`
	DeletedSegmentsRetentionPeriod types.String `tfsdk:"deleted_segments_retention_period"`
}

type UpsertConfig struct {
	Mode                         types.String `tfsdk:"mode"`
	PartialUpsertStrategy        types.Map    `tfsdk:"partial_upsert_strategies"`
	DefaultPartialUpsertStrategy types.String `tfsdk:"default_partial_upsert_strategy"`
	DeletedRecordColumn          types.String `tfsdk:"delete_record_column"`
	DeletedKeysTTL               types.Int64  `tfsdk:"deleted_keys_ttl"`
	MetadataTTL                  types.Int64  `tfsdk:"metadata_ttl"`
	HashFunction                 types.String `tfsdk:"hash_function"`
	EnableSnapshot               types.Bool   `tfsdk:"enable_snapshot"`
	EnablePreLoad                types.Bool   `tfsdk:"enable_preload"`
	DropOutOfOrderRecord         types.Bool   `tfsdk:"drop_out_of_order_record"`
	OutOfOrderRecordColumn       types.String `tfsdk:"out_of_order_record_column"`
	MetadataManagerClass         types.String `tfsdk:"metadata_manager_class"`
	MetadataManagerConfigs       types.Map    `tfsdk:"metadata_manager_configs"`
}

type SegmentPartitionConfig struct {
	ColumnPartitionMap map[string]map[string]string `tfsdk:"column_partition_map"`
}

type TimestampConfig struct {
	Granularities []string `tfsdk:"granularities"`
}

type FiendIndexInverted struct {
	Enabled types.String `tfsdk:"enabled"`
}

type FieldIndexes struct {
	Inverted *FiendIndexInverted `tfsdk:"inverted"`
}

type FieldConfig struct {
	Name            types.String     `tfsdk:"name"`
	EncodingType    types.String     `tfsdk:"encoding_type"`
	IndexType       types.String     `tfsdk:"index_type"`
	IndexTypes      []string         `tfsdk:"index_types"`
	TimestampConfig *TimestampConfig `tfsdk:"timestamp_config"`
	Indexes         *FieldIndexes    `tfsdk:"indexes"`
}

type TableIndexConfig struct {
	SortedColumn                               types.List              `tfsdk:"sorted_column"`
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
	AggregateMetrics                           types.Bool              `tfsdk:"aggregate_metrics"`
	StarTreeIndexConfigs                       []*StarTreeIndexConfigs `tfsdk:"star_tree_index_configs"`
	SegmentPartitionConfig                     *SegmentPartitionConfig `tfsdk:"segment_partition_config"`
	NoDictionaryColumns                        types.List              `tfsdk:"no_dictionary_columns"`
	RangeIndexColumns                          types.List              `tfsdk:"range_index_columns"`
	OnHeapDictionaryColumns                    types.List              `tfsdk:"on_heap_dictionary_columns"`
	VarLengthDictionaryColumns                 types.List              `tfsdk:"var_length_dictionary_columns"`
	BloomFilterColumns                         types.List              `tfsdk:"bloom_filter_columns"`
	RangeIndexVersion                          types.Int64             `tfsdk:"range_index_version"`
}

type AggregationConfig struct {
	AggregationFunction types.String `tfsdk:"aggregation_function"`
	ColumnName          types.String `tfsdk:"column_name"`
	CompressionCodec    types.String `tfsdk:"compression_codec"`
}

type RoutingConfig struct {
	SegmentPrunerTypes   types.List   `tfsdk:"segment_pruner_types"`
	InstanceSelectorType types.String `tfsdk:"instance_selector_type"`
}

type StarTreeIndexConfigs struct {
	MaxLeafRecords                  types.Int64          `tfsdk:"max_leaf_records"`
	SkipStarNodeCreationForDimNames types.List           `tfsdk:"skip_star_node_creation_for_dim_names"`
	DimensionsSplitOrder            types.List           `tfsdk:"dimensions_split_order"`
	FunctionColumnPairs             types.List           `tfsdk:"function_column_pairs"`
	AggregationConfigs              []*AggregationConfig `tfsdk:"aggregation_configs"`
}

type TransformConfig struct {
	ColumnName        types.String `tfsdk:"column_name"`
	TransformFunction types.String `tfsdk:"transform_function"`
}

type FilterConfig struct {
	FilterFunction types.String `tfsdk:"filter_function"`
}

type IngestionConfig struct {
	SegmentTimeValueCheck types.Bool             `tfsdk:"segment_time_value_check"`
	RowTimeValueCheck     types.Bool             `tfsdk:"row_time_value_check"`
	ContinueOnError       types.Bool             `tfsdk:"continue_on_error"`
	StreamIngestionConfig *StreamIngestionConfig `tfsdk:"stream_ingestion_config"`
	TransformConfigs      []*TransformConfig     `tfsdk:"transform_configs"`
	FilterConfig          *FilterConfig          `tfsdk:"filter_config"`
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
