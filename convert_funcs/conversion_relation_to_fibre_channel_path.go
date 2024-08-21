package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateFvRsFcPathAtt(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvRsFcPathAttResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["description"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["target_dn"].(string); ok && v != "" {
		data.TDn = types.StringValue(v)
	}
	if v, ok := attributes["vsan"].(string); ok && v != "" {
		data.Vsan = types.StringValue(v)
	}
	if v, ok := attributes["vsan_mode"].(string); ok && v != "" {
		data.VsanMode = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvRsFcPathAtt(attributes["annotations"])
	planTagTag := convertToTagTagFvRsFcPathAtt(attributes["tags"])

	if status == "deleted" {

		provider.SetFvRsFcPathAttId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "fvRsFcPathAtt", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}

	}

	newAciFvRsFcPathAtt := provider.GetFvRsFcPathAttCreateJsonPayload(ctx, &diags, true, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvRsFcPathAtt.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetFvRsFcPathAttId(ctx, data)
	attrs := payload["fvRsFcPathAtt"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToTagAnnotationFvRsFcPathAtt(resources interface{}) []provider.TagAnnotationFvRsFcPathAttResourceModel {
	var planResources []provider.TagAnnotationFvRsFcPathAttResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvRsFcPathAttResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvRsFcPathAtt(resources interface{}) []provider.TagTagFvRsFcPathAttResourceModel {
	var planResources []provider.TagTagFvRsFcPathAttResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvRsFcPathAttResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
