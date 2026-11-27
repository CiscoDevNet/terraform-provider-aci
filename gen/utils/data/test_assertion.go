package data

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
)

// TestAssertionValue describes one expected value within a nested set element.
type TestAssertionValue struct {
	AttributePath string
	ExpectedValue string
	PairAddress   string
	PairAttribute string
}

// TestAssertionParent identifies an enclosing set element for a descendant assertion.
type TestAssertionParent struct {
	AttributePath string
	Values        []TestAssertionValue
}

type testAssertionType int

const (
	testAssertionAttributeSet testAssertionType = iota
	testAssertionAttributeValue
	testAssertionAttributePair
	testAssertionSetCount
	testAssertionSetElementValue
	testAssertionSetElementPair
	testAssertionNestedSetElement
	testAssertionAttributeAbsent
	testAssertionAttributeNotNull
	testAssertionSetPair
)

// TestAssertion is one ordered state assertion for a generated acceptance test.
type TestAssertion struct {
	assertionType      testAssertionType
	AttributePath      string
	ExpectedValue      string
	ExpectedExpression string
	SemanticStringType string
	ExpectedCount      int
	PairAddress        string
	PairAttribute      string
	PairFormat         string
	NestedValues       []TestAssertionValue
	Parents            []TestAssertionParent
	Versions           []*Versions
}

func (a TestAssertion) IsAttributeSet() bool {
	return a.assertionType == testAssertionAttributeSet
}

func (a TestAssertion) IsAttributeValue() bool {
	return a.assertionType == testAssertionAttributeValue
}

func (a TestAssertion) IsAttributePair() bool {
	return a.assertionType == testAssertionAttributePair
}

func (a TestAssertion) IsSetPair() bool {
	return a.assertionType == testAssertionSetPair
}

func (a TestAssertion) IsSetCount() bool {
	return a.assertionType == testAssertionSetCount
}

func (a TestAssertion) IsSetElementValue() bool {
	return a.assertionType == testAssertionSetElementValue
}

func (a TestAssertion) IsSetElementPair() bool {
	return a.assertionType == testAssertionSetElementPair
}

func (a TestAssertion) IsNestedSetElement() bool {
	return a.assertionType == testAssertionNestedSetElement
}

func (a TestAssertion) IsAttributeAbsent() bool {
	return a.assertionType == testAssertionAttributeAbsent
}

func (a TestAssertion) IsAttributeNotNull() bool {
	return a.assertionType == testAssertionAttributeNotNull
}

