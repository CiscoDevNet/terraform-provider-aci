package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciPrivateLinkLabelfortheserviceEPg_Basic(t *testing.T) {
	var private_link_labelfortheservice_epg models.PrivateLinkLabelfortheserviceEPg
	fv_tenant_name := acctest.RandString(5)
	cloud_app_name := acctest.RandString(5)
	cloud_svc_epg_name := acctest.RandString(5)
	cloud_private_link_label_name := acctest.RandString(5)
	description := "private_link_labelfortheservice_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPrivateLinkLabelfortheserviceEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPrivateLinkLabelfortheserviceEPgConfig_basic(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_private_link_label_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPrivateLinkLabelfortheserviceEPgExists("aci_private_link_labelfortheservice_epg.foo_private_link_labelfortheservice_epg", &private_link_labelfortheservice_epg),
					testAccCheckAciPrivateLinkLabelfortheserviceEPgAttributes(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_private_link_label_name, description, &private_link_labelfortheservice_epg),
				),
			},
		},
	})
}

func TestAccAciPrivateLinkLabelfortheserviceEPg_Update(t *testing.T) {
	var private_link_labelfortheservice_epg models.PrivateLinkLabelfortheserviceEPg
	fv_tenant_name := acctest.RandString(5)
	cloud_app_name := acctest.RandString(5)
	cloud_svc_epg_name := acctest.RandString(5)
	cloud_private_link_label_name := acctest.RandString(5)
	description := "private_link_labelfortheservice_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPrivateLinkLabelfortheserviceEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPrivateLinkLabelfortheserviceEPgConfig_basic(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_private_link_label_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPrivateLinkLabelfortheserviceEPgExists("aci_private_link_labelfortheservice_epg.foo_private_link_labelfortheservice_epg", &private_link_labelfortheservice_epg),
					testAccCheckAciPrivateLinkLabelfortheserviceEPgAttributes(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_private_link_label_name, description, &private_link_labelfortheservice_epg),
				),
			},
			{
				Config: testAccCheckAciPrivateLinkLabelfortheserviceEPgConfig_basic(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_private_link_label_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPrivateLinkLabelfortheserviceEPgExists("aci_private_link_labelfortheservice_epg.foo_private_link_labelfortheservice_epg", &private_link_labelfortheservice_epg),
					testAccCheckAciPrivateLinkLabelfortheserviceEPgAttributes(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_private_link_label_name, description, &private_link_labelfortheservice_epg),
				),
			},
		},
	})
}

func testAccCheckAciPrivateLinkLabelfortheserviceEPgConfig_basic(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_private_link_label_name string) string {
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

	resource "aci_private_link_labelfortheservice_epg" "foo_private_link_labelfortheservice_epg" {
		name 		= "%s"
		description = "private_link_labelfortheservice_epg created while acceptance testing"
		cloud_service_epg_dn = aci_cloud_service_epg.foo_cloud_service_epg.id
	}

	`, fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_private_link_label_name)
}

func testAccCheckAciPrivateLinkLabelfortheserviceEPgExists(name string, private_link_labelfortheservice_epg *models.PrivateLinkLabelfortheserviceEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Private Link Label for the service EPg %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Private Link Label for the service EPg dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		private_link_labelfortheservice_epgFound := models.PrivateLinkLabelfortheserviceEPgFromContainer(cont)
		if private_link_labelfortheservice_epgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Private Link Label for the service EPg %s not found", rs.Primary.ID)
		}
		*private_link_labelfortheservice_epg = *private_link_labelfortheservice_epgFound
		return nil
	}
}

func testAccCheckAciPrivateLinkLabelfortheserviceEPgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_private_link_labelfortheservice_epg" {
			cont, err := client.Get(rs.Primary.ID)
			private_link_labelfortheservice_epg := models.PrivateLinkLabelfortheserviceEPgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Private Link Label for the service EPg %s Still exists", private_link_labelfortheservice_epg.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciPrivateLinkLabelfortheserviceEPgAttributes(fv_tenant_name, cloud_app_name, cloud_svc_epg_name, cloud_private_link_label_name, description string, private_link_labelfortheservice_epg *models.PrivateLinkLabelfortheserviceEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cloud_private_link_label_name != GetMOName(private_link_labelfortheservice_epg.DistinguishedName) {
			return fmt.Errorf("Bad cloudprivate_link_label %s", GetMOName(private_link_labelfortheservice_epg.DistinguishedName))
		}

		if cloud_svc_epg_name != GetMOName(GetParentDn(private_link_labelfortheservice_epg.DistinguishedName, private_link_labelfortheservice_epg.Rn)) {
			return fmt.Errorf(" Bad cloudsvc_epg %s", GetMOName(GetParentDn(private_link_labelfortheservice_epg.DistinguishedName, private_link_labelfortheservice_epg.Rn)))
		}
		if description != private_link_labelfortheservice_epg.Description {
			return fmt.Errorf("Bad private_link_labelfortheservice_epg Description %s", private_link_labelfortheservice_epg.Description)
		}
		return nil
	}
}
