package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type TerraformPlan struct {
	PlannedValues struct {
		RootModule struct {
			Resources []Resource `json:"resources"`
		} `json:"root_module"`
	} `json:"planned_values"`
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

type Resource struct {
	Type   string                 `json:"type"`
	Name   string                 `json:"name"`
	Values map[string]interface{} `json:"values"`
}

type ResourceChange struct {
	Address string `json:"address"`
	Type    string `json:"type"`
	Change  struct {
		Actions []string               `json:"actions"`
		Before  map[string]interface{} `json:"before"`
		After   map[string]interface{} `json:"after"`
	} `json:"change"`
}

type Item struct {
	Attributes map[string]interface{} `json:"attributes"`
	Children   []map[string]Item      `json:"children,omitempty"`
}

func (i Item) MarshalJSON() ([]byte, error) {
	type Alias Item
	alias := struct {
		Attributes map[string]interface{} `json:"attributes"`
		Children   []map[string]Item      `json:"children,omitempty"`
	}{
		Attributes: i.Attributes,
		Children:   i.Children,
	}
	return json.Marshal(alias)
}

func runTerraformCommands() (string, error) {
	planFile := "plan.bin"
	jsonFile := "plan.json"

	cmdPlan := exec.Command("terraform", "plan", "-out="+planFile)
	if err := cmdPlan.Run(); err != nil {
		return "", fmt.Errorf("failed to run terraform plan: %w", err)
	}

	cmdShow := exec.Command("terraform", "show", "-json", planFile)
	output, err := os.Create(jsonFile)
	if err != nil {
		return "", fmt.Errorf("failed to create json file: %w", err)
	}
	defer output.Close()

	cmdShow.Stdout = output
	if err := cmdShow.Run(); err != nil {
		return "", fmt.Errorf("failed to run terraform show: %w", err)
	}

	return jsonFile, nil
}

func outputToFile(outputFile string, data interface{}) error {
	outputData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to convert data to JSON: %w", err)
	}

	err = os.WriteFile(outputFile, outputData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func processResources(terraformPlan TerraformPlan) []map[string]Item {
	var itemList []map[string]Item

	resourceStatusMap := make(map[string]string)
	for _, resourceChange := range terraformPlan.ResourceChanges {
		if len(resourceChange.Change.Actions) > 0 {
			resourceStatusMap[resourceChange.Address] = resourceChange.Change.Actions[0]
		}
	}

	for _, resourceChange := range terraformPlan.ResourceChanges {
		if len(resourceChange.Change.Actions) > 0 && resourceChange.Change.Actions[0] == "delete" {
			resourceType := resourceChange.Type
			item := createItem(resourceType, resourceChange.Change.Before, "deleted")
			if item != nil {
				itemList = append(itemList, item)
			}
		}
	}

	for _, resource := range terraformPlan.PlannedValues.RootModule.Resources {
		status := resourceStatusMap[fmt.Sprintf("%s.%s", resource.Type, resource.Name)]
		item := createItem(resource.Type, resource.Values, status)
		if item != nil {
			itemList = append(itemList, item)
		}
	}

	return itemList
}

func createItem(resourceType string, resourceValues map[string]interface{}, status string) map[string]Item {
	attributes := make(map[string]interface{})

	// Extract nested objects
	nestedObjects := extractNestedObjects(resourceValues)

	for key, val := range resourceValues {
		if key == "source" {
			attributes["src"] = val
		} else {
			attributes[key] = val
		}
	}
	if status == "deleted" {
		attributes["status"] = status
	}

	var item map[string]Item
	switch resourceType {
	case "aci_tenant":
		tenantAttributes := make(map[string]interface{})
		if name, exists := attributes["name"].(string); exists {
			tenantAttributes["dn"] = fmt.Sprintf("uni/tn-%s", name)
			tenantAttributes["name"] = name
		}
		if descr, exists := attributes["description"].(string); exists && descr != "" {
			tenantAttributes["descr"] = descr
		}
		if annotation, exists := attributes["annotation"].(string); exists && annotation != "" {
			tenantAttributes["annotation"] = annotation
		}
		if nameAlias, exists := attributes["name_alias"].(string); exists && nameAlias != "" {
			tenantAttributes["nameAlias"] = nameAlias
		}
		if status, exists := attributes["status"].(string); exists {
			tenantAttributes["status"] = status
		}
		item = map[string]Item{
			"fvTenant": {
				Attributes: tenantAttributes,
				Children:   createChildrenFromAttributes(tenantAttributes),
			},
		}
	case "aci_netflow_monitor_policy":
		netflowAttributes := make(map[string]interface{})
		if name, exists := attributes["name"].(string); exists {
			netflowAttributes["dn"] = fmt.Sprintf("%s/monitorpol-%s", attributes["parent_dn"].(string), name)
			netflowAttributes["id"] = netflowAttributes["dn"]
			netflowAttributes["name"] = name
		}
		if descr, exists := attributes["description"].(string); exists && descr != "" {
			netflowAttributes["descr"] = descr
		}
		if annotation, exists := attributes["annotation"].(string); exists && annotation != "" {
			netflowAttributes["annotation"] = annotation
		}
		if nameAlias, exists := attributes["name_alias"].(string); exists && nameAlias != "" {
			netflowAttributes["nameAlias"] = nameAlias
		}
		if ownerKey, exists := attributes["owner_key"].(string); exists && ownerKey != "" {
			netflowAttributes["ownerKey"] = ownerKey
		}
		if ownerTag, exists := attributes["owner_tag"].(string); exists && ownerTag != "" {
			netflowAttributes["ownerTag"] = ownerTag
		}
		if status, exists := attributes["status"].(string); exists {
			netflowAttributes["status"] = status
		}
		children := createChildrenFromAttributes(netflowAttributes)
		item = map[string]Item{
			"netflowMonitorPol": {
				Attributes: netflowAttributes,
				Children:   children,
			},
		}
	case "aci_annotation":
		annotationAttributes := make(map[string]interface{})
		if key, exists := attributes["key"].(string); exists && key != "" {
			annotationAttributes["key"] = key
		}
		if value, exists := attributes["value"].(string); exists && value != "" {
			annotationAttributes["value"] = value
		}
		if parentDN, exists := attributes["parent_dn"].(string); exists && parentDN != "" {
			annotationAttributes["dn"] = fmt.Sprintf("%s/annotation-%s", parentDN, attributes["key"])
			annotationAttributes["parent_dn"] = parentDN
		}
		if status, exists := attributes["status"].(string); exists {
			annotationAttributes["status"] = status
		}

		item = map[string]Item{
			"tagAnnotation": {
				Attributes: annotationAttributes,
			},
		}

	case "aci_endpoint_tag_ip":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["ctxName"].(string); exists && val != "" {
			attributes["ctxName"] = val
		}
		if val, exists := resourceValues["id"].(string); exists && val != "" {
			attributes["id"] = val
		}
		if val, exists := resourceValues["ip"].(string); exists && val != "" {
			attributes["ip"] = val
		}
		if val, exists := resourceValues["name"].(string); exists && val != "" {
			attributes["name"] = val
		}
		if val, exists := resourceValues["nameAlias"].(string); exists && val != "" {
			attributes["nameAlias"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"fvEpIpTag": {
				Attributes: attributes,
			},
		}

	case "aci_endpoint_tag_mac":
		endpointTagMacAttributes := make(map[string]interface{})
		var parentDN, mac, bdName string
		if val, exists := attributes["parent_dn"].(string); exists {
			parentDN = val
			endpointTagMacAttributes["parent_dn"] = val
		}
		if val, exists := attributes["mac"].(string); exists {
			mac = val
			endpointTagMacAttributes["mac"] = val
		}
		if val, exists := attributes["bd_name"].(string); exists {
			bdName = val
			endpointTagMacAttributes["bdName"] = val
		}
		endpointTagMacAttributes["dn"] = fmt.Sprintf("%s/eptags/epmactag-%s-[%s]", parentDN, mac, bdName)
		if name, exists := attributes["name"].(string); exists {
			endpointTagMacAttributes["name"] = name
		}

		if nameAlias, exists := attributes["name_alias"].(string); exists && nameAlias != "" {
			endpointTagMacAttributes["nameAlias"] = nameAlias
		}
		if annotation, exists := attributes["annotation"].(string); exists && annotation != "" {
			endpointTagMacAttributes["annotation"] = annotation
		}
		if status, exists := attributes["status"].(string); exists {
			endpointTagMacAttributes["status"] = status
		}
		children := createChildrenFromAttributes(endpointTagMacAttributes)
		item = map[string]Item{
			"fvEpMacTag": {
				Attributes: endpointTagMacAttributes,
				Children:   children,
			},
		}

	case "aci_external_management_network_instance_profile":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["description"].(string); exists && val != "" {
			attributes["descr"] = val
		}
		if val, exists := resourceValues["name"].(string); exists && val != "" {
			attributes["name"] = val

			attributes["dn"] = fmt.Sprintf("uni/tn-%s/extmgmt-%s/instp-%s", val, val, val)

		}
		if val, exists := resourceValues["name_alias"].(string); exists && val != "" {
			attributes["nameAlias"] = val
		}
		if val, exists := resourceValues["priority"].(string); exists && val != "" {
			attributes["priority"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"mgmtInstP": {
				Attributes: attributes,
			},
		}

	case "aci_external_management_network_subnet":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["descr"].(string); exists && val != "" {
			attributes["descr"] = val
		}
		if val, exists := resourceValues["ip"].(string); exists && val != "" {
			attributes["ip"] = val
		}
		if val, exists := resourceValues["name"].(string); exists && val != "" {
			attributes["name"] = val
		}
		if val, exists := resourceValues["nameAlias"].(string); exists && val != "" {
			attributes["nameAlias"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"mgmtSubnet": {
				Attributes: attributes,
				Children:   []map[string]Item{},
			},
		}

	case "aci_l3out_consumer_label":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["descr"].(string); exists && val != "" {
			attributes["descr"] = val
		}
		if val, exists := resourceValues["name"].(string); exists && val != "" {
			attributes["name"] = val
		}
		if val, exists := resourceValues["nameAlias"].(string); exists && val != "" {
			attributes["nameAlias"] = val
		}
		if val, exists := resourceValues["owner"].(string); exists && val != "" {
			attributes["owner"] = val
		}
		if val, exists := resourceValues["ownerKey"].(string); exists && val != "" {
			attributes["ownerKey"] = val
		}
		if val, exists := resourceValues["ownerTag"].(string); exists && val != "" {
			attributes["ownerTag"] = val
		}
		if val, exists := resourceValues["tag"].(string); exists && val != "" {
			attributes["tag"] = val
		}
		if val, exists := resourceValues["parent_dn"].(string); exists && val != "" {
			attributes["dn"] = fmt.Sprintf("%s/consLbl-%s", val, attributes["name"])
			attributes["parent_dn"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"l3extConsLbl": {
				Attributes: attributes,
			},
		}

	case "aci_l3out_node_sid_profile":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["descr"].(string); exists && val != "" {
			attributes["descr"] = val
		}
		if val, exists := resourceValues["loopbackAddr"].(string); exists && val != "" {
			attributes["loopbackAddr"] = val
		}
		if val, exists := resourceValues["name"].(string); exists && val != "" {
			attributes["name"] = val
		}
		if val, exists := resourceValues["nameAlias"].(string); exists && val != "" {
			attributes["nameAlias"] = val
		}
		if val, exists := resourceValues["sidoffset"].(string); exists && val != "" {
			attributes["sidoffset"] = val
		}
		if val, exists := resourceValues["parent_dn"].(string); exists && val != "" {
			attributes["dn"] = fmt.Sprintf("%s/sidP-%s", val, attributes["name"])
			attributes["parent_dn"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"mplsNodeSidP": {
				Attributes: attributes,
			},
		}

	case "aci_l3out_provider_label":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["descr"].(string); exists && val != "" {
			attributes["descr"] = val
		}
		if val, exists := resourceValues["name"].(string); exists && val != "" {
			attributes["name"] = val
		}
		if val, exists := resourceValues["nameAlias"].(string); exists && val != "" {
			attributes["nameAlias"] = val
		}
		if val, exists := resourceValues["ownerKey"].(string); exists && val != "" {
			attributes["ownerKey"] = val
		}
		if val, exists := resourceValues["ownerTag"].(string); exists && val != "" {
			attributes["ownerTag"] = val
		}
		if val, exists := resourceValues["tag"].(string); exists && val != "" {
			attributes["tag"] = val
		}
		if val, exists := resourceValues["parent_dn"].(string); exists && val != "" {
			attributes["dn"] = fmt.Sprintf("%s/provLbl-%s", val, attributes["name"])
			attributes["parent_dn"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"l3extProvLbl": {
				Attributes: attributes,
			},
		}

	case "aci_l3out_redistribute_policy":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}

		if val, exists := resourceValues["route_control_profile_name"].(string); exists && val != "" {
			attributes["tnRtctrlProfileName"] = val
		}

		if val, exists := resourceValues["source"].(string); exists && val != "" {
			attributes["src"] = val
		}
		if val, exists := resourceValues["parent_dn"].(string); exists && val != "" {
			attributes["dn"] = fmt.Sprintf("%s/rsredistributePol-[%s]-%s", val, attributes["tnRtctrlProfileName"], attributes["src"])
			attributes["parent_dn"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		delete(attributes, "route_control_profile_name")

		item = map[string]Item{
			"l3extRsRedistributePol": {
				Attributes: attributes,
			},
		}

	case "aci_out_of_band_contract":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["descr"].(string); exists && val != "" {
			attributes["descr"] = val
		}
		if val, exists := resourceValues["intent"].(string); exists && val != "" {
			attributes["intent"] = val
		}
		if val, exists := resourceValues["name"].(string); exists && val != "" {
			attributes["name"] = val
		}
		if val, exists := resourceValues["nameAlias"].(string); exists && val != "" {
			attributes["nameAlias"] = val
		}
		if val, exists := resourceValues["ownerKey"].(string); exists && val != "" {
			attributes["ownerKey"] = val
		}
		if val, exists := resourceValues["ownerTag"].(string); exists && val != "" {
			attributes["ownerTag"] = val
		}
		if val, exists := resourceValues["prio"].(string); exists && val != "" {
			attributes["prio"] = val
		}
		if val, exists := resourceValues["scope"].(string); exists && val != "" {
			attributes["scope"] = val
		}
		if val, exists := resourceValues["targetDscp"].(string); exists && val != "" {
			attributes["targetDscp"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"vzOOBBrCP": {
				Attributes: attributes,
			},
		}

	case "aci_pim_route_map_entry":
		if val, exists := resourceValues["action"].(string); exists && val != "" {
			attributes["action"] = val
		}
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["descr"].(string); exists && val != "" {
			attributes["descr"] = val
		}
		if val, exists := resourceValues["grp"].(string); exists && val != "" {
			attributes["grp"] = val
		}
		if val, exists := resourceValues["name"].(string); exists && val != "" {
			attributes["name"] = val
		}
		if val, exists := resourceValues["nameAlias"].(string); exists && val != "" {
			attributes["nameAlias"] = val
		}
		if val, exists := resourceValues["order"].(string); exists && val != "" {
			attributes["order"] = val
		}
		if val, exists := resourceValues["rp"].(string); exists && val != "" {
			attributes["rp"] = val
		}
		if val, exists := resourceValues["src"].(string); exists && val != "" {
			attributes["src"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"pimRouteMapEntry": {
				Attributes: attributes,
			},
		}

	case "aci_pim_route_map_policy":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["descr"].(string); exists && val != "" {
			attributes["descr"] = val
		}
		if val, exists := resourceValues["name"].(string); exists && val != "" {
			attributes["name"] = val
		}
		if val, exists := resourceValues["nameAlias"].(string); exists && val != "" {
			attributes["nameAlias"] = val
		}
		if val, exists := resourceValues["ownerKey"].(string); exists && val != "" {
			attributes["ownerKey"] = val
		}
		if val, exists := resourceValues["ownerTag"].(string); exists && val != "" {
			attributes["ownerTag"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"pimRouteMapPol": {
				Attributes: attributes,
			},
		}

	case "aci_relation_to_consumed_out_of_band_contract":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["prio"].(string); exists && val != "" {
			attributes["prio"] = val
		}
		if val, exists := resourceValues["out_of_band_contract_name"].(string); exists && val != "" {
			attributes["tnVzOOBBrCPName"] = val
		}
		if val, exists := resourceValues["parent_dn"].(string); exists && val != "" {
			attributes["dn"] = fmt.Sprintf("%s/rsOoBCons-%s", val, attributes["tnVzOOBBrCPName"])
			attributes["parent_dn"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"mgmtRsOoBCons": {
				Attributes: attributes,
			},
		}

	case "aci_relation_to_fallback_route_group":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["tDn"].(string); exists && val != "" {
			attributes["tDn"] = val
		}
		if val, exists := resourceValues["parent_dn"].(string); exists && val != "" {
			attributes["dn"] = fmt.Sprintf("%s/rsOutToFBRGroup-%s", val, attributes["tDn"])
			attributes["parent_dn"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"l3extRsOutToFBRGroup": {
				Attributes: attributes,
			},
		}

	case "aci_relation_to_netflow_exporter":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["netflow_exporter_policy_name"].(string); exists && val != "" {
			attributes["tnNetflowExporterPolName"] = val
		}
		if val, exists := resourceValues["parent_dn"].(string); exists && val != "" {
			attributes["dn"] = fmt.Sprintf("%s/rsMonitorToExporter-%s", val, attributes["tnNetflowExporterPolName"])
			attributes["parent_dn"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"netflowRsMonitorToExporter": {
				Attributes: attributes,
			},
		}

	case "aci_relation_to_netflow_record":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["netflow_record_policy_name"].(string); exists && val != "" {
			attributes["tnNetflowRecordPolName"] = val
		}
		if val, exists := resourceValues["parent_dn"].(string); exists && val != "" {
			attributes["dn"] = fmt.Sprintf("%s/rsMonitorToRecord-%s", val, attributes["tnNetflowRecordPolName"])
			attributes["parent_dn"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"netflowRsMonitorToRecord": {
				Attributes: attributes,
			},
		}

	case "aci_tag":
		if val, exists := resourceValues["key"].(string); exists && val != "" {
			attributes["key"] = val
		}
		if val, exists := resourceValues["value"].(string); exists && val != "" {
			attributes["value"] = val
		}
		if val, exists := resourceValues["parent_dn"].(string); exists && val != "" {
			attributes["dn"] = fmt.Sprintf("%s/tag-%s", val, attributes["key"])
			attributes["parent_dn"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"tagTag": {
				Attributes: attributes,
			},
		}

	case "aci_vrf_fallback_route_group_member":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["descr"].(string); exists && val != "" {
			attributes["descr"] = val
		}
		if val, exists := resourceValues["name"].(string); exists && val != "" {
			attributes["name"] = val
		}
		if val, exists := resourceValues["nameAlias"].(string); exists && val != "" {
			attributes["nameAlias"] = val
		}
		if val, exists := resourceValues["rnhAddr"].(string); exists && val != "" {
			attributes["rnhAddr"] = val
		}
		if val, exists := resourceValues["parent_dn"].(string); exists && val != "" {
			attributes["dn"] = fmt.Sprintf("%s/member-%s", val, attributes["name"])
			attributes["parent_dn"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"fvFBRMember": {
				Attributes: attributes,
			},
		}

	case "aci_vrf_fallback_route_group":
		if val, exists := resourceValues["annotation"].(string); exists && val != "" {
			attributes["annotation"] = val
		}
		if val, exists := resourceValues["descr"].(string); exists && val != "" {
			attributes["descr"] = val
		}
		if val, exists := resourceValues["name"].(string); exists && val != "" {
			attributes["name"] = val
		}
		if val, exists := resourceValues["nameAlias"].(string); exists && val != "" {
			attributes["nameAlias"] = val
		}
		if val, exists := resourceValues["parent_dn"].(string); exists && val != "" {
			attributes["dn"] = fmt.Sprintf("%s/group-%s", val, attributes["name"])
			attributes["parent_dn"] = val
		}
		if status, exists := attributes["status"].(string); exists {
			attributes["status"] = status
		}

		item = map[string]Item{
			"fvFBRGroup": {
				Attributes: attributes,
			},
		}
	}

	// Add nested objects back to the attributes
	for key, val := range nestedObjects {
		attributes[key] = val
	}

	children := createChildrenFromAttributes(attributes)
	if len(children) > 0 {
		for resourceType := range item {
			resource := item[resourceType]
			resource.Children = append(resource.Children, children...)
			item[resourceType] = resource
		}
	}

	// Remove attributes that should not appear in the final output
	delete(attributes, "annotations")
	delete(attributes, "tags")
	delete(attributes, "relation_to_netflow_exporters")
	delete(attributes, "relation_to_netflow_record")
	delete(attributes, "relation_to_consumed_out_of_band_contracts")

	return item
}

func extractNestedObjects(attributes map[string]interface{}) map[string]interface{} {
	nestedObjects := make(map[string]interface{})
	for key, val := range attributes {
		if _, ok := val.([]interface{}); ok {
			nestedObjects[key] = val
		}
	}
	for key := range nestedObjects {
		delete(attributes, key)
	}
	return nestedObjects
}

func createChildrenFromAttributes(attributes map[string]interface{}) []map[string]Item {
	var children []map[string]Item

	if annotations, exists := attributes["annotations"].([]interface{}); exists {
		for _, annotation := range annotations {
			if annotationMap, ok := annotation.(map[string]interface{}); ok {
				child := createTagAnnotation(annotationMap)
				children = append(children, child)
			}
		}
	}

	if tags, exists := attributes["tags"].([]interface{}); exists {
		for _, tag := range tags {
			if tagMap, ok := tag.(map[string]interface{}); ok {
				child := createTag(tagMap)
				children = append(children, child)
			}
		}
	}

	if netflowExporters, exists := attributes["relation_to_netflow_exporters"].([]interface{}); exists {
		children = append(children, extractRelations("netflowRsMonitorToExporter", netflowExporters)...)
	}

	if netflowRecords, exists := attributes["relation_to_netflow_record"].([]interface{}); exists {
		children = append(children, extractRelations("netflowRsMonitorToRecord", netflowRecords)...)
	}

	if rel_oobc, exists := attributes["relation_to_consumed_out_of_band_contracts"].([]interface{}); exists {
		children = append(children, extractRelations("mgmtRsOoBCons", rel_oobc)...)
	}

	return children
}

func createTagAnnotation(annotation map[string]interface{}) map[string]Item {
	attributes := make(map[string]interface{})
	if key, exists := annotation["key"].(string); exists && key != "" {
		attributes["key"] = key
	}
	if value, exists := annotation["value"].(string); exists && value != "" {
		attributes["value"] = value
	}
	if parentDN, exists := annotation["parent_dn"].(string); exists && parentDN != "" {
		attributes["dn"] = fmt.Sprintf("%s/annotation-%s", parentDN, attributes["key"])
		attributes["parent_dn"] = parentDN
	}

	return map[string]Item{
		"tagAnnotation": {
			Attributes: attributes,
		},
	}
}

func createTag(tag map[string]interface{}) map[string]Item {
	attributes := make(map[string]interface{})
	if key, exists := tag["key"].(string); exists && key != "" {
		attributes["key"] = key
	}
	if value, exists := tag["value"].(string); exists && value != "" {
		attributes["value"] = value
	}
	if parentDN, exists := tag["parent_dn"].(string); exists && parentDN != "" {
		attributes["dn"] = fmt.Sprintf("%s/tag-%s", parentDN, attributes["key"])
		attributes["parent_dn"] = parentDN
	}

	return map[string]Item{
		"tagTag": {
			Attributes: attributes,
		},
	}
}

func extractRelations(relationType string, relations interface{}) []map[string]Item {
	var children []map[string]Item
	if relList, ok := relations.([]interface{}); ok {
		for _, rel := range relList {
			if relMap, ok := rel.(map[string]interface{}); ok {
				attributes := make(map[string]interface{})
				for key, value := range relMap {
					attributes[key] = value
				}

				switch relationType {
				case "netflowRsMonitorToExporter":
					attributes["tnNetflowExporterPolName"] = attributes["netflow_exporter_policy_name"]
					delete(attributes, "netflow_exporter_policy_name")
				case "netflowRsMonitorToRecord":
					attributes["tnNetflowRecordPolName"] = attributes["netflow_record_policy_name"]
					delete(attributes, "netflow_record_policy_name")
				case "mgmtRsOoBCons":
					attributes["tnVzOOBBrCPName"] = attributes["out_of_band_contract_name"]
					delete(attributes, "out_of_band_contract_name")
				}

				children = append(children, map[string]Item{
					relationType: {
						Attributes: attributes,
					},
				})
			}
		}
	}
	return children
}

func constructTree(itemList []map[string]Item) []map[string]Item {
	dnMap := make(map[string]*Item)
	var root []map[string]Item

	// First pass: Create all resources and store them in a map by their dn
	for _, item := range itemList {
		for resourceType, resource := range item {
			dn, ok := resource.Attributes["dn"].(string)
			if !ok {
				log.Printf("Missing dn in resource of type %s", resourceType)
				continue
			}
			dnMap[dn] = &resource
			resource.Attributes["resourceType"] = resourceType
		}
	}

	// Second pass: Establish parent-child relationships using parent_dn
	for _, resource := range dnMap {
		parentDn, parentExists := resource.Attributes["parent_dn"].(string)
		if parentExists {
			if parent, exists := dnMap[parentDn]; exists {
				if parent.Children == nil {
					parent.Children = []map[string]Item{}
				}
				childType := resource.Attributes["resourceType"].(string)
				parent.Children = append(parent.Children, map[string]Item{childType: *resource})
			} else {
				root = append(root, map[string]Item{resource.Attributes["resourceType"].(string): *resource})
			}
		} else {
			root = append(root, map[string]Item{resource.Attributes["resourceType"].(string): *resource})
		}
	}

	return root
}

func replaceDnWithResourceType(data []map[string]Item) []map[string]Item {
	var result []map[string]Item

	for _, item := range data {
		for resourceType, resource := range item {
			if resource.Children != nil {
				resource.Children = replaceDnWithResourceType(resource.Children)
			}
			newItem := map[string]Item{resourceType: resource}
			result = append(result, newItem)
		}
	}

	return result
}

func removeDnAndResourceTypeFromAttributes(data []map[string]Item) {
	for _, item := range data {
		for _, resource := range item {
			delete(resource.Attributes, "resourceType")
			if resource.Children != nil {
				removeDnAndResourceTypeFromAttributes(resource.Children)
			}
		}
	}
}

func removeParentDnFromAttributes(data []map[string]Item) {
	for _, item := range data {
		for _, resource := range item {
			delete(resource.Attributes, "parent_dn")
			if resource.Children != nil {
				removeParentDnFromAttributes(resource.Children)
			}
		}
	}
}

func countTenants(itemList []map[string]Item) int {
	count := 0
	for _, item := range itemList {
		if _, exists := item["fvTenant"]; exists {
			count++
		}
	}
	return count
}

func unwrapOuterArray(data []map[string]Item) interface{} {
	if len(data) == 1 {
		return data[0]
	}
	return data
}

func countResources(itemList []map[string]Item) int {
	count := 0
	for _, item := range itemList {
		count++
		for _, resource := range item {
			if resource.Children != nil {
				count += countResources(resource.Children)
			}
		}
	}
	return count
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run converter.go <output.json>")
		os.Exit(1)
	}

	outputFile := os.Args[1]

	for attempts := 1; attempts <= 10; attempts++ {
		jsonFile, err := runTerraformCommands()
		if err != nil {
			log.Fatalf("Error running Terraform commands: %v", err)
		}

		inputData, err := os.ReadFile(jsonFile)
		if err != nil {
			log.Fatalf("Error reading input file: %v", err)
		}

		var terraformPlan TerraformPlan
		err = json.Unmarshal(inputData, &terraformPlan)
		if err != nil {
			log.Fatalf("Error parsing input file: %v", err)
		}

		itemList := processResources(terraformPlan)
		itemCount := countResources(itemList)

		tree := constructTree(itemList)
		finalTree := replaceDnWithResourceType(tree)
		removeDnAndResourceTypeFromAttributes(finalTree)
		removeParentDnFromAttributes(finalTree)

		finalTreeCount := countResources(finalTree)

		if itemCount == finalTreeCount {
			tenantCount := countTenants(finalTree)

			var output interface{}
			if tenantCount > 1 {
				output = map[string]interface{}{
					"polUni": map[string]interface{}{
						"attributes": map[string]interface{}{},
						"children":   finalTree,
					},
				}
			} else {
				output = finalTree
			}

			if tenantCount <= 1 {
				output = unwrapOuterArray(finalTree)
			}

			err = outputToFile(outputFile, output)
			if err != nil {
				log.Fatalf("Error writing output file: %v", err)
			}

			fmt.Printf("ACI Payload written to %s\n", outputFile)
			return
		}

		log.Printf("Attempt %d: Item count %d does not match final tree count %d. Retrying...\n", attempts, itemCount, finalTreeCount)
	}

	log.Fatalf("Failed to create matching item list and final tree after 10 attempts")
}
