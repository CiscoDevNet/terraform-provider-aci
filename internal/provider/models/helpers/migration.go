package helpers

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	aciValidators "github.com/CiscoDevNet/terraform-provider-aci/v2/internal/validators"
)

// LegacyAttributeSchema describes compatibility fields and historical nested member mappings.
type LegacyAttributeSchema struct {
	Type            attr.Type
	Replacement     string
	Frozen          bool
	RequiresReplace bool
	Attributes      map[string]LegacyAttributeSchema
	Required        bool
	Sensitive       bool
	ReadOnly        bool
	Removed         bool
}

func (s LegacyAttributeSchema) Upgrade(ctx context.Context, diagnostics *diag.Diagnostics, value attr.Value, target attr.Type) attr.Value {
	// Preserve null/unknown containers before traversing their historical member mappings.
	null, err := target.ValueFromTerraform(ctx, tftypes.NewValue(target.TerraformType(ctx), nil))
	if err != nil {
		diagnostics.AddError("Unable to convert legacy value", err.Error())
		return nil
	}
	if s.Removed || value.IsNull() {
		return null
	}
	if value.IsUnknown() {
		unknown, err := target.ValueFromTerraform(ctx, tftypes.NewValue(target.TerraformType(ctx), tftypes.UnknownValue))
		if err != nil {
			diagnostics.AddError("Unable to convert legacy value", err.Error())
			return null
		}
		return unknown
	}
	if _, objectTarget := target.(types.ObjectType); objectTarget {
		if set, historicalSet := value.(types.Set); historicalSet {
			if elementType, objectElements := set.ElementType(ctx).(types.ObjectType); objectElements {
				// Keep historical member names/types until the recursive mappings have converted them.
				value = LegacyChildFromSet(ctx, diagnostics, set, elementType.AttrTypes)
				return s.Upgrade(ctx, diagnostics, value, target)
			}
		}
	}
	if len(s.Attributes) == 0 {
		converted := LegacyValue[attr.Value](ctx, diagnostics, value, target)
		if converted == nil {
			return null
		}
		return converted
	}
	switch target := target.(type) {
	case types.ObjectType:
		object, ok := value.(types.Object)
		if !ok {
			break
		}
		return s.upgradeObject(ctx, diagnostics, object, target)
	case types.SetType:
		var elements []attr.Value
		switch value := value.(type) {
		case types.Set:
			elements = value.Elements()
		case types.List:
			elements = value.Elements()
		default:
			diagnostics.AddError("Unable to convert legacy collection", "Expected a historical list or set of child objects.")
			return null
		}
		converted := make([]attr.Value, 0, len(elements))
		for _, element := range elements {
			converted = append(converted, s.Upgrade(ctx, diagnostics, element, target.ElemType))
		}
		if diagnostics.HasError() {
			return null
		}
		list, valueDiagnostics := types.ListValue(target.ElemType, converted)
		diagnostics.Append(valueDiagnostics...)
		if valueDiagnostics.HasError() {
			return null
		}
		return LegacyValue[attr.Value](ctx, diagnostics, list, target)
	}
	diagnostics.AddError("Unable to convert legacy child", fmt.Sprintf("Cannot convert historical %T to %T with child mappings.", value, target))
	return null
}

func (s LegacyAttributeSchema) upgradeObject(ctx context.Context, diagnostics *diag.Diagnostics, value types.Object, target types.ObjectType) types.Object {
	// Current metadata supplies unchanged members; only recorded differences replace their historical names/types.
	attributes := map[string]attr.Value{}
	for name, attributeType := range target.AttrTypes {
		null, err := attributeType.ValueFromTerraform(ctx, tftypes.NewValue(attributeType.TerraformType(ctx), nil))
		if err != nil {
			diagnostics.AddError("Unable to convert legacy child", err.Error())
			return types.ObjectNull(target.AttrTypes)
		}
		attributes[name] = null
	}
	previous := value.Attributes()
	for name, oldValue := range previous {
		mapping := s.Attributes[name]
		if mapping.Removed || mapping.Replacement != "" && mapping.Replacement != name {
			continue
		}
		if attributeType, exists := target.AttrTypes[name]; exists {
			attributes[name] = mapping.Upgrade(ctx, diagnostics, oldValue, attributeType)
		}
	}
	for name, oldValue := range previous {
		mapping := s.Attributes[name]
		replacement := mapping.Replacement
		if mapping.Removed || replacement == "" || replacement == name {
			continue
		}
		if attributeType, exists := target.AttrTypes[replacement]; exists {
			current, hasCurrent := previous[replacement]
			if !hasCurrent || current.IsNull() || current.IsUnknown() {
				attributes[replacement] = mapping.Upgrade(ctx, diagnostics, oldValue, attributeType)
			}
		}
		if attributeType, retained := target.AttrTypes[name]; retained {
			attributes[name] = LegacyValue[attr.Value](ctx, diagnostics, oldValue, attributeType)
		}
	}
	if diagnostics.HasError() {
		return types.ObjectNull(target.AttrTypes)
	}
	result, valueDiagnostics := types.ObjectValue(target.AttrTypes, attributes)
	diagnostics.Append(valueDiagnostics...)
	return result
}

