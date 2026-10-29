package data

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// whitespaceRegex collapses consecutive whitespace characters (spaces, tabs, newlines) into a single space.
var whitespaceRegex = regexp.MustCompile(`\s+`)

// DefaultValue describes a single default value for a property.
// Mirrors the ValidValue pattern and is designed to be extended with additional fields.
type DefaultValue struct {
	// The localName of the default value (e.g. "enabled").
	Value string
	// The APIC versions for which this default value applies.
	// Nil means the default applies to all versions.
	Versions *Versions
}

type PropertyDocumentation struct {
	// The default values of the property in APIC, expressed as localNames.
	// Stored as a list to track default value changes across APIC versions.
	DefaultValues []DefaultValue
	// A generic explanation of the property and its usage.
	// When applicable, a reference to classes the property points to and which resources/datasources are used for this is included.
	// When version is higher than the class version, a property specific version is included.
	Description string
	// Notes rendered with -> prefix in the property documentation.
	Notes []string
	// Warnings rendered with ~> prefix in the property documentation.
	Warnings []string
	// SupportedVersions is rendered when this property is not available for the
	// complete version range of its owning class.
	SupportedVersions string
}

// ResourceDocumentationDefaultValues returns the defaults displayed in resource documentation.
// The annotation default is provider-driven and is therefore not available in APIC metadata.
func (p *Property) ResourceDocumentationDefaultValues() []DefaultValue {
	if p.PropertyName == "annotation" {
		return []DefaultValue{{Value: "orchestrator:terraform"}}
	}

	return p.Documentation.DefaultValues
}

func (p *Property) setDocumentation() error {
	genLogger.Debugf("Setting Documentation for property '%s'.", p.PropertyName)

	p.Documentation.setDescription(p)

	p.Documentation.setNotes(p)

	p.Documentation.setWarnings(p)

	err := p.Documentation.setDefaultValues(p)
	if err != nil {
		return err
	}

	genLogger.Debugf("Successfully set Documentation for property '%s'.", p.PropertyName)
	return nil
}

func (d *PropertyDocumentation) setDescription(p *Property) {
	genLogger.Debugf("Setting Documentation Description for property '%s'.", p.PropertyName)

	if p.propertyDefinition.Documentation.Description != "" {
		d.Description = p.propertyDefinition.Documentation.Description
	} else if comment, ok := p.metaDetails["comment"].([]any); ok && len(comment) > 0 {
		parts := make([]string, 0, len(comment))
		for _, entry := range comment {
			if s, ok := entry.(string); ok {
				parts = append(parts, s)
			}
		}
		d.Description = whitespaceRegex.ReplaceAllString(strings.Join(parts, " "), " ")
	} else if label, ok := p.metaDetails["label"].(string); ok {
		d.Description = label
	}

	genLogger.Debugf("Successfully set Documentation Description for property '%s'. Description: %s", p.PropertyName, d.Description)
}

func (c *Class) setPropertyDocumentationReferences(ds *DataStore) {
	// Mention only artifacts that users can actually address in configuration;
	// a class name alone is insufficient evidence that a registry page exists.
	for _, propertyName := range c.PropertiesAll {
		property := c.Properties[propertyName]
		if property == nil || property.Reference == nil || property.Reference.ExampleClass == nil {
			continue
		}
		description := strings.ToLower(property.Documentation.Description)
		if strings.Contains(description, "can be referenced from") || strings.Contains(description, "registry.terraform.io/") {
			continue
		}

		attributeName := property.Reference.ExampleAttributeName
		var references []string
		if resourceName, url := ds.getTerraformRegistryDocumentationReference(property.Reference.ExampleClass.String(), ResourceArtifact); resourceName != "" {
			references = append(references, fmt.Sprintf(
				"a [resource](%s) with `aci_%s.example.%s`",
				url,
				resourceName,
				attributeName,
			))
		}
		if resourceName, url := ds.getTerraformRegistryDocumentationReference(property.Reference.ExampleClass.String(), DatasourceArtifact); resourceName != "" {
			references = append(references, fmt.Sprintf(
				"a [datasource](%s) with `data.aci_%s.example.%s`",
				url,
				resourceName,
				attributeName,
			))
		}
		if len(references) == 0 {
			continue
		}

		referenceDescription := "This attribute can be referenced from " + strings.Join(references, " or from ") + "."
		property.Documentation.Description = strings.TrimSpace(property.Documentation.Description + " " + referenceDescription)
	}
}

