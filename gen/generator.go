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
	"strings"
	"text/template"
	"unicode"

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
	testVarsPath            = "./gen/testvars"
	providerExamplePath     = "./examples/provider/provider.tf"
	resourcesExamplesPath   = "./examples/resources"
	datasourcesExamplesPath = "./examples/data-sources"
	docsPath                = "./docs"
	resourcesDocsPath       = "./docs/resources"
	datasourcesDocsPath     = "./docs/data-sources"
	providerPath            = "./internal/provider/"
)

const providerName = "aci"
const pubhupDevnetBaseUrl = "https://pubhub.devnetcloud.com/media/model-doc-latest/docs"

// Function map used during template rendering in order to call functions from the template
// The map contains a key which is the name of the function used in the template and a value which is the function itself
// The functions itself are defined in the current file
var templateFuncs = template.FuncMap{
	"snakeCase":                    Underscore,
	"validatorString":              ValidatorString,
	"containsString":               ContainsString,
	"listToString":                 ListToString,
	"overwriteProperty":            GetOverwriteAttributeName,
	"overwritePropertyValue":       GetOverwriteAttributeValue,
	"createTestValue":              func(val string) string { return fmt.Sprintf("test_%s", val) },
	"createNonExistingValue":       func(val string) string { return fmt.Sprintf("non_existing_%s", val) },
	"getParentTestDependencies":    GetParentTestDependencies,
	"getTargetTestDependencies":    GetTargetTestDependencies,
	"fromInterfacesToString":       FromInterfacesToString,
	"containsNoneAttributeValue":   ContainsNoneAttributeValue,
	"definedInMap":                 DefinedInMap,
	"add":                          func(val1, val2 int) int { return val1 + val2 },
	"lookupTestValue":              LookupTestValue,
	"lookupChildTestValue":         LookupChildTestValue,
	"createParentDnValue":          CreateParentDnValue,
	"getResourceName":              GetResourceName,
	"getResourceNameAsDescription": GetResourceNameAsDescription,
	"capitalize":                   Capitalize,
	"getTestConfigVariableName":    GetTestConfigVariableName,
	"getDevnetDocForClass":         GetDevnetDocForClass,
}

// Global variables used for unique resource name setting based on label from meta data
var labels = []string{"dns_provider", "filter_entry"}
var duplicateLabels = []string{}
var resourceNames = map[string]string{}
var rnPrefix = map[string]string{}
var targetRelationalPropertyClasses = map[string]string{}
var alwaysIncludeChildren = []string{"tag:Annotation", "tag:Tag"}
var excludeChildResourceNamesFromDocs = []string{"", "annotation", "tag"}

func GetResourceNameAsDescription(s string, definitions Definitions) string {
	resourceName := cases.Title(language.English).String(strings.ReplaceAll(s, "_", " "))
	for k, v := range definitions.Properties["global"].(map[interface{}]interface{})["resource_name_doc_overwrite"].(map[interface{}]interface{}) {
		resourceName = strings.ReplaceAll(resourceName, k.(string), v.(string))
	}
	return resourceName
}

func GetDevnetDocForClass(className string) string {
	return fmt.Sprintf("[%s](%s/app/index.html#/objects/%s/overview)", className, pubhupDevnetBaseUrl, className)
}

func Capitalize(s string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(s[:1]), s[1:])
}

