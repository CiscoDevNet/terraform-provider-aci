package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCloudVpnGatewayDataSource_Basic(t *testing.T) {
	resourceName := "aci_cloud_vpn_gateway.test"
	dataSourceName := "data.aci_cloud_vpn_gateway.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	cidr, _ := acctest.RandIpAddress("10.203.0.0/16")
	cidr = fmt.Sprintf("%s/16", cidr)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudVpnGatewayDSWithoutRequired(rName, rName, cidr, rName, "cloud_context_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudVpnGatewayDSWithoutRequired(rName, rName, cidr, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudVpnGatewayConfigDataSource(rName, rName, cidr, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "cloud_context_profile_dn", resourceName, "cloud_context_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "num_instances", resourceName, "num_instances"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cloud_router_profile_type", resourceName, "cloud_router_profile_type"),
				),
			},
			{
				Config:      CreateAccCloudVpnGatewayDataSourceUpdate(rName, rName, cidr, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccCloudVpnGatewayDSWithInvalidParentDn(rName, rName, cidr, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccCloudVpnGatewayDataSourceUpdatedResource(rName, rName, cidr, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccCloudVpnGatewayConfigDataSource(fvTenantName, cloudCtxProfileName, cidr, rName string) string {
	fmt.Println("=== STEP  testing cloud_vpn_gateway Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "us-west-1"
		cloud_vendor = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	
	resource "aci_cloud_vpn_gateway" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		name  = "%s"
	}

	data "aci_cloud_vpn_gateway" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		name  = aci_cloud_vpn_gateway.test.name
		depends_on = [ aci_cloud_vpn_gateway.test ]
	}
	`, fvTenantName, fvTenantName, cloudCtxProfileName, cidr, rName)
	return resource
}

func CreateCloudVpnGatewayDSWithoutRequired(fvTenantName, cloudCtxProfileName, cidr, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_vpn_gateway Data Source without ", attrName)
	rBlock := `
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "us-west-1"
		cloud_vendor = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	
	resource "aci_cloud_vpn_gateway" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "cloud_context_profile_dn":
		rBlock += `
	data "aci_cloud_vpn_gateway" "test" {
	#	cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		name  = aci_cloud_vpn_gateway.test.name
		depends_on = [ aci_cloud_vpn_gateway.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_cloud_vpn_gateway" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
	#	name  = aci_cloud_vpn_gateway.test.name
		depends_on = [ aci_cloud_vpn_gateway.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, fvTenantName, cloudCtxProfileName, cidr, rName)
}

func CreateAccCloudVpnGatewayDSWithInvalidParentDn(fvTenantName, cloudCtxProfileName, cidr, rName string) string {
	fmt.Println("=== STEP  testing cloud_vpn_gateway Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "us-west-1"
		cloud_vendor = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	
	resource "aci_cloud_vpn_gateway" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		name  = "%s"
	}

	data "aci_cloud_vpn_gateway" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		name  = "${aci_cloud_vpn_gateway.test.name}_invalid"
		depends_on = [ aci_cloud_vpn_gateway.test ]
	}
	`, fvTenantName, fvTenantName, cloudCtxProfileName, cidr, rName)
	return resource
}

func CreateAccCloudVpnGatewayDataSourceUpdate(fvTenantName, cloudCtxProfileName, cidr, rName, key, value string) string {
	fmt.Println("=== STEP  testing cloud_vpn_gateway Data Source with random attribute")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "us-west-1"
		cloud_vendor = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	
	resource "aci_cloud_vpn_gateway" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		name  = "%s"
	}


	data "aci_cloud_vpn_gateway" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		name  = aci_cloud_vpn_gateway.test.name
		%s = "%s"
		depends_on = [ aci_cloud_vpn_gateway.test ]
	}
	`, fvTenantName, fvTenantName, cloudCtxProfileName, cidr, rName, key, value)
	return resource
}

func CreateAccCloudVpnGatewayDataSourceUpdatedResource(fvTenantName, cloudCtxProfileName, cidr, rName, key, value string) string {
	fmt.Println("=== STEP  testing cloud_vpn_gateway Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "us-west-1"
		cloud_vendor = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	
	resource "aci_cloud_vpn_gateway" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_cloud_vpn_gateway" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		name  = aci_cloud_vpn_gateway.test.name
		depends_on = [ aci_cloud_vpn_gateway.test ]
	}
	`, fvTenantName, fvTenantName, cloudCtxProfileName, cidr, rName, key, value)
	return resource
}
