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

func TestAccAciAccessCredentialtomanagethecloudresources_Basic(t *testing.T) {
	var access_credentialtomanagethecloudresources models.CloudCredentials
	fv_tenant_name := acctest.RandString(5)
	cloud_credentials_name := acctest.RandString(5)
	description := "access_credentialtomanagethecloudresources created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessCredentialtomanagethecloudresourcesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessCredentialtomanagethecloudresourcesConfig_basic(fv_tenant_name, cloud_credentials_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessCredentialtomanagethecloudresourcesExists("aci_cloud_credentials.fooaccess_credentialtomanagethecloudresources", &access_credentialtomanagethecloudresources),
					testAccCheckAciAccessCredentialtomanagethecloudresourcesAttributes(fv_tenant_name, cloud_credentials_name, description, &access_credentialtomanagethecloudresources),
				),
			},
		},
	})
}

func TestAccAciAccessCredentialtomanagethecloudresources_Update(t *testing.T) {
	var access_credentialtomanagethecloudresources models.CloudCredentials
	fv_tenant_name := acctest.RandString(5)
	cloud_credentials_name := acctest.RandString(5)
	description := "access_credentialtomanagethecloudresources created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessCredentialtomanagethecloudresourcesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessCredentialtomanagethecloudresourcesConfig_basic(fv_tenant_name, cloud_credentials_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessCredentialtomanagethecloudresourcesExists("aci_cloud_credentials.fooaccess_credentialtomanagethecloudresources", &access_credentialtomanagethecloudresources),
					testAccCheckAciAccessCredentialtomanagethecloudresourcesAttributes(fv_tenant_name, cloud_credentials_name, description, &access_credentialtomanagethecloudresources),
				),
			},
			{
				Config: testAccCheckAciAccessCredentialtomanagethecloudresourcesConfig_basic(fv_tenant_name, cloud_credentials_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessCredentialtomanagethecloudresourcesExists("aci_cloud_credentials.fooaccess_credentialtomanagethecloudresources", &access_credentialtomanagethecloudresources),
					testAccCheckAciAccessCredentialtomanagethecloudresourcesAttributes(fv_tenant_name, cloud_credentials_name, description, &access_credentialtomanagethecloudresources),
				),
			},
		},
	})
}

func testAccCheckAciAccessCredentialtomanagethecloudresourcesConfig_basic(fv_tenant_name, cloud_credentials_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_cloud_credentials" "fooaccess_credentialtomanagethecloudresources" {
		name 		= "%s"
		description = "access_credentialtomanagethecloudresources created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	`, fv_tenant_name, cloud_credentials_name)
}

func testAccCheckAciAccessCredentialtomanagethecloudresourcesExists(name string, access_credentialtomanagethecloudresources *models.CloudCredentials) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Credential to manage the cloud resources %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Credential to manage the cloud resources dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		access_credentialtomanagethecloudresourcesFound := models.CloudCredentialsFromContainer(cont)
		if access_credentialtomanagethecloudresourcesFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Credential to manage the cloud resources %s not found", rs.Primary.ID)
		}
		*access_credentialtomanagethecloudresources = *access_credentialtomanagethecloudresourcesFound
		return nil
	}
}

func testAccCheckAciAccessCredentialtomanagethecloudresourcesDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_credentials" {
			cont, err := client.Get(rs.Primary.ID)
			access_credentialtomanagethecloudresources := models.CloudCredentialsFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Credential to manage the cloud resources %s Still exists", access_credentialtomanagethecloudresources.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciAccessCredentialtomanagethecloudresourcesAttributes(fv_tenant_name, cloud_credentials_name, description string, access_credentialtomanagethecloudresources *models.CloudCredentials) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cloud_credentials_name != GetMOName(access_credentialtomanagethecloudresources.DistinguishedName) {
			return fmt.Errorf("Bad cloud_credentials %s", GetMOName(access_credentialtomanagethecloudresources.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(access_credentialtomanagethecloudresources.DistinguishedName)) {
			return fmt.Errorf(" Bad fv_tenant %s", GetMOName(GetParentDn(access_credentialtomanagethecloudresources.DistinguishedName)))
		}
		if description != access_credentialtomanagethecloudresources.Description {
			return fmt.Errorf("Bad access_credentialtomanagethecloudresources Description %s", access_credentialtomanagethecloudresources.Description)
		}
		return nil
	}
}
