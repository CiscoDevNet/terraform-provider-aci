package data

import (
	"fmt"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils"
)

type Property struct {
	// The name of the property in the resource and datasource schemas.
	AttributeName string
	// Indicates if the property is computed in the resource schemas.
	Computed bool
	// Indicates if the property has a custom type.
	// Custom types are used for valid values that have a string and an integer value pointing to the same value. (ex. ssh and 22)
	CustomType bool
	// Indicates if the property is deprecated in the resource and datasource schemas.
	// Deprecated properties include a warning the resource and datasource schemas.
	Deprecated bool
	// Documentation specific information for the property.
	Documentation PropertyDocumentation
	// Migration specific information for the property.
	// This is a map that contains the migration value details of the attribute for a specific schema version.
	MigrationValues map[int]MigrationValue
	// Indicates if a property is optional in the resource and datasource schemas.
	// Exposed as a separate bool because it directly maps to a Terraform schema construct, which makes templating easier.
	Optional bool
	// When a property that points to another class, this is the class to which the property points to.
	PointsToClass string
	// The name of the property in the APIC class.
	PropertyName string
	// Indicates if the property is read-only in the resource schemas.
	// Exposed as a separate bool because it directly maps to a Terraform schema construct, which makes templating easier.
	ReadOnly bool
	// Indicates if the property is required in the resource and datasource schemas.
	// Exposed as a separate bool because it directly maps to a Terraform schema construct, which makes templating easier.
	Required bool
	// Indicates if the property is sensitive in the resource and datasource schemas.
	// Exposed as a separate bool because it directly maps to a Terraform schema construct, which makes templating easier.
	Sensitive bool
	// The supported APIC versions for the property.
	// Each version range is separated by a comma, ex "4.2(7f)-4.2(7w),5.2(1g)-".
	// The first version is the minimum version and the second version is the maximum version.
	// A dash at the end of a range (ex. 4.2(7f)-) indicates that the class is supported from the first version to the latest version.
	// TODO: Add DeprecatedVersions field when meta file exposes deprecation information.
	SupportedVersions *Versions
	// Test specific information for the property.
	// This is used to generate the test cases and examples for the property.
	// TODO: re-evaluate the structure when creating example and test templates.
	TestValues []TestValue
	// Validation specific information for the property.
	// In the meta file for the class this is a regex statement that is used to validate the property.
	// TODO: re-evaluate the structure when creating resource templates. We might want to create a separate struct type for each type of validation.
	Validators []Validator
	// Specifies the valid values for the property when only certain values are allowed as input.
	ValidValues []ValidValue
	// The ValueTypeEnum type is used to indicate the type of the property in the resource and datasource schemas.
	ValueType ValueTypeEnum
	// The global meta definition containing global overrides. Unexported because it is only used internally by setter methods.
	globalDefinition GlobalMetaDefinition
	// The meta file details for the property. Unexported because it is only used internally by setter methods.
	metaDetails map[string]interface{}
	// The property definition overrides from the class definition file. Unexported because it is only used internally by setter methods.
	propertyDefinition PropertyDefinition
}

type MigrationValue struct {
	// The name of the property in the legacy resource schema.
	AttributeName string
	// Indicates if a property is computed in the legacy resource schema.
	Computed bool
	// Indicates if a property is optional in the legacy resource schema.
	Optional bool
	// Indicates if a property is required in the legacy resource schema.
	Required bool
	// The type of the legacy attribute.
	Type ValueTypeEnum
}

type PropertyDocumentation struct {
	// The default values of the property in APIC.
	// Defauls values is a list of valid to be able to determine if the default value is changed in versions of APIC.
	DefaultValues []ValidValue
	// A generic explanation of the property and its usage.
	// When applicable, a reference to classes the property points to and which resources/datasources are used for this is included.
	// When version is higher than the class version, a property specific version is included.
	Description string
}

type RegexStatement struct {
	// The regex string.
	Regex string
	// The type of the regex statement.
	Type RegexStatementTypeEnum
}

type RegexStatementTypeEnum int

// The enumeration options of the RegexStatement Type.
const (
	// Include indicates that the value must match the regex statement.
	Include RegexStatementTypeEnum = iota + 1
	// TODO: re-evaluate the possible regex statements type options.
)

type TestValue struct {
	// The changed value of the property to be used in the test when a property is allowed to be changed without destruction of the resource.
	Changed []string
	// The initial value of the property to be used in the test.
	// This is set to the default value of the property in APIC when it is not a required value.
	Initial []string
}

type Validator struct {
	// If the property has a range of values, these are the minimum and maximum values of the range.
	Max float64
	Min float64
	// If the property has one or more regex statements it requires to match.
	RegexList []RegexStatement
}

