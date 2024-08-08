

func createVzOOBBrCP(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.VzOOBBrCPResourceModel{}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["descr"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["intent"].(string); ok && v != "" {
		data.Intent = types.StringValue(v)
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
	if v, ok := attributes["prio"].(string); ok && v != "" {
		data.Prio = types.StringValue(v)
	}
	if v, ok := attributes["scope"].(string); ok && v != "" {
		data.Scope = types.StringValue(v)
	}
	if v, ok := attributes["target_dscp"].(string); ok && v != "" {
		data.TargetDscp = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationVzOOBBrCP(attributes["annotations"])
	planTagTag := convertToTagTagVzOOBBrCP(attributes["tags"])

	newAciVzOOBBrCP := provider.GetVzOOBBrCPCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciVzOOBBrCP.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetVzOOBBrCPId(ctx, data)
	attrs := payload["vzOOBBrCP"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
func convertToTagAnnotationVzOOBBrCP(resources interface{}) []provider.TagAnnotationVzOOBBrCPResourceModel {
	var planResources []provider.TagAnnotationVzOOBBrCPResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationVzOOBBrCPResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagVzOOBBrCP(resources interface{}) []provider.TagTagVzOOBBrCPResourceModel {
	var planResources []provider.TagTagVzOOBBrCPResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagVzOOBBrCPResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
