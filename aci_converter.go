package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// Represents the structure for a Terraform plan input
type TerraformPlan struct {
	PlannedValues struct {
		RootModule struct {
			Resources []Resource `json:"resources"`
		} `json:"root_module"`
	} `json:"planned_values"`
}

// Represents each resource in the Terraform plan input
type Resource struct {
	Address         string                 `json:"address"`
	Mode            string                 `json:"mode"`
	Type            string                 `json:"type"`
	Name            string                 `json:"name"`
	ProviderName    string                 `json:"provider_name"`
	SchemaVersion   int                    `json:"schema_version"`
	Values          map[string]interface{} `json:"values"`
	SensitiveValues map[string]interface{} `json:"sensitive_values"`
}

// omitempty is used to not include children{[]} under a resource with no children
type Item struct {
	Attributes map[string]interface{} `json:"attributes"`
	Children   []map[string]Item      `json:"children,omitempty"`
}

// Writes data to a JSON file
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

// Writes data to a JSON file after removing outer array brackets
func outputToFileWithoutArrayBrackets(outputFile string, data []map[string]Item) error {
	var outputStr strings.Builder

	for i, item := range data {
		outputData, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to convert data to JSON: %w", err)
		}

		itemStr := string(outputData)
		if len(itemStr) > 2 && itemStr[0] == '{' && itemStr[len(itemStr)-1] == '}' {
			itemStr = itemStr[1 : len(itemStr)-1]
		}

		outputStr.WriteString(itemStr)
		if i < len(data)-1 {
			outputStr.WriteString(",\n")
		}
	}

	outputStrFinal := "{" + outputStr.String() + "}"

	err := os.WriteFile(outputFile, []byte(outputStrFinal), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// Processes the resources in the Terraform plan
func processResources(terraformPlan TerraformPlan) []map[string]Item {
	var itemList []map[string]Item

	for _, resource := range terraformPlan.PlannedValues.RootModule.Resources {
		attributes := extractAttributes(resource.Values)
		item := map[string]Item{}

		itemType, mappedAttrs := mapResourceTypeAndAttributes(resource.Type, attributes)
		if itemType != "" {
			item[itemType] = Item{Attributes: mappedAttrs}
			itemList = append(itemList, item)
		}
	}

	return itemList
}

// Maps  resource type and attributes
// Can be updated to fix format errors in attributes{} for resources in output
func mapResourceTypeAndAttributes(resourceType string, attributes map[string]interface{}) (string, map[string]interface{}) {
	switch resourceType {
	case "aci_tenant":
		dn := fmt.Sprintf("uni/tn-%s", attributes["name"])
		attributes["dn"] = dn
		attributes = mapAttributes(attributes, map[string]string{
			"description": "descr",
			"name_alias":  "nameAlias",
		})
		return "fvTenant", attributes
	case "aci_annotation":
		if parentDn, ok := attributes["parent_dn"].(string); ok {
			attributes["dn"] = fmt.Sprintf("%s/tagAnnotation-%s", parentDn, attributes["key"])
			delete(attributes, "parent_dn")
		}
		return "tagAnnotation", attributes
	case "aci_endpoint_tag_ip":
		return "fvEpIpTag", attributes
	case "aci_endpoint_tag_mac":
		return "fvEpMacTag", attributes
	case "aci_external_management_network_instance_profile":
		return "mgmtInstP", attributes
	case "aci_external_management_network_subnet":
		return "mgmtSubnet", attributes
	case "aci_l3out_consumer_label":
		return "l3extConsLbl", attributes
	case "aci_l3out_node_sid_profile":
		return "mplsNodeSidP", attributes
	case "aci_l3out_provider_label":
		return "l3extProvLbl", attributes
	case "aci_l3out_redistribute_policy":
		return "l3extRsRedistributePol", attributes
	case "aci_netflow_monitor_policy":
		return "netflowMonitorPol", attributes
	case "aci_out_of_band_contract":
		return "vzOOBBrCP", attributes
	case "aci_pim_route_map_entry":
		return "pimRouteMapEntry", attributes
	case "aci_pim_route_map_policy":
		return "pimRouteMapPol", attributes
	case "aci_relation_to_consumed_out_of_band_contract":
		return "mgmtRsOoBCons", attributes
	case "aci_relation_to_fallback_route_group":
		return "l3extRsOutToFBRGroup", attributes
	case "aci_relation_to_netflow_exporter":
		return "netflowRsMonitorToExporter", attributes
	case "aci_tag":
		return "tagTag", attributes
	case "aci_vrf_fallback_route_group_member":
		return "fvFBRMember", attributes
	case "aci_vrf_fallback_route_group":
		return "fvFBRGroup", attributes
	default:
		return "", nil
	}
}

// Maps attribute keys to new keys
func mapAttributes(attributes map[string]interface{}, mappings map[string]string) map[string]interface{} {
	mapped := make(map[string]interface{})
	for key, value := range attributes {
		if newKey, exists := mappings[key]; exists {
			mapped[newKey] = value
		} else {
			mapped[key] = value
		}
	}
	return mapped
}

// Creates a map of attributes for a given resource from input values
func extractAttributes(values map[string]interface{}) map[string]interface{} {
	attributes := make(map[string]interface{})
	for key, value := range values {
		if valueStr, ok := value.(string); ok && valueStr != "" {
			attributes[key] = valueStr
		}
	}
	return attributes
}

// Cnstructs tree hierarchy for resources
func constructTree(items []map[string]Item) []map[string]Item {
	dnMap := make(map[string]*Item)
	var root []map[string]Item

	for _, item := range items {
		for resourceType, resource := range item {
			dn, ok := resource.Attributes["dn"].(string)
			if !ok {
				log.Fatalf("Missing dn attribute in resource of type %s", resourceType)
			}
			dnMap[dn] = &resource
			resource.Attributes["resourceType"] = resourceType
		}
	}

	for dn, resource := range dnMap {
		parts := strings.Split(dn, "/")
		if len(parts) < 2 {
			log.Printf("Invalid DN: %s", dn)
			continue
		}
		parentDn := strings.Join(parts[:len(parts)-1], "/")

		if parent, exists := dnMap[parentDn]; exists {
			if parent.Children == nil {
				parent.Children = []map[string]Item{}
			}
			childType := resource.Attributes["resourceType"].(string)
			parent.Children = append(parent.Children, map[string]Item{childType: *resource})
		} else {
			root = append(root, map[string]Item{resource.Attributes["resourceType"].(string): *resource})
		}
	}

	return root
}

// DN is needed to map child resources to parent in tree construction
// DN is delete after tree construction to match ACI Payload format
func removeDnFromTagAnnotations(data []map[string]Item) {
	for _, item := range data {
		for resourceType, resource := range item {
			if resourceType == "tagAnnotation" {
				delete(resource.Attributes, "dn")
			}
			if resource.Children != nil {
				removeDnFromTagAnnotations(resource.Children)
			}
		}
	}
}

// ConstructTree() requires TreeNode keys to be the DN
// This function traverses tree and replaces DN with resource type to match ACI Payload format
func replaceDnWithResourceType(data []map[string]Item) []map[string]Item {
	var result []map[string]Item

	for _, item := range data {
		for _, resource := range item {
			if resource.Children != nil {
				resource.Children = replaceDnWithResourceType(resource.Children)
			}
			resourceType := resource.Attributes["resourceType"].(string)
			newItem := map[string]Item{resourceType: resource}
			delete(resource.Attributes, "resourceType")
			result = append(result, newItem)
		}
	}

	return result
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run aci_converter.go <input.json> <output.json>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	var terraformPlan TerraformPlan
	err = json.Unmarshal(inputData, &terraformPlan)
	if err != nil {
		log.Fatalf("Error parsing input file: %v", err)
	}

	itemList := processResources(terraformPlan)
	tree := constructTree(itemList)
	removeDnFromTagAnnotations(tree)
	finalTree := replaceDnWithResourceType(tree)

	err = outputToFileWithoutArrayBrackets(outputFile, finalTree)

	fmt.Printf("ACI Payload written to %s\n", outputFile)
}
