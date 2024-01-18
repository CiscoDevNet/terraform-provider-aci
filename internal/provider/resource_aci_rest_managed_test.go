package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAciRestManaged_tenant(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccAciRestManagedConfig_tenant(name, "Create description"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.descr", "Create description"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "annotation", "non-default:value"),
				),
			},
			{
				Config:             testAccAciRestManagedConfig_tenant(name, "Updated description"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.descr", "Updated description"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "annotation", "non-default:value"),
				),
			},
			// Import testing
			{
				ResourceName: "aci_rest_managed.fvTenant",
				ImportState:  true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "annotation", "non-default:value"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "class_name", "fvTenant"),
				),
			},
		},
	})
}

func TestAccAciRestManaged_connPref(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfig_connPref("ooband"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.testConnPref", "content.interfacePref", "ooband"),
					resource.TestCheckResourceAttr("aci_rest_managed.testConnPref", "dn", "uni/fabric/connectivityPrefs"),
					resource.TestCheckResourceAttr("aci_rest_managed.testConnPref", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.testConnPref", "class_name", "mgmtConnectivityPrefs"),
				),
			},
			{
				Config: testAccAciRestManagedConfig_connPref("inband"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.testConnPref", "content.interfacePref", "inband"),
					resource.TestCheckResourceAttr("aci_rest_managed.testConnPref", "dn", "uni/fabric/connectivityPrefs"),
					resource.TestCheckResourceAttr("aci_rest_managed.testConnPref", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.testConnPref", "class_name", "mgmtConnectivityPrefs"),
				),
			},
		},
	})
}

func TestAccAciRestManaged_noContent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfig_connPref(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.testConnPref", "dn", "uni/fabric/connectivityPrefs"),
					resource.TestCheckResourceAttr("aci_rest_managed.testConnPref", "class_name", "mgmtConnectivityPrefs"),
				),
			},
		},
	})
}

func TestAccAciRestManaged_tenantVrf(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfig_tenantVrf(name, "testtenVrf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.testtenVrf", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.testtenVrf", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.testtenVrf", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.testtenVrf", "class_name", "fvTenant"),
					resource.TestCheckResourceAttr("aci_rest_managed.testtenVrf", "child.0.class_name", "fvCtx"),
					resource.TestCheckResourceAttr("aci_rest_managed.testtenVrf", "child.0.rn", "ctx-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.testtenVrf", "child.0.content.name", name),
				),
			},
		},
	})
}

func TestAccAciRestManaged_import(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfig_import(name, "import"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.import", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "class_name", "fvTenant"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("uni/tn-%s", name),
				ResourceName:  "aci_rest_managed.import",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.import", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "class_name", "fvTenant"),
				),
			},
			{
				Config: testAccAciRestManagedConfig_importWithChild(name, "import"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.import", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "class_name", "fvTenant"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.0.class_name", "fvCtx"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.0.rn", "ctx-VRF1"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.0.content.name", "VRF1"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.#", "1"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("uni/tn-%s:ctx-VRF1", name),
				ResourceName:  "aci_rest_managed.import",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.import", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "class_name", "fvTenant"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.0.class_name", "fvCtx"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.0.rn", "ctx-VRF1"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.0.content.name", "VRF1"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.#", "1"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("uni/tn-%s:ctx-VRF1,ctx-VRF2", name),
				ResourceName:  "aci_rest_managed.import",
				ExpectError:   regexp.MustCompile("Import Failed"),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("uni/tn-%s:ctx-VRF1:ctx-VRF2", name),
				ResourceName:  "aci_rest_managed.import",
				ExpectError:   regexp.MustCompile("Unexpected Import Identifier"),
			},
			{
				Config: testAccAciRestManagedConfig_importWithMultipleChildren(name, "import"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.import", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "class_name", "fvTenant"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.0.class_name", "fvAp"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.0.rn", "ap-AP1"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.0.content.name", "AP1"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.1.class_name", "fvCtx"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.1.rn", "ctx-VRF1"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.1.content.name", "VRF1"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.2.class_name", "fvCtx"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.2.rn", "ctx-VRF2"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.2.content.name", "VRF2"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.#", "3"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("uni/tn-%s:ctx-VRF1,ctx-VRF2,ap-AP1", name),
				ResourceName:  "aci_rest_managed.import",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.import", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "class_name", "fvTenant"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.0.class_name", "fvAp"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.0.rn", "ap-AP1"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.0.content.name", "AP1"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.1.class_name", "fvCtx"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.1.rn", "ctx-VRF1"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.1.content.name", "VRF1"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.2.class_name", "fvCtx"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.2.rn", "ctx-VRF2"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.2.content.name", "VRF2"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "child.#", "2"),
				),
			},
		},
	})
}

