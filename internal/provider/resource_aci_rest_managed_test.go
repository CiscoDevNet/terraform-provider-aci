package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAciRestManaged_tenant(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfig_escapeHtml(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.aaaPreLoginBanner", "content.message", "&&& WARNING &&& VERIFYING THE CONVERSION OF HTML CHARACTERS."),
					resource.TestCheckResourceAttr("aci_rest_managed.aaaPreLoginBanner", "escape_html", "false"),
					resource.TestCheckResourceAttr("aci_rest_managed.aaaPreLoginBanner", "dn", "uni/userext/preloginbanner"),
					resource.TestCheckResourceAttr("aci_rest_managed.aaaPreLoginBanner", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_rest_managed.aaaPreLoginBanner", "class_name", "aaaPreLoginBanner"),
				),
			},
		},
	})
}

func TestAccAciRestManaged_escapeHtmlTrue(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccAciRestManagedConfig_escapeHtmlTrue(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.aaaPreLoginBanner", "content.message", "&&& WARNING &&& VERIFYING THE CONVERSION OF HTML CHARACTERS."),
					resource.TestCheckResourceAttr("aci_rest_managed.aaaPreLoginBanner", "escape_html", "true"),
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
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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
			{
				Config: testAccAciRestManagedConfig_importMultipleChildrenWithImportJsonString(name, "import_json_str"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "dn", "uni/tn-"+name+"/eptags"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.#", "4"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "annotation", "orchestrator:from_resource"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "class_name", "fvEpTags"),
				),
			},
			{

				ImportState:   true,
				ImportStateId: fmt.Sprintf(`{ "parentDn": "uni/tn-%s/eptags" }`, name),
				ResourceName:  "aci_rest_managed.eptags_import_json_str",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "content.annotation", "orchestrator:from_resource"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "dn", "uni/tn-"+name+"/eptags"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "annotation", "orchestrator:from_resource"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "class_name", "fvEpTags"),
				),
			},
			{

				ImportState:   true,
				ImportStateId: fmt.Sprintf(`{ "parentDn": "uni/tn-%s/eptags", "childRns": [ "annotationKey-[~!$([])_+-={};:|,.]", "annotationKey-[tagAnnotation2]", "epmactag-90:B5:B8:42:D1:88-[default]", "epiptag-[2001:10:1::1]-default" ] }`, name),
				ResourceName:  "aci_rest_managed.eptags_import_json_str",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "content.annotation", "orchestrator:from_resource"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "dn", "uni/tn-"+name+"/eptags"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.#", "4"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "annotation", "orchestrator:from_resource"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "class_name", "fvEpTags"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.0.class_name", "tagAnnotation"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.0.rn", "annotationKey-[~!$([])_+-={};:|,.]"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.0.content.key", "~!$([])_+-={};:|,."),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.0.content.value", "tagAnnotation1"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.1.class_name", "tagAnnotation"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.1.rn", "annotationKey-[tagAnnotation2]"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.1.content.key", "tagAnnotation2"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.1.content.value", "tagAnnotation2"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.2.class_name", "fvEpMacTag"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.2.rn", "epmactag-90:B5:B8:42:D1:88-[default]"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.2.content.bdName", "default"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.2.content.mac", "90:B5:B8:42:D1:88"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.3.class_name", "fvEpIpTag"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.3.rn", "epiptag-[2001:10:1::1]-default"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.3.content.ip", "2001:10:1::1"),
					resource.TestCheckResourceAttr("aci_rest_managed.eptags_import_json_str", "child.3.content.bdName", "default"),
				),
			},
			{

				ImportState:   true,
				ImportStateId: fmt.Sprintf(`{ "parentDn": "uni/tn-%s/BD-bd_test", "childRns": [ "rgexpmac-00:22:BD:F8:11:FF", "rgexpmac-00:22:BD:F8:12:FF" ] }`, name),
				ResourceName:  "aci_rest_managed.bd",
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("aci_rest_managed.bd", "dn", "uni/tn-"+name+"/BD-bd_test"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd", "child.#", "2"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd", "class_name", "fvBD"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd", "name", "bd_test"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd", "child.0.class_name", "fvRogueExceptionMac"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd", "child.0.mac", "00:22:BD:F8:11:FF"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd", "child.1.class_name", "fvRogueExceptionMac"),
					resource.TestCheckResourceAttr("aci_rest_managed.bd", "child.1.mac", "00:22:BD:F8:12:FF"),
				),
			},
			{

				ImportState:   true,
				ImportStateId: fmt.Sprintf(`{ "parentDn": "uni/tn-%s/BD-bd_test/rgexpmac-00:22:BD:F8:11:FF" }`, name),
				ResourceName:  "aci_rest_managed.rogue_exception_mac_11",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.rogue_exception_mac_11", "dn", "uni/tn-"+name+"/BD-bd_test/rgexpmac-00:22:BD:F8:11:FF"),
					resource.TestCheckResourceAttr("aci_rest_managed.rogue_exception_mac_11", "class_name", "fvRogueExceptionMac"),
					resource.TestCheckResourceAttr("aci_rest_managed.rogue_exception_mac_11", "mac", "00:22:BD:F8:11:FF"),
				),
			},
		},
	})
}

