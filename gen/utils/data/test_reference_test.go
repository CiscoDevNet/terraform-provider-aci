package data

import (
	"fmt"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestCollectChildDrivenDependenciesMatchesTerraformObject(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_attribute_projection", Input: "aci_contract.test.name", Expected: "aci_contract.test_1.name"},
		{Name: "test_placeholder", Input: "{{aci_contract.test.id}}", Expected: "aci_contract.test_1.id"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			dependency := &TestDependency{
				Class:         testClassName("vzBrCP"),
				Reference:     "aci_contract.test.id",
				ReferenceType: ResourceReference,
				Role:          Target,
			}
			child := dataSourceTestClass("exampleChildRelation", "example_child_relation")
			child.TestDependencies = []*TestDependency{dependency}
			child.canonicalizeTestDependencies()
			parent := dataSourceTestClass("exampleParent", "example_parent")
			parent.TestChildren = []*TestChild{{
				Class: child.Name,
				Instances: []TestChildInstance{{Properties: map[string]TestValueScenario{
					"contract_name": dataSourceTestScenario(testCase.Input.(string), testCase.Input.(string), ReferenceValue, true),
				}}},
			}}
			store := &DataStore{Classes: map[string]Class{child.Name.String(): child}}

			parent.collectChildDrivenDependencies(store)

			if !assert.Len(t, parent.TestDependencies, 1) {
				return
			}
			imported := parent.TestDependencies[0]
			assert.NotSame(t, dependency, imported)
			assert.Equal(t, "aci_contract.test_1.id", dependency.Reference)
			assert.Equal(t, "aci_contract.test.id", dependency.SourceReference)
			assert.Equal(t, "aci_contract.test_1.id", imported.Reference)
			assert.Equal(t, "aci_contract.test.id", imported.SourceReference)
			assert.Equal(t, 1, imported.TestInstance)
			assert.Equal(t, testCase.Expected.(string), parent.TestChildren[0].Instances[0].Properties["contract_name"].Entries[0].ConfigValue)
		})
	}
}

func TestImportTestDependencyPreservesIgnoredChanges(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	root := dataSourceTestClass("exampleRoot", "example_root")
	source := &TestDependency{
		Class: testClassName("l3extOut"), Reference: "aci_l3_outside.test.id", ReferenceType: ResourceReference,
		IgnoreChanges: []string{"relation_l3extrs_redistribute_pol"},
	}
	imported, _ := root.importTestDependency(source)
	assert.Equal(t, source.IgnoreChanges, imported.IgnoreChanges)
	imported.IgnoreChanges[0] = "other_attribute"
	assert.Equal(t, []string{"relation_l3extrs_redistribute_pol"}, source.IgnoreChanges, "imported fixtures must own their ignore list")
}

func TestCanonicalSameResourceDependencyNames(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_parent", Input: Parent, Expected: "parent"},
		{Name: "test_target", Input: Target, Expected: "target"},
		{Name: "test_prerequisite", Input: UndefinedRole, Expected: "dependency"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			role := testCase.Input.(TestDependencyRoleEnum)
			prefix := testCase.Expected.(string)
			class := dataSourceTestClass("fvAEPg", "application_epg", ResourceArtifact)
			class.Properties = map[string]*Property{
				"name": dataSourceTestProperty("name", "name", dataSourceTestValues("name_1", "name_2", "name_1", true)),
			}
			class.TestDependencies = []*TestDependency{{
				Class:         class.Name,
				Reference:     "aci_application_epg.legacy.id",
				ReferenceType: ResourceReference,
				Role:          role,
			}}

			class.canonicalizeTestDependencies()

			if !assert.Len(t, class.TestDependencies, 1) {
				return
			}
			dependency := class.TestDependencies[0]
			assert.Equal(t, fmt.Sprintf("aci_application_epg.%s_1.id", prefix), dependency.Reference)
			assert.Equal(t, fmt.Sprintf("%s_name_1", prefix), dependency.ConfigOverrides["name"])
		})
	}
}

