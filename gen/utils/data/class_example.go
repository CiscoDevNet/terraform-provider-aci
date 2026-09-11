package data

import (
	"fmt"
	"slices"
	"strings"
)

func (c *Class) setExampleContext(ds *DataStore) error {
	// Choose one deterministic, resolvable parent placement per class.
	c.ExampleContext = ExampleContext{
		Platform:     c.PlatformType,
		ParentDnMode: ParentDnFromTestValues,
	}
	if len(c.ParentDnVariants) > 0 {
		variant := c.ParentDnVariants[0]
		c.ExampleContext = ExampleContext{
			ParentClass:  variant.ParentClass,
			Platform:     variant.TestPlatform,
			ParentDnMode: ParentDnFromSchemaDefault,
		}
		return nil
	}

	if c.ParentDnProperty() == nil {
		return nil
	}
	if c.ParentDnHasSchemaDefault() {
		c.ExampleContext.ParentDnMode = ParentDnFromSchemaDefault
		return nil
	}
	parentClasses := c.Parents
	if c.ClassDefinition.ExampleParentClass != "" {
		parentClass, err := NewClassName(c.ClassDefinition.ExampleParentClass)
		if err != nil {
			return fmt.Errorf("class '%s': invalid example_parent_class '%s': %w", c.Name, c.ClassDefinition.ExampleParentClass, err)
		}
		parentClasses = []*ClassName{parentClass}
	}
	for _, parentClass := range parentClasses {
		parent, ok := ds.Classes[parentClass.String()]
		if !ok || !parent.HasResourceArtifact() {
			continue
		}
		c.ExampleContext.ParentClass = parentClass
		c.ExampleContext.ParentValue = newExampleReference(ResourceReference, parent.ResourceName, "id")
		c.ExampleContext.ParentDnMode = ParentDnFromReference
		return nil
	}
	for _, dependency := range c.TestDependencies {
		if dependency.Role != Parent || dependency.ReferenceType == StaticReference {
			continue
		}
		c.ExampleContext.ParentClass = dependency.Class
		c.ExampleContext.ParentValue = dependency.exampleValue(ds)
		c.ExampleContext.ParentDnMode = ParentDnFromReference
		return nil
	}

	// Some system-scoped classes have no generated parent resource. Their normalized
	// parent_dn test value supplies the literal value for this single fallback example.
	return nil
}

func (c *Class) propertyExampleValue(ds *DataStore, property *Property, context ExampleContext, variant ExampleVariantEnum) ExampleValue {
	// Full examples vary optional attributes only. Required and identifying values stay
	// aligned with the minimum example so both configurations address the same object.
	if property.ReadOnly {
		return nil
	}
	if value, defined := property.exampleDefinitionValue(ds, c.Name, variant); defined {
		return value
	}

	if property.PropertyName == "parentDn" {
		switch context.ParentDnMode {
		case ParentDnFromTestValues:
			return property.exampleValueFromTestScenario(ds, c.Name, c.propertyExampleScenario(property, variant))
		case ParentDnFromSchemaDefault:
			// Keep parent_dn omitted so Terraform applies the schema default.
			return nil
		case ParentDnFromReference:
			return context.ParentValue
		default:
			panic(fmt.Sprintf("unsupported example parent-DN mode %d", context.ParentDnMode))
		}
	}

	scenario := c.propertyExampleScenario(property, variant)
	return c.resolvePropertyExampleValue(ds, property, context, scenario)
}

func (c *Class) resolvePropertyExampleValue(ds *DataStore, property *Property, context ExampleContext, scenario TestValueScenario) ExampleValue {
	// Resolve semantic references before converting ordinary test values so a
	// reference never has to be reconstructed from its HCL expression.
	if property.PropertyName == "tDn" && c.Relation.RelationalClass && c.Relation.Type == Explicit {
		if dependency := c.exampleTargetDependency(context); dependency != nil {
			return dependency.exampleValue(ds)
		}
		return property.exampleValueFromTestScenario(ds, c.Name, scenario)
	}

	if property.Reference != nil && property.Reference.ExampleClass != nil && property.exampleScenarioIncluded(scenario) {
		if resourceName := ds.getResourceName(property.Reference.ExampleClass.String()); resourceName != "" {
			return newExampleReference(ResourceReference, resourceName, property.Reference.ExampleAttributeName)
		}
	}

	return c.normalizeGeneratedExampleValue(property, property.exampleValueFromTestScenario(ds, c.Name, scenario))
}

