package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreatePimRouteMapPol(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.PimRouteMapPolResourceModel{}
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
	planTagAnnotation := convertToTagAnnotationPimRouteMapPol(attributes["annotations"])
	planTagTag := convertToTagTagPimRouteMapPol(attributes["tags"])

	// Handle deletion logic
	if status == "deleted" {

		provider.SetPimRouteMapPolId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "pimRouteMapPol", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}

	}

	newAciPimRouteMapPol := provider.GetPimRouteMapPolCreateJsonPayload(ctx, &diags, true, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciPimRouteMapPol.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetPimRouteMapPolId(ctx, data)
	attrs := payload["pimRouteMapPol"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToTagAnnotationPimRouteMapPol(resources interface{}) []provider.TagAnnotationPimRouteMapPolResourceModel {
	var planResources []provider.TagAnnotationPimRouteMapPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationPimRouteMapPolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagPimRouteMapPol(resources interface{}) []provider.TagTagPimRouteMapPolResourceModel {
	var planResources []provider.TagTagPimRouteMapPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagPimRouteMapPolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
