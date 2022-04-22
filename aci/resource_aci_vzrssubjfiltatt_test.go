package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSubjectFilter_Basic(t *testing.T) {
	var subject_filter models.SubjectFilter
	fv_tenant_name := acctest.RandString(5)
	vz_br_cp_name := acctest.RandString(5)
	vz_subj_name := acctest.RandString(5)
	vz_rs_subj_filt_att_name := acctest.RandString(5)
	description := "subject_filter created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubjectFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubjectFilterConfig_basic(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_subj_filt_att_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubjectFilterExists("aci_subject_filter.foosubject_filter", &subject_filter),
					testAccCheckAciSubjectFilterAttributes(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_subj_filt_att_name, description, &subject_filter),
				),
			},
		},
	})
}

func TestAccAciSubjectFilter_Update(t *testing.T) {
	var subject_filter models.SubjectFilter
	fv_tenant_name := acctest.RandString(5)
	vz_br_cp_name := acctest.RandString(5)
	vz_subj_name := acctest.RandString(5)
	vz_rs_subj_filt_att_name := acctest.RandString(5)
	description := "subject_filter created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubjectFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubjectFilterConfig_basic(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_subj_filt_att_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubjectFilterExists("aci_subject_filter.foosubject_filter", &subject_filter),
					testAccCheckAciSubjectFilterAttributes(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_subj_filt_att_name, description, &subject_filter),
				),
			},
			{
				Config: testAccCheckAciSubjectFilterConfig_basic(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_subj_filt_att_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubjectFilterExists("aci_subject_filter.foosubject_filter", &subject_filter),
					testAccCheckAciSubjectFilterAttributes(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_subj_filt_att_name, description, &subject_filter),
				),
			},
		},
	})
}

func testAccCheckAciSubjectFilterConfig_basic(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_subj_filt_att_name string) string {
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

	resource "aci_subject_filter" "foosubject_filter" {
		name 		= "%s"
		description = "subject_filter created while acceptance testing"
		contract_subject_dn = aci_contract_subject.foocontract_subject.id
	}

	`, fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_subj_filt_att_name)
}

func testAccCheckAciSubjectFilterExists(name string, subject_filter *models.SubjectFilter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Subject Filter %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subject Filter dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		subject_filterFound := models.SubjectFilterFromContainer(cont)
		if subject_filterFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Subject Filter %s not found", rs.Primary.ID)
		}
		*subject_filter = *subject_filterFound
		return nil
	}
}

func testAccCheckAciSubjectFilterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_subject_filter" {
			cont, err := client.Get(rs.Primary.ID)
			subject_filter := models.SubjectFilterFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Subject Filter %s Still exists", subject_filter.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSubjectFilterAttributes(fv_tenant_name, vz_br_cp_name, vz_subj_name, vz_rs_subj_filt_att_name, description string, subject_filter *models.SubjectFilter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vz_rs_subj_filt_att_name != GetMOName(subject_filter.DistinguishedName) {
			return fmt.Errorf("Bad vz_rs_subj_filt_att %s", GetMOName(subject_filter.DistinguishedName))
		}

		if vz_subj_name != GetMOName(GetParentDn(subject_filter.DistinguishedName)) {
			return fmt.Errorf(" Bad vz_subj %s", GetMOName(GetParentDn(subject_filter.DistinguishedName)))
		}
		if description != subject_filter.Description {
			return fmt.Errorf("Bad subject_filter Description %s", subject_filter.Description)
		}
		return nil
	}
}
