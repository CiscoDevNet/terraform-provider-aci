

func createFhsTrustCtrlPol(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FhsTrustCtrlPolResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["descr"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["has_dhcpv4_server"].(string); ok && v != "" {
		data.HasDhcpv4Server = types.StringValue(v)
	}
	if v, ok := attributes["has_dhcpv6_server"].(string); ok && v != "" {
		data.HasDhcpv6Server = types.StringValue(v)
	}
	if v, ok := attributes["has_ipv6_router"].(string); ok && v != "" {
		data.HasIpv6Router = types.StringValue(v)
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
	if v, ok := attributes["trust_arp"].(string); ok && v != "" {
		data.TrustArp = types.StringValue(v)
	}
	if v, ok := attributes["trust_nd"].(string); ok && v != "" {
		data.TrustNd = types.StringValue(v)
	}
	if v, ok := attributes["trust_ra"].(string); ok && v != "" {
		data.TrustRa = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFhsTrustCtrlPol(attributes["annotations"])
	planTagTag := convertToTagTagFhsTrustCtrlPol(attributes["tags"])

	newAciFhsTrustCtrlPol := provider.GetFhsTrustCtrlPolCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFhsTrustCtrlPol.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFhsTrustCtrlPolId(ctx, data)
	attrs := payload["fhsTrustCtrlPol"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}
func convertToTagAnnotationFhsTrustCtrlPol(resources interface{}) []provider.TagAnnotationFhsTrustCtrlPolResourceModel {
	var planResources []provider.TagAnnotationFhsTrustCtrlPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFhsTrustCtrlPolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFhsTrustCtrlPol(resources interface{}) []provider.TagTagFhsTrustCtrlPolResourceModel {
	var planResources []provider.TagTagFhsTrustCtrlPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFhsTrustCtrlPolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
