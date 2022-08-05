package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCloudCIDRPoolDataSource_Basic(t *testing.T) {
	resourceName := "aci_cloud_cidr_pool.test"
	dataSourceName := "data.aci_cloud_cidr_pool.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	addr, _ := acctest.RandIpAddress("10.5.0.0/19")
	addr = fmt.Sprintf("%s/19", addr)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudCIDRPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudCIDRPoolDSWithoutRequired(rName, rName, addr, "cloud_context_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudCIDRPoolDSWithoutRequired(rName, rName, addr, "addr"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudCIDRPoolConfigDataSource(rName, rName, addr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "cloud_context_profile_dn", resourceName, "cloud_context_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "addr", resourceName, "addr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "primary", resourceName, "primary"),
				),
			},
			{
				Config:      CreateAccCloudCIDRPoolDataSourceUpdate(rName, rName, addr, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccCloudCIDRPoolDSWithInvalidParentDn(rName, rName, addr),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccCloudCIDRPoolDataSourceUpdatedResource(rName, rName, addr, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccCloudCIDRPoolConfigDataSource(fvTenantName, cloudCtxProfileName, addr string) string {
	fmt.Println("=== STEP  testing cloud_cidr_pool Data Source with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_context_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	resource "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		addr  = "%s"
	}

	data "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		addr  = aci_cloud_cidr_pool.test.addr
		depends_on = [ aci_cloud_cidr_pool.test ]
	}
	`, fvTenantName, fvTenantName, cloudCtxProfileName, primaryCidrVal, regionVal, cloudVendorVal, addr)
	return resource
}

func CreateCloudCIDRPoolDSWithoutRequired(fvTenantName, cloudCtxProfileName, addr, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_cidr_pool Data Source without ", attrName)
	rBlock := `

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_vrf" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_context_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	resource "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		addr  = "%s"
	}
	`
	switch attrName {
	case "cloud_context_profile_dn":
		rBlock += `
	data "aci_cloud_cidr_pool" "test" {
	#	cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		addr  = "%s"
		depends_on = [ aci_cloud_cidr_pool.test ]
	}
		`
	case "addr":
		rBlock += `
	data "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
	#	addr  = "%s"
		depends_on = [ aci_cloud_cidr_pool.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, fvTenantName, cloudCtxProfileName, primaryCidrVal, regionVal, cloudVendorVal, addr)
}

func CreateAccCloudCIDRPoolDSWithInvalidParentDn(fvTenantName, cloudCtxProfileName, addr string) string {
	fmt.Println("=== STEP  testing cloud_cidr_pool Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_context_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	resource "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		addr  = "%s"
	}

	data "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn  = "${aci_cloud_context_profile.test.id}_invalid"
		addr  = aci_cloud_cidr_pool.test.addr
		depends_on = [ aci_cloud_cidr_pool.test ]
	}
	`, fvTenantName, fvTenantName, cloudCtxProfileName, primaryCidrVal, regionVal, cloudVendorVal, addr)
	return resource
}

func CreateAccCloudCIDRPoolDataSourceUpdate(fvTenantName, cloudCtxProfileName, addr, key, value string) string {
	fmt.Println("=== STEP  testing cloud_cidr_pool Data Source with random attribute")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_context_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	resource "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		addr  = "%s"
	}

	data "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		addr  = aci_cloud_cidr_pool.test.addr
		%s = "%s"
		depends_on = [ aci_cloud_cidr_pool.test ]
	}
	`, fvTenantName, fvTenantName, cloudCtxProfileName, primaryCidrVal, regionVal, cloudVendorVal, addr, key, value)
	return resource
}

func CreateAccCloudCIDRPoolDataSourceUpdatedResource(fvTenantName, cloudCtxProfileName, addr, key, value string) string {
	fmt.Println("=== STEP  testing cloud_cidr_pool Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_vrf" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_context_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	resource "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		addr  = "%s"
		%s = "%s"
	}

	data "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn  = aci_cloud_cidr_pool.test.cloud_context_profile_dn
		addr  = aci_cloud_cidr_pool.test.addr
		depends_on = [ aci_cloud_cidr_pool.test ]
	}
	`, fvTenantName, fvTenantName, cloudCtxProfileName, primaryCidrVal, regionVal, cloudVendorVal, addr, key, value)
	return resource
}