// DataSourceTestAssertions returns deterministic checks for the data source's
// direct properties and configured child hierarchy.
func (c *Class) DataSourceTestAssertions(ds *DataStore) []TestAssertion {
	renderer := newTestConfigurationRenderer(c, ds)
	assertions := make([]TestAssertion, 0, len(c.PropertiesAll))
	for _, propertyName := range c.PropertiesAll {
		property := c.Properties[propertyName]
		if property == nil || property.IgnoreInTest || property.IgnoreInDataSourceTest || property.Sensitive {
			continue
		}
		runtimeValue := false
		if property.TestValues != nil && property.ValueType != Set && property.ValueType != Object {
			if entry, included := property.TestValues.Create.firstIncludedEntry(); included {
				rendered := renderer.testPropertyScenario(c, property, newTestValueScenario([]TestValueEntry{entry})).Entries[0]
				runtimeValue = entry.AssertStrategy == AssertValueStrategy && rendered.ValueType == ExpressionValue && (entry.ValueType != ExpressionValue || entry.AssertValue == entry.ConfigValue)
			}
		}
		if c.HasResourceArtifact() && runtimeValue {
			assertions = append(assertions, TestAssertion{
				assertionType: testAssertionAttributePair,
				AttributePath: property.AttributeName,
				PairAddress:   renderer.subjectResourceAddress().String(),
				PairAttribute: property.AttributeName,
				Versions:      testAssertionVersions(nil, c.SupportedVersions, property.SupportedVersions),
			})
			continue
		}
		assertions = append(assertions, testPropertyAssertions(c.SupportedVersions, property.AttributeName, property, property.TestValues, nil)...)
	}
	assertions = append(assertions, testChildAssertions(ds, c, c.TestChildren, "", nil, nil, DatasourceArtifact)...)
	for _, legacy := range c.testLegacyAttributes() {
		if legacy.Status == Removed {
			continue
		}
		if legacy.Sensitive {
			assertions = append(assertions, TestAssertion{assertionType: testAssertionAttributeAbsent, AttributePath: legacy.AttributeName})
			continue
		}
		targetClass := c
		var versions []*Versions
		if legacy.ChildName != nil {
			childClass := ds.Classes[legacy.ChildName.String()]
			targetClass = &childClass
			versions = testAssertionVersions(versions, c.SupportedVersions, childClass.SupportedVersions)
		}
		versions = testAssertionVersions(versions, targetClass.SupportedVersions, targetClass.testPropertySupportedVersions(legacy.ChildAttribute))
		if legacy.Status == Frozen {
			if slices.ContainsFunc(c.NestedLegacyAttributes, func(alias MigrationAttribute) bool { return alias.AttributeName == legacy.AttributeName }) {
				// Class-owned frozen fields have no APIC value or prior datasource state to restore.
				assertions = append(assertions, TestAssertion{assertionType: testAssertionAttributeAbsent, AttributePath: legacy.AttributeName})
			} else if legacy.Type == StringAttribute {
				assertions = append(assertions, TestAssertion{assertionType: testAssertionAttributeValue, AttributePath: legacy.AttributeName, ExpectedValue: ""})
			} else {
				assertions = append(assertions, TestAssertion{assertionType: testAssertionSetCount, AttributePath: legacy.AttributeName, ExpectedCount: 0})
			}
		} else if legacy.IsListToSet() && c.HasResourceArtifact() {
			assertions = append(assertions, TestAssertion{assertionType: testAssertionAttributePair, AttributePath: legacy.AttributeName + ".#", PairAddress: renderer.subjectResourceAddress().String(), PairAttribute: legacy.AttributeName + ".#", Versions: versions})
			for _, property := range c.ModelProperties() {
				if property.AttributeName != legacy.ChildAttribute || property.TestValues == nil {
					continue
				}
				versions := testAssertionVersions(nil, c.SupportedVersions, property.SupportedVersions)
				for index := range property.TestValues.Create.includedEntries() {
					path := fmt.Sprintf("%s.%d", legacy.AttributeName, index)
					assertions = append(assertions, TestAssertion{assertionType: testAssertionAttributePair, AttributePath: path, PairAddress: renderer.subjectResourceAddress().String(), PairAttribute: path, Versions: versions})
				}
			}
		} else if c.HasResourceArtifact() {
			assertion := TestAssertion{assertionType: testAssertionAttributePair, AttributePath: legacy.AttributeName, PairAddress: renderer.subjectResourceAddress().String(), PairAttribute: legacy.AttributeName, Versions: versions}
			if legacy.Block || legacy.Type == SetAttribute {
				assertion.assertionType = testAssertionSetPair
			} else if legacy.Named && legacy.Type == StringAttribute {
				child := findResourceTestChild(c.TestChildren, legacy.ChildName)
				childClass := ds.Classes[legacy.ChildName.String()]
				if child != nil && len(child.Instances) > 0 && !slices.Contains(childClass.TestConfig.IgnoreTests, ChildIgnoreTest) {
					// APIC can resolve named-relation tDn between reads; verify the alias exists even while it is empty.
					assertion = TestAssertion{assertionType: testAssertionAttributeNotNull, AttributePath: legacy.AttributeName, Versions: versions}
				}
			}
			assertions = append(assertions, assertion)
		} else if legacy.ChildName == nil {
			assertions = append(assertions, TestAssertion{assertionType: testAssertionAttributePair, AttributePath: legacy.AttributeName, PairAddress: fmt.Sprintf("data.aci_%s.test", c.ResourceName), PairAttribute: legacy.ChildAttribute, Versions: versions})
		}
	}
	return renderer.renderCountedDependencyAssertions(assertions)
}

