

func createNetflowRsMonitorToExporter(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.NetflowRsMonitorToExporterResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["tn_netflow_exporter_pol_name"].(string); ok && v != "" {
		data.TnNetflowExporterPolName = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationNetflowRsMonitorToExporter(attributes["annotations"])
	planTagTag := convertToTagTagNetflowRsMonitorToExporter(attributes["tags"])

	newAciNetflowRsMonitorToExporter := provider.GetNetflowRsMonitorToExporterCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciNetflowRsMonitorToExporter.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetNetflowRsMonitorToExporterId(ctx, data)
	attrs := payload["netflowRsMonitorToExporter"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

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
