package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var primaryCidrVal string = "10.1.0.1/16"
var regionVal string = "us-west-1"
var cloudVendorVal string = "aws"

func TestAccAciCloudCIDRPool_Basic(t *testing.T) {
	var cloud_cidr_pool_default models.CloudCIDRPool
	var cloud_cidr_pool_updated models.CloudCIDRPool
	resourceName := "aci_cloud_cidr_pool.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	addr, _ := acctest.RandIpAddress("10.2.0.0/16")
	addr = fmt.Sprintf("%s/16", addr)
	addrUpdated, _ := acctest.RandIpAddress("10.3.0.0/17")
	addrUpdated = fmt.Sprintf("%s/17", addrUpdated)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudCIDRPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudCIDRPoolWithoutRequired(rName, rName, addr, "cloud_context_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudCIDRPoolWithoutRequired(rName, rName, addr, "addr"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudCIDRPoolConfig(rName, rName, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudCIDRPoolExists(resourceName, &cloud_cidr_pool_default),
					resource.TestCheckResourceAttr(resourceName, "cloud_context_profile_dn", fmt.Sprintf("uni/tn-%s/ctxprofile-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "addr", addr),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "primary", "no"),
				),
			},
			{
				Config: CreateAccCloudCIDRPoolConfigWithOptionalValues(rName, rName, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudCIDRPoolExists(resourceName, &cloud_cidr_pool_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_context_profile_dn", fmt.Sprintf("uni/tn-%s/ctxprofile-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "addr", addr),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_cloud_cidr_pool"),
					resource.TestCheckResourceAttr(resourceName, "primary", "no"),
					testAccCheckAciCloudCIDRPoolIdEqual(&cloud_cidr_pool_default, &cloud_cidr_pool_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccCloudCIDRPoolWithInavalidIP(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+`),
			},

			{
				Config:      CreateAccCloudCIDRPoolRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudCIDRPoolConfigWithRequiredParams(rNameUpdated, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudCIDRPoolExists(resourceName, &cloud_cidr_pool_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_context_profile_dn", fmt.Sprintf("uni/tn-%s/ctxprofile-%s", rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "addr", addr),
					testAccCheckAciCloudCIDRPoolIdNotEqual(&cloud_cidr_pool_default, &cloud_cidr_pool_updated),
				),
			},
			{
				Config: CreateAccCloudCIDRPoolConfig(rName, rName, addr),
			},
			{
				Config: CreateAccCloudCIDRPoolConfigWithRequiredParams(rName, addrUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudCIDRPoolExists(resourceName, &cloud_cidr_pool_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_context_profile_dn", fmt.Sprintf("uni/tn-%s/ctxprofile-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "addr", addrUpdated),
					testAccCheckAciCloudCIDRPoolIdNotEqual(&cloud_cidr_pool_default, &cloud_cidr_pool_updated),
				),
			},
		},
	})
}

func TestAccAciCloudCIDRPool_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	addr, _ := acctest.RandIpAddress("10.4.0.0/18")
	addr = fmt.Sprintf("%s/18", addr)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudCIDRPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudCIDRPoolConfig(rName, rName, addr),
			},
			{
				Config:      CreateAccCloudCIDRPoolWithInValidParentDn(rName, addr),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccCloudCIDRPoolUpdatedAttr(rName, rName, addr, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudCIDRPoolUpdatedAttr(rName, rName, addr, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudCIDRPoolUpdatedAttr(rName, rName, addr, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudCIDRPoolUpdatedAttr(rName, rName, addr, "primary", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccCloudCIDRPoolUpdatedAttr(rName, rName, addr, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCloudCIDRPoolConfig(rName, rName, addr),
			},
		},
	})
}

func TestAccAciCloudCIDRPool_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	addr1, _ := acctest.RandIpAddress("10.6.0.0/20")
	addr1 = fmt.Sprintf("%s/20", addr1)
	addr2, _ := acctest.RandIpAddress("10.7.0.0/21")
	addr2 = fmt.Sprintf("%s/21", addr2)
	addr3, _ := acctest.RandIpAddress("10.8.0.0/22")
	addr3 = fmt.Sprintf("%s/22", addr3)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudCIDRPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudCIDRPoolConfigMultiple(rName, addr1, addr2, addr3),
			},
		},
	})
}

func testAccCheckAciCloudCIDRPoolExists(name string, cloud_cidr_pool *models.CloudCIDRPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud CIDR Pool %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud CIDR Pool dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_cidr_poolFound := models.CloudCIDRPoolFromContainer(cont)
		if cloud_cidr_poolFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud CIDR Pool %s not found", rs.Primary.ID)
		}
		*cloud_cidr_pool = *cloud_cidr_poolFound
		return nil
	}
}

func testAccCheckAciCloudCIDRPoolDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing cloud_cidr_pool destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_cidr_pool" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_cidr_pool := models.CloudCIDRPoolFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud CIDR Pool %s Still exists", cloud_cidr_pool.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudCIDRPoolIdEqual(m1, m2 *models.CloudCIDRPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("cloud_cidr_pool DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciCloudCIDRPoolIdNotEqual(m1, m2 *models.CloudCIDRPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("cloud_cidr_pool DNs are equal")
		}
		return nil
	}
}

func CreateCloudCIDRPoolWithoutRequired(fvTenantName, cloudCtxProfileName, addr, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_cidr_pool creation without ", attrName)
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
	
	`
	switch attrName {
	case "cloud_context_profile_dn":
		rBlock += `
	resource "aci_cloud_cidr_pool" "test" {
	#	cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		addr  = "%s"
	}
		`
	case "addr":
		rBlock += `
	resource "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
	#	addr  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, fvTenantName, cloudCtxProfileName, primaryCidrVal, regionVal, cloudVendorVal, addr)
}

func CreateAccCloudCIDRPoolConfigWithRequiredParams(rName, addr string) string {
	fmt.Printf("=== STEP  testing cloud_cidr_pool creation with parent resource name %s and cloud_cidr_pool addr %s\n", rName, addr)
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
	`, rName, rName, rName, primaryCidrVal, regionVal, cloudVendorVal, addr)
	return resource
}

func CreateAccCloudCIDRPoolConfigMultiple(rName, addr1, addr2, addr3 string) string {
	fmt.Println("=== STEP  testing multiple cloud_cidr_pool creation with required arguments only")
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
	}

	resource "aci_cloud_cidr_pool" "test1" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		addr  = "%s"
	}

	resource "aci_cloud_cidr_pool" "test2" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		addr  = "%s"
	}
	`, rName, rName, rName, primaryCidrVal, regionVal, cloudVendorVal, addr1, addr2, addr3)
	return resource
}

func CreateAccCloudCIDRPoolConfig(fvTenantName, cloudCtxProfileName, addr string) string {
	fmt.Println("=== STEP  testing cloud_cidr_pool creation with required arguments only")
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
	}
	`, fvTenantName, fvTenantName, cloudCtxProfileName, primaryCidrVal, regionVal, cloudVendorVal, addr)
	return resource
}

func CreateAccCloudCIDRPoolWithInValidParentDn(rName, addr string) string {
	fmt.Println("=== STEP  Negative Case: testing cloud_cidr_pool creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn  = aci_tenant.test.id
		addr  = "%s"	
	}
	`, rName, addr)
	return resource
}

func CreateAccCloudCIDRPoolWithInavalidIP(rName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_cidr_pool creation with invalid ip")
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
		cloud_context_profile_dn  = "${aci_cloud_context_profile.test.id}"
		addr  = "%s"
	}
	`, rName, rName, rName, primaryCidrVal, regionVal, cloudVendorVal, rName)
	return resource
}

func CreateAccCloudCIDRPoolConfigWithOptionalValues(fvTenantName, cloudCtxProfileName, addr string) string {
	fmt.Println("=== STEP  Basic: testing cloud_cidr_pool creation with optional parameters")
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
		cloud_context_profile_dn  = "${aci_cloud_context_profile.test.id}"
		addr  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_cidr_pool"
	}
	`, fvTenantName, fvTenantName, cloudCtxProfileName, primaryCidrVal, regionVal, cloudVendorVal, addr)

	return resource
}

func CreateAccCloudCIDRPoolRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing cloud_cidr_pool updation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_cloud_cidr_pool" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_cloud_cidr_pool"
		primary = "no"
	}
	`)

	return resource
}

func CreateAccCloudCIDRPoolUpdatedAttr(fvTenantName, cloudCtxProfileName, addr, attribute, value string) string {
	fmt.Printf("=== STEP  testing cloud_cidr_pool attribute: %s = %s \n", attribute, value)
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
		%s = "%s"
	}
	`, fvTenantName, fvTenantName, cloudCtxProfileName, primaryCidrVal, regionVal, cloudVendorVal, addr, attribute, value)
	return resource
}
