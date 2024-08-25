package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateQosDppPol(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.QosDppPolResourceModel{}
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
	planTagAnnotation := convertToTagAnnotationQosDppPol(attributes["annotations"])
	planTagTag := convertToTagTagQosDppPol(attributes["tags"])

	if status == "deleted" {
		provider.SetQosDppPolId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "qosDppPol", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}
	}

	newAciQosDppPol := provider.GetQosDppPolCreateJsonPayload(ctx, &diags, true, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciQosDppPol.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetQosDppPolId(ctx, data)
	attrs := payload["qosDppPol"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToTagAnnotationQosDppPol(resources interface{}) []provider.TagAnnotationQosDppPolResourceModel {
	var planResources []provider.TagAnnotationQosDppPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationQosDppPolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagQosDppPol(resources interface{}) []provider.TagTagQosDppPolResourceModel {
	var planResources []provider.TagTagQosDppPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagQosDppPolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
