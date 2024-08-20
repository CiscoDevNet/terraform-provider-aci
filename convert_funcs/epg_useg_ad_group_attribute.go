package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateFvIdGroupAttr(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvIdGroupAttrResourceModel{}
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
	if v, ok := attributes["selector"].(string); ok && v != "" {
		data.Selector = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvIdGroupAttr(attributes["annotations"])
	planTagTag := convertToTagTagFvIdGroupAttr(attributes["tags"])

	// Handle deletion logic
	if status == "deleted" {

		provider.SetFvIdGroupAttrId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "fvIdGroupAttr", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}

	}

	newAciFvIdGroupAttr := provider.GetFvIdGroupAttrCreateJsonPayload(ctx, &diags, true, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvIdGroupAttr.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetFvIdGroupAttrId(ctx, data)
	attrs := payload["fvIdGroupAttr"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToTagAnnotationFvIdGroupAttr(resources interface{}) []provider.TagAnnotationFvIdGroupAttrResourceModel {
	var planResources []provider.TagAnnotationFvIdGroupAttrResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvIdGroupAttrResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvIdGroupAttr(resources interface{}) []provider.TagTagFvIdGroupAttrResourceModel {
	var planResources []provider.TagTagFvIdGroupAttrResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvIdGroupAttrResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
