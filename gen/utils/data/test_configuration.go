package data

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/netip"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils"
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

// ResourceTestScenario is one step in the generated resource's lifecycle.
type ResourceTestScenario struct {
	Name           string
	Description    string
	Configuration  string
	Assertions     []TestAssertion
	Action         string
	ChangeVersions [][]string
	PlanOnly       bool
	Import         bool
	ExternalDelete bool
	Remove         bool
	IDChanges      bool
	ExpectError    string
	ExistingObject bool
	ForgetExisting bool
}

type resourceTestBuilder struct {
	class     *Class
	dataStore *DataStore
	scenarios []ResourceTestScenario
	state     *Class
}

func (c *Class) ResourceTestDependenciesHCL(ds *DataStore) (string, error) {
	class := c.prepareResourceTest(ds)
	renderer := newTestConfigurationRenderer(&class, ds)
	renderer.resourceTest = true
	return renderer.dependenciesHCL()
}

func (c *Class) ResourceTestScenarios(ds *DataStore) ([]ResourceTestScenario, error) {
	// Project the existing buckets without changing canonical values used by other artifacts.
	class := c.prepareResourceTest(ds)
	c = &class
	builder := resourceTestBuilder{class: c, dataStore: ds}
	minimum := c.resourceTestCopy()
	minimum.TestChildren = builder.minimumChildren(minimum.TestChildren)
	expected := minimum.resourceTestCopy()
	for _, name := range c.PropertiesAll {
		property := minimum.Properties[name]
		if property.TestValues == nil || property.ReadOnly {
			continue
		}
		if !property.Required && !property.CreateOnly && !property.resourceIncludedIn("minimum") && !c.IsIdentifyingProperty(name) {
			property.TestValues.Create = TestValueScenario{}
			expected.Properties[name].TestValues.Create = resourceDefaultExpectation(property)
		}
	}
	minimumConfiguration, minimumExpected := builder.projectScenario("Minimum", &minimum, &expected)
	if err := builder.addScenario("Minimum", "Create Minimum Object and Verify Defaults", minimumConfiguration, minimumExpected, "Create"); err != nil {
		return nil, err
	}
	if err := builder.existingObject(minimumConfiguration); err != nil {
		return nil, err
	}
	current := c.resourceTestCopy()
	if err := builder.add("Full", "Apply Full Configuration", &current, &current, ""); err != nil {
		return nil, err
	}
	current = builder.state.resourceTestCopy()
	if err := builder.unchanged("FullUnchanged", "Repeat Full Configuration (Expect Empty Plan)", &current); err != nil {
		return nil, err
	}
	if err := builder.legacy(&current); err != nil {
		return nil, err
	}
	updated := current.resourceTestCopy()
	for _, name := range c.PropertiesAll {
		property := updated.Properties[name]
		if resourceMutableProperty(&updated, property) && property.TestValues.Update.Defined {
			property.TestValues.Create = property.TestValues.Update.clone()
		}
	}
	if changed, err := builder.change("Update", "Update Mutable Attributes in Place", &updated, "Update"); err != nil {
		return nil, err
	} else if changed {
		current = updated
		if err := builder.unchanged("UpdateUnchanged", "Repeat Updated Configuration (Expect Empty Plan)", &current); err != nil {
			return nil, err
		}
	}

	omitted := current.resourceTestCopy()
	expected = current.resourceTestCopy()
	for _, name := range c.PropertiesAll {
		property := omitted.Properties[name]
		if !property.Required && !property.resourceIncludedIn("omit_optional") && resourceMutableProperty(&omitted, property) {
			property.TestValues.Create = TestValueScenario{}
			if name == "annotation" {
				expected.Properties[name].TestValues.Create = resourceDefaultExpectation(property)
			} else if !property.Computed {
				expected.Properties[name].TestValues.Create = TestValueScenario{}
			}
		}
	}
	if err := builder.add("OmitOptional", "Omit Optional Attributes and Verify Retained Values", &omitted, &expected, ""); err != nil {
		return nil, err
	}
	// Reapply explicit values before resets so expected state does not depend on an omitted annotation.
	if err := builder.add("RestoreOptional", "Restore Explicit Optional Attributes", &current, &current, ""); err != nil {
		return nil, err
	}
	reset := current.resourceTestCopy()
	for _, name := range c.PropertiesAll {
		property := reset.Properties[name]
		if resourceMutableProperty(&reset, property) && !property.Required && (!property.Sensitive || property.resourceIncludedIn("reset") && property.propertyDefinition.TestConfig.Default.Defined) {
			if scenario, ok := property.resourceResetScenario(); ok {
				property.TestValues.Create = scenario
			}
		}
	}
	if changed, err := builder.change("Reset", "Explicitly Reset Optional Attributes", &reset, "Update"); err != nil {
		return nil, err
	} else if changed {
		current = reset
	}

	// Follow the legacy ordering so apply steps persist computed allocations before import comparison.
	if err := builder.add("Import", "Import by Distinguished Name", &current, &current, ""); err != nil {
		return nil, err
	}
	builder.last().Import = true

	if err := builder.children(&current); err != nil {
		return nil, err
	}
	if c.AllowDelete {
		if err := builder.add("ExternalDeletion", "Delete Remotely and Verify Terraform Plans Recreation", &current, &current, "Create"); err != nil {
			return nil, err
		}
		builder.last().ExternalDelete = true
		builder.last().PlanOnly = true
		if err := builder.add("Recreate", "Recreate the Externally Deleted Object", &current, &current, "Create"); err != nil {
			return nil, err
		}
	}

	replaced := current.resourceTestCopy()
	for _, name := range c.IdentifiedBy {
		property := replaced.Properties[name]
		if property != nil && property.RequiresReplace && property.TestValues != nil && property.TestValues.Update.Defined {
			property.TestValues.Create = property.TestValues.Update.clone()
		}
	}
	builder.compatibleParentTargets(&replaced)
	if changed, err := builder.change("ReplaceIdentifier", "Replace After Changing an Identifier", &replaced, "Replace"); err != nil {
		return nil, err
	} else if changed {
		builder.last().IDChanges = true
		current = replaced
	}
	if parent := current.ParentDnProperty(); parent != nil && parent.RequiresReplace && parent.TestValues != nil {
		replaced = current.resourceTestCopy()
		replaced.Properties["parentDn"].TestValues.Create = parent.TestValues.ForceNew.clone()
		builder.compatibleParentTargets(&replaced)
		if changed, err := builder.change("ReplaceParent", "Replace Under an Alternate Parent", &replaced, "Replace"); err != nil {
			return nil, err
		} else if changed {
			builder.last().IDChanges = true
		}
	}
	if err := builder.add("Destroy", "Remove Resource and Verify Deletion Behavior", nil, nil, "Destroy"); err != nil {
		return nil, err
	}
	builder.last().Remove = true
	return builder.scenarios, nil
}

