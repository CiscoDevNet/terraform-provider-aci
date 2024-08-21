package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateMgmtInstP(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.MgmtInstPResourceModel{}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["description"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["name"].(string); ok && v != "" {
		data.Name = types.StringValue(v)
	}
	if v, ok := attributes["name_alias"].(string); ok && v != "" {
		data.NameAlias = types.StringValue(v)
	}
	if v, ok := attributes["priority"].(string); ok && v != "" {
		data.Prio = types.StringValue(v)
	}
	planMgmtRsOoBCons := convertToMgmtRsOoBConsMgmtInstP(attributes["relation_to_consumed_out_of_band_contracts"])
	planTagAnnotation := convertToTagAnnotationMgmtInstP(attributes["annotations"])
	planTagTag := convertToTagTagMgmtInstP(attributes["tags"])

	if status == "deleted" {

		provider.SetMgmtInstPId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "mgmtInstP", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}

	}

	newAciMgmtInstP := provider.GetMgmtInstPCreateJsonPayload(ctx, &diags, true, data, planMgmtRsOoBCons, planMgmtRsOoBCons, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciMgmtInstP.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetMgmtInstPId(ctx, data)
	attrs := payload["mgmtInstP"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToMgmtRsOoBConsMgmtInstP(resources interface{}) []provider.MgmtRsOoBConsMgmtInstPResourceModel {
	var planResources []provider.MgmtRsOoBConsMgmtInstPResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.MgmtRsOoBConsMgmtInstPResourceModel{
				Annotation:      types.StringValue(resourceMap["annotation"].(string)),
				Prio:            types.StringValue(resourceMap["priority"].(string)),
				TnVzOOBBrCPName: types.StringValue(resourceMap["out_of_band_contract_name"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagAnnotationMgmtInstP(resources interface{}) []provider.TagAnnotationMgmtInstPResourceModel {
	var planResources []provider.TagAnnotationMgmtInstPResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationMgmtInstPResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagMgmtInstP(resources interface{}) []provider.TagTagMgmtInstPResourceModel {
	var planResources []provider.TagTagMgmtInstPResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagMgmtInstPResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
