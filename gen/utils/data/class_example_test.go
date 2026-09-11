package data

import (
	"reflect"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

func TestSetExampleDataBuildsIndependentExampleArtifacts(t *testing.T) {
	t.Parallel()

	name := &Property{
		PropertyName:  "name",
		AttributeName: "name",
		Required:      true,
		ValueType:     String,
		TestValues: &TestValues{
			Default: newTestValueScenario([]TestValueEntry{{ConfigValue: "minimum_name", ConfigInclude: true}}),
			Create:  newTestValueScenario([]TestValueEntry{{ConfigValue: "full_name", ConfigInclude: true}}),
		},
	}
	description := &Property{
		PropertyName:  "descr",
		AttributeName: "description",
		Optional:      true,
		ValueType:     String,
		TestValues: &TestValues{
			Default: newTestValueScenario([]TestValueEntry{{ConfigValue: "", ConfigInclude: false, AssertValue: "server_default"}}),
			Create:  newTestValueScenario([]TestValueEntry{{ConfigValue: "full_description", ConfigInclude: true}}),
		},
	}
	requiredMode := &Property{
		PropertyName:  "mode",
		AttributeName: "mode",
		Required:      true,
		ValueType:     String,
		TestValues: &TestValues{
			Default: newTestValueScenario([]TestValueEntry{{ConfigValue: "minimum_mode", ConfigInclude: true}}),
			Create:  newTestValueScenario([]TestValueEntry{{ConfigValue: "full_mode", ConfigInclude: true}}),
		},
	}
	scope := &Property{
		PropertyName:  "scope",
		AttributeName: "scope",
		Required:      true,
		ValueType:     String,
		TestValues: &TestValues{
			Default: newTestValueScenario([]TestValueEntry{{ConfigValue: "uni/tn-test", ConfigInclude: true}}),
			Create:  newTestValueScenario([]TestValueEntry{{ConfigValue: "uni/tn-test", ConfigInclude: true}}),
		},
		Reference: &PropertyReference{
			Classes:              []*ClassName{testClassName("fvBD")},
			ExampleClass:         testClassName("fvBD"),
			ExampleAttributeName: "id",
		},
	}
	ignored := &Property{
		PropertyName:  "ignored",
		AttributeName: "ignored",
		Optional:      true,
		IgnoreInTest:  true,
		ValueType:     String,
		ValidValues:   ValidValues{"enabled": {LocalName: "enabled"}},
	}
	members := &Property{
		PropertyName:  "members",
		AttributeName: "members",
		Optional:      true,
		ValueType:     Set,
		TestValues: &TestValues{
			Default: newTestValueScenario([]TestValueEntry{{ConfigInclude: false}}),
			Create:  newTestValueScenario(nil),
		},
	}

	childClass := Class{
		Name:                             testClassName("tagAnnotation"),
		ResourceNameNested:               "annotations",
		RequiredAsChild:                  true,
		IsSingleNestedWhenDefinedAsChild: false,
		Properties: map[string]*Property{
			"key": {
				PropertyName:  "key",
				AttributeName: "key",
				Required:      true,
				ValueType:     String,
			},
			"value": {
				PropertyName:  "value",
				AttributeName: "value",
				Optional:      true,
				ValueType:     String,
			},
		},
		PropertiesAll: []string{"key", "value"},
	}

	resourceClass := Class{
		Name:            testClassName("exampleClass"),
		ResourceName:    "example_class",
		Artifacts:       []ArtifactEnum{ResourceArtifact, DatasourceArtifact},
		ClassDefinition: ClassDefinition{ExampleFiles: []string{"example.pem"}},
		IdentifiedBy:    []string{"name"},
		Properties: map[string]*Property{
			"name":    name,
			"descr":   description,
			"scope":   scope,
			"mode":    requiredMode,
			"ignored": ignored,
			"members": members,
		},
		PropertiesAll: []string{"name", "descr", "scope", "mode", "ignored", "members"},
		TestChildren: []*TestChild{{
			Class: childClass.Name,
			Instances: []TestChildInstance{{Properties: map[string]TestValueScenario{
				"key":   newTestValueScenario([]TestValueEntry{{ConfigValue: "key_1", ConfigInclude: true}}),
				"value": newTestValueScenario([]TestValueEntry{{ConfigValue: "value_1", ConfigInclude: true}}),
			}}},
		}},
	}
	targetClass := Class{
		Name:         testClassName("fvBD"),
		ResourceName: "bridge_domain",
		Artifacts:    []ArtifactEnum{ResourceArtifact},
	}
	ds := &DataStore{Classes: map[string]Class{
		resourceClass.Name.String(): resourceClass,
		childClass.Name.String():    childClass,
		targetClass.Name.String():   targetClass,
	}}

	if err := ds.setExampleData(); err != nil {
		t.Fatal(err)
	}
	resolved := ds.Classes[resourceClass.Name.String()]
	if !resolved.ResourceExamples.HasDistinctFullExample() {
		t.Fatal("all-attributes example should differ from the minimum example")
	}
	if !reflect.DeepEqual(resolved.ResourceExamples.Files, []string{"example.pem"}) {
		t.Fatalf("example files %#v, expected [example.pem]", resolved.ResourceExamples.Files)
	}

	minimum := resolved.ResourceExamples.Minimum
	assertExampleAttributes(t, minimum.Body.Attributes, map[string]string{
		"name":  `"minimum_name"`,
		"scope": "aci_bridge_domain.example.id",
		"mode":  `"minimum_mode"`,
	})
	if len(minimum.Body.Children) != 1 {
		t.Fatalf("unexpected minimum children: %#v", minimum.Body.Children)
	}
	assertExampleAttributes(t, minimum.Body.Children[0].Body.Attributes, map[string]string{"key": `"example_key"`})

	full := resolved.ResourceExamples.Full
	assertExampleAttributes(t, full.Body.Attributes, map[string]string{
		"name":        `"minimum_name"`,
		"description": `"full_description"`,
		"scope":       "aci_bridge_domain.example.id",
		"mode":        `"minimum_mode"`,
		"ignored":     `"enabled"`,
		"members":     "[]",
	})
	assertExampleAttributes(t, full.Body.Children[0].Body.Attributes, map[string]string{
		"key":   `"example_key"`,
		"value": `"example_value"`,
	})

	assertExampleAttributes(t, resolved.DataSourceExample.Body.Attributes, map[string]string{"name": `"minimum_name"`})
	if actual := scope.TestValues.Create.Entries[0].ConfigValue; actual != "uni/tn-test" {
		t.Fatalf("example normalization changed test input to %q", actual)
	}
}

func TestDataSourceExampleDoesNotUseFullOnlyOverride(t *testing.T) {
	t.Parallel()

	name := &Property{
		PropertyName:  "name",
		AttributeName: "name",
		Optional:      true,
		ValueType:     String,
		TestValues: &TestValues{
			Default: newTestValueScenario([]TestValueEntry{{ConfigInclude: false}}),
			Create:  newTestValueScenario([]TestValueEntry{{ConfigValue: "create_name", ConfigInclude: true}}),
		},
		propertyDefinition: PropertyDefinition{ExampleConfig: PropertyExampleConfigDefinition{
			Full: ExampleValueScenarioDefinition{
				Defined: true,
				Entries: []ExampleValueEntryDefinition{{Value: "full_only_name"}},
			},
		}},
	}
	class := Class{
		Name:          testClassName("exampleClass"),
		Artifacts:     []ArtifactEnum{DatasourceArtifact},
		IdentifiedBy:  []string{"name"},
		Properties:    map[string]*Property{"name": name},
		PropertiesAll: []string{"name"},
	}

	class.setDataSourceExample(&DataStore{})
	if len(class.DataSourceExample.Body.Attributes) != 0 {
		t.Fatalf("datasource example unexpectedly used full-only values: %#v", class.DataSourceExample.Body.Attributes)
	}
}

func TestChildExampleUsesResolvedInstancePropertyValues(t *testing.T) {
	t.Parallel()

	mode := &Property{
		PropertyName:  "mode",
		AttributeName: "mode",
		Optional:      true,
		ValueType:     String,
		TestValues: &TestValues{
			Default: newTestValueScenario([]TestValueEntry{{ConfigInclude: false}}),
			Create:  newTestValueScenario([]TestValueEntry{{ConfigValue: "class_example", ConfigInclude: true}}),
		},
	}
	child := Class{
		Name:          testClassName("exampleChild"),
		Properties:    map[string]*Property{"mode": mode},
		PropertiesAll: []string{"mode"},
	}
	ds := &DataStore{Classes: map[string]Class{child.Name.String(): child}}
	if err := ds.setExampleData(); err != nil {
		t.Fatal(err)
	}
	child = ds.Classes[child.Name.String()]
	parent := Class{Name: testClassName("exampleParent")}
	instance := TestChildInstance{Properties: map[string]TestValueScenario{
		"mode": newTestValueScenario([]TestValueEntry{{ConfigValue: "test_instance", ConfigInclude: true}}),
	}}

	object := child.buildChildExampleObject(ds, parent.Name, instance, FullExample)
	assertExampleAttributes(t, object.Attributes, map[string]string{"mode": `"test_instance"`})
}

func TestChildExampleDefinitionOverridesInstancePropertyValue(t *testing.T) {
	t.Parallel()

	mode := &Property{
		PropertyName:  "mode",
		AttributeName: "mode",
		Optional:      true,
		ValueType:     String,
		propertyDefinition: PropertyDefinition{ExampleConfig: PropertyExampleConfigDefinition{
			Full: ExampleValueScenarioDefinition{
				Defined: true,
				Entries: []ExampleValueEntryDefinition{{Value: "authored_example"}},
			},
		}},
	}
	child := Class{
		Name:          testClassName("exampleChild"),
		Properties:    map[string]*Property{"mode": mode},
		PropertiesAll: []string{"mode"},
	}
	instance := TestChildInstance{Properties: map[string]TestValueScenario{
		"mode": newTestValueScenario([]TestValueEntry{{ConfigValue: "test_instance", ConfigInclude: true}}),
	}}

	object := child.buildChildExampleObject(&DataStore{}, testClassName("exampleParent"), instance, FullExample)
	assertExampleAttributes(t, object.Attributes, map[string]string{"mode": `"authored_example"`})
}

func TestChildMinimumExampleIncludesExplicitOptionalProperty(t *testing.T) {
	t.Parallel()

	mode := &Property{
		PropertyName:  "mode",
		AttributeName: "mode",
		Optional:      true,
		ValueType:     String,
		propertyDefinition: PropertyDefinition{ExampleConfig: PropertyExampleConfigDefinition{
			Minimum: ExampleValueScenarioDefinition{
				Defined: true,
				Entries: []ExampleValueEntryDefinition{{Value: "minimum_example"}},
			},
		}},
	}
	child := Class{
		Name:          testClassName("exampleChild"),
		Properties:    map[string]*Property{"mode": mode},
		PropertiesAll: []string{"mode"},
	}

	object := child.buildChildExampleObject(&DataStore{}, testClassName("exampleParent"), TestChildInstance{}, MinimumExample)
	assertExampleAttributes(t, object.Attributes, map[string]string{"mode": `"minimum_example"`})
}

func TestDependencyReferenceUsesDefinitionOnlyResourceName(t *testing.T) {
	t.Parallel()

	fvBD := testClassName("fvBD")
	ds := &DataStore{ClassDefinitions: map[string]ClassDefinition{
		fvBD.String(): {ResourceName: "bridge_domain"},
	}}
	dependency := &TestDependency{
		Class:         fvBD,
		Reference:     "aci_bridge_domain.test.id",
		ReferenceType: ResourceReference,
	}

	if actual := dependency.exampleValue(ds).HCL(); actual != "aci_bridge_domain.example.id" {
		t.Fatalf("definition-only dependency example reference %q", actual)
	}
}

func TestExampleContextUsesParentDependencyReference(t *testing.T) {
	t.Parallel()

	l3extOut := testClassName("l3extOut")
	class := Class{
		Name: testClassName("exampleChild"),
		Properties: map[string]*Property{
			"parentDn": {PropertyName: "parentDn", AttributeName: "parent_dn", Required: true, ValueType: String},
		},
		TestDependencies: []*TestDependency{{
			Class:         l3extOut,
			Reference:     "aci_l3_outside.test.id",
			ReferenceType: ResourceReference,
			Role:          Parent,
		}},
	}
	ds := &DataStore{
		ClassDefinitions: map[string]ClassDefinition{
			"l3extOut": {ResourceName: "l3_outside"},
		},
		ctx: NewContext(),
	}

	if err := class.setExampleContext(ds); err != nil {
		t.Fatal(err)
	}
	if class.ExampleContext.ParentDnMode != ParentDnFromReference {
		t.Fatalf("parent DN mode %v", class.ExampleContext.ParentDnMode)
	}
	if actual := class.ExampleContext.ParentValue.HCL(); actual != "aci_l3_outside.example.id" {
		t.Fatalf("parent reference %q", actual)
	}
}

type exampleParentSelectionInput struct {
	Parents  []string
	Override string
}

type exampleParentSelectionExpected struct {
	ParentClass string
	Reference   string
	Error       bool
}

func TestExampleContextSelectsOneClassParent(t *testing.T) {
	t.Parallel()

	testCases := []test.TestCase{
		{
			Name:  "test_first_normalized_parent_with_resource",
			Input: exampleParentSelectionInput{Parents: []string{"fvAp", "fvTenant"}},
			Expected: exampleParentSelectionExpected{
				ParentClass: "fvTenant",
				Reference:   "aci_tenant.example.id",
			},
		},
		{
			Name: "test_configured_parent_replaces_normalized_order",
			Input: exampleParentSelectionInput{
				Parents:  []string{"fvTenant", "fvAEPg"},
				Override: "fvAEPg",
			},
			Expected: exampleParentSelectionExpected{
				ParentClass: "fvAEPg",
				Reference:   "aci_application_epg.example.id",
			},
		},
		{
			Name:     "test_invalid_configured_parent",
			Input:    exampleParentSelectionInput{Override: "invalidclass"},
			Expected: exampleParentSelectionExpected{Error: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(exampleParentSelectionInput)
			expected := testCase.Expected.(exampleParentSelectionExpected)

			class := Class{
				Name: testClassName("exampleClass"),
				ClassDefinition: ClassDefinition{
					ExampleParentClass: input.Override,
				},
				Properties: map[string]*Property{
					"parentDn": {PropertyName: "parentDn", AttributeName: "parent_dn"},
				},
			}
			for _, parent := range input.Parents {
				class.Parents = append(class.Parents, testClassName(parent))
			}
			ds := &DataStore{Classes: map[string]Class{
				"fvAp":     {Name: testClassName("fvAp")},
				"fvTenant": {Name: testClassName("fvTenant"), ResourceName: "tenant", Artifacts: []ArtifactEnum{ResourceArtifact}},
				"fvAEPg":   {Name: testClassName("fvAEPg"), ResourceName: "application_epg", Artifacts: []ArtifactEnum{ResourceArtifact}},
			}}

			err := class.setExampleContext(ds)
			if expected.Error {
				if err == nil {
					t.Fatal("expected parent selection error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if class.ExampleContext.ParentClass == nil || class.ExampleContext.ParentClass.String() != expected.ParentClass {
				t.Fatalf("selected parent %v, expected %s", class.ExampleContext.ParentClass, expected.ParentClass)
			}
			if actual := class.ExampleContext.ParentValue.HCL(); actual != expected.Reference {
				t.Fatalf("parent reference %q, expected %q", actual, expected.Reference)
			}
		})
	}
}

func TestPolymorphicRelationExampleSelectsTargetForParent(t *testing.T) {
	t.Parallel()

	fvAEPg := testClassName("fvAEPg")
	fvESg := testClassName("fvESg")
	vmmDomP := testClassName("vmmDomP")
	physDomP := testClassName("physDomP")
	parentDn := &Property{
		PropertyName:  "parentDn",
		AttributeName: "parent_dn",
		Required:      true,
		ValueType:     String,
	}
	targetDn := &Property{
		PropertyName:  "tDn",
		AttributeName: "target_dn",
		Required:      true,
		ValueType:     String,
	}
	relationClass := Class{
		Name:         testClassName("fvRsDomAtt"),
		ResourceName: "relation_to_domain",
		Artifacts:    []ArtifactEnum{ResourceArtifact},
		Parents: []*ClassName{
			fvAEPg,
			fvESg,
		},
		Properties:    map[string]*Property{"parentDn": parentDn, "tDn": targetDn},
		PropertiesAll: []string{"parentDn", "tDn"},
		Relation: Relation{
			RelationalClass: true,
			Type:            Explicit,
			ToClasses:       []*ClassName{vmmDomP, physDomP},
		},
		TestDependencies: []*TestDependency{
			{Class: fvAEPg, Role: Parent, TargetClasses: []*ClassName{vmmDomP}},
			{Class: fvESg, Role: Parent, TargetClasses: []*ClassName{physDomP}},
			{Class: vmmDomP, Role: Target, Reference: "aci_vmm_domain.test.id", ReferenceType: ResourceReference},
			{Class: physDomP, Role: Target, Reference: "uni/phys-test", ReferenceType: StaticReference},
		},
	}
	ds := &DataStore{Classes: map[string]Class{
		relationClass.Name.String(): relationClass,
		fvAEPg.String():             {Name: fvAEPg, ResourceName: "application_epg", Artifacts: []ArtifactEnum{ResourceArtifact}},
		fvESg.String():              {Name: fvESg, ResourceName: "endpoint_security_group", Artifacts: []ArtifactEnum{ResourceArtifact}},
		vmmDomP.String():            {Name: vmmDomP, ResourceName: "vmm_domain", Artifacts: []ArtifactEnum{ResourceArtifact}},
		physDomP.String():           {Name: physDomP, ResourceName: "physical_domain", Artifacts: []ArtifactEnum{ResourceArtifact}},
	}}

	if err := ds.setExampleData(); err != nil {
		t.Fatal(err)
	}
	example := ds.Classes[relationClass.Name.String()].ResourceExamples.Minimum
	assertExampleAttributes(t, example.Body.Attributes, map[string]string{
		"parent_dn": "aci_application_epg.example.id",
		"target_dn": "aci_vmm_domain.example.id",
	})
}

func TestDependencyReferenceUsesGeneratedParentPlacement(t *testing.T) {
	t.Parallel()

	fvTenant := testClassName("fvTenant")
	fvAp := testClassName("fvAp")
	fvAEPg := testClassName("fvAEPg")
	parentDn := &Property{PropertyName: "parentDn", AttributeName: "parent_dn"}
	ds := &DataStore{Classes: map[string]Class{
		fvTenant.String(): {Name: fvTenant, ResourceName: "tenant", Artifacts: []ArtifactEnum{ResourceArtifact}},
		fvAp.String(): {
			Name:          fvAp,
			ResourceName:  "application_profile",
			Artifacts:     []ArtifactEnum{ResourceArtifact},
			Properties:    map[string]*Property{"parentDn": parentDn},
			PropertiesAll: []string{"parentDn"},
			Parents:       []*ClassName{fvTenant},
		},
		fvAEPg.String(): {
			Name:          fvAEPg,
			ResourceName:  "application_epg",
			Artifacts:     []ArtifactEnum{ResourceArtifact},
			Properties:    map[string]*Property{"parentDn": parentDn},
			PropertiesAll: []string{"parentDn"},
			Parents:       []*ClassName{fvAp},
		},
	}}
	dependency := &TestDependency{
		Class:         fvAEPg,
		Reference:     "aci_application_epg.test.id",
		ReferenceType: ResourceReference,
		Dependencies: []*TestDependency{{
			Class:         fvAp,
			Reference:     "aci_application_profile.test.id",
			ReferenceType: ResourceReference,
		}},
	}

	if actual := dependency.exampleValue(ds).HCL(); actual != "aci_application_epg.example.id" {
		t.Fatalf("dependency example reference %q", actual)
	}
}

func TestDependencyReferenceUsesSingleGeneratedLabel(t *testing.T) {
	t.Parallel()

	pkiEp := testClassName("pkiEp")
	fvTenant := testClassName("fvTenant")
	pkiTP := testClassName("pkiTP")
	ds := &DataStore{Classes: map[string]Class{
		fvTenant.String(): {Name: fvTenant, ResourceName: "tenant", Artifacts: []ArtifactEnum{ResourceArtifact}},
		pkiTP.String(): {
			Name:         pkiTP,
			ResourceName: "certificate_authority",
			Artifacts:    []ArtifactEnum{ResourceArtifact},
			ParentDnVariants: []ParentDnVariant{
				{ParentClass: pkiEp},
				{ParentClass: fvTenant},
			},
		},
	}}
	dependency := &TestDependency{
		Class:         pkiTP,
		Reference:     "aci_certificate_authority.test.id",
		ReferenceType: ResourceReference,
		Dependencies:  []*TestDependency{{Class: fvTenant}},
	}

	if actual := dependency.exampleValue(ds).HCL(); actual != "aci_certificate_authority.example.id" {
		t.Fatalf("dependency example reference %q", actual)
	}
}

func TestExampleContextUsesCanonicalParentDnVariant(t *testing.T) {
	t.Parallel()

	pkiEp := testClassName("pkiEp")
	fvTenant := testClassName("fvTenant")
	class := Class{
		ParentDnVariants: []ParentDnVariant{
			{ParentClass: pkiEp, ParentDn: "uni/userext/pkiext", TestPlatform: Apic},
			{ParentClass: fvTenant, RnPrepend: "certstore", TestPlatform: Cloud},
		},
	}

	if err := class.setExampleContext(&DataStore{}); err != nil {
		t.Fatal(err)
	}
	context := class.ExampleContext
	if !context.ParentClass.Equal(pkiEp) {
		t.Fatalf("selected parent class %v", context.ParentClass)
	}
	if context.ParentDnMode != ParentDnFromSchemaDefault || context.ParentValue != nil {
		t.Fatalf("canonical parent context %#v", context)
	}
	if context.Platform != Apic {
		t.Fatalf("selected platform %q", context.Platform)
	}
}

func TestDefaultParentDnExamplesUseImplicitPlacement(t *testing.T) {
	t.Parallel()

	fvTenant := testClassName("fvTenant")
	parentDn := &Property{
		PropertyName:  "parentDn",
		AttributeName: "parent_dn",
		Optional:      true,
		Computed:      true,
		ValueType:     String,
		Documentation: PropertyDocumentation{DefaultValues: []DefaultValue{{Value: "uni/infra"}}},
		TestValues: &TestValues{
			Default: newTestValueScenario([]TestValueEntry{{ConfigValue: "uni/infra", ConfigInclude: true}}),
			Create:  newTestValueScenario([]TestValueEntry{{ConfigValue: "uni/infra", ConfigInclude: true}}),
		},
	}
	name := &Property{
		PropertyName:  "name",
		AttributeName: "name",
		Required:      true,
		ValueType:     String,
		TestValues: &TestValues{
			Default: newTestValueScenario([]TestValueEntry{{ConfigValue: "example_name", ConfigInclude: true}}),
			Create:  newTestValueScenario([]TestValueEntry{{ConfigValue: "example_name", ConfigInclude: true}}),
		},
	}
	class := Class{
		Name:          testClassName("exampleClass"),
		ResourceName:  "example_class",
		Artifacts:     []ArtifactEnum{ResourceArtifact, DatasourceArtifact},
		IdentifiedBy:  []string{"name"},
		Parents:       []*ClassName{fvTenant},
		Properties:    map[string]*Property{"name": name, "parentDn": parentDn},
		PropertiesAll: []string{"name", "parentDn"},
		TestDependencies: []*TestDependency{{
			Class:         fvTenant,
			Role:          Parent,
			Reference:     "aci_tenant.test.id",
			ReferenceType: ResourceReference,
		}},
	}
	ds := &DataStore{Classes: map[string]Class{
		class.Name.String(): class,
		fvTenant.String():   {Name: fvTenant, ResourceName: "tenant", Artifacts: []ArtifactEnum{ResourceArtifact}},
	}}

	if err := ds.setExampleData(); err != nil {
		t.Fatal(err)
	}
	resolved := ds.Classes[class.Name.String()]
	if resolved.ResourceExamples.HasDistinctFullExample() {
		t.Fatal("all-attributes example should be omitted when it matches the minimum example")
	}
	assertExampleAttributes(t, resolved.ResourceExamples.Minimum.Body.Attributes, map[string]string{"name": `"example_name"`})
	assertExampleAttributes(t, resolved.ResourceExamples.Full.Body.Attributes, map[string]string{"name": `"example_name"`})
	assertExampleAttributes(t, resolved.DataSourceExample.Body.Attributes, map[string]string{"name": `"example_name"`})
	reference := newExampleReference(ResourceReference, class.ResourceName, "id")
	if actual := reference.HCL(); actual != "aci_example_class.example.id" {
		t.Fatalf("generated example reference %q", actual)
	}
}

func TestNestedExplicitRelationSelectsTargetFromRelationModel(t *testing.T) {
	t.Parallel()

	poeIfPol := testClassName("poeIfPol")
	fvTenant := testClassName("fvTenant")
	fvAp := testClassName("fvAp")
	fvAEPg := testClassName("fvAEPg")
	targetDn := &Property{PropertyName: "tDn", AttributeName: "target_dn", Required: true, ValueType: String}
	childClass := Class{
		Name:               testClassName("poeRsPoeEpg"),
		ResourceNameNested: "relation_to_application_epg",
		Properties:         map[string]*Property{"tDn": targetDn},
		PropertiesAll:      []string{"tDn"},
		Relation: Relation{
			RelationalClass: true,
			Type:            Explicit,
			ToClasses:       []*ClassName{fvAEPg},
		},
		TestDependencies: []*TestDependency{
			{Class: fvTenant, Role: Target, Reference: "aci_tenant.test.id", ReferenceType: ResourceReference},
			{Class: fvAp, Role: Target, Reference: "aci_application_profile.test.id", ReferenceType: ResourceReference},
			{
				Class:         fvAEPg,
				Role:          Target,
				Reference:     "aci_application_epg.test.id",
				ReferenceType: ResourceReference,
				Dependencies:  []*TestDependency{{Class: fvAp}},
			},
		},
	}
	parentDn := &Property{PropertyName: "parentDn", AttributeName: "parent_dn"}
	ds := &DataStore{Classes: map[string]Class{
		fvAp.String(): {Name: fvAp, ResourceName: "application_profile", Artifacts: []ArtifactEnum{ResourceArtifact}},
		fvAEPg.String(): {
			Name:          fvAEPg,
			ResourceName:  "application_epg",
			Artifacts:     []ArtifactEnum{ResourceArtifact},
			Properties:    map[string]*Property{"parentDn": parentDn},
			PropertiesAll: []string{"parentDn"},
			Parents:       []*ClassName{fvAp},
		},
	}}
	parent := Class{Name: poeIfPol}
	instance := TestChildInstance{Properties: map[string]TestValueScenario{
		"target_dn": newTestValueScenario([]TestValueEntry{{ConfigValue: "aci_tenant.test.id", ConfigInclude: true, ValueType: ReferenceValue}}),
	}}

	object := childClass.buildChildExampleObject(ds, parent.Name, instance, FullExample)
	assertExampleAttributes(t, object.Attributes, map[string]string{
		"target_dn": "aci_application_epg.example.id",
	})
}

func TestNormalizeGeneratedExampleLiteral(t *testing.T) {
	t.Parallel()

	globalDefinition := GlobalMetaDefinition{PropertyExampleValueOverrides: map[string]string{
		"accessControlAllowOrigins": "https://example.com",
		"referer":                   "https://example.com",
	}}
	class := Class{
		Name:         testClassName("tagAnnotation"),
		ResourceName: "annotation",
		Documentation: ClassDocumentation{
			Label: "Annotation",
		},
	}
	testCases := []test.TestCase{
		{Name: "test_annotation", Input: normalizeGeneratedExampleLiteralInput{Property: &Property{PropertyName: "annotation", AttributeName: "annotation"}, Value: "annotation"}, Expected: "example_annotation"},
		{Name: "test_name", Input: normalizeGeneratedExampleLiteralInput{Property: &Property{PropertyName: "name", AttributeName: "name"}, Value: "test_name"}, Expected: "annotation"},
		{Name: "test_username", Input: normalizeGeneratedExampleLiteralInput{Property: &Property{PropertyName: "username", AttributeName: "username"}, Value: "username_1"}, Expected: "example_username"},
		{Name: "test_password", Input: normalizeGeneratedExampleLiteralInput{Property: &Property{PropertyName: "password", AttributeName: "password"}, Value: "password_1"}, Expected: "example_password"},
		{Name: "test_key", Input: normalizeGeneratedExampleLiteralInput{Property: &Property{PropertyName: "key", AttributeName: "key"}, Value: "key_1"}, Expected: "example_key"},
		{Name: "test_value", Input: normalizeGeneratedExampleLiteralInput{Property: &Property{PropertyName: "value", AttributeName: "value"}, Value: "test_value"}, Expected: "example_value"},
		{Name: "test_name_alias", Input: normalizeGeneratedExampleLiteralInput{Property: &Property{PropertyName: "nameAlias", AttributeName: "name_alias"}, Value: "name_alias_1"}, Expected: "example_name_alias"},
		{Name: "test_owner_key", Input: normalizeGeneratedExampleLiteralInput{Property: &Property{PropertyName: "ownerKey", AttributeName: "owner_key"}, Value: "owner_key_1"}, Expected: "example_owner_key"},
		{Name: "test_owner_tag", Input: normalizeGeneratedExampleLiteralInput{Property: &Property{PropertyName: "ownerTag", AttributeName: "owner_tag"}, Value: "owner_tag_1"}, Expected: "example_owner_tag"},
		{Name: "test_description", Input: normalizeGeneratedExampleLiteralInput{Property: &Property{PropertyName: "descr", AttributeName: "description"}, Value: "description_1"}, Expected: "Example Annotation"},
		{Name: "test_global_override", Input: normalizeGeneratedExampleLiteralInput{Property: &Property{PropertyName: "accessControlAllowOrigins", AttributeName: "allow_origins"}, Value: "allow_origins_1"}, Expected: "https://example.com"},
		{Name: "test_second_global_override", Input: normalizeGeneratedExampleLiteralInput{Property: &Property{PropertyName: "referer", AttributeName: "referer"}, Value: "referer_1"}, Expected: "https://example.com"},
		{Name: "test_explicit_literal", Input: normalizeGeneratedExampleLiteralInput{Property: &Property{PropertyName: "accessControlAllowOrigins", AttributeName: "allow_origins"}, Value: "https://authored.example"}, Expected: "https://authored.example"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(normalizeGeneratedExampleLiteralInput)
			input.Property.globalDefinition = globalDefinition
			actual := class.normalizeGeneratedExampleLiteral(input.Property, input.Value)
			if actual != testCase.Expected.(string) {
				t.Fatalf("normalized generated example value %q, expected %q", actual, testCase.Expected)
			}
		})
	}
}

type normalizeGeneratedExampleLiteralInput struct {
	Property *Property
	Value    string
}

func assertExampleAttributes(t *testing.T, attributes []ExampleAttribute, expected map[string]string) {
	t.Helper()
	actual := make(map[string]string, len(attributes))
	for _, attribute := range attributes {
		actual[attribute.Name] = attribute.Value.HCL()
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("example attributes %#v, expected %#v", actual, expected)
	}
}
