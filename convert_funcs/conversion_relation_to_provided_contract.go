package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateFvRsProv(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvRsProvResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["match_criteria"].(string); ok && v != "" {
		data.MatchT = types.StringValue(v)
	}
	if v, ok := attributes["priority"].(string); ok && v != "" {
		data.Prio = types.StringValue(v)
	}
	if v, ok := attributes["contract_name"].(string); ok && v != "" {
		data.TnVzBrCPName = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvRsProv(attributes["annotations"])
	planTagTag := convertToTagTagFvRsProv(attributes["tags"])

	if status == "deleted" {

		provider.SetFvRsProvId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "fvRsProv", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}

	}

	newAciFvRsProv := provider.GetFvRsProvCreateJsonPayload(ctx, &diags, true, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvRsProv.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetFvRsProvId(ctx, data)
	attrs := payload["fvRsProv"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToTagAnnotationFvRsProv(resources interface{}) []provider.TagAnnotationFvRsProvResourceModel {
	var planResources []provider.TagAnnotationFvRsProvResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvRsProvResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvRsProv(resources interface{}) []provider.TagTagFvRsProvResourceModel {
	var planResources []provider.TagTagFvRsProvResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvRsProvResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
