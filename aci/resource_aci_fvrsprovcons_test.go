package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciContractProvider_Basic(t *testing.T) {
	var contract_provider models.ContractProvider
	description := "contract_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciContractProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractProviderExists("aci_contract_provider.foocontract_provider", &contract_provider),
					testAccCheckAciContractProviderAttributes(description, &contract_provider),
				),
			},
			{
				ResourceName:      "aci_contract_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciContractProvider_update(t *testing.T) {
	var contract_provider models.ContractProvider
	description := "contract_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciContractProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractProviderExists("aci_contract_provider.foocontract_provider", &contract_provider),
					testAccCheckAciContractProviderAttributes(description, &contract_provider),
				),
			},
			{
				Config: testAccCheckAciContractProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractProviderExists("aci_contract_provider.foocontract_provider", &contract_provider),
					testAccCheckAciContractProviderAttributes(description, &contract_provider),
				),
			},
		},
	})
}

func testAccCheckAciContractProviderConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_contract_provider" "foocontract_provider" {
		application_epg_dn = aci_application_epg.example.id
		description        = "%s"
		tnVzBrCPName       = "example"
		annotation         = "example"
		match_t            = "All"
		prio               = "unspecified"
		tn_vz_br_cp_name   = "example"
	}
	`, description)
}

func testAccCheckAciContractProviderExists(name string, contract_provider *models.ContractProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Contract Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Contract Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		contract_providerFound := models.ContractProviderFromContainer(cont)
		if contract_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Contract Provider %s not found", rs.Primary.ID)
		}
		*contract_provider = *contract_providerFound
		return nil
	}
}

func testAccCheckAciContractProviderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_contract_provider" {
			cont, err := client.Get(rs.Primary.ID)
			contract_provider := models.ContractProviderFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Contract Provider %s Still exists", contract_provider.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciContractProviderAttributes(description string, contract_provider *models.ContractProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != contract_provider.Description {
			return fmt.Errorf("Bad contract_provider Description %s", contract_provider.Description)
		}

		if "example" != contract_provider.TnVzBrCPName {
			return fmt.Errorf("Bad contract_provider tn_vz_br_cp_name %s", contract_provider.TnVzBrCPName)
		}

		if "example" != contract_provider.Annotation {
			return fmt.Errorf("Bad contract_provider annotation %s", contract_provider.Annotation)
		}

		if "All" != contract_provider.MatchT {
			return fmt.Errorf("Bad contract_provider match_t %s", contract_provider.MatchT)
		}

		if "unspecified" != contract_provider.Prio {
			return fmt.Errorf("Bad contract_provider prio %s", contract_provider.Prio)
		}

		if "example" != contract_provider.TnVzBrCPName {
			return fmt.Errorf("Bad contract_provider tn_vz_br_cp_name %s", contract_provider.TnVzBrCPName)
		}

		return nil
	}
}
