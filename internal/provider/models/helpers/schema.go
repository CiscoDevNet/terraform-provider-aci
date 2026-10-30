package helpers

import (
	"maps"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/planmodifiers"
	aciValidators "github.com/CiscoDevNet/terraform-provider-aci/v2/internal/validators"
)

const distinguishedNameDescription = "The distinguished name (DN) of the object."

// ParentDNAttributeConfig defines the optional default and accepted values for
// a top-level parent_dn attribute.
type ParentDNAttributeConfig struct {
	DefaultValue string
	ValidValues  []string
}

// PropertyType identifies the Terraform Framework attribute shape generated
// for a class property.
type PropertyType uint8

const (
	// StringProperty constructs a Terraform string attribute.
	StringProperty PropertyType = iota
	// SetProperty constructs a Terraform set attribute.
	SetProperty
	// ObjectProperty constructs a Terraform object attribute.
	ObjectProperty
)

// StringRange defines the numeric string range accepted alongside named valid
// values by semantic-equality properties.
type StringRange struct {
	Min int64
	Max int64
}

// PropertySchema is the neutral definition used to construct both resource
// and data-source attributes for one class property.
type PropertySchema struct {
	Type                 PropertyType
	Required             bool
	Optional             bool
	Computed             bool
	ReadOnly             bool
	Sensitive            bool
	RequiresReplace      bool
	Description          string
	CustomStringType     basetypes.StringTypable
	SetElementType       attr.Type
	ObjectAttributeTypes map[string]attr.Type
	ValidValues          []string
	StringRange          *StringRange
}

type requiredPropertyRendering uint8

const (
	// requiredPropertyDirect uses Terraform's native Required flag.
	requiredPropertyDirect requiredPropertyRendering = iota
	// requiredPropertyWithValidator preserves required semantics when a
	// computed schema context requires the property to be optional and computed.
	requiredPropertyWithValidator
	// requiredPropertyWithObjectValidator leaves validation to an optional
	// singleton child so an empty object can still represent deletion.
	requiredPropertyWithObjectValidator
)

// resourceSchemaContext carries root behavior through recursively constructed
// children without adding context-specific fields to the neutral class model.
type resourceSchemaContext struct {
	includeRequiresReplace    bool
	stateManagedByParent      bool
	rootHasReadOnlyProperties bool
	requiredPropertyRendering requiredPropertyRendering
}

func (p PropertySchema) resourceStringValidators(requiredRendering requiredPropertyRendering) []validator.String {
	validators := make([]validator.String, 0, 2)
	if len(p.ValidValues) > 0 {
		validValues := stringvalidator.OneOf(p.ValidValues...)
		if p.StringRange == nil {
			validators = append(validators, validValues)
		} else {
			validators = append(validators, stringvalidator.Any(
				validValues,
				aciValidators.InBetweenFromString(float64(p.StringRange.Min), float64(p.StringRange.Max)),
			))
		}
	}
	if p.Required && requiredRendering == requiredPropertyWithValidator {
		validators = append(validators, aciValidators.StringNotNull())
	}
	if len(validators) == 0 {
		return nil
	}
	return validators
}

func (p PropertySchema) dataSourceStringValidators() []validator.String {
	if len(p.ValidValues) == 0 {
		return nil
	}

	return []validator.String{stringvalidator.OneOf(p.ValidValues...)}
}

func (p PropertySchema) setValidators() []validator.Set {
	if len(p.ValidValues) == 0 {
		return nil
	}

	return []validator.Set{
		setvalidator.SizeAtMost(len(p.ValidValues)),
		setvalidator.ValueStringsAre(stringvalidator.OneOf(p.ValidValues...)),
	}
}

