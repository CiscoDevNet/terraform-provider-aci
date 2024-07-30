package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// EXCECUTES Terraform plan -out=plan.bin
// EXECUTES Terraform show -json plan.bin > plan.json
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

	for key, val := range resourceValues {
		attributes[key] = val
	}
	if status == "deleted" {
		attributes["status"] = status
	}

	var item map[string]Item
	switch resourceType {
	case "aci_tenant":
		item = createAciTenant(attributes)
	default:
		item = createProviderItem(resourceType, attributes)
	}

	if resourceType == "aci_tenant" {
		children := createChildrenFromAttributes(attributes)
		if len(children) > 0 {
			for resourceType := range item {
				resource := item[resourceType]
				resource.Children = append(resource.Children, children...)
				item[resourceType] = resource
			}
		}
	}

	return item
}

func createAciTenant(attributes map[string]interface{}) map[string]Item {
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
	return map[string]Item{
		"fvTenant": {
			Attributes: tenantAttributes,
			Children:   []map[string]Item{},
		},
	}
}

func createProviderItem(resourceType string, attributes map[string]interface{}) map[string]Item {
	var item map[string]Item

	switch resourceType {
	case "aci_endpoint_tag_ip":
		item = createAciEndpointTagIP(attributes)
	case "aci_external_management_network_instance_profile":
		item = createAciExternalManagementNetworkInstanceProfile(attributes)
	case "aci_vrf_fallback_route_group":
		item = createAciVrfFallbackRouteGroup(attributes)
	case "aci_endpoint_tag_mac":
		item = createAciEndpointTagMac(attributes)
	case "aci_netflow_monitor_policy":
		item = createAciNetflowMonitorPolicy(attributes)
	default:
		return nil
	}

	return item
}

func createAciEndpointTagIP(attributes map[string]interface{}) map[string]Item {
	endpointTagIPAttributes := make(map[string]interface{})
	if val, exists := attributes["annotation"].(string); exists && val != "" {
		endpointTagIPAttributes["annotation"] = val
	}
	if val, exists := attributes["ctxName"].(string); exists && val != "" {
		endpointTagIPAttributes["ctxName"] = val
	}
	if val, exists := attributes["id"].(string); exists && val != "" {
		endpointTagIPAttributes["id"] = val
	}
	if val, exists := attributes["ip"].(string); exists && val != "" {
		endpointTagIPAttributes["ip"] = val
	}
	if val, exists := attributes["name"].(string); exists && val != "" {
		endpointTagIPAttributes["name"] = val
	}
	if val, exists := attributes["nameAlias"].(string); exists && val != "" {
		endpointTagIPAttributes["nameAlias"] = val
	}
	if status, exists := attributes["status"].(string); exists {
		endpointTagIPAttributes["status"] = status
	}

	endpointTagIPAttributes["dn"] = fmt.Sprintf("%s/eptags/epiptag-[%s]-%s", attributes["parent_dn"], endpointTagIPAttributes["ip"], endpointTagIPAttributes["ctxName"])

	item := map[string]Item{
		"fvEpIpTag": {
			Attributes: endpointTagIPAttributes,
			Children:   []map[string]Item{},
		},
	}

	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvEpIpTagResourceModel{}

	// Process TagAnnotation child objects
	var planAnnotations []provider.TagAnnotationFvEpIpTagResourceModel
	if annotations, exists := attributes["annotations"].([]interface{}); exists {
		for _, annotation := range annotations {
			annotationMap := annotation.(map[string]interface{})
			planAnnotations = append(planAnnotations, provider.TagAnnotationFvEpIpTagResourceModel{
				Key:   types.StringValue(annotationMap["key"].(string)),
				Value: types.StringValue(annotationMap["value"].(string)),
			})
		}
	}

	var stateAnnotations []provider.TagAnnotationFvEpIpTagResourceModel
	tagAnnotationChildren := provider.GetFvEpIpTagTagAnnotationChildPayloads(ctx, &diags, data, planAnnotations, stateAnnotations)

	for _, child := range tagAnnotationChildren {
		epIpTagItem := item["fvEpIpTag"]
		epIpTagItem.Children = append(epIpTagItem.Children, map[string]Item{"tagAnnotation": {Attributes: child}})
		item["fvEpIpTag"] = epIpTagItem

	}

	return item
}

