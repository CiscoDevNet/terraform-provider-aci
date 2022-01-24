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

func TestAccAciDuoProviderGroup_Basic(t *testing.T) {
	var duo_provider_group_default models.DuoProviderGroup
	var duo_provider_group_updated models.DuoProviderGroup
	resourceName := "aci_duo_provider_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDuoProviderGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateDuoProviderGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccDuoProviderGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "auth_choice", "CiscoAVPair"),
					resource.TestCheckResourceAttr(resourceName, "ldap_group_map_ref", ""),
					resource.TestCheckResourceAttr(resourceName, "provider_type", "radius"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.0", "auto"),
				),
			},
			{
				Config: CreateAccDuoProviderGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_duo_provider_group"),
					resource.TestCheckResourceAttr(resourceName, "auth_choice", "LdapGroupMap"),
					resource.TestCheckResourceAttr(resourceName, "ldap_group_map_ref", "DuoEmpGroupMap"),
					resource.TestCheckResourceAttr(resourceName, "provider_type", "ldap"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.0", "passcode"),
					testAccCheckAciDuoProviderGroupIdEqual(&duo_provider_group_default, &duo_provider_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccDuoProviderGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccDuoProviderGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccDuoProviderGroupConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciDuoProviderGroupIdNotEqual(&duo_provider_group_default, &duo_provider_group_updated),
				),
			},
		},
	})
}

func TestAccAciDuoProviderGroup_Update(t *testing.T) {
	var duo_provider_group_default models.DuoProviderGroup
	var duo_provider_group_updated models.DuoProviderGroup
	resourceName := "aci_duo_provider_group.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDuoProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDuoProviderGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_default),
				),
			},

			{

				Config: CreateAccDuoProviderGroupUpdatedAttrList(rName, "sec_fac_auth_methods", StringListtoString([]string{"auto", "passcode"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_updated),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.0", "auto"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.1", "passcode"),
				),
			},
			{

				Config: CreateAccDuoProviderGroupUpdatedAttrList(rName, "sec_fac_auth_methods", StringListtoString([]string{"passcode"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_updated),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.0", "passcode"),
				),
			},
			{

				Config: CreateAccDuoProviderGroupUpdatedAttrList(rName, "sec_fac_auth_methods", StringListtoString([]string{"auto", "passcode", "phone"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_updated),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.0", "auto"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.1", "passcode"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.2", "phone"),
				),
			},
			{

				Config: CreateAccDuoProviderGroupUpdatedAttrList(rName, "sec_fac_auth_methods", StringListtoString([]string{"passcode", "phone"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_updated),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.0", "passcode"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.1", "phone"),
				),
			},
			{

				Config: CreateAccDuoProviderGroupUpdatedAttrList(rName, "sec_fac_auth_methods", StringListtoString([]string{"phone"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_updated),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.0", "phone"),
				),
			},
			{

				Config: CreateAccDuoProviderGroupUpdatedAttrList(rName, "sec_fac_auth_methods", StringListtoString([]string{"auto", "passcode", "phone", "push"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_updated),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.0", "auto"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.1", "passcode"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.2", "phone"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.3", "push"),
				),
			},
			{

				Config: CreateAccDuoProviderGroupUpdatedAttrList(rName, "sec_fac_auth_methods", StringListtoString([]string{"passcode", "phone", "push"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_updated),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.0", "passcode"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.1", "phone"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.2", "push"),
				),
			},
			{

				Config: CreateAccDuoProviderGroupUpdatedAttrList(rName, "sec_fac_auth_methods", StringListtoString([]string{"phone", "push"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_updated),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.0", "phone"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.1", "push"),
				),
			},
			{

				Config: CreateAccDuoProviderGroupUpdatedAttrList(rName, "sec_fac_auth_methods", StringListtoString([]string{"push"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_updated),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.0", "push"),
				),
			},
			{
				Config: CreateAccDuoProviderGroupUpdatedAttrList(rName, "sec_fac_auth_methods", StringListtoString([]string{"push", "phone", "passcode", "auto"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists(resourceName, &duo_provider_group_updated),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.0", "push"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.1", "phone"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.2", "passcode"),
					resource.TestCheckResourceAttr(resourceName, "sec_fac_auth_methods.3", "auto"),
				),
			},
			{
				Config: CreateAccDuoProviderGroupConfig(rName),
			},
		},
	})
}

func TestAccAciDuoProviderGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDuoProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDuoProviderGroupConfig(rName),
			},

			{
				Config:      CreateAccDuoProviderGroupUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDuoProviderGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDuoProviderGroupUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccDuoProviderGroupUpdatedAttr(rName, "auth_choice", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccDuoProviderGroupUpdatedAttr(rName, "ldap_group_map_ref", acctest.RandString(513)),
				ExpectError: regexp.MustCompile(`failed validation for value ''`),
			},

			{
				Config:      CreateAccDuoProviderGroupUpdatedAttr(rName, "provider_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccDuoProviderGroupUpdatedAttrList(rName, "sec_fac_auth_methods", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccDuoProviderGroupUpdatedAttrList(rName, "sec_fac_auth_methods", StringListtoString([]string{"auto", "auto"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccDuoProviderGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccDuoProviderGroupConfig(rName),
			},
		},
	})
}

func TestAccAciDuoProviderGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDuoProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDuoProviderGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciDuoProviderGroupExists(name string, duo_provider_group *models.DuoProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Duo Provider Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Duo Provider Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		duo_provider_groupFound := models.DuoProviderGroupFromContainer(cont)
		if duo_provider_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Duo Provider Group %s not found", rs.Primary.ID)
		}
		*duo_provider_group = *duo_provider_groupFound
		return nil
	}
}

func testAccCheckAciDuoProviderGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing duo_provider_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_duo_provider_group" {
			cont, err := client.Get(rs.Primary.ID)
			duo_provider_group := models.DuoProviderGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Duo Provider Group %s Still exists", duo_provider_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciDuoProviderGroupIdEqual(m1, m2 *models.DuoProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("duo_provider_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciDuoProviderGroupIdNotEqual(m1, m2 *models.DuoProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("duo_provider_group DNs are equal")
		}
		return nil
	}
}

func CreateDuoProviderGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing duo_provider_group creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_duo_provider_group" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccDuoProviderGroupConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing duo_provider_group creation with name =", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccDuoProviderGroupConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing duo_provider_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccDuoProviderGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing duo_provider_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccDuoProviderGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple duo_provider_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccDuoProviderGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing duo_provider_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_duo_provider_group"
		auth_choice = "LdapGroupMap"
		ldap_group_map_ref = "DuoEmpGroupMap"
		provider_type = "ldap"
		sec_fac_auth_methods = ["passcode"]
		
	}
	`, rName)

	return resource
}

func CreateAccDuoProviderGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing duo_provider_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_duo_provider_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_duo_provider_group"
		auth_choice = "LdapGroupMap"
		ldap_group_map_ref = ""
		provider_type = "ldap"
		sec_fac_auth_methods = ["passcode"]
		
	}
	`)

	return resource
}

func CreateAccDuoProviderGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing duo_provider_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccDuoProviderGroupUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing duo_provider_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
