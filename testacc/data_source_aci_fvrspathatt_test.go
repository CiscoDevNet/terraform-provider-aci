package acctest

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciStaticPathDataSource_Basic(t *testing.T) {
	resourceName := "aci_epg_to_static_path.test"
	dataSourceName := "data.aci_epg_to_static_path.test"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciStaticPathDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccStaticPathDSWithoutEpg(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccStaticPathDSWithoutTdn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccStaticPathConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "application_epg_dn", resourceName, "application_epg_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tdn", resourceName, "tdn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "encap", resourceName, "encap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "instr_imedcy", resourceName, "instr_imedcy"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mode", resourceName, "mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "primary_encap", resourceName, "primary_encap"),
				),
			},
			{
				Config:      CreateAccStaticPathDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccStaticPathDSWithInvalidEpg(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccStaticPathDataSourceUpdate(rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccStaticPathDataSourceUpdateRandomAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing epg_to_static_path data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_application_epg" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_application_epg.test.id
		tdn = "%s"
		encap = "vlan-201"
	}

	data "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_epg_to_static_path.test.application_epg_dn
		tdn = aci_epg_to_static_path.test.tdn
		%s = "%s"
	}
	`, rName, rName, rName, tdn1, attribute, value)
	return resource
}

func CreateAccStaticPathDataSourceUpdate(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing epg_to_static_path data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_application_epg" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_application_epg.test.id
		tdn = "%s"
		encap = "vlan-201"
		%s = "%s"
	}

	data "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_epg_to_static_path.test.application_epg_dn
		tdn = aci_epg_to_static_path.test.tdn
	}
	`, rName, rName, rName, tdn1, attribute, value)
	return resource
}

func CreateAccStaticPathConfigDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: testing epg_to_static_path creation for data source test")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_application_epg" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_application_epg.test.id
		tdn = "%s"
		encap = "vlan-201"
	}

	data "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_epg_to_static_path.test.application_epg_dn
		tdn = aci_epg_to_static_path.test.tdn
	}
	`, rName, rName, rName, tdn1)
	return resource
}

func CreateAccStaticPathDSWithInvalidEpg(rName string) string {
	fmt.Println("=== STEP  Basic: testing epg_to_static_path reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_application_epg" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_application_epg.test.id
		tdn = "%s"
		encap = "vlan-201"
	}

	data "aci_epg_to_static_path" "test" {
		application_epg_dn = "${aci_epg_to_static_path.test.application_epg_dn}xyz"
		tdn = aci_epg_to_static_path.test.tdn
	}
	`, rName, rName, rName, tdn1)
	return resource
}

func CreateAccStaticPathDSWithoutEpg() string {
	fmt.Println("=== STEP  Basic: testing epg_to_static_path reading without giving application_epg_dn")
	resource := fmt.Sprintf(`
	data "aci_epg_to_static_path" "test" {
		tdn = "%s"
	}
	`, tdn1)
	return resource
}

func CreateAccStaticPathDSWithoutTdn(rName string) string {
	fmt.Println("=== STEP  Basic: testing epg_to_static_path reading without giving tdn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_application_epg" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	data "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_application_epg.test.id
	}
	`, rName, rName, rName)
	return resource
}
