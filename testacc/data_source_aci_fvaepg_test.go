package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciApplicationEPGDataSource_Basic(t *testing.T) {
	resourceName := "aci_application_epg.test"
	dataSourceName := "data.aci_application_epg.test"
	rName := acctest.RandString(5)
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationEPGDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccApplicationEPGDSWithoutApplicationProfile(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccApplicationEPGDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccApplicationEPGConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
					resource.TestCheckResourceAttrPair(dataSourceName, "exception_tag", resourceName, "exception_tag"),
					resource.TestCheckResourceAttrPair(dataSourceName, "flood_on_encap", resourceName, "flood_on_encap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fwd_ctrl", resourceName, "fwd_ctrl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "has_mcast_source", resourceName, "has_mcast_source"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_attr_based_epg", resourceName, "is_attr_based_epg"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tanant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_t", resourceName, "match_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pc_enf_pref", resourceName, "pc_enf_pref"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pref_gr_memb", resourceName, "pref_gr_memb"),
					resource.TestCheckResourceAttrPair(dataSourceName, "shutdown", resourceName, "shutdown"),
				),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedConfigDataSourceRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccApplicationEPGUpdatedConfigDataSource(rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
			{
				Config:      CreateAccApplicationEPGDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
		},
	})
}

func CreateAccApplicationEPGUpdatedConfigDataSourceRandomAttr(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg data source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	data "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = aci_application_epg.test.name
		%s = "%s"
	}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccApplicationEPGUpdatedConfigDataSource(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg data source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
		%s = "%s"
	}

	data "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = aci_application_epg.test.name
	}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccApplicationEPGConfigDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg data source")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
		description = "test_description"
	}

	data "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = aci_application_epg.test.name
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccApplicationEPGDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing Application EPG reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	data "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "${aci_application_epg.test.name}abc"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccApplicationEPGDSWithoutApplicationProfile(rName string) string {
	fmt.Println("=== STEP  Basic: testing Application EPG reading without giving application_profile_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	data "aci_application_epg" "test" {
		name = "%s"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccApplicationEPGDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing Application EPG reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	data "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
	}
	`, rName, rName, rName)
	return resource
}
