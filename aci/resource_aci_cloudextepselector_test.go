package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciCloudEndpointSelectorforExternalEPgs_Basic(t *testing.T) {
	var cloud_endpoint_selectorfor_external_epgs models.CloudEndpointSelectorforExternalEPgs
	description := "cloud_endpoint_selectorfor_external_epgs created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudEndpointSelectorforExternalEPgsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudEndpointSelectorforExternalEPgsConfig_basic(description, "0.0.0.0/0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorforExternalEPgsExists("aci_cloud_endpoint_selectorfor_external_epgs.foocloud_endpoint_selectorfor_external_epgs", &cloud_endpoint_selectorfor_external_epgs),
					testAccCheckAciCloudEndpointSelectorforExternalEPgsAttributes(description, "0.0.0.0/0", &cloud_endpoint_selectorfor_external_epgs),
				),
			},
		},
	})
}

func TestAccAciCloudEndpointSelectorforExternalEPgs_update(t *testing.T) {
	var cloud_endpoint_selectorfor_external_epgs models.CloudEndpointSelectorforExternalEPgs
	description := "cloud_endpoint_selectorfor_external_epgs created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudEndpointSelectorforExternalEPgsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudEndpointSelectorforExternalEPgsConfig_basic(description, "0.0.0.0/0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorforExternalEPgsExists("aci_cloud_endpoint_selectorfor_external_epgs.foocloud_endpoint_selectorfor_external_epgs", &cloud_endpoint_selectorfor_external_epgs),
					testAccCheckAciCloudEndpointSelectorforExternalEPgsAttributes(description, "0.0.0.0/0", &cloud_endpoint_selectorfor_external_epgs),
				),
			},
			{
				Config: testAccCheckAciCloudEndpointSelectorforExternalEPgsConfig_basic(description, "10.0.0.0/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorforExternalEPgsExists("aci_cloud_endpoint_selectorfor_external_epgs.foocloud_endpoint_selectorfor_external_epgs", &cloud_endpoint_selectorfor_external_epgs),
					testAccCheckAciCloudEndpointSelectorforExternalEPgsAttributes(description, "10.0.0.0/24", &cloud_endpoint_selectorfor_external_epgs),
				),
			},
		},
	})
}

func testAccCheckAciCloudEndpointSelectorforExternalEPgsConfig_basic(description, subnet string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		description = "Tenant created while acceptance testing"
		name        = "demo_tenant"
	}

	resource "aci_cloud_applicationcontainer" "foocloud_applicationcontainer" {
		tenant_dn   = aci_tenant.footenant.id
		name        = "demo_app"
		annotation  = "tag_app"
	}

	resource "aci_cloud_external_epg" "foocloud_external_epg" {
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.foocloud_applicationcontainer.id
		name                          = "cloud_ext_epg"
	}

	resource "aci_cloud_endpoint_selectorfor_external_epgs" "foocloud_endpoint_selectorfor_external_epgs" {
		cloud_external_epg_dn  = aci_cloud_external_epg.foocloud_external_epg.id
		description            = "%s"
		name                   = "ext_ep_selector"
		annotation             = "tag_ext_selector"
		is_shared              = "yes"
		name_alias             = "alias_select"
		subnet                 = "%s"
	}
	  
	`, description, subnet)
}

func testAccCheckAciCloudEndpointSelectorforExternalEPgsExists(name string, cloud_endpoint_selectorfor_external_epgs *models.CloudEndpointSelectorforExternalEPgs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Endpoint Selector for External EPgs %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Endpoint Selector for External EPgs dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_endpoint_selectorfor_external_epgsFound := models.CloudEndpointSelectorforExternalEPgsFromContainer(cont)
		if cloud_endpoint_selectorfor_external_epgsFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Endpoint Selector for External EPgs %s not found", rs.Primary.ID)
		}
		*cloud_endpoint_selectorfor_external_epgs = *cloud_endpoint_selectorfor_external_epgsFound
		return nil
	}
}

func testAccCheckAciCloudEndpointSelectorforExternalEPgsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_endpoint_selectorfor_external_epgs" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_endpoint_selectorfor_external_epgs := models.CloudEndpointSelectorforExternalEPgsFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Endpoint Selector for External EPgs %s Still exists", cloud_endpoint_selectorfor_external_epgs.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudEndpointSelectorforExternalEPgsAttributes(description, subnet string, cloud_endpoint_selectorfor_external_epgs *models.CloudEndpointSelectorforExternalEPgs) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_endpoint_selectorfor_external_epgs.Description {
			return fmt.Errorf("Bad cloud_endpoint_selectorfor_external_epgs Description %s", cloud_endpoint_selectorfor_external_epgs.Description)
		}

		if "ext_ep_selector" != cloud_endpoint_selectorfor_external_epgs.Name {
			return fmt.Errorf("Bad cloud_endpoint_selectorfor_external_epgs name %s", cloud_endpoint_selectorfor_external_epgs.Name)
		}

		if "tag_ext_selector" != cloud_endpoint_selectorfor_external_epgs.Annotation {
			return fmt.Errorf("Bad cloud_endpoint_selectorfor_external_epgs annotation %s", cloud_endpoint_selectorfor_external_epgs.Annotation)
		}

		if "yes" != cloud_endpoint_selectorfor_external_epgs.IsShared {
			return fmt.Errorf("Bad cloud_endpoint_selectorfor_external_epgs is_shared %s", cloud_endpoint_selectorfor_external_epgs.IsShared)
		}

		if "alias_select" != cloud_endpoint_selectorfor_external_epgs.NameAlias {
			return fmt.Errorf("Bad cloud_endpoint_selectorfor_external_epgs name_alias %s", cloud_endpoint_selectorfor_external_epgs.NameAlias)
		}

		if subnet != cloud_endpoint_selectorfor_external_epgs.Subnet {
			return fmt.Errorf("Bad cloud_endpoint_selectorfor_external_epgs subnet %s", cloud_endpoint_selectorfor_external_epgs.Subnet)
		}

		return nil
	}
}
