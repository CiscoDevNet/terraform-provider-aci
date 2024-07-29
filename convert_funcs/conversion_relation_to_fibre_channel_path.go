

func createFvRsFcPathAtt(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvRsFcPathAttResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["descr"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["t_dn"].(string); ok && v != "" {
		data.TDn = types.StringValue(v)
	}
	if v, ok := attributes["vsan"].(string); ok && v != "" {
		data.Vsan = types.StringValue(v)
	}
	if v, ok := attributes["vsan_mode"].(string); ok && v != "" {
		data.VsanMode = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvRsFcPathAtt(attributes["annotations"])
	planTagTag := convertToTagTagFvRsFcPathAtt(attributes["tags"])

	newAciFvRsFcPathAtt := provider.GetFvRsFcPathAttCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvRsFcPathAtt.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFvRsFcPathAttId(ctx, data)
	attrs := payload["fvRsFcPathAtt"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
func convertToTagAnnotationFvRsFcPathAtt(resources interface{}) []provider.TagAnnotationFvRsFcPathAttResourceModel {
	var planResources []provider.TagAnnotationFvRsFcPathAttResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvRsFcPathAttResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvRsFcPathAtt(resources interface{}) []provider.TagTagFvRsFcPathAttResourceModel {
	var planResources []provider.TagTagFvRsFcPathAttResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvRsFcPathAttResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
