package data

import "fmt"

// PlatformTypeEnum represents the APIC platform type. The default value is Apic.
type PlatformTypeEnum int

const (
	// Apic indicates that the class is available on the on-premises version of APIC. This is the default value.
	Apic PlatformTypeEnum = iota
	// Both indicates that the class is available on both the on-premises and cloud versions of APIC.
	Both
	// Cloud indicates that the class is available on the cloud version of APIC.
	Cloud
)

func (p PlatformTypeEnum) String() string {
	switch p {
	case Both:
		return "both"
	case Cloud:
		return "cloud"
	default:
		return "apic"
	}
}

func (p *PlatformTypeEnum) UnmarshalText(text []byte) error {
	switch string(text) {
	case "", "apic":
		*p = Apic
	case "cloud":
		*p = Cloud
	case "both":
		*p = Both
	default:
		return fmt.Errorf("unknown platform_type %q (expected one of: apic, cloud, both)", string(text))
	}
	return nil
}

// RelationshipTypeEnum identifies how a relational class connects from one MO to another.
type RelationshipTypeEnum int

const (
	// UndefinedRelationshipType is the zero value: no relationship type assigned.
	UndefinedRelationshipType RelationshipTypeEnum = iota
	// Named indicates that the relationship is a named relation.
	Named
	// Explicit indicates that the relationship is an explicit relation.
	Explicit
)

func (r RelationshipTypeEnum) String() string {
	switch r {
	case Named:
		return "named"
	case Explicit:
		return "explicit"
	default:
		return ""
	}
}

func (r *RelationshipTypeEnum) UnmarshalText(text []byte) error {
	switch string(text) {
	case "":
		*r = UndefinedRelationshipType
	case "named":
		*r = Named
	case "explicit":
		*r = Explicit
	default:
		return fmt.Errorf("unknown relationship type %q (expected one of: named, explicit)", string(text))
	}
	return nil
}

// TestDependencyRoleEnum indicates how the dependency is consumed in HCL.
type TestDependencyRoleEnum int

const (
	// UndefinedRole is the zero value: no role assigned. Valid for nested
	// dependencies (depth > 0) which are pure prerequisites; invalid at depth 0
	// where Parent or Target is required.
	UndefinedRole TestDependencyRoleEnum = iota
	// Parent means the dependency provides the parent_dn attribute.
	Parent
	// Target means the dependency provides the target_dn attribute (relation classes).
	Target
)

func (r TestDependencyRoleEnum) String() string {
	switch r {
	case Parent:
		return "parent"
	case Target:
		return "target"
	default:
		return ""
	}
}

func (r *TestDependencyRoleEnum) UnmarshalText(text []byte) error {
	switch string(text) {
	case "":
		*r = UndefinedRole
	case "parent":
		*r = Parent
	case "target":
		*r = Target
	default:
		return fmt.Errorf("unknown role %q (expected one of: parent, target)", string(text))
	}
	return nil
}

// ReferenceTypeEnum indicates how to interpret a TestDependency.Reference value.
type ReferenceTypeEnum int

const (
	// UndefinedReferenceType is the zero value: no type assigned.
	UndefinedReferenceType ReferenceTypeEnum = iota
	// StaticReference is a hardcoded DN string (e.g. "uni/vmmp-VMware/dom-domain_1").
	StaticReference
	// ResourceReference is a Terraform resource attribute path (e.g. "aci_tenant.test.id").
	// This is the default when reference_type is omitted.
	ResourceReference
	// DataSourceReference is a Terraform data source attribute path (e.g. "data.aci_tenant.test.id").
	DataSourceReference
)

func (r ReferenceTypeEnum) String() string {
	switch r {
	case StaticReference:
		return "static"
	case ResourceReference:
		return "resource"
	case DataSourceReference:
		return "data_source"
	default:
		return ""
	}
}

func (r *ReferenceTypeEnum) UnmarshalText(text []byte) error {
	switch string(text) {
	case "":
		// Documented default when reference_type is omitted.
		*r = ResourceReference
	case "static":
		*r = StaticReference
	case "resource":
		*r = ResourceReference
	case "data_source":
		*r = DataSourceReference
	default:
		return fmt.Errorf("unknown reference_type %q (expected one of: static, resource, data_source)", string(text))
	}
	return nil
}

// RegexStatementTypeEnum identifies the kind of regex match constraint in a Validator.
type RegexStatementTypeEnum int

