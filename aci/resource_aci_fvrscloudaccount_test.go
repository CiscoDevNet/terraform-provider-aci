package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciTenanttoaccountassociation_Basic(t *testing.T) {
	var tenanttoaccountassociation models.Tenanttoaccountassociation
	fv_tenant_name := acctest.RandString(5)
	fv_rs_cloud_account_name := acctest.RandString(5)
	description := "tenanttoaccountassociation created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTenanttoaccountassociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTenanttoaccountassociationConfig_basic(fv_tenant_name, fv_rs_cloud_account_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTenanttoaccountassociationExists("aci_tenanttoaccountassociation.footenanttoaccountassociation", &tenanttoaccountassociation),
					testAccCheckAciTenanttoaccountassociationAttributes(fv_tenant_name, fv_rs_cloud_account_name, description, &tenanttoaccountassociation),
				),
			},
		},
	})
}

func TestAccAciTenanttoaccountassociation_Update(t *testing.T) {
	var tenanttoaccountassociation models.Tenanttoaccountassociation
	fv_tenant_name := acctest.RandString(5)
	fv_rs_cloud_account_name := acctest.RandString(5)
	description := "tenanttoaccountassociation created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTenanttoaccountassociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTenanttoaccountassociationConfig_basic(fv_tenant_name, fv_rs_cloud_account_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTenanttoaccountassociationExists("aci_tenanttoaccountassociation.footenanttoaccountassociation", &tenanttoaccountassociation),
					testAccCheckAciTenanttoaccountassociationAttributes(fv_tenant_name, fv_rs_cloud_account_name, description, &tenanttoaccountassociation),
				),
			},
			{
				Config: testAccCheckAciTenanttoaccountassociationConfig_basic(fv_tenant_name, fv_rs_cloud_account_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTenanttoaccountassociationExists("aci_tenanttoaccountassociation.footenanttoaccountassociation", &tenanttoaccountassociation),
					testAccCheckAciTenanttoaccountassociationAttributes(fv_tenant_name, fv_rs_cloud_account_name, description, &tenanttoaccountassociation),
				),
			},
		},
	})
}

func testAccCheckAciTenanttoaccountassociationConfig_basic(fv_tenant_name, fv_rs_cloud_account_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_tenanttoaccountassociation" "footenanttoaccountassociation" {
		name 		= "%s"
		description = "tenanttoaccountassociation created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	`, fv_tenant_name, fv_rs_cloud_account_name)
}

func testAccCheckAciTenanttoaccountassociationExists(name string, tenanttoaccountassociation *models.Tenanttoaccountassociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Tenant to account association %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Tenant to account association dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tenanttoaccountassociationFound := models.TenanttoaccountassociationFromContainer(cont)
		if tenanttoaccountassociationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Tenant to account association %s not found", rs.Primary.ID)
		}
		*tenanttoaccountassociation = *tenanttoaccountassociationFound
		return nil
	}
}

func testAccCheckAciTenanttoaccountassociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tenanttoaccountassociation" {
			cont, err := client.Get(rs.Primary.ID)
			tenanttoaccountassociation := models.TenanttoaccountassociationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Tenant to account association %s Still exists", tenanttoaccountassociation.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTenanttoaccountassociationAttributes(fv_tenant_name, fv_rs_cloud_account_name, description string, tenanttoaccountassociation *models.Tenanttoaccountassociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if fv_rs_cloud_account_name != GetMOName(tenanttoaccountassociation.DistinguishedName) {
			return fmt.Errorf("Bad fv_rs_cloud_account %s", GetMOName(tenanttoaccountassociation.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(tenanttoaccountassociation.DistinguishedName)) {
			return fmt.Errorf(" Bad fv_tenant %s", GetMOName(GetParentDn(tenanttoaccountassociation.DistinguishedName)))
		}
		if description != tenanttoaccountassociation.Description {
			return fmt.Errorf("Bad tenanttoaccountassociation Description %s", tenanttoaccountassociation.Description)
		}
		return nil
	}
}