func TestAccAciRestManaged_importWithIpv6(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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
			{
				Config:      testAccAciRestManagedData_vrf(name),
				ExpectError: regexp.MustCompile("Failed to read aci_rest_managed data source"),
			},
		},
	})
}

func TestAccAciRestManaged_ignoreChildAnnotation(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccAciRestManagedConfig_ignoreChildAnnotation(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.ignore_child_annotation", "child.#", "1"),
					resource.TestCheckResourceAttr("aci_rest_managed.ignore_child_annotation", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.ignore_child_annotation", "content.name", name),
					resource.TestCheckResourceAttr("aci_rest_managed.ignore_child_annotation", "child.0.class_name", "tagAnnotation"),
					resource.TestCheckResourceAttr("aci_rest_managed.ignore_child_annotation", "child.0.rn", "annotationKey-["+name+"]"),
					resource.TestCheckResourceAttr("aci_rest_managed.ignore_child_annotation", "child.0.content.key", name),
					resource.TestCheckResourceAttr("aci_rest_managed.ignore_child_annotation", "child.0.content.value", "value"),
					resource.TestCheckNoResourceAttr("aci_rest_managed.ignore_child_annotation", "child.0.content.annotation"),
				),
			},
		},
	})
}

func TestAccAciRestManaged_globalAllowExistingOnCreate(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccAciRestManagedConfig_globalAllowExisting(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant_2", "dn", "uni/tn-"+name),
				),
			},
		},
	})

	setEnvVariable(t, "ACI_ALLOW_EXISTING_ON_CREATE", "false")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccAciRestManagedConfig_globalAllowExisting(name),
				ExpectError: regexp.MustCompile("Object Already Exists"),
			},
		},
	})

	setEnvVariable(t, "ACI_ALLOW_EXISTING_ON_CREATE", "true")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccAciRestManagedConfig_globalAllowExisting(name),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant", "dn", "uni/tn-"+name),
					resource.TestCheckResourceAttr("aci_rest_managed.fvTenant_2", "dn", "uni/tn-"+name),
				),
			},
		},
	})
}

func TestAccAciRestManaged_globalAnnotation(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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

	setEnvVariable(t, "ACI_ANNOTATION", "orchestrator:from_env")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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
	setEnvVariable(t, "ACI_ANNOTATION", "")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(7g)-") },
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

func TestAccAciRestManaged_undeletableObject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "2.0(1m)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccAciRestManagedConfig_undeletableObject(),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					// Validate the attributes for fvFabricExtConnP
					resource.TestCheckResourceAttr("aci_rest_managed.fvFabricExtConnP", "dn", "uni/tn-infra/fabricExtConnP-1"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvFabricExtConnP", "class_name", "fvFabricExtConnP"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvFabricExtConnP", "content.id", "1"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvFabricExtConnP", "content.name", "IPN"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvFabricExtConnP", "content.rt", "extended:as2-nn4:5:16"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvFabricExtConnP", "content.siteId", "0"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvFabricExtConnP", "content.%", "4"),

					// Validate the attributes for fvPeeringP
					resource.TestCheckResourceAttr("aci_rest_managed.fvPeeringP", "dn", "uni/tn-infra/fabricExtConnP-1/peeringP"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvPeeringP", "class_name", "fvPeeringP"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvPeeringP", "content.type", "automatic_with_rr"),
					resource.TestCheckResourceAttr("aci_rest_managed.fvPeeringP", "content.%", "1"),
				),
			},
		},
	})
}

var onDestroyTestName string

func init() {
	onDestroyTestName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
}

