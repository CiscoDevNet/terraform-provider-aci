package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/convert_funcs"
	"github.com/CiscoDevNet/terraform-provider-aci/v2/dict"
	"github.com/CiscoDevNet/terraform-provider-aci/v2/test"
)

const (
	payloadFile = "payload.json"
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

func runTerraform(terraformDir string) (string, error) {
	planBin := "plan.bin"
	planJSON := "plan.json"

	if err := os.Chdir(terraformDir); err != nil {
		return "", fmt.Errorf("failed to change directory: %w", err)
	}

	if err := exec.Command("terraform", "init").Run(); err != nil {
		return "", fmt.Errorf("failed to run terraform init: %w", err)
	}

	if err := exec.Command("terraform", "plan", "-out="+planBin).Run(); err != nil {
		return "", fmt.Errorf("failed to run terraform plan: %w", err)
	}

	output, err := os.Create(planJSON)
	if err != nil {
		return "", fmt.Errorf("failed to create JSON file: %w", err)
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

	os.Remove("plan.bin")
	os.Remove("plan.json")

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

type TreeNode struct {
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	Children   map[string]*TreeNode   `json:"children,omitempty"`
	ClassName  string                 `json:"-"`
}

func NewTreeNode(className string, attributes map[string]interface{}) *TreeNode {
	return &TreeNode{
		ClassName:  className,
		Attributes: attributes,
		Children:   make(map[string]*TreeNode),
	}
}

func constructTree(resources []map[string]interface{}) map[string]interface{} {
	rootNode := NewTreeNode("uni", map[string]interface{}{"dn": "uni"})

	dnMap := make(map[string]*TreeNode)
	dnMap["uni"] = rootNode

	for _, resourceList := range resources {
		for resourceType, resourceData := range resourceList {
			resourceAttributes := resourceData.(map[string]interface{})
			attributes := safeMapInterface(resourceAttributes, "attributes")
			dn := safeString(attributes, "dn")

			var children []map[string]interface{}
			if rawChildren, ok := resourceAttributes["children"].([]interface{}); ok {
				for _, child := range rawChildren {
					if childMap, ok := child.(map[string]interface{}); ok {
						children = append(children, childMap)
					}
				}
			}

			buildParentPath(dnMap, rootNode, resourceType, dn, attributes, children)
		}
	}

	removeParentDn(rootNode)

	tenants := filterTenants(rootNode)
	if len(tenants) > 1 {
		return map[string]interface{}{
			"polUni": map[string]interface{}{
				"attributes": map[string]interface{}{},
				"children":   tenants,
			},
		}
	}

	return map[string]interface{}{
		"imdata": tenants,
	}
}

func buildParentPath(dnMap map[string]*TreeNode, rootNode *TreeNode, resourceType, dn string, attributes map[string]interface{}, children []map[string]interface{}) {
	if dn == "" && resourceType == "" {
		return
	}

	cursor := rootNode
	if dn != "" {
		cursor = traverseOrCreatePath(dnMap, rootNode, resourceType, dn)
	}

	var leafNode *TreeNode
	if existingLeafNode, exists := dnMap[dn]; exists {
		for key, value := range attributes {
			existingLeafNode.Attributes[key] = value
		}
		leafNode = existingLeafNode
	} else {
		leafNode = NewTreeNode(resourceType, attributes)
		cursor.Children[dn] = leafNode
		dnMap[dn] = leafNode
	}

	for _, child := range children {
		for childClassName, childData := range child {
			childAttributes := safeMapInterface(childData.(map[string]interface{}), "attributes")
			childDn := safeString(childAttributes, "dn")

			childKey := childDn
			if childDn == "" {
				childKey = generateUniqueKeyForNonDnNode(childClassName, childAttributes)
			}

			if _, exists := leafNode.Children[childKey]; !exists {
				childNode := NewTreeNode(childClassName, childAttributes)
				leafNode.Children[childKey] = childNode
				dnMap[childKey] = childNode
			}

			if grandChildren, ok := childData.(map[string]interface{})["children"].([]interface{}); ok {
				for _, grandchild := range grandChildren {
					if grandchildMap, ok := grandchild.(map[string]interface{}); ok {
						buildParentPath(dnMap, leafNode, childClassName, childDn, safeMapInterface(grandchildMap, "attributes"), []map[string]interface{}{grandchildMap})
					}
				}
			}
		}
	}
}

func filterTenants(node *TreeNode) []map[string]interface{} {
	tenants := []map[string]interface{}{}
	for _, child := range node.Children {
		if child.ClassName == "fvTenant" {
			tenants = append(tenants, exportTree(child))
		}
	}
	return tenants
}

func generateUniqueKeyForNonDnNode(resourceType string, attributes map[string]interface{}) string {
	return fmt.Sprintf("%s-%v", resourceType, attributes["name"])
}

func traverseOrCreatePath(dnMap map[string]*TreeNode, rootNode *TreeNode, resourceType, dn string) *TreeNode {
	pathSegments := strings.Split(dn, "/")
	cursor := rootNode

	classNames := parseClassNames(pathSegments, resourceType)

	for i := 1; i < len(pathSegments); i++ {
		className := classNames[i-1]
		currentDn := strings.Join(pathSegments[:i+1], "/")

		if existingNode, exists := dnMap[currentDn]; exists {
			cursor = existingNode
		} else {
			newNode := NewTreeNode(className, map[string]interface{}{
				"dn": currentDn,
			})
			cursor.Children[currentDn] = newNode
			cursor = newNode
			dnMap[currentDn] = newNode
		}
	}

	return cursor
}

func parseClassNames(pathSegments []string, resourceType string) []string {
	classNames := []string{resourceType}
	for i := len(pathSegments) - 2; i >= 0; i-- {
		prefix := strings.Split(pathSegments[i], "-")[0]
		if pathSegments[i] == "uni" {
			break
		}
		className := dict.GetDnToAciClassMap(classNames[len(classNames)-1], prefix)
		classNames = append(classNames, className)
	}

	for i, j := 0, len(classNames)-1; i < j; i, j = i+1, j-1 {
		classNames[i], classNames[j] = classNames[j], classNames[i]
	}
	return classNames
}

func exportTree(node *TreeNode) map[string]interface{} {
	if node == nil {
		return nil
	}

	result := map[string]interface{}{
		node.ClassName: map[string]interface{}{
			"attributes": node.Attributes,
		},
	}

	if len(node.Children) > 0 {
		children := []map[string]interface{}{}
		for _, child := range node.Children {
			children = append(children, exportTree(child))
		}
		result[node.ClassName].(map[string]interface{})["children"] = children
	}

	return result
}

func removeParentDn(node *TreeNode) {
	if node == nil {
		return
	}

	delete(node.Attributes, "parent_dn")

	for _, child := range node.Children {
		removeParentDn(child)
	}
}

func safeMapInterface(data map[string]interface{}, key string) map[string]interface{} {
	if value, ok := data[key].(map[string]interface{}); ok {
		return value
	}
	return nil
}

func safeString(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}

func main() {
	testFlag := flag.Bool("test", false, "Run the test to POST payload to APIC and validate change")
	terraformDir := flag.String("dir", ".", "Path to the directory containing Terraform file")

	flag.Parse()

	planJSON, err := runTerraform(*terraformDir)
	if err != nil {
		log.Fatalf("Error running Terraform: %v", err)
	}

	plan, err := readPlan(planJSON)
	if err != nil {
		log.Fatalf("Error reading plan: %v", err)
	}

	payloadList := createPayloadList(plan)

	aciTree := constructTree(payloadList)

	err = writeToFile(payloadFile, aciTree)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("ACI Payload written to %s\n", payloadFile)

	if *testFlag {
		fmt.Println("Testing...")

		if _, err := os.Stat(payloadFile); os.IsNotExist(err) {
			fmt.Printf("Expected %s not found: %v\n", payloadFile, err)
			os.Exit(3)
		}

		payload, err := os.ReadFile(payloadFile)
		if err != nil {
			fmt.Printf("Failed to read %s: %v\n", payloadFile, err)
			os.Exit(4)
		}

		token, err := test.GetAPICLoginToken()
		if err != nil {
			fmt.Printf("Failed to obtain login token: %v\n", err)
			os.Exit(5)
		}

		err = test.PostToAPIC(token, payload)
		if err != nil {
			fmt.Printf("Failed to post payload to APIC: %v\n", err)
			os.Exit(6)
		}

		emptyPlan, err := test.CheckTerraformPlan()
		if err != nil {
			fmt.Printf("Terraform plan failed: %v\n", err)
			os.Exit(7)
		}

		if !emptyPlan {
			fmt.Println("Terraform plan detected changes, the conversion was not successful")

			os.Remove("plan.bin")
			os.Remove("plan.json")
			os.Exit(8)
		}

		os.Remove("plan.bin")
		os.Remove("plan.json")
		fmt.Println("Test passed: The conversion was successful")
	}
}
