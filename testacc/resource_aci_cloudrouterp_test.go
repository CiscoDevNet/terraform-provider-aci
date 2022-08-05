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

func TestAccAciCloudVpnGateway_Basic(t *testing.T) {
	var cloud_vpn_gateway_default models.CloudVpnGateway
	var cloud_vpn_gateway_updated models.CloudVpnGateway
	resourceName := "aci_cloud_vpn_gateway.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	cidr, _ := acctest.RandIpAddress("10.200.0.0/16")
	cidr = fmt.Sprintf("%s/16", cidr)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudVpnGatewayWithoutRequired(rName, rName, rName, cidr, "cloud_context_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudVpnGatewayWithoutRequired(rName, rName, rName, cidr, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudVpnGatewayConfig(rName, rName, cidr, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudVpnGatewayExists(resourceName, &cloud_vpn_gateway_default),
					resource.TestCheckResourceAttr(resourceName, "cloud_context_profile_dn", fmt.Sprintf("uni/tn-%s/ctxprofile-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "num_instances", "1"),
					resource.TestCheckResourceAttr(resourceName, "cloud_router_profile_type", "vpn-gw"),
				),
			},
			{
				Config: CreateAccCloudVpnGatewayConfigWithOptionalValues(rName, rName, cidr, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudVpnGatewayExists(resourceName, &cloud_vpn_gateway_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_context_profile_dn", fmt.Sprintf("uni/tn-%s/ctxprofile-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_cloud_vpn_gateway"),
					resource.TestCheckResourceAttr(resourceName, "cloud_router_profile_type", "vpn-gw"),
					resource.TestCheckResourceAttr(resourceName, "num_instances", "1"),
					testAccCheckAciCloudVpnGatewayIdEqual(&cloud_vpn_gateway_default, &cloud_vpn_gateway_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccCloudVpnGatewayConfigUpdatedName(rName, rName, cidr, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccCloudVpnGatewayRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudVpnGatewayConfigWithRequiredParams(rNameUpdated, cidr, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudVpnGatewayExists(resourceName, &cloud_vpn_gateway_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_context_profile_dn", fmt.Sprintf("uni/tn-%s/ctxprofile-%s", rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciCloudVpnGatewayIdNotEqual(&cloud_vpn_gateway_default, &cloud_vpn_gateway_updated),
				),
			},
			{
				Config: CreateAccCloudVpnGatewayConfig(rName, rName, cidr, rName),
			},
			{
				Config: CreateAccCloudVpnGatewayConfigWithRequiredParams(rName, cidr, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudVpnGatewayExists(resourceName, &cloud_vpn_gateway_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_context_profile_dn", fmt.Sprintf("uni/tn-%s/ctxprofile-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciCloudVpnGatewayIdNotEqual(&cloud_vpn_gateway_default, &cloud_vpn_gateway_updated),
				),
			},
		},
	})
}

func TestAccAciCloudVpnGateway_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	cidr, _ := acctest.RandIpAddress("10.202.0.0/16")
	cidr = fmt.Sprintf("%s/16", cidr)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudVpnGatewayConfig(rName, rName, cidr, rName),
			},
			{
				Config:      CreateAccCloudVpnGatewayWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccCloudVpnGatewayUpdatedAttr(rName, rName, cidr, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudVpnGatewayUpdatedAttr(rName, rName, cidr, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudVpnGatewayUpdatedAttr(rName, rName, cidr, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccCloudVpnGatewayUpdatedAttr(rName, rName, cidr, rName, "num_instances", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccCloudVpnGatewayUpdatedAttr(rName, rName, cidr, rName, "cloud_router_profile_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccCloudVpnGatewayUpdatedAttr(rName, rName, cidr, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCloudVpnGatewayConfig(rName, rName, cidr, rName),
			},
		},
	})
}

func testAccCheckAciCloudVpnGatewayExists(name string, cloud_vpn_gateway *models.CloudVpnGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Vpn Gateway %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Vpn Gateway dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_vpn_gatewayFound := models.CloudVpnGatewayFromContainer(cont)
		if cloud_vpn_gatewayFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Vpn Gateway %s not found", rs.Primary.ID)
		}
		*cloud_vpn_gateway = *cloud_vpn_gatewayFound
		return nil
	}
}

func testAccCheckAciCloudVpnGatewayDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing cloud_vpn_gateway destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_vpn_gateway" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_vpn_gateway := models.CloudVpnGatewayFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Vpn Gateway %s Still exists", cloud_vpn_gateway.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudVpnGatewayIdEqual(m1, m2 *models.CloudVpnGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("cloud_vpn_gateway DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciCloudVpnGatewayIdNotEqual(m1, m2 *models.CloudVpnGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("cloud_vpn_gateway DNs are equal")
		}
		return nil
	}
}

func CreateCloudVpnGatewayWithoutRequired(fvTenantName, cloudCtxProfileName, rName, cidr, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_vpn_gateway creation without ", attrName)
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
	
	`
	switch attrName {
	case "cloud_context_profile_dn":
		rBlock += `
	resource "aci_cloud_vpn_gateway" "test" {
	#	cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_cloud_vpn_gateway" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, fvTenantName, cloudCtxProfileName, cidr, rName)
}

func CreateAccCloudVpnGatewayConfigWithRequiredParams(prName, cidr, rName string) string {
	fmt.Printf("=== STEP  testing cloud_vpn_gateway creation with parent resource name %s and resource name %s\n", prName, rName)
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
	`, prName, prName, prName, cidr, rName)
	return resource
}
func CreateAccCloudVpnGatewayConfigUpdatedName(fvTenantName, cloudCtxProfileName, cidr, rName string) string {
	fmt.Println("=== STEP  testing cloud_vpn_gateway creation with invalid name = ", rName)
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
	`, fvTenantName, fvTenantName, cloudCtxProfileName, cidr, rName)
	return resource
}

func CreateAccCloudVpnGatewayConfigInfraTenant(rName, cidr string) string {
	fmt.Println("=== STEP  testing cloud_vpn_gateway creation for cloud_router_profile_type = host-router")
	resource := fmt.Sprintf(`
	data "aci_tenant" "test" {
		name 		= "infra"
	}

	resource "aci_vrf" "test" {
		name = "%s"
		tenant_dn = data.aci_tenant.test.id
	}
	
	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = data.aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "us-west-1"
		cloud_vendor = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	
	resource "aci_cloud_vpn_gateway" "test" {
		cloud_context_profile_dn  = aci_cloud_context_profile.test.id
		name  = "%s"
		cloud_router_profile_type = "host-router"
	}
	`, rName, rName, cidr, rName)
	return resource
}

func CreateAccCloudVpnGatewayConfig(fvTenantName, cloudCtxProfileName, cidr, rName string) string {
	fmt.Println("=== STEP  testing cloud_vpn_gateway creation with required arguments only")
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
	`, fvTenantName, fvTenantName, cloudCtxProfileName, cidr, rName)
	return resource
}

func CreateAccCloudVpnGatewayWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing cloud_vpn_gateway creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_cloud_vpn_gateway" "test" {
		cloud_context_profile_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccCloudVpnGatewayConfigWithOptionalValues(fvTenantName, cloudCtxProfileName, cidr, rName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_vpn_gateway creation with optional parameters")
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
		cloud_context_profile_dn  = "${aci_cloud_context_profile.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_vpn_gateway"
		num_instances = "1"
		cloud_router_profile_type = "vpn-gw"
		
	}
	`, fvTenantName, fvTenantName, cloudCtxProfileName, cidr, rName)

	return resource
}

func CreateAccCloudVpnGatewayRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing cloud_vpn_gateway updation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_cloud_vpn_gateway" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_vpn_gateway"
		cloud_router_profile_type = "host-router"
		
	}
	`)

	return resource
}

func CreateAccCloudVpnGatewayUpdatedAttr(fvTenantName, cloudCtxProfileName, cidr, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing cloud_vpn_gateway attribute: %s = %s \n", attribute, value)
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
	`, fvTenantName, fvTenantName, cloudCtxProfileName, cidr, rName, attribute, value)
	return resource
}
