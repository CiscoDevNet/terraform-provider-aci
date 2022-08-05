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

func TestAccAciCloudSubnet_Basic(t *testing.T) {
	var cloud_subnet_default models.CloudSubnet
	var cloud_subnet_updated models.CloudSubnet
	resourceName := "aci_cloud_subnet.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("45.3.0.0/16")
	ipUpdated := fmt.Sprintf("%s/17", ip)
	ip = fmt.Sprintf("%s/16", ip)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudSubnetDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateCloudSubnetWithoutRequired(rName, ip, "ip"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudSubnetWithoutRequired(rName, ip, "cloud_cidr_pool_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudSubnetConfig(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists(resourceName, &cloud_subnet_default),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "cloud_cidr_pool_dn", fmt.Sprintf("uni/tn-%s/ctxprofile-%s/cidr-[%s]", rName, rName, ip)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "name", ""),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "private"),
					resource.TestCheckResourceAttr(resourceName, "usage", "user"),
					resource.TestCheckResourceAttr(resourceName, "zone", "uni/clouddomp/provp-aws/region-us-east-1/zone-us-east-1a"),
				),
			},
			{
				Config: CreateAccCloudSubnetConfigWithOptionalValues(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists(resourceName, &cloud_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "cloud_cidr_pool_dn", fmt.Sprintf("uni/tn-%s/ctxprofile-%s/cidr-[%s]", rName, rName, ip)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_cloud_subnet"),
					resource.TestCheckResourceAttr(resourceName, "name", "test_cloud_subnet_name"),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "public"),
					resource.TestCheckResourceAttr(resourceName, "usage", "gateway"),
					resource.TestCheckResourceAttr(resourceName, "zone", "uni/clouddomp/provp-aws/region-us-east-1/zone-us-east-1a"),
					testAccCheckAciCloudSubnetIdEqual(&cloud_subnet_default, &cloud_subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccCloudSubnetWithInavalidIP(rName, ip),
				ExpectError: regexp.MustCompile(`Invalid RN`),
			},
			{
				Config:      CreateAccCloudSubnetRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccCloudSubnetConfigWithRequiredParams(rName, ip, ipUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists(resourceName, &cloud_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ip", ipUpdated),
					resource.TestCheckResourceAttr(resourceName, "cloud_cidr_pool_dn", fmt.Sprintf("uni/tn-%s/ctxprofile-%s/cidr-[%s]", rName, rName, ip)),
					testAccCheckAciCloudSubnetIdNotEqual(&cloud_subnet_default, &cloud_subnet_updated),
				),
			},
			{
				Config: CreateAccCloudSubnetConfig(rName, ip),
			},
			{
				Config: CreateAccCloudSubnetConfigWithRequiredParams(rNameUpdated, ip, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists(resourceName, &cloud_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "cloud_cidr_pool_dn", fmt.Sprintf("uni/tn-%s/ctxprofile-%s/cidr-[%s]", rNameUpdated, rNameUpdated, ip)),
					testAccCheckAciCloudSubnetIdNotEqual(&cloud_subnet_default, &cloud_subnet_updated),
				),
			},
		},
	})
}

func TestAccAciCloudSubnet_Update(t *testing.T) {
	var cloud_subnet_default models.CloudSubnet
	var cloud_subnet_updated models.CloudSubnet
	resourceName := "aci_cloud_subnet.test"
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("45.4.0.0/16")
	ip = fmt.Sprintf("%s/16", ip)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudSubnetConfig(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists(resourceName, &cloud_subnet_default),
				),
			},
			{

				Config: CreateAccCloudSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"private", "public"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists(resourceName, &cloud_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "private"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "public"),
					testAccCheckAciCloudSubnetIdEqual(&cloud_subnet_default, &cloud_subnet_updated),
				),
			},
			{

				Config: CreateAccCloudSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"private", "public", "shared"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists(resourceName, &cloud_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "private"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "public"),
					resource.TestCheckResourceAttr(resourceName, "scope.2", "shared"),
					testAccCheckAciCloudSubnetIdEqual(&cloud_subnet_default, &cloud_subnet_updated),
				),
			},
			{

				Config: CreateAccCloudSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"public", "shared"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists(resourceName, &cloud_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "public"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "shared"),
					testAccCheckAciCloudSubnetIdEqual(&cloud_subnet_default, &cloud_subnet_updated),
				),
			},
			{

				Config: CreateAccCloudSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"shared"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists(resourceName, &cloud_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "shared"),
					testAccCheckAciCloudSubnetIdEqual(&cloud_subnet_default, &cloud_subnet_updated),
				),
			},
			{
				Config: CreateAccCloudSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"shared", "public", "private"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists(resourceName, &cloud_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "shared"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "public"),
					resource.TestCheckResourceAttr(resourceName, "scope.2", "private"),
					testAccCheckAciCloudSubnetIdEqual(&cloud_subnet_default, &cloud_subnet_updated),
				),
			},
			{

				Config: CreateAccCloudSubnetUpdatedAttr(rName, ip, "usage", "infra-router"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists(resourceName, &cloud_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "usage", "infra-router"),
					testAccCheckAciCloudSubnetIdEqual(&cloud_subnet_default, &cloud_subnet_updated),
				),
			},
		},
	})
}

