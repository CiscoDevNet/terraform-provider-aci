package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSubnet_Basic(t *testing.T) {
	var subnet models.Subnet
	description := "subnet created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubnetConfig_basic(description, "unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists("aci_subnet.foosubnet", &subnet),
					testAccCheckAciSubnetAttributes(description, "unspecified", &subnet),
				),
			},
		},
	})
}

func TestAccAciSubnet_update(t *testing.T) {
	var subnet models.Subnet
	description := "subnet created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubnetConfig_basic(description, "unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists("aci_subnet.foosubnet", &subnet),
					testAccCheckAciSubnetAttributes(description, "unspecified", &subnet),
				),
			},
			{
				Config: testAccCheckAciSubnetConfig_basic(description, "nd"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists("aci_subnet.foosubnet", &subnet),
					testAccCheckAciSubnetAttributes(description, "nd", &subnet),
				),
			},
		},
	})
}

func testAccCheckAciSubnetConfig_basic(description, Ctrl string) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "tenant_for_subnet" {
		name        = "tenant_for_subnet"
		description = "This tenant is created by terraform ACI provider"
	}
	resource "aci_bridge_domain" "bd_for_subnet" {
		tenant_dn = aci_tenant.tenant_for_subnet.id
		name      = "bd_for_subnet"
	}

	resource "aci_subnet" "foosubnet" {
		parent_dn   = aci_bridge_domain.bd_for_subnet.id
		description = "%s"
		ip          = "10.0.3.28/27"
		annotation  = "tag_subnet"
		ctrl        = ["nd"]
		name_alias  = "alias_subnet"
		preferred   = "no"
		scope       = ["private"]
		virtual     = "yes"
	}
	`, description)
}

func testAccCheckAciSubnetExists(name string, subnet *models.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Subnet %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subnet dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		subnetFound := models.SubnetFromContainer(cont)
		if subnetFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Subnet %s not found", rs.Primary.ID)
		}
		*subnet = *subnetFound
		return nil
	}
}

func testAccCheckAciSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_subnet" {
			cont, err := client.Get(rs.Primary.ID)
			subnet := models.SubnetFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Subnet %s Still exists", subnet.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSubnetAttributes(description, Ctrl string, subnet *models.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != subnet.Description {
			return fmt.Errorf("Bad subnet Description %s", subnet.Description)
		}

		if "10.0.3.28/27" != subnet.Ip {
			return fmt.Errorf("Bad subnet ip %s", subnet.Ip)
		}

		if "tag_subnet" != subnet.Annotation {
			return fmt.Errorf("Bad subnet annotation %s", subnet.Annotation)
		}

		if "alias_subnet" != subnet.NameAlias {
			return fmt.Errorf("Bad subnet name_alias %s", subnet.NameAlias)
		}

		if "no" != subnet.Preferred {
			return fmt.Errorf("Bad subnet preferred %s", subnet.Preferred)
		}

		if "yes" != subnet.Virtual {
			return fmt.Errorf("Bad subnet virtual %s", subnet.Virtual)
		}

		return nil
	}
}

// MS NLB Check only for IGMP
func TestAccAciSubnetEPGsMsNlbIgmp_Basic(t *testing.T) {
	var subnet models.Subnet
	var ms_nlb models.NlbEndpoint

	ms_nlb_group_value := "224.0.0.1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetEPGsMsNlbIgmpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubnetEPGsMsNlbIgmpConfig_basic(
					ms_nlb_group_value,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetEPGsMsNlbIgmpExists(
						"aci_subnet.foo_epg_subnet_msnlb_mcast_igmp",
						&subnet,
						&ms_nlb,
					),
					testAccCheckAciSubnetEPGsMsNlbIgmpAttributes(
						ms_nlb_group_value,
						&ms_nlb,
					),
				),
			},
		},
	})
}

func TestAccAciSubnetEPGsMsNlbIgmp_update(t *testing.T) {
	var subnet models.Subnet
	var ms_nlb models.NlbEndpoint

	ms_nlb_group_value := "224.0.0.2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetEPGsMsNlbIgmpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubnetEPGsMsNlbIgmpConfig_basic(
					ms_nlb_group_value,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetEPGsMsNlbIgmpExists(
						"aci_subnet.foo_epg_subnet_msnlb_mcast_igmp",
						&subnet,
						&ms_nlb,
					),
					testAccCheckAciSubnetEPGsMsNlbIgmpAttributes(
						ms_nlb_group_value,
						&ms_nlb,
					),
				),
			},
			{
				Config: testAccCheckAciSubnetEPGsMsNlbIgmpConfig_basic(
					ms_nlb_group_value,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetEPGsMsNlbIgmpExists(
						"aci_subnet.foo_epg_subnet_msnlb_mcast_igmp",
						&subnet,
						&ms_nlb,
					),
					testAccCheckAciSubnetEPGsMsNlbIgmpAttributes(
						ms_nlb_group_value,
						&ms_nlb,
					),
				),
			},
		},
	})
}

func testAccCheckAciSubnetEPGsMsNlbIgmpConfig_basic(
	ms_nlb_group_value string,
) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "tenant_for_subnet" {
		name        = "tenant_for_subnet"
		description = "This tenant is created by terraform ACI provider"
	}

	resource "aci_application_profile" "foo_app_profile" {
		tenant_dn   = aci_tenant.tenant_for_subnet.id
		name        = "foo_app_profile"
		annotation  = "tag"
		description = "from terraform"
		name_alias  = "test_ap"
		prio        = "unspecified"
	  }

	  resource "aci_application_epg" "foo_epg" {
		application_profile_dn = aci_application_profile.foo_app_profile.id
		name                   = "foo_epg"
		description            = "from terraform"
		annotation             = "tag_epg"
		exception_tag          = "0"
		flood_on_encap         = "disabled"
		fwd_ctrl               = "none"
		has_mcast_source       = "no"
		is_attr_based_epg      = "no"
		match_t                = "AtleastOne"
		name_alias             = "alias_epg"
		pc_enf_pref            = "unenforced"
		pref_gr_memb           = "exclude"
		prio                   = "unspecified"
		shutdown               = "no"
	  }

	  resource "aci_subnet" "foo_epg_subnet_msnlb_mcast_igmp" {
		parent_dn   = aci_application_epg.foo_epg.id
		ip          = "10.0.3.29/32"
		scope       = ["private"]
		description = "from terraform"
		ctrl        = ["no-default-gateway"]
		preferred   = "no"
		virtual     = "yes"
		msnlb {
		  mode  = "mode-mcast-igmp"
		  group = "%s"
		}
	  }

	`, ms_nlb_group_value)
}

