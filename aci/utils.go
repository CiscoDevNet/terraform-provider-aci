package aci

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
	// annotation is one of the base attributes for all terraform resources.
	// we use annotation to check whether function is used for Terraform execution or Terraformer import.
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

func preparePayload(className string, inputMap map[string]string, children []interface{}) (*container.Container, error) {
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
	cont.Array(className, "children")
	if children != nil {
		for _, child := range children {
			childMap := child.(map[string]interface{})
			childClassName := childMap["class_name"].(string)
			childContent := childMap["content"].(map[string]string)

			childCont := container.New()
			childCont.Object(childClassName)
			childCont.Object(childClassName, "attributes")

			for attr, value := range childContent {
				childCont.Set(value, childClassName, "attributes", attr)
			}
			cont.ArrayAppend(childCont.Data(), className, "children")
		}
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

func validateIntBetweenFromString(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected input type of %s to be a string containing an integer", k))
			return warnings, errors
		}

		vint, err := strconv.Atoi(v)

		if err != nil {
			errors = append(errors, err)
			return
		}

		if vint < min || vint > max {
			errors = append(errors, fmt.Errorf("expected %s to be in the range (%d - %d), got %s", k, min, max, v))
			return warnings, errors
		}

		return warnings, errors
	}
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func allowEmpty(err error, allow bool) error {
	if allow && err.Error() == "Error retrieving Object: Object may not exists" {
		return nil
	} else {
		return err
	}
}

// getOldObjectsNotInNew returns elements that are in oldSet but not in newSet, based on the given keyName.
func getOldObjectsNotInNew(keyName string, oldSet, newSet *schema.Set) (oldObjects []interface{}) {
	for _, oldMap := range oldSet.List() {
		found := false
		for _, newMap := range newSet.List() {
			if oldMap.(map[string]interface{})[keyName] == newMap.(map[string]interface{})[keyName] {
				found = true
				break
			}
		}
		if !found {
			oldObjects = append(oldObjects, oldMap)
		}
	}
	return oldObjects
}

// Return the name of the object based on the given Distinguished Name
func getTargetObjectName(paramMap map[string]interface{}, targetDn, targetName string) string {
	if paramMap[targetDn] != "" {
		return GetMOName(paramMap[targetDn].(string))
	} else {
		return paramMap[targetName].(string)
	}
}

// Compares two list-of-strings and sends a new list-of-strings with only unique strings from the first list.
// Reversing the two lists generates a new list with different outputs
func getStringsFromListNotInOtherList(previousValueList interface{}, newValueList interface{}) (generatedList []interface{}) {
	for _, oldValue := range previousValueList.([]interface{}) {
		found := false
		for _, newValue := range newValueList.([]interface{}) {
			if oldValue == newValue {
				found = true
				break
			}
		}
		if !found {
			generatedList = append(generatedList, oldValue)
		}
	}
	return generatedList
}

func errorForObjectNotFound(err error, dn string, d *schema.ResourceData) diag.Diagnostics {
	if strings.HasSuffix(err.Error(), "not found") {
		log.Printf("[WARN] %s, removing from state: %s", err, dn)
		d.SetId("")
		return nil
	} else {
		return diag.FromErr(err)
	}
}
