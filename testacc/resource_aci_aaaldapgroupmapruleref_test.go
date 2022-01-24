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

func TestAccAciLDAPGroupMapruleref_Basic(t *testing.T) {
	var ldap_group_map_rule_to_group_map_default models.LDAPGroupMapruleref
	var ldap_group_map_rule_to_group_map_updated models.LDAPGroupMapruleref
	resourceName := "aci_ldap_group_map_rule_to_group_map.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPGroupMaprulerefDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateLDAPGroupMaprulerefWithoutRequired(rName, rName, "ldap_group_map_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLDAPGroupMaprulerefWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLDAPGroupMaprulerefConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMaprulerefExists(resourceName, &ldap_group_map_rule_to_group_map_default),
					resource.TestCheckResourceAttr(resourceName, "ldap_group_map_dn", fmt.Sprintf("uni/userext/duoext/ldapgroupmap-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccLDAPGroupMaprulerefConfigWithOptionalValues(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMaprulerefExists(resourceName, &ldap_group_map_rule_to_group_map_updated),
					resource.TestCheckResourceAttr(resourceName, "ldap_group_map_dn", fmt.Sprintf("uni/userext/duoext/ldapgroupmap-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_ldap_group_map_rule_to_group_map"),
					testAccCheckAciLDAPGroupMaprulerefIdEqual(&ldap_group_map_rule_to_group_map_default, &ldap_group_map_rule_to_group_map_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccLDAPGroupMaprulerefConfigUpdatedName(rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccLDAPGroupMaprulerefRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLDAPGroupMaprulerefConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMaprulerefExists(resourceName, &ldap_group_map_rule_to_group_map_updated),
					resource.TestCheckResourceAttr(resourceName, "ldap_group_map_dn", fmt.Sprintf("uni/userext/duoext/ldapgroupmap-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciLDAPGroupMaprulerefIdNotEqual(&ldap_group_map_rule_to_group_map_default, &ldap_group_map_rule_to_group_map_updated),
				),
			},
			{
				Config: CreateAccLDAPGroupMaprulerefConfig(rName, rName),
			},
			{
				Config: CreateAccLDAPGroupMaprulerefConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMaprulerefExists(resourceName, &ldap_group_map_rule_to_group_map_updated),
					resource.TestCheckResourceAttr(resourceName, "ldap_group_map_dn", fmt.Sprintf("uni/userext/duoext/ldapgroupmap-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciLDAPGroupMaprulerefIdNotEqual(&ldap_group_map_rule_to_group_map_default, &ldap_group_map_rule_to_group_map_updated),
				),
			},
		},
	})
}

func TestAccAciLDAPGroupMapruleref_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPGroupMaprulerefDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLDAPGroupMaprulerefConfig(rName, rName),
			},
			{
				Config:      CreateAccLDAPGroupMaprulerefWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLDAPGroupMaprulerefUpdatedAttr(rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLDAPGroupMaprulerefUpdatedAttr(rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLDAPGroupMaprulerefUpdatedAttr(rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLDAPGroupMaprulerefUpdatedAttr(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLDAPGroupMaprulerefConfig(rName, rName),
			},
		},
	})
}

func TestAccAciLDAPGroupMapruleref_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	aaaLdapGroupMapName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPGroupMaprulerefDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLDAPGroupMaprulerefConfigMultiple(aaaLdapGroupMapName, rName),
			},
		},
	})
}

func testAccCheckAciLDAPGroupMaprulerefExists(name string, ldap_group_map_rule_to_group_map *models.LDAPGroupMapruleref) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LDAP Group Mapruleref %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LDAP Group Mapruleref dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ldap_group_map_rule_to_group_mapFound := models.LDAPGroupMaprulerefFromContainer(cont)
		if ldap_group_map_rule_to_group_mapFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LDAP Group Mapruleref %s not found", rs.Primary.ID)
		}
		*ldap_group_map_rule_to_group_map = *ldap_group_map_rule_to_group_mapFound
		return nil
	}
}

