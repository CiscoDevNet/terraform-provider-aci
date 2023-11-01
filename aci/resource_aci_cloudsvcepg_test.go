package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciCloudServiceEPg_Basic(t *testing.T) {
	var cloud_service_epg models.CloudServiceEPg
	fv_tenant_name := acctest.RandString(5)
	cloud_app_name := acctest.RandString(5)
	cloud_svc_epg_name := acctest.RandString(5)
	description := "cloud_service_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudServiceEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudServiceEPgConfig_basic(fv_tenant_name, cloud_app_name, cloud_svc_epg_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudServiceEPgExists("aci_cloud_service_epg.foo_cloud_service_epg", &cloud_service_epg),
					testAccCheckAciCloudServiceEPgAttributes(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, description, &cloud_service_epg),
				),
			},
		},
	})
}

func TestAccAciCloudServiceEPg_Update(t *testing.T) {
	var cloud_service_epg models.CloudServiceEPg
	fv_tenant_name := acctest.RandString(5)
	cloud_app_name := acctest.RandString(5)
	cloud_svc_epg_name := acctest.RandString(5)
	description := "cloud_service_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudServiceEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudServiceEPgConfig_basic(fv_tenant_name, cloud_app_name, cloud_svc_epg_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudServiceEPgExists("aci_cloud_service_epg.foo_cloud_service_epg", &cloud_service_epg),
					testAccCheckAciCloudServiceEPgAttributes(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, description, &cloud_service_epg),
				),
			},
			{
				Config: testAccCheckAciCloudServiceEPgConfig_basic(fv_tenant_name, cloud_app_name, cloud_svc_epg_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudServiceEPgExists("aci_cloud_service_epg.foo_cloud_service_epg", &cloud_service_epg),
					testAccCheckAciCloudServiceEPgAttributes(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, description, &cloud_service_epg),
				),
			},
		},
	})
}

func testAccCheckAciCloudServiceEPgConfig_basic(fv_tenant_name, cloud_app_name, cloud_svc_epg_name string) string {
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

	`, fv_tenant_name, cloud_app_name, cloud_svc_epg_name)
}

func testAccCheckAciCloudServiceEPgExists(name string, cloud_service_epg *models.CloudServiceEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Service EPg %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Service EPg dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_service_epgFound := models.CloudServiceEPgFromContainer(cont)
		if cloud_service_epgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Service EPg %s not found", rs.Primary.ID)
		}
		*cloud_service_epg = *cloud_service_epgFound
		return nil
	}
}

func testAccCheckAciCloudServiceEPgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_service_epg" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_service_epg := models.CloudServiceEPgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Service EPg %s Still exists", cloud_service_epg.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudServiceEPgAttributes(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, description string, cloud_service_epg *models.CloudServiceEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cloud_svc_epg_name != GetMOName(cloud_service_epg.DistinguishedName) {
			return fmt.Errorf("Bad cloudsvc_epg %s", GetMOName(cloud_service_epg.DistinguishedName))
		}

		if cloud_app_name != GetMOName(GetParentDn(cloud_service_epg.DistinguishedName, cloud_service_epg.Rn)) {
			return fmt.Errorf(" Bad cloudapp %s", GetMOName(GetParentDn(cloud_service_epg.DistinguishedName, cloud_service_epg.Rn)))
		}
		if description != cloud_service_epg.Description {
			return fmt.Errorf("Bad cloud_service_epg Description %s", cloud_service_epg.Description)
		}
		return nil
	}
}
