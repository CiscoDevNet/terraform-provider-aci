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

func TestAccAciAccessPortBlock_Basic(t *testing.T) {
	var access_port_block_default models.AccessPortBlock
	var access_port_block_updated models.AccessPortBlock
	resourceName := "aci_access_port_block.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	infraAccPortPName := makeTestVariable(acctest.RandString(5))
	infraHPortSName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessPortBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccessPortBlockWithoutRequired(infraAccPortPName, infraHPortSName, rName, "access_port_selector_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccessPortBlockWithoutRequired(infraAccPortPName, infraHPortSName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAccessPortBlockConfig(infraAccPortPName, infraHPortSName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortBlockExists(resourceName, &access_port_block_default),
					resource.TestCheckResourceAttr(resourceName, "access_port_selector_dn", fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-ALL", infraAccPortPName, infraHPortSName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "from_card", "1"),
					resource.TestCheckResourceAttr(resourceName, "from_port", "1"),
					resource.TestCheckResourceAttr(resourceName, "to_card", "1"),
					resource.TestCheckResourceAttr(resourceName, "to_port", "1"),
				),
			},
			{
				Config: CreateAccAccessPortBlockConfigWithOptionalValues(infraAccPortPName, infraHPortSName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortBlockExists(resourceName, &access_port_block_updated),
					resource.TestCheckResourceAttr(resourceName, "access_port_selector_dn", fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-ALL", infraAccPortPName, infraHPortSName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_access_port_block"),
					resource.TestCheckResourceAttr(resourceName, "from_card", "2"),
					resource.TestCheckResourceAttr(resourceName, "from_port", "2"),
					resource.TestCheckResourceAttr(resourceName, "to_card", "2"),
					resource.TestCheckResourceAttr(resourceName, "to_port", "2"),

					testAccCheckAciAccessPortBlockIdEqual(&access_port_block_default, &access_port_block_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccAccessPortBlockConfigUpdatedName(infraAccPortPName, infraHPortSName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccAccessPortBlockRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAccessPortBlockConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortBlockExists(resourceName, &access_port_block_updated),
					resource.TestCheckResourceAttr(resourceName, "access_port_selector_dn", fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-ALL", rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciAccessPortBlockIdNotEqual(&access_port_block_default, &access_port_block_updated),
				),
			},
			{
				Config: CreateAccAccessPortBlockConfig(infraAccPortPName, infraHPortSName, rName),
			},
			{
				Config: CreateAccAccessPortBlockConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortBlockExists(resourceName, &access_port_block_updated),
					resource.TestCheckResourceAttr(resourceName, "access_port_selector_dn", fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-ALL", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciAccessPortBlockIdNotEqual(&access_port_block_default, &access_port_block_updated),
				),
			},
		},
	})
}

