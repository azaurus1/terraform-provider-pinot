package tf_schema

import (
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SegmentsConfig() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
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
			"schema_name": schema.StringAttribute{
				Description: "The name of the used schema for table",
				Optional:    true,
			},
			"time_type": schema.StringAttribute{
				Description: "The time type for the segments.",
				Optional:    true,
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
			"peer_segment_download_scheme": schema.StringAttribute{
				Description: "",
				Optional:    true,
			},
			"segment_push_type": schema.StringAttribute{
				Description: "",
				Optional:    true,
			},
			"minimize_data_movement": schema.BoolAttribute{
				Description: "",
				Optional:    true,
			},
			"completion_config": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"completion_mode": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func Tenants() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
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
	}
}

func TableIndexConfig() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
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
			"star_tree_index_configs": StarTreeIndexConfigs(),
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
			"auto_generated_inverted_index": schema.BoolAttribute{
				Description: "",
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
			"inverted_index_columns": schema.ListAttribute{
				Description: "The inverted index columns for the table.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"json_index_columns": schema.ListAttribute{
				Description: "The json index columns for the table.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"range_index_version": schema.Int64Attribute{
				Description: "The range index version for the table.",
				Optional:    true,
			},
		},
	}
}

func UpsertConfig() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "The upsert configuration for the table.",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"mode": schema.StringAttribute{
				Description: "The upsert mode for the table.",
				Optional:    true,
			},
			"partial_upsert_strategies": schema.MapAttribute{
				Description: "The partial upsert strategies for the table.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"metadata_manager_class": schema.StringAttribute{
				Description: "The metadata manager class name for the table.",
				Optional:    true,
			},
			"metadata_manager_configs": schema.MapAttribute{
				Description: "The metadata manager configs for the table.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"enable_preload": schema.BoolAttribute{
				Description: "The enable preload for the table.",
				Optional:    true,
			},
			"metadata_ttl": schema.Int64Attribute{
				Description: "The metadata ttl for the table.",
				Optional:    true,
			},
			"deleted_keys_ttl": schema.Int64Attribute{
				Description: "The deleted keys ttl for the table.",
				Optional:    true,
			},
			"delete_record_column": schema.StringAttribute{
				Description: "The deleted record column for the table.",
				Optional:    true,
			},
			"out_of_order_record_column": schema.StringAttribute{
				Description: "The out of order record column for the table.",
				Optional:    true,
			},
			"drop_out_of_order_record": schema.BoolAttribute{
				Description: "The drop out of order record for the table.",
				Optional:    true,
			},
			"default_partial_upsert_strategy": schema.StringAttribute{
				Description: "The default partial upsert strategy for the table.",
				Optional:    true,
			},
			"hash_function": schema.StringAttribute{
				Description: "The hash function for the table.",
				Optional:    true,
			},
			"enable_snapshot": schema.BoolAttribute{
				Description: "The enable snapshot for the table.",
				Optional:    true,
			},
		},
	}
}