func (s ClassSchema) SetFromLegacyConfig(ctx context.Context, diagnostics *diag.Diagnostics, config, plan, state types.Object) types.Object {
	// Normalize deprecated inputs before reusing the schema's recursive state reconciliation.
	plan, normalizedConfig := s.legacyPlanObject(ctx, diagnostics, config, plan, state)
	if diagnostics.HasError() {
		return plan
	}
	reconciled, reconciliationDiagnostics := s.reconciliationSchema().Reconcile(ctx, normalizedConfig, plan, state)
	diagnostics.Append(reconciliationDiagnostics...)
	return reconciled.(types.Object)
}

func (s ClassSchema) legacyConfigObject(ctx context.Context, diagnostics *diag.Diagnostics, value types.Object) types.Object {
	if value.IsNull() || value.IsUnknown() || objectIsEmpty(value) {
		return value
	}
	attributes := value.Attributes()
	attributeTypes := value.AttributeTypes(ctx)
	for name, legacy := range s.NestedLegacyAttributes {
		if legacy.Frozen {
			continue
		}
		configured := attributes[name]
		if configured.IsNull() {
			// An omitted computed alias must not overwrite a changed canonical value with stale state.
			unknown, err := attributeTypes[name].ValueFromTerraform(ctx, tftypes.NewValue(attributeTypes[name].TerraformType(ctx), tftypes.UnknownValue))
			if err != nil {
				diagnostics.AddError("Unable to construct legacy child plan", err.Error())
				continue
			}
			attributes[name] = unknown
			continue
		}
		current := attributes[legacy.Replacement]
		if current.IsNull() || current.IsUnknown() {
			attributes[legacy.Replacement] = LegacyValue[attr.Value](ctx, diagnostics, configured, attributeTypes[legacy.Replacement])
		}
	}
	if diagnostics.HasError() {
		return value
	}
	result, valueDiagnostics := types.ObjectValue(attributeTypes, attributes)
	diagnostics.Append(valueDiagnostics...)
	return result
}

func (s ClassSchema) legacyPlanObject(ctx context.Context, diagnostics *diag.Diagnostics, config, plan, state types.Object) (types.Object, types.Object) {
	if config.IsNull() || config.IsUnknown() || objectIsEmpty(config) {
		return plan, config
	}
	attributeTypes := config.AttributeTypes(ctx)
	configured := config.Attributes()
	planned := plan.Attributes()
	previous := state.Attributes()
	if plan.IsNull() || plan.IsUnknown() {
		planned = config.Attributes()
		for name, value := range planned {
			if value.IsNull() {
				unknown, err := attributeTypes[name].ValueFromTerraform(ctx, tftypes.NewValue(attributeTypes[name].TerraformType(ctx), tftypes.UnknownValue))
				if err != nil {
					diagnostics.AddError("Unable to construct legacy child plan", err.Error())
					return plan, config
				}
				planned[name] = unknown
			}
		}
	}
	for name, legacy := range s.NestedLegacyAttributes {
		if !configured[name].IsNull() {
			if !configured[name].IsUnknown() {
				planned[name] = configured[name]
			}
			if !legacy.Frozen {
				if !configured[legacy.Replacement].IsNull() {
					diagnostics.AddError("Conflicting legacy configuration", fmt.Sprintf("Configure either %s or %s, not both.", name, legacy.Replacement))
					return plan, config
				}
				value := configured[name]
				if value.IsUnknown() && !planned[name].IsUnknown() {
					value = planned[name]
				}
				planned[legacy.Replacement] = LegacyValue[attr.Value](ctx, diagnostics, value, attributeTypes[legacy.Replacement])
			}
			continue
		}
		if prior, exists := previous[name]; exists && (legacy.Frozen || legacy.Sensitive || planned[legacy.Replacement].Equal(previous[legacy.Replacement])) {
			planned[name] = prior
		} else if !legacy.Frozen && !legacy.Sensitive {
			// Attribute plan modifiers may already have copied the old alias before its canonical value changed.
			unknown, err := attributeTypes[name].ValueFromTerraform(ctx, tftypes.NewValue(attributeTypes[name].TerraformType(ctx), tftypes.UnknownValue))
			if err != nil {
				diagnostics.AddError("Unable to construct legacy child plan", err.Error())
				return plan, config
			}
			planned[name] = unknown
		}
	}
	normalized := s.legacyConfigObject(ctx, diagnostics, config)
	if diagnostics.HasError() {
		return plan, config
	}
	normalizedAttributes := normalized.Attributes()
	for name, child := range s.Children {
		childSchema := child.Schema()
		if child.Cardinality == SingleChild {
			planned[name], normalizedAttributes[name] = childSchema.legacyPlanObject(ctx, diagnostics, configured[name].(types.Object), planned[name].(types.Object), previousObject(previous[name], attributeTypes[name].(types.ObjectType)))
			continue
		}
		prior := types.SetNull(attributeTypes[name].(types.SetType).ElemType)
		if value, exists := previous[name]; exists {
			prior = value.(types.Set)
		}
		planned[name], normalizedAttributes[name] = childSchema.legacyPlanSet(ctx, diagnostics, configured[name].(types.Set), planned[name].(types.Set), prior)
	}
	if diagnostics.HasError() {
		return plan, config
	}
	result, valueDiagnostics := types.ObjectValue(attributeTypes, planned)
	diagnostics.Append(valueDiagnostics...)
	normalized, valueDiagnostics = types.ObjectValue(attributeTypes, normalizedAttributes)
	diagnostics.Append(valueDiagnostics...)
	return result, normalized
}