func testAccCheckAciSubnetEPGsMsNlbIgmpExists(
	name string,
	subnet *models.Subnet,
	ms_nlb *models.NlbEndpoint,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Subnet %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subnet dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		subnetFound := models.SubnetFromContainer(cont)
		if subnetFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Subnet %s not found", rs.Primary.ID)
		}
		*subnet = *subnetFound

		// fvEpNlb - Beginning of Read
		ms_nlb_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnfvEpNlb)
		ms_nlb_cont, err := client.Get(ms_nlb_dn)
		if err != nil {
			return err
		}

		ms_nlb_found := models.NlbEndpointFromContainer(ms_nlb_cont)
		if ms_nlb_found.DistinguishedName != ms_nlb_dn {
			return fmt.Errorf("MSNLB %s object not found", ms_nlb_dn)
		}
		*ms_nlb = *ms_nlb_found
		// fvEpNlb - Read finished successfully

		return nil
	}
}

func testAccCheckAciSubnetEPGsMsNlbIgmpDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_subnet" {
			cont, err := client.Get(rs.Primary.ID)
			subnet := models.SubnetFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Subnet %s Still exists", subnet.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSubnetEPGsMsNlbIgmpAttributes(
	ms_nlb_group_value string,
	ms_nlb *models.NlbEndpoint,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if ms_nlb_group_value != ms_nlb.Group {
			return fmt.Errorf("Bad MSNLB Group Address %s", ms_nlb.Group)
		}

		if "00:00:00:00:00:00" != ms_nlb.Mac {
			return fmt.Errorf("Bad MSNLB default MAC Address %s", ms_nlb.Mac)
		}

		if "mode-mcast-igmp" != ms_nlb.Mode {
			return fmt.Errorf("MSNLB Mode is not matching with IGMP: %s", ms_nlb.Mac)
		}

		return nil
	}
}

// MS NLB Check only for IGMP

// EP Reachability
func TestAccAciSubnetEPGsEpReachability_Basic(t *testing.T) {
	var subnet models.Subnet
	var ep_reachability models.EpReachability
	var next_hop_ep_reachability models.NexthopEpPReachability

	next_hop_addr_value := "10.0.3.20"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetEPGsEpReachabilityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubnetEPGsEpReachabilityConfig_basic(
					next_hop_addr_value,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetEPGsEpReachabilityExists(
						"aci_subnet.foo_epg_subnet_ep_reachability",
						next_hop_addr_value,
						&subnet,
						&ep_reachability,
						&next_hop_ep_reachability,
					),
					testAccCheckAciSubnetEPGsEpReachabilityAttributes(
						next_hop_addr_value,
						&next_hop_ep_reachability,
					),
				),
			},
		},
	})
}

