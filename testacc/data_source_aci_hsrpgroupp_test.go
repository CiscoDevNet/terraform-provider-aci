package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outHsrpInterfaceGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_hsrp_interface_group.test"
	dataSourceName := "data.aci_l3out_hsrp_interface_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outHsrpInterfaceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outHsrpInterfaceGroupDSWithoutRequired(rName, rName, rName, rName, rName, "l3out_hsrp_interface_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outHsrpInterfaceGroupDSWithoutRequired(rName, rName, rName, rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupDSConfig(rName, rName, rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrPair(dataSourceName, "l3out_hsrp_interface_profile_dn", resourceName, "l3out_hsrp_interface_profile_dn"),
						resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
						resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
						resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
						resource.TestCheckResourceAttrPair(dataSourceName, "config_issues", resourceName, "config_issues"),
						resource.TestCheckResourceAttrPair(dataSourceName, "group_af", resourceName, "group_af"),
						resource.TestCheckResourceAttrPair(dataSourceName, "group_id", resourceName, "group_id"),
						resource.TestCheckResourceAttrPair(dataSourceName, "group_name", resourceName, "group_name"),
						resource.TestCheckResourceAttrPair(dataSourceName, "ip", resourceName, "ip"),
						resource.TestCheckResourceAttrPair(dataSourceName, "ip_obtain_mode", resourceName, "ip_obtain_mode"),
						resource.TestCheckResourceAttrPair(dataSourceName, "mac", resourceName, "mac"),
						resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					),
				),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupDSConfigRandomAttr(rName, rName, rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupDSConfigUpdatedResource(rName, rName, rName, rName, rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					),
				),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupDSConfigInvalidName(rName, rName, rName, rName, rName),
				ExpectError: regexp.MustCompile(`Object may not exists`),
			},
		},
	})
}

func CreateAccL3outHsrpInterfaceGroupDSConfigInvalidName(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group Data Source creation with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}

  	resource "aci_logical_interface_profile" "test" {
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
    	ip = "10.20.30.40"
	}

	data "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_group.test.l3out_hsrp_interface_profile_dn
		name  = "${aci_l3out_hsrp_interface_group.test.name}xyz"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)
	return resource
}

func CreateAccL3outHsrpInterfaceGroupDSConfigUpdatedResource(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, key, value string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group Data Source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}

  	resource "aci_logical_interface_profile" "test" {
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
    	ip = "10.20.30.40"
		%s = "%s"
	}

	data "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_group.test.l3out_hsrp_interface_profile_dn
		name  = aci_l3out_hsrp_interface_group.test.name
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, key, value)
	return resource
}

func CreateAccL3outHsrpInterfaceGroupDSConfigRandomAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, key, value string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group Data Source with random attribute")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}

  	resource "aci_logical_interface_profile" "test" {
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
    	ip = "10.20.30.40"
	}

	data "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_group.test.l3out_hsrp_interface_profile_dn
		name  = aci_l3out_hsrp_interface_group.test.name
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, key, value)
	return resource
}

func CreateAccL3outHsrpInterfaceGroupDSConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group Data Source with required arguments")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}

  	resource "aci_logical_interface_profile" "test" {
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
    	ip = "10.20.30.40"
	}

	data "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_group.test.l3out_hsrp_interface_profile_dn
		name  = aci_l3out_hsrp_interface_group.test.name
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)
	return resource
}

func CreateL3outHsrpInterfaceGroupDSWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_interface_group Data Source creation without ", attrName)
	rBlock := `
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}

  	resource "aci_logical_interface_profile" "test" {
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
		ip = "10.20.30.40"
	}
	
	`
	switch attrName {
	case "l3out_hsrp_interface_profile_dn":
		rBlock += `
	data "aci_l3out_hsrp_interface_group" "test" {
	#	l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_group.test.l3out_hsrp_interface_profile_dn
		name  = aci_l3out_hsrp_interface_group.test.name
	}
		`
	case "name":
		rBlock += `
	data "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_group.test.l3out_hsrp_interface_profile_dn
	#	name  = aci_l3out_hsrp_interface_group.test.name
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)
}