func previousObject(value attr.Value, objectType types.ObjectType) types.Object {
	if value == nil {
		return types.ObjectNull(objectType.AttrTypes)
	}
	return value.(types.Object)
}

func (s ClassSchema) legacyPlanSet(ctx context.Context, diagnostics *diag.Diagnostics, config, plan, state types.Set) (types.Set, types.Set) {
	if config.IsNull() || config.IsUnknown() {
		return plan, config
	}
	objects := func(value types.Set) []types.Object {
		var result []types.Object
		for _, element := range value.Elements() {
			object := element.(types.Object)
			if object.IsNull() || object.IsUnknown() {
				continue
			}
			result = append(result, s.legacyConfigObject(ctx, diagnostics, object))
		}
		return result
	}
	plannedObjects, previousObjects := objects(plan), objects(state)
	plannedMatched, previousMatched := make([]bool, len(plannedObjects)), make([]bool, len(previousObjects))
	elementType := config.ElementType(ctx).(types.ObjectType)
	var planned, normalized []attr.Value
	matching := s.reconciliationSchema()
	for _, element := range config.Elements() {
		object := element.(types.Object)
		if object.IsNull() || object.IsUnknown() {
			planned, normalized = append(planned, object), append(normalized, object)
			continue
		}
		identifiers := s.legacyConfigObject(ctx, diagnostics, object)
		plannedObject := types.ObjectUnknown(elementType.AttrTypes)
		previousObject := types.ObjectNull(elementType.AttrTypes)
		if len(s.Identifiers) > 0 {
			candidate, found, matchingDiagnostics := matching.MatchingObject(ctx, identifiers, plannedObjects, plannedMatched)
			diagnostics.Append(matchingDiagnostics...)
			if found {
				plannedObject = candidate
			}
			candidate, found, matchingDiagnostics = matching.MatchingObject(ctx, identifiers, previousObjects, previousMatched)
			diagnostics.Append(matchingDiagnostics...)
			if found {
				previousObject = candidate
			}
		}
		plannedObject, normalizedObject := s.legacyPlanObject(ctx, diagnostics, object, plannedObject, previousObject)
		planned, normalized = append(planned, plannedObject), append(normalized, normalizedObject)
	}
	if diagnostics.HasError() {
		return plan, config
	}
	plannedSet, valueDiagnostics := types.SetValue(elementType, planned)
	diagnostics.Append(valueDiagnostics...)
	normalizedSet, valueDiagnostics := types.SetValue(elementType, normalized)
	diagnostics.Append(valueDiagnostics...)
	return plannedSet, normalizedSet
}

