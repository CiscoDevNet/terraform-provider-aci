package data

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils"
)

// validSubCategories contains the allowed sub-category values for Terraform Registry sidebar grouping.
var validSubCategories = []string{
	"AAA",
	"Access Policies",
	"Application Management",
	"Cloud",
	"Contract",
	"Fabric Access Policies",
	"Fabric Inventory",
	"Fabric Policies",
	"Firmware",
	"Generic",
	"Import/Export",
	"L2Out",
	"L3Out",
	"L4-L7",
	"Monitoring",
	"Multi-Site",
	"Networking",
	"Node Management",
	"Scheduler",
	"System Settings",
	"Tenant Infra Policies",
	"Tenant Policies",
	"Virtual Networking",
}

// validateSubCategory validates that the given sub-category value is one of the allowed sub-categories.
func validateSubCategory(class *Class) error {
	if class.ClassDefinition.Documentation.SubCategory == "" {
		if class.IsSingleNestedWhenDefinedAsChild {
			return nil
		}
		return fmt.Errorf("sub_category not specified for class '%s': add documentation.sub_category to the class definition file", class.Name.full)
	}
	if !slices.Contains(validSubCategories, class.ClassDefinition.Documentation.SubCategory) {
		return fmt.Errorf("invalid sub_category '%s'", class.ClassDefinition.Documentation.SubCategory)
	}
	return nil
}

type ClassDocumentation struct {
	// A markdown link to the DevNet documentation page for the class, used to reference the class in the documentation.
	// Format: "[<className>](https://<host>/app/index.html#/objects/<className>/overview)".
	ClassName string
	// Child classes that have their own standalone resource but are not inline children.
	Children []string
	// A deprecation warning message for the documentation when the class is deprecated.
	DeprecationWarning string
	// The human-readable documentation label derived from class.ResourceName, with
	// DocumentationLabelOverrides applied (e.g., "Application EPG", "BGP Timers").
	// Resolved as ClassDocumentationDefinition.Label when set, otherwise humanized from class.ResourceName.
	Label string
	// The description line for the resource documentation (e.g., "Manages ACI Application EPG").
	ResourceDescription string
	// The description line for the datasource documentation (e.g., "Data source for ACI Application EPG").
	DatasourceDescription string
	// The description used when this class appears as a nested child in a parent resource.
	DescriptionWhenDefinedAsChild string
	// Completed target descriptions rendered below a multi-target child relation.
	RelationTargetDescriptions []string
	// DN format strings from meta file (e.g., "uni/tn-{name}").
	DnFormats []string
	// ImportDnFormat is the actual DN format used by import examples. It follows
	// the selected example parent or parent_dn default before falling back to the
	// first actual format and never contains a display-only notice.
	ImportDnFormat string
	// ParentDnNotice follows the sampled parent lists when the complete parent
	// set is too large to render. ParentDnResources contains links to parent
	// classes backed by a Terraform resource and ParentDnClasses contains
	// APIC-class links for the remaining parents.
	ParentDnNotice    string
	ParentDnResources []string
	ParentDnClasses   []string
	// A migration warning message for the documentation when the class has been migrated from a previous provider version.
	MigrationWarning string
	// The supported APIC versions string for display in documentation.
	SupportedVersions string
	// Notes rendered with -> prefix in the resource documentation.
	// Resolved as the shared ClassDocumentationDefinition.Notes followed by ClassDocumentationDefinition.Resource.Notes.
	ResourceNotes []string
	// Warnings rendered with ~> prefix in the resource documentation.
	// Resolved as the shared ClassDocumentationDefinition.Warnings followed by ClassDocumentationDefinition.Resource.Warnings.
	ResourceWarnings []string
	// Notes rendered with -> prefix in the datasource documentation.
	// Resolved as the shared ClassDocumentationDefinition.Notes followed by ClassDocumentationDefinition.Datasource.Notes.
	DatasourceNotes []string
	// Warnings rendered with ~> prefix in the datasource documentation.
	// Resolved as the shared ClassDocumentationDefinition.Warnings followed by ClassDocumentationDefinition.Datasource.Warnings.
	DatasourceWarnings []string
	// Sub-category for Terraform Registry sidebar grouping.
	SubCategory string
	// GUI locations in APIC (e.g., "Tenants -> Networking -> VRFs").
	UiLocations []string
}

// NestedDocumentationClass identifies one embedded class in depth-first
// documentation order. Depth is used only to indent the corresponding
// Markdown list entries.
type NestedDocumentationClass struct {
	Class             *Class
	Depth             int
	SupportedVersions string
}