func TestAccAciSubnetEPGsEpReachability_update(t *testing.T) {
	var subnet models.Subnet
	var ep_reachability models.EpReachability
	var next_hop_ep_reachability models.NexthopEpPReachability

	next_hop_addr_value := "10.0.3.20"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetEPGsEpReachabilityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubnetEPGsEpReachabilityConfig_basic(
					next_hop_addr_value,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetEPGsEpReachabilityExists(
						"aci_subnet.foo_epg_subnet_ep_reachability",
						next_hop_addr_value,
						&subnet,
						&ep_reachability,
						&next_hop_ep_reachability,
					),
					testAccCheckAciSubnetEPGsEpReachabilityAttributes(
						next_hop_addr_value,
						&next_hop_ep_reachability,
					),
				),
			},
			{
				Config: testAccCheckAciSubnetEPGsEpReachabilityConfig_basic(
					next_hop_addr_value,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetEPGsEpReachabilityExists(
						"aci_subnet.foo_epg_subnet_ep_reachability",
						next_hop_addr_value,
						&subnet,
						&ep_reachability,
						&next_hop_ep_reachability,
					),
					testAccCheckAciSubnetEPGsEpReachabilityAttributes(
						next_hop_addr_value,
						&next_hop_ep_reachability,
					),
				),
			},
		},
	})
}

func testAccCheckAciSubnetEPGsEpReachabilityConfig_basic(
	next_hop_addr_value string,
) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "tenant_for_subnet" {
		name        = "tenant_for_subnet"
		description = "This tenant is created by terraform ACI provider"
	}

	resource "aci_application_profile" "foo_app_profile" {
		tenant_dn   = aci_tenant.tenant_for_subnet.id
		name        = "foo_app_profile"
		annotation  = "tag"
		description = "from terraform"
		name_alias  = "test_ap"
		prio        = "unspecified"
	  }

	  resource "aci_application_epg" "foo_epg" {
		application_profile_dn = aci_application_profile.foo_app_profile.id
		name                   = "foo_epg"
		description            = "from terraform"
		annotation             = "tag_epg"
		exception_tag          = "0"
		flood_on_encap         = "disabled"
		fwd_ctrl               = "none"
		has_mcast_source       = "no"
		is_attr_based_epg      = "no"
		match_t                = "AtleastOne"
		name_alias             = "alias_epg"
		pc_enf_pref            = "unenforced"
		pref_gr_memb           = "exclude"
		prio                   = "unspecified"
		shutdown               = "no"
	  }

	  resource "aci_subnet" "foo_epg_subnet_ep_reachability" {
		parent_dn     = aci_application_epg.foo_epg.id
		ip            = "10.0.3.29/32"
		scope         = ["private"]
		description   = "from terraform"
		ctrl          = ["no-default-gateway"]
		preferred     = "no"
		virtual       = "yes"
		next_hop_addr = "%s"
	  }

	`, next_hop_addr_value)
}

func testAccCheckAciSubnetEPGsEpReachabilityExists(
	name,
	next_hop_addr_value string,
	subnet *models.Subnet,
	ep_reachability *models.EpReachability,
	next_hop_ep_reachability *models.NexthopEpPReachability,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Subnet %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subnet dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		subnetFound := models.SubnetFromContainer(cont)
		if subnetFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Subnet %s not found", rs.Primary.ID)
		}
		*subnet = *subnetFound

		// fvEpReachability and ipNexthopEpP - Beginning of Read
		next_hop_addr_dn := rs.Primary.ID + fmt.Sprintf("/epReach/"+models.RnipNexthopEpP, next_hop_addr_value)

		next_hop_addr_cont, err := client.Get(next_hop_addr_dn)
		if err != nil {
			return err
		}

		next_hop_addr_found := models.NexthopEpPReachabilityFromContainer(next_hop_addr_cont)
		if next_hop_addr_found.DistinguishedName != next_hop_addr_dn {
			return fmt.Errorf("EP Reachability Next Hop Address Object %s not found", next_hop_addr_dn)
		}
		*next_hop_ep_reachability = *next_hop_addr_found
		// fvEpReachability and ipNexthopEpP - Beginning of Read

		return nil
	}
}

func testAccCheckAciSubnetEPGsEpReachabilityDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_subnet" {
			cont, err := client.Get(rs.Primary.ID)
			subnet := models.SubnetFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Subnet %s Still exists", subnet.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSubnetEPGsEpReachabilityAttributes(
	next_hop_addr_value string,
	next_hop_ep_reachability *models.NexthopEpPReachability,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if next_hop_addr_value != next_hop_ep_reachability.NhAddr {
			return fmt.Errorf("Bad EP Reachability Next Hop Address %s", next_hop_ep_reachability.NhAddr)
		}

		return nil
	}
}

// EP Reachability

// Anycast MAC
func TestAccAciSubnetEPGsAnycastMac_Basic(t *testing.T) {
	var subnet models.Subnet
	var anycast_mac models.AnycastEndpoint

	anycast_mac_value := "00:00:00:11:11:11"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetEPGsAnycastMacDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubnetEPGsAnycastMacConfig_basic(
					anycast_mac_value,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetEPGsAnycastMacExists(
						"aci_subnet.foo_epg_subnet_ep_reachability",
						anycast_mac_value,
						&subnet,
						&anycast_mac,
					),
					testAccCheckAciSubnetEPGsAnycastMacAttributes(
						anycast_mac_value,
						&anycast_mac,
					),
				),
			},
		},
	})
}

func TestAccAciSubnetEPGsAnycastMac_update(t *testing.T) {
	var subnet models.Subnet
	var anycast_mac models.AnycastEndpoint

	anycast_mac_value := "11:11:11:00:00:00"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetEPGsAnycastMacDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubnetEPGsAnycastMacConfig_basic(
					anycast_mac_value,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetEPGsAnycastMacExists(
						"aci_subnet.foo_epg_subnet_ep_reachability",
						anycast_mac_value,
						&subnet,
						&anycast_mac,
					),
					testAccCheckAciSubnetEPGsAnycastMacAttributes(
						anycast_mac_value,
						&anycast_mac,
					),
				),
			},
			{
				Config: testAccCheckAciSubnetEPGsAnycastMacConfig_basic(
					anycast_mac_value,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetEPGsAnycastMacExists(
						"aci_subnet.foo_epg_subnet_ep_reachability",
						anycast_mac_value,
						&subnet,
						&anycast_mac,
					),
					testAccCheckAciSubnetEPGsAnycastMacAttributes(
						anycast_mac_value,
						&anycast_mac,
					),
				),
			},
		},
	})
}

func testAccCheckAciSubnetEPGsAnycastMacConfig_basic(
	anycast_mac_value string,
) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "tenant_for_subnet" {
		name        = "tenant_for_subnet"
		description = "This tenant is created by terraform ACI provider"
	}

	resource "aci_application_profile" "foo_app_profile" {
		tenant_dn   = aci_tenant.tenant_for_subnet.id
		name        = "foo_app_profile"
		annotation  = "tag"
		description = "from terraform"
		name_alias  = "test_ap"
		prio        = "unspecified"
	  }

	  resource "aci_application_epg" "foo_epg" {
		application_profile_dn = aci_application_profile.foo_app_profile.id
		name                   = "foo_epg"
		description            = "from terraform"
		annotation             = "tag_epg"
		exception_tag          = "0"
		flood_on_encap         = "disabled"
		fwd_ctrl               = "none"
		has_mcast_source       = "no"
		is_attr_based_epg      = "no"
		match_t                = "AtleastOne"
		name_alias             = "alias_epg"
		pc_enf_pref            = "unenforced"
		pref_gr_memb           = "exclude"
		prio                   = "unspecified"
		shutdown               = "no"
	  }

	  resource "aci_subnet" "foo_epg_subnet_ep_reachability" {
		parent_dn     = aci_application_epg.foo_epg.id
		ip            = "10.0.3.29/32"
		scope         = ["private"]
		description   = "from terraform"
		ctrl          = ["no-default-gateway"]
		preferred     = "no"
		virtual       = "yes"
		anycast_mac = "%s"
	  }

	`, anycast_mac_value)
}