func (s ClassSchema) SetLegacyFromValues(ctx context.Context, diagnostics *diag.Diagnostics, value types.Object, fallback *types.Object) types.Object {
	// Populate class-owned aliases after normal APIC decoding, retaining frozen and unreadable values from their own state.
	if value.IsNull() || value.IsUnknown() {
		return value
	}
	attributes := value.Attributes()
	attributeTypes := value.AttributeTypes(ctx)
	for name, legacy := range s.NestedLegacyAttributes {
		if legacy.Frozen || legacy.Sensitive && !s.Properties[legacy.Replacement].ReadOnly {
			if fallback != nil && !fallback.IsNull() && !fallback.IsUnknown() {
				if prior := fallback.Attributes()[name]; !prior.IsUnknown() {
					attributes[name] = prior
				}
				if property, exists := s.Properties[legacy.Replacement]; exists && property.Sensitive && !property.ReadOnly {
					prior := fallback.Attributes()[legacy.Replacement]
					if !prior.IsUnknown() {
						attributes[legacy.Replacement] = prior
					}
				}
			}
			continue
		}
		attributes[name] = LegacyValue[attr.Value](ctx, diagnostics, attributes[legacy.Replacement], attributeTypes[name])
	}
	previous := map[string]attr.Value{}
	if fallback != nil && !fallback.IsNull() && !fallback.IsUnknown() {
		previous = fallback.Attributes()
	}
	for name, child := range s.Children {
		childSchema := child.Schema()
		if child.Cardinality == SingleChild {
			prior := previousObject(previous[name], attributeTypes[name].(types.ObjectType))
			attributes[name] = childSchema.SetLegacyFromValues(ctx, diagnostics, attributes[name].(types.Object), &prior)
			continue
		}
		current := attributes[name].(types.Set)
		if current.IsNull() || current.IsUnknown() {
			continue
		}
		var priorObjects []types.Object
		if previous[name] != nil {
			for _, element := range previous[name].(types.Set).Elements() {
				priorObjects = append(priorObjects, element.(types.Object))
			}
		}
		matched := make([]bool, len(priorObjects))
		var elements []attr.Value
		for _, element := range current.Elements() {
			var prior *types.Object
			if len(childSchema.Identifiers) > 0 {
				previous, found, matchingDiagnostics := childSchema.reconciliationSchema().MatchingObject(ctx, element.(types.Object), priorObjects, matched)
				diagnostics.Append(matchingDiagnostics...)
				if found {
					prior = &previous
				}
			}
			elements = append(elements, childSchema.SetLegacyFromValues(ctx, diagnostics, element.(types.Object), prior))
		}
		children, valueDiagnostics := types.SetValue(current.ElementType(ctx), elements)
		diagnostics.Append(valueDiagnostics...)
		attributes[name] = children
	}
	if diagnostics.HasError() {
		return value
	}
	result, valueDiagnostics := types.ObjectValue(attributeTypes, attributes)
	diagnostics.Append(valueDiagnostics...)
	return result
}

func (s LegacyAttributeSchema) ResourceBlock() resourceschema.Block {
	attributes := make(map[string]resourceschema.Attribute, len(s.Attributes))
	for name, member := range s.Attributes {
		attribute := member.ResourceAttribute().(resourceschema.StringAttribute)
		attribute.PlanModifiers = nil
		attribute.Validators = []validator.String{stringvalidator.ConflictsWith(path.MatchRoot(s.Replacement))}
		if member.Required {
			attribute.Validators = append(attribute.Validators, aciValidators.StringNotNull())
		}
		attributes[name] = attribute
	}
	return resourceschema.SetNestedBlock{NestedObject: resourceschema.NestedBlockObject{Attributes: attributes}}
}

func (s LegacyAttributeSchema) DataSourceBlock() datasourceschema.Block {
	attributes := make(map[string]datasourceschema.Attribute, len(s.Attributes))
	for name, member := range s.Attributes {
		attributes[name] = member.DataSourceAttribute()
	}
	return datasourceschema.SetNestedBlock{NestedObject: datasourceschema.NestedBlockObject{Attributes: attributes}}
}

func (s LegacyAttributeSchema) ResourceAttribute() resourceschema.Attribute {
	return s.resourceAttribute(path.MatchRoot(s.Replacement))
}

func (s LegacyAttributeSchema) resourceAttribute(replacementPath path.Expression) resourceschema.Attribute {
	message := deprecatedPropertyMessage
	if s.Replacement != "" {
		message = fmt.Sprintf("Use attribute `%s` instead.", s.Replacement)
	}
	if s.Frozen {
		message = "This deprecated attribute retains configuration and state for compatibility but no longer configures an APIC object."
	}
	switch attributeType := s.Type.(type) {
	case basetypes.StringTypable:
		attribute := resourceschema.StringAttribute{Optional: !s.ReadOnly, Computed: true, Sensitive: s.Sensitive, DeprecationMessage: message}
		if !s.Type.Equal(types.StringType) {
			attribute.CustomType = attributeType
		}
		if s.Frozen {
			attribute.PlanModifiers = ResourceStringPlanModifiers(false)
		}
		if s.RequiresReplace && !s.Frozen {
			attribute.PlanModifiers = append(attribute.PlanModifiers, stringplanmodifier.RequiresReplaceIfConfigured())
		}
		if s.Replacement != "" && !s.Frozen {
			attribute.Validators = []validator.String{stringvalidator.ConflictsWith(replacementPath)}
		}
		return attribute
	case basetypes.ListType:
		attribute := resourceschema.ListAttribute{Optional: !s.ReadOnly, Computed: true, Sensitive: s.Sensitive, ElementType: types.StringType, DeprecationMessage: message}
		if s.Replacement != "" && !s.Frozen {
			// The shared framework conflict validator implements both collection interfaces.
			attribute.Validators = []validator.List{setvalidator.ConflictsWith(replacementPath).(validator.List)}
		}
		return attribute
	default:
		attribute := resourceschema.SetAttribute{
			Optional: !s.ReadOnly, Computed: true, Sensitive: s.Sensitive, ElementType: types.StringType, DeprecationMessage: message,
		}
		if s.Frozen {
			attribute.PlanModifiers = []planmodifier.Set{setplanmodifier.UseNonNullStateForUnknown()}
		}
		if s.Replacement != "" && !s.Frozen {
			attribute.Validators = []validator.Set{setvalidator.ConflictsWith(replacementPath)}
		}
		return attribute
	}
}

