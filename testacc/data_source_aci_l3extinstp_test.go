package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciExternalNetworkInstanceProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_external_network_instance_profile.test"
	dataSourceName := "data.aci_external_network_instance_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciExternalNetworkInstanceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccExternalNetworkInstanceProfileDSWithoutL3Out(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileConfigDS(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "exception_tag", resourceName, "exception_tag"),
					resource.TestCheckResourceAttrPair(dataSourceName, "flood_on_encap", resourceName, "flood_on_encap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_t", resourceName, "match_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pref_gr_memb", resourceName, "pref_gr_memb"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
					resource.TestCheckResourceAttrPair(dataSourceName, "target_dscp", resourceName, "target_dscp"),
				),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileDataSourceUpdate(rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileDSWithInvalidL3OutDn(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
		},
	})
}

func CreateAccExternalNetworkInstanceProfileDataSourceUpdate(rName, key, value string) string {
	fmt.Printf("=== STEP  Basic: testing external_network_instance_profile data source update for %s = %s\n", key, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_external_network_instance_profile" "test" {
        l3_outside_dn  = aci_l3_outside.test.id
		name = "%s"
		%s = "%s"
    }

	data "aci_external_network_instance_profile" "test" {
        l3_outside_dn  = aci_external_network_instance_profile.test.l3_outside_dn
		name = aci_external_network_instance_profile.test.name
    }
	`, rName, rName, rName, key, value)
	return resource
}

func CreateAccExternalNetworkInstanceProfileDSWithInvalidL3OutDn(rName string) string {
	fmt.Println("=== STEP  testing external_network_instance_profile data source creation with invalid l3_outside_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_external_network_instance_profile" "test" {
        l3_outside_dn  = aci_l3_outside.test.id
		name = "%s"
    }

	data "aci_external_network_instance_profile" "test" {
        l3_outside_dn  = "${aci_external_network_instance_profile.test.l3_outside_dn}xyz"
		name = aci_external_network_instance_profile.test.name
    }
	`, rName, rName, rName)
	return resource
}

func CreateAccExternalNetworkInstanceProfileDataSourceUpdateRandomAttr(rName, key, value string) string {
	fmt.Printf("=== STEP  Basic: testing external_network_instance_profile data source creation with %s = %s in data source\n", key, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_external_network_instance_profile" "test" {
        l3_outside_dn  = aci_l3_outside.test.id
		name = "%s"
    }

	data "aci_external_network_instance_profile" "test" {
        l3_outside_dn  = aci_external_network_instance_profile.test.l3_outside_dn
		name = aci_external_network_instance_profile.test.name
		%s = "%s"
    }
	`, rName, rName, rName, key, value)
	return resource
}

func CreateAccExternalNetworkInstanceProfileConfigDS(rName string) string {
	fmt.Println("=== STEP  Basic: testing external_network_instance_profile data source creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_external_network_instance_profile" "test" {
        l3_outside_dn  = aci_l3_outside.test.id
		name = "%s"
    }

	data "aci_external_network_instance_profile" "test" {
        l3_outside_dn  = aci_external_network_instance_profile.test.l3_outside_dn
		name = aci_external_network_instance_profile.test.name
    }
	`, rName, rName, rName)
	return resource
}

func CreateAccExternalNetworkInstanceProfileDSWithoutL3Out(rName string) string {
	fmt.Println("=== STEP  Basic: testing external_network_instance_profile data source creation without creating l3_outside")
	resource := fmt.Sprintf(`
	data "aci_external_network_instance_profile" "test" {
        name = "%s"
    }
	`, rName)
	return resource
}

func CreateAccExternalNetworkInstanceProfileDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing external_network_instance_profile data source creation without name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_external_network_instance_profile" "test" {
        l3_outside_dn  = aci_l3_outside.test.id
		name = "%s"
    }

	data "aci_external_network_instance_profile" "test" {
        l3_outside_dn  = aci_external_network_instance_profile.test.l3_outside_dn
    }
	`, rName, rName, rName)
	return resource
}
