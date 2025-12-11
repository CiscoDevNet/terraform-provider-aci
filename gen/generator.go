/*
Generates terraform provider code based on templates.

The code assumes that the following directories with content exist in the current working directory (the directory where the generate.go file is located):

- ./definitions (contains the manually created YAML files with the ACI class definitions for overriding the meta data retrieved from APIC)
- ./meta (contains the JSON files with the ACI class metadata which are retrieved from APIC)
- ./templates (contains the Go templates used to generate the full provider code)
	- provider.go.tmpl (the template used to generate the provider.go file in the ../internal/provider directory)
	- index.md.tmpl (the template used to generate the index (provider) documentation file in the ../docs directory)
	- testvars.yaml.tmpl (the template used to generate test variables in the ../testvars directory)
	- annotation_unsupported.go.tmpl (the template used to generate the list of classes that do not support annotation in the ../internal/provider directory)

	- resource.go.tmpl (the template used to generate the resource_*.go files in the ../internal/provider directory)
	- resource.md.tmpl (the template used to generate the *.md files in the ../docs/resources directory)
	- resource_test.go.tmpl (the templates used to generate the *_test.go files in the ../internal/provider directory)
	- resource_test_example.go.tmpl (the templates used to generate the example files used in the documentation which is auto generated with in the ../examples directory)

	- data_source.go.tmpl (the template used to generate the data_source_*.go files in the ../internal/provider directory)
	- data_source.md.tmpl (the template used to generate the *.md files in the ../docs/data-sources directory)
	- data_source_test.go.tmpl (the templates used to generate the *_test.go files in the ../internal/provider directory)
	- data_source_example.go.tmpl (the templates used to generate the example files used in the documentation which is auto generated in the ../examples directory)

- ./testvars (contains the manually created YAML files with the test variables used in the *_test.go and example files)

Usage:
	go run generate.go
*/

package main

