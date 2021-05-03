package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciContract_Basic(t *testing.T) {
	var contract models.Contract
	description := "contract created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciContractConfig_basic(description, "tenant"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists("aci_contract.foocontract", &contract),
					testAccCheckAciContractAttributes(description, "tenant", &contract),
				),
			},
		},
	})
}

func TestAccAciContract_update(t *testing.T) {
	var contract models.Contract
	description := "contract created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciContractConfig_basic(description, "tenant"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists("aci_contract.foocontract", &contract),
					testAccCheckAciContractAttributes(description, "tenant", &contract),
				),
			},
			{
				Config: testAccCheckAciContractConfig_basic(description, "global"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists("aci_contract.foocontract", &contract),
					testAccCheckAciContractAttributes(description, "global", &contract),
				),
			},
		},
	})
}

func testAccCheckAciContractConfig_basic(description, scope string) string {
	return fmt.Sprintf(`

	resource "aci_contract" "foocontract" {
		tenant_dn   = "${aci_tenant.example.id}"
		description = "%s"
		name        = "demo_contract"
		annotation  = "tag_contract"
		name_alias  = "alias_contract"
		prio        = "level1"
		scope       = "%s"
		target_dscp = "unspecified"
	}
	  
	`, description, scope)
}

func testAccCheckAciContractExists(name string, contract *models.Contract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Contract %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Contract dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		contractFound := models.ContractFromContainer(cont)
		if contractFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Contract %s not found", rs.Primary.ID)
		}
		*contract = *contractFound
		return nil
	}
}

func testAccCheckAciContractDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_contract" {
			cont, err := client.Get(rs.Primary.ID)
			contract := models.ContractFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Contract %s Still exists", contract.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciContractAttributes(description, scope string, contract *models.Contract) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != contract.Description {
			return fmt.Errorf("Bad contract Description %s", contract.Description)
		}

		if "demo_contract" != contract.Name {
			return fmt.Errorf("Bad contract name %s", contract.Name)
		}

		if "tag_contract" != contract.Annotation {
			return fmt.Errorf("Bad contract annotation %s", contract.Annotation)
		}

		if "alias_contract" != contract.NameAlias {
			return fmt.Errorf("Bad contract name_alias %s", contract.NameAlias)
		}

		if "level1" != contract.Prio {
			return fmt.Errorf("Bad contract prio %s", contract.Prio)
		}

		if scope != contract.Scope {
			return fmt.Errorf("Bad contract scope %s", contract.Scope)
		}

		if "unspecified" != contract.TargetDscp {
			return fmt.Errorf("Bad contract target_dscp %s", contract.TargetDscp)
		}

		return nil
	}
}
