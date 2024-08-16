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

func CreateNetflowRsMonitorToExporter(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.NetflowRsMonitorToExporterResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["netflow_exporter_policy_name"].(string); ok && v != "" {
		data.TnNetflowExporterPolName = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationNetflowRsMonitorToExporter(attributes["annotations"])
	planTagTag := convertToTagTagNetflowRsMonitorToExporter(attributes["tags"])

	newAciNetflowRsMonitorToExporter := provider.GetNetflowRsMonitorToExporterCreateJsonPayload(ctx, &diags, true, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciNetflowRsMonitorToExporter.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetNetflowRsMonitorToExporterId(ctx, data)
	attrs := payload["netflowRsMonitorToExporter"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
func convertToTagAnnotationNetflowRsMonitorToExporter(resources interface{}) []provider.TagAnnotationNetflowRsMonitorToExporterResourceModel {
	var planResources []provider.TagAnnotationNetflowRsMonitorToExporterResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationNetflowRsMonitorToExporterResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagNetflowRsMonitorToExporter(resources interface{}) []provider.TagTagNetflowRsMonitorToExporterResourceModel {
	var planResources []provider.TagTagNetflowRsMonitorToExporterResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagNetflowRsMonitorToExporterResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
