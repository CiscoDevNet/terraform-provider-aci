package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ReconciliationSchema contains only the class information needed to reconcile
// configured, planned, and prior-state child values.
type ReconciliationSchema struct {
	Identifiers []string
	Children    map[string]ReconciliationChild
}

// ReconciliationChild identifies the shape and class schema of a child value.
type ReconciliationChild struct {
	Repeated bool
	Schema   ReconciliationSchema
}

type reconciliationMode uint8

const (
	singletonReconciliation reconciliationMode = iota
	repeatedReconciliation
)

func (s ReconciliationSchema) reconcileObject(
	ctx context.Context,
	configValue types.Object,
	planValue types.Object,
	stateValue types.Object,
	mode reconciliationMode,
) (types.Object, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	configAttributes := configValue.Attributes()
	planAttributes := planValue.Attributes()
	stateAttributes := stateValue.Attributes()

	for name, planAttribute := range planAttributes {
		configAttribute, configAttributeExists := configAttributes[name]
		stateAttribute, stateAttributeExists := stateAttributes[name]
		if !configAttributeExists || !stateAttributeExists {
			diagnostics.AddError(
				"Unable to Reconcile Nested Child",
				"Expected attribute "+name+" in the configured, planned, and prior-state child values.",
			)
			return planValue, diagnostics
		}

		child, isChild := s.Children[name]
		if !isChild {
			planAttributes[name] = reconcileProperty(configAttribute, planAttribute, stateAttribute, mode)
			continue
		}

		childValue, childDiagnostics := child.reconcile(
			ctx,
			name,
			configAttribute,
			planAttribute,
			stateAttribute,
		)
		diagnostics.Append(childDiagnostics...)
		if childDiagnostics.HasError() {
			return planValue, diagnostics
		}
		planAttributes[name] = childValue
	}

	updatedPlan, objectDiagnostics := types.ObjectValue(planValue.AttributeTypes(ctx), planAttributes)
	diagnostics.Append(objectDiagnostics...)
	if objectDiagnostics.HasError() {
		return planValue, diagnostics
	}

	return updatedPlan, diagnostics
}

func reconcileProperty(configValue, planValue, stateValue attr.Value, mode reconciliationMode) attr.Value {
	if mode == repeatedReconciliation && stateValue.IsNull() {
		return stateValue
	}

	if planValue.IsUnknown() && configValue.IsNull() && !stateValue.IsNull() && !stateValue.IsUnknown() {
		return stateValue
	}

	return planValue
}

func (c ReconciliationChild) reconcile(
	ctx context.Context,
	name string,
	configValue attr.Value,
	planValue attr.Value,
	stateValue attr.Value,
) (attr.Value, diag.Diagnostics) {
	if c.Repeated {
		configSet, configIsSet := configValue.(types.Set)
		planSet, planIsSet := planValue.(types.Set)
		stateSet, stateIsSet := stateValue.(types.Set)
		if !configIsSet || !planIsSet || !stateIsSet {
			return planValue, unexpectedChildTypeDiagnostics(name, "set")
		}

		return c.Schema.reconcileRepeatedValue(ctx, configSet, planSet, stateSet, true)
	}

	configObject, configIsObject := configValue.(types.Object)
	planObject, planIsObject := planValue.(types.Object)
	stateObject, stateIsObject := stateValue.(types.Object)
	if !configIsObject || !planIsObject || !stateIsObject {
		return planValue, unexpectedChildTypeDiagnostics(name, "object")
	}

	return c.Schema.reconcileSingletonValue(ctx, configObject, planObject, stateObject, true)
}

func (s ReconciliationSchema) reconcileSingletonValue(
	ctx context.Context,
	configValue types.Object,
	planValue types.Object,
	stateValue types.Object,
	resourceExists bool,
) (types.Object, diag.Diagnostics) {
	if !resourceExists {
		return planValue, nil
	}

	// An explicitly configured empty object represents child deletion and must
	// win even when Terraform marks the complete planned object as unknown.
	if !configValue.IsNull() && !configValue.IsUnknown() && objectIsEmpty(configValue) {
		return configValue, nil
	}

	if planValue.IsUnknown() {
		if stateValue.IsNull() || stateValue.IsUnknown() || configValue.IsUnknown() {
			return planValue, nil
		}
		return stateValue, nil
	}

	if planValue.IsNull() || configValue.IsNull() || configValue.IsUnknown() {
		return planValue, nil
	}

	if stateValue.IsNull() || stateValue.IsUnknown() {
		return planValue, nil
	}

	return s.reconcileObject(ctx, configValue, planValue, stateValue, singletonReconciliation)
}

func (s ReconciliationSchema) reconcileRepeatedValue(
	ctx context.Context,
	configValue types.Set,
	planValue types.Set,
	stateValue types.Set,
	resourceExists bool,
) (types.Set, diag.Diagnostics) {
	if !resourceExists {
		return planValue, nil
	}

	if planValue.IsUnknown() {
		if stateValue.IsUnknown() {
			return planValue, nil
		}
		if stateValue.IsNull() || !configValue.IsUnknown() {
			return stateValue, nil
		}
		return planValue, nil
	}

	if planValue.IsNull() || configValue.IsNull() || configValue.IsUnknown() || stateValue.IsNull() || stateValue.IsUnknown() {
		return planValue, nil
	}

	if len(s.Identifiers) == 0 {
		return planValue, nil
	}

	return s.reconcileSetElements(ctx, configValue, planValue, stateValue)
}

