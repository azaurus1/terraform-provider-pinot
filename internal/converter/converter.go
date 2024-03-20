package converter

import (
	"github.com/azaurus1/go-pinot-api/model"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-pinot/internal/models"
)

func SetStateFromTable(state *models.TableResourceModel, table *model.Table) {

	state.TableName = types.StringValue(table.TableName)
	state.TableType = types.StringValue(table.TableType)

	state.TenantsConfig = &models.TenantsConfig{
		Broker: types.StringValue(table.Tenants.Broker),
		Server: types.StringValue(table.Tenants.Server),
	}

	state.SegmentsConfig = &models.SegmentsConfig{
		TimeType:           types.StringValue(table.SegmentsConfig.TimeType),
		Replication:        types.StringValue(table.SegmentsConfig.Replication),
		TimeColumnName:     types.StringValue(table.SegmentsConfig.TimeColumnName),
		RetentionTimeUnit:  types.StringValue(table.SegmentsConfig.RetentionTimeUnit),
		RetentionTimeValue: types.StringValue(table.SegmentsConfig.RetentionTimeValue),
	}

	var starTreeIndexConfigs []*models.StarTreeIndexConfigs
	for _, starConfig := range table.TableIndexConfig.StarTreeIndexConfigs {
		starTreeIndexConfigs = append(starTreeIndexConfigs, &models.StarTreeIndexConfigs{
			MaxLeafRecords:                  types.Int64Value(int64(starConfig.MaxLeafRecords)),
			DimensionsSplitOrder:            starConfig.DimensionsSplitOrder,
			FunctionColumnPairs:             starConfig.FunctionColumnPairs,
			SkipStarNodeCreationForDimNames: starConfig.SkipStarNodeCreationForDimensions,
		})
	}

	state.TableIndexConfig = &models.TableIndexConfig{
		LoadMode:            types.StringValue(table.TableIndexConfig.LoadMode),
		NullHandlingEnabled: types.BoolValue(table.TableIndexConfig.NullHandlingEnabled),
		CreateInvertedIndexDuringSegmentGeneration: types.BoolValue(table.TableIndexConfig.CreateInvertedIndexDuringSegmentGeneration),
		EnableDefaultStarTree:                      types.BoolValue(table.TableIndexConfig.EnableDefaultStarTree),
		OptimizeDictionary:                         types.BoolValue(table.TableIndexConfig.OptimizeDictionary),
		OptimizeDictionaryForMetrics:               types.BoolValue(table.TableIndexConfig.OptimizeDictionaryForMetrics),
		NoDictionarySizeRatioThreshold:             types.Float64Value(table.TableIndexConfig.NoDictionarySizeRatioThreshold),
		StarTreeIndexConfigs:                       starTreeIndexConfigs,
		SortedColumn:                               table.TableIndexConfig.SortedColumn,
		AggregateMetrics:                           types.BoolValue(table.TableIndexConfig.AggregateMetrics),
		SegmentPartitionConfig:                     convertSegmentPartitionConfig(table),
		OnHeapDictionaryColumns:                    table.TableIndexConfig.OnHeapDictionaryColumns,
		VarLengthDictionaryColumns:                 table.TableIndexConfig.VarLengthDictionaryColumns,
		RangeIndexColumns:                          table.TableIndexConfig.RangeIndexColumns,
		NoDictionaryColumns:                        table.TableIndexConfig.NoDictionaryColumns,
		BloomFilterColumns:                         table.TableIndexConfig.BloomFilterColumns,
		RangeIndexVersion:                          types.Int64Value(int64(table.TableIndexConfig.RangeIndexVersion)),
	}

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
