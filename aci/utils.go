package aci

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func setRelationAttribute(d *schema.ResourceData, relation string, data interface{}) {
	if _, ok := d.GetOk("annotation"); ok {
		if _, ok := d.GetOk(relation); ok {
			d.Set(relation, data)
		} else {
			switch data.(type) {
			case string:
				d.Set(relation, "")
			case []string:
				d.Set(relation, make([]string, 0, 1))
			}
		}
	} else {
		d.Set(relation, data)
	}
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

	if !client.ValidateRelationDn {
		return nil
	}
	log.Printf("relation Dns being validated: %v", dns)
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

func validateCommaSeparatedStringInSlice(valid []string, ignoreCase bool, zeroVal string) schema.SchemaValidateFunc {
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
			if val == zeroVal && len(vals) > 1 {
				es = append(es, fmt.Errorf("%s should't be used along with other values in %s", zeroVal, k))
				break
			}
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

func validateColonSeparatedTimeStamp() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		res, err := regexp.MatchString(`^(\d)+(:)(\d){0,2}(:)(\d){0,2}(:)(\d){0,2}(.)(\d){3}$`, v)
		if !res {
			log.Printf("err: %v\n", err)
			es = append(es, fmt.Errorf("Invalid Time Stamp"))
		}
		return
	}
}

func validateNameAttribute() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		c := strings.Contains(v, " ")
		if c {
			es = append(es, fmt.Errorf("property name failed validation for '%s'", v))
		}
		return
	}
}

func validateRemoteFilePath() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		if len(v) >= 1 && v[0] != '/' {
			es = append(es, fmt.Errorf("The first character of remote_path should be '/'"))
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

func G(cont *container.Container, key string) string {
	return StripQuotes(cont.S(key).String())
}

func StripQuotes(word string) string {
	if strings.HasPrefix(word, "\"") && strings.HasSuffix(word, "\"") {
		return strings.TrimSuffix(strings.TrimPrefix(word, "\""), "\"")
	}
	return word
}

func checkAtleastOne() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		if v < "1" {
			es = append(es, fmt.Errorf("Property is out of range"))
		}
		return
	}
}

func validateIntRange(a, b int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		vint, err := strconv.Atoi(v)

		if err != nil {
			es = append(es, err)
			return
		}

		if vint < a || vint > b {
			es = append(es, fmt.Errorf("property is out of range"))
			return
		}
		return
	}
}