// Indent returns four spaces for each nesting level plus additionalDepth.
func (c NestedDocumentationClass) Indent(additionalDepth int) string {
	return strings.Repeat("    ", c.Depth+additionalDepth)
}

// NestedDocumentationClasses resolves all embedded children in deterministic
// parent-first order without storing a second class hierarchy.
func (c *Class) NestedDocumentationClasses(ds *DataStore) []NestedDocumentationClass {
	return nestedDocumentationClasses(c, ds, 0)
}

// RequiredNestedDocumentationClasses resolves the complete subtrees rooted at
// children that are required in the owning resource schema.
func (c *Class) RequiredNestedDocumentationClasses(ds *DataStore) []NestedDocumentationClass {
	return nestedDocumentationClassesForSection(c, ds, true)
}

// OptionalNestedDocumentationClasses resolves the complete subtrees rooted at
// children that are optional in the owning resource schema.
func (c *Class) OptionalNestedDocumentationClasses(ds *DataStore) []NestedDocumentationClass {
	return nestedDocumentationClassesForSection(c, ds, false)
}

// RequiredDocumentationProperties returns required APIC properties in schema
// order. parentDn is rendered separately because its required/optional status
// depends on the top-level schema default and it is never part of a child model.
func (c *Class) RequiredDocumentationProperties() []*Property {
	return c.documentationProperties(c.PropertiesRequired)
}

// OptionalDocumentationProperties returns optional APIC properties in schema
// order, excluding the separately rendered parentDn property.
func (c *Class) OptionalDocumentationProperties() []*Property {
	return c.documentationProperties(c.PropertiesOptional)
}

// ReadOnlyDocumentationProperties returns read-only APIC properties in schema
// order, excluding the separately rendered parentDn property.
func (c *Class) ReadOnlyDocumentationProperties() []*Property {
	return c.documentationProperties(c.PropertiesReadOnly)
}

func (c *Class) documentationProperties(propertyNames []string) []*Property {
	properties := make([]*Property, 0, len(propertyNames))
	for _, propertyName := range propertyNames {
		if propertyName == "parentDn" {
			continue
		}
		properties = append(properties, c.Properties[propertyName])
	}
	return properties
}

func nestedDocumentationClassesForSection(parent *Class, ds *DataStore, required bool) []NestedDocumentationClass {
	var nestedClasses []NestedDocumentationClass
	for _, childName := range parent.Children {
		childClass, ok := ds.Classes[childName.String()]
		if !ok || childClass.RequiredAsChild != required {
			continue
		}
		child := &childClass
		nestedClasses = append(nestedClasses, nestedDocumentationClass(child, parent, 0))
		nestedClasses = append(nestedClasses, nestedDocumentationClasses(child, ds, 1)...)
	}
	return nestedClasses
}

func nestedDocumentationClasses(parent *Class, ds *DataStore, depth int) []NestedDocumentationClass {
	var nestedClasses []NestedDocumentationClass
	for _, childName := range parent.Children {
		childClass, ok := ds.Classes[childName.String()]
		if !ok {
			continue
		}
		child := &childClass
		nestedClasses = append(nestedClasses, nestedDocumentationClass(child, parent, depth))
		nestedClasses = append(nestedClasses, nestedDocumentationClasses(child, ds, depth+1)...)
	}
	return nestedClasses
}

func nestedDocumentationClass(child, parent *Class, depth int) NestedDocumentationClass {
	supportedVersions := ""
	if child.SupportedVersions != nil && !child.SupportedVersions.Covers(parent.SupportedVersions) {
		supportedVersions = documentationVersions(child.SupportedVersions)
	}
	return NestedDocumentationClass{
		Class:             child,
		Depth:             depth,
		SupportedVersions: supportedVersions,
	}
}

func (c *Class) setDocumentation(ds *DataStore) error {
	genLogger.Debugf("Setting Documentation for class '%s'.", c.Name)

	c.Documentation.setChildren(c, ds)

	c.Documentation.setDeprecationWarning(c)

	c.setPropertyDocumentationDescriptions(ds)
	c.setPropertyDocumentationReferences(ds)
	c.setPropertyDocumentationSupportedVersions()

	c.Documentation.setDescription(c)

	c.Documentation.setDescriptionWhenDefinedAsChild(c, ds)

	c.Documentation.setDnFormats(c)

	c.Documentation.setParentDns(c, ds)

	c.Documentation.setMigrationWarning(c)

	c.Documentation.setSupportedVersions(c)

	c.Documentation.setNotes(c)

	err := c.Documentation.setSubCategory(c)
	if err != nil {
		return err
	}

	err = c.Documentation.setUiLocations(c)
	if err != nil {
		return err
	}

	c.Documentation.setWarnings(c)

	genLogger.Debugf("Successfully set Documentation for class '%s'.", c.Name)
	return nil
}

