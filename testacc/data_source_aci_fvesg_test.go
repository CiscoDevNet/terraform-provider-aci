package acctest

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciEndpointSecurityGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_endpoint_security_group.test"
	dataSourceName := "data.aci_endpoint_security_group.test"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndpointSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccEndpointSecurityGroupDSWithoutAP(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEndpointSecurityGroupConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "application_profile_dn", resourceName, "application_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "flood_on_encap", resourceName, "flood_on_encap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_t", resourceName, "match_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pc_enf_pref", resourceName, "pc_enf_pref"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pref_gr_memb", resourceName, "pref_gr_memb"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
				),
			},
			{
				Config:      CreateAccEndpointSecurityGroupDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccEndpointSecurityGroupDataSourceUpdate(rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccEndpointSecurityGroupDataSourceUpdate(rName, key, value string) string {
	fmt.Printf("=== STEP  Basic: testing endpoint_security_group data source update for attribute: %s = %s \n", key, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	 }
	  
	resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	  
	resource "aci_endpoint_security_group" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
		%s = "%s"
	}

	data "aci_endpoint_security_group" "test" {
		name = aci_endpoint_security_group.test.name
		application_profile_dn = aci_endpoint_security_group.test.application_profile_dn
	}
	`, rName, rName, rName, key, value)
	return resource
}

func CreateAccEndpointSecurityGroupDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_security_group reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	 }
	  
	resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	  
	resource "aci_endpoint_security_group" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	data "aci_endpoint_security_group" "test" {
		name = "${aci_endpoint_security_group.test.name}xyz"
		application_profile_dn = aci_endpoint_security_group.test.application_profile_dn
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccEndpointSecurityGroupDataSourceUpdateRandomAttr(rName, key, val string) string {
	fmt.Printf("=== STEP  Basic: testing endpoint_security_group data source update for attribute: %s = %s \n", key, val)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	 }
	  
	resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	  
	resource "aci_endpoint_security_group" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	data "aci_endpoint_security_group" "test" {
		name = aci_endpoint_security_group.test.name
		application_profile_dn = aci_endpoint_security_group.test.application_profile_dn
		%s = "%s"
	}
	`, rName, rName, rName, key, val)
	return resource
}

func CreateAccEndpointSecurityGroupConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing endpoint_security_group data source creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	 }
	  
	resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	  
	resource "aci_endpoint_security_group" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	data "aci_endpoint_security_group" "test" {
		name = aci_endpoint_security_group.test.name
		application_profile_dn = aci_endpoint_security_group.test.application_profile_dn
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccEndpointSecurityGroupDSWithoutAP(rName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_security_group data source creation without creating application_profile_dn")
	resource := fmt.Sprintf(`
	data "aci_endpoint_security_group" "test" {
		name = "%s"
	  }
	`, rName)
	return resource
}

func CreateAccEndpointSecurityGroupDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_security_group creation without name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	  
	 resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	 }
	  
	 resource "aci_endpoint_security_group" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	 }

	 data "aci_endpoint_security_group" "test" {
		application_profile_dn = aci_endpoint_security_group.test.application_profile_dn
	  }
	`, rName, rName, rName)
	return resource
}
