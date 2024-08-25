package dict

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetDnToAciClassMap(childClass string, parentPrefix string) string {
	rnMapping := make(map[string]map[string]string)

	resp, err := http.Get("https://pubhub.devnetcloud.com/media/model-doc-latest/docs/doc/jsonmeta/aci-meta.json")
	if err != nil {
		fmt.Printf("Error fetching metadata from URL: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error fetching metadata: received non-200 status code %d\n", resp.StatusCode)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return ""
	}

	var metaData map[string]interface{}
	err = json.Unmarshal(body, &metaData)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return ""
	}

	classes, ok := metaData["classes"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid format for classes in metadata")
		return ""
	}

	for aciClassName, aciClassInfoRaw := range classes {
		aciClassInfo, ok := aciClassInfoRaw.(map[string]interface{})
		if !ok {
			continue
		}

		rnFormat, ok := aciClassInfo["rnFormat"].(string)
		if !ok {
			continue
		}
		rnPrefix := strings.Split(rnFormat, "-")[0]

		rnMap, ok := aciClassInfo["rnMap"].(map[string]interface{})
		if !ok {
			continue
		}

		for _, childClassRaw := range rnMap {
			childClass, ok := childClassRaw.(string)
			if !ok {
				continue
			}

			if _, exists := rnMapping[childClass]; !exists {
				rnMapping[childClass] = map[string]string{}
			}
			rnMapping[childClass][rnPrefix] = aciClassName
		}
	}

	if class, found := rnMapping[childClass][parentPrefix]; found {
		return class
	}

	return ""
}