func TestAccAciRestManaged_onDestroyCreate(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(1g)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfig_onDestroy(),
				Check: resource.ComposeTestCheckFunc(
					// Check parent attributes
					resource.TestCheckResourceAttr("aci_rest_managed.tenant_with_destroy", "dn", "uni/tn-"+onDestroyTestName),
					resource.TestCheckResourceAttr("aci_rest_managed.tenant_with_destroy", "class_name", "fvTenant"),
					resource.TestCheckResourceAttr("aci_rest_managed.tenant_with_destroy", "content.name", onDestroyTestName),
					resource.TestCheckResourceAttr("aci_rest_managed.tenant_with_destroy", "content.descr", "Production Tenant"),

					// Check regular children (3 children)
					resource.TestCheckResourceAttr("aci_rest_managed.tenant_with_destroy", "child.#", "3"),

					// Check content_on_destroy
					resource.TestCheckResourceAttr("aci_rest_managed.tenant_with_destroy", "content_on_destroy.descr", "Decommissioned Tenant"),

					// Check child_on_destroy (2 children with destroy config)
					resource.TestCheckResourceAttr("aci_rest_managed.tenant_with_destroy", "child_on_destroy.#", "2"),
				),
			},
		},
	})
}

func TestAccAciRestManaged_onDestroyImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(1g)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        testAccAciRestManagedConfig_onDestroyMinimal(),
				ImportState:   true,
				ImportStateId: fmt.Sprintf("uni/tn-%s:ap-web,ctx-vrf1", onDestroyTestName),
				ResourceName:  "aci_rest_managed.tenant_with_destroy",
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 resource state, got %d", len(states))
					}

					state := states[0]

					// Validate parent
					if state.Attributes["dn"] != "uni/tn-"+onDestroyTestName {
						return fmt.Errorf("dn mismatch")
					}

					if state.Attributes["class_name"] != "fvTenant" {
						return fmt.Errorf("class_name mismatch")
					}

					if state.Attributes["content.descr"] != "Decommissioned Tenant" {
						return fmt.Errorf("content.descr mismatch: expected 'Decommissioned Tenant', got '%s'", state.Attributes["content.descr"])
					}

					// Check children count
					childCount := 2
					if state.Attributes["child.#"] != "2" {
						return fmt.Errorf("expected 2 children, got %s", state.Attributes["child.#"])
					}

					// Validate each child using helper
					if err := findAndValidateChildForRestManaged(state, childCount, "ap-web", "fvAp", map[string]string{
						"name":  "web",
						"descr": "Archived Web App",
					}); err != nil {
						return err
					}

					if err := findAndValidateChildForRestManaged(state, childCount, "ctx-vrf1", "fvCtx", map[string]string{
						"name":  "vrf1",
						"descr": "Archived VRF",
					}); err != nil {
						return err
					}

					return nil
				},
			},
			{
				Config:        testAccAciRestManagedConfig_onDestroyMinimal(),
				ImportState:   true,
				ImportStateId: fmt.Sprintf("uni/tn-%s:ap-database", onDestroyTestName),
				ResourceName:  "aci_rest_managed.tenant_with_destroy",
				ExpectError:   regexp.MustCompile("Unable to find specified child 'ap-database'"),
				Destroy:       true,
			},
		},
	})
}

