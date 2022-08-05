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

func TestAccAciNodeMgmtEpg_Basic(t *testing.T) {
	var node_mgmt_epg_default models.InBandManagementEPg
	var node_mgmt_epg_updated models.InBandManagementEPg
	var node_oob_mgmt_epg_updated models.OutOfBandManagementEPg
	resourceName := "aci_node_mgmt_epg.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	nodeType := "in_band"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeMgmtEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateNodeMgmtEpgWithoutRequired(nodeType, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateNodeMgmtEpgWithoutRequired(nodeType, rName, "type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccNodeMgmtEpgConfig(nodeType, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeMgmtEpgExists(resourceName, &node_mgmt_epg_default),
					resource.TestCheckResourceAttr(resourceName, "management_profile_dn", "uni/tn-mgmt/mgmtp-default"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "encap", "unknown"),
					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtleastOne"),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "exclude"),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
				),
			},
			{
				Config: CreateAccNodeMgmtEpgConfigWithOptionalValues(nodeType, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeMgmtEpgExists(resourceName, &node_mgmt_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "management_profile_dn", "uni/tn-mgmt/mgmtp-default"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_node_mgmt_epg"),

					resource.TestCheckResourceAttr(resourceName, "encap", "unknown"),

					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "enabled"),

					resource.TestCheckResourceAttr(resourceName, "match_t", "All"),

					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "include"),

					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),

					testAccCheckAciNodeMgmtEpgIdEqual(&node_mgmt_epg_default, &node_mgmt_epg_updated),
				),
			},
			{
				Config:      CreateAccNodeMgmtEpgConfigUpdatedName(nodeType, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccNodeMgmtEpgRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccNodeMgmtEpgConfigWithRequiredParams("out_of_band", rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeMgmtOoBEpgExists(resourceName, &node_oob_mgmt_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "out_of_band"),
					testAccCheckAciNodeOoBMgmtEpgIdNotEqual(&node_mgmt_epg_default, &node_oob_mgmt_epg_updated),
				),
			},
			{
				Config: CreateAccNodeMgmtEpgConfig(nodeType, rName),
			},
			{
				Config: CreateAccNodeMgmtEpgConfigWithRequiredParams(nodeType, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeMgmtEpgExists(resourceName, &node_mgmt_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "type", nodeType),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciNodeMgmtEpgIdNotEqual(&node_mgmt_epg_default, &node_mgmt_epg_updated),
				),
			},
		},
	})
}

func TestAccAciNodeMgmtEpg_Update(t *testing.T) {
	var node_mgmt_epg_default models.InBandManagementEPg
	var node_mgmt_epg_updated models.InBandManagementEPg
	resourceName := "aci_node_mgmt_epg.test"
	rName := makeTestVariable(acctest.RandString(5))
	nodeType := "in_band"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeMgmtEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccNodeMgmtEpgConfig(nodeType, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeMgmtEpgExists(resourceName, &node_mgmt_epg_default),
				),
			},

			{
				Config: CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "match_t", "AtmostOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeMgmtEpgExists(resourceName, &node_mgmt_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtmostOne"),
					testAccCheckAciNodeMgmtEpgIdEqual(&node_mgmt_epg_default, &node_mgmt_epg_updated),
				),
			},
			{
				Config: CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "match_t", "None"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeMgmtEpgExists(resourceName, &node_mgmt_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "None"),
					testAccCheckAciNodeMgmtEpgIdEqual(&node_mgmt_epg_default, &node_mgmt_epg_updated),
				),
			},
			{
				Config: CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeMgmtEpgExists(resourceName, &node_mgmt_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciNodeMgmtEpgIdEqual(&node_mgmt_epg_default, &node_mgmt_epg_updated),
				),
			},
			{
				Config: CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeMgmtEpgExists(resourceName, &node_mgmt_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciNodeMgmtEpgIdEqual(&node_mgmt_epg_default, &node_mgmt_epg_updated),
				),
			},
			{
				Config: CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeMgmtEpgExists(resourceName, &node_mgmt_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciNodeMgmtEpgIdEqual(&node_mgmt_epg_default, &node_mgmt_epg_updated),
				),
			},
			{
				Config: CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeMgmtEpgExists(resourceName, &node_mgmt_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciNodeMgmtEpgIdEqual(&node_mgmt_epg_default, &node_mgmt_epg_updated),
				),
			},
			{
				Config: CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeMgmtEpgExists(resourceName, &node_mgmt_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciNodeMgmtEpgIdEqual(&node_mgmt_epg_default, &node_mgmt_epg_updated),
				),
			},
			{
				Config: CreateAccNodeMgmtEpgConfig(nodeType, rName),
			},
		},
	})
}

func TestAccAciNodeMgmtEpg_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	nodeType := "in_band"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeMgmtEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccNodeMgmtEpgConfig(nodeType, rName),
			},
			{
				Config:      CreateAccNodeMgmtEpgWithInValidParentDn(nodeType, rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "encap", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "flood_on_encap", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "match_t", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "pref_gr_memb", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, "prio", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccNodeMgmtEpgConfig(nodeType, rName),
			},
		},
	})
}

