package data

import (
	"fmt"
	"maps"
	"net/netip"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// DataSourceTestDependenciesHCL renders every resource needed by the selected
// data-source test flow. The returned configuration is complete within the
// generated test file and does not refer to legacy test constants.
func (c *Class) DataSourceTestDependenciesHCL(ds *DataStore) (string, error) {
	return newTestConfigurationRenderer(c, ds).dependenciesHCL()
}

// DataSourceTestExistingHCL renders the full resource configuration and the
// data source that reads it. TestValues.Create is the all-attributes scenario.
func (c *Class) DataSourceTestExistingHCL(ds *DataStore) (string, error) {
	return newTestConfigurationRenderer(c, ds).existingDataSourceHCL()
}

// DataSourceTestMissingHCL renders a lookup for an object that was not created.
// It is empty when the class does not have a safe negative lookup flow.
func (c *Class) DataSourceTestMissingHCL(ds *DataStore) (string, error) {
	return newTestConfigurationRenderer(c, ds).missingDataSourceHCL()
}

type testConfigurationRenderer struct {
	class      *Class
	dataStore  *DataStore
	dependency *TestDependency
}

type testResourceScenario int

const (
	fullTestResourceScenario testResourceScenario = iota
	dependencyTestResourceScenario
	alternateDependencyTestResourceScenario
)

func newTestConfigurationRenderer(class *Class, dataStore *DataStore) testConfigurationRenderer {
	return testConfigurationRenderer{class: class, dataStore: dataStore}
}

func (p *Property) canUseTestPrefix(class *Class) bool {
	// Automatic namespacing only changes a freely named RN component, never an IP, MAC, or enum.
	if p == nil || p.ReadOnly || !class.IsIdentifyingProperty(p.PropertyName) || len(p.ValidValues) > 0 {
		return false
	}
	uiType, _ := p.metaDetails["uitype"].(string)
	if uiType == "number" || p.ValueType != String {
		return false
	}
	return p.PropertyName == "name" || p.PrefixEnabled != nil && *p.PrefixEnabled
}

func (r testConfigurationRenderer) parentDependency(class *Class, dependency *TestDependency) *TestDependency {
	find := func(reference string) *TestDependency {
		if dependency != nil {
			if parent := findTestDependencyByObjectReference(dependency.Dependencies, reference); parent != nil {
				return parent
			}
		}
		return r.findDependency(reference)
	}
	if class != nil {
		if parent := class.ParentDnProperty(); parent != nil {
			if dependency != nil {
				if value, ok := testPropertyOverride(parent, dependency.ConfigOverrides); ok {
					return find(value)
				}
			}
			if parent.TestValues != nil {
				scenario := parent.TestValues.Create
				if dependency != nil {
					resourceScenario := dependencyTestResourceScenario
					if dependency.TestInstance > 1 {
						resourceScenario = alternateDependencyTestResourceScenario
					}
					scenario = testResourcePropertyScenario(class, parent, resourceScenario)
				}
				if entry, ok := scenario.firstIncludedEntry(); ok && entry.ValueType == ReferenceValue {
					return find(entry.ConfigValue)
				}
			}
		}
	}
	if dependency != nil {
		// Definition-only resources carry their parent inputs explicitly instead of normalized properties.
		for _, name := range slices.Sorted(maps.Keys(dependency.ConfigOverrides)) {
			if (name == "parent_dn" || strings.HasSuffix(name, "_dn")) && !strings.HasPrefix(name, "relation_") {
				if parent := find(dependency.ConfigOverrides[name]); parent != nil {
					return parent
				}
			}
		}
	}
	return nil
}

func (r testConfigurationRenderer) dependencyIsolated(dependency *TestDependency, visited map[*TestDependency]bool) bool {
	if dependency == nil || dependency.ReferenceType != ResourceReference || visited[dependency] {
		return false
	}
	visited[dependency] = true
	class, loaded := r.dataStore.Classes[dependency.Class.String()]
	if loaded {
		for _, property := range class.Properties {
			if property.canUseTestPrefix(&class) {
				// An explicit exclusion asserts that this object's existing name already provides isolation.
				return true
			}
		}
		return r.dependencyIsolated(r.parentDependency(&class, dependency), visited)
	}
	if dependency.ConfigOverrides["name"] != "" {
		return true
	}
	return r.dependencyIsolated(r.parentDependency(nil, dependency), visited)
}

func (r testConfigurationRenderer) shouldPrefixProperty(class *Class, property *Property) bool {
	if !property.canUseTestPrefix(class) {
		return false
	}
	if r.dependency != nil {
		if r.dependency.ReferenceType != ResourceReference {
			return false
		}
		if r.dependency.PrefixEnabled != nil {
			return *r.dependency.PrefixEnabled
		}
	}
	if property.PrefixEnabled != nil {
		return *property.PrefixEnabled
	}
	return !r.dependencyIsolated(r.parentDependency(class, r.dependency), make(map[*TestDependency]bool))
}

func (r testConfigurationRenderer) shouldPrefixLegacyName(dependency *TestDependency) bool {
	if dependency.ReferenceType != ResourceReference {
		return false
	}
	if dependency.PrefixEnabled != nil {
		return *dependency.PrefixEnabled
	}
	definition := r.dataStore.ClassDefinitions[dependency.Class.String()]
	if enabled := definition.Properties["name"].TestConfig.Prefix.Enabled; enabled != nil {
		return *enabled
	}
	return !r.dependencyIsolated(r.parentDependency(nil, dependency), make(map[*TestDependency]bool))
}

func prefixedTestEntry(entry TestValueEntry) TestValueEntry {
	if entry.ValueType == StringValue {
		entry.ConfigValue = fmt.Sprintf("format(\"%%s-%%s\", var.test_prefix, %s)", testValueHCL(entry))
		entry.ValueType = ExpressionValue
	}
	return entry
}

func (r testConfigurationRenderer) testPropertyScenario(class *Class, property *Property, scenario TestValueScenario) TestValueScenario {
	// Copy before projecting runtime expressions so public examples retain the canonical values.
	scenario = scenario.clone()
	prefix := r.shouldPrefixProperty(class, property)
	for index, entry := range scenario.Entries {
		if prefix {
			entry = prefixedTestEntry(entry)
		}
		scenario.Entries[index] = r.prefixedLiteralDN(entry)
	}
	return scenario
}

func (r testConfigurationRenderer) prefixedLiteralDN(entry TestValueEntry) TestValueEntry {
	if entry.ValueType != StringValue || !strings.HasPrefix(entry.ConfigValue, "uni/") {
		return entry
	}
	// Rewrite only exact DN prefixes belonging to named roots this fixture creates.
	// Fixed infrastructure DNs and external static targets remain untouched.
	visited := make(map[*TestDependency]bool)
	var rewrite func([]*TestDependency) bool
	rewrite = func(dependencies []*TestDependency) bool {
		for _, dependency := range dependencies {
			if dependency == nil || visited[dependency] {
				continue
			}
			visited[dependency] = true
			if projected, matched := r.prefixedDependencyDN(entry, dependency); matched {
				entry = projected
				return true
			}
			if rewrite(dependency.Dependencies) {
				return true
			}
		}
		return false
	}
	rewrite(r.selectedDependencies())
	return entry
}

func (r testConfigurationRenderer) prefixedDependencyDN(entry TestValueEntry, dependency *TestDependency) (TestValueEntry, bool) {
	class, loaded := r.dataStore.Classes[dependency.Class.String()]
	if !loaded || len(class.IdentifiedBy) != 1 {
		return entry, false
	}
	dependencyRenderer := r
	dependencyRenderer.dependency = dependency
	name := class.Properties["name"]
	if name == nil || name.TestValues == nil || !dependencyRenderer.shouldPrefixProperty(&class, name) {
		return entry, false
	}
	scenario := dependencyTestResourceScenario
	if dependency.TestInstance > 1 {
		scenario = alternateDependencyTestResourceScenario
	}
	value, included := testResourcePropertyScenario(&class, name, scenario).firstIncludedEntry()
	if override, ok := testPropertyOverride(name, dependency.ConfigOverrides); ok {
		value.ConfigValue = override
		value.ValueType = testOverrideValueType(override)
		included = true
	}
	parentDN := strings.TrimSuffix(class.FixedParentDn(), "/")
	root := fmt.Sprintf("%s/%s", parentDN, strings.ReplaceAll(class.RnFormat, "{name}", value.ConfigValue))
	if !included || value.ValueType != StringValue || entry.ConfigValue != root && !strings.HasPrefix(entry.ConfigValue, root+"/") {
		return entry, false
	}
	address, ok := parseTerraformTestReference(dependency.Reference)
	if !ok {
		return entry, false
	}
	format := fmt.Sprintf("%s/%s%s", parentDN, strings.ReplaceAll(class.RnFormat, "{name}", "%s"), strings.TrimPrefix(entry.ConfigValue, root))
	entry.ConfigValue = fmt.Sprintf("format(%q, %s)", format, countedTestDependencyReference(address.withAttribute("name")))
	entry.ValueType = ExpressionValue
	return entry, true
}

func (r testConfigurationRenderer) dependenciesHCL() (string, error) {
	var contents strings.Builder
	rendered := make(map[string]bool)
	rendering := make(map[string]bool)
	for _, dependency := range r.selectedDependencies() {
		if err := r.writeDependency(&contents, dependency, rendered, rendering); err != nil {
			return "", err
		}
	}
	return formatTestHCL(contents.String())
}

func (r testConfigurationRenderer) existingDataSourceHCL() (string, error) {
	var contents strings.Builder
	if r.class.HasResourceArtifact() {
		address := r.subjectResourceAddress()
		if err := r.writeClassBlock(&contents, address, r.class, fullTestResourceScenario, nil, r.class.TestChildren, "", nil); err != nil {
			return "", err
		}
	}
	if err := r.writeDataSourceBlock(&contents, "test", false); err != nil {
		return "", err
	}
	return formatTestHCL(contents.String())
}

func (r testConfigurationRenderer) missingDataSourceHCL() (string, error) {
	if !r.hasMissingDataSource() {
		return "", nil
	}

	var contents strings.Builder
	if err := r.writeDataSourceBlock(&contents, "test_non_existing", true); err != nil {
		return "", err
	}
	return formatTestHCL(contents.String())
}

func (r testConfigurationRenderer) writeDependency(contents *strings.Builder, dependency *TestDependency, rendered, rendering map[string]bool) error {
	if dependency == nil || dependency.ReferenceType == StaticReference {
		return nil
	}
	if rendered[dependency.Reference] || rendering[dependency.Reference] {
		return nil
	}
	rendering[dependency.Reference] = true
	defer delete(rendering, dependency.Reference)
	for _, prerequisite := range dependency.Dependencies {
		if err := r.writeDependency(contents, prerequisite, rendered, rendering); err != nil {
			return err
		}
	}

	dependencyClass, ok := r.dataStore.Classes[dependency.Class.String()]
	if !ok {
		if err := r.writeLegacyTestDependency(contents, dependency); err != nil {
			return fmt.Errorf("render test dependency %q: %w", dependency.Reference, err)
		}
		rendered[dependency.Reference] = true
		return nil
	}
	address, ok := parseTerraformTestReference(dependency.Reference)
	if !ok {
		return fmt.Errorf("render test dependency %q: expected a Terraform resource or data-source reference", dependency.Reference)
	}
	scenario := dependencyTestResourceScenario
	if dependency.TestInstance > 1 {
		scenario = alternateDependencyTestResourceScenario
	}
	children := sortedTestDependencyChildren(dependency.Children)
	for _, reference := range r.classBlockReferences(&dependencyClass, address.BlockType, scenario, dependency.ConfigOverrides, children) {
		prerequisite := findTestDependencyByObjectReference(dependencyClass.TestDependencies, reference)
		if prerequisite == nil {
			continue
		}
		if err := r.writeDependency(contents, prerequisite, rendered, rendering); err != nil {
			return err
		}
	}
	dependencyRenderer := r
	dependencyRenderer.dependency = dependency
	if err := dependencyRenderer.writeClassBlock(contents, address, &dependencyClass, scenario, dependency.ConfigOverrides, children, r.dependencyCountExpression(&dependencyClass), orderOnlyDependencyAddresses(dependency.Dependencies)); err != nil {
		return fmt.Errorf("render test dependency %q: %w", dependency.Reference, err)
	}
	rendered[dependency.Reference] = true
	return nil
}

func (r testConfigurationRenderer) classBlockReferences(class *Class, blockType string, resourceScenario testResourceScenario, overrides map[string]string, children []*TestChild) []string {
	var references []string
	appendReference := func(value string, valueType ValueRenderTypeEnum) {
		if valueType == ReferenceValue && !slices.Contains(references, value) {
			references = append(references, value)
		}
	}
	appendScenarioReferences := func(scenario TestValueScenario) {
		for _, entry := range scenario.Entries {
			if entry.ConfigInclude {
				appendReference(entry.ConfigValue, entry.ValueType)
			}
		}
	}

	for _, propertyName := range class.PropertiesAll {
		property := class.Properties[propertyName]
		if property == nil || property.ReadOnly || property.IgnoreInTest || property.IgnoreInDataSourceTest {
			continue
		}
		if blockType == "data" && property.PropertyName != "parentDn" && !class.IsIdentifyingProperty(property.PropertyName) {
			continue
		}
		if override, ok := testPropertyOverride(property, overrides); ok {
			appendReference(override, testOverrideValueType(override))
			continue
		}
		if property.TestValues == nil {
			continue
		}
		if blockType == "data" {
			appendScenarioReferences(property.TestValues.Create)
			continue
		}
		appendScenarioReferences(testResourcePropertyScenario(class, property, resourceScenario))
	}
	r.collectTestChildReferences(children, &references)
	return references
}

func (r testConfigurationRenderer) writeLegacyTestDependency(contents *strings.Builder, dependency *TestDependency) error {
	address, ok := parseTerraformTestReference(dependency.Reference)
	if !ok {
		return fmt.Errorf("legacy class %q requires a Terraform resource or data-source reference", dependency.Class)
	}
	if len(dependency.ConfigOverrides) == 0 {
		return fmt.Errorf("legacy class %q requires complete config_overrides", dependency.Class)
	}
	if len(dependency.Children) > 0 {
		return fmt.Errorf("legacy class %q cannot render child overrides without loaded class metadata", dependency.Class)
	}

	fmt.Fprintf(contents, "%s %q %q {\n", address.BlockType, address.Type, address.Name)
	contents.WriteString("count = 1\n")
	if dependencies := orderOnlyDependencyAddresses(dependency.Dependencies); len(dependencies) > 0 {
		fmt.Fprintf(contents, "depends_on = [%s]\n", strings.Join(dependencies, ", "))
	}
	for _, attributeName := range slices.Sorted(maps.Keys(dependency.ConfigOverrides)) {
		entry := TestValueEntry{
			ConfigValue:   dependency.ConfigOverrides[attributeName],
			ConfigInclude: true,
			ValueType:     testOverrideValueType(dependency.ConfigOverrides[attributeName]),
		}
		if attributeName == "name" && r.shouldPrefixLegacyName(dependency) {
			entry = prefixedTestEntry(entry)
		}
		// SDKv2 validates this relation during planning, when .id is still unknown.
		// Its configured .name keeps the DN known and follows any prefix override.
		if dependency.Class.String() == "vnsLDevVip" && attributeName == "relation_vns_rs_al_dev_to_phys_dom_p" {
			for _, target := range dependency.Dependencies {
				if target.Class.String() != "physDomP" || entry.ConfigValue != "uni/phys-"+target.ConfigOverrides["name"] {
					continue
				}
				address, ok := parseTerraformTestReference(target.Reference)
				if ok {
					entry.ConfigValue = fmt.Sprintf("format(\"uni/phys-%%s\", %s)", countedTestDependencyReference(address.withAttribute("name")))
					entry.ValueType = ExpressionValue
				}
			}
		}
		entry = r.prefixedLiteralDN(entry)
		fmt.Fprintf(contents, "%s = %s\n", attributeName, testValueHCL(entry))
	}
	contents.WriteString("}\n")
	return nil
}

func (r testConfigurationRenderer) writeClassBlock(
	contents *strings.Builder,
	address terraformTestAddress,
	class *Class,
	resourceScenario testResourceScenario,
	overrides map[string]string,
	children []*TestChild,
	countExpression string,
	dependsOn []string,
) error {
	fmt.Fprintf(contents, "%s %q %q {\n", address.BlockType, address.Type, address.Name)
	if countExpression != "" {
		fmt.Fprintf(contents, "count = %s\n", countExpression)
	}
	if address.BlockType == "data" {
		if err := r.writeLookupAttributes(contents, class, false, overrides); err != nil {
			return err
		}
		contents.WriteString("}\n")
		return nil
	}

	if len(dependsOn) > 0 {
		fmt.Fprintf(contents, "depends_on = [%s]\n", strings.Join(dependsOn, ", "))
	}
	usedOverrides := make(map[string]bool, len(overrides))
	for _, propertyName := range class.PropertiesAll {
		property := class.Properties[propertyName]
		if property == nil || property.ReadOnly || property.IgnoreInTest || property.IgnoreInDataSourceTest {
			continue
		}
		if override, ok := testPropertyOverride(property, overrides); ok {
			usedOverrides[property.AttributeName] = true
			usedOverrides[property.PropertyName] = true
			entry := TestValueEntry{
				ConfigValue:   override,
				ConfigInclude: true,
				ValueType:     testOverrideValueType(override),
			}
			writeTestAttribute(contents, property.AttributeName, property.ValueType, r.testPropertyScenario(class, property, newTestValueScenario([]TestValueEntry{entry})), testConfigurationVersions(class.SupportedVersions, property.SupportedVersions))
			continue
		}
		if property.TestValues == nil {
			continue
		}
		scenario := testResourcePropertyScenario(class, property, resourceScenario)
		scenario = r.testPropertyScenario(class, property, scenario)
		writeTestAttribute(contents, property.AttributeName, property.ValueType, scenario, testConfigurationVersions(class.SupportedVersions, property.SupportedVersions))
	}
	for key := range overrides {
		if usedOverrides[key] {
			continue
		}
		return fmt.Errorf("config override %q does not match a property on class %q", key, class.Name)
	}
	if err := r.writeTestChildren(contents, class, children); err != nil {
		return err
	}
	contents.WriteString("}\n")
	return nil
}

func orderOnlyDependencyAddresses(dependencies []*TestDependency) []string {
	var addresses []string
	for _, dependency := range dependencies {
		if dependency == nil || !dependency.OrderOnly || dependency.ReferenceType == StaticReference {
			continue
		}
		address, ok := parseTerraformTestReference(dependency.Reference)
		if !ok || slices.Contains(addresses, address.String()) {
			continue
		}
		addresses = append(addresses, address.String())
	}
	return addresses
}

func testResourcePropertyScenario(class *Class, property *Property, resourceScenario testResourceScenario) TestValueScenario {
	switch resourceScenario {
	case dependencyTestResourceScenario:
		return property.TestValues.Default
	case alternateDependencyTestResourceScenario:
		if class.IsIdentifyingProperty(property.PropertyName) && property.TestValues.Update.Defined {
			return property.TestValues.Update
		}
		if len(class.IdentifiedBy) == 0 && property.PropertyName == "parentDn" && property.TestValues.ForceNew.Defined {
			return property.TestValues.ForceNew
		}
		return property.TestValues.Default
	default:
		return property.TestValues.Create
	}
}

func (r testConfigurationRenderer) writeTestChildren(contents *strings.Builder, parentClass *Class, children []*TestChild) error {
	for _, child := range children {
		if child == nil || len(child.Instances) == 0 {
			continue
		}
		childClass, ok := r.dataStore.Classes[child.Class.String()]
		if !ok {
			return fmt.Errorf("render test child %q: class is not loaded", child.Class)
		}
		if slices.Contains(childClass.TestConfig.IgnoreTests, ChildIgnoreTest) {
			continue
		}

		versions := testConfigurationVersions(parentClass.SupportedVersions, childClass.SupportedVersions)
		fmt.Fprintf(contents, "%s = ", childClass.ResourceNameNested)
		if versions != "" {
			fmt.Fprintf(contents, "%s ? ", testVersionExpression(versions))
		}
		if childClass.IsSingleNestedWhenDefinedAsChild {
			contents.WriteString("{\n")
			if err := r.writeTestChildInstance(contents, &childClass, child.Instances[0]); err != nil {
				return err
			}
			contents.WriteString("}")
			if versions != "" {
				contents.WriteString(" : null")
			}
			contents.WriteString("\n")
			continue
		}

		contents.WriteString("[\n")
		for _, instance := range child.Instances {
			contents.WriteString("{\n")
			if err := r.writeTestChildInstance(contents, &childClass, instance); err != nil {
				return err
			}
			contents.WriteString("},\n")
		}
		contents.WriteString("]")
		if versions != "" {
			contents.WriteString(" : null")
		}
		contents.WriteString("\n")
	}
	return nil
}

func (r testConfigurationRenderer) writeTestChildInstance(contents *strings.Builder, class *Class, instance TestChildInstance) error {
	for _, propertyName := range class.PropertiesAll {
		property := class.Properties[propertyName]
		if property == nil || property.PropertyName == "parentDn" || property.ReadOnly || property.IgnoreInTest || property.IgnoreInDataSourceTest {
			continue
		}
		scenario, ok := instance.Properties[property.AttributeName]
		if !ok {
			continue
		}
		writeTestAttribute(contents, property.AttributeName, property.ValueType, scenario, testConfigurationVersions(class.SupportedVersions, property.SupportedVersions))
	}
	return r.writeTestChildren(contents, class, instance.Children)
}

func (r testConfigurationRenderer) writeDataSourceBlock(contents *strings.Builder, instanceName string, missing bool) error {
	fmt.Fprintf(contents, "data %q %q {\n", "aci_"+r.class.ResourceName, instanceName)
	if err := r.writeLookupAttributes(contents, r.class, missing, nil); err != nil {
		return err
	}
	if !missing && r.class.HasResourceArtifact() {
		fmt.Fprintf(contents, "depends_on = [%s]\n", r.subjectResourceAddress().String())
	}
	contents.WriteString("}\n")
	return nil
}

func (r testConfigurationRenderer) writeLookupAttributes(contents *strings.Builder, class *Class, missing bool, overrides map[string]string) error {
	missingIdentifier := ""
	missingScenario := TestValueScenario{}
	if missing && len(class.IdentifiedBy) > 0 {
		missingIdentifier, missingScenario, _ = r.missingIdentifierScenario()
	}
	usedOverrides := make(map[string]bool, len(overrides))

	for _, propertyName := range class.PropertiesAll {
		property := class.Properties[propertyName]
		if property == nil || (property.PropertyName != "parentDn" && !class.IsIdentifyingProperty(property.PropertyName)) {
			continue
		}
		if override, ok := testPropertyOverride(property, overrides); ok {
			usedOverrides[property.AttributeName] = true
			usedOverrides[property.PropertyName] = true
			scenario := newTestValueScenario([]TestValueEntry{{
				ConfigValue:   override,
				ConfigInclude: true,
				ValueType:     testOverrideValueType(override),
			}})
			writeTestAttribute(contents, property.AttributeName, property.ValueType, r.testPropertyScenario(class, property, scenario), "")
			continue
		}
		if property.TestValues == nil {
			if property.PropertyName == "parentDn" && class.ParentDnHasSchemaDefault() {
				continue
			}
			return fmt.Errorf("class %q lookup property %q has no test values", class.Name, property.AttributeName)
		}

		scenario := property.TestValues.Create
		switch {
		case missingIdentifier == property.PropertyName:
			scenario = missingScenario
		case missing && len(class.IdentifiedBy) == 0 && property.PropertyName == "parentDn":
			scenario, _ = r.missingParentScenario()
		}
		writeTestAttribute(contents, property.AttributeName, property.ValueType, r.testPropertyScenario(class, property, scenario), "")
	}
	for key := range overrides {
		if usedOverrides[key] {
			continue
		}
		return fmt.Errorf("config override %q does not match a lookup property on class %q", key, class.Name)
	}
	return nil
}

func (r testConfigurationRenderer) hasMissingDataSource() bool {
	if !r.class.HasResourceArtifact() {
		return false
	}
	if len(r.class.IdentifiedBy) > 0 {
		_, _, ok := r.missingIdentifierScenario()
		return ok
	}
	_, ok := r.missingParentScenario()
	return ok
}

func (r testConfigurationRenderer) missingParentScenario() (TestValueScenario, bool) {
	// An explicit missing parent also covers implicit singletons that APIC
	// creates automatically beneath every existing parent.
	parentDn := r.class.ParentDnProperty()
	if parentDn == nil || parentDn.TestValues == nil {
		return TestValueScenario{}, false
	}
	create, createOK := parentDn.TestValues.Create.firstIncludedEntry()
	if !createOK {
		return TestValueScenario{}, false
	}
	if missing, ok := parentDn.TestValues.DataSourceMissing.firstIncludedEntry(); ok && (create.ConfigValue != missing.ConfigValue || create.ValueType != missing.ValueType) {
		return newTestValueScenario([]TestValueEntry{missing}), true
	}
	forceNew, forceNewOK := parentDn.TestValues.ForceNew.firstIncludedEntry()
	if !forceNewOK || create.ConfigValue == forceNew.ConfigValue || forceNew.ValueType != ReferenceValue {
		return TestValueScenario{}, false
	}
	dependency := r.findDependency(forceNew.ConfigValue)
	if dependency == nil || dependency.ReferenceType == StaticReference {
		return TestValueScenario{}, false
	}
	return newTestValueScenario([]TestValueEntry{forceNew}), true
}

func (r testConfigurationRenderer) missingIdentifierScenario() (string, TestValueScenario, bool) {
	for _, propertyName := range r.class.PropertiesAll {
		if !r.class.IsIdentifyingProperty(propertyName) {
			continue
		}
		property := r.class.Properties[propertyName]
		if property == nil || property.TestValues == nil {
			continue
		}
		create, ok := property.TestValues.Create.firstIncludedEntry()
		if !ok {
			continue
		}
		if missing, ok := property.TestValues.DataSourceMissing.firstIncludedEntry(); ok && (create.ConfigValue != missing.ConfigValue || create.ValueType != missing.ValueType) {
			return propertyName, newTestValueScenario([]TestValueEntry{missing}), true
		}
		if propertyUsesDNTestValue(property) {
			continue
		}
		if update, ok := property.TestValues.Update.firstIncludedEntry(); ok && (create.ConfigValue != update.ConfigValue || create.ValueType != update.ValueType) {
			return propertyName, newTestValueScenario([]TestValueEntry{update}), true
		}
		if missing, ok := nonExistingTestEntry(property, create); ok {
			return propertyName, newTestValueScenario([]TestValueEntry{missing}), true
		}
	}
	return "", TestValueScenario{}, false
}

func propertyUsesDNTestValue(property *Property) bool {
	return property.PropertyName == "tDn" || property.Reference != nil && property.Reference.ExampleAttributeName == "id"
}

func nonExistingTestEntry(property *Property, create TestValueEntry) (TestValueEntry, bool) {
	if property.ValueType == Set || property.ValueType == Object {
		return TestValueEntry{}, false
	}
	if value, ok := nonExistingNumericTestValue(property, create.ConfigValue); ok {
		return TestValueEntry{ConfigValue: value, ConfigInclude: true, ValueType: StringValue}, true
	}
	if property.ValueType == IpAddress {
		if value, ok := nonExistingIPTestValue(create.ConfigValue); ok {
			return TestValueEntry{ConfigValue: value, ConfigInclude: true, ValueType: StringValue}, true
		}
		return TestValueEntry{}, false
	}
	for _, value := range property.ValidValues.LocalNamesList() {
		if value != create.ConfigValue {
			return TestValueEntry{ConfigValue: value, ConfigInclude: true, ValueType: StringValue}, true
		}
	}
	if len(property.ValidValues) > 0 {
		return TestValueEntry{}, false
	}
	value := property.AttributeName + "_2"
	if value == create.ConfigValue {
		value = property.AttributeName + "_3"
	}
	return TestValueEntry{ConfigValue: value, ConfigInclude: true, ValueType: StringValue}, true
}

func nonExistingNumericTestValue(property *Property, value string) (string, bool) {
	uiType, _ := property.metaDetails["uitype"].(string)
	if uiType != "number" {
		return "", false
	}
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return "", false
	}
	if limits := property.numericRange(); limits != nil && number >= limits[1] {
		if number <= limits[0] {
			return "", false
		}
		return strconv.FormatInt(number-1, 10), true
	}
	return strconv.FormatInt(number+1, 10), true
}

