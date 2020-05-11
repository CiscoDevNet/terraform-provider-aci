package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciCloudEndpointSelector_Basic(t *testing.T) {
	var cloud_endpoint_selector models.CloudEndpointSelector
	description := "cloud_endpoint_selector created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudEndpointSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudEndpointSelectorConfig_basic(description, "custom:Name=='admin-ep2'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorExists("aci_cloud_endpoint_selector.foocloud_endpoint_selector", &cloud_endpoint_selector),
					testAccCheckAciCloudEndpointSelectorAttributes(description, "custom:Name=='admin-ep2'", &cloud_endpoint_selector),
				),
			},
			{
				ResourceName:      "aci_cloud_endpoint_selector",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciCloudEndpointSelector_update(t *testing.T) {
	var cloud_endpoint_selector models.CloudEndpointSelector
	description := "cloud_endpoint_selector created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudEndpointSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudEndpointSelectorConfig_basic(description, "custom:Name=='admin-ep2'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorExists("aci_cloud_endpoint_selector.foocloud_endpoint_selector", &cloud_endpoint_selector),
					testAccCheckAciCloudEndpointSelectorAttributes(description, "custom:Name=='admin-ep2'", &cloud_endpoint_selector),
				),
			},
			{
				Config: testAccCheckAciCloudEndpointSelectorConfig_basic(description, "custom:Name=='admin-ep1'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorExists("aci_cloud_endpoint_selector.foocloud_endpoint_selector", &cloud_endpoint_selector),
					testAccCheckAciCloudEndpointSelectorAttributes(description, "custom:Name=='admin-ep1'", &cloud_endpoint_selector),
				),
			},
		},
	})
}

func testAccCheckAciCloudEndpointSelectorConfig_basic(description, match_expression string) string {
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

	resource "aci_cloud_epg" "foocloud_e_pg" {
		cloud_applicationcontainer_dn = "${aci_cloud_applicationcontainer.foocloud_applicationcontainer.id}"
		name                          = "cloud_epg"
	}

	resource "aci_cloud_endpoint_selector" "foocloud_endpoint_selector" {
		cloud_e_pg_dn    = "${aci_cloud_epg.foocloud_e_pg.id}"
		description      = "%s"
		name             = "ep_select"
		annotation       = "tag_ep"
		match_expression = "%s"
		name_alias       = "alias_ep"
	}
	  
	`, description, match_expression)
}

func testAccCheckAciCloudEndpointSelectorExists(name string, cloud_endpoint_selector *models.CloudEndpointSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Endpoint Selector %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Endpoint Selector dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_endpoint_selectorFound := models.CloudEndpointSelectorFromContainer(cont)
		if cloud_endpoint_selectorFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Endpoint Selector %s not found", rs.Primary.ID)
		}
		*cloud_endpoint_selector = *cloud_endpoint_selectorFound
		return nil
	}
}

func testAccCheckAciCloudEndpointSelectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_endpoint_selector" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_endpoint_selector := models.CloudEndpointSelectorFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Endpoint Selector %s Still exists", cloud_endpoint_selector.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudEndpointSelectorAttributes(description, match_expression string, cloud_endpoint_selector *models.CloudEndpointSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_endpoint_selector.Description {
			return fmt.Errorf("Bad cloud_endpoint_selector Description %s", cloud_endpoint_selector.Description)
		}

		if "ep_select" != cloud_endpoint_selector.Name {
			return fmt.Errorf("Bad cloud_endpoint_selector name %s", cloud_endpoint_selector.Name)
		}

		if "tag_ep" != cloud_endpoint_selector.Annotation {
			return fmt.Errorf("Bad cloud_endpoint_selector annotation %s", cloud_endpoint_selector.Annotation)
		}

		if match_expression != cloud_endpoint_selector.MatchExpression {
			return fmt.Errorf("Bad cloud_endpoint_selector match_expression %s", cloud_endpoint_selector.MatchExpression)
		}

		if "alias_ep" != cloud_endpoint_selector.NameAlias {
			return fmt.Errorf("Bad cloud_endpoint_selector name_alias %s", cloud_endpoint_selector.NameAlias)
		}

		return nil
	}
}
