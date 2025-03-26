package converter

import (
	"context"
	"log"
	"strconv"
	"terraform-provider-pinot/internal/models"

	"github.com/azaurus1/go-pinot-api/model"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ToTable(plan *models.TableResourceModel) (*model.Table, diag.Diagnostics) {

	var diags diag.Diagnostics

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	table := model.Table{
		TableName:       plan.TableName.ValueString(),
		TableType:       plan.TableType.ValueString(),
		IsDimTable:      plan.IsDimTable.ValueBool(),
		Tenants:         ToTenantsConfig(plan),
		SegmentsConfig:  ToSegmentsConfig(plan),
		IngestionConfig: ToIngestionConfig(plan),
	}

	tableIndexConfig, resultDiags := ToTableIndexConfig(ctx, plan)
	if resultDiags.HasError() {
		diags.Append(resultDiags...)
	}
	table.TableIndexConfig = *tableIndexConfig

	if plan.Metadata != nil {
		table.Metadata = ToMetadata(plan)
	}
	if plan.TierConfigs != nil {
		table.TierConfigs = ToTierConfigs(plan)
	}
	if plan.FieldConfigList != nil {
		table.FieldConfigList = ToFieldConfigList(plan)
	}

	if plan.UpsertConfig != nil {

		upsertConfig, upsertDiags := ToUpsertConfig(ctx, plan.UpsertConfig)
		if upsertDiags.HasError() {
			diags.Append(upsertDiags...)
		}

		table.UpsertConfig = upsertConfig
	}

	if plan.Routing != nil {

		routingConfig, routingDiags := ToRoutingConfig(ctx, plan.Routing)
		if routingDiags.HasError() {
			diags.Append(routingDiags...)
		}

		table.Routing = routingConfig
	}

	if plan.Task != nil {
		table.Task = ToTask(plan.Task)
	}

	if plan.InstanceAssignmentConfigMap != nil {
		table.InstanceAssignmentConfigMap = ToInstanceAssignmentConfigMap(plan.InstanceAssignmentConfigMap)
	}

	return &table, diags
}

func ToUpsertConfig(ctx context.Context, stateConfig *models.UpsertConfig) (*model.UpsertConfig, diag.Diagnostics) {

	var diags diag.Diagnostics

	partialUpsertStrategies, resultDiags := toMap(ctx, stateConfig.PartialUpsertStrategy)
	if resultDiags.HasError() {
		diags.Append(resultDiags...)
		return nil, diags
	}

	metadataManagerConfigs, resultDiags := toMap(ctx, stateConfig.MetadataManagerConfigs)
	if resultDiags.HasError() {
		diags.Append(resultDiags...)
		return nil, diags
	}

	upsertConfig := model.UpsertConfig{
		Mode:                         stateConfig.Mode.ValueString(),
		PartialUpsertStrategies:      partialUpsertStrategies,
		DeleteRecordColumn:           stateConfig.DeletedRecordColumn.ValueString(),
		MetadataTTL:                  float64(stateConfig.MetadataTTL.ValueInt64()),
		DeletedKeysTTL:               float64(stateConfig.DeletedKeysTTL.ValueInt64()),
		HashFunction:                 stateConfig.HashFunction.ValueString(),
		EnableSnapshot:               stateConfig.EnableSnapshot.ValueBoolPointer(),
		EnablePreload:                stateConfig.EnablePreLoad.ValueBoolPointer(),
		DropOutOfOrderRecord:         stateConfig.DropOutOfOrderRecord.ValueBoolPointer(),
		DefaultPartialUpsertStrategy: stateConfig.DefaultPartialUpsertStrategy.ValueString(),
		MetadataManagerConfigs:       metadataManagerConfigs,
		MetadataManagerClass:         stateConfig.MetadataManagerClass.ValueString(),
	}

	if stateConfig.OutOfOrderRecordColumn.IsNull() {
		upsertConfig.OutOfOrderRecordColumn = stateConfig.OutOfOrderRecordColumn.ValueString()
	}

	if stateConfig.MetadataManagerClass.IsNull() {
		upsertConfig.MetadataManagerClass = stateConfig.MetadataManagerClass.ValueString()
	}

	return &upsertConfig, diags

}

func ToRoutingConfig(ctx context.Context, stateConfig *models.RoutingConfig) (*model.RoutingConfig, diag.Diagnostics) {

	var diags diag.Diagnostics

	segmentPrunerTypes, resultDiags := toList[string](ctx, stateConfig.SegmentPrunerTypes)
	if resultDiags.HasError() {
		diags.Append(resultDiags...)
		return nil, diags
	}

	routingConfig := model.RoutingConfig{
		SegmentPrunerTypes:   segmentPrunerTypes,
		InstanceSelectorType: stateConfig.InstanceSelectorType.ValueString(),
	}

	return &routingConfig, diags

}

func ToFieldConfigList(plan *models.TableResourceModel) []model.FieldConfig {

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

func ToTableIndexConfig(ctx context.Context, plan *models.TableResourceModel) (*model.TableIndexConfig, diag.Diagnostics) {

	var diagnostics diag.Diagnostics

	sortedCols, resultDiags := toList[string](ctx, plan.TableIndexConfig.SortedColumn)
	if resultDiags.HasError() {
		diagnostics.Append(resultDiags...)
	}

	rangeIndexCols, resultDiags := toList[string](ctx, plan.TableIndexConfig.RangeIndexColumns)
	if resultDiags.HasError() {
		diagnostics.Append(resultDiags...)
	}

	noDictionaryColumns, resultDiags := toList[string](ctx, plan.TableIndexConfig.NoDictionaryColumns)
	if resultDiags.HasError() {
		diagnostics.Append(resultDiags...)
	}

	onHeapDictionaryCols, resultDiags := toList[string](ctx, plan.TableIndexConfig.OnHeapDictionaryColumns)
	if resultDiags.HasError() {
		diagnostics.Append(resultDiags...)
	}

	varLengthDictionaryCols, resultDiags := toList[string](ctx, plan.TableIndexConfig.VarLengthDictionaryColumns)
	if resultDiags.HasError() {
		diagnostics.Append(resultDiags...)
	}

	bloomFilterCols, resultDiags := toList[string](ctx, plan.TableIndexConfig.BloomFilterColumns)
	if resultDiags.HasError() {
		diagnostics.Append(resultDiags...)
	}

	invertedIndexColumns, resultDiags := toList[string](ctx, plan.TableIndexConfig.InvertedIndexColumns)
	if resultDiags.HasError() {
		diagnostics.Append(resultDiags...)
	}

	jsonIndexColumns, resultDiags := toList[string](ctx, plan.TableIndexConfig.JsonIndexColumns)
	if resultDiags.HasError() {
		diagnostics.Append(resultDiags...)
	}

	tableConfig := model.TableIndexConfig{
		CreateInvertedIndexDuringSegmentGeneration: plan.TableIndexConfig.CreateInvertedIndexDuringSegmentGeneration.ValueBool(),
		SortedColumn:                   sortedCols,
		RangeIndexColumns:              rangeIndexCols,
		NoDictionaryColumns:            noDictionaryColumns,
		OnHeapDictionaryColumns:        onHeapDictionaryCols,
		VarLengthDictionaryColumns:     varLengthDictionaryCols,
		BloomFilterColumns:             bloomFilterCols,
		InvertedIndexColumns:           invertedIndexColumns,
		JsonIndexColumns:               jsonIndexColumns,
		SegmentPartitionConfig:         ToSegmentPartitionConfig(plan),
		EnableDefaultStarTree:          plan.TableIndexConfig.EnableDefaultStarTree.ValueBool(),
		EnableDynamicStarTreeCreation:  plan.TableIndexConfig.EnableDynamicStarTree.ValueBool(),
		LoadMode:                       plan.TableIndexConfig.LoadMode.ValueString(),
		ColumnMinMaxValueGeneratorMode: plan.TableIndexConfig.ColumnMinMaxValueGeneratorMode.ValueString(),
		NullHandlingEnabled:            plan.TableIndexConfig.NullHandlingEnabled.ValueBool(),
		AggregateMetrics:               plan.TableIndexConfig.AggregateMetrics.ValueBool(),
		AutoGeneratedInvertedIndex:     plan.TableIndexConfig.AutoGeneratedInvertedIndex.ValueBool(),
		OptimizeDictionary:             plan.TableIndexConfig.OptimizeDictionary.ValueBool(),
		OptimizeDictionaryForMetrics:   plan.TableIndexConfig.OptimizeDictionaryForMetrics.ValueBool(),
		NoDictionarySizeRatioThreshold: plan.TableIndexConfig.NoDictionarySizeRatioThreshold.ValueFloat64(),
		SegmentNameGeneratorType:       plan.TableIndexConfig.SegmentNameGeneratorType.ValueString(),
		RangeIndexVersion:              int(plan.TableIndexConfig.RangeIndexVersion.ValueInt64()),
	}

	if plan.TableIndexConfig.StarTreeIndexConfigs != nil {

		starTreeIndexConfigs, starResultDiags := ToStarTreeIndexConfig(ctx, plan)
		if starResultDiags.HasError() {
			diagnostics.Append(starResultDiags...)
		}

		tableConfig.StarTreeIndexConfigs = starTreeIndexConfigs
	}

	return &tableConfig, diagnostics
}

func ToSegmentPartitionConfig(plan *models.TableResourceModel) *model.SegmentPartitionConfig {

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

func ToStarTreeIndexConfig(ctx context.Context, plan *models.TableResourceModel) ([]*model.StarTreeIndexConfig, diag.Diagnostics) {

	var diagnostics diag.Diagnostics

	var starTreeConfigs []*model.StarTreeIndexConfig
	for _, starConfig := range plan.TableIndexConfig.StarTreeIndexConfigs {

		dimensionsSplitOrder, resultDiags := toList[string](ctx, starConfig.DimensionsSplitOrder)
		if resultDiags.HasError() {
			diagnostics.Append(resultDiags...)
		}

		functionColumnPairs, resultDiags := toList[string](ctx, starConfig.FunctionColumnPairs)
		if resultDiags.HasError() {
			diagnostics.Append(resultDiags...)
		}

		skipStarNodeCreationForDimensions, resultDiags := toList[string](ctx, starConfig.SkipStarNodeCreationForDimNames)
		if resultDiags.HasError() {
			diagnostics.Append(resultDiags...)
		}

		starTreeConfigs = append(starTreeConfigs, &model.StarTreeIndexConfig{
			MaxLeafRecords:                    int(starConfig.MaxLeafRecords.ValueInt64()),
			DimensionsSplitOrder:              dimensionsSplitOrder,
			FunctionColumnPairs:               functionColumnPairs,
			SkipStarNodeCreationForDimensions: skipStarNodeCreationForDimensions,
		})
	}

	if diagnostics.HasError() {
		return nil, diagnostics
	}

	return starTreeConfigs, diagnostics
}

func ToSegmentsConfig(plan *models.TableResourceModel) model.TableSegmentsConfig {

	segmentsConfig := model.TableSegmentsConfig{
		TimeType:           plan.SegmentsConfig.TimeType.ValueString(),
		Replication:        plan.SegmentsConfig.Replication.ValueString(),
		TimeColumnName:     plan.SegmentsConfig.TimeColumnName.ValueString(),
		RetentionTimeUnit:  plan.SegmentsConfig.RetentionTimeUnit.ValueString(),
		RetentionTimeValue: plan.SegmentsConfig.RetentionTimeValue.ValueString(),
	}

	if plan.SegmentsConfig.ReplicasPerPartition.ValueString() != "" {
		segmentsConfig.ReplicasPerPartition = plan.SegmentsConfig.ReplicasPerPartition.ValueString()
	}

	if plan.SegmentsConfig.SchemaName.ValueString() != "" {
		segmentsConfig.SchemaName = plan.SegmentsConfig.SchemaName.ValueString()
	}

	if plan.SegmentsConfig.DeletedSegmentsRetentionPeriod.ValueString() != "" {
		segmentsConfig.DeletedSegmentsRetentionPeriod = plan.SegmentsConfig.DeletedSegmentsRetentionPeriod.ValueString()
	}

	if plan.SegmentsConfig.PeerSegmentDownloadScheme.ValueString() != "" {
		segmentsConfig.PeerSegmentDownloadScheme = plan.SegmentsConfig.PeerSegmentDownloadScheme.ValueString()
	}

	if plan.SegmentsConfig.SegmentPushType.ValueString() != "" {
		segmentsConfig.SegmentPushType = plan.SegmentsConfig.SegmentPushType.ValueString()
	}

	if plan.SegmentsConfig.MinimizeDataMovement.ValueBool() {
		segmentsConfig.MinimizeDataMovement = plan.SegmentsConfig.MinimizeDataMovement.ValueBool()
	}

	if plan.SegmentsConfig.CompletionConfig != nil && plan.SegmentsConfig.CompletionConfig.CompletionMode.ValueString() != "" {
		segmentsConfig.CompletionConfig = &model.CompletionConfig{
			CompletionMode: plan.SegmentsConfig.CompletionConfig.CompletionMode.ValueString(),
		}
	}

	return segmentsConfig

}

func ToMetadata(plan *models.TableResourceModel) *model.TableMetadata {
	if plan.Metadata == nil {
		return nil
	}

	return &model.TableMetadata{
		CustomConfigs: plan.Metadata.CustomConfigs,
	}
}

func ToIngestionConfig(plan *models.TableResourceModel) *model.TableIngestionConfig {

	if plan.IngestionConfig == nil {
		return nil
	}

	ingestionConfig := model.TableIngestionConfig{
		SegmentTimeValueCheck: plan.IngestionConfig.SegmentTimeValueCheck.ValueBoolPointer(),
		RowTimeValueCheck:     plan.IngestionConfig.RowTimeValueCheck.ValueBool(),
		ContinueOnError:       plan.IngestionConfig.ContinueOnError.ValueBoolPointer(),
	}

	if plan.IngestionConfig.TransformConfigs != nil {

		var transformConfigs []model.TransformConfig

		for _, transformConfig := range plan.IngestionConfig.TransformConfigs {
			transformConfigs = append(transformConfigs, model.TransformConfig{
				ColumnName:        transformConfig.ColumnName.ValueString(),
				TransformFunction: transformConfig.TransformFunction.ValueString(),
			})
		}

		ingestionConfig.TransformConfigs = transformConfigs
	}

	if plan.IngestionConfig.StreamIngestionConfig.StreamConfigMaps != nil {

		var streamConfigs []model.StreamConfig

		for _, streamConfig := range plan.IngestionConfig.StreamIngestionConfig.StreamConfigMaps {
			streamConfigs = append(streamConfigs, model.StreamConfig{
				AuthenticationType:                                               streamConfig.AuthenticationType.ValueString(),
				StreamType:                                                       streamConfig.StreamType.ValueString(),
				SslKeyPassword:                                                   streamConfig.SslKeyPassword.ValueString(),
				RealtimeSegmentFlushThresholdRows:                                streamConfig.RealtimeSegmentFlushThresholdRows.ValueString(),
				RealtimeSegmentFlushThresholdTime:                                streamConfig.RealtimeSegmentFlushThresholdTime.ValueString(),
				RealtimeSegmentFlushThresholdSegmentRows:                         streamConfig.RealtimeSegmentFlushThresholdSegmentRows.ValueString(),
				RealtimeSegmentFlushThresholdSegmentTime:                         streamConfig.RealtimeSegmentFlushThresholdSegmentTime.ValueString(),
				RealtimeSegmentServerUploadToDeepStore:                           streamConfig.RealtimeSegmentServerUploadToDeepStore.ValueString(),
				SecurityProtocol:                                                 streamConfig.SecurityProtocol.ValueString(),
				SslKeystoreLocation:                                              streamConfig.SslKeystoreLocation.ValueString(),
				SslKeystorePassword:                                              streamConfig.SslKeystorePassword.ValueString(),
				SslTruststoreLocation:                                            streamConfig.SslTruststoreLocation.ValueString(),
				SslTruststorePassword:                                            streamConfig.SslTruststorePassword.ValueString(),
				StreamKafkaBrokerList:                                            streamConfig.StreamKafkaBrokerList.ValueString(),
				StreamKafkaConsumerFactoryClassName:                              streamConfig.StreamKafkaConsumerFactoryClassName.ValueString(),
				StreamKafkaConsumerPropAutoOffsetReset:                           streamConfig.StreamKafkaConsumerPropAutoOffsetReset.ValueString(),
				StreamKafkaConsumerType:                                          streamConfig.StreamKafkaConsumerType.ValueString(),
				StreamKafkaDecoderClassName:                                      streamConfig.StreamKafkaDecoderClassName.ValueString(),
				StreamKafkaDecoderPropDescriptorFile:                             streamConfig.StreamKafkaDecoderPropDescriptorFile.ValueString(),
				StreamKafkaDecoderPropProtoClassName:                             streamConfig.StreamKafkaDecoderPropProtoClassName.ValueString(),
				StreamKafkaDecoderPropSchemaRegistryRestUrl:                      streamConfig.StreamKafkaDecoderPropSchemaRegistryRestUrl.ValueString(),
				StreamKafkaDecoderPropSchemaRegistrySchemaName:                   streamConfig.StreamKafkaDecoderPropSchemaRegistrySchemaName.ValueString(),
				StreamKafkaDecoderPropSchemaRegistrySslTruststoreLocation:        streamConfig.StreamKafkaDecoderPropSchemaRegistrySslTruststoreLocation.ValueString(),
				StreamKafkaDecoderPropSchemaRegistrySslTruststorePassword:        streamConfig.StreamKafkaDecoderPropSchemaRegistrySslTruststorePassword.ValueString(),
				StreamKafkaDecoderPropSchemaRegistrySslTruststoreType:            streamConfig.StreamKafkaDecoderPropSchemaRegistrySslTruststoreType.ValueString(),
				StreamKafkaDecoderPropSchemaRegistrySslProtocol:                  streamConfig.StreamKafkaDecoderPropSchemaRegistrySslProtocol.ValueString(),
				StreamKafkaMetadataPopulate:                                      streamConfig.StreamKafkaMetadataPopulate.ValueString(),
				StreamKafkaTopicName:                                             streamConfig.StreamKafkaTopicName.ValueString(),
				AccessKey:                                                        streamConfig.AccessKey.ValueString(),
				KeySerializer:                                                    streamConfig.KeySerializer.ValueString(),
				MaxRecordsToFetch:                                                streamConfig.MaxRecordsToFetch.ValueString(),
				RealtimeSegmentCommitTimeoutSeconds:                              streamConfig.RealtimeSegmentCommitTimeoutSeconds.ValueString(),
				RealtimeSegmentFlushAutotuneInitialRows:                          streamConfig.RealtimeSegmentFlushAutotuneInitialRows.ValueString(),
				RealtimeSegmentFlushDesiredSize:                                  streamConfig.RealtimeSegmentFlushDesiredSize.ValueString(),
				RealtimeSegmentFlushThresholdSegmentSize:                         streamConfig.RealtimeSegmentFlushThresholdSegmentSize.ValueString(),
				Region:                                                           streamConfig.Region.ValueString(),
				SaslJaasConfig:                                                   streamConfig.SaslJaasConfig.ValueString(),
				SaslMechanism:                                                    streamConfig.SaslMechanism.ValueString(),
				SecretKey:                                                        streamConfig.SecretKey.ValueString(),
				ShardIteratorType:                                                streamConfig.ShardIteratorType.ValueString(),
				SslKeystoreType:                                                  streamConfig.SslKeystoreType.ValueString(),
				SslTruststoreType:                                                streamConfig.SslTruststoreType.ValueString(),
				StreamKafkaBufferSize:                                            streamConfig.StreamKafkaBufferSize.ValueString(),
				StreamKafkaFetchTimeoutMillis:                                    streamConfig.StreamKafkaFetchTimeoutMillis.ValueString(),
				StreamKafkaConnectionTimeoutMillis:                               streamConfig.StreamKafkaConnectionTimeoutMillis.ValueString(),
				StreamKafkaDecoderPropBasicAuthCredentialsSource:                 streamConfig.StreamKafkaDecoderPropBasicAuthCredentialsSource.ValueString(),
				StreamKafkaDecoderPropFormat:                                     streamConfig.StreamKafkaDecoderPropFormat.ValueString(),
				StreamKafkaDecoderPropSchemaRegistryBasicAuthUserInfo:            streamConfig.StreamKafkaDecoderPropSchemaRegistryBasicAuthUserInfo.ValueString(),
				StreamKafkaDecoderPropSchemaRegistryBasicAuthCredentialsSource:   streamConfig.StreamKafkaDecoderPropSchemaRegistryBasicAuthCredentialsSource.ValueString(),
				StreamKafkaFetcherMinBytes:                                       streamConfig.StreamKafkaFetcherMinBytes.ValueString(),
				StreamKafkaFetcherSize:                                           streamConfig.StreamKafkaFetcherSize.ValueString(),
				StreamKafkaHlcGroupId:                                            streamConfig.StreamKafkaHlcGroupId.ValueString(),
				StreamKafkaIdleTimeoutMillis:                                     streamConfig.StreamKafkaIdleTimeoutMillis.ValueString(),
				StreamKafkaIsolationLevel:                                        streamConfig.StreamKafkaIsolationLevel.ValueString(),
				StreamKafkaSchemaRegistryUrl:                                     streamConfig.StreamKafkaSchemaRegistryUrl.ValueString(),
				StreamKafkaSocketTimeout:                                         streamConfig.StreamKafkaSocketTimeout.ValueString(),
				StreamKafkaSslCertificateType:                                    streamConfig.StreamKafkaSslCertificateType.ValueString(),
				StreamKafkaSslClientCertificate:                                  streamConfig.StreamKafkaSslClientCertificate.ValueString(),
				StreamKafkaSslClientKey:                                          streamConfig.StreamKafkaSslClientKey.ValueString(),
				StreamKafkaSslClientKeyAlgorithm:                                 streamConfig.StreamKafkaSslClientKeyAlgorithm.ValueString(),
				StreamKafkaSslServerCertificate:                                  streamConfig.StreamKafkaSslServerCertificate.ValueString(),
				StreamKinesisConsumerFactoryClassName:                            streamConfig.StreamKinesisConsumerFactoryClassName.ValueString(),
				StreamKinesisConsumerType:                                        streamConfig.StreamKinesisConsumerType.ValueString(),
				StreamKinesisDecoderClassName:                                    streamConfig.StreamKinesisDecoderClassName.ValueString(),
				StreamKinesisFetchTimeoutMillis:                                  streamConfig.StreamKinesisFetchTimeoutMillis.ValueString(),
				StreamKinesisTopicName:                                           streamConfig.StreamKinesisTopicName.ValueString(),
				TopicConsumptionRateLimit:                                        streamConfig.TopicConsumptionRateLimit.ValueString(),
				ValueSerializer:                                                  streamConfig.ValueSerializer.ValueString(),
				StreamKafkaDecoderPropSchemaRegistrySslKeystoreLocation:          streamConfig.StreamKafkaDecoderPropSchemaRegistrySslKeystoreLocation.ValueString(),
				StreamKafkaDecoderPropSchemaRegistrySslKeystorePassword:          streamConfig.StreamKafkaDecoderPropSchemaRegistrySslKeystorePassword.ValueString(),
				StreamKafkaDecoderPropSchemaRegistrySslKeystoreType:              streamConfig.StreamKafkaDecoderPropSchemaRegistrySslKeystoreType.ValueString(),
				StreamKafkaZkBrokerUrl:                                           streamConfig.StreamKafkaZkBrokerUrl.ValueString(),
				StreamKinesisDecoderPropSchemaRegistryBasicAuthUserInfo:          streamConfig.StreamKinesisDecoderPropSchemaRegistryBasicAuthUserInfo.ValueString(),
				StreamKinesisDecoderPropSchemaRegistryBasicAuthCredentialsSource: streamConfig.StreamKinesisDecoderPropSchemaRegistryBasicAuthCredentialsSource.ValueString(),
				StreamKinesisDecoderPropSchemaRegistryRestUrl:                    streamConfig.StreamKinesisDecoderPropSchemaRegistryRestUrl.ValueString(),
				StreamKinesisDecoderPropSchemaRegistrySchemaName:                 streamConfig.StreamKinesisDecoderPropSchemaRegistrySchemaName.ValueString(),
				StreamKinesisDecoderPropSchemaRegistrySslKeystoreLocation:        streamConfig.StreamKinesisDecoderPropSchemaRegistrySslKeystoreLocation.ValueString(),
				StreamKinesisDecoderPropSchemaRegistrySslKeystorePassword:        streamConfig.StreamKinesisDecoderPropSchemaRegistrySslKeystorePassword.ValueString(),
				StreamKinesisDecoderPropSchemaRegistrySslKeystoreType:            streamConfig.StreamKinesisDecoderPropSchemaRegistrySslKeystoreType.ValueString(),
				StreamKinesisDecoderPropSchemaRegistrySslTruststoreLocation:      streamConfig.StreamKinesisDecoderPropSchemaRegistrySslTruststoreLocation.ValueString(),
				StreamKinesisDecoderPropSchemaRegistrySslTruststorePassword:      streamConfig.StreamKinesisDecoderPropSchemaRegistrySslTruststorePassword.ValueString(),
				StreamKinesisDecoderPropSchemaRegistrySslTruststoreType:          streamConfig.StreamKinesisDecoderPropSchemaRegistrySslTruststoreType.ValueString(),
				StreamKinesisDecoderPropSchemaRegistrySslProtocol:                streamConfig.StreamKinesisDecoderPropSchemaRegistrySslProtocol.ValueString(),
				StreamPulsarAudience:                                             streamConfig.StreamPulsarAudience.ValueString(),
				StreamPulsarAuthenticationToken:                                  streamConfig.StreamPulsarAuthenticationToken.ValueString(),
				StreamPulsarBootstrapServers:                                     streamConfig.StreamPulsarBootstrapServers.ValueString(),
				StreamPulsarCredsFilePath:                                        streamConfig.StreamPulsarCredsFilePath.ValueString(),
				StreamPulsarConsumerFactoryClassName:                             streamConfig.StreamPulsarConsumerFactoryClassName.ValueString(),
				StreamPulsarConsumerPropAutoOffsetReset:                          streamConfig.StreamPulsarConsumerPropAutoOffsetReset.ValueString(),
				StreamPulsarConsumerType:                                         streamConfig.StreamPulsarConsumerType.ValueString(),
				StreamPulsarDecoderClassName:                                     streamConfig.StreamPulsarDecoderClassName.ValueString(),
				StreamPulsarFetchTimeoutMillis:                                   streamConfig.StreamPulsarFetchTimeoutMillis.ValueString(),
				StreamPulsarIssuerUrl:                                            streamConfig.StreamPulsarIssuerUrl.ValueString(),
				StreamPulsarMetadataPopulate:                                     streamConfig.StreamPulsarMetadataPopulate.ValueString(),
				StreamPulsarMetadataFields:                                       streamConfig.StreamPulsarMetadataFields.ValueString(),
				StreamPulsarDecoderPropSchemaRegistryBasicAuthUserInfo:           streamConfig.StreamPulsarDecoderPropSchemaRegistryBasicAuthUserInfo.ValueString(),
				StreamPulsarDecoderPropSchemaRegistryBasicAuthCredentialsSource:  streamConfig.StreamPulsarDecoderPropSchemaRegistryBasicAuthCredentialsSource.ValueString(),
				StreamPulsarDecoderPropSchemaRegistryRestUrl:                     streamConfig.StreamPulsarDecoderPropSchemaRegistryRestUrl.ValueString(),
				StreamPulsarDecoderPropSchemaRegistrySchemaName:                  streamConfig.StreamPulsarDecoderPropSchemaRegistrySchemaName.ValueString(),
				StreamPulsarDecoderPropSchemaRegistrySslKeystoreLocation:         streamConfig.StreamPulsarDecoderPropSchemaRegistrySslKeystoreLocation.ValueString(),
				StreamPulsarDecoderPropSchemaRegistrySslKeystorePassword:         streamConfig.StreamPulsarDecoderPropSchemaRegistrySslKeystorePassword.ValueString(),
				StreamPulsarDecoderPropSchemaRegistrySslKeystoreType:             streamConfig.StreamPulsarDecoderPropSchemaRegistrySslKeystoreType.ValueString(),
				StreamPulsarDecoderPropSchemaRegistrySslTruststoreLocation:       streamConfig.StreamPulsarDecoderPropSchemaRegistrySslTruststoreLocation.ValueString(),
				StreamPulsarDecoderPropSchemaRegistrySslTruststorePassword:       streamConfig.StreamPulsarDecoderPropSchemaRegistrySslTruststorePassword.ValueString(),
				StreamPulsarDecoderPropSchemaRegistrySslTruststoreType:           streamConfig.StreamPulsarDecoderPropSchemaRegistrySslTruststoreType.ValueString(),
				StreamPulsarDecoderPropSchemaRegistrySslProtocol:                 streamConfig.StreamPulsarDecoderPropSchemaRegistrySslProtocol.ValueString(),
				StreamPulsarTlsTrustCertsFilePath:                                streamConfig.StreamPulsarTlsTrustCertsFilePath.ValueString(),
				StreamPulsarTopicName:                                            streamConfig.StreamPulsarTopicName.ValueString(),
			})
		}

		ingestionConfig.StreamIngestionConfig = &model.StreamIngestionConfig{
			ColumnMajorSegmentBuilderEnabled: plan.IngestionConfig.StreamIngestionConfig.ColumnMajorSegmentBuilderEnabled.ValueBoolPointer(),
			TrackFilteredMessageOffsets:      plan.IngestionConfig.StreamIngestionConfig.TrackFilteredMessageOffsets.ValueBoolPointer(),
			StreamConfigMaps:                 streamConfigs,
		}
	}

	if plan.IngestionConfig.FilterConfig != nil {
		ingestionConfig.FilterConfig = &model.FilterConfig{
			FilterFunction: plan.IngestionConfig.FilterConfig.FilterFunction.ValueString(),
		}
	}

	return &ingestionConfig
}

func ToTierConfigs(plan *models.TableResourceModel) []*model.TierConfig {

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

func ToTenantsConfig(plan *models.TableResourceModel) model.TableTenant {
	var tagOverrideConfig *model.TenantTagOverrideConfig

	if plan.TenantsConfig.TagOverrideConfig != nil {
		tagOverrideConfig = &model.TenantTagOverrideConfig{
			RealtimeConsuming: plan.TenantsConfig.TagOverrideConfig.RealtimeConsuming.ValueString(),
			RealtimeCompleted: plan.TenantsConfig.TagOverrideConfig.RealtimeCompleted.ValueString(),
		}
	}

	return model.TableTenant{
		Broker:            plan.TenantsConfig.Broker.ValueString(),
		Server:            plan.TenantsConfig.Server.ValueString(),
		TagOverrideConfig: tagOverrideConfig,
	}
}

func ToTask(plan *models.Task) *model.Task {

	if plan == nil {
		return nil
	}

	return &model.Task{
		TaskTypeConfigsMap: plan.TaskTypeConfigsMap, // Changed to match new field name
	}
}

func toMap(ctx context.Context, inputMap types.Map) (map[string]string, diag.Diagnostics) {

	res := make(map[string]types.String, len(inputMap.Elements()))

	resultDiags := inputMap.ElementsAs(ctx, &res, false)
	if resultDiags.HasError() {
		return nil, resultDiags
	}

	return convertValuesToString(res), resultDiags

}

func toList[T any](ctx context.Context, inputList types.List) ([]T, diag.Diagnostics) {
	var res []T
	resultDiags := inputList.ElementsAs(ctx, &res, false)
	return res, resultDiags
}

func convertValuesToString(m map[string]types.String) map[string]string {
	res := make(map[string]string, len(m))
	for k, v := range m {
		res[k] = v.ValueString()
	}
	return res
}

func ToInstanceAssignmentConfigMap(instance *models.InstanceAssignmentConfigMap) *model.InstanceAssignmentConfigMap {

	if instance == nil {
		return nil
	}

	var consumingInstanceAssingment *model.InstanceAssignment
	if instance.Consuming != nil {
		consumingInstanceAssingment = &model.InstanceAssignment{
			TagPoolConfig:               ToTagPoolConfigInstanceAssignment(instance.Consuming),
			ReplicaGroupPartitionConfig: ToReplicaGroupPartitionConfig(instance.Consuming),
			PartitionSelector:           instance.Consuming.PartitionSelector.ValueString(),
			MinimizeDataMovement:        instance.Consuming.MinimizeDataMovement.ValueBool(),
		}
	}
	var completedInstanceAssingment *model.InstanceAssignment
	if instance.Completed != nil {
		completedInstanceAssingment = &model.InstanceAssignment{
			TagPoolConfig:               ToTagPoolConfigInstanceAssignment(instance.Completed),
			ReplicaGroupPartitionConfig: ToReplicaGroupPartitionConfig(instance.Completed),
			PartitionSelector:           instance.Completed.PartitionSelector.ValueString(),
			MinimizeDataMovement:        instance.Completed.MinimizeDataMovement.ValueBool(),
		}
	}
	var offlineInstanceAssingment *model.InstanceAssignment
	if instance.Offline != nil {
		offlineInstanceAssingment = &model.InstanceAssignment{
			TagPoolConfig:               ToTagPoolConfigInstanceAssignment(instance.Offline),
			ReplicaGroupPartitionConfig: ToReplicaGroupPartitionConfig(instance.Offline),
			PartitionSelector:           instance.Offline.PartitionSelector.ValueString(),
			MinimizeDataMovement:        instance.Offline.MinimizeDataMovement.ValueBool(),
		}
	}

	return &model.InstanceAssignmentConfigMap{
		Consuming: consumingInstanceAssingment,
		Completed: completedInstanceAssingment,
		Offline:   offlineInstanceAssingment,
	}
}

func ToTagPoolConfigInstanceAssignment(instance *models.InstanceAssignment) *model.TagPoolConfigInstanceAssignment {
	if instance.TagPoolConfig == nil {
		return nil
	}
	return &model.TagPoolConfigInstanceAssignment{
		Tag:       instance.TagPoolConfig.Tag.ValueString(),
		NumPools:  instance.TagPoolConfig.NumPools.ValueInt64(),
		PoolBased: instance.TagPoolConfig.PoolBased.ValueBool(),
	}
}

func ToReplicaGroupPartitionConfig(instance *models.InstanceAssignment) *model.ReplicaGroupPartitionInstanceAssignment {
	if instance.ReplicaGroupPartitionConfig == nil {
		return nil
	}
	return &model.ReplicaGroupPartitionInstanceAssignment{
		ReplicaGroupBased:           instance.ReplicaGroupPartitionConfig.ReplicaGroupBased.ValueBool(),
		NumInstances:                instance.ReplicaGroupPartitionConfig.NumInstances.ValueInt64(),
		NumReplicaGroups:            instance.ReplicaGroupPartitionConfig.NumReplicaGroups.ValueInt64(),
		NumInstancesPerReplicaGroup: instance.ReplicaGroupPartitionConfig.NumInstancesPerReplicaGroup.ValueInt64(),
		NumPartitions:               instance.ReplicaGroupPartitionConfig.NumPartitions.ValueInt64(),
		NumInstancesPerPartitions:   instance.ReplicaGroupPartitionConfig.NumInstancesPerPartitions.ValueInt64(),
		PartitionColumn:             instance.ReplicaGroupPartitionConfig.PartitionColumn.ValueString(),
		MinimizeDataMovement:        instance.ReplicaGroupPartitionConfig.MinimizeDataMovement.ValueBool(),
	}
}
