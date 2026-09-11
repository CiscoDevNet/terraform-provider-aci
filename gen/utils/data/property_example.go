package data

func (p *Property) exampleValue(values []ExampleValue) ExampleValue {
	// Preserve a set as one collection value while scalar properties use their first value.
	if p.ValueType == Set {
		return ExampleCollection(values)
	}
	if len(values) == 0 {
		return ExampleLiteral("")
	}
	return values[0]
}

func (p *Property) exampleScenarioIncluded(scenario TestValueScenario) bool {
	// Inclusion can be decided before values are converted, allowing semantic
	// references to bypass their string-shaped acceptance-test representation.
	if p.ValueType == Set && scenario.Defined && len(scenario.Entries) == 0 {
		return true
	}
	for _, entry := range scenario.Entries {
		if entry.ConfigInclude {
			return true
		}
	}
	return false
}

func (p *Property) exampleValueFromTestScenario(ds *DataStore, className *ClassName, scenario TestValueScenario) ExampleValue {
	if !scenario.Defined {
		return nil
	}
	values := make([]ExampleValue, 0, len(scenario.Entries))
	included := p.ValueType == Set && len(scenario.Entries) == 0
	for _, entry := range scenario.Entries {
		if !entry.ConfigInclude {
			continue
		}
		included = true
		if entry.ValueType == ReferenceValue {
			ds.ctx.Diagnostics.AddError("Class '%s': property '%s' has a reference test value without semantic example-reference data.", className, p.AttributeName)
			continue
		}
		values = append(values, entry.exampleValue())
	}
	if !included {
		return nil
	}
	return p.exampleValue(values)
}

func (e TestValueEntry) exampleValue() ExampleValue {
	// Reference values are resolved from class semantics before reaching this conversion.
	switch e.ValueType {
	case ExpressionValue:
		return ExampleExpression(e.ConfigValue)
	default:
		return ExampleLiteral(e.ConfigValue)
	}
}

func (p *Property) exampleDefinitionValue(ds *DataStore, className *ClassName, variant ExampleVariantEnum) (ExampleValue, bool) {
	// Explicit example scenarios replace inferred values and already carry the
	// structured information required to construct references.
	definition := p.propertyDefinition.ExampleConfig.Minimum
	if variant == FullExample && p.propertyDefinition.ExampleConfig.Full.Defined {
		definition = p.propertyDefinition.ExampleConfig.Full
	}
	if !definition.Defined {
		return nil, false
	}

	values := make([]ExampleValue, 0, len(definition.Entries))
	for _, entry := range definition.Entries {
		if entry.Reference != nil {
			referenceClass, err := NewClassName(entry.Reference.ClassName)
			if err != nil {
				ds.ctx.Diagnostics.AddError("Class '%s': property '%s' has invalid example reference class '%s': %v", className, p.AttributeName, entry.Reference.ClassName, err)
				continue
			}
			resourceName := entry.Reference.ResourceName
			if resourceName == "" {
				resourceName = ds.getResourceName(referenceClass.String())
			}
			if resourceName == "" {
				ds.ctx.Diagnostics.AddError("Class '%s': property '%s' example reference class '%s' has no Terraform resource name.", className, p.AttributeName, referenceClass)
				continue
			}
			attributeName := entry.Reference.AttributeName
			if attributeName == "" {
				attributeName = "id"
			}
			values = append(values, newExampleReference(ResourceReference, resourceName, attributeName))
			continue
		}
		switch entry.ValueType {
		case ReferenceValue:
			ds.ctx.Diagnostics.AddError("Class '%s': property '%s' uses a string example reference; use the structured reference field instead.", className, p.AttributeName)
		case ExpressionValue:
			values = append(values, ExampleExpression(entry.Value))
		default:
			values = append(values, ExampleLiteral(entry.Value))
		}
	}
	return p.exampleValue(values), true
}