func TestAccAciNodeMgmtEpg_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeMgmtEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccNodeMgmtEpgConfigMultiple("in_band", rName),
			},
			{
				Config: CreateAccNodeMgmtEpgConfigMultiple("out_of_band", rName),
			},
		},
	})
}

func testAccCheckAciNodeMgmtEpgExists(name string, node_mgmt_epg *models.InBandManagementEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Node Mgmt Epg %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Node Mgmt Epg dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		node_mgmt_epgFound := models.InBandManagementEPgFromContainer(cont)
		if node_mgmt_epgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Node Mgmt Epg %s not found", rs.Primary.ID)
		}
		*node_mgmt_epg = *node_mgmt_epgFound
		return nil
	}
}

func testAccCheckAciNodeMgmtOoBEpgExists(name string, node_mgmt_epg *models.OutOfBandManagementEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Node Mgmt Epg %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Node Mgmt Epg dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		node_mgmt_epgFound := models.OutOfBandManagementEPgFromContainer(cont)
		if node_mgmt_epgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Node Mgmt Epg %s not found", rs.Primary.ID)
		}
		*node_mgmt_epg = *node_mgmt_epgFound
		return nil
	}
}

func testAccCheckAciNodeMgmtEpgDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing node_mgmt_epg destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_node_mgmt_epg" {
			cont, err := client.Get(rs.Primary.ID)
			node_mgmt_epg := models.InBandManagementEPgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Node Mgmt Epg %s Still exists", node_mgmt_epg.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciNodeMgmtEpgIdEqual(m1, m2 *models.InBandManagementEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("node_mgmt_epg DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciNodeMgmtEpgIdNotEqual(m1, m2 *models.InBandManagementEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("node_mgmt_epg DNs are equal")
		}
		return nil
	}
}

func testAccCheckAciNodeOoBMgmtEpgIdNotEqual(m1 *models.InBandManagementEPg, m2 *models.OutOfBandManagementEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("node_mgmt_epg DNs are equal")
		}
		return nil
	}
}

func CreateNodeMgmtEpgWithoutRequired(nodeType, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing node_mgmt_epg creation without ", attrName)
	rBlock := `
	`
	switch attrName {
	case "type":
		rBlock += `
	resource "aci_node_mgmt_epg" "test" {
	#	type = "%s"
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_node_mgmt_epg" "test" {
		type = "%s"
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, nodeType, rName)
}

func CreateAccNodeMgmtEpgConfigWithRequiredParams(nodeType, rName string) string {
	fmt.Println("=== STEP  testing node_mgmt_epg creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		type = "%s"
		name  = "%s"
	}
	`, nodeType, rName)
	return resource
}
func CreateAccNodeMgmtEpgConfigUpdatedName(nodeType, rName string) string {
	fmt.Println("=== STEP  testing node_mgmt_epg creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		type = "%s"
		name  = "%s"
	}
	`, nodeType, rName)
	return resource
}

func CreateAccNodeMgmtEpgConfig(nodeType, rName string) string {
	fmt.Println("=== STEP  testing node_mgmt_epg creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_node_mgmt_epg" "test" {
		type = "%s"
		name  = "%s"
	}
	`, nodeType, rName)
	return resource
}

func CreateAccNodeMgmtEpgConfigMultiple(nodeType, rName string) string {
	fmt.Println("=== STEP  testing multiple node_mgmt_epg creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		type = "%s"
		name  = "%s_${count.index}"
		count = 5
	}
	`, nodeType, rName)
	return resource
}

func CreateAccNodeMgmtEpgWithInValidParentDn(nodeType, rName string) string {
	fmt.Println("=== STEP  Negative Case: testing node_mgmt_epg creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_node_mgmt_epg" "test" {
		management_profile_dn  = aci_tenant.test.id
		type = "%s"
		name = "%s"	
	}
	`, rName, nodeType, rName)
	return resource
}

func CreateAccNodeMgmtEpgConfigWithOptionalValues(nodeType, rName string) string {
	fmt.Println("=== STEP  Basic: testing node_mgmt_epg creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		type = "%s"
		name = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_node_mgmt_epg"
		encap = ""
		flood_on_encap = "enabled"
		match_t = "All"
		pref_gr_memb = "include"
		prio = "level1"
		
	}
	`, nodeType, rName)

	return resource
}

func CreateAccNodeMgmtEpgRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing node_mgmt_epg updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_node_mgmt_epg" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_node_mgmt_epg"
		encap = ""
		flood_on_encap = "enabled"
		match_t = "All"
		pref_gr_memb = "include"
		prio = "level1"
		
	}
	`)

	return resource
}

func CreateAccNodeMgmtEpgUpdatedAttr(nodeType, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing node_mgmt_epg attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		type = "%s"
		name = "%s"
		%s = "%s"
	}
	`, nodeType, rName, attribute, value)
	return resource
}
