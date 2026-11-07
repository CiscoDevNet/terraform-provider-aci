package data

import (
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
)

// TestAssertion is one ordered state assertion for a generated acceptance test.
type TestAssertion struct {
	assertionType testAssertionType
	AttributePath string
	ExpectedValue string
	ExpectedCount int
	PairAddress   string
	PairAttribute string
	NestedValues  []TestAssertionValue
	Parents       []TestAssertionParent
	Versions      []string
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
				runtimeValue = rendered.ValueType == ExpressionValue && (entry.ValueType != ExpressionValue || entry.AssertValue == entry.ConfigValue)
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
	assertions = append(assertions, testChildAssertions(ds, c, c.TestChildren, "", nil, nil)...)
	return renderer.renderCountedDependencyAssertions(assertions)
}

func testPropertyAssertions(ownerVersions *Versions, attributePath string, property *Property, values *TestValues, inheritedVersions []string) []TestAssertion {
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

func testAttributeAssertion(attributePath string, entry TestValueEntry, versions []string) TestAssertion {
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

func testSetElementAssertion(attributePath string, entry TestValueEntry, versions []string) TestAssertion {
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
	Versions []string
	Values   []TestAssertionValue
}

type nestedTestAssertion struct {
	AttributePath string
	Identity      []TestAssertionValue
	Parents       []TestAssertionParent
	Groups        []nestedTestAssertionGroup
}

func (a *nestedTestAssertion) addValue(versions []string, value TestAssertionValue) {
	groupIndex := slices.IndexFunc(a.Groups, func(group nestedTestAssertionGroup) bool {
		return slices.Equal(group.Versions, versions)
	})
	if groupIndex == -1 {
		a.Groups = append(a.Groups, nestedTestAssertionGroup{Versions: append([]string(nil), versions...)})
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

func testChildAssertions(ds *DataStore, parentClass *Class, children []*TestChild, pathPrefix string, inheritedVersions []string, parents []TestAssertionParent) []TestAssertion {
	var assertions []TestAssertion
	for _, child := range children {
		if child == nil || len(child.Instances) == 0 {
			continue
		}
		childClass, ok := ds.Classes[child.Class.String()]
		if !ok || slices.Contains(childClass.TestConfig.IgnoreTests, ChildIgnoreTest) {
			continue
		}
		childPath := joinTestAttributePath(pathPrefix, childClass.ResourceNameNested)
		childVersions := testAssertionVersions(inheritedVersions, parentClass.SupportedVersions, childClass.SupportedVersions)
		if childClass.IsSingleNestedWhenDefinedAsChild {
			assertions = append(assertions, testSingleChildAssertions(ds, &childClass, child.Instances[0], childPath, childVersions, parents)...)
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
			assertions = append(assertions, testRepeatedChildInstanceAssertions(ds, &childClass, instance, childPath+".*", childVersions, parents)...)
		}
	}
	return assertions
}

func testSingleChildAssertions(ds *DataStore, childClass *Class, instance TestChildInstance, childPath string, versions []string, parents []TestAssertionParent) []TestAssertion {
	if len(parents) > 0 {
		return testSingleChildWithinSetAssertions(ds, childClass, instance, childPath, versions, parents)
	}

	var assertions []TestAssertion
	for _, propertyName := range childClass.PropertiesAll {
		property := childClass.Properties[propertyName]
		if property == nil || property.PropertyName == "parentDn" || property.IgnoreInTest || property.IgnoreInDataSourceTest || property.Sensitive {
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
	assertions = append(assertions, testChildAssertions(ds, childClass, instance.Children, childPath, versions, parents)...)
	return assertions
}

func testSingleChildWithinSetAssertions(ds *DataStore, childClass *Class, instance TestChildInstance, childPath string, versions []string, parents []TestAssertionParent) []TestAssertion {
	parent := parents[len(parents)-1]
	relativePrefix := strings.TrimPrefix(childPath, parent.AttributePath+".")
	nested := nestedTestAssertion{AttributePath: parent.AttributePath, Identity: parent.Values, Parents: parents[:len(parents)-1]}
	assertions := nested.collectProperties(childClass, instance, childPath, relativePrefix, versions)
	assertions = append(nested.assertions(), assertions...)
	assertions = append(assertions, testChildAssertions(ds, childClass, instance.Children, childPath, versions, parents)...)
	return assertions
}

func testRepeatedChildInstanceAssertions(ds *DataStore, childClass *Class, instance TestChildInstance, childPath string, versions []string, parents []TestAssertionParent) []TestAssertion {
	nested := nestedTestAssertion{AttributePath: childPath, Parents: parents}
	assertions := nested.collectProperties(childClass, instance, childPath, "", versions)
	assertions = append(nested.assertions(), assertions...)
	childParents := append(slices.Clone(parents), TestAssertionParent{AttributePath: childPath, Values: nested.Identity})
	assertions = append(assertions, testChildAssertions(ds, childClass, instance.Children, childPath, versions, childParents)...)
	return assertions
}

func (a *nestedTestAssertion) collectProperties(childClass *Class, instance TestChildInstance, childPath, relativePrefix string, versions []string) []TestAssertion {
	var assertions []TestAssertion
	for _, propertyName := range childClass.PropertiesAll {
		property := childClass.Properties[propertyName]
		if property == nil || property.PropertyName == "parentDn" || property.IgnoreInTest || property.IgnoreInDataSourceTest || property.Sensitive || property.ReadOnly {
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
			a.addValue(propertyVersions, value)
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

func testAssertionVersions(inherited []string, ownerVersions, valueVersions *Versions) []string {
	versions := append([]string(nil), inherited...)
	if valueVersions == nil || valueVersions.Covers(ownerVersions) || slices.Contains(versions, valueVersions.Raw()) {
		return versions
	}
	return append(versions, valueVersions.Raw())
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
			if versions != "" && !dependencyVersionsCoveredByAssertion(dependencyClass.SupportedVersions, assertion.Versions) {
				assertion.Versions = append(assertion.Versions, versions)
			}
		}
		if assertion.IsAttributePair() || assertion.IsSetElementPair() {
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

func dependencyVersionsCoveredByAssertion(dependencyVersions *Versions, assertionVersions []string) bool {
	// Assertion gates are combined with AND, so one range covered by the dependency is sufficient.
	return slices.ContainsFunc(assertionVersions, func(raw string) bool {
		versions, err := NewVersions(raw)
		return err == nil && dependencyVersions.Covers(versions)
	})
}