func TestAccAciRestManaged_tagTag(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccAciRestManagedConfig_tagTag(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.testTag", "content.value", "test"),
					resource.TestCheckResourceAttr("aci_rest_managed.testTag", "dn", "uni/fabric/connectivityPrefs/tagKey-"+name),
					resource.TestCheckNoResourceAttr("aci_rest_managed.testTag", "annotation"),
				),
			},
		},
	})
}

// step 1 - create a tenant with no children
// step 2 - update same tenant with child vrf
// step 3 - Update tenant again and remove children
func TestAccAciRestManaged_tenantChildren(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccAciRestManagedConfig_tenant(name, "Create description"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.descr", "Create description"),
				),
			},
			{
				Config:             testAccAciRestManagedConfig_tenantVrf(name, "fvTenant"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "child.0.content.name", name),
					resource.TestCheckNoResourceAttr("aci_rest_managed.fvTenant", "content.descr"),
				),
			},
			{
				Config:             testAccAciRestManagedConfig_tenant(name, "Removed children"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.descr", "Removed children"),
				),
			},
		},
	})

}

func TestAccAciRestManaged_globalAnnotation(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccAciRestManagedConfig_globalAnnotation(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
				),
			},
		},
	})

	setGlobalAnnotationEnvVariable(t, "orchestrator:from_env")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccAciRestManagedConfig_globalAnnotation(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "annotation", "orchestrator:from_env"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
				),
			},
			{
				Config:             testAccAciRestManagedConfig_globalAnnotationOverwriteResource(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "annotation", "orchestrator:from_resource"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
				),
			},
			{
				Config:      testAccAciRestManagedConfig_globalAnnotationResourceAndContent(name),
				ExpectError: regexp.MustCompile("Annotation not supported in content"),
			},
			{
				Config:      testAccAciRestManagedConfig_globalAnnotationOverwriteFromContent(name),
				ExpectError: regexp.MustCompile("Annotation not supported in content"),
			},
			{
				Config:             testAccAciRestManagedConfig_globalAnnotationOverwriteResourceOverwriteFromContentNull(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "annotation", "orchestrator:from_resource"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
				),
			},
			{
				Config:             testAccAciRestManagedConfig_globalAnnotationOverwriteFromContentNull(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "annotation", "orchestrator:from_env"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
				),
			},
			{
				Config:             testAccAciRestManagedConfig_globalAnnotation(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "annotation", "orchestrator:from_env"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
				),
			},
		},
	})

}

func TestAccAciRestManaged_undeletableClass(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccAciRestManagedConfig_undeletableClass(),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.mcpInstPol", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.mcpInstPol", "dn", "uni/infra/mcpInstP-default"),
					resource.TestCheckResourceAttr("aci_rest_managed.mcpInstPol", "content.adminSt", "disabled"),
					resource.TestCheckResourceAttr("aci_rest_managed.mcpInstPol", "content.loopProtectAct", "port-disable"),
					resource.TestCheckResourceAttr("aci_rest_managed.mcpInstPol", "content.key", "test"),
					resource.TestCheckResourceAttr("aci_rest_managed.mcpInstPol", "content.%", "3"),
				),
			},
			{
				Config:             testAccAciRestManagedConfig_ignoreKeyChange(),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.mcpInstPol", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.mcpInstPol", "dn", "uni/infra/mcpInstP-default"),
					resource.TestCheckResourceAttr("aci_rest_managed.mcpInstPol", "content.adminSt", "disabled"),
					resource.TestCheckResourceAttr("aci_rest_managed.mcpInstPol", "content.loopProtectAct", "port-disable"),
					resource.TestCheckResourceAttr("aci_rest_managed.mcpInstPol", "content.%", "3"),
				),
			},
		},
	})
}

func TestAccAciRestManaged_explicitNull(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccAciRestManagedConfig_null(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
					// content is 2 because description exists in state but is set to null which type is not testable is TestCheckResourceAttr
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.%", "2"),
				),
			},
			{
				Config:             testAccAciRestManagedConfig_notNull(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.descr", "non_null_description"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.%", "2"),
				),
			},
			{
				Config:             testAccAciRestManagedConfig_null(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
					// content is 2 because description exists in state but is set to null which type is not testable is TestCheckResourceAttr
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.%", "2"),
				),
			},
			{
				Config:             testAccAciRestManagedConfig_notDefined(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "child.#", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.%", "1"),
					resource.TestCheckNoResourceAttr("aci_rest_managed.fvTenant", "content.descr"),
				),
			},
		},
	})

}

