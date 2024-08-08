

func createMgmtRsOoBCons(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.MgmtRsOoBConsResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["prio"].(string); ok && v != "" {
		data.Prio = types.StringValue(v)
	}
	if v, ok := attributes["tn_vz_oob_br_cp_name"].(string); ok && v != "" {
		data.TnVzOOBBrCPName = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationMgmtRsOoBCons(attributes["annotations"])
	planTagTag := convertToTagTagMgmtRsOoBCons(attributes["tags"])

	newAciMgmtRsOoBCons := provider.GetMgmtRsOoBConsCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciMgmtRsOoBCons.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetMgmtRsOoBConsId(ctx, data)
	attrs := payload["mgmtRsOoBCons"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
func convertToTagAnnotationMgmtRsOoBCons(resources interface{}) []provider.TagAnnotationMgmtRsOoBConsResourceModel {
	var planResources []provider.TagAnnotationMgmtRsOoBConsResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationMgmtRsOoBConsResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagMgmtRsOoBCons(resources interface{}) []provider.TagTagMgmtRsOoBConsResourceModel {
	var planResources []provider.TagTagMgmtRsOoBConsResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagMgmtRsOoBConsResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
