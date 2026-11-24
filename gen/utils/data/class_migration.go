package data

import (
	"fmt"
	"slices"
	"strings"
)

// MigrationSchema is the current metadata with one version's decoding overrides applied.
type MigrationSchema struct {
	PriorSchemaVersion int
	Attributes         []MigrationAttribute
}

func (s MigrationSchema) Blocks() []MigrationAttribute {
	var blocks []MigrationAttribute
	for _, attribute := range s.Attributes {
		if attribute.Block {
			blocks = append(blocks, attribute)
		}
	}
	return blocks
}

// MigrationAttribute keeps historical decoding separate from current alias exposure.
type MigrationAttribute struct {
	AttributeName          string
	FieldName              string
	Type                   LegacyAttributeTypeEnum
	ElementType            LegacyAttributeTypeEnum
	Required               bool
	Optional               bool
	Computed               bool
	Sensitive              bool
	Status                 LegacyStatusEnum
	TargetField            string
	ChildName              *ClassName
	ChildAttribute         string
	IdentifyingAttributes  []string
	ResponseAttribute      string
	Named                  bool
	SchemaExpression       string
	InheritedModelType     string
	InheritedAttributeType string
	InheritedNullValue     string
	TargetModelType        string
	TargetAttributeType    string
	RequiresReplace        bool
	Repeated               bool
	Block                  bool
	Attributes             []MigrationAttribute
}

func (a MigrationAttribute) IsAlias() bool {
	return a.FieldName != a.TargetField || a.ChildName != nil && !a.IsChildTypeChange()
}

func (a MigrationAttribute) IsChildTypeChange() bool {
	return a.ChildName != nil && a.Type == SetNestedAttribute && !a.Block
}

func (a MigrationAttribute) IsListToSet() bool {
	return a.Type == ListAttribute && a.TargetModelType == "types.Set"
}

func (a MigrationAttribute) ModelValueType() string {
	if a.InheritedModelType != "" {
		return a.InheritedModelType
	}
	switch a.Type {
	case StringAttribute:
		return "types.String"
	case SetAttribute, SetNestedAttribute:
		return "types.Set"
	case ListAttribute:
		return "types.List"
	default:
		return ""
	}
}

func (a MigrationAttribute) ModelAttributeType() string {
	if a.InheritedAttributeType != "" {
		return a.InheritedAttributeType
	}
	if a.Block {
		return fmt.Sprintf("types.SetType{ElemType: %s}", a.ElementAttributeType())
	}
	switch a.Type {
	case StringAttribute:
		return "types.StringType"
	case SetAttribute:
		return "types.SetType{ElemType: types.StringType}"
	case ListAttribute:
		return "types.ListType{ElemType: types.StringType}"
	default:
		return ""
	}
}

func (a MigrationAttribute) ElementAttributeType() string {
	if !a.Block {
		return "types.StringType"
	}
	var attributes strings.Builder
	for _, member := range a.Attributes {
		fmt.Fprintf(&attributes, "%q: %s,", member.AttributeName, member.ModelAttributeType())
	}
	return fmt.Sprintf("types.ObjectType{AttrTypes: map[string]attr.Type{%s}}", attributes.String())
}

func (a MigrationAttribute) SchemaAttributeType() string {
	if a.Type == StringAttribute {
		return "StringAttribute"
	}
	return strings.TrimPrefix(a.ModelValueType(), "types.") + "Attribute"
}

func (a MigrationAttribute) ModelCustomStringType() string {
	if a.Type == StringAttribute && a.ModelValueType() != "types.String" {
		return a.ModelAttributeType()
	}
	return ""
}

func (a MigrationAttribute) NullValue() string {
	if a.InheritedNullValue != "" {
		return a.InheritedNullValue
	}
	if a.Block {
		return fmt.Sprintf("types.SetNull(%s)", a.ElementAttributeType())
	}
	switch a.Type {
	case StringAttribute:
		return "types.StringNull()"
	case SetAttribute:
		return "types.SetNull(types.StringType)"
	case ListAttribute:
		return "types.ListNull(types.StringType)"
	default:
		return ""
	}
}