func (s LegacyAttributeSchema) DataSourceAttribute() datasourceschema.Attribute {
	switch attributeType := s.Type.(type) {
	case basetypes.StringTypable:
		attribute := datasourceschema.StringAttribute{Computed: true, Sensitive: s.Sensitive, DeprecationMessage: deprecatedPropertyMessage}
		if !s.Type.Equal(types.StringType) {
			attribute.CustomType = attributeType
		}
		return attribute
	case basetypes.ListType:
		return datasourceschema.ListAttribute{Computed: true, Sensitive: s.Sensitive, ElementType: types.StringType, DeprecationMessage: deprecatedPropertyMessage}
	default:
		return datasourceschema.SetAttribute{Computed: true, Sensitive: s.Sensitive, ElementType: types.StringType, DeprecationMessage: deprecatedPropertyMessage}
	}
}

func LegacyValue[T attr.Value](ctx context.Context, diagnostics *diag.Diagnostics, value attr.Value, attributeType attr.Type) T {
	// Convert through Terraform's wire type to retain null/unknown values and custom string semantics.
	var result T
	raw, err := value.ToTerraformValue(ctx)
	if err != nil {
		diagnostics.AddError("Unable to convert legacy value", err.Error())
		return result
	}
	raw, err = legacyCollectionValue(raw, attributeType.TerraformType(ctx))
	if err != nil {
		diagnostics.AddError("Unable to convert legacy collection", err.Error())
		return result
	}
	converted, err := attributeType.ValueFromTerraform(ctx, raw)
	if err != nil {
		diagnostics.AddError("Unable to convert legacy value", err.Error())
		return result
	}
	result, ok := converted.(T)
	if !ok {
		diagnostics.AddError("Unable to convert legacy value", fmt.Sprintf("Unexpected converted type %T", converted))
	}
	return result
}

func legacyCollectionValue(value tftypes.Value, target tftypes.Type) (tftypes.Value, error) {
	// List/set changes preserve unknown members; only fully known duplicates can be removed safely.
	listToSet := false
	switch source := value.Type().(type) {
	case tftypes.List:
		destination, ok := target.(tftypes.Set)
		if !ok || !source.ElementType.Equal(destination.ElementType) {
			return value, nil
		}
		listToSet = true
	case tftypes.Set:
		destination, ok := target.(tftypes.List)
		if !ok || !source.ElementType.Equal(destination.ElementType) {
			return value, nil
		}
	default:
		return value, nil
	}
	if value.IsNull() {
		return tftypes.NewValue(target, nil), nil
	}
	if !value.IsKnown() {
		return tftypes.NewValue(target, tftypes.UnknownValue), nil
	}
	var elements []tftypes.Value
	if err := value.As(&elements); err != nil {
		return tftypes.Value{}, err
	}
	if listToSet {
		unique := make([]tftypes.Value, 0, len(elements))
		for _, element := range elements {
			if element.IsFullyKnown() && slices.ContainsFunc(unique, element.Equal) {
				continue
			}
			unique = append(unique, element)
		}
		elements = unique
	}
	return tftypes.NewValue(target, elements), nil
}

func LegacyStringList(ctx context.Context, diagnostics *diag.Diagnostics, current types.Set, legacy types.List) types.List {
	// Keep configured list order when its values still match; otherwise use the old sorted APIC read representation.
	if !legacy.IsNull() && !legacy.IsUnknown() {
		converted := LegacyValue[types.Set](ctx, diagnostics, legacy, types.SetType{ElemType: types.StringType})
		if converted.Equal(current) {
			return legacy
		}
	}
	converted := LegacyValue[types.List](ctx, diagnostics, current, types.ListType{ElemType: types.StringType})
	return LegacySortedStringList(ctx, diagnostics, converted)
}

func LegacySortedStringList(ctx context.Context, diagnostics *diag.Diagnostics, value types.List) types.List {
	// SDKv2 stored a sorted list; retain unresolved members until their order can be determined.
	if value.IsNull() || value.IsUnknown() {
		return value
	}
	for _, element := range value.Elements() {
		if element.IsNull() || element.IsUnknown() {
			return value
		}
	}
	var values []string
	diagnostics.Append(value.ElementsAs(ctx, &values, false)...)
	if diagnostics.HasError() {
		return types.ListNull(types.StringType)
	}
	slices.Sort(values)
	result, valueDiagnostics := types.ListValueFrom(ctx, types.StringType, values)
	diagnostics.Append(valueDiagnostics...)
	return result
}