func TestCanonicalSameResourceDependencyNameOverrides(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_terraform_attribute_override", Input: map[string]string{"display_name": "authored_name"}, Expected: "authored_name"},
		{Name: "test_apic_property_override", Input: map[string]string{"name": "authored_name"}, Expected: "authored_name"},
		{Name: "test_no_override_uses_generated_name", Input: map[string]string{}, Expected: "target_name_1"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			class := dataSourceTestClass("fvAEPg", "application_epg", ResourceArtifact)
			property := dataSourceTestProperty("name", "display_name", dataSourceTestValues("name_1", "name_2", "name_1", true))
			class.Properties = map[string]*Property{"name": property}
			class.TestDependencies = []*TestDependency{{
				Class: class.Name, Reference: "aci_application_epg.legacy.id", ReferenceType: ResourceReference,
				Role: Target, ConfigOverrides: testCase.Input.(map[string]string),
			}}

			class.canonicalizeTestDependencies()

			actual, ok := testPropertyOverride(property, class.TestDependencies[0].ConfigOverrides)
			assert.True(t, ok)
			assert.Equal(t, testCase.Expected, actual)
		})
	}
}

func TestImportedSameResourceDependencyNames(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_generated_name_follows_new_instance", Input: true, Expected: "target_name_2"},
		{Name: "test_explicit_name_is_preserved", Input: false, Expected: "authored_name"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			class := dataSourceTestClass("fvAEPg", "application_epg", ResourceArtifact)
			class.Properties = map[string]*Property{"name": dataSourceTestProperty("name", "name", nil)}
			class.TestDependencies = []*TestDependency{{Class: class.Name, Reference: "aci_application_epg.first.id", ReferenceType: ResourceReference, Role: Target}}
			class.canonicalizeTestDependencies()
			name := "authored_name"
			if testCase.Input.(bool) {
				name = "target_name_1"
			}
			source := &TestDependency{
				Class: class.Name, SourceReference: "aci_application_epg.second.id", Reference: "aci_application_epg.target_1.id",
				ReferenceType: ResourceReference, Role: Target, TestInstance: 1,
				ConfigOverrides: map[string]string{"name": name}, nameGenerated: testCase.Input.(bool),
			}

			imported, _ := class.importTestDependency(source)

			assert.Equal(t, "aci_application_epg.target_2.id", imported.Reference)
			assert.Equal(t, testCase.Expected, imported.ConfigOverrides["name"])
			assert.Equal(t, name, source.ConfigOverrides["name"])
		})
	}
}

func TestCollectChildDrivenDependenciesPreservesAssertions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_different_attribute", Input: "aci_contract.test.name", Expected: "aci_contract.test_1.name"},
		{Name: "test_placeholder_attribute", Input: "{{ aci_contract.test.name }}", Expected: "aci_contract.test_1.name"},
		{Name: "test_literal_expectation", Input: "normalized_value", Expected: "normalized_value"},
		{Name: "test_empty_expectation", Input: "", Expected: ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			child := dataSourceTestClass("exampleChildRelation", "example_child_relation")
			child.TestDependencies = []*TestDependency{{Class: testClassName("vzBrCP"), Reference: "aci_contract.test.id", ReferenceType: ResourceReference, Role: Target}}
			child.canonicalizeTestDependencies()
			parent := dataSourceTestClass("exampleParent", "example_parent")
			parent.TestChildren = []*TestChild{{Class: child.Name, Instances: []TestChildInstance{{Properties: map[string]TestValueScenario{
				"contract": dataSourceTestScenario("{{ aci_contract.test.id }}", testCase.Input.(string), ReferenceValue, true),
			}}}}}

			parent.collectChildDrivenDependencies(&DataStore{Classes: map[string]Class{child.Name.String(): child}})

			entry := parent.TestChildren[0].Instances[0].Properties["contract"].Entries[0]
			assert.Equal(t, "aci_contract.test_1.id", entry.ConfigValue)
			assert.Equal(t, testCase.Expected, entry.AssertValue)
		})
	}
}