func nonExistingIPTestValue(value string) (string, bool) {
	if prefix, err := netip.ParsePrefix(value); err == nil {
		next := prefix.Addr().Next()
		if !next.IsValid() {
			return "", false
		}
		return netip.PrefixFrom(next, prefix.Bits()).String(), true
	}
	address, err := netip.ParseAddr(value)
	if err != nil {
		return "", false
	}
	next := address.Next()
	if !next.IsValid() {
		return "", false
	}
	return next.String(), true
}

func (r testConfigurationRenderer) selectedDependencies() []*TestDependency {
	selected := make(map[*TestDependency]bool)
	primaryParent := r.primaryParentDependency()
	if primaryParent != nil {
		selected[primaryParent] = true
	}
	for _, dependency := range r.semanticReferenceDependencies() {
		selected[dependency] = true
	}
	references := r.configurationReferences()
	for _, reference := range references {
		if dependency := r.findDependency(reference); dependency != nil {
			selected[dependency] = true
		}
	}

	dependencies := make([]*TestDependency, 0, len(selected))
	for _, dependency := range r.class.TestDependencies {
		if selected[dependency] {
			dependencies = append(dependencies, dependency)
			delete(selected, dependency)
		}
	}
	for _, reference := range references {
		dependency := r.findDependency(reference)
		if dependency != nil && selected[dependency] {
			dependencies = append(dependencies, dependency)
			delete(selected, dependency)
		}
	}
	return dependencies
}

