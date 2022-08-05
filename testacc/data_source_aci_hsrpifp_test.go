package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outHSRPInterfaceProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_hsrp_interface_profile.test"
	dataSourceName := "data.aci_l3out_hsrp_interface_profile.test"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outHSRPInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outHSRPInterfaceProfileDSWithoutRequired(rName, rName, rName, rName, "logical_interface_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outHSRPInterfaceProfileConfigDataSource(rName, rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "logical_interface_profile_dn", resourceName, "logical_interface_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "version", resourceName, "version"),
				),
			},
			{
				Config:      CreateAccL3outHSRPInterfaceProfileDataSourceUpdate(rName, rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL3outHSRPInterfaceProfileDSWithInvalidParentDn(rName, rName, rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL3outHSRPInterfaceProfileDataSourceUpdateResource(rName, rName, rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateL3outHSRPInterfaceProfileDSWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, attrName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_profile Data Source without ", attrName)
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
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}
	
	`
	switch attrName {
	case "logical_interface_profile_dn":
		rBlock += `
	data "aci_l3out_hsrp_interface_profile" "test" {
	#	logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}
		`

	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName)
}

func CreateAccL3outHSRPInterfaceProfileConfigDataSource(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_profile Data Source with required arguments only")
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
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}

	data "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		depends_on = [
			aci_l3out_hsrp_interface_profile.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName)
	return resource
}

func CreateAccL3outHSRPInterfaceProfileDSWithInvalidParentDn(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_profile Data Source Invalid logical_interface_profile_dn")
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
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}

	data "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = "${aci_logical_interface_profile.test.id}xyz"
		depends_on = [
			aci_l3out_hsrp_interface_profile.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName)
	return resource
}

func CreateAccL3outHSRPInterfaceProfileDataSourceUpdateResource(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, key, value string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_profile Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}

	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		%s = "%s"
	}

	data "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		depends_on = [
			aci_l3out_hsrp_interface_profile.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, key, value)
	return resource
}

func CreateAccL3outHSRPInterfaceProfileDataSourceUpdate(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, key, value string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_profile Data Source with random attributes")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}

	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}

	data "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		%s = "%s"
		depends_on = [
			aci_l3out_hsrp_interface_profile.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, key, value)
	return resource
}
