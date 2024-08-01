package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Item represents the structure of the data
type Item struct {
	Attributes map[string]interface{} `json:"attributes"`
	Children   []map[string]Item      `json:"children,omitempty"`
}

// Helper function to check if a slice contains a specific string
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Function to remove duplicate entries from a slice of strings
func removeDuplicates(strList []string) []string {
	list := []string{}
	for _, item := range strList {
		if !contains(list, item) {
			list = append(list, item)
		}
	}
	return list
}

// Recursive function to clean the item structure
func cleanItem(item map[string]Item) map[string]Item {
	cleanedItem := make(map[string]Item)
	for key, val := range item {
		val.Attributes = cleanAttributes(val.Attributes)
		if val.Children != nil {
			for i, child := range val.Children {
				for childKey, childVal := range child {
					val.Children[i][childKey] = cleanItem(map[string]Item{childKey: childVal})[childKey]
				}
			}
		}
		cleanedItem[key] = val
	}
	return cleanedItem
}

// Function to clean redundant attributes
func cleanAttributes(attributes map[string]interface{}) map[string]interface{} {
	cleanedAttributes := make(map[string]interface{})
	for key, val := range attributes {
		if nestedAttrs, ok := val.(map[string]interface{}); ok && key == "attributes" {
			for nestedKey, nestedVal := range nestedAttrs {
				cleanedAttributes[nestedKey] = nestedVal
			}
		} else {
			cleanedAttributes[key] = val
		}
	}
	return cleanedAttributes
}

func main() {
	// Input JSON with redundant attributes
	inputJSON := `
	{
	  "fvTenant": {
		"attributes": {
		  "annotation": "orchestrator:terraform",
		  "descr": "This tenant is created by terraform",
		  "dn": "uni/tn-example_tenant",
		  "name": "example_tenant",
		  "resourceType": "fvTenant"
		},
		"children": [
		  {
			"fvFBRGroup": {
			  "attributes": {
				"annotation": "annotation",
				"descr": "description",
				"dn": "uni/tn-example_tenant/fbrg-fallback_route_group",
				"name": "fallback_route_group",
				"nameAlias": "name_alias",
				"parent_dn": "uni/tn-example_tenant",
				"resourceType": "fvFBRGroup"
			  },
			  "children": [
				{
				  "fvFBRMember": {
					"attributes": {
					  "fvFBRMember": {
						"attributes": {
						  "annotation": "annotation_1",
						  "descr": "description_1",
						  "name": "name_1",
						  "nameAlias": "name_alias_1",
						  "rnhAddr": "2.2.2.2"
						}
					  }
					}
				  }
				},
				{
				  "tagAnnotation": {
					"attributes": {
					  "tagAnnotation": {
						"attributes": {
						  "key": "key_0",
						  "value": "value_1"
						}
					  }
					}
				  }
				},
				{
				  "tagTag": {
					"attributes": {
					  "tagTag": {
						"attributes": {
						  "key": "key_0",
						  "value": "value_1"
						}
					  }
					}
				  }
				}
			  ]
			}
		  },
		  {
			"netflowMonitorPol": {
			  "attributes": {
				"annotation": "annotation",
				"descr": "description",
				"dn": "uni/tn-example_tenant/monitorpol-netfow_monitor",
				"name": "netfow_monitor",
				"nameAlias": "name_alias",
				"ownerKey": "owner_key",
				"ownerTag": "owner_tag",
				"parent_dn": "uni/tn-example_tenant",
				"resourceType": "netflowMonitorPol"
			  },
			  "children": [
				{
				  "netflowRsMonitorToExporter": {
					"attributes": {
					  "netflowRsMonitorToExporter": {
						"attributes": {
						  "annotation": "annotation_1",
						  "tnNetflowExporterPolName": "aci_netflow_exporter_policy.example.name"
						}
					  }
					}
				  }
				},
				{
				  "netflowRsMonitorToRecord": {
					"attributes": {
					  "netflowRsMonitorToRecord": {
						"attributes": {
						  "annotation": "annotation_1",
						  "tnNetflowRecordPolName": "aci_netflow_record_policy.example.name"
						}
					  }
					}
				  }
				},
				{
				  "tagAnnotation": {
					"attributes": {
					  "tagAnnotation": {
						"attributes": {
						  "key": "key_0",
						  "value": "value_1"
						}
					  }
					}
				  }
				},
				{
				  "tagTag": {
					"attributes": {
					  "tagTag": {
						"attributes": {
						  "key": "key_0",
						  "value": "value_1"
						}
					  }
					}
				  }
				}
			  ]
			}
		  }
		]
	  }
	}`

	// Unmarshal the input JSON
	var input map[string]Item
	err := json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		log.Fatalf("Error parsing input JSON: %v", err)
	}

	// Clean the input
	output := cleanItem(input)

	// Marshal the output to JSON
	outputJSON, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling output JSON: %v", err)
	}

	// Print the cleaned output
	fmt.Println(string(outputJSON))
}