func (d *ClassDocumentation) setClassName(class *Class) {
	genLogger.Debugf("Setting Documentation ClassName for class '%s'.", class.Name.full)

	d.ClassName = getDevnetObjectDocumentationLink(class.Name.full)

	genLogger.Debugf("Successfully set Documentation ClassName for class '%s'. ClassName: %s", class.Name.full, d.ClassName)
}

func getDevnetObjectDocumentationLink(className string) string {
	return fmt.Sprintf("[%s](https://%s/app/index.html#/objects/%s/overview)", className, constPubhubDevnetHost, className)
}

func getTerraformRegistryDocumentationLink(resourceName string, artifact ArtifactEnum) string {
	return getTerraformRegistryDocumentationLinkAtPath(resourceName, resourceName, artifact)
}

func getTerraformRegistryDocumentationLinkAtPath(resourceName, registryPath string, artifact ArtifactEnum) string {
	url := getTerraformRegistryDocumentationURL(registryPath, artifact)
	if url == "" {
		return ""
	}

	return fmt.Sprintf("[%s_%s](%s)", constProviderName, resourceName, url)
}

func getTerraformRegistryDocumentationURL(registryPath string, artifact ArtifactEnum) string {
	var baseURL string
	switch artifact {
	case ResourceArtifact:
		baseURL = constRegistryResourceBaseUrl
	case DatasourceArtifact:
		baseURL = constRegistryDatasourceBaseUrl
	default:
		return ""
	}

	return fmt.Sprintf("%s/%s", baseURL, registryPath)
}

func (ds *DataStore) getTerraformRegistryDocumentationReference(className string, artifact ArtifactEnum) (string, string) {
	resourceName := ds.getResourceNameForArtifact(className, artifact)
	if resourceName == "" {
		return "", ""
	}

	registryPath := resourceName
	if definition, ok := ds.ClassDefinitions[className]; ok && definition.Documentation.RegistryPath != "" {
		registryPath = definition.Documentation.RegistryPath
	}

	return resourceName, getTerraformRegistryDocumentationURL(registryPath, artifact)
}

func (ds *DataStore) getTerraformRegistryDocumentationLink(className string, artifact ArtifactEnum) string {
	resourceName, url := ds.getTerraformRegistryDocumentationReference(className, artifact)
	if url == "" {
		return ""
	}

	return fmt.Sprintf("[%s_%s](%s)", constProviderName, resourceName, url)
}

func getResourceDocumentationLink(resourceName string) string {
	return getTerraformRegistryDocumentationLink(resourceName, ResourceArtifact)
}

func getDatasourceDocumentationLink(resourceName string) string {
	return getTerraformRegistryDocumentationLink(resourceName, DatasourceArtifact)
}

