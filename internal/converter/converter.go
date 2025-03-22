package converter

import (
	"context"
	"terraform-provider-pinot/internal/models"

	pinot_api "github.com/azaurus1/go-pinot-api/model"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SetStateFromTable(ctx context.Context, state *models.TableResourceModel, table *pinot_api.Table) diag.Diagnostics {

	var diags diag.Diagnostics

	state.TableName = types.StringValue(table.TableName)
	state.TableType = types.StringValue(table.TableType)
	state.IsDimTable = types.BoolValue(table.IsDimTable)

	state.TenantsConfig = convertTenantConfig(table)
	state.SegmentsConfig = convertSegmentsConfig(table)

	if table.IngestionConfig != nil {
		state.IngestionConfig = convertIngestionConfig(table)
	}

	if table.TierConfigs != nil {
		state.TierConfigs = convertTierConfigs(table)
	}

	if table.Metadata != nil {
		state.Metadata = convertMetadata(table)
	}

	tableIndexConfig, resultDiags := convertTableIndexConfig(ctx, table)
	if resultDiags.HasError() {
		diags.Append(resultDiags...)
	}

	state.TableIndexConfig = tableIndexConfig

	// Routing Config
	if table.Routing != nil {
		routingConfig, routingDiags := convertRoutingConfig(ctx, table)
		if routingDiags.HasError() {
			diags.Append(routingDiags...)
		}
		state.Routing = routingConfig
	}

	// Upsert Config
	if table.UpsertConfig != nil {
		upsertConfig, upsertDiags := convertUpsertConfig(ctx, table)
		if upsertDiags.HasError() {
			diags.Append(upsertDiags...)
		}
		state.UpsertConfig = upsertConfig
	}

	// Field Config List
	if table.FieldConfigList != nil {
		state.FieldConfigList = convertFieldConfigList(table)
	}

	if table.Task != nil {
		state.Task = convertTaskConfig(table)
	}

	return diags
}

func convertMetadata(table *pinot_api.Table) *models.Metadata {
	return &models.Metadata{
		CustomConfigs: table.Metadata.CustomConfigs,
	}
}

func convertTierConfigs(table *pinot_api.Table) []*models.TierConfig {
	var tierConfigs []*models.TierConfig
	for _, tierConfig := range table.TierConfigs {
		tierConfigs = append(tierConfigs, &models.TierConfig{
			Name:                types.StringValue(tierConfig.Name),
			StorageType:         types.StringValue(tierConfig.StorageType),
			SegmentSelectorType: types.StringValue(tierConfig.SegmentSelectorType),
			SegmentAge:          types.StringValue(tierConfig.SegmentAge),
			ServerTag:           types.StringValue(tierConfig.ServerTag),
		})
	}

	return tierConfigs
}

func convertTenantConfig(table *pinot_api.Table) *models.TenantsConfig {
	return &models.TenantsConfig{
		Broker: types.StringValue(table.Tenants.Broker),
		Server: types.StringValue(table.Tenants.Server),
	}
}

func convertIngestionConfig(table *pinot_api.Table) *models.IngestionConfig {

	var transformConfigs []*models.TransformConfig

	if table.IngestionConfig.TransformConfigs != nil {
		for _, transformConfig := range table.IngestionConfig.TransformConfigs {
			transformConfigs = append(transformConfigs, &models.TransformConfig{
				ColumnName:        types.StringValue(transformConfig.ColumnName),
				TransformFunction: types.StringValue(transformConfig.TransformFunction),
			})
		}
	}

	var filterConfig *models.FilterConfig
	if table.IngestionConfig.FilterConfig != nil {
		filterConfig = &models.FilterConfig{
			FilterFunction: types.StringValue(table.IngestionConfig.FilterConfig.FilterFunction),
		}
	}

	var streamConfigs []*models.StreamConfig

	if table.IngestionConfig.StreamIngestionConfig.StreamConfigMaps != nil {
		for _, streamConfig := range table.IngestionConfig.StreamIngestionConfig.StreamConfigMaps {
			streamConfigs = append(streamConfigs, &models.StreamConfig{
				AuthenticationType:                                               types.StringValue(streamConfig.AuthenticationType),
				StreamType:                                                       types.StringValue(streamConfig.StreamType),
				SslKeyPassword:                                                   types.StringValue(streamConfig.SslKeyPassword),
				RealtimeSegmentFlushThresholdRows:                                types.StringValue(streamConfig.RealtimeSegmentFlushThresholdRows),
				RealtimeSegmentFlushThresholdTime:                                types.StringValue(streamConfig.RealtimeSegmentFlushThresholdTime),
				RealtimeSegmentFlushThresholdSegmentRows:                         types.StringValue(streamConfig.RealtimeSegmentFlushThresholdSegmentRows),
				RealtimeSegmentFlushThresholdSegmentTime:                         types.StringValue(streamConfig.RealtimeSegmentFlushThresholdSegmentTime),
				RealtimeSegmentServerUploadToDeepStore:                           types.StringValue(streamConfig.RealtimeSegmentServerUploadToDeepStore),
				SecurityProtocol:                                                 types.StringValue(streamConfig.SecurityProtocol),
				SslKeystoreLocation:                                              types.StringValue(streamConfig.SslKeystoreLocation),
				SslKeystorePassword:                                              types.StringValue(streamConfig.SslKeystorePassword),
				SslTruststoreLocation:                                            types.StringValue(streamConfig.SslTruststoreLocation),
				SslTruststorePassword:                                            types.StringValue(streamConfig.SslTruststorePassword),
				StreamKafkaBrokerList:                                            types.StringValue(streamConfig.StreamKafkaBrokerList),
				StreamKafkaConsumerFactoryClassName:                              types.StringValue(streamConfig.StreamKafkaConsumerFactoryClassName),
				StreamKafkaConsumerPropAutoOffsetReset:                           types.StringValue(streamConfig.StreamKafkaConsumerPropAutoOffsetReset),
				StreamKafkaConsumerType:                                          types.StringValue(streamConfig.StreamKafkaConsumerType),
				StreamKafkaDecoderClassName:                                      types.StringValue(streamConfig.StreamKafkaDecoderClassName),
				StreamKafkaDecoderPropDescriptorFile:                             types.StringValue(streamConfig.StreamKafkaDecoderPropDescriptorFile),
				StreamKafkaDecoderPropProtoClassName:                             types.StringValue(streamConfig.StreamKafkaDecoderPropProtoClassName),
				StreamKafkaDecoderPropSchemaRegistryRestUrl:                      types.StringValue(streamConfig.StreamKafkaDecoderPropSchemaRegistryRestUrl),
				StreamKafkaDecoderPropSchemaRegistrySchemaName:                   types.StringValue(streamConfig.StreamKafkaDecoderPropSchemaRegistrySchemaName),
				StreamKafkaDecoderPropSchemaRegistrySslTruststoreLocation:        types.StringValue(streamConfig.StreamKafkaDecoderPropSchemaRegistrySslTruststoreLocation),
				StreamKafkaDecoderPropSchemaRegistrySslTruststorePassword:        types.StringValue(streamConfig.StreamKafkaDecoderPropSchemaRegistrySslTruststorePassword),
				StreamKafkaDecoderPropSchemaRegistrySslTruststoreType:            types.StringValue(streamConfig.StreamKafkaDecoderPropSchemaRegistrySslTruststoreType),
				StreamKafkaDecoderPropSchemaRegistrySslProtocol:                  types.StringValue(streamConfig.StreamKafkaDecoderPropSchemaRegistrySslProtocol),
				StreamKafkaMetadataPopulate:                                      types.StringValue(streamConfig.StreamKafkaMetadataPopulate),
				StreamKafkaTopicName:                                             types.StringValue(streamConfig.StreamKafkaTopicName),
				AccessKey:                                                        types.StringValue(streamConfig.AccessKey),
				KeySerializer:                                                    types.StringValue(streamConfig.KeySerializer),
				MaxRecordsToFetch:                                                types.StringValue(streamConfig.MaxRecordsToFetch),
				RealtimeSegmentCommitTimeoutSeconds:                              types.StringValue(streamConfig.RealtimeSegmentCommitTimeoutSeconds),
				RealtimeSegmentFlushAutotuneInitialRows:                          types.StringValue(streamConfig.RealtimeSegmentFlushAutotuneInitialRows),
				RealtimeSegmentFlushDesiredSize:                                  types.StringValue(streamConfig.RealtimeSegmentFlushDesiredSize),
				RealtimeSegmentFlushThresholdSegmentSize:                         types.StringValue(streamConfig.RealtimeSegmentFlushThresholdSegmentSize),
				Region:                                                           types.StringValue(streamConfig.Region),
				SaslJaasConfig:                                                   types.StringValue(streamConfig.SaslJaasConfig),
				SaslMechanism:                                                    types.StringValue(streamConfig.SaslMechanism),
				SecretKey:                                                        types.StringValue(streamConfig.SecretKey),
				ShardIteratorType:                                                types.StringValue(streamConfig.ShardIteratorType),
				SslKeystoreType:                                                  types.StringValue(streamConfig.SslKeystoreType),
				SslTruststoreType:                                                types.StringValue(streamConfig.SslTruststoreType),
				StreamKafkaBufferSize:                                            types.StringValue(streamConfig.StreamKafkaBufferSize),
				StreamKafkaFetchTimeoutMillis:                                    types.StringValue(streamConfig.StreamKafkaFetchTimeoutMillis),
				StreamKafkaConnectionTimeoutMillis:                               types.StringValue(streamConfig.StreamKafkaConnectionTimeoutMillis),
				StreamKafkaDecoderPropBasicAuthCredentialsSource:                 types.StringValue(streamConfig.StreamKafkaDecoderPropBasicAuthCredentialsSource),
				StreamKafkaDecoderPropFormat:                                     types.StringValue(streamConfig.StreamKafkaDecoderPropFormat),
				StreamKafkaDecoderPropSchemaRegistryBasicAuthUserInfo:            types.StringValue(streamConfig.StreamKafkaDecoderPropSchemaRegistryBasicAuthUserInfo),
				StreamKafkaDecoderPropSchemaRegistryBasicAuthCredentialsSource:   types.StringValue(streamConfig.StreamKafkaDecoderPropSchemaRegistryBasicAuthCredentialsSource),
				StreamKafkaFetcherMinBytes:                                       types.StringValue(streamConfig.StreamKafkaFetcherMinBytes),
				StreamKafkaFetcherSize:                                           types.StringValue(streamConfig.StreamKafkaFetcherSize),
				StreamKafkaHlcGroupId:                                            types.StringValue(streamConfig.StreamKafkaHlcGroupId),
				StreamKafkaIdleTimeoutMillis:                                     types.StringValue(streamConfig.StreamKafkaIdleTimeoutMillis),
				StreamKafkaIsolationLevel:                                        types.StringValue(streamConfig.StreamKafkaIsolationLevel),
				StreamKafkaSchemaRegistryUrl:                                     types.StringValue(streamConfig.StreamKafkaSchemaRegistryUrl),
				StreamKafkaSocketTimeout:                                         types.StringValue(streamConfig.StreamKafkaSocketTimeout),
				StreamKafkaSslCertificateType:                                    types.StringValue(streamConfig.StreamKafkaSslCertificateType),
				StreamKafkaSslClientCertificate:                                  types.StringValue(streamConfig.StreamKafkaSslClientCertificate),
				StreamKafkaSslClientKey:                                          types.StringValue(streamConfig.StreamKafkaSslClientKey),
				StreamKafkaSslClientKeyAlgorithm:                                 types.StringValue(streamConfig.StreamKafkaSslClientKeyAlgorithm),
				StreamKafkaSslServerCertificate:                                  types.StringValue(streamConfig.StreamKafkaSslServerCertificate),
				StreamKinesisConsumerFactoryClassName:                            types.StringValue(streamConfig.StreamKinesisConsumerFactoryClassName),
				StreamKinesisConsumerType:                                        types.StringValue(streamConfig.StreamKinesisConsumerType),
				StreamKinesisDecoderClassName:                                    types.StringValue(streamConfig.StreamKinesisDecoderClassName),
				StreamKinesisFetchTimeoutMillis:                                  types.StringValue(streamConfig.StreamKinesisFetchTimeoutMillis),
				StreamKinesisTopicName:                                           types.StringValue(streamConfig.StreamKinesisTopicName),
				TopicConsumptionRateLimit:                                        types.StringValue(streamConfig.TopicConsumptionRateLimit),
				ValueSerializer:                                                  types.StringValue(streamConfig.ValueSerializer),
				StreamKafkaDecoderPropSchemaRegistrySslKeystoreLocation:          types.StringValue(streamConfig.StreamKafkaDecoderPropSchemaRegistrySslKeystoreLocation),
				StreamKafkaDecoderPropSchemaRegistrySslKeystorePassword:          types.StringValue(streamConfig.StreamKafkaDecoderPropSchemaRegistrySslKeystorePassword),
				StreamKafkaDecoderPropSchemaRegistrySslKeystoreType:              types.StringValue(streamConfig.StreamKafkaDecoderPropSchemaRegistrySslKeystoreType),
				StreamKafkaZkBrokerUrl:                                           types.StringValue(streamConfig.StreamKafkaZkBrokerUrl),
				StreamKinesisDecoderPropSchemaRegistryBasicAuthUserInfo:          types.StringValue(streamConfig.StreamKinesisDecoderPropSchemaRegistryBasicAuthUserInfo),
				StreamKinesisDecoderPropSchemaRegistryBasicAuthCredentialsSource: types.StringValue(streamConfig.StreamKinesisDecoderPropSchemaRegistryBasicAuthCredentialsSource),
				StreamKinesisDecoderPropSchemaRegistryRestUrl:                    types.StringValue(streamConfig.StreamKinesisDecoderPropSchemaRegistryRestUrl),
				StreamKinesisDecoderPropSchemaRegistrySchemaName:                 types.StringValue(streamConfig.StreamKinesisDecoderPropSchemaRegistrySchemaName),
				StreamKinesisDecoderPropSchemaRegistrySslKeystoreLocation:        types.StringValue(streamConfig.StreamKinesisDecoderPropSchemaRegistrySslKeystoreLocation),
				StreamKinesisDecoderPropSchemaRegistrySslKeystorePassword:        types.StringValue(streamConfig.StreamKinesisDecoderPropSchemaRegistrySslKeystorePassword),
				StreamKinesisDecoderPropSchemaRegistrySslKeystoreType:            types.StringValue(streamConfig.StreamKinesisDecoderPropSchemaRegistrySslKeystoreType),
				StreamKinesisDecoderPropSchemaRegistrySslTruststoreLocation:      types.StringValue(streamConfig.StreamKinesisDecoderPropSchemaRegistrySslTruststoreLocation),
				StreamKinesisDecoderPropSchemaRegistrySslTruststorePassword:      types.StringValue(streamConfig.StreamKinesisDecoderPropSchemaRegistrySslTruststorePassword),
				StreamKinesisDecoderPropSchemaRegistrySslTruststoreType:          types.StringValue(streamConfig.StreamKinesisDecoderPropSchemaRegistrySslTruststoreType),
				StreamKinesisDecoderPropSchemaRegistrySslProtocol:                types.StringValue(streamConfig.StreamKinesisDecoderPropSchemaRegistrySslProtocol),
				StreamPulsarAudience:                                             types.StringValue(streamConfig.StreamPulsarAudience),
				StreamPulsarAuthenticationToken:                                  types.StringValue(streamConfig.StreamPulsarAuthenticationToken),
				StreamPulsarBootstrapServers:                                     types.StringValue(streamConfig.StreamPulsarBootstrapServers),
				StreamPulsarCredsFilePath:                                        types.StringValue(streamConfig.StreamPulsarCredsFilePath),
				StreamPulsarConsumerFactoryClassName:                             types.StringValue(streamConfig.StreamPulsarConsumerFactoryClassName),
				StreamPulsarConsumerPropAutoOffsetReset:                          types.StringValue(streamConfig.StreamPulsarConsumerPropAutoOffsetReset),
				StreamPulsarConsumerType:                                         types.StringValue(streamConfig.StreamPulsarConsumerType),
				StreamPulsarDecoderClassName:                                     types.StringValue(streamConfig.StreamPulsarDecoderClassName),
				StreamPulsarFetchTimeoutMillis:                                   types.StringValue(streamConfig.StreamPulsarFetchTimeoutMillis),
				StreamPulsarIssuerUrl:                                            types.StringValue(streamConfig.StreamPulsarIssuerUrl),
				StreamPulsarMetadataPopulate:                                     types.StringValue(streamConfig.StreamPulsarMetadataPopulate),
				StreamPulsarMetadataFields:                                       types.StringValue(streamConfig.StreamPulsarMetadataFields),
				StreamPulsarDecoderPropSchemaRegistryBasicAuthUserInfo:           types.StringValue(streamConfig.StreamPulsarDecoderPropSchemaRegistryBasicAuthUserInfo),
				StreamPulsarDecoderPropSchemaRegistryBasicAuthCredentialsSource:  types.StringValue(streamConfig.StreamPulsarDecoderPropSchemaRegistryBasicAuthCredentialsSource),
				StreamPulsarDecoderPropSchemaRegistryRestUrl:                     types.StringValue(streamConfig.StreamPulsarDecoderPropSchemaRegistryRestUrl),
				StreamPulsarDecoderPropSchemaRegistrySchemaName:                  types.StringValue(streamConfig.StreamPulsarDecoderPropSchemaRegistrySchemaName),
				StreamPulsarDecoderPropSchemaRegistrySslKeystoreLocation:         types.StringValue(streamConfig.StreamPulsarDecoderPropSchemaRegistrySslKeystoreLocation),
				StreamPulsarDecoderPropSchemaRegistrySslKeystorePassword:         types.StringValue(streamConfig.StreamPulsarDecoderPropSchemaRegistrySslKeystorePassword),
				StreamPulsarDecoderPropSchemaRegistrySslKeystoreType:             types.StringValue(streamConfig.StreamPulsarDecoderPropSchemaRegistrySslKeystoreType),
				StreamPulsarDecoderPropSchemaRegistrySslTruststoreLocation:       types.StringValue(streamConfig.StreamPulsarDecoderPropSchemaRegistrySslTruststoreLocation),
				StreamPulsarDecoderPropSchemaRegistrySslTruststorePassword:       types.StringValue(streamConfig.StreamPulsarDecoderPropSchemaRegistrySslTruststorePassword),
				StreamPulsarDecoderPropSchemaRegistrySslTruststoreType:           types.StringValue(streamConfig.StreamPulsarDecoderPropSchemaRegistrySslTruststoreType),
				StreamPulsarDecoderPropSchemaRegistrySslProtocol:                 types.StringValue(streamConfig.StreamPulsarDecoderPropSchemaRegistrySslProtocol),
				StreamPulsarTlsTrustCertsFilePath:                                types.StringValue(streamConfig.StreamPulsarTlsTrustCertsFilePath),
				StreamPulsarTopicName:                                            types.StringValue(streamConfig.StreamPulsarTopicName),
			})
		}
	}

	return &models.IngestionConfig{
		SegmentTimeValueCheck: types.BoolPointerValue(table.IngestionConfig.SegmentTimeValueCheck),
		RowTimeValueCheck:     types.BoolValue(table.IngestionConfig.RowTimeValueCheck),
		ContinueOnError:       types.BoolPointerValue(table.IngestionConfig.ContinueOnError),
		StreamIngestionConfig: &models.StreamIngestionConfig{
			ColumnMajorSegmentBuilderEnabled: types.BoolPointerValue(table.IngestionConfig.StreamIngestionConfig.ColumnMajorSegmentBuilderEnabled),
			TrackFilteredMessageOffsets:      types.BoolPointerValue(table.IngestionConfig.StreamIngestionConfig.TrackFilteredMessageOffsets),
			StreamConfigMaps:                 streamConfigs,
		},
		TransformConfigs: transformConfigs,
		FilterConfig:     filterConfig,
	}
}

func convertSegmentPartitionConfig(table *pinot_api.Table) *models.SegmentPartitionConfig {

	segmentPartitionConfig := &models.SegmentPartitionConfig{}

	if table.TableIndexConfig.SegmentPartitionConfig != nil {

		segmentPartitionMapConfig := map[string]map[string]string{}
		for key, value := range table.TableIndexConfig.SegmentPartitionConfig.ColumnPartitionMap {
			segmentPartitionMapConfig[key] = map[string]string{
				"functionName":  value.FunctionName,
				"numPartitions": string(rune(value.NumPartitions)),
			}
		}

		segmentPartitionConfig.ColumnPartitionMap = segmentPartitionMapConfig

	} else {
		return nil
	}

	return segmentPartitionConfig
}

func convertTableIndexConfig(ctx context.Context, table *pinot_api.Table) (*models.TableIndexConfig, diag.Diagnostics) {

	var diags diag.Diagnostics

	noDictionaryColumns, resultDiags := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.NoDictionaryColumns)
	if resultDiags.HasError() {
		diags.Append(resultDiags...)
	}

	onHeapDictionaryColumns, resultDiags := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.OnHeapDictionaryColumns)
	if resultDiags.HasError() {
		diags.Append(resultDiags...)
	}

	varLengthDictionaryColumns, resultDiags := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.VarLengthDictionaryColumns)
	if resultDiags.HasError() {
		diags.Append(resultDiags...)
	}

	rangeIndexColumns, resultDiags := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.RangeIndexColumns)
	if resultDiags.HasError() {
		diags.Append(resultDiags...)
	}

	bloomFilterColumns, resultDiags := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.BloomFilterColumns)
	if resultDiags.HasError() {
		diags.Append(resultDiags...)
	}

	sortedColumn, resultDiags := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.SortedColumn)
	if resultDiags.HasError() {
		diags.Append(resultDiags...)
	}

	invertedIndexColumns, resultDiags := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.InvertedIndexColumns)
	if resultDiags.HasError() {
		diags.Append(resultDiags...)
	}

	jsonIndexColumns, resultDiags := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.JsonIndexColumns)
	if resultDiags.HasError() {
		diags.Append(resultDiags...)
	}

	indexConfig := models.TableIndexConfig{
		LoadMode:                 types.StringValue(table.TableIndexConfig.LoadMode),
		NullHandlingEnabled:      types.BoolValue(table.TableIndexConfig.NullHandlingEnabled),
		SegmentNameGeneratorType: types.StringValue(table.TableIndexConfig.SegmentNameGeneratorType),
		CreateInvertedIndexDuringSegmentGeneration: types.BoolValue(table.TableIndexConfig.CreateInvertedIndexDuringSegmentGeneration),
		EnableDefaultStarTree:                      types.BoolValue(table.TableIndexConfig.EnableDefaultStarTree),
		OptimizeDictionary:                         types.BoolValue(table.TableIndexConfig.OptimizeDictionary),
		OptimizeDictionaryForMetrics:               types.BoolValue(table.TableIndexConfig.OptimizeDictionaryForMetrics),
		AutoGeneratedInvertedIndex:                 types.BoolValue(table.TableIndexConfig.AutoGeneratedInvertedIndex),
		NoDictionarySizeRatioThreshold:             types.Float64Value(table.TableIndexConfig.NoDictionarySizeRatioThreshold),
		StarTreeIndexConfigs:                       convertStarTreeIndexConfigs(ctx, table),
		AggregateMetrics:                           types.BoolValue(table.TableIndexConfig.AggregateMetrics),
		SegmentPartitionConfig:                     convertSegmentPartitionConfig(table),
		RangeIndexVersion:                          types.Int64Value(int64(table.TableIndexConfig.RangeIndexVersion)),
		SortedColumn:                               sortedColumn,
		OnHeapDictionaryColumns:                    onHeapDictionaryColumns,
		VarLengthDictionaryColumns:                 varLengthDictionaryColumns,
		RangeIndexColumns:                          rangeIndexColumns,
		NoDictionaryColumns:                        noDictionaryColumns,
		BloomFilterColumns:                         bloomFilterColumns,
		InvertedIndexColumns:                       invertedIndexColumns,
		JsonIndexColumns:                           jsonIndexColumns,
	}

	return &indexConfig, diags
}