func (c *Class) setMigrationData(ds *DataStore) {
	// Published schema versions are independent of compatibility implementation progress.
	if c.ClassDefinition.SchemaVersion != nil {
		c.SchemaVersion = *c.ClassDefinition.SchemaVersion
	}
	if len(c.StateUpgrades) == 0 {
		return
	}
	dataSourceOnly := c.HasDatasourceArtifact() && !c.HasResourceArtifact()
	genLogger.Debugf("Resolving state-upgrade differences for class '%s'.", c.Name)
	versions := make([]int, 0, len(c.StateUpgrades))
	for version := range c.StateUpgrades {
		versions = append(versions, version)
	}
	slices.Sort(versions)
	aliases := map[string]MigrationAttribute{}
	overrides := map[int][]MigrationAttribute{}
	flattened := map[int]map[string]bool{}
	for _, version := range versions {
		entry := c.StateUpgrades[version]
		if version < 0 || !dataSourceOnly && version >= c.SchemaVersion {
			ds.ctx.Diagnostics.AddError("Class '%s': schema version %d must precede schema_version %d", c.Name, version, c.SchemaVersion)
			continue
		}
		flattened[version] = map[string]bool{}
		add := func(key string, definition AttributeUpgradeDefinition, child *Class) {
			name := definition.AttributeName
			if name == "" {
				name = key
			}
			if previous, exists := aliases[name]; exists {
				if definition.Status == UndefinedLegacyStatus {
					definition.Status = previous.Status
				}
				// Orphaned retained fields have no current metadata to supply their old shape.
				owner := c
				if child != nil {
					owner = child
				}
				if owner.Properties[key] == nil {
					if definition.Type == UndefinedLegacyAttributeType {
						definition.Type = previous.Type
					}
					if definition.ElementType == UndefinedLegacyAttributeType {
						definition.ElementType = previous.ElementType
					}
					if definition.Restriction == UndefinedLegacyRestriction {
						switch {
						case previous.Required:
							definition.Restriction = LegacyRequired
						case previous.Optional && previous.Computed:
							definition.Restriction = LegacyOptionalComputed
						case previous.Optional:
							definition.Restriction = LegacyOptional
						case previous.Computed:
							definition.Restriction = LegacyReadOnly
						}
					}
				}
			}
			attribute := c.migrationAttribute(ds, key, definition, child)
			if attribute.IsAlias() {
				aliases[attribute.AttributeName] = attribute
			}
			overrides[version] = append(overrides[version], attribute)
		}
		for key, definition := range entry.Attributes {
			add(key, definition, nil)
		}
		for name, definition := range entry.Children {
			child, exists := ds.Classes[name]
			if !exists || len(definition.Children) > 0 {
				ds.ctx.Diagnostics.AddError("Class '%s': historical child %q needs nested-schema migration support before it can be rendered", c.Name, name)
				continue
			}
			if definition.AttributeName != "" {
				attribute := MigrationAttribute{
					AttributeName: definition.AttributeName, FieldName: "Legacy" + migrationFieldName(definition.AttributeName),
					TargetField: child.Name.Capitalized(), ChildName: child.Name, Type: SetNestedAttribute,
					Block: true, Repeated: !child.IsSingleNestedWhenDefinedAsChild, Status: definition.Status,
				}
				if attribute.Status == UndefinedLegacyStatus {
					if previous, exists := aliases[attribute.AttributeName]; exists {
						attribute.Status = previous.Status
					} else {
						attribute.Status = LegacyDeprecated
					}
				}
				for key, member := range definition.Attributes {
					resolved := child.migrationAttribute(ds, key, member, nil)
					resolved.ResponseAttribute = key
					// Legacy named-relation members receive DNs, just like the old flat relationship attributes.
					if child.Relation.Type == Named {
						for _, target := range child.Relation.ToClasses {
							if key == target.namedRelationPropertyName() {
								resolved.Named, resolved.ResponseAttribute = true, "tDn"
								attribute.Named, attribute.ChildAttribute = true, resolved.ChildAttribute
							}
						}
					}
					attribute.Attributes = append(attribute.Attributes, resolved)
				}
				for _, identifier := range child.IdentifiedBy {
					attribute.IdentifyingAttributes = append(attribute.IdentifyingAttributes, child.Properties[identifier].AttributeName)
				}
				if !attribute.Named && len(attribute.IdentifyingAttributes) > 0 {
					attribute.ChildAttribute = attribute.IdentifyingAttributes[0]
				}
				slices.SortFunc(attribute.Attributes, func(a, b MigrationAttribute) int { return strings.Compare(a.AttributeName, b.AttributeName) })
				aliases[attribute.AttributeName] = attribute
				overrides[version] = append(overrides[version], attribute)
				flattened[version][child.ResourceNameNested] = true
				continue
			}
			if definition.Type == SetNestedAttribute && child.IsSingleNestedWhenDefinedAsChild {
				attribute := c.currentMigrationAttributes(ds)[child.ResourceNameNested]
				attribute.ChildName, attribute.Type = child.Name, SetNestedAttribute
				attribute.TargetModelType, attribute.TargetAttributeType = attribute.InheritedModelType, attribute.InheritedAttributeType
				attribute.InheritedModelType = "types.Set"
				attribute.InheritedAttributeType = fmt.Sprintf("types.SetType{ElemType: %s}", attribute.TargetAttributeType)
				attribute.SchemaExpression = fmt.Sprintf("modelHelpers.LegacyChildSetAttribute(current.Attributes[%q].(resourceschema.SingleNestedAttribute))", child.ResourceNameNested)
				overrides[version] = append(overrides[version], attribute)
				continue
			}
			if definition.Type != UndefinedLegacyAttributeType {
				ds.ctx.Diagnostics.AddError("Class '%s': historical child %q needs in-place nested-override support before it can be rendered", c.Name, name)
				continue
			}
			for key, attribute := range definition.Attributes {
				if attribute.AttributeName == "" {
					ds.ctx.Diagnostics.AddError("Class '%s': child %q needs in-place nested-override support before it can be rendered", c.Name, name)
					continue
				}
				flattened[version][child.ResourceNameNested] = true
				add(key, attribute, &child)
			}
		}
	}
	for _, attribute := range aliases {
		c.LegacyAttributes = append(c.LegacyAttributes, attribute)
	}
	slices.SortFunc(c.LegacyAttributes, func(a, b MigrationAttribute) int { return strings.Compare(a.AttributeName, b.AttributeName) })
	for _, property := range c.Properties {
		for version, value := range property.StateUpgradeValues {
			if alias, exists := aliases[value.AttributeName]; exists {
				value.Status = alias.Status
				property.StateUpgradeValues[version] = value
			}
		}
	}
	// Datasource aliases are refreshed from APIC; they do not register resource state upgraders.
	if dataSourceOnly {
		return
	}
	for _, version := range versions {
		// The latest metadata supplies unchanged fields. Missing old JSON fields decode to null.
		attributes := c.currentMigrationAttributes(ds)
		for _, alias := range c.LegacyAttributes {
			if alias.Status == Removed {
				continue
			}
			field := alias
			field.ChildName, field.TargetField = nil, ""
			field.Required, field.Optional, field.Computed = false, true, true
			field.SchemaExpression = fmt.Sprintf("current.Attributes[%q]", field.AttributeName)
			if field.Block {
				field.SchemaExpression = fmt.Sprintf("current.Blocks[%q]", field.AttributeName)
			}
			attributes[field.AttributeName] = field
		}
		for name := range flattened[version] {
			delete(attributes, name)
		}
		for _, attribute := range overrides[version] {
			if attribute.ChildName == nil && attribute.ChildAttribute != "" {
				delete(attributes, attribute.ChildAttribute)
			}
			if attribute.IsAlias() {
				attribute.Status = aliases[attribute.AttributeName].Status
			}
			attributes[attribute.AttributeName] = attribute
		}
		schema := MigrationSchema{PriorSchemaVersion: version}
		for _, attribute := range attributes {
			schema.Attributes = append(schema.Attributes, attribute)
		}
		slices.SortFunc(schema.Attributes, func(a, b MigrationAttribute) int { return strings.Compare(a.AttributeName, b.AttributeName) })
		c.MigrationSchemas = append(c.MigrationSchemas, schema)
	}
}

