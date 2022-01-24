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

func TestAccAciTACACSSource_Basic(t *testing.T) {
	var tacacs_source_default models.TACACSSource
	var tacacs_source_updated models.TACACSSource
	resourceName := "aci_tacacs_source.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSSourceDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateTACACSSourceWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTACACSSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists(resourceName, &tacacs_source_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", "uni/fabric/moncommon"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "incl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "incl.0", "audit"),
					resource.TestCheckResourceAttr(resourceName, "incl.1", "session"),
					resource.TestCheckResourceAttr(resourceName, "min_sev", "info"),
				),
			},
			{
				Config: CreateAccTACACSSourceConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists(resourceName, &tacacs_source_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", "uni/fabric/moncommon"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_tacacs_source"),
					resource.TestCheckResourceAttr(resourceName, "incl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "incl.0", "audit"),
					resource.TestCheckResourceAttr(resourceName, "min_sev", "cleared"),
					testAccCheckAciTACACSSourceIdEqual(&tacacs_source_default, &tacacs_source_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccTACACSSourceConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccTACACSSourceRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccTACACSSourceConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists(resourceName, &tacacs_source_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciTACACSSourceIdNotEqual(&tacacs_source_default, &tacacs_source_updated),
				),
			},
		},
	})
}

