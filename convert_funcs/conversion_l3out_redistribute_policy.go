package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateL3extRsRedistributePol(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.L3extRsRedistributePolResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["source"].(string); ok && v != "" {
		data.Src = types.StringValue(v)
	}
	if v, ok := attributes["route_control_profile_name"].(string); ok && v != "" {
		data.TnRtctrlProfileName = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationL3extRsRedistributePol(attributes["annotations"])
	planTagTag := convertToTagTagL3extRsRedistributePol(attributes["tags"])

	if status == "deleted" {

		provider.SetL3extRsRedistributePolId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "l3extRsRedistributePol", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}

	}

	newAciL3extRsRedistributePol := provider.GetL3extRsRedistributePolCreateJsonPayload(ctx, &diags, true, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciL3extRsRedistributePol.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetL3extRsRedistributePolId(ctx, data)
	attrs := payload["l3extRsRedistributePol"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToTagAnnotationL3extRsRedistributePol(resources interface{}) []provider.TagAnnotationL3extRsRedistributePolResourceModel {
	var planResources []provider.TagAnnotationL3extRsRedistributePolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationL3extRsRedistributePolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagL3extRsRedistributePol(resources interface{}) []provider.TagTagL3extRsRedistributePolResourceModel {
	var planResources []provider.TagTagL3extRsRedistributePolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagL3extRsRedistributePolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
