package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCloudContextProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_cloud_context_profile.test"
	dataSourceName := "data.aci_cloud_context_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	cidr, _ := acctest.RandIpAddress("30.1.0.0/16")
	cidr = fmt.Sprintf("%s/16", cidr)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudContextProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudContextProfileDSWithoutRequired(rName, rName, cidr, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudContextProfileDSWithoutRequired(rName, rName, cidr, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudContextProfileConfigDataSource(rName, rName, cidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config:      CreateAccCloudContextProfileDataSourceUpdate(rName, rName, cidr, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccCloudContextProfileDSWithInvalidName(rName, rName, cidr),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccCloudContextProfileDataSourceUpdatedResource(rName, rName, cidr, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccCloudContextProfileConfigDataSource(fvTenantName, rName, cidr string) string {
	fmt.Println("=== STEP  testing cloud_context_profile Data Source with required arguments only")
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
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	data "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_cloud_context_profile.test.name
		depends_on = [ aci_cloud_context_profile.test ]
	}
	`, fvTenantName, rName, rName, cidr, region, cloudVendor)
	return resource
}

func CreateCloudContextProfileDSWithoutRequired(fvTenantName, rName, cidr, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_context_profile Data Source without ", attrName)
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
		region = "%s"
		cloud_vendor = "%S"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_cloud_context_profile" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = aci_cloud_context_profile.test.name
		depends_on = [ aci_cloud_context_profile.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = aci_cloud_context_profile.test.name
		depends_on = [ aci_cloud_context_profile.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName, rName, cidr, region, cloudVendor)
}

func CreateAccCloudContextProfileDSWithInvalidName(fvTenantName, rName, cidr string) string {
	fmt.Println("=== STEP  testing cloud_context_profile Data Source with invalid name")
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
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	data "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_cloud_context_profile.test.name}_invalid"
		depends_on = [ aci_cloud_context_profile.test ]
	}
	`, fvTenantName, rName, rName, cidr, region, cloudVendor)
	return resource
}

func CreateAccCloudContextProfileDataSourceUpdate(fvTenantName, rName, cidr, key, value string) string {
	fmt.Println("=== STEP  testing cloud_context_profile Data Source with random attribute")
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
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	data "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_cloud_context_profile.test.name
		%s = "%s"
		depends_on = [ aci_cloud_context_profile.test ]
	}
	`, fvTenantName, rName, rName, cidr, region, cloudVendor, key, value)
	return resource
}

func CreateAccCloudContextProfileDataSourceUpdatedResource(fvTenantName, rName, cidr, key, value string) string {
	fmt.Println("=== STEP  testing cloud_context_profile Data Source with updated resource")
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
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
		%s = "%s"
	}

	data "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_cloud_context_profile.test.name
		depends_on = [ aci_cloud_context_profile.test ]
	}
	`, fvTenantName, rName, rName, cidr, region, cloudVendor, key, value)
	return resource
}
