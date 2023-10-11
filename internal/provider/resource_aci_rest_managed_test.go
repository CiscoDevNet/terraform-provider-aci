package provider

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAciRestManaged_tenant(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfig_tenant(name, "Create description"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.descr", "Create description"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-" + name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "annotation", "orchestrator:terraform"),
				),
			},
			{
				Config: testAccAciRestManagedConfig_tenant(name, "Updated description"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "content.descr", "Updated description"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-" + name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "annotation", "orchestrator:terraform"),
				),
			},
		},
	})
}

func TestAccAciRestManaged_connPref(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
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
		PreCheck:          func() { testAccPreCheck(t) },
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfig_tenantVrf(name),
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

func testAccAciRestManagedConfig_tenant(name string, description string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "fvTenant" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
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

func testAccAciRestManagedConfig_tenantVrf(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "testtenVrf" {
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
	`, name)
}