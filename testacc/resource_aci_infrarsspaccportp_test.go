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

func TestAccAciSpinePortSelector_Basic(t *testing.T) {
	var spine_port_selector_default models.InterfaceProfile
	var spine_port_selector_updated models.InterfaceProfile
	resourceName := "aci_spine_port_selector.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	tDn := makeTestVariable(acctest.RandString(5))
	tDnUpdated := makeTestVariable(acctest.RandString(5))
	infraSpinePName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpinePortSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateSpinePortSelectorWithoutRequired(infraSpinePName, tDn, "spine_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSpinePortSelectorWithoutRequired(infraSpinePName, tDn, "tdn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSpinePortSelectorConfig(infraSpinePName, tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpinePortSelectorExists(resourceName, &spine_port_selector_default),
					resource.TestCheckResourceAttr(resourceName, "spine_profile_dn", fmt.Sprintf("uni/infra/spprof-%s", infraSpinePName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/infra/spaccportprof-%s", tDn)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
				),
			},
			{
				Config: CreateAccSpinePortSelectorConfigWithOptionalValues(infraSpinePName, tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpinePortSelectorExists(resourceName, &spine_port_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "spine_profile_dn", fmt.Sprintf("uni/infra/spprof-%s", infraSpinePName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/infra/spaccportprof-%s", tDn)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),

					testAccCheckAciSpinePortSelectorIdEqual(&spine_port_selector_default, &spine_port_selector_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      CreateAccSpinePortSelectorRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSpinePortSelectorConfigWithRequiredParams(rNameUpdated, tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpinePortSelectorExists(resourceName, &spine_port_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "spine_profile_dn", fmt.Sprintf("uni/infra/spprof-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/infra/spaccportprof-%s", tDn)),
					testAccCheckAciSpinePortSelectorIdNotEqual(&spine_port_selector_default, &spine_port_selector_updated),
				),
			},
			{
				Config: CreateAccSpinePortSelectorConfig(infraSpinePName, tDn),
			},
			{
				Config: CreateAccSpinePortSelectorConfigWithRequiredParams(rName, tDnUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpinePortSelectorExists(resourceName, &spine_port_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "spine_profile_dn", fmt.Sprintf("uni/infra/spprof-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/infra/spaccportprof-%s", tDnUpdated)),
					testAccCheckAciSpinePortSelectorIdNotEqual(&spine_port_selector_default, &spine_port_selector_updated),
				),
			},
		},
	})
}

func TestAccAciSpinePortSelector_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	tDn := makeTestVariable(acctest.RandString(5))
	infraSpinePName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpinePortSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSpinePortSelectorConfig(infraSpinePName, tDn),
			},
			{
				Config:      CreateAccSpinePortSelectorWithInValidParentDn(rName, tDn),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccSpinePortSelectorUpdatedAttr(infraSpinePName, tDn, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccSpinePortSelectorUpdatedAttr(infraSpinePName, tDn, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSpinePortSelectorConfig(infraSpinePName, tDn),
			},
		},
	})
}

func TestAccAciSpinePortSelector_MultipleCreateDelete(t *testing.T) {
	tDn := makeTestVariable(acctest.RandString(5))
	infraSpinePName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpinePortSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSpinePortSelectorConfigMultiple(infraSpinePName, tDn),
			},
		},
	})
}

func testAccCheckAciSpinePortSelectorExists(name string, spine_port_selector *models.InterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Spine Port Selector %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spine Port Selector dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		spine_port_selectorFound := models.InterfaceProfileFromContainer(cont)
		if spine_port_selectorFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spine Port Selector %s not found", rs.Primary.ID)
		}
		*spine_port_selector = *spine_port_selectorFound
		return nil
	}
}

func testAccCheckAciSpinePortSelectorDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing spine_port_selector destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_spine_port_selector" {
			cont, err := client.Get(rs.Primary.ID)
			spine_port_selector := models.InterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spine Port Selector %s Still exists", spine_port_selector.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSpinePortSelectorIdEqual(m1, m2 *models.InterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("spine_port_selector DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSpinePortSelectorIdNotEqual(m1, m2 *models.InterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("spine_port_selector DNs are equal")
		}
		return nil
	}
}

func CreateSpinePortSelectorWithoutRequired(infraSpinePName, tDn, attrName string) string {
	fmt.Println("=== STEP  Basic: testing spine_port_selector creation without ", attrName)
	rBlock := `
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "spine_profile_dn":
		rBlock += `
	resource "aci_spine_port_selector" "test" {
	#	spine_profile_dn  = aci_spine_profile.test.id
		tdn  = "%s"
	}
		`
	case "tdn":
		rBlock += `
	resource "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
	#	tdn  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, infraSpinePName, tDn)
}

func CreateAccSpinePortSelectorConfigWithRequiredParams(infraSpinePName, tDn string) string {
	fmt.Println("=== STEP  testing spine_port_selector creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_interface_profile" "test" {
		name        = "%s"
	}
	
	resource "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		tdn  = aci_spine_interface_profile.test.id
	}
	`, infraSpinePName, tDn)
	return resource
}

func CreateAccSpinePortSelectorConfig(infraSpinePName, tDn string) string {
	fmt.Println("=== STEP  testing spine_port_selector creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	resource "aci_spine_interface_profile" "test" {
		name        = "%s"
	}
	
	resource "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		tdn  = aci_spine_interface_profile.test.id
	}
	`, infraSpinePName, tDn)
	return resource
}

func CreateAccSpinePortSelectorConfigMultiple(infraSpinePName, tDn string) string {
	fmt.Println("=== STEP  testing multiple spine_port_selector creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	resource "aci_spine_interface_profile" "test" {
		count = 5
		name  = "%s_${count.index}"
	}

	resource "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		tdn  = aci_spine_interface_profile.test[count.index].id
		count = 5
	}
	`, infraSpinePName, tDn)
	return resource
}

func CreateAccSpinePortSelectorWithInValidParentDn(rName, tDn string) string {
	fmt.Println("=== STEP  Negative Case: testing spine_port_selector creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_spine_interface_profile" "test" {
		name        = "%s"
	}

	resource "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_tenant.test.id
		tdn  = aci_spine_interface_profile.test.id	
	}
	`, rName, tDn)
	return resource
}

func CreateAccSpinePortSelectorConfigWithOptionalValues(infraSpinePName, tDn string) string {
	fmt.Println("=== STEP  Basic: testing spine_port_selector creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	resource "aci_spine_interface_profile" "test" {
		name        = "%s"
	}
	
	resource "aci_spine_port_selector" "test" {
		spine_profile_dn  = "${aci_spine_profile.test.id}"
		tdn  = aci_spine_interface_profile.test.id
		annotation = "orchestrator:terraform_testacc"
		
	}
	`, infraSpinePName, tDn)

	return resource
}

func CreateAccSpinePortSelectorRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing spine_port_selector updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_spine_port_selector" "test" {
		annotation = "orchestrator:terraform_testacc"
		
	}
	`)

	return resource
}

func CreateAccSpinePortSelectorUpdatedAttr(infraSpinePName, tDn, attribute, value string) string {
	fmt.Printf("=== STEP  testing spine_port_selector attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	resource "aci_spine_interface_profile" "test" {
		name        = "%s"
	}
	resource "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		tdn  = aci_spine_interface_profile.test.id
		%s = "%s"
	}
	`, infraSpinePName, tDn, attribute, value)
	return resource
}