func convertStarTreeIndexConfigs(ctx context.Context, table *pinot_api.Table) []*models.StarTreeIndexConfigs {

	var starTreeIndexConfigs []*models.StarTreeIndexConfigs

	for _, starConfig := range table.TableIndexConfig.StarTreeIndexConfigs {

		dimensionSplitOrder, _ := types.ListValueFrom(ctx, types.StringType, starConfig.DimensionsSplitOrder)
		functionColumnPairs, _ := types.ListValueFrom(ctx, types.StringType, starConfig.FunctionColumnPairs)
		skipStarNodeCreationDimensionColumns, _ := types.ListValueFrom(ctx, types.StringType, starConfig.SkipStarNodeCreationForDimensions)

		starTreeIndexConfigs = append(starTreeIndexConfigs, &models.StarTreeIndexConfigs{
			MaxLeafRecords:                  types.Int64Value(int64(starConfig.MaxLeafRecords)),
			DimensionsSplitOrder:            dimensionSplitOrder,
			FunctionColumnPairs:             functionColumnPairs,
			SkipStarNodeCreationForDimNames: skipStarNodeCreationDimensionColumns,
		})

	}

	return starTreeIndexConfigs
}

func convertSegmentsConfig(table *pinot_api.Table) *models.SegmentsConfig {

	segmentsConfig := models.SegmentsConfig{
		TimeType:           types.StringValue(table.SegmentsConfig.TimeType),
		Replication:        types.StringValue(table.SegmentsConfig.Replication),
		TimeColumnName:     types.StringValue(table.SegmentsConfig.TimeColumnName),
		RetentionTimeUnit:  types.StringValue(table.SegmentsConfig.RetentionTimeUnit),
		RetentionTimeValue: types.StringValue(table.SegmentsConfig.RetentionTimeValue),
	}

	if table.SegmentsConfig.ReplicasPerPartition != "" {
		segmentsConfig.ReplicasPerPartition = types.StringValue(table.SegmentsConfig.ReplicasPerPartition)
	}

	if table.SegmentsConfig.SchemaName != "" {
		segmentsConfig.SchemaName = types.StringValue(table.SegmentsConfig.SchemaName)
	}

	if table.SegmentsConfig.DeletedSegmentsRetentionPeriod != "" {
		segmentsConfig.DeletedSegmentsRetentionPeriod = types.StringValue(table.SegmentsConfig.DeletedSegmentsRetentionPeriod)
	}

	if table.SegmentsConfig.PeerSegmentDownloadScheme != "" {
		segmentsConfig.PeerSegmentDownloadScheme = types.StringValue(table.SegmentsConfig.PeerSegmentDownloadScheme)
	}

	if table.SegmentsConfig.SegmentPushType != "" {
		segmentsConfig.SegmentPushType = types.StringValue(table.SegmentsConfig.SegmentPushType)
	}

	if table.SegmentsConfig.MinimizeDataMovement {
		segmentsConfig.MinimizeDataMovement = types.BoolValue(table.SegmentsConfig.MinimizeDataMovement)
	}

	if table.SegmentsConfig.CompletionConfig != nil && table.SegmentsConfig.CompletionConfig.CompletionMode != "" {
		segmentsConfig.CompletionConfig = &models.CompletionConfig{
			CompletionMode: types.StringValue(table.SegmentsConfig.CompletionConfig.CompletionMode),
		}
	}

	return &segmentsConfig
}

