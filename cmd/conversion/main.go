package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/convert_funcs"
)

// Represents parts of the Terraform Plan we are interested in
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

// Converts Plan json into bytes, then unmarshals into Plan struct
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
	outputData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to convert data to JSON: %w", err)
	}

	if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func createItem(resourceType string, values map[string]interface{}, status string) map[string]interface{} {
	if create, exists := convert_funcs.ResourceMap[resourceType]; exists {
		item := create(values)
		if status == "deleted" && item != nil {
			if attributes, ok := item[resourceType].(map[string]interface{}); ok {
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

// Ranges through resources to create each item, adds to []map[string]interface{}
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

func constructTree(data []map[string]interface{}) map[string]interface{} {
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

	tree := make(map[string]interface{})
	for _, resource := range resourceMap {
		for key, obj := range resource {
			if key == "unknownParent" {
				children := obj.(map[string]interface{})["children"].([]map[string]interface{})
				for _, child := range children {
					for childKey, childValue := range child {
						tree[childKey] = childValue
					}
				}
			} else {
				if attributes, ok := obj.(map[string]interface{})["attributes"].(map[string]interface{}); ok {
					if parentDn, ok := attributes["parent_dn"].(string); !ok || parentDn == "" {
						tree[key] = obj
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

	err = writeToFile(outputFile, aciPayload)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("ACI Payload written to %s\n", outputFile)
}
