package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciTabooContract_Basic(t *testing.T) {
	var taboo_contract models.TabooContract
	description := "taboo_contract created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTabooContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTabooContractConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTabooContractExists("aci_taboo_contract.footaboo_contract", &taboo_contract),
					testAccCheckAciTabooContractAttributes(description, &taboo_contract),
				),
			},
		},
	})
}

func TestAccAciTabooContract_update(t *testing.T) {
	var taboo_contract models.TabooContract
	description := "taboo_contract created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTabooContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTabooContractConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTabooContractExists("aci_taboo_contract.footaboo_contract", &taboo_contract),
					testAccCheckAciTabooContractAttributes(description, &taboo_contract),
				),
			},
			{
				Config: testAccCheckAciTabooContractConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTabooContractExists("aci_taboo_contract.footaboo_contract", &taboo_contract),
					testAccCheckAciTabooContractAttributes(description, &taboo_contract),
				),
			},
		},
	})
}

func testAccCheckAciTabooContractConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "example"{
		description = "Tenant created while acceptance testing"
		name        = "demo_tenant"
	}

	resource "aci_taboo_contract" "footaboo_contract" {
		  tenant_dn  = aci_tenant.example.id
		  description = "%s"
		
		  name        = "example"
		  annotation  = "example"
		  name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciTabooContractExists(name string, taboo_contract *models.TabooContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Taboo Contract %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Taboo Contract dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		taboo_contractFound := models.TabooContractFromContainer(cont)
		if taboo_contractFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Taboo Contract %s not found", rs.Primary.ID)
		}
		*taboo_contract = *taboo_contractFound
		return nil
	}
}

func testAccCheckAciTabooContractDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_taboo_contract" {
			cont, err := client.Get(rs.Primary.ID)
			taboo_contract := models.TabooContractFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Taboo Contract %s Still exists", taboo_contract.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciTabooContractAttributes(description string, taboo_contract *models.TabooContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != taboo_contract.Description {
			return fmt.Errorf("Bad taboo_contract Description %s", taboo_contract.Description)
		}

		if "example" != taboo_contract.Name {
			return fmt.Errorf("Bad taboo_contract name %s", taboo_contract.Name)
		}

		if "example" != taboo_contract.Annotation {
			return fmt.Errorf("Bad taboo_contract annotation %s", taboo_contract.Annotation)
		}

		if "example" != taboo_contract.NameAlias {
			return fmt.Errorf("Bad taboo_contract name_alias %s", taboo_contract.NameAlias)
		}

		return nil
	}
}