func TestAccAciAccessPortBlock_Update(t *testing.T) {
	var access_port_block_default models.AccessPortBlock
	var access_port_block_updated models.AccessPortBlock
	resourceName := "aci_access_port_block.test"
	rName := makeTestVariable(acctest.RandString(5))

	infraAccPortPName := makeTestVariable(acctest.RandString(5))
	infraHPortSName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessPortBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAccessPortBlockUpdatedCardAttr(infraAccPortPName, infraHPortSName, rName, "50", "55"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortBlockExists(resourceName, &access_port_block_updated),
					resource.TestCheckResourceAttr(resourceName, "from_card", "50"),
					resource.TestCheckResourceAttr(resourceName, "to_card", "55"),
					testAccCheckAciAccessPortBlockIdNotEqual(&access_port_block_default, &access_port_block_updated),
				),
			},
			{
				Config: CreateAccAccessPortBlockUpdatedCardAttr(infraAccPortPName, infraHPortSName, rName, "100", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortBlockExists(resourceName, &access_port_block_updated),
					resource.TestCheckResourceAttr(resourceName, "from_card", "100"),
					resource.TestCheckResourceAttr(resourceName, "to_card", "100"),
					testAccCheckAciAccessPortBlockIdNotEqual(&access_port_block_default, &access_port_block_updated),
				),
			},
			{
				Config: CreateAccAccessPortBlockUpdatedPortAttr(infraAccPortPName, infraHPortSName, rName, "50", "55"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortBlockExists(resourceName, &access_port_block_updated),
					resource.TestCheckResourceAttr(resourceName, "from_port", "50"),
					resource.TestCheckResourceAttr(resourceName, "to_port", "55"),
					testAccCheckAciAccessPortBlockIdNotEqual(&access_port_block_default, &access_port_block_updated),
				),
			},
			{
				Config: CreateAccAccessPortBlockUpdatedPortAttr(infraAccPortPName, infraHPortSName, rName, "127", "127"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortBlockExists(resourceName, &access_port_block_updated),
					resource.TestCheckResourceAttr(resourceName, "from_port", "127"),
					resource.TestCheckResourceAttr(resourceName, "to_port", "127"),
					testAccCheckAciAccessPortBlockIdNotEqual(&access_port_block_default, &access_port_block_updated),
				),
			},
		},
	})
}

func TestAccAciAccessPortBlock_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	infraAccPortPName := makeTestVariable(acctest.RandString(5))
	infraHPortSName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessPortBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAccessPortBlockConfig(infraAccPortPName, infraHPortSName, rName),
			},
			{
				Config:      CreateAccAccessPortBlockWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "from_card", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "from_card", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "from_card", "101"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "from_port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "from_port", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "from_port", "128"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedPortAttr(infraAccPortPName, infraHPortSName, rName, "55", "50"),
				ExpectError: regexp.MustCompile(`cannot be less than`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedCardAttr(infraAccPortPName, infraHPortSName, rName, "55", "50"),
				ExpectError: regexp.MustCompile(`cannot be less than`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "to_card", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "to_card", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "to_card", "101"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "to_port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "to_port", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, "to_port", "128"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccAccessPortBlockConfig(infraAccPortPName, infraHPortSName, rName),
			},
		},
	})
}

func TestAccAciAccessPortBlock_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	infraAccPortPName := makeTestVariable(acctest.RandString(5))
	infraHPortSName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessPortBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAccessPortBlockConfigMultiple(infraAccPortPName, infraHPortSName, rName),
			},
		},
	})
}

func testAccCheckAciAccessPortBlockExists(name string, access_port_block *models.AccessPortBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Access Port Block %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Access Port Block dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		access_port_blockFound := models.AccessPortBlockFromContainer(cont)
		if access_port_blockFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Access Port Block %s not found", rs.Primary.ID)
		}
		*access_port_block = *access_port_blockFound
		return nil
	}
}

func testAccCheckAciAccessPortBlockDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing access_port_block destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_access_port_block" {
			cont, err := client.Get(rs.Primary.ID)
			access_port_block := models.AccessPortBlockFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Access Port Block %s Still exists", access_port_block.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciAccessPortBlockIdEqual(m1, m2 *models.AccessPortBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("access_port_block DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciAccessPortBlockIdNotEqual(m1, m2 *models.AccessPortBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("access_port_block DNs are equal")
		}
		return nil
	}
}

