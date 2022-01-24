package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outHSRPSecondaryVIPDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_hsrp_secondary_vip.test"
	dataSourceName := "data.aci_l3out_hsrp_secondary_vip.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.5.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outHSRPSecondaryVIPDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outHSRPSecondaryVIPDSWithoutRequired(rName, rName, rName, rName, rName, ip, "l3out_hsrp_interface_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outHSRPSecondaryVIPDSWithoutRequired(rName, rName, rName, rName, rName, ip, "ip"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPConfigDataSource(rName, rName, rName, rName, rName, ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "l3out_hsrp_interface_group_dn", resourceName, "l3out_hsrp_interface_group_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ip", resourceName, "ip"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "config_issues", resourceName, "config_issues"),
				),
			},
			{
				Config:      CreateAccL3outHSRPSecondaryVIPDataSourceUpdate(rName, rName, rName, rName, rName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL3outHSRPSecondaryVIPDSWithInvalidParentDn(rName, rName, rName, rName, rName, ip),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPDataSourceUpdateResource(rName, rName, rName, rName, rName, ip, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateL3outHSRPSecondaryVIPDSWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip, attrName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_secondary_vip Data Source without ", attrName)
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
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.test.id
		name = "%s"
		ip_obtain_mode = "learn"
	}

	resource "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s"
	}

	`
	switch attrName {
	case "l3out_hsrp_interface_group_dn":
		rBlock += `
	data "aci_l3out_hsrp_secondary_vip" "test" {
	#	l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_secondary_vip.test.l3out_hsrp_interface_group_dn
		ip  = aci_l3out_hsrp_interface_group.test.ip
	}
		`
	case "ip":
		rBlock += `
	data "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_secondary_vip.test.l3out_hsrp_interface_group_dn
	#	ip  = aci_l3out_hsrp_interface_group.test.ip
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip)
}

func CreateAccL3outHSRPSecondaryVIPConfigDataSource(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_secondary_vip Data Source with required arguments only")
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
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.test.id
		name = "%s"
		ip_obtain_mode = "learn"
	}

	resource "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s"
	}

	data "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_secondary_vip.test.l3out_hsrp_interface_group_dn
		ip  = aci_l3out_hsrp_secondary_vip.test.ip
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip)
	return resource
}
func CreateAccL3outHSRPSecondaryVIPDSWithInvalidParentDn(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_secondary_vip Data Source with Invalid l3out_hsrp_interface_group_dn")
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
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.test.id
		name = "%s"
		ip_obtain_mode = "learn"
	}

	resource "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s"
	}

	data "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = "${aci_l3out_hsrp_secondary_vip.test.l3out_hsrp_interface_group_dn}invalid"
		ip  = aci_l3out_hsrp_secondary_vip.test.ip
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip)
	return resource
}

func CreateAccL3outHSRPSecondaryVIPDataSourceUpdateResource(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip, key, value string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_secondary_vip Data Source with updated resource")
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
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.test.id
		name = "%s"
		ip_obtain_mode = "learn"
	}

	resource "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s"
		%s = "%s"
	}

	data "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_secondary_vip.test.l3out_hsrp_interface_group_dn
		ip  = aci_l3out_hsrp_secondary_vip.test.ip
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip, key, value)
	return resource
}

func CreateAccL3outHSRPSecondaryVIPDataSourceUpdate(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip, key, value string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_secondary_vip Data Source invalid attribute")
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
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.test.id
		name = "%s"
		ip_obtain_mode = "learn"
	}

	resource "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s"
	}

	data "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_secondary_vip.test.l3out_hsrp_interface_group_dn
		ip  = aci_l3out_hsrp_secondary_vip.test.ip
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip, key, value)
	return resource
}
