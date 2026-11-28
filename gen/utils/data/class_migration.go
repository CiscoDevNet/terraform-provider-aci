package data

import (
	"fmt"
	"maps"
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
	Nested                 bool
	Attributes             []MigrationAttribute
	LegacyAttributes       []MigrationAttribute
}

func (a MigrationAttribute) IsAlias() bool {
	return a.FieldName != a.TargetField || a.ChildName != nil && !a.IsChildTypeChange() && !a.Nested
}

func (a MigrationAttribute) IsChildTypeChange() bool {
	return a.ChildName != nil && a.Type == SetNestedAttribute && !a.Block && !a.Nested
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
	if a.Nested {
		objectType := a.ElementAttributeType()
		if a.Type == SingleNestedAttribute {
			return objectType
		}
		return fmt.Sprintf("types.SetType{ElemType: %s}", objectType)
	}
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
	if !a.Block && !a.Nested {
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

func (c *Class) HasCurrentLegacyAttributes() bool {
	// Removed fields belong only to historical decoders, not current compatibility models or test scenarios.
	for _, attribute := range c.RootLegacyAttributes() {
		if attribute.Status != Removed {
			return true
		}
	}
	return false
}

func (c *Class) RootLegacyAttributes() []MigrationAttribute {
	// Class-owned compatibility fields are embedded only once, even when the resource has the same alias.
	var attributes []MigrationAttribute
	for _, attribute := range c.LegacyAttributes {
		if !slices.ContainsFunc(c.NestedLegacyAttributes, func(nested MigrationAttribute) bool { return nested.AttributeName == attribute.AttributeName }) {
			attributes = append(attributes, attribute)
		}
	}
	return attributes
}

func (c *Class) HasLegacyConfig() bool {
	return c.HasCurrentLegacyAttributes() || c.HasNestedLegacyChanges
}

func (c *Class) validateLegacyAliases(ctx *Context) {
	// Current configuration supports one previous name; older names remain available only to historical decoders.
	names := map[string]string{}
	for _, alias := range append(c.RootLegacyAttributes(), c.NestedLegacyAttributes...) {
		if alias.Status != LegacyDeprecated || alias.Computed && !alias.Optional && !alias.Required {
			continue
		}
		target := alias.ChildAttribute
		if alias.ChildName != nil {
			target = fmt.Sprintf("%s.%s", alias.ChildName, target)
		}
		if target == "" {
			continue
		}
		if previous, exists := names[target]; exists && previous != alias.AttributeName {
			ctx.Diagnostics.AddError("Class '%s': property %q has multiple active deprecated names %q and %q; mark the superseded name as removed in state_upgrades", c.Name, target, previous, alias.AttributeName)
			continue
		}
		names[target] = alias.AttributeName
	}
}

func (c *Class) LegacyLookupAttributes() []MigrationAttribute {
	// Datasource compatibility inputs are limited to the fields used to construct its DN.
	var attributes []MigrationAttribute
	for _, alias := range append(c.RootLegacyAttributes(), c.NestedLegacyAttributes...) {
		if alias.Status != LegacyDeprecated || alias.ChildName != nil || alias.Type != StringAttribute {
			continue
		}
		for _, name := range c.PropertiesAll {
			property := c.Properties[name]
			if property.AttributeName == alias.ChildAttribute && (property.PropertyName == "parentDn" || c.IsIdentifyingProperty(property.PropertyName)) {
				attributes = append(attributes, alias)
				break
			}
		}
	}
	return attributes
}

func (ds *DataStore) setNestedLegacyAttributes() {
	// Reused class models have one current shape; gather their active aliases from historical parent subtrees.
	aliases := map[string]map[string]MigrationAttribute{}
	var collect func([]MigrationAttribute)
	collect = func(attributes []MigrationAttribute) {
		for _, attribute := range attributes {
			if !attribute.Nested {
				continue
			}
			class := ds.Classes[attribute.ChildName.String()]
			for _, alias := range attribute.LegacyAttributes {
				if alias.Status == LegacyDeprecated {
					// Historical parent restrictions describe decoding; the owning class defines current alias inputs.
					for _, property := range class.ModelProperties() {
						if property.AttributeName == alias.ChildAttribute {
							alias.Required, alias.Optional, alias.Computed = false, !property.ReadOnly, true
							break
						}
					}
				}
				for _, own := range class.LegacyAttributes {
					if own.AttributeName == alias.AttributeName && (own.Status != alias.Status || alias.Status != Removed && (own.ModelAttributeType() != alias.ModelAttributeType() || own.ChildAttribute != alias.ChildAttribute)) {
						ds.ctx.Diagnostics.AddError("Class '%s': nested alias %q disagrees with its resource compatibility definition", class.Name, alias.AttributeName)
					}
				}
				if aliases[class.Name.String()] == nil {
					aliases[class.Name.String()] = map[string]MigrationAttribute{}
				}
				if previous, exists := aliases[class.Name.String()][alias.AttributeName]; exists {
					if previous.Status != alias.Status || alias.Status != Removed && (previous.ModelAttributeType() != alias.ModelAttributeType() || previous.ChildAttribute != alias.ChildAttribute || previous.Sensitive != alias.Sensitive) {
						ds.ctx.Diagnostics.AddError("Class '%s': nested alias %q has inconsistent current exposure across parent histories", class.Name, alias.AttributeName)
					}
				}
				aliases[class.Name.String()][alias.AttributeName] = alias
				if alias.Status == Removed || slices.ContainsFunc(class.NestedLegacyAttributes, func(previous MigrationAttribute) bool { return previous.AttributeName == alias.AttributeName }) {
					continue
				}
				class.NestedLegacyAttributes = append(class.NestedLegacyAttributes, alias)
			}
			slices.SortFunc(class.NestedLegacyAttributes, func(a, b MigrationAttribute) int { return strings.Compare(a.AttributeName, b.AttributeName) })
			ds.Classes[class.Name.String()] = class
			collect(attribute.Attributes)
		}
	}
	for _, class := range ds.Classes {
		for _, schema := range class.MigrationSchemas {
			collect(schema.Attributes)
		}
	}
	var containsLegacy func(Class, map[string]bool) bool
	containsLegacy = func(class Class, visited map[string]bool) bool {
		if len(class.NestedLegacyAttributes) > 0 {
			return true
		}
		if visited[class.Name.String()] {
			return false
		}
		visited[class.Name.String()] = true
		for _, name := range class.Children {
			if containsLegacy(ds.Classes[name.String()], visited) {
				return true
			}
		}
		return false
	}
	for name, class := range ds.Classes {
		class.HasNestedLegacyChanges = containsLegacy(class, map[string]bool{})
		ds.Classes[name] = class
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
	nested := c.nestedMigrationAttributes(ds, versions)
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
		names := map[string]bool{}
		addOverride := func(attribute MigrationAttribute, explicitName bool) {
			// Explicit duplicates are checked on the definition tree; inherited names become known only during normalization.
			if previousExplicit, exists := names[attribute.AttributeName]; exists && (!previousExplicit || !explicitName) {
				ds.ctx.Diagnostics.AddError("Class '%s': schema version %d: duplicate historical attribute name %q after resolving metadata names", c.Name, version, attribute.AttributeName)
			}
			names[attribute.AttributeName] = explicitName
			overrides[version] = append(overrides[version], attribute)
		}
		add := func(key string, definition AttributeUpgradeDefinition, child *Class) {
			storedAlias := false
			name := definition.AttributeName
			if name == "" {
				name = key
			}
			if previous, exists := aliases[name]; exists {
				if definition.Status == UndefinedLegacyStatus {
					definition.Status = previous.Status
				}
				storedAlias = entry.Source == FromFramework && child == nil && key == name && definition.AttributeName == name && definition.Status == Removed
				if storedAlias && previous.Block {
					if len(definition.Attributes) > 0 || len(definition.Children) > 0 {
						ds.ctx.Diagnostics.AddError("Class '%s': schema version %d: retired historical block %q cannot override members or nested children without nested-schema migration support", c.Name, version, name)
						return
					}
					if definition.Type != UndefinedLegacyAttributeType && definition.Type != previous.Type ||
						definition.ElementType != UndefinedLegacyAttributeType && definition.ElementType != previous.ElementType ||
						definition.ValueType != UndefinedValueType {
						ds.ctx.Diagnostics.AddError("Class '%s': schema version %d: retired historical block %q cannot change its inherited type without additional migration support", c.Name, version, name)
						return
					}
					// Keep the old block's members without replacing the modern child stored alongside it.
					previous.Status = definition.Status
					previous.ChildName, previous.TargetField, previous.ChildAttribute = nil, "", ""
					previous.TargetModelType, previous.TargetAttributeType = "", ""
					aliases[name] = previous
					addOverride(previous, true)
					return
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
			if storedAlias {
				// Retire the stored alias, even when its name is also an APIC property key such as prio.
				attribute.ChildName, attribute.TargetField, attribute.ChildAttribute = nil, "", ""
				attribute.TargetModelType, attribute.TargetAttributeType = "", ""
			}
			if attribute.IsAlias() {
				aliases[attribute.AttributeName] = attribute
			}
			addOverride(attribute, definition.AttributeName != "")
		}
		for key, definition := range entry.Attributes {
			add(key, definition, nil)
		}
		for name, definition := range entry.Children {
			if attribute, exists := nested[version][name]; exists {
				addOverride(attribute, definition.AttributeName != "")
				continue
			}
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
				addOverride(attribute, true)
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
				addOverride(attribute, false)
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

func (c *Class) nestedMigrationAttributes(ds *DataStore, versions []int) map[int]map[string]MigrationAttribute {
	// Existing SDKv2 child entries describe flat relations/blocks. Framework history describes nested attributes.
	result := map[int]map[string]MigrationAttribute{}
	histories := map[string]map[int]StateUpgradeDefinition{}
	for _, version := range versions {
		entry := c.StateUpgrades[version]
		for name, definition := range entry.Children {
			if entry.Source != FromFramework {
				continue
			}
			child, exists := ds.Classes[name]
			if !exists {
				continue
			}
			if definition.Type == SetNestedAttribute && child.IsSingleNestedWhenDefinedAsChild && len(definition.Attributes) == 0 && len(definition.Children) == 0 {
				// Preserve the published NetFlow set-to-singleton conversion.
				continue
			}
			if definition.Type != UndefinedLegacyAttributeType && definition.Type != SingleNestedAttribute && definition.Type != SetNestedAttribute {
				continue
			}
			if histories[name] == nil {
				histories[name] = map[int]StateUpgradeDefinition{}
			}
			histories[name][version] = StateUpgradeDefinition{Source: FromFramework, Attributes: definition.Attributes, Children: definition.Children}
		}
	}
	for name, history := range histories {
		child := ds.Classes[name]
		current, exists := c.currentMigrationAttributes(ds)[child.ResourceNameNested]
		if !exists {
			ds.ctx.Diagnostics.AddError("Class '%s': historical child %q has no current nested target; removing entire children needs additional migration support", c.Name, name)
			continue
		}
		child.Properties = maps.Clone(child.Properties)
		for name, property := range child.Properties {
			propertyCopy := *property
			propertyCopy.StateUpgradeValues = maps.Clone(property.StateUpgradeValues)
			child.Properties[name] = &propertyCopy
		}
		child.ClassDefinition.SchemaVersion = &c.SchemaVersion
		child.StateUpgrades, child.MigrationSchemas, child.LegacyAttributes = history, nil, nil
		child.Artifacts = []ArtifactEnum{ResourceArtifact}
		child.setMigrationData(ds)
		for _, schema := range child.MigrationSchemas {
			definition := c.StateUpgrades[schema.PriorSchemaVersion].Children[name]
			attribute := current
			attribute.Nested, attribute.ChildName = true, child.Name
			attribute.Type = definition.Type
			if attribute.Type == UndefinedLegacyAttributeType {
				attribute.Type = SetNestedAttribute
				if child.IsSingleNestedWhenDefinedAsChild {
					attribute.Type = SingleNestedAttribute
				}
			}
			if attribute.Type == SingleNestedAttribute && !child.IsSingleNestedWhenDefinedAsChild {
				ds.ctx.Diagnostics.AddError("Class '%s': historical child %q needs singleton-to-set migration support before it can be rendered", c.Name, name)
				continue
			}
			if definition.Status != UndefinedLegacyStatus {
				ds.ctx.Diagnostics.AddError("Class '%s': historical child %q needs nested child-name migration support for a container status", c.Name, name)
				continue
			}
			attribute.Required, attribute.Optional, attribute.Computed = child.RequiredAsChild, !child.RequiredAsChild, !child.RequiredAsChild
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
			attribute.InheritedModelType = "types.Set"
			if attribute.Type == SingleNestedAttribute {
				attribute.InheritedModelType = "types.Object"
			}
			attribute.TargetModelType, attribute.TargetAttributeType = current.InheritedModelType, current.InheritedAttributeType
			attribute.LegacyAttributes = child.LegacyAttributes
			for _, member := range schema.Attributes {
				if member.AttributeName != "id" && member.TargetField != "ParentDN" {
					attribute.Attributes = append(attribute.Attributes, member)
				}
			}
			if definition.AttributeName != "" && definition.AttributeName != attribute.AttributeName {
				ds.ctx.Diagnostics.AddError("Class '%s': historical child %q needs nested child-name migration support before it can be rendered", c.Name, name)
				continue
			}
			attribute.SchemaExpression = attribute.nestedSchemaExpression()
			if result[schema.PriorSchemaVersion] == nil {
				result[schema.PriorSchemaVersion] = map[string]MigrationAttribute{}
			}
			result[schema.PriorSchemaVersion][name] = attribute
		}
	}
	return result
}

func (a MigrationAttribute) nestedSchemaExpression() string {
	// Historical schemas need the old wire shape; current validators and plan modifiers do not participate in decoding.
	flags := fmt.Sprintf("Required: %t, Optional: %t, Computed: %t,", a.Required, a.Optional, a.Computed)
	if a.Type == SingleNestedAttribute {
		return fmt.Sprintf("resourceschema.ObjectAttribute{%s AttributeTypes: %s.AttrTypes}", flags, a.ElementAttributeType())
	}
	return fmt.Sprintf("resourceschema.SetAttribute{%s ElementType: %s}", flags, a.ElementAttributeType())
}

func (a MigrationAttribute) UpgradeSchemaExpression() string {
	var attributes strings.Builder
	for _, member := range a.Attributes {
		fmt.Fprintf(&attributes, "%q: %s,", member.AttributeName, member.UpgradeSchemaExpression())
	}
	return fmt.Sprintf("modelHelpers.LegacyAttributeSchema{Replacement: %q, Removed: %t, Attributes: map[string]modelHelpers.LegacyAttributeSchema{%s}}", a.ChildAttribute, a.Status == Removed && a.TargetField == "" && a.ChildAttribute == "", attributes.String())
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
		if definition.ValueType == SemanticEquality && (property == nil || property.ValueType != SemanticEquality) {
			ds.ctx.Diagnostics.AddError("Class '%s': historical attribute %q needs a historical semantic-string definition before it can be rendered", c.Name, key)
		} else if attribute.Type != StringAttribute || historical.ValueType.SchemaPropertyType() != "StringProperty" {
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
	if child == nil && property != nil && (property.ValueType.SchemaPropertyType() == "StringProperty" && attribute.Type != StringAttribute || property.ValueType == Set && attribute.Type != ListAttribute && attribute.Type != SetAttribute) {
		ds.ctx.Diagnostics.AddError("Class '%s': historical attribute %q needs an explicit conversion from %q to %q", c.Name, key, attribute.Type, property.ValueType)
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