func (d *ClassDocumentation) setChildren(class *Class, ds *DataStore) {
	genLogger.Debugf("Setting Documentation Children for class '%s'.", class.Name.full)

	// Build a set of children already embedded as nested attributes in this resource
	// so they can be excluded from the documentation children list.
	childrenIncludedInResource := make(map[string]struct{}, len(class.Children))
	for _, child := range class.Children {
		childrenIncludedInResource[child.full] = struct{}{}
	}

	documentationIncludes := make(map[string]struct{}, len(class.ClassDefinition.Documentation.IncludeChildren))
	for _, childName := range class.ClassDefinition.Documentation.IncludeChildren {
		normalizedChildName, err := sanitizeClassName(childName)
		if err != nil {
			genLogger.Warnf("Skipping invalid documentation include child class name '%s' for class '%s': %s", childName, class.Name.full, err)
			continue
		}
		documentationIncludes[normalizedChildName] = struct{}{}
	}

	documentationExcludes := make(map[string]struct{}, len(class.ClassDefinition.Documentation.ExcludeChildren))
	for _, childName := range class.ClassDefinition.Documentation.ExcludeChildren {
		normalizedChildName, err := sanitizeClassName(childName)
		if err != nil {
			genLogger.Warnf("Skipping invalid documentation exclude child class name '%s' for class '%s': %s", childName, class.Name.full, err)
			continue
		}
		documentationExcludes[normalizedChildName] = struct{}{}
	}

	linksByChild := make(map[string]string)
	rnMap, hasRnMap := class.MetaFileContent["rnMap"].(map[string]any)
	if hasRnMap {
		for _, classNameInterface := range rnMap {
			className, ok := classNameInterface.(string)
			if !ok {
				genLogger.Warnf("Skipping non-string child class name in rnMap for class '%s'", class.Name.full)
				continue
			}
			childName, err := sanitizeClassName(className)
			if err != nil {
				genLogger.Warnf("Skipping invalid child class name in rnMap for class '%s': %s", class.Name.full, err)
				continue
			}
			_, forceInclude := documentationIncludes[childName]
			if link, ok := documentationChildLink(childName, forceInclude, class, ds, documentationExcludes, childrenIncludedInResource); ok {
				linksByChild[childName] = link
			}
		}
	} else {
		genLogger.Debugf("No rnMap available for class '%s'; checking documentation include children.", class.Name.full)
	}

	for childName := range documentationIncludes {
		if link, ok := documentationChildLink(childName, true, class, ds, documentationExcludes, childrenIncludedInResource); ok {
			linksByChild[childName] = link
		}
	}

	var links []string
	if hasRnMap || len(documentationIncludes) > 0 {
		links = make([]string, 0, len(linksByChild))
	}
	for _, link := range linksByChild {
		links = append(links, link)
	}

	slices.Sort(links)
	d.Children = slices.Compact(links)

	genLogger.Debugf("Successfully set Documentation Children for class '%s'. Children: %v", class.Name.full, d.Children)
}

func documentationChildLink(
	childName string,
	forceInclude bool,
	class *Class,
	ds *DataStore,
	documentationExcludes map[string]struct{},
	childrenIncludedInResource map[string]struct{},
) (string, bool) {
	if !forceInclude {
		if _, excluded := documentationExcludes[childName]; excluded {
			return "", false
		}
		if slices.Contains(class.ClassDefinition.ExcludeChildren, childName) {
			return "", false
		}
		if _, embedded := childrenIncludedInResource[childName]; embedded {
			return "", false
		}
	}

	link := ds.getTerraformRegistryDocumentationLink(childName, ResourceArtifact)
	if link == "" {
		return "", false
	}

	return link, true
}

func (d *ClassDocumentation) setDeprecationWarning(class *Class) {
	genLogger.Debugf("Setting Documentation DeprecationWarning for class '%s'.", class.Name.full)

	var warningParts []string

	switch {
	case class.Hidden:
		if class.HiddenVersions != nil {
			warningParts = append(warningParts, fmt.Sprintf("The %s class is no longer accepted by the APIC API as of version %s.", d.ClassName, class.HiddenVersions))
		} else {
			warningParts = append(warningParts, fmt.Sprintf("The %s class is no longer accepted by the APIC API.", d.ClassName))
		}
	case class.Deprecated:
		if class.DeprecatedVersions != nil {
			warningParts = append(warningParts, fmt.Sprintf("The %s class is deprecated by the APIC API as of version %s and may be removed in a future APIC release.", d.ClassName, class.DeprecatedVersions))
		} else {
			warningParts = append(warningParts, fmt.Sprintf("The %s class is deprecated by the APIC API and may be removed in a future APIC release.", d.ClassName))
		}
	}

	if class.Hidden || class.Deprecated {
		warningParts = append(warningParts, fmt.Sprintf("The `aci_%s` resource and datasource are deprecated and will be removed in a future release of the provider.", class.ResourceName))
	}

	d.DeprecationWarning = strings.Join(warningParts, " ")

	genLogger.Debugf("Successfully set Documentation DeprecationWarning for class '%s'. DeprecationWarning: %s", class.Name.full, d.DeprecationWarning)
}

func (d *ClassDocumentation) setLabel(class *Class, ds *DataStore) {
	genLogger.Debugf("Setting Documentation Label for class '%s'.", class.Name.full)

	if class.ClassDefinition.Documentation.Label != "" {
		d.Label = class.ClassDefinition.Documentation.Label
	} else {
		d.Label = documentationLabel(class.ResourceName, ds.GlobalMetaDefinition.DocumentationLabelOverrides)
	}

	genLogger.Debugf("Successfully set Documentation Label for class '%s'. Label: %s", class.Name.full, d.Label)
}

