package helpers

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Attributes is a validated APIC attributes object together with its class context.
type Attributes struct {
	className string
	values    map[string]*container.Container
}

// AttributesFromObject returns the validated attributes from an APIC class object.
func AttributesFromObject(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	object *container.Container,
	className string,
) (Attributes, bool) {
	attributes := Attributes{className: className}
	if object == nil {
		diagnostics.AddError(
			fmt.Sprintf("Malformed APIC %s object", className),
			fmt.Sprintf("APIC returned a nil %s object.", className),
		)
		return attributes, false
	}

	attributesContainer := object.Search("attributes")
	if attributesContainer == nil {
		diagnostics.AddError(
			fmt.Sprintf("Malformed APIC %s object", className),
			fmt.Sprintf("APIC object %s does not contain an attributes object.", className),
		)
		return attributes, false
	}

	values, err := attributesContainer.ChildrenMap()
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("Malformed APIC %s attributes", className),
			fmt.Sprintf("APIC attributes for class %q must be an object: %v.", className, err),
		)
		return attributes, false
	}

	attributes.values = values
	return attributes, true
}

// ChildObjectsByClass groups direct APIC children by the known class names requested by a model.
func ChildObjectsByClass(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	object *container.Container,
	parentClassName string,
	childClassNames []string,
) map[string][]*container.Container {
	objectsByClass := make(map[string][]*container.Container, len(childClassNames))

	childrenContainer := object.Search("children")
	if childrenContainer == nil {
		return objectsByClass
	}
	if _, ok := childrenContainer.Data().([]any); !ok {
		diagnostics.AddError(
			fmt.Sprintf("Malformed APIC %s children", parentClassName),
			fmt.Sprintf("APIC children for class %q must be an array, got %T.", parentClassName, childrenContainer.Data()),
		)
		return objectsByClass
	}

	children, err := childrenContainer.Children()
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("Malformed APIC %s children", parentClassName),
			fmt.Sprintf("APIC children for class %q could not be read: %v.", parentClassName, err),
		)
		return objectsByClass
	}

	knownClasses := make(map[string]struct{}, len(childClassNames))
	for _, className := range childClassNames {
		knownClasses[className] = struct{}{}
	}

	for childIndex, childEnvelope := range children {
		childClasses, childErr := childEnvelope.ChildrenMap()
		if childErr != nil {
			diagnostics.AddError(
				"Malformed APIC child envelope",
				fmt.Sprintf("APIC child %d of class %q must be an object: %v.", childIndex, parentClassName, childErr),
			)
			continue
		}

		for className, childObject := range childClasses {
			if _, known := knownClasses[className]; known {
				objectsByClass[className] = append(objectsByClass[className], childObject)
			}
		}
	}

	return objectsByClass
}

// DecodeSingletonChild converts at most one APIC child object to a Terraform object value.
func DecodeSingletonChild[T any](
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	objects []*container.Container,
	fallbackValue *types.Object,
	attributeTypes map[string]attr.Type,
	parentClassName string,
	childClassName string,
	newNullModel func() T,
	decode func(context.Context, *container.Container, *T) (T, diag.Diagnostics),
	target *types.Object,
) {
	var fallbackModel *T
	if fallbackValue != nil &&
		!fallbackValue.IsNull() &&
		!fallbackValue.IsUnknown() &&
		!objectIsEmpty(*fallbackValue) {
		fallback := newNullModel()
		fallbackDiagnostics := fallbackValue.As(ctx, &fallback, basetypes.ObjectAsOptions{})
		diagnostics.Append(fallbackDiagnostics...)
		if !fallbackDiagnostics.HasError() {
			fallbackModel = &fallback
		}
	}

	model := newNullModel()
	if len(objects) > 1 {
		tflog.Warn(ctx, "APIC returned multiple objects for a singleton child class", map[string]any{
			"child_class":  childClassName,
			"parent_class": parentClassName,
			"count":        len(objects),
		})
		diagnostics.AddWarning(
			"APIC model invariant violation",
			fmt.Sprintf(
				"APIC returned %d %q children below %q, but the model permits one; only the first object was decoded.",
				len(objects),
				childClassName,
				parentClassName,
			),
		)
	}
	if len(objects) > 0 {
		decodedModel, childDiagnostics := decode(ctx, objects[0], fallbackModel)
		diagnostics.Append(childDiagnostics...)
		model = decodedModel
	}

	decodedValue, valueDiagnostics := types.ObjectValueFrom(ctx, attributeTypes, model)
	diagnostics.Append(valueDiagnostics...)
	if !valueDiagnostics.HasError() {
		*target = decodedValue
	}
}