// setPropertyDocumentationDescriptions applies global description overrides
// and resolves class-label placeholders in the resulting descriptions.
//
// The override is keyed by APIC PropertyName (not the snake_case schema
// AttributeName). Per the section 9.1 contract this branch sits between the
// per-class override and the meta-comment fallback: a property carrying its
// own ClassDefinition.Properties[name].Documentation.Description keeps that
// value; otherwise the global template wins over whatever Property.setDescription
// derived from meta comment / meta label.
//
// Must run after ClassDocumentation.setLabel so the substitution target exists.
func (c *Class) setPropertyDocumentationDescriptions(ds *DataStore) {
	overrides := ds.GlobalMetaDefinition.PropertyDocumentationOverrides
	for name, property := range c.Properties {
		if template, ok := overrides[name]; ok && property.propertyDefinition.Documentation.Description == "" {
			property.Documentation.Description = template
			genLogger.Debugf("Applied global PropertyDocumentationOverride for class '%s' property '%s'. Description: %s", c.Name.full, name, property.Documentation.Description)
		}

		if strings.Contains(property.Documentation.Description, "%s") {
			property.Documentation.Description = strings.ReplaceAll(property.Documentation.Description, "%s", c.Documentation.Label)
		}
		if property.Documentation.Description != "" && !strings.HasSuffix(property.Documentation.Description, ".") {
			property.Documentation.Description += "."
		}
	}
}

func (c *Class) setPropertyDocumentationSupportedVersions() {
	for _, property := range c.Properties {
		if property.SupportedVersions != nil && !property.SupportedVersions.Covers(c.SupportedVersions) {
			property.Documentation.SupportedVersions = documentationVersions(property.SupportedVersions)
		}
	}
}

func (d *PropertyDocumentation) setNotes(p *Property) {
	genLogger.Debugf("Setting Documentation Notes for property '%s'.", p.PropertyName)

	d.Notes = p.propertyDefinition.Documentation.Notes

	genLogger.Debugf("Successfully set Documentation Notes for property '%s'. Notes: %v", p.PropertyName, d.Notes)
}

func (d *PropertyDocumentation) setWarnings(p *Property) {
	genLogger.Debugf("Setting Documentation Warnings for property '%s'.", p.PropertyName)

	d.Warnings = p.propertyDefinition.Documentation.Warnings

	genLogger.Debugf("Successfully set Documentation Warnings for property '%s'. Warnings: %v", p.PropertyName, d.Warnings)
}

func (d *PropertyDocumentation) setDefaultValues(p *Property) error {
	genLogger.Debugf("Setting Documentation DefaultValues for property '%s'.", p.PropertyName)

	if len(p.propertyDefinition.DefaultValues) > 0 {
		d.DefaultValues = make([]DefaultValue, 0, len(p.propertyDefinition.DefaultValues))
		for value, versionStr := range p.propertyDefinition.DefaultValues {
			dv := DefaultValue{Value: value}
			if versionStr != "" {
				parsedVersions, err := NewVersions(versionStr)
				if err != nil {
					return fmt.Errorf("failed to parse default value versions for property '%s', value '%s': %w", p.PropertyName, value, err)
				}
				dv.Versions = parsedVersions
			}
			d.DefaultValues = append(d.DefaultValues, dv)
		}
		slices.SortFunc(d.DefaultValues, func(a, b DefaultValue) int {
			return strings.Compare(a.Value, b.Value)
		})
		genLogger.Debugf("Successfully set Documentation DefaultValues for property '%s' from definition override. DefaultValues: %v", p.PropertyName, d.DefaultValues)
		return nil
	}

	var defaultValue string
	switch value := p.metaDetails["default"].(type) {
	case string:
		if entry, ok := p.ValidValues[value]; ok {
			defaultValue = entry.LocalName
		} else if p.ValueType == Set && value == "none" {
			defaultValue = ""
		} else {
			defaultValue = value
		}
	case float64:
		defaultValue = fmt.Sprintf("%g", value)
	default:
		genLogger.Debugf("No default value found for property '%s'.", p.PropertyName)
		return nil
	}

	d.DefaultValues = []DefaultValue{{Value: defaultValue}}

	genLogger.Debugf("Successfully set Documentation DefaultValues for property '%s'. DefaultValues: %v", p.PropertyName, d.DefaultValues)
	return nil
}
