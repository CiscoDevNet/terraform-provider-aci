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

func TestAccAciBDDHCPLabel_Basic(t *testing.T) {
	var bd_dhcp_label_default models.BDDHCPLabel
	var bd_dhcp_label_updated models.BDDHCPLabel
	resourceName := "aci_bd_dhcp_label.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	fvBDName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBDDHCPLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBDDHCPLabelWithoutRequired(fvTenantName, fvBDName, rName, "bridge_domain_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBDDHCPLabelWithoutRequired(fvTenantName, fvBDName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBDDHCPLabelConfig(fvTenantName, fvBDName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBDDHCPLabelExists(resourceName, &bd_dhcp_label_default),
					resource.TestCheckResourceAttr(resourceName, "bridge_domain_dn", fmt.Sprintf("uni/tn-%s/BD-%s", fvTenantName, fvBDName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "owner", "infra"),
				),
			},
			{
				Config: CreateAccBDDHCPLabelConfigWithOptionalValues(fvTenantName, fvBDName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBDDHCPLabelExists(resourceName, &bd_dhcp_label_updated),
					resource.TestCheckResourceAttr(resourceName, "bridge_domain_dn", fmt.Sprintf("uni/tn-%s/BD-%s", fvTenantName, fvBDName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_bd_dhcp_label"),
					resource.TestCheckResourceAttr(resourceName, "owner", "tenant"),
					resource.TestCheckResourceAttr(resourceName, "tag", "alice-blue"),
					testAccCheckAciBDDHCPLabelIdEqual(&bd_dhcp_label_default, &bd_dhcp_label_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccBDDHCPLabelConfigUpdatedName(fvTenantName, fvBDName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccBDDHCPLabelRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBDDHCPLabelConfigWithRequiredParams(rName, rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBDDHCPLabelExists(resourceName, &bd_dhcp_label_updated),
					resource.TestCheckResourceAttr(resourceName, "bridge_domain_dn", fmt.Sprintf("uni/tn-%s/BD-%s", rName, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciBDDHCPLabelIdNotEqual(&bd_dhcp_label_default, &bd_dhcp_label_updated),
				),
			},
			{
				Config: CreateAccBDDHCPLabelConfig(fvTenantName, fvBDName, rName),
			},
			{
				Config: CreateAccBDDHCPLabelConfigWithRequiredParams(rName, rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBDDHCPLabelExists(resourceName, &bd_dhcp_label_updated),
					resource.TestCheckResourceAttr(resourceName, "bridge_domain_dn", fmt.Sprintf("uni/tn-%s/BD-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciBDDHCPLabelIdNotEqual(&bd_dhcp_label_default, &bd_dhcp_label_updated),
				),
			},
		},
	})
}

func TestAccAciBDDHCPLabel_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	fvBDName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBDDHCPLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBDDHCPLabelConfig(fvTenantName, fvBDName, rName),
			},
			{
				Config:      CreateAccBDDHCPLabelWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBDDHCPLabelUpdatedAttr(fvTenantName, fvBDName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBDDHCPLabelUpdatedAttr(fvTenantName, fvBDName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBDDHCPLabelUpdatedAttr(fvTenantName, fvBDName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccBDDHCPLabelUpdatedAttr(fvTenantName, fvBDName, rName, "owner", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccBDDHCPLabelUpdatedAttr(fvTenantName, fvBDName, rName, "tag", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccBDDHCPLabelUpdatedAttr(fvTenantName, fvBDName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccBDDHCPLabelConfig(fvTenantName, fvBDName, rName),
			},
		},
	})
}

func TestAccAciBDDHCPLabel_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	fvBDName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBDDHCPLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBDDHCPLabelConfigMultiple(fvTenantName, fvBDName, rName),
			},
		},
	})
}

func testAccCheckAciBDDHCPLabelExists(name string, bd_dhcp_label *models.BDDHCPLabel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Bd DHCP Label %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Bd DHCP Label dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bd_dhcp_labelFound := models.BDDHCPLabelFromContainer(cont)
		if bd_dhcp_labelFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Bd DHCP Label %s not found", rs.Primary.ID)
		}
		*bd_dhcp_label = *bd_dhcp_labelFound
		return nil
	}
}

func testAccCheckAciBDDHCPLabelDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing bd_dhcp_label destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_bd_dhcp_label" {
			cont, err := client.Get(rs.Primary.ID)
			bd_dhcp_label := models.BDDHCPLabelFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Bd DHCP Label %s Still exists", bd_dhcp_label.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBDDHCPLabelIdEqual(m1, m2 *models.BDDHCPLabel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("bd_dhcp_label DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciBDDHCPLabelIdNotEqual(m1, m2 *models.BDDHCPLabel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("bd_dhcp_label DNs are equal")
		}
		return nil
	}
}

func CreateBDDHCPLabelWithoutRequired(fvTenantName, fvBDName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bd_dhcp_label creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_bridge_domain" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	`
	switch attrName {
	case "bridge_domain_dn":
		rBlock += `
	resource "aci_bd_dhcp_label" "test" {
	#	bridge_domain_dn  = aci_bridge_domain.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, fvBDName, rName)
}

func CreateAccBDDHCPLabelConfigWithRequiredParams(fvTenantName, fvBDName, rName string) string {
	fmt.Println("=== STEP  testing bd_dhcp_label creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bridge_domain" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = "%s"
	}
	`, fvTenantName, fvBDName, rName)
	return resource
}
func CreateAccBDDHCPLabelConfigUpdatedName(fvTenantName, fvBDName, rName string) string {
	fmt.Println("=== STEP  testing bd_dhcp_label creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bridge_domain" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = "%s"
	}
	`, fvTenantName, fvBDName, rName)
	return resource
}

func CreateAccBDDHCPLabelConfig(fvTenantName, fvBDName, rName string) string {
	fmt.Println("=== STEP  testing bd_dhcp_label creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bridge_domain" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = "%s"
	}
	`, fvTenantName, fvBDName, rName)
	return resource
}

func CreateAccBDDHCPLabelConfigMultiple(fvTenantName, fvBDName, rName string) string {
	fmt.Println("=== STEP  testing multiple bd_dhcp_label creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bridge_domain" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, fvBDName, rName)
	return resource
}

func CreateAccBDDHCPLabelWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing bd_dhcp_label creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccBDDHCPLabelConfigWithOptionalValues(fvTenantName, fvBDName, rName string) string {
	fmt.Println("=== STEP  Basic: testing bd_dhcp_label creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bridge_domain" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = "${aci_bridge_domain.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_bd_dhcp_label"
		owner = "tenant"
		tag = "alice-blue"
		
	}
	`, fvTenantName, fvBDName, rName)

	return resource
}

func CreateAccBDDHCPLabelRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing bd_dhcp_label updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_bd_dhcp_label" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_bd_dhcp_label"
		owner = "tenant"
		tag = "alice-blue"
		
	}
	`)

	return resource
}

func CreateAccBDDHCPLabelUpdatedAttr(fvTenantName, fvBDName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing bd_dhcp_label attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bridge_domain" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, fvBDName, rName, attribute, value)
	return resource
}
