package aci

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciRestManaged_tenant(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRestManagedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfigTenant(name, "Create description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRestManagedCheckObject("aci_rest_managed.fvTenant"),
				),
			},
			{
				Config: testAccAciRestManagedConfigTenant(name, "Updated description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRestManagedCheckObject("aci_rest_managed.fvTenant"),
				),
			},
		},
	})
}

func TestAccAciRestManaged_connPref(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRestManagedStillExists,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfigConnPref("ooband"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRestManagedCheckObject("aci_rest_managed.mgmtConnectivityPrefs"),
				),
			},
			{
				Config: testAccAciRestManagedConfigConnPref("inband"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRestManagedCheckObject("aci_rest_managed.mgmtConnectivityPrefs"),
				),
			},
		},
	})
}

func TestAccAciRestManaged_noContent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRestManagedStillExists,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfigConnPref(""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRestManagedCheckObject("aci_rest_managed.mgmtConnectivityPrefs"),
				),
			},
		},
	})
}

func testAccAciRestManagedConfigTenant(name string, description string) string {
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

func testAccAciRestManagedConfigConnPref(status string) string {
	if status != "" {
		return fmt.Sprintf(`
		resource "aci_rest_managed" "mgmtConnectivityPrefs" {
			dn = "uni/fabric/connectivityPrefs"
			class_name = "mgmtConnectivityPrefs"
			content = {
				interfacePref = "%[1]s"
			}
		}
		`, status)
	} else {
		return `
		resource "aci_rest_managed" "mgmtConnectivityPrefs" {
			dn = "uni/fabric/connectivityPrefs"
			class_name = "mgmtConnectivityPrefs"
		}
		`
	}
}

func testAccCheckAciRestManagedCheckObject(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Resource aci_rest_managed %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No aci_rest_managed dn attribute was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		className := rs.Primary.Attributes["class_name"]
		dn := models.StripQuotes(models.StripSquareBrackets(cont.Search("imdata", className, "attributes", "dn").String()))

		if dn != rs.Primary.ID {
			return fmt.Errorf("APIC object %s not found", rs.Primary.ID)
		}

		for key, value := range rs.Primary.Attributes {
			if strings.Contains(key, "content.") && key != "content.%" {
				attr := key[8:]
				v := models.StripQuotes(models.StripSquareBrackets(cont.Search("imdata", className, "attributes", attr).String()))
				if v != value {
					return fmt.Errorf("APIC object %s, expected: %s, got: %s", rs.Primary.ID, value, v)
				}
			}
		}

		return nil
	}
}

func testAccCheckAciRestManagedDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_rest_managed" {
			_, err := client.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Resource aci_rest_managed %s still exists", rs.Primary.ID)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciRestManagedStillExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_rest_managed" {
			_, err := client.Get(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error retrieving resource aci_rest_managed %s", rs.Primary.ID)
			}

		} else {
			continue
		}
	}

	return nil
}