func TestImportTestDependencyPreservesExistingPrerequisites(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	annotation := func(reference string) map[string]*TestChild {
		return map[string]*TestChild{"tagAnnotation": {
			Class: testClassName("tagAnnotation"),
			Instances: []TestChildInstance{{Properties: map[string]TestValueScenario{
				"value": dataSourceTestScenario(reference, reference, ReferenceValue, true),
			}}},
		}}
	}
	existingTenant := &TestDependency{
		Class: testClassName("fvTenant"), SourceReference: "aci_tenant.first.id", Reference: "aci_tenant.test_1.id",
		ReferenceType: ResourceReference, TestInstance: 1,
	}
	existingVRF := &TestDependency{
		Class: testClassName("fvCtx"), SourceReference: "aci_vrf.shared.id", Reference: "aci_vrf.test_1.id",
		ReferenceType: ResourceReference, TestInstance: 1, OrderOnly: true,
		ConfigOverrides: map[string]string{"parent_dn": existingTenant.Reference}, Dependencies: []*TestDependency{existingTenant},
		Children: annotation("aci_tenant.test_1.name"),
	}
	root := dataSourceTestClass("exampleRoot", "example_root")
	root.TestDependencies = []*TestDependency{existingVRF}
	sourceTenant := &TestDependency{
		Class: existingTenant.Class, SourceReference: "aci_tenant.second.id", Reference: "aci_tenant.test_1.id",
		ReferenceType: ResourceReference, TestInstance: 1,
	}
	sourceVRF := &TestDependency{
		Class: existingVRF.Class, SourceReference: existingVRF.SourceReference, Reference: existingVRF.Reference,
		ReferenceType: ResourceReference, TestInstance: 1, OrderOnly: true,
	}
	source := &TestDependency{
		Class: testClassName("fvAp"), SourceReference: "aci_application_profile.fresh.id", Reference: "aci_application_profile.test_1.id",
		ReferenceType: ResourceReference, Role: Target, TestInstance: 1,
		ConfigOverrides: map[string]string{"parent_dn": sourceTenant.Reference}, Dependencies: []*TestDependency{sourceTenant, sourceVRF},
		Children: annotation("aci_tenant.test_1.name"),
	}

	imported, _ := root.importTestDependency(source)

	assert.Same(t, existingVRF, imported.Dependencies[1], "the shared prerequisite remains one fixture object")
	assert.Same(t, existingTenant, existingVRF.Dependencies[0])
	assert.Equal(t, "aci_tenant.test_2.id", imported.Dependencies[0].Reference)
	testCases := []test.TestCase{
		{Name: "test_existing_parent_reference", Input: existingVRF.ConfigOverrides["parent_dn"], Expected: "aci_tenant.test_1.id"},
		{Name: "test_existing_child_reference", Input: existingVRF.Children["tagAnnotation"].Instances[0].Properties["value"].Entries[0].ConfigValue, Expected: "aci_tenant.test_1.name"},
		{Name: "test_existing_child_assertion", Input: existingVRF.Children["tagAnnotation"].Instances[0].Properties["value"].Entries[0].AssertValue, Expected: "aci_tenant.test_1.name"},
		{Name: "test_new_parent_reference", Input: imported.ConfigOverrides["parent_dn"], Expected: "aci_tenant.test_2.id"},
		{Name: "test_new_child_reference", Input: imported.Children["tagAnnotation"].Instances[0].Properties["value"].Entries[0].ConfigValue, Expected: "aci_tenant.test_2.name"},
		{Name: "test_new_child_assertion", Input: imported.Children["tagAnnotation"].Instances[0].Properties["value"].Entries[0].AssertValue, Expected: "aci_tenant.test_2.name"},
		{Name: "test_source_configuration_unchanged", Input: source.ConfigOverrides["parent_dn"], Expected: "aci_tenant.test_1.id"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, testCase.Expected, testCase.Input)
		})
	}
}