func documentationLabel(resourceName string, overrides map[string]string) string {
	label := cases.Title(language.English).String(strings.ReplaceAll(resourceName, "_", " "))
	// Apply multi-word substitutions as phrases and single-word substitutions
	// only to complete words so short acronyms do not alter longer words.
	for key, replacement := range overrides {
		if strings.Contains(key, " ") {
			label = strings.ReplaceAll(label, key, replacement)
		} else if slices.Contains(strings.Split(label, " "), key) {
			label = strings.ReplaceAll(label, key, replacement)
		}
	}
	return label
}

func documentationClassLabel(className string, ds *DataStore) string {
	if class, ok := ds.Classes[className]; ok && class.Documentation.Label != "" {
		return class.Documentation.Label
	}
	if definition, ok := ds.ClassDefinitions[className]; ok {
		if definition.Documentation.Label != "" {
			return definition.Documentation.Label
		}
		if definition.ResourceName != "" {
			return documentationLabel(definition.ResourceName, ds.GlobalMetaDefinition.DocumentationLabelOverrides)
		}
	}
	return className
}

func (d *ClassDocumentation) setDescription(class *Class) {
	genLogger.Debugf("Setting Documentation Description for class '%s'.", class.Name.full)

	docDef := class.ClassDefinition.Documentation

	// Build the resource and datasource description lines. Keep the standard description without
	// terminal punctuation unless more text follows, matching the existing generated documentation.
	resourceParts := []string{fmt.Sprintf("Manages ACI %s", d.Label)}
	if docDef.Description != "" || docDef.Resource.Description != "" {
		resourceParts[0] += "."
	}
	if docDef.Description != "" {
		resourceParts = append(resourceParts, docDef.Description)
	}
	if docDef.Resource.Description != "" {
		resourceParts = append(resourceParts, docDef.Resource.Description)
	}
	d.ResourceDescription = strings.Join(resourceParts, " ")

	datasourceParts := []string{fmt.Sprintf("Data source for ACI %s", d.Label)}
	if docDef.Description != "" || docDef.Datasource.Description != "" {
		datasourceParts[0] += "."
	}
	if docDef.Description != "" {
		datasourceParts = append(datasourceParts, docDef.Description)
	}
	if docDef.Datasource.Description != "" {
		datasourceParts = append(datasourceParts, docDef.Datasource.Description)
	}
	d.DatasourceDescription = strings.Join(datasourceParts, " ")

	genLogger.Debugf("Successfully set Documentation Description for class '%s'. ResourceDescription: %s, DatasourceDescription: %s", class.Name.full, d.ResourceDescription, d.DatasourceDescription)
}

