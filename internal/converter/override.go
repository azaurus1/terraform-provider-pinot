package converter

import (
	"context"
	"github.com/azaurus1/go-pinot-api/model"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-pinot/internal/models"
)

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
		Mode:                    stateConfig.Mode.ValueString(),
		PartialUpsertStrategies: partialUpsertStrategies,
		DeleteRecordColumn:      stateConfig.DeletedRecordColumn.ValueString(),
		// DeletedKeysTTL:          stateConfig.DeletedKeysTTL.ValueString()
		HashFunction:           stateConfig.HashFunction.ValueString(),
		EnableSnapshot:         stateConfig.EnableSnapshot.ValueBoolPointer(),
		EnablePreLoad:          stateConfig.EnablePreLoad.ValueBoolPointer(),
		UpsertTTL:              stateConfig.UpsertTTL.ValueString(),
		DropOutOfOrderRecords:  stateConfig.DropOutOfOrderRecords.ValueBoolPointer(),
		OutOfOrderRecordColumn: stateConfig.OutOfOrderRecordColumn.ValueString(),
		MetadataManagerClass:   stateConfig.MetadataManagerClass.ValueString(),
		MetadataManagerConfigs: metadataManagerConfigs,
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
