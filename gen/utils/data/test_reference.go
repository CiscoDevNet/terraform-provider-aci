package data

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type testReferenceAllocation struct {
	address  terraformTestAddress
	instance int
}

type testReferenceAllocator struct {
	subjectType string
	allocations map[string]testReferenceAllocation
	counters    map[string]int
	aliases     map[string]terraformTestAddress
}

func newTestReferenceAllocator(class *Class) *testReferenceAllocator {
	return &testReferenceAllocator{
		subjectType: "aci_" + class.ResourceName,
		allocations: make(map[string]testReferenceAllocation),
		counters:    make(map[string]int),
		aliases:     make(map[string]terraformTestAddress),
	}
}

func newSeededTestReferenceAllocator(class *Class) *testReferenceAllocator {
	allocator := newTestReferenceAllocator(class)
	visited := make(map[*TestDependency]bool)
	var seed func([]*TestDependency)
	seed = func(dependencies []*TestDependency) {
		for _, dependency := range dependencies {
			if dependency == nil || visited[dependency] {
				continue
			}
			visited[dependency] = true
			allocator.seed(dependency)
			seed(dependency.Dependencies)
		}
	}
	seed(class.TestDependencies)
	return allocator
}

func (a *testReferenceAllocator) seed(dependency *TestDependency) {
	// Reserve existing assignments before importing more dependencies into the same fixture.
	if dependency.ReferenceType == StaticReference {
		return
	}
	sourceAddress, sourceOK := parseTestDependencyReference(dependency.SourceReference, dependency.ReferenceType)
	canonicalAddress, canonicalOK := parseTerraformTestReference(dependency.Reference)
	if !sourceOK || !canonicalOK {
		return
	}
	canonicalAddress.Attribute = ""
	a.allocations[sourceAddress.String()] = testReferenceAllocation{
		address:  canonicalAddress,
		instance: dependency.TestInstance,
	}
	a.aliases[sourceAddress.String()] = canonicalAddress
	a.aliases[canonicalAddress.String()] = canonicalAddress
	prefix := strings.TrimSuffix(canonicalAddress.Name, fmt.Sprintf("_%d", dependency.TestInstance))
	counterKey := canonicalAddress.BlockType + "." + canonicalAddress.Type + "." + prefix
	a.counters[counterKey] = max(a.counters[counterKey], dependency.TestInstance)
}

func (a *testReferenceAllocator) allocate(dependency *TestDependency) {
	if dependency == nil || dependency.ReferenceType == StaticReference {
		return
	}
	if dependency.SourceReference == "" {
		dependency.SourceReference = dependency.Reference
	}
	sourceAddress, ok := parseTestDependencyReference(dependency.SourceReference, dependency.ReferenceType)
	if !ok {
		return
	}
	sourceKey := sourceAddress.String()
	if allocation, exists := a.allocations[sourceKey]; exists {
		if currentAddress, parsed := parseTestDependencyReference(dependency.Reference, dependency.ReferenceType); parsed {
			a.aliases[currentAddress.String()] = allocation.address
		}
		dependency.TestInstance = allocation.instance
		dependency.Reference = allocation.address.withAttribute(sourceAddress.Attribute)
		a.aliases[allocation.address.String()] = allocation.address
		return
	}

	prefix := "test"
	if sourceAddress.BlockType == "resource" && sourceAddress.Type == a.subjectType {
		switch dependency.Role {
		case Parent:
			prefix = "parent"
		case Target:
			prefix = "target"
		default:
			prefix = "dependency"
		}
	}
	counterKey := sourceAddress.BlockType + "." + sourceAddress.Type + "." + prefix
	a.counters[counterKey]++
	canonicalAddress := terraformTestAddress{
		BlockType: sourceAddress.BlockType,
		Type:      sourceAddress.Type,
		Name:      fmt.Sprintf("%s_%d", prefix, a.counters[counterKey]),
	}
	allocation := testReferenceAllocation{address: canonicalAddress, instance: a.counters[counterKey]}
	a.allocations[sourceKey] = allocation
	a.aliases[sourceKey] = canonicalAddress
	if currentAddress, parsed := parseTestDependencyReference(dependency.Reference, dependency.ReferenceType); parsed {
		a.aliases[currentAddress.String()] = canonicalAddress
	}
	dependency.TestInstance = allocation.instance
	dependency.Reference = canonicalAddress.withAttribute(sourceAddress.Attribute)
}

