//go:build ignore

/*
Generates terraform provider code based on templates.

The code assumes that the following directories with content exist in the current working directory (the directory where the generate.go file is located):

- ./definitions (contains the manually created YAML files with the ACI class definitions for overriding the meta data retrieved from APIC)
- ./meta (contains the JSON files with the ACI class metadata which are retrieved from APIC)
- ./templates (contains the Go templates used to generate the full provider code)
	- provider.go.tmpl (the template used to generate the provider.go file in the ../internal/provider directory)
	- index.md.tmpl (the template used to generate the index (provider) documentation file in the ../docs directory)
	- resource.go.tmpl (the template used to generate the resource_*.go files in the ../internal/provider directory)
	- resource.md.tmpl (the template used to generate the *.md files in the ../docs/resources directory)
	- data_source.go.tmpl (the template used to generate the data_source_*.go files in the ../internal/provider directory)
	- data_source.md.tmpl (the template used to generate the *.md files in the ../docs/data-sources directory)
	- *_test.go.tmpl (the templates used to generate the *_test.go files in the ../internal/provider directory)
	- *_example.go.tmpl (the templates used to generate the example files used in the documentation which is auto generated with tfplugindocs in the ../examples directory)
- ./testVars (contains the manually created YAML files with the test variables used in the *_test.go and example files)

Usage:
	go run generate.go
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/buxizhizhoum/inflection"
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Paths to the directories containing the files used for generating the provider code and outputting the generated code
const (
	definitionsPath         = "./gen/definitions"
	metaPath                = "./gen/meta"
	templatePath            = "./gen/templates"
	testVarsPath            = "./gen/testVars"
	providerExamplePath     = "./examples/provider/provider.tf"
	resourcesExamplesPath   = "./examples/resources"
	datasourcesExamplesPath = "./examples/data-sources"
	docsPath                = "./docs"
	resourcesDocsPath       = "./docs/resources"
	datasourcesDocsPath     = "./docs/data-sources"
	providerPath            = "./internal/provider/"
)

const providerName = "aci"

// Function map used during template rendering in order to call functions from the template
// The map contains a key which is the name of the function used in the template and a value which is the function itself
// The functions itself are defined in the current file
var templateFuncs = template.FuncMap{
	"snakeCase":                      inflection.Underscore,
	"validatorString":                ValidatorString,
	"listToString":                   ListToString,
	"overwriteProperty":              GetOverwriteAttributeName,
	"fromInterfacesToString":         FromInterfacesToString,
	"containsNoneAttributeValue":     ContainsNoneAttributeValue,
	"containsStringAttributeValue":   ContainsStringAttributeValue,
	"add":                            func(val1, val2 int) int { return val1 + val2 },
	"lookupTestValue":                LookupTestValue,
	"createParentDnValue":            CreateParentDnValue,
	"getResourceNameAsDescription":   GetResourceNameAsDescription,
	"capitalize":                     Capitalize,
	"removeAnnotationFromProperties": RemoveAnnotationFromProperties,
}

// Global variables used for unique resource name setting based on label from meta data
var labels = []string{"dns_provider", "filter_entry"}
var duplicateLabels = []string{}
var resourceNames = map[string]string{}

// Global variable used to determine if a class is defined in a parent resource
// During testing this is required to determine if test step expects non empty plan states
var nestedClasses = []string{}

func GetResourceNameAsDescription(s string) string {
	return cases.Title(language.English).String(strings.ReplaceAll(s, "_", " "))
}

func RemoveAnnotationFromProperties(properties map[string]Property) map[string]Property {

	cleanedProperties := make(map[string]Property)
	for propertyName, property := range properties {
		if propertyName != "annotation" {
			cleanedProperties[propertyName] = property
		}
	}
	return cleanedProperties
}

func Capitalize(s string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(s[:1]), s[1:])
}

func ValidatorString(stringList []string) string {
	sort.Strings(stringList)
	return fmt.Sprintf("\"%s\"", strings.Join(stringList, "\", \""))
}

func ListToString(stringList []string) string {
	sort.Strings(stringList)
	return fmt.Sprintf("%s", strings.Join(stringList, ", "))
}

// Creates a parent dn value for the resources and datasources in the example files
func CreateParentDnValue(className, caller string, definitions Definitions) string {
	resourceName := GetResourceName(className, definitions)
	return fmt.Sprintf("%s_%s.%s.id", providerName, resourceName, caller)
}

// Retrieves a value for a attribute of a aci class when defined in the testVars YAML file of the class
// Returns "foo" if no value is defined in the testVars YAML file
func LookupTestValue(classPkgName, propertyName string, testVars map[string]interface{}, definitions Definitions) string {
	lookupValue := "foo"
	propertyName = GetOverwriteAttributeName(classPkgName, propertyName, definitions)
	_, ok := testVars["all"]
	if ok {
		val, ok := testVars["all"].(interface{}).(map[interface{}]interface{})[propertyName]
		if ok {
			lookupValue = val.(string)
		}
	}
	_, ok = testVars["resource_required"]
	if ok {
		val, ok := testVars["resource_required"].(interface{}).(map[interface{}]interface{})[propertyName]
		if ok {
			lookupValue = val.(string)
		}
	}
	return lookupValue
}

// Determins if a list of values contains the value "none"
func ContainsNoneAttributeValue(values []string) bool {
	if slices.Contains(values, "none") {
		return true
	}
	return false
}

func ContainsStringAttributeValue(s string, values []interface{}) bool {
	for _, value := range values {
		if s == value.(string) {
			return true
		}
	}
	return false
}

// Create a string from a list of interfaces
// Error handling is not implemented because the identifiedBy interfaces are known to be strings
// TODO rewrite to generic function with error handling
func FromInterfacesToString(identifiedBy []interface{}) string {
	var identifiers []string
	for _, identifier := range identifiedBy {
		identifiers = append(identifiers, identifier.(string))
	}
	return fmt.Sprintf("\"%s\"", strings.Join(identifiers, "\", \"")) // TODO similar return used in the templateFunc validatorString could be combined
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
		classModel.setClassModel(metaPath, false, definitions, []string{}, pkgNames)
		classModels[pkgName] = classModel
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
	return testVarsMap, nil
}

// Retrieves the property and classs overwrite definitions from the definitions YAML files
func getDefinitions() Definitions {

	definitions := Definitions{}

	files, err := os.ReadDir(definitionsPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if path.Ext(file.Name()) != ".yaml" {
			continue
		}
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
	}

	return definitions
}

// Remove all files in a directory except when the file contains the string "provider_test.go" in the name
func cleanDirectory(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		if !slices.Contains([]string{"provider_test.go", "test_constants.go"}, name) {
			err = os.RemoveAll(filepath.Join(dir, name))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Container function to clean all directories properly
func cleanDirectories() {
	cleanDirectory(providerPath)
	cleanDirectory(resourcesDocsPath)
	cleanDirectory(datasourcesDocsPath)

	// The *ExamplesPath directories are removed and recreated to ensure all previously rendered files are removed
	// The provider example file is not removed because it contains static provider configuration
	os.RemoveAll(resourcesExamplesPath)
	os.Mkdir(resourcesExamplesPath, 0755)
	os.RemoveAll(datasourcesExamplesPath)
	os.Mkdir(datasourcesExamplesPath, 0755)
}

func getExampleCode(filePath string) []byte {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return content
}

func main() {

	cleanDirectories()
	definitions := getDefinitions()
	classModels := getClassModels(definitions)

	renderTemplate("provider.go.tmpl", "provider.go", providerPath, classModels)
	renderTemplate("index.md.tmpl", "index.md", docsPath, ProviderModel{Example: string(getExampleCode(providerExamplePath))})
	for _, model := range classModels {

		// Only render resources and datasources when the class has a unique identifier or is marked as include in the classes definitions YAML file
		if len(model.IdentifiedBy) > 0 || model.Include {
			model.ResourceName = GetResourceName(model.PkgName, definitions)
			setDocumentationData(&model, definitions)

			// Check if the testVars file exists, if not create it with generic values that might need to manually adjusted
			// When manual adjustment is needed the code needs to be regenerated
			_, err := os.Stat(fmt.Sprintf("%s/%s.yaml", testVarsPath, model.PkgName))
			if err != nil && os.IsNotExist(err) {
				renderTemplate("testvars.yaml.tmpl", fmt.Sprintf("%s.yaml", model.PkgName), testVarsPath, model)
			}

			testVarsMap, err := getTestVars(model)
			if err != nil {
				panic(err)
			}
			model.TestVars = testVarsMap

			childMap := make(map[string]Model, 0)
			for childName, childModel := range model.Children {
				childModel.ResourceName = GetResourceName(childModel.PkgName, definitions)
				testVarsMap, err := getTestVars(childModel)
				if err != nil {
					panic(err)
				}
				childModel.TestVars = testVarsMap
				childMap[childName] = childModel
			}
			model.Children = childMap

			renderTemplate("resource.go.tmpl", fmt.Sprintf("resource_%s_%s.go", providerName, model.ResourceName), providerPath, model)
			renderTemplate("datasource.go.tmpl", fmt.Sprintf("data_source_%s_%s.go", providerName, model.ResourceName), providerPath, model)

			os.Mkdir(fmt.Sprintf("%s/%s_%s", resourcesExamplesPath, providerName, model.ResourceName), 0755)
			renderTemplate("resource_example.tf.tmpl", fmt.Sprintf("%s_%s/resource.tf", providerName, model.ResourceName), resourcesExamplesPath, model)
			model.Example = string(hclwrite.Format(getExampleCode(fmt.Sprintf("%s/%s_%s/resource.tf", resourcesExamplesPath, providerName, model.ResourceName))))
			renderTemplate("resource.md.tmpl", fmt.Sprintf("%s.md", model.ResourceName), resourcesDocsPath, model)

			os.Mkdir(fmt.Sprintf("%s/%s_%s", datasourcesExamplesPath, providerName, model.ResourceName), 0755)
			renderTemplate("datasource_example.tf.tmpl", fmt.Sprintf("%s_%s/data-source.tf", providerName, model.ResourceName), datasourcesExamplesPath, model)
			model.Example = string(hclwrite.Format(getExampleCode(fmt.Sprintf("%s/%s_%s/data-source.tf", datasourcesExamplesPath, providerName, model.ResourceName))))
			renderTemplate("datasource.md.tmpl", fmt.Sprintf("%s.md", model.ResourceName), datasourcesDocsPath, model)

			if model.TestVars != nil {
				renderTemplate("resource_test.go.tmpl", fmt.Sprintf("resource_%s_%s_test.go", providerName, model.ResourceName), providerPath, model)
				renderTemplate("datasource_test.go.tmpl", fmt.Sprintf("data_source_%s_%s_test.go", providerName, model.ResourceName), providerPath, model)
			}
		}
	}

}

// A Model that represents the provider
type ProviderModel struct {
	Example string
}

// A Model represents a ACI class
// All information is retrieved directly or deduced from the metadata
type Model struct {
	PkgName                  string
	Label                    string
	Name                     string
	RnFormat                 string
	Comment                  string
	ResourceClassName        string
	ResourceName             string
	Example                  string
	SubCategory              string
	UiLocation               string
	RelationshipClass        string
	RelationshipResourceName string
	ChildClasses             []string
	ContainedBy              []string
	DocumentationDnFormats   []string
	DocumentationParentDns   []string
	DocumentationExamples    []string
	IdentifiedBy             []interface{}
	DnFormats                []interface{}
	Properties               map[string]Property
	Children                 map[string]Model
	Parents                  []string
	Configuration            map[string]interface{}
	TestVars                 map[string]interface{}
	Definitions              Definitions
	// Below booleans are used during template rendering to determine correct rendering the go code
	AllowDelete               bool
	AllowChildDelete          bool
	HasBitmask                bool
	HasChild                  bool
	HasParent                 bool
	HasAnnotation             bool
	HasValidValues            bool
	HasChildWithoutIdentifier bool
	HasNaming                 bool
	Include                   bool
}

// A Property represents a ACI class property
// All information is retrieved directly or deduced from the metadata
type Property struct {
	Name               string
	PropertyName       string
	SnakeCaseName      string
	ResourceClassName  string
	PkgName            string
	ValueType          string
	Label              string
	Comment            string
	DefaultValue       string
	ValidValues        []string
	IdentifiedBy       []interface{}
	IdentifyProperties []Property
	Validators         []interface{}
	// Below booleans are used during template rendering to determine correct rendering the go code
	IsNaming   bool
	CreateOnly bool
}

// A Definitions represents the ACI class and property definitions as defined in the definitions YAML files
type Definitions struct {
	Classes    map[string]interface{}
	Properties map[string]interface{}
}

// Reads the class details from the meta file and sets all details to the Model
func (m *Model) setClassModel(metaPath string, child bool, definitions Definitions, parents, pkgNames []string) {
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
		m.SetClassLabel(classDetails, child)
		m.SetClassName(classDetails)
		m.SetClassRnFormat(classDetails)
		m.SetClassDnFormats(classDetails)
		m.SetClassIdentifiers(classDetails)
		m.SetClassInclude()
		m.SetClassAllowDelete(classDetails)
		m.SetClassContainedByAndParent(classDetails, parents, pkgNames)
		m.SetClassComment(classDetails)
		m.SetClassProperties(classDetails)
		m.SetClassChildren(classDetails, pkgNames)
	}

	/*
		Checks if the setClassModel is a child class to prevent more than one level of nesting
			- Correct: Parent -> Child
			- Incorrect: Parent -> Child -> Grandchild
		// TODO add grandchild logic
	*/
	if !child {
		if len(m.ChildClasses) > 0 {
			m.HasChild = true
			m.Children = make(map[string]Model)
			for _, child := range m.ChildClasses {
				childModel := Model{PkgName: child}
				childModel.setClassModel(metaPath, true, definitions, []string{m.PkgName}, pkgNames)
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
			}
		} else {
			m.HasChild = false
		}
	} else {
		if !slices.Contains(nestedClasses, m.PkgName) {
			nestedClasses = append(nestedClasses, m.PkgName)
		}
	}
}

