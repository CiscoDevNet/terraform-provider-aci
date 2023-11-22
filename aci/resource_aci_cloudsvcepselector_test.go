package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciCloudServiceEndpointSelector_Basic(t *testing.T) {
	var cloud_service_endpoint_selector models.CloudServiceEndpointSelector
	fv_tenant_name := acctest.RandString(5)
	cloud_app_name := acctest.RandString(5)
	cloud_svc_epg_name := acctest.RandString(5)
	cloud_svc_epselector_name := acctest.RandString(5)
	description := "cloud_service_endpoint_selector created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudServiceEndpointSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudServiceEndpointSelectorConfig_basic(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_svc_epselector_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudServiceEndpointSelectorExists("aci_cloud_service_endpoint_selector.foo_cloud_service_endpoint_selector", &cloud_service_endpoint_selector),
					testAccCheckAciCloudServiceEndpointSelectorAttributes(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_svc_epselector_name, description, &cloud_service_endpoint_selector),
				),
			},
		},
	})
}

func TestAccAciCloudServiceEndpointSelector_Update(t *testing.T) {
	var cloud_service_endpoint_selector models.CloudServiceEndpointSelector
	fv_tenant_name := acctest.RandString(5)
	cloud_app_name := acctest.RandString(5)
	cloud_svc_epg_name := acctest.RandString(5)
	cloud_svc_epselector_name := acctest.RandString(5)
	description := "cloud_service_endpoint_selector created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudServiceEndpointSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudServiceEndpointSelectorConfig_basic(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_svc_epselector_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudServiceEndpointSelectorExists("aci_cloud_service_endpoint_selector.foo_cloud_service_endpoint_selector", &cloud_service_endpoint_selector),
					testAccCheckAciCloudServiceEndpointSelectorAttributes(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_svc_epselector_name, description, &cloud_service_endpoint_selector),
				),
			},
			{
				Config: testAccCheckAciCloudServiceEndpointSelectorConfig_basic(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_svc_epselector_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudServiceEndpointSelectorExists("aci_cloud_service_endpoint_selector.foo_cloud_service_endpoint_selector", &cloud_service_endpoint_selector),
					testAccCheckAciCloudServiceEndpointSelectorAttributes(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_svc_epselector_name, description, &cloud_service_endpoint_selector),
				),
			},
		},
	})
}

func testAccCheckAciCloudServiceEndpointSelectorConfig_basic(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_svc_epselector_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_cloud_applicationcontainer" "foo_cloud_applicationcontainer" {
		name 		= "%s"
		description = "cloud_applicationcontainer created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	resource "aci_cloud_service_epg" "foo_cloud_service_epg" {
		name 		= "%s"
		description = "cloud_service_epg created while acceptance testing"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.foo_cloud_applicationcontainer.id
	}

	resource "aci_cloud_service_endpoint_selector" "foo_cloud_service_endpoint_selector" {
		name 		= "%s"
		description = "cloud_service_endpoint_selector created while acceptance testing"
		cloud_service_epg_dn = aci_cloud_service_epg.foo_cloud_service_epg.id
	}

	`, fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_svc_epselector_name)
}

func testAccCheckAciCloudServiceEndpointSelectorExists(name string, cloud_service_endpoint_selector *models.CloudServiceEndpointSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Service Endpoint Selector %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Service Endpoint Selector dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_service_endpoint_selectorFound := models.CloudServiceEndpointSelectorFromContainer(cont)
		if cloud_service_endpoint_selectorFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Service Endpoint Selector %s not found", rs.Primary.ID)
		}
		*cloud_service_endpoint_selector = *cloud_service_endpoint_selectorFound
		return nil
	}
}

func testAccCheckAciCloudServiceEndpointSelectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_service_endpoint_selector" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_service_endpoint_selector := models.CloudServiceEndpointSelectorFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Service Endpoint Selector %s Still exists", cloud_service_endpoint_selector.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudServiceEndpointSelectorAttributes(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_svc_epselector_name, description string, cloud_service_endpoint_selector *models.CloudServiceEndpointSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cloud_svc_epselector_name != GetMOName(cloud_service_endpoint_selector.DistinguishedName) {
			return fmt.Errorf("Bad cloudsvc_epselector %s", GetMOName(cloud_service_endpoint_selector.DistinguishedName))
		}

		if cloud_svc_epg_name != GetMOName(GetParentDn(cloud_service_endpoint_selector.DistinguishedName, cloud_service_endpoint_selector.Rn)) {
			return fmt.Errorf(" Bad cloudsvc_epg %s", GetMOName(GetParentDn(cloud_service_endpoint_selector.DistinguishedName, cloud_service_endpoint_selector.Rn)))
		}
		if description != cloud_service_endpoint_selector.Description {
			return fmt.Errorf("Bad cloud_service_endpoint_selector Description %s", cloud_service_endpoint_selector.Description)
		}
		return nil
	}
}
