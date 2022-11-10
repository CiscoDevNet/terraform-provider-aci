package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciContractSubject_Basic(t *testing.T) {
	var contract_subject models.ContractSubject
	description := "contract_subject created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractSubjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciContractSubjectConfig_basic(description, "AtleastOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists("aci_contract_subject.foocontract_subject", &contract_subject),
					testAccCheckAciContractSubjectAttributes(description, "AtleastOne", &contract_subject),
				),
			},
		},
	})
}

func TestAccAciContractSubject_update(t *testing.T) {
	var contract_subject models.ContractSubject
	description := "contract_subject created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractSubjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciContractSubjectConfig_basic(description, "AtleastOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists("aci_contract_subject.foocontract_subject", &contract_subject),
					testAccCheckAciContractSubjectAttributes(description, "AtleastOne", &contract_subject),
				),
			},
			{
				Config: testAccCheckAciContractSubjectConfig_basic(description, "All"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists("aci_contract_subject.foocontract_subject", &contract_subject),
					testAccCheckAciContractSubjectAttributes(description, "All", &contract_subject),
				),
			},
		},
	})
}

func testAccCheckAciContractSubjectConfig_basic(description, cons_match_t string) string {
	return fmt.Sprintf(`

	resource "aci_contract_subject" "foocontract_subject" {
		contract_dn   = "uni/tn-test_rutvik_tenant/brc-demo_contract"
		description   = "%s"
		name          = "demo_subject"
		annotation    = "tag_subject"
		cons_match_t  = "%s"
		name_alias    = "alias_subject"
		prio          = "level1"
		prov_match_t  = "AtleastOne"
		rev_flt_ports = "yes"
		target_dscp   = "CS0"
	}
	  
	`, description, cons_match_t)
}

func testAccCheckAciContractSubjectExists(name string, contract_subject *models.ContractSubject) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Contract Subject %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Contract Subject dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		contract_subjectFound := models.ContractSubjectFromContainer(cont)
		if contract_subjectFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Contract Subject %s not found", rs.Primary.ID)
		}
		*contract_subject = *contract_subjectFound
		return nil
	}
}

func testAccCheckAciContractSubjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_contract_subject" {
			cont, err := client.Get(rs.Primary.ID)
			contract_subject := models.ContractSubjectFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Contract Subject %s Still exists", contract_subject.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciContractSubjectAttributes(description, cons_match_t string, contract_subject *models.ContractSubject) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != contract_subject.Description {
			return fmt.Errorf("Bad contract_subject Description %s", contract_subject.Description)
		}

		if "demo_subject" != contract_subject.Name {
			return fmt.Errorf("Bad contract_subject name %s", contract_subject.Name)
		}

		if "tag_subject" != contract_subject.Annotation {
			return fmt.Errorf("Bad contract_subject annotation %s", contract_subject.Annotation)
		}

		if cons_match_t != contract_subject.ConsMatchT {
			return fmt.Errorf("Bad contract_subject cons_match_t %s", contract_subject.ConsMatchT)
		}

		if "alias_subject" != contract_subject.NameAlias {
			return fmt.Errorf("Bad contract_subject name_alias %s", contract_subject.NameAlias)
		}

		if "level1" != contract_subject.Prio {
			return fmt.Errorf("Bad contract_subject prio %s", contract_subject.Prio)
		}

		if "AtleastOne" != contract_subject.ProvMatchT {
			return fmt.Errorf("Bad contract_subject prov_match_t %s", contract_subject.ProvMatchT)
		}

		if "yes" != contract_subject.RevFltPorts {
			return fmt.Errorf("Bad contract_subject rev_flt_ports %s", contract_subject.RevFltPorts)
		}

		if "CS0" != contract_subject.TargetDscp {
			return fmt.Errorf("Bad contract_subject target_dscp %s", contract_subject.TargetDscp)
		}

		return nil
	}
}
