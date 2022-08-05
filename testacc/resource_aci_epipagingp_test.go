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

func TestAccAciIPAgingPolicy_Basic(t *testing.T) {
	var endpoint_ip_aging_profile_default models.IPAgingPolicy
	var endpoint_ip_aging_profile_updated models.IPAgingPolicy
	resourceName := "aci_endpoint_ip_aging_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciIPAgingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccIPAgingPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciIPAgingPolicyExists(resourceName, &endpoint_ip_aging_profile_default),
				),
				// all default values vary based on server
			},
			{
				Config: CreateAccIPAgingPolicyConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciIPAgingPolicyExists(resourceName, &endpoint_ip_aging_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_endpoint_ip_aging_profile"),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "enabled"),

					testAccCheckAciIPAgingPolicyIdEqual(&endpoint_ip_aging_profile_default, &endpoint_ip_aging_profile_updated),
				),
			},
			{
				Config: CreateAccIPAgingPolicyUpdatedAttr("admin_st", "disabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciIPAgingPolicyExists(resourceName, &endpoint_ip_aging_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "disabled"),
					testAccCheckAciIPAgingPolicyIdEqual(&endpoint_ip_aging_profile_default, &endpoint_ip_aging_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciIPAgingPolicy_Negative(t *testing.T) {

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciIPAgingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccIPAgingPolicyConfig(),
			},
			{
				Config:      CreateAccIPAgingPolicyUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccIPAgingPolicyUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccIPAgingPolicyUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccIPAgingPolicyUpdatedAttr("admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccIPAgingPolicyUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccIPAgingPolicyConfig(),
			},
		},
	})
}

func testAccCheckAciIPAgingPolicyExists(name string, endpoint_ip_aging_profile *models.IPAgingPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Endpoint Ip Aging Profile not found")
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Endpoint Ip Aging Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		endpoint_ip_aging_profileFound := models.IPAgingPolicyFromContainer(cont)
		if endpoint_ip_aging_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Endpoint Ip Aging Profile %s not found", rs.Primary.ID)
		}
		*endpoint_ip_aging_profile = *endpoint_ip_aging_profileFound
		return nil
	}
}

func testAccCheckAciIPAgingPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing endpoint_ip_aging_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_endpoint_ip_aging_profile" {
			cont, err := client.Get(rs.Primary.ID)
			endpoint_ip_aging_profile := models.IPAgingPolicyFromContainer(cont)
			if err != nil {
				return fmt.Errorf("Endpoint Ip Aging Profile %s Still exists", endpoint_ip_aging_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciIPAgingPolicyIdEqual(m1, m2 *models.IPAgingPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("endpoint_ip_aging_profile DNs are not equal default = %s : updated = %s ", m1.DistinguishedName, m2.DistinguishedName)
		}
		return nil
	}
}

func testAccCheckAciIPAgingPolicyIdNotEqual(m1, m2 *models.IPAgingPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("endpoint_ip_aging_profile DNs are equal")
		}
		return nil
	}
}

func CreateAccIPAgingPolicyConfig() string {
	fmt.Println("=== STEP  testing endpoint_ip_aging_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {

	}
	`)
	return resource
}

func CreateAccIPAgingPolicyConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing endpoint_ip_aging_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_endpoint_ip_aging_profile"
		admin_st = "enabled"
	}
	`)

	return resource
}

func CreateAccIPAgingPolicyUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing endpoint_ip_aging_profile attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
