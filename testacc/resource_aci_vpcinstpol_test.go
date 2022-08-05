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

func TestAccAciVPCDomainPolicy_Basic(t *testing.T) {
	var vpc_domain_policy_default models.VPCDomainPolicy
	var vpc_domain_policy_updated models.VPCDomainPolicy
	resourceName := "aci_vpc_domain_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVPCDomainPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateVPCDomainPolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVPCDomainPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCDomainPolicyExists(resourceName, &vpc_domain_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "dead_intvl", "200"),
				),
			},
			{
				Config: CreateAccVPCDomainPolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCDomainPolicyExists(resourceName, &vpc_domain_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_vpc_domain_policy"),
					resource.TestCheckResourceAttr(resourceName, "dead_intvl", "6"),

					testAccCheckAciVPCDomainPolicyIdEqual(&vpc_domain_policy_default, &vpc_domain_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccVPCDomainPolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccVPCDomainPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccVPCDomainPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCDomainPolicyExists(resourceName, &vpc_domain_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciVPCDomainPolicyIdNotEqual(&vpc_domain_policy_default, &vpc_domain_policy_updated),
				),
			},
		},
	})
}

func TestAccAciVPCDomainPolicy_Update(t *testing.T) {
	var vpc_domain_policy_default models.VPCDomainPolicy
	var vpc_domain_policy_updated models.VPCDomainPolicy
	resourceName := "aci_vpc_domain_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVPCDomainPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVPCDomainPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCDomainPolicyExists(resourceName, &vpc_domain_policy_default),
				),
			},
			{
				Config: CreateAccVPCDomainPolicyUpdatedAttr(rName, "dead_intvl", "600"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCDomainPolicyExists(resourceName, &vpc_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "dead_intvl", "600"),
					testAccCheckAciVPCDomainPolicyIdEqual(&vpc_domain_policy_default, &vpc_domain_policy_updated),
				),
			},
			{
				Config: CreateAccVPCDomainPolicyUpdatedAttr(rName, "dead_intvl", "297"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCDomainPolicyExists(resourceName, &vpc_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "dead_intvl", "297"),
					testAccCheckAciVPCDomainPolicyIdEqual(&vpc_domain_policy_default, &vpc_domain_policy_updated),
				),
			},

			{
				Config: CreateAccVPCDomainPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciVPCDomainPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVPCDomainPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVPCDomainPolicyConfig(rName),
			},

			{
				Config:      CreateAccVPCDomainPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVPCDomainPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVPCDomainPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccVPCDomainPolicyUpdatedAttr(rName, "dead_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVPCDomainPolicyUpdatedAttr(rName, "dead_intvl", "4"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccVPCDomainPolicyUpdatedAttr(rName, "dead_intvl", "601"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccVPCDomainPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccVPCDomainPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciVPCDomainPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVPCDomainPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVPCDomainPolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciVPCDomainPolicyExists(name string, vpc_domain_policy *models.VPCDomainPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VPC Domain Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC Domain Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vpc_domain_policyFound := models.VPCDomainPolicyFromContainer(cont)
		if vpc_domain_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VPC Domain Policy %s not found", rs.Primary.ID)
		}
		*vpc_domain_policy = *vpc_domain_policyFound
		return nil
	}
}

func testAccCheckAciVPCDomainPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vpc_domain_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vpc_domain_policy" {
			cont, err := client.Get(rs.Primary.ID)
			vpc_domain_policy := models.VPCDomainPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VPC Domain Policy %s Still exists", vpc_domain_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVPCDomainPolicyIdEqual(m1, m2 *models.VPCDomainPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("vpc_domain_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciVPCDomainPolicyIdNotEqual(m1, m2 *models.VPCDomainPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("vpc_domain_policy DNs are equal")
		}
		return nil
	}
}

func CreateVPCDomainPolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vpc_domain_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_vpc_domain_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccVPCDomainPolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing vpc_domain_policy creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_domain_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccVPCDomainPolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing vpc_domain_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_domain_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccVPCDomainPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing vpc_domain_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_domain_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccVPCDomainPolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple vpc_domain_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_domain_policy" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccVPCDomainPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing vpc_domain_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_domain_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vpc_domain_policy"
		dead_intvl = "6"
		
	}
	`, rName)

	return resource
}

func CreateAccVPCDomainPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vpc_domain_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_vpc_domain_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vpc_domain_policy"
		dead_intvl = "6"
		
	}
	`)

	return resource
}

func CreateAccVPCDomainPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vpc_domain_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_domain_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
