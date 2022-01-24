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

func TestAccAciStaticNodeMgmtAddress_Basic(t *testing.T) {
	var static_node_mgmt_address_default models.InbandStaticNode
	var static_node_mgmt_address_updated models.InbandStaticNode
	var static_node_mgmt_address_oob_updated models.OutofbandStaticNode
	resourceName := "aci_static_node_mgmt_address.test"

	rNameUpdated := makeTestVariable(acctest.RandString(5))
	addrType := "in_band"
	tDn := "topology/pod-1/node-1"

	mgmtInBName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciStaticNodeMgmtAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateStaticNodeMgmtAddressWithoutRequired(mgmtInBName, addrType, tDn, "management_epg_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateStaticNodeMgmtAddressWithoutRequired(mgmtInBName, addrType, tDn, "t_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccStaticNodeMgmtAddressConfig(mgmtInBName, addrType, tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciStaticNodeMgmtAddressExists(resourceName, &static_node_mgmt_address_default),
					resource.TestCheckResourceAttr(resourceName, "management_epg_dn", fmt.Sprintf("uni/tn-mgmt/mgmtp-default/inb-%s", mgmtInBName)),
					resource.TestCheckResourceAttr(resourceName, "t_dn", tDn),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "addr", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "gw", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "v6_addr", "::"),
					resource.TestCheckResourceAttr(resourceName, "v6_gw", "::"),
					resource.TestCheckResourceAttr(resourceName, "type", addrType),
				),
			},
			{
				Config: CreateAccStaticNodeMgmtAddressConfigWithOptionalValues(mgmtInBName, addrType, tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciStaticNodeMgmtAddressExists(resourceName, &static_node_mgmt_address_updated),
					resource.TestCheckResourceAttr(resourceName, "management_epg_dn", fmt.Sprintf("uni/tn-mgmt/mgmtp-default/inb-%s", mgmtInBName)),
					resource.TestCheckResourceAttr(resourceName, "t_dn", tDn),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "addr", "1.0.0.1/24"),
					resource.TestCheckResourceAttr(resourceName, "gw", "1.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "type", addrType),
					// resource.TestCheckResourceAttr(resourceName, "v6_addr", "2001:0db8:85a3:0000:0000:8a2e:0370:7334/64"),
					// resource.TestCheckResourceAttr(resourceName, "v6_gw", "2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
					testAccCheckAciStaticNodeMgmtAddressIdEqual(&static_node_mgmt_address_default, &static_node_mgmt_address_updated),
				),
			},

			{
				Config:      CreateAccStaticNodeMgmtAddressRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccStaticNodeMgmtAddressConfigWithRequiredParams(rNameUpdated, addrType, tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciStaticNodeMgmtAddressExists(resourceName, &static_node_mgmt_address_updated),
					resource.TestCheckResourceAttr(resourceName, "management_epg_dn", fmt.Sprintf("uni/tn-mgmt/mgmtp-default/inb-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "t_dn", tDn),
					resource.TestCheckResourceAttr(resourceName, "type", addrType),
					testAccCheckAciStaticNodeMgmtAddressIdNotEqual(&static_node_mgmt_address_default, &static_node_mgmt_address_updated),
				),
			},
			{
				Config: CreateAccStaticNodeMgmtAddressConfig(mgmtInBName, "out_of_band", tDn),
			},
			{
				Config: CreateAccStaticNodeMgmtAddressConfigWithRequiredParams(rNameUpdated, "out_of_band", tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciStaticNodeMgmtAddressOoBExists(resourceName, &static_node_mgmt_address_oob_updated),
					resource.TestCheckResourceAttr(resourceName, "management_epg_dn", fmt.Sprintf("uni/tn-mgmt/mgmtp-default/oob-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "type", "out_of_band"),
					resource.TestCheckResourceAttr(resourceName, "t_dn", tDn),
					testAccCheckAciStaticNodeMgmtAddressOoBIdNotEqual(&static_node_mgmt_address_default, &static_node_mgmt_address_oob_updated),
				),
			},
		},
	})
}

func TestAccAciStaticNodeMgmtAddress_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	tDn := "topology/pod-1/node-1"
	addrType := "in_band"
	mgmtInBName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciStaticNodeMgmtAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccStaticNodeMgmtAddressConfig(mgmtInBName, addrType, tDn),
			},
			{
				Config:      CreateAccStaticNodeMgmtAddressWithInValidParentDn(rName, tDn),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccStaticNodeMgmtAddressUpdatedAttr(mgmtInBName, addrType, tDn, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccStaticNodeMgmtAddressUpdatedAttr(mgmtInBName, addrType, tDn, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccStaticNodeMgmtAddressUpdatedAttr(mgmtInBName, addrType, tDn, "addr", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccStaticNodeMgmtAddressUpdatedAttr(mgmtInBName, addrType, tDn, "gw", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccStaticNodeMgmtAddressUpdatedAttr(mgmtInBName, addrType, tDn, "v6_addr", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccStaticNodeMgmtAddressUpdatedAttr(mgmtInBName, addrType, tDn, "v6_gw", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccStaticNodeMgmtAddressUpdatedAttr(mgmtInBName, addrType, tDn, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccStaticNodeMgmtAddressConfig(mgmtInBName, addrType, tDn),
			},
		},
	})
}

func testAccCheckAciStaticNodeMgmtAddressExists(name string, static_node_mgmt_address *models.InbandStaticNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Static Node Mgmt Address %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Static Node Mgmt Address dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		static_node_mgmt_addressFound := models.InbandStaticNodeFromContainer(cont)
		if static_node_mgmt_addressFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Static Node Mgmt Address %s not found", rs.Primary.ID)
		}
		*static_node_mgmt_address = *static_node_mgmt_addressFound
		return nil
	}
}

func testAccCheckAciStaticNodeMgmtAddressOoBExists(name string, static_node_mgmt_address *models.OutofbandStaticNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Static Node Mgmt Address %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Static Node Mgmt Address dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		static_node_mgmt_addressFound := models.OutofbandStaticNodeFromContainer(cont)
		if static_node_mgmt_addressFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Static Node Mgmt Address %s not found", rs.Primary.ID)
		}
		*static_node_mgmt_address = *static_node_mgmt_addressFound
		return nil
	}
}

