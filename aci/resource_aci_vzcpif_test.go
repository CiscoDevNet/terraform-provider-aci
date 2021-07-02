package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciImportedContract_Basic(t *testing.T) {
	var imported_contract models.ImportedContract
	description := "imported_contract created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciImportedContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciImportedContractConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciImportedContractExists("aci_imported_contract.fooimported_contract", &imported_contract),
					testAccCheckAciImportedContractAttributes(description, &imported_contract),
				),
			},
		},
	})
}

func TestAccAciImportedContract_update(t *testing.T) {
	var imported_contract models.ImportedContract
	description := "imported_contract created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciImportedContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciImportedContractConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciImportedContractExists("aci_imported_contract.fooimported_contract", &imported_contract),
					testAccCheckAciImportedContractAttributes(description, &imported_contract),
				),
			},
			{
				Config: testAccCheckAciImportedContractConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciImportedContractExists("aci_imported_contract.fooimported_contract", &imported_contract),
					testAccCheckAciImportedContractAttributes(description, &imported_contract),
				),
			},
		},
	})
}

func testAccCheckAciImportedContractConfig_basic(description string) string {
	return fmt.Sprintf(`

	

	resource "aci_imported_contract" "fooimported_contract" {
		  tenant_dn  = "${aci_tenant.example.id}"
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciImportedContractExists(name string, imported_contract *models.ImportedContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Imported Contract %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Imported Contract dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		imported_contractFound := models.ImportedContractFromContainer(cont)
		if imported_contractFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Imported Contract %s not found", rs.Primary.ID)
		}
		*imported_contract = *imported_contractFound
		return nil
	}
}

func testAccCheckAciImportedContractDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_imported_contract" {
			cont, err := client.Get(rs.Primary.ID)
			imported_contract := models.ImportedContractFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Imported Contract %s Still exists", imported_contract.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciImportedContractAttributes(description string, imported_contract *models.ImportedContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != imported_contract.Description {
			return fmt.Errorf("Bad imported_contract Description %s", imported_contract.Description)
		}

		if "example" != imported_contract.Name {
			return fmt.Errorf("Bad imported_contract name %s", imported_contract.Name)
		}

		if "example" != imported_contract.Annotation {
			return fmt.Errorf("Bad imported_contract annotation %s", imported_contract.Annotation)
		}

		if "example" != imported_contract.NameAlias {
			return fmt.Errorf("Bad imported_contract name_alias %s", imported_contract.NameAlias)
		}

		return nil
	}
}