func LegacyChildSetAttribute(current resourceschema.SingleNestedAttribute) resourceschema.SetNestedAttribute {
	// A historical collection uses the same member schema while retaining its original outer shape.
	return resourceschema.SetNestedAttribute{
		Required: current.Required, Optional: current.Optional, Computed: current.Computed,
		NestedObject: resourceschema.NestedAttributeObject{Attributes: current.Attributes},
	}
}

func LegacyChildFromSet(ctx context.Context, diagnostics *diag.Diagnostics, previous types.Set, attributeTypes map[string]attr.Type) types.Object {
	// Preserve the published NetFlow upgrade: keep the first child, or null for an unset/empty collection.
	if previous.IsUnknown() {
		return types.ObjectUnknown(attributeTypes)
	}
	if len(previous.Elements()) == 0 {
		return types.ObjectNull(attributeTypes)
	}
	return LegacyValue[types.Object](ctx, diagnostics, previous.Elements()[0], types.ObjectType{AttrTypes: attributeTypes})
}

// LegacyRelationName preserves Terraform null/unknown while adapting a legacy DN.
func LegacyRelationName(_ context.Context, _ *diag.Diagnostics, value types.String) types.String {
	if value.IsNull() || value.IsUnknown() {
		return value
	}
	parts := strings.Split(value.ValueString(), "/")
	if len(parts) == 1 {
		return value
	}
	_, name, _ := strings.Cut(parts[len(parts)-1], "-")
	return types.StringValue(name)
}

func LegacyChildObject(ctx context.Context, diagnostics *diag.Diagnostics, attributeTypes map[string]attr.Type, values map[string]attr.Value) types.Object {
	// Historical flat inputs populate a child whose other attributes must remain typed nulls.
	attributes := make(map[string]attr.Value, len(attributeTypes))
	maps.Copy(attributes, values)
	for name, attributeType := range attributeTypes {
		if _, exists := attributes[name]; exists {
			continue
		}
		value, err := attributeType.ValueFromTerraform(ctx, tftypes.NewValue(attributeType.TerraformType(ctx), nil))
		if err != nil {
			diagnostics.AddError("Unable to construct legacy child state", fmt.Sprintf("Attribute %q: %s", name, err))
			return types.ObjectNull(attributeTypes)
		}
		attributes[name] = value
	}
	object, objectDiagnostics := types.ObjectValue(attributeTypes, attributes)
	diagnostics.Append(objectDiagnostics...)
	return object
}

// LegacyChildAttributes preserves unrelated planned/read values when only an old alias is configured.
func LegacyChildAttributes(ctx context.Context, diagnostics *diag.Diagnostics, value types.Object, fallback *types.Object, attributeTypes map[string]attr.Type) map[string]attr.Value {
	var planned map[string]attr.Value
	if !value.IsNull() && !value.IsUnknown() && !objectIsEmpty(value) {
		planned = maps.Clone(value.Attributes())
	}
	if fallback != nil && !fallback.IsNull() && !fallback.IsUnknown() && !objectIsEmpty(*fallback) {
		if planned == nil {
			return maps.Clone(fallback.Attributes())
		}
		for name, previous := range fallback.Attributes() {
			if planned[name].IsNull() || planned[name].IsUnknown() {
				planned[name] = previous
			}
		}
	}
	if planned != nil {
		return planned
	}
	attributes := make(map[string]attr.Value, len(attributeTypes))
	for name, attributeType := range attributeTypes {
		unknown, err := attributeType.ValueFromTerraform(ctx, tftypes.NewValue(attributeType.TerraformType(ctx), tftypes.UnknownValue))
		if err != nil {
			diagnostics.AddError("Unable to construct legacy child plan", fmt.Sprintf("Attribute %q: %s", name, err))
			continue
		}
		attributes[name] = unknown
	}
	return attributes
}

// LegacyChildRelation maps a legacy relation DN to a singleton child's identifying property.
type LegacyChildRelation struct {
	AttributeTypes map[string]attr.Type
	AttributeName  string
	Named          bool
}

func (s LegacyChildRelation) Plan(ctx context.Context, diagnostics *diag.Diagnostics, legacy types.String, planned types.Object, fallback *types.Object, annotation attr.Value) types.Object {
	// Preserve the child's other planned/read values while replacing only the legacy relation's identifier.
	attributes := LegacyChildAttributes(ctx, diagnostics, planned, fallback, s.AttributeTypes)
	if s.Named {
		legacy = LegacyRelationName(ctx, diagnostics, legacy)
	}
	attributes[s.AttributeName] = legacy
	if value, exists := attributes["annotation"]; exists && value.IsUnknown() && annotation != nil {
		attributes["annotation"] = annotation
	}
	object, objectDiagnostics := types.ObjectValue(s.AttributeTypes, attributes)
	diagnostics.Append(objectDiagnostics...)
	return object
}