func testAccCheckAciSubnetEPGsAnycastMacExists(
	name,
	anycast_mac_value string,
	subnet *models.Subnet,
	anycast_mac *models.AnycastEndpoint,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Subnet %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subnet dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		subnetFound := models.SubnetFromContainer(cont)
		if subnetFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Subnet %s not found", rs.Primary.ID)
		}
		*subnet = *subnetFound

		// fvEpAnycast - Beginning of Read
		anycast_mac_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnfvEpAnycast, anycast_mac_value)
		anycast_mac_cont, err := client.Get(anycast_mac_dn)
		if err != nil {
			return err
		}

		anycast_mac_found := models.AnycastEndpointFromContainer(anycast_mac_cont)
		if anycast_mac_found.DistinguishedName != anycast_mac_dn {
			return fmt.Errorf("Anycast MAC Object %s not found", anycast_mac_dn)
		}
		*anycast_mac = *anycast_mac_found
		// fvEpAnycast - Beginning of Read

		return nil
	}
}

func testAccCheckAciSubnetEPGsAnycastMacDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_subnet" {
			cont, err := client.Get(rs.Primary.ID)
			subnet := models.SubnetFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Subnet %s Still exists", subnet.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSubnetEPGsAnycastMacAttributes(
	anycast_mac_value string,
	anycast_mac *models.AnycastEndpoint,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if anycast_mac_value != anycast_mac.Mac {
			return fmt.Errorf("Bad Anycast MAC Address %s", anycast_mac.Mac)
		}

		return nil
	}
}

// Anycast MAC