func createAciExternalManagementNetworkInstanceProfile(attributes map[string]interface{}) map[string]Item {
	profileAttributes := make(map[string]interface{})
	if name, exists := attributes["name"].(string); exists {
		profileAttributes["dn"] = fmt.Sprintf("%s/extmgmt-default/instp-%s", attributes["parent_dn"], name)
		profileAttributes["name"] = name
	}
	if descr, exists := attributes["description"].(string); exists && descr != "" {
		profileAttributes["descr"] = descr
	}
	if annotation, exists := attributes["annotation"].(string); exists && annotation != "" {
		profileAttributes["annotation"] = annotation
	}
	if nameAlias, exists := attributes["name_alias"].(string); exists && nameAlias != "" {
		profileAttributes["nameAlias"] = nameAlias
	}
	if status, exists := attributes["status"].(string); exists {
		profileAttributes["status"] = status
	}

	item := map[string]Item{
		"mgmtInstP": {
			Attributes: profileAttributes,
			Children:   []map[string]Item{},
		},
	}

	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.MgmtInstPResourceModel{}

	var planOoBCons []provider.MgmtRsOoBConsMgmtInstPResourceModel
	if oobCons, exists := attributes["oob_cons"].([]interface{}); exists {
		for _, oobCon := range oobCons {
			oobConMap := oobCon.(map[string]interface{})
			planOoBCons = append(planOoBCons, provider.MgmtRsOoBConsMgmtInstPResourceModel{
				Annotation: types.StringValue(oobConMap["annotation"].(string)),
			})
		}
	}

	var stateOoBCons []provider.MgmtRsOoBConsMgmtInstPResourceModel
	oobConsChildren := provider.GetMgmtInstPMgmtRsOoBConsChildPayloads(ctx, &diags, data, planOoBCons, stateOoBCons)
	for _, child := range oobConsChildren {
		children := item["mgmtInstP"].Children
		children = append(children, map[string]Item{"mgmtRsOoBCons": {Attributes: child}})
		item["mgmtInstP"] = Item{
			Attributes: item["mgmtInstP"].Attributes,
			Children:   children,
		}
	}

	return item
}

func createAciVrfFallbackRouteGroup(attributes map[string]interface{}) map[string]Item {

	vrfAttributes := make(map[string]interface{})

	if val, exists := attributes["parent_dn"].(string); exists {
		vrfAttributes["parent_dn"] = val
	}

	if name, exists := attributes["name"].(string); exists {
		vrfAttributes["name"] = name
		vrfAttributes["dn"] = fmt.Sprintf("%s/fbrg-%s", vrfAttributes["parent_dn"], name)
	}
	if descr, exists := attributes["description"].(string); exists && descr != "" {
		vrfAttributes["descr"] = descr
	}
	if annotation, exists := attributes["annotation"].(string); exists && annotation != "" {
		vrfAttributes["annotation"] = annotation
	}
	if nameAlias, exists := attributes["name_alias"].(string); exists && nameAlias != "" {
		vrfAttributes["nameAlias"] = nameAlias
	}
	if status, exists := attributes["status"].(string); exists {
		vrfAttributes["status"] = status
	}

	item := map[string]Item{
		"fvFBRGroup": {
			Attributes: vrfAttributes,
			Children:   []map[string]Item{},
		},
	}

	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvFBRGroupResourceModel{}

	var planMembers []provider.FvFBRMemberFvFBRGroupResourceModel

	if members, exists := attributes["vrf_fallback_route_group_members"].([]interface{}); exists {
		for _, member := range members {
			memberMap := member.(map[string]interface{})
			planMembers = append(planMembers, provider.FvFBRMemberFvFBRGroupResourceModel{
				Annotation: types.StringValue(memberMap["annotation"].(string)),
				Descr:      types.StringValue(memberMap["description"].(string)),
				Name:       types.StringValue(memberMap["name"].(string)),
				NameAlias:  types.StringValue(memberMap["name_alias"].(string)),
				RnhAddr:    types.StringValue(memberMap["fallback_member"].(string)),
			})
		}
	}

	var stateMembers []provider.FvFBRMemberFvFBRGroupResourceModel
	memberChildren := provider.GetFvFBRGroupFvFBRMemberChildPayloads(ctx, &diags, data, planMembers, stateMembers)
	for _, child := range memberChildren {
		children := item["fvFBRGroup"].Children
		children = append(children, map[string]Item{"fvFBRMember": {Attributes: child}})
		item["fvFBRGroup"] = Item{
			Attributes: item["fvFBRGroup"].Attributes,
			Children:   children,
		}
	}

	var planAnnotations []provider.TagAnnotationFvFBRGroupResourceModel
	if annotations, exists := attributes["annotations"].([]interface{}); exists {
		for _, annotation := range annotations {
			annotationMap := annotation.(map[string]interface{})
			planAnnotations = append(planAnnotations, provider.TagAnnotationFvFBRGroupResourceModel{
				Key:   types.StringValue(annotationMap["key"].(string)),
				Value: types.StringValue(annotationMap["value"].(string)),
			})
		}
	}
	var stateAnnotations []provider.TagAnnotationFvFBRGroupResourceModel
	tagAnnotationChildren := provider.GetFvFBRGroupTagAnnotationChildPayloads(ctx, &diags, data, planAnnotations, stateAnnotations)

	for _, child := range tagAnnotationChildren {
		children := item["fvFBRGroup"].Children
		children = append(children, map[string]Item{"tagAnnotation": {Attributes: child}})
		item["fvFBRGroup"] = Item{
			Attributes: item["fvFBRGroup"].Attributes,
			Children:   children,
		}
	}

	var planTags []provider.TagTagFvFBRGroupResourceModel
	if tags, exists := attributes["tags"].([]interface{}); exists {
		for _, tag := range tags {
			tagMap := tag.(map[string]interface{})
			planTags = append(planTags, provider.TagTagFvFBRGroupResourceModel{
				Key:   types.StringValue(tagMap["key"].(string)),
				Value: types.StringValue(tagMap["value"].(string)),
			})
		}
	}

	var stateTags []provider.TagTagFvFBRGroupResourceModel
	tagTagChildren := provider.GetFvFBRGroupTagTagChildPayloads(ctx, &diags, data, planTags, stateTags)
	for _, child := range tagTagChildren {
		children := item["fvFBRGroup"].Children
		children = append(children, map[string]Item{"tagTag": {Attributes: child}})
		item["fvFBRGroup"] = Item{
			Attributes: item["fvFBRGroup"].Attributes,
			Children:   children,
		}
	}

	return item
}

