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

// Work in progress... getAciClass() needs to be generated and referenced outside of main
func getAciClass(prefix string) string {
	mapping := map[string]string{
		"tn":        "fvTenant",
		"epg":       "fvAEPg",
		"ap":        "fvAp",
		"BD":        "fvBD",
		"subnet":    "fvSubnet",
		"instP":     "l3extInstP",
		"extsubnet": "l3extSubnet",
		"ctx":       "fvCtx",
	}

	if class, found := mapping[prefix]; found {
		return class
	}
	return ""
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
		payload := createFunc(values)

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
			data = append(data, payload)
		}
	}

	return data
}

// Work in progress... based on constructTree() implementation from ansible-nd/plugins/module_utils/ndi.py
func constructTree(resources []map[string]interface{}) map[string]interface{} {

	rootNode := map[string]interface{}{
		"attributes": nil,
		"children":   map[string]interface{}{},
	}

	for _, resourceList := range resources {
		for resourceType, resourceData := range resourceList {
			resourceAttributes, valid := resourceData.(map[string]interface{})
			if !valid {
				continue
			}

			attributes, valid := resourceAttributes["attributes"].(map[string]interface{})
			if !valid {
				continue
			}

			dn, valid := attributes["dn"].(string)
			if !valid {
				continue
			}

			pathSegments := parsePath(dn)
			currentNode := rootNode
			currentDn := ""

			for i, segment := range pathSegments {
				if i == 0 && segment == "uni" {
					continue
				}

				currentDn = strings.TrimPrefix(currentDn+"/"+segment, "/")

				prefix := strings.Split(segment, "-")[0]
				className := getAciClass(prefix)
				if className == "" {
					className = segment
				}

				childNodes, valid := currentNode["children"].(map[string]interface{})
				if !valid {
					childNodes = map[string]interface{}{}
					currentNode["children"] = childNodes
				}

				if _, exists := childNodes[className]; !exists {
					childNodes[className] = map[string]interface{}{
						"attributes": nil,
						"children":   []map[string]interface{}{},
					}
				}

				nextNode, valid := childNodes[className].(map[string]interface{})
				if !valid {
					continue
				}
				currentNode = nextNode
			}

			childList, valid := currentNode["children"].([]map[string]interface{})
			if !valid {
				childList = []map[string]interface{}{}
			}

			childList = append(childList, map[string]interface{}{
				resourceType: resourceAttributes,
			})

			currentNode["children"] = childList
		}
	}

	return map[string]interface{}{
		"uni": rootNode,
	}
}

func parsePath(dn string) []string {
	var pathSegments []string
	var segmentBuffer string
	inBracket := false

	for _, char := range dn {
		switch char {
		case '/':
			if !inBracket {
				if segmentBuffer != "" {
					pathSegments = append(pathSegments, segmentBuffer)
					segmentBuffer = ""
				}
			} else {
				segmentBuffer += string(char)
			}
		case '[':
			inBracket = true
			segmentBuffer += string(char)
		case ']':
			inBracket = false
			segmentBuffer += string(char)
		default:
			segmentBuffer += string(char)
		}
	}

	if segmentBuffer != "" {
		pathSegments = append(pathSegments, segmentBuffer)
	}

	return pathSegments
}

func main() {
	if len(os.Args) != 1 {
		fmt.Println("Usage: no arguments needed")
		os.Exit(1)
	}

	outputFile := "payload.json"

	// Run Terraform commands, generate the plan JSON
	planJSON, err := runTerraform()
	if err != nil {
		log.Fatalf("Error running Terraform: %v", err)
	}

	// Read the plan file and unmarshal it into a Plan struct
	plan, err := readPlan(planJSON)
	if err != nil {
		log.Fatalf("Error reading plan: %v", err)
	}

	// Create the payload list from the plan
	payloadList := createPayloadList(plan)

	// Construct the tree from the payload list
	aciPayload := constructTree(payloadList)

	// Writes the final ACI payload to the output file
	err = writeToFile(outputFile, aciPayload)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("ACI Payload written to %s\n", outputFile)
}