func (r testConfigurationRenderer) semanticReferenceDependencies() []*TestDependency {
	var dependencies []*TestDependency
	for _, propertyName := range r.class.PropertiesAll {
		property := r.class.Properties[propertyName]
		if property == nil || property.Reference == nil || property.TestValues == nil || !property.TestValues.Create.Defined {
			continue
		}
		for _, referenceClass := range property.Reference.Classes {
			for _, dependency := range r.class.TestDependencies {
				if dependency.Role == Target && dependency.Class.Equal(referenceClass) && !slices.Contains(dependencies, dependency) {
					dependencies = append(dependencies, dependency)
					break
				}
			}
		}
	}
	return dependencies
}

func (r testConfigurationRenderer) primaryParentDependency() *TestDependency {
	if parentDn := r.class.ParentDnProperty(); parentDn != nil && parentDn.TestValues != nil {
		if entry, ok := parentDn.TestValues.Create.firstIncludedEntry(); ok && entry.ValueType == ReferenceValue {
			if dependency := r.findDependency(entry.ConfigValue); dependency != nil {
				return dependency
			}
		}
	}
	for _, dependency := range r.class.TestDependencies {
		if dependency.Role == Parent {
			return dependency
		}
	}
	return nil
}

func (r testConfigurationRenderer) configurationReferences() []string {
	var references []string
	appendScenarioReferences := func(scenario TestValueScenario) {
		for _, entry := range scenario.Entries {
			if entry.ConfigInclude && entry.ValueType == ReferenceValue && !slices.Contains(references, entry.ConfigValue) {
				references = append(references, entry.ConfigValue)
			}
		}
	}
	for _, propertyName := range r.class.PropertiesAll {
		property := r.class.Properties[propertyName]
		if property == nil || property.TestValues == nil || property.ReadOnly || property.IgnoreInTest || property.IgnoreInDataSourceTest {
			continue
		}
		appendScenarioReferences(property.TestValues.Create)
	}
	r.collectTestChildReferences(r.class.TestChildren, &references)
	if _, missingScenario, ok := r.missingIdentifierScenario(); ok {
		appendScenarioReferences(missingScenario)
	} else if len(r.class.IdentifiedBy) == 0 && r.class.HasResourceArtifact() {
		if missingParent, ok := r.missingParentScenario(); ok {
			appendScenarioReferences(missingParent)
		}
	}
	return references
}

