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

func TestAccAciSwitchAssociation_Basic(t *testing.T) {
	var leaf_selector_default models.SwitchAssociation
	var leaf_selector_updated models.SwitchAssociation
	resourceName := "aci_leaf_selector.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSwitchAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateSwitchAssociationWithoutRequired(rName, rName, "ALL", "leaf_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSwitchAssociationWithoutRequired(rName, rName, "ALL", "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSwitchAssociationWithoutRequired(rName, rName, "ALL", "switch_association_type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSwitchAssociationConfig(rName, rName, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchAssociationExists(resourceName, &leaf_selector_default),
					resource.TestCheckResourceAttr(resourceName, "leaf_profile_dn", fmt.Sprintf("uni/infra/nprof-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "switch_association_type", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccSwitchAssociationConfigWithOptionalValues(rName, rName, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchAssociationExists(resourceName, &leaf_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "leaf_profile_dn", fmt.Sprintf("uni/infra/nprof-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "switch_association_type", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_leaf_selector"),
					testAccCheckAciSwitchAssociationIdEqual(&leaf_selector_default, &leaf_selector_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSwitchAssociationConfigUpdatedName(rName, acctest.RandString(65), "ALL"),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccSwitchAssociationRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSwitchAssociationConfigWithRequiredParams(rNameUpdated, rName, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchAssociationExists(resourceName, &leaf_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "leaf_profile_dn", fmt.Sprintf("uni/infra/nprof-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "switch_association_type", "ALL"),
					testAccCheckAciSwitchAssociationIdNotEqual(&leaf_selector_default, &leaf_selector_updated),
				),
			},
			{
				Config: CreateAccSwitchAssociationConfig(rName, rName, "ALL"),
			},
			{
				Config: CreateAccSwitchAssociationConfigWithRequiredParams(rName, rNameUpdated, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchAssociationExists(resourceName, &leaf_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "leaf_profile_dn", fmt.Sprintf("uni/infra/nprof-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "switch_association_type", "ALL"),
					testAccCheckAciSwitchAssociationIdNotEqual(&leaf_selector_default, &leaf_selector_updated),
				),
			},
			{
				Config: CreateAccSwitchAssociationConfig(rName, rName, "ALL"),
			},
			{
				Config: CreateAccSwitchAssociationConfigWithRequiredParams(rName, rName, "range"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchAssociationExists(resourceName, &leaf_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "leaf_profile_dn", fmt.Sprintf("uni/infra/nprof-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "switch_association_type", "range"),
					testAccCheckAciSwitchAssociationIdNotEqual(&leaf_selector_default, &leaf_selector_updated),
				),
			},
		},
	})
}

func TestAccAciSwitchAssociation_Update(t *testing.T) {
	var leaf_selector_default models.SwitchAssociation
	resourceName := "aci_leaf_selector.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSwitchAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSwitchAssociationConfig(rName, rName, "ALL_IN_POD"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchAssociationExists(resourceName, &leaf_selector_default),
					resource.TestCheckResourceAttr(resourceName, "leaf_profile_dn", fmt.Sprintf("uni/infra/nprof-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "switch_association_type", "ALL_IN_POD"),
				),
			},
		},
	})
}

func TestAccAciSwitchAssociation_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSwitchAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSwitchAssociationConfig(rName, rName, "ALL"),
			},
			{
				Config:      CreateAccSwitchAssociationConfig(rName, rName, rName),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccSwitchAssociationWithInValidParentDn(rName, "ALL"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccSwitchAssociationUpdatedAttr(rName, rName, "ALL", "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSwitchAssociationUpdatedAttr(rName, rName, "ALL", "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSwitchAssociationUpdatedAttr(rName, rName, "ALL", "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSwitchAssociationUpdatedAttr(rName, rName, "ALL", randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSwitchAssociationConfig(rName, rName, "ALL"),
			},
		},
	})
}

func TestAccAciSwitchAssociation_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSwitchAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSwitchAssociationConfigMultiple(rName, rName, "ALL"),
			},
		},
	})
}

