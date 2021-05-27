package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciCloudApplicationcontainer_Basic(t *testing.T) {
	var cloud_applicationcontainer models.CloudApplicationcontainer
	description := "cloud_applicationcontainer created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudApplicationcontainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudApplicationcontainerConfig_basic(description, "alias_app"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudApplicationcontainerExists("aci_cloud_applicationcontainer.foocloud_applicationcontainer", &cloud_applicationcontainer),
					testAccCheckAciCloudApplicationcontainerAttributes(description, "alias_app", &cloud_applicationcontainer),
				),
			},
		},
	})
}

func TestAccAciCloudApplicationcontainer_update(t *testing.T) {
	var cloud_applicationcontainer models.CloudApplicationcontainer
	description := "cloud_applicationcontainer created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudApplicationcontainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudApplicationcontainerConfig_basic(description, "alias_app"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudApplicationcontainerExists("aci_cloud_applicationcontainer.foocloud_applicationcontainer", &cloud_applicationcontainer),
					testAccCheckAciCloudApplicationcontainerAttributes(description, "alias_app", &cloud_applicationcontainer),
				),
			},
			{
				Config: testAccCheckAciCloudApplicationcontainerConfig_basic(description, "updated_app"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudApplicationcontainerExists("aci_cloud_applicationcontainer.foocloud_applicationcontainer", &cloud_applicationcontainer),
					testAccCheckAciCloudApplicationcontainerAttributes(description, "updated_app", &cloud_applicationcontainer),
				),
			},
		},
	})
}

func testAccCheckAciCloudApplicationcontainerConfig_basic(description, name_alias string) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "footenant" {
		description = "Tenant created while acceptance testing"
		name        = "demo_tenant"
	}

	resource "aci_cloud_applicationcontainer" "foocloud_applicationcontainer" {
		tenant_dn   = aci_tenant.footenant.id
		description = "%s"
		name        = "demo_app"
		annotation  = "tag_app"
		name_alias  = "%s"
	}
	  
	`, description, name_alias)
}

func testAccCheckAciCloudApplicationcontainerExists(name string, cloud_applicationcontainer *models.CloudApplicationcontainer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Application container %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Application container dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_applicationcontainerFound := models.CloudApplicationcontainerFromContainer(cont)
		if cloud_applicationcontainerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Application container %s not found", rs.Primary.ID)
		}
		*cloud_applicationcontainer = *cloud_applicationcontainerFound
		return nil
	}
}

func testAccCheckAciCloudApplicationcontainerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_applicationcontainer" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_applicationcontainer := models.CloudApplicationcontainerFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Application container %s Still exists", cloud_applicationcontainer.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudApplicationcontainerAttributes(description, name_alias string, cloud_applicationcontainer *models.CloudApplicationcontainer) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_applicationcontainer.Description {
			return fmt.Errorf("Bad cloud_applicationcontainer Description %s", cloud_applicationcontainer.Description)
		}

		if "demo_app" != cloud_applicationcontainer.Name {
			return fmt.Errorf("Bad cloud_applicationcontainer name %s", cloud_applicationcontainer.Name)
		}

		if "tag_app" != cloud_applicationcontainer.Annotation {
			return fmt.Errorf("Bad cloud_applicationcontainer annotation %s", cloud_applicationcontainer.Annotation)
		}

		if name_alias != cloud_applicationcontainer.NameAlias {
			return fmt.Errorf("Bad cloud_applicationcontainer name_alias %s", cloud_applicationcontainer.NameAlias)
		}

		return nil
	}
}
