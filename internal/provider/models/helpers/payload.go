package helpers

import (
	"context"
	"encoding/json"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// BuildSingletonChildPayload constructs the create, update, or delete payload for a singleton child.
func BuildSingletonChildPayload[T any](
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	desiredValue types.Object,
	priorValue *types.Object,
	childClassName string,
	cannotDeleteSummary string,
	defaultAnnotation string,
	build func(*T, context.Context, *T, bool, string) (map[string]any, diag.Diagnostics),
	buildDelete func(*T, context.Context, *diag.Diagnostics) map[string]any,
) ([]map[string]any, bool) {
	if desiredValue.IsNull() || desiredValue.IsUnknown() {
		return nil, true
	}

	var priorModel *T
	if priorValue != nil &&
		!priorValue.IsNull() &&
		!priorValue.IsUnknown() &&
		!objectIsEmpty(*priorValue) {
		priorModel = new(T)
		priorDiagnostics := priorValue.As(ctx, priorModel, basetypes.ObjectAsOptions{})
		diagnostics.Append(priorDiagnostics...)
		if priorDiagnostics.HasError() {
			return nil, false
		}
	}

	if !objectIsEmpty(desiredValue) {
		desiredModel := new(T)
		desiredDiagnostics := desiredValue.As(ctx, desiredModel, basetypes.ObjectAsOptions{})
		diagnostics.Append(desiredDiagnostics...)
		if desiredDiagnostics.HasError() {
			return nil, false
		}

		payloadObject, childDiagnostics := build(desiredModel, ctx, priorModel, true, defaultAnnotation)
		diagnostics.Append(childDiagnostics...)
		if childDiagnostics.HasError() {
			return nil, false
		}

		return []map[string]any{{childClassName: payloadObject}}, true
	}

	if priorModel == nil {
		return nil, true
	}
	if buildDelete == nil {
		diagnostics.AddError(
			cannotDeleteSummary,
			"deletion of child is only possible upon deletion of the parent",
		)
		return nil, false
	}

	payloadObject := buildDelete(priorModel, ctx, diagnostics)
	return []map[string]any{{childClassName: payloadObject}}, !diagnostics.HasError()
}

// BuildRepeatedChildPayloads reconciles repeated children by RN and constructs their payloads.
func BuildRepeatedChildPayloads[T any](
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	desiredValue types.Set,
	priorValue *types.Set,
	childClassName string,
	cannotDeleteSummary string,
	defaultAnnotation string,
	buildRN func(*T) string,
	build func(*T, context.Context, *T, bool, string) (map[string]any, diag.Diagnostics),
	buildDelete func(*T, context.Context, *diag.Diagnostics) map[string]any,
) ([]map[string]any, bool) {
	if desiredValue.IsNull() || desiredValue.IsUnknown() {
		return nil, true
	}

	var desiredModels []T
	desiredDiagnostics := desiredValue.ElementsAs(ctx, &desiredModels, false)
	diagnostics.Append(desiredDiagnostics...)
	if desiredDiagnostics.HasError() {
		return nil, false
	}

	var priorModels []T
	if priorValue != nil && !priorValue.IsNull() && !priorValue.IsUnknown() {
		priorDiagnostics := priorValue.ElementsAs(ctx, &priorModels, false)
		diagnostics.Append(priorDiagnostics...)
		if priorDiagnostics.HasError() {
			return nil, false
		}
	}

	payloads := make([]map[string]any, 0, len(desiredModels)+len(priorModels))
	for desiredIndex := range desiredModels {
		desiredModel := &desiredModels[desiredIndex]
		var priorModel *T
		for priorIndex := range priorModels {
			if buildRN(&priorModels[priorIndex]) == buildRN(desiredModel) {
				priorModel = &priorModels[priorIndex]
				break
			}
		}

		payloadObject, childDiagnostics := build(desiredModel, ctx, priorModel, true, defaultAnnotation)
		diagnostics.Append(childDiagnostics...)
		if childDiagnostics.HasError() {
			return nil, false
		}
		payloads = append(payloads, map[string]any{childClassName: payloadObject})
	}

	success := true
	for priorIndex := range priorModels {
		priorModel := &priorModels[priorIndex]
		found := false
		for desiredIndex := range desiredModels {
			if buildRN(&desiredModels[desiredIndex]) == buildRN(priorModel) {
				found = true
				break
			}
		}
		if found {
			continue
		}

		if buildDelete == nil {
			diagnostics.AddError(
				cannotDeleteSummary,
				"deletion of child is only possible upon deletion of the parent",
			)
			success = false
			continue
		}
		payloadObject := buildDelete(priorModel, ctx, diagnostics)
		payloads = append(payloads, map[string]any{childClassName: payloadObject})
		if diagnostics.HasError() {
			success = false
		}
	}

	return payloads, success
}

// NewPayloadObject constructs the attributes and optional children body of an APIC class payload.
func NewPayloadObject(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	attributes map[string]any,
	children []map[string]any,
	includeChildren bool,
) map[string]any {
	payloadObject := map[string]any{"attributes": attributes}
	if includeChildren {
		payloadObject["children"] = children
	}

	return payloadObject
}

// NewNestedDeletePayloadObject marks a nested APIC object for deletion.
func NewNestedDeletePayloadObject(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	attributes map[string]any,
) map[string]any {
	attributes["status"] = "deleted"
	return NewPayloadObject(ctx, diagnostics, attributes, []map[string]any{}, true)
}

// NewPayloadContainer converts a payload map to the container accepted by the APIC client.
func NewPayloadContainer(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	payload map[string]any,
	description string,
) *container.Container {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		diagnostics.AddError(
			"Marshalling of "+description+" failed",
			err.Error()+". Please report this issue to the provider developers.",
		)
		return nil
	}

	jsonPayload, err := container.ParseJSON(payloadBytes)
	if err != nil {
		diagnostics.AddError(
			"Construction of "+description+" failed",
			err.Error()+". Please report this issue to the provider developers.",
		)
		return nil
	}

	return jsonPayload
}

// NewDeletePayload constructs a top-level APIC delete payload for a class DN.
func NewDeletePayload(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	className string,
	dn string,
) *container.Container {
	return NewPayloadContainer(
		ctx,
		diagnostics,
		map[string]any{
			className: map[string]any{
				"attributes": map[string]any{
					"dn":     dn,
					"status": "deleted",
				},
			},
		},
		"JSON delete payload",
	)
}