type testConfigurationRenderer struct {
	class        *Class
	dataStore    *DataStore
	dependency   *TestDependency
	resourceTest bool
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
	if r.resourceTest {
		scenario = r.resourceReferenceScenario(class, scenario)
	}
	return scenario
}

func (r testConfigurationRenderer) prefixedLiteralDN(entry TestValueEntry) TestValueEntry {
	format, reference, matched := r.prefixedLiteralDNParts(entry)
	if matched {
		entry.ConfigValue = fmt.Sprintf("format(%q, %s)", format, reference)
		entry.ValueType = ExpressionValue
	}
	return entry
}

func (r testConfigurationRenderer) prefixedLiteralDNParts(entry TestValueEntry) (string, string, bool) {
	if entry.ValueType != StringValue || !strings.HasPrefix(entry.ConfigValue, "uni/") {
		return "", "", false
	}
	// Configuration and independent resource assertions must follow the same prerequisite graph.
	visited := make(map[*TestDependency]bool)
	var find func([]*TestDependency) (string, string, bool)
	find = func(dependencies []*TestDependency) (string, string, bool) {
		for _, dependency := range dependencies {
			if dependency == nil || visited[dependency] {
				continue
			}
			visited[dependency] = true
			if format, reference, matched := r.prefixedDependencyDNParts(entry, dependency); matched {
				return format, reference, true
			}
			if format, reference, matched := find(dependency.Dependencies); matched {
				return format, reference, true
			}
		}
		return "", "", false
	}
	return find(r.selectedDependencies())
}

