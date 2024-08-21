package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateInfraHPathS(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.InfraHPathSResourceModel{}
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
	planInfraRsHPathAtt := convertToInfraRsHPathAttInfraHPathS(attributes["relation_to_host_path"])
	planInfraRsPathToAccBaseGrp := convertToInfraRsPathToAccBaseGrpInfraHPathS(attributes["relation_to_access_interface_policy_group"])
	planTagAnnotation := convertToTagAnnotationInfraHPathS(attributes["annotations"])
	planTagTag := convertToTagTagInfraHPathS(attributes["tags"])

	if status == "deleted" {

		provider.SetInfraHPathSId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "infraHPathS", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}

	}

	newAciInfraHPathS := provider.GetInfraHPathSCreateJsonPayload(ctx, &diags, true, data, planInfraRsHPathAtt, planInfraRsHPathAtt, planInfraRsPathToAccBaseGrp, planInfraRsPathToAccBaseGrp, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciInfraHPathS.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetInfraHPathSId(ctx, data)
	attrs := payload["infraHPathS"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToInfraRsHPathAttInfraHPathS(resources interface{}) []provider.InfraRsHPathAttInfraHPathSResourceModel {
	var planResources []provider.InfraRsHPathAttInfraHPathSResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.InfraRsHPathAttInfraHPathSResourceModel{
				Annotation: types.StringValue(resourceMap["annotation"].(string)),
				TDn:        types.StringValue(resourceMap["target_dn"].(string)),
			})
		}
	}
	return planResources
}
func convertToInfraRsPathToAccBaseGrpInfraHPathS(resources interface{}) []provider.InfraRsPathToAccBaseGrpInfraHPathSResourceModel {
	var planResources []provider.InfraRsPathToAccBaseGrpInfraHPathSResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.InfraRsPathToAccBaseGrpInfraHPathSResourceModel{
				Annotation: types.StringValue(resourceMap["annotation"].(string)),
				TDn:        types.StringValue(resourceMap["target_dn"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagAnnotationInfraHPathS(resources interface{}) []provider.TagAnnotationInfraHPathSResourceModel {
	var planResources []provider.TagAnnotationInfraHPathSResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationInfraHPathSResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagInfraHPathS(resources interface{}) []provider.TagTagInfraHPathSResourceModel {
	var planResources []provider.TagTagInfraHPathSResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagInfraHPathSResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
