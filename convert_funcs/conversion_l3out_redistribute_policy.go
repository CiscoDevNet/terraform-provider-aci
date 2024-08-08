

func createL3extRsRedistributePol(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.L3extRsRedistributePolResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["src"].(string); ok && v != "" {
		data.Src = types.StringValue(v)
	}
	if v, ok := attributes["tn_rtctrl_profile_name"].(string); ok && v != "" {
		data.TnRtctrlProfileName = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationL3extRsRedistributePol(attributes["annotations"])
	planTagTag := convertToTagTagL3extRsRedistributePol(attributes["tags"])

	newAciL3extRsRedistributePol := provider.GetL3extRsRedistributePolCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciL3extRsRedistributePol.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetL3extRsRedistributePolId(ctx, data)
	attrs := payload["l3extRsRedistributePol"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
func convertToTagAnnotationL3extRsRedistributePol(resources interface{}) []provider.TagAnnotationL3extRsRedistributePolResourceModel {
	var planResources []provider.TagAnnotationL3extRsRedistributePolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationL3extRsRedistributePolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagL3extRsRedistributePol(resources interface{}) []provider.TagTagL3extRsRedistributePolResourceModel {
	var planResources []provider.TagTagL3extRsRedistributePolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagL3extRsRedistributePolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
