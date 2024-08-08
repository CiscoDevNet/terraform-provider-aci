

func createTagAnnotation(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.TagAnnotationResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["key"].(string); ok && v != "" {
		data.Key = types.StringValue(v)
	}
	if v, ok := attributes["value"].(string); ok && v != "" {
		data.Value = types.StringValue(v)
	}

	newAciTagAnnotation := provider.GetTagAnnotationCreateJsonPayload(ctx, &diags, data)

	jsonPayload := newAciTagAnnotation.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetTagAnnotationId(ctx, data)
	attrs := payload["tagAnnotation"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