func (c *Class) canonicalizeTestDependencies() {
	allocator := newTestReferenceAllocator(c)
	visited := make(map[*TestDependency]bool)
	var allocate func([]*TestDependency)
	allocate = func(dependencies []*TestDependency) {
		for _, dependency := range dependencies {
			if dependency == nil || visited[dependency] {
				continue
			}
			visited[dependency] = true
			allocator.allocate(dependency)
			c.setSameResourceDependencyName(dependency)
			allocate(dependency.Dependencies)
		}
	}
	allocate(c.TestDependencies)
	c.rewriteTestReferences(allocator.aliases)
}

func parseTestDependencyReference(reference string, referenceType ReferenceTypeEnum) (terraformTestAddress, bool) {
	if address, ok := parseTerraformTestReference(reference); ok {
		return address, true
	}
	if referenceType != ResourceReference {
		return terraformTestAddress{}, false
	}
	parts := strings.Split(reference, ".")
	if len(parts) < 3 || strings.Contains(parts[0], "/") {
		return terraformTestAddress{}, false
	}
	resourceType := parts[0]
	if !strings.HasPrefix(resourceType, "aci_") {
		resourceType = "aci_" + resourceType
	}
	return terraformTestAddress{
		BlockType: "resource",
		Type:      resourceType,
		Name:      parts[1],
		Attribute: strings.Join(parts[2:], "."),
	}, true
}

func (address terraformTestAddress) withAttribute(attribute string) string {
	if attribute == "" {
		return address.String()
	}
	return address.String() + "." + attribute
}

func rewriteTestReference(reference string, referenceType ReferenceTypeEnum, aliases map[string]terraformTestAddress) string {
	address, ok := parseTestDependencyReference(reference, referenceType)
	if !ok {
		return reference
	}
	canonicalAddress, exists := aliases[address.String()]
	if !exists {
		return reference
	}
	return canonicalAddress.withAttribute(address.Attribute)
}

func (c *Class) rewriteTestReferences(aliases map[string]terraformTestAddress) {
	rewriteTestDependencyReferences(c.TestDependencies, aliases)
	for _, property := range c.Properties {
		property.rewriteTestReferences(aliases)
	}
}

func rewriteTestDependencyReferences(dependencies []*TestDependency, aliases map[string]terraformTestAddress) {
	visited := make(map[*TestDependency]bool)
	var rewriteDependencies func([]*TestDependency)
	rewriteDependencies = func(dependencies []*TestDependency) {
		for _, dependency := range dependencies {
			if dependency == nil || visited[dependency] {
				continue
			}
			visited[dependency] = true
			dependency.rewriteReferences(aliases)
			rewriteDependencies(dependency.Dependencies)
		}
	}
	rewriteDependencies(dependencies)
}

func (dependency *TestDependency) rewriteReferences(aliases map[string]terraformTestAddress) {
	for name, value := range dependency.ConfigOverrides {
		dependency.ConfigOverrides[name] = rewriteTestReference(value, referenceValueType(value), aliases)
	}
	for _, child := range dependency.Children {
		rewriteTestChildReferences([]*TestChild{child}, aliases)
	}
}

func (property *Property) rewriteTestReferences(aliases map[string]terraformTestAddress) {
	if property == nil || property.TestValues == nil {
		return
	}
	property.TestValues.Create.rewriteReferences(aliases)
	property.TestValues.DataSourceMissing.rewriteReferences(aliases)
	property.TestValues.Update.rewriteReferences(aliases)
	property.TestValues.Default.rewriteReferences(aliases)
	property.TestValues.ForceNew.rewriteReferences(aliases)
	property.TestValues.Legacy.rewriteReferences(aliases)
}

func (scenario *TestValueScenario) rewriteReferences(aliases map[string]terraformTestAddress) {
	for index := range scenario.Entries {
		entry := &scenario.Entries[index]
		if entry.ValueType != ReferenceValue {
			continue
		}
		entry.ConfigValue = rewriteTestReference(entry.ConfigValue, ResourceReference, aliases)
		entry.AssertValue = rewriteTestReference(entry.AssertValue, ResourceReference, aliases)
	}
}

