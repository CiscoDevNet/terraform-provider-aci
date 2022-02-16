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

func TestAccAciLDAPGroupMap_Basic(t *testing.T) {
	var ldap_group_map_default models.LDAPGroupMap
	var ldap_group_map_updated models.LDAPGroupMap
	resourceName := "aci_ldap_group_map.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPGroupMapDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLDAPGroupMapWithoutRequired(rName, "duo", "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLDAPGroupMapWithoutRequired(rName, "duo", "type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLDAPGroupMapConfig(rName, "duo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapExists(resourceName, &ldap_group_map_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "type", "duo"),
				),
			},
			{
				Config: CreateAccLDAPGroupMapConfigWithOptionalValues(rName, "duo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapExists(resourceName, &ldap_group_map_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_ldap_group_map"),
					resource.TestCheckResourceAttr(resourceName, "type", "duo"),
					testAccCheckAciLDAPGroupMapIdEqual(&ldap_group_map_default, &ldap_group_map_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"type"},
			},
			{
				Config:      CreateAccLDAPGroupMapConfigUpdatedName(acctest.RandString(65), "duo"),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccLDAPGroupMapRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLDAPGroupMapConfigWithRequiredParams(rName, rName),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config: CreateAccLDAPGroupMapConfigWithRequiredParams(rNameUpdated, "duo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapExists(resourceName, &ldap_group_map_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "type", "duo"),
					testAccCheckAciLDAPGroupMapIdNotEqual(&ldap_group_map_default, &ldap_group_map_updated),
				),
			},
			{
				Config: CreateAccLDAPGroupMapConfig(rName, "duo"),
			},
			{
				Config: CreateAccLDAPGroupMapConfigWithRequiredParams(rName, "ldap"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapExists(resourceName, &ldap_group_map_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "ldap"),
					testAccCheckAciLDAPGroupMapIdNotEqual(&ldap_group_map_default, &ldap_group_map_updated),
				),
			},
		},
	})
}

func TestAccAciLDAPGroupMap_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPGroupMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLDAPGroupMapConfig(rName, "duo"),
			},

			{
				Config:      CreateAccLDAPGroupMapUpdatedAttr(rName, "duo", "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLDAPGroupMapUpdatedAttr(rName, "duo", "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLDAPGroupMapUpdatedAttr(rName, "duo", "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLDAPGroupMapUpdatedAttr(rName, "duo", randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLDAPGroupMapConfig(rName, "duo"),
			},
		},
	})
}

func TestAccAciLDAPGroupMap_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPGroupMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLDAPGroupMapConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciLDAPGroupMapExists(name string, ldap_group_map *models.LDAPGroupMap) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LDAP Group Map %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LDAP Group Map dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ldap_group_mapFound := models.LDAPGroupMapFromContainer(cont)
		if ldap_group_mapFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LDAP Group Map %s not found", rs.Primary.ID)
		}
		*ldap_group_map = *ldap_group_mapFound
		return nil
	}
}

func testAccCheckAciLDAPGroupMapDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing ldap_group_map destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ldap_group_map" {
			cont, err := client.Get(rs.Primary.ID)
			ldap_group_map := models.LDAPGroupMapFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LDAP Group Map %s Still exists", ldap_group_map.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLDAPGroupMapIdEqual(m1, m2 *models.LDAPGroupMap) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("ldap_group_map DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLDAPGroupMapIdNotEqual(m1, m2 *models.LDAPGroupMap) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("ldap_group_map DNs are equal")
		}
		return nil
	}
}

func CreateLDAPGroupMapWithoutRequired(rName, groupType, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ldap_group_map creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_ldap_group_map" "test" {
	#	name  = "%s"
		type = "%s"
	}
		`
	case "type":
		rBlock += `
	resource "aci_ldap_group_map" "test" {
		name  = "%s"
	#	type = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, groupType)
}

func CreateAccLDAPGroupMapConfigWithRequiredParams(rName, groupType string) string {
	fmt.Printf("=== STEP  testing ldap_group_map creation with type %s and name %s\n", groupType, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
	
		name  = "%s"
		type = "%s"
	}
	`, rName, groupType)
	return resource
}
func CreateAccLDAPGroupMapConfigUpdatedName(rName, groupType string) string {
	fmt.Println("=== STEP  testing ldap_group_map creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
	
		name  = "%s"
		type = "%s"
	}
	`, rName, groupType)
	return resource
}

func CreateAccLDAPGroupMapConfig(rName, groupType string) string {
	fmt.Println("=== STEP  testing ldap_group_map creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
	
		name  = "%s"
		type = "%s"
	}
	`, rName, groupType)
	return resource
}

func CreateAccLDAPGroupMapConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple ldap_group_map creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
	
		name  = "%s_${count.index}"
		type = "duo"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccLDAPGroupMapConfigWithOptionalValues(rName, groupType string) string {
	fmt.Println("=== STEP  Basic: testing ldap_group_map creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ldap_group_map"
		type = "%s"
	}
	`, rName, groupType)

	return resource
}

func CreateAccLDAPGroupMapRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing ldap_group_map updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_ldap_group_map" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ldap_group_map"
		
	}
	`)

	return resource
}

func CreateAccLDAPGroupMapUpdatedAttr(rName, groupType, attribute, value string) string {
	fmt.Printf("=== STEP  testing ldap_group_map attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
	
		name  = "%s"
		type = "%s"
		%s = "%s"
	}
	`, rName, groupType, attribute, value)
	return resource
}
