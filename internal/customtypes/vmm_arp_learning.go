package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.StringTypable = VMMArpLearningType{}

// VMMArpLearningType preserves the APIC equivalence between an empty value and disabled.
type VMMArpLearningType struct {
	basetypes.StringType
}

func (t VMMArpLearningType) Equal(other attr.Type) bool {
	_, ok := other.(VMMArpLearningType)
	return ok
}

func (t VMMArpLearningType) String() string {
	return "customtypes.VMMArpLearningType"
}

func (t VMMArpLearningType) ValueFromString(_ context.Context, value basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return VMMArpLearningValue{StringValue: value}, nil
}

func (t VMMArpLearningType) ValueFromTerraform(ctx context.Context, value tftypes.Value) (attr.Value, error) {
	attributeValue, err := t.StringType.ValueFromTerraform(ctx, value)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attributeValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attributeValue)
	}

	return VMMArpLearningValue{StringValue: stringValue}, nil
}

func (t VMMArpLearningType) ValueType(context.Context) attr.Value {
	return VMMArpLearningValue{}
}

var _ basetypes.StringValuableWithSemanticEquals = VMMArpLearningValue{}

// VMMArpLearningValue preserves its raw value while comparing empty and disabled equally.
type VMMArpLearningValue struct {
	basetypes.StringValue
}

func (v VMMArpLearningValue) Equal(other attr.Value) bool {
	otherValue, ok := other.(VMMArpLearningValue)
	return ok && v.StringValue.Equal(otherValue.StringValue)
}

func (v VMMArpLearningValue) Type(context.Context) attr.Type {
	return VMMArpLearningType{}
}

func (v VMMArpLearningValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	newValue, ok := newValuable.(VMMArpLearningValue)
	if !ok {
		return false, unexpectedSemanticStringValueDiagnostics(v, newValuable)
	}

	return v.CanonicalString() == newValue.CanonicalString(), nil
}

// CanonicalString returns disabled for the APIC empty representation.
func (v VMMArpLearningValue) CanonicalString() string {
	if v.ValueString() == "" {
		return "disabled"
	}
	return v.ValueString()
}

func NewVMMArpLearningNull() VMMArpLearningValue {
	return VMMArpLearningValue{StringValue: basetypes.NewStringNull()}
}

func NewVMMArpLearningUnknown() VMMArpLearningValue {
	return VMMArpLearningValue{StringValue: basetypes.NewStringUnknown()}
}

func NewVMMArpLearningValue(value string) VMMArpLearningValue {
	return VMMArpLearningValue{StringValue: basetypes.NewStringValue(value)}
}