func (d *ClassDocumentation) setDescriptionWhenDefinedAsChild(class *Class, ds *DataStore) {
	genLogger.Debugf("Setting Documentation DescriptionWhenDefinedAsChild for class '%s'.", class.Name.full)
	d.RelationTargetDescriptions = nil

	nestingType := "list"
	nestedLabel := utils.Plural(d.Label)
	if class.IsSingleNestedWhenDefinedAsChild {
		nestingType = "map"
		nestedLabel = d.Label
	}

	header := fmt.Sprintf("`%s` - (%s)", class.ResourceNameNested, nestingType)
	childDescription := fmt.Sprintf("%s (ACI object %s)", nestedLabel, d.ClassName)

	var sentence string
	if class.Relation.RelationalClass {
		if len(class.Relation.ToClasses) > 1 {
			for _, toClassName := range documentationRelationTargets(class) {
				toClassFull := toClassName.String()
				targetDescription := getDevnetObjectDocumentationLink(toClassFull)
				if resourceLink := ds.getTerraformRegistryDocumentationLink(toClassFull, ResourceArtifact); resourceLink != "" {
					targetDescription += fmt.Sprintf(" which can be configured using the %s resource", resourceLink)
				} else {
					targetDescription += ". Currently there is no resource available to configure this ACI object"
				}
				d.RelationTargetDescriptions = append(d.RelationTargetDescriptions, targetDescription)
			}
			sentence = fmt.Sprintf("A %s of %s.", nestingType, childDescription)
			d.DescriptionWhenDefinedAsChild = header + " " + sentence
			genLogger.Debugf("Successfully set Documentation DescriptionWhenDefinedAsChild for class '%s'. DescriptionWhenDefinedAsChild: %s", class.Name.full, d.DescriptionWhenDefinedAsChild)
			return
		}

		targetParts := make([]string, 0, len(class.Relation.ToClasses))
		resourceParts := make([]string, 0, len(class.Relation.ToClasses))
		for _, toClassName := range class.Relation.ToClasses {
			toClassFull := toClassName.String()
			toClassLink := getDevnetObjectDocumentationLink(toClassFull)
			toLabel := documentationClassLabel(toClassFull, ds)
			targetParts = append(targetParts, fmt.Sprintf("%s (ACI Object %s)", toLabel, toClassLink))
			if resourceLink := ds.getTerraformRegistryDocumentationLink(toClassFull, ResourceArtifact); resourceLink != "" {
				resourceParts = append(resourceParts, resourceLink)
			}
		}
		targets := strings.Join(targetParts, ", ")
		// Relation cardinality controls only the nested schema shape. A relation target
		// remains independently configurable when it has a resource artifact.
		if len(resourceParts) > 0 && len(resourceParts) == len(class.Relation.ToClasses) {
			resourceWord := "resource"
			if len(resourceParts) > 1 {
				resourceWord = "resources"
			}
			sentence = fmt.Sprintf("A %s of %s pointing to %s which can be configured using the %s %s.", nestingType, childDescription, targets, strings.Join(resourceParts, ", "), resourceWord)
		} else {
			sentence = fmt.Sprintf("A %s of %s pointing to %s.", nestingType, childDescription, targets)
		}
	} else if class.HasResourceArtifact() {
		registryPath := class.ResourceName
		if class.ClassDefinition.Documentation.RegistryPath != "" {
			registryPath = class.ClassDefinition.Documentation.RegistryPath
		}
		resourceLink := getTerraformRegistryDocumentationLinkAtPath(class.ResourceName, registryPath, ResourceArtifact)
		sentence = fmt.Sprintf("A %s of %s which can also be configured using a separate %s resource.", nestingType, childDescription, resourceLink)
	} else {
		sentence = fmt.Sprintf("A %s of %s.", nestingType, childDescription)
	}

	d.DescriptionWhenDefinedAsChild = header + " " + sentence

	genLogger.Debugf("Successfully set Documentation DescriptionWhenDefinedAsChild for class '%s'. DescriptionWhenDefinedAsChild: %s", class.Name.full, d.DescriptionWhenDefinedAsChild)
}

func documentationRelationTargets(class *Class) []*ClassName {
	definedTargets := class.ClassDefinition.RelationInfo.ToClasses
	if len(definedTargets) != len(class.Relation.ToClasses) {
		return class.Relation.ToClasses
	}

	orderedTargets := make([]*ClassName, 0, len(definedTargets))
	for _, definedTarget := range definedTargets {
		target, err := NewClassName(definedTarget)
		if err != nil {
			return class.Relation.ToClasses
		}
		orderedTargets = append(orderedTargets, target)
	}
	return orderedTargets
}

func (d *ClassDocumentation) setDnFormats(class *Class) {
	genLogger.Debugf("Setting Documentation DnFormats for class '%s'.", class.Name.full)

	d.ImportDnFormat = ""
	d.DnFormats = class.documentationDnFormats()
	if len(d.DnFormats) == 0 {
		genLogger.Debugf("No dnFormats available for class '%s'.", class.Name.full)
	}

	if dnFormats := actualDnFormats(d.DnFormats); len(dnFormats) > 0 {
		d.ImportDnFormat = dnFormats[0]
	}
	if len(class.ClassDefinition.Documentation.DnFormats) == 0 && len(d.DnFormats) > constMaxDnFormatsToDisplay {
		notice := fmt.Sprintf("Too many DN formats to display, see model documentation for all possible parents of %s.", d.ClassName)
		d.DnFormats = append([]string{notice}, d.DnFormats[:constMaxDnFormatsToDisplay]...)
	}

	genLogger.Debugf("Successfully set Documentation DnFormats for class '%s'. DnFormats: %v, ImportDnFormat: %s", class.Name.full, d.DnFormats, d.ImportDnFormat)
}

func (c *Class) documentationDnFormats() []string {
	if override := c.ClassDefinition.Documentation.DnFormats; len(override) > 0 {
		dnFormats := slices.Clone(override)
		slices.Sort(dnFormats)
		return dnFormats
	}

	var dnFormats []string
	if rawFormats, ok := c.MetaFileContent["dnFormats"].([]any); ok {
		for _, raw := range rawFormats {
			if dnFormat, ok := raw.(string); ok {
				dnFormats = append(dnFormats, dnFormat)
			}
		}
	}
	slices.Sort(dnFormats)
	return dnFormats
}