// DecodeRepeatedChildren converts APIC child objects to a Terraform set value.
func DecodeRepeatedChildren[T any](
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	objects []*container.Container,
	attributeTypes map[string]attr.Type,
	decode func(context.Context, *container.Container, *T) (T, diag.Diagnostics),
	target *types.Set,
) {
	models := make([]T, 0, len(objects))
	for _, object := range objects {
		model, childDiagnostics := decode(ctx, object, nil)
		diagnostics.Append(childDiagnostics...)
		models = append(models, model)
	}

	decodedValue, valueDiagnostics := types.SetValueFrom(
		ctx,
		types.ObjectType{AttrTypes: attributeTypes},
		models,
	)
	diagnostics.Append(valueDiagnostics...)
	if !valueDiagnostics.HasError() {
		*target = decodedValue
	}
}

// ModelFromResponse resolves one expected APIC class object and delegates its concrete model decoding.
func ModelFromResponse[T any](
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	response *container.Container,
	className string,
	fallbackModel *T,
	decode func(context.Context, *container.Container, *T) (T, diag.Diagnostics),
) (*T, string) {
	if response == nil {
		diagnostics.AddError(
			"Malformed APIC response",
			fmt.Sprintf("APIC returned a nil response while reading class %s.", className),
		)
		return nil, ""
	}

	imdataContainer := response.Search("imdata")
	if imdataContainer == nil {
		diagnostics.AddError(
			"Malformed APIC response",
			fmt.Sprintf("APIC response for class %s does not contain imdata.", className),
		)
		return nil, ""
	}
	if _, ok := imdataContainer.Data().([]any); !ok {
		diagnostics.AddError(
			"Malformed APIC response",
			fmt.Sprintf("APIC imdata for class %q must be an array, got %T.", className, imdataContainer.Data()),
		)
		return nil, ""
	}

	imdata, err := imdataContainer.Children()
	if err != nil {
		diagnostics.AddError(
			"Malformed APIC response",
			fmt.Sprintf("APIC imdata for class %q could not be read: %v.", className, err),
		)
		return nil, ""
	}

	expectedObjects := make([]*container.Container, 0, 1)
	for objectIndex, objectEnvelope := range imdata {
		classes, objectErr := objectEnvelope.ChildrenMap()
		if objectErr != nil {
			diagnostics.AddError(
				"Malformed APIC object envelope",
				fmt.Sprintf("APIC imdata object %d must be an object: %v.", objectIndex, objectErr),
			)
			continue
		}
		if object, ok := classes[className]; ok {
			expectedObjects = append(expectedObjects, object)
		}
	}

	if len(expectedObjects) == 0 {
		return nil, ""
	}
	if len(expectedObjects) > 1 {
		diagnostics.AddError(
			"Unexpected APIC response cardinality",
			fmt.Sprintf("APIC returned %d objects for class %q; expected at most one.", len(expectedObjects), className),
		)
		return nil, ""
	}

	object := expectedObjects[0]
	attributes, ok := AttributesFromObject(ctx, diagnostics, object, className)
	if !ok {
		return nil, ""
	}

	dnAttribute, found := attributes.values["dn"]
	if !found {
		diagnostics.AddError(
			"Missing APIC object identity",
			fmt.Sprintf("APIC object %s does not contain the dn attribute.", className),
		)
		return nil, ""
	}
	dn, ok := dnAttribute.Data().(string)
	if !ok {
		diagnostics.AddError(
			"Malformed APIC object identity",
			fmt.Sprintf("APIC dn for class %q must be a string, got %T.", className, dnAttribute.Data()),
		)
		return nil, ""
	}

	model, modelDiagnostics := decode(ctx, object, fallbackModel)
	diagnostics.Append(modelDiagnostics...)
	return &model, dn
}
