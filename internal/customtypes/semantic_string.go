package customtypes

import (
	"context"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.StringTypable = SemanticStringType{}

type semanticStringDefinition struct {
	valueMap map[string]string
}

// SemanticStringType is a string type whose wire and local-name values compare
// equally according to a property-owned value map.
type SemanticStringType struct {
	basetypes.StringType
	definition *semanticStringDefinition
}

// NewSemanticStringType creates an immutable semantic string type from a map
// of APIC wire values to their canonical local names.
func NewSemanticStringType(valueMap map[string]string) SemanticStringType {
	return SemanticStringType{
		definition: &semanticStringDefinition{valueMap: maps.Clone(valueMap)},
	}
}

func (t SemanticStringType) Equal(other attr.Type) bool {
	otherType, ok := other.(SemanticStringType)
	return ok && semanticStringDefinitionsEqual(t.definition, otherType.definition)
}

func (t SemanticStringType) String() string {
	return "customtypes.SemanticStringType"
}

func (t SemanticStringType) ValueFromString(_ context.Context, value basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return SemanticStringValue{StringValue: value, definition: t.definition}, nil
}

func (t SemanticStringType) ValueFromTerraform(ctx context.Context, value tftypes.Value) (attr.Value, error) {
	attributeValue, err := t.StringType.ValueFromTerraform(ctx, value)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attributeValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attributeValue)
	}

	return SemanticStringValue{StringValue: stringValue, definition: t.definition}, nil
}

func (t SemanticStringType) ValueType(context.Context) attr.Value {
	return SemanticStringValue{definition: t.definition}
}

// Null returns a null value associated with this configured semantic type.
func (t SemanticStringType) Null() SemanticStringValue {
	return SemanticStringValue{StringValue: basetypes.NewStringNull(), definition: t.definition}
}

// Unknown returns an unknown value associated with this configured semantic type.
func (t SemanticStringType) Unknown() SemanticStringValue {
	return SemanticStringValue{StringValue: basetypes.NewStringUnknown(), definition: t.definition}
}

// Value returns a known value associated with this configured semantic type.
func (t SemanticStringType) Value(value string) SemanticStringValue {
	return SemanticStringValue{StringValue: basetypes.NewStringValue(value), definition: t.definition}
}

var _ basetypes.StringValuableWithSemanticEquals = SemanticStringValue{}

// SemanticStringValue retains its property mapping so values created through
// schemas, null initialization, and APIC response decoding behave identically.
type SemanticStringValue struct {
	basetypes.StringValue
	definition *semanticStringDefinition
}

func (v SemanticStringValue) Equal(other attr.Value) bool {
	otherValue, ok := other.(SemanticStringValue)
	return ok &&
		semanticStringDefinitionsEqual(v.definition, otherValue.definition) &&
		v.StringValue.Equal(otherValue.StringValue)
}

func (v SemanticStringValue) Type(context.Context) attr.Type {
	return SemanticStringType{definition: v.definition}
}

func (v SemanticStringValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	newValue, ok := newValuable.(SemanticStringValue)
	if !ok {
		return false, unexpectedSemanticStringValueDiagnostics(v, newValuable)
	}
	if !semanticStringDefinitionsEqual(v.definition, newValue.definition) {
		return false, diag.Diagnostics{diag.NewErrorDiagnostic(
			"Semantic Equality Check Error",
			"Semantic string values with different value maps cannot be compared.",
		)}
	}

	return v.CanonicalString() == newValue.CanonicalString(), nil
}

// CanonicalString returns the local name for a mapped APIC wire value and the
// raw Terraform value when no mapping exists.
func (v SemanticStringValue) CanonicalString() string {
	if v.definition != nil {
		if canonicalValue, ok := v.definition.valueMap[v.ValueString()]; ok {
			return canonicalValue
		}
	}
	return v.ValueString()
}

func semanticStringDefinitionsEqual(left, right *semanticStringDefinition) bool {
	if left == right {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	return maps.Equal(left.valueMap, right.valueMap)
}

func unexpectedSemanticStringValueDiagnostics(expected, actual any) diag.Diagnostics {
	return diag.Diagnostics{diag.NewErrorDiagnostic(
		"Semantic Equality Check Error",
		"An unexpected value type was received while performing semantic equality checks. "+
			"Please report this to the provider developers.\n\n"+
			fmt.Sprintf("Expected Value Type: %T\nGot Value Type: %T", expected, actual),
	)}
}
