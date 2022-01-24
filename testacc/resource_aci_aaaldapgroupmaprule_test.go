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

func TestAccAciLDAPGroupMapRule_Basic(t *testing.T) {
	var ldap_group_map_rule_default models.LDAPGroupMapRule
	var ldap_group_map_rule_updated models.LDAPGroupMapRule
	resourceName := "aci_ldap_group_map_rule.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPGroupMapRuleDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLDAPGroupMapRuleWithoutRequired(rName, "duo", "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLDAPGroupMapRuleWithoutRequired(rName, "duo", "type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLDAPGroupMapRuleConfig(rName, "duo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapRuleExists(resourceName, &ldap_group_map_rule_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "groupdn", ""),
					resource.TestCheckResourceAttr(resourceName, "type", "duo"),
				),
			},
			{
				Config: CreateAccLDAPGroupMapRuleConfigWithOptionalValues(rName, "duo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapRuleExists(resourceName, &ldap_group_map_rule_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_ldap_group_map_rule"),
					resource.TestCheckResourceAttr(resourceName, "groupdn", "groupdn_example"),
					resource.TestCheckResourceAttr(resourceName, "type", "duo"),
					testAccCheckAciLDAPGroupMapRuleIdEqual(&ldap_group_map_rule_default, &ldap_group_map_rule_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccLDAPGroupMapRuleConfigWithRequiredParams(rName, rName),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccLDAPGroupMapRuleConfigUpdatedName(acctest.RandString(65), "duo"),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccLDAPGroupMapRuleRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccLDAPGroupMapRuleConfigWithRequiredParams(rNameUpdated, "duo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapRuleExists(resourceName, &ldap_group_map_rule_updated),
					resource.TestCheckResourceAttr(resourceName, "type", "duo"),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciLDAPGroupMapRuleIdNotEqual(&ldap_group_map_rule_default, &ldap_group_map_rule_updated),
				),
			},
			{
				Config: CreateAccLDAPGroupMapRuleConfig(rName, "duo"),
			},
			{
				Config: CreateAccLDAPGroupMapRuleConfigWithRequiredParams(rName, "ldap"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapRuleExists(resourceName, &ldap_group_map_rule_updated),
					resource.TestCheckResourceAttr(resourceName, "type", "ldap"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciLDAPGroupMapRuleIdNotEqual(&ldap_group_map_rule_default, &ldap_group_map_rule_updated),
				),
			},
		},
	})
}

func TestAccAciLDAPGroupMapRule_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPGroupMapRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLDAPGroupMapRuleConfig(rName, "duo"),
			},

			{
				Config:      CreateAccLDAPGroupMapRuleUpdatedAttr(rName, "duo", "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLDAPGroupMapRuleUpdatedAttr(rName, "duo", "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLDAPGroupMapRuleUpdatedAttr(rName, "duo", "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLDAPGroupMapRuleUpdatedAttr(rName, "duo", "groupdn", acctest.RandString(128)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLDAPGroupMapRuleUpdatedAttr(rName, "duo", randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLDAPGroupMapRuleConfig(rName, "duo"),
			},
		},
	})
}

func TestAccAciLDAPGroupMapRule_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPGroupMapRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLDAPGroupMapRuleConfigMultiple(rName, "duo"),
			},
		},
	})
}

func testAccCheckAciLDAPGroupMapRuleExists(name string, ldap_group_map_rule *models.LDAPGroupMapRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LDAP Group Map Rule %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LDAP Group Map Rule dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ldap_group_map_ruleFound := models.LDAPGroupMapRuleFromContainer(cont)
		if ldap_group_map_ruleFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LDAP Group Map Rule %s not found", rs.Primary.ID)
		}
		*ldap_group_map_rule = *ldap_group_map_ruleFound
		return nil
	}
}

func testAccCheckAciLDAPGroupMapRuleDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing ldap_group_map_rule destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ldap_group_map_rule" {
			cont, err := client.Get(rs.Primary.ID)
			ldap_group_map_rule := models.LDAPGroupMapRuleFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LDAP Group Map Rule %s Still exists", ldap_group_map_rule.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLDAPGroupMapRuleIdEqual(m1, m2 *models.LDAPGroupMapRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("ldap_group_map_rule DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLDAPGroupMapRuleIdNotEqual(m1, m2 *models.LDAPGroupMapRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("ldap_group_map_rule DNs are equal")
		}
		return nil
	}
}

func CreateLDAPGroupMapRuleWithoutRequired(rName, ruleType, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ldap_group_map_rule creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_ldap_group_map_rule" "test" {
	#	name  = "%s"
		type = "%s"
	}
		`
	case "type":
		rBlock += `
	resource "aci_ldap_group_map_rule" "test" {
		name  = "%s"
	#	type = "%s"
	}
	`
	}
	return fmt.Sprintf(rBlock, rName, ruleType)
}

func CreateAccLDAPGroupMapRuleConfigWithRequiredParams(rName, ruleType string) string {
	fmt.Printf("=== STEP  testing ldap_group_map_rule creation with name %s and type %s\n", rName, ruleType)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map_rule" "test" {
		name  = "%s"
		type = "%s"
	}
	`, rName, ruleType)
	return resource
}
func CreateAccLDAPGroupMapRuleConfigUpdatedName(rName, ruleType string) string {
	fmt.Println("=== STEP  testing ldap_group_map_rule creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map_rule" "test" {
		name  = "%s"
		type = "%s"
	}
	`, rName, ruleType)
	return resource
}

func CreateAccLDAPGroupMapRuleConfig(rName, ruleType string) string {
	fmt.Println("=== STEP  testing ldap_group_map_rule creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map_rule" "test" {
		name  = "%s"
		type = "%s"
	}
	`, rName, ruleType)
	return resource
}

func CreateAccLDAPGroupMapRuleConfigMultiple(rName, ruleType string) string {
	fmt.Println("=== STEP  testing multiple ldap_group_map_rule creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map_rule" "test" {
		name  = "%s_${count.index}"
		type = "%s"
		count = 5
	}
	`, rName, ruleType)
	return resource
}

func CreateAccLDAPGroupMapRuleConfigWithOptionalValues(rName, ruleType string) string {
	fmt.Println("=== STEP  Basic: testing ldap_group_map_rule creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map_rule" "test" {
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ldap_group_map_rule"
		groupdn = "groupdn_example"
		type = "%s"
	}
	`, rName, ruleType)

	return resource
}

func CreateAccLDAPGroupMapRuleRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing ldap_group_map_rule updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_ldap_group_map_rule" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ldap_group_map_rule"
		groupdn = "groupdn_example"
	}
	`)

	return resource
}

func CreateAccLDAPGroupMapRuleUpdatedAttr(rName, ruleType, attribute, value string) string {
	fmt.Printf("=== STEP  testing ldap_group_map_rule attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map_rule" "test" {
		name  = "%s"
		type = "%s"
		%s = "%s"
	}
	`, rName, ruleType, attribute, value)
	return resource
}
