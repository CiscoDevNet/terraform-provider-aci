package aci

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func checkTDn(client *client.Client, dns []string) error {
	flag := false
	var errMessage string

	for _, dn := range dns {
		_, err := client.Get(dn)
		if err != nil {
			if flag == false {
				flag = true
			}
			errMessage = fmt.Sprintf("%s\nRelation target dn %s not found", errMessage, dn)
		}
	}

	if flag == true {
		return fmt.Errorf(errMessage)
	}

	return nil
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

func GetParentDn(dn string, rn string) string {
	arr := strings.Split(dn, rn)
	return arr[0]
}

func stripQuotes(word string) string {
	if strings.HasPrefix(word, "\"") && strings.HasSuffix(word, "\"") {
		return strings.TrimSuffix(strings.TrimPrefix(word, "\""), "\"")
	}
	return word
}

func validateCommaSeparatedStringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		// modified validation.StringInSlice function.
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		vals := strings.Split(v, ",")
		elemap := make(map[string]bool)
		for _, val := range vals {
			if !elemap[val] {
				match := false
				for _, str := range valid {
					if val == str || (ignoreCase && strings.ToLower(val) == strings.ToLower(str)) {
						match = true
					}
				}
				if !match {
					es = append(es, fmt.Errorf("expected %s to be one of %v, got %s", k, valid, val))
				}
			} else {
				es = append(es, fmt.Errorf("unexpected duplicate values in %s : %s", k, val))
			}
			elemap[val] = true
		}
		return
	}
}

func suppressBitMaskDiffFunc() func(k, old, new string, d *schema.ResourceData) bool {
	return func(k, old, new string, d *schema.ResourceData) bool {
		oldList := strings.Split(old, ",")
		newList := strings.Split(new, ",")
		sort.Strings(oldList)
		sort.Strings(newList)

		return reflect.DeepEqual(oldList, newList)
	}
}