func (r testConfigurationRenderer) prefixedDependencyDNParts(entry TestValueEntry, dependency *TestDependency) (string, string, bool) {
	class, loaded := r.dataStore.Classes[dependency.Class.String()]
	if !loaded || len(class.IdentifiedBy) != 1 {
		return "", "", false
	}
	dependencyRenderer := r
	dependencyRenderer.dependency = dependency
	name := class.Properties["name"]
	if name == nil || name.TestValues == nil || !dependencyRenderer.shouldPrefixProperty(&class, name) {
		return "", "", false
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
		return "", "", false
	}
	address, ok := parseTerraformTestReference(dependency.Reference)
	if !ok {
		return "", "", false
	}
	format := fmt.Sprintf("%s/%s%s", parentDN, strings.ReplaceAll(class.RnFormat, "{name}", "%s"), strings.TrimPrefix(entry.ConfigValue, root))
	return format, countedTestDependencyReference(address.withAttribute("name")), true
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
		if property == nil || property.ReadOnly || property.IgnoreInTest || !r.resourceTest && property.IgnoreInDataSourceTest {
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
	writeTestDependencyLifecycle(contents, dependency)
	contents.WriteString("}\n")
	return nil
}

func writeTestDependencyLifecycle(contents *strings.Builder, dependency *TestDependency) {
	if dependency == nil || dependency.ReferenceType != ResourceReference || len(dependency.IgnoreChanges) == 0 {
		return
	}
	contents.WriteString("lifecycle {\n")
	contents.WriteString("# This dependency reads these attributes into state, but standalone resources manage them.\n")
	contents.WriteString("# Ignore the dependency's copy to avoid competing configuration changes after refresh.\n")
	fmt.Fprintf(contents, "ignore_changes = [%s]\n", strings.Join(dependency.IgnoreChanges, ", "))
	contents.WriteString("}\n")
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
		if property == nil || property.ReadOnly || property.IgnoreInTest || !r.resourceTest && property.IgnoreInDataSourceTest {
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
	writeTestDependencyLifecycle(contents, r.dependency)
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
		if child == nil || len(child.Instances) == 0 && !r.resourceTest {
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
		if len(child.Instances) == 0 {
			if childClass.IsSingleNestedWhenDefinedAsChild {
				contents.WriteString("{}")
			} else {
				contents.WriteString("[]")
			}
			if versions != "" {
				contents.WriteString(" : null")
			}
			contents.WriteString("\n")
			continue
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
		if property == nil || property.PropertyName == "parentDn" || property.ReadOnly || property.IgnoreInTest || !r.resourceTest && property.IgnoreInDataSourceTest {
			continue
		}
		scenario, ok := instance.Properties[property.AttributeName]
		if !ok {
			continue
		}
		if r.resourceTest {
			scenario = r.resourceReferenceScenario(class, scenario)
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
		if property == nil || property.TestValues == nil || property.ReadOnly || property.IgnoreInTest || !r.resourceTest && property.IgnoreInDataSourceTest {
			continue
		}
		appendScenarioReferences(property.TestValues.Create)
		if r.resourceTest {
			appendScenarioReferences(property.TestValues.Update)
			appendScenarioReferences(property.TestValues.Default)
			appendScenarioReferences(property.TestValues.ForceNew)
		}
	}
	r.collectTestChildReferences(r.class.TestChildren, &references)
	if r.resourceTest {
		r.collectResourceChildUpdateReferences(r.class.TestChildren, &references)
		if parent := r.class.ParentDnProperty(); parent != nil && parent.TestValues != nil {
			// Include the compatible target chosen for the alternate parent, rather than every unused target fixture.
			replacement := r.class.resourceTestCopy()
			replacement.Properties["parentDn"].TestValues.Create = parent.TestValues.ForceNew.clone()
			builder := resourceTestBuilder{class: r.class, dataStore: r.dataStore}
			builder.compatibleParentTargets(&replacement)
			for _, name := range replacement.PropertiesAll {
				if property := replacement.Properties[name]; property.TestValues != nil {
					appendScenarioReferences(property.TestValues.Create)
				}
			}
			r.collectTestChildReferences(replacement.TestChildren, &references)
		}
		return references
	}
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

func (c *Class) resourceTestCopy() Class {
	copy := *c
	copy.TestDependencies = slices.Clone(c.TestDependencies)
	copy.Properties = make(map[string]*Property, len(c.Properties))
	for name, property := range c.Properties {
		p := *property
		if property.TestValues != nil {
			values := *property.TestValues
			values.Create = values.Create.clone()
			p.TestValues = &values
		}
		copy.Properties[name] = &p
	}
	copy.TestChildren = make([]*TestChild, len(c.TestChildren))
	for index, child := range c.TestChildren {
		copy.TestChildren[index] = child.clone()
	}
	return copy
}

func resourceMutableProperty(class *Class, property *Property) bool {
	return property != nil && property.TestValues != nil && !property.ReadOnly && !property.IgnoreInTest && !property.CreateOnly && !property.RequiresReplace && !class.IsIdentifyingProperty(property.PropertyName) && property.PropertyName != "parentDn"
}

func resourceDefaultExpectation(property *Property) TestValueScenario {
	// APIC cannot supply an initial secret, even when its schema preserves known values after omission.
	if property.Sensitive {
		return TestValueScenario{}
	}
	scenario := property.TestValues.Default.clone()
	if property.ValueType == Set {
		scenario.Entries = slices.DeleteFunc(scenario.Entries, func(entry TestValueEntry) bool { return entry.AssertValue == "" })
	}
	for index := range scenario.Entries {
		scenario.Entries[index].ConfigInclude = true
	}
	return scenario
}

func (p *Property) resourceResetScenario() (TestValueScenario, bool) {
	// Only use documented or explicitly configured defaults; inferred empty strings can violate APIC validation.
	if len(p.Documentation.DefaultValues) == 0 && !p.propertyDefinition.TestConfig.Default.Defined && p.PropertyName != "annotation" {
		return TestValueScenario{}, false
	}
	scenario := resourceDefaultExpectation(p)
	if p.Sensitive && p.resourceIncludedIn("reset") && p.propertyDefinition.TestConfig.Default.Defined {
		// Clearing a secret requires an explicitly authored reset value and inclusion.
		scenario = p.TestValues.Default.clone()
		for index := range scenario.Entries {
			scenario.Entries[index].ConfigInclude = true
		}
	}
	if p.PropertyName == "annotation" {
		return newTestValueScenario([]TestValueEntry{{ConfigValue: "orchestrator:terraform", AssertValue: "orchestrator:terraform", ConfigInclude: true, ValueType: StringValue}}), true
	}
	for index := range scenario.Entries {
		entry := &scenario.Entries[index]
		if entry.ValueType != StringValue || len(p.ValidValues) > 0 && !slices.Contains(p.ValidValues.LocalNamesList(), entry.AssertValue) {
			return TestValueScenario{}, false
		}
		entry.ConfigValue = entry.AssertValue
	}
	return scenario, scenario.Defined
}

func (b *resourceTestBuilder) last() *ResourceTestScenario {
	return &b.scenarios[len(b.scenarios)-1]
}

func (b *resourceTestBuilder) existingObject(minimum *Class) error {
	// Keep the original minimum resource managed while testing a second address for the same APIC object.
	minimumHCL := b.last().Configuration
	renderer := newTestConfigurationRenderer(minimum, b.dataStore)
	renderer.resourceTest = true
	address := renderer.subjectResourceAddress()
	address.Name = "existing_on_create"
	var duplicate strings.Builder
	if err := renderer.writeClassBlock(&duplicate, address, minimum, fullTestResourceScenario, nil, minimum.TestChildren, "", nil); err != nil {
		return err
	}
	for _, scenario := range []ResourceTestScenario{
		{Name: "RejectExisting", Description: "Reject Existing Object", ExistingObject: true, ExpectError: "Object Already Exists"},
		{Name: "AllowExisting", Description: "Allow Existing Object", ExistingObject: true, Action: "Create"},
		{Name: "ForgetExisting", Description: "Forget Temporary State Entry", ForgetExisting: true},
	} {
		var contents strings.Builder
		contents.WriteString(minimumHCL)
		fmt.Fprintf(&contents, "\n\nprovider \"aci\" {\nallow_existing_on_create = %t\n}\n", scenario.ExpectError == "")
		if scenario.ForgetExisting {
			// Forget the duplicate state entry without deleting the object still owned by test_1.
			fmt.Fprintf(&contents, "\nremoved {\nfrom = %s\nlifecycle {\ndestroy = false\n}\n}\n", address.String())
		} else {
			contents.WriteString("\n")
			contents.WriteString(duplicate.String())
		}
		formatted, err := formatTestHCL(contents.String())
		if err != nil {
			return err
		}
		scenario.Configuration = formatted
		b.scenarios = append(b.scenarios, scenario)
	}
	return nil
}

func (b *resourceTestBuilder) add(name, description string, configuration, expected *Class, action string) error {
	configuration, expected = b.projectScenario(name, configuration, expected)
	return b.addScenario(name, description, configuration, expected, action)
}

func (b *resourceTestBuilder) addScenario(name, description string, configuration, expected *Class, action string) error {
	scenario := ResourceTestScenario{Name: name, Description: description, Action: action}
	if configuration != nil {
		renderer := newTestConfigurationRenderer(configuration, b.dataStore)
		renderer.resourceTest = true
		var contents strings.Builder
		if err := renderer.writeClassBlock(&contents, renderer.subjectResourceAddress(), configuration, fullTestResourceScenario, nil, configuration.TestChildren, "", nil); err != nil {
			return fmt.Errorf("resource test %s/%s: %w", b.class.Name, name, err)
		}
		formatted, err := formatTestHCL(contents.String())
		if err != nil {
			return err
		}
		scenario.Configuration = formatted
		scenario.Assertions = expected.resourceTestAssertions(b.dataStore)
	}
	b.scenarios = append(b.scenarios, scenario)
	b.state = expected
	return nil
}

func (b *resourceTestBuilder) unchanged(name, description string, current *Class) error {
	if err := b.add(name, description, current, current, ""); err != nil {
		return err
	}
	b.last().PlanOnly = true
	return nil
}

func (b *resourceTestBuilder) change(name, description string, after *Class, action string) (bool, error) {
	configuration, expected := b.projectScenario(name, after, after)
	versions := resourceTestChanges(b.state, expected, b.dataStore)
	if len(versions) == 0 {
		return false, nil
	}
	if action == "Replace" {
		// Every new identifying reference must exist before a replacement can be attempted.
		var requiredVersions []string
		renderer := newTestConfigurationRenderer(after, b.dataStore)
		for _, name := range after.PropertiesAll {
			property := after.Properties[name]
			if property.TestValues == nil || !property.RequiresReplace {
				continue
			}
			for _, entry := range property.TestValues.Create.includedEntries() {
				if entry.ValueType == ReferenceValue {
					if dependency := renderer.findDependency(entry.ConfigValue); dependency != nil {
						if dependencyClass, ok := b.dataStore.Classes[dependency.Class.String()]; ok {
							requiredVersions = testAssertionVersions(requiredVersions, after.SupportedVersions, dependencyClass.SupportedVersions)
						}
					}
				}
			}
		}
		for _, group := range versions {
			for _, value := range group {
				if !slices.Contains(requiredVersions, value) {
					requiredVersions = append(requiredVersions, value)
				}
			}
		}
		versions = [][]string{requiredVersions}
	}
	if err := b.addScenario(name, description, configuration, expected, action); err != nil {
		return false, err
	}
	// Later steps start from expected state, including retained values after omission.
	*after = *expected
	if !slices.ContainsFunc(versions, func(group []string) bool { return len(group) == 0 }) {
		b.last().ChangeVersions = versions
	}
	return true, nil
}

func (p *Property) resourceIncludedIn(scenario string) bool {
	return slices.Contains(p.propertyDefinition.TestConfig.IncludeIn, scenario)
}

func (p *Property) resourceScenarioPresence(scenario string) (include, exclude bool) {
	include = p.resourceIncludedIn(scenario)
	exclude = slices.Contains(p.propertyDefinition.TestConfig.ExcludeFrom, scenario)
	if !include && !exclude && strings.HasSuffix(scenario, "_unchanged") {
		// A repeat checks the preceding configuration, including its presence overrides.
		return p.resourceScenarioPresence(strings.TrimSuffix(scenario, "_unchanged"))
	}
	return include, exclude
}

func (b *resourceTestBuilder) projectScenario(name string, configuration, expected *Class) (*Class, *Class) {
	if configuration == nil {
		return nil, expected
	}
	configured := configuration.resourceTestCopy()
	wanted := expected.resourceTestCopy()
	scenario := utils.Underscore(name)
	apply := func(property *Property, value, expectation, previous TestValueScenario) (TestValueScenario, TestValueScenario) {
		include, exclude := property.resourceScenarioPresence(scenario)
		if exclude {
			expectation = TestValueScenario{}
			if property.Computed {
				expectation = previous.clone()
				if !expectation.Defined {
					expectation = resourceDefaultExpectation(property)
				}
			}
			return TestValueScenario{}, expectation
		}
		if !include {
			return value, expectation
		}
		if scenario == "reset" && !property.CreateOnly {
			if reset, ok := property.resourceResetScenario(); ok {
				return reset, reset.clone()
			}
		}
		if !value.Defined {
			value = expectation.clone()
			if !value.Defined {
				value = property.TestValues.Create.clone()
			}
			for index := range value.Entries {
				value.Entries[index].ConfigInclude = true
			}
		}
		return value, value.clone()
	}
	for _, propertyName := range configured.PropertiesAll {
		property := configured.Properties[propertyName]
		if property == nil || property.ReadOnly || property.TestValues == nil {
			continue
		}
		previous := TestValueScenario{}
		if b.state != nil {
			if old := b.state.Properties[propertyName]; old != nil && old.TestValues != nil {
				previous = old.TestValues.Create
			}
		}
		property.TestValues.Create, wanted.Properties[propertyName].TestValues.Create = apply(property, property.TestValues.Create, wanted.Properties[propertyName].TestValues.Create, previous)
	}
	var children func([]*TestChild, []*TestChild, []*TestChild)
	children = func(configured, wanted, previous []*TestChild) {
		for _, child := range configured {
			childClass := b.dataStore.Classes[child.Class.String()]
			expectedChild := findResourceTestChild(wanted, child.Class)
			oldChild := findResourceTestChild(previous, child.Class)
			for index := range child.Instances {
				instance := &child.Instances[index]
				expectedInstance := &expectedChild.Instances[index]
				var old *TestChildInstance
				if oldChild != nil {
					for oldIndex := range oldChild.Instances {
						candidate := &oldChild.Instances[oldIndex]
						if slices.ContainsFunc(childClass.IdentifiedBy, func(name string) bool {
							attribute := childClass.Properties[name].AttributeName
							return !reflect.DeepEqual(candidate.Properties[attribute], instance.Properties[attribute])
						}) {
							continue
						}
						old = candidate
						break
					}
				}
				for _, name := range childClass.PropertiesAll {
					property := childClass.Properties[name]
					if property == nil || property.ReadOnly || property.TestValues == nil {
						continue
					}
					before := TestValueScenario{}
					if old != nil {
						before = old.Properties[property.AttributeName]
					}
					instance.Properties[property.AttributeName], expectedInstance.Properties[property.AttributeName] = apply(property, instance.Properties[property.AttributeName], expectedInstance.Properties[property.AttributeName], before)
				}
				var oldChildren []*TestChild
				if old != nil {
					oldChildren = old.Children
				}
				children(instance.Children, expectedInstance.Children, oldChildren)
			}
		}
	}
	var previousChildren []*TestChild
	if b.state != nil {
		previousChildren = b.state.TestChildren
	}
	children(configured.TestChildren, wanted.TestChildren, previousChildren)
	return &configured, &wanted
}

func (b *resourceTestBuilder) minimumChildren(children []*TestChild) []*TestChild {
	// Required children and their required descendants belong in the minimum configuration.
	var required []*TestChild
	for _, child := range children {
		childClass, ok := b.dataStore.Classes[child.Class.String()]
		if !ok || !childClass.RequiredAsChild {
			continue
		}
		for index := range child.Instances {
			child.Instances[index].Children = b.minimumChildren(child.Instances[index].Children)
		}
		required = append(required, child)
	}
	return required
}

func (b *resourceTestBuilder) children(current *Class) error {
	if len(current.TestChildren) == 0 {
		return nil
	}
	full := current.resourceTestCopy()
	cleared := current.resourceTestCopy()
	b.changeChildren(&cleared, "clear")
	if _, err := b.change("ClearChildren", "Explicitly Clear Deletable Children", &cleared, "Update"); err != nil {
		return err
	}
	if _, err := b.change("AddChildren", "Add Children and Verify Descendants", &full, "Update"); err != nil {
		return err
	}
	updated := full.resourceTestCopy()
	b.changeChildren(&updated, "update")
	if _, err := b.change("UpdateChildren", "Update Child Attributes Without Replacing the Parent", &updated, "Update"); err != nil {
		return err
	}
	trimmed := updated.resourceTestCopy()
	b.changeChildren(&trimmed, "trim")
	if _, err := b.change("RemoveChild", "Remove One Child Instance", &trimmed, "Update"); err != nil {
		return err
	}
	reidentified := trimmed.resourceTestCopy()
	b.changeChildren(&reidentified, "identifier")
	if _, err := b.change("ChangeChildIdentifier", "Change Child Identifiers Without Replacing the Parent", &reidentified, "Update"); err != nil {
		return err
	}
	if err := b.add("RestoreChildren", "Restore Full Child Configuration", &full, &full, ""); err != nil {
		return err
	}
	reordered := full.resourceTestCopy()
	b.changeChildren(&reordered, "reorder")
	if err := b.unchanged("ReorderChildren", "Reorder Child Sets (Expect Empty Plan)", &reordered); err != nil {
		return err
	}
	omitted := full.resourceTestCopy()
	omitted.TestChildren = b.minimumChildren(omitted.TestChildren)
	if err := b.add("OmitChildren", "Omit Children and Verify Retained State", &omitted, &full, ""); err != nil {
		return err
	}
	if err := b.rejectChildRemovals(&full); err != nil {
		return err
	}
	*current = full
	return nil
}

func (b *resourceTestBuilder) rejectChildRemovals(full *Class) error {
	// Remove one object at a time, including descendants, then verify the failed update left the full configuration intact.
	invalid := full.resourceTestCopy()
	var visit func(*[]*TestChild, string, []string) error
	visit = func(children *[]*TestChild, path string, inherited []string) error {
		for childIndex, child := range *children {
			childClass, ok := b.dataStore.Classes[child.Class.String()]
			if !ok || slices.Contains(childClass.TestConfig.IgnoreTests, ChildIgnoreTest) || len(child.Instances) == 0 {
				continue
			}
			versions := testAssertionVersions(inherited, b.class.SupportedVersions, childClass.SupportedVersions)
			if childClass.RequiredAsChild || !childClass.AllowDelete {
				name := "RejectDelete" + path + childClass.Name.Capitalized()
				instances := child.Instances
				originalChildren := *children
				child.Instances = nil
				description := "Reject Removal of Non-deletable " + childClass.Documentation.Label
				if childClass.RequiredAsChild {
					// Omit the required argument so schema validation catches removal before payload construction.
					*children = slices.Concat(originalChildren[:childIndex], originalChildren[childIndex+1:])
					description = "Reject Omission of Required " + childClass.Documentation.Label
				}
				err := b.add(name, description, &invalid, full, "")
				*children = originalChildren
				child.Instances = instances
				if err != nil {
					return err
				}
				b.last().ExpectError = fmt.Sprintf("%s object defined by %s cannot be deleted", childClass.Name.Capitalized(), childClass.ResourceNameNested)
				if childClass.RequiredAsChild {
					b.last().ExpectError = regexp.QuoteMeta(fmt.Sprintf("The argument %q is required", childClass.ResourceNameNested))
				}
				if len(versions) > 0 {
					b.last().ChangeVersions = [][]string{versions}
				}
				if err := b.unchanged(name+"Unchanged", "Verify Rejected Removal Left the Object Unchanged", full); err != nil {
					return err
				}
			}
			for index := range child.Instances {
				if err := visit(&child.Instances[index].Children, fmt.Sprintf("%s%s%d", path, childClass.Name.Capitalized(), index+1), versions); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return visit(&invalid.TestChildren, "", nil)
}

func (b *resourceTestBuilder) changeChildren(class *Class, operation string) {
	var change func([]*TestChild, []*TestChild)
	change = func(children, canonical []*TestChild) {
		for _, child := range children {
			childClass, ok := b.dataStore.Classes[child.Class.String()]
			if !ok || slices.Contains(childClass.TestConfig.IgnoreTests, ChildIgnoreTest) {
				continue
			}
			original := findResourceTestChild(canonical, child.Class)
			if operation == "clear" && childClass.AllowDelete && !childClass.RequiredAsChild {
				child.Instances = nil
				continue
			}
			if operation == "trim" && childClass.AllowDelete && len(child.Instances) > 1 {
				child.Instances = child.Instances[:1]
			}
			if operation == "reorder" && !childClass.IsSingleNestedWhenDefinedAsChild {
				slices.Reverse(child.Instances)
			}
			for index := range child.Instances {
				instance := &child.Instances[index]
				if index == 0 && operation == "update" {
					for _, name := range childClass.PropertiesAll {
						property := childClass.Properties[name]
						if resourceMutableProperty(&childClass, property) && property.TestValues.Update.Defined {
							instance.Properties[property.AttributeName] = class.resourceChildUpdateScenario(&childClass, property)
						}
					}
				}
				if operation == "identifier" && childClass.AllowDelete && !childClass.IsSingleNestedWhenDefinedAsChild {
					// Preserve the replacement's canonical descendants so their identifiers remain distinct.
					if original != nil && len(original.Instances) > 1 {
						*instance = original.clone().Instances[1]
						continue
					}
				}
				if original != nil && index < len(original.Instances) {
					change(instance.Children, original.Instances[index].Children)
				}
			}
		}
	}
	change(class.TestChildren, b.class.TestChildren)
}

func findResourceTestChild(children []*TestChild, name *ClassName) *TestChild {
	for _, child := range children {
		if child.Class.Equal(name) {
			return child
		}
	}
	return nil
}

func resourceTestChanges(before, after *Class, ds *DataStore) [][]string {
	var groups [][]string
	add := func(versions []string) {
		if !slices.ContainsFunc(groups, func(group []string) bool { return slices.Equal(group, versions) }) {
			groups = append(groups, versions)
		}
	}
	for _, name := range before.PropertiesAll {
		property := before.Properties[name]
		other := after.Properties[name]
		if property == nil || property.ReadOnly || property.IgnoreInTest || property.TestValues == nil || other.TestValues == nil {
			continue
		}
		if !reflect.DeepEqual(property.TestValues.Create, other.TestValues.Create) {
			add(testAssertionVersions(nil, before.SupportedVersions, property.SupportedVersions))
		}
	}
	var compareChildren func([]*TestChild, []*TestChild, []string)
	compareChildren = func(previous, next []*TestChild, inherited []string) {
		for _, child := range next {
			childClass, ok := ds.Classes[child.Class.String()]
			if !ok || slices.Contains(childClass.TestConfig.IgnoreTests, ChildIgnoreTest) {
				continue
			}
			old := findResourceTestChild(previous, child.Class)
			versions := testAssertionVersions(inherited, before.SupportedVersions, childClass.SupportedVersions)
			if old == nil || len(old.Instances) != len(child.Instances) {
				add(versions)
				continue
			}
			for index, instance := range child.Instances {
				for _, name := range childClass.PropertiesAll {
					property := childClass.Properties[name]
					if property == nil || property.ReadOnly || property.IgnoreInTest {
						continue
					}
					if !reflect.DeepEqual(old.Instances[index].Properties[property.AttributeName], instance.Properties[property.AttributeName]) {
						add(testAssertionVersions(versions, childClass.SupportedVersions, property.SupportedVersions))
					}
				}
				compareChildren(old.Instances[index].Children, instance.Children, versions)
			}
		}
	}
	compareChildren(before.TestChildren, after.TestChildren, nil)
	return groups
}

func (c *Class) prepareResourceTest(ds *DataStore) Class {
	class := c.resourceTestCopy()
	var collect func([]*TestChild)
	collect = func(children []*TestChild) {
		for _, child := range children {
			childClass, ok := ds.Classes[child.Class.String()]
			if !ok || slices.Contains(childClass.TestConfig.IgnoreTests, ChildIgnoreTest) {
				continue
			}
			for _, name := range childClass.PropertiesAll {
				if property := childClass.Properties[name]; resourceMutableProperty(&childClass, property) {
					class.resourceChildUpdateScenario(&childClass, property)
				}
			}
			for _, instance := range child.Instances {
				collect(instance.Children)
			}
		}
	}
	collect(class.TestChildren)
	if parent := class.ParentDnProperty(); parent != nil && parent.TestValues != nil {
		replacement := class.resourceTestCopy()
		replacement.Properties["parentDn"].TestValues.Create = parent.TestValues.ForceNew.clone()
		builder := resourceTestBuilder{class: &class, dataStore: ds}
		builder.compatibleParentTargets(&replacement)
		// Prepare replacement targets once so all scenarios share the same dependency configuration.
		class.TestDependencies = replacement.TestDependencies
	}
	return class
}

// MigrationTestCase supplies saved-state values and current-state expectations derived from normalized definitions.
type MigrationTestCase struct {
	Name     string
	Version  int
	State    string
	Expected string
}

func (c *Class) MigrationTestCases(ds *DataStore) ([]MigrationTestCase, error) {
	var cases []MigrationTestCase
	for _, schema := range c.MigrationSchemas {
		for _, variant := range []string{"values", "null_relation", "empty_relation", "missing_optional_fields"} {
			prior := map[string]any{}
			expected := migrationNullState(c, ds)
			expected["id"] = nil
			if parent := c.ParentDnProperty(); parent != nil {
				expected[parent.AttributeName] = nil
			}
			for _, alias := range c.LegacyAttributes {
				if alias.Status != Removed {
					expected[alias.AttributeName] = nil
				}
			}
			for _, attribute := range schema.Attributes {
				if attribute.IsChildTypeChange() {
					child := ds.Classes[attribute.ChildName.String()]
					childState := migrationNullState(&child, ds)
					for _, property := range child.ModelProperties() {
						if property.ValueType == String && (variant != "missing_optional_fields" || property.Required) {
							childState[property.AttributeName] = attribute.AttributeName + "_" + property.AttributeName + "_value_1"
						}
					}
					prior[attribute.AttributeName] = []any{childState}
					expected[attribute.AttributeName] = childState
					if variant == "null_relation" || variant == "empty_relation" {
						prior[attribute.AttributeName] = nil
						if variant == "empty_relation" {
							prior[attribute.AttributeName] = []any{}
						}
						expected[attribute.AttributeName] = nil
					}
					continue
				}
				if variant == "missing_optional_fields" && attribute.SchemaExpression != "" && !attribute.IsAlias() && !attribute.Required && attribute.AttributeName != "id" {
					continue
				}
				// Distinct markers reveal assignments to the wrong property or relationship.
				marker := attribute.AttributeName + "_value_1"
				var value any = marker
				if attribute.InheritedModelType == "types.Set" || attribute.InheritedModelType == "types.Object" {
					value = nil
				}
				if attribute.Type == SetAttribute || attribute.Type == ListAttribute {
					value = []string{marker, attribute.AttributeName + "_value_2"}
				}
				if attribute.IsListToSet() {
					value = []string{attribute.AttributeName + "_value_2", marker}
				}
				if attribute.AttributeName == "id" {
					value = "uni/tn-example"
				}
				if attribute.Named {
					value = "uni/tn-example/policy-" + marker
				}
				if attribute.Repeated {
					value = []string{"uni/tn-example/policy-" + marker, "uni/tn-example/policy-" + attribute.AttributeName + "_value_2"}
				}
				if attribute.Block {
					// Two objects expose dropped children and accidental overwriting of an earlier member.
					var instances []any
					for index := range 2 {
						members := map[string]any{}
						for _, member := range attribute.Attributes {
							memberMarker := fmt.Sprintf("%s_%s_value_%d", attribute.AttributeName, member.AttributeName, index+1)
							members[member.AttributeName] = memberMarker
							if member.Named {
								members[member.AttributeName] = "uni/tn-example/policy-" + memberMarker
							}
						}
						instances = append(instances, members)
					}
					value = instances
				}
				if attribute.ChildName != nil || attribute.IsListToSet() {
					if variant == "null_relation" {
						value = nil
					}
					if variant == "empty_relation" {
						if attribute.Repeated || attribute.Block || attribute.IsListToSet() {
							value = []any{}
						} else {
							value = ""
						}
					}
				}
				prior[attribute.AttributeName] = value
				if attribute.IsAlias() && attribute.Status != Removed {
					expected[attribute.AttributeName] = value
				}
				if attribute.IsListToSet() {
					if values, ok := value.([]string); ok {
						ordered := slices.Clone(values)
						slices.Sort(ordered)
						value = ordered
						expected[attribute.AttributeName] = value
					}
					expected[attribute.ChildAttribute] = value
					continue
				}
				if attribute.ChildName != nil {
					child := ds.Classes[attribute.ChildName.String()]
					if attribute.Repeated {
						children := []any{}
						if variant != "null_relation" && variant != "empty_relation" {
							if attribute.Block {
								for index := range 2 {
									childState := migrationNullState(&child, ds)
									for _, member := range attribute.Attributes {
										childState[member.ChildAttribute] = fmt.Sprintf("%s_%s_value_%d", attribute.AttributeName, member.AttributeName, index+1)
									}
									children = append(children, childState)
								}
							} else {
								for index, dn := range value.([]string) {
									childState := migrationNullState(&child, ds)
									var target any = dn
									if attribute.Named {
										target = fmt.Sprintf("%s_value_%d", attribute.AttributeName, index+1)
									}
									childState[attribute.ChildAttribute] = target
									children = append(children, childState)
								}
							}
						}
						expected[child.ResourceNameNested] = children
					} else {
						childState := migrationNullState(&child, ds)
						if attribute.Named {
							// Published SDKv2 upgraders converted an unset singleton relation to an empty name.
							if value == nil || value == "" {
								value = ""
							} else {
								value = marker
							}
						}
						childState[attribute.ChildAttribute] = value
						expected[child.ResourceNameNested] = childState
					}
				} else if attribute.TargetField != "" {
					expected[attribute.ChildAttribute] = value
				}
			}
			inputJSON, err := json.MarshalIndent(prior, "", "  ")
			if err != nil {
				return nil, err
			}
			expectedJSON, err := json.MarshalIndent(expected, "", "  ")
			if err != nil {
				return nil, err
			}
			cases = append(cases, MigrationTestCase{Name: fmt.Sprintf("v%d_%s", schema.PriorSchemaVersion, variant), Version: schema.PriorSchemaVersion, State: string(inputJSON), Expected: string(expectedJSON)})
		}
	}
	return cases, nil
}

func migrationNullState(c *Class, ds *DataStore) map[string]any {
	state := map[string]any{}
	for _, property := range c.ModelProperties() {
		state[property.AttributeName] = nil
	}
	for _, name := range c.Children {
		state[ds.Classes[name.String()].ResourceNameNested] = nil
	}

	return state
}

func (b *resourceTestBuilder) legacy(current *Class) error {
	if !b.class.HasCurrentLegacyAttributes() {
		return nil
	}
	base := *b.last()
	parse := func() (*hclwrite.File, error) {
		file, diagnostics := hclwrite.ParseConfig([]byte(base.Configuration), "legacy.tf", hcl.InitialPos)
		if diagnostics.HasErrors() {
			return nil, fmt.Errorf("legacy test configuration: %s", diagnostics)
		}
		return file, nil
	}
	legacyFile, err := parse()
	if err != nil {
		return err
	}
	legacyBody := legacyFile.Body().Blocks()[0].Body()
	conflictFile, err := parse()
	if err != nil {
		return err
	}
	checks := slices.Clone(base.Assertions)
	hasDeprecatedConfiguration := false
	renderer := newTestConfigurationRenderer(current, b.dataStore)
	renderer.resourceTest = true
	for _, alias := range b.class.LegacyAttributes {
		if alias.Status != LegacyDeprecated {
			continue
		}
		if alias.ChildName == nil {
			if attribute := legacyBody.GetAttribute(alias.ChildAttribute); attribute != nil {
				tokens := attribute.Expr().BuildTokens(nil)
				legacyBody.RemoveAttribute(alias.ChildAttribute)
				legacyBody.SetAttributeRaw(alias.AttributeName, tokens)
				conflictFile.Body().Blocks()[0].Body().SetAttributeRaw(alias.AttributeName, tokens)
				hasDeprecatedConfiguration = true
				if alias.IsListToSet() {
					for _, property := range current.ModelProperties() {
						if property.AttributeName != alias.ChildAttribute || property.TestValues == nil {
							continue
						}
						entries := property.TestValues.Create.includedEntries()
						versions := testAssertionVersions(nil, current.SupportedVersions, property.SupportedVersions)
						checks = append(checks, TestAssertion{assertionType: testAssertionSetCount, AttributePath: alias.AttributeName, ExpectedCount: len(entries), Versions: versions})
						for index, entry := range entries {
							checks = append(checks, testAttributeAssertion(fmt.Sprintf("%s.%d", alias.AttributeName, index), entry, versions))
						}
					}
				} else {
					checks = append(checks, TestAssertion{assertionType: testAssertionAttributePair, AttributePath: alias.AttributeName, PairAddress: renderer.subjectResourceAddress().String(), PairAttribute: alias.ChildAttribute})
				}
			}
			continue
		}
		child := b.dataStore.Classes[alias.ChildName.String()]
		// Convert only children configured by the current fixture; excluded children have no rendered dependencies.
		if legacyBody.GetAttribute(child.ResourceNameNested) == nil {
			continue
		}
		fragment, assertions, err := b.legacyChildConfiguration(current, alias, renderer)
		if err != nil {
			return err
		}
		if fragment == "" {
			continue
		}
		parsed, diagnostics := hclwrite.ParseConfig([]byte(fragment), "legacy-values.tf", hcl.InitialPos)
		if diagnostics.HasErrors() {
			return fmt.Errorf("legacy value: %s", diagnostics)
		}
		legacyBody.RemoveAttribute(child.ResourceNameNested)
		legacyBody.AppendUnstructuredTokens(parsed.Body().BuildTokens(nil))
		conflictFile.Body().Blocks()[0].Body().AppendUnstructuredTokens(parsed.Body().BuildTokens(nil))
		hasDeprecatedConfiguration = true
		checks = append(checks, renderer.renderCountedDependencyAssertions(assertions)...)
	}
	add := func(name, description string, file *hclwrite.File, assertions []TestAssertion) {
		b.scenarios = append(b.scenarios, ResourceTestScenario{Name: name, Description: description, Configuration: string(hclwrite.Format(file.Bytes())), Assertions: assertions})
	}
	if hasDeprecatedConfiguration {
		// A conflict is meaningful only when the fixture actually uses a deprecated input.
		add("Legacy", "Configure Through Legacy Attributes", legacyFile, checks)
		add("LegacyUnchanged", "Repeat Legacy Configuration (Expect Empty Plan)", legacyFile, checks)
		b.last().PlanOnly = true
		b.scenarios = append(b.scenarios, ResourceTestScenario{Name: "CurrentAfterLegacy", Description: "Switch to Current Attributes Without Replacement", Configuration: base.Configuration, Assertions: base.Assertions, PlanOnly: true})
		add("LegacyConflict", "Reject Simultaneous Legacy and Current Inputs", conflictFile, nil)
		b.last().ExpectError = "Invalid Attribute Combination|Conflicting legacy configuration"
	}
	frozen, err := parse()
	if err != nil {
		return err
	}
	frozenChecks := slices.Clone(base.Assertions)
	clearedChecks := slices.Clone(base.Assertions)
	hasFrozen := false
	for _, alias := range b.class.LegacyAttributes {
		if alias.Status != Frozen {
			continue
		}
		hasFrozen = true
		if alias.Type == StringAttribute {
			frozen.Body().Blocks()[0].Body().SetAttributeValue(alias.AttributeName, cty.StringVal("example"))
			frozenChecks = append(frozenChecks, TestAssertion{assertionType: testAssertionAttributeValue, AttributePath: alias.AttributeName, ExpectedValue: "example"})
			clearedChecks = append(clearedChecks, TestAssertion{assertionType: testAssertionAttributeValue, AttributePath: alias.AttributeName, ExpectedValue: ""})
		} else {
			frozen.Body().Blocks()[0].Body().SetAttributeValue(alias.AttributeName, cty.SetVal([]cty.Value{cty.StringVal("example")}))
			frozenChecks = append(frozenChecks, TestAssertion{assertionType: testAssertionSetCount, AttributePath: alias.AttributeName, ExpectedCount: 1})
			clearedChecks = append(clearedChecks, TestAssertion{assertionType: testAssertionSetCount, AttributePath: alias.AttributeName, ExpectedCount: 0})
		}
	}
	if hasFrozen {
		add("Frozen", "Retain Frozen Attribute Configuration", frozen, frozenChecks)
		add("FrozenUnchanged", "Repeat Frozen Configuration (Expect Empty Plan)", frozen, frozenChecks)
		b.last().PlanOnly = true
		for _, alias := range b.class.LegacyAttributes {
			if alias.Status == Frozen {
				if alias.Type == StringAttribute {
					frozen.Body().Blocks()[0].Body().SetAttributeValue(alias.AttributeName, cty.StringVal(""))
				} else {
					frozen.Body().Blocks()[0].Body().SetAttributeValue(alias.AttributeName, cty.SetValEmpty(cty.String))
				}
			}
		}
		add("FrozenCleared", "Clear Frozen Attribute State", frozen, clearedChecks)
	}
	b.state = current
	return nil
}

func (b *resourceTestBuilder) legacyChildConfiguration(current *Class, alias MigrationAttribute, renderer testConfigurationRenderer) (string, []TestAssertion, error) {
	// Reuse the full fixture's normalized references and values; only its legacy shape changes.
	childClass := b.dataStore.Classes[alias.ChildName.String()]
	var fixture *TestChild
	for _, child := range current.TestChildren {
		if child.Class.Equal(alias.ChildName) {
			fixture = child
			break
		}
	}
	if fixture == nil || len(fixture.Instances) == 0 {
		return "", nil, nil
	}
	versions := testAssertionVersions(nil, current.SupportedVersions, childClass.SupportedVersions)
	var contents strings.Builder
	var assertions []TestAssertion
	if alias.Block {
		for _, instance := range fixture.Instances {
			fmt.Fprintf(&contents, "%s {\n", alias.AttributeName)
			for _, member := range alias.Attributes {
				scenario := renderer.resourceReferenceScenario(&childClass, instance.Properties[member.ChildAttribute])
				if member.Named {
					for index, entry := range scenario.Entries {
						if entry.ValueType == ReferenceValue {
							address, ok := parseTerraformTestReference(entry.ConfigValue)
							if !ok {
								return "", nil, fmt.Errorf("class %s: invalid legacy target reference %q", current.Name, entry.ConfigValue)
							}
							scenario.Entries[index].ConfigValue = address.withAttribute("id")
						}
					}
				}
				writeTestAttribute(&contents, member.AttributeName, String, scenario, "")
			}
			contents.WriteString("}\n")
		}
		assertions = append(assertions, TestAssertion{assertionType: testAssertionSetCount, AttributePath: alias.AttributeName, ExpectedCount: len(fixture.Instances), Versions: versions})
		return contents.String(), assertions, nil
	}
	var property *Property
	for _, candidate := range childClass.ModelProperties() {
		if candidate.AttributeName == alias.ChildAttribute {
			property = candidate
			break
		}
	}
	scenario := TestValueScenario{Defined: true}
	for _, instance := range fixture.Instances {
		entry, exists := instance.Properties[alias.ChildAttribute].firstIncludedEntry()
		if !exists {
			continue
		}
		if alias.Named {
			if entry.ValueType == ReferenceValue {
				address, ok := parseTerraformTestReference(entry.ConfigValue)
				if !ok {
					return "", nil, fmt.Errorf("class %s: invalid legacy target reference %q", current.Name, entry.ConfigValue)
				}
				entry.ConfigValue = address.withAttribute("id")
			} else if property != nil && property.TestValues != nil {
				if legacy, ok := property.TestValues.Legacy.firstIncludedEntry(); ok {
					entry = legacy
				}
			}
		}
		assertion := TestAssertion{assertionType: testAssertionAttributeValue, AttributePath: alias.AttributeName, ExpectedValue: entry.AssertValue, Versions: versions}
		if entry.ValueType == StringValue {
			assertion.ExpectedValue = entry.ConfigValue
		}
		if alias.Repeated {
			assertion.assertionType = testAssertionSetElementValue
			assertion.AttributePath += ".*"
		}
		if entry.ValueType == ReferenceValue {
			address, _ := parseTerraformTestReference(entry.ConfigValue)
			assertion.PairAddress, assertion.PairAttribute = address.String(), address.Attribute
			if alias.Repeated {
				assertion.assertionType = testAssertionSetElementPair
			} else {
				assertion.assertionType = testAssertionAttributePair
			}
		}
		if entry.ValueType != ExpressionValue {
			assertions = append(assertions, assertion)
		}
		scenario.Entries = append(scenario.Entries, entry)
	}
	valueType := String
	if alias.Repeated {
		valueType = Set
		assertions = append(assertions, TestAssertion{assertionType: testAssertionSetCount, AttributePath: alias.AttributeName, ExpectedCount: len(scenario.Entries), Versions: versions})
	}
	scenario = renderer.resourceReferenceScenario(&childClass, scenario)
	writeTestAttribute(&contents, alias.AttributeName, valueType, scenario, testConfigurationVersions(current.SupportedVersions, childClass.SupportedVersions))
	return contents.String(), assertions, nil
}