func testAccCheckAciStaticNodeMgmtAddressDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing static_node_mgmt_address destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_static_node_mgmt_address" {
			cont, err := client.Get(rs.Primary.ID)
			static_node_mgmt_address := models.InbandStaticNodeFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Static Node Mgmt Address %s Still exists", static_node_mgmt_address.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciStaticNodeMgmtAddressIdEqual(m1, m2 *models.InbandStaticNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("static_node_mgmt_address DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciStaticNodeMgmtAddressIdNotEqual(m1, m2 *models.InbandStaticNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("static_node_mgmt_address DNs are equal")
		}
		return nil
	}
}

func testAccCheckAciStaticNodeMgmtAddressOoBIdNotEqual(m1 *models.InbandStaticNode, m2 *models.OutofbandStaticNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("static_node_mgmt_address DNs are equal")
		}
		return nil
	}
}

func CreateStaticNodeMgmtAddressWithoutRequired(mgmtInBName, addrType, tDn, attrName string) string {
	fmt.Println("=== STEP  Basic: testing static_node_mgmt_address creation without ", attrName)
	rBlock := `
	
	resource "aci_node_mgmt_epg" "test" {
		type = "in_band"
		name  = "%s"
	}
	
	`
	switch attrName {
	case "management_epg_dn":
		rBlock += `
	resource "aci_static_node_mgmt_address" "test" {
	#	management_epg_dn  = aci_node_mgmt_epg.test.id
		type = "%s"
		t_dn  = "%s"
	}
		`
	case "t_dn":
		rBlock += `
	resource "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		type = "%s"
	#	t_dn  = "%s"
	}
		`
	case "type":
		rBlock += `
	resource "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
	#	type = "%s"
		t_dn  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, mgmtInBName, addrType, tDn)
}

func CreateAccStaticNodeMgmtAddressConfigWithRequiredParams(mgmtInBName, addrType, tDn string) string {
	fmt.Printf("=== STEP  testing static_node_mgmt_address creation with parent resource name %s, type %s and tdn %s\n", mgmtInBName, addrType, tDn)
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		type = "%s"
		name  = "%s"
	}
	
	resource "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		type = "%s"
		t_dn  = "%s"
	}
	`, addrType, mgmtInBName, addrType, tDn)
	return resource
}

func CreateAccStaticNodeMgmtAddressConfig(mgmtInBName, addrType, tDn string) string {
	fmt.Println("=== STEP  testing static_node_mgmt_address creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		type = "%s"
		name  = "%s"
	}
	
	resource "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		type = "%s"
		t_dn  = "%s"
	}
	`, addrType, mgmtInBName, addrType, tDn)
	return resource
}

func CreateAccStaticNodeMgmtAddressWithInValidParentDn(rName, tDn string) string {
	fmt.Println("=== STEP  Negative Case: testing static_node_mgmt_address creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_tenant.test.id
		type = "%s"
		t_dn  = "%s"	
	}
	`, rName, "in_band", tDn)
	return resource
}

func CreateAccStaticNodeMgmtAddressConfigWithOptionalValues(mgmtInBName, addrType, tDn string) string {
	fmt.Println("=== STEP  Basic: testing static_node_mgmt_address creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		type = "in_band"
		name  = "%s"
	}
	
	resource "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = "${aci_node_mgmt_epg.test.id}"
		type = "%s"
		t_dn  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		addr = "1.0.0.1/24"
		gw = "1.0.0.1"
		// v6_addr = "2001:0db8:85a3:0000:0000:8a2e:0370:7334/64"
		// v6_gw = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
		
	}
	`, mgmtInBName, addrType, tDn)

	return resource
}

func CreateAccStaticNodeMgmtAddressRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing static_node_mgmt_address updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_static_node_mgmt_address" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		addr = ""
		gw = ""
		v6_addr = ""
		v6_gw = ""
		
	}
	`)

	return resource
}

func CreateAccStaticNodeMgmtAddressUpdatedAttr(mgmtInBName, addrType, tDn, attribute, value string) string {
	fmt.Printf("=== STEP  testing static_node_mgmt_address attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		type = "in_band"
		name  = "%s"
	}
	
	resource "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		type = "%s"
		t_dn  = "%s"
		%s = "%s"
	}
	`, mgmtInBName, addrType, tDn, attribute, value)
	return resource
}
