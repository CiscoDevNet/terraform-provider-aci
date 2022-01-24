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

func TestAccAciMgmtZone_Basic(t *testing.T) {
	var mgmt_zone_default models.InBManagedNodesZone
	var mgmt_zone_updated models.InBManagedNodesZone
	var mgmt_zone_oob_updated models.OOBManagedNodesZone
	resourceName := "aci_mgmt_zone.test"
	zoneType := "in_band"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	mgmtGrpName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMgmtZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateMgmtZoneWithoutRequired(mgmtGrpName, zoneType, rName, "managed_node_connectivity_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMgmtZoneWithoutRequired(mgmtGrpName, zoneType, rName, "type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMgmtZoneWithoutRequired(mgmtGrpName, zoneType, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMgmtZoneConfig(mgmtGrpName, zoneType, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMgmtZoneExists(resourceName, &mgmt_zone_default),
					resource.TestCheckResourceAttr(resourceName, "managed_node_connectivity_group_dn", fmt.Sprintf("uni/infra/funcprof/grp-%s", mgmtGrpName)),
					resource.TestCheckResourceAttr(resourceName, "type", zoneType),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccMgmtZoneConfigWithOptionalValues(mgmtGrpName, zoneType, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMgmtZoneExists(resourceName, &mgmt_zone_updated),
					resource.TestCheckResourceAttr(resourceName, "managed_node_connectivity_group_dn", fmt.Sprintf("uni/infra/funcprof/grp-%s", mgmtGrpName)),
					resource.TestCheckResourceAttr(resourceName, "type", zoneType),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_mgmt_zone"),

					testAccCheckAciMgmtZoneIdEqual(&mgmt_zone_default, &mgmt_zone_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccMgmtZoneRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccMgmtZoneConfigWithRequiredParams(rNameUpdated, zoneType, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccMgmtZoneConfigWithRequiredParams(rNameUpdated, acctest.RandString(5), rName),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config: CreateAccMgmtZoneConfigWithRequiredParams(rNameUpdated, zoneType, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMgmtZoneExists(resourceName, &mgmt_zone_updated),
					resource.TestCheckResourceAttr(resourceName, "managed_node_connectivity_group_dn", fmt.Sprintf("uni/infra/funcprof/grp-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "type", zoneType),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciMgmtZoneIdNotEqual(&mgmt_zone_default, &mgmt_zone_updated),
				),
			},
			{
				Config: CreateAccMgmtZoneConfigWithRequiredParams(rName, "out_of_band", rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMgmtOoBZoneExists(resourceName, &mgmt_zone_oob_updated),
					resource.TestCheckResourceAttr(resourceName, "managed_node_connectivity_group_dn", fmt.Sprintf("uni/infra/funcprof/grp-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "out_of_band"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciMgmtOoBZoneIdNotEqual(&mgmt_zone_default, &mgmt_zone_oob_updated),
				),
			},
			{
				Config: CreateAccMgmtZoneConfigWithRequiredParams(rName, zoneType, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMgmtZoneExists(resourceName, &mgmt_zone_updated),
					resource.TestCheckResourceAttr(resourceName, "managed_node_connectivity_group_dn", fmt.Sprintf("uni/infra/funcprof/grp-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "type", zoneType),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciMgmtZoneIdNotEqual(&mgmt_zone_default, &mgmt_zone_updated),
				),
			},
			{
				Config: CreateAccMgmtZoneConfig(mgmtGrpName, zoneType, rName),
			},
		},
	})
}

func TestAccAciMgmtZone_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	zoneType := "in_band"
	mgmtGrpName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMgmtZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMgmtZoneConfig(mgmtGrpName, zoneType, rName),
			},
			{
				Config:      CreateAccMgmtZoneWithInValidParentDn(rName, zoneType),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMgmtZoneUpdatedAttr(mgmtGrpName, zoneType, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMgmtZoneUpdatedAttr(mgmtGrpName, zoneType, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMgmtZoneUpdatedAttr(mgmtGrpName, zoneType, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccMgmtZoneUpdatedAttr(mgmtGrpName, zoneType, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccMgmtZoneConfig(mgmtGrpName, zoneType, rName),
			},
		},
	})
}

func testAccCheckAciMgmtZoneExists(name string, mgmt_zone *models.InBManagedNodesZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Mgmt Zone %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Mgmt Zone dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		mgmt_zoneFound := models.InBManagedNodesZoneFromContainer(cont)
		if mgmt_zoneFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Mgmt Zone %s not found", rs.Primary.ID)
		}
		*mgmt_zone = *mgmt_zoneFound
		return nil
	}
}

func testAccCheckAciMgmtOoBZoneExists(name string, mgmt_zone *models.OOBManagedNodesZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Mgmt Zone %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Mgmt Zone dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		mgmt_zoneFound := models.OOBManagedNodesZoneFromContainer(cont)
		if mgmt_zoneFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Mgmt Zone %s not found", rs.Primary.ID)
		}
		*mgmt_zone = *mgmt_zoneFound
		return nil
	}
}

func testAccCheckAciMgmtZoneDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing mgmt_zone destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_mgmt_zone" {
			cont, err := client.Get(rs.Primary.ID)
			mgmt_zone := models.InBManagedNodesZoneFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Mgmt Zone %s Still exists", mgmt_zone.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciMgmtZoneIdEqual(m1, m2 *models.InBManagedNodesZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("mgmt_zone DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciMgmtZoneIdNotEqual(m1, m2 *models.InBManagedNodesZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("mgmt_zone DNs are equal")
		}
		return nil
	}
}

func testAccCheckAciMgmtOoBZoneIdNotEqual(m1 *models.InBManagedNodesZone, m2 *models.OOBManagedNodesZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("mgmt_zone DNs are equal")
		}
		return nil
	}
}