func testPropertyAssertions(ownerVersions *Versions, attributePath string, property *Property, values *TestValues, inheritedVersions []*Versions) []TestAssertion {
	versions := testAssertionVersions(inheritedVersions, ownerVersions, property.SupportedVersions)
	if property.ReadOnly || values == nil {
		if property.ValueType == Set {
			attributePath += ".#"
		}
		return []TestAssertion{{assertionType: testAssertionAttributeSet, AttributePath: attributePath, Versions: versions}}
	}

	entries := values.Create.includedEntries()
	if property.ValueType == Set {
		assertions := []TestAssertion{{
			assertionType: testAssertionSetCount,
			AttributePath: attributePath,
			ExpectedCount: len(entries),
			Versions:      versions,
		}}
		for _, entry := range entries {
			assertions = append(assertions, testSetElementAssertion(attributePath+".*", entry, versions))
		}
		return assertions
	}

	assertions := make([]TestAssertion, 0, len(entries))
	for _, entry := range entries {
		assertions = append(assertions, testAttributeAssertion(attributePath, entry, versions))
	}
	if len(assertions) == 0 && property.Computed {
		assertions = append(assertions, TestAssertion{assertionType: testAssertionAttributeSet, AttributePath: attributePath, Versions: versions})
	}
	return assertions
}

func testAttributeAssertion(attributePath string, entry TestValueEntry, versions []*Versions) TestAssertion {
	if entry.AssertStrategy == AssertPresenceOnly {
		return TestAssertion{assertionType: testAssertionAttributeSet, AttributePath: attributePath, ExpectedValue: entry.AssertValue, Versions: versions}
	}
	if entry.ValueType == ReferenceValue {
		if address, ok := parseTerraformTestReference(entry.AssertValue); ok {
			return TestAssertion{
				assertionType: testAssertionAttributePair,
				AttributePath: attributePath,
				PairAddress:   address.String(),
				PairAttribute: address.Attribute,
				Versions:      versions,
			}
		}
	}
	return TestAssertion{assertionType: testAssertionAttributeValue, AttributePath: attributePath, ExpectedValue: entry.AssertValue, Versions: versions}
}

func testSetElementAssertion(attributePath string, entry TestValueEntry, versions []*Versions) TestAssertion {
	if entry.AssertStrategy == AssertPresenceOnly {
		return TestAssertion{assertionType: testAssertionAttributeSet, AttributePath: strings.TrimSuffix(attributePath, ".*") + ".#", Versions: versions}
	}
	if entry.ValueType == ReferenceValue {
		if address, ok := parseTerraformTestReference(entry.AssertValue); ok {
			return TestAssertion{
				assertionType: testAssertionSetElementPair,
				AttributePath: attributePath,
				PairAddress:   address.String(),
				PairAttribute: address.Attribute,
				Versions:      versions,
			}
		}
	}
	return TestAssertion{assertionType: testAssertionSetElementValue, AttributePath: attributePath, ExpectedValue: entry.AssertValue, Versions: versions}
}

type nestedTestAssertionGroup struct {
	Versions []*Versions
	Values   []TestAssertionValue
}

type nestedTestAssertion struct {
	AttributePath string
	Identity      []TestAssertionValue
	Parents       []TestAssertionParent
	Groups        []nestedTestAssertionGroup
}

func (a *nestedTestAssertion) addValue(versions []*Versions, value TestAssertionValue) {
	groupIndex := slices.IndexFunc(a.Groups, func(group nestedTestAssertionGroup) bool {
		return slices.EqualFunc(group.Versions, versions, (*Versions).Equal)
	})
	if groupIndex == -1 {
		a.Groups = append(a.Groups, nestedTestAssertionGroup{Versions: append([]*Versions(nil), versions...)})
		groupIndex = len(a.Groups) - 1
	}
	a.Groups[groupIndex].Values = append(a.Groups[groupIndex].Values, value)
}

func (a nestedTestAssertion) assertions() []TestAssertion {
	assertions := make([]TestAssertion, 0, len(a.Groups))
	for _, group := range a.Groups {
		values := append([]TestAssertionValue(nil), group.Values...)
		for identityIndex := len(a.Identity) - 1; identityIndex >= 0; identityIndex-- {
			identity := a.Identity[identityIndex]
			if !slices.ContainsFunc(values, func(value TestAssertionValue) bool {
				return value.AttributePath == identity.AttributePath
			}) {
				values = append([]TestAssertionValue{identity}, values...)
			}
		}
		if !slices.ContainsFunc(values, func(value TestAssertionValue) bool { return value.ExpectedValue != "" || value.PairAddress != "" }) {
			continue
		}
		assertions = append(assertions, TestAssertion{
			assertionType: testAssertionNestedSetElement,
			AttributePath: a.AttributePath,
			NestedValues:  values,
			Parents:       a.Parents,
			Versions:      group.Versions,
		})
	}
	return assertions
}

