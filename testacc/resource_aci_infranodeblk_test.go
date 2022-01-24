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

func TestAccAciNodeBlock_Basic(t *testing.T) {
	var node_block_default models.NodeBlock
	var node_block_updated models.NodeBlock
	resourceName := "aci_node_block.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateNodeBlockWithoutRequired(rName, rName, rName, "switch_association_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateNodeBlockWithoutRequired(rName, rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccNodeBlockConfig(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockExists(resourceName, &node_block_default),
					resource.TestCheckResourceAttr(resourceName, "switch_association_dn", fmt.Sprintf("uni/infra/nprof-%s/leaves-%s-typ-ALL", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "from_", "1"),
					resource.TestCheckResourceAttr(resourceName, "to_", "1"),
				),
			},
			{
				Config: CreateAccNodeBlockConfigWithOptionalValues(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockExists(resourceName, &node_block_updated),
					resource.TestCheckResourceAttr(resourceName, "switch_association_dn", fmt.Sprintf("uni/infra/nprof-%s/leaves-%s-typ-ALL", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_node_block"),
					resource.TestCheckResourceAttr(resourceName, "from_", "1600"),
					resource.TestCheckResourceAttr(resourceName, "to_", "1600"),
					testAccCheckAciNodeBlockIdEqual(&node_block_default, &node_block_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccNodeBlockConfigUpdatedName(rName, rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccNodeBlockRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccNodeBlockConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockExists(resourceName, &node_block_updated),
					resource.TestCheckResourceAttr(resourceName, "switch_association_dn", fmt.Sprintf("uni/infra/nprof-%s/leaves-%s-typ-ALL", rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciNodeBlockIdNotEqual(&node_block_default, &node_block_updated),
				),
			},
			{
				Config: CreateAccNodeBlockConfig(rName, rName, rName),
			},
			{
				Config: CreateAccNodeBlockConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockExists(resourceName, &node_block_updated),
					resource.TestCheckResourceAttr(resourceName, "switch_association_dn", fmt.Sprintf("uni/infra/nprof-%s/leaves-%s-typ-ALL", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciNodeBlockIdNotEqual(&node_block_default, &node_block_updated),
				),
			},
		},
	})
}

func TestAccAciNodeBlock_Update(t *testing.T) {
	var node_block_default models.NodeBlock
	var node_block_updated models.NodeBlock
	resourceName := "aci_node_block.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccNodeBlockConfig(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockExists(resourceName, &node_block_default),
				),
			},
			{
				Config: CreateAccNodeBlockUpdatedAttr(rName, rName, rName, "to_", "8000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockExists(resourceName, &node_block_updated),
					resource.TestCheckResourceAttr(resourceName, "to_", "8000"),
					testAccCheckAciNodeBlockIdEqual(&node_block_default, &node_block_updated),
				),
			},
			{
				Config: CreateAccNodeBlockUpdatedAttr(rName, rName, rName, "from_", "8000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockExists(resourceName, &node_block_updated),
					resource.TestCheckResourceAttr(resourceName, "from_", "8000"),
					testAccCheckAciNodeBlockIdEqual(&node_block_default, &node_block_updated),
				),
			},
		},
	})
}

func TestAccAciNodeBlock_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccNodeBlockConfig(rName, rName, rName),
			},
			{
				Config:      CreateAccNodeBlockWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccNodeBlockUpdatedAttr(rName, rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccNodeBlockUpdatedAttr(rName, rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccNodeBlockUpdatedAttr(rName, rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccNodeBlockUpdatedAttr(rName, rName, rName, "from_", "0"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccNodeBlockUpdatedAttr(rName, rName, rName, "from_", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccNodeBlockUpdatedAttr(rName, rName, rName, "to_", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccNodeBlockUpdatedAttr(rName, rName, rName, "from_", "16001"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccNodeBlockUpdatedAttr(rName, rName, rName, "to_", "0"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccNodeBlockUpdatedAttr(rName, rName, rName, "to_", "16001"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccNodeBlockUpdatedAttr(rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccNodeBlockUpdatedAttr(rName, rName, rName, "to_", "100"),
			},
			{
				Config:      CreateAccNodeBlockUpdatedAttr(rName, rName, rName, "from_", "200"),
				ExpectError: regexp.MustCompile(`to_ cannot be less than from_`),
			},
			{
				Config: CreateAccNodeBlockConfig(rName, rName, rName),
			},
		},
	})
}

func TestAccAciNodeBlock_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccNodeBlockConfigMultiple(rName, rName, rName),
			},
		},
	})
}

func testAccCheckAciNodeBlockExists(name string, node_block *models.NodeBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Node Block %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Node Block dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		node_blockFound := models.NodeBlockFromContainerBLK(cont)
		if node_blockFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Node Block %s not found", rs.Primary.ID)
		}
		*node_block = *node_blockFound
		return nil
	}
}

func testAccCheckAciNodeBlockDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing node_block destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_node_block" {
			cont, err := client.Get(rs.Primary.ID)
			node_block := models.NodeBlockFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Node Block %s Still exists", node_block.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciNodeBlockIdEqual(m1, m2 *models.NodeBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("node_block DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciNodeBlockIdNotEqual(m1, m2 *models.NodeBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("node_block DNs are equal")
		}
		return nil
	}
}

func CreateNodeBlockWithoutRequired(mgmtNodeGrpName, infrazoneNodeGrpName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing node_block creation without ", attrName)
	rBlock := `
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}

  	resource "aci_leaf_selector" "test" {
    	name = "%s"
    	leaf_profile_dn = aci_leaf_profile.test.id
    	switch_association_type = "ALL"
  	}
	
	`
	switch attrName {
	case "switch_association_dn":
		rBlock += `
	resource "aci_node_block" "test" {
	#	switch_association_dn  = aci_leaf_selector.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_node_block" "test" {
		switch_association_dn  = aci_leaf_selector.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, mgmtNodeGrpName, infrazoneNodeGrpName, rName)
}

func CreateAccNodeBlockConfigWithRequiredParams(prName, rName string) string {
	fmt.Printf("=== STEP  testing node_block creation with parent resource name %s and resource name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}

  	resource "aci_leaf_selector" "test" {
    	name = "%s"
    	leaf_profile_dn = aci_leaf_profile.test.id
    	switch_association_type = "ALL"
  	}
	
	resource "aci_node_block" "test" {
		switch_association_dn  = aci_leaf_selector.test.id
		name  = "%s"
	}
	`, prName, prName, rName)
	return resource
}
func CreateAccNodeBlockConfigUpdatedName(mgmtNodeGrpName, infrazoneNodeGrpName, rName string) string {
	fmt.Println("=== STEP  testing node_block creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}

  	resource "aci_leaf_selector" "test" {
    	name = "%s"
    	leaf_profile_dn = aci_leaf_profile.test.id
    	switch_association_type = "ALL"
  	}
	
	resource "aci_node_block" "test" {
		switch_association_dn  = aci_leaf_selector.test.id
		name  = "%s"
	}
	`, mgmtNodeGrpName, infrazoneNodeGrpName, rName)
	return resource
}

func CreateAccNodeBlockConfig(mgmtNodeGrpName, infrazoneNodeGrpName, rName string) string {
	fmt.Println("=== STEP  testing node_block creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}

  	resource "aci_leaf_selector" "test" {
    	name = "%s"
    	leaf_profile_dn = aci_leaf_profile.test.id
    	switch_association_type = "ALL"
  	}
	
	resource "aci_node_block" "test" {
		switch_association_dn  = aci_leaf_selector.test.id
		name  = "%s"
	}
	`, mgmtNodeGrpName, infrazoneNodeGrpName, rName)
	return resource
}

func CreateAccNodeBlockConfigMultiple(mgmtNodeGrpName, infrazoneNodeGrpName, rName string) string {
	fmt.Println("=== STEP  testing multiple node_block creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}

  	resource "aci_leaf_selector" "test" {
    	name = "%s"
    	leaf_profile_dn = aci_leaf_profile.test.id
    	switch_association_type = "ALL"
  	}
	
	resource "aci_node_block" "test" {
		switch_association_dn  = aci_leaf_selector.test.id
		name  = "%s_${count.index}"
		from_ = (count.index+1)*10
		to_ = (count.index+1)*10+5
		count = 5
	}
	`, mgmtNodeGrpName, infrazoneNodeGrpName, rName)
	return resource
}

func CreateAccNodeBlockWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing node_block creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_node_block" "test" {
		switch_association_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccNodeBlockConfigWithOptionalValues(mgmtNodeGrpName, infrazoneNodeGrpName, rName string) string {
	fmt.Println("=== STEP  Basic: testing node_block creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}

  	resource "aci_leaf_selector" "test" {
    	name = "%s"
    	leaf_profile_dn = aci_leaf_profile.test.id
    	switch_association_type = "ALL"
  	}
	
	resource "aci_node_block" "test" {
		switch_association_dn  = "${aci_leaf_selector.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_node_block"
		from_ = "1600"
		to_ = "1600"
		
	}
	`, mgmtNodeGrpName, infrazoneNodeGrpName, rName)

	return resource
}

func CreateAccNodeBlockRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing node_block updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_node_block" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_node_block"
		from_ = "2"
		to_ = "2"
	}
	`)

	return resource
}

func CreateAccNodeBlockUpdatedAttr(mgmtNodeGrpName, infrazoneNodeGrpName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing node_block attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}

  	resource "aci_leaf_selector" "test" {
    	name = "%s"
    	leaf_profile_dn = aci_leaf_profile.test.id
    	switch_association_type = "ALL"
  	}
	
	resource "aci_node_block" "test" {
		switch_association_dn  = aci_leaf_selector.test.id
		name  = "%s"
		%s = "%s"
	}
	`, mgmtNodeGrpName, infrazoneNodeGrpName, rName, attribute, value)
	return resource
}
