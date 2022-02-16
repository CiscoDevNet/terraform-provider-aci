package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/terraform-providers/terraform-provider-aci/aci"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciCloudContextProfile_Basic(t *testing.T) {
	var cloud_context_profile_default models.CloudContextProfile
	var cloud_context_profile_updated models.CloudContextProfile
	resourceName := "aci_cloud_context_profile.test"
	cidr, _ := acctest.RandIpAddress("30.2.0.0/16")
	cidr = fmt.Sprintf("%s/16", cidr)
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudContextProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudContextProfileWithoutRequired(rName, rName, cidr, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudContextProfileWithoutRequired(rName, rName, cidr, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudContextProfileWithoutRequired(rName, rName, cidr, "primary_cidr"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudContextProfileWithoutRequired(rName, rName, cidr, "region"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudContextProfileWithoutRequired(rName, rName, cidr, "cloud_vendor"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudContextProfileWithoutRsToCtx(rName, cidr, "relation_cloud_rs_to_ctx"),
				ExpectError: regexp.MustCompile(`Invalid Configuration`),
			},
			{
				Config: CreateAccCloudContextProfileConfig(rName, rName, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudContextProfileExists(resourceName, &cloud_context_profile_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "cloud_vendor", cloudVendor),
					resource.TestCheckResourceAttr(resourceName, "primary_cidr", cidr),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "relation_cloud_rs_to_ctx", fmt.Sprintf("uni/tn-%s/ctx-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "regular"),
				),
			},
			{
				Config: CreateAccCloudContextProfileConfigWithOptionalValues(rName, rName, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudContextProfileExists(resourceName, &cloud_context_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_cloud_context_profile"),
					resource.TestCheckResourceAttr(resourceName, "cloud_vendor", cloudVendor),
					resource.TestCheckResourceAttr(resourceName, "primary_cidr", cidr),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "relation_cloud_rs_to_ctx", fmt.Sprintf("uni/tn-%s/ctx-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "shadow"),
					resource.TestCheckResourceAttr(resourceName, "hub_network", hubNetwork),
					testAccCheckAciCloudContextProfileIdEqual(&cloud_context_profile_default, &cloud_context_profile_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_vendor", "hub_network", "relation_cloud_rs_to_ctx"},
			},
			{
				Config:      CreateAccCloudContextProfileConfigUpdatedName(rName, acctest.RandString(65), cidr),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccCloudContextProfileConfigWithInvalidCidr(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccCloudContextProfileConfigWithInvalidRegion(rName, cidr),
				ExpectError: regexp.MustCompile(`Invalid Configuration`),
			},
			{
				Config:      CreateAccCloudContextProfileConfigWithInvalidVendor(rName, cidr),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccCloudContextProfileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudContextProfileConfigWithRequiredParams(rNameUpdated, rName, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudContextProfileExists(resourceName, &cloud_context_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "relation_cloud_rs_to_ctx", fmt.Sprintf("uni/tn-%s/ctx-%s", rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciCloudContextProfileIdNotEqual(&cloud_context_profile_default, &cloud_context_profile_updated),
				),
			},
			{
				Config: CreateAccCloudContextProfileConfig(rName, rName, cidr),
			},
			{
				Config: CreateAccCloudContextProfileConfigWithRequiredParams(rName, rNameUpdated, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudContextProfileExists(resourceName, &cloud_context_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "relation_cloud_rs_to_ctx", fmt.Sprintf("uni/tn-%s/ctx-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciCloudContextProfileIdNotEqual(&cloud_context_profile_default, &cloud_context_profile_updated),
				),
			},
		},
	})
}

func TestAccAciCloudContextProfile_Update(t *testing.T) {
	var cloud_context_profile_default models.CloudContextProfile
	var cloud_context_profile_updated models.CloudContextProfile
	resourceName := "aci_cloud_context_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	cidr, _ := acctest.RandIpAddress("30.3.0.0/16")
	cidr = fmt.Sprintf("%s/16", cidr)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudContextProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudContextProfileConfig(rName, rName, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudContextProfileExists(resourceName, &cloud_context_profile_default),
				),
			},

			{
				Config: CreateAccCloudContextProfileUpdatedAttr(rName, rName, cidr, "type", "hosted"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudContextProfileExists(resourceName, &cloud_context_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "type", "hosted"),
					testAccCheckAciCloudContextProfileIdEqual(&cloud_context_profile_default, &cloud_context_profile_updated),
				),
			},
			{
				Config: CreateAccCloudContextProfileUpdatedAttr(rName, rName, cidr, "type", "container-overlay"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudContextProfileExists(resourceName, &cloud_context_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "type", "container-overlay"),
					testAccCheckAciCloudContextProfileIdEqual(&cloud_context_profile_default, &cloud_context_profile_updated),
				),
			},
			{
				Config: CreateAccCloudContextProfileConfig(rName, rName, cidr),
			},
		},
	})
}

func TestAccAciCloudContextProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	cidr, _ := acctest.RandIpAddress("30.4.0.0/16")
	cidr = fmt.Sprintf("%s/16", cidr)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudContextProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudContextProfileConfig(rName, rName, cidr),
			},
			{
				Config:      CreateAccCloudContextProfileWithInValidParentDn(rName, cidr),
				ExpectError: regexp.MustCompile(`Invalid DN`),
			},
			{
				Config:      CreateAccCloudContextProfileUpdatedAttr(rName, rName, cidr, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudContextProfileUpdatedAttr(rName, rName, cidr, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudContextProfileUpdatedAttr(rName, rName, cidr, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccCloudContextProfileUpdatedAttr(rName, rName, cidr, "type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccCloudContextProfileUpdatedAttr(rName, rName, cidr, "hub_network", randomValue),
				ExpectError: regexp.MustCompile(`Relation target dn (.)+ not found`),
			},
			{
				Config:      CreateAccCloudContextProfileUpdatedAttr(rName, rName, cidr, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCloudContextProfileConfig(rName, rName, cidr),
			},
		},
	})
}

func testAccCheckAciCloudContextProfileExists(name string, cloud_context_profile *models.CloudContextProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Context Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Context Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cloud_context_profileFound, err := aci.GetRemoteCloudContextProfile(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if cloud_context_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Context Profile %s not found", rs.Primary.ID)
		}
		*cloud_context_profile = *cloud_context_profileFound
		return nil
	}
}

func testAccCheckAciCloudContextProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing cloud_context_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_context_profile" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_context_profile := models.CloudContextProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Context Profile %s Still exists", cloud_context_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudContextProfileIdEqual(m1, m2 *models.CloudContextProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("cloud_context_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciCloudContextProfileIdNotEqual(m1, m2 *models.CloudContextProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("cloud_context_profile DNs are equal")
		}
		return nil
	}
}

func CreateCloudContextProfileWithoutRsToCtx(rName, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_context_profile creation without relation_cloud_rs_to_ctx")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
	}
	`, rName, rName, ip, region, cloudVendor)
	return resource
}

func CreateCloudContextProfileWithoutRequired(fvTenantName, rName, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_context_profile creation without ", attrName)
	rBlock := `
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_cloud_context_profile" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
		`
	case "name":
		rBlock += `
	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
		`
	case "primary_cidr":
		rBlock += `
	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	#	primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
		`
	case "region":
		rBlock += `
	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
	#	region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	`
	case "cloud_vendor":
		rBlock += `
	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "%s"
	#	cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName, rName, ip, region, cloudVendor)
}

func CreateAccCloudContextProfileConfigWithRequiredParams(fvTenantName, rName, ip string) string {
	fmt.Printf("=== STEP  testing cloud_context_profile creation with parent resource name %s, resource name %s and cidr %s\n", fvTenantName, rName, ip)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	`, fvTenantName, fvTenantName, rName, ip, region, cloudVendor)
	return resource
}

func CreateAccCloudContextProfileConfigWithInvalidVendor(rName, ip string) string {
	fmt.Println("=== STEP  testing cloud_context_profile creation with invalid cloud_vendor")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	`, rName, rName, rName, ip, region, rName)
	return resource
}

func CreateAccCloudContextProfileConfigWithInvalidRegion(rName, ip string) string {
	fmt.Println("=== STEP  testing cloud_context_profile creation with invalid region")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	`, rName, rName, rName, ip, rName, cloudVendor)
	return resource
}

func CreateAccCloudContextProfileConfigWithInvalidCidr(rName string) string {
	fmt.Println("=== STEP  testing cloud_context_profile creation with invalid primary_cidr")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	`, rName, rName, rName, rName, region, cloudVendor)
	return resource
}

func CreateAccCloudContextProfileConfigUpdatedName(fvTenantName, rName, ip string) string {
	fmt.Println("=== STEP  testing cloud_context_profile creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	`, fvTenantName, fvTenantName, rName, ip, region, cloudVendor)
	return resource
}

