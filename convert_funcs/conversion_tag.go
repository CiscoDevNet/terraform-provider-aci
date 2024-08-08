

func createTagTag(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.TagTagResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["key"].(string); ok && v != "" {
		data.Key = types.StringValue(v)
	}
	if v, ok := attributes["value"].(string); ok && v != "" {
		data.Value = types.StringValue(v)
	}

	newAciTagTag := provider.GetTagTagCreateJsonPayload(ctx, &diags, data)

	jsonPayload := newAciTagTag.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetTagTagId(ctx, data)
	attrs := payload["tagTag"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