func testChildAssertions(ds *DataStore, parentClass *Class, children []*TestChild, pathPrefix string, inheritedVersions []*Versions, parents []TestAssertionParent, artifact ArtifactEnum) []TestAssertion {
	var assertions []TestAssertion
	for _, child := range children {
		if child == nil || len(child.Instances) == 0 && artifact != ResourceArtifact {
			continue
		}
		childClass, ok := ds.Classes[child.Class.String()]
		if !ok || slices.Contains(childClass.TestConfig.IgnoreTests, ChildIgnoreTest) {
			continue
		}
		childPath := joinTestAttributePath(pathPrefix, childClass.ResourceNameNested)
		childVersions := testAssertionVersions(inheritedVersions, parentClass.SupportedVersions, childClass.SupportedVersions)
		if childClass.IsSingleNestedWhenDefinedAsChild {
			if len(child.Instances) == 0 {
				assertions = append(assertions, TestAssertion{assertionType: testAssertionAttributeAbsent, AttributePath: childPath, Parents: parents, Versions: childVersions})
				continue
			}
			assertions = append(assertions, testSingleChildAssertions(ds, &childClass, child.Instances[0], childPath, childVersions, parents, artifact)...)
			continue
		}

		if len(parents) == 0 {
			assertions = append(assertions, TestAssertion{
				assertionType: testAssertionSetCount,
				AttributePath: childPath,
				ExpectedCount: len(child.Instances),
				Versions:      childVersions,
			})
		} else {
			// Check the count on the selected parent, including when it has no children.
			parent := parents[len(parents)-1]
			values := append(slices.Clone(parent.Values), TestAssertionValue{
				AttributePath: strings.TrimPrefix(childPath, parent.AttributePath+".") + ".#",
				ExpectedValue: strconv.Itoa(len(child.Instances)),
			})
			assertions = append(assertions, TestAssertion{
				assertionType: testAssertionNestedSetElement,
				AttributePath: parent.AttributePath,
				NestedValues:  values,
				Parents:       parents[:len(parents)-1],
				Versions:      childVersions,
			})
		}
		for _, instance := range child.Instances {
			assertions = append(assertions, testRepeatedChildInstanceAssertions(ds, &childClass, instance, childPath+".*", childVersions, parents, artifact)...)
		}
	}
	return assertions
}

func testSingleChildAssertions(ds *DataStore, childClass *Class, instance TestChildInstance, childPath string, versions []*Versions, parents []TestAssertionParent, artifact ArtifactEnum) []TestAssertion {
	if len(parents) > 0 {
		return testSingleChildWithinSetAssertions(ds, childClass, instance, childPath, versions, parents, artifact)
	}
	class, instance := childClass.nestedLegacyTestInstance(instance)
	childClass = &class

	var assertions []TestAssertion
	for _, propertyName := range childClass.PropertiesAll {
		property := childClass.Properties[propertyName]
		if property == nil || property.PropertyName == "parentDn" || property.IgnoreInTest || artifact == DatasourceArtifact && (property.IgnoreInDataSourceTest || property.Sensitive) {
			continue
		}
		scenario, exists := instance.Properties[property.AttributeName]
		if !exists && !property.ReadOnly && !property.Computed {
			continue
		}
		values := &TestValues{Create: scenario}
		if !exists {
			values = nil
		}
		assertions = append(assertions, testPropertyAssertions(childClass.SupportedVersions, joinTestAttributePath(childPath, property.AttributeName), property, values, versions)...)
	}
	assertions = append(assertions, testChildAssertions(ds, childClass, instance.Children, childPath, versions, parents, artifact)...)
	return assertions
}

func testSingleChildWithinSetAssertions(ds *DataStore, childClass *Class, instance TestChildInstance, childPath string, versions []*Versions, parents []TestAssertionParent, artifact ArtifactEnum) []TestAssertion {
	class, instance := childClass.nestedLegacyTestInstance(instance)
	childClass = &class
	parent := parents[len(parents)-1]
	relativePrefix := strings.TrimPrefix(childPath, parent.AttributePath+".")
	nested := nestedTestAssertion{AttributePath: parent.AttributePath, Identity: parent.Values, Parents: parents[:len(parents)-1]}
	assertions := nested.collectProperties(childClass, instance, childPath, relativePrefix, versions, artifact)
	assertions = append(nested.assertions(), assertions...)
	assertions = append(assertions, testChildAssertions(ds, childClass, instance.Children, childPath, versions, parents, artifact)...)
	return assertions
}