func (c *Class) currentMigrationAttributes(ds *DataStore) map[string]MigrationAttribute {
	attributes := map[string]MigrationAttribute{}
	for _, property := range c.ModelProperties() {
		kind := StringAttribute
		if property.ValueType == Set {
			kind = SetAttribute
		}
		attributes[property.AttributeName] = MigrationAttribute{
			AttributeName: property.AttributeName, FieldName: property.ModelFieldName(c.Name), TargetField: property.ModelFieldName(c.Name), ChildAttribute: property.AttributeName,
			Type: kind, ElementType: StringAttribute, Required: property.Required, Optional: property.Optional, Computed: property.Computed,
			InheritedModelType: property.ModelValueType(c.Name), InheritedAttributeType: property.ModelAttributeType(c.Name),
			SchemaExpression: fmt.Sprintf("current.Attributes[%q]", property.AttributeName),
		}
	}
	attributes["id"] = MigrationAttribute{AttributeName: "id", FieldName: "ID", TargetField: "ID", ChildAttribute: "id", Type: StringAttribute, Computed: true, SchemaExpression: `current.Attributes["id"]`}
	if parent := c.ParentDnProperty(); parent != nil {
		attributes[parent.AttributeName] = MigrationAttribute{AttributeName: parent.AttributeName, FieldName: "ParentDN", TargetField: "ParentDN", ChildAttribute: parent.AttributeName, Type: StringAttribute, Required: parent.Required, Optional: parent.Optional, Computed: parent.Computed, SchemaExpression: fmt.Sprintf("current.Attributes[%q]", parent.AttributeName)}
	}
	for _, name := range c.Children {
		child := ds.Classes[name.String()]
		valueType, attributeType := "types.Set", fmt.Sprintf("types.SetType{ElemType: types.ObjectType{AttrTypes: %sModelAttributeTypes()}}", child.Name.Capitalized())
		if child.IsSingleNestedWhenDefinedAsChild {
			valueType, attributeType = "types.Object", fmt.Sprintf("types.ObjectType{AttrTypes: %sModelAttributeTypes()}", child.Name.Capitalized())
		}
		attributes[child.ResourceNameNested] = MigrationAttribute{AttributeName: child.ResourceNameNested, FieldName: child.Name.Capitalized(), TargetField: child.Name.Capitalized(), ChildAttribute: child.ResourceNameNested, Optional: true, Computed: true, InheritedModelType: valueType, InheritedAttributeType: attributeType, SchemaExpression: fmt.Sprintf("current.Attributes[%q]", child.ResourceNameNested)}
	}
	return attributes
}