func TestResourceTestParentReplacementTargets(t *testing.T) {
	test.InitializeTest(t)
	_, file, _, _ := runtime.Caller(0)
	t.Chdir(filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..")))
	ctx := NewContext()
	store, err := NewDataStore(ctx)
	if !assert.NoError(t, err) || !assert.NoError(t, ctx.Diagnostics.Error()) {
		return
	}
	class := store.Classes["l3extConsLbl"]
	originalChildren := class.resourceTestCopy().TestChildren
	prepared := class.prepareResourceTest(store)
	replacement := prepared.resourceTestCopy()
	replacement.Properties["parentDn"].TestValues.Create = replacement.Properties["parentDn"].TestValues.ForceNew.clone()
	builder := resourceTestBuilder{class: &prepared, dataStore: store}
	builder.compatibleParentTargets(&replacement)
	renderer := newTestConfigurationRenderer(&replacement, store)
	parent := renderer.primaryParentDependency()
	if !assert.NotNil(t, parent) {
		return
	}
	dependencies, err := class.ResourceTestDependenciesHCL(store)
	if !assert.NoError(t, err) {
		return
	}
	targets := 0
	names := make(map[string]bool)
	for _, child := range replacement.TestChildren {
		if child.Class.String() != "l3extRsLblToInstP" {
			continue
		}
		for _, instance := range child.Instances {
			entry, ok := instance.Properties["target_dn"].firstIncludedEntry()
			if !assert.True(t, ok) {
				continue
			}
			dependency := renderer.findDependency(entry.ConfigValue)
			if !assert.NotNil(t, dependency) {
				continue
			}
			targetClass := store.Classes[dependency.Class.String()]
			targetParent := renderer.parentDependency(&targetClass, dependency)
			if assert.NotNil(t, targetParent) {
				assert.Equal(t, parent.Reference, targetParent.Reference)
			}
			address, _ := parseTerraformTestReference(dependency.Reference)
			assert.Contains(t, dependencies, `resource "aci_external_epg" "`+address.Name+`"`)
			assert.NotEmpty(t, dependency.ConfigOverrides["name"])
			names[dependency.ConfigOverrides["name"]] = true
			targets++
		}
	}
	assert.Equal(t, 2, targets)
	assert.Len(t, names, 2)
	assert.Equal(t, originalChildren, class.resourceTestCopy().TestChildren)
	// Tenant-scoped route-control profiles keep their existing references when the L3Out changes.
	profiles := 0
	for index, child := range replacement.TestChildren {
		if child.Class.String() == "l3extRsLblToProfile" {
			assert.Equal(t, originalChildren[index], child)
			profiles++
		}
	}
	assert.Equal(t, 1, profiles)
}

func TestResourceTestReplacementDependencyReferences(t *testing.T) {
	test.InitializeTest(t)
	_, file, _, _ := runtime.Caller(0)
	t.Chdir(filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..")))
	ctx := NewContext()
	store, err := NewDataStore(ctx)
	if !assert.NoError(t, err) || !assert.NoError(t, ctx.Diagnostics.Error()) {
		return
	}
	class := store.Classes["fvBD"]
	prepared := class.prepareResourceTest(store)
	original := prepared.resourceTestCopy()
	replacement := prepared.resourceTestCopy()
	replacement.Properties["parentDn"].TestValues.Create = replacement.Properties["parentDn"].TestValues.ForceNew.clone()
	builder := resourceTestBuilder{class: &prepared, dataStore: store}
	builder.compatibleParentTargets(&replacement)
	renderer := newTestConfigurationRenderer(&replacement, store)
	var vrfReference string
	for _, child := range replacement.TestChildren {
		if child.Class.String() == "fvRsCtx" {
			for _, value := range child.Instances[0].Properties {
				entry, ok := value.firstIncludedEntry()
				if ok && entry.ValueType == ReferenceValue {
					vrfReference = entry.ConfigValue
				}
			}
		}
	}
	if !assert.NotEmpty(t, vrfReference) {
		return
	}
	vrf := renderer.findDependency(vrfReference)
	if !assert.NotNil(t, vrf) {
		return
	}
	count := 0
	for _, child := range replacement.TestChildren {
		if child.Class.String() != "fvRsBDToOut" {
			continue
		}
		for _, instance := range child.Instances {
			for _, value := range instance.Properties {
				entry, ok := value.firstIncludedEntry()
				if !ok || entry.ValueType != ReferenceValue {
					continue
				}
				dependency := renderer.findDependency(entry.ConfigValue)
				if !assert.NotNil(t, dependency) {
					continue
				}
				assert.Contains(t, dependency.ConfigOverrides, "relation_l3ext_rs_ectx")
				assert.Equal(t, vrf.Reference, dependency.ConfigOverrides["relation_l3ext_rs_ectx"])
				assert.True(t, slices.Contains(dependency.Dependencies, vrf), "replacement L3Out must depend on the replacement VRF")
				count++
			}
		}
	}
	assert.Equal(t, 2, count)
	assert.Equal(t, original, prepared.resourceTestCopy(), "original fixtures must remain unchanged")
}
