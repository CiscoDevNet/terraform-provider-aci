

func createFvESg(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvESgResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["descr"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["exception_tag"].(string); ok && v != "" {
		data.ExceptionTag = types.StringValue(v)
	}
	if v, ok := attributes["match_t"].(string); ok && v != "" {
		data.MatchT = types.StringValue(v)
	}
	if v, ok := attributes["name"].(string); ok && v != "" {
		data.Name = types.StringValue(v)
	}
	if v, ok := attributes["name_alias"].(string); ok && v != "" {
		data.NameAlias = types.StringValue(v)
	}
	if v, ok := attributes["pc_enf_pref"].(string); ok && v != "" {
		data.PcEnfPref = types.StringValue(v)
	}
	if v, ok := attributes["pc_tag"].(string); ok && v != "" {
		data.PcTag = types.StringValue(v)
	}
	if v, ok := attributes["pref_gr_memb"].(string); ok && v != "" {
		data.PrefGrMemb = types.StringValue(v)
	}
	if v, ok := attributes["shutdown"].(string); ok && v != "" {
		data.Shutdown = types.StringValue(v)
	}
	planFvRsCons := convertToFvRsConsFvESg(attributes["relation_to_consumed_contracts"])
	planFvRsConsIf := convertToFvRsConsIfFvESg(attributes["relation_to_imported_contracts"])
	planFvRsIntraEpg := convertToFvRsIntraEpgFvESg(attributes["relation_to_intra_epg_contracts"])
	planFvRsProv := convertToFvRsProvFvESg(attributes["relation_to_provided_contracts"])
	planFvRsScope := convertToFvRsScopeFvESg(attributes["relation_to_vrf"])
	planFvRsSecInherited := convertToFvRsSecInheritedFvESg(attributes["relation_to_contract_masters"])
	planTagAnnotation := convertToTagAnnotationFvESg(attributes["annotations"])
	planTagTag := convertToTagTagFvESg(attributes["tags"])

	newAciFvESg := provider.GetFvESgCreateJsonPayload(ctx, &diags, data, planFvRsCons, planFvRsCons, planFvRsConsIf, planFvRsConsIf, planFvRsIntraEpg, planFvRsIntraEpg, planFvRsProv, planFvRsProv, planFvRsScope, planFvRsScope, planFvRsSecInherited, planFvRsSecInherited, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvESg.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFvESgId(ctx, data)
	attrs := payload["fvESg"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
func convertToFvRsConsFvESg(resources interface{}) []provider.FvRsConsFvESgResourceModel {
	var planResources []provider.FvRsConsFvESgResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.FvRsConsFvESgResourceModel{
				Annotation:   types.StringValue(resourceMap["annotation"].(string)),
				Prio:         types.StringValue(resourceMap["prio"].(string)),
				TnVzBrCPName: types.StringValue(resourceMap["tn_vz_br_cp_name"].(string)),
			})
		}
	}
	return planResources
}
func convertToFvRsConsIfFvESg(resources interface{}) []provider.FvRsConsIfFvESgResourceModel {
	var planResources []provider.FvRsConsIfFvESgResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.FvRsConsIfFvESgResourceModel{
				Annotation:   types.StringValue(resourceMap["annotation"].(string)),
				Prio:         types.StringValue(resourceMap["prio"].(string)),
				TnVzCPIfName: types.StringValue(resourceMap["tn_vz_cp_if_name"].(string)),
			})
		}
	}
	return planResources
}
func convertToFvRsIntraEpgFvESg(resources interface{}) []provider.FvRsIntraEpgFvESgResourceModel {
	var planResources []provider.FvRsIntraEpgFvESgResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.FvRsIntraEpgFvESgResourceModel{
				Annotation:   types.StringValue(resourceMap["annotation"].(string)),
				TnVzBrCPName: types.StringValue(resourceMap["tn_vz_br_cp_name"].(string)),
			})
		}
	}
	return planResources
}
func convertToFvRsProvFvESg(resources interface{}) []provider.FvRsProvFvESgResourceModel {
	var planResources []provider.FvRsProvFvESgResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.FvRsProvFvESgResourceModel{
				Annotation:   types.StringValue(resourceMap["annotation"].(string)),
				MatchT:       types.StringValue(resourceMap["match_t"].(string)),
				Prio:         types.StringValue(resourceMap["prio"].(string)),
				TnVzBrCPName: types.StringValue(resourceMap["tn_vz_br_cp_name"].(string)),
			})
		}
	}
	return planResources
}
func convertToFvRsScopeFvESg(resources interface{}) []provider.FvRsScopeFvESgResourceModel {
	var planResources []provider.FvRsScopeFvESgResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.FvRsScopeFvESgResourceModel{
				Annotation:  types.StringValue(resourceMap["annotation"].(string)),
				TnFvCtxName: types.StringValue(resourceMap["tn_fv_ctx_name"].(string)),
			})
		}
	}
	return planResources
}
func convertToFvRsSecInheritedFvESg(resources interface{}) []provider.FvRsSecInheritedFvESgResourceModel {
	var planResources []provider.FvRsSecInheritedFvESgResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.FvRsSecInheritedFvESgResourceModel{
				Annotation: types.StringValue(resourceMap["annotation"].(string)),
				TDn:        types.StringValue(resourceMap["t_dn"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagAnnotationFvESg(resources interface{}) []provider.TagAnnotationFvESgResourceModel {
	var planResources []provider.TagAnnotationFvESgResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvESgResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvESg(resources interface{}) []provider.TagTagFvESgResourceModel {
	var planResources []provider.TagTagFvESgResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvESgResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