func (d *ClassDocumentation) setImportDnFormat(class *Class, ds *DataStore) {
	d.ImportDnFormat = class.documentationImportDnFormat(ds, make(map[string]struct{}))
}

func (c *Class) documentationImportDnFormat(ds *DataStore, visited map[string]struct{}) string {
	// Select from the complete source list because the rendered DN formats can be
	// shortened to a representative sample.
	dnFormats := actualDnFormats(c.documentationDnFormats())
	if len(dnFormats) == 0 {
		return ""
	}
	fallback := dnFormats[0]

	if c.Name != nil {
		className := c.Name.String()
		if _, found := visited[className]; found {
			return fallback
		}
		visited[className] = struct{}{}
		defer delete(visited, className)
	}

	var parentDnFormat string
	if c.ExampleContext.ParentDnMode == ParentDnFromSchemaDefault {
		if c.DefaultParentDn != nil {
			parentDnFormat = c.DefaultParentDn.ParentDn
		}
		if parentDnFormat == "" {
			if parentDn := c.ParentDnProperty(); parentDn != nil {
				for _, defaultValue := range parentDn.Documentation.DefaultValues {
					if defaultValue.Value != "" {
						parentDnFormat = defaultValue.Value
						break
					}
				}
			}
		}
	}
	if parentDnFormat == "" && c.ExampleContext.ParentClass != nil {
		if parent, ok := ds.Classes[c.ExampleContext.ParentClass.String()]; ok {
			parentDnFormat = parent.documentationImportDnFormat(ds, visited)
		}
	}
	if parentDnFormat == "" {
		return fallback
	}

	expectedDnFormat := strings.TrimSuffix(parentDnFormat, "/") + "/" + c.RnFormat
	for _, dnFormat := range dnFormats {
		if documentationDnFormatMatches(dnFormat, expectedDnFormat) {
			return expectedDnFormat
		}
	}
	return fallback
}

func documentationDnFormatMatches(pattern, value string) bool {
	if pattern == value {
		return true
	}

	var expression strings.Builder
	expression.WriteByte('^')
	for {
		start := strings.IndexByte(pattern, '{')
		if start == -1 {
			expression.WriteString(regexp.QuoteMeta(pattern))
			break
		}
		end := strings.IndexByte(pattern[start:], '}')
		if end == -1 {
			expression.WriteString(regexp.QuoteMeta(pattern))
			break
		}
		expression.WriteString(regexp.QuoteMeta(pattern[:start]))
		expression.WriteString(`[^/]+`)
		pattern = pattern[start+end+1:]
	}
	expression.WriteByte('$')

	matches, err := regexp.MatchString(expression.String(), value)
	return err == nil && matches
}

func actualDnFormats(dnFormats []string) []string {
	actualFormats := make([]string, 0, len(dnFormats))
	for _, dnFormat := range dnFormats {
		if !strings.HasPrefix(dnFormat, "Too many DN formats") {
			actualFormats = append(actualFormats, dnFormat)
		}
	}
	return actualFormats
}

func (d *ClassDocumentation) setMigrationWarning(class *Class) {
	genLogger.Debugf("Setting Documentation MigrationWarning for class '%s'.", class.Name.full)

	if class.MigrationSource != UndefinedMigrationSource {
		d.MigrationWarning = "This resource has been migrated to the terraform plugin protocol version 6, refer to the [migration guide](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/guides/migration) for more details and implications for already managed resources."
	}

	genLogger.Debugf("Successfully set Documentation MigrationWarning for class '%s'. MigrationWarning: %s", class.Name.full, d.MigrationWarning)
}

