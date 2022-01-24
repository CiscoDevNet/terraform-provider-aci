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

func TestAccAciSwitchSpineAssociation_Basic(t *testing.T) {
	var spine_switch_association_default models.SwitchSpineAssociation
	var spine_switch_association_updated models.SwitchSpineAssociation
	resourceName := "aci_spine_switch_association.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	infraSpinePName := makeTestVariable(acctest.RandString(5))
	spine_switch_association_type := "ALL"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSwitchSpineAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateSwitchSpineAssociationWithoutRequired(infraSpinePName, rName, spine_switch_association_type, "spine_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSwitchSpineAssociationWithoutRequired(infraSpinePName, rName, spine_switch_association_type, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSwitchSpineAssociationWithoutRequired(infraSpinePName, rName, spine_switch_association_type, "spine_switch_association_type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSwitchSpineAssociationConfig(infraSpinePName, rName, spine_switch_association_type),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchSpineAssociationExists(resourceName, &spine_switch_association_default),
					resource.TestCheckResourceAttr(resourceName, "spine_profile_dn", fmt.Sprintf("uni/infra/spprof-%s", infraSpinePName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "spine_switch_association_type", spine_switch_association_type),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccSwitchSpineAssociationConfigWithOptionalValues(infraSpinePName, rName, spine_switch_association_type),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchSpineAssociationExists(resourceName, &spine_switch_association_updated),
					resource.TestCheckResourceAttr(resourceName, "spine_profile_dn", fmt.Sprintf("uni/infra/spprof-%s", infraSpinePName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "spine_switch_association_type", spine_switch_association_type),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_spine_switch_association"),

					testAccCheckAciSwitchSpineAssociationIdEqual(&spine_switch_association_default, &spine_switch_association_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSwitchSpineAssociationConfigUpdatedName(infraSpinePName, acctest.RandString(65), spine_switch_association_type),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccSwitchSpineAssociationRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSwitchSpineAssociationConfigWithRequiredParams(rNameUpdated, rName, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchSpineAssociationExists(resourceName, &spine_switch_association_updated),
					resource.TestCheckResourceAttr(resourceName, "spine_profile_dn", fmt.Sprintf("uni/infra/spprof-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "spine_switch_association_type", spine_switch_association_type),
					testAccCheckAciSwitchSpineAssociationIdNotEqual(&spine_switch_association_default, &spine_switch_association_updated),
				),
			},
			{
				Config: CreateAccSwitchSpineAssociationConfig(infraSpinePName, rName, spine_switch_association_type),
			},
			{
				Config: CreateAccSwitchSpineAssociationConfigWithRequiredParams(rName, rNameUpdated, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchSpineAssociationExists(resourceName, &spine_switch_association_updated),
					resource.TestCheckResourceAttr(resourceName, "spine_profile_dn", fmt.Sprintf("uni/infra/spprof-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "spine_switch_association_type", "ALL"),
					testAccCheckAciSwitchSpineAssociationIdNotEqual(&spine_switch_association_default, &spine_switch_association_updated),
				),
			},
			{
				Config: CreateAccSwitchSpineAssociationConfigWithRequiredParams(rName, rName, "range"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchSpineAssociationExists(resourceName, &spine_switch_association_updated),
					resource.TestCheckResourceAttr(resourceName, "spine_profile_dn", fmt.Sprintf("uni/infra/spprof-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "spine_switch_association_type", "range"),
					testAccCheckAciSwitchSpineAssociationIdNotEqual(&spine_switch_association_default, &spine_switch_association_updated),
				),
			},
			{
				Config: CreateAccSwitchSpineAssociationConfigWithRequiredParams(rName, rName, "ALL_IN_POD"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchSpineAssociationExists(resourceName, &spine_switch_association_updated),
					resource.TestCheckResourceAttr(resourceName, "spine_profile_dn", fmt.Sprintf("uni/infra/spprof-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "spine_switch_association_type", "ALL_IN_POD"),
					testAccCheckAciSwitchSpineAssociationIdNotEqual(&spine_switch_association_default, &spine_switch_association_updated),
				),
			},
		},
	})
}

func TestAccAciSwitchSpineAssociation_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	spine_switch_association_type := "ALL"
	infraSpinePName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSwitchSpineAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSwitchSpineAssociationConfig(infraSpinePName, rName, spine_switch_association_type),
			},
			{
				Config:      CreateAccSwitchSpineAssociationWithInValidParentDn(rName, spine_switch_association_type),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccSwitchSpineAssociationUpdatedAttr(infraSpinePName, rName, spine_switch_association_type, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSwitchSpineAssociationUpdatedAttr(infraSpinePName, rName, spine_switch_association_type, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSwitchSpineAssociationUpdatedAttr(infraSpinePName, rName, spine_switch_association_type, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccSwitchSpineAssociationUpdatedAttr(infraSpinePName, rName, spine_switch_association_type, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSwitchSpineAssociationConfig(infraSpinePName, rName, spine_switch_association_type),
			},
		},
	})
}

func TestAccAciSwitchSpineAssociation_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	spine_switch_association_type := "ALL"
	infraSpinePName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSwitchSpineAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSwitchSpineAssociationConfigMultiple(infraSpinePName, rName, spine_switch_association_type),
			},
		},
	})
}

func testAccCheckAciSwitchSpineAssociationExists(name string, spine_switch_association *models.SwitchSpineAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Spine Switch Association %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spine Switch Association dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		spine_switch_associationFound := models.SwitchSpineAssociationFromContainer(cont)
		if spine_switch_associationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spine Switch Association %s not found", rs.Primary.ID)
		}
		*spine_switch_association = *spine_switch_associationFound
		return nil
	}
}

func testAccCheckAciSwitchSpineAssociationDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing spine_switch_association destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_spine_switch_association" {
			cont, err := client.Get(rs.Primary.ID)
			spine_switch_association := models.SwitchSpineAssociationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spine Switch Association %s Still exists", spine_switch_association.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSwitchSpineAssociationIdEqual(m1, m2 *models.SwitchSpineAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("spine_switch_association DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSwitchSpineAssociationIdNotEqual(m1, m2 *models.SwitchSpineAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("spine_switch_association DNs are equal")
		}
		return nil
	}
}

func CreateSwitchSpineAssociationWithoutRequired(infraSpinePName, rName, spine_switch_association_type, attrName string) string {
	fmt.Println("=== STEP  Basic: testing spine_switch_association creation without ", attrName)
	rBlock := `
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "spine_profile_dn":
		rBlock += `
	resource "aci_spine_switch_association" "test" {
	#	spine_profile_dn  = aci_spine_profile.test.id
		name  = "%s"	
		spine_switch_association_type  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
	#	name  = "%s"
		spine_switch_association_type  = "%s"
	}
		`
	case "spine_switch_association_type":
		rBlock += `
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = "%s"
	#	spine_switch_association_type  = "%s"
	}
	`
	}
	return fmt.Sprintf(rBlock, infraSpinePName, rName, spine_switch_association_type)
}

func CreateAccSwitchSpineAssociationConfigWithRequiredParams(infraSpinePName, rName, spine_switch_association_type string) string {
	fmt.Println("=== STEP  testing spine_switch_association creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = "%s"
		spine_switch_association_type  = "%s"
	}
	`, infraSpinePName, rName, spine_switch_association_type)
	return resource
}
func CreateAccSwitchSpineAssociationConfigUpdatedName(infraSpinePName, rName, spine_switch_association_type string) string {
	fmt.Println("=== STEP  testing spine_switch_association creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = "%s"
		spine_switch_association_type  = "%s"
	}
	`, infraSpinePName, rName, spine_switch_association_type)
	return resource
}

func CreateAccSwitchSpineAssociationConfig(infraSpinePName, rName, spine_switch_association_type string) string {
	fmt.Println("=== STEP  testing spine_switch_association creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = "%s"
		spine_switch_association_type  = "%s"
	}
	`, infraSpinePName, rName, spine_switch_association_type)
	return resource
}

func CreateAccSwitchSpineAssociationConfigMultiple(infraSpinePName, rName, spine_switch_association_type string) string {
	fmt.Println("=== STEP  testing multiple spine_switch_association creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = "%s_${count.index}"
		spine_switch_association_type  = "%s"
		count = 5
	}
	`, infraSpinePName, rName, spine_switch_association_type)
	return resource
}

func CreateAccSwitchSpineAssociationWithInValidParentDn(rName, spine_switch_association_type string) string {
	fmt.Println("=== STEP  Negative Case: testing spine_switch_association creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_tenant.test.id
		name  = "%s"
		spine_switch_association_type  = "%s"	
	}
	`, rName, rName, spine_switch_association_type)
	return resource
}

func CreateAccSwitchSpineAssociationConfigWithOptionalValues(infraSpinePName, rName, spine_switch_association_type string) string {
	fmt.Println("=== STEP  Basic: testing spine_switch_association creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = "${aci_spine_profile.test.id}"
		name  = "%s"
		spine_switch_association_type  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_spine_switch_association"
		
	}
	`, infraSpinePName, rName, spine_switch_association_type)

	return resource
}

func CreateAccSwitchSpineAssociationRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing spine_switch_association updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_spine_switch_association" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_spine_switch_association"
		
	}
	`)

	return resource
}

func CreateAccSwitchSpineAssociationUpdatedAttr(infraSpinePName, rName, spine_switch_association_type, attribute, value string) string {
	fmt.Printf("=== STEP  testing spine_switch_association attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = "%s"
		spine_switch_association_type  = "%s"
		%s = "%s"
	}
	`, infraSpinePName, rName, spine_switch_association_type, attribute, value)
	return resource
}