func (p PropertySchema) resourceAttribute(context resourceSchemaContext) resourceschema.Attribute {
	required := p.Required
	optional := p.Optional
	computed := p.Computed
	if p.Required && context.requiredPropertyRendering != requiredPropertyDirect {
		required = false
		optional = true
		computed = true
	}

	switch p.Type {
	case SetProperty:
		attribute := resourceschema.SetAttribute{
			ElementType:         p.SetElementType,
			Required:            required,
			Optional:            optional,
			Computed:            computed,
			Sensitive:           p.Sensitive,
			Validators:          p.setValidators(),
			MarkdownDescription: p.Description,
		}
		if !p.ReadOnly && !context.stateManagedByParent {
			attribute.PlanModifiers = ResourceSetPlanModifiers(context.includeRequiresReplace && p.RequiresReplace)
		}
		return attribute
	case ObjectProperty:
		attribute := resourceschema.ObjectAttribute{
			AttributeTypes:      p.ObjectAttributeTypes,
			Required:            required,
			Optional:            optional,
			Computed:            computed,
			Sensitive:           p.Sensitive,
			MarkdownDescription: p.Description,
		}
		if !p.ReadOnly && !context.stateManagedByParent {
			attribute.PlanModifiers = []planmodifier.Object{
				objectplanmodifier.UseNonNullStateForUnknown(),
			}
		}
		return attribute
	default:
		attribute := resourceschema.StringAttribute{
			CustomType:          p.CustomStringType,
			Required:            required,
			Optional:            optional,
			Computed:            computed,
			Sensitive:           p.Sensitive,
			Validators:          p.resourceStringValidators(context.requiredPropertyRendering),
			MarkdownDescription: p.Description,
		}
		if !p.ReadOnly && !context.stateManagedByParent {
			attribute.PlanModifiers = ResourceStringPlanModifiers(context.includeRequiresReplace && p.RequiresReplace)
		}
		return attribute
	}
}

// ResourceAttribute constructs a direct resource schema attribute. Class
// schema construction supplies additional context for nested attributes and
// classes containing read-only properties.
func (p PropertySchema) ResourceAttribute(includeRequiresReplace bool) resourceschema.Attribute {
	return p.resourceAttribute(resourceSchemaContext{
		includeRequiresReplace:    includeRequiresReplace,
		requiredPropertyRendering: requiredPropertyDirect,
	})
}

// DataSourceAttribute constructs either a computed child value or a required
// top-level lookup attribute from the same property definition.
func (p PropertySchema) DataSourceAttribute(required bool) datasourceschema.Attribute {
	switch p.Type {
	case SetProperty:
		attribute := datasourceschema.SetAttribute{
			ElementType:         p.SetElementType,
			Required:            required,
			Computed:            !required,
			Sensitive:           p.Sensitive,
			MarkdownDescription: p.Description,
		}
		if required {
			attribute.Validators = p.setValidators()
		}
		return attribute
	case ObjectProperty:
		return datasourceschema.ObjectAttribute{
			AttributeTypes:      p.ObjectAttributeTypes,
			Required:            required,
			Computed:            !required,
			Sensitive:           p.Sensitive,
			MarkdownDescription: p.Description,
		}
	default:
		attribute := datasourceschema.StringAttribute{
			CustomType:          p.CustomStringType,
			Required:            required,
			Computed:            !required,
			Sensitive:           p.Sensitive,
			MarkdownDescription: p.Description,
		}
		if required {
			attribute.Validators = p.dataSourceStringValidators()
		}
		return attribute
	}
}

// ChildCardinality identifies whether a child is represented as one nested
// object or a set of nested objects.
type ChildCardinality uint8

const (
	// SingleChild constructs a single nested child attribute.
	SingleChild ChildCardinality = iota
	// SetChild constructs a set nested child attribute.
	SetChild
)

// ClassSchemaFactory resolves a child class schema without creating package
// initialization dependencies between generated class definitions.
type ClassSchemaFactory func() ClassSchema

// ChildSchema defines the parent-owned name, cardinality, and configuration
// behavior of one child class.
type ChildSchema struct {
	Cardinality ChildCardinality
	Required    bool
	Description string
	Schema      ClassSchemaFactory
}

