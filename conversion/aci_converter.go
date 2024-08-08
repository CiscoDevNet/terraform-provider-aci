package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Represents parts of Terraform Plan we are interested in
type Plan struct {
	PlannedValues struct {
		RootModule struct {
			Resources []Resource `json:"resources"`
		} `json:"root_module"`
	} `json:"planned_values"`
	Changes []Change `json:"resource_changes"`
}

type Resource struct {
	Type   string                 `json:"type"`
	Name   string                 `json:"name"`
	Values map[string]interface{} `json:"values"`
}

type Change struct {
	Type   string `json:"type"`
	Change struct {
		Actions []string               `json:"actions"`
		Before  map[string]interface{} `json:"before"`
	} `json:"change"`
}

type createFunc func(map[string]interface{}) map[string]interface{}

// Map to query create() functions
var resourceMap = map[string]createFunc{
	"aci_endpoint_tag_ip":                              createEndpointTagIP,
	"aci_vrf_fallback_route_group":                     createVrfFallbackRouteGroup,
	"aci_pim_route_map_policy":                         createPimRouteMapPol,
	"aci_l3out_consumer_label":                         createL3extConsLbl,
	"aci_external_management_network_instance_profile": createMgmtInstP,
	"aci_external_management_subnet":                   createMgmtSubnet,
	"aci_endpoint_tag_mac":                             createFvEpMacTag,
	"aci_netflow_monitor_policy":                       createNetflowMonitorPol,
	"aci_l3out_redistribute_policy":                    createL3extRsRedistributePol,
	"aci_out_of_band_contract":                         createVzOOBBrCP,
	"aci_relation_to_consumed_out_of_band_contract":    createMgmtRsOoBCons,
	"aci_relation_to_netflow_exporter":                 createNetflowRsMonitorToExporter,
	"aci_pim_route_map_entry":                          createPimRouteMapEntry,
	"aci_l3out_node_sid_profile":                       createMplsNodeSidP,
}

func resourceTypeToMapKey(resourceType string) string {
	switch resourceType {
	case "aci_endpoint_tag_ip":
		return "fvEpIpTag"
	case "aci_external_management_network_instance_profile":
		return "mgmtInstP"
	case "aci_l3out_consumer_label":
		return " l3extConsLbl"
	case "aci_vrf_fallback_route_group":
		return "fvFBRGroup"
	case "aci_pim_route_map_policy":
		return "pimRouteMapPol"
	case "aci_netflow_monitor_pol":
		return "netflowMonitorPol"
	case "aci_external_management_subnet":
		return "mgmtSubnet"
	case "aci_l3out_redistribute_policy":
		return " l3extRsRedistributePol"
	case "aci_out_of_band_contract":
		return "VzOOBBrCP"
	case "aci_relation_to_consumed_out_of_band_contract":
		return "mgmtRsOoBCons"
	case "aci_relation_to_netflow_exporter":
		return "netflowRsMonitorToExporter"
	case "aci_l3out_node_sid_profile":
		return "mplsNodeSidP"

	default:
		return resourceType
	}
}

// Executes Terraform commands to create input JSON
func runTerraform() (string, error) {
	planBin := "plan.bin"
	planJson := "plan.json"

	if err := exec.Command("terraform", "plan", "-out="+planBin).Run(); err != nil {
		return "", fmt.Errorf("failed to run terraform plan: %w", err)
	}

	output, err := os.Create(planJson)
	if err != nil {
		return "", fmt.Errorf("failed to create json file: %w", err)
	}
	defer output.Close()

	cmdShow := exec.Command("terraform", "show", "-json", planBin)
	cmdShow.Stdout = output
	if err := cmdShow.Run(); err != nil {
		return "", fmt.Errorf("failed to run terraform show: %w", err)
	}

	return planJson, nil
}

// Converts Plan json into bytes, then unmarshalls into Plan struct
func readPlan(jsonFile string) (Plan, error) {
	var plan Plan
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return plan, fmt.Errorf("failed to read input file: %w", err)
	}

	if err := json.Unmarshal(data, &plan); err != nil {
		return plan, fmt.Errorf("failed to parse input file: %w", err)
	}

	return plan, nil
}

