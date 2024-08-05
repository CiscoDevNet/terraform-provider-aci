package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

var resourceMap = map[string]createFunc{
	"aci_endpoint_tag_ip":                              createEndpointTagIP,
	"aci_external_management_network_instance_profile": createNetworkInstanceProfile,
	"aci_vrf_fallback_route_group":                     createVrfFallbackRouteGroup,
	"aci_tenant":                                       createTenant,
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run aci_converter.go <output.json>")
		os.Exit(1)
	}

	outputFile := os.Args[1]

	planJSON, err := runTerraform()
	if err != nil {
		log.Fatalf("Error running Terraform: %v", err)
	}

	plan, err := readPlan(planJSON)
	if err != nil {
		log.Fatalf("Error reading plan: %v", err)
	}

	jsonDump := processResources(plan)

	aciPayload := constructTree(jsonDump)

	err = writeToFile(outputFile, aciPayload)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("ACI Payload written to %s\n", outputFile)
}

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

func writeToFile(outputFile string, data interface{}) error {
	outputData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to convert data to JSON: %w", err)
	}

	if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func processResources(plan Plan) []map[string]interface{} {
	var data []map[string]interface{}

	for _, change := range plan.Changes {
		if len(change.Change.Actions) > 0 && change.Change.Actions[0] == "delete" {
			item := createItem(change.Type, change.Change.Before, "deleted")
			if item != nil {
				data = append(data, item)
			}
		}
	}

	for _, resource := range plan.PlannedValues.RootModule.Resources {
		item := createItem(resource.Type, resource.Values, "")
		if item != nil {
			data = append(data, item)
		}
	}

	return data
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

func resourceTypeToMapKey(resourceType string) string {
	switch resourceType {
	case "aci_tenant":
		return "fvTenant"
	case "aci_endpoint_tag_ip":
		return "fvEpIpTag"
	case "aci_external_management_network_instance_profile":
		return "mgmtInstP"
	case "aci_vrf_fallback_route_group":
		return "fvFBRGroup"
	default:
		return resourceType
	}
}

func constructTree(data []map[string]interface{}) []map[string]interface{} {
	var tree []map[string]interface{}
	for _, item := range data {
		for key, value := range item {
			if obj, ok := value.(map[string]interface{}); ok {
				if attributes, ok := obj["attributes"].(map[string]interface{}); ok {
					if dn, ok := attributes["dn"].(string); ok {
						tree = addToTree(tree, key, obj, dn)
					}
				}
			}
		}
	}
	return tree
}

func addToTree(tree []map[string]interface{}, key string, obj map[string]interface{}, dn string) []map[string]interface{} {
	parts := strings.Split(dn, "/")
	tree = addToBranch(tree, key, obj, parts, 1)
	return tree
}

func addToBranch(branch []map[string]interface{}, key string, obj map[string]interface{}, parts []string, index int) []map[string]interface{} {
	if index >= len(parts) {
		return branch
	}

	parentDn := strings.Join(parts[:index+1], "/")
	parentKey, parentIndex := findParentKey(branch, parentDn)
	if parentKey == "" {
		branch = append(branch, map[string]interface{}{key: obj})
		return branch
	}

	if parent, ok := branch[parentIndex][parentKey].(map[string]interface{}); ok {
		if children, ok := parent["children"].([]map[string]interface{}); ok {
			parent["children"] = append(children, map[string]interface{}{key: obj})
		} else {
			parent["children"] = []map[string]interface{}{map[string]interface{}{key: obj}}
		}
	}
	return branch
}

func findParentKey(branch []map[string]interface{}, dn string) (string, int) {
	for i, item := range branch {
		for key, value := range item {
			if obj, ok := value.(map[string]interface{}); ok {
				if attributes, ok := obj["attributes"].(map[string]interface{}); ok {
					if objDn, ok := attributes["dn"].(string); ok {
						if objDn == dn {
							return key, i
						}
					}
				}
			}
		}
	}
	return "", -1
}

func parseCustomJSON(jsonPayload []byte) (map[string]interface{}, error) {
	var customData map[string]interface{}
	err := json.Unmarshal(jsonPayload, &customData)
	if err != nil {
		return nil, err
	}
	return customData, nil
}

func createTenant(attributes map[string]interface{}) map[string]interface{} {
	attrs := map[string]interface{}{
		"dn":         fmt.Sprintf("uni/tn-%s", attributes["name"]),
		"name":       attributes["name"],
		"descr":      attributes["description"],
		"annotation": attributes["annotation"],
		"nameAlias":  attributes["name_alias"],
	}

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return map[string]interface{}{
		"fvTenant": map[string]interface{}{
			"attributes": attrs,
		},
	}
}

func createEndpointTagIP(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvEpIpTagResourceModel{}

	data.ParentDn = types.StringValue(attributes["parent_dn"].(string))
	data.Annotation = types.StringValue(attributes["annotation"].(string))
	data.CtxName = types.StringValue(attributes["vrf_name"].(string))
	data.FvEpIpTagId = types.StringValue(attributes["id_attribute"].(string))
	data.Ip = types.StringValue(attributes["ip"].(string))
	data.Name = types.StringValue(attributes["name"].(string))
	data.NameAlias = types.StringValue(attributes["name_alias"].(string))

	planAnnotations := convertToEpIpTagAnnotations(attributes["annotations"])
	planTags := convertToEpIpTagTags(attributes["tags"])

	stateAnnotations := planAnnotations
	stateTags := planTags

	newEndpointTag := provider.GetFvEpIpTagCreateJsonPayload(ctx, &diags, data, planAnnotations, stateAnnotations, planTags, stateTags)

	payload_bin := newEndpointTag.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(payload_bin)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFvEpIpTagId(ctx, data)
	attrs := payload["fvEpIpTag"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
}

func createNetworkInstanceProfile(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.MgmtInstPResourceModel{}

	data.Name = types.StringValue(attributes["name"].(string))
	data.Descr = types.StringValue(attributes["description"].(string))
	data.Prio = types.StringValue(attributes["priority"].(string))
	data.Annotation = types.StringValue(attributes["annotation"].(string))
	data.NameAlias = types.StringValue(attributes["name_alias"].(string))

	planOoBCons := convertToOoBCons(attributes["relation_to_consumed_out_of_band_contracts"])
	stateOoBCons := planOoBCons

	planAnnotations := convertToMgmtInstPAnnotations(attributes["annotations"])
	stateAnnotations := planAnnotations

	planTags := convertToMgmtInstPTags(attributes["tags"])
	stateTags := planTags

	newExtItem := provider.GetMgmtInstPCreateJsonPayload(ctx, &diags, data, planOoBCons, stateOoBCons, planAnnotations, stateAnnotations, planTags, stateTags)

	jsonPayload := newExtItem.EncodeJSON(container.EncodeOptIndent("", "  "))
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

func createVrfFallbackRouteGroup(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvFBRGroupResourceModel{}

	data.ParentDn = types.StringValue(attributes["parent_dn"].(string))
	data.Name = types.StringValue(attributes["name"].(string))
	data.Descr = types.StringValue(attributes["description"].(string))
	data.Annotation = types.StringValue(attributes["annotation"].(string))
	data.NameAlias = types.StringValue(attributes["name_alias"].(string))

	planMembers := convertToFBRMembers(attributes["vrf_fallback_route_group_members"])
	planAnnotations := convertToFvFBRGroupAnnotations(attributes["annotations"])
	planTags := convertToFvFBRGroupTags(attributes["tags"])

	newAciVrf := provider.GetFvFBRGroupCreateJsonPayload(ctx, &diags, data, planMembers, planMembers, planAnnotations, planAnnotations, planTags, planTags)

	jsonPayload := newAciVrf.EncodeJSON(container.EncodeOptIndent("", "  "))
	payload, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	provider.SetFvFBRGroupId(ctx, data)
	attrs := payload["fvFBRGroup"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	if status, ok := attributes["status"].(string); ok && status != "" {
		attrs["status"] = status
	}

	return payload
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

func convertToOoBCons(oobCons interface{}) []provider.MgmtRsOoBConsMgmtInstPResourceModel {
	var planOoBCons []provider.MgmtRsOoBConsMgmtInstPResourceModel
	if oobCons, ok := oobCons.([]interface{}); ok {
		for _, oobCon := range oobCons {
			oobConMap := oobCon.(map[string]interface{})
			planOoBCons = append(planOoBCons, provider.MgmtRsOoBConsMgmtInstPResourceModel{
				Annotation:      types.StringValue(oobConMap["annotation"].(string)),
				Prio:            types.StringValue(oobConMap["priority"].(string)),
				TnVzOOBBrCPName: types.StringValue(oobConMap["out_of_band_contract_name"].(string)),
			})
		}
	}
	return planOoBCons
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