func (c ChildSchema) resourceAttribute(attributeName string, parentContext resourceSchemaContext) resourceschema.Attribute {
	childSchema := c.Schema()
	childContext := resourceSchemaContext{
		stateManagedByParent:      true,
		rootHasReadOnlyProperties: parentContext.rootHasReadOnlyProperties,
		requiredPropertyRendering: requiredPropertyDirect,
	}
	if c.Cardinality == SingleChild && !c.Required {
		childContext.requiredPropertyRendering = requiredPropertyWithObjectValidator
	} else if c.Cardinality == SetChild && parentContext.rootHasReadOnlyProperties {
		childContext.requiredPropertyRendering = requiredPropertyWithValidator
	}
	attributes := childSchema.resourceAttributes(childContext)

	if c.Cardinality == SingleChild {
		attribute := resourceschema.SingleNestedAttribute{
			Required:            c.Required,
			Optional:            !c.Required,
			Computed:            !c.Required,
			MarkdownDescription: c.Description,
			Attributes:          attributes,
		}
		if !parentContext.stateManagedByParent {
			attribute.PlanModifiers = []planmodifier.Object{
				planmodifiers.ReconcileSingletonChild(childSchema.reconciliationSchema()),
			}
		}
		if !c.Required {
			requiredProperties := childSchema.requiredPropertyNames()
			if len(requiredProperties) > 0 {
				attribute.Validators = []validator.Object{
					aciValidators.RequiredAttributesWhenNotEmpty(attributeName, requiredProperties),
				}
			}
		}
		return attribute
	}

	attribute := resourceschema.SetNestedAttribute{
		Required:            c.Required,
		Optional:            !c.Required,
		Computed:            !c.Required,
		MarkdownDescription: c.Description,
		NestedObject:        resourceschema.NestedAttributeObject{Attributes: attributes},
	}
	if !parentContext.stateManagedByParent {
		attribute.PlanModifiers = []planmodifier.Set{
			planmodifiers.ReconcileRepeatedChild(childSchema.reconciliationSchema()),
		}
	}
	return attribute
}

// DataSourceAttribute constructs the computed data-source attribute for the
// child class.
func (c ChildSchema) DataSourceAttribute() datasourceschema.Attribute {
	attributes := c.Schema().DataSourceAttributes()
	if c.Cardinality == SingleChild {
		return datasourceschema.SingleNestedAttribute{
			Computed:            true,
			MarkdownDescription: c.Description,
			Attributes:          attributes,
		}
	}

	return datasourceschema.SetNestedAttribute{
		Computed:            true,
		MarkdownDescription: c.Description,
		NestedObject:        datasourceschema.NestedAttributeObject{Attributes: attributes},
	}
}

// ParentDNSchema defines an optional top-level parent_dn attribute.
type ParentDNSchema struct {
	AttributeName string
	Config        ParentDNAttributeConfig
}

// ClassSchema is the generated, framework-neutral schema definition for one
// APIC class. It constructs concrete resource and data-source schemas on
// demand.
type ClassSchema struct {
	ResourceDescription   string
	DataSourceDescription string
	Properties            map[string]PropertySchema
	Children              map[string]ChildSchema
	Identifiers           []string
	ParentDN              *ParentDNSchema
	AnnotationAttribute   string
}

func (s ClassSchema) hasReadOnlyProperties() bool {
	for _, property := range s.Properties {
		if property.ReadOnly {
			return true
		}
	}
	return false
}

func (s ClassSchema) requiredPropertyNames() []string {
	properties := make([]string, 0, len(s.Properties))
	for name, property := range s.Properties {
		if property.Required {
			properties = append(properties, name)
		}
	}
	slices.Sort(properties)
	return properties
}

func (s ClassSchema) topLevelResourceContext(includeRequiresReplace bool) resourceSchemaContext {
	context := resourceSchemaContext{
		includeRequiresReplace:    includeRequiresReplace,
		rootHasReadOnlyProperties: s.hasReadOnlyProperties(),
		requiredPropertyRendering: requiredPropertyDirect,
	}
	if context.rootHasReadOnlyProperties {
		context.requiredPropertyRendering = requiredPropertyWithValidator
	}
	return context
}