func writeToFile(outputFile string, data map[string]interface{}) error {
	var dataList []map[string]interface{}
	for key, value := range data {
		dataList = append(dataList, map[string]interface{}{key: value})
	}

	traverseAndRemoveAttributes(dataList)

	updatedData := make(map[string]interface{})
	for _, item := range dataList {
		for key, value := range item {
			updatedData[key] = value
		}
	}

	outputData, err := json.MarshalIndent(updatedData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to convert data to JSON: %w", err)
	}

	if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func createItem(resourceType string, values map[string]interface{}, status string) map[string]interface{} {
	if create, exists := resourceMap[resourceType]; exists {
		item := create(values)
		if status == "deleted" && item != nil {
			if attributes, ok := item[resourceTypeToMapKey(resourceType)].(map[string]interface{}); ok {
				if attrs, ok := attributes["attributes"].(map[string]interface{}); ok {
					if status != "" {
						attrs["status"] = status
					}
				}
			}
		}
		return item
	}
	return nil
}

// Ranges through resource created/deleted and resources, creates each item, adds to []map[string]interface{}
func createItemList(plan Plan) []map[string]interface{} {
	var data []map[string]interface{}

	// Create deleted items
	for _, change := range plan.Changes {
		if len(change.Change.Actions) > 0 && change.Change.Actions[0] == "delete" {
			item := createItem(change.Type, change.Change.Before, "deleted")
			if item != nil {
				data = append(data, item)
			}
		}
	}

	// Create-created items and include parent_dn if exists for tree construction
	for _, resource := range plan.PlannedValues.RootModule.Resources {
		item := createItem(resource.Type, resource.Values, "")
		if item != nil {
			for _, value := range item {
				if obj, ok := value.(map[string]interface{}); ok {
					if attributes, ok := obj["attributes"].(map[string]interface{}); ok {
						if parentDn, ok := resource.Values["parent_dn"].(string); ok && parentDn != "" {
							attributes["parent_dn"] = parentDn
						}
					}
				}
			}
			data = append(data, item)
		}
	}

	return data
}

func constructTree(data []map[string]interface{}) []map[string]interface{} {
	resourceMap := make(map[string]map[string]interface{})
	parentChildMap := make(map[string][]map[string]interface{})

	for _, item := range data {
		for key, value := range item {
			if obj, ok := value.(map[string]interface{}); ok {
				if attributes, ok := obj["attributes"].(map[string]interface{}); ok {
					if dn, ok := attributes["dn"].(string); ok {
						resourceMap[dn] = map[string]interface{}{key: obj}
						if parentDn, ok := attributes["parent_dn"].(string); ok && parentDn != "" {
							parentChildMap[parentDn] = append(parentChildMap[parentDn], map[string]interface{}{key: obj})
						}
					}
				}
			}
		}
	}

	for parentDn, children := range parentChildMap {
		if parent, exists := resourceMap[parentDn]; exists {
			attachChildren(parent, children)
		} else {
			// If parent doesn't exist, create placeholder item
			missingParent := map[string]interface{}{
				"unknownParent": map[string]interface{}{
					"attributes": map[string]interface{}{
						"dn": parentDn,
					},
					"children": children,
				},
			}
			resourceMap[parentDn] = missingParent
		}
	}

	var tree []map[string]interface{}
	for _, resource := range resourceMap {
		for key, obj := range resource {
			if key == "unknownParent" {
				children := obj.(map[string]interface{})["children"].([]map[string]interface{})
				tree = append(tree, children...)
			} else {
				if attributes, ok := obj.(map[string]interface{})["attributes"].(map[string]interface{}); ok {
					if parentDn, ok := attributes["parent_dn"].(string); !ok || parentDn == "" {
						tree = append(tree, resource)
					}
				}
			}
		}
	}

	return tree
}

func attachChildren(parent map[string]interface{}, children []map[string]interface{}) {
	for _, value := range parent {
		if obj, ok := value.(map[string]interface{}); ok {
			if objChildren, ok := obj["children"].([]map[string]interface{}); ok {
				obj["children"] = append(objChildren, children...)
			} else {
				obj["children"] = children
			}
		}
	}
}

func removeAttributes(obj map[string]interface{}) {
	if attributes, ok := obj["attributes"].(map[string]interface{}); ok {
		delete(attributes, "dn")
		delete(attributes, "parent_dn")
	}
	if children, ok := obj["children"].([]map[string]interface{}); ok {
		for _, child := range children {
			for _, childObj := range child {
				removeAttributes(childObj.(map[string]interface{}))
			}
		}
	}
}

func traverseAndRemoveAttributes(data []map[string]interface{}) {
	for _, item := range data {
		for _, value := range item {
			removeAttributes(value.(map[string]interface{}))
		}
	}
}

func removeArrayBrackets(jsonData []byte) ([]byte, error) {
	var temp interface{}
	if err := json.Unmarshal(jsonData, &temp); err != nil {
		return nil, err
	}

	jsonMap, ok := temp.([]interface{})
	if !ok || len(jsonMap) == 0 {
		return nil, fmt.Errorf("invalid JSON structure")
	}

	var finalJSON []byte
	for _, element := range jsonMap {
		elementJSON, err := json.MarshalIndent(element, "", "  ")
		if err != nil {
			return nil, err
		}
		finalJSON = append(finalJSON, elementJSON...)
		finalJSON = append(finalJSON, '\n')
	}

	return finalJSON, nil
}

func parseBinToMap(jsonPayload []byte) (map[string]interface{}, error) {
	var customData map[string]interface{}
	err := json.Unmarshal(jsonPayload, &customData)
	if err != nil {
		return nil, err
	}
	return customData, nil
}

func createEndpointTagIP(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvEpIpTagResourceModel{}

	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["vrf_name"].(string); ok && v != "" {
		data.CtxName = types.StringValue(v)
	}
	if v, ok := attributes["id_attribute"].(string); ok && v != "" {
		data.FvEpIpTagId = types.StringValue(v)
	}
	if v, ok := attributes["ip"].(string); ok && v != "" {
		data.Ip = types.StringValue(v)
	}
	if v, ok := attributes["name"].(string); ok && v != "" {
		data.Name = types.StringValue(v)
	}
	if v, ok := attributes["name_alias"].(string); ok && v != "" {
		data.NameAlias = types.StringValue(v)
	}

	planAnnotations := convertToEpIpTagAnnotations(attributes["annotations"])
	planTags := convertToEpIpTagTags(attributes["tags"])

	stateAnnotations := planAnnotations
	stateTags := planTags

	newEndpointTag := provider.GetFvEpIpTagCreateJsonPayload(ctx, &diags, data, planAnnotations, stateAnnotations, planTags, stateTags)

	payload_bin := newEndpointTag.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(payload_bin)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFvEpIpTagId(ctx, data)
	attrs := payload["fvEpIpTag"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if v, ok := attributes["status"].(string); ok && v != "" {
		attrs["status"] = v
	}

	return payload
}

func createVrfFallbackRouteGroup(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvFBRGroupResourceModel{}

	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["name"].(string); ok && v != "" {
		data.Name = types.StringValue(v)
	}
	if v, ok := attributes["description"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["name_alias"].(string); ok && v != "" {
		data.NameAlias = types.StringValue(v)
	}

	planMembers := convertToFBRMembers(attributes["vrf_fallback_route_group_members"])
	planAnnotations := convertToFvFBRGroupAnnotations(attributes["annotations"])
	planTags := convertToFvFBRGroupTags(attributes["tags"])

	newAciVrf := provider.GetFvFBRGroupCreateJsonPayload(ctx, &diags, data, planMembers, planMembers, planAnnotations, planAnnotations, planTags, planTags)

	jsonPayload := newAciVrf.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFvFBRGroupId(ctx, data)
	attrs := payload["fvFBRGroup"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if v, ok := attributes["status"].(string); ok && v != "" {
		attrs["status"] = v
	}

	return payload
}

func convertToFvFBRGroupAnnotations(annotations interface{}) []provider.TagAnnotationFvFBRGroupResourceModel {
	var planAnnotations []provider.TagAnnotationFvFBRGroupResourceModel
	if annotations, ok := annotations.([]interface{}); ok {
		for _, annotation := range annotations {
			annotationMap := annotation.(map[string]interface{})
			planAnnotations = append(planAnnotations, provider.TagAnnotationFvFBRGroupResourceModel{
				Key:   types.StringValue(annotationMap["key"].(string)),
				Value: types.StringValue(annotationMap["value"].(string)),
			})
		}
	}
	return planAnnotations
}

func convertToFvFBRGroupTags(tags interface{}) []provider.TagTagFvFBRGroupResourceModel {
	var planTags []provider.TagTagFvFBRGroupResourceModel
	if tags, ok := tags.([]interface{}); ok {
		for _, tag := range tags {
			tagMap := tag.(map[string]interface{})
			planTags = append(planTags, provider.TagTagFvFBRGroupResourceModel{
				Key:   types.StringValue(tagMap["key"].(string)),
				Value: types.StringValue(tagMap["value"].(string)),
			})
		}
	}
	return planTags
}

func convertToFBRMembers(members interface{}) []provider.FvFBRMemberFvFBRGroupResourceModel {
	var planMembers []provider.FvFBRMemberFvFBRGroupResourceModel
	if members, ok := members.([]interface{}); ok {
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
	return planMembers
}

func convertToEpIpTagAnnotations(annotations interface{}) []provider.TagAnnotationFvEpIpTagResourceModel {
	var planAnnotations []provider.TagAnnotationFvEpIpTagResourceModel
	if annotations, ok := annotations.([]interface{}); ok {
		for _, annotation := range annotations {
			annotationMap := annotation.(map[string]interface{})
			planAnnotations = append(planAnnotations, provider.TagAnnotationFvEpIpTagResourceModel{
				Key:   types.StringValue(annotationMap["key"].(string)),
				Value: types.StringValue(annotationMap["value"].(string)),
			})
		}
	}
	return planAnnotations
}

func convertToEpIpTagTags(tags interface{}) []provider.TagTagFvEpIpTagResourceModel {
	var planTags []provider.TagTagFvEpIpTagResourceModel
	if tags, ok := tags.([]interface{}); ok {
		for _, tag := range tags {
			tagMap := tag.(map[string]interface{})
			planTags = append(planTags, provider.TagTagFvEpIpTagResourceModel{
				Key:   types.StringValue(tagMap["key"].(string)),
				Value: types.StringValue(tagMap["value"].(string)),
			})
		}
	}
	return planTags
}

func convertToMgmtInstPAnnotations(annotations interface{}) []provider.TagAnnotationMgmtInstPResourceModel {
	var planAnnotations []provider.TagAnnotationMgmtInstPResourceModel
	if annotations, ok := annotations.([]interface{}); ok {
		for _, annotation := range annotations {
			annotationMap := annotation.(map[string]interface{})
			planAnnotations = append(planAnnotations, provider.TagAnnotationMgmtInstPResourceModel{
				Key:   types.StringValue(annotationMap["key"].(string)),
				Value: types.StringValue(annotationMap["value"].(string)),
			})
		}
	}
	return planAnnotations
}

func convertToMgmtInstPTags(tags interface{}) []provider.TagTagMgmtInstPResourceModel {
	var planTags []provider.TagTagMgmtInstPResourceModel
	if tags, ok := tags.([]interface{}); ok {
		for _, tag := range tags {
			tagMap := tag.(map[string]interface{})
			planTags = append(planTags, provider.TagTagMgmtInstPResourceModel{
				Key:   types.StringValue(tagMap["key"].(string)),
				Value: types.StringValue(tagMap["value"].(string)),
			})
		}
	}
	return planTags
}

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
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetTagAnnotationId(ctx, data)
	attrs := payload["tagAnnotation"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}

	return payload
}

func createFvEpIpTag(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvEpIpTagResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["ctx_name"].(string); ok && v != "" {
		data.CtxName = types.StringValue(v)
	}
	if v, ok := attributes["id"].(string); ok && v != "" {
		data.FvEpIpTagId = types.StringValue(v)
	}
	if v, ok := attributes["ip"].(string); ok && v != "" {
		data.Ip = types.StringValue(v)
	}
	if v, ok := attributes["name"].(string); ok && v != "" {
		data.Name = types.StringValue(v)
	}
	if v, ok := attributes["name_alias"].(string); ok && v != "" {
		data.NameAlias = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvEpIpTag(attributes["annotations"])
	planTagTag := convertToTagTagFvEpIpTag(attributes["tags"])

	newAciFvEpIpTag := provider.GetFvEpIpTagCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvEpIpTag.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFvEpIpTagId(ctx, data)
	attrs := payload["fvEpIpTag"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}

	return payload
}
func convertToTagAnnotationFvEpIpTag(resources interface{}) []provider.TagAnnotationFvEpIpTagResourceModel {
	var planResources []provider.TagAnnotationFvEpIpTagResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvEpIpTagResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvEpIpTag(resources interface{}) []provider.TagTagFvEpIpTagResourceModel {
	var planResources []provider.TagTagFvEpIpTagResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvEpIpTagResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

func createFvEpMacTag(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvEpMacTagResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["bd_name"].(string); ok && v != "" {
		data.BdName = types.StringValue(v)
	}
	if v, ok := attributes["id"].(string); ok && v != "" {
		data.FvEpMacTagId = types.StringValue(v)
	}
	if v, ok := attributes["mac"].(string); ok && v != "" {
		data.Mac = types.StringValue(v)
	}
	if v, ok := attributes["name"].(string); ok && v != "" {
		data.Name = types.StringValue(v)
	}
	if v, ok := attributes["name_alias"].(string); ok && v != "" {
		data.NameAlias = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvEpMacTag(attributes["annotations"])
	planTagTag := convertToTagTagFvEpMacTag(attributes["tags"])

	newAciFvEpMacTag := provider.GetFvEpMacTagCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvEpMacTag.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFvEpMacTagId(ctx, data)
	attrs := payload["fvEpMacTag"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}
	return payload
}
func convertToTagAnnotationFvEpMacTag(resources interface{}) []provider.TagAnnotationFvEpMacTagResourceModel {
	var planResources []provider.TagAnnotationFvEpMacTagResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvEpMacTagResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvEpMacTag(resources interface{}) []provider.TagTagFvEpMacTagResourceModel {
	var planResources []provider.TagTagFvEpMacTagResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvEpMacTagResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

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
	planMgmtRsOoBCons := convertToMgmtRsOoBConsMgmtInstP(attributes["out_of_band_contract_name"])
	planTagAnnotation := convertToTagAnnotationMgmtInstP(attributes["annotations"])
	planTagTag := convertToTagTagMgmtInstP(attributes["tags"])

	newAciMgmtInstP := provider.GetMgmtInstPCreateJsonPayload(ctx, &diags, data, planMgmtRsOoBCons, planMgmtRsOoBCons, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciMgmtInstP.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetMgmtInstPId(ctx, data)
	attrs := payload["mgmtInstP"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
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

func createMgmtSubnet(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.MgmtSubnetResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["descr"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["ip"].(string); ok && v != "" {
		data.Ip = types.StringValue(v)
	}
	if v, ok := attributes["name"].(string); ok && v != "" {
		data.Name = types.StringValue(v)
	}
	if v, ok := attributes["name_alias"].(string); ok && v != "" {
		data.NameAlias = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationMgmtSubnet(attributes["annotations"])
	planTagTag := convertToTagTagMgmtSubnet(attributes["tags"])

	newAciMgmtSubnet := provider.GetMgmtSubnetCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciMgmtSubnet.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetMgmtSubnetId(ctx, data)
	attrs := payload["mgmtSubnet"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}
	return payload
}
func convertToTagAnnotationMgmtSubnet(resources interface{}) []provider.TagAnnotationMgmtSubnetResourceModel {
	var planResources []provider.TagAnnotationMgmtSubnetResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationMgmtSubnetResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagMgmtSubnet(resources interface{}) []provider.TagTagMgmtSubnetResourceModel {
	var planResources []provider.TagTagMgmtSubnetResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagMgmtSubnetResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

func createL3extConsLbl(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.L3extConsLblResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
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
	if v, ok := attributes["owner"].(string); ok && v != "" {
		data.Owner = types.StringValue(v)
	}
	if v, ok := attributes["owner_key"].(string); ok && v != "" {
		data.OwnerKey = types.StringValue(v)
	}
	if v, ok := attributes["owner_tag"].(string); ok && v != "" {
		data.OwnerTag = types.StringValue(v)
	}
	if v, ok := attributes["tag"].(string); ok && v != "" {
		data.Tag = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationL3extConsLbl(attributes["annotations"])
	planTagTag := convertToTagTagL3extConsLbl(attributes["tags"])

	newAciL3extConsLbl := provider.GetL3extConsLblCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciL3extConsLbl.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetL3extConsLblId(ctx, data)
	attrs := payload["l3extConsLbl"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}

	return payload
}
func convertToTagAnnotationL3extConsLbl(resources interface{}) []provider.TagAnnotationL3extConsLblResourceModel {
	var planResources []provider.TagAnnotationL3extConsLblResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationL3extConsLblResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagL3extConsLbl(resources interface{}) []provider.TagTagL3extConsLblResourceModel {
	var planResources []provider.TagTagL3extConsLblResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagL3extConsLblResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

func createMplsNodeSidP(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.MplsNodeSidPResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["descr"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["loopback_addr"].(string); ok && v != "" {
		data.LoopbackAddr = types.StringValue(v)
	}
	if v, ok := attributes["name"].(string); ok && v != "" {
		data.Name = types.StringValue(v)
	}
	if v, ok := attributes["name_alias"].(string); ok && v != "" {
		data.NameAlias = types.StringValue(v)
	}
	if v, ok := attributes["sidoffset"].(string); ok && v != "" {
		data.Sidoffset = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationMplsNodeSidP(attributes["annotations"])
	planTagTag := convertToTagTagMplsNodeSidP(attributes["tags"])

	newAciMplsNodeSidP := provider.GetMplsNodeSidPCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciMplsNodeSidP.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetMplsNodeSidPId(ctx, data)
	attrs := payload["mplsNodeSidP"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}
	return payload
}
func convertToTagAnnotationMplsNodeSidP(resources interface{}) []provider.TagAnnotationMplsNodeSidPResourceModel {
	var planResources []provider.TagAnnotationMplsNodeSidPResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationMplsNodeSidPResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagMplsNodeSidP(resources interface{}) []provider.TagTagMplsNodeSidPResourceModel {
	var planResources []provider.TagTagMplsNodeSidPResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagMplsNodeSidPResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

func createL3extProvLbl(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.L3extProvLblResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
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
	if v, ok := attributes["owner_key"].(string); ok && v != "" {
		data.OwnerKey = types.StringValue(v)
	}
	if v, ok := attributes["owner_tag"].(string); ok && v != "" {
		data.OwnerTag = types.StringValue(v)
	}
	if v, ok := attributes["tag"].(string); ok && v != "" {
		data.Tag = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationL3extProvLbl(attributes["annotations"])
	planTagTag := convertToTagTagL3extProvLbl(attributes["tags"])

	newAciL3extProvLbl := provider.GetL3extProvLblCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciL3extProvLbl.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetL3extProvLblId(ctx, data)
	attrs := payload["l3extProvLbl"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}
	return payload
}
func convertToTagAnnotationL3extProvLbl(resources interface{}) []provider.TagAnnotationL3extProvLblResourceModel {
	var planResources []provider.TagAnnotationL3extProvLblResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationL3extProvLblResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagL3extProvLbl(resources interface{}) []provider.TagTagL3extProvLblResourceModel {
	var planResources []provider.TagTagL3extProvLblResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagL3extProvLblResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

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
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetL3extRsRedistributePolId(ctx, data)
	attrs := payload["l3extRsRedistributePol"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
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

func createNetflowMonitorPol(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.NetflowMonitorPolResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
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
	if v, ok := attributes["owner_key"].(string); ok && v != "" {
		data.OwnerKey = types.StringValue(v)
	}
	if v, ok := attributes["owner_tag"].(string); ok && v != "" {
		data.OwnerTag = types.StringValue(v)
	}
	planNetflowRsMonitorToExporter := convertToNetflowRsMonitorToExporterNetflowMonitorPol(attributes["relation_to_netflow_exporters"])
	planNetflowRsMonitorToRecord := convertToNetflowRsMonitorToRecordNetflowMonitorPol(attributes["relation_to_netflow_record"])
	planTagAnnotation := convertToTagAnnotationNetflowMonitorPol(attributes["annotations"])
	planTagTag := convertToTagTagNetflowMonitorPol(attributes["tags"])

	newAciNetflowMonitorPol := provider.GetNetflowMonitorPolCreateJsonPayload(ctx, &diags, data, planNetflowRsMonitorToExporter, planNetflowRsMonitorToExporter, planNetflowRsMonitorToRecord, planNetflowRsMonitorToRecord, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciNetflowMonitorPol.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetNetflowMonitorPolId(ctx, data)
	attrs := payload["netflowMonitorPol"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if status, ok := attributes["parent_dn"].(string); ok && status != "" {
		attrs["parent_dn"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}

	return payload
}

func convertToNetflowRsMonitorToExporterNetflowMonitorPol(resources interface{}) []provider.NetflowRsMonitorToExporterNetflowMonitorPolResourceModel {
	var planResources []provider.NetflowRsMonitorToExporterNetflowMonitorPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.NetflowRsMonitorToExporterNetflowMonitorPolResourceModel{
				Annotation:               types.StringValue(resourceMap["annotation"].(string)),
				TnNetflowExporterPolName: types.StringValue(resourceMap["netflow_exporter_policy_name"].(string)),
			})
		}
	}
	return planResources
}

func convertToNetflowRsMonitorToRecordNetflowMonitorPol(resources interface{}) []provider.NetflowRsMonitorToRecordNetflowMonitorPolResourceModel {
	var planResources []provider.NetflowRsMonitorToRecordNetflowMonitorPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.NetflowRsMonitorToRecordNetflowMonitorPolResourceModel{
				Annotation:             types.StringValue(resourceMap["annotation"].(string)),
				TnNetflowRecordPolName: types.StringValue(resourceMap["netflow_record_policy_name"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagAnnotationNetflowMonitorPol(resources interface{}) []provider.TagAnnotationNetflowMonitorPolResourceModel {
	var planResources []provider.TagAnnotationNetflowMonitorPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationNetflowMonitorPolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagNetflowMonitorPol(resources interface{}) []provider.TagTagNetflowMonitorPolResourceModel {
	var planResources []provider.TagTagNetflowMonitorPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagNetflowMonitorPolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

func createVzOOBBrCP(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.VzOOBBrCPResourceModel{}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["descr"].(string); ok && v != "" {
		data.Descr = types.StringValue(v)
	}
	if v, ok := attributes["intent"].(string); ok && v != "" {
		data.Intent = types.StringValue(v)
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
	if v, ok := attributes["prio"].(string); ok && v != "" {
		data.Prio = types.StringValue(v)
	}
	if v, ok := attributes["scope"].(string); ok && v != "" {
		data.Scope = types.StringValue(v)
	}
	if v, ok := attributes["target_dscp"].(string); ok && v != "" {
		data.TargetDscp = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationVzOOBBrCP(attributes["annotations"])
	planTagTag := convertToTagTagVzOOBBrCP(attributes["tags"])

	newAciVzOOBBrCP := provider.GetVzOOBBrCPCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciVzOOBBrCP.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetVzOOBBrCPId(ctx, data)
	attrs := payload["vzOOBBrCP"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if status, ok := attributes["parent_dn"].(string); ok && status != "" {
		attrs["parent_dn"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}

	return payload
}
func convertToTagAnnotationVzOOBBrCP(resources interface{}) []provider.TagAnnotationVzOOBBrCPResourceModel {
	var planResources []provider.TagAnnotationVzOOBBrCPResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationVzOOBBrCPResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagVzOOBBrCP(resources interface{}) []provider.TagTagVzOOBBrCPResourceModel {
	var planResources []provider.TagTagVzOOBBrCPResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagVzOOBBrCPResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

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
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetPimRouteMapEntryId(ctx, data)
	attrs := payload["pimRouteMapEntry"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
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

func createPimRouteMapPol(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.PimRouteMapPolResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
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
	if v, ok := attributes["owner_key"].(string); ok && v != "" {
		data.OwnerKey = types.StringValue(v)
	}
	if v, ok := attributes["owner_tag"].(string); ok && v != "" {
		data.OwnerTag = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationPimRouteMapPol(attributes["annotations"])
	planTagTag := convertToTagTagPimRouteMapPol(attributes["tags"])

	newAciPimRouteMapPol := provider.GetPimRouteMapPolCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciPimRouteMapPol.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetPimRouteMapPolId(ctx, data)
	attrs := payload["pimRouteMapPol"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}

	return payload
}
func convertToTagAnnotationPimRouteMapPol(resources interface{}) []provider.TagAnnotationPimRouteMapPolResourceModel {
	var planResources []provider.TagAnnotationPimRouteMapPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationPimRouteMapPolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagPimRouteMapPol(resources interface{}) []provider.TagTagPimRouteMapPolResourceModel {
	var planResources []provider.TagTagPimRouteMapPolResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagPimRouteMapPolResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

func createMgmtRsOoBCons(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.MgmtRsOoBConsResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["prio"].(string); ok && v != "" {
		data.Prio = types.StringValue(v)
	}
	if v, ok := attributes["tn_vz_oob_br_cp_name"].(string); ok && v != "" {
		data.TnVzOOBBrCPName = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationMgmtRsOoBCons(attributes["annotations"])
	planTagTag := convertToTagTagMgmtRsOoBCons(attributes["tags"])

	newAciMgmtRsOoBCons := provider.GetMgmtRsOoBConsCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciMgmtRsOoBCons.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetMgmtRsOoBConsId(ctx, data)
	attrs := payload["mgmtRsOoBCons"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}

	return payload
}
func convertToTagAnnotationMgmtRsOoBCons(resources interface{}) []provider.TagAnnotationMgmtRsOoBConsResourceModel {
	var planResources []provider.TagAnnotationMgmtRsOoBConsResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationMgmtRsOoBConsResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagMgmtRsOoBCons(resources interface{}) []provider.TagTagMgmtRsOoBConsResourceModel {
	var planResources []provider.TagTagMgmtRsOoBConsResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagMgmtRsOoBConsResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

func createL3extRsOutToFBRGroup(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.L3extRsOutToFBRGroupResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["t_dn"].(string); ok && v != "" {
		data.TDn = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationL3extRsOutToFBRGroup(attributes["annotations"])
	planTagTag := convertToTagTagL3extRsOutToFBRGroup(attributes["tags"])

	newAciL3extRsOutToFBRGroup := provider.GetL3extRsOutToFBRGroupCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciL3extRsOutToFBRGroup.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetL3extRsOutToFBRGroupId(ctx, data)
	attrs := payload["l3extRsOutToFBRGroup"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}

	return payload
}
func convertToTagAnnotationL3extRsOutToFBRGroup(resources interface{}) []provider.TagAnnotationL3extRsOutToFBRGroupResourceModel {
	var planResources []provider.TagAnnotationL3extRsOutToFBRGroupResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationL3extRsOutToFBRGroupResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagL3extRsOutToFBRGroup(resources interface{}) []provider.TagTagL3extRsOutToFBRGroupResourceModel {
	var planResources []provider.TagTagL3extRsOutToFBRGroupResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagL3extRsOutToFBRGroupResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

func createNetflowRsMonitorToExporter(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.NetflowRsMonitorToExporterResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["annotation"].(string); ok && v != "" {
		data.Annotation = types.StringValue(v)
	}
	if v, ok := attributes["tn_netflow_exporter_pol_name"].(string); ok && v != "" {
		data.TnNetflowExporterPolName = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationNetflowRsMonitorToExporter(attributes["annotations"])
	planTagTag := convertToTagTagNetflowRsMonitorToExporter(attributes["tags"])

	newAciNetflowRsMonitorToExporter := provider.GetNetflowRsMonitorToExporterCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciNetflowRsMonitorToExporter.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetNetflowRsMonitorToExporterId(ctx, data)
	attrs := payload["netflowRsMonitorToExporter"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}

	return payload
}
func convertToTagAnnotationNetflowRsMonitorToExporter(resources interface{}) []provider.TagAnnotationNetflowRsMonitorToExporterResourceModel {
	var planResources []provider.TagAnnotationNetflowRsMonitorToExporterResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationNetflowRsMonitorToExporterResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagNetflowRsMonitorToExporter(resources interface{}) []provider.TagTagNetflowRsMonitorToExporterResourceModel {
	var planResources []provider.TagTagNetflowRsMonitorToExporterResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagNetflowRsMonitorToExporterResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

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
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetTagTagId(ctx, data)
	attrs := payload["tagTag"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}
	return payload
}

func createFvFBRMember(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvFBRMemberResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
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
	if v, ok := attributes["rnh_addr"].(string); ok && v != "" {
		data.RnhAddr = types.StringValue(v)
	}
	planTagAnnotation := convertToTagAnnotationFvFBRMember(attributes["annotations"])
	planTagTag := convertToTagTagFvFBRMember(attributes["tags"])

	newAciFvFBRMember := provider.GetFvFBRMemberCreateJsonPayload(ctx, &diags, data, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvFBRMember.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFvFBRMemberId(ctx, data)
	attrs := payload["fvFBRMember"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}
	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}
	return payload
}
func convertToTagAnnotationFvFBRMember(resources interface{}) []provider.TagAnnotationFvFBRMemberResourceModel {
	var planResources []provider.TagAnnotationFvFBRMemberResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvFBRMemberResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvFBRMember(resources interface{}) []provider.TagTagFvFBRMemberResourceModel {
	var planResources []provider.TagTagFvFBRMemberResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvFBRMemberResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

func createFvFBRGroup(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvFBRGroupResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
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
	planFvFBRMember := convertToFvFBRMemberFvFBRGroup(attributes["vrf_fallback_route_group_members"])
	planTagAnnotation := convertToTagAnnotationFvFBRGroup(attributes["annotations"])
	planTagTag := convertToTagTagFvFBRGroup(attributes["tags"])

	newAciFvFBRGroup := provider.GetFvFBRGroupCreateJsonPayload(ctx, &diags, data, planFvFBRMember, planFvFBRMember, planTagAnnotation, planTagAnnotation, planTagTag, planTagTag)

	jsonPayload := newAciFvFBRGroup.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseBinToMap(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFvFBRGroupId(ctx, data)
	attrs := payload["fvFBRGroup"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	if parent, ok := attributes["parent_dn"].(string); ok && parent != "" {
		attrs["parent_dn"] = parent
	}
	return payload
}
func convertToFvFBRMemberFvFBRGroup(resources interface{}) []provider.FvFBRMemberFvFBRGroupResourceModel {
	var planResources []provider.FvFBRMemberFvFBRGroupResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.FvFBRMemberFvFBRGroupResourceModel{
				Annotation: types.StringValue(resourceMap["annotation"].(string)),
				Descr:      types.StringValue(resourceMap["descr"].(string)),
				Name:       types.StringValue(resourceMap["name"].(string)),
				NameAlias:  types.StringValue(resourceMap["name_alias"].(string)),
				RnhAddr:    types.StringValue(resourceMap["rnh_addr"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagAnnotationFvFBRGroup(resources interface{}) []provider.TagAnnotationFvFBRGroupResourceModel {
	var planResources []provider.TagAnnotationFvFBRGroupResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagAnnotationFvFBRGroupResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}
func convertToTagTagFvFBRGroup(resources interface{}) []provider.TagTagFvFBRGroupResourceModel {
	var planResources []provider.TagTagFvFBRGroupResourceModel
	if resources, ok := resources.([]interface{}); ok {
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			planResources = append(planResources, provider.TagTagFvFBRGroupResourceModel{
				Key:   types.StringValue(resourceMap["key"].(string)),
				Value: types.StringValue(resourceMap["value"].(string)),
			})
		}
	}
	return planResources
}

func main() {
	if len(os.Args) != 1 {
		fmt.Println("Usage: go run aci_converter.go")
		os.Exit(1)
	}

	outputFile := "payload.json"

	planJSON, err := runTerraform()
	if err != nil {
		log.Fatalf("Error running Terraform: %v", err)
	}

	plan, err := readPlan(planJSON)
	if err != nil {
		log.Fatalf("Error reading plan: %v", err)
	}

	jsonDump := createItemList(plan)

	aciPayload := constructTree(jsonDump)

	aciPayloadMap := make(map[string]interface{})
	for _, item := range aciPayload {
		for key, value := range item {
			aciPayloadMap[key] = value
		}
	}

	err = writeToFile(outputFile, aciPayloadMap)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("ACI Payload written to %s\n", outputFile)
}
