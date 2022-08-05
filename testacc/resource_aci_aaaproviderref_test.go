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

func TestAccAciProviderGroupMember_Basic(t *testing.T) {
	var login_domain_provider_default models.ProviderGroupMember
	var login_domain_provider_updated models.ProviderGroupMember
	resourceName := "aci_login_domain_provider.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	aaaDuoProviderGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciProviderGroupMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateProviderGroupMemberWithoutRequired(aaaDuoProviderGroupName, rName, "parent_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateProviderGroupMemberWithoutRequired(aaaDuoProviderGroupName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccProviderGroupMemberConfig(aaaDuoProviderGroupName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciProviderGroupMemberExists(resourceName, &login_domain_provider_default),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/userext/duoext/duoprovidergroup-%s", aaaDuoProviderGroupName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "order", "0"),
				),
			},
			{
				Config: CreateAccProviderGroupMemberConfigWithOptionalValues(aaaDuoProviderGroupName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciProviderGroupMemberExists(resourceName, &login_domain_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/userext/duoext/duoprovidergroup-%s", aaaDuoProviderGroupName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_login_domain_provider"),
					resource.TestCheckResourceAttr(resourceName, "order", "lowest-available"),
					testAccCheckAciProviderGroupMemberIdEqual(&login_domain_provider_default, &login_domain_provider_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"order"},
			},
			{
				Config:      CreateAccProviderGroupMemberConfigUpdatedName(aaaDuoProviderGroupName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccProviderGroupMemberRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccProviderGroupMemberConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciProviderGroupMemberExists(resourceName, &login_domain_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/userext/duoext/duoprovidergroup-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciProviderGroupMemberIdNotEqual(&login_domain_provider_default, &login_domain_provider_updated),
				),
			},
			{
				Config: CreateAccProviderGroupMemberConfig(aaaDuoProviderGroupName, rName),
			},
			{
				Config: CreateAccProviderGroupMemberConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciProviderGroupMemberExists(resourceName, &login_domain_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/userext/duoext/duoprovidergroup-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciProviderGroupMemberIdNotEqual(&login_domain_provider_default, &login_domain_provider_updated),
				),
			},
		},
	})
}

func TestAccAciProviderGroupMember_Update(t *testing.T) {
	var login_domain_provider_default models.ProviderGroupMember
	var login_domain_provider_updated models.ProviderGroupMember
	resourceName := "aci_login_domain_provider.test"
	rName := makeTestVariable(acctest.RandString(5))

	aaaDuoProviderGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciProviderGroupMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccProviderGroupMemberConfig(aaaDuoProviderGroupName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciProviderGroupMemberExists(resourceName, &login_domain_provider_default),
				),
			},
			{
				Config: CreateAccProviderGroupMemberUpdatedAttr(aaaDuoProviderGroupName, rName, "order", "16"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciProviderGroupMemberExists(resourceName, &login_domain_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "order", "16"),
					testAccCheckAciProviderGroupMemberIdEqual(&login_domain_provider_default, &login_domain_provider_updated),
				),
			},
			{
				Config: CreateAccProviderGroupMemberUpdatedAttr(aaaDuoProviderGroupName, rName, "order", "8"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciProviderGroupMemberExists(resourceName, &login_domain_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "order", "8"),
					testAccCheckAciProviderGroupMemberIdEqual(&login_domain_provider_default, &login_domain_provider_updated),
				),
			},

			{
				Config: CreateAccProviderGroupMemberConfig(aaaDuoProviderGroupName, rName),
			},
		},
	})
}

func TestAccAciProviderGroupMember_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	aaaDuoProviderGroupName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciProviderGroupMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccProviderGroupMemberConfig(aaaDuoProviderGroupName, rName),
			},
			{
				Config:      CreateAccProviderGroupMemberWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccProviderGroupMemberUpdatedAttr(aaaDuoProviderGroupName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccProviderGroupMemberUpdatedAttr(aaaDuoProviderGroupName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccProviderGroupMemberUpdatedAttr(aaaDuoProviderGroupName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccProviderGroupMemberUpdatedAttr(aaaDuoProviderGroupName, rName, "order", randomValue),
				ExpectError: regexp.MustCompile(`invalid syntax`),
			},
			{
				Config:      CreateAccProviderGroupMemberUpdatedAttr(aaaDuoProviderGroupName, rName, "order", "-1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccProviderGroupMemberUpdatedAttr(aaaDuoProviderGroupName, rName, "order", "17"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccProviderGroupMemberUpdatedAttr(aaaDuoProviderGroupName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccProviderGroupMemberConfig(aaaDuoProviderGroupName, rName),
			},
		},
	})
}