func testAccCheckAciLDAPGroupMaprulerefDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing ldap_group_map_rule_to_group_map destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ldap_group_map_rule_to_group_map" {
			cont, err := client.Get(rs.Primary.ID)
			ldap_group_map_rule_to_group_map := models.LDAPGroupMaprulerefFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LDAP Group Mapruleref %s Still exists", ldap_group_map_rule_to_group_map.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLDAPGroupMaprulerefIdEqual(m1, m2 *models.LDAPGroupMapruleref) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("ldap_group_map_rule_to_group_map DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLDAPGroupMaprulerefIdNotEqual(m1, m2 *models.LDAPGroupMapruleref) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("ldap_group_map_rule_to_group_map DNs are equal")
		}
		return nil
	}
}

func CreateLDAPGroupMaprulerefWithoutRequired(aaaLdapGroupMapName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ldap_group_map_rule_to_group_map creation without ", attrName)
	rBlock := `
	
	resource "aci_ldap_group_map" "test" {
		name 		= "%s"
		type = "duo"
	}
	
	`
	switch attrName {
	case "ldap_group_map_dn":
		rBlock += `
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
	#	ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, aaaLdapGroupMapName, rName)
}

func CreateAccLDAPGroupMaprulerefConfigWithRequiredParams(aaaLdapGroupMapName, rName string) string {
	fmt.Println("=== STEP  testing ldap_group_map_rule_to_group_map creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
		name 		= "%s"
		type = "duo"
	}
	
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = "%s"
	}
	`, aaaLdapGroupMapName, rName)
	return resource
}
func CreateAccLDAPGroupMaprulerefConfigUpdatedName(aaaLdapGroupMapName, rName string) string {
	fmt.Println("=== STEP  testing ldap_group_map_rule_to_group_map creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
		name 		= "%s"
		type = "duo"
	}
	
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = "%s"
	}
	`, aaaLdapGroupMapName, rName)
	return resource
}

func CreateAccLDAPGroupMaprulerefConfig(aaaLdapGroupMapName, rName string) string {
	fmt.Println("=== STEP  testing ldap_group_map_rule_to_group_map creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
		name 		= "%s"
		type = "duo"
	}
	
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = "%s"
	}
	`, aaaLdapGroupMapName, rName)
	return resource
}

func CreateAccLDAPGroupMaprulerefConfigMultiple(aaaLdapGroupMapName, rName string) string {
	fmt.Println("=== STEP  testing multiple ldap_group_map_rule_to_group_map creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
		name 		= "%s"
		type = "duo"
	}
	
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, aaaLdapGroupMapName, rName)
	return resource
}

func CreateAccLDAPGroupMaprulerefWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing ldap_group_map_rule_to_group_map creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccLDAPGroupMaprulerefConfigWithOptionalValues(aaaLdapGroupMapName, rName string) string {
	fmt.Println("=== STEP  Basic: testing ldap_group_map_rule_to_group_map creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
		name 		= "%s"
		type = "duo"
	}
	
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = "${aci_ldap_group_map.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ldap_group_map_rule_to_group_map"
		
	}
	`, aaaLdapGroupMapName, rName)

	return resource
}

func CreateAccLDAPGroupMaprulerefRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing ldap_group_map_rule_to_group_map updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ldap_group_map_rule_to_group_map"
		
	}
	`)

	return resource
}

func CreateAccLDAPGroupMaprulerefUpdatedAttr(aaaLdapGroupMapName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing ldap_group_map_rule_to_group_map attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
		name 		= "%s"
		type = "duo"
	}
	
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = "%s"
		%s = "%s"
	}
	`, aaaLdapGroupMapName, rName, attribute, value)
	return resource
}

func CreateAccLDAPGroupMaprulerefUpdatedAttrList(aaaLdapGroupMapName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing ldap_group_map_rule_to_group_map attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = "%s"
		%s = %s
	}
	`, aaaLdapGroupMapName, rName, attribute, value)
	return resource
}