// LegacyChildSet maps legacy relations to current children and matches their identifying attributes.
type LegacyChildSet struct {
	AttributeTypes        map[string]attr.Type
	AttributeName         string
	IdentifyingAttributes []string
	Named                 bool
	Attributes            map[string]string
	ResponseAttributes    map[string]string
	LegacyType            types.ObjectType
}

func (s LegacyChildSet) Upgrade(ctx context.Context, diagnostics *diag.Diagnostics, legacy types.Set) types.Set {
	return s.convert(ctx, diagnostics, legacy, nil, nil, nil)
}

func (s LegacyChildSet) Plan(ctx context.Context, diagnostics *diag.Diagnostics, legacy, planned types.Set, fallback *types.Set, annotation attr.Value) types.Set {
	return s.convert(ctx, diagnostics, legacy, &planned, fallback, annotation)
}

func (s LegacyChildSet) convert(ctx context.Context, diagnostics *diag.Diagnostics, legacy types.Set, planned, fallback *types.Set, annotation attr.Value) types.Set {
	elementType := types.ObjectType{AttrTypes: s.AttributeTypes}
	if legacy.IsUnknown() {
		return types.SetUnknown(elementType)
	}
	// The published SDKv2 upgrade turned an unset DN collection into an empty child set.
	children := make([]attr.Value, 0, len(legacy.Elements()))
	for _, value := range legacy.Elements() {
		values := map[string]attr.Value{}
		if s.Attributes != nil {
			if value.IsUnknown() {
				children = append(children, types.ObjectUnknown(s.AttributeTypes))
				continue
			}
			for oldName, currentName := range s.Attributes {
				legacyValue := value.(types.Object).Attributes()[oldName]
				if s.Named && currentName == s.AttributeName {
					legacyValue = LegacyRelationName(ctx, diagnostics, legacyValue.(types.String))
				}
				values[currentName] = LegacyValue[attr.Value](ctx, diagnostics, legacyValue, s.AttributeTypes[currentName])
			}
		} else {
			identifier := value.(types.String)
			if s.Named {
				identifier = LegacyRelationName(ctx, diagnostics, identifier)
			}
			values[s.AttributeName] = LegacyValue[attr.Value](ctx, diagnostics, identifier, s.AttributeTypes[s.AttributeName])
		}
		if diagnostics.HasError() {
			return types.SetNull(elementType)
		}
		attributes := values
		if planned != nil {
			previous := s.match(fallback, values)
			attributes = LegacyChildAttributes(ctx, diagnostics, *s.match(planned, values), previous, s.AttributeTypes)
			for name, value := range values {
				if value.IsUnknown() && !attributes[name].IsUnknown() && !attributes[name].IsNull() {
					continue
				}
				attributes[name] = value
			}
			if value, exists := attributes["annotation"]; exists && value.IsUnknown() && annotation != nil {
				attributes["annotation"] = annotation
			}
		}
		children = append(children, LegacyChildObject(ctx, diagnostics, s.AttributeTypes, attributes))
	}
	result, valueDiagnostics := types.SetValue(elementType, children)
	diagnostics.Append(valueDiagnostics...)
	return result
}

func (s LegacyChildSet) Legacy(ctx context.Context, diagnostics *diag.Diagnostics, current types.Set) types.Set {
	// Decode only members of the historical block, leaving newer child fields out of compatibility state.
	if current.IsUnknown() {
		return types.SetUnknown(s.LegacyType)
	}
	values := make([]attr.Value, 0, len(current.Elements()))
	for _, child := range current.Elements() {
		if child.IsUnknown() {
			values = append(values, types.ObjectUnknown(s.LegacyType.AttrTypes))
			continue
		}
		attributes := map[string]attr.Value{}
		for oldName, currentName := range s.Attributes {
			attributes[oldName] = LegacyValue[attr.Value](ctx, diagnostics, child.(types.Object).Attributes()[currentName], s.LegacyType.AttrTypes[oldName])
		}
		values = append(values, LegacyChildObject(ctx, diagnostics, s.LegacyType.AttrTypes, attributes))
	}
	result, valueDiagnostics := types.SetValue(s.LegacyType, values)
	diagnostics.Append(valueDiagnostics...)
	return result
}

func (s LegacyChildSet) LegacyFromResponse(ctx context.Context, diagnostics *diag.Diagnostics, response *container.Container, parentClass, childClass string) types.Set {
	// Keep APIC's resolved target DN in legacy blocks while the modern model exposes the target name.
	values := []attr.Value{}
	if response != nil {
		objects, ok := response.Search("imdata").Data().([]any)
		if ok && len(objects) > 0 {
			parent := response.Search("imdata").Index(0).Search(parentClass)
			for _, child := range ChildObjectsByClass(ctx, diagnostics, parent, parentClass, []string{childClass})[childClass] {
				attributes, ok := AttributesFromObject(ctx, diagnostics, child, childClass)
				if !ok {
					continue
				}
				legacy := make(map[string]attr.Value, len(s.ResponseAttributes))
				for name, property := range s.ResponseAttributes {
					value := types.StringNull()
					DecodeStringAttribute(ctx, diagnostics, attributes, property, &value)
					legacy[name] = value
				}
				values = append(values, LegacyChildObject(ctx, diagnostics, s.LegacyType.AttrTypes, legacy))
			}
		}
	}
	result, valueDiagnostics := types.SetValue(s.LegacyType, values)
	diagnostics.Append(valueDiagnostics...)
	return result
}