const (
	// Include indicates that the value must match the regex statement.
	Include RegexStatementTypeEnum = iota + 1
)

func (r RegexStatementTypeEnum) String() string {
	switch r {
	case Include:
		return "include"
	default:
		return ""
	}
}

func (r *RegexStatementTypeEnum) UnmarshalText(text []byte) error {
	switch string(text) {
	case "include":
		*r = Include
	default:
		return fmt.Errorf("unknown regex statement type %q (expected one of: include)", string(text))
	}
	return nil
}

// ValueRenderTypeEnum controls how a TestValueEntry is rendered in HCL configuration.
type ValueRenderTypeEnum int

const (
	// StringValue renders as a quoted string: attribute = "value". Default when value_type is omitted.
	StringValue ValueRenderTypeEnum = iota + 1
	// ReferenceValue renders as an unquoted reference expression: attribute = aci_tenant.test.id
	ReferenceValue
)

func (v ValueRenderTypeEnum) String() string {
	switch v {
	case StringValue:
		return "string"
	case ReferenceValue:
		return "reference"
	default:
		return ""
	}
}

func (v *ValueRenderTypeEnum) UnmarshalText(text []byte) error {
	switch string(text) {
	case "", "string":
		*v = StringValue
	case "reference":
		*v = ReferenceValue
	default:
		return fmt.Errorf("unknown value_type %q (expected one of: string, reference)", string(text))
	}
	return nil
}

// ValueTypeEnum identifies the data shape of a property and is the single dispatch point
// for type-specific schema and template behavior. New named custom types should be added
// here as constants and registered in UnmarshalText so the same definition `value_type`
// override key can select them.
type ValueTypeEnum int

const (
	// String indicates that the property is a plain string value.
	String ValueTypeEnum = iota + 1
	// Set indicates that the property is a set value (driven by meta uitype "bitmask").
	Set
	// IpAddress indicates that the property is an IP address (IPv4 or IPv6); driven by
	// meta `validateAsIPv4OrIPv6`. Renders with the IP-address custom type for parsing,
	// validation, and semantic-equality (e.g. zero-padding normalization).
	IpAddress
	// SemanticEquality indicates that the property has both ValidValues and Validators,
	// meaning the wire form (e.g. "22") and the human form (e.g. "ssh") must compare equal.
	// Templates render this with the semantic-equality custom type.
	SemanticEquality
)

func (v ValueTypeEnum) String() string {
	switch v {
	case String:
		return "string"
	case Set:
		return "set"
	case IpAddress:
		return "ip_address"
	case SemanticEquality:
		return "semantic_equality"
	default:
		return ""
	}
}

func (v *ValueTypeEnum) UnmarshalText(text []byte) error {
	switch string(text) {
	case "":
		*v = 0
	case "string":
		*v = String
	case "set":
		*v = Set
	case "ip_address":
		*v = IpAddress
	case "semantic_equality":
		*v = SemanticEquality
	default:
		return fmt.Errorf("unknown value_type %q (expected one of: string, set, ip_address, semantic_equality)", string(text))
	}
	return nil
}

// RestrictionEnum controls the schema behavior of a property as declared by the class definition.
type RestrictionEnum int

const (
	// UndefinedRestriction is the zero value: derive behavior from the meta file.
	UndefinedRestriction RestrictionEnum = iota
	// Required marks the property as required in the resource and datasource schemas.
	Required
	// Optional marks the property as optional in the resource and datasource schemas.
	Optional
	// ReadOnly includes a non-configurable property as a computed-only attribute.
	ReadOnly
	// Exclude omits the property entirely from generated schemas.
	Exclude
)

func (r RestrictionEnum) String() string {
	switch r {
	case Required:
		return "required"
	case Optional:
		return "optional"
	case ReadOnly:
		return "read_only"
	case Exclude:
		return "exclude"
	default:
		return ""
	}
}

func (r *RestrictionEnum) UnmarshalText(text []byte) error {
	switch string(text) {
	case "":
		*r = UndefinedRestriction
	case "required":
		*r = Required
	case "optional":
		*r = Optional
	case "read_only":
		*r = ReadOnly
	case "exclude":
		*r = Exclude
	default:
		return fmt.Errorf("unknown restriction %q (expected one of: required, optional, read_only, exclude)", string(text))
	}
	return nil
}
