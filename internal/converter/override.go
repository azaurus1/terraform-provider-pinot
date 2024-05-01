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
		DeletedKeysTTL:               float64(stateConfig.DeletedKeysTTL.ValueInt64()),
		HashFunction:                 stateConfig.HashFunction.ValueString(),
		EnableSnapshot:               stateConfig.EnableSnapshot.ValueBoolPointer(),
		EnablePreLoad:                stateConfig.EnablePreLoad.ValueBoolPointer(),
		UpsertTTL:                    stateConfig.UpsertTTL.ValueString(),
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

	noDictionaryCols, resultDiags := toList[string](ctx, plan.TableIndexConfig.NoDictionaryColumns)
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

	tableConfig := model.TableIndexConfig{
		CreateInvertedIndexDuringSegmentGeneration: plan.TableIndexConfig.CreateInvertedIndexDuringSegmentGeneration.ValueBool(),
		SortedColumn:                   sortedCols,
		RangeIndexColumns:              rangeIndexCols,
		NoDictionaryColumns:            noDictionaryCols,
		OnHeapDictionaryColumns:        onHeapDictionaryCols,
		VarLengthDictionaryColumns:     varLengthDictionaryCols,
		BloomFilterColumns:             bloomFilterCols,
		SegmentPartitionConfig:         ToSegmentPartitionConfig(plan),
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

	if plan.SegmentsConfig.DeletedSegmentsRetentionPeriod.ValueString() != "" {
		segmentsConfig.DeletedSegmentsRetentionPeriod = plan.SegmentsConfig.DeletedSegmentsRetentionPeriod.ValueString()
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
		SegmentTimeValueCheck: plan.IngestionConfig.SegmentTimeValueCheck.ValueBool(),
		RowTimeValueCheck:     plan.IngestionConfig.RowTimeValueCheck.ValueBool(),
		ContinueOnError:       plan.IngestionConfig.ContinueOnError.ValueBool(),
		StreamIngestionConfig: &model.StreamIngestionConfig{
			StreamConfigMaps: plan.IngestionConfig.StreamIngestionConfig.StreamConfigMaps,
		},
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
	return model.TableTenant{
		Broker: plan.TenantsConfig.Broker.ValueString(),
		Server: plan.TenantsConfig.Server.ValueString(),
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
