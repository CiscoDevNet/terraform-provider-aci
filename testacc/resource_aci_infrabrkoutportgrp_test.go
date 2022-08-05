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

func TestAccAciLeafBreakoutPortGroup_Basic(t *testing.T) {
	var leaf_breakout_port_group_default models.LeafBreakoutPortGroup
	var leaf_breakout_port_group_updated models.LeafBreakoutPortGroup
	resourceName := "aci_leaf_breakout_port_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafBreakoutPortGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLeafBreakoutPortGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLeafBreakoutPortGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafBreakoutPortGroupExists(resourceName, &leaf_breakout_port_group_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "brkout_map", "none"),
				),
			},
			{
				// in this step all optional attribute expect realational attribute are given for the same resource and then compared
				Config: CreateAccLeafBreakoutPortGroupConfigWithOptionalValues(rName), // configuration to update optional filelds
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafBreakoutPortGroupExists(resourceName, &leaf_breakout_port_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_leaf_breakout_port_group"),
					resource.TestCheckResourceAttr(resourceName, "brkout_map", "100g-2x"),

					testAccCheckAciLeafBreakoutPortGroupIdEqual(&leaf_breakout_port_group_default, &leaf_breakout_port_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccLeafBreakoutPortGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccLeafBreakoutPortGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccLeafBreakoutPortGroupConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafBreakoutPortGroupExists(resourceName, &leaf_breakout_port_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciLeafBreakoutPortGroupIdNotEqual(&leaf_breakout_port_group_default, &leaf_breakout_port_group_updated),
				),
			},
		},
	})
}

func TestAccAciLeafBreakoutPortGroup_Update(t *testing.T) {
	var leaf_breakout_port_group_default models.LeafBreakoutPortGroup
	var leaf_breakout_port_group_updated models.LeafBreakoutPortGroup
	resourceName := "aci_leaf_breakout_port_group.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafBreakoutPortGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLeafBreakoutPortGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafBreakoutPortGroupExists(resourceName, &leaf_breakout_port_group_default),
				),
			},

			{
				Config: CreateAccLeafBreakoutPortGroupUpdatedAttr(rName, "brkout_map", "100g-4x"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafBreakoutPortGroupExists(resourceName, &leaf_breakout_port_group_updated),
					resource.TestCheckResourceAttr(resourceName, "brkout_map", "100g-4x"),
					testAccCheckAciLeafBreakoutPortGroupIdEqual(&leaf_breakout_port_group_default, &leaf_breakout_port_group_updated),
				),
			},
			{
				Config: CreateAccLeafBreakoutPortGroupUpdatedAttr(rName, "brkout_map", "10g-4x"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafBreakoutPortGroupExists(resourceName, &leaf_breakout_port_group_updated),
					resource.TestCheckResourceAttr(resourceName, "brkout_map", "10g-4x"),
					testAccCheckAciLeafBreakoutPortGroupIdEqual(&leaf_breakout_port_group_default, &leaf_breakout_port_group_updated),
				),
			},
			{
				Config: CreateAccLeafBreakoutPortGroupUpdatedAttr(rName, "brkout_map", "25g-4x"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafBreakoutPortGroupExists(resourceName, &leaf_breakout_port_group_updated),
					resource.TestCheckResourceAttr(resourceName, "brkout_map", "25g-4x"),
					testAccCheckAciLeafBreakoutPortGroupIdEqual(&leaf_breakout_port_group_default, &leaf_breakout_port_group_updated),
				),
			},
			{
				Config: CreateAccLeafBreakoutPortGroupUpdatedAttr(rName, "brkout_map", "50g-8x"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafBreakoutPortGroupExists(resourceName, &leaf_breakout_port_group_updated),
					resource.TestCheckResourceAttr(resourceName, "brkout_map", "50g-8x"),
					testAccCheckAciLeafBreakoutPortGroupIdEqual(&leaf_breakout_port_group_default, &leaf_breakout_port_group_updated),
				),
			},
			{
				Config: CreateAccLeafBreakoutPortGroupConfig(rName),
			},
		},
	})
}

func TestAccAciLeafBreakoutPortGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafBreakoutPortGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLeafBreakoutPortGroupConfig(rName),
			},

			{
				Config:      CreateAccLeafBreakoutPortGroupUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafBreakoutPortGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafBreakoutPortGroupUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafBreakoutPortGroupUpdatedAttr(rName, "brkout_map", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLeafBreakoutPortGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLeafBreakoutPortGroupConfig(rName),
			},
		},
	})
}

func TestAccAciLeafBreakoutPortGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafBreakoutPortGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLeafBreakoutPortGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciLeafBreakoutPortGroupExists(name string, leaf_breakout_port_group *models.LeafBreakoutPortGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Leaf Breakout Port Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Leaf Breakout Port Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		leaf_breakout_port_groupFound := models.LeafBreakoutPortGroupFromContainer(cont)
		if leaf_breakout_port_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Leaf Breakout Port Group %s not found", rs.Primary.ID)
		}
		*leaf_breakout_port_group = *leaf_breakout_port_groupFound
		return nil
	}
}

func testAccCheckAciLeafBreakoutPortGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing leaf_breakout_port_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_leaf_breakout_port_group" {
			cont, err := client.Get(rs.Primary.ID)
			leaf_breakout_port_group := models.LeafBreakoutPortGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Leaf Breakout Port Group %s Still exists", leaf_breakout_port_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLeafBreakoutPortGroupIdEqual(m1, m2 *models.LeafBreakoutPortGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("leaf_breakout_port_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLeafBreakoutPortGroupIdNotEqual(m1, m2 *models.LeafBreakoutPortGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("leaf_breakout_port_group DNs are equal")
		}
		return nil
	}
}

func CreateLeafBreakoutPortGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_breakout_port_group creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_leaf_breakout_port_group" "test" {
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLeafBreakoutPortGroupConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing leaf_breakout_port_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_breakout_port_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafBreakoutPortGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing leaf_breakout_port_group multiple creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_breakout_port_group" "test" {
		count = 5
		name  = "%s_${count.index}"
	}
	`, rName)
	return resource
}

func CreateAccLeafBreakoutPortGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing leaf_breakout_port_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_breakout_port_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafBreakoutPortGroupConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing leaf_breakout_port_group creation with Updated Name")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_breakout_port_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafBreakoutPortGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_breakout_port_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_breakout_port_group" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_leaf_breakout_port_group"
		brkout_map = "100g-2x"
	}
	`, rName)

	return resource
}

func CreateAccLeafBreakoutPortGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing leaf_breakout_port_group update without required parameters")
	resource := `
	resource "aci_leaf_breakout_port_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_leaf_breakout_port_group"
		brkout_map = "100g-2x"
	}
	`

	return resource
}

func CreateAccLeafBreakoutPortGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing leaf_breakout_port_group attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_breakout_port_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
