package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSubnetDataSource_Basic(t *testing.T) {
	resourceName := "aci_subnet.test"
	dataSourceName := "data.aci_subnet.test"
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.20.0.0/16")
	ip = fmt.Sprintf("%s/16", ip)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz") // creating random string of 5 characters (to give as random parameter)
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccSubnetDSWithoutParentDn(rName, ip),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccSubnetDSWithoutIP(rName, ip),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSubnetConfigDataSource(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.#", resourceName, "ctrl.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.0", resourceName, "ctrl.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "preferred", resourceName, "preferred"),
					resource.TestCheckResourceAttrPair(dataSourceName, "scope.#", resourceName, "scope.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "scope.0", resourceName, "scope.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "virtual", resourceName, "virtual"),
					resource.TestCheckResourceAttrPair(dataSourceName, "parent_dn", resourceName, "parent_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ip", resourceName, "ip"),
				),
			},
			{
				Config:      CreateAccSubnetDataSourceUpdateRandomAttr(rName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccSubnetDSWithInvalidParentDn(rName, ip),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccSubnetDataSourceUpdate(rName, ip, "description", "description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccSubnetDataSourceUpdateRandomAttr(rName, ip, key, value string) string {
	fmt.Printf("=== STEP  Basic: testing subnet data source update for attribute: %s = %s \n", key, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_subnet" "test"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
	}

	data "aci_subnet" "test" {
		parent_dn = aci_bridge_domain.test.id
		ip = aci_subnet.test.ip
		%s = "%s"
	}
	`, rName, rName, ip, key, value)
	return resource
}

func CreateAccSubnetDataSourceUpdate(rName, ip, key, value string) string {
	fmt.Printf("=== STEP  Basic: testing subnet data source update for attribute: %s = %s \n", key, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_subnet" "test"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
	}

	data "aci_subnet" "test" {
		parent_dn = aci_bridge_domain.test.id
		ip = aci_subnet.test.ip
		%s = "%s"
	}
	`, rName, rName, ip, key, value)
	return resource
}

func CreateAccSubnetConfigDataSource(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation for data source test")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_subnet" "test"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
	}

	data "aci_subnet" "test" {
		parent_dn = aci_bridge_domain.test.id
		ip = aci_subnet.test.ip
	}
	`, rName, rName, ip)
	return resource
}

func CreateAccSubnetDSWithInvalidParentDn(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing Subnet reading with invalid IP")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_subnet" "test"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
	}

	data "aci_subnet" "test" {
		parent_dn = "${aci_subnet.test.parent_dn}xyz"
		ip = "%s"
	}
	`, rName, rName, ip, ip)
	return resource
}

func CreateAccSubnetDSWithoutParentDn(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing subnet reading without giving parent_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_subnet" "test"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
	}

	data "aci_subnet" "test" {
		ip = aci_subnet.test.ip
	}
	`, rName, rName, ip)
	return resource
}

func CreateAccSubnetDSWithoutIP(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing subnet reading without giving IP")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_subnet" "test"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
	}

	data "aci_subnet" "test" {
		parent_dn = aci_bridge_domain.test.id
	}
	`, rName, rName, ip)
	return resource
}