func (r testConfigurationRenderer) collectTestChildReferences(children []*TestChild, references *[]string) {
	for _, child := range children {
		childClass, ok := r.dataStore.Classes[child.Class.String()]
		if ok && slices.Contains(childClass.TestConfig.IgnoreTests, ChildIgnoreTest) {
			continue
		}
		for _, instance := range child.Instances {
			for _, propertyName := range slices.Sorted(maps.Keys(instance.Properties)) {
				for _, entry := range instance.Properties[propertyName].Entries {
					if entry.ConfigInclude && entry.ValueType == ReferenceValue && !slices.Contains(*references, entry.ConfigValue) {
						*references = append(*references, entry.ConfigValue)
					}
				}
			}
			r.collectTestChildReferences(instance.Children, references)
		}
	}
}

func (r testConfigurationRenderer) findDependency(reference string) *TestDependency {
	return findTestDependencyByObjectReference(r.class.TestDependencies, reference)
}

func (r testConfigurationRenderer) subjectResourceAddress() terraformTestAddress {
	return terraformTestAddress{BlockType: "resource", Type: "aci_" + r.class.ResourceName, Name: "test_1"}
}

func writeTestAttribute(contents *strings.Builder, attributeName string, valueType ValueTypeEnum, scenario TestValueScenario, versions string) {
	if !scenario.Defined {
		return
	}
	entries := scenario.includedEntries()
	if valueType == Set {
		var value strings.Builder
		value.WriteString("[")
		for index, entry := range entries {
			if index > 0 {
				value.WriteString(", ")
			}
			value.WriteString(testValueHCL(entry))
		}
		value.WriteString("]")
		writeTestAssignment(contents, attributeName, value.String(), versions)
		return
	}
	if len(entries) == 0 {
		return
	}
	value := testValueHCL(entries[0])
	if len(entries) > 1 && slices.ContainsFunc(entries, func(entry TestValueEntry) bool { return entry.Versions != nil }) {
		value = "null"
		for index := len(entries) - 1; index >= 0; index-- {
			entry := entries[index]
			if entry.Versions == nil {
				value = testValueHCL(entry)
				continue
			}
			value = fmt.Sprintf("%s ? %s : %s", testVersionExpression(entry.Versions.Raw()), testValueHCL(entry), value)
		}
	}
	writeTestAssignment(contents, attributeName, value, versions)
}