func testRepeatedChildInstanceAssertions(ds *DataStore, childClass *Class, instance TestChildInstance, childPath string, versions []*Versions, parents []TestAssertionParent, artifact ArtifactEnum) []TestAssertion {
	class, instance := childClass.nestedLegacyTestInstance(instance)
	childClass = &class
	nested := nestedTestAssertion{AttributePath: childPath, Parents: parents}
	assertions := nested.collectProperties(childClass, instance, childPath, "", versions, artifact)
	assertions = append(nested.assertions(), assertions...)
	childParents := append(slices.Clone(parents), TestAssertionParent{AttributePath: childPath, Values: nested.Identity})
	assertions = append(assertions, testChildAssertions(ds, childClass, instance.Children, childPath, versions, childParents, artifact)...)
	return assertions
}

func (c *Class) nestedLegacyTestInstance(instance TestChildInstance) (Class, TestChildInstance) {
	// Project compatibility outputs through the existing property assertions without changing canonical fixtures.
	if len(c.NestedLegacyAttributes) == 0 {
		return *c, instance
	}
	class := *c
	class.PropertiesAll = slices.Clone(c.PropertiesAll)
	class.Properties = maps.Clone(c.Properties)
	instance.Properties = maps.Clone(instance.Properties)
	for _, alias := range c.NestedLegacyAttributes {
		if alias.Status == Removed {
			continue
		}
		property := Property{Computed: true, ValueType: String}
		for _, candidate := range c.ModelProperties() {
			if candidate.AttributeName == alias.ChildAttribute {
				property = *candidate
				if _, configured := instance.Properties[alias.AttributeName]; !configured {
					if scenario, exists := instance.Properties[alias.ChildAttribute]; exists {
						instance.Properties[alias.AttributeName] = scenario.clone()
					}
				}
				break
			}
		}
		if alias.Status == Frozen {
			if _, configured := instance.Properties[alias.AttributeName]; !configured {
				continue
			}
		}
		// An old Terraform name can equal an APIC property key that now has a different Terraform name.
		propertyName := "__legacy_" + alias.AttributeName
		property.PropertyName, property.AttributeName = propertyName, alias.AttributeName
		property.Sensitive = alias.Sensitive
		if alias.Type != StringAttribute {
			property.ValueType = Set
		}
		class.PropertiesAll = append(class.PropertiesAll, propertyName)
		class.Properties[propertyName] = &property
	}
	return class, instance
}

func (a *nestedTestAssertion) collectProperties(childClass *Class, instance TestChildInstance, childPath, relativePrefix string, versions []*Versions, artifact ArtifactEnum) []TestAssertion {
	var assertions []TestAssertion
	for _, propertyName := range childClass.PropertiesAll {
		property := childClass.Properties[propertyName]
		if property == nil || property.PropertyName == "parentDn" || property.IgnoreInTest || artifact == DatasourceArtifact && (property.IgnoreInDataSourceTest || property.Sensitive) || property.ReadOnly {
			continue
		}
		scenario, exists := instance.Properties[property.AttributeName]
		if !exists {
			continue
		}
		entries := scenario.includedEntries()
		propertyVersions := testAssertionVersions(versions, childClass.SupportedVersions, property.SupportedVersions)
		fullPath := joinTestAttributePath(childPath, property.AttributeName)
		relativePath := joinTestAttributePath(relativePrefix, property.AttributeName)
		if property.ValueType == Set {
			// Include empty sets, and bind the count to the same owner as its members.
			a.addValue(propertyVersions, TestAssertionValue{AttributePath: relativePath + ".#", ExpectedValue: strconv.Itoa(len(entries))})
			for _, entry := range entries {
				assertions = append(assertions, testSetElementAssertion(fullPath+".*", entry, propertyVersions))
			}
			continue
		}
		for _, entry := range entries {
			scalar := testAttributeAssertion(relativePath, entry, propertyVersions)
			value := TestAssertionValue{
				AttributePath: relativePath,
				ExpectedValue: scalar.ExpectedValue,
				PairAddress:   scalar.PairAddress,
				PairAttribute: scalar.PairAttribute,
			}
			if scalar.IsAttributeSet() {
				// Check presence separately while retaining this child's identity and ancestor scope.
				scalar.AttributePath = fullPath
				assertions = append(assertions, scalar)
			} else {
				a.addValue(propertyVersions, value)
			}
			// A singleton retains the enclosing set identity; repeated children add their own identifiers.
			if !childClass.IsSingleNestedWhenDefinedAsChild && childClass.IsIdentifyingProperty(property.PropertyName) {
				a.Identity = append(a.Identity, value)
			}
		}
	}
	parents := append(slices.Clone(a.Parents), TestAssertionParent{AttributePath: a.AttributePath, Values: a.Identity})
	for index := range assertions {
		assertions[index].Parents = parents
	}
	return assertions
}

