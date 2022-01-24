package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciVRFDataSource_Basic(t *testing.T) {
	resourceName := "aci_vrf.test"
	dataSourceName := "data.aci_vrf.test"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccVRFDSWithoutTenant(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccVRFDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVRFConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bd_enforced_enable", resourceName, "bd_enforced_enable"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ip_data_plane_learning", resourceName, "ip_data_plane_learning"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "knw_mcast_act", resourceName, "knw_mcast_act"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pc_enf_dir", resourceName, "pc_enf_dir"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pc_enf_pref", resourceName, "pc_enf_pref"),
				),
			},
			{
				Config:      CreateAccVRFDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccVRFDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccVRFDataSourceUpdate(rName, "description", "description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccVRFDataSourceUpdateRandomAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing vrf data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_vrf" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	data "aci_vrf" "test" {
		name = aci_vrf.test.name
		tenant_dn = aci_vrf.test.tenant_dn
		%s = "%s"
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccVRFDataSourceUpdate(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing vrf data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_vrf" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		%s = "%s"
	}
	data "aci_vrf" "test" {
		name = aci_vrf.test.name
		tenant_dn = aci_vrf.test.tenant_dn
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccVRFConfigDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: testing vrf creation for data source test")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = aci_vrf.test.name
	}
	`, rName, rName)
	return resource
}

func CreateAccVRFDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing vrf reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "${aci_vrf.test.name}xyz"
	}
	`, rName, rName)
	return resource
}

func CreateAccVRFDSWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing vrf reading without giving tenant_dn")
	resource := fmt.Sprintf(`
	data "aci_vrf" "test" {
		name = "%s"
	}
	`, rName)
	return resource
}

func CreateAccVRFDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing vrf reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	data "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
	}
	`, rName)
	return resource
}