func TestAccAciProviderGroupMember_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	aaaDuoProviderGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciProviderGroupMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccProviderGroupMemberConfigMultiple(aaaDuoProviderGroupName, rName),
			},
		},
	})
}

func testAccCheckAciProviderGroupMemberExists(name string, login_domain_provider *models.ProviderGroupMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Provider Group Member %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Provider Group Member dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		login_domain_providerFound := models.ProviderGroupMemberFromContainer(cont)
		if login_domain_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Provider Group Member %s not found", rs.Primary.ID)
		}
		*login_domain_provider = *login_domain_providerFound
		return nil
	}
}

func testAccCheckAciProviderGroupMemberDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing login_domain_provider destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_login_domain_provider" {
			cont, err := client.Get(rs.Primary.ID)
			login_domain_provider := models.ProviderGroupMemberFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Provider Group Member %s Still exists", login_domain_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciProviderGroupMemberIdEqual(m1, m2 *models.ProviderGroupMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("login_domain_provider DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciProviderGroupMemberIdNotEqual(m1, m2 *models.ProviderGroupMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("login_domain_provider DNs are equal")
		}
		return nil
	}
}

func CreateProviderGroupMemberWithoutRequired(aaaDuoProviderGroupName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing login_domain_provider creation without ", attrName)
	rBlock := `
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "parent_dn":
		rBlock += `
	resource "aci_login_domain_provider" "test" {
	#	parent_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, aaaDuoProviderGroupName, rName)
}

func CreateAccProviderGroupMemberConfigWithRequiredParams(aaaDuoProviderGroupName, rName string) string {
	fmt.Printf("=== STEP  testing login_domain_provider creation with parent resource name %s and name %s\n", aaaDuoProviderGroupName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}
	`, aaaDuoProviderGroupName, rName)
	return resource
}
func CreateAccProviderGroupMemberConfigUpdatedName(aaaDuoProviderGroupName, rName string) string {
	fmt.Println("=== STEP  testing login_domain_provider creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}
	`, aaaDuoProviderGroupName, rName)
	return resource
}

func CreateAccProviderGroupMemberConfig(aaaDuoProviderGroupName, rName string) string {
	fmt.Println("=== STEP  testing login_domain_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}
	`, aaaDuoProviderGroupName, rName)
	return resource
}

func CreateAccProviderGroupMemberConfigMultiple(aaaDuoProviderGroupName, rName string) string {
	fmt.Println("=== STEP  testing multiple login_domain_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = "%s_${count.index}"
		order = "${count.index}"
		count = 5
	}
	`, aaaDuoProviderGroupName, rName)
	return resource
}

func CreateAccProviderGroupMemberWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing login_domain_provider creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_login_domain_provider" "test" {
		parent_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccProviderGroupMemberConfigWithOptionalValues(aaaDuoProviderGroupName, rName string) string {
	fmt.Println("=== STEP  Basic: testing login_domain_provider creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		parent_dn  = "${aci_duo_provider_group.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_login_domain_provider"
		order = "lowest-available"
		
	}
	`, aaaDuoProviderGroupName, rName)

	return resource
}

func CreateAccProviderGroupMemberRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing login_domain_provider updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_login_domain_provider" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_login_domain_provider"
		order = "1"
		
	}
	`)

	return resource
}

func CreateAccProviderGroupMemberUpdatedAttr(aaaDuoProviderGroupName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing login_domain_provider attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = "%s"
		%s = "%s"
	}
	`, aaaDuoProviderGroupName, rName, attribute, value)
	return resource
}

func CreateAccProviderGroupMemberUpdatedAttrList(aaaDuoProviderGroupName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing login_domain_provider attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = "%s"
		%s = %s
	}
	`, aaaDuoProviderGroupName, rName, attribute, value)
	return resource
}
