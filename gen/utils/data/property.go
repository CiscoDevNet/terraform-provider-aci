package data

import "fmt"

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
	// The APIC versions in which the property is deprecated.
	DeprecatedVersions []VersionRange
	// Documentation specific information for the property.
	Documentation PropertyDocumentation
	// Migration specific information for the property.
	// This is a map that contains the migration value details of the attribute for a specific schema version.
	MigrationValues map[int]MigrationValue
	// Indicates if a property is optional in the resource and datasource schemas.
	Optional bool
	// When a property that points to another class, this is the class to which the property points to.
	PointsToClass string
	// The name of the property in the APIC class.
	PropertyName string
	// Indicates if the property is read-only in the resource schemas.
	ReadOnly bool
	// Indicates if the property is required in the resource and datasource schemas.
	Required bool
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
	// The supported APIC versions for the property.
	// Each version range is separated by a comma, ex "4.2(7f)-4.2(7w),5.2(1g)-".
	// The first version is the minimum version and the second version is the maximum version.
	// A dash at the end of a range (ex. 4.2(7f)-) indicates that the class is supported from the first version to the latest version.
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

type MigrationValue struct {
	// Indicates if a property is computed in the legacy resource schema.
	Computed bool
	// The name of the property in the legacy resource schema.
	AttributeName string
	// Indicates if a property is optional in the legacy resource schema.
	Optional bool
	// Indicates if a property is required in the legacy resource schema.
	Required bool
	// The type of the legacy attribute.
	Type ValueTypeEnum
}

type TestValue struct {
	// The changed value of the property to be used in the test when a property is allowed to be changed without destruction of the resource.
	Changed []string
	// The initial value of the property to be used in the test.
	// This is set to the default value of the property in APIC when it is not a required value.
	Initial []string
}

type ValidValue struct {
	// The valid value of the property.
	Value    string
	Versions []VersionRange
}

type Validator struct {
	// If the property has a range of values, these are the minimum and maximum values of the range.
	Max float64
	Min float64
	// If the property has one or more regex statements it requires to match.
	RegexList []RegexStatement
}

type RegexStatementTypeEnum int

// The enumeration options of the RegexStatement Type.
const (
	// Include indicates that the value must match the regex statement.
	Include RegexStatementTypeEnum = iota + 1
	// TODO: re-evaluate the possible regex statements type options.
)

type RegexStatement struct {
	// The regex string.
	Regex string
	// The type of the regex statement.
	Type RegexStatementTypeEnum
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

func NewProperty(name string, details map[string]interface{}) *Property {
	genLogger.Trace(fmt.Sprintf("Creating new property struct for property: %s.", name))
	property := &Property{PropertyName: name}

	// TODO: add function to set AttributeName

	// TODO: add function to set Computed

	// TODO: add function to set CustomType

	// TODO: add placeholder function for Deprecated

	// TODO: add placeholder function for DeprecatedVersions

	// TODO: add function to set Documentation

	// TODO: add function to set MigrationValues

	// TODO: add function to set Optional

	// TODO: add function to set PointsToClass

	// TODO: add function to set ReadOnly

	// TODO: add function to set Required

	// TODO: add function to set TestValues

	// TODO: add function to set Validators

	// TODO: add function to set ValidValues

	// TODO: add function to set ValueType

	// TODO: add function to set Versions

	return property
}
