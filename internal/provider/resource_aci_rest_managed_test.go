package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
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
