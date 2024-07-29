

func createL3extRsOutToFBRGroup(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.L3extRsOutToFBRGroupResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["t_dn"].(string); ok && v != "" {
		data.TDn = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationL3extRsOutToFBRGroup(attributes["annotations"])
	planTagTag := convertToTagTagL3extRsOutToFBRGroup(attributes["tags"])

	newAciL3extRsOutToFBRGroup := provider.GetL3extRsOutToFBRGroupCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciL3extRsOutToFBRGroup.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetL3extRsOutToFBRGroupId(ctx, data)
	attrs := payload["l3extRsOutToFBRGroup"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
func convertToTagAnnotationL3extRsOutToFBRGroup(resources interface{}) []provider.TagAnnotationL3extRsOutToFBRGroupResourceModel {
	var planResources []provider.TagAnnotationL3extRsOutToFBRGroupResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationL3extRsOutToFBRGroupResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagL3extRsOutToFBRGroup(resources interface{}) []provider.TagTagL3extRsOutToFBRGroupResourceModel {
	var planResources []provider.TagTagL3extRsOutToFBRGroupResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagL3extRsOutToFBRGroupResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