func (s ReconciliationSchema) reconcileSetElements(
	ctx context.Context,
	configValue types.Set,
	planValue types.Set,
	stateValue types.Set,
) (types.Set, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	configObjects, configDiagnostics := objectElements(configValue.Elements(), "configured")
	diagnostics.Append(configDiagnostics...)
	planObjects, planDiagnostics := objectElements(planValue.Elements(), "planned")
	diagnostics.Append(planDiagnostics...)
	stateObjects, stateDiagnostics := objectElements(stateValue.Elements(), "prior-state")
	diagnostics.Append(stateDiagnostics...)
	if diagnostics.HasError() {
		return planValue, diagnostics
	}

	configMatched := make([]bool, len(configObjects))
	stateMatched := make([]bool, len(stateObjects))
	planElements := planValue.Elements()
	for index, planObject := range planObjects {
		configObject, configFound, matchDiagnostics := s.matchingObject(ctx, planObject, configObjects, configMatched)
		diagnostics.Append(matchDiagnostics...)
		if matchDiagnostics.HasError() {
			return planValue, diagnostics
		}
		if !configFound {
			continue
		}

		stateObject, stateFound, matchDiagnostics := s.matchingObject(ctx, planObject, stateObjects, stateMatched)
		diagnostics.Append(matchDiagnostics...)
		if matchDiagnostics.HasError() {
			return planValue, diagnostics
		}
		if !stateFound {
			continue
		}

		updatedObject, objectDiagnostics := s.reconcileObject(
			ctx,
			configObject,
			planObject,
			stateObject,
			repeatedReconciliation,
		)
		diagnostics.Append(objectDiagnostics...)
		if objectDiagnostics.HasError() {
			return planValue, diagnostics
		}
		planElements[index] = updatedObject
	}

	updatedPlan, setDiagnostics := types.SetValue(planValue.ElementType(ctx), planElements)
	diagnostics.Append(setDiagnostics...)
	if setDiagnostics.HasError() {
		return planValue, diagnostics
	}

	return updatedPlan, diagnostics
}

func objectElements(elements []attr.Value, description string) ([]types.Object, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	objects := make([]types.Object, len(elements))
	for index, element := range elements {
		object, ok := element.(types.Object)
		if !ok {
			diagnostics.AddError(
				"Unable to Reconcile Nested Child",
				"Expected every "+description+" repeated child element to be a Terraform object value.",
			)
			return nil, diagnostics
		}
		objects[index] = object
	}
	return objects, diagnostics
}

func (s ReconciliationSchema) matchingObject(
	ctx context.Context,
	target types.Object,
	candidates []types.Object,
	matched []bool,
) (types.Object, bool, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	for index, candidate := range candidates {
		if matched[index] {
			continue
		}

		equal, identifierDiagnostics := s.identifiersEqual(ctx, target, candidate)
		diagnostics.Append(identifierDiagnostics...)
		if identifierDiagnostics.HasError() {
			return types.Object{}, false, diagnostics
		}
		if equal {
			matched[index] = true
			return candidate, true, diagnostics
		}
	}

	return types.Object{}, false, diagnostics
}

func (s ReconciliationSchema) identifiersEqual(
	ctx context.Context,
	left types.Object,
	right types.Object,
) (bool, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	leftAttributes := left.Attributes()
	rightAttributes := right.Attributes()
	for _, identifier := range s.Identifiers {
		leftValue, leftExists := leftAttributes[identifier]
		rightValue, rightExists := rightAttributes[identifier]
		if !leftExists || !rightExists {
			diagnostics.AddError(
				"Unable to Reconcile Nested Child",
				"Expected identifier "+identifier+" in both repeated child values.",
			)
			return false, diagnostics
		}

		if leftValue.IsNull() || leftValue.IsUnknown() || rightValue.IsNull() || rightValue.IsUnknown() {
			return false, diagnostics
		}

		equal, valueDiagnostics := identifierValuesEqual(ctx, leftValue, rightValue)
		diagnostics.Append(valueDiagnostics...)
		if valueDiagnostics.HasError() || !equal {
			return false, diagnostics
		}
	}

	return true, diagnostics
}

func identifierValuesEqual(ctx context.Context, left, right attr.Value) (bool, diag.Diagnostics) {
	if left.Equal(right) {
		return true, nil
	}

	leftSemantic, leftSupportsSemanticEquality := left.(basetypes.StringValuableWithSemanticEquals)
	rightString, rightIsString := right.(basetypes.StringValuable)
	if leftSupportsSemanticEquality && rightIsString {
		return leftSemantic.StringSemanticEquals(ctx, rightString)
	}

	rightSemantic, rightSupportsSemanticEquality := right.(basetypes.StringValuableWithSemanticEquals)
	leftString, leftIsString := left.(basetypes.StringValuable)
	if rightSupportsSemanticEquality && leftIsString {
		return rightSemantic.StringSemanticEquals(ctx, leftString)
	}

	return false, nil
}

func objectIsEmpty(value types.Object) bool {
	for _, attribute := range value.Attributes() {
		if !attribute.IsNull() {
			return false
		}
	}
	return true
}

func unexpectedChildTypeDiagnostics(name, expectedType string) diag.Diagnostics {
	return diag.Diagnostics{
		diag.NewErrorDiagnostic(
			"Unable to Reconcile Nested Child",
			"Expected child "+name+" to be represented by Terraform "+expectedType+" values.",
		),
	}
}