func (c *Class) propertyExampleScenario(property *Property, variant ExampleVariantEnum) TestValueScenario {
	// Select the test-value source that corresponds to the requested public example.
	switch variant {
	case MinimumExample:
		if property.TestValues == nil {
			return TestValueScenario{}
		}
		return property.TestValues.Default
	case FullExample:
		if property.Required || slices.Contains(c.IdentifiedBy, property.PropertyName) {
			if property.TestValues == nil {
				return TestValueScenario{}
			}
			return property.TestValues.Default
		}
		scenario := TestValueScenario{}
		if property.TestValues != nil {
			scenario = property.TestValues.Create
		}
		if property.IgnoreInTest && (!scenario.Defined || len(scenario.Entries) == 0) {
			if generated := property.generateTestValues(); generated != nil {
				scenario = generated.Create
			}
		}
		return scenario
	default:
		panic(fmt.Sprintf("unsupported example variant %d", variant))
	}
}

func (d *TestDependency) exampleValue(ds *DataStore) ExampleValue {
	resourceName := ""
	if d.Class != nil {
		resourceName = ds.getResourceName(d.Class.String())
	}
	if d.ReferenceType == StaticReference {
		if d.Class != nil {
			if target, ok := ds.Classes[d.Class.String()]; ok && target.HasResourceArtifact() && resourceName != "" {
				return newExampleReference(ResourceReference, resourceName, "id")
			}
		}
		return ExampleLiteral(d.Reference)
	}
	if resourceName == "" {
		ds.ctx.Diagnostics.AddError("Dependency class '%s' has no Terraform resource name for public-example reference '%s'.", d.Class, d.Reference)
		return ExampleLiteral(d.Reference)
	}
	return newExampleReference(d.ReferenceType, resourceName, "id")
}

func (c *Class) exampleTargetDependency(context ExampleContext) *TestDependency {
	candidates := c.Relation.ToClasses
	if context.ParentClass != nil {
		for _, dependency := range c.TestDependencies {
			if dependency.Role == Parent && dependency.Class.Equal(context.ParentClass) && len(dependency.TargetClasses) > 0 {
				candidates = dependency.TargetClasses
				break
			}
		}
	}

	for _, targetClass := range candidates {
		for _, dependency := range c.TestDependencies {
			if dependency.Role == Target && dependency.Class.Equal(targetClass) {
				return dependency
			}
		}
	}
	for _, dependency := range c.TestDependencies {
		if dependency.Role == Target {
			return dependency
		}
	}
	return nil
}

func (c *Class) setResourceExamples(ds *DataStore) {
	// Assemble complete minimum and full resource artifacts.
	c.ResourceExamples = ResourceExamples{}
	if !c.HasResourceArtifact() {
		return
	}

	c.ResourceExamples = ResourceExamples{
		Minimum: Example{
			Platform: c.ExampleContext.Platform,
			Body:     c.buildTopLevelExampleObject(ds, MinimumExample),
		},
		Full: Example{
			Platform: c.ExampleContext.Platform,
			Body:     c.buildTopLevelExampleObject(ds, FullExample),
		},
		Files: append([]string(nil), c.ClassDefinition.ExampleFiles...),
	}
}

func (c *Class) buildTopLevelExampleObject(ds *DataStore, variant ExampleVariantEnum) ExampleObject {
	// Keep parent_dn first and then follow schema order.
	object := ExampleObject{}
	if parentDn := c.ParentDnProperty(); parentDn != nil {
		if value := c.propertyExampleValue(ds, parentDn, c.ExampleContext, variant); value != nil {
			object.Attributes = append(object.Attributes, ExampleAttribute{Name: parentDn.AttributeName, Value: value})
		}
	}
	for _, propertyName := range c.PropertiesAll {
		property := c.Properties[propertyName]
		if property.PropertyName == "parentDn" {
			continue
		}
		value := c.propertyExampleValue(ds, property, c.ExampleContext, variant)
		if value == nil {
			continue
		}
		object.Attributes = append(object.Attributes, ExampleAttribute{Name: property.AttributeName, Value: value})
	}
	object.Children = c.buildExampleChildren(ds, c.TestChildren, variant)
	return object
}