func CreateAccessPortBlockWithoutRequired(infraAccPortPName, infraHPortSName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing access_port_block creation without ", attrName)
	rBlock := `
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
		access_port_selector_type = "ALL"
	}
	
	`
	switch attrName {
	case "access_port_selector_dn":
		rBlock += `
	resource "aci_access_port_block" "test" {
	#	access_port_selector_dn  = aci_access_port_selector.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_access_port_block" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, infraAccPortPName, infraHPortSName, rName)
}

func CreateAccAccessPortBlockConfigWithRequiredParams(prName, rName string) string {
	fmt.Printf("=== STEP  testing access_port_block creation with parent resource name %s and resource name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		access_port_selector_type = "ALL"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
	}
	
	resource "aci_access_port_block" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
		name  = "%s"
	}
	`, prName, prName, rName)
	return resource
}
func CreateAccAccessPortBlockConfigUpdatedName(infraAccPortPName, infraHPortSName, rName string) string {
	fmt.Println("=== STEP  testing access_port_block creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		access_port_selector_type = "ALL"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
	}
	
	resource "aci_access_port_block" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
		name  = "%s"
	}
	`, infraAccPortPName, infraHPortSName, rName)
	return resource
}

func CreateAccAccessPortBlockConfig(infraAccPortPName, infraHPortSName, rName string) string {
	fmt.Println("=== STEP  testing access_port_block creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		access_port_selector_type  = "ALL"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
	}
	
	resource "aci_access_port_block" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
		name  = "%s"
	}
	`, infraAccPortPName, infraHPortSName, rName)
	return resource
}

func CreateAccAccessPortBlockConfigMultiple(infraAccPortPName, infraHPortSName, rName string) string {
	fmt.Println("=== STEP  testing multiple access_port_block creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		access_port_selector_type  = "ALL"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
	}
	
	resource "aci_access_port_block" "test" {
		count = 5
		access_port_selector_dn  = aci_access_port_selector.test.id
		name  = "%s_${count.index+1}"
		from_port = (count.index+1)*10
		to_port = (count.index+1)*10+5
	}
	`, infraAccPortPName, infraHPortSName, rName)
	return resource
}

func CreateAccAccessPortBlockWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing access_port_block creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_access_port_block" "test" {
		access_port_selector_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccAccessPortBlockConfigWithOptionalValues(infraAccPortPName, infraHPortSName, rName string) string {
	fmt.Println("=== STEP  Basic: testing access_port_block creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		access_port_selector_type  = "ALL"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
	}
	
	resource "aci_access_port_block" "test" {
		access_port_selector_dn  = "${aci_access_port_selector.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_access_port_block"
		from_card = "2"
		from_port = "2"
		to_card = "2"
		to_port = "2"
		
	}
	`, infraAccPortPName, infraHPortSName, rName)

	return resource
}

func CreateAccAccessPortBlockRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing access_port_block updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_access_port_block" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_access_port_block"
		from_card = "2"
		from_port = "2"
		to_card = "2"
		to_port = "2"
		
	}
	`)

	return resource
}

func CreateAccAccessPortBlockUpdatedAttr(infraAccPortPName, infraHPortSName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing access_port_block attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		access_port_selector_type  = "ALL"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
	}
	
	resource "aci_access_port_block" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
		name  = "%s"
		%s = "%s"
	}
	`, infraAccPortPName, infraHPortSName, rName, attribute, value)
	return resource
}

func CreateAccAccessPortBlockUpdatedCardAttr(infraAccPortPName, infraHPortSName, rName, from, to string) string {
	fmt.Printf("=== STEP  testing access_port_block  from_card = \"%s\" and to_card = \"%s\" \n", from, to)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		access_port_selector_type  = "ALL"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
	}
	
	resource "aci_access_port_block" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
		name  = "%s"
		from_card = "%s"
		to_card = "%s"
	}
	`, infraAccPortPName, infraHPortSName, rName, from, to)
	return resource
}

func CreateAccAccessPortBlockUpdatedPortAttr(infraAccPortPName, infraHPortSName, rName, from, to string) string {
	fmt.Printf("=== STEP  testing access_port_block  from_port = \"%s\" and to_port = \"%s\" \n", from, to)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		access_port_selector_type  = "ALL"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
	}
	
	resource "aci_access_port_block" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
		name  = "%s"
		from_port = "%s"
		to_port = "%s"
	}
	`, infraAccPortPName, infraHPortSName, rName, from, to)
	return resource
}