type ValidValue struct {
	// The valid value of the property.
	Value    string
	Versions []VersionRange
}

type ValueTypeEnum int

// The enumeration options of the ValueType.
const (
	// String indicates that the property is a string value.
	String ValueTypeEnum = iota + 1
	// Set indicates that the property is a set value.
	Set
	// List indicates that the property is a list value.
	List
)

func NewProperty(name string, details map[string]interface{}, definition PropertyDefinition, globalDefinition GlobalMetaDefinition) *Property {
	genLogger.Trace(fmt.Sprintf("Creating new property struct for property: %s.", name))

	property := &Property{
		PropertyName:       name,
		globalDefinition:   globalDefinition,
		metaDetails:        details,
		propertyDefinition: definition,
	}

	property.setPropertyData()

	genLogger.Trace(fmt.Sprintf("Successfully created new property struct for property: %s.", name))

	return property
}

func (p *Property) setPropertyData() {
	genLogger.Debug(fmt.Sprintf("Setting property data for property '%s'.", p.PropertyName))

	// TODO: add function to set AttributeName
	p.setAttributeName()

	// TODO: add function to set CustomType
	p.setCustomType()

	// TODO: add placeholder function for Deprecated
	p.setDeprecated()

	// TODO: add placeholder function for DeprecatedVersions
	p.setDeprecatedVersions()

	// TODO: add function to set Documentation
	p.setDocumentation()

	// TODO: add function to set MigrationValues
	p.setMigrationValues()

	// TODO: add function to set PointsToClass
	p.setPointsToClass()

	p.setRequired()

	// setOptional is called after setRequired because it depends on p.Required.
	p.setOptional()

	p.setReadOnly()

	// setComputed is called after setRequired because it depends on p.Required.
	p.setComputed()

	p.setSensitive()

	// TODO: add function to set TestValues
	p.setTestValues()

	// TODO: add function to set Validators
	p.setValidators()

	// TODO: add function to set ValidValues
	p.setValidValues()

	// TODO: add function to set ValueType
	p.setValueType()

	// TODO: add function to set Versions
	p.setVersions()

	genLogger.Debug(fmt.Sprintf("Successfully set property data for property '%s'.", p.PropertyName))
}

func (p *Property) setAttributeName() {
	// Determine the attribute name of the property.
	// Priority: per-class definition override > global attribute name override > default snake_case derivation.
	genLogger.Debug(fmt.Sprintf("Setting AttributeName for property '%s'.", p.PropertyName))

	if p.propertyDefinition.AttributeName != "" {
		p.AttributeName = p.propertyDefinition.AttributeName
	} else if override, ok := p.globalDefinition.AttributeNameOverrides[p.PropertyName]; ok {
		p.AttributeName = override
	} else {
		p.AttributeName = utils.Underscore(p.PropertyName)
	}

	genLogger.Debug(fmt.Sprintf("Successfully set AttributeName '%s' for property '%s'.", p.AttributeName, p.PropertyName))
}

func (p *Property) setComputed() {
	// Determine if the property is computed.
	// By default all properties are computed except required properties,
	// because optional properties can have server-side defaults and read-only properties are always computed.
	genLogger.Debug(fmt.Sprintf("Setting Computed for property '%s'.", p.PropertyName))

	p.Computed = !p.Required

	genLogger.Debug(fmt.Sprintf("Successfully set Computed '%t' for property '%s'.", p.Computed, p.PropertyName))
}

func (p *Property) setCustomType() {
	// Determine if the property has a custom type.
	genLogger.Debug(fmt.Sprintf("Setting CustomType for property '%s'.", p.PropertyName))
	genLogger.Debug(fmt.Sprintf("Successfully set CustomType for property '%s'.", p.PropertyName))
}

func (p *Property) setDeprecated() {
	// Determine if the property is deprecated.
	genLogger.Debug(fmt.Sprintf("Setting Deprecated for property '%s'.", p.PropertyName))
	genLogger.Debug(fmt.Sprintf("Successfully set Deprecated for property '%s'.", p.PropertyName))
}

func (p *Property) setDeprecatedVersions() {
	// Determine the APIC versions in which the property is deprecated.
	genLogger.Debug(fmt.Sprintf("Setting DeprecatedVersions for property '%s'.", p.PropertyName))
	genLogger.Debug(fmt.Sprintf("Successfully set DeprecatedVersions for property '%s'.", p.PropertyName))
}

func (p *Property) setDocumentation() {
	// Determine the documentation specific information for the property.
	genLogger.Debug(fmt.Sprintf("Setting Documentation for property '%s'.", p.PropertyName))
	genLogger.Debug(fmt.Sprintf("Successfully set Documentation for property '%s'.", p.PropertyName))
}

