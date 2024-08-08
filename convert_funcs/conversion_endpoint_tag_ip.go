

func createFvEpIpTag(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvEpIpTagResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["ctx_name"].(string); ok && v != "" {
		data.CtxName = types.StringValue(v)
	}
	if v, ok := attributes["id"].(string); ok && v != "" {
		data.FvEpIpTagId = types.StringValue(v)
	}
	if v, ok := attributes["ip"].(string); ok && v != "" {
		data.Ip = types.StringValue(v)
	}
	if v, ok := attributes["name"].(string); ok && v != "" {
		data.Name = types.StringValue(v)
	}
	if v, ok := attributes["name_alias"].(string); ok && v != "" {
		data.NameAlias = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvEpIpTag(attributes["annotations"])
	planTagTag := convertToTagTagFvEpIpTag(attributes["tags"])

	newAciFvEpIpTag := provider.GetFvEpIpTagCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvEpIpTag.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFvEpIpTagId(ctx, data)
	attrs := payload["fvEpIpTag"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
func convertToTagAnnotationFvEpIpTag(resources interface{}) []provider.TagAnnotationFvEpIpTagResourceModel {
	var planResources []provider.TagAnnotationFvEpIpTagResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvEpIpTagResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvEpIpTag(resources interface{}) []provider.TagTagFvEpIpTagResourceModel {
	var planResources []provider.TagTagFvEpIpTagResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvEpIpTagResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