// Explicit Destroy of resource with on_destroy attributes
func TestAccAciRestManaged_explpicitDestroy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "5.2(1g)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfig_onDestroyMinimal(),
			},
			{
				Config:  testAccAciRestManagedConfig_onDestroyMinimal(),
				Destroy: true,
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

func testAccAciRestManagedData_vrf(name string) string {
	return fmt.Sprintf(`
	data "aci_rest_managed" "%[1]s" {
		dn = "uni/tn-%[1]s/ctx-%[1]s"
	}
	`, name)
}

func testAccAciRestManagedConfig_ignoreChildAnnotation(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "ignore_child_annotation" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
		}
		child {
			rn         = "annotationKey-[%[1]s]"
			class_name = "tagAnnotation"
			content = {
				key = "%[1]s"
				value = "value"
			}
		}
	}
	`, name)
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

func testAccAciRestManagedConfig_importMultipleChildrenWithImportJsonString(name string, resource string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "tenant_%[2]s" {
		dn         = "uni/tn-%[1]s"
		class_name = "fvTenant"
		annotation = "orchestrator:from_resource"
		content    = {
			name = "%[1]s"
		}
	}

	resource "aci_rest_managed" "bd" {
		dn         = "uni/tn-%[1]s/BD-bd_test"
		class_name = "fvBD"
		content    = {
			name = "bd_test"
		}
		depends_on = [aci_rest_managed.tenant_%[2]s]
	}

	resource "aci_rest_managed" "rogue_exception_mac_11" {
		dn         = "uni/tn-%[1]s/BD-bd_test/rgexpmac-00:22:BD:F8:11:FF"
		class_name = "fvRogueExceptionMac"
		content    = {
			mac = "00:22:BD:F8:11:FF"
		}
		depends_on = [aci_rest_managed.bd]
	}

	resource "aci_rest_managed" "rogue_exception_mac_12" {
		dn         = "uni/tn-%[1]s/BD-bd_test/rgexpmac-00:22:BD:F8:12:FF"
		class_name = "fvRogueExceptionMac"
		content    = {
			mac = "00:22:BD:F8:12:FF"
		}
		depends_on = [aci_rest_managed.bd]
	}

	resource "aci_rest_managed" "eptags_%[2]s" {
		dn         = "${aci_rest_managed.tenant_%[2]s.id}/eptags"
		class_name = "fvEpTags"
		annotation = "orchestrator:from_resource"
		content    = {}
		child {
			rn         = "annotationKey-[~!$([])_+-={};:|,.]"
			class_name = "tagAnnotation"
			content    = {
				key   = "~!$([])_+-={};:|,."
				value = "tagAnnotation1"
			}
		}
		child {
			rn         = "annotationKey-[tagAnnotation2]"
			class_name = "tagAnnotation"
			content    = {
				key   = "tagAnnotation2"
				value = "tagAnnotation2"
			}
		}
		child {
			rn         = "epmactag-90:B5:B8:42:D1:88-[default]"
			class_name = "fvEpMacTag"
			content    = {
				mac	= "90:B5:B8:42:D1:88"
			}
		}
		child {
			rn         = "epiptag-[2001:10:1::1]-default"
			class_name = "fvEpIpTag"
			content    = {
				ip	= "2001:10:1::1"
			}
		}
		depends_on = [aci_rest_managed.tenant_%[2]s]
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

func testAccAciRestManagedConfig_globalAllowExisting(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "fvTenant" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
		}
	}
	resource "aci_rest_managed" "fvTenant_2" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
		}
		depends_on = [aci_rest_managed.fvTenant]
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
			message = "&&& WARNING &&& VERIFYING THE CONVERSION OF HTML CHARACTERS."
		}
	}
	`
}

func testAccAciRestManagedConfig_escapeHtmlTrue() string {
	return `
	resource "aci_rest_managed" "aaaPreLoginBanner" {
		dn          = "uni/userext/preloginbanner"
		class_name  = "aaaPreLoginBanner"
		escape_html = true
		content = {
			message = "&&& WARNING &&& VERIFYING THE CONVERSION OF HTML CHARACTERS."
		}
	}
	`
}

func testAccAciRestManagedConfig_undeletableObject() string {
	return `
	resource "aci_rest_managed" "fvFabricExtConnP" {
		dn         = "uni/tn-infra/fabricExtConnP-1"
		class_name = "fvFabricExtConnP"
		content = {
		  id     = "1"
		  name   = "IPN"
		  rt     = "extended:as2-nn4:5:16"
		  siteId = "0"
		}
	  }
	  resource "aci_rest_managed" "fvPeeringP" {
		dn          = "${aci_rest_managed.fvFabricExtConnP.dn}/peeringP"
		class_name  = "fvPeeringP"
		escape_html = false
		content = {
		  type     = "automatic_with_rr"
		}
	  }
	`
}

func testAccAciRestManagedConfig_onDestroy() string {
	return fmt.Sprintf(`
    resource "aci_rest_managed" "tenant_with_destroy" {
        dn         = "uni/tn-%[1]s"
        class_name = "fvTenant"
        content = {
            name  = "%[1]s"
            descr = "Production Tenant"
        }

        child {
            rn         = "ap-web"
            class_name = "fvAp"
            content = {
                name  = "web"
                descr = "Web Application"
            }
        }

        child {
            rn         = "ap-database"
            class_name = "fvAp"
            content = {
                name  = "database"
                descr = "Database Application"
            }
        }

        child {
            rn         = "ctx-vrf1"
            class_name = "fvCtx"
            content = {
                name  = "vrf1"
                descr = "Production VRF"
            }
        }

        content_on_destroy = {
            descr = "Decommissioned Tenant"
        }

        child_on_destroy {
            rn         = "ap-web"
            class_name = "fvAp"
            content = {
                descr = "Archived Web App"
            }
        }

        child_on_destroy {
            rn         = "ctx-vrf1"
            class_name = "fvCtx"
            content = {
                descr = "Archived VRF"
            }
        }
    }
    `, onDestroyTestName)
}

func testAccAciRestManagedConfig_onDestroyMinimal() string {
	return fmt.Sprintf(`
    resource "aci_rest_managed" "tenant_with_destroy" {
        dn         = "uni/tn-%[1]s"
        class_name = "fvTenant"
        content = {
            name = "%[1]s"
        }
    }
    `, onDestroyTestName)
}
