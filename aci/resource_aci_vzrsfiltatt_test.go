package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciFilter_Basic(t *testing.T) {
	var filter models.Filter
	fv_tenant_name := acctest.RandString(5)
	vz_br_cp_name := acctest.RandString(5)
	vz_subj_name := acctest.RandString(5)
	vz_rs_filt_att_name := acctest.RandString(5)
	description := "filter created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFilterConfig_basic(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_filt_att_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists("aci_filter.foofilter", &filter),
					testAccCheckAciFilterAttributes(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_filt_att_name, description, &filter),
				),
			},
		},
	})
}

func TestAccAciFilter_Update(t *testing.T) {
	var filter models.Filter
	fv_tenant_name := acctest.RandString(5)
	vz_br_cp_name := acctest.RandString(5)
	vz_subj_name := acctest.RandString(5)
	vz_rs_filt_att_name := acctest.RandString(5)
	description := "filter created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFilterConfig_basic(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_filt_att_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists("aci_filter.foofilter", &filter),
					testAccCheckAciFilterAttributes(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_filt_att_name, description, &filter),
				),
			},
			{
				Config: testAccCheckAciFilterConfig_basic(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_filt_att_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists("aci_filter.foofilter", &filter),
					testAccCheckAciFilterAttributes(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_filt_att_name, description, &filter),
				),
			},
		},
	})
}

func testAccCheckAciFilterConfig_basic(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_filt_att_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_contract" "foocontract" {
		name 		= "%s"
		description = "contract created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	resource "aci_contract_subject" "foocontract_subject" {
		name 		= "%s"
		description = "contract_subject created while acceptance testing"
		contract_dn = aci_contract.foocontract.id
	}

	resource "aci_filter" "foofilter" {
		name 		= "%s"
		description = "filter created while acceptance testing"
		contract_subject_dn = aci_contract_subject.foocontract_subject.id
	}

	`, fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_filt_att_name)
}

func testAccCheckAciFilterExists(name string, filter *models.Filter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Filter %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Filter dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		filterFound := models.FilterFromContainer(cont)
		if filterFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Filter %s not found", rs.Primary.ID)
		}
		*filter = *filterFound
		return nil
	}
}

func testAccCheckAciFilterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_filter" {
			cont, err := client.Get(rs.Primary.ID)
			filter := models.FilterFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Filter %s Still exists", filter.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFilterAttributes(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_filt_att_name, description string, filter *models.Filter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vz_rs_filt_att_name != GetMOName(filter.DistinguishedName) {
			return fmt.Errorf("Bad vz_rs_filt_att %s", GetMOName(filter.DistinguishedName))
		}

		if vz_subj_name != GetMOName(GetParentDn(filter.DistinguishedName)) {
			return fmt.Errorf(" Bad vz_subj %s", GetMOName(GetParentDn(filter.DistinguishedName)))
		}
		if description != filter.Description {
			return fmt.Errorf("Bad filter Description %s", filter.Description)
		}
		return nil
	}
}
