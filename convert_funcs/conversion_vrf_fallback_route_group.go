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

func CreateFvFBRGroup(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvFBRGroupResourceModel{}
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
	planFvFBRMember := convertToFvFBRMemberFvFBRGroup(attributes["vrf_fallback_route_group_members"])
	planFvFBRoute := convertToFvFBRouteFvFBRGroup(attributes["vrf_fallback_routes"])
	planTagAnnotation := convertToTagAnnotationFvFBRGroup(attributes["annotations"])
	planTagTag := convertToTagTagFvFBRGroup(attributes["tags"])

	newAciFvFBRGroup := provider.GetFvFBRGroupCreateJsonPayload(ctx, &diags, true, data, planFvFBRMember, planFvFBRMember, planFvFBRoute, planFvFBRoute, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvFBRGroup.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetFvFBRGroupId(ctx, data)
	attrs := payload["fvFBRGroup"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToFvFBRMemberFvFBRGroup(resources interface{}) []provider.FvFBRMemberFvFBRGroupResourceModel {
	var planResources []provider.FvFBRMemberFvFBRGroupResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.FvFBRMemberFvFBRGroupResourceModel{
				Annotation: types.StringValue(resourceMap["annotation"].(string)),
				Descr:      types.StringValue(resourceMap["description"].(string)),
				Name:       types.StringValue(resourceMap["name"].(string)),
				NameAlias:  types.StringValue(resourceMap["name_alias"].(string)),
				RnhAddr:    types.StringValue(resourceMap["fallback_member"].(string)),
			})
		}
	}
	return planResources
}
func convertToFvFBRouteFvFBRGroup(resources interface{}) []provider.FvFBRouteFvFBRGroupResourceModel {
	var planResources []provider.FvFBRouteFvFBRGroupResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.FvFBRouteFvFBRGroupResourceModel{
				Annotation: types.StringValue(resourceMap["annotation"].(string)),
				Descr:      types.StringValue(resourceMap["description"].(string)),
				FbrPrefix:  types.StringValue(resourceMap["prefix_address"].(string)),
				Name:       types.StringValue(resourceMap["name"].(string)),
				NameAlias:  types.StringValue(resourceMap["name_alias"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagAnnotationFvFBRGroup(resources interface{}) []provider.TagAnnotationFvFBRGroupResourceModel {
	var planResources []provider.TagAnnotationFvFBRGroupResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvFBRGroupResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvFBRGroup(resources interface{}) []provider.TagTagFvFBRGroupResourceModel {
	var planResources []provider.TagTagFvFBRGroupResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvFBRGroupResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