func rewriteTestChildReferences(children []*TestChild, aliases map[string]terraformTestAddress) {
	for _, child := range children {
		if child == nil {
			continue
		}
		for instanceIndex := range child.Instances {
			instance := &child.Instances[instanceIndex]
			for name, scenario := range instance.Properties {
				scenario.rewriteReferences(aliases)
				instance.Properties[name] = scenario
			}
			rewriteTestChildReferences(instance.Children, aliases)
		}
	}
}

func referenceValueType(value string) ReferenceTypeEnum {
	if strings.HasPrefix(value, "data.") {
		return DataSourceReference
	}
	return ResourceReference
}

func (dependency *TestDependency) sourceKey() string {
	if dependency == nil {
		return ""
	}
	reference := dependency.SourceReference
	if reference == "" {
		reference = dependency.Reference
	}
	if address, ok := parseTestDependencyReference(reference, dependency.ReferenceType); ok {
		return dependency.ReferenceType.String() + ":" + address.String()
	}
	return dependency.ReferenceType.String() + ":" + dependency.Class.String() + ":" + reference
}

func (source *TestDependency) clone() *TestDependency {
	// Copy this node's configuration; dependency import reconnects its prerequisite graph.
	clone := &TestDependency{
		PrefixEnabled:   source.PrefixEnabled,
		Class:           source.Class,
		SourceReference: source.SourceReference,
		Reference:       source.Reference,
		TestInstance:    source.TestInstance,
		ReferenceType:   source.ReferenceType,
		Role:            source.Role,
		OrderOnly:       source.OrderOnly,
		IgnoreChanges:   append([]string(nil), source.IgnoreChanges...),
		TargetClasses:   append([]*ClassName(nil), source.TargetClasses...),
		nameGenerated:   source.nameGenerated,
	}
	if len(source.ConfigOverrides) > 0 {
		clone.ConfigOverrides = maps.Clone(source.ConfigOverrides)
	}
	if len(source.Children) > 0 {
		clone.Children = make(map[string]*TestChild, len(source.Children))
		for name, child := range source.Children {
			clone.Children[name] = child.clone()
		}
	}
	return clone
}

func (source *TestChild) clone() *TestChild {
	if source == nil {
		return nil
	}
	clone := &TestChild{Class: source.Class, Instances: make([]TestChildInstance, len(source.Instances))}
	for index, sourceInstance := range source.Instances {
		instance := TestChildInstance{Properties: make(map[string]TestValueScenario, len(sourceInstance.Properties))}
		for name, scenario := range sourceInstance.Properties {
			instance.Properties[name] = scenario.clone()
		}
		instance.Children = make([]*TestChild, 0, len(sourceInstance.Children))
		for _, child := range sourceInstance.Children {
			instance.Children = append(instance.Children, child.clone())
		}
		clone.Instances[index] = instance
	}
	return clone
}

type terraformTestAddress struct {
	BlockType string
	Type      string
	Name      string
	Attribute string
}

// String returns the object address used to identify a dependency, without its attribute.
func (a terraformTestAddress) String() string {
	if a.BlockType == "data" {
		return "data." + a.Type + "." + a.Name
	}
	return a.Type + "." + a.Name
}

func parseTerraformTestReference(reference string) (terraformTestAddress, bool) {
	parts := strings.Split(reference, ".")
	if len(parts) >= 4 && parts[0] == "data" {
		return terraformTestAddress{BlockType: "data", Type: parts[1], Name: parts[2], Attribute: strings.Join(parts[3:], ".")}, true
	}
	if len(parts) >= 3 && strings.HasPrefix(parts[0], "aci_") {
		return terraformTestAddress{BlockType: "resource", Type: parts[0], Name: parts[1], Attribute: strings.Join(parts[2:], ".")}, true
	}
	return terraformTestAddress{}, false
}

func findTestDependencyByObjectReference(dependencies []*TestDependency, reference string) *TestDependency {
	address, ok := parseTerraformTestReference(reference)
	if !ok {
		return nil
	}
	matches := func(candidate string, referenceType ReferenceTypeEnum) bool {
		candidateAddress, parsed := parseTestDependencyReference(candidate, referenceType)
		return parsed && candidateAddress.BlockType == address.BlockType && candidateAddress.Type == address.Type && candidateAddress.Name == address.Name
	}
	visited := make(map[*TestDependency]bool)
	var find func([]*TestDependency) *TestDependency
	find = func(dependencies []*TestDependency) *TestDependency {
		for _, dependency := range dependencies {
			if visited[dependency] {
				continue
			}
			visited[dependency] = true
			if matches(dependency.Reference, dependency.ReferenceType) || matches(dependency.SourceReference, dependency.ReferenceType) {
				return dependency
			}
			if found := find(dependency.Dependencies); found != nil {
				return found
			}
		}
		return nil
	}
	return find(dependencies)
}

