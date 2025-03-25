package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type TableResourceModel struct {
	TableName                   types.String                 `tfsdk:"table_name"`
	Table                       types.String                 `tfsdk:"table"`
	TableType                   types.String                 `tfsdk:"table_type"`
	SegmentsConfig              *SegmentsConfig              `tfsdk:"segments_config"`
	TenantsConfig               *TenantsConfig               `tfsdk:"tenants"`
	TableIndexConfig            *TableIndexConfig            `tfsdk:"table_index_config"`
	UpsertConfig                *UpsertConfig                `tfsdk:"upsert_config"`
	IngestionConfig             *IngestionConfig             `tfsdk:"ingestion_config"`
	TierConfigs                 []*TierConfig                `tfsdk:"tier_configs"`
	IsDimTable                  types.Bool                   `tfsdk:"is_dim_table"`
	Metadata                    *Metadata                    `tfsdk:"metadata"`
	FieldConfigList             []*FieldConfig               `tfsdk:"field_config_list"`
	Routing                     *RoutingConfig               `tfsdk:"routing"`
	Task                        *Task                        `tfsdk:"task"`
	InstanceAssignmentConfigMap *InstanceAssignmentConfigMap `tfsdk:"instance_assignment_config_map"`
}

type TenantsConfig struct {
	Broker            types.String             `tfsdk:"broker"`
	Server            types.String             `tfsdk:"server"`
	TagOverrideConfig *TenantTagOverrideConfig `tfsdk:"tag_override_config"`
}

type TenantTagOverrideConfig struct {
	RealtimeConsuming types.String `tfsdk:"realtime_consuming"`
	RealtimeCompleted types.String `tfsdk:"realtime_completed"`
}

type SegmentsConfig struct {
	TimeType                       types.String      `tfsdk:"time_type"`
	Replication                    types.String      `tfsdk:"replication"`
	ReplicasPerPartition           types.String      `tfsdk:"replicas_per_partition"`
	TimeColumnName                 types.String      `tfsdk:"time_column_name"`
	RetentionTimeUnit              types.String      `tfsdk:"retention_time_unit"`
	RetentionTimeValue             types.String      `tfsdk:"retention_time_value"`
	DeletedSegmentsRetentionPeriod types.String      `tfsdk:"deleted_segments_retention_period"`
	SchemaName                     types.String      `tfsdk:"schema_name"`
	CompletionConfig               *CompletionConfig `tfsdk:"completion_config"`
	PeerSegmentDownloadScheme      types.String      `tfsdk:"peer_segment_download_scheme"`
	SegmentPushType                types.String      `tfsdk:"segment_push_type"`
	MinimizeDataMovement           types.Bool        `tfsdk:"minimize_data_movement"`
}

type CompletionConfig struct {
	CompletionMode types.String `tfsdk:"completion_mode"`
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
	AutoGeneratedInvertedIndex                 types.Bool              `tfsdk:"auto_generated_inverted_index"`
	StarTreeIndexConfigs                       []*StarTreeIndexConfigs `tfsdk:"star_tree_index_configs"`
	SegmentPartitionConfig                     *SegmentPartitionConfig `tfsdk:"segment_partition_config"`
	NoDictionaryColumns                        types.List              `tfsdk:"no_dictionary_columns"`
	RangeIndexColumns                          types.List              `tfsdk:"range_index_columns"`
	OnHeapDictionaryColumns                    types.List              `tfsdk:"on_heap_dictionary_columns"`
	VarLengthDictionaryColumns                 types.List              `tfsdk:"var_length_dictionary_columns"`
	BloomFilterColumns                         types.List              `tfsdk:"bloom_filter_columns"`
	InvertedIndexColumns                       types.List              `tfsdk:"inverted_index_columns"`
	JsonIndexColumns                           types.List              `tfsdk:"json_index_columns"`
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
	StreamConfigMaps                 []*StreamConfig `tfsdk:"stream_config_maps"`
	ColumnMajorSegmentBuilderEnabled types.Bool      `tfsdk:"column_major_segment_builder_enabled"`
	TrackFilteredMessageOffsets      types.Bool      `tfsdk:"track_filtered_message_offsets"`
}