func CreateAccCloudContextProfileConfig(fvTenantName, rName, ip string) string {
	fmt.Println("=== STEP  testing cloud_context_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	`, fvTenantName, fvTenantName, rName, ip, region, cloudVendor)
	return resource
}

func CreateAccCloudContextProfileWithInValidParentDn(rName, ip string) string {
	fmt.Println("=== STEP  Negative Case: testing cloud_context_profile creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_vrf.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
	}
	`, rName, rName, rName, ip, region, cloudVendor)
	return resource
}

func CreateAccCloudContextProfileConfigWithOptionalValues(fvTenantName, rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing cloud_context_profile creation with optional parameters")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_cloud_context_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		primary_cidr = "%s"
		region = "%s"
		cloud_vendor = "%s"
		relation_cloud_rs_to_ctx = aci_vrf.test.id
		name_alias = "test_cloud_context_profile"
		annotation = "orchestrator:terraform_testacc"
		description = "created while acceptance testing"
		type = "shadow"
		hub_network = "%s"
	}
	`, fvTenantName, rName, rName, ip, region, cloudVendor, hubNetwork)

	return resource
}

func CreateAccCloudContextProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing cloud_context_profile updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_cloud_context_profile" "test" {
		name_alias = "test_cloud_context_profile"
		annotation = "orchestrator:terraform_testacc"
		description = "created while acceptance testing"
		type = "shadow"
		hub_network = "%s"
	}
	`, hubNetwork)

	return resource
}

func CreateAccCloudContextProfileUpdatedAttr(fvTenantName, rName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing cloud_context_profile attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
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
	`, fvTenantName, fvTenantName, rName, ip, region, cloudVendor, attribute, value)
	return resource
}