func (s ClassSchema) reconciliationSchema() planmodifiers.ReconciliationSchema {
	children := make(map[string]planmodifiers.ReconciliationChild, len(s.Children))
	for name, child := range s.Children {
		children[name] = planmodifiers.ReconciliationChild{
			Repeated: child.Cardinality == SetChild,
			Schema:   child.Schema().reconciliationSchema(),
		}
	}

	return planmodifiers.ReconciliationSchema{
		Identifiers: slices.Clone(s.Identifiers),
		Children:    children,
	}
}

func (s ClassSchema) resourcePropertyAttributes(context resourceSchemaContext) map[string]resourceschema.Attribute {
	attributes := make(map[string]resourceschema.Attribute, len(s.Properties))
	for name, property := range s.Properties {
		attributes[name] = property.resourceAttribute(context)
	}
	return attributes
}

// ResourcePropertyAttributes constructs the direct resource properties for
// the class.
func (s ClassSchema) ResourcePropertyAttributes(includeRequiresReplace bool) map[string]resourceschema.Attribute {
	return s.resourcePropertyAttributes(s.topLevelResourceContext(includeRequiresReplace))
}

func (s ClassSchema) resourceChildAttributes(context resourceSchemaContext) map[string]resourceschema.Attribute {
	attributes := make(map[string]resourceschema.Attribute, len(s.Children))
	for name, child := range s.Children {
		attributes[name] = child.resourceAttribute(name, context)
	}
	return attributes
}

// ResourceChildAttributes constructs the resource attributes owned by the
// class for its child classes.
func (s ClassSchema) ResourceChildAttributes() map[string]resourceschema.Attribute {
	return s.resourceChildAttributes(s.topLevelResourceContext(false))
}

func (s ClassSchema) resourceAttributes(context resourceSchemaContext) map[string]resourceschema.Attribute {
	attributes := s.resourcePropertyAttributes(context)
	maps.Copy(attributes, s.resourceChildAttributes(context))
	return attributes
}

// ResourceAttributes combines the class's direct properties and children.
func (s ClassSchema) ResourceAttributes(includeRequiresReplace bool) map[string]resourceschema.Attribute {
	return s.resourceAttributes(s.topLevelResourceContext(includeRequiresReplace))
}

// DataSourcePropertyAttributes constructs the computed direct properties for
// the class when it is read from a data source.
func (s ClassSchema) DataSourcePropertyAttributes() map[string]datasourceschema.Attribute {
	attributes := make(map[string]datasourceschema.Attribute, len(s.Properties))
	for name, property := range s.Properties {
		attributes[name] = property.DataSourceAttribute(false)
	}
	return attributes
}

// DataSourceChildAttributes constructs the computed data-source attributes
// owned by the class for its child classes.
func (s ClassSchema) DataSourceChildAttributes() map[string]datasourceschema.Attribute {
	attributes := make(map[string]datasourceschema.Attribute, len(s.Children))
	for name, child := range s.Children {
		attributes[name] = child.DataSourceAttribute()
	}
	return attributes
}

// DataSourceAttributes combines the class's direct properties and children.
func (s ClassSchema) DataSourceAttributes() map[string]datasourceschema.Attribute {
	attributes := s.DataSourcePropertyAttributes()
	maps.Copy(attributes, s.DataSourceChildAttributes())
	return attributes
}

// ResourceSchema constructs the complete top-level resource schema.
func (s ClassSchema) ResourceSchema(defaultAnnotation string) resourceschema.Schema {
	attributes := s.ResourceAttributes(true)
	attributes["id"] = ResourceIDAttribute()
	if s.ParentDN != nil {
		attributes[s.ParentDN.AttributeName] = ResourceParentDNAttribute(s.ParentDN.Config)
	}
	if s.AnnotationAttribute != "" {
		attributes[s.AnnotationAttribute] = ResourceAnnotationAttribute(
			defaultAnnotation,
			s.Properties[s.AnnotationAttribute].Description,
		)
	}

	return resourceschema.Schema{
		MarkdownDescription: s.ResourceDescription,
		Attributes:          attributes,
	}
}