func (m *Model) SetClassLabel(classDetails interface{}, child bool) {
	m.Label = cleanLabel(classDetails.(map[string]interface{})["label"].(string))

	if !child {
		if slices.Contains(labels, m.Label) || m.Label == "" {
			if !slices.Contains(duplicateLabels, m.Label) {
				duplicateLabels = append(duplicateLabels, m.Label)
			}
			resourceNames[m.PkgName] = inflection.Underscore(m.PkgName)
		} else {
			labels = append(labels, m.Label)
			resourceNames[m.PkgName] = m.Label
		}
	}

}

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
	return inflection.Underscore(returnLabel)
}

func (m *Model) SetClassName(classDetails interface{}) {
	m.Name = classDetails.(map[string]interface{})["className"].(string)
	m.ResourceName = GetResourceName(m.PkgName, m.Definitions)
}

func (m *Model) SetClassRnFormat(classDetails interface{}) {
	m.RnFormat = GetOverwriteRnPrepend(classDetails.(map[string]interface{})["rnFormat"].(string), m.PkgName, m.Definitions)
	if strings.HasPrefix(m.RnFormat, "rs") {
		toMo := classDetails.(map[string]interface{})["relationInfo"].(map[string]interface{})["toMo"].(string)
		m.RelationshipClass = strings.Replace(toMo, ":", "", 1)
		m.RelationshipResourceName = GetResourceName(m.RelationshipClass, m.Definitions)
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
}

func (m *Model) SetClassChildren(classDetails interface{}, pkgNames []string) {
	childClasses := []string{}
	notAddedFromRnMap := []string{}
	rnMap := classDetails.(map[string]interface{})["rnMap"].(map[string]interface{})
	for rn, className := range rnMap {
		// TODO check if this condition is correct since there might be cases where that we should exclude
		if !strings.HasSuffix(rn, "-") || strings.HasPrefix(rn, "rs") {
			pkgName := strings.ReplaceAll(className.(string), ":", "")
			if slices.Contains(pkgNames, pkgName) {
				childClasses = append(childClasses, pkgName)
			}
		} else {
			notAddedFromRnMap = append(notAddedFromRnMap, strings.ReplaceAll(className.(string), ":", ""))
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

func (m *Model) SetClassInclude() {
	if classDetails, ok := m.Definitions.Classes[m.PkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "include" {
				m.Include = value.(bool)
			} else {
				m.Include = false
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

func (m *Model) SetClassContainedByAndParent(classDetails interface{}, parents, pkgNames []string) {

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

// Construct a property map for the class, that contains all details of the property that will be used during the rendering of the template
func (m *Model) SetClassProperties(classDetails interface{}) {

	properties := make(map[string]Property)

	for propertyName, propertyValue := range classDetails.(map[string]interface{})["properties"].(map[string]interface{}) {

		if propertyValue.(map[string]interface{})["isConfigurable"] == true {

			if ignoreProperty(propertyName, m.PkgName, m.Definitions) {
				continue
			}

			property := Property{
				Name:              fmt.Sprintf("%s%s", strings.ToUpper(propertyName[0:1]), propertyName[1:]),
				PropertyName:      propertyName,
				SnakeCaseName:     inflection.Underscore(propertyName),
				ResourceClassName: strings.ToUpper(m.PkgName[:1]) + m.PkgName[1:],
				PkgName:           m.PkgName,
				IdentifiedBy:      m.IdentifiedBy,
				ValueType:         propertyValue.(map[string]interface{})["uitype"].(string),
				Label:             propertyValue.(map[string]interface{})["label"].(string),
				IsNaming:          propertyValue.(map[string]interface{})["isNaming"].(bool),
				CreateOnly:        propertyValue.(map[string]interface{})["createOnly"].(bool),
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
				for _, details := range propertyValue.(map[string]interface{})["validValues"].([]interface{}) {
					validValue := details.(map[string]interface{})["localName"].(string)
					if validValue != "defaultValue" {
						property.ValidValues = append(property.ValidValues, validValue)
					}
				}
				if len(property.ValidValues) > 0 {
					m.HasValidValues = true
				}
			}

			val, ok = propertyValue.(map[string]interface{})["default"]
			if ok {
				if reflect.TypeOf(val).String() == "string" {
					property.DefaultValue = val.(string)
				} else if reflect.TypeOf(val).String() == "float64" {
					property.DefaultValue = fmt.Sprintf("%f", val.(float64))
				} else {
					log.Fatal(fmt.Sprintf("Reflect type %s not not defined. Define in SetClassProperties function.", reflect.TypeOf(val).String()))
				}
			}

			properties[propertyName] = property

		}

	}

	m.Properties = properties
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
								return fmt.Sprintf(v.(string), GetResourceNameAsDescription(GetResourceName(classPkgName, definitions)))
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

/*
Determine if a attribute name in terraform configuration should be overwritten by a attribute name overwrite in the properties.yaml file
Precendence order is:
 1. class level from properties.yaml
 2. global level from properties.yaml
 3. meta data property attribute name
*/
func GetOverwriteAttributeName(classPkgName, propertyName string, definitions Definitions) string {
	precedenceList := []string{classPkgName, "global"}
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
func GetOverwriteRnPrepend(rnFormat, classPkgName string, definitions Definitions) string {

	if v, ok := definitions.Classes[classPkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "rn_prepend" {
				rnFormat = fmt.Sprintf("%s/%s", value.(string), rnFormat)
			}
		}
	}
	return rnFormat
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

// Set variables that are used during the rendering of the example and documentation templates
func setDocumentationData(m *Model, definitions Definitions) {
	UiLocation := "Generic"
	SubCategory := "Generic"
	if v, ok := definitions.Classes[m.PkgName]; ok {
		for key, value := range v.(map[interface{}]interface{}) {
			if key.(string) == "ui_location" {
				UiLocation = value.(string)
			} else if key.(string) == "sub_category" {
				SubCategory = value.(string)
			}
		}
	}
	m.UiLocation = UiLocation
	m.SubCategory = SubCategory

	resourcesFound := [][]string{}
	resourcesNotFound := []string{}
	for _, containedClassName := range m.ContainedBy {
		resourceName := GetResourceName(containedClassName, definitions)
		if resourceName != "" {
			resourcesFound = append(resourcesFound, []string{resourceName, containedClassName})
		} else {
			resourcesNotFound = append(resourcesNotFound, containedClassName)
		}
	}

	docsParentDnAmount := m.Configuration["docs_parent_dn_amount"].(int)

	if len(m.DocumentationDnFormats) > docsParentDnAmount {
		m.DocumentationDnFormats = m.DocumentationDnFormats[0:docsParentDnAmount]
		m.DocumentationDnFormats = append(m.DocumentationDnFormats, "Too many DN formats to display, see model documentation for all possible parents.")
	}

	if len(resourcesFound) > docsParentDnAmount {
		for _, resourceDetails := range resourcesFound[0:docsParentDnAmount] {
			m.DocumentationParentDns = append(m.DocumentationParentDns, fmt.Sprintf("`%s_%s` (class: %s)", providerName, resourceDetails[0], resourceDetails[1]))
		}
		m.DocumentationParentDns = append(m.DocumentationParentDns, "Too many parent DNs to display, see model documentation for all possible parents.")
	} else {
		for _, resourceDetails := range resourcesFound {
			m.DocumentationParentDns = append(m.DocumentationParentDns, fmt.Sprintf("`%s_%s` (class: %s)", providerName, resourceDetails[0], resourceDetails[1]))
		}
	}

	if len(resourcesNotFound) != 0 {
		if len(resourcesNotFound) > docsParentDnAmount-len(resourcesFound) {
			// TODO catch default classes and add to documentation
			resourcesNotFound = resourcesNotFound[0:(docsParentDnAmount - len(resourcesFound))]
			m.DocumentationParentDns = append(m.DocumentationParentDns, "Too many classes to display, see model documentation for all possible classes.")
		} else {
			resourceDetails := strings.Join(resourcesNotFound, "`\n    - `")
			m.DocumentationParentDns = append(m.DocumentationParentDns, fmt.Sprintf("The distinquised name (DN) of classes below can be used but currently there is no available resource for it:\n    - `%s` ", resourceDetails))
		}
	}

	// TODO add overwrite to provide which documentation examples to be included
	docsExampleAmount := m.Configuration["docs_examples_amount"].(int)

	if len(m.ContainedBy) >= docsExampleAmount {
		for _, resourceDetails := range resourcesFound[0 : docsExampleAmount-1] {
			m.DocumentationExamples = append(m.DocumentationExamples, resourceDetails[1])
		}
	} else {
		for _, resourceDetails := range resourcesFound {
			m.DocumentationExamples = append(m.DocumentationExamples, resourceDetails[1])
		}
	}
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
