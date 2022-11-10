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

func TestAccAciPBRBackupPolicy_Basic(t *testing.T) {
	var pbr_backup_policy models.PBRBackupPolicy
	fv_tenant_name := acctest.RandString(5)
	vns_backup_pol_name := acctest.RandString(5)
	description := "pbr_backup_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPBRBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPBRBackupPolicyConfig_basic(fv_tenant_name, vns_backup_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPBRBackupPolicyExists("aci_service_redirect_backup_policy.foopbr_backup_policy", &pbr_backup_policy),
					testAccCheckAciPBRBackupPolicyAttributes(fv_tenant_name, vns_backup_pol_name, description, &pbr_backup_policy),
				),
			},
		},
	})
}

func TestAccAciPBRBackupPolicy_Update(t *testing.T) {
	var pbr_backup_policy models.PBRBackupPolicy
	fv_tenant_name := acctest.RandString(5)
	vns_backup_pol_name := acctest.RandString(5)
	description := "pbr_backup_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPBRBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPBRBackupPolicyConfig_basic(fv_tenant_name, vns_backup_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPBRBackupPolicyExists("aci_service_redirect_backup_policy.foopbr_backup_policy", &pbr_backup_policy),
					testAccCheckAciPBRBackupPolicyAttributes(fv_tenant_name, vns_backup_pol_name, description, &pbr_backup_policy),
				),
			},
			{
				Config: testAccCheckAciPBRBackupPolicyConfig_basic(fv_tenant_name, vns_backup_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPBRBackupPolicyExists("aci_service_redirect_backup_policy.foopbr_backup_policy", &pbr_backup_policy),
					testAccCheckAciPBRBackupPolicyAttributes(fv_tenant_name, vns_backup_pol_name, description, &pbr_backup_policy),
				),
			},
		},
	})
}

func testAccCheckAciPBRBackupPolicyConfig_basic(fv_tenant_name, vns_backup_pol_name string) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "footenant" {
		name        = "%s"
		description = "tenant created while acceptance testing"
	}

	resource "aci_service_redirect_backup_policy" "foopbr_backup_policy" {
		name        = "%s"
		description = "pbr_backup_policy created while acceptance testing"
		tenant_dn   = aci_tenant.footenant.id
	}
	`, fv_tenant_name, vns_backup_pol_name)
}

func testAccCheckAciPBRBackupPolicyExists(name string, pbr_backup_policy *models.PBRBackupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("PBR Backup Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PBR Backup Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		pbr_backup_policyFound := models.PBRBackupPolicyFromContainer(cont)
		if pbr_backup_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("PBR Backup Policy %s not found", rs.Primary.ID)
		}
		*pbr_backup_policy = *pbr_backup_policyFound
		return nil
	}
}

func testAccCheckAciPBRBackupPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_service_redirect_backup_policy" {
			cont, err := client.Get(rs.Primary.ID)
			pbr_backup_policy := models.PBRBackupPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("PBR Backup Policy %s Still exists", pbr_backup_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciPBRBackupPolicyAttributes(fv_tenant_name, vns_backup_pol_name, description string, pbr_backup_policy *models.PBRBackupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vns_backup_pol_name != GetMOName(pbr_backup_policy.DistinguishedName) {
			return fmt.Errorf("Bad vns_backup_pol %s", GetMOName(pbr_backup_policy.DistinguishedName))
		}
		if description != pbr_backup_policy.Description {
			return fmt.Errorf("Bad pbr_backup_policy Description %s", pbr_backup_policy.Description)
		}
		return nil
	}
}