// DataSourceSchema constructs the complete top-level data-source schema.
func (s ClassSchema) DataSourceSchema() datasourceschema.Schema {
	attributes := s.DataSourceAttributes()
	attributes["id"] = DataSourceIDAttribute()
	if s.ParentDN != nil {
		attributes[s.ParentDN.AttributeName] = DataSourceParentDNAttribute(s.ParentDN.Config)
	}
	for _, identifier := range s.Identifiers {
		attributes[identifier] = s.Properties[identifier].DataSourceAttribute(true)
	}

	return datasourceschema.Schema{
		MarkdownDescription: s.DataSourceDescription,
		Attributes:          attributes,
	}
}

// ResourceIDAttribute returns the common computed resource ID attribute.
func ResourceIDAttribute() resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Computed:            true,
		MarkdownDescription: distinguishedNameDescription,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseNonNullStateForUnknown(),
		},
	}
}

// DataSourceIDAttribute returns the common computed data-source ID attribute.
func DataSourceIDAttribute() datasourceschema.StringAttribute {
	return datasourceschema.StringAttribute{
		Computed:            true,
		MarkdownDescription: distinguishedNameDescription,
	}
}

// ResourceParentDNAttribute returns a required parent_dn attribute unless a
// default is supplied, in which case the attribute is optional and computed.
func ResourceParentDNAttribute(config ParentDNAttributeConfig) resourceschema.StringAttribute {
	attribute := resourceschema.StringAttribute{
		MarkdownDescription: distinguishedNameDescription,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseNonNullStateForUnknown(),
			stringplanmodifier.RequiresReplace(),
		},
	}

	if config.DefaultValue == "" {
		attribute.Required = true
	} else {
		attribute.Optional = true
		attribute.Computed = true
		attribute.Default = stringdefault.StaticString(config.DefaultValue)
	}

	if len(config.ValidValues) > 0 {
		attribute.Validators = []validator.String{
			stringvalidator.OneOf(config.ValidValues...),
		}
	}

	return attribute
}

// DataSourceParentDNAttribute returns a required parent_dn lookup attribute
// unless the data source applies a known default during its read operation.
func DataSourceParentDNAttribute(config ParentDNAttributeConfig) datasourceschema.StringAttribute {
	attribute := datasourceschema.StringAttribute{
		MarkdownDescription: distinguishedNameDescription,
	}

	if config.DefaultValue == "" {
		attribute.Required = true
	} else {
		attribute.Optional = true
	}

	if len(config.ValidValues) > 0 {
		attribute.Validators = []validator.String{
			stringvalidator.OneOf(config.ValidValues...),
		}
	}

	return attribute
}

// ResourceAnnotationAttribute returns the top-level resource annotation
// attribute with the provider-resolved static default.
func ResourceAnnotationAttribute(defaultAnnotation, markdownDescription string) resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			planmodifiers.PreserveStringStateForUnknown(),
		},
		Default:             stringdefault.StaticString(defaultAnnotation),
		MarkdownDescription: markdownDescription,
	}
}

// ResourceStringPlanModifiers returns the common string plan modifiers used by
// generated resource properties.
func ResourceStringPlanModifiers(requiresReplace bool) []planmodifier.String {
	modifiers := []planmodifier.String{
		planmodifiers.PreserveStringStateForUnknown(),
	}
	if requiresReplace {
		modifiers = append(modifiers, stringplanmodifier.RequiresReplace())
	}
	return modifiers
}

// ResourceSetPlanModifiers returns the common set plan modifiers used by
// generated resource properties.
func ResourceSetPlanModifiers(requiresReplace bool) []planmodifier.Set {
	modifiers := []planmodifier.Set{
		planmodifiers.PreserveSetStateForUnknown(),
	}
	if requiresReplace {
		modifiers = append(modifiers, setplanmodifier.RequiresReplace())
	}
	return modifiers
}
