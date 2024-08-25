package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateL3extConsLbl(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.L3extConsLblResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
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
	if v, ok := attributes["owner_key"].(string); ok && v != "" {
		data.OwnerKey = types.StringValue(v)
	}
	if v, ok := attributes["owner_tag"].(string); ok && v != "" {
		data.OwnerTag = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationL3extConsLbl(attributes["annotations"])
	planTagTag := convertToTagTagL3extConsLbl(attributes["tags"])

	if status == "deleted" {
		provider.SetL3extConsLblId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "l3extConsLbl", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}
	}

	newAciL3extConsLbl := provider.GetL3extConsLblCreateJsonPayload(ctx, &diags, true, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciL3extConsLbl.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetL3extConsLblId(ctx, data)
	attrs := payload["l3extConsLbl"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToTagAnnotationL3extConsLbl(resources interface{}) []provider.TagAnnotationL3extConsLblResourceModel {
	var planResources []provider.TagAnnotationL3extConsLblResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationL3extConsLblResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagL3extConsLbl(resources interface{}) []provider.TagTagL3extConsLblResourceModel {
	var planResources []provider.TagTagL3extConsLblResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagL3extConsLblResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
