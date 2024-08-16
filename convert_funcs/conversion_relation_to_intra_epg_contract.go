package convert_funcs

import (
	"context"
	"encoding/json"

	//"github.com/CiscoDevNet/terraform-provider-aci/v2/convert_funcs"
	//"github.com/CiscoDevNet/terraform-provider-aci/v2/convert_funcs"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateFvRsIntraEpg(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvRsIntraEpgResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["contract_name"].(string); ok && v != "" {
		data.TnVzBrCPName = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvRsIntraEpg(attributes["annotations"])
	planTagTag := convertToTagTagFvRsIntraEpg(attributes["tags"])

	newAciFvRsIntraEpg := provider.GetFvRsIntraEpgCreateJsonPayload(ctx, &diags, true, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvRsIntraEpg.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetFvRsIntraEpgId(ctx, data)
	attrs := payload["fvRsIntraEpg"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToTagAnnotationFvRsIntraEpg(resources interface{}) []provider.TagAnnotationFvRsIntraEpgResourceModel {
	var planResources []provider.TagAnnotationFvRsIntraEpgResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvRsIntraEpgResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvRsIntraEpg(resources interface{}) []provider.TagTagFvRsIntraEpgResourceModel {
	var planResources []provider.TagTagFvRsIntraEpgResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvRsIntraEpgResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