func ContainsString(s, sub string) bool {
	if strings.Contains(s, sub) {
		return true
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

func ListToString(stringList []string) string {
	sort.Strings(stringList)
	return fmt.Sprintf("%s", strings.Join(stringList, ","))
}

// Creates a parent dn value for the resources and datasources in the example files
func CreateParentDnValue(className, caller string, definitions Definitions) string {
	resourceName := GetResourceName(className, definitions)
	return fmt.Sprintf("%s_%s.%s.id", providerName, resourceName, caller)
}

// Retrieves a value for a attribute of a aci class when defined in the testVars YAML file of the class
// Returns "test_value" if no value is defined in the testVars YAML file
func LookupTestValue(classPkgName, propertyName string, testVars map[string]interface{}, definitions Definitions) string {
	lookupValue := "test_value"
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

// Retrieves a value for a attribute of a aci class when defined in the testVars YAML file of the class
// Returns "test_value_for_child" if no value is defined in the testVars YAML file
func LookupChildTestValue(classPkgName, childResourceName, propertyName string, testVars map[string]interface{}, testValueIndex int, definitions Definitions) string {
	propertyName = GetOverwriteAttributeName(classPkgName, propertyName, definitions)
	overwritePropertyValue := GetOverwriteAttributeValue(classPkgName, propertyName, "", "test_values_for_parent", testValueIndex, definitions)
	if overwritePropertyValue != "" {
		return overwritePropertyValue
	}
	_, ok := testVars["children"]
	if ok {
		val, ok := testVars["children"].(map[interface{}]interface{})[childResourceName].([]interface{})[0].(map[interface{}]interface{})[propertyName]
		if ok {
			return val.(string)
		}
	} else {
		return fmt.Sprintf("%s_%d", propertyName, testValueIndex)
	}
	return "test_value_for_child"
}

func ContainsNoneAttributeValue(values []string) bool {
	if slices.Contains(values, "none") {
		return true
	}
	return false
}

func DefinedInMap(s string, values map[interface{}]interface{}) bool {
	if _, ok := values[s]; ok {
		return true
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
	testVarsMap["targetResourceName"] = model.TargetResourceName
	testVarsMap["targetResourceClassName"] = model.TargetResourceClassName
	testVarsMap["targetResourceParentClassName"] = ""
	testVarsMap["targetDn"] = model.TargetDn
	targets, targetsOk := testVarsMap["targets"]
	if targetsOk && targets != nil {
		if targets.([]interface{})[0].(map[interface{}]interface{})["parent_dependency"] != nil {
			testVarsMap["targetResourceParentClassName"] = targets.([]interface{})[0].(map[interface{}]interface{})["parent_dependency"].(string)
		}

		if targets.([]interface{})[0].(map[interface{}]interface{})["target_dn"] != "" {
			testVarsMap["targetDn"] = targets.([]interface{})[0].(map[interface{}]interface{})["target_dn"].(string)
		}
	}
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

// Remove all files in a directory except when the files that do not match the ignore list
func cleanDirectory(dir string, ignores []string) {
	d, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	for _, name := range names {
		if !slices.Contains(ignores, name) {
			err = os.RemoveAll(filepath.Join(dir, name))
			if err != nil {
				panic(err)
			}
		}
	}
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
		d, err := os.Open(dirPath)
		if err != nil {
			panic(err)
		}
		defer d.Close()
		names, err := d.Readdirnames(-1)
		if err != nil {
			panic(err)
		}
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
	cleanDirectory(docsPath, []string{"resources", "data-sources"})
	cleanDirectory(providerPath, []string{"provider_test.go", "utils.go", "test_constants.go", "resource_aci_rest_managed.go", "resource_aci_rest_managed_test.go", "data_source_aci_rest_managed.go", "data_source_aci_rest_managed_test.go", "annotation_unsupported.go"})
	cleanDirectory(resourcesDocsPath, []string{})
	cleanDirectory(datasourcesDocsPath, []string{})
	cleanDirectory(testVarsPath, []string{})

	// The *ExamplesPath directories are removed and recreated to ensure all previously rendered files are removed
	// The provider example file is not removed because it contains static provider configuration
	os.RemoveAll(resourcesExamplesPath)
	os.Mkdir(resourcesExamplesPath, 0755)
	os.RemoveAll(datasourcesExamplesPath)
	os.Mkdir(datasourcesExamplesPath, 0755)

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

// When GEN_CLASSES environment variable is set, the class metadata is retrieved from the APIC or devnet docs and stored in the meta directory.
func getClassMetadata() {

	classNames := os.Getenv("GEN_CLASSES")

	if classNames != "" {
		var name, nameSpace, url string
		classNameList := strings.Split(classNames, ",")
		for _, className := range classNameList {

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
func genAnnotationUnsupported() []string {
	classes := []string{}
	_, gen_annotation_unsupported := os.LookupEnv("GEN_ANNOTATION_UNSUPPORTED")
	if gen_annotation_unsupported {
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
	getClassMetadata()
	cleanDirectories()

	definitions := getDefinitions()
	classModels := getClassModels(definitions)
	annotationUnsupported := genAnnotationUnsupported()

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

		// Only render resources and datasources when the class has a unique identifier or is marked as include in the classes definitions YAML file
		if len(model.IdentifiedBy) > 0 || model.Include {

			// All classmodels have been read, thus now the model, child and relational resources names can be set
			// When done before additional files would need to be opened and read which would slow down the generation process
			model.ResourceName = GetResourceName(model.PkgName, definitions)
			model.RelationshipResourceName = GetResourceName(model.RelationshipClass, definitions)
			childMap := make(map[string]Model, 0)
			for childName, childModel := range model.Children {
				childModel.ChildResourceName = GetResourceName(childModel.PkgName, definitions)
				if len(childModel.IdentifiedBy) > 0 {
					// TODO add logic to determine the naming for plural child resources
					childModel.ResourceNameDocReference = childModel.ChildResourceName
					childModel.ResourceName = fmt.Sprintf("%ss", childModel.ChildResourceName)
				} else {
					childModel.ResourceName = childModel.ChildResourceName
				}
				childModel.RelationshipResourceName = GetResourceName(childModel.RelationshipClass, definitions)
				childMap[childName] = childModel
			}
			model.Children = childMap

			// Set the documentation specific information for the resource
			// This is done to ensure references can be made to parent/child resources and output amounts can be restricted
			setDocumentationData(&model, definitions)

			// Render the testvars file for the resource
			// First generate run would not mean the file is correct from beginning since some testvars would need to be manually overwritten in the properties definitions YAML file
			SetModelTargetValues(model.PkgName, &model, classModels, definitions)
			renderTemplate("testvars.yaml.tmpl", fmt.Sprintf("%s.yaml", model.PkgName), testVarsPath, model)
			testVarsMap, err := getTestVars(model)
			if err != nil {
				panic(err)
			}
			model.TestVars = testVarsMap
			renderTemplate("resource.go.tmpl", fmt.Sprintf("resource_%s_%s.go", providerName, model.ResourceName), providerPath, model)
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

			os.Mkdir(fmt.Sprintf("%s/%s_%s", datasourcesExamplesPath, providerName, model.ResourceName), 0755)
			renderTemplate("provider_example.tf.tmpl", fmt.Sprintf("%s_%s/provider.tf", providerName, model.ResourceName), datasourcesExamplesPath, model)
			renderTemplate("datasource_example.tf.tmpl", fmt.Sprintf("%s_%s/data-source.tf", providerName, model.ResourceName), datasourcesExamplesPath, model)
			// Leverage the hclwrite package to format the example code
			model.ExampleDataSource = string(hclwrite.Format(getExampleCode(fmt.Sprintf("%s/%s_%s/data-source.tf", datasourcesExamplesPath, providerName, model.ResourceName))))
			renderTemplate("datasource.md.tmpl", fmt.Sprintf("%s.md", model.ResourceName), datasourcesDocsPath, model)
			renderTemplate("resource_test.go.tmpl", fmt.Sprintf("resource_%s_%s_test.go", providerName, model.ResourceName), providerPath, model)
			renderTemplate("datasource_test.go.tmpl", fmt.Sprintf("data_source_%s_%s_test.go", providerName, model.ResourceName), providerPath, model)

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
	PkgName                   string
	Label                     string
	Name                      string
	RnFormat                  string
	RnPrepend                 string
	Comment                   string
	ResourceClassName         string
	ResourceName              string
	ResourceNameDocReference  string
	ChildResourceName         string
	ExampleDataSource         string
	ExampleResource           string
	ExampleResourceFull       string
	SubCategory               string
	RelationshipClass         string
	RelationshipResourceName  string
	Versions                  string
	ChildClasses              []string
	ContainedBy               []string
	Contains                  []string
	DocumentationDnFormats    []string
	DocumentationParentDns    []string
	DocumentationExamples     []string
	TargetResourceClassName   string
	TargetResourceName        string
	TargetDn                  string
	TargetProperties          map[string]Property
	TargetNamedProperties     map[string]Property
	DocumentationChildren     []string
	ResourceNotes             []string
	ResourceWarnings          []string
	DatasourceNotes           []string
	DatasourceWarnings        []string
	Parents                   []string
	UiLocations               []string
	IdentifiedBy              []interface{}
	DnFormats                 []interface{}
	Properties                map[string]Property
	NamedProperties           map[string]Property
	Children                  map[string]Model
	Configuration             map[string]interface{}
	TestVars                  map[string]interface{}
	Definitions               Definitions
	ResourceNameAsDescription string
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
	HasOptionalProperties     bool
	HasOnlyRequiredProperties bool
	HasNamedProperties        bool
	HasChildNamedProperties   bool
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
	Versions           string
	NamedPropertyClass string
	ValidValues        []string
	IdentifiedBy       []interface{}
	Validators         []interface{}
	IdentifyProperties []Property
	// Below booleans are used during template rendering to determine correct rendering the go code
	IsNaming   bool
	CreateOnly bool
	IsRequired bool
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
		m.SetClassContainedByAndParent(classDetails, parents)
		m.SetClassContains(classDetails)
		m.SetClassComment(classDetails)
		m.SetClassVersions(classDetails)
		m.SetClassProperties(classDetails)
		m.SetClassChildren(classDetails, pkgNames)
		m.SetResourceNotesAndWarnigns(m.PkgName, definitions)
		m.SetResourceNameAsDescription(m.PkgName, definitions)
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
				if childModel.HasNamedProperties {
					m.HasNamedProperties = true
					m.HasChildNamedProperties = true
				}
			}
		} else {
			m.HasChild = false
		}
	}
}

func (m *Model) SetClassLabel(classDetails interface{}, child bool) {
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
	if strings.HasPrefix(m.RnFormat, "rs") {
		toMo := classDetails.(map[string]interface{})["relationInfo"].(map[string]interface{})["toMo"].(string)
		m.RelationshipClass = strings.Replace(toMo, ":", "", 1)
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
	rnMap := classDetails.(map[string]interface{})["rnMap"].(map[string]interface{})
	for rn, className := range rnMap {
		// TODO check if this condition is correct since there might be cases where that we should exclude
		if !strings.HasSuffix(rn, "-") || strings.HasPrefix(rn, "rs") || slices.Contains(alwaysIncludeChildren, className.(string)) {
			pkgName := strings.ReplaceAll(className.(string), ":", "")
			if slices.Contains(pkgNames, pkgName) {
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
}

func (m *Model) SetResourceNameAsDescription(classPkgName string, definitions Definitions) {
	m.ResourceNameAsDescription = GetResourceNameAsDescription(GetResourceName(classPkgName, definitions), definitions)
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
		m.Versions = formatVersion(versions.(string))
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

	for propertyName, propertyValue := range classDetails.(map[string]interface{})["properties"].(map[string]interface{}) {

		if propertyValue.(map[string]interface{})["isConfigurable"] == true {

			if ignoreProperty(propertyName, m.PkgName, m.Definitions) {
				continue
			}

			property := Property{
				Name:              fmt.Sprintf("%s%s", strings.ToUpper(propertyName[0:1]), propertyName[1:]),
				PropertyName:      propertyName,
				SnakeCaseName:     Underscore(propertyName),
				ResourceClassName: strings.ToUpper(m.PkgName[:1]) + m.PkgName[1:],
				PkgName:           m.PkgName,
				IdentifiedBy:      m.IdentifiedBy,
				ValueType:         propertyValue.(map[string]interface{})["uitype"].(string),
				Label:             propertyValue.(map[string]interface{})["label"].(string),
				IsNaming:          propertyValue.(map[string]interface{})["isNaming"].(bool),
				CreateOnly:        propertyValue.(map[string]interface{})["createOnly"].(bool),
			}

			if requiredProperty(propertyName, m.PkgName, m.Definitions) || property.IsNaming == true {
				property.IsRequired = true
				requiredCount += 1
			}

			if property.IsRequired == false {
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

			versions, ok := classDetails.(map[string]interface{})["versions"]
			if ok {
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

			properties[propertyName] = property

		}

	}

	m.Properties = properties
	m.NamedProperties = namedProperties
	if requiredCount == len(properties) {
		m.HasOnlyRequiredProperties = true
	}
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
func GetOverwriteAttributeValue(classPkgName, propertyName, propertyValue, testType string, valueIndex int, definitions Definitions) string {
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
							return v.(string)
						}
					}

				}
			}
		}
	}
	return propertyValue
}

func GetParentTestDependencies(classPkgName string, index int, definitions Definitions) map[string]interface{} {
	parentDependency := ""
	classInParent := false
	parentDependencyName := ""
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
				}
			}
		}
	}
	return map[string]interface{}{"parent_dependency": parentDependency, "class_in_parent": classInParent, "parent_dependency_name": parentDependencyName}
}

func GetTargetTestDependencies(classPkgName string, index int, definitions Definitions) map[string]interface{} {
	parentDependency := ""
	grandParentDependency := ""
	parentDependencyDn := ""
	targetDn := ""
	parentClassName := ""
	if classDetails, ok := definitions.Properties[classPkgName]; ok {
		for key, value := range classDetails.(map[interface{}]interface{}) {
			if key.(string) == "targets" {
				if len(value.([]interface{})) > index {
					targetMap := value.([]interface{})[index].(map[interface{}]interface{})

					if pd, ok := targetMap["parent_dependency"]; ok {
						parentDependency = pd.(string)
					}

					if gpd, ok := targetMap["grandparent_class_name"]; ok {
						grandParentDependency = gpd.(string)
					}

					if pd_dn, ok := targetMap["parent_dependency_dn"]; ok {
						parentDependencyDn = pd_dn.(string)
					}

					if doc_dn, ok := targetMap["target_dn"]; ok {
						targetDn = doc_dn.(string)
					}

					if pcn, ok := targetMap["class_name"]; ok {
						parentClassName = pcn.(string)
					}
				}
			}
		}
	}

	return map[string]interface{}{
		"parent_dependency":      parentDependency,
		"parent_dependency_dn":   parentDependencyDn,
		"grandparent_class_name": grandParentDependency,
		"target_dn":              targetDn,
		"parent_class_name":      parentClassName,
	}
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

func SetModelTargetValues(PkgName string, model *Model, classModels map[string]Model, definitions Definitions) {
	targetTestDependencies := GetTargetTestDependencies(model.PkgName, 0, definitions)
	parentClassName := targetTestDependencies["parent_class_name"].(string)
	model.TargetProperties = classModels[parentClassName].Properties
	model.TargetNamedProperties = classModels[parentClassName].NamedProperties
	model.TargetResourceClassName = classModels[parentClassName].PkgName
	model.TargetDn = targetTestDependencies["target_dn"].(string)
	model.TargetResourceName = GetResourceName(model.TargetResourceClassName, definitions)
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
		if resourceName != "" {
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
			m.DocumentationParentDns = append(m.DocumentationParentDns, fmt.Sprintf("The distinquised name (DN) of classes below can be used but currently there is no available resource for it:\n%s", resourceDetails))
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
		if !match {
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
