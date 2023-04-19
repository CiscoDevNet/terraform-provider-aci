package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciContractInterfaceRelationship_Basic(t *testing.T) {
	var contract_interface models.ContractInterfaceRelationship
	fv_tenant_name := acctest.RandString(5)
	fv_ap_name := acctest.RandString(5)
	fv_ae_pg_name := acctest.RandString(5)
	fv_rs_cons_if_name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractInterfaceRelationshipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciContractInterfaceRelationshipConfig_basic(fv_tenant_name, fv_ap_name, fv_ae_pg_name, fv_rs_cons_if_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractInterfaceRelationshipExists("aci_contract_interface.foocontract_interface", &contract_interface),
					testAccCheckAciContractInterfaceRelationshipAttributes(fv_tenant_name, fv_ap_name, fv_ae_pg_name, fv_rs_cons_if_name, &contract_interface),
				),
			},
		},
	})
}

func TestAccAciContractInterfaceRelationship_Update(t *testing.T) {
	var contract_interface models.ContractInterfaceRelationship
	fv_tenant_name := acctest.RandString(5)
	fv_ap_name := acctest.RandString(5)
	fv_ae_pg_name := acctest.RandString(5)
	fv_rs_cons_if_name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractInterfaceRelationshipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciContractInterfaceRelationshipConfig_basic(fv_tenant_name, fv_ap_name, fv_ae_pg_name, fv_rs_cons_if_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractInterfaceRelationshipExists("aci_contract_interface.foocontract_interface", &contract_interface),
					testAccCheckAciContractInterfaceRelationshipAttributes(fv_tenant_name, fv_ap_name, fv_ae_pg_name, fv_rs_cons_if_name, &contract_interface),
				),
			},
			{
				Config: testAccCheckAciContractInterfaceRelationshipConfig_basic(fv_tenant_name, fv_ap_name, fv_ae_pg_name, fv_rs_cons_if_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractInterfaceRelationshipExists("aci_contract_interface.foocontract_interface", &contract_interface),
					testAccCheckAciContractInterfaceRelationshipAttributes(fv_tenant_name, fv_ap_name, fv_ae_pg_name, fv_rs_cons_if_name, &contract_interface),
				),
			},
		},
	})
}

func testAccCheckAciContractInterfaceRelationshipConfig_basic(fv_tenant_name, fv_ap_name, fv_ae_pg_name, fv_rs_cons_if_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_application_profile" "fooapplication_profile" {
		name 		= "%s"
		description = "application_profile created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	resource "aci_application_epg" "fooapplication_epg" {
		name 		= "%s"
		description = "application_epg created while acceptance testing"
		application_profile_dn = aci_application_profile.fooapplication_profile.id
	}

	resource "aci_contract_interface" "foocontract_interface" {
		name 		= "%s"
		application_epg_dn = aci_application_epg.fooapplication_epg.id
	}

	`, fv_tenant_name, fv_ap_name, fv_ae_pg_name, fv_rs_cons_if_name)
}

func testAccCheckAciContractInterfaceRelationshipExists(name string, contract_interface *models.ContractInterfaceRelationship) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Contract Interface %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Contract Interface dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		contract_interfaceFound := models.ContractInterfaceRelationshipFromContainer(cont)
		if contract_interfaceFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Contract Interface %s not found", rs.Primary.ID)
		}
		*contract_interface = *contract_interfaceFound
		return nil
	}
}

func testAccCheckAciContractInterfaceRelationshipDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_contract_interface" {
			cont, err := client.Get(rs.Primary.ID)
			contract_interface := models.ContractInterfaceRelationshipFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Contract Interface %s Still exists", contract_interface.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciContractInterfaceRelationshipAttributes(fv_tenant_name, fv_ap_name, fv_ae_pg_name, fv_rs_cons_if_name string, contract_interface *models.ContractInterfaceRelationship) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if fv_rs_cons_if_name != GetMOName(contract_interface.DistinguishedName) {
			return fmt.Errorf("Bad fv_rs_cons_if %s", GetMOName(contract_interface.DistinguishedName))
		}

		if fv_ae_pg_name != GetMOName(GetParentDn(contract_interface.DistinguishedName)) {
			return fmt.Errorf(" Bad fvae_pg %s", GetMOName(GetParentDn(contract_interface.DistinguishedName)))
		}
		return nil
	}
}
