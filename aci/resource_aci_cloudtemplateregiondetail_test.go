package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciAdditionalconfigforregion_Basic(t *testing.T) {
	var additionalconfigforregion models.Additionalconfigforregion
	fv_tenant_name := acctest.RandString(5)
	cloudtemplate_infra_network_name := acctest.RandString(5)
	cloudtemplate_stats_name := acctest.RandString(5)
	cloud_region_name_name := acctest.RandString(5)
	cloudtemplate_region_detail_name := acctest.RandString(5)
	description := "additionalconfigforregion created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAdditionalconfigforregionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAdditionalconfigforregionConfig_basic(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_stats_name, cloud_region_name_name, cloudtemplate_region_detail_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAdditionalconfigforregionExists("aci_additionalconfigforregion.foo_additionalconfigforregion", &additionalconfigforregion),
					testAccCheckAciAdditionalconfigforregionAttributes(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_stats_name, cloud_region_name_name, cloudtemplate_region_detail_name, description, &additionalconfigforregion),
				),
			},
		},
	})
}

func TestAccAciAdditionalconfigforregion_Update(t *testing.T) {
	var additionalconfigforregion models.Additionalconfigforregion
	fv_tenant_name := acctest.RandString(5)
	cloudtemplate_infra_network_name := acctest.RandString(5)
	cloudtemplate_stats_name := acctest.RandString(5)
	cloud_region_name_name := acctest.RandString(5)
	cloudtemplate_region_detail_name := acctest.RandString(5)
	description := "additionalconfigforregion created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAdditionalconfigforregionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAdditionalconfigforregionConfig_basic(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_stats_name, cloud_region_name_name, cloudtemplate_region_detail_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAdditionalconfigforregionExists("aci_additionalconfigforregion.foo_additionalconfigforregion", &additionalconfigforregion),
					testAccCheckAciAdditionalconfigforregionAttributes(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_stats_name, cloud_region_name_name, cloudtemplate_region_detail_name, description, &additionalconfigforregion),
				),
			},
			{
				Config: testAccCheckAciAdditionalconfigforregionConfig_basic(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_stats_name, cloud_region_name_name, cloudtemplate_region_detail_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAdditionalconfigforregionExists("aci_additionalconfigforregion.foo_additionalconfigforregion", &additionalconfigforregion),
					testAccCheckAciAdditionalconfigforregionAttributes(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_stats_name, cloud_region_name_name, cloudtemplate_region_detail_name, description, &additionalconfigforregion),
				),
			},
		},
	})
}

func testAccCheckAciAdditionalconfigforregionConfig_basic(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_stats_name, cloud_region_name_name, cloudtemplate_region_detail_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_infra_network_template" "foo_infra_network_template" {
		name 		= "%s"
		description = "infra_network_template created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	resource "aci_cloud_statistics" "foo_cloud_statistics" {
		name 		= "%s"
		description = "cloud_statistics created while acceptance testing"
		infra_network_template_dn = aci_infra_network_template.foo_infra_network_template.id
	}

	resource "aci_cloud_providerand_region_names" "foo_cloud_providerand_region_names" {
		name 		= "%s"
		description = "cloud_providerand_region_names created while acceptance testing"
		cloud_statistics_dn = aci_cloud_statistics.foo_cloud_statistics.id
	}

	resource "aci_additionalconfigforregion" "foo_additionalconfigforregion" {
		name 		= "%s"
		description = "additionalconfigforregion created while acceptance testing"
		cloud_providerand_region_names_dn = aci_cloud_providerand_region_names.foo_cloud_providerand_region_names.id
	}

	`, fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_stats_name, cloud_region_name_name, cloudtemplate_region_detail_name)
}

func testAccCheckAciAdditionalconfigforregionExists(name string, additionalconfigforregion *models.Additionalconfigforregion) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Additional config for region %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Additional config for region dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		additionalconfigforregionFound := models.AdditionalconfigforregionFromContainer(cont)
		if additionalconfigforregionFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Additional config for region %s not found", rs.Primary.ID)
		}
		*additionalconfigforregion = *additionalconfigforregionFound
		return nil
	}
}

func testAccCheckAciAdditionalconfigforregionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_additionalconfigforregion" {
			cont, err := client.Get(rs.Primary.ID)
			additionalconfigforregion := models.AdditionalconfigforregionFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Additional config for region %s Still exists", additionalconfigforregion.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciAdditionalconfigforregionAttributes(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_stats_name, cloud_region_name_name, cloudtemplate_region_detail_name, description string, additionalconfigforregion *models.Additionalconfigforregion) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cloudtemplate_region_detail_name != GetMOName(additionalconfigforregion.DistinguishedName) {
			return fmt.Errorf("Bad cloudtemplateregion_detail %s", GetMOName(additionalconfigforregion.DistinguishedName))
		}

		if cloud_region_name_name != GetMOName(GetParentDn(additionalconfigforregion.DistinguishedName, additionalconfigforregion.Rn)) {
			return fmt.Errorf(" Bad cloudregion_name %s", GetMOName(GetParentDn(additionalconfigforregion.DistinguishedName, additionalconfigforregion.Rn)))
		}
		if description != additionalconfigforregion.Description {
			return fmt.Errorf("Bad additionalconfigforregion Description %s", additionalconfigforregion.Description)
		}
		return nil
	}
}
