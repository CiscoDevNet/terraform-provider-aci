package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciBgpRouteControlProfileDS_Basic(t *testing.T) {
	resourceName := "aci_bgp_route_control_profile.test"
	dataSourceName := "data.aci_bgp_route_control_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBgpRouteControlProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBgpRouteControlProfileDSWithoutRequired(rName, rName, "parent_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBgpRouteControlProfileDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBgpRouteControlProfileDSConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "parent_dn", resourceName, "parent_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "route_control_profile_type", resourceName, "route_control_profile_type"),
				),
			},
			{
				Config:      CreateAccBgpRouteControlProfileDSUpdateRandomAttr(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccBgpRouteControlProfileDSWithInvalidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`Object may not exists`),
			},
			{
				Config: CreateAccBgpRouteControlProfileDSUpdate(rName, rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateBgpRouteControlProfileDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: Creating bgp_route_control_profile Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		name  = "%s"
		description = "created while acceptance testing"
	}
	
	`
	switch attrName {
	case "parent_dn":
		rBlock += `
	data "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		# name  = aci_bgp_route_control_profile.test.name
		depends_on = [aci_bgp_route_control_profile.test]
	}
	`
	case "name":
		rBlock += `
	data "aci_bgp_route_control_profile" "test" {
		#parent_dn  = aci_tenant.test.id
		name  = aci_bgp_route_control_profile.test.name
		depends_on = [aci_bgp_route_control_profile.test]
	}
	`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBgpRouteControlProfileDSUpdateRandomAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Println("=== STEP  Testing bgp_route_control_profile Data Source update with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		name  = "%s"
	}
	data "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		name  = aci_bgp_route_control_profile.test.name 
		%s = "%s"
		depends_on = [aci_bgp_route_control_profile.test]
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}

func CreateAccBgpRouteControlProfileDSConfig(fvTenantName, rName string) string {
<<<<<<< HEAD
<<<<<<< HEAD
	fmt.Println("=== STEP  Testing bgp_route_control_profile Data Source creation with required arguments only")
=======
	fmt.Println("=== STEP  Testing bgp_route_control_profile Data Source creation with required arguments only")
>>>>>>> f70b9a12 (updated typo of argument)
=======
	fmt.Println("=== STEP  Testing bgp_route_control_profile Data Source creation with required arguments only")
>>>>>>> f70b9a12ba23ca0b20bcf43464445aa3601ca1b3
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		name  = "%s"
	}
	data "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		name  = aci_bgp_route_control_profile.test.name 
		depends_on = [aci_bgp_route_control_profile.test]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBgpRouteControlProfileDSConfigL3Outside(ParentName, rName string) string {
	fmt.Println("=== STEP  Testing bgp_route_control_profile Data Source creation when parent is l3outside")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        name           = "%s"
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_l3_outside.test.id
		name  = "%s"
	}
	data "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_l3_outside.test.id
		name  = aci_bgp_route_control_profile.test.name 
		depends_on = [aci_bgp_route_control_profile.test]
	}
	`, ParentName, ParentName, rName)
	return resource
}

func CreateAccBgpRouteControlProfileDSWithInvalidParentDn(ParentName, rName string) string {
	fmt.Println("=== STEP  Testing bgp_route_control_profile Data Source creation with invalid parent_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		name  = "%s"
	}
	data "aci_bgp_route_control_profile" "test" {
		parent_dn  = "${aci_tenant.test.id}xyz"
		name  = aci_bgp_route_control_profile.test.name 
		depends_on = [aci_bgp_route_control_profile.test]
	}
	`, ParentName, rName)
	return resource
}

func CreateAccBgpRouteControlProfileDSUpdate(fvTenantName, rName, attribute, value string) string {
	fmt.Println("=== STEP  Testing bgp_route_control_profile Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	data "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		name  = aci_bgp_route_control_profile.test.name 
		depends_on = [aci_bgp_route_control_profile.test]
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