func (p *Property) setMigrationValues() {
	// Determine the migration values for the property.
	genLogger.Debug(fmt.Sprintf("Setting MigrationValues for property '%s'.", p.PropertyName))
	genLogger.Debug(fmt.Sprintf("Successfully set MigrationValues for property '%s'.", p.PropertyName))
}

func (p *Property) setOptional() {
	// Determine if the property is optional.
	// A property is optional when the definition restriction is "optional",
	// or when the meta file indicates isConfigurable is true and the property is not required and the restriction is not "read_only".
	genLogger.Debug(fmt.Sprintf("Setting Optional for property '%s'.", p.PropertyName))

	if p.propertyDefinition.Restriction == "optional" {
		p.Optional = true
	} else if p.propertyDefinition.Restriction != "read_only" && p.metaDetails["isConfigurable"] == true && !p.Required {
		p.Optional = true
	}

	genLogger.Debug(fmt.Sprintf("Successfully set Optional '%t' for property '%s'.", p.Optional, p.PropertyName))
}

func (p *Property) setPointsToClass() {
	// Determine the class to which the property points.
	genLogger.Debug(fmt.Sprintf("Setting PointsToClass for property '%s'.", p.PropertyName))
	genLogger.Debug(fmt.Sprintf("Successfully set PointsToClass for property '%s'.", p.PropertyName))
}

func (p *Property) setReadOnly() {
	// Determine if the property is read-only.
	// A property is read-only only when the definition restriction is "read_only".
	// This is used to include isConfigurable=false properties as read-only attributes in the schema.
	genLogger.Debug(fmt.Sprintf("Setting ReadOnly for property '%s'.", p.PropertyName))

	if p.propertyDefinition.Restriction == "read_only" {
		p.ReadOnly = true
	}

	genLogger.Debug(fmt.Sprintf("Successfully set ReadOnly '%t' for property '%s'.", p.ReadOnly, p.PropertyName))
}

func (p *Property) setRequired() {
	// Determine if the property is required.
	// A property is required when the definition restriction is "required",
	// or when the meta file indicates isConfigurable and isNaming are both true.
	genLogger.Debug(fmt.Sprintf("Setting Required for property '%s'.", p.PropertyName))

	if p.propertyDefinition.Restriction == "required" {
		p.Required = true
	} else if p.metaDetails["isConfigurable"] == true && p.metaDetails["isNaming"] == true {
		p.Required = true
	}

	genLogger.Debug(fmt.Sprintf("Successfully set Required '%t' for property '%s'.", p.Required, p.PropertyName))
}

func (p *Property) setSensitive() {
	// Determine if the property is sensitive.
	// A property is sensitive when the definition override is true,
	// or when the meta file indicates secure is true.
	genLogger.Debug(fmt.Sprintf("Setting Sensitive for property '%s'.", p.PropertyName))

	if p.propertyDefinition.Sensitive {
		p.Sensitive = true
	} else if p.metaDetails["secure"] == true {
		p.Sensitive = true
	}

	genLogger.Debug(fmt.Sprintf("Successfully set Sensitive '%t' for property '%s'.", p.Sensitive, p.PropertyName))
}

func (p *Property) setTestValues() {
	// Determine the test values for the property.
	genLogger.Debug(fmt.Sprintf("Setting TestValues for property '%s'.", p.PropertyName))
	genLogger.Debug(fmt.Sprintf("Successfully set TestValues for property '%s'.", p.PropertyName))
}

func (p *Property) setValidators() {
	// Determine the validators for the property.
	genLogger.Debug(fmt.Sprintf("Setting Validators for property '%s'.", p.PropertyName))
	genLogger.Debug(fmt.Sprintf("Successfully set Validators for property '%s'.", p.PropertyName))
}

func (p *Property) setValidValues() {
	// Determine the valid values for the property.
	genLogger.Debug(fmt.Sprintf("Setting ValidValues for property '%s'.", p.PropertyName))
	genLogger.Debug(fmt.Sprintf("Successfully set ValidValues for property '%s'.", p.PropertyName))
}

func (p *Property) setValueType() {
	// Determine the value type of the property.
	genLogger.Debug(fmt.Sprintf("Setting ValueType for property '%s'.", p.PropertyName))
	genLogger.Debug(fmt.Sprintf("Successfully set ValueType for property '%s'.", p.PropertyName))
}

func (p *Property) setVersions() {
	// Determine the supported APIC versions for the property.
	genLogger.Debug(fmt.Sprintf("Setting Versions for property '%s'.", p.PropertyName))
	genLogger.Debug(fmt.Sprintf("Successfully set Versions for property '%s'.", p.PropertyName))
}
