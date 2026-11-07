package data

import (
	"fmt"
	"maps"
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
