package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCloudSubnetDataSource_Basic(t *testing.T) {
	resourceName := "aci_cloud_subnet.test"
	dataSourceName := "data.aci_cloud_subnet.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("45.1.0.0/16")
	ip = fmt.Sprintf("%s/16", ip)
	ipOther, _ := acctest.RandIpAddress("45.2.0.0/17")
	ipOther = fmt.Sprintf("%s/17", ipOther)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudSubnetDSWithoutRequired(rName, ip, "ip"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudSubnetDSWithoutRequired(rName, ip, "cloud_cidr_pool_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudSubnetConfigDataSource(rName, ip),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "ip", resourceName, "ip"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cloud_cidr_pool_dn", resourceName, "cloud_cidr_pool_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "scope.#", resourceName, "scope.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "scope.0", resourceName, "scope.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "usage", resourceName, "usage"),
				),
			},
			{
				Config:      CreateAccCloudSubnetDataSourceUpdate(rName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccCloudSubnetDSWithInvalidParentDn(rName, ip, ipOther),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccCloudSubnetDataSourceUpdatedResource(rName, ip, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccCloudSubnetConfigDataSource(rName, ip string) string {
	fmt.Println("=== STEP  testing cloud_subnet Data Source with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn = aci_tenant.test.id
		primary_cidr = "%s"
		name = "%s"
		region = "us-east-1"
		cloud_vendor = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	resource "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn = aci_cloud_context_profile.test.id
		addr = "%s"
	}

	resource "aci_cloud_subnet" "test" {
		cloud_cidr_pool_dn = aci_cloud_cidr_pool.test.id
		ip  = "%s"
		zone = "uni/clouddomp/provp-aws/region-us-east-1/zone-us-east-1a"
	}

	data "aci_cloud_subnet" "test" {
		ip  = aci_cloud_subnet.test.ip
		cloud_cidr_pool_dn = aci_cloud_subnet.test.cloud_cidr_pool_dn
		depends_on = [ aci_cloud_subnet.test ]
	}
	`, rName, rName, ip, rName, ip, ip)
	return resource
}

func CreateCloudSubnetDSWithoutRequired(rName, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_subnet Data Source without ", attrName)
	rBlock := `
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn = aci_tenant.test.id
		primary_cidr = "%s"
		name = "%s"
		region = "us-east-1"
		cloud_vendor = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	resource "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn = aci_cloud_context_profile.test.id
		addr = "%s"
	}

	resource "aci_cloud_subnet" "test" {
		cloud_cidr_pool_dn = aci_cloud_cidr_pool.test.id
		ip  = "%s"
		zone = "uni/clouddomp/provp-aws/region-us-east-1/zone-us-east-1a"
	}
	`
	switch attrName {
	case "ip":
		rBlock += `
	data "aci_cloud_subnet" "test" {
	#	ip  = "%s"
		cloud_cidr_pool_dn = aci_cloud_subnet.test.cloud_cidr_pool_dn
		depends_on = [ aci_cloud_subnet.test ]
	}
		`
	case "cloud_cidr_pool_dn":
		rBlock += `
	data "aci_cloud_subnet" "test" {
		ip  = aci_cloud_subnet.test.ip
	#	cloud_cidr_pool_dn = aci_cloud_subnet.test.cloud_cidr_pool_dn
		depends_on = [ aci_cloud_subnet.test ]
	}
		`

	}
	return fmt.Sprintf(rBlock, rName, rName, ip, ip, ip)
}

func CreateAccCloudSubnetDSWithInvalidParentDn(rName, ip, ipOther string) string {
	fmt.Println("=== STEP  testing cloud_subnet Data Source with invalid ip")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn = aci_tenant.test.id
		primary_cidr = "%s"
		name = "%s"
		region = "us-east-1"
		cloud_vendor = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	resource "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn = aci_cloud_context_profile.test.id
		addr = "%s"
	}

	resource "aci_cloud_subnet" "test" {
		cloud_cidr_pool_dn = aci_cloud_cidr_pool.test.id
		ip  = "%s"
		zone = "uni/clouddomp/provp-aws/region-us-east-1/zone-us-east-1a"
	}

	data "aci_cloud_subnet" "test" {
		ip  = "%s"
		cloud_cidr_pool_dn = aci_cloud_subnet.test.cloud_cidr_pool_dn
		depends_on = [ aci_cloud_subnet.test ]
	}
	`, rName, rName, ip, rName, ip, ip, ipOther)
	return resource
}

func CreateAccCloudSubnetDataSourceUpdate(rName, ip, key, value string) string {
	fmt.Println("=== STEP  testing cloud_subnet Data Source with random attribute")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn = aci_tenant.test.id
		primary_cidr = "%s"
		name = "%s"
		region = "us-east-1"
		cloud_vendor = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	resource "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn = aci_cloud_context_profile.test.id
		addr = "%s"
	}

	resource "aci_cloud_subnet" "test" {
		cloud_cidr_pool_dn = aci_cloud_cidr_pool.test.id
		ip  = "%s"
		zone = "uni/clouddomp/provp-aws/region-us-east-1/zone-us-east-1a"
	}

	data "aci_cloud_subnet" "test" {
		ip  = aci_cloud_subnet.test.ip
		cloud_cidr_pool_dn = aci_cloud_subnet.test.cloud_cidr_pool_dn
		%s = "%s"
		depends_on = [ aci_cloud_subnet.test ]
	}
	`, rName, rName, ip, rName, ip, ip, key, value)
	return resource
}

func CreateAccCloudSubnetDataSourceUpdatedResource(rName, ip, key, value string) string {
	fmt.Println("=== STEP  testing cloud_subnet Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn = aci_tenant.test.id
		primary_cidr = "%s"
		name = "%s"
		region = "us-east-1"
		cloud_vendor = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	resource "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn = aci_cloud_context_profile.test.id
		addr = "%s"
	}

	resource "aci_cloud_subnet" "test" {
		cloud_cidr_pool_dn = aci_cloud_cidr_pool.test.id
		ip  = "%s"
		%s = "%s"
		zone = "uni/clouddomp/provp-aws/region-us-east-1/zone-us-east-1a"
	}

	data "aci_cloud_subnet" "test" {
		ip  = aci_cloud_subnet.test.ip
		cloud_cidr_pool_dn = aci_cloud_subnet.test.cloud_cidr_pool_dn
		depends_on = [ aci_cloud_subnet.test ]
	}
	`, rName, rName, ip, rName, ip, ip, key, value)
	return resource
}
