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

func mergeMultipleMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			if existing, ok := result[k].(map[string]interface{}); ok {
				if vMap, ok := v.(map[string]interface{}); ok {
					result[k] = mergeMultipleMaps(existing, vMap)
				} else {
					result[k] = v
				}
			} else {
				result[k] = v
			}
		}
	}
	return result
}

type CustomJSON map[string]interface{}

func (c *CustomJSON) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*c = parseMap(raw)
	return nil
}

func parseMap(raw map[string]interface{}) map[string]interface{} {
	parsed := make(map[string]interface{})
	for k, v := range raw {
		switch value := v.(type) {
		case map[string]interface{}:
			parsed[k] = parseMap(value)
		case []interface{}:
			parsed[k] = parseSlice(value)
		default:
			parsed[k] = value
		}
	}
	return parsed
}

func parseSlice(raw []interface{}) []interface{} {
	parsed := make([]interface{}, len(raw))
	for i, v := range raw {
		switch value := v.(type) {
		case map[string]interface{}:
			parsed[i] = parseMap(value)
		case []interface{}:
			parsed[i] = parseSlice(value)
		default:
			parsed[i] = value
		}
	}
	return parsed
}

func parseCustomJSON(jsonPayload []byte) (map[string]interface{}, error) {
	var customData CustomJSON
	err := json.Unmarshal(jsonPayload, &customData)
	if err != nil {
		return nil, err
	}
	return customData, nil
}

func printMap(data map[string]interface{}, indent int) {
	for k, v := range data {
		printIndent(indent)
		switch val := v.(type) {
		case map[string]interface{}:
			fmt.Printf("%s: {\n", k)
			printMap(val, indent+2)
			printIndent(indent)
			fmt.Println("}")
		case []interface{}:
			fmt.Printf("%s: [\n", k)
			printSlice(val, indent+2)
			printIndent(indent)
			fmt.Println("]")
		default:
			fmt.Printf("%s: %v\n", k, val)
		}
	}
}

func printSlice(data []interface{}, indent int) {
	for _, v := range data {
		printIndent(indent)
		switch val := v.(type) {
		case map[string]interface{}:
			fmt.Println("{")
			printMap(val, indent+2)
			printIndent(indent)
			fmt.Println("},")
		case []interface{}:
			fmt.Println("[")
			printSlice(val, indent+2)
			printIndent(indent)
			fmt.Println("],")
		default:
			fmt.Printf("%v,\n", val)
		}
	}
}

func printIndent(indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print(" ")
	}
}

type createItemFunc func(map[string]interface{}) map[string]interface{}

