package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFilterEntryDataSource_Basic(t *testing.T) {
	resourceName := "aci_filter_entry.test"
	dataSourceName := "data.aci_filter_entry.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFilterEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccFilterEntryDSWithoutFilter(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccFilterEntryDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFilterEntryConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "apply_to_frag", resourceName, "apply_to_frag"),
					resource.TestCheckResourceAttrPair(dataSourceName, "arp_opc", resourceName, "arp_opc"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ether_t", resourceName, "ether_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "icmpv4_t", resourceName, "icmpv4_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "icmpv6_t", resourceName, "icmpv6_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_dscp", resourceName, "match_dscp"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prot", resourceName, "prot"),
					resource.TestCheckResourceAttrPair(dataSourceName, "stateful", resourceName, "stateful"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tcp_rules", resourceName, "tcp_rules"),
					resource.TestCheckResourceAttrPair(dataSourceName, "filter_dn", resourceName, "filter_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
				),
			},
			{
				Config:      CreateAccFilterEntryUpdatedConfigDataSourceRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFilterEntryUpdatedConfigDataSource(rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
			{
				Config:      CreateAccFilterEntryDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
		},
	})
}

func CreateAccFilterEntryUpdatedConfigDataSourceRandomAttr(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing filter entry data source with Random Attribute")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name = "%s"
	}

	data "aci_filter_entry" "test" {
		filter_dn = aci_filter.test.id
		name = aci_filter_entry.test.name
		%s = "%s"
	}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccFilterEntryUpdatedConfigDataSource(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing filter entry data source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name = "%s"
		%s = "%s"
	}

	data "aci_filter_entry" "test" {
		filter_dn = aci_filter.test.id
		name = aci_filter_entry.test.name
	}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccFilterEntryConfigDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: Testing filter entry data source")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name = "%s"
		description = "test_description"
	}

	data "aci_filter_entry" "test" {
		filter_dn = aci_filter.test.id
		name = aci_filter_entry.test.name
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccFilterEntryDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter entry reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name = "%s"
	}

	data "aci_filter_entry" "test" {
		filter_dn = aci_filter.test.id
		name = "${aci_filter_entry.test.name}abc"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccFilterEntryDSWithoutFilter(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter entry reading without giving filter_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name = "%s"
	}

	data "aci_filter_entry" "test" {
		name = "%s"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccFilterEntryDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter entry reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name = "%s"
	}

	data "aci_filter_entry" "test" {
		filter_dn = aci_filter.test.id
	}
	`, rName, rName, rName)
	return resource
}
