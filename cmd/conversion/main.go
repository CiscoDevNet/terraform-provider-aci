package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/convert_funcs"
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

func runTerraform() (string, error) {
	planBin := "plan.bin"
	planJSON := "plan.json"

	if err := exec.Command("terraform", "plan", "-out="+planBin).Run(); err != nil {
		return "", fmt.Errorf("failed to run terraform plan: %w", err)
	}

	output, err := os.Create(planJSON)
	if err != nil {
		return "", fmt.Errorf("failed to create json file: %w", err)
	}
	defer output.Close()

	cmdShow := exec.Command("terraform", "show", "-json", planBin)
	cmdShow.Stdout = output
	if err := cmdShow.Run(); err != nil {
		return "", fmt.Errorf("failed to run terraform show: %w", err)
	}

	return planJSON, nil
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

func createPayload(resourceType string, values map[string]interface{}, status string) map[string]interface{} {
	if createFunc, exists := convert_funcs.ResourceMap[resourceType]; exists {
		payload := createFunc(values, status)
		return payload
	}
	return nil
}

func createPayloadList(plan Plan) []map[string]interface{} {
	var data []map[string]interface{}

	for _, change := range plan.Changes {
		if len(change.Change.Actions) > 0 && change.Change.Actions[0] == "delete" {
			payload := createPayload(change.Type, change.Change.Before, "deleted")
			if payload != nil {
				data = append(data, payload)
			}
		}
	}

	for _, resource := range plan.PlannedValues.RootModule.Resources {
		payload := createPayload(resource.Type, resource.Values, "")
		if payload != nil {
			for _, value := range payload {
				if obj, ok := value.(map[string]interface{}); ok {
					if attributes, ok := obj["attributes"].(map[string]interface{}); ok {
						if parentDn, ok := resource.Values["parent_dn"].(string); ok && parentDn != "" {
							attributes["parent_dn"] = parentDn
						}
					}
				}
			}
			data = append(data, payload)
		}
	}

	return data
}

// Work in progress... "panic: assignment to entry in nil map" for cases where AciClassMap() returns "" in this function
func constructTree(resources []map[string]interface{}) map[string]interface{} {
	rootNode := map[string]interface{}{
		"children": []map[string]interface{}{},
	}

	nodeMap := map[string]map[string]interface{}{
		"uni": rootNode,
	}

	for _, resourceList := range resources {
		for resourceType, resourceData := range resourceList {
			resourceAttributes := resourceData.(map[string]interface{})
			attributes := resourceAttributes["attributes"].(map[string]interface{})
			dn := attributes["dn"].(string)

			parentDn, hasParentDn := attributes["parent_dn"].(string)
			if !hasParentDn || parentDn == "" {
				pathSegments := strings.Split(dn, "/")
				if len(pathSegments) > 1 {
					parentDn = strings.Join(pathSegments[:len(pathSegments)-1], "/")

				} else {
					parentDn = "uni"
				}
			}

			if _, parentExists := nodeMap[parentDn]; !parentExists {
				createParentPath(nodeMap, parentDn)
			}

			currentNode, nodeExists := nodeMap[dn]
			if !nodeExists {
				currentNode = map[string]interface{}{
					"attributes": attributes,
					"children":   resourceAttributes["children"],
				}
			} else {

				existingChildren, ok := currentNode["children"].([]map[string]interface{})
				if !ok || existingChildren == nil {
					existingChildren = []map[string]interface{}{}
				}

				newChildren, ok := resourceAttributes["children"].([]map[string]interface{})
				if ok && newChildren != nil {
					currentNode["children"] = append(existingChildren, newChildren...)
				} else {
					currentNode["children"] = existingChildren
				}
			}

			parentNode := nodeMap[parentDn]

			className := convert_funcs.GetAciClass(strings.Split(dn, "/")[len(strings.Split(dn, "/"))-1])
			if className == "" {
				className = resourceType
			}

			parentChildren, ok := parentNode["children"].([]map[string]interface{})
			if !ok || parentChildren == nil {
				parentChildren = []map[string]interface{}{}
			}

			parentNode["children"] = append(parentChildren, map[string]interface{}{className: currentNode})

			nodeMap[dn] = currentNode
		}
	}

	return map[string]interface{}{
		"uni": rootNode,
	}
}

func createParentPath(nodeMap map[string]map[string]interface{}, parentDn string) {
	pathSegments := strings.Split(parentDn, "/")
	currentDn := "uni"
	var lastValidNode map[string]interface{}

	for _, segment := range pathSegments[1:] {
		currentDn += "/" + segment

		className := convert_funcs.GetAciClass(strings.Split(segment, "-")[0])

		if className == "" {
			className = segment
		}

		if _, exists := nodeMap[currentDn]; !exists {
			newNode := map[string]interface{}{
				"attributes": map[string]interface{}{
					"dn": currentDn,
				},
				"children": []map[string]interface{}{}, // Ensure children is always initialized
			}

			// Attach this node to the last valid node or the root if at the top level
			if lastValidNode == nil {
				if _, ok := nodeMap["uni"]["children"]; !ok {
					nodeMap["uni"]["children"] = []map[string]interface{}{}
				}
				nodeMap["uni"]["children"] = append(nodeMap["uni"]["children"].([]map[string]interface{}), map[string]interface{}{className: newNode})
			} else {
				if _, ok := lastValidNode["children"]; !ok {
					lastValidNode["children"] = []map[string]interface{}{}
				}
				lastValidNode["children"] = append(lastValidNode["children"].([]map[string]interface{}), map[string]interface{}{className: newNode})
			}

			nodeMap[currentDn] = newNode
		}

		// Update the last valid node reference
		lastValidNode = nodeMap[currentDn]
	}
}

func main() {

	if len(os.Args) != 1 {
		fmt.Println("Usage: no arguments needed")
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

	payloadList := createPayloadList(plan)

	jsonData, err := json.MarshalIndent(payloadList, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonData))

	aciPayload := constructTree(payloadList)

	err = writeToFile(outputFile, aciPayload)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("ACI Payload written to %s\n", outputFile)
}
