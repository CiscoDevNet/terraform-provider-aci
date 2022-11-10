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

func TestAccAciActiveDirectory_Basic(t *testing.T) {
	var active_directory models.CloudActiveDirectory
	fv_tenant_name := acctest.RandString(5)
	cloud_ad_name := acctest.RandString(5)
	description := "active_directory created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciActiveDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciActiveDirectoryConfig_basic(fv_tenant_name, cloud_ad_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActiveDirectoryExists("aci_cloud_ad.fooactive_directory", &active_directory),
					testAccCheckAciActiveDirectoryAttributes(fv_tenant_name, cloud_ad_name, description, &active_directory),
				),
			},
		},
	})
}

func TestAccAciActiveDirectory_Update(t *testing.T) {
	var active_directory models.CloudActiveDirectory
	fv_tenant_name := acctest.RandString(5)
	cloud_ad_name := acctest.RandString(5)
	description := "active_directory created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciActiveDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciActiveDirectoryConfig_basic(fv_tenant_name, cloud_ad_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActiveDirectoryExists("aci_cloud_ad.fooactive_directory", &active_directory),
					testAccCheckAciActiveDirectoryAttributes(fv_tenant_name, cloud_ad_name, description, &active_directory),
				),
			},
			{
				Config: testAccCheckAciActiveDirectoryConfig_basic(fv_tenant_name, cloud_ad_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActiveDirectoryExists("aci_cloud_ad.fooactive_directory", &active_directory),
					testAccCheckAciActiveDirectoryAttributes(fv_tenant_name, cloud_ad_name, description, &active_directory),
				),
			},
		},
	})
}

func testAccCheckAciActiveDirectoryConfig_basic(fv_tenant_name, cloud_ad_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_cloud_ad" "fooactive_directory" {
		name 		= "%s"
		description = "active_directory created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	`, fv_tenant_name, cloud_ad_name)
}

func testAccCheckAciActiveDirectoryExists(name string, active_directory *models.CloudActiveDirectory) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Active Directory %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Active Directory dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		active_directoryFound := models.CloudActiveDirectoryFromContainer(cont)
		if active_directoryFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Active Directory %s not found", rs.Primary.ID)
		}
		*active_directory = *active_directoryFound
		return nil
	}
}

func testAccCheckAciActiveDirectoryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_ad" {
			cont, err := client.Get(rs.Primary.ID)
			active_directory := models.CloudActiveDirectoryFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Active Directory %s Still exists", active_directory.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciActiveDirectoryAttributes(fv_tenant_name, cloud_ad_name, description string, active_directory *models.CloudActiveDirectory) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cloud_ad_name != GetMOName(active_directory.DistinguishedName) {
			return fmt.Errorf("Bad cloudad %s", GetMOName(active_directory.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(active_directory.DistinguishedName)) {
			return fmt.Errorf(" Bad fv_tenant %s", GetMOName(GetParentDn(active_directory.DistinguishedName)))
		}
		if description != active_directory.Description {
			return fmt.Errorf("Bad active_directory Description %s", active_directory.Description)
		}
		return nil
	}
}