func createAciEndpointTagMac(attributes map[string]interface{}) map[string]Item {
	endpointTagMacAttributes := make(map[string]interface{})

	if val, exists := attributes["parent_dn"].(string); exists && val != "" {
		endpointTagMacAttributes["parent_dn"] = val
	}
	if val, exists := attributes["annotation"].(string); exists && val != "" {
		endpointTagMacAttributes["annotation"] = val
	}
	if val, exists := attributes["bd_name"].(string); exists && val != "" {
		endpointTagMacAttributes["bd_name"] = val
	}
	if val, exists := attributes["id_attribute"].(string); exists && val != "" {
		endpointTagMacAttributes["id_attribute"] = val
	}
	if val, exists := attributes["mac"].(string); exists && val != "" {
		endpointTagMacAttributes["mac"] = val
	}
	if val, exists := attributes["name"].(string); exists && val != "" {
		endpointTagMacAttributes["name"] = val
	}
	if val, exists := attributes["name_alias"].(string); exists && val != "" {
		endpointTagMacAttributes["name_alias"] = val
	}
	if status, exists := attributes["status"].(string); exists {
		endpointTagMacAttributes["status"] = status
	}

	endpointTagMacAttributes["dn"] = fmt.Sprintf("%s/eptags/epmactag-[%s]-[%s]", endpointTagMacAttributes["parent_dn"], endpointTagMacAttributes["mac"], endpointTagMacAttributes["bd_name"])

	item := map[string]Item{
		"fvEpMacTag": {
			Attributes: endpointTagMacAttributes,
			Children:   []map[string]Item{},
		},
	}

	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvEpMacTagResourceModel{}

	var planAnnotations []provider.TagAnnotationFvEpMacTagResourceModel
	if annotations, exists := attributes["annotations"].([]interface{}); exists {
		for _, annotation := range annotations {
			annotationMap := annotation.(map[string]interface{})
			planAnnotations = append(planAnnotations, provider.TagAnnotationFvEpMacTagResourceModel{
				Key:   types.StringValue(annotationMap["key"].(string)),
				Value: types.StringValue(annotationMap["value"].(string)),
			})
		}
	}

	var stateAnnotations []provider.TagAnnotationFvEpMacTagResourceModel
	tagAnnotationChildren := provider.GetFvEpMacTagTagAnnotationChildPayloads(ctx, &diags, data, planAnnotations, stateAnnotations)
	for _, child := range tagAnnotationChildren {
		children := item["fvEpMacTag"].Children
		children = append(children, map[string]Item{"tagAnnotation": {Attributes: child}})
		item["fvEpMacTag"] = Item{
			Attributes: item["fvEpMacTag"].Attributes,
			Children:   children,
		}
	}

	return item
}