func convertRoutingConfig(ctx context.Context, table *pinot_api.Table) (*models.RoutingConfig, diag.Diagnostics) {

	segmentPrunerTypes, resultDiags := types.ListValueFrom(ctx, types.StringType, table.Routing.SegmentPrunerTypes)
	if resultDiags.HasError() {
		return nil, resultDiags
	}

	routingConfig := models.RoutingConfig{
		InstanceSelectorType: types.StringValue(table.Routing.InstanceSelectorType),
		SegmentPrunerTypes:   segmentPrunerTypes,
	}

	return &routingConfig, resultDiags
}

func convertUpsertConfig(ctx context.Context, table *pinot_api.Table) (*models.UpsertConfig, diag.Diagnostics) {

	partialUpsertStrategies, resultDiags := types.MapValueFrom(ctx, types.StringType, table.UpsertConfig.PartialUpsertStrategies)
	if resultDiags.HasError() {
		return nil, resultDiags
	}

	metadataManagerConfigs, resultDiags := types.MapValueFrom(ctx, types.StringType, table.UpsertConfig.MetadataManagerConfigs)
	if resultDiags.HasError() {
		return nil, resultDiags
	}

	upsertConfig := models.UpsertConfig{
		Mode:                         types.StringValue(table.UpsertConfig.Mode),
		PartialUpsertStrategy:        partialUpsertStrategies,
		MetadataTTL:                  types.Int64Value(int64(table.UpsertConfig.MetadataTTL)),
		DeletedKeysTTL:               types.Int64Value(int64(table.UpsertConfig.DeletedKeysTTL)),
		HashFunction:                 types.StringValue(table.UpsertConfig.HashFunction),
		EnableSnapshot:               types.BoolPointerValue(table.UpsertConfig.EnableSnapshot),
		EnablePreLoad:                types.BoolPointerValue(table.UpsertConfig.EnablePreload),
		DropOutOfOrderRecord:         types.BoolPointerValue(table.UpsertConfig.DropOutOfOrderRecord),
		DefaultPartialUpsertStrategy: types.StringValue(table.UpsertConfig.DefaultPartialUpsertStrategy),
		MetadataManagerConfigs:       metadataManagerConfigs,
		MetadataManagerClass:         types.StringValue(table.UpsertConfig.MetadataManagerClass),
	}

	if table.UpsertConfig.DeleteRecordColumn != "" {
		upsertConfig.DeletedRecordColumn = types.StringValue(table.UpsertConfig.DeleteRecordColumn)
	}

	if table.UpsertConfig.MetadataManagerClass != "" {
		upsertConfig.MetadataManagerClass = types.StringValue(table.UpsertConfig.MetadataManagerClass)
	}

	if table.UpsertConfig.OutOfOrderRecordColumn != "" {
		upsertConfig.OutOfOrderRecordColumn = types.StringValue(table.UpsertConfig.OutOfOrderRecordColumn)
	}

	return &upsertConfig, resultDiags
}

