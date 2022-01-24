package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciHsrpInterfacePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_hsrp_interface_policy.test"
	dataSourceName := "data.aci_hsrp_interface_policy.test"
	rName := acctest.RandString(5)
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciHsrpInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccHsrpInterfacePolicyDSWithoutTenant(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccHsrpInterfacePolicyDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccHsrpInterfacePolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl", resourceName, "ctrl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "reload_delay", resourceName, "reload_delay"),
					resource.TestCheckResourceAttrPair(dataSourceName, "delay", resourceName, "delay"),
				),
			},
			{
				Config:      CreateAccHsrpInterfacePolicyUpdatedConfigDataSourceRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccHsrpInterfacePolicyUpdatedConfigDataSource(rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
			{
				Config:      CreateAccHsrpInterfacePolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
		},
	})
}

func CreateAccHsrpInterfacePolicyUpdatedConfigDataSourceRandomAttr(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing Hrsp Interface Policy data source with random attribute")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_hsrp_interface_policy" "test" {
		tenant_dn = aci_hsrp_interface_policy.test.tenant_dn
		name = aci_hsrp_interface_policy.test.name
		%s = "%s"
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccHsrpInterfacePolicyUpdatedConfigDataSource(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing Hsrp Interface Policy data source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"

	}

	data "aci_hsrp_interface_policy" "test" {
		tenant_dn = aci_hsrp_interface_policy.test.tenant_dn
		name = aci_hsrp_interface_policy.test.name
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccHsrpInterfacePolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: Testing Hsrp Interface Policy data source")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_hsrp_interface_policy" "test" {
		tenant_dn = aci_hsrp_interface_policy.test.tenant_dn
		name = aci_hsrp_interface_policy.test.name
	}
	`, rName, rName)
	return resource
}

func CreateAccHsrpInterfacePolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing Hsrp Interface Policy reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_hsrp_interface_policy" "test" {
		tenant_dn = "${aci_hsrp_interface_policy.test.tenant_dn}abc"
		name = aci_hsrp_interface_policy.test.name
	}
	`, rName, rName)
	return resource
}

func CreateAccHsrpInterfacePolicyDSWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing Hsrp Interface Policy reading without giving tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_hsrp_interface_policy" "test" {
		name = aci_hsrp_interface_policy.test.name
	}
	`, rName, rName)
	return resource
}

func CreateAccHsrpInterfacePolicyDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing Hsrp Interface Policy reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_hsrp_interface_policy" "test" {
		tenant_dn = aci_hsrp_interface_policy.test.tenant_dn
	}
	`, rName, rName)
	return resource
}