func createAciNetflowMonitorPolicy(attributes map[string]interface{}) map[string]Item {
	monitorPolicyAttributes := make(map[string]interface{})

	if val, exists := attributes["annotation"].(string); exists {
		monitorPolicyAttributes["annotation"] = val
	}

	if val, exists := attributes["parent_dn"].(string); exists {
		monitorPolicyAttributes["parent_dn"] = val
	}

	if val, exists := attributes["description"].(string); exists {
		monitorPolicyAttributes["descr"] = val
	}
	if val, exists := attributes["name"].(string); exists {
		monitorPolicyAttributes["name"] = val
		monitorPolicyAttributes["dn"] = fmt.Sprintf("%s/monitorpol-%s", monitorPolicyAttributes["parent_dn"], val)
	}
	if val, exists := attributes["name_alias"].(string); exists {
		monitorPolicyAttributes["nameAlias"] = val
	}
	if val, exists := attributes["owner_key"].(string); exists {
		monitorPolicyAttributes["ownerKey"] = val
	}
	if val, exists := attributes["owner_tag"].(string); exists {
		monitorPolicyAttributes["ownerTag"] = val
	}
	if status, exists := attributes["status"].(string); exists {
		monitorPolicyAttributes["status"] = status
	}

	item := map[string]Item{
		"netflowMonitorPol": {
			Attributes: monitorPolicyAttributes,
			Children:   []map[string]Item{},
		},
	}

	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.NetflowMonitorPolResourceModel{}

	var planExporters []provider.NetflowRsMonitorToExporterNetflowMonitorPolResourceModel
	if exporters, exists := attributes["relation_to_netflow_exporters"].([]interface{}); exists {
		for _, exporter := range exporters {
			exporterMap := exporter.(map[string]interface{})
			planExporters = append(planExporters, provider.NetflowRsMonitorToExporterNetflowMonitorPolResourceModel{
				Annotation:               types.StringValue(exporterMap["annotation"].(string)),
				TnNetflowExporterPolName: types.StringValue(exporterMap["netflow_exporter_policy_name"].(string)),
			})
		}
	}

	var stateExporters []provider.NetflowRsMonitorToExporterNetflowMonitorPolResourceModel
	exporterChildren := provider.GetNetflowMonitorPolNetflowRsMonitorToExporterChildPayloads(ctx, &diags, data, planExporters, stateExporters)

	for _, child := range exporterChildren {
		children := item["netflowMonitorPol"].Children
		children = append(children, map[string]Item{"netflowRsMonitorToExporter": {Attributes: child}})
		item["netflowMonitorPol"] = Item{
			Attributes: item["netflowMonitorPol"].Attributes,
			Children:   children,
		}

	}

	var planRecords []provider.NetflowRsMonitorToRecordNetflowMonitorPolResourceModel
	if records, exists := attributes["relation_to_netflow_record"].([]interface{}); exists {
		for _, record := range records {
			recordMap := record.(map[string]interface{})
			planRecords = append(planRecords, provider.NetflowRsMonitorToRecordNetflowMonitorPolResourceModel{
				Annotation:             types.StringValue(recordMap["annotation"].(string)),
				TnNetflowRecordPolName: types.StringValue(recordMap["netflow_record_policy_name"].(string)),
			})
		}
	}

	var stateRecords []provider.NetflowRsMonitorToRecordNetflowMonitorPolResourceModel
	recordChildren := provider.GetNetflowMonitorPolNetflowRsMonitorToRecordChildPayloads(ctx, &diags, data, planRecords, stateRecords)

	for _, child := range recordChildren {
		children := item["netflowMonitorPol"].Children
		children = append(children, map[string]Item{"netflowRsMonitorToRecord": {Attributes: child}})
		item["netflowMonitorPol"] = Item{
			Attributes: item["netflowMonitorPol"].Attributes,
			Children:   children,
		}
	}

	for _, child := range recordChildren {
		children := item["netflowMonitorPol"].Children
		children = append(children, map[string]Item{"netflowRsMonitorToRecord": {Attributes: child}})
		item["netflowMonitorPol"] = Item{
			Attributes: item["netflowMonitorPol"].Attributes,
			Children:   children,
		}
	}

	var planAnnotations []provider.TagAnnotationNetflowMonitorPolResourceModel
	if annotations, exists := attributes["annotations"].([]interface{}); exists {
		for _, annotation := range annotations {
			annotationMap := annotation.(map[string]interface{})
			planAnnotations = append(planAnnotations, provider.TagAnnotationNetflowMonitorPolResourceModel{
				Key:   types.StringValue(annotationMap["key"].(string)),
				Value: types.StringValue(annotationMap["value"].(string)),
			})
		}
	}

	var stateAnnotations []provider.TagAnnotationNetflowMonitorPolResourceModel
	tagAnnotationChildren := provider.GetNetflowMonitorPolTagAnnotationChildPayloads(ctx, &diags, data, planAnnotations, stateAnnotations)

	for _, child := range tagAnnotationChildren {
		children := item["netflowMonitorPol"].Children
		children = append(children, map[string]Item{"tagAnnotation": {Attributes: child}})
		item["netflowMonitorPol"] = Item{
			Attributes: item["netflowMonitorPol"].Attributes,
			Children:   children,
		}
	}

	var planTags []provider.TagTagNetflowMonitorPolResourceModel
	if tags, exists := attributes["tags"].([]interface{}); exists {
		for _, tag := range tags {
			tagMap := tag.(map[string]interface{})
			planTags = append(planTags, provider.TagTagNetflowMonitorPolResourceModel{
				Key:   types.StringValue(tagMap["key"].(string)),
				Value: types.StringValue(tagMap["value"].(string)),
			})
		}
	}

	var stateTags []provider.TagTagNetflowMonitorPolResourceModel
	tagTagChildren := provider.GetNetflowMonitorPolTagTagChildPayloads(ctx, &diags, data, planTags, stateTags)
	for _, child := range tagTagChildren {
		children := item["netflowMonitorPol"].Children
		children = append(children, map[string]Item{"tagTag": {Attributes: child}})
		item["netflowMonitorPol"] = Item{
			Attributes: item["netflowMonitorPol"].Attributes,
			Children:   children,
		}
	}
	return item
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

func removeChildAttributes(attributes map[string]interface{}, childKeys []string) {
	for _, key := range childKeys {
		delete(attributes, key)
	}
}

func constructTree(itemList []map[string]Item) []map[string]Item {
	dnMap := make(map[string]*Item)
	var root []map[string]Item

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

func removeDuplicateHeaders(item map[string]Item) map[string]Item {
	cleanedItem := make(map[string]Item)
	for key, val := range item {
		if attributes, exists := val.Attributes["attributes"]; exists {
			if nestedAttributes, ok := attributes.(map[string]interface{}); ok {
				val.Attributes = nestedAttributes
			}
		}
		if val.Children != nil {
			for i, child := range val.Children {
				for childKey, childVal := range child {
					val.Children[i][childKey] = removeDuplicateHeaders(map[string]Item{childKey: childVal})[childKey]
				}
			}
		}
		cleanedItem[key] = val
	}
	return cleanedItem
}

func cleanFinalTree(tree []map[string]Item) []map[string]Item {
	for i, item := range tree {
		for key, val := range item {
			tree[i][key] = removeDuplicateHeaders(map[string]Item{key: val})[key]
		}
	}
	return tree
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

func removeNestedObjectsWithoutResources(data []map[string]Item) {
	for i := len(data) - 1; i >= 0; i-- {
		for _, resource := range data[i] {
			if resource.Children != nil {
				removeNestedObjectsWithoutResources(resource.Children)
				if len(resource.Children) == 0 {
					data = append(data[:i], data[i+1:]...)
				}
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

	for attempts := 0; attempts < 10; attempts++ {
		tree := constructTree(itemList)
		finalTree := replaceDnWithResourceType(tree)
		//removeDnAndResourceTypeFromAttributes(finalTree)
		//removeParentDnFromAttributes(finalTree)
		//removeNestedObjectsWithoutResources(finalTree)

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

			output = unwrapOuterArray(finalTree)

			err = outputToFile(outputFile, output)
			if err != nil {
				log.Fatalf("Error writing output file: %v", err)
			}

			fmt.Printf("ACI Payload written to %s\n", outputFile)

			return

		}

	}
	log.Fatalf("Failed to create matching item list and final tree after 10 attempts")
}
