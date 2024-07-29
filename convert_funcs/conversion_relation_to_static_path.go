

func createFvRsPathAtt(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvRsPathAttResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["descr"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["encap"].(string); ok && v != "" {
		data.Encap = types.StringValue(v)
	}
	if v, ok := attributes["instr_imedcy"].(string); ok && v != "" {
		data.InstrImedcy = types.StringValue(v)
	}
	if v, ok := attributes["mode"].(string); ok && v != "" {
		data.Mode = types.StringValue(v)
	}
	if v, ok := attributes["primary_encap"].(string); ok && v != "" {
		data.PrimaryEncap = types.StringValue(v)
	}
	if v, ok := attributes["t_dn"].(string); ok && v != "" {
		data.TDn = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvRsPathAtt(attributes["annotations"])
	planTagTag := convertToTagTagFvRsPathAtt(attributes["tags"])

	newAciFvRsPathAtt := provider.GetFvRsPathAttCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvRsPathAtt.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFvRsPathAttId(ctx, data)
	attrs := payload["fvRsPathAtt"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

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