func (s LegacyChildSet) Unknown(ctx context.Context, diagnostics *diag.Diagnostics, current types.Set) types.Set {
	// A legacy block cannot be computed as a whole, so preserve the modern set's size with unknown member values.
	// An unknown modern set has no known elements; keep the unconfigured legacy block empty until final planning.
	values := make([]attr.Value, 0, len(current.Elements()))
	for range current.Elements() {
		attributes := map[string]attr.Value{}
		for name, attributeType := range s.LegacyType.AttrTypes {
			value, err := attributeType.ValueFromTerraform(ctx, tftypes.NewValue(attributeType.TerraformType(ctx), tftypes.UnknownValue))
			if err != nil {
				diagnostics.AddError("Unable to construct legacy block plan", err.Error())
				continue
			}
			attributes[name] = value
		}
		values = append(values, LegacyChildObject(ctx, diagnostics, s.LegacyType.AttrTypes, attributes))
	}
	result, valueDiagnostics := types.SetValue(s.LegacyType, values)
	diagnostics.Append(valueDiagnostics...)
	return result
}

func (s LegacyChildSet) match(values *types.Set, identifiers map[string]attr.Value) *types.Object {
	// Historical blocks can identify a child with several fields, such as policy name and filter type.
	names := s.IdentifyingAttributes
	if len(names) == 0 {
		names = []string{s.AttributeName}
	}
	if values != nil {
		for _, value := range values.Elements() {
			child := value.(types.Object)
			if child.IsNull() || child.IsUnknown() {
				continue
			}
			matches := true
			attributes := child.Attributes()
			for _, name := range names {
				identifier := identifiers[name]
				if identifier.IsNull() || identifier.IsUnknown() || !attributes[name].Equal(identifier) {
					matches = false
					break
				}
			}
			if matches {
				return &child
			}
		}
	}
	missing := types.ObjectNull(s.AttributeTypes)
	return &missing
}

func LegacyTargetDNs(ctx context.Context, diagnostics *diag.Diagnostics, response *container.Container, parentClass, childClass string) types.Set {
	values := []attr.Value{}
	if response != nil {
		objects, ok := response.Search("imdata").Data().([]any)
		if ok && len(objects) > 0 {
			parent := response.Search("imdata").Index(0).Search(parentClass)
			for _, child := range ChildObjectsByClass(ctx, diagnostics, parent, parentClass, []string{childClass})[childClass] {
				attributes, ok := AttributesFromObject(ctx, diagnostics, child, childClass)
				if !ok {
					continue
				}
				if target, exists := attributes.values["tDn"]; exists && target.Data() != nil {
					dn, ok := target.Data().(string)
					if !ok {
						diagnostics.AddError("Malformed APIC relation target", fmt.Sprintf("APIC %s tDn must be a string.", childClass))
						continue
					}
					values = append(values, types.StringValue(dn))
				}
			}
		}
	}
	result, valueDiagnostics := types.SetValue(types.StringType, values)
	diagnostics.Append(valueDiagnostics...)
	return result
}

// LegacyTargetDN uses APIC's resolved target DN, which can differ from the relation's parent tenant.
func LegacyTargetDN(ctx context.Context, diagnostics *diag.Diagnostics, response *container.Container, parentClass, childClass string) types.String {
	if response == nil {
		return types.StringNull()
	}
	imdata := response.Search("imdata")
	if imdata == nil {
		return types.StringNull()
	}
	objects, ok := imdata.Data().([]any)
	if !ok || len(objects) == 0 {
		return types.StringNull()
	}
	object := response.Search("imdata").Index(0).Search(parentClass)
	children := ChildObjectsByClass(ctx, diagnostics, object, parentClass, []string{childClass})[childClass]
	if len(children) == 0 {
		return types.StringNull()
	}
	attributes, ok := AttributesFromObject(ctx, diagnostics, children[0], childClass)
	if !ok {
		return types.StringNull()
	}
	value, exists := attributes.values["tDn"]
	if !exists || value.Data() == nil {
		return types.StringNull()
	}
	dn, ok := value.Data().(string)
	if !ok {
		diagnostics.AddError("Malformed APIC relation target", fmt.Sprintf("APIC %s tDn must be a string.", childClass))
		return types.StringNull()
	}
	return types.StringValue(dn)
}