func testAccAciRestManagedConfig_tenant(name string, description string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "fvTenant" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		annotation = "non-default:value"
		content = {
			name = "%[1]s"
			descr = "%[2]s"
			nameAlias = "Testacc_Tenant"
		}
	}
	`, name, description)
}

func testAccAciRestManagedConfig_connPref(status string) string {
	if status != "" {
		return fmt.Sprintf(`
		resource "aci_rest_managed" "testConnPref" {
			dn = "uni/fabric/connectivityPrefs"
			class_name = "mgmtConnectivityPrefs"
			content = {
				interfacePref = "%[1]s"
			}
		}
		`, status)
	} else {
		return `
		resource "aci_rest_managed" "testConnPref" {
			dn = "uni/fabric/connectivityPrefs"
			class_name = "mgmtConnectivityPrefs"
		}
		`
	}
}

func testAccAciRestManagedConfig_tenantVrf(name string, resource string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "%[2]s" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
		}

		child {
			rn         = "ctx-%[1]s"
			class_name = "fvCtx"
			content = {
			  name = "%[1]s"
			  description = null
			}
		}
	}
	`, name, resource)
}

func testAccAciRestManagedConfig_import(name string, resource string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "%[2]s" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
		}
	}
	`, name, resource)
}

func testAccAciRestManagedConfig_importWithChild(name string, resource string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "%[2]s" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
		}
		child {
			rn         = "ctx-VRF1"
			class_name = "fvCtx"
			content = {
				name = "VRF1"
			}
		}
	}
	`, name, resource)
}

func testAccAciRestManagedConfig_importWithMultipleChildren(name string, resource string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "%[2]s" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
		}
		child {
			rn         = "ctx-VRF1"
			class_name = "fvCtx"
			content = {
				name = "VRF1"
			}
		}
		child {
			rn         = "ctx-VRF2"
			class_name = "fvCtx"
			content = {
				name = "VRF2"
			}
		}
		child {
			rn         = "ap-AP1"
			class_name = "fvAp"
			content = {
				name = "AP1"
			}
		}
	}
	`, name, resource)
}

func testAccAciRestManagedConfig_tagTag(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "testTag" {
		dn = "uni/fabric/connectivityPrefs/tagKey-%[1]s"
		class_name = "tagTag"
		content = {
			value = "test"
		}
	}
	`, name)
}

func testAccAciRestManagedConfig_globalAnnotation(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "fvTenant" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
		}
	}
	`, name)
}

func testAccAciRestManagedConfig_globalAnnotationOverwriteResource(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "fvTenant" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		annotation = "orchestrator:from_resource"
		content = {
			name = "%[1]s"
		}
	}
	`, name)
}

func testAccAciRestManagedConfig_globalAnnotationResourceAndContent(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "fvTenant" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		annotation = "orchestrator:from_resource"
		content = {
			name = "%[1]s"
			annotation = "orchestrator:from_content"
		}
	}
	`, name)
}

func testAccAciRestManagedConfig_globalAnnotationOverwriteFromContent(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "fvTenant" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
			annotation = "orchestrator:from_content"
		}
	}
	`, name)
}

func testAccAciRestManagedConfig_globalAnnotationOverwriteFromContentNull(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "fvTenant" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
			annotation = null
		}
	}
	`, name)
}

func testAccAciRestManagedConfig_globalAnnotationOverwriteResourceOverwriteFromContentNull(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "fvTenant" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		annotation = "orchestrator:from_resource"
		content = {
			name = "%[1]s"
			annotation = null
		}
	}
	`, name)
}

func testAccAciRestManagedConfig_undeletableClass() string {
	return `
	resource "aci_rest_managed" "mcpInstPol" {
		dn         = "uni/infra/mcpInstP-default"
		class_name = "mcpInstPol"
		content = {
		  adminSt        = "disabled"
		  key            = "test"
		  loopProtectAct = "port-disable"
		}
	} 
	`
}

func testAccAciRestManagedConfig_ignoreKeyChange() string {
	return `
	resource "aci_rest_managed" "mcpInstPol" {
		dn         = "uni/infra/mcpInstP-default"
		class_name = "mcpInstPol"
		content = {
		  adminSt        = "disabled"
		  key            = "test"
		  loopProtectAct = "port-disable"
		}
		lifecycle {
			ignore_changes = [content["key"]]
		}
	}
	`
}

func testAccAciRestManagedConfig_null(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "fvTenant" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
			descr = null
		}
	}
	`, name)
}

func testAccAciRestManagedConfig_notNull(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "fvTenant" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
			descr = "non_null_description"
		}
	}
	`, name)
}

func testAccAciRestManagedConfig_notDefined(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "fvTenant" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
		}
	}
	`, name)
}
