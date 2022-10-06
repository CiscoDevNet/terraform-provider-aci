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

func TestAccAciCloudAccount_Basic(t *testing.T) {
	var account models.Account
	fv_tenant_name := acctest.RandString(5)
	cloud_account_name := acctest.RandString(5)
	description := "account created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudAccountConfig_basic(fv_tenant_name, cloud_account_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudAccountExists("aci_account.fooaccount", &account),
					testAccCheckAciCloudAccountAttributes(fv_tenant_name, cloud_account_name, description, &account),
				),
			},
		},
	})
}

func TestAccAciCloudAccount_Update(t *testing.T) {
	var account models.Account
	fv_tenant_name := acctest.RandString(5)
	cloud_account_name := acctest.RandString(5)
	description := "account created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudAccountConfig_basic(fv_tenant_name, cloud_account_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudAccountExists("aci_account.fooaccount", &account),
					testAccCheckAciCloudAccountAttributes(fv_tenant_name, cloud_account_name, description, &account),
				),
			},
			{
				Config: testAccCheckAciCloudAccountConfig_basic(fv_tenant_name, cloud_account_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudAccountExists("aci_account.fooaccount", &account),
					testAccCheckAciCloudAccountAttributes(fv_tenant_name, cloud_account_name, description, &account),
				),
			},
		},
	})
}

func testAccCheckAciCloudAccountConfig_basic(fv_tenant_name, cloud_account_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_account" "fooaccount" {
		name 		= "%s"
		description = "account created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	`, fv_tenant_name, cloud_account_name)
}

func testAccCheckAciCloudAccountExists(name string, account *models.Account) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Account %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Account dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		accountFound := models.AccountFromContainer(cont)
		if accountFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Account %s not found", rs.Primary.ID)
		}
		*account = *accountFound
		return nil
	}
}

func testAccCheckAciCloudAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_account" {
			cont, err := client.Get(rs.Primary.ID)
			account := models.AccountFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Account %s Still exists", account.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudAccountAttributes(fv_tenant_name, cloud_account_name, description string, account *models.Account) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cloud_account_name != GetMOName(account.DistinguishedName) {
			return fmt.Errorf("Bad cloud_account %s", GetMOName(account.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(account.DistinguishedName)) {
			return fmt.Errorf(" Bad fv_tenant %s", GetMOName(GetParentDn(account.DistinguishedName)))
		}
		if description != account.Description {
			return fmt.Errorf("Bad account Description %s", account.Description)
		}
		return nil
	}
}
