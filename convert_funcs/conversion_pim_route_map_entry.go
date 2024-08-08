

func createPimRouteMapEntry(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.PimRouteMapEntryResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["action"].(string); ok && v != "" {
		data.Action = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["descr"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["grp"].(string); ok && v != "" {
		data.Grp = types.StringValue(v)
	}
	if v, ok := attributes["name"].(string); ok && v != "" {
		data.Name = types.StringValue(v)
	}
	if v, ok := attributes["name_alias"].(string); ok && v != "" {
		data.NameAlias = types.StringValue(v)
	}
	if v, ok := attributes["order"].(string); ok && v != "" {
		data.Order = types.StringValue(v)
	}
	if v, ok := attributes["rp"].(string); ok && v != "" {
		data.Rp = types.StringValue(v)
	}
	if v, ok := attributes["src"].(string); ok && v != "" {
		data.Src = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationPimRouteMapEntry(attributes["annotations"])
	planTagTag := convertToTagTagPimRouteMapEntry(attributes["tags"])

	newAciPimRouteMapEntry := provider.GetPimRouteMapEntryCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciPimRouteMapEntry.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetPimRouteMapEntryId(ctx, data)
	attrs := payload["pimRouteMapEntry"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
func convertToTagAnnotationPimRouteMapEntry(resources interface{}) []provider.TagAnnotationPimRouteMapEntryResourceModel {
	var planResources []provider.TagAnnotationPimRouteMapEntryResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationPimRouteMapEntryResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagPimRouteMapEntry(resources interface{}) []provider.TagTagPimRouteMapEntryResourceModel {
	var planResources []provider.TagTagPimRouteMapEntryResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagPimRouteMapEntryResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