func joinTestAttributePath(prefix, attribute string) string {
	if prefix == "" {
		return attribute
	}
	return prefix + "." + attribute
}

func testAssertionVersions(inherited []*Versions, ownerVersions, valueVersions *Versions) []*Versions {
	versions := append([]*Versions(nil), inherited...)
	if valueVersions == nil || valueVersions.Covers(ownerVersions) || slices.ContainsFunc(versions, valueVersions.Equal) {
		return versions
	}
	return append(versions, valueVersions)
}

func (r testConfigurationRenderer) renderCountedDependencyAssertions(assertions []TestAssertion) []TestAssertion {
	// terraform-plugin-testing addresses the first counted state instance with a `.0` suffix.
	for index := range assertions {
		assertion := &assertions[index]
		renderReference := func(address *string, attribute string) {
			dependency := r.findDependency(*address + "." + attribute)
			if dependency == nil {
				return
			}
			*address += ".0"
			dependencyClass, ok := r.dataStore.Classes[dependency.Class.String()]
			if !ok {
				return
			}
			versions := testConfigurationVersions(r.class.SupportedVersions, dependencyClass.SupportedVersions)
			if versions != nil && !dependencyVersionsCoveredByAssertion(dependencyClass.SupportedVersions, assertion.Versions) {
				assertion.Versions = append(assertion.Versions, versions)
			}
		}
		if assertion.IsAttributePair() || assertion.IsSetElementPair() || assertion.IsSetPair() {
			renderReference(&assertion.PairAddress, assertion.PairAttribute)
		}
		for valueIndex := range assertion.NestedValues {
			value := &assertion.NestedValues[valueIndex]
			if value.PairAddress != "" {
				renderReference(&value.PairAddress, value.PairAttribute)
			}
		}
		// Parent selections are shared while assertions are built. Copy them before
		// adding state indices so one assertion cannot modify another's references.
		assertion.Parents = append([]TestAssertionParent(nil), assertion.Parents...)
		for parentIndex := range assertion.Parents {
			parent := &assertion.Parents[parentIndex]
			parent.Values = slices.Clone(parent.Values)
			for valueIndex := range parent.Values {
				value := &parent.Values[valueIndex]
				if value.PairAddress != "" {
					renderReference(&value.PairAddress, value.PairAttribute)
				}
			}
		}
	}
	return assertions
}

func dependencyVersionsCoveredByAssertion(dependencyVersions *Versions, assertionVersions []*Versions) bool {
	// Assertion gates are combined with AND, so one range covered by the dependency is sufficient.
	return slices.ContainsFunc(assertionVersions, func(versions *Versions) bool {
		return dependencyVersions.Covers(versions)
	})
}

func (a TestAssertion) ExpectedGo() string {
	if a.ExpectedExpression != "" {
		return a.ExpectedExpression
	}
	return fmt.Sprintf("%q", a.ExpectedValue)
}

func (v TestAssertionValue) ExpectedGo() string {
	return fmt.Sprintf("%q", v.ExpectedValue)
}

func (a TestAssertion) UsesFormatting() bool {
	return strings.Contains(a.ExpectedExpression, "fmt.")
}

