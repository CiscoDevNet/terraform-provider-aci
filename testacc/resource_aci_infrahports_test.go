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

func TestAccAciAccessPortSelector_Basic(t *testing.T) {
	var access_port_selector_default models.AccessPortSelector
	var access_port_selector_updated models.AccessPortSelector
	resourceName := "aci_access_port_selector.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	infraAccPortPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessPortSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccessPortSelectorWithoutRequired(infraAccPortPName, rName, "ALL", "leaf_interface_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccessPortSelectorWithoutRequired(infraAccPortPName, rName, "ALL", "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAccessPortSelectorConfig(infraAccPortPName, rName, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortSelectorExists(resourceName, &access_port_selector_default),
					resource.TestCheckResourceAttr(resourceName, "leaf_interface_profile_dn", fmt.Sprintf("uni/infra/accportprof-%s", infraAccPortPName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "access_port_selector_type", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccAccessPortSelectorConfigWithOptionalValues(infraAccPortPName, rName, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortSelectorExists(resourceName, &access_port_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "leaf_interface_profile_dn", fmt.Sprintf("uni/infra/accportprof-%s", infraAccPortPName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "access_port_selector_type", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_access_port_selector"),

					testAccCheckAciAccessPortSelectorIdEqual(&access_port_selector_default, &access_port_selector_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccAccessPortSelectorConfigUpdatedName(infraAccPortPName, acctest.RandString(65), "ALL"),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccAccessPortSelectorConfig(infraAccPortPName, rName, acctest.RandString(5)),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccAccessPortSelectorRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAccessPortSelectorConfigWithRequiredParams(rNameUpdated, rName, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortSelectorExists(resourceName, &access_port_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "leaf_interface_profile_dn", fmt.Sprintf("uni/infra/accportprof-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "access_port_selector_type", "ALL"),
					testAccCheckAciAccessPortSelectorIdNotEqual(&access_port_selector_default, &access_port_selector_updated),
				),
			},
			{
				Config: CreateAccAccessPortSelectorConfig(infraAccPortPName, rName, "ALL"),
			},
			{
				Config: CreateAccAccessPortSelectorConfigWithRequiredParams(infraAccPortPName, rNameUpdated, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortSelectorExists(resourceName, &access_port_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "leaf_interface_profile_dn", fmt.Sprintf("uni/infra/accportprof-%s", infraAccPortPName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciAccessPortSelectorIdNotEqual(&access_port_selector_default, &access_port_selector_updated),
				),
			},
			{
				Config: CreateAccAccessPortSelectorConfigWithRequiredParams(infraAccPortPName, rName, "range"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortSelectorExists(resourceName, &access_port_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "leaf_interface_profile_dn", fmt.Sprintf("uni/infra/accportprof-%s", infraAccPortPName)),
					resource.TestCheckResourceAttr(resourceName, "access_port_selector_type", "range"),
					testAccCheckAciAccessPortSelectorIdNotEqual(&access_port_selector_default, &access_port_selector_updated),
				),
			},
		},
	})
}

func TestAccAciAccessPortSelector_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	infraAccPortPName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessPortSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAccessPortSelectorConfig(infraAccPortPName, rName, "ALL"),
			},
			{
				Config:      CreateAccAccessPortSelectorWithInValidParentDn(rName, "ALL"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccAccessPortSelectorUpdatedAttr(infraAccPortPName, rName, "ALL", "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAccessPortSelectorUpdatedAttr(infraAccPortPName, rName, "ALL", "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAccessPortSelectorUpdatedAttr(infraAccPortPName, rName, "ALL", "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccAccessPortSelectorUpdatedAttr(infraAccPortPName, rName, "ALL", randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccAccessPortSelectorConfig(infraAccPortPName, rName, "ALL"),
			},
		},
	})
}

func TestAccAciAccessPortSelector_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	infraAccPortPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessPortSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAccessPortSelectorConfigMultiple(infraAccPortPName, rName, "ALL"),
			},
		},
	})
}

func testAccCheckAciAccessPortSelectorExists(name string, access_port_selector *models.AccessPortSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Access Port Selector %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Access Port Selector dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		access_port_selectorFound := models.AccessPortSelectorFromContainer(cont)
		if access_port_selectorFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Access Port Selector %s not found", rs.Primary.ID)
		}
		*access_port_selector = *access_port_selectorFound
		return nil
	}
}

func testAccCheckAciAccessPortSelectorDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing access_port_selector destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_access_port_selector" {
			cont, err := client.Get(rs.Primary.ID)
			access_port_selector := models.AccessPortSelectorFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Access Port Selector %s Still exists", access_port_selector.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciAccessPortSelectorIdEqual(m1, m2 *models.AccessPortSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("access_port_selector DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciAccessPortSelectorIdNotEqual(m1, m2 *models.AccessPortSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("access_port_selector DNs are equal")
		}
		return nil
	}
}

func CreateAccessPortSelectorWithoutRequired(infraAccPortPName, rName, access_port_selector_type, attrName string) string {
	fmt.Println("=== STEP  Basic: testing access_port_selector creation without ", attrName)
	rBlock := `
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "leaf_interface_profile_dn":
		rBlock += `
	resource "aci_access_port_selector" "test" {
	#	leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "%s"	
		access_port_selector_type  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
	#	name  = "%s"
		access_port_selector_type  = "%s"
	}
		`
	case "access_port_selector_type":
		rBlock += `
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "%s"
	#	access_port_selector_type  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, infraAccPortPName, rName, access_port_selector_type)
}

func CreateAccAccessPortSelectorConfigWithRequiredParams(infraAccPortPName, rName, access_port_selector_type string) string {
	fmt.Println("=== STEP  testing access_port_selector creation with Updated required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "%s"
		access_port_selector_type  = "%s"
	}
	`, infraAccPortPName, rName, access_port_selector_type)
	return resource
}
func CreateAccAccessPortSelectorConfigUpdatedName(infraAccPortPName, rName, access_port_selector_type string) string {
	fmt.Println("=== STEP  testing access_port_selector creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "%s"
		access_port_selector_type  = "%s"
	}
	`, infraAccPortPName, rName, access_port_selector_type)
	return resource
}

func CreateAccAccessPortSelectorConfigUpdatedType(infraAccPortPName, rName, access_port_selector_type string) string {
	fmt.Println("=== STEP  testing access_port_selector creation with invalid access_port_selector_type = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "%s"
		access_port_selector_type  = "%s"
	}
	`, infraAccPortPName, rName, access_port_selector_type)
	return resource
}

func CreateAccAccessPortSelectorConfig(infraAccPortPName, rName, access_port_selector_type string) string {
	fmt.Printf("=== STEP  testing access_port_selector creation with name %s and access port selector type %s required arguments only\n", rName, access_port_selector_type)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "%s"
		access_port_selector_type  = "%s"
	}
	`, infraAccPortPName, rName, access_port_selector_type)
	return resource
}

func CreateAccAccessPortSelectorConfigMultiple(infraAccPortPName, rName, access_port_selector_type string) string {
	fmt.Println("=== STEP  testing multiple access_port_selector creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "%s_${count.index}"
		access_port_selector_type  = "%s"
		count = 5
	}
	`, infraAccPortPName, rName, access_port_selector_type)
	return resource
}

func CreateAccAccessPortSelectorWithInValidParentDn(rName, access_port_selector_type string) string {
	fmt.Println("=== STEP  Negative Case: testing access_port_selector creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_tenant.test.id
		name  = "%s"
		access_port_selector_type  = "%s"	
	}
	`, rName, rName, access_port_selector_type)
	return resource
}

func CreateAccAccessPortSelectorConfigWithOptionalValues(infraAccPortPName, rName, access_port_selector_type string) string {
	fmt.Println("=== STEP  Basic: testing access_port_selector creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = "${aci_leaf_interface_profile.test.id}"
		name  = "%s"
		access_port_selector_type  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_access_port_selector"
		
	}
	`, infraAccPortPName, rName, access_port_selector_type)

	return resource
}

func CreateAccAccessPortSelectorRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing access_port_selector updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_access_port_selector" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_access_port_selector"
	}
	`)

	return resource
}

func CreateAccAccessPortSelectorUpdatedAttr(infraAccPortPName, rName, access_port_selector_type, attribute, value string) string {
	fmt.Printf("=== STEP  testing access_port_selector attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "%s"
		access_port_selector_type  = "%s"
		%s = "%s"
	}
	`, infraAccPortPName, rName, access_port_selector_type, attribute, value)
	return resource
}