func TestAccAciTACACSSource_Update(t *testing.T) {
	var tacacs_source_default models.TACACSSource
	var tacacs_source_updated models.TACACSSource
	resourceName := "aci_tacacs_source.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists(resourceName, &tacacs_source_default),
				),
			},
			{

				Config: CreateAccTACACSSourceUpdatedAttrList(rName, "incl", StringListtoString([]string{"audit"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists(resourceName, &tacacs_source_updated),
					resource.TestCheckResourceAttr(resourceName, "incl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "incl.0", "audit"),
				),
			},
			{

				Config: CreateAccTACACSSourceUpdatedAttrList(rName, "incl", StringListtoString([]string{"events"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists(resourceName, &tacacs_source_updated),
					resource.TestCheckResourceAttr(resourceName, "incl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "incl.0", "events"),
				),
			},
			{

				Config: CreateAccTACACSSourceUpdatedAttrList(rName, "incl", StringListtoString([]string{"audit", "events", "faults", "session"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists(resourceName, &tacacs_source_updated),
					resource.TestCheckResourceAttr(resourceName, "incl.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "incl.0", "audit"),
					resource.TestCheckResourceAttr(resourceName, "incl.1", "events"),
					resource.TestCheckResourceAttr(resourceName, "incl.2", "faults"),
					resource.TestCheckResourceAttr(resourceName, "incl.3", "session"),
				),
			},

			{
				Config: CreateAccTACACSSourceUpdatedAttrList(rName, "incl", StringListtoString([]string{"session", "faults", "events", "audit"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists(resourceName, &tacacs_source_updated),
					resource.TestCheckResourceAttr(resourceName, "incl.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "incl.0", "session"),
					resource.TestCheckResourceAttr(resourceName, "incl.1", "faults"),
					resource.TestCheckResourceAttr(resourceName, "incl.2", "events"),
					resource.TestCheckResourceAttr(resourceName, "incl.3", "audit"),
				),
			},
			{
				Config: CreateAccTACACSSourceUpdatedAttr(rName, "min_sev", "critical"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists(resourceName, &tacacs_source_updated),
					resource.TestCheckResourceAttr(resourceName, "min_sev", "critical"),
					testAccCheckAciTACACSSourceIdEqual(&tacacs_source_default, &tacacs_source_updated),
				),
			},
			{
				Config: CreateAccTACACSSourceUpdatedAttr(rName, "min_sev", "major"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists(resourceName, &tacacs_source_updated),
					resource.TestCheckResourceAttr(resourceName, "min_sev", "major"),
					testAccCheckAciTACACSSourceIdEqual(&tacacs_source_default, &tacacs_source_updated),
				),
			},
			{
				Config: CreateAccTACACSSourceUpdatedAttr(rName, "min_sev", "minor"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists(resourceName, &tacacs_source_updated),
					resource.TestCheckResourceAttr(resourceName, "min_sev", "minor"),
					testAccCheckAciTACACSSourceIdEqual(&tacacs_source_default, &tacacs_source_updated),
				),
			},
			{
				Config: CreateAccTACACSSourceUpdatedAttr(rName, "min_sev", "warning"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists(resourceName, &tacacs_source_updated),
					resource.TestCheckResourceAttr(resourceName, "min_sev", "warning"),
					testAccCheckAciTACACSSourceIdEqual(&tacacs_source_default, &tacacs_source_updated),
				),
			},
			{
				Config: CreateAccTACACSSourceConfig(rName),
			},
		},
	})
}

func TestAccAciTACACSSource_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSSourceConfig(rName),
			},
			{
				Config:      CreateAccTACACSSourceWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccTACACSSourceUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSSourceUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSSourceUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSSourceUpdatedAttrList(rName, "incl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccTACACSSourceUpdatedAttrList(rName, "incl", StringListtoString([]string{"audit", "audit"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},

			{
				Config:      CreateAccTACACSSourceUpdatedAttr(rName, "min_sev", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccTACACSSourceUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccTACACSSourceConfig(rName),
			},
		},
	})
}

func TestAccAciTACACSSource_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSSourceConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciTACACSSourceExists(name string, tacacs_source *models.TACACSSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("TACACS Source %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No TACACS Source dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tacacs_sourceFound := models.TACACSSourceFromContainer(cont)
		if tacacs_sourceFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("TACACS Source %s not found", rs.Primary.ID)
		}
		*tacacs_source = *tacacs_sourceFound
		return nil
	}
}

func testAccCheckAciTACACSSourceDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing tacacs_source destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tacacs_source" {
			cont, err := client.Get(rs.Primary.ID)
			tacacs_source := models.TACACSSourceFromContainer(cont)
			if err == nil {
				return fmt.Errorf("TACACS Source %s Still exists", tacacs_source.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTACACSSourceIdEqual(m1, m2 *models.TACACSSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("tacacs_source DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciTACACSSourceIdNotEqual(m1, m2 *models.TACACSSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("tacacs_source DNs are equal")
		}
		return nil
	}
}

func CreateTACACSSourceWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_source creation without ", attrName)
	rBlock := `
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_tacacs_source" "test" {
		
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccTACACSSourceWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing tacacs_source creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tacacs_source" "test" {
		parent_dn   = "uni/tn-common"
  		name        = "%s"
	}
	`, rName)
	return resource
}

func CreateAccTACACSSourceConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing tacacs_source creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_source" "test" {
  		name        = "%s"
	}
	`, rName)
	return resource
}
func CreateAccTACACSSourceConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing tacacs_source creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_source" "test" {
		
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccTACACSSourceConfig(rName string) string {
	fmt.Println("=== STEP  testing tacacs_source creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_source" "test" {
		
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccTACACSSourceConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple tacacs_source creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_source" "test" {
		
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccTACACSSourceConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_source creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_source" "test" {
		
		name  = "%s"
		parent_dn   = "uni/fabric/moncommon"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tacacs_source"
		incl = ["audit"]
		min_sev = "cleared"
		
	}
	`, rName)

	return resource
}

func CreateAccTACACSSourceRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing tacacs_source updation without required parameters")
	resource := `
	resource "aci_tacacs_source" "test" {
		parent_dn   = "uni/fabric/moncommon"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tacacs_source"
		incl = ["all"]
		min_sev = "cleared"
		
	}
	`

	return resource
}

func CreateAccTACACSSourceUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing tacacs_source attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_source" "test" {
		
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccTACACSSourceUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing tacacs_source attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_source" "test" {
		
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