func convertTaskConfig(table *pinot_api.Table) *models.Task {

	if table.Task == nil {
		return nil
	}

	return &models.Task{
		TaskTypeConfigsMap: table.Task.TaskTypeConfigsMap, // Changed to match new field name
	}
}

func convertFieldConfigList(table *pinot_api.Table) []*models.FieldConfig {

	if table.FieldConfigList == nil {
		return nil
	}

	var fieldConfigs []*models.FieldConfig
	for _, fieldConfig := range table.FieldConfigList {

		fc := &models.FieldConfig{
			Name:         types.StringValue(fieldConfig.Name),
			EncodingType: types.StringValue(fieldConfig.EncodingType),
			IndexType:    types.StringValue(fieldConfig.IndexType),
			IndexTypes:   fieldConfig.IndexTypes,
		}

		if fieldConfig.TimestampConfig != nil {
			fc.TimestampConfig = &models.TimestampConfig{
				Granularities: fieldConfig.TimestampConfig.Granularities,
			}
		}

		if fieldConfig.Indexes != nil {
			fc.Indexes = &models.FieldIndexes{
				Inverted: &models.FiendIndexInverted{
					Enabled: types.StringValue(fieldConfig.Indexes.Inverted.Enabled),
				},
			}
		}

		fieldConfigs = append(fieldConfigs, fc)
	}
	return fieldConfigs
}
