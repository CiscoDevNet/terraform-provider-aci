package aci

import (
	"fmt"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/container"
)

func toStrMap(inputMap map[string]interface{}) map[string]string {
	rt := make(map[string]string)
	for key, value := range inputMap {
		rt[key] = value.(string)
	}

	return rt
}

func toStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

func preparePayload(className string, inputMap map[string]string) (*container.Container, error) {
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
			}
		}
	}`, className))
	cont, err := container.ParseJSON(containerJSON)

	if err != nil {
		return nil, err
	}
	for key, value := range inputMap {
		cont.Set(value, className, "attributes", key)
	}
	return cont, nil

}

func GetMOName(dn string) string {
	arr := strings.Split(dn, "/")
	// Get the last element
	last_ele := arr[len(arr)-1]
	// split on -
	dash_split := strings.Split(last_ele, "-")
	// join except first element as that will be rn
	return strings.Join(dash_split[1:], "-")

	// re := regexp.MustCompile(".*/\\S+-(\\S+.*)$")
	// match := re.FindStringSubmatch(dn)
	// return match[1]

}

// func GetParentDn(childDn string) string {
// 	arr := strings.Split(childDn, "/")
// 	// in case of cidr blocks we have extra / in the ip range so let's catch it and remove. This will have extra part.
// 	if strings.Contains(childDn, "]") && strings.Contains(childDn, "[") {
// 		dnWithRn := strings.Split(childDn, "[")

// 		slashedArr := strings.Split(dnWithRn[0], "/")
// 		return strings.Join(slashedArr[:len(slashedArr)-1], "/")

// 	}

// 	return strings.Join(arr[:len(arr)-1], "/")

// }