func IngestionConfig() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
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
					"column_major_segment_builder_enabled": schema.BoolAttribute{
						Optional: true,
					},
					"track_filtered_message_offsets": schema.BoolAttribute{
						Optional: true,
					},
					"stream_config_maps": schema.ListNestedAttribute{
						Description: "stream configuration",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"access_key": schema.StringAttribute{
									Optional: true,
								},
								"authentication_type": schema.StringAttribute{
									Optional: true,
								},
								"key_serializer": schema.StringAttribute{
									Optional: true,
								},
								"max_records_to_fetch": schema.StringAttribute{
									Optional: true,
								},
								"realtime_segment_commit_timeout_seconds": schema.StringAttribute{
									Optional: true,
								},
								"realtime_segment_flush_autotune_initial_rows": schema.StringAttribute{
									Optional: true,
								},
								"realtime_segment_flush_desired_size": schema.StringAttribute{
									Optional: true,
								},
								"realtime_segment_flush_threshold_rows": schema.StringAttribute{
									Optional: true,
								},
								"realtime_segment_flush_threshold_time": schema.StringAttribute{
									Optional: true,
								},
								"realtime_segment_flush_threshold_segment_rows": schema.StringAttribute{
									Optional: true,
								},
								"realtime_segment_flush_threshold_segment_size": schema.StringAttribute{
									Optional: true,
								},
								"realtime_segment_flush_threshold_segment_time": schema.StringAttribute{
									Optional: true,
								},
								"realtime_segment_server_upload_to_deep_store": schema.StringAttribute{
									Optional: true,
								},
								"region": schema.StringAttribute{
									Optional: true,
								},
								"sasl_jaas_config": schema.StringAttribute{
									Optional: true,
								},
								"sasl_mechanism": schema.StringAttribute{
									Optional: true,
								},
								"secret_key": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"security_protocol": schema.StringAttribute{
									Optional: true,
								},
								"shard_iterator_type": schema.StringAttribute{
									Optional: true,
								},
								"stream_type": schema.StringAttribute{
									Optional: true,
								},
								"ssl_key_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"ssl_keystore_location": schema.StringAttribute{
									Optional: true,
								},
								"ssl_keystore_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"ssl_keystore_type": schema.StringAttribute{
									Optional: true,
								},
								"ssl_truststore_location": schema.StringAttribute{
									Optional: true,
								},
								"ssl_truststore_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"ssl_truststore_type": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_broker_list": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_buffer_size": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_consumer_factory_class_name": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_consumer_prop_auto_offset_reset": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_consumer_type": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_fetch_timeout_millis": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_connection_timeout_millis": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_class_name": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_prop_basic_auth_credentials_source": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_prop_descriptor_file": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_prop_proto_class_name": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_prop_format": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_prop_schema_registry_basic_auth_user_info": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_prop_schema_registry_basic_auth_credentials_source": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_prop_schema_registry_rest_url": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_prop_schema_registry_schema_name": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_prop_schema_registry_ssl_keystore_location": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_prop_schema_registry_ssl_keystore_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"stream_kafka_decoder_prop_schema_registry_ssl_keystore_type": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_prop_schema_registry_ssl_truststore_location": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_prop_schema_registry_ssl_truststore_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"stream_kafka_decoder_prop_schema_registry_ssl_truststore_type": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_decoder_prop_schema_registry_ssl_protocol": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_fetcher_min_bytes": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_fetcher_size": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_hlc_group_id": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_idle_timeout_millis": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_isolation_level": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_metadata_populate": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_schema_registry_url": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_socket_timeout": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_ssl_certificate_type": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_ssl_client_certificate": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_ssl_client_key": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_ssl_client_key_algorithm": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_ssl_server_certificate": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_topic_name": schema.StringAttribute{
									Optional: true,
								},
								"stream_kafka_zk_broker_url": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_consumer_factory_class_name": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_consumer_type": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_decoder_class_name": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_fetch_timeout_millis": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_decoder_prop_schema_registry_basic_auth_user_info": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_decoder_prop_schema_registry_basic_auth_credentials_source": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_decoder_prop_schema_registry_rest_url": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_decoder_prop_schema_registry_schema_name": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_decoder_prop_schema_registry_ssl_keystore_location": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_decoder_prop_schema_registry_ssl_keystore_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"stream_kinesis_decoder_prop_schema_registry_ssl_keystore_type": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_decoder_prop_schema_registry_ssl_truststore_location": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_decoder_prop_schema_registry_ssl_truststore_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"stream_kinesis_decoder_prop_schema_registry_ssl_truststore_type": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_decoder_prop_schema_registry_ssl_protocol": schema.StringAttribute{
									Optional: true,
								},
								"stream_kinesis_topic_name": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_audience": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_authentication_token": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"stream_pulsar_bootstrap_servers": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_creds_file_path": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_consumer_factory_class_name": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_consumer_prop_auto_offset_reset": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_consumer_type": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_decoder_class_name": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_fetch_timeout_millis": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_issuer_url": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_metadata_populate": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_metadata_fields": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_decoder_prop_schema_registry_basic_auth_user_info": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_decoder_prop_schema_registry_basic_auth_credentials_source": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_decoder_prop_schema_registry_rest_url": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_decoder_prop_schema_registry_schema_name": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_decoder_prop_schema_registry_ssl_keystore_location": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_decoder_prop_schema_registry_ssl_keystore_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"stream_pulsar_decoder_prop_schema_registry_ssl_keystore_type": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_decoder_prop_schema_registry_ssl_truststore_location": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_decoder_prop_schema_registry_ssl_truststore_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"stream_pulsar_decoder_prop_schema_registry_ssl_truststore_type": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_decoder_prop_schema_registry_ssl_protocol": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_tls_trust_certs_file_path": schema.StringAttribute{
									Optional: true,
								},
								"stream_pulsar_topic_name": schema.StringAttribute{
									Optional: true,
								},
								"topic_consumption_rate_limit": schema.StringAttribute{
									Optional: true,
								},
								"value_serializer": schema.StringAttribute{
									Optional: true,
								},
							},
						},
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
			"filter_config": schema.SingleNestedAttribute{
				Description: "filter configuration",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"filter_function": schema.StringAttribute{
						Description: "filter function",
						Optional:    true,
					},
				},
			},
		},
	}
}

func TierConfigs() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
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
	}
}

func FieldConfigList() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
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
	}
}

func Routing() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "routing configuration for the table",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"segment_pruner_types": schema.ListAttribute{
				Description: "segment pruner types",
				Optional:    true,
				ElementType: types.StringType,
			},
			"instance_selector_type": schema.StringAttribute{
				Description: "instance selector type",
				Optional:    true,
			},
		},
	}
}

func Metadata() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "metadata for the table",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"custom_configs": schema.MapAttribute{
				Description: "custom configs",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func StarTreeIndexConfigs() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
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
	}
}

func TaskConfig() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"task_type_configs_map": schema.MapAttribute{
				Description: "The task configs for table",
				Optional:    true,
				ElementType: types.MapType{
					ElemType: types.StringType,
				},
			},
		},
	}
}