func (d *ClassDocumentation) setParentDns(class *Class, ds *DataStore) {
	genLogger.Debugf("Setting Documentation ParentDns for class '%s'.", class.Name.full)

	parents := class.Parents
	if len(class.ParentDnVariants) > 0 {
		parents = make([]*ClassName, 0, len(class.ParentDnVariants))
		seen := make(map[string]struct{}, len(class.ParentDnVariants))
		for _, variant := range class.ParentDnVariants {
			if variant.ParentClass == nil {
				continue
			}
			parentName := variant.ParentClass.String()
			if _, ok := seen[parentName]; ok {
				continue
			}
			seen[parentName] = struct{}{}
			parents = append(parents, variant.ParentClass)
		}
	}

	d.ParentDnNotice = ""
	d.ParentDnResources = nil
	d.ParentDnClasses = nil

	var resourceEntries, classOnlyEntries []string
	for _, parent := range parents {
		parentLink := getDevnetObjectDocumentationLink(parent.full)
		if resourceLink := ds.getTerraformRegistryDocumentationLink(parent.full, ResourceArtifact); resourceLink != "" {
			resourceEntries = append(resourceEntries, fmt.Sprintf("%s (%s)", resourceLink, parentLink))
		} else {
			classOnlyEntries = append(classOnlyEntries, parentLink)
		}
	}

	// The source parent list is deterministic, so the partitioned slices preserve
	// its order without an additional sort.

	if len(parents) > constMaxParentDnsToDisplay {
		d.ParentDnNotice = fmt.Sprintf("Additional parent options are available; see model documentation for all possible parents of %s.", d.ClassName)
		resourceCount := min(len(resourceEntries), constMaxParentDnsToDisplay)
		d.ParentDnResources = resourceEntries[:resourceCount]
		remaining := constMaxParentDnsToDisplay - resourceCount
		if remaining > 0 {
			d.ParentDnClasses = classOnlyEntries[:min(len(classOnlyEntries), remaining)]
		}
	} else {
		d.ParentDnResources = resourceEntries
		d.ParentDnClasses = classOnlyEntries
	}

	genLogger.Debugf("Successfully set Documentation ParentDns for class '%s'. Notice: %s, Resources: %v, Classes: %v", class.Name.full, d.ParentDnNotice, d.ParentDnResources, d.ParentDnClasses)
}

func (d *ClassDocumentation) setSupportedVersions(class *Class) {
	genLogger.Debugf("Setting Documentation SupportedVersions for class '%s'.", class.Name.full)

	if class.SupportedVersions != nil {
		d.SupportedVersions = documentationVersions(class.SupportedVersions)
	}

	genLogger.Debugf("Successfully set Documentation SupportedVersions for class '%s'. SupportedVersions: %s", class.Name.full, d.SupportedVersions)
}

func (d *ClassDocumentation) setNotes(class *Class) {
	genLogger.Debugf("Setting Documentation Notes for class '%s'.", class.Name.full)

	docDef := class.ClassDefinition.Documentation
	d.ResourceNotes = slices.Concat(docDef.Notes, docDef.Resource.Notes)
	d.DatasourceNotes = slices.Concat(docDef.Notes, docDef.Datasource.Notes)

	genLogger.Debugf("Successfully set Documentation Notes for class '%s'. ResourceNotes: %v, DatasourceNotes: %v", class.Name.full, d.ResourceNotes, d.DatasourceNotes)
}

func (d *ClassDocumentation) setUiLocations(class *Class) error {
	genLogger.Debugf("Setting Documentation UiLocations for class '%s'.", class.Name.full)

	if len(class.ClassDefinition.Documentation.UiLocations) == 0 {
		if class.IsSingleNestedWhenDefinedAsChild {
			return nil
		}
		return fmt.Errorf("class '%s': ui_locations not specified: add documentation.ui_locations to the class definition file", class.Name.full)
	}

	d.UiLocations = class.ClassDefinition.Documentation.UiLocations

	genLogger.Debugf("Successfully set Documentation UiLocations for class '%s'. UiLocations: %v", class.Name.full, d.UiLocations)
	return nil
}

func (d *ClassDocumentation) setWarnings(class *Class) {
	genLogger.Debugf("Setting Documentation Warnings for class '%s'.", class.Name.full)

	docDef := class.ClassDefinition.Documentation
	d.ResourceWarnings = slices.Concat(docDef.Warnings, docDef.Resource.Warnings)
	d.DatasourceWarnings = slices.Concat(docDef.Warnings, docDef.Datasource.Warnings)

	genLogger.Debugf("Successfully set Documentation Warnings for class '%s'. ResourceWarnings: %v, DatasourceWarnings: %v", class.Name.full, d.ResourceWarnings, d.DatasourceWarnings)
}

func (d *ClassDocumentation) setSubCategory(class *Class) error {
	genLogger.Debugf("Setting Documentation SubCategory for class '%s'.", class.Name.full)

	if err := validateSubCategory(class); err != nil {
		return fmt.Errorf("class '%s': %w", class.Name.full, err)
	}

	d.SubCategory = class.ClassDefinition.Documentation.SubCategory

	genLogger.Debugf("Successfully set Documentation SubCategory '%s' for class '%s'.", d.SubCategory, class.Name.full)
	return nil
}