func (r testConfigurationRenderer) collectResourceChildUpdateReferences(children []*TestChild, references *[]string) {
	for _, child := range children {
		childClass, ok := r.dataStore.Classes[child.Class.String()]
		if !ok || slices.Contains(childClass.TestConfig.IgnoreTests, ChildIgnoreTest) {
			continue
		}
		for _, name := range childClass.PropertiesAll {
			property := childClass.Properties[name]
			if !resourceMutableProperty(&childClass, property) {
				continue
			}
			for _, entry := range r.class.resourceChildUpdateScenario(&childClass, property).includedEntries() {
				if entry.ValueType == ReferenceValue && !slices.Contains(*references, entry.ConfigValue) {
					*references = append(*references, entry.ConfigValue)
				}
			}
		}
		for _, instance := range child.Instances {
			r.collectResourceChildUpdateReferences(instance.Children, references)
		}
	}
}

func (b *resourceTestBuilder) compatibleParentTargets(class *Class) {
	renderer := newTestConfigurationRenderer(class, b.dataStore)
	parent := renderer.primaryParentDependency()
	if parent == nil {
		return
	}
	for _, name := range class.PropertiesAll {
		property := class.Properties[name]
		if property == nil || property.TestValues == nil || name == "parentDn" {
			continue
		}
		entry, ok := property.TestValues.Create.firstIncludedEntry()
		if !ok || entry.ValueType != ReferenceValue {
			continue
		}
		dependency := renderer.findDependency(entry.ConfigValue)
		if dependency == nil || dependency.Role != Target || slices.ContainsFunc(parent.TargetClasses, dependency.Class.Equal) {
			continue
		}
		for _, target := range class.TestDependencies {
			if target.Role == Target && slices.ContainsFunc(parent.TargetClasses, target.Class.Equal) {
				entry = testValueEntryForDependency(target)
				if class.Relation.Type == Named {
					entry.ConfigValue = targetReferenceToName(entry.ConfigValue)
					entry.AssertValue = targetReferenceToName(entry.AssertValue)
				}
				property.TestValues.Create = newTestValueScenario([]TestValueEntry{entry})
				break
			}
		}
	}
	b.reparentTestTargets(class, parent)
}