func TestAccAciCloudSubnet_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("45.5.0.0/16")
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudSubnetConfig(rName, ip),
			},
			{
				Config:      CreateAccCloudSubnetConfigWithInvalidParentDn(rName, ip),
				ExpectError: regexp.MustCompile(`Invalid DN`),
			},
			{
				Config:      CreateAccCloudSubnetUpdatedAttr(rName, ip, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudSubnetUpdatedAttr(rName, ip, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudSubnetUpdatedAttr(rName, ip, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudSubnetUpdatedAttr(rName, ip, "name", acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccCloudSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccCloudSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"private", "private"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccCloudSubnetUpdatedAttr(rName, ip, "usage", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccCloudSubnetWithInvalidZone(rName, ip, randomValue),
				ExpectError: regexp.MustCompile(`Relation target dn (.)+ not found`),
			},
			{
				Config:      CreateAccCloudSubnetUpdatedAttr(rName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCloudSubnetConfig(rName, ip),
			},
		},
	})
}

func TestAccAciCloudSubnet_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudSubnetConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciCloudSubnetExists(name string, cloud_subnet *models.CloudSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Subnet %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Subnet dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_subnetFound := models.CloudSubnetFromContainer(cont)
		if cloud_subnetFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Subnet %s not found", rs.Primary.ID)
		}
		*cloud_subnet = *cloud_subnetFound
		return nil
	}
}

func testAccCheckAciCloudSubnetDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing cloud_subnet destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_subnet" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_subnet := models.CloudSubnetFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Subnet %s Still exists", cloud_subnet.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudSubnetIdEqual(m1, m2 *models.CloudSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("cloud_subnet DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciCloudSubnetIdNotEqual(m1, m2 *models.CloudSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("cloud_subnet DNs are equal")
		}
		return nil
	}
}

func CreateCloudSubnetWithoutRequired(rName, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_subnet creation without ", attrName)
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
	`
	switch attrName {
	case "ip":
		rBlock += `
	resource "aci_cloud_subnet" "test" {
		cloud_cidr_pool_dn = aci_cloud_cidr_pool.test.id
	#	ip  = "%s"
	}
		`
	case "cloud_cidr_pool_dn":
		rBlock += `
	resource "aci_cloud_subnet" "test" {
	#	cloud_cidr_pool_dn = aci_cloud_cidr_pool.test.id
		ip  = "%s"
	}
		`

	}
	return fmt.Sprintf(rBlock, rName, rName, ip, rName, ip)
}

func CreateAccCloudSubnetConfigWithRequiredParams(rName, ip, ip_subnet string) string {
	fmt.Printf("=== STEP  testing cloud_subnet creation with parent resource name %s, parent resource ip %s and subnet ip %s\n", rName, ip, ip_subnet)
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
	`, rName, rName, ip, rName, ip, ip_subnet)
	return resource
}

func CreateAccCloudSubnetConfig(rName, ip string) string {
	fmt.Println("=== STEP  testing cloud_subnet creation with required arguments and zone")
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
	`, rName, rName, ip, rName, ip, ip)
	return resource
}

func CreateAccCloudSubnetConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple cloud_subnet creation with required arguments only")
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
		primary_cidr = "45.6.0.0/16"
		name = "%s"
		region = "us-east-1"
		cloud_vendor = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}

	resource "aci_cloud_cidr_pool" "test" {
		cloud_context_profile_dn = aci_cloud_context_profile.test.id
		addr = "45.6.0.0/16"
	}

	resource "aci_cloud_subnet" "test" {
		cloud_cidr_pool_dn = aci_cloud_cidr_pool.test.id
		ip  = "45.6.${count.index+1}.0/24"
		count = 5
		zone = "uni/clouddomp/provp-aws/region-us-east-1/zone-us-east-1a"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccCloudSubnetWithInavalidIP(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing cloud_subnet creation with invalid IP")
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
		name = "test_cloud_subnet_name"
		ip  = "%s"
		zone = "uni/clouddomp/provp-aws/region-us-east-1/zone-us-east-1a"
	}
	`, rName, rName, ip, rName, ip, rName)

	return resource
}

func CreateAccCloudSubnetConfigWithOptionalValues(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing cloud_subnet creation with optional parameters")
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
		name = "test_cloud_subnet_name"
		ip  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_subnet"
		scope = ["public"]
		usage = "gateway"
		zone = "uni/clouddomp/provp-aws/region-us-east-1/zone-us-east-1a"
	}
	`, rName, rName, ip, rName, ip, ip)

	return resource
}

func CreateAccCloudSubnetRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing cloud_subnet updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_cloud_subnet" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_subnet"
		scope = ["public"]
		usage = "gateway"
	}
	`)

	return resource
}

func CreateAccCloudSubnetConfigWithInvalidParentDn(rName, ip string) string {
	fmt.Println("=== STEP  testing cloud_subnet with invalid parent dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_cloud_subnet" "test" {
		cloud_cidr_pool_dn = aci_tenant.test.id
		ip  = "%s"
		zone = "uni/clouddomp/provp-aws/region-us-east-1/zone-us-east-1a"
	}
	`, rName, ip)
	return resource
}

func CreateAccCloudSubnetWithInvalidZone(rName, ip, value string) string {
	fmt.Printf("=== STEP  testing cloud_subnet attribute: zone = %s \n", value)
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
		zone = "%s"
	}
	`, rName, rName, ip, rName, ip, ip, value)
	return resource
}

func CreateAccCloudSubnetUpdatedAttr(rName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing cloud_subnet attribute: %s = %s \n", attribute, value)
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
		%s = "%s"
	}
	`, rName, rName, ip, rName, ip, ip, attribute, value)
	return resource
}

func CreateAccCloudSubnetUpdatedAttrList(rName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing cloud_subnet attribute: %s = %s \n", attribute, value)
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
		name = "test_cloud_subnet_name"
		ip  = "%s"
		zone = "uni/clouddomp/provp-aws/region-us-east-1/zone-us-east-1a"
		%s = %s
	}
	`, rName, rName, ip, rName, ip, ip, attribute, value)
	return resource
}
