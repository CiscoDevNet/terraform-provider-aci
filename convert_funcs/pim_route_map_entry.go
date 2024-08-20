package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreatePimRouteMapEntry(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.PimRouteMapEntryResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["action"].(string); ok && v != "" {
		data.Action = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["description"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["group_ip"].(string); ok && v != "" {
		data.Grp = types.StringValue(v)
	}
	if v, ok := attributes["name"].(string); ok && v != "" {
		data.Name = types.StringValue(v)
	}
	if v, ok := attributes["name_alias"].(string); ok && v != "" {
		data.NameAlias = types.StringValue(v)
	}
	if v, ok := attributes["order"].(string); ok && v != "" {
		data.Order = types.StringValue(v)
	}
	if v, ok := attributes["rendezvous_point_ip"].(string); ok && v != "" {
		data.Rp = types.StringValue(v)
	}
	if v, ok := attributes["source_ip"].(string); ok && v != "" {
		data.Src = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationPimRouteMapEntry(attributes["annotations"])
	planTagTag := convertToTagTagPimRouteMapEntry(attributes["tags"])

	// Handle deletion logic
	if status == "deleted" {

		provider.SetPimRouteMapEntryId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "pimRouteMapEntry", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}

	}

	newAciPimRouteMapEntry := provider.GetPimRouteMapEntryCreateJsonPayload(ctx, &diags, true, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciPimRouteMapEntry.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetPimRouteMapEntryId(ctx, data)
	attrs := payload["pimRouteMapEntry"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToTagAnnotationPimRouteMapEntry(resources interface{}) []provider.TagAnnotationPimRouteMapEntryResourceModel {
	var planResources []provider.TagAnnotationPimRouteMapEntryResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationPimRouteMapEntryResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagPimRouteMapEntry(resources interface{}) []provider.TagTagPimRouteMapEntryResourceModel {
	var planResources []provider.TagTagPimRouteMapEntryResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagPimRouteMapEntryResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