var createItemFuncMap = map[string]createItemFunc{
	"aci_endpoint_tag_ip":                              createAciEndpointTagIP,
	"aci_external_management_network_instance_profile": createAciExternalManagementNetworkInstanceProfile,
	"aci_vrf_fallback_route_group":                     createAciVrfFallbackRouteGroup,
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

func processResources(terraformPlan TerraformPlan) map[string]interface{} {
	mergedData := make(map[string]interface{})

	resourceStatusMap := make(map[string]string)
	for _, resourceChange := range terraformPlan.ResourceChanges {
		if len(resourceChange.Change.Actions) > 0 {
			resourceStatusMap[resourceChange.Address] = resourceChange.Change.Actions[0]
		}
	}

	for _, resourceChange := range terraformPlan.ResourceChanges {
		if len(resourceChange.Change.Actions) > 0 && resourceChange.Change.Actions[0] == "delete" {
			resourceType := resourceChange.Type
			beforeAttributes := resourceChange.Change.Before
			if beforeAttributes == nil {
				beforeAttributes = make(map[string]interface{})
			}
			beforeAttributes["status"] = "deleted"
			item := createItem(resourceType, beforeAttributes, "deleted")
			if item != nil {
				mergedData = mergeMultipleMaps(mergedData, item)
			}
		}
	}

	for _, resource := range terraformPlan.PlannedValues.RootModule.Resources {
		status := resourceStatusMap[fmt.Sprintf("%s.%s", resource.Type, resource.Name)]
		item := createItem(resource.Type, resource.Values, status)
		if item != nil {
			mergedData = mergeMultipleMaps(mergedData, item)
		}
	}

	return mergedData
}

func createItem(resourceType string, resourceValues map[string]interface{}, status string) map[string]interface{} {
	attributes := make(map[string]interface{})

	for key, val := range resourceValues {
		attributes[key] = val
	}

	var item map[string]interface{}
	switch resourceType {
	case "aci_tenant":
		item = createAciTenant(attributes)
	case "aci_endpoint_tag_ip":
		item = createAciEndpointTagIP(attributes)
	case "aci_external_management_network_instance_profile":
		item = createAciExternalManagementNetworkInstanceProfile(attributes)
	case "aci_vrf_fallback_route_group":
		item = createAciVrfFallbackRouteGroup(attributes)
	default:
		item = createProviderItem(resourceType, attributes)
	}

	if status == "deleted" && item != nil {
		if attributes, ok := item[resourceTypeToMapKey(resourceType)].(map[string]interface{}); ok {
			if attrs, ok := attributes["attributes"].(map[string]interface{}); ok {
				attrs["status"] = status
			}
		}
	}

	return item
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

func createProviderItem(resourceType string, attributes map[string]interface{}) map[string]interface{} {
	if createFunc, exists := createItemFuncMap[resourceType]; exists {
		return createFunc(attributes)
	}
	return nil
}

func createAciTenant(attributes map[string]interface{}) map[string]interface{} {
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
	return map[string]interface{}{
		"fvTenant": map[string]interface{}{
			"attributes": tenantAttributes,
		},
	}
}

func createAciEndpointTagIP(attributes map[string]interface{}) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.FvEpIpTagResourceModel{}

	if val, exists := attributes["annotation"].(string); exists && val != "" {
		data.Annotation = types.StringValue(val)
	}
	if val, exists := attributes["vrf_name"].(string); exists && val != "" {
		data.CtxName = types.StringValue(val)
	}
	if val, exists := attributes["id_attribute"].(string); exists && val != "" {
		data.Id = types.StringValue(val)
	}
	if val, exists := attributes["ip"].(string); exists && val != "" {
		data.Ip = types.StringValue(val)
	}
	if val, exists := attributes["name"].(string); exists && val != "" {
		data.Name = types.StringValue(val)
	}
	if val, exists := attributes["name_alias"].(string); exists && val != "" {
		data.NameAlias = types.StringValue(val)
	}

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

	stateAnnotations := planAnnotations

	var planTags []provider.TagTagFvEpIpTagResourceModel
	if tags, exists := attributes["tags"].([]interface{}); exists {
		for _, tag := range tags {
			tagMap := tag.(map[string]interface{})
			planTags = append(planTags, provider.TagTagFvEpIpTagResourceModel{
				Key:   types.StringValue(tagMap["key"].(string)),
				Value: types.StringValue(tagMap["value"].(string)),
			})
		}
	}

	stateTags := planTags

	newEndpointTag := provider.GetFvEpIpTagCreateJsonPayload(ctx, &diags, data, planAnnotations, stateAnnotations, planTags, stateTags)

	jsonPayload := newEndpointTag.EncodeJSON(container.EncodeOptIndent("", "  "))

	customData, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	// Adding status to customData if exists
	if status, exists := attributes["status"].(string); exists {
		if attributes, ok := customData["fvEpIpTag"].(map[string]interface{}); ok {
			if attrs, ok := attributes["attributes"].(map[string]interface{}); ok {
				attrs["status"] = status
			}
		}
	}

	printMap(customData, 0)

	return customData
}

func createAciExternalManagementNetworkInstanceProfile(attributes map[string]interface{}) map[string]interface{} {
	var data = provider.MgmtInstPResourceModel{}

	profileAttributes := make(map[string]interface{})
	if name, exists := attributes["name"].(string); exists {
		data.Name = types.StringValue(name)
		profileAttributes["dn"] = fmt.Sprintf("%s/extmgmt-default/instp-%s", attributes["parent_dn"], name)
	}
	if descr, exists := attributes["description"].(string); exists && descr != "" {
		data.Descr = types.StringValue(descr)
	}
	if prio, exists := attributes["priority"].(string); exists && prio != "" {
		data.Prio = types.StringValue(prio)
	}
	if annotation, exists := attributes["annotation"].(string); exists && annotation != "" {
		data.Annotation = types.StringValue(annotation)
	}
	if nameAlias, exists := attributes["name_alias"].(string); exists && nameAlias != "" {
		data.NameAlias = types.StringValue(nameAlias)
	}

	ctx := context.Background()
	var diags diag.Diagnostics

	var planOoBCons []provider.MgmtRsOoBConsMgmtInstPResourceModel
	if oobCons, exists := attributes["relation_to_consumed_out_of_band_contracts"].([]interface{}); exists {
		for _, oobCon := range oobCons {
			oobConMap := oobCon.(map[string]interface{})
			planOoBCons = append(planOoBCons, provider.MgmtRsOoBConsMgmtInstPResourceModel{
				Annotation:      types.StringValue(oobConMap["annotation"].(string)),
				Prio:            types.StringValue(oobConMap["priority"].(string)),
				TnVzOOBBrCPName: types.StringValue(oobConMap["out_of_band_contract_name"].(string)),
			})
		}
	}

	stateOoBCons := planOoBCons

	var planAnnotations []provider.TagAnnotationMgmtInstPResourceModel
	if annotations, exists := attributes["annotations"].([]interface{}); exists {
		for _, annotation := range annotations {
			annotationMap := annotation.(map[string]interface{})
			planAnnotations = append(planAnnotations, provider.TagAnnotationMgmtInstPResourceModel{
				Key:   types.StringValue(annotationMap["key"].(string)),
				Value: types.StringValue(annotationMap["value"].(string)),
			})
		}
	}

	stateAnnotations := planAnnotations

	var planTags []provider.TagTagMgmtInstPResourceModel
	if tags, exists := attributes["tags"].([]interface{}); exists {
		for _, tag := range tags {
			tagMap := tag.(map[string]interface{})
			planTags = append(planTags, provider.TagTagMgmtInstPResourceModel{
				Key:   types.StringValue(tagMap["key"].(string)),
				Value: types.StringValue(tagMap["value"].(string)),
			})
		}
	}

	stateTags := planTags

	newExtItem := provider.GetMgmtInstPCreateJsonPayload(ctx, &diags, &data, planOoBCons, stateOoBCons, planAnnotations, stateAnnotations, planTags, stateTags)

	jsonPayload := newExtItem.EncodeJSON(container.EncodeOptIndent("", "  "))

	customData, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	// Adding status to customData if exists
	if status, exists := attributes["status"].(string); exists {
		if attributes, ok := customData["mgmtInstP"].(map[string]interface{}); ok {
			if attrs, ok := attributes["attributes"].(map[string]interface{}); ok {
				attrs["status"] = status
			}
		}
	}

	printMap(customData, 0)

	return customData
}

func createAciVrfFallbackRouteGroup(attributes map[string]interface{}) map[string]interface{} {

	var data = provider.FvFBRGroupResourceModel{}

	// Set attributes
	if val, exists := attributes["parent_dn"].(string); exists {
		data.ParentDn = types.StringValue(val)
	}
	if name, exists := attributes["name"].(string); exists {
		data.Name = types.StringValue(name)
	}
	if descr, exists := attributes["description"].(string); exists && descr != "" {
		data.Descr = types.StringValue(descr)
	}
	if annotation, exists := attributes["annotation"].(string); exists && annotation != "" {
		data.Annotation = types.StringValue(annotation)
	}
	if nameAlias, exists := attributes["name_alias"].(string); exists && nameAlias != "" {
		data.NameAlias = types.StringValue(nameAlias)
	}

	ctx := context.Background()
	var diags diag.Diagnostics

	// Set plan members
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

	// Set plan annotations
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

	// Set plan tags
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

	// Get the JSON payload
	newAciVrf := provider.GetFvFBRGroupCreateJsonPayload(ctx, &diags, &data, planMembers, planMembers, planAnnotations, planAnnotations, planTags, planTags)

	// Encode the JSON payload
	jsonPayload := newAciVrf.EncodeJSON(container.EncodeOptIndent("", "  "))

	// Parse the custom JSON
	customData, err := parseCustomJSON(jsonPayload)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	// Adding status to customData if exists
	if status, exists := attributes["status"].(string); exists {
		if attributes, ok := customData["fvFBRGroup"].(map[string]interface{}); ok {
			if attrs, ok := attributes["attributes"].(map[string]interface{}); ok {
				attrs["status"] = status
			}
		}
	}

	// Print the map (for debugging purposes)
	printMap(customData, 0)

	provider.SetFvFBRGroupId(ctx, &data)

	fmt.Printf("Current DN: %s\n", data.Id.ValueString())

	return customData

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

	mergedData := processResources(terraformPlan)

	err = outputToFile(outputFile, mergedData)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("ACI Payload written to %s\n", outputFile)
}
