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

func TestAccAciRestManaged_escapeHtml(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfig_escapeHtml(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.aaaPreLoginBanner", "content.message", "<<< WARNING >>>  VERIFYING THE CONVERSION OF HTML CHARACTERS."),
					resource.TestCheckResourceAttr("aci_rest_managed.aaaPreLoginBanner", "escape_html", "false"),
					resource.TestCheckResourceAttr("aci_rest_managed.aaaPreLoginBanner", "dn", "uni/userext/preloginbanner"),
					resource.TestCheckResourceAttr("aci_rest_managed.aaaPreLoginBanner", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.aaaPreLoginBanner", "class_name", "aaaPreLoginBanner"),
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
				Config:        testAccAciRestManagedConfig_import(name, "import"),
				ImportState:   true,
				ImportStateId: "uni/tn-non-existent",
				ResourceName:  "aci_rest_managed.import",
				ExpectError:   regexp.MustCompile("Cannot import non-existent remote object"),
			},
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
					resource.TestCheckNoResourceAttr("aci_rest_managed.import", "content.status"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("fvTenant:uni/tn-%s", name),
				ResourceName:  "aci_rest_managed.import",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.import", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.import", "class_name", "fvTenant"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("polUni:uni/tn-%s", name),
				ResourceName:  "aci_rest_managed.import",
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
					resource.TestCheckNoResourceAttr("aci_rest_managed.import", "child.0.content.status"),
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
				ImportStateId: fmt.Sprintf("polUni:uni/tn-%s:ctx-VRF1,ctx-VRF2", name),
				ResourceName:  "aci_rest_managed.import",
				ExpectError:   regexp.MustCompile("Import Failed"),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("polUni:uni/tn-%s:ctx-VRF1:ctx-VRF2", name),
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
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("polUni:uni/tn-%s:ctx-VRF1,ctx-VRF2,ap-AP1", name),
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

func TestAccAciRestManaged_importWithIpv6(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfig_importWithIpv6(name, "import", "2001:1:2::5/28", "2001:1:2::5/28"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.tenant_import", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.tenant_import", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.tenant_import", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.tenant_import", "class_name", "fvTenant"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import", "dn", "uni/tn-"+name+"/BD-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import", "class_name", "fvBD"),
					resource.TestCheckResourceAttr("aci_rest_managed.subnet_import", "content.ip", "2001:1:2::5/28"),
					resource.TestCheckResourceAttr("aci_rest_managed.subnet_import", "dn", "uni/tn-"+name+"/BD-"+name+"/subnet-[2001:1:2::5/28]"),
					resource.TestCheckResourceAttr("aci_rest_managed.subnet_import", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.subnet_import", "class_name", "fvSubnet"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "content.name", name+"_2"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "dn", "uni/tn-"+name+"/BD-"+name+"_2"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "class_name", "fvBD"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("fvSubnet:uni/tn-%s/BD-%s/subnet-[2001:1:2::5/28]", name, name),
				ResourceName:  "aci_rest_managed.bd_import",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.subnet_import", "content.ip", "2001:1:2::5/28"),
					resource.TestCheckResourceAttr("aci_rest_managed.subnet_import", "dn", "uni/tn-"+name+"/BD-"+name+"/subnet-[2001:1:2::5/28]"),
					resource.TestCheckResourceAttr("aci_rest_managed.subnet_import", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.subnet_import", "class_name", "fvSubnet"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("uni/tn-%s/BD-%s/subnet-[2001:1:2::5/28]", name, name),
				ResourceName:  "aci_rest_managed.bd_import",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.subnet_import", "content.ip", "2001:1:2::5/28"),
					resource.TestCheckResourceAttr("aci_rest_managed.subnet_import", "dn", "uni/tn-"+name+"/BD-"+name+"/subnet-[2001:1:2::5/28]"),
					resource.TestCheckResourceAttr("aci_rest_managed.subnet_import", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.subnet_import", "class_name", "fvSubnet"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("fvBD:uni/tn-%s/BD-%s:rsctx,subnet-[2001:1:2::5/28]", name, name),
				ResourceName:  "aci_rest_managed.bd_import_2",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "content.name", name+"_2"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "dn", "uni/tn-"+name+"/BD-"+name+"_2"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "class_name", "fvBD"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "child.0.class_name", "fvRsCtx"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "child.0.rn", "rsctx"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "child.0.content.tnFvCtxName", "VRF2"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "child.1.class_name", "fvSubnet"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "child.1.rn", "subnet-[2001:1:2::5/28]"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "child.1.content.ip", "2001:1:2::5/28"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("uni/tn-%s/BD-%s:subnet-[2001:1:2::5/28],rsctx", name, name),
				ResourceName:  "aci_rest_managed.bd_import_2",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "content.name", name+"_2"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "dn", "uni/tn-"+name+"/BD-"+name+"_2"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "class_name", "fvBD"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "child.0.class_name", "fvSubnet"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "child.0.rn", "subnet-[2001:1:2::5/28]"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "child.0.content.ip", "2001:1:2::5/28"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "child.1.class_name", "fvRsCtx"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "child.1.rn", "rsctx"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd_import_2", "child.1.content.tnFvCtxName", "VRF2"),
				),
			},
		},
	})
}

func TestAccAciRestManaged_importWithBracket(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfig_importWithBracket(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "content.allocMode", "static"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "dn", "uni/infra/vlanns-["+name+"]-static"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "class_name", "fvnsVlanInstP"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("fvnsVlanInstP:uni/infra/vlanns-[%s]-static", name),
				ResourceName:  "aci_rest_managed.bracket",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "content.allocMode", "static"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "dn", "uni/infra/vlanns-["+name+"]-static"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "class_name", "fvnsVlanInstP"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("uni/infra/vlanns-[%s]-static", name),
				ResourceName:  "aci_rest_managed.bracket",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "content.allocMode", "static"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "dn", "uni/infra/vlanns-["+name+"]-static"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "class_name", "fvnsVlanInstP"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("fvnsVlanInstP:uni/infra/vlanns-[%s]-static:from-[vlan-200]-to-[vlan-2200]", name),
				ResourceName:  "aci_rest_managed.bracket",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "content.allocMode", "static"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "dn", "uni/infra/vlanns-["+name+"]-static"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "class_name", "fvnsVlanInstP"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "child.0.class_name", "fvnsEncapBlk"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "child.0.rn", "from-[vlan-200]-to-[vlan-2200]"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "child.0.content.from", "vlan-200"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "child.0.content.to", "vlan-2200"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: fmt.Sprintf("uni/infra/vlanns-[%s]-static:from-[vlan-200]-to-[vlan-2200]", name),
				ResourceName:  "aci_rest_managed.bracket",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "content.allocMode", "static"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "dn", "uni/infra/vlanns-["+name+"]-static"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "class_name", "fvnsVlanInstP"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "child.0.class_name", "fvnsEncapBlk"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "child.0.rn", "from-[vlan-200]-to-[vlan-2200]"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "child.0.content.from", "vlan-200"),
					resource.TestCheckResourceAttr("aci_rest_managed.bracket", "child.0.content.to", "vlan-2200"),
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
	setGlobalAnnotationEnvVariable(t, "")
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
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "annotation", ""),
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

func testAccAciRestManagedConfig_importWithIpv6(name string, resource string, ip string, ip2 string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "tenant_%[2]s" {
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
	}

	resource "aci_rest_managed" "bd_%[2]s" {
		dn = "${aci_rest_managed.tenant_%[2]s.id}/BD-%[1]s"
		class_name = "fvBD"
		content = {
			name  = "%[1]s"
		}
		child {
			rn         = "rsctx"
			class_name = "fvRsCtx"
			content = {
				tnFvCtxName = "VRF1"
			}
		}
	}

	resource "aci_rest_managed" "subnet_%[2]s" {
		dn         = "${aci_rest_managed.bd_%[2]s.id}/subnet-[%[3]s]"
		class_name = "fvSubnet"
		content = {
		  ip           = "%[3]s"
		  scope        = "private"
		  ipDPLearning = "enabled"
		  ctrl         = "nd"
		}
	}

	resource "aci_rest_managed" "bd_%[2]s_2" {
		dn = "${aci_rest_managed.tenant_%[2]s.id}/BD-%[1]s_2"
		class_name = "fvBD"
		content = {
			name  = "%[1]s_2"
		}
		child {
			rn         = "rsctx"
			class_name = "fvRsCtx"
			content = {
				tnFvCtxName = "VRF2"
			}
		}
		child {
			rn         = "subnet-[%[4]s]"
			class_name = "fvSubnet"
			content = {
				ip = "%[4]s"
				scope        = "private"
				ipDPLearning = "enabled"
				ctrl         = "nd"
			}
		}
	}

	`, name, resource, ip, ip2)
}

func testAccAciRestManagedConfig_importWithBracket(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "bracket" {
		dn         = "uni/infra/vlanns-[%[1]s]-static"
		class_name = "fvnsVlanInstP"
		content = {
			name      = "%[1]s"
			allocMode = "static"
		}
		child {
			rn         = "from-[vlan-200]-to-[vlan-2200]"
			class_name = "fvnsEncapBlk"
			content = {
				from = "vlan-200"
				to   = "vlan-2200"
			}
		}
	}
	`, name)
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

func testAccAciRestManagedConfig_escapeHtml() string {
	return `
	resource "aci_rest_managed" "aaaPreLoginBanner" {
		dn          = "uni/userext/preloginbanner"
		class_name  = "aaaPreLoginBanner"
		escape_html = false
		content = {
			message = "<<< WARNING >>>  VERIFYING THE CONVERSION OF HTML CHARACTERS."
		}
	}
	`
}
