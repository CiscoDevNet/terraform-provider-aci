package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateFvRsNodeAtt(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvRsNodeAttResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["description"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["encapsulation"].(string); ok && v != "" {
		data.Encap = types.StringValue(v)
	}
	if v, ok := attributes["deployment_immediacy"].(string); ok && v != "" {
		data.InstrImedcy = types.StringValue(v)
	}
	if v, ok := attributes["mode"].(string); ok && v != "" {
		data.Mode = types.StringValue(v)
	}
	if v, ok := attributes["target_dn"].(string); ok && v != "" {
		data.TDn = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvRsNodeAtt(attributes["annotations"])
	planTagTag := convertToTagTagFvRsNodeAtt(attributes["tags"])

	if status == "deleted" {

		provider.SetFvRsNodeAttId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "fvRsNodeAtt", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}

	}

	newAciFvRsNodeAtt := provider.GetFvRsNodeAttCreateJsonPayload(ctx, &diags, true, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvRsNodeAtt.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetFvRsNodeAttId(ctx, data)
	attrs := payload["fvRsNodeAtt"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToTagAnnotationFvRsNodeAtt(resources interface{}) []provider.TagAnnotationFvRsNodeAttResourceModel {
	var planResources []provider.TagAnnotationFvRsNodeAttResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvRsNodeAttResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvRsNodeAtt(resources interface{}) []provider.TagTagFvRsNodeAttResourceModel {
	var planResources []provider.TagTagFvRsNodeAttResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvRsNodeAttResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
