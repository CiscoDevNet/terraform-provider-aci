package aci

import (
	"fmt"
	"log"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL1L2RedirectDestTraffic_Basic(t *testing.T) {
	var l1_l2_redirect_dest_traffic models.L1L2RedirectDestTraffic
	fv_tenant_name := acctest.RandString(5)
	vns_backup_pol_name := acctest.RandString(5)
	vns_l1_l2_redirect_dest_name := acctest.RandString(5)
	description := "Destination of L1/L2 redirected traffic created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL1L2RedirectDestTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL1L2RedirectDestTrafficConfig_basic(fv_tenant_name, vns_backup_pol_name, vns_l1_l2_redirect_dest_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL1L2RedirectDestTrafficExists("aci_pbr_l1_l2_destination.foo_l1_l2_redirect_dest_traffic", &l1_l2_redirect_dest_traffic),
					testAccCheckAciL1L2RedirectDestTrafficAttributes(fv_tenant_name, vns_backup_pol_name, vns_l1_l2_redirect_dest_name, description, &l1_l2_redirect_dest_traffic),
				),
			},
		},
	})
}

func TestAccAciL1L2RedirectDestTraffic_Update(t *testing.T) {
	var l1_l2_redirect_dest_traffic models.L1L2RedirectDestTraffic
	fv_tenant_name := acctest.RandString(5)
	vns_backup_pol_name := acctest.RandString(5)
	vns_l1_l2_redirect_dest_name := acctest.RandString(5)
	description := "Destination of L1/L2 redirected traffic created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL1L2RedirectDestTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL1L2RedirectDestTrafficConfig_basic(fv_tenant_name, vns_backup_pol_name, vns_l1_l2_redirect_dest_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL1L2RedirectDestTrafficExists("aci_pbr_l1_l2_destination.foo_l1_l2_redirect_dest_traffic", &l1_l2_redirect_dest_traffic),
					testAccCheckAciL1L2RedirectDestTrafficAttributes(fv_tenant_name, vns_backup_pol_name, vns_l1_l2_redirect_dest_name, description, &l1_l2_redirect_dest_traffic),
				),
			},
			{
				Config: testAccCheckAciL1L2RedirectDestTrafficConfig_basic(fv_tenant_name, vns_backup_pol_name, vns_l1_l2_redirect_dest_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL1L2RedirectDestTrafficExists("aci_pbr_l1_l2_destination.foo_l1_l2_redirect_dest_traffic", &l1_l2_redirect_dest_traffic),
					testAccCheckAciL1L2RedirectDestTrafficAttributes(fv_tenant_name, vns_backup_pol_name, vns_l1_l2_redirect_dest_name, description, &l1_l2_redirect_dest_traffic),
				),
			},
		},
	})
}

func testAccCheckAciL1L2RedirectDestTrafficConfig_basic(fv_tenant_name, vns_backup_pol_name, vns_l1_l2_redirect_dest_name string) string {
	return fmt.Sprintf(`
		resource "aci_tenant" "footenant" {
			name        = "%s"
			description = "tenant created while acceptance testing"
		}

		# Concrete Interface Setup
		resource "aci_l4_l7_device" "l4_l7_device" {
			tenant_dn   = aci_tenant.footenant.id
			name        = "tf_l4_l7_device"
			device_type = "CLOUD"
		}

		resource "aci_concrete_device" "concrete_device" {
			l4_l7_device_dn = aci_l4_l7_device.l4_l7_device.id
			name            = "tf_concrete_device"
		}

		resource "aci_concrete_interface" "concrete_interface" {
			concrete_device_dn = aci_concrete_device.concrete_device.id
			name               = "tf_concrete_interface"
		}

		resource "aci_service_redirect_backup_policy" "foop_br_backup_policy" {
			name        = "%s"
			description = "pbr_backup_policy created while acceptance testing"
			tenant_dn   = aci_tenant.footenant.id
		}

		resource "aci_pbr_l1_l2_destination" "foo_l1_l2_redirect_dest_traffic" {
			dest_name                = "%s"
			description              = "Destination of L1/L2 redirected traffic created while acceptance testing"
			policy_based_redirect_dn = aci_service_redirect_backup_policy.foop_br_backup_policy.id
			relation_vns_rs_to_c_if  = aci_concrete_interface.concrete_interface.id
		}
	`, fv_tenant_name, vns_backup_pol_name, vns_l1_l2_redirect_dest_name)
}

func testAccCheckAciL1L2RedirectDestTrafficExists(name string, l1_l2_redirect_dest_traffic *models.L1L2RedirectDestTraffic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L1/L2 Redirect Destination Traffic %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L1/L2 Redirect Destination Traffic dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l1_l2_redirect_dest_trafficFound := models.L1L2RedirectDestTrafficFromContainer(cont)
		if l1_l2_redirect_dest_trafficFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L1/L2 Redirect Destination Traffic %s not found", rs.Primary.ID)
		}
		*l1_l2_redirect_dest_traffic = *l1_l2_redirect_dest_trafficFound
		return nil
	}
}

func testAccCheckAciL1L2RedirectDestTrafficDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_pbr_l1_l2_destination" {
			cont, err := client.Get(rs.Primary.ID)
			l1_l2_redirect_dest_traffic := models.L1L2RedirectDestTrafficFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L1/L2 Redirect Destination Traffic %s Still exists", l1_l2_redirect_dest_traffic.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL1L2RedirectDestTrafficAttributes(fv_tenant_name, vns_backup_pol_name, vns_l1_l2_redirect_dest_name, description string, l1_l2_redirect_dest_traffic *models.L1L2RedirectDestTraffic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vns_l1_l2_redirect_dest_name != GetMOName(l1_l2_redirect_dest_traffic.DistinguishedName) {
			return fmt.Errorf("Bad L1/L2 Redirect Destination Traffic %s", GetMOName(l1_l2_redirect_dest_traffic.DistinguishedName))
		}
		if vns_backup_pol_name != GetMOName(GetParentDn(l1_l2_redirect_dest_traffic.DistinguishedName, fmt.Sprintf("/%s", fmt.Sprintf(models.RnvnsL1L2RedirectDest, vns_l1_l2_redirect_dest_name)))) {
			log.Printf("[DEBUG] vns_backup_pol_name: %s, Parent DN: %s", vns_backup_pol_name, GetParentDn(l1_l2_redirect_dest_traffic.DistinguishedName, fmt.Sprintf(models.RnvnsL1L2RedirectDest, vns_l1_l2_redirect_dest_name)))
			return fmt.Errorf("Bad Policy-Based Redirect Name %s", GetMOName(GetParentDn(l1_l2_redirect_dest_traffic.DistinguishedName, fmt.Sprintf("/%s", fmt.Sprintf(models.RnvnsL1L2RedirectDest, vns_l1_l2_redirect_dest_name)))))
		}

		if vns_l1_l2_redirect_dest_name != GetMOName(l1_l2_redirect_dest_traffic.DistinguishedName) {
			return fmt.Errorf("Bad L1/L2 Redirect Destination Traffic destName %s", GetMOName(l1_l2_redirect_dest_traffic.DistinguishedName))
		}

		if description != l1_l2_redirect_dest_traffic.Description {
			return fmt.Errorf("Bad L1/L2 Redirect Destination Traffic Description %s", l1_l2_redirect_dest_traffic.Description)
		}
		return nil
	}
}
