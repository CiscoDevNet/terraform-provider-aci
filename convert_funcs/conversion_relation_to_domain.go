

func createFvRsDomAtt(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvRsDomAttResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["binding_type"].(string); ok && v != "" {
		data.BindingType = types.StringValue(v)
	}
	if v, ok := attributes["class_pref"].(string); ok && v != "" {
		data.ClassPref = types.StringValue(v)
	}
	if v, ok := attributes["custom_epg_name"].(string); ok && v != "" {
		data.CustomEpgName = types.StringValue(v)
	}
	if v, ok := attributes["delimiter"].(string); ok && v != "" {
		data.Delimiter = types.StringValue(v)
	}
	if v, ok := attributes["encap"].(string); ok && v != "" {
		data.Encap = types.StringValue(v)
	}
	if v, ok := attributes["encap_mode"].(string); ok && v != "" {
		data.EncapMode = types.StringValue(v)
	}
	if v, ok := attributes["epg_cos"].(string); ok && v != "" {
		data.EpgCos = types.StringValue(v)
	}
	if v, ok := attributes["epg_cos_pref"].(string); ok && v != "" {
		data.EpgCosPref = types.StringValue(v)
	}
	if v, ok := attributes["instr_imedcy"].(string); ok && v != "" {
		data.InstrImedcy = types.StringValue(v)
	}
	if v, ok := attributes["ipam_dhcp_override"].(string); ok && v != "" {
		data.IpamDhcpOverride = types.StringValue(v)
	}
	if v, ok := attributes["ipam_enabled"].(string); ok && v != "" {
		data.IpamEnabled = types.StringValue(v)
	}
	if v, ok := attributes["ipam_gateway"].(string); ok && v != "" {
		data.IpamGateway = types.StringValue(v)
	}
	if v, ok := attributes["lag_policy_name"].(string); ok && v != "" {
		data.LagPolicyName = types.StringValue(v)
	}
	if v, ok := attributes["netflow_dir"].(string); ok && v != "" {
		data.NetflowDir = types.StringValue(v)
	}
	if v, ok := attributes["netflow_pref"].(string); ok && v != "" {
		data.NetflowPref = types.StringValue(v)
	}
	if v, ok := attributes["num_ports"].(string); ok && v != "" {
		data.NumPorts = types.StringValue(v)
	}
	if v, ok := attributes["port_allocation"].(string); ok && v != "" {
		data.PortAllocation = types.StringValue(v)
	}
	if v, ok := attributes["primary_encap"].(string); ok && v != "" {
		data.PrimaryEncap = types.StringValue(v)
	}
	if v, ok := attributes["primary_encap_inner"].(string); ok && v != "" {
		data.PrimaryEncapInner = types.StringValue(v)
	}
	if v, ok := attributes["res_imedcy"].(string); ok && v != "" {
		data.ResImedcy = types.StringValue(v)
	}
	if v, ok := attributes["secondary_encap_inner"].(string); ok && v != "" {
		data.SecondaryEncapInner = types.StringValue(v)
	}
	if v, ok := attributes["switching_mode"].(string); ok && v != "" {
		data.SwitchingMode = types.StringValue(v)
	}
	if v, ok := attributes["t_dn"].(string); ok && v != "" {
		data.TDn = types.StringValue(v)
	}
	if v, ok := attributes["untagged"].(string); ok && v != "" {
		data.Untagged = types.StringValue(v)
	}
	if v, ok := attributes["vnet_only"].(string); ok && v != "" {
		data.VnetOnly = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvRsDomAtt(attributes["annotations"])
	planTagTag := convertToTagTagFvRsDomAtt(attributes["tags"])

	newAciFvRsDomAtt := provider.GetFvRsDomAttCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvRsDomAtt.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFvRsDomAttId(ctx, data)
	attrs := payload["fvRsDomAtt"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
func convertToTagAnnotationFvRsDomAtt(resources interface{}) []provider.TagAnnotationFvRsDomAttResourceModel {
	var planResources []provider.TagAnnotationFvRsDomAttResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvRsDomAttResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvRsDomAtt(resources interface{}) []provider.TagTagFvRsDomAttResourceModel {
	var planResources []provider.TagTagFvRsDomAttResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvRsDomAttResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