import (
	"bytes"
	"cmp"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"

	providerFunctions "github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Paths to the directories containing the files used for generating the provider code and outputting the generated code
const (
	definitionsPath         = "./gen/definitions"
	metaPath                = "./gen/meta"
	templatePath            = "./gen/templates"
	testVarsPath            = "./gen/testvars"
	providerExamplePath     = "./examples/provider/provider.tf"
	resourcesExamplesPath   = "./examples/resources"
	datasourcesExamplesPath = "./examples/data-sources"
	docsPath                = "./docs"
	resourcesDocsPath       = "./docs/resources"
	listResourcesDocsPath   = "./docs/list-resources"
	datasourcesDocsPath     = "./docs/data-sources"
	providerPath            = "./internal/provider/"
)

const providerName = "aci"
const pubhupDevnetBaseUrl = "https://pubhub.devnetcloud.com/media/model-doc-latest/docs"

// Function map used during template rendering in order to call functions from the template
// The map contains a key which is the name of the function used in the template and a value which is the function itself
// The functions itself are defined in the current file
var templateFuncs = template.FuncMap{
	"snakeCase":                                Underscore,
	"validatorString":                          ValidatorString,
	"validatorStringCustomType":                ValidatorStringCustomType,
	"containsString":                           ContainsString,
	"listToString":                             ListToString,
	"overwriteProperty":                        GetOverwriteAttributeName,
	"overwritePropertyValue":                   GetOverwriteAttributeValue,
	"createTestValue":                          func(val string) string { return fmt.Sprintf("test_%s", val) },
	"createNonExistingValue":                   func(val string) string { return fmt.Sprintf("non_existing_%s", val) },
	"getParentTestDependencies":                GetParentTestDependencies,
	"getTestTargetDn":                          GetTestTargetDn,
	"getDefaultValues":                         GetDefaultValues,
	"fromInterfacesToString":                   FromInterfacesToString,
	"containsNoneAttributeValue":               ContainsNoneAttributeValue,
	"legacyAttributeContainsNoneValue":         LegacyAttributeContainsNoneAttributeValue,
	"definedInMap":                             DefinedInMap,
	"getValueFromMap":                          GetValueFromMap,
	"add":                                      func(val1, val2 int) int { return val1 + val2 },
	"subtract":                                 func(val1, val2 int) int { return val1 - val2 },
	"mod":                                      func(val1, val2 int) int { return val1 % val2 },
	"isInterfaceSlice":                         IsInterfaceSlice,
	"replace":                                  Replace,
	"lookupTestValue":                          LookupTestValue,
	"getTestValueOverwrite":                    GetTestValueOverwrite,
	"lookupChildTestValue":                     LookupChildTestValue,
	"createParentDnValue":                      CreateParentDnValue,
	"getResourceName":                          GetResourceName,
	"getResourceNameAsDescription":             GetResourceNameAsDescription,
	"capitalize":                               Capitalize,
	"decapitalize":                             Decapitalize,
	"getTestConfigVariableName":                GetTestConfigVariableName,
	"getDevnetDocForClass":                     GetDevnetDocForClass,
	"getMigrationType":                         GetMigrationType,
	"isLegacyAttribute":                        IsLegacyAttribute,
	"isLegacyChild":                            IsLegacyChild,
	"getLegacyChildAttribute":                  GetLegacyChildAttribute,
	"getConflictingAttributeName":              GetConflictingAttributeName,
	"getAttributeNameForDeprecationMessage":    GetAttributeNameForDeprecationMessage,
	"getPropertyNameForLegacyAttribute":        GetPropertyNameForLegacyAttribute,
	"isNewAttributeStringType":                 IsNewAttributeStringType,
	"isNewNamedClassAttribute":                 IsNewNamedClassAttribute,
	"isRequiredInTestValue":                    IsRequiredInTestValue,
	"getChildAttributesFromBlocks":             GetChildAttributesFromBlocks,
	"getNewChildAttributes":                    GetNewChildAttributes,
	"containsRequired":                         ContainsRequired,
	"hasPrefix":                                strings.HasPrefix,
	"hasCustomTypeDocs":                        HasCustomTypeDocs,
	"containsSingleNestedChildren":             ContainsSingleNestedChildren,
	"containsDeletableSingleNestedChildren":    ContainsDeletableSingleNestedChildren,
	"containsNotDeletableSingleNestedChildren": ContainsNotDeletableSingleNestedChildren,
	"isResourceClass":                          IsResourceClass,
	"getOldType":                               GetOldType,
	"contains":                                 strings.Contains,
	"hasKey":                                   HasKey,
	"definedInList":                            DefinedInList,
	"keyExists":                                KeyExists,
	"isNewNamedClassAttributeMatch":            IsNewNamedClassAttributeMatch,
	"getRnFormat":                              GetRnFormat,
	"identifierIsCustomType":                   IdentifierIsCustomType,
	"isListEmpty":                              func(stringList []string) bool { return len(stringList) == 0 },
	"addToTemplateProperties":                  AddToTemplateProperties,
	"addToChild":                               AddToChildInTestTemplate,
	"checkDeletableChild":                      CheckDeletableChild,
	"emptyChild":                               EmptyChild,
	"excludeForNullInSetCheck":                 ExcludeForNullInSetCheck,
	"getTestTargetValue":                       GetTestTargetValue,
	"isReference":                              IsReference,
	"getLegacyPropertyTestValue":               GetTestValue,
	"getLegacyBlockTestValue":                  GetBlockTestValue,
	"getDeprecatedExplanation":                 GetDeprecatedExplanation,
	"getIgnoredExplanation":                    GetIgnoredExplanation,
	"getCustomTestDependency":                  GetCustomTestDependency,
	"getIgnoreInLegacy":                        GetIgnoreInLegacy,
	"isSensitiveAttribute":                     IsSensitiveAttribute,
	"getChildClassNames":                       GetChildClassNames,
	"getRnForListTesting":                      GetRnForListTesting,
}

func IsSensitiveAttribute(attributeName string, properties map[string]Property) bool {

	if attributeValue, ok := properties[attributeName]; ok {
		if attributeValue.ValueType == "password" {
			return true
		}
	}
	return false
}

func GetChildClassNames(model Model, childClassNames []string) []string {

	if childClassNames == nil {
		childClassNames = []string{}
	}

	for _, child := range model.Children {
		if !slices.Contains(childClassNames, child.PkgName) {
			childClassNames = append(childClassNames, child.PkgName)
		}

		if child.HasChild {
			childClassNames = GetChildClassNames(child, childClassNames)
		}

	}
	return childClassNames
}

func GetDeprecatedExplanation(attributeName, replacedByAttributeName string) string {
	return fmt.Sprintf("Attribute '%s' is deprecated, please refer to '%s' instead. The attribute will be removed in the next major version of the provider.", attributeName, replacedByAttributeName)
}

func GetIgnoredExplanation(attributeName string) string {
	return fmt.Sprintf("Attribute `%s` is deprecated. The configuration was not functioning as intended because the Managed Object (MO) created by the pre-migrated resource was either configured incorrectly or exposed without any implemented functionality on the APIC. The MO for this attribute is no longer created on the APIC, but the existing MO will remain present until the resource is destroyed. This attribute will be removed in the next major version of the provider.", attributeName)
}

func getChildTarget(model Model, childKey, childValue string, tDn bool) string {

	for index, childDependency := range model.ChildTestDependencies {
		targetResourceName := childDependency.TargetResourceName
		if !tDn {
			if childKey == fmt.Sprintf("%s_name", targetResourceName) {
				childKey = "name"
			}
			if result, ok := childDependency.Properties[childKey]; ok && childValue == result {
				return fmt.Sprintf(`aci_%s.test_%s_%d.id`, targetResourceName, targetResourceName, index%2)
			}
		} else {
			if childDependency.Source == childKey || DefinedInList(childDependency.SharedClasses, childKey) {
				if !childDependency.Static {
					return fmt.Sprintf(`aci_%s.test_%s_%d.id`, targetResourceName, targetResourceName, index%2)
				} else {
					return childDependency.TargetDn
				}
			}
		}
	}
	return childValue
}

func GetBlockTestValue(className, attributeName string, model Model) string {
	for _, child := range model.Children {
		if child.PkgName == className {
			if attributeName == "TargetDn" {
				return getChildTarget(model, className, "target_dn_2", true)
			}
			for _, property := range child.Properties {
				if property.Name == attributeName {
					if len(property.ValidValues) > 0 {
						childKey := GetOverwriteAttributeName(model.PkgName, property.SnakeCaseName, model.Definitions)
						return GetOverwriteAttributeValue(child.PkgName, childKey, property.ValidValues[0], "legacy", 0, model.Definitions).(string)
					} else if property.PropertyName == "tDn" {
						return getChildTarget(model, className, "target_dn_1", true)
					}
					childKey := GetOverwriteAttributeName(child.PkgName, property.SnakeCaseName, model.Definitions)
					childValue := GetOverwriteAttributeValue(child.PkgName, childKey, fmt.Sprintf("%s_1", childKey), "default", 0, model.Definitions).(string)
					return getChildTarget(model, childKey, childValue, false)
				}
			}
		}
	}
	return attributeName
}

func GetTestValue(name string, model Model) string {

	for _, property := range model.Properties {
		if property.Name == name {
			childKey := GetOverwriteAttributeName(model.PkgName, property.SnakeCaseName, model.Definitions)
			if len(property.ValidValues) > 0 {
				return GetOverwriteAttributeValue(model.PkgName, childKey, property.ValidValues[0], "legacy", 0, model.Definitions).(string)
			}
			return GetOverwriteAttributeValue(model.PkgName, childKey, fmt.Sprintf("%s_1", childKey), "legacy", 0, model.Definitions).(string)
		}
	}

	for _, child := range model.Children {
		if Capitalize(child.PkgName) == name {
			properties := []string{}
			for _, property := range child.Properties {
				if strings.HasSuffix(property.PropertyName, "Name") {
					childKey := GetOverwriteAttributeName(child.PkgName, property.SnakeCaseName, model.Definitions)
					childValue := GetOverwriteAttributeValue(child.PkgName, childKey, fmt.Sprintf("%s_1", childKey), "default", 0, model.Definitions).(string)
					properties = append(properties, getChildTarget(model, childKey, childValue, false))
				} else if property.PropertyName == "tDn" {
					return getChildTarget(model, child.PkgName, "", true)
				}
			}
			// If there are multiple properties that has suffix of Name, return the one that matches the child resource name
			if len(properties) > 1 {
				for _, property := range properties {
					if strings.Contains(property, string(child.ChildResourceName[strings.Index(child.ChildResourceName, "_to_")+4])) {
						return property
					}
				}
			} else if len(properties) == 1 {
				return properties[0]
			}
		}
	}
	if name == "ParentDn" {
		if len(model.ContainedBy) > 0 {
			parentResource := GetResourceName(model.ContainedBy[0], model.Definitions)
			return fmt.Sprintf(`aci_%s.test.id`, parentResource)
		}
	}
	return name
}

func GetTestTargetValue(targets []interface{}, key string, value interface{}) interface{} {

	for index, target := range targets {

		var resourceName string
		if targetResourceName, ok := target.(map[interface{}]interface{})["target_resource_name"]; ok {
			resourceName = targetResourceName.(string)
		}

		if properties, ok := target.(map[interface{}]interface{})["properties"]; ok {

			if key == fmt.Sprintf("%s_name", resourceName) {
				key = "name"
			}
			if result, ok := properties.(map[interface{}]interface{})[key]; ok && value == result {
				return fmt.Sprintf(`aci_%s.test_%s_%d.%s`, resourceName, resourceName, index%2, key)
			}
		}
	}
	return value
}

func ExcludeForNullInSetCheck(resourceClassName string) bool {
	// Function to exclude TagTag and TagAnnotation from the null check in the Set function
	// Done to reduce the amount of functions created which are not needed for these classes
	// During refactor to struct per class which is reused in children this is not needed anymore
	var childClasses []string
	for _, child := range alwaysIncludeChildren {
		childClasses = append(childClasses, Capitalize(strings.ReplaceAll(child, ":", "")))
	}
	return !slices.Contains(childClasses, resourceClassName)
}

func ContainsRequired(properties map[string]Property) bool {
	for _, property := range properties {
		if property.IsRequired {
			return true
		}
	}
	return false
}

func ContainsSingleNestedChildren(children map[string]Model) bool {
	for _, child := range children {
		if len(child.IdentifiedBy) == 0 || child.MaxOneClassAllowed {
			return true
		}
		if child.Children != nil && ContainsSingleNestedChildren(child.Children) {
			return true
		}
	}
	return false
}

func ContainsNotDeletableSingleNestedChildren(children map[string]Model) bool {
	for _, child := range children {
		if (len(child.IdentifiedBy) == 0 || child.MaxOneClassAllowed) && !child.AllowDelete {
			return true
		}
		if child.Children != nil && ContainsNotDeletableSingleNestedChildren(child.Children) {
			return true
		}
	}
	return false
}

func ContainsDeletableSingleNestedChildren(children map[string]Model) bool {
	for _, child := range children {
		if (len(child.IdentifiedBy) == 0 || child.MaxOneClassAllowed) && child.AllowDelete {
			return true
		}
		if child.Children != nil && ContainsDeletableSingleNestedChildren(child.Children) {
			return true
		}
	}
	return false
}

func CheckDeletableChild(children map[interface{}]interface{}) bool {
	for _, value := range children {
		if val, ok := value.([]interface{}); ok {
			for _, item := range val {
				if child, ok := item.(map[interface{}]interface{}); ok {
					if deletable, exists := child["deletable_child"]; exists && deletable == true {
						return true
					}
					if nestedChildren, exists := child["children"]; exists {
						if nestedChildrenMap, ok := nestedChildren.(map[interface{}]interface{}); ok {
							if CheckDeletableChild(nestedChildrenMap) {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func EmptyChild() map[interface{}]interface{} {
	return make(map[interface{}]interface{})
}

func Replace(oldValue, newValue, inputString string) string {
	return strings.Replace(inputString, oldValue, newValue, -1)
}

func GetMigrationType(valueType []string) string {
	if len(valueType) == 2 && valueType[0] == "set" && valueType[1] == "string" {
		return "Set"
	} else if len(valueType) == 2 && valueType[0] == "list" && valueType[1] == "string" {
		return "List"
	}
	return "String"
}

func GetValueFromMap(key string, searchMap map[string]string) string {
	if value, ok := searchMap[key]; ok {
		return value
	}
	return ""
}

func IsLegacyAttribute(name string, legacyAttributes map[string]LegacyAttribute) bool {
	if _, ok := legacyAttributes[name]; ok {
		return true
	}
	return false
}

func IsLegacyChild(child string, children []string) bool {
	return slices.Contains(children, child)
}

func IsNewAttributeStringType(attributeName string) bool {
	return len(strings.Split(attributeName, ".")) == 1
}

func IsNewNamedClassAttribute(attributeName string) bool {
	return strings.HasSuffix(attributeName, "_name")
}

func IsNewNamedClassAttributeMatch(attributeName, resourceName string) bool {
	// Function to determine the correct named attribute when multiple named attributes are present but only 1 in legacy configuration
	return strings.Contains(resourceName, attributeName[:len(attributeName)-5])
}

func GetAttributeNameForDeprecationMessage(attribute LegacyAttribute, model Model) string {
	if childModel, ok := model.Children[attribute.ReplacedBy.ClassName]; ok && len(strings.Split(attribute.ReplacedBy.AttributeName, ".")) > 1 {
		return fmt.Sprintf("%s.%s", GetOverwriteAttributeName(childModel.PkgName, childModel.ResourceName, model.Definitions), strings.Split(attribute.ReplacedBy.AttributeName, ".")[1])
	}
	return attribute.AttributeName
}

func GetConflictingAttributeName(attributeName string) string {
	return strings.Split(attributeName, ".")[0]
}

func GetPropertyNameForLegacyAttribute(name string, legacyAttributes map[string]LegacyAttribute) string {
	for _, legacyAttribute := range legacyAttributes {
		if legacyAttribute.ReplacedBy.AttributeName == name {
			return legacyAttribute.AttributeName
		}
	}
	return ""
}

func GetRnFormat(rnformat string, identifiers []interface{}) string {
	if len(identifiers) > 0 {
		for _, identifier := range identifiers {
			rnformat = strings.ReplaceAll(rnformat, fmt.Sprintf("{%s}", identifier), "%s")
		}
	}
	return rnformat
}

func IdentifierIsCustomType(identifier string, properties map[string]Property) bool {
	for propertyName, property := range properties {
		if propertyName == identifier {
			return property.HasCustomType
		}
	}
	return false
}

func GetLegacyChildAttribute(className, overwriteProperty string, property Property, legacyAttributes map[string]LegacyAttribute, legacyBlocks []LegacyBlock) string {

	for _, legacyBlock := range legacyBlocks {
		if legacyBlock.ClassName == className {
			// Temporary fix for the issue with the legacy block where the wrong value is returned because the attribute name is also used in the child name and other attribute name in the legacy block
			//  example for this is the relation_from_vrf_to_bgp_address_family_context
			// 	tn_bgp_ctx_af_pol_name: relation_from_vrf_to_bgp_address_family_context.bgp_address_family_context_name
			//  af: relation_from_vrf_to_bgp_address_family_context.address_family
			// when matching for address_family the order of the attribute is of importance, because it was matching on any attribute that contains the string address_family
			// thus when relation_from_vrf_to_bgp_address_family_context.bgp_address_family_context_name it would match on the first attribute in loop, which could be the wrong attribute
			for _, legacyAttribute := range legacyBlock.Attributes {
				attributeName := strings.Split(legacyAttribute.ReplacedBy.AttributeName, ".")
				if len(attributeName) > 1 && attributeName[1] == overwriteProperty {
					return legacyAttribute.Name
				}
			}
			for _, legacyAttribute := range legacyBlock.Attributes {
				if strings.Contains(legacyAttribute.ReplacedBy.AttributeName, overwriteProperty) {
					return legacyAttribute.Name
				}
			}
			return ""
		}
	}

	for _, legacyAttribute := range legacyAttributes {
		if strings.Contains(legacyAttribute.ReplacedBy.AttributeName, overwriteProperty) {
			return Capitalize(className)
		}
	}
	return ""
}

func GetChildAttributesFromBlocks(className string, legacyBlocks []LegacyBlock) map[string]LegacyAttribute {
	legacyAttributes := map[string]LegacyAttribute{}
	for _, legacyBlock := range legacyBlocks {
		if legacyBlock.ClassName == className {
			legacyAttributes = legacyBlock.Attributes
		}
	}
	return legacyAttributes
}

func GetNewChildAttributes(legacyAttributes map[string]LegacyAttribute, properties map[string]Property) []Property {
	result := []Property{}
	for _, property := range properties {
		found := false
		for _, attribute := range legacyAttributes {
			if property.Name == attribute.Name {
				found = true
				break
			}
		}
		if !found && property.Name != "Annotation" {
			result = append(result, property)
		}
	}

	// return result sorted to guarantee consistent output
	slices.SortFunc(result, func(a, b Property) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return result
}

// Global variables used for unique resource name setting based on label from meta data
var labels = []string{"dns_provider", "filter_entry"}
var duplicateLabels = []string{}
var resourceNames = map[string]string{}
var targetRelationalPropertyClasses = map[string]string{}
var alwaysIncludeChildren = []string{"tag:Annotation", "tag:Tag"}
var excludeChildResourceNamesFromDocs = []string{"", "annotation", "tag"}
var classesWithoutResource = []string{"fabricPathEp"}
var resourceIdentifier []interface{}
var alwaysMultiLineTypes = []string{"pkiCert", "pkiPrivateKey"}

func IsResourceClass(className string) bool {
	return !slices.Contains(classesWithoutResource, className)
}

func GetResourceNameAsDescription(s string, definitions Definitions) string {
	resourceName := cases.Title(language.English).String(strings.ReplaceAll(s, "_", " "))
	for k, v := range definitions.Properties["global"].(map[interface{}]interface{})["resource_name_doc_overwrite"].(map[interface{}]interface{}) {
		matchList := strings.Split(resourceName, " ")
		// Always replace when the key is containing of multiple words
		// Replace only individual word on exact match of key, in order to prevent partial word replacement
		if len(strings.Split(k.(string), " ")) > 1 {
			resourceName = strings.ReplaceAll(resourceName, k.(string), v.(string))
		} else if len(matchList) >= 1 && slices.Contains(matchList, k.(string)) {
			resourceName = strings.ReplaceAll(resourceName, k.(string), v.(string))
		}
	}
	return resourceName
}

func GetDevnetDocForClass(className string) string {
	return fmt.Sprintf("[%s](%s/app/index.html#/objects/%s/overview)", className, pubhupDevnetBaseUrl, className)
}

func HasKey(dict map[interface{}]interface{}, key string) bool {
	_, ok := dict[key]
	return ok
}

func KeyExists(m map[string][]string, key string) bool {
	_, exists := m[key]
	return exists
}

func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(s[:1]), s[1:])
}

func Decapitalize(s string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", strings.ToLower(s[:1]), s[1:])
}

func DefinedInList(list interface{}, item string) bool {
	listStr, ok := list.([]interface{})
	if !ok {
		return false
	}

	for _, v := range listStr {
		if v == item {
			return true
		}
	}

	return false
}

func ContainsString(s, sub string) bool {
	return strings.Contains(s, sub)
}

func IsReference(s interface{}) bool {
	if str, ok := s.(string); ok {
		return strings.HasPrefix(str, "aci_") || strings.HasPrefix(str, "data.aci_")
	}
	return false
}

// Reused from https://github.com/buxizhizhoum/inflection/blob/master/inflection.go#L8 to avoid importing the whole package
func Underscore(s string) string {
	for _, reStr := range []string{`([A-Z]+)([A-Z][a-z])`, `([a-z\d])([A-Z])`} {
		re := regexp.MustCompile(reStr)
		s = re.ReplaceAllString(s, "${1}_${2}")
	}
	return strings.ToLower(s)
}

func ValidatorString(stringList []string) string {
	sort.Strings(stringList)
	return fmt.Sprintf("\"%s\"", strings.Join(stringList, "\", \""))
}

func ValidatorStringCustomType(stringList []string, stringMap map[string]string) string {
	newStringList := make([]string, len(stringList))
	copy(newStringList, stringList)
	for key := range stringMap {
		if _, err := strconv.ParseFloat(key, 64); err != nil {
			newStringList = append(newStringList, key)
		}
	}
	return ValidatorString(newStringList)
}

func ListToString(stringList []string) string {
	sort.Strings(stringList)
	return strings.Join(stringList, ",")
}

func isMultiLine(propertyName, classPkgName string, definitions Definitions, modelProperties map[string]Property) bool {
	for _, property := range modelProperties {
		if propertyName == property.SnakeCaseName && slices.Contains(alwaysMultiLineTypes, property.ModelType) {
			return true
		}
	}
	precedenceList := []string{classPkgName, "global"}
	for _, precedence := range precedenceList {
		if classDetails, ok := definitions.Properties[precedence]; ok {
			for key, value := range classDetails.(map[interface{}]interface{}) {
				if key.(string) == "multi_line" {
					for _, v := range value.([]interface{}) {
						if v.(string) == propertyName {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func processMultiLine(multiLineValue string) string {
	cert := strings.ReplaceAll(multiLineValue, "\\n", "\n")
	return fmt.Sprintf(
		`<<EOT
%s
EOT`, cert)
}

func isAttributeATerraformReference(attributeValue string) bool {
	referencePattern := `^(aci_|data\.aci_)\w*\.\w*\.\w*$`
	re := regexp.MustCompile(referencePattern)
	return re.MatchString(attributeValue)
}

// Creates a parent dn value for the resources and datasources in the example files
func CreateParentDnValue(className, caller string, definitions Definitions) string {
	resourceName := GetResourceName(className, definitions)
	return fmt.Sprintf("%s_%s.%s.id", providerName, resourceName, caller)
}

// Retrieves a value for a attribute of a aci class when defined in the testVars YAML file of the class
// Returns "test_value" if no value is defined in the testVars YAML file
func LookupTestValue(classPkgName, originalPropertyName string, testVars map[string]interface{}, definitions Definitions, modelProperties map[string]Property) interface{} {
	var lookupValue interface{} = "test_value"
	propertyName := GetOverwriteAttributeName(classPkgName, originalPropertyName, definitions)

	if allVars, ok := testVars["all"].(map[interface{}]interface{}); ok {
		if val, ok := allVars[propertyName]; ok {
			switch val := val.(type) {
			case string:
				if isMultiLine(originalPropertyName, classPkgName, definitions, modelProperties) {
					lookupValue = processMultiLine(val)
				} else if isAttributeATerraformReference(val) {
					lookupValue = fmt.Sprintf(`%s`, val)
				} else {
					lookupValue = fmt.Sprintf(`"%s"`, val)
				}
			case []interface{}:
				lookupValue = formatSlice(val)
			}
		}

		if versionMismatch, ok := allVars["version_mismatch"].(map[interface{}]interface{}); ok {
			for _, versionVars := range versionMismatch {
				if versionVarsMap, ok := versionVars.(map[interface{}]interface{}); ok {
					if val, ok := versionVarsMap[propertyName]; ok {
						switch val := val.(type) {
						case string:
							if isMultiLine(originalPropertyName, classPkgName, definitions, modelProperties) {
								lookupValue = processMultiLine(val)
							} else if isAttributeATerraformReference(val) {
								lookupValue = fmt.Sprintf(`%s`, val)
							} else {
								lookupValue = fmt.Sprintf(`"%s"`, val)
							}
						case []interface{}:
							lookupValue = formatSlice(val)
						}
					}
				}
			}
		}

		if resourceVars, ok := testVars["resource_required"].(map[interface{}]interface{}); ok {
			if val, ok := resourceVars[propertyName]; ok {
				if strVal, ok := val.(string); ok {
					if isMultiLine(originalPropertyName, classPkgName, definitions, modelProperties) {
						lookupValue = processMultiLine(strVal)
					} else if isAttributeATerraformReference(strVal) {
						lookupValue = fmt.Sprintf(`%s`, strVal)
					} else {
						lookupValue = fmt.Sprintf(`"%s"`, strVal)
					}
				}
			}
		}
	}

	if propertyName == "target_dn" {
		targetResourceName := ""
		resourceName := GetResourceName(classPkgName, definitions)
		if strings.HasPrefix(resourceName, "relation_from_") {
			definitions := getDefinitions().Properties["resource_name_overwrite"].(map[interface{}]interface{})
			if definitions[resourceName] == nil {
				targetResourceName = definitions[strings.TrimSuffix(resourceName, "s")].(string)
			} else {
				targetResourceName = definitions[resourceName].(string)
			}
		} else {
			targetResourceName = strings.TrimPrefix(GetResourceName(classPkgName, definitions), "relation_to_")
		}

		targets, ok := testVars["targets"].([]interface{})
		if ok {
			for _, target := range targets {
				if targetResourceName == target.(map[interface{}]interface{})["relation_resource_name"].(string) {
					return target.(map[interface{}]interface{})["target_dn_ref"].(string)
				}
			}
		}
	}

	// Referencing is done based on target_dn logic
	// This lookup is created as a workaround to reference in an examples on non target_dn attributes
	// Redesign of testing / example creation logic should be done to cover this reference use-case
	if testVars["overwriteTestValue"].(bool) {
		if classDetails, ok := definitions.Properties[classPkgName]; ok {
			for key, value := range classDetails.(map[interface{}]interface{}) {
				if key.(string) == "example_value_overwrite" {
					for k, v := range value.(map[interface{}]interface{}) {
						if k.(string) == propertyName {
							return v.(string)
						}
					}
				}
			}
		}
	}

	return lookupValue
}

func GetTestValueOverwrite(classPkgName, propertyName, propertyValue string, definitions Definitions) string {
	if classDetails, ok := definitions.Properties[classPkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "example_value_overwrite" {
				for k, v := range value.(map[interface{}]interface{}) {
					if k.(string) == propertyName {
						return v.(string)
					}
				}
			}
		}
	}
	return propertyValue
}

// Retrieves a value for a attribute of a aci class when defined in the testVars YAML file of the class
// Returns "test_value_for_child" if no value is defined in the testVars YAML file
func LookupChildTestValue(classPkgName, childResourceName, propertyName string, testVars map[string]interface{}, testValueIndex int, definitions Definitions) interface{} {
	propertyName = GetOverwriteAttributeName(classPkgName, propertyName, definitions)

	overwritePropertyValue := GetOverwriteAttributeValue(classPkgName, propertyName, "", "test_values_for_parent", testValueIndex, definitions)
	if overwritePropertyValue != "" {
		return fmt.Sprintf(`"%s"`, overwritePropertyValue)
	}

	if children, ok := testVars["children"].(map[interface{}]interface{}); ok {
		result := searchChildResources(children, childResourceName, propertyName)
		if result != nil {
			return result
		} else {
			return fmt.Sprintf(`"%s"`, "test_value_for_child")
		}
	}

	return fmt.Sprintf(`"%s_%d"`, propertyName, testValueIndex)
}

func searchChildResources(children map[interface{}]interface{}, childResourceName, propertyName string) interface{} {
	childResources, ok := children[childResourceName].([]interface{})
	if ok && len(childResources) > 0 {
		for _, childResourceInterface := range childResources {
			if childResource, ok := childResourceInterface.(map[interface{}]interface{}); ok {
				if val, ok := childResource[propertyName]; ok {
					switch val := val.(type) {
					case string:
						return fmt.Sprintf(`"%s"`, val)
					case []interface{}:
						return formatSlice(val)
					}
				}

				if versionMismatch, ok := childResource["version_mismatch"].(map[interface{}]interface{}); ok {
					for _, versionVars := range versionMismatch {
						if versionVarsMap, ok := versionVars.(map[interface{}]interface{}); ok {
							if val, ok := versionVarsMap[propertyName]; ok {
								switch val := val.(type) {
								case string:
									return fmt.Sprintf(`"%s"`, val)
								case []interface{}:
									return formatSlice(val)
								}
							}
						}
					}
				}
			}
		}
	} else {
		for _, childResourcesInterfaces := range children {
			if childResources, ok := childResourcesInterfaces.([]interface{}); ok {
				for _, childResourceInterface := range childResources {
					if childResource, ok := childResourceInterface.(map[interface{}]interface{}); ok {
						if nestedChildren, ok := childResource["children"].(map[interface{}]interface{}); ok {
							result := searchChildResources(nestedChildren, childResourceName, propertyName)
							if result != nil {
								return result
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func LegacyAttributeContainsNoneAttributeValue(legacyAttribute LegacyAttribute, properties map[string]Property) bool {
	return ContainsNoneAttributeValue(properties[Decapitalize(legacyAttribute.Name)].ValidValues)
}

func ContainsNoneAttributeValue(values []string) bool {
	return slices.Contains(values, "none")
}

func DefinedInMap(s string, values interface{}) bool {
	if values != nil {
		if _, ok := values.(map[interface{}]interface{})[s]; ok {
			return true
		}
	}
	return false
}

// Create a string from a list of interfaces
func FromInterfacesToString(identifiedBy []interface{}) string {
	var identifiers []string
	for _, identifier := range identifiedBy {
		identifiers = append(identifiers, identifier.(string))
	}
	return fmt.Sprintf("\"%s\"", strings.Join(identifiers, "\", \""))
}

func DictForTemplates(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid number of arguments passed to the dict")
	}
	dict := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// AddToTemplateProperties creates a copy of the model and updates it with the new fields provided in the form of key value pairs
func AddToTemplateProperties(model Model, values ...interface{}) (*Model, error) {
	// Create a copy of the model
	newModel := model
	updates, err := DictForTemplates(values...)
	if err != nil {
		return nil, err
	}

	if newModel.TemplateProperties == nil {
		newModel.TemplateProperties = make(map[string]interface{})
	}

	for k, v := range updates {
		newModel.TemplateProperties[k] = v
	}

	return &newModel, nil
}

// AddToChildInTestTemplate is used within the test templates for applying indentation in the test config
func AddToChildInTestTemplate(child map[interface{}]interface{}, values ...interface{}) (map[interface{}]interface{}, error) {

	newChild := make(map[interface{}]interface{})
	childValue := make(map[interface{}]interface{})
	for k, v := range child {
		childValue[k] = v
	}
	newChild["childValue"] = childValue

	updates, err := DictForTemplates(values...)
	if err != nil {
		return nil, err
	}

	if _, ok := newChild["TemplateProperties"]; !ok {
		newChild["TemplateProperties"] = make(map[string]interface{})
	}

	for k, v := range updates {
		newChild["TemplateProperties"].(map[string]interface{})[k] = v
	}

	return newChild, nil
}

// Renders the templates and writes a file to the output directory
func renderTemplate(templateName, outputFileName, outputPath string, outputData interface{}) {
	templateData, err := os.ReadFile(fmt.Sprintf("%s/%s", templatePath, templateName))
	if err != nil {
		panic(err)
	}
	var buffer bytes.Buffer
	tmpl := template.Must(template.New("").Funcs(templateFuncs).Parse(string(templateData)))
	// The templates have a different data structure, thus based on the template name the output data is casted to the correct type
	if templateName == "provider.go.tmpl" {
		err = tmpl.Execute(&buffer, outputData.(map[string]Model))
	} else if templateName == "index.md.tmpl" {
		err = tmpl.Execute(&buffer, outputData.(ProviderModel))
	} else if strings.Contains(templateName, "_test.go.tmpl") {
		err = tmpl.Execute(&buffer, outputData.(Model).TestVars)
	} else if strings.Contains(templateName, "annotation_unsupported.go.tmpl") {
		err = tmpl.Execute(&buffer, outputData.([]string))
	} else if strings.Contains(templateName, "custom_type.go.tmpl") {
		err = tmpl.Execute(&buffer, outputData.(Property))
	} else {
		err = tmpl.Execute(&buffer, outputData.(Model))
	}
	if err != nil {
		panic(err)
	}
	bytes := buffer.Bytes()
	if strings.Contains(templateName, "go.tmpl") {
		bytes, err = format.Source(buffer.Bytes())
		if err != nil {
			// If the template is not valid Go code, write the code as text to a file for debugging purposes
			// This is done because the error message from the format.Source function is not very helpful
			os.WriteFile(fmt.Sprintf("%s/failed_render.go", outputPath), buffer.Bytes(), 0644)
			panic(err)
		}
	}

	// check if file already exists because all files generated should not exist and cleaned during the cleanDirectories functions
	// when a file already exists, the generation process is stopped to prevent overwriting existing files (main use case is to prevent overwriting migrated documentation files)
	_, error := os.Stat(fmt.Sprintf("%s/%s", outputPath, outputFileName))
	if error == nil {
		panic(fmt.Sprintf("File %s already exists", fmt.Sprintf("%s/%s", outputPath, outputFileName)))
	}

	outputFile, err := os.Create(fmt.Sprintf("%s/%s", outputPath, outputFileName))
	if err != nil {
		panic(err)
	}
	outputFile.Write(bytes)
}

// Creates a map of models for the resources and datasources from the meta data and definitions
func getClassModels(definitions Definitions) map[string]Model {
	files, err := os.ReadDir(metaPath)
	if err != nil {
		panic(err)
	}
	classModels := make(map[string]Model)
	pkgNames := []string{}
	for _, file := range files {
		if path.Ext(file.Name()) != ".json" {
			continue
		}
		pkgNames = append(pkgNames, strings.TrimSuffix(file.Name(), path.Ext(file.Name())))
	}
	for _, pkgName := range pkgNames {
		classModel := Model{PkgName: pkgName}
		classModel.setClassModel(metaPath, false, definitions, []string{}, pkgNames, nil, nil)
		classModels[pkgName] = classModel

		rnName := make(map[string]string)
		rnName[pkgName] = classModel.RnFormat
		resourceIdentifier = append(resourceIdentifier, rnName)
	}
	return classModels
}

// Retrieves the testVars for a model from the testVars YAML file
func getTestVars(model Model) (map[string]interface{}, error) {
	testVarsMap := make(map[string]interface{})
	testVars, err := os.ReadFile(fmt.Sprintf("%s/%s.yaml", testVarsPath, model.PkgName))
	if err != nil {
		return nil, nil
	}
	err = yaml.Unmarshal([]byte(testVars), &testVarsMap)
	if err != nil {
		return nil, err
	}
	// Adds the resource name and resource class name to the testVars map to be used in test template rendering
	testVarsMap["resourceName"] = model.ResourceName
	testVarsMap["resourceClassName"] = model.ResourceClassName
	testVarsMap["pkgName"] = model.PkgName
	// allow access to the model from testvars
	testVarsMap["model"] = model
	// bool to prevent/allow overwriting during value lookup for list resource testing
	testVarsMap["overwriteTestValue"] = true
	return testVarsMap, nil
}

func GetRnForListTesting(targetClasses interface{}, testVars map[string]interface{}) string {
	model := testVars["model"].(Model)
	rn := model.RnFormat
	for _, property := range model.Properties {
		if property.IsNaming {
			var testValue string
			if property.PropertyName == "tDn" {
				testValue = GetTestTargetDn(testVars["targets"].([]interface{}), GetTargetResourceName(model.ResourceName), testValue, false, targetClasses, 0, false)
			} else {
				testVars["overwriteTestValue"] = false
				testValue = LookupTestValue(model.PkgName, property.PropertyName, testVars, model.Definitions, model.Properties).(string)
				if testValue == "test_value" {
					testValue = LookupTestValue(model.PkgName, property.SnakeCaseName, testVars, model.Definitions, model.Properties).(string)
				}
				if testValue == "test_value" {
					testValue = fmt.Sprintf("test_%s", property.SnakeCaseName)
				}
				testVars["overwriteTestValue"] = true
			}
			if len(testValue) > 0 && testValue[0] == '"' {
				testValue = testValue[1 : len(testValue)-1]
			}
			rn = strings.ReplaceAll(rn, fmt.Sprintf("{%s}", property.PropertyName), testValue)
		}
	}

	if len(model.ContainedBy) == 1 {
		defaultParentValue := GetDefaultValues(model.PkgName, "parent_dn", model.Definitions)
		if defaultParentValue != "" && strings.HasPrefix(defaultParentValue, "uni/") {
			rn = fmt.Sprintf("%s/%s", defaultParentValue[4:], rn)
		}
	}

	return rn
}

// Retrieves the property and classs overwrite definitions from the definitions YAML files
func getDefinitions() Definitions {
	definitions := Definitions{}
	files, err := os.ReadDir(definitionsPath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if path.Ext(file.Name()) == ".yaml" {
			definitionMap := make(map[string]interface{})
			definition, err := os.ReadFile(fmt.Sprintf("%s/%s", definitionsPath, file.Name()))
			if err != nil {
				panic(err)
			}
			err = yaml.Unmarshal([]byte(definition), &definitionMap)
			if err != nil {
				panic(err)
			}
			if file.Name() == "classes.yaml" {
				definitions.Classes = definitionMap
			} else if file.Name() == "properties.yaml" {
				definitions.Properties = definitionMap
			}
		} else if path.Ext(file.Name()) == ".json" && strings.Contains(file.Name(), "schema-git-commit-") {
			definitionMap := make(map[string]interface{})
			definition, err := os.ReadFile(fmt.Sprintf("%s/%s", definitionsPath, file.Name()))
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal([]byte(definition), &definitionMap)
			if err != nil {
				panic(err)
			}
			definitions.Migration = definitionMap
		}
	}
	return definitions
}

// Remove all files in a directory except when the files that do not match the ignore list
func cleanDirectory(dir string, ignores []string) {
	names := getFileNames(dir)
	for _, name := range names {
		if !slices.Contains(ignores, name) {
			err := os.RemoveAll(filepath.Join(dir, name))
			if err != nil {
				panic(err)
			}
		}
	}
}

// Returns all the files names in a directory
func getFileNames(dir string) []string {
	d, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	return names
}

// One time function used to migrate the legacy documentation to the new documentation structure
//   - can be removed when legacy documentation is no longer present and no changes
//   - un-generated resources should be updated in old (renamed) documentation location
func migrateLegacyDocumentation() {
	dirPaths := []string{"./legacy-docs/docs/r", "./legacy-docs/docs/d"}
	for _, dirPath := range dirPaths {
		newDirPath := datasourcesDocsPath
		if dirPath == "./legacy-docs/docs/r" {
			newDirPath = resourcesDocsPath
		}
		names := getFileNames(dirPath)
		for _, name := range names {
			newName := strings.Replace(name, ".html.markdown", ".md", 1)
			source, err := os.Open(filepath.Join(dirPath, name))
			if err != nil {
				panic(err)
			}
			defer source.Close()
			destination, err := os.Create(filepath.Join(newDirPath, newName))
			if err != nil {
				panic(err)
			}
			defer destination.Close()
			_, err = io.Copy(destination, source)
			if err != nil {
				panic(err)
			}
		}
	}

}

// Container function to clean all directories properly
func cleanDirectories() {
	cleanDirectory(docsPath, []string{"list-resources", "resources", "data-sources", "guides"})
	cleanDirectory(providerPath, []string{"provider_test.go", "utils.go", "test_constants.go", "resource_aci_rest_managed.go", "resource_aci_rest_managed_test.go", "data_source_aci_rest_managed.go", "data_source_aci_rest_managed_test.go", "annotation_unsupported.go", "data_source_aci_system.go", "data_source_aci_system_test.go", "function_compare_versions.go", "function_compare_versions_test.go"})
	cleanDirectory(resourcesDocsPath, []string{})
	cleanDirectory(listResourcesDocsPath, []string{})
	cleanDirectory(datasourcesDocsPath, []string{"system.md"})
	cleanDirectory(testVarsPath, []string{})
	cleanDirectory("./internal/custom_types", []string{})

	// The *ExamplesPath directories are removed and recreated to ensure all previously rendered files are removed
	// The provider example file is not removed because it contains static provider configuration
	os.RemoveAll(resourcesExamplesPath)
	os.Mkdir(resourcesExamplesPath, 0755)
	cleanDirectory(datasourcesExamplesPath, []string{"aci_system"})

	// Migrate legacy documentation directory ( /website/docs ) format ( .html.markdown ) to new documentation directory ( /docs ) and format ( .md )
	migrateLegacyDocumentation()
}

// Retrieves the example from the example file to be used in the resource and datasource documentation
func getExampleCode(filePath string) []byte {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return content
}

// When RE_GEN_CLASSES environment variable is set, the existing class metadata is retrieved from APIC or the latest devnet docs and stored in the meta directory.
func reGenerateClassMetadata() {
	reGenClasses, err := strconv.ParseBool(os.Getenv("RE_GEN_CLASSES"))
	if err != nil {
		return
	}
	if reGenClasses {
		names := getFileNames(metaPath)
		classNames := strings.Join(names, ",")
		getClassMetadata(classNames)
	}
}

// When GEN_CLASSES environment variable is set, the class metadata is retrieved from the APIC or the latest devnet docs and stored in the meta directory.
func getClassMetadata(classNames string) {
	if classNames != "" {
		var name, nameSpace, url string
		classNameList := strings.Split(classNames, ",")
		for _, className := range classNameList {
			if strings.HasSuffix(className, ".json") {
				className = strings.Replace(className, ".json", "", 1)
			}
			for index, character := range className {
				if unicode.IsUpper(character) {
					nameSpace = className[:index]
					name = className[index:]
					break
				}
			}

			host := os.Getenv("GEN_HOST")
			if host == "" {
				url = fmt.Sprintf("%s/doc/jsonmeta/%s/%s.json", pubhupDevnetBaseUrl, nameSpace, name)
			} else {
				url = fmt.Sprintf("https://%s/doc/jsonmeta/%s/%s.json", host, nameSpace, name)
			}

			client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
			res, err := client.Get(url)
			if err != nil {
				panic(err)
			}

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}

			outputFile, err := os.Create(fmt.Sprintf("%s/%s.json", metaPath, className))
			if err != nil {
				panic(err)
			}
			outputFile.Write(resBody)
		}
	}
}

// When GEN_ANNOTATION_UNSUPPORTED environment variable is set, the list of classes that don't support annotation are retrieved and annotation_unsupported.go is generated.
func generateAnnotationUnsupported() []string {
	classes := []string{}
	genAnnotationUnsupported, err := strconv.ParseBool(os.Getenv("GEN_ANNOTATION_UNSUPPORTED"))
	if err != nil {
		return classes
	}
	if genAnnotationUnsupported {
		var url string
		host := os.Getenv("GEN_HOST")
		if host == "" {
			url = fmt.Sprintf("%s/doc/jsonmeta/aci-meta.json", pubhupDevnetBaseUrl)
		} else {
			url = fmt.Sprintf("https://%s/acimeta/aci-meta.json", host)
		}

		client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		res, err := client.Get(url)
		if err != nil {
			panic(err)
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		var result Metadata
		err = json.Unmarshal(resBody, &result)
		if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
		}

		for class, meta := range result.Classes {
			if meta.IsConfigurable && meta.Properties.Annotation == nil {
				classes = append(classes, class)
			}
		}
		sort.Strings(classes)
	}
	return classes
}

func main() {
	reGenerateClassMetadata()
	getClassMetadata(os.Getenv("GEN_CLASSES"))
	cleanDirectories()

	definitions := getDefinitions()
	classModels := getClassModels(definitions)
	annotationUnsupported := generateAnnotationUnsupported()

	renderTemplate("provider.go.tmpl", "provider.go", providerPath, classModels)
	renderTemplate("index.md.tmpl", "index.md", docsPath, ProviderModel{Example: string(getExampleCode(providerExamplePath))})
	if len(annotationUnsupported) > 0 {
		err := os.Remove(filepath.Join(providerPath, "annotation_unsupported.go"))
		if err != nil {
			panic(err)
		}
		renderTemplate("annotation_unsupported.go.tmpl", "annotation_unsupported.go", providerPath, annotationUnsupported)
	}

	for _, model := range classModels {
		if (len(model.IdentifiedBy) == 0 && !model.Include) || model.Exclude {
			classesWithoutResource = append(classesWithoutResource, model.PkgName)
		}
	}

	for _, model := range classModels {
		// Only render resources and datasources when the class has a unique identifier or is marked as include in the classes definitions YAML file
		// And if the class has a unique identifier, only render when the class is also not a Rs object and is not set to max one class allowed
		// TODO might need to modify this last condition in the future
		if (len(model.IdentifiedBy) > 0 && !(strings.HasPrefix(model.RnFormat, "rs") && model.MaxOneClassAllowed) && !model.Exclude) || model.Include {

			// All classmodels have been read, thus now the model, child and relational resources names can be set
			// When done before additional files would need to be opened and read which would slow down the generation process
			model.ResourceName = GetResourceName(model.PkgName, definitions)
			for _, relationshipClass := range model.RelationshipClasses {
				model.RelationshipResourceNames = append(model.RelationshipResourceNames, GetResourceName(relationshipClass, definitions))
			}
			model.Children = SetChildClassNames(definitions, &model, model.Children)

			if model.VersionMismatched != nil {
				sortVersionMismatched(model.VersionMismatched)
			}

			// Set the documentation specific information for the resource
			// This is done to ensure references can be made to parent/child resources and output amounts can be restricted
			setDocumentationData(&model, definitions)

			// Render the testvars file for the resource
			// First generate run would not mean the file is correct from beginning since some testvars would need to be manually overwritten in the properties definitions YAML file
			model.SetModelTestDependencies(classModels, definitions)
			renderTemplate("testvars.yaml.tmpl", fmt.Sprintf("%s.yaml", model.PkgName), testVarsPath, model)
			testVarsMap, err := getTestVars(model)
			if err != nil {
				panic(err)
			}
			model.TestVars = testVarsMap
			for propertyName, property := range model.Properties {
				if property.HasCustomType {
					renderTemplate("custom_type.go.tmpl", fmt.Sprintf("%s_%s.go", model.PkgName, propertyName), "./internal/custom_types", property)
				}
			}
			renderTemplate("resource.go.tmpl", fmt.Sprintf("resource_%s_%s.go", providerName, model.ResourceName), providerPath, model)
			renderTemplate("resource_list.go.tmpl", fmt.Sprintf("resource_%s_%s_list.go", providerName, model.ResourceName), providerPath, model)
			renderTemplate("datasource.go.tmpl", fmt.Sprintf("data_source_%s_%s.go", providerName, model.ResourceName), providerPath, model)

			os.Mkdir(fmt.Sprintf("%s/%s_%s", resourcesExamplesPath, providerName, model.ResourceName), 0755)
			renderTemplate("provider_example.tf.tmpl", fmt.Sprintf("%s_%s/provider.tf", providerName, model.ResourceName), resourcesExamplesPath, model)
			renderTemplate("resource_example.tf.tmpl", fmt.Sprintf("%s_%s/resource.tf", providerName, model.ResourceName), resourcesExamplesPath, model)
			if !model.HasOnlyRequiredProperties {
				renderTemplate("resource_example_all_attributes.tf.tmpl", fmt.Sprintf("%s_%s/resource-all-attributes.tf", providerName, model.ResourceName), resourcesExamplesPath, model)
				model.ExampleResourceFull = string(hclwrite.Format(getExampleCode(fmt.Sprintf("%s/%s_%s/resource-all-attributes.tf", resourcesExamplesPath, providerName, model.ResourceName))))
			}
			// Leverage the hclwrite package to format the example code
			model.ExampleResource = string(hclwrite.Format(getExampleCode(fmt.Sprintf("%s/%s_%s/resource.tf", resourcesExamplesPath, providerName, model.ResourceName))))
			renderTemplate("resource.md.tmpl", fmt.Sprintf("%s.md", model.ResourceName), resourcesDocsPath, model)
			renderTemplate("resource_list.md.tmpl", fmt.Sprintf("%s.md", model.ResourceName), listResourcesDocsPath, model)

			os.Mkdir(fmt.Sprintf("%s/%s_%s", datasourcesExamplesPath, providerName, model.ResourceName), 0755)
			renderTemplate("provider_example.tf.tmpl", fmt.Sprintf("%s_%s/provider.tf", providerName, model.ResourceName), datasourcesExamplesPath, model)
			renderTemplate("datasource_example.tf.tmpl", fmt.Sprintf("%s_%s/data-source.tf", providerName, model.ResourceName), datasourcesExamplesPath, model)
			// Leverage the hclwrite package to format the example code
			model.ExampleDataSource = string(hclwrite.Format(getExampleCode(fmt.Sprintf("%s/%s_%s/data-source.tf", datasourcesExamplesPath, providerName, model.ResourceName))))
			renderTemplate("datasource.md.tmpl", fmt.Sprintf("%s.md", model.ResourceName), datasourcesDocsPath, model)
			renderTemplate("resource_test.go.tmpl", fmt.Sprintf("resource_%s_%s_test.go", providerName, model.ResourceName), providerPath, model)
			renderTemplate("datasource_test.go.tmpl", fmt.Sprintf("data_source_%s_%s_test.go", providerName, model.ResourceName), providerPath, model)
			renderTemplate("resource_list_test.go.tmpl", fmt.Sprintf("resource_%s_%s_list_test.go", providerName, model.ResourceName), providerPath, model)
		}
	}
}

// A Model that represents the provider
type ProviderModel struct {
	Example string
}

// A Model that represents the ACI Metadata for purposes of finding annotation supported classes.
type Metadata struct {
	Classes map[string]struct {
		IsConfigurable bool `json:"isConfigurable"`
		Properties     struct {
			Annotation interface{} `json:"annotation"`
		} `json:"properties"`
	} `json:"classes"`
}

// A Model represents a ACI class
// All information is retrieved directly or deduced from the metadata
type Model struct {
	PkgName                     string
	Label                       string
	Name                        string
	RnFormat                    string
	RnPrepend                   string
	Comment                     string
	ResourceClassName           string
	ResourceName                string
	ResourceNameDocReference    string
	ChildResourceName           string
	ExampleDataSource           string
	ExampleResource             string
	ExampleResourceFull         string
	SubCategory                 string
	RelationshipClasses         []string
	RelationshipClass           string
	MultiRelationshipClass      bool
	RelationshipResourceNames   []string
	RelationshipResourceName    string
	Versions                    string
	RawVersions                 string
	ChildClasses                []string
	ContainedBy                 []string
	Contains                    []string
	DocumentationDnFormats      []string
	DocumentationParentDns      []string
	DocumentationExamples       []string
	TestDependencies            []TestDependency
	ChildTestDependencies       []TestDependency
	DocumentationChildren       []string
	ResourceNotes               []string
	ResourceWarnings            []string
	DatasourceNotes             []string
	DatasourceWarnings          []string
	Parents                     []string
	UiLocations                 []string
	IdentifiedBy                []interface{}
	MaxOneClassAllowed          bool
	DnFormats                   []interface{}
	Properties                  map[string]Property
	NamedProperties             map[string]Property
	RequiredPropertiesNames     []string
	LegacyAttributes            map[string]LegacyAttribute
	LegacyBlocks                []LegacyBlock
	LegacySchemaVersion         int
	LegacyChildren              []string
	MigrationClassTypes         map[string]string
	Children                    map[string]Model
	Configuration               map[string]interface{}
	TestVars                    map[string]interface{}
	Definitions                 Definitions
	ResourceNameAsDescription   string
	TypeChanges                 map[int][]TypeChange
	SchemaVersion               int
	TestType                    string
	MultiParentFormats          map[string]MultiParentFormat
	MultiParentFormatsTestTypes map[string]string
	ClassVersion                string
	ParentName                  string
	ParentHierarchy             string
	TargetResourceClassName     string
	TargetResourceName          string
	TargetDn                    string
	TargetProperties            map[string]Property
	TargetNamedProperties       map[string]Property
	DirectParent                *Model
	// Below booleans are used during template rendering to determine correct rendering the go code
	AllowDelete                   bool
	AllowChildDelete              bool
	HasBitmask                    bool
	HasChild                      bool
	HasParent                     bool
	HasAnnotation                 bool
	HasValidValues                bool
	HasChildWithoutIdentifier     bool
	HasNaming                     bool
	HasOptionalProperties         bool
	HasOnlyRequiredProperties     bool
	HasNamedProperties            bool
	HasChildNamedProperties       bool
	Include                       bool
	HasReadOnlyProperties         bool
	HasCustomTypeProperties       bool
	Exclude                       bool
	VersionMismatched             map[string][]string
	TemplateProperties            map[string]interface{}
	RequiredAsChild               bool
	ContainsDefaultParentDn       bool
	DataSourceHasNoNameIdentifier bool
}

type TypeChange struct {
	OldType   string
	Attribute string
}

type TestDependency struct {
	Source                string
	SharedClasses         interface{}
	ClassName             string
	ParentDependency      string
	ParentDependencyDnRef string
	ParentDnKey           string
	TargetDn              string
	TargetDnRef           string
	TargetDnOverwriteDocs string
	TargetResourceName    string
	RelationResourceName  string
	Static                bool
	Properties            map[string]string
}

type LegacyBlock struct {
	Name        string
	NestingMode string
	ClassName   string
	Attributes  map[string]LegacyAttribute
}

type LegacyAttribute struct {
	Name            string
	AttributeName   string
	ValueType       []string
	ReplacedBy      ReplacementAttribute
	Optional        bool
	Computed        bool
	Required        bool
	NeedsCustomType bool
}

type ReplacementAttribute struct {
	ClassName     string
	AttributeName string
}

// A Property represents a ACI class property
// All information is retrieved directly or deduced from the metadata
type Property struct {
	Name                     string
	PropertyName             string
	SnakeCaseName            string
	ResourceClassName        string
	PkgName                  string
	ValueType                string
	Label                    string
	Comment                  string
	DefaultValue             string
	Versions                 string
	NamedPropertyClass       string
	IgnoreInTestExampleValue interface{}
	ValidValuesMap           map[string]string
	RawVersion               string
	ModelType                string
	ValidValues              []string
	IdentifiedBy             []interface{}
	Validators               []interface{}
	IdentifyProperties       []Property
	// Below booleans are used during template rendering to determine correct rendering the go code
	IsNaming      bool
	CreateOnly    bool
	IsRequired    bool
	IgnoreInTest  bool
	ReadOnly      bool
	HasCustomType bool
}

// A Definitions represents the ACI class and property definitions as defined in the definitions YAML files
type Definitions struct {
	Classes    map[string]interface{}
	Properties map[string]interface{}
	Migration  map[string]interface{}
}

// Reads the class details from the meta file and sets all details to the Model
func (m *Model) setClassModel(metaPath string, isChildIteration bool, definitions Definitions, parents, pkgNames, mainParentChildren, parentHierarchyList []string) {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s.json", metaPath, m.PkgName))
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var classInfo map[string]interface{}
	err = json.Unmarshal(fileContent, &classInfo)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	m.ResourceClassName = strings.ToUpper(m.PkgName[:1]) + m.PkgName[1:]
	m.Definitions = definitions
	m.Configuration = GetClassConfiguration(m.PkgName, definitions)

	for _, classDetails := range classInfo {
		m.SetClassLabel(classDetails)
		m.SetClassName(classDetails)
		m.SetRelationshipClasses(definitions)
		m.SetClassRnFormat(classDetails)
		m.SetClassRnFormatList(classDetails)
		m.SetClassDnFormats(classDetails)
		m.SetClassIdentifiers(classDetails)
		m.SetClassInclude()
		m.SetClassExclude()
		m.SetClassAllowDelete(classDetails)
		m.SetClassContainedByAndParent(classDetails, parents)
		m.SetClassContains(classDetails)
		m.SetClassComment(classDetails)
		m.SetClassVersions(classDetails)
		m.SetClassProperties(classDetails)
		m.SetClassChildren(classDetails, pkgNames, mainParentChildren)
		if len(parents) != 0 {
			m.SetParentName(parents)
		}
		m.SetResourceNameAsDescription(m.PkgName, definitions)
		m.SetTestType(classDetails, definitions)
		m.SetTestApplicableFromVersion(classDetails)
		m.SetRequiredAsChild(m.PkgName, definitions)
		m.SetDataSourceHasNoNameIdentifier()
	}

	/*
		Checks if the setClassModel is a child class to prevent more than one level of nesting
			- Correct: Parent -> Child
			- Incorrect: Parent -> Child -> Grandchild
		// TODO add grandchild logic
	*/
	m.ParentHierarchy = fmt.Sprintf("%s", strings.Join(reverseList(parentHierarchyList), ""))

	if len(parentHierarchyList) == 0 {
		parentHierarchyList = []string{m.ResourceClassName}
	} else {
		parentHierarchyList = append(parentHierarchyList, m.ResourceClassName)
	}

	if len(m.ChildClasses) > 0 {
		mainParentChildren := append(mainParentChildren, m.ChildClasses...)
		m.HasChild = true
		m.Children = make(map[string]Model)
		for _, child := range m.ChildClasses {
			childModel := Model{PkgName: child}
			childModel.setDirectParent(m)
			childModel.setClassModel(metaPath, true, definitions, []string{m.PkgName}, pkgNames, mainParentChildren, parentHierarchyList)
			m.Children[child] = childModel
			if childModel.HasValidValues {
				m.HasValidValues = true
			}
			if len(childModel.IdentifiedBy) == 0 {
				m.HasChildWithoutIdentifier = true
			}
			if childModel.AllowDelete {
				m.AllowChildDelete = true
			}
			if childModel.HasBitmask {
				m.HasBitmask = true
			}
			if childModel.HasNamedProperties {
				m.HasNamedProperties = true
				m.HasChildNamedProperties = true
			}
		}
	} else {
		m.HasChild = false
	}

	version, changes := isMigrationResource(m.PkgName, definitions)
	if version {
		m.SetMigrationVersion(definitions)
	}
	if changes {
		m.SetMigrationClassTypes(definitions)
		m.SetLegacyChildren(definitions)
		m.SetLegacyAttributes(definitions)
	}

	m.SetStateUpgradeTypeChanges(definitions)
	m.SetResourceNotesAndWarnigns(m.PkgName, definitions)

}

func (m *Model) SetStateUpgradeTypeChanges(definitions Definitions) {
	if v, ok := definitions.Classes[m.PkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "type_changes" {
				m.TypeChanges = make(map[int][]TypeChange)
				for _, value := range value.([]interface{}) {
					var version int
					change := TypeChange{}
					for key, value := range value.(map[interface{}]interface{}) {
						if key.(string) == "old_type" {
							change.OldType = value.(string)
						}
						if key.(string) == "attribute" {
							change.Attribute = value.(string)
						}
						if key.(string) == "version" {
							version = value.(int)
						}
					}
					m.TypeChanges[version] = append(m.TypeChanges[version], change)
					if version >= m.SchemaVersion {
						m.SchemaVersion = version + 1
					}
				}
			}
		}
	}
}

func GetOldType(attributeName string, typeChanges []TypeChange) string {
	for _, change := range typeChanges {
		if change.Attribute == attributeName {
			return change.OldType
		}
	}
	return ""
}

func (m *Model) SetClassLabel(classDetails interface{}) {
	m.Label = cleanLabel(classDetails.(map[string]interface{})["label"].(string))
	if slices.Contains(labels, m.Label) || m.Label == "" {
		if !slices.Contains(duplicateLabels, m.Label) {
			duplicateLabels = append(duplicateLabels, m.Label)
		}
		if _, ok := resourceNames[m.PkgName]; !ok {
			resourceNames[m.PkgName] = Underscore(m.PkgName)
		}
	} else {
		labels = append(labels, m.Label)
		resourceNames[m.PkgName] = m.Label
	}
}

// Remove duplicates from a slice of interfaces
func uniqueInterfaceSlice(interfaceSlice []interface{}) []interface{} {
	keys := make(map[interface{}]bool)
	result := []interface{}{}
	for _, entry := range interfaceSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			result = append(result, entry)
		}
	}
	return result
}

// Remove duplicates from a slice of strings
func uniqueStringSlice(stringSlice []string) []string {
	slices.Sort[[]string](stringSlice)
	return slices.Compact[[]string, string](stringSlice)
}

func cleanLabel(label string) string {
	// Remove all characters that are not allowed in a go variable name
	cleanCharacters := []string{"-", " ", "/", "(", ")", ",", ".", ":", ";", "'", "\"", "!", "?", "&", "%", "$", "#", "@", "+", "="}
	for _, character := range cleanCharacters {
		label = strings.ReplaceAll(label, character, "_")
	}

	// Block attributes are not allowed to start with a number and since the resource name of child classes is based on the label, we need to replace the 802_1x with dot1x
	label = strings.ReplaceAll(label, "802_1x", "dot1x")

	// Remove consecutive underscores from the label
	var returnLabel string
	for i, c := range label {
		if len(label) == i || i == 0 {
			returnLabel = fmt.Sprintf("%s%s", returnLabel, string(c))
		} else if !(string(label[i-1]) == "_" && string(label[i]) == "_") {
			returnLabel = fmt.Sprintf("%s%s", returnLabel, string(c))
		}
	}

	// Remove all capital letters from the label and convert to snake case
	return Underscore(returnLabel)
}

func (m *Model) SetClassName(classDetails interface{}) {
	m.Name = classDetails.(map[string]interface{})["className"].(string)
}

func (m *Model) SetClassRnFormat(classDetails interface{}) {
	m.GetOverwriteRnFormat(classDetails.(map[string]interface{})["rnFormat"].(string))
	if strings.HasPrefix(m.RnFormat, "rs") && len(m.RelationshipClasses) == 0 {
		toMo := classDetails.(map[string]interface{})["relationInfo"].(map[string]interface{})["toMo"].(string)
		m.RelationshipClasses = []string{strings.Replace(toMo, ":", "", 1)}
	}
}

func (m *Model) SetRelationshipClasses(definitions Definitions) {
	overwriteExampleClasses := []interface{}{}
	if v, ok := definitions.Classes[m.PkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "relationship_classes" {
				overwriteExampleClasses = value.([]interface{})
			}
			if key == "multi_relationship_class" {
				// Used in documentation when a relationship can point to multple classes
				m.MultiRelationshipClass = true
			}
		}
	}

	for _, className := range overwriteExampleClasses {
		m.RelationshipClasses = append(m.RelationshipClasses, className.(string))
	}

}

func (m *Model) SetClassDnFormats(classDetails interface{}) {
	for _, dnFormat := range GetOverwriteDnFormats(classDetails.(map[string]interface{})["dnFormats"].([]interface{}), m.PkgName, m.Definitions) {
		m.DnFormats = append(m.DnFormats, dnFormat.(string))
		m.DocumentationDnFormats = append(m.DocumentationDnFormats, dnFormat.(string))
	}
	m.DnFormats = uniqueInterfaceSlice(m.DnFormats)
	m.DocumentationDnFormats = uniqueStringSlice(m.DocumentationDnFormats)
	sort.Strings(m.DocumentationDnFormats)
}

func (m *Model) SetClassIdentifiers(classDetails interface{}) {
	m.IdentifiedBy = uniqueInterfaceSlice(classDetails.(map[string]interface{})["identifiedBy"].([]interface{}))
	m.setMax1Entry()
}

func (m *Model) setMax1Entry() {
	if v, ok := m.Definitions.Classes[m.PkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "max_one_class_allowed" {
				m.MaxOneClassAllowed = value.(bool)
			}
		}
	}
}

func (m *Model) SetClassChildren(classDetails interface{}, pkgNames, mainParentChildren []string) {
	childClasses := []string{}
	excludeChildClasses := []string{}
	if classDetails, ok := m.Definitions.Classes[m.PkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "exclude_children" {
				for _, child := range value.([]interface{}) {
					if !slices.Contains(childClasses, child.(string)) {
						excludeChildClasses = append(excludeChildClasses, child.(string))
					}
				}
			}
		}
	}

	rnMap := classDetails.(map[string]interface{})["rnMap"].(map[string]interface{})
	for rn, className := range rnMap {
		// TODO check if this condition is correct since there might be cases where that we should exclude
		if !strings.HasSuffix(rn, "-") || strings.HasPrefix(rn, "rs") || slices.Contains(alwaysIncludeChildren, className.(string)) {
			pkgName := strings.ReplaceAll(className.(string), ":", "")
			if slices.Contains(pkgNames, pkgName) && !slices.Contains(excludeChildClasses, pkgName) {
				childClasses = append(childClasses, pkgName)
			}
		}
	}
	if classDetails, ok := m.Definitions.Classes[m.PkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "children" {
				for _, child := range value.([]interface{}) {
					if !slices.Contains(childClasses, child.(string)) {
						childClasses = append(childClasses, child.(string))
					}
				}
			}
		}
	}
	m.ChildClasses = uniqueStringSlice(childClasses)
}

func SetChildClassNames(definitions Definitions, model *Model, children map[string]Model) map[string]Model {
	childMap := make(map[string]Model, 0)
	for childName, childModel := range children {
		childModel.ChildResourceName = GetResourceName(childModel.PkgName, definitions)
		childModel.ResourceNameDocReference = childModel.ChildResourceName
		if len(childModel.IdentifiedBy) > 0 && !childModel.MaxOneClassAllowed && !strings.HasSuffix(childModel.ChildResourceName, "s") {
			// TODO add logic to determine the naming for plural child resources
			childModel.ResourceName = fmt.Sprintf("%ss", childModel.ChildResourceName)
		} else {
			childModel.ResourceName = childModel.ChildResourceName
		}
		for _, relationshipClass := range childModel.RelationshipClasses {
			childModel.RelationshipResourceNames = append(childModel.RelationshipResourceNames, GetResourceName(relationshipClass, definitions))
		}

		if len(childModel.VersionMismatched) > 0 {
			sortVersionMismatched(childModel.VersionMismatched)
			updateVersionMismatchedWithChildren(model, childModel.VersionMismatched)
		}
		childModel.Children = SetChildClassNames(definitions, model, childModel.Children)

		childMap[childName] = childModel
	}
	return childMap
}

func (m *Model) SetClassInclude() {
	if classDetails, ok := m.Definitions.Classes[m.PkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "include" {
				m.Include = value.(bool)
			}
		}
	}
}

func (m *Model) SetClassExclude() {
	if classDetails, ok := m.Definitions.Classes[m.PkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "exclude" {
				m.Exclude = value.(bool)
			} else {
				m.Exclude = false
			}
		}
	}
}

func (m *Model) SetClassAllowDelete(classDetails interface{}) {
	if classDetails.(map[string]interface{})["isCreatableDeletable"].(string) == "never" || !AllowClassDelete(m.PkgName, m.Definitions) {
		m.AllowDelete = false
	} else {
		m.AllowDelete = true
	}
}

func (m *Model) SetParentName(classPkgName []string) {
	m.ParentName = classPkgName[0]
}

func (m *Model) setDirectParent(parentModel *Model) {
	m.DirectParent = parentModel
}

func reverseList(items []string) []string {
	reversedList := make([]string, len(items))
	for i, item := range items {
		reversedList[len(items)-1-i] = item
	}
	return reversedList
}

// Determine if a class is allowed to be deleted as defined in the classes.yaml file
// Flag created to ensure classes that only classes allowed to be deleted are deleted
func (m *Model) SetResourceNotesAndWarnigns(classPkgName string, definitions Definitions) {
	if classDetails, ok := definitions.Classes[classPkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "resource_notes" {
				for _, note := range value.([]interface{}) {
					m.ResourceNotes = append(m.ResourceNotes, note.(string))
				}
			}
			if key.(string) == "resource_warnings" {
				for _, note := range value.([]interface{}) {
					m.ResourceWarnings = append(m.ResourceWarnings, note.(string))
				}
			}
			if key.(string) == "datasource_notes" {
				for _, note := range value.([]interface{}) {
					m.DatasourceNotes = append(m.DatasourceNotes, note.(string))
				}
			}
			if key.(string) == "datasource_warnings" {
				for _, note := range value.([]interface{}) {
					m.DatasourceWarnings = append(m.DatasourceWarnings, note.(string))
				}
			}
		}
	}

	if len(m.LegacyAttributes) > 0 {
		m.ResourceWarnings = append(
			m.ResourceWarnings,
			"This resource has been migrated to the terraform plugin protocol version 6, refer to the [migration guide](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/guides/migration) for more details and implications for already managed resources.",
		)
	}
}

func (m *Model) SetResourceNameAsDescription(classPkgName string, definitions Definitions) {
	m.ResourceNameAsDescription = GetResourceNameAsDescription(GetResourceName(classPkgName, definitions), definitions)
}

func (m *Model) SetTestType(classDetails interface{}, definitions Definitions) {
	m.TestType = GetOverwriteTestType(m.PkgName, definitions)
	if m.TestType == "" {
		if platformFlavors, ok := classDetails.(map[string]interface{})["platformFlavors"].([]interface{}); ok {
			for _, value := range platformFlavors {
				if value.(string) == "capic" {
					m.TestType = "cloud"
				} else if value.(string) == "apic" {
					m.TestType = "apic"
				}
			}
			if m.TestType == "" {
				m.TestType = "both"
			}
		}
	}
}

// Determine if a class is allowed to be deleted as defined in the classes.yaml file
// Flag created to ensure classes that only classes allowed to be deleted are deleted
func AllowClassDelete(classPkgName string, definitions Definitions) bool {
	if classDetails, ok := definitions.Classes[classPkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "allow_delete" {
				return value.(bool)
			}
		}
	}
	return true
}

func (m *Model) SetClassContainedByAndParent(classDetails interface{}, parents []string) {
	// Do not include polUni because is excluded because the parent tree in the provider will start at a child of polUni
	// - example fvTenant is contained by polUni, but we want the parent tree to start at fvTenant
	containedExcludes := m.Configuration["contained_by_excludes"].([]string)
	for className := range GetOverwriteContainedBy(classDetails, m.PkgName, m.Definitions) {
		containedclassName := strings.ReplaceAll(className, ":", "")
		if !slices.Contains(containedExcludes, containedclassName) {
			m.ContainedBy = append(m.ContainedBy, containedclassName)
		}
		if slices.Contains(containedExcludes, containedclassName) && containedclassName != "polUni" {
			parents = append(parents, containedclassName)
		}
	}
	m.ContainedBy = uniqueStringSlice(m.ContainedBy)
	sort.Strings(m.ContainedBy)
	if len(m.ContainedBy) == 0 {
		m.HasParent = false
	} else {
		m.HasParent = true
	}
}

func (m *Model) SetClassContains(classDetails interface{}) {
	for className := range classDetails.(map[string]interface{})["contains"].(map[string]interface{}) {
		containclassName := strings.ReplaceAll(className, ":", "")
		m.Contains = append(m.Contains, containclassName)
	}
	m.Contains = uniqueStringSlice(m.Contains)
	sort.Strings(m.Contains)
}

func (m *Model) SetClassComment(classDetails interface{}) {
	var comment string
	metaComment, ok := classDetails.(map[string]interface{})["comment"]
	if ok {
		for _, details := range metaComment.([]interface{}) {
			comment = comment + details.(string)
		}
	}
	m.Comment = comment
}

func (m *Model) SetClassVersions(classDetails interface{}) {
	versions, ok := classDetails.(map[string]interface{})["versions"]
	if ok {
		m.RawVersions = versions.(string)
		m.Versions = formatVersion(versions.(string))
	}
}

func (m *Model) SetTestApplicableFromVersion(classDetails interface{}) {
	m.ClassVersion = GetOverwriteClassVersion(m.PkgName, m.Definitions)
	if m.ClassVersion == "" {
		versions, ok := classDetails.(map[string]interface{})["versions"]
		if ok {
			m.ClassVersion = versions.(string)
		} else {
			m.ClassVersion = "unknown"
		}
	}
}

func (m *Model) SetRequiredAsChild(classPkgName string, definitions Definitions) {
	if classDetails, ok := definitions.Classes[classPkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "required_as_child" {
				m.RequiredAsChild = value.(bool)
				m.AllowDelete = false
			}
		}
	}
}

func formatVersion(versions string) string {
	if versions[len(versions)-1:] == "-" {
		versions = fmt.Sprintf("%s and later.", versions[:len(versions)-1])
	}

	return strings.ReplaceAll(strings.ReplaceAll(versions, ",", ", "), "-", " to ")

}

// Construct a property map for the class, that contains all details of the property that will be used during the rendering of the template
func (m *Model) SetClassProperties(classDetails interface{}) {

	properties := make(map[string]Property)
	namedProperties := make(map[string]Property)
	requiredCount := 0
	readOnlyProperties := readOnlyProperties(m.PkgName, m.Definitions)

	if len(readOnlyProperties) > 0 {
		m.HasReadOnlyProperties = true
	}

	for propertyName, propertyValue := range classDetails.(map[string]interface{})["properties"].(map[string]interface{}) {

		if propertyValue.(map[string]interface{})["isConfigurable"] == true || slices.Contains(readOnlyProperties, propertyName) {

			if ignoreProperty(propertyName, m.PkgName, m.Definitions) {
				continue
			}
			ignoreInTest, ignoreInTestExampleValue := ignoreTestProperty(propertyName, m.PkgName, m.Definitions)

			property := Property{
				Name:                     fmt.Sprintf("%s%s", strings.ToUpper(propertyName[0:1]), propertyName[1:]),
				PropertyName:             propertyName,
				SnakeCaseName:            Underscore(propertyName),
				ResourceClassName:        strings.ToUpper(m.PkgName[:1]) + m.PkgName[1:],
				PkgName:                  m.PkgName,
				IdentifiedBy:             m.IdentifiedBy,
				ValueType:                getOverwritePropertyType(propertyName, m.PkgName, propertyValue.(map[string]interface{})["uitype"].(string), m.Definitions),
				Label:                    propertyValue.(map[string]interface{})["label"].(string),
				IsNaming:                 propertyValue.(map[string]interface{})["isNaming"].(bool),
				CreateOnly:               propertyValue.(map[string]interface{})["createOnly"].(bool),
				IgnoreInTest:             ignoreInTest,
				IgnoreInTestExampleValue: ignoreInTestExampleValue,
				ReadOnly:                 slices.Contains(readOnlyProperties, propertyName),
				HasCustomType:            false,
			}

			if requiredProperty(GetOverwriteAttributeName(m.PkgName, propertyName, m.Definitions), m.PkgName, m.Definitions) || property.IsNaming {
				property.IsRequired = true
				requiredCount += 1
				m.RequiredPropertiesNames = append(m.RequiredPropertiesNames, GetOverwriteAttributeName(m.PkgName, property.SnakeCaseName, m.Definitions))
			}

			if !property.IsRequired {
				m.HasOptionalProperties = true
			}

			if property.ValueType == "bitmask" {
				m.HasBitmask = true
			}

			if property.IsNaming {
				m.HasNaming = true
			}

			commentOverwrite := getOverwritePropertyComment(propertyName, m.PkgName, m.Definitions)
			if commentOverwrite != "" {
				property.Comment = commentOverwrite
			} else {
				val, ok := propertyValue.(map[string]interface{})["comment"]
				if ok {
					space := regexp.MustCompile(`\s+`)
					var comment string
					for _, details := range val.([]interface{}) {
						comment = comment + details.(string)
					}
					property.Comment = space.ReplaceAllString(comment, " ")
				} else {
					property.Comment = property.Label
				}
			}
			// Not all comments end with a dot, this is added to ensure the comment is correctly formatted.
			// TODO consider a generic comment clean up function
			if len(property.Comment) > 0 && property.Comment[len(property.Comment)-1:] != "." {
				property.Comment = fmt.Sprintf("%s.", property.Comment)
			}

			if property.Name == "Annotation" {
				m.HasAnnotation = true
				property.DefaultValue = "orchestrator:terraform"
			}

			val, ok := propertyValue.(map[string]interface{})["validators"]
			if ok {
				property.Validators = val.([]interface{})
			}

			if propertyValue.(map[string]interface{})["validValues"] != nil {
				removedValidValuesList := GetValidValuesToRemove(m.PkgName, propertyName, m.Definitions)
				property.ValidValuesMap = make(map[string]string)
				for _, details := range propertyValue.(map[string]interface{})["validValues"].([]interface{}) {
					validValueName := details.(map[string]interface{})["localName"].(string)
					validValueKey := details.(map[string]interface{})["value"].(string)
					if validValueName != "defaultValue" && !IsInSlice(removedValidValuesList, validValueName) {
						property.ValidValues = append(property.ValidValues, validValueName)
						property.ValidValuesMap[validValueKey] = validValueName
					}
				}
				addValidValuesList := GetValidValuesToAdd(m.PkgName, propertyName, m.Definitions)
				for _, validValueName := range addValidValuesList {
					property.ValidValues = append(property.ValidValues, validValueName.(string))
				}
				if len(property.ValidValues) > 0 {
					m.HasValidValues = true
				}
				if len(property.ValidValuesMap) > 0 && len(property.Validators) > 0 {
					property.HasCustomType = true
					m.HasCustomTypeProperties = true
				}
			}

			defaultValueOverwrite := GetDefaultValues(m.PkgName, propertyName, m.Definitions)
			if defaultValueOverwrite != "" {
				property.DefaultValue = defaultValueOverwrite
			} else {
				val, ok = propertyValue.(map[string]interface{})["default"]
				if ok {
					if reflect.TypeOf(val).String() == "string" {
						// Check if the default value not a valid value in type bitmask where the value is none
						// In the MO this will default the attribute to an empty string
						if property.ValueType == "bitmask" && !slices.Contains(property.ValidValues, val.(string)) && val.(string) == "none" {
							property.DefaultValue = ""
						} else {
							property.DefaultValue = val.(string)
						}
					} else if reflect.TypeOf(val).String() == "float64" {
						property.DefaultValue = fmt.Sprintf("%g", val.(float64))
					} else {
						log.Fatal(fmt.Sprintf("Reflect type %s not defined. Define in SetClassProperties function.", reflect.TypeOf(val).String()))
					}
				}
			}

			versions, ok := propertyValue.(map[string]interface{})["versions"]
			if ok {
				// Check if version is smaller then parent class version if so set to parent class version
				lowClassVersion := providerFunctions.ParseVersion(strings.Split(strings.Split(m.RawVersions, ",")[0], "-")[0]).Version
				lowPropertyVersion := providerFunctions.ParseVersion(strings.Split(strings.Split(versions.(string), ",")[0], "-")[0]).Version

				if providerFunctions.IsVersionGreaterOrEqual(*lowClassVersion, *lowPropertyVersion) {
					versions = m.RawVersions
				}
				property.Versions = formatVersion(versions.(string))
			}

			// The targetRelationalPropertyClasses map is used to store the class name of a named relational property
			// This is stored because in order to determine the correct overwrite property name, the resouce name of the target class is needed
			// The resource name of the target class could be unknown at this point, thus the class name is stored and the resource name is determined later
			// Chose to do this way to prevent constant opening and closing of files to retrieve the resource name
			if strings.HasPrefix(propertyName, "tn") && strings.HasSuffix(propertyName, "Name") {
				// Get the class name and covert to start with lower case for resource name lookup.
				// - Example: tnVzOOBBrCPName -> VzOOBBrCP -> vzOOBBrCP
				className := propertyName[2 : len(propertyName)-4]
				namedClassName := strings.ToLower(className[:1]) + className[1:]
				property.NamedPropertyClass = namedClassName
				targetRelationalPropertyClasses[property.SnakeCaseName] = namedClassName
				namedProperties[propertyName] = property
				m.HasNamedProperties = true
			}

			/* The piece of code below modifies the documentation by adding the version of APIC in which the attributes were introduced. */
			// if propertyValue.(map[string]interface{})["versions"] != nil {
			// 	property.Versions = formatVersion(propertyValue.(map[string]interface{})["versions"].(string))
			// }

			if propertyValue.(map[string]interface{})["versions"] != nil {
				property.RawVersion = propertyValue.(map[string]interface{})["versions"].(string)
			}

			if !property.IgnoreInTest {
				updateVersionMismatched(m, m.Versions, property.RawVersion, property.PropertyName)
			}

			if propertyValue.(map[string]interface{})["modelType"] != nil {
				property.ModelType = strings.Replace(propertyValue.(map[string]interface{})["modelType"].(string), ":", "", -1)
			}

			properties[propertyName] = property

		}

	}

	m.Properties = properties
	m.NamedProperties = namedProperties
	if requiredCount == len(properties) {
		m.HasOnlyRequiredProperties = true
	}
	defaultParentEntry := GetDefaultValues(m.PkgName, "parent_dn", m.Definitions)
	if defaultParentEntry != "" && len(m.ContainedBy) == 1 {
		m.ContainsDefaultParentDn = true
	}
}

func ignoreTestProperty(propertyName, classPkgName string, definitions Definitions) (bool, interface{}) {

	precedenceList := []string{classPkgName, "global"}
	for _, precedence := range precedenceList {
		if classDetails, ok := definitions.Properties[precedence]; ok {
			for key, value := range classDetails.(map[interface{}]interface{}) {
				if key.(string) == "ignore_properties_in_test" {
					for k, v := range value.(map[interface{}]interface{}) {
						if k.(string) == propertyName {
							switch val := v.(type) {
							case []interface{}:
								return true, formatSlice(val)
							default:
								return true, fmt.Sprintf(`"%s"`, val)
							}
						}
					}
				}
			}
		}
	}
	return false, ""
}

func formatSlice(slice []interface{}) string {
	formattedSlice := make([]string, len(slice))
	for i, v := range slice {
		formattedSlice[i] = fmt.Sprintf("\"%v\"", v)
	}
	return fmt.Sprintf("[%v]", strings.Join(formattedSlice, ", "))
}

func updateVersionMismatched(model *Model, classVersion, propertyVersion, propertyName string) {
	classVersionResult := providerFunctions.ParseVersion(classVersion)
	propertyVersionResult := providerFunctions.ParseVersion(propertyVersion)

	if model.VersionMismatched == nil {
		model.VersionMismatched = make(map[string][]string)
	}

	if propertyVersionResult.Error == "unknown" {
		model.VersionMismatched["unknown"] = append(model.VersionMismatched["unknown"], propertyName)
	} else if providerFunctions.IsVersionGreater(*propertyVersionResult.Version, *classVersionResult.Version) && !providerFunctions.IsVersionLesser(*propertyVersionResult.Version, *providerFunctions.ParseVersion("4.0(0a)").Version) {
		model.VersionMismatched[propertyVersion] = append(model.VersionMismatched[propertyVersion], propertyName)
	}
}

func updateVersionMismatchedWithChildren(model *Model, childVersionMap map[string][]string) {

	if model.VersionMismatched == nil {
		model.VersionMismatched = make(map[string][]string)
	}

	for childVersion, properties := range childVersionMap {
		childVersionResult := providerFunctions.ParseVersion(childVersion)
		if childVersionResult.Error == "unknown" {
			for _, property := range properties {
				model.VersionMismatched["unknown"] = append(model.VersionMismatched["unknown"], property)
			}
		} else {
			for _, property := range properties {
				model.VersionMismatched[childVersion] = append(model.VersionMismatched[childVersion], property)
			}
		}
	}
}

func sortVersionMismatched(versionMismatched map[string][]string) {
	for version, properties := range versionMismatched {
		sort.Strings(properties)
		versionMismatched[version] = properties
	}
}

func (m *Model) SetMigrationClassTypes(definitions Definitions) {

	migrationClassTypes := map[string]string{}
	resource_name := fmt.Sprintf("%s_%s", providerName, resourceNames[m.PkgName])
	legacyResource := definitions.Migration["provider_schemas"].(map[string]interface{})["registry.terraform.io/ciscodevnet/aci"].(map[string]interface{})["resource_schemas"].(map[string]interface{})[resource_name]
	if blockTypes, ok := legacyResource.(map[string]interface{})["block"].(map[string]interface{})["block_types"]; ok {
		for blockName, blockDetails := range blockTypes.(map[string]interface{}) {
			className := m.GetClassFromMigrationClassMapping(definitions, blockName, false)
			attributeType := fmt.Sprintf("block,%s", blockDetails.(map[string]interface{})["nesting_mode"].(string))
			migrationClassTypes[className] = attributeType
		}
	}

	attributes := legacyResource.(map[string]interface{})["block"].(map[string]interface{})["attributes"].(map[string]interface{})
	for attributeName, attributeValue := range attributes {
		className := m.GetClassFromMigrationClassMapping(definitions, attributeName, true)
		if className != m.PkgName {
			switch v := attributeValue.(map[string]interface{})["type"].(type) {
			case string:
				migrationClassTypes[className] = v
			case []interface{}:
				attributeTypes := []string{}
				for _, value := range v {
					attributeTypes = append(attributeTypes, value.(string))
				}
				migrationClassTypes[className] = strings.Join(attributeTypes, ",")
			}
		}
	}

	m.MigrationClassTypes = migrationClassTypes
}

func (m *Model) SetLegacyChildren(definitions Definitions) {

	legacyChildren := []string{}
	if v, ok := definitions.Classes[m.PkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "migration_blocks" {
				for className := range value.(map[interface{}]interface{}) {
					if className.(string) != m.PkgName {
						legacyChildren = append(legacyChildren, className.(string))
					}
				}
			}
		}
	}
	m.LegacyChildren = legacyChildren
}

func (m *Model) SetMigrationVersion(definitions Definitions) {
	m.LegacySchemaVersion = 1
	if v, ok := definitions.Classes[m.PkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "migration_version" {
				m.LegacySchemaVersion = value.(int)
			}
		}
	}
}

func (m *Model) SetLegacyAttributes(definitions Definitions) {

	m.LegacyAttributes = make(map[string]LegacyAttribute)
	resourceName := fmt.Sprintf("%s_%s", providerName, GetResourceName(m.PkgName, definitions))

	legacyResource := definitions.Migration["provider_schemas"].(map[string]interface{})["registry.terraform.io/ciscodevnet/aci"].(map[string]interface{})["resource_schemas"].(map[string]interface{})[resourceName]

	if legacyResource != nil {

		attributeNames := []string{}
		for _, property := range m.Properties {
			attributeNames = append(attributeNames, property.SnakeCaseName)
		}

		if blockTypes, ok := legacyResource.(map[string]interface{})["block"].(map[string]interface{})["block_types"]; ok {
			for blockName, blockDetails := range blockTypes.(map[string]interface{}) {

				block := LegacyBlock{
					Name:        blockName,
					NestingMode: blockDetails.(map[string]interface{})["nesting_mode"].(string),
					Attributes:  make(map[string]LegacyAttribute),
				}
				block.ClassName = m.GetClassForBlockName(definitions, block.Name)
				attributes := blockDetails.(map[string]interface{})["block"].(map[string]interface{})["attributes"].(map[string]interface{})
				for attributeName, attributeValue := range attributes {
					legacyAttribute, propertyName := m.GetLegacyAttribute(attributeName, block.ClassName, attributeValue, definitions)

					if legacyAttribute.AttributeName == "target_dn" {
						legacyAttribute.Name = "TargetDn"
					} else if strings.HasSuffix(legacyAttribute.AttributeName, "_dn") {
						legacyAttribute.Name = "TDn"
					}

					childClass := m.Children[block.ClassName]
					for _, property := range childClass.Properties {
						if GetOverwriteAttributeName(m.PkgName, legacyAttribute.AttributeName, definitions) == GetOverwriteAttributeName(m.PkgName, property.SnakeCaseName, definitions) {
							legacyAttribute.Name = property.Name
							break
						}
					}
					block.Attributes[propertyName] = legacyAttribute
				}

				m.LegacyBlocks = append(m.LegacyBlocks, block)
			}
		}

		// sort LegacyBlocks to guarantee consistent output
		slices.SortFunc(m.LegacyBlocks, func(a, b LegacyBlock) int {
			return cmp.Compare(a.ClassName, b.ClassName)
		})

		attributes := legacyResource.(map[string]interface{})["block"].(map[string]interface{})["attributes"].(map[string]interface{})
		for attributeName, attributeValue := range attributes {
			legacyAttribute, propertyName := m.GetLegacyAttribute(attributeName, m.GetClassFromMigrationClassMapping(definitions, attributeName, true), attributeValue, definitions)
			if attributeName == propertyName {
				legacyAttribute.Name = fmt.Sprintf("Ignored_%s", legacyAttribute.Name)
			}
			if legacyAttribute.ReplacedBy.ClassName != "" && legacyAttribute.ReplacedBy.ClassName != m.PkgName {
				legacyAttribute.Name = Capitalize(legacyAttribute.ReplacedBy.ClassName)
			}
			m.LegacyAttributes[propertyName] = legacyAttribute
		}

	}

}

func (m *Model) SetDataSourceHasNoNameIdentifier() {
	if classDetails, ok := m.Definitions.Classes[m.PkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "data_source_has_no_name_identifier" {
				m.DataSourceHasNoNameIdentifier = value.(bool)
			}
		}
	}
}

func (m *Model) GetLegacyAttribute(attributeName, className string, attributeValue interface{}, definitions Definitions) (LegacyAttribute, string) {
	optional := false
	if attributeValue.(map[string]interface{})["optional"] != nil {
		optional = attributeValue.(map[string]interface{})["optional"].(bool)
	}

	computed := false
	if attributeValue.(map[string]interface{})["computed"] != nil {
		computed = attributeValue.(map[string]interface{})["computed"].(bool)
	}

	required := false
	if attributeValue.(map[string]interface{})["required"] != nil {
		required = attributeValue.(map[string]interface{})["required"].(bool)
	}

	valueType := []string{}
	if value, found := attributeValue.(map[string]interface{})["type"].(string); found {
		valueType = []string{value}
	} else {
		for _, value := range attributeValue.(map[string]interface{})["type"].([]interface{}) {
			valueType = append(valueType, value.(string))
		}
	}

	replacedBy := m.GetOverwriteAttributeMigration(definitions, attributeName, className)

	var propertyName string
	if replacedBy != nil {
		propertyName = GetOverwriteAttributeName(m.PkgName, replacedBy.(ReplacementAttribute).AttributeName, definitions)
	} else {
		propertyName = attributeName
	}

	if propertyName == "id_attribute" || propertyName == "id" {
		propertyName = "Id"
	} else if propertyName == "parent_dn" {
		propertyName = "ParentDn"
	} else {
		for _, property := range m.Properties {
			if GetOverwriteAttributeName(m.PkgName, property.SnakeCaseName, definitions) == propertyName {
				propertyName = property.Name
				break
			}
		}
	}

	needsCustomType := false
	for _, property := range m.Properties {
		if propertyName == property.Name && len(property.ValidValuesMap) > 0 && len(property.Validators) > 0 {
			needsCustomType = true
			break
		}

	}

	legacyAttribute := LegacyAttribute{
		Name:            propertyName,
		AttributeName:   attributeName,
		ValueType:       valueType,
		Optional:        optional,
		Computed:        computed,
		Required:        required,
		NeedsCustomType: needsCustomType,
	}

	if replacedBy != nil {
		legacyAttribute.ReplacedBy = replacedBy.(ReplacementAttribute)
	}

	return legacyAttribute, propertyName
}

/*
Determine if a property comment from meta data should be overwritten by a comment overwrite in the properties.yaml file
Precendence order is:
 1. class level from properties.yaml
 2. global level from properties.yaml
*/
func getOverwritePropertyComment(propertyName, classPkgName string, definitions Definitions) string {
	precedenceList := []string{classPkgName, "global"}
	for _, precedence := range precedenceList {
		if classDetails, ok := definitions.Properties[precedence]; ok {
			for key, value := range classDetails.(map[interface{}]interface{}) {
				if key.(string) == "documentation" {
					for k, v := range value.(map[interface{}]interface{}) {
						if k.(string) == propertyName {
							if strings.Contains(v.(string), "%s") {
								return fmt.Sprintf(v.(string), GetResourceNameAsDescription(GetResourceName(classPkgName, definitions), definitions))
							}
							return v.(string)
						}
					}
				}
			}
		}
	}
	return ""
}

func getOverwritePropertyType(propertyName, classPkgName, uiType string, definitions Definitions) string {
	precedenceList := []string{classPkgName, "global"}
	for _, precedence := range precedenceList {
		if classDetails, ok := definitions.Properties[precedence]; ok {
			for key, value := range classDetails.(map[interface{}]interface{}) {
				if key.(string) == "type_overwrites" {
					for k, v := range value.(map[interface{}]interface{}) {
						if k.(string) == propertyName {
							return v.(string)
						}
					}
				}
			}
		}
	}
	return uiType
}

/*
Determine if a property should be ignored as defined in the properties.yaml file
Precendence order is:
 1. class level from properties.yaml
 2. global level from properties.yaml
*/
func ignoreProperty(propertyName, classPkgName string, definitions Definitions) bool {

	precedenceList := []string{classPkgName, "global"}
	for _, precedence := range precedenceList {
		if classDetails, ok := definitions.Properties[precedence]; ok {
			for key, value := range classDetails.(map[interface{}]interface{}) {
				if key.(string) == "ignores" {
					for _, v := range value.([]interface{}) {
						if v.(string) == propertyName {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func readOnlyProperties(classPkgName string, definitions Definitions) []string {
	readOnlyProperties := []string{}
	if classDetails, ok := definitions.Properties[classPkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "read_only_properties" {
				for _, v := range value.([]interface{}) {
					readOnlyProperties = append(readOnlyProperties, v.(string))
				}
			}
		}
	}
	return readOnlyProperties
}

/*
Determine if a property should be required as defined in the properties.yaml file
Precendence order is:
 1. class level from properties.yaml
 2. global level from properties.yaml
*/
func requiredProperty(propertyName, classPkgName string, definitions Definitions) bool {
	precedenceList := []string{classPkgName, "global"}
	for _, precedence := range precedenceList {
		if classDetails, ok := definitions.Properties[precedence]; ok {
			for key, value := range classDetails.(map[interface{}]interface{}) {
				if key.(string) == "resource_required" {
					for _, v := range value.([]interface{}) {
						if v.(string) == propertyName {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

/*
Determine if a attribute name in terraform configuration should be overwritten by a attribute name overwrite in the properties.yaml file
Precendence order is:
 1. class level from properties.yaml
 2. global level from properties.yaml
 3. meta data property attribute name
*/
func GetOverwriteAttributeName(classPkgName, propertyName string, definitions Definitions) string {
	precedenceList := []string{classPkgName, "global"}

	// Overwrite tn..Name attributes with the lookup resource name prepended to _name whe resourcename is known
	// Resource names can be provided in the class definitions YAML file for classed that are not yet generated by the provider
	if strings.HasPrefix(propertyName, "tn_") && strings.HasSuffix(propertyName, "_name") {
		resourceName := GetResourceName(targetRelationalPropertyClasses[propertyName], definitions)
		if resourceName != "" {
			return fmt.Sprintf("%s_name", resourceName)
		}
	}

	for _, precedence := range precedenceList {
		if classDetails, ok := definitions.Properties[precedence]; ok {
			for key, value := range classDetails.(map[interface{}]interface{}) {
				if key.(string) == "overwrites" {
					for k, v := range value.(map[interface{}]interface{}) {
						if k.(string) == propertyName {
							return v.(string)
						}
					}
				}
			}
		}
	}
	return propertyName
}

/*
Determine if a attribute value in testvars should be overwritten by a attribute value overwrite in the properties.yaml file
Precendence order is:
 1. class level from properties.yaml
*/
func GetOverwriteAttributeValue(classPkgName, propertyName, propertyValue, testType string, valueIndex int, definitions Definitions) interface{} {

	if classDetails, ok := definitions.Properties[classPkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "test_values" {
				for test_type, test_type_values := range value.(map[interface{}]interface{}) {
					test_type_value := make(map[interface{}]interface{})
					if test_type.(string) == "test_values_for_parent" && testType == "test_values_for_parent" {
						test_type_value = test_type_values.([]interface{})[valueIndex].(map[interface{}]interface{})
					} else if test_type.(string) == testType {
						test_type_value = test_type_values.(map[interface{}]interface{})
					}
					for k, v := range test_type_value {
						if k.(string) == propertyName {
							switch v := v.(type) {
							case string:
								return v
							case []interface{}:
								return v
							}
						}
					}

				}
			}
		}
	}
	if propertyValue == propertyName {
		index := valueIndex
		if testType == "all" {
			index = valueIndex + 1
		} else if testType == "default" {
			index = valueIndex + 2
		}
		return fmt.Sprintf("%s_%d", propertyValue, index)
	}
	return propertyValue
}

func GetIgnoreInLegacy(classPkgName string, definitions Definitions) []string {
	result := []string{}
	if classDetails, ok := definitions.Properties[classPkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "test_values" {
				for test_type, test_type_values := range value.(map[interface{}]interface{}) {
					if test_type.(string) == "ignore_in_legacy" {
						for _, test_type_value := range test_type_values.([]interface{}) {
							result = append(result, test_type_value.(string))
						}
					}
				}
			}
		}
	}
	return result
}

// Determine if possible parent classes in terraform configuration should be overwritten by contained_by from the classes.yaml file
func GetOverwriteContainedBy(classDetails interface{}, classPkgName string, definitions Definitions) map[string]interface{} {
	containedBy := make(map[string]interface{})
	if v, ok := definitions.Classes[classPkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "contained_by" {
				for _, containedByValue := range value.([]interface{}) {
					containedBy[containedByValue.(string)] = ""
				}
			}
		}
	}
	if len(containedBy) == 0 {
		return classDetails.(map[string]interface{})["containedBy"].(map[string]interface{})
	} else {
		return containedBy
	}
}

// Determine if a reformat in terraform configuration should be prepended with a rn from the classes.yaml file
func (m *Model) GetOverwriteRnFormat(rnFormat string) {
	m.RnFormat = rnFormat
	if v, ok := m.Definitions.Classes[m.PkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "rn_prepend" {
				m.RnFormat = fmt.Sprintf("%s/%s", value.(string), rnFormat)
				m.RnPrepend = value.(string)
			}
		}
	}
}

func GetOverwriteTestType(classPkgName string, definitions Definitions) string {
	if v, ok := definitions.Classes[classPkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "test_type" {
				return value.(string)
			}
		}
	}
	return ""
}

func GetOverwriteClassVersion(classPkgName string, definitions Definitions) string {
	if v, ok := definitions.Classes[classPkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "class_version" {
				return value.(string)
			}
		}
	}
	return ""
}

type MultiParentFormat struct {
	RnPrepend    string
	WrapperClass string
	TestType     string
	ContainedBy  string
	RnFormat     string
}

func GetMultiParentFormats(classPkgName string, definitions Definitions) map[string]MultiParentFormat {
	multiParentFormats := make(map[string]MultiParentFormat)
	if v, ok := definitions.Classes[classPkgName]; ok {
		classMap, ok := v.(map[interface{}]interface{})
		if !ok {
			return multiParentFormats
		}
		multiParentsValue, ok := classMap["multi_parents"]
		if !ok {
			return multiParentFormats
		}
		multiParentsSlice, ok := multiParentsValue.([]interface{})
		if !ok {
			return multiParentFormats
		}
		for _, entry := range multiParentsSlice {
			entryMap, ok := entry.(map[interface{}]interface{})
			if !ok {
				continue
			}
			var multiParentEntry MultiParentFormat
			for key, value := range entryMap {
				switch key {
				case "rn_prepend":
					multiParentEntry.RnPrepend, ok = value.(string)
				case "contained_by":
					multiParentEntry.ContainedBy, ok = value.(string)
				case "test_type":
					multiParentEntry.TestType, ok = value.(string)
				case "wrapper_class":
					multiParentEntry.WrapperClass, ok = value.(string)
				}
				if !ok {
					break
				}
			}
			if ok {
				identifier := GetOverwriteResourceIdentifier(multiParentEntry.ContainedBy, definitions)
				pattern := regexp.MustCompile(`^(.*?)-[\[{]`)
				if identifier == "" {
					for _, item := range resourceIdentifier {
						if rnMap, ok := item.(map[string]string); ok {
							if value, found := rnMap[multiParentEntry.ContainedBy]; found {
								matches := pattern.FindStringSubmatch(value)
								if len(matches) > 1 {
									identifier = matches[1] + "-"
								} else {
									identifier = ""
								}
								break
							}
						}
					}
				}
				multiParentFormats[identifier] = MultiParentFormat{
					ContainedBy:  multiParentEntry.ContainedBy,
					RnPrepend:    multiParentEntry.RnPrepend,
					WrapperClass: multiParentEntry.WrapperClass,
					TestType:     multiParentEntry.TestType,
				}
			}
		}
		defaultParentEntry := GetDefaultValues(classPkgName, "parent_dn", definitions)
		if defaultParentEntry != "" {
			defaultMultiParentFormat := MultiParentFormat{
				ContainedBy: "",
				RnPrepend:   defaultParentEntry,
			}
			multiParentFormats["default"] = defaultMultiParentFormat
		}
	}
	return multiParentFormats
}

func GetOverwriteResourceIdentifier(classPkgName string, definitions Definitions) string {
	if v, ok := definitions.Classes[classPkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "resource_identifier" {
				return value.(string)
			}
		}
	}
	return ""
}

func (m *Model) SetClassRnFormatList(classDetails interface{}) {
	rnFormat := classDetails.(map[string]interface{})["rnFormat"].(string)
	getMultiParentFormats := GetMultiParentFormats(m.PkgName, m.Definitions)
	for key, format := range getMultiParentFormats {
		format.RnFormat = rnFormat
		getMultiParentFormats[key] = format
	}
	m.MultiParentFormats = getMultiParentFormats
}

// Determine if possible dn formats in terraform documentation should be overwritten by dn formats from the classes.yaml file
func GetOverwriteDnFormats(dnFormats []interface{}, classPkgName string, definitions Definitions) []interface{} {
	if v, ok := definitions.Classes[classPkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "dn_formats" {
				dnFormats = value.([]interface{})
			}
		}
	}
	return dnFormats
}

// GetTestConfigVariableName generates a test configuration variable name based on the provided class name, middle string, and parent class name.
func GetTestConfigVariableName(className string, midString string, parentClassName string) string {
	if className != "" {
		return fmt.Sprintf("testConfig%s%s%s", className, midString, Capitalize(parentClassName))
	}
	return ""
}

// Determine if possible dn formats in terraform documentation should be overwritten by dn formats from the classes.yaml file
func GetOverwriteExampleClasses(classPkgName string, definitions Definitions) []interface{} {
	overwriteExampleClasses := []interface{}{}
	if v, ok := definitions.Classes[classPkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "example_classes" {
				overwriteExampleClasses = value.([]interface{})
			}
		}
	}
	return overwriteExampleClasses
}

func GetParentTestDependencies(classPkgName string, index int, definitions Definitions) map[string]interface{} {
	parentDependency := ""
	classInParent := false
	parentDependencyName := ""
	targetClasses := []string{}

	if classDetails, ok := definitions.Properties[classPkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "parents" {
				if len(value.([]interface{})) > index {
					parentMap := value.([]interface{})[index].(map[interface{}]interface{})
					if pd, ok := parentMap["parent_dependency"]; ok {
						parentDependency = pd.(string)
					}
					if cip, ok := parentMap["class_in_parent"]; ok {
						classInParent = cip.(bool)
					}
					if pdn, ok := parentMap["parent_dependency_name"]; ok {
						parentDependencyName = pdn.(string)
					}
					if tc, ok := parentMap["target_classes"]; ok {
						results := []string{}
						for _, targetClass := range tc.([]interface{}) {
							results = append(results, targetClass.(string))
						}
						sort.Strings(results)
						targetClasses = results
					}
				}
			}
		}
	}
	return map[string]interface{}{"parent_dependency": parentDependency, "class_in_parent": classInParent, "parent_dependency_name": parentDependencyName, "target_classes": targetClasses}
}

func (m *Model) SetModelTestDependencies(classModels map[string]Model, definitions Definitions) {

	testDependencies := []TestDependency{}
	if classDetails, ok := definitions.Properties[m.PkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "targets" {
				for index, v := range value.([]interface{}) {
					targetMap := v.(map[interface{}]interface{})
					if className, ok := targetMap["class_name"]; ok && !slices.Contains(m.getExcludeTargets(), className.(string)) {
						testDependencies = append(testDependencies, getTestDependency(m.PkgName, className.(string), targetMap, definitions, index))
					}
				}
			}
		}
	}

	childTestDependencies := []TestDependency{}
	for _, child := range m.ChildClasses {
		if classDetails, ok := definitions.Properties[child]; ok {
			for key, value := range classDetails.(map[interface{}]interface{}) {
				if key.(string) == "targets" {
					for index, v := range value.([]interface{}) {
						targetMap := v.(map[interface{}]interface{})
						if className, ok := targetMap["class_name"]; ok && !slices.Contains(m.getExcludeTargets(), className.(string)) {
							childTestDependencies = append(childTestDependencies, getTestDependency(child, className.(string), targetMap, definitions, index))
						}
					}
				}
			}
		}
	}

	m.TestDependencies = testDependencies
	m.ChildTestDependencies = childTestDependencies

}

func (m *Model) getExcludeTargets() []string {
	excludeTargets := []string{}
	if v, ok := m.Definitions.Properties[m.PkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "exclude_targets" {
				for _, v := range value.([]interface{}) {
					excludeTargets = append(excludeTargets, v.(string))
				}
			}
		}
	}
	return excludeTargets
}

func getTestDependency(sourceClassName, className string, targetMap map[interface{}]interface{}, definitions Definitions, index int) TestDependency {

	testDependency := TestDependency{}
	testDependency.Source = sourceClassName
	testDependency.ClassName = className
	testDependency.TargetResourceName = GetResourceName(className, definitions)
	testDependency.RelationResourceName = testDependency.TargetResourceName
	testDependency.ParentDnKey = "parent_dn"

	if relationResourceName, ok := targetMap["relation_resource_name"]; ok {
		testDependency.RelationResourceName = relationResourceName.(string)
	}

	if sharedClasses, ok := targetMap["shared_classes"]; ok {
		testDependency.SharedClasses = sharedClasses
	}

	if parentDependency, ok := targetMap["parent_dependency"]; ok {
		testDependency.ParentDependency = parentDependency.(string)
		if parentDependencyDn, ok := targetMap["parent_dependency_dn_ref"]; ok {
			testDependency.ParentDependencyDnRef = parentDependencyDn.(string)
		} else {
			testDependency.ParentDependencyDnRef = fmt.Sprintf("%s_%s.test.id", providerName, GetResourceName(parentDependency.(string), definitions))
		}
	}

	if targetDn, ok := targetMap["target_dn"]; ok {
		testDependency.TargetDn = targetDn.(string)
		testDependency.TargetDnRef = fmt.Sprintf("%s_%s.test_%s_%d.id", providerName, GetResourceName(className, definitions), GetResourceName(className, definitions), index%2)
	}

	if properties, ok := targetMap["properties"]; ok {
		testDependency.Properties = make(map[string]string)
		for name, value := range properties.(map[interface{}]interface{}) {
			testDependency.Properties[name.(string)] = value.(string)
		}
	}

	if overwriteKey, ok := targetMap["overwrite_parent_dn_key"]; ok {
		testDependency.ParentDnKey = overwriteKey.(string)
	}

	testDependency.Static = false
	if static, ok := targetMap["static"]; ok {
		testDependency.Static = static.(bool)
	}

	if targetDnOverwriteDocs, ok := targetMap["target_dn_overwrite_docs"]; ok {
		testDependency.TargetDnOverwriteDocs = targetDnOverwriteDocs.(string)
	}

	return testDependency
}

func GetTargetResourceName(resourceName string) string {
	definitions := getDefinitions().Properties["resource_name_overwrite"].(map[interface{}]interface{})
	if definitions[resourceName] != nil {
		return definitions[resourceName].(string)
	} else if definitions[resourceName] == nil {
		if definitions[strings.TrimSuffix(resourceName, "s")] != nil {
			return definitions[strings.TrimSuffix(resourceName, "s")].(string)
		} else {
			return strings.TrimPrefix(resourceName, "relation_to_")
		}
	}
	return resourceName
}

func GetTestTargetDn(targets []interface{}, resourceName, targetDnValue string, reference bool, targetClasses interface{}, index int, overwriteDocs bool) string {
	var filteredTargets []interface{}
	targetResourceName := GetTargetResourceName(resourceName)

	if targetClasses != nil {
		// CHANGE logic here when allowing for multiple target classes in single resource
		targetClass := targetClasses.([]interface{})[0].(string)

		for _, target := range targets {
			if targetClass == target.(map[interface{}]interface{})["class_name"].(string) {
				filteredTargets = append(filteredTargets, target)
			}
		}
	} else {
		filteredTargets = targets
	}

	for _, target := range filteredTargets {

		targetRelationResourceName := target.(map[interface{}]interface{})["relation_resource_name"].(string)

		static := false
		if v, ok := target.(map[interface{}]interface{})["static"]; ok {
			static = v.(bool)
		}

		if targetResourceName == targetRelationResourceName || strings.TrimSuffix(targetResourceName, "s") == targetRelationResourceName {
			if index > 0 {
				index = index - 1
			} else {
				if reference && !static {
					return target.(map[interface{}]interface{})["target_dn_ref"].(string)
				}
				if v, ok := target.(map[interface{}]interface{})["target_dn_overwrite_docs"].(string); ok && overwriteDocs && v != "" {
					return v
				}
				return target.(map[interface{}]interface{})["target_dn"].(string)
			}
		}
	}

	if targetDnValue == "" && overwriteDocs {
		for _, target := range filteredTargets {
			if v, ok := target.(map[interface{}]interface{})["target_dn_overwrite_docs"].(string); ok && overwriteDocs && v != "" {
				return v
			}
		}
	} else if targetDnValue == "test_t_dn" {
		for _, target := range filteredTargets {
			if reference {
				if v, ok := target.(map[interface{}]interface{})["target_dn_ref"].(string); ok && v != "" {
					return v
				}
			} else {
				if v, ok := target.(map[interface{}]interface{})["target_dn"].(string); ok && v != "" {
					return v
				}
			}

		}
	}

	return targetDnValue
}

func (m *Model) GetClassForBlockName(definitions Definitions, blockName string) string {
	if v, ok := definitions.Classes[m.PkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "migration_blocks" {
				for className, classValues := range value.(map[interface{}]interface{}) {
					for migrationKey := range classValues.(map[interface{}]interface{}) {
						if migrationKey.(string) == blockName {
							return className.(string)
						}
					}
				}
			}
		}
	}
	panic(fmt.Sprintf("Class name not found for block %s", blockName))
}

func (m *Model) GetClassFromMigrationClassMapping(definitions Definitions, attributeName string, set bool) string {
	if v, ok := definitions.Classes[m.PkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "migration_blocks" {
				for className, attributes := range value.(map[interface{}]interface{}) {
					if set {
						if className.(string) != m.PkgName && len(attributes.(map[interface{}]interface{})) == 1 {
							for name := range attributes.(map[interface{}]interface{}) {
								if name.(string) == attributeName {
									return className.(string)
								}
							}
						}
					} else {
						for name := range attributes.(map[interface{}]interface{}) {
							if name.(string) == attributeName {
								return className.(string)
							}
						}
					}
				}
			}
		}
	}
	return m.PkgName
}

func (m *Model) GetOverwriteAttributeMigration(definitions Definitions, attributeName, className string) interface{} {
	if v, ok := definitions.Classes[m.PkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "migration_blocks" {
				for classNameDefinition, classValues := range value.(map[interface{}]interface{}) {
					if classNameDefinition.(string) == className {
						for migrationKey, migrationValue := range classValues.(map[interface{}]interface{}) {
							if migrationKey.(string) == attributeName {
								return ReplacementAttribute{AttributeName: migrationValue.(string), ClassName: className}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func isMigrationResource(classPkgName string, definitions Definitions) (bool, bool) {
	version := false
	changes := false
	if v, ok := definitions.Classes[classPkgName]; ok {
		for key := range v.(map[interface{}]interface{}) {
			if key.(string) == "migration_blocks" {
				changes = true
			} else if key.(string) == "migration_version" {
				version = true
			}
		}
	}
	return version, changes
}

// Set variables that are used during the rendering of the example and documentation templates
func setDocumentationData(m *Model, definitions Definitions) {
	UiLocations := []string{}
	SubCategory := "Generic"
	if v, ok := definitions.Classes[m.PkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "ui_locations" {
				for _, UiLocation := range value.([]interface{}) {
					UiLocations = append(UiLocations, UiLocation.(string))
				}
			} else if key.(string) == "sub_category" {
				SubCategory = value.(string)
			}
		}
	}
	if len(UiLocations) == 0 {
		UiLocations = append(UiLocations, "Generic")
	}
	m.UiLocations = UiLocations
	m.SubCategory = SubCategory

	resourcesFound := [][]string{}
	resourcesNotFound := []string{}
	for _, containedClassName := range m.ContainedBy {
		resourceName := GetResourceName(containedClassName, definitions)
		if resourceName != "" && !slices.Contains(classesWithoutResource, containedClassName) {
			resourcesFound = append(resourcesFound, []string{resourceName, containedClassName})
		} else {
			resourcesNotFound = append(resourcesNotFound, containedClassName)
		}
	}

	docsParentDnAmount := m.Configuration["docs_parent_dn_amount"].(int)

	if len(m.DocumentationDnFormats) > docsParentDnAmount {
		m.DocumentationDnFormats = append([]string{fmt.Sprintf("Too many DN formats to display, see model documentation for all possible parents of %s.", GetDevnetDocForClass(m.PkgName))}, m.DocumentationDnFormats[0:docsParentDnAmount]...)
	}

	if len(resourcesFound) > docsParentDnAmount {
		for _, resourceDetails := range resourcesFound[0:docsParentDnAmount] {
			m.DocumentationParentDns = append(m.DocumentationParentDns, fmt.Sprintf("[%s_%s](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/%s) (%s)", providerName, resourceDetails[0], resourceDetails[0], GetDevnetDocForClass(resourceDetails[1])))
		}
		m.DocumentationParentDns = append([]string{fmt.Sprintf("Too many parent DNs to display, see model documentation for all possible parents of %s.", GetDevnetDocForClass(m.PkgName))}, m.DocumentationParentDns...)
	} else {
		for _, resourceDetails := range resourcesFound {
			m.DocumentationParentDns = append(m.DocumentationParentDns, fmt.Sprintf("[%s_%s](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/%s) (%s)", providerName, resourceDetails[0], resourceDetails[0], GetDevnetDocForClass(resourceDetails[1])))

		}
	}

	if len(resourcesNotFound) != 0 && len(resourcesFound) < docsParentDnAmount {
		if len(resourcesNotFound) > docsParentDnAmount-len(resourcesFound) {
			// TODO catch default classes and add to documentation
			resourcesNotFound = resourcesNotFound[0:(docsParentDnAmount - len(resourcesFound))]
			m.DocumentationParentDns = append(m.DocumentationParentDns, fmt.Sprintf("Too many classes to display, see model documentation for all possible classes of %s.", GetDevnetDocForClass(m.PkgName)))
		} else {
			var resourceDetails string
			for _, resource := range resourcesNotFound {
				resourceDetails = fmt.Sprintf("%s    - %s\n", resourceDetails, GetDevnetDocForClass(resource))
			}
			m.DocumentationParentDns = append(m.DocumentationParentDns, fmt.Sprintf("The distinguished name (DN) of classes below can be used but currently there is no available resource for it:\n%s", resourceDetails))
		}
	}

	getMultiParentFormats := GetMultiParentFormats(m.PkgName, definitions)
	if len(getMultiParentFormats) > 0 {
		m.DocumentationParentDns = nil
		for _, format := range getMultiParentFormats {
			if format.ContainedBy != "" {
				m.DocumentationParentDns = append(m.DocumentationParentDns, fmt.Sprintf("[%s_%s](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/%s) (%s)", providerName, GetResourceName(format.ContainedBy, definitions), format.ContainedBy, GetDevnetDocForClass(format.ContainedBy)))
			}
		}
	}

	// TODO add overwrite to provide which documentation examples to be included
	docsExampleAmount := m.Configuration["docs_examples_amount"].(int)
	if len(m.ContainedBy) > docsExampleAmount {
		overwriteExampleClasses := GetOverwriteExampleClasses(m.PkgName, definitions)
		if len(overwriteExampleClasses) > 0 {
			for _, exampleClass := range overwriteExampleClasses {
				m.DocumentationExamples = append(m.DocumentationExamples, exampleClass.(string))
			}
		} else {
			for _, resourceDetails := range resourcesFound[0:docsExampleAmount] {
				m.DocumentationExamples = append(m.DocumentationExamples, resourceDetails[1])
			}
		}
	} else {
		for _, resourceDetails := range resourcesFound {
			m.DocumentationExamples = append(m.DocumentationExamples, resourceDetails[1])
		}
	}

	// Add child class references to documentation when resource name is known
	for _, child := range m.Contains {
		match, _ := regexp.MatchString("[Rs][A-Z][^\r\n\t\f\v]", child) // match all Rs classes
		if !match && !slices.Contains(m.ChildClasses, child) {
			resourceName := GetResourceName(child, definitions)
			if !slices.Contains(excludeChildResourceNamesFromDocs, resourceName) { // exclude anotation children since they will be included into the resource when possible
				m.DocumentationChildren = append(m.DocumentationChildren, fmt.Sprintf("[%s_%s](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/%s)", providerName, resourceName, resourceName))
			}
		}
	}
	sort.Strings(m.DocumentationChildren)
}

// Determine if a resource / datasource name in terraform configuration should be overwritten by a resource name overwrite in the classes.yaml file
// When a manual overwrite is not set the className (transformed from camel case to snake case) will be used
// - example fvTenant will become fv_tenant
// TODO determine way to handle the creation of resource name
func GetResourceName(classPkgName string, definitions Definitions) string {
	if v, ok := definitions.Classes[classPkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "resource_name" {
				return value.(string)
			}
		}
	}
	return resourceNames[classPkgName]
}

func GetClassConfiguration(classPkgName string, definitions Definitions) map[string]interface{} {
	classConfiguration := make(map[string]interface{})
	reversePrecedenceList := []string{"global", classPkgName}
	for _, precedence := range reversePrecedenceList {
		if classDetails, ok := definitions.Classes[precedence]; ok {
			for key, value := range classDetails.(map[interface{}]interface{}) {
				if key.(string) == "contained_by_excludes" {
					containedByExcludes := []string{}
					for _, containedByExcludeValue := range value.(interface{}).([]interface{}) {
						containedByExcludes = append(containedByExcludes, containedByExcludeValue.(string))
					}
					classConfiguration["contained_by_excludes"] = containedByExcludes
				} else if key.(string) == "docs_examples_amount" {
					classConfiguration["docs_examples_amount"] = value.(int)
				} else if key.(string) == "docs_parent_dn_amount" {
					classConfiguration["docs_parent_dn_amount"] = value.(int)
				}
			}
		}
	}
	return classConfiguration
}

/*
Determine if a attribute default value in terraform configuration should be overwritten/or defined by a attribute default value overwrite in the properties.yaml file.
It can be assigned to "parent_dn".
Precendence order is:
 1. class level from properties.yaml
 2. global level from properties.yaml
 3. meta data property default value
*/
func GetDefaultValues(classPkgName, propertyName string, definitions Definitions) string {
	precedenceList := []string{classPkgName, "global"}
	for _, precedence := range precedenceList {
		if classDetails, ok := definitions.Properties[precedence]; ok {
			for key, value := range classDetails.(map[interface{}]interface{}) {
				if key.(string) == "default_values" {
					for k, v := range value.(map[interface{}]interface{}) {
						if k.(string) == propertyName {
							switch v := v.(type) {
							case string:
								return v
							case float64:
								return strconv.FormatFloat(v, 'f', -1, 64)
							}
						}
					}
				}
			}
		}
	}
	return ""
}

func IsInterfaceSlice(input interface{}) bool {
	_, ok := input.([]interface{})
	return ok
}

func IsInSlice(slice []interface{}, element interface{}) bool {
	if len(slice) == 0 {
		return false
	}
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

func GetValidValuesToRemove(classPkgName, propertyName string, definitions Definitions) []interface{} {
	precedenceList := []string{classPkgName, "global"}
	removedValidValuesSlice := make([]interface{}, 0)
	for _, precedence := range precedenceList {
		if classDetails, ok := definitions.Properties[precedence]; ok {
			for key, value := range classDetails.(map[interface{}]interface{}) {
				if key.(string) == "remove_valid_values" {
					for k, v := range value.(map[interface{}]interface{}) {
						if k.(string) == propertyName {
							removedValidValuesSlice = append(removedValidValuesSlice, v.([]interface{})...)
						}
					}
				}
			}
		}
	}
	return removedValidValuesSlice
}

func GetValidValuesToAdd(classPkgName, propertyName string, definitions Definitions) []interface{} {
	precedenceList := []string{classPkgName, "global"}
	addValidValuesSlice := make([]interface{}, 0)
	for _, precedence := range precedenceList {
		if classDetails, ok := definitions.Properties[precedence]; ok {
			for key, value := range classDetails.(map[interface{}]interface{}) {
				if key.(string) == "add_valid_values" {
					for k, v := range value.(map[interface{}]interface{}) {
						if k.(string) == propertyName {
							addValidValuesSlice = append(addValidValuesSlice, v.([]interface{})...)
						}
					}
				}
			}
		}
	}
	return addValidValuesSlice
}

func IsRequiredInTestValue(classPkgName, propertyName string, definitions Definitions, testType string) bool {
	if classDetails, ok := definitions.Properties[classPkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "test_values" {
				for test_type, test_type_values := range value.(map[interface{}]interface{}) {
					if test_type.(string) == testType {
						for k := range test_type_values.(map[interface{}]interface{}) {
							if k.(string) == propertyName {
								return true
							}
						}
					}
				}

			}
		}
	}
	return false
}

func HasCustomTypeDocs(classPkgName, propertyName string, definitions Definitions) bool {
	if classDetails, ok := definitions.Properties[classPkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "ignore_custom_type_docs" {
				for _, property := range value.([]interface{}) {
					if property.(string) == propertyName {
						return false
					}
				}

			}
		}
	}
	return true
}

func GetCustomTestDependency(classPkgName string, definitions Definitions) string {
	if classDetails, ok := definitions.Properties[classPkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "custom_test_dependency_name" {
				return value.(string)
			}
		}
	}
	return ""
}
