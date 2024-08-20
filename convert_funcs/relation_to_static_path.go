package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateFvRsPathAtt(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvRsPathAttResourceModel{}
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
	if v, ok := attributes["primary_encapsulation"].(string); ok && v != "" {
		data.PrimaryEncap = types.StringValue(v)
	}
	if v, ok := attributes["target_dn"].(string); ok && v != "" {
		data.TDn = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvRsPathAtt(attributes["annotations"])
	planTagTag := convertToTagTagFvRsPathAtt(attributes["tags"])

	// Handle deletion logic
	if status == "deleted" {

		provider.SetFvRsPathAttId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "fvRsPathAtt", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}

	}

	newAciFvRsPathAtt := provider.GetFvRsPathAttCreateJsonPayload(ctx, &diags, true, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvRsPathAtt.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetFvRsPathAttId(ctx, data)
	attrs := payload["fvRsPathAtt"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToTagAnnotationFvRsPathAtt(resources interface{}) []provider.TagAnnotationFvRsPathAttResourceModel {
	var planResources []provider.TagAnnotationFvRsPathAttResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvRsPathAttResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvRsPathAtt(resources interface{}) []provider.TagTagFvRsPathAttResourceModel {
	var planResources []provider.TagTagFvRsPathAttResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvRsPathAttResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