func (c *Class) buildExampleChildren(ds *DataStore, testChildren []*TestChild, variant ExampleVariantEnum) []ExampleChild {
	// Project normalized test children into public-example objects.
	children := make([]ExampleChild, 0, len(testChildren))
	for _, testChild := range testChildren {
		childClass, ok := ds.Classes[testChild.Class.String()]
		if !ok || (variant == MinimumExample && !childClass.RequiredAsChild) || len(testChild.Instances) == 0 {
			continue
		}

		// Public examples intentionally show one representative child instance. Test data
		// may contain a second instance solely to exercise update and set identity behavior.
		exampleChild := ExampleChild{
			Name:     childClass.ResourceNameNested,
			IsSingle: childClass.IsSingleNestedWhenDefinedAsChild,
			Body:     childClass.buildChildExampleObject(ds, c.Name, testChild.Instances[0], variant),
		}
		children = append(children, exampleChild)
	}
	return children
}

func (c *Class) buildChildExampleObject(ds *DataStore, parentClass *ClassName, instance TestChildInstance, variant ExampleVariantEnum) ExampleObject {
	// Apply this child class's own example rules in its parent context.
	object := ExampleObject{}
	context := ExampleContext{ParentClass: parentClass}
	for _, propertyName := range c.PropertiesAll {
		property := c.Properties[propertyName]
		if property.PropertyName == "parentDn" || property.ReadOnly {
			continue
		}

		value, defined := property.exampleDefinitionValue(ds, c.Name, variant)
		if !defined {
			if variant == MinimumExample && !property.Required {
				continue
			}
			scenario, exists := instance.Properties[property.AttributeName]
			if !exists {
				scenario = c.propertyExampleScenario(property, variant)
			}
			value = c.resolvePropertyExampleValue(ds, property, context, scenario)
		}
		if value == nil {
			continue
		}
		object.Attributes = append(object.Attributes, ExampleAttribute{Name: property.AttributeName, Value: value})
	}
	object.Children = c.buildExampleChildren(ds, instance.Children, variant)
	return object
}

func (c *Class) normalizeGeneratedExampleValue(property *Property, value ExampleValue) ExampleValue {
	// Rename only automatically generated literals. Explicit values, references, and
	// expressions retain their authored representation.
	switch value := value.(type) {
	case ExampleLiteral:
		return ExampleLiteral(c.normalizeGeneratedExampleLiteral(property, string(value)))
	case ExampleCollection:
		for index, entry := range value {
			literal, ok := entry.(ExampleLiteral)
			if ok {
				value[index] = ExampleLiteral(c.normalizeGeneratedExampleLiteral(property, string(literal)))
			}
		}
		return value
	default:
		return value
	}
}

func (c *Class) normalizeGeneratedExampleLiteral(property *Property, value string) string {
	// Convert only recognizable generated placeholders into readable documentation
	// values, leaving explicitly authored literals untouched.
	// annotation is the one auto-generated string whose Create value has no `_1` suffix.
	generated := value == property.AttributeName+"_1" ||
		value == "test_"+property.AttributeName ||
		(property.AttributeName == "annotation" && value == "annotation")
	if !generated {
		return value
	}
	if override, ok := property.globalDefinition.PropertyExampleValueOverrides[property.PropertyName]; ok {
		return override
	}

	baseName := c.ResourceName
	if baseName == "" {
		baseName = c.ResourceNameNested
	}
	if baseName == "" && c.Name != nil {
		baseName = strings.ToLower(c.Name.String())
	}
	switch property.AttributeName {
	case "name":
		return baseName
	case "description":
		label := c.Documentation.Label
		if label == "" {
			label = strings.ReplaceAll(baseName, "_", " ")
		}
		return "Example " + label
	default:
		return "example_" + property.AttributeName
	}
}

func (c *Class) setDataSourceExample(ds *DataStore) {
	// Emit only the values needed to identify an existing object.
	c.DataSourceExample = Example{}
	if !c.HasDatasourceArtifact() {
		return
	}

	body := ExampleObject{}
	if parentDn := c.ParentDnProperty(); parentDn != nil {
		if value := c.propertyExampleValue(ds, parentDn, c.ExampleContext, MinimumExample); value != nil {
			body.Attributes = append(body.Attributes, ExampleAttribute{Name: parentDn.AttributeName, Value: value})
		}
	}
	for _, propertyName := range c.PropertiesAll {
		if propertyName == "parentDn" || !slices.Contains(c.IdentifiedBy, propertyName) {
			continue
		}
		property := c.Properties[propertyName]
		value := c.propertyExampleValue(ds, property, c.ExampleContext, MinimumExample)
		if value == nil {
			continue
		}
		body.Attributes = append(body.Attributes, ExampleAttribute{Name: property.AttributeName, Value: value})
	}
	c.DataSourceExample = Example{
		Platform: c.ExampleContext.Platform,
		Body:     body,
	}
}
