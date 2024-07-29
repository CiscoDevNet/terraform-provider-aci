

func createNetflowMonitorPol(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.NetflowMonitorPolResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["descr"].(string); ok && v != "" {
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
	planNetflowRsMonitorToExporter := convertToNetflowRsMonitorToExporterNetflowMonitorPol(attributes["relation_to_netflow_exporters"])
	planNetflowRsMonitorToRecord := convertToNetflowRsMonitorToRecordNetflowMonitorPol(attributes["relation_to_netflow_record"])
	planTagAnnotation := convertToTagAnnotationNetflowMonitorPol(attributes["annotations"])
	planTagTag := convertToTagTagNetflowMonitorPol(attributes["tags"])

	newAciNetflowMonitorPol := provider.GetNetflowMonitorPolCreateJsonPayload(ctx, &diags, data, planNetflowRsMonitorToExporter, planNetflowRsMonitorToExporter, planNetflowRsMonitorToRecord, planNetflowRsMonitorToRecord, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciNetflowMonitorPol.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetNetflowMonitorPolId(ctx, data)
	attrs := payload["netflowMonitorPol"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
func convertToNetflowRsMonitorToExporterNetflowMonitorPol(resources interface{}) []provider.NetflowRsMonitorToExporterNetflowMonitorPolResourceModel {
	var planResources []provider.NetflowRsMonitorToExporterNetflowMonitorPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.NetflowRsMonitorToExporterNetflowMonitorPolResourceModel{
				Annotation:               types.StringValue(resourceMap["annotation"].(string)),
				TnNetflowExporterPolName: types.StringValue(resourceMap["tn_netflow_exporter_pol_name"].(string)),
			})
		}
	}
	return planResources
}
func convertToNetflowRsMonitorToRecordNetflowMonitorPol(resources interface{}) []provider.NetflowRsMonitorToRecordNetflowMonitorPolResourceModel {
	var planResources []provider.NetflowRsMonitorToRecordNetflowMonitorPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.NetflowRsMonitorToRecordNetflowMonitorPolResourceModel{
				Annotation:             types.StringValue(resourceMap["annotation"].(string)),
				TnNetflowRecordPolName: types.StringValue(resourceMap["tn_netflow_record_pol_name"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagAnnotationNetflowMonitorPol(resources interface{}) []provider.TagAnnotationNetflowMonitorPolResourceModel {
	var planResources []provider.TagAnnotationNetflowMonitorPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationNetflowMonitorPolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagNetflowMonitorPol(resources interface{}) []provider.TagTagNetflowMonitorPolResourceModel {
	var planResources []provider.TagTagNetflowMonitorPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagNetflowMonitorPolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