func (c *Class) migrationAttribute(ds *DataStore, key string, definition AttributeUpgradeDefinition, child *Class) MigrationAttribute {
	owner := c
	if child != nil {
		owner = child
	}
	property := owner.Properties[key]
	attribute := MigrationAttribute{AttributeName: definition.AttributeName, Type: definition.Type, ElementType: definition.ElementType, Status: definition.Status}
	if property != nil {
		attribute.TargetField, attribute.ChildAttribute = property.ModelFieldName(owner.Name), property.AttributeName
		attribute.TargetModelType, attribute.TargetAttributeType = property.ModelValueType(owner.Name), property.ModelAttributeType(owner.Name)
		attribute.RequiresReplace = child == nil && property.RequiresReplace
		if child == nil && key == "parentDn" {
			attribute.TargetField, attribute.RequiresReplace = "ParentDN", true
		}
		attribute.Required, attribute.Optional, attribute.Computed = property.Required, property.Optional, property.Computed
		attribute.Sensitive = property.Sensitive
		if attribute.Type == UndefinedLegacyAttributeType {
			attribute.Type = StringAttribute
			if property.ValueType == Set {
				attribute.Type, attribute.ElementType = SetAttribute, StringAttribute
			}
		}
		if property.ValueType == Set && attribute.ElementType == UndefinedLegacyAttributeType {
			attribute.ElementType = StringAttribute
		}
		if attribute.AttributeName == "" {
			attribute.AttributeName = property.AttributeName
		}
	} else if key == "id" && child == nil {
		attribute.TargetField, attribute.ChildAttribute, attribute.Computed = "ID", "id", true
		if attribute.Type == UndefinedLegacyAttributeType {
			attribute.Type = StringAttribute
		}
	} else if definition.Status != Frozen && definition.Status != Removed {
		ds.ctx.Diagnostics.AddError("Class '%s': historical attribute %q has no current target or frozen/removed status", c.Name, key)
	}
	if attribute.AttributeName == "" {
		attribute.AttributeName = key
	}
	if definition.ValueType != UndefinedValueType {
		// Only explicitly recorded historical custom types affect the retained alias.
		historical := Property{PropertyName: key}
		if property != nil {
			historical = *property
		}
		historical.ValueType = definition.ValueType
		if attribute.Type != StringAttribute || historical.ValueType.SchemaPropertyType() != "StringProperty" {
			ds.ctx.Diagnostics.AddError("Class '%s': historical attribute %q value_type requires a string attribute", c.Name, key)
		} else {
			attribute.InheritedModelType = historical.ModelValueType(owner.Name)
			attribute.InheritedAttributeType = historical.ModelAttributeType(owner.Name)
			attribute.InheritedNullValue = historical.ModelNullValue(owner.Name)
		}
	}
	switch definition.Restriction {
	case LegacyRequired:
		attribute.Required, attribute.Optional, attribute.Computed = true, false, false
	case LegacyOptional:
		attribute.Required, attribute.Optional, attribute.Computed = false, true, false
	case LegacyOptionalComputed:
		attribute.Required, attribute.Optional, attribute.Computed = false, true, true
	case LegacyReadOnly:
		attribute.Required, attribute.Optional, attribute.Computed = false, false, true
	}
	if child != nil {
		attribute.ChildName, attribute.Named = child.Name, child.Relation.Type == Named
		attribute.Repeated = !child.IsSingleNestedWhenDefinedAsChild
	}
	attribute.FieldName = attribute.TargetField
	if child != nil || attribute.TargetField == "" || (property != nil && attribute.AttributeName != property.AttributeName) {
		attribute.FieldName = "Legacy" + migrationFieldName(attribute.AttributeName)
	}
	if attribute.IsAlias() && attribute.Status == UndefinedLegacyStatus {
		attribute.Status = LegacyDeprecated
	}
	if attribute.Type != StringAttribute && attribute.Type != SetAttribute && attribute.Type != ListAttribute {
		ds.ctx.Diagnostics.AddError("Class '%s': historical attribute %q needs migration support for type %q", c.Name, key, attribute.Type)
	}
	if attribute.Type != StringAttribute && attribute.ElementType != StringAttribute {
		ds.ctx.Diagnostics.AddError("Class '%s': historical collection %q requires a supported element_type", c.Name, key)
	}
	if child != nil && !((attribute.Type == StringAttribute && !attribute.Repeated) || (attribute.Type == SetAttribute && attribute.Repeated)) {
		ds.ctx.Diagnostics.AddError("Class '%s': historical attribute %q needs collection-to-child migration support before it can be rendered", c.Name, key)
	}
	return attribute
}

func migrationFieldName(name string) string {
	parts := strings.Split(name, "_")
	for index, part := range parts {
		if part != "" {
			parts[index] = fmt.Sprintf("%s%s", strings.ToUpper(part[:1]), part[1:])
		}
	}
	return strings.Join(parts, "")
}