type StreamConfig struct {
	AccessKey                                                        types.String `tfsdk:"access_key"`
	AuthenticationType                                               types.String `tfsdk:"authentication_type"`
	KeySerializer                                                    types.String `tfsdk:"key_serializer"`
	MaxRecordsToFetch                                                types.String `tfsdk:"max_records_to_fetch"`
	RealtimeSegmentCommitTimeoutSeconds                              types.String `tfsdk:"realtime_segment_commit_timeout_seconds"`
	RealtimeSegmentFlushAutotuneInitialRows                          types.String `tfsdk:"realtime_segment_flush_autotune_initial_rows"`
	RealtimeSegmentFlushDesiredSize                                  types.String `tfsdk:"realtime_segment_flush_desired_size"`
	RealtimeSegmentFlushThresholdRows                                types.String `tfsdk:"realtime_segment_flush_threshold_rows"`
	RealtimeSegmentFlushThresholdTime                                types.String `tfsdk:"realtime_segment_flush_threshold_time"`
	RealtimeSegmentFlushThresholdSegmentRows                         types.String `tfsdk:"realtime_segment_flush_threshold_segment_rows"`
	RealtimeSegmentFlushThresholdSegmentSize                         types.String `tfsdk:"realtime_segment_flush_threshold_segment_size"`
	RealtimeSegmentFlushThresholdSegmentTime                         types.String `tfsdk:"realtime_segment_flush_threshold_segment_time"`
	RealtimeSegmentServerUploadToDeepStore                           types.String `tfsdk:"realtime_segment_server_upload_to_deep_store"`
	Region                                                           types.String `tfsdk:"region"`
	SaslJaasConfig                                                   types.String `tfsdk:"sasl_jaas_config"`
	SaslMechanism                                                    types.String `tfsdk:"sasl_mechanism"`
	SecretKey                                                        types.String `tfsdk:"secret_key"`
	SecurityProtocol                                                 types.String `tfsdk:"security_protocol"`
	ShardIteratorType                                                types.String `tfsdk:"shard_iterator_type"`
	StreamType                                                       types.String `tfsdk:"stream_type"`
	SslKeyPassword                                                   types.String `tfsdk:"ssl_key_password"`
	SslKeystoreLocation                                              types.String `tfsdk:"ssl_keystore_location"`
	SslKeystorePassword                                              types.String `tfsdk:"ssl_keystore_password"`
	SslKeystoreType                                                  types.String `tfsdk:"ssl_keystore_type"`
	SslTruststoreLocation                                            types.String `tfsdk:"ssl_truststore_location"`
	SslTruststorePassword                                            types.String `tfsdk:"ssl_truststore_password"`
	SslTruststoreType                                                types.String `tfsdk:"ssl_truststore_type"`
	StreamKafkaBrokerList                                            types.String `tfsdk:"stream_kafka_broker_list"`
	StreamKafkaBufferSize                                            types.String `tfsdk:"stream_kafka_buffer_size"`
	StreamKafkaConsumerFactoryClassName                              types.String `tfsdk:"stream_kafka_consumer_factory_class_name"`
	StreamKafkaConsumerPropAutoOffsetReset                           types.String `tfsdk:"stream_kafka_consumer_prop_auto_offset_reset"`
	StreamKafkaConsumerType                                          types.String `tfsdk:"stream_kafka_consumer_type"`
	StreamKafkaFetchTimeoutMillis                                    types.String `tfsdk:"stream_kafka_fetch_timeout_millis"`
	StreamKafkaConnectionTimeoutMillis                               types.String `tfsdk:"stream_kafka_connection_timeout_millis"`
	StreamKafkaDecoderClassName                                      types.String `tfsdk:"stream_kafka_decoder_class_name"`
	StreamKafkaDecoderPropBasicAuthCredentialsSource                 types.String `tfsdk:"stream_kafka_decoder_prop_basic_auth_credentials_source"`
	StreamKafkaDecoderPropDescriptorFile                             types.String `tfsdk:"stream_kafka_decoder_prop_descriptor_file"`
	StreamKafkaDecoderPropProtoClassName                             types.String `tfsdk:"stream_kafka_decoder_prop_proto_class_name"`
	StreamKafkaDecoderPropFormat                                     types.String `tfsdk:"stream_kafka_decoder_prop_format"`
	StreamKafkaDecoderPropSchemaRegistryBasicAuthUserInfo            types.String `tfsdk:"stream_kafka_decoder_prop_schema_registry_basic_auth_user_info"`
	StreamKafkaDecoderPropSchemaRegistryBasicAuthCredentialsSource   types.String `tfsdk:"stream_kafka_decoder_prop_schema_registry_basic_auth_credentials_source"`
	StreamKafkaDecoderPropSchemaRegistryRestUrl                      types.String `tfsdk:"stream_kafka_decoder_prop_schema_registry_rest_url"`
	StreamKafkaDecoderPropSchemaRegistrySchemaName                   types.String `tfsdk:"stream_kafka_decoder_prop_schema_registry_schema_name"`
	StreamKafkaDecoderPropSchemaRegistrySslKeystoreLocation          types.String `tfsdk:"stream_kafka_decoder_prop_schema_registry_ssl_keystore_location"`
	StreamKafkaDecoderPropSchemaRegistrySslKeystorePassword          types.String `tfsdk:"stream_kafka_decoder_prop_schema_registry_ssl_keystore_password"`
	StreamKafkaDecoderPropSchemaRegistrySslKeystoreType              types.String `tfsdk:"stream_kafka_decoder_prop_schema_registry_ssl_keystore_type"`
	StreamKafkaDecoderPropSchemaRegistrySslTruststoreLocation        types.String `tfsdk:"stream_kafka_decoder_prop_schema_registry_ssl_truststore_location"`
	StreamKafkaDecoderPropSchemaRegistrySslTruststorePassword        types.String `tfsdk:"stream_kafka_decoder_prop_schema_registry_ssl_truststore_password"`
	StreamKafkaDecoderPropSchemaRegistrySslTruststoreType            types.String `tfsdk:"stream_kafka_decoder_prop_schema_registry_ssl_truststore_type"`
	StreamKafkaDecoderPropSchemaRegistrySslProtocol                  types.String `tfsdk:"stream_kafka_decoder_prop_schema_registry_ssl_protocol"`
	StreamKafkaFetcherMinBytes                                       types.String `tfsdk:"stream_kafka_fetcher_min_bytes"`
	StreamKafkaFetcherSize                                           types.String `tfsdk:"stream_kafka_fetcher_size"`
	StreamKafkaHlcGroupId                                            types.String `tfsdk:"stream_kafka_hlc_group_id"`
	StreamKafkaIdleTimeoutMillis                                     types.String `tfsdk:"stream_kafka_idle_timeout_millis"`
	StreamKafkaIsolationLevel                                        types.String `tfsdk:"stream_kafka_isolation_level"`
	StreamKafkaMetadataPopulate                                      types.String `tfsdk:"stream_kafka_metadata_populate"`
	StreamKafkaSchemaRegistryUrl                                     types.String `tfsdk:"stream_kafka_schema_registry_url"`
	StreamKafkaSocketTimeout                                         types.String `tfsdk:"stream_kafka_socket_timeout"`
	StreamKafkaSslCertificateType                                    types.String `tfsdk:"stream_kafka_ssl_certificate_type"`
	StreamKafkaSslClientCertificate                                  types.String `tfsdk:"stream_kafka_ssl_client_certificate"`
	StreamKafkaSslClientKey                                          types.String `tfsdk:"stream_kafka_ssl_client_key"`
	StreamKafkaSslClientKeyAlgorithm                                 types.String `tfsdk:"stream_kafka_ssl_client_key_algorithm"`
	StreamKafkaSslServerCertificate                                  types.String `tfsdk:"stream_kafka_ssl_server_certificate"`
	StreamKafkaTopicName                                             types.String `tfsdk:"stream_kafka_topic_name"`
	StreamKafkaZkBrokerUrl                                           types.String `tfsdk:"stream_kafka_zk_broker_url"`
	StreamKinesisConsumerFactoryClassName                            types.String `tfsdk:"stream_kinesis_consumer_factory_class_name"`
	StreamKinesisConsumerType                                        types.String `tfsdk:"stream_kinesis_consumer_type"`
	StreamKinesisDecoderClassName                                    types.String `tfsdk:"stream_kinesis_decoder_class_name"`
	StreamKinesisFetchTimeoutMillis                                  types.String `tfsdk:"stream_kinesis_fetch_timeout_millis"`
	StreamKinesisDecoderPropSchemaRegistryBasicAuthUserInfo          types.String `tfsdk:"stream_kinesis_decoder_prop_schema_registry_basic_auth_user_info"`
	StreamKinesisDecoderPropSchemaRegistryBasicAuthCredentialsSource types.String `tfsdk:"stream_kinesis_decoder_prop_schema_registry_basic_auth_credentials_source"`
	StreamKinesisDecoderPropSchemaRegistryRestUrl                    types.String `tfsdk:"stream_kinesis_decoder_prop_schema_registry_rest_url"`
	StreamKinesisDecoderPropSchemaRegistrySchemaName                 types.String `tfsdk:"stream_kinesis_decoder_prop_schema_registry_schema_name"`
	StreamKinesisDecoderPropSchemaRegistrySslKeystoreLocation        types.String `tfsdk:"stream_kinesis_decoder_prop_schema_registry_ssl_keystore_location"`
	StreamKinesisDecoderPropSchemaRegistrySslKeystorePassword        types.String `tfsdk:"stream_kinesis_decoder_prop_schema_registry_ssl_keystore_password"`
	StreamKinesisDecoderPropSchemaRegistrySslKeystoreType            types.String `tfsdk:"stream_kinesis_decoder_prop_schema_registry_ssl_keystore_type"`
	StreamKinesisDecoderPropSchemaRegistrySslTruststoreLocation      types.String `tfsdk:"stream_kinesis_decoder_prop_schema_registry_ssl_truststore_location"`
	StreamKinesisDecoderPropSchemaRegistrySslTruststorePassword      types.String `tfsdk:"stream_kinesis_decoder_prop_schema_registry_ssl_truststore_password"`
	StreamKinesisDecoderPropSchemaRegistrySslTruststoreType          types.String `tfsdk:"stream_kinesis_decoder_prop_schema_registry_ssl_truststore_type"`
	StreamKinesisDecoderPropSchemaRegistrySslProtocol                types.String `tfsdk:"stream_kinesis_decoder_prop_schema_registry_ssl_protocol"`
	StreamKinesisTopicName                                           types.String `tfsdk:"stream_kinesis_topic_name"`
	StreamPulsarAudience                                             types.String `tfsdk:"stream_pulsar_audience"`
	StreamPulsarAuthenticationToken                                  types.String `tfsdk:"stream_pulsar_authentication_token"`
	StreamPulsarBootstrapServers                                     types.String `tfsdk:"stream_pulsar_bootstrap_servers"`
	StreamPulsarCredsFilePath                                        types.String `tfsdk:"stream_pulsar_creds_file_path"`
	StreamPulsarConsumerFactoryClassName                             types.String `tfsdk:"stream_pulsar_consumer_factory_class_name"`
	StreamPulsarConsumerPropAutoOffsetReset                          types.String `tfsdk:"stream_pulsar_consumer_prop_auto_offset_reset"`
	StreamPulsarConsumerType                                         types.String `tfsdk:"stream_pulsar_consumer_type"`
	StreamPulsarDecoderClassName                                     types.String `tfsdk:"stream_pulsar_decoder_class_name"`
	StreamPulsarFetchTimeoutMillis                                   types.String `tfsdk:"stream_pulsar_fetch_timeout_millis"`
	StreamPulsarIssuerUrl                                            types.String `tfsdk:"stream_pulsar_issuer_url"`
	StreamPulsarMetadataPopulate                                     types.String `tfsdk:"stream_pulsar_metadata_populate"`
	StreamPulsarMetadataFields                                       types.String `tfsdk:"stream_pulsar_metadata_fields"`
	StreamPulsarDecoderPropSchemaRegistryBasicAuthUserInfo           types.String `tfsdk:"stream_pulsar_decoder_prop_schema_registry_basic_auth_user_info"`
	StreamPulsarDecoderPropSchemaRegistryBasicAuthCredentialsSource  types.String `tfsdk:"stream_pulsar_decoder_prop_schema_registry_basic_auth_credentials_source"`
	StreamPulsarDecoderPropSchemaRegistryRestUrl                     types.String `tfsdk:"stream_pulsar_decoder_prop_schema_registry_rest_url"`
	StreamPulsarDecoderPropSchemaRegistrySchemaName                  types.String `tfsdk:"stream_pulsar_decoder_prop_schema_registry_schema_name"`
	StreamPulsarDecoderPropSchemaRegistrySslKeystoreLocation         types.String `tfsdk:"stream_pulsar_decoder_prop_schema_registry_ssl_keystore_location"`
	StreamPulsarDecoderPropSchemaRegistrySslKeystorePassword         types.String `tfsdk:"stream_pulsar_decoder_prop_schema_registry_ssl_keystore_password"`
	StreamPulsarDecoderPropSchemaRegistrySslKeystoreType             types.String `tfsdk:"stream_pulsar_decoder_prop_schema_registry_ssl_keystore_type"`
	StreamPulsarDecoderPropSchemaRegistrySslTruststoreLocation       types.String `tfsdk:"stream_pulsar_decoder_prop_schema_registry_ssl_truststore_location"`
	StreamPulsarDecoderPropSchemaRegistrySslTruststorePassword       types.String `tfsdk:"stream_pulsar_decoder_prop_schema_registry_ssl_truststore_password"`
	StreamPulsarDecoderPropSchemaRegistrySslTruststoreType           types.String `tfsdk:"stream_pulsar_decoder_prop_schema_registry_ssl_truststore_type"`
	StreamPulsarDecoderPropSchemaRegistrySslProtocol                 types.String `tfsdk:"stream_pulsar_decoder_prop_schema_registry_ssl_protocol"`
	StreamPulsarTlsTrustCertsFilePath                                types.String `tfsdk:"stream_pulsar_tls_trust_certs_file_path"`
	StreamPulsarTopicName                                            types.String `tfsdk:"stream_pulsar_topic_name"`
	TopicConsumptionRateLimit                                        types.String `tfsdk:"topic_consumption_rate_limit"`
	ValueSerializer                                                  types.String `tfsdk:"value_serializer"`
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

type Task struct {
	TaskTypeConfigsMap map[string]map[string]string `tfsdk:"task_type_configs_map"` // Changed from TaskTypeConfigMap
}

type InstanceAssignmentConfigMap struct {
	Consuming *InstanceAssignment `tfsdk:"consuming"`
	Completed *InstanceAssignment `tfsdk:"completed"`
	Offline   *InstanceAssignment `tfsdk:"offline"`
}

type InstanceAssignment struct {
	TagPoolConfig               *TagPoolConfigInstanceAssignment         `tfsdk:"tag_pool_config"`
	ReplicaGroupPartitionConfig *ReplicaGroupPartitionInstanceAssignment `tfsdk:"replica_group_partition_config"`
	PartitionSelector           types.String                             `tfsdk:"partition_selector"`
	MinimizeDataMovement        types.Bool                               `tfsdk:"minimize_data_movement"`
}

type TagPoolConfigInstanceAssignment struct {
	Tag       types.String `tfsdk:"tag"`
	PoolBased types.Bool   `tfsdk:"pool_based"`
	NumPools  types.Int64  `tfsdk:"num_pools"`
}

type ReplicaGroupPartitionInstanceAssignment struct {
	ReplicaGroupBased           types.Bool   `tfsdk:"replica_group_based"`
	NumInstances                types.Int64  `tfsdk:"num_instances"`
	NumReplicaGroups            types.Int64  `tfsdk:"num_replica_groups"`
	NumInstancesPerReplicaGroup types.Int64  `tfsdk:"num_instances_per_replica_group"`
	NumPartitions               types.Int64  `tfsdk:"num_partitions"`
	NumInstancesPerPartitions   types.Int64  `tfsdk:"num_instances_per_partition"`
	PartitionColumn             types.String `tfsdk:"partition_column"`
	MinimizeDataMovement        types.Bool   `tfsdk:"minimize_data_movement"`
}
