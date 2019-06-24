package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAciCloudEndpointSelectorforExternalEPgs_Basic(t *testing.T) {
	var cloud_endpoint_selectorfor_external_e_pgs models.CloudEndpointSelectorforExternalEPgs
	description := "cloud_endpoint_selectorfor_external_e_pgs created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudEndpointSelectorforExternalEPgsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudEndpointSelectorforExternalEPgsConfig_basic(description, "0.0.0.0/0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorforExternalEPgsExists("aci_cloud_endpoint_selectorfor_external_e_pgs.foocloud_endpoint_selectorfor_external_e_pgs", &cloud_endpoint_selectorfor_external_e_pgs),
					testAccCheckAciCloudEndpointSelectorforExternalEPgsAttributes(description, "0.0.0.0/0", &cloud_endpoint_selectorfor_external_e_pgs),
				),
			},
			{
				ResourceName:      "aci_cloud_endpoint_selectorfor_external_e_pgs",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciCloudEndpointSelectorforExternalEPgs_update(t *testing.T) {
	var cloud_endpoint_selectorfor_external_e_pgs models.CloudEndpointSelectorforExternalEPgs
	description := "cloud_endpoint_selectorfor_external_e_pgs created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudEndpointSelectorforExternalEPgsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudEndpointSelectorforExternalEPgsConfig_basic(description, "0.0.0.0/0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorforExternalEPgsExists("aci_cloud_endpoint_selectorfor_external_e_pgs.foocloud_endpoint_selectorfor_external_e_pgs", &cloud_endpoint_selectorfor_external_e_pgs),
					testAccCheckAciCloudEndpointSelectorforExternalEPgsAttributes(description, "0.0.0.0/0", &cloud_endpoint_selectorfor_external_e_pgs),
				),
			},
			{
				Config: testAccCheckAciCloudEndpointSelectorforExternalEPgsConfig_basic(description, "10.0.0.0/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorforExternalEPgsExists("aci_cloud_endpoint_selectorfor_external_e_pgs.foocloud_endpoint_selectorfor_external_e_pgs", &cloud_endpoint_selectorfor_external_e_pgs),
					testAccCheckAciCloudEndpointSelectorforExternalEPgsAttributes(description, "10.0.0.0/24", &cloud_endpoint_selectorfor_external_e_pgs),
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
		tenant_dn   = "${aci_tenant.footenant.id}"
		name        = "demo_app"
		annotation  = "tag_app"
	}

	resource "aci_cloud_external_e_pg" "foocloud_external_e_pg" {
		cloud_applicationcontainer_dn = "${aci_cloud_applicationcontainer.foocloud_applicationcontainer.id}"
		name                          = "cloud_ext_epg"
	}

	resource "aci_cloud_endpoint_selectorfor_external_e_pgs" "foocloud_endpoint_selectorfor_external_e_pgs" {
		cloud_external_e_pg_dn = "${aci_cloud_external_e_pg.foocloud_external_e_pg.id}"
		description            = "%s"
		name                   = "ext_ep_selector"
		annotation             = "tag_ext_selector"
		is_shared              = "yes"
		name_alias             = "alias_select"
		subnet                 = "%s"
	}
	  
	`, description, subnet)
}

func testAccCheckAciCloudEndpointSelectorforExternalEPgsExists(name string, cloud_endpoint_selectorfor_external_e_pgs *models.CloudEndpointSelectorforExternalEPgs) resource.TestCheckFunc {
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

		cloud_endpoint_selectorfor_external_e_pgsFound := models.CloudEndpointSelectorforExternalEPgsFromContainer(cont)
		if cloud_endpoint_selectorfor_external_e_pgsFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Endpoint Selector for External EPgs %s not found", rs.Primary.ID)
		}
		*cloud_endpoint_selectorfor_external_e_pgs = *cloud_endpoint_selectorfor_external_e_pgsFound
		return nil
	}
}

func testAccCheckAciCloudEndpointSelectorforExternalEPgsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_endpoint_selectorfor_external_e_pgs" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_endpoint_selectorfor_external_e_pgs := models.CloudEndpointSelectorforExternalEPgsFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Endpoint Selector for External EPgs %s Still exists", cloud_endpoint_selectorfor_external_e_pgs.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudEndpointSelectorforExternalEPgsAttributes(description, subnet string, cloud_endpoint_selectorfor_external_e_pgs *models.CloudEndpointSelectorforExternalEPgs) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_endpoint_selectorfor_external_e_pgs.Description {
			return fmt.Errorf("Bad cloud_endpoint_selectorfor_external_e_pgs Description %s", cloud_endpoint_selectorfor_external_e_pgs.Description)
		}

		if "ext_ep_selector" != cloud_endpoint_selectorfor_external_e_pgs.Name {
			return fmt.Errorf("Bad cloud_endpoint_selectorfor_external_e_pgs name %s", cloud_endpoint_selectorfor_external_e_pgs.Name)
		}

		if "tag_ext_selector" != cloud_endpoint_selectorfor_external_e_pgs.Annotation {
			return fmt.Errorf("Bad cloud_endpoint_selectorfor_external_e_pgs annotation %s", cloud_endpoint_selectorfor_external_e_pgs.Annotation)
		}

		if "yes" != cloud_endpoint_selectorfor_external_e_pgs.IsShared {
			return fmt.Errorf("Bad cloud_endpoint_selectorfor_external_e_pgs is_shared %s", cloud_endpoint_selectorfor_external_e_pgs.IsShared)
		}

		if "alias_select" != cloud_endpoint_selectorfor_external_e_pgs.NameAlias {
			return fmt.Errorf("Bad cloud_endpoint_selectorfor_external_e_pgs name_alias %s", cloud_endpoint_selectorfor_external_e_pgs.NameAlias)
		}

		if subnet != cloud_endpoint_selectorfor_external_e_pgs.Subnet {
			return fmt.Errorf("Bad cloud_endpoint_selectorfor_external_e_pgs subnet %s", cloud_endpoint_selectorfor_external_e_pgs.Subnet)
		}

		return nil
	}
}