func (c *Class) resourceTestAssertions(ds *DataStore) []TestAssertion {
	renderer := newTestConfigurationRenderer(c, ds)
	renderer.resourceTest = true
	var assertions []TestAssertion
	for _, name := range c.PropertiesAll {
		property := c.Properties[name]
		if property == nil || property.IgnoreInTest {
			continue
		}
		if property.ReadOnly || property.TestValues == nil {
			assertions = append(assertions, testPropertyAssertions(c.SupportedVersions, property.AttributeName, property, nil, nil)...)
			continue
		}
		scenario := property.TestValues.Create
		if !scenario.Defined && (!property.Computed || property.Sensitive) {
			assertions = append(assertions, TestAssertion{assertionType: testAssertionAttributeAbsent, AttributePath: property.AttributeName, Versions: testAssertionVersions(nil, c.SupportedVersions, property.SupportedVersions)})
			continue
		}
		if property.ValueType == Set {
			assertions = append(assertions, testPropertyAssertions(c.SupportedVersions, property.AttributeName, property, &TestValues{Create: scenario}, nil)...)
			continue
		}
		for _, entry := range scenario.includedEntries() {
			versions := testAssertionVersions(nil, c.SupportedVersions, property.SupportedVersions)
			versions = testAssertionVersions(versions, property.SupportedVersions, entry.Versions)
			assertion := testAttributeAssertion(property.AttributeName, entry, versions)
			if entry.AssertStrategy == AssertPresenceOnly {
				assertions = append(assertions, assertion)
				continue
			}
			if assertion.IsAttributeValue() && property.ValueType == VMMArpLearning {
				assertion.SemanticStringType = property.ModelCustomStringType(c.Name)
			}
			if name == "annotation" && entry.ConfigValue == "" && entry.AssertValue == "" {
				assertion.ExpectedExpression = "globalAnnotation"
			} else if renderer.shouldPrefixProperty(c, property) && entry.ValueType == StringValue {
				assertion.ExpectedExpression = fmt.Sprintf("fmt.Sprintf(\"%%s-%%s\", testPrefix, %q)", entry.AssertValue)
			} else if entry.ValueType == ExpressionValue && entry.AssertValue == entry.ConfigValue {
				// An expression without an explicit expected value has no independent literal assertion.
				assertion.assertionType = testAssertionAttributeSet
			} else if entry.ValueType == StringValue && strings.HasPrefix(entry.ConfigValue, "uni/") {
				if format, reference, ok := renderer.prefixedLiteralDNParts(entry); ok {
					if address, parsed := parseTerraformTestReference(reference); parsed {
						assertion.assertionType = testAssertionAttributePair
						assertion.PairAddress = strings.TrimSuffix(address.String(), "[0]")
						assertion.PairAttribute = address.Attribute
						assertion.PairFormat = format
					}
				}
			}
			assertions = append(assertions, assertion)
		}
	}
	assertions = append(assertions, testChildAssertions(ds, c, c.TestChildren, "", nil, nil, ResourceArtifact)...)
	return renderer.renderCountedDependencyAssertions(assertions)
}

func (c *Class) ResourceTestImportIgnore(ds *DataStore) []string {
	var ignored []string
	var collect func(*Class, string, map[string]bool)
	collect = func(class *Class, prefix string, visiting map[string]bool) {
		if visiting[class.Name.String()] {
			return
		}
		visiting = maps.Clone(visiting)
		visiting[class.Name.String()] = true
		for _, name := range class.PropertiesAll {
			if property := class.Properties[name]; property != nil && property.Sensitive {
				ignored = append(ignored, joinTestAttributePath(prefix, property.AttributeName))
			}
		}
		for _, name := range class.Children {
			if child, ok := ds.Classes[name.String()]; ok {
				collect(&child, joinTestAttributePath(prefix, child.ResourceNameNested), visiting)
			}
		}
	}
	collect(c, "", map[string]bool{})
	for _, legacy := range c.testLegacyAttributes() {
		if legacy.Sensitive && legacy.Status != Removed {
			ignored = append(ignored, legacy.AttributeName)
		}
	}
	return ignored
}

func (ds *DataStore) ResourceTestNonDeletableTypes() []string {
	var names []string
	for _, class := range ds.Classes {
		if class.HasResourceArtifact() && !class.AllowDelete {
			names = append(names, "aci_"+class.ResourceName)
		}
	}
	slices.Sort(names)
	return names
}