func CreateMgmtZoneWithoutRequired(mgmtGrpName, zoneType, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing mgmt_zone creation without ", attrName)
	rBlock := `
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "managed_node_connectivity_group_dn":
		rBlock += `
	resource "aci_mgmt_zone" "test" {
	#	managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		type = "%s"
		name = "%s"
	
	}
		`
	case "type":
		rBlock += `
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
	#	type = "%s"
		name = "%s"
	
	}
		`
	case "name":
		rBlock += `
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		type = "%s"
	#	name = "%s"
	
	}
		`
	}
	return fmt.Sprintf(rBlock, mgmtGrpName, zoneType, rName)
}

func CreateAccMgmtZoneConfigWithRequiredParams(mgmtGrpName, zoneType, rName string) string {
	fmt.Printf("=== STEP  testing mgmt_zone creation with parent resource name %s, zone type %s and resource name %s\n", mgmtGrpName, zoneType, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		type = "%s"
		name = "%s"
	}
	`, mgmtGrpName, zoneType, rName)
	return resource
}
func CreateAccMgmtZoneConfig(mgmtGrpName, zoneType, rName string) string {
	fmt.Println("=== STEP  testing mgmt_zone creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		type = "%s"
		name = "%s"
	}
	`, mgmtGrpName, zoneType, rName)
	return resource
}

func CreateAccMgmtZoneWithInValidParentDn(rName, zoneType string) string {
	fmt.Println("=== STEP  Negative Case: testing mgmt_zone creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_tenant.test.id
		type = "%s"
		name = "%s"
	}
	`, rName, zoneType, rName)
	return resource
}

func CreateAccMgmtZoneConfigWithOptionalValues(mgmtGrpName, zoneType, rName string) string {
	fmt.Println("=== STEP  Basic: testing mgmt_zone creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = "${aci_managed_node_connectivity_group.test.id}"
		type = "%s"
		name = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_mgmt_zone"
		
	}
	`, mgmtGrpName, zoneType, rName)

	return resource
}

func CreateAccMgmtZoneRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing mgmt_zone updation without required parameters")
	resource := `
	resource "aci_mgmt_zone" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_mgmt_zone"
		
	}
	`

	return resource
}

func CreateAccMgmtZoneUpdatedAttr(mgmtGrpName, zoneType, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing mgmt_zone attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		type = "%s"
		name = "%s"
		%s = "%s"
	}
	`, mgmtGrpName, zoneType, rName, attribute, value)
	return resource
}