func writeTestAssignment(contents *strings.Builder, attributeName, value, versions string) {
	if versions == "" {
		fmt.Fprintf(contents, "%s = %s\n", attributeName, value)
		return
	}
	fmt.Fprintf(contents, "%s = %s ? %s : null\n", attributeName, testVersionExpression(versions), value)
}

func testVersionExpression(versions string) string {
	return fmt.Sprintf("provider::aci::compare_versions(var.apic_version, %q, %q)", "inside", versions)
}

func testConfigurationVersions(ownerVersions, valueVersions *Versions) string {
	if valueVersions == nil || valueVersions.Covers(ownerVersions) {
		return ""
	}
	return valueVersions.Raw()
}

func (r testConfigurationRenderer) dependencyCountExpression(class *Class) string {
	// Keep every dependency counted while omitting classes unavailable on the active APIC version.
	versions := testConfigurationVersions(r.class.SupportedVersions, class.SupportedVersions)
	if versions == "" {
		return "1"
	}
	return testVersionExpression(versions) + " ? 1 : 0"
}

func countedTestDependencyReference(reference string) string {
	// All non-static dependencies are counted, including those whose count is always one.
	address, ok := parseTerraformTestReference(reference)
	if !ok {
		return reference
	}
	countedReference := address.String() + "[0]"
	if address.Attribute != "" {
		countedReference += "." + address.Attribute
	}
	return countedReference
}

