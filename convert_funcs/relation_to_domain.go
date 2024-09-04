package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateFvRsDomAtt(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvRsDomAttResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["custom_epg_name"].(string); ok && v != "" {
		data.CustomEpgName = types.StringValue(v)
	}
	if v, ok := attributes["delimiter"].(string); ok && v != "" {
		data.Delimiter = types.StringValue(v)
	}
	if v, ok := attributes["encapsulation"].(string); ok && v != "" {
		data.Encap = types.StringValue(v)
	}
	if v, ok := attributes["ipam_dhcp_override"].(string); ok && v != "" {
		data.IpamDhcpOverride = types.StringValue(v)
	}
	if v, ok := attributes["ipam_gateway"].(string); ok && v != "" {
		data.IpamGateway = types.StringValue(v)
	}
	if v, ok := attributes["lag_policy_name"].(string); ok && v != "" {
		data.LagPolicyName = types.StringValue(v)
	}
	if v, ok := attributes["primary_encapsulation"].(string); ok && v != "" {
		data.PrimaryEncap = types.StringValue(v)
	}
	if v, ok := attributes["primary_encapsulation_inner"].(string); ok && v != "" {
		data.PrimaryEncapInner = types.StringValue(v)
	}
	if v, ok := attributes["secondary_encapsulation_inner"].(string); ok && v != "" {
		data.SecondaryEncapInner = types.StringValue(v)
	}
	if v, ok := attributes["target_dn"].(string); ok && v != "" {
		data.TDn = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvRsDomAtt(attributes["annotations"])
	planTagTag := convertToTagTagFvRsDomAtt(attributes["tags"])

	if status == "deleted" {
		provider.SetFvRsDomAttId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "fvRsDomAtt", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}
	}

	newAciFvRsDomAtt := provider.GetFvRsDomAttCreateJsonPayload(ctx, &diags, true, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvRsDomAtt.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetFvRsDomAttId(ctx, data)
	attrs := payload["fvRsDomAtt"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToTagAnnotationFvRsDomAtt(resources interface{}) []provider.TagAnnotationFvRsDomAttResourceModel {
	var planResources []provider.TagAnnotationFvRsDomAttResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvRsDomAttResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvRsDomAtt(resources interface{}) []provider.TagTagFvRsDomAttResourceModel {
	var planResources []provider.TagTagFvRsDomAttResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvRsDomAttResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
