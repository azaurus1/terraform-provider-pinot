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

	state.TierConfigs = convertTierConfigs(table)
	state.TierConfigs = convertTierConfigs(table)
	state.Metadata = convertMetadata(table)

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

	return &models.IngestionConfig{
		SegmentTimeValueCheck: types.BoolValue(table.IngestionConfig.SegmentTimeValueCheck),
		RowTimeValueCheck:     types.BoolValue(table.IngestionConfig.RowTimeValueCheck),
		ContinueOnError:       types.BoolValue(table.IngestionConfig.ContinueOnError),
		StreamIngestionConfig: &models.StreamIngestionConfig{
			StreamConfigMaps: table.IngestionConfig.StreamIngestionConfig.StreamConfigMaps,
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

	if table.SegmentsConfig.DeletedSegmentsRetentionPeriod != "" {
		segmentsConfig.DeletedSegmentsRetentionPeriod = types.StringValue(table.SegmentsConfig.DeletedSegmentsRetentionPeriod)
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
		MetadataTTL:                  types.Int64Value(table.UpsertConfig.MetadataTTL),
		DeletedKeysTTL:               types.Int64Value(int64(table.UpsertConfig.DeletedKeysTTL)),
		HashFunction:                 types.StringValue(table.UpsertConfig.HashFunction),
		EnableSnapshot:               types.BoolPointerValue(table.UpsertConfig.EnableSnapshot),
		EnablePreLoad:                types.BoolPointerValue(table.UpsertConfig.EnablePreLoad),
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
