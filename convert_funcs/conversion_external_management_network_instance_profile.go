

func createMgmtInstP(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.MgmtInstPResourceModel{}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["descr"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["name"].(string); ok && v != "" {
		data.Name = types.StringValue(v)
	}
	if v, ok := attributes["name_alias"].(string); ok && v != "" {
		data.NameAlias = types.StringValue(v)
	}
	if v, ok := attributes["prio"].(string); ok && v != "" {
		data.Prio = types.StringValue(v)
	}
	planMgmtRsOoBCons := convertToMgmtRsOoBConsMgmtInstP(attributes["relation_to_consumed_out_of_band_contracts"])
	planTagAnnotation := convertToTagAnnotationMgmtInstP(attributes["annotations"])
	planTagTag := convertToTagTagMgmtInstP(attributes["tags"])

	newAciMgmtInstP := provider.GetMgmtInstPCreateJsonPayload(ctx, &diags, data, planMgmtRsOoBCons, planMgmtRsOoBCons, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciMgmtInstP.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetMgmtInstPId(ctx, data)
	attrs := payload["mgmtInstP"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
func convertToMgmtRsOoBConsMgmtInstP(resources interface{}) []provider.MgmtRsOoBConsMgmtInstPResourceModel {
	var planResources []provider.MgmtRsOoBConsMgmtInstPResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.MgmtRsOoBConsMgmtInstPResourceModel{
				Annotation:      types.StringValue(resourceMap["annotation"].(string)),
				Prio:            types.StringValue(resourceMap["prio"].(string)),
				TnVzOOBBrCPName: types.StringValue(resourceMap["tn_vz_oob_br_cp_name"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagAnnotationMgmtInstP(resources interface{}) []provider.TagAnnotationMgmtInstPResourceModel {
	var planResources []provider.TagAnnotationMgmtInstPResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationMgmtInstPResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagMgmtInstP(resources interface{}) []provider.TagTagMgmtInstPResourceModel {
	var planResources []provider.TagTagMgmtInstPResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagMgmtInstPResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