func testAccCheckAciSwitchAssociationExists(name string, leaf_selector *models.SwitchAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Switch Association %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Switch Association dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		leaf_selectorFound := models.SwitchAssociationFromContainer(cont)
		if leaf_selectorFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Switch Association %s not found", rs.Primary.ID)
		}
		*leaf_selector = *leaf_selectorFound
		return nil
	}
}

func testAccCheckAciSwitchAssociationDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing leaf_selector destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_leaf_selector" {
			cont, err := client.Get(rs.Primary.ID)
			leaf_selector := models.SwitchAssociationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Switch Association %s Still exists", leaf_selector.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSwitchAssociationIdEqual(m1, m2 *models.SwitchAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("leaf_selector DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSwitchAssociationIdNotEqual(m1, m2 *models.SwitchAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("leaf_selector DNs are equal")
		}
		return nil
	}
}

func CreateSwitchAssociationWithoutRequired(infraNodePName, rName, switch_association_type, attrName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_selector creation without ", attrName)
	rBlock := `
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "leaf_profile_dn":
		rBlock += `
	resource "aci_leaf_selector" "test" {
	#	leaf_profile_dn  = aci_leaf_profile.test.id
		name  = "%s"	
		switch_association_type  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
	#	name  = "%s"
		switch_association_type  = "%s"
	}
		`
	case "switch_association_type":
		rBlock += `
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = "%s"
	#	switch_association_type  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, infraNodePName, rName, switch_association_type)
}

func CreateAccSwitchAssociationConfigWithRequiredParams(infraNodePName, rName, switch_association_type string) string {
	fmt.Printf("=== STEP  testing leaf_selector creation with parent resource name %s, resource name %s and switch_association_type %s\n", infraNodePName, rName, switch_association_type)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = "%s"
		switch_association_type  = "%s"
	}
	`, infraNodePName, rName, switch_association_type)
	return resource
}
func CreateAccSwitchAssociationConfigUpdatedName(infraNodePName, rName, switch_association_type string) string {
	fmt.Println("=== STEP  testing leaf_selector creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = "%s"
		switch_association_type  = "%s"
	}
	`, infraNodePName, rName, switch_association_type)
	return resource
}

func CreateAccSwitchAssociationConfig(infraNodePName, rName, switch_association_type string) string {
	fmt.Println("=== STEP  testing leaf_selector creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = "%s"
		switch_association_type  = "%s"
	}
	`, infraNodePName, rName, switch_association_type)
	return resource
}

func CreateAccSwitchAssociationConfigMultiple(infraNodePName, rName, switch_association_type string) string {
	fmt.Println("=== STEP  testing multiple leaf_selector creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = "%s_${count.index}"
		switch_association_type  = "%s"
		count = 5
	}
	`, infraNodePName, rName, switch_association_type)
	return resource
}

func CreateAccSwitchAssociationWithInValidParentDn(rName, switch_association_type string) string {
	fmt.Println("=== STEP  Negative Case: testing leaf_selector creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_tenant.test.id
		name  = "%s"
		switch_association_type  = "%s"	
	}
	`, rName, rName, switch_association_type)
	return resource
}

func CreateAccSwitchAssociationConfigWithOptionalValues(infraNodePName, rName, switch_association_type string) string {
	fmt.Println("=== STEP  Basic: testing leaf_selector creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = "${aci_leaf_profile.test.id}"
		name  = "%s"
		switch_association_type  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_leaf_selector"
	}
	`, infraNodePName, rName, switch_association_type)

	return resource
}

func CreateAccSwitchAssociationRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing leaf_selector updation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_leaf_selector" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_leaf_selector"		
	}
	`)

	return resource
}

func CreateAccSwitchAssociationUpdatedAttr(infraNodePName, rName, switch_association_type, attribute, value string) string {
	fmt.Printf("=== STEP  testing leaf_selector attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = "%s"
		switch_association_type  = "%s"
		%s = "%s"
	}
	`, infraNodePName, rName, switch_association_type, attribute, value)
	return resource
}
