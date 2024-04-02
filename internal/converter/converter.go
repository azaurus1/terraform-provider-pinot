package converter

import (
	"context"
	"github.com/azaurus1/go-pinot-api/model"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-pinot/internal/models"
)

func SetStateFromTable(ctx context.Context, state *models.TableResourceModel, table *model.Table) diag.Diagnostics {

	var diags diag.Diagnostics

	state.TableName = types.StringValue(table.TableName)
	state.TableType = types.StringValue(table.TableType)

	state.TenantsConfig = &models.TenantsConfig{
		Broker: types.StringValue(table.Tenants.Broker),
		Server: types.StringValue(table.Tenants.Server),
	}

	state.SegmentsConfig = convertSegmentsConfig(table)
	state.TableIndexConfig = convertTableIndexConfig(ctx, table)

	var ingestionTransformConfigs []*models.TransformConfig
	for _, transformConfig := range table.IngestionConfig.TransformConfigs {
		ingestionTransformConfigs = append(ingestionTransformConfigs, &models.TransformConfig{
			ColumnName:        types.StringValue(transformConfig.ColumnName),
			TransformFunction: types.StringValue(transformConfig.TransformFunction),
		})
	}

	state.IngestionConfig = &models.IngestionConfig{
		SegmentTimeValueCheck: types.BoolValue(table.IngestionConfig.SegmentTimeValueCheck),
		RowTimeValueCheck:     types.BoolValue(table.IngestionConfig.RowTimeValueCheck),
		ContinueOnError:       types.BoolValue(table.IngestionConfig.ContinueOnError),
		StreamIngestionConfig: &models.StreamIngestionConfig{
			StreamConfigMaps: table.IngestionConfig.StreamIngestionConfig.StreamConfigMaps,
		},
		TransformConfigs: ingestionTransformConfigs,
	}

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

	state.TierConfigs = tierConfigs

	state.IsDimTable = types.BoolValue(table.IsDimTable)

	state.Metadata = &models.Metadata{
		CustomConfigs: table.Metadata.CustomConfigs,
	}

	// Routing Config
	routingConfig, routingDiags := convertRoutingConfig(ctx, table)
	if routingDiags.HasError() {
		diags.Append(routingDiags...)
		return diags
	}
	state.Routing = routingConfig

	// Upsert Config
	upsertConfig, upsertDiags := convertUpsertConfig(ctx, table)
	if upsertDiags.HasError() {
		diags.Append(upsertDiags...)
		return diags
	}
	state.UpsertConfig = upsertConfig

	return diags
}

func convertSegmentPartitionConfig(table *model.Table) *models.SegmentPartitionConfig {

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

func convertTableIndexConfig(ctx context.Context, table *model.Table) *models.TableIndexConfig {

	noDictionaryColumns, _ := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.NoDictionaryColumns)
	onHeapDictionaryColumns, _ := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.OnHeapDictionaryColumns)
	sortedColumn, _ := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.SortedColumn)
	varLengthDictionaryColumns, _ := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.VarLengthDictionaryColumns)
	bloomFilterColumns, _ := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.BloomFilterColumns)
	rangeIndexColumns, _ := types.ListValueFrom(ctx, types.StringType, table.TableIndexConfig.RangeIndexColumns)

	indexConfig := models.TableIndexConfig{
		LoadMode:            types.StringValue(table.TableIndexConfig.LoadMode),
		NullHandlingEnabled: types.BoolValue(table.TableIndexConfig.NullHandlingEnabled),
		CreateInvertedIndexDuringSegmentGeneration: types.BoolValue(table.TableIndexConfig.CreateInvertedIndexDuringSegmentGeneration),
		EnableDefaultStarTree:                      types.BoolValue(table.TableIndexConfig.EnableDefaultStarTree),
		OptimizeDictionary:                         types.BoolValue(table.TableIndexConfig.OptimizeDictionary),
		OptimizeDictionaryForMetrics:               types.BoolValue(table.TableIndexConfig.OptimizeDictionaryForMetrics),
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
	}

	return &indexConfig
}

func convertStarTreeIndexConfigs(ctx context.Context, table *model.Table) []*models.StarTreeIndexConfigs {

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

func convertSegmentsConfig(table *model.Table) *models.SegmentsConfig {

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

	if table.SegmentsConfig.DeletedSegmentsRetentionPeriod != "" {
		segmentsConfig.DeletedSegmentsRetentionPeriod = types.StringValue(table.SegmentsConfig.DeletedSegmentsRetentionPeriod)
	}

	return &segmentsConfig
}

func convertRoutingConfig(ctx context.Context, table *model.Table) (*models.RoutingConfig, diag.Diagnostics) {

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

func convertUpsertConfig(ctx context.Context, table *model.Table) (*models.UpsertConfig, diag.Diagnostics) {

	partialUpsertStrategies, resultDiags := types.MapValueFrom(ctx, types.StringType, table.UpsertConfig.PartialUpsertStrategies)
	if resultDiags.HasError() {
		return nil, resultDiags
	}

	metadataManagerConfigs, resultDiags := types.MapValueFrom(ctx, types.StringType, table.UpsertConfig.MetadataManagerConfigs)
	if resultDiags.HasError() {
		return nil, resultDiags
	}

	upsertConfig := models.UpsertConfig{
		Mode:                   types.StringValue(table.UpsertConfig.Mode),
		PartialUpsertStrategy:  partialUpsertStrategies,
		DeletedKeysTTL:         types.StringValue(strconv.Itoa(table.UpsertConfig.DeletedKeysTTL)),
		HashFunction:           types.StringValue(table.UpsertConfig.HashFunction),
		EnableSnapshot:         types.BoolPointerValue(table.UpsertConfig.EnableSnapshot),
		EnablePreLoad:          types.BoolPointerValue(table.UpsertConfig.EnablePreLoad),
		UpsertTTL:              types.StringValue(table.UpsertConfig.UpsertTTL),
		DropOutOfOrderRecords:  types.BoolPointerValue(table.UpsertConfig.DropOutOfOrderRecords),
		OutOfOrderRecordColumn: types.StringValue(table.UpsertConfig.OutOfOrderRecordColumn),
		MetadataManagerClass:   types.StringValue(table.UpsertConfig.MetadataManagerClass),
		MetadataManagerConfigs: metadataManagerConfigs,
	}

	return &upsertConfig, resultDiags
}