func testValueHCL(entry TestValueEntry) string {
	if entry.ValueType == ReferenceValue {
		return countedTestDependencyReference(entry.ConfigValue)
	}
	if entry.ValueType == ExpressionValue {
		return entry.ConfigValue
	}
	return string(hclwrite.TokensForValue(cty.StringVal(entry.ConfigValue)).Bytes())
}

func testPropertyOverride(property *Property, overrides map[string]string) (string, bool) {
	if value, ok := overrides[property.AttributeName]; ok {
		return value, true
	}
	value, ok := overrides[property.PropertyName]
	return value, ok
}

func testOverrideValueType(value string) ValueRenderTypeEnum {
	if _, ok := parseTerraformTestReference(value); ok {
		return ReferenceValue
	}
	trimmedValue := strings.TrimSpace(value)
	if strings.HasPrefix(trimmedValue, "[") || strings.HasPrefix(trimmedValue, "{") {
		return ExpressionValue
	}
	return StringValue
}

func sortedTestDependencyChildren(childrenByClass map[string]*TestChild) []*TestChild {
	children := make([]*TestChild, 0, len(childrenByClass))
	for _, className := range slices.Sorted(maps.Keys(childrenByClass)) {
		children = append(children, childrenByClass[className])
	}
	return children
}

func formatTestHCL(contents string) (string, error) {
	if strings.TrimSpace(contents) == "" {
		return "", nil
	}
	formatted := hclwrite.Format([]byte(contents))
	if _, diagnostics := hclsyntax.ParseConfig(formatted, "generated-test.tf", hcl.InitialPos); diagnostics.HasErrors() {
		return "", fmt.Errorf("format generated test configuration: %s", diagnostics.Error())
	}
	if strings.Contains(string(formatted), "`") {
		return "", fmt.Errorf("generated test configuration contains a backtick and cannot be embedded in a Go raw string")
	}
	return strings.TrimSpace(string(formatted)), nil
}