func (b *resourceTestBuilder) reparentTestTargets(class *Class, parent *TestDependency) {
	renderer := newTestConfigurationRenderer(class, b.dataStore)
	previous := newTestConfigurationRenderer(b.class, b.dataStore).primaryParentDependency()
	if previous == nil || previous.Reference == parent.Reference {
		return
	}
	oldParent, oldOK := parseTerraformTestReference(previous.Reference)
	newParent, newOK := parseTerraformTestReference(parent.Reference)
	if !oldOK || !newOK {
		return
	}
	parentAliases := map[string]terraformTestAddress{oldParent.String(): newParent}
	aliases := make(map[string]terraformTestAddress)
	var importedTargets []*TestDependency
	bind := func(scenario TestValueScenario) {
		for _, entry := range scenario.includedEntries() {
			if entry.ValueType != ReferenceValue {
				continue
			}
			dependency := renderer.findDependency(entry.ConfigValue)
			if dependency == nil || dependency.ReferenceType != ResourceReference {
				continue
			}
			targetClass, loaded := b.dataStore.Classes[dependency.Class.String()]
			var targetSchema *Class
			if loaded {
				targetSchema = &targetClass
			}
			targetParent := renderer.parentDependency(targetSchema, dependency)
			if targetParent == nil || targetParent.Reference != previous.Reference {
				continue
			}
			// Targets sharing the old fixture parent follow its replacement; other placements are preserved.
			address, _ := parseTerraformTestReference(dependency.Reference)
			sourceAddress := address
			sourceAddress.Name = fmt.Sprintf("%s_parent_%s", address.Name, newParent.Name)
			clone := dependency.clone()
			clone.SourceReference = sourceAddress.withAttribute(address.Attribute)
			clone.Reference = clone.SourceReference
			clone.Dependencies = append([]*TestDependency(nil), dependency.Dependencies...)
			clone.Dependencies = slices.DeleteFunc(clone.Dependencies, func(prerequisite *TestDependency) bool {
				return prerequisite.Reference == previous.Reference || prerequisite.Reference == parent.Reference
			})
			clone.Dependencies = append(clone.Dependencies, parent)
			if clone.ConfigOverrides == nil {
				clone.ConfigOverrides = make(map[string]string)
			}
			if loaded {
				if property := targetClass.ParentDnProperty(); property != nil {
					clone.ConfigOverrides[property.AttributeName] = parent.Reference
				}
				resourceScenario := dependencyTestResourceScenario
				if dependency.TestInstance > 1 {
					resourceScenario = alternateDependencyTestResourceScenario
				}
				// New fixture addresses must retain distinct naming values from the original instances.
				for _, name := range targetClass.IdentifiedBy {
					property := targetClass.Properties[name]
					if property == nil || property.TestValues == nil {
						continue
					}
					if _, overridden := testPropertyOverride(property, clone.ConfigOverrides); overridden {
						continue
					}
					if value, ok := testResourcePropertyScenario(&targetClass, property, resourceScenario).firstIncludedEntry(); ok {
						clone.ConfigOverrides[property.AttributeName] = value.ConfigValue
					}
				}
			}
			clone.rewriteReferences(parentAliases)
			imported, _ := class.importTestDependency(clone)
			importedTargets = append(importedTargets, imported)
			importedAddress, _ := parseTerraformTestReference(imported.Reference)
			aliases[address.String()] = importedAddress
		}
	}
	for _, name := range class.PropertiesAll {
		if property := class.Properties[name]; property != nil && name != "parentDn" && property.TestValues != nil {
			bind(property.TestValues.Create)
		}
	}
	var visit func([]*TestChild)
	visit = func(children []*TestChild) {
		for _, child := range children {
			for _, instance := range child.Instances {
				for _, name := range slices.Sorted(maps.Keys(instance.Properties)) {
					bind(instance.Properties[name])
				}
				visit(instance.Children)
			}
		}
	}
	visit(class.TestChildren)
	// A target may refer to another target visited later, such as an L3Out's VRF.
	// Rewrite only replacement nodes after the complete alias map is available.
	for _, dependency := range importedTargets {
		dependency.rewriteReferences(aliases)
		for index, prerequisite := range dependency.Dependencies {
			reference := rewriteTestReference(prerequisite.Reference, prerequisite.ReferenceType, aliases)
			if reference != prerequisite.Reference {
				dependency.Dependencies[index] = renderer.findDependency(reference)
			}
		}
	}
	for _, property := range class.Properties {
		if property.TestValues != nil {
			property.TestValues.Create.rewriteReferences(aliases)
		}
	}
	rewriteTestChildReferences(class.TestChildren, aliases)
}

func (c *Class) resourceChildUpdateScenario(childClass *Class, property *Property) TestValueScenario {
	scenario := property.TestValues.Update.clone()
	for index := range scenario.Entries {
		entry := &scenario.Entries[index]
		if entry.ValueType != ReferenceValue {
			continue
		}
		if source := findTestDependencyByObjectReference(childClass.TestDependencies, entry.ConfigValue); source != nil {
			_, aliases := c.importTestDependency(source)
			entry.ConfigValue = rewriteTestReference(entry.ConfigValue, referenceValueType(entry.ConfigValue), aliases)
			entry.AssertValue = rewriteTestReference(entry.AssertValue, referenceValueType(entry.AssertValue), aliases)
		}
	}
	return scenario
}

func (r testConfigurationRenderer) resourceReferenceScenario(class *Class, scenario TestValueScenario) TestValueScenario {
	scenario = scenario.clone()
	for index := range scenario.Entries {
		entry := &scenario.Entries[index]
		if entry.ValueType != ReferenceValue {
			continue
		}
		dependency := r.findDependency(entry.ConfigValue)
		if dependency == nil {
			continue
		}
		if target, ok := r.dataStore.Classes[dependency.Class.String()]; ok {
			if versions := testConfigurationVersions(class.SupportedVersions, target.SupportedVersions); versions != nil {
				entry.ConfigValue = fmt.Sprintf("%s ? %s : null", testVersionExpression(versions), testValueHCL(*entry))
				entry.ValueType = ExpressionValue
			}
		}
	}
	return scenario
}
