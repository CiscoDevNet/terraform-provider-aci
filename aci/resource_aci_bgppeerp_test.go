package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBgpPeerConnectivityProfile_Basic(t *testing.T) {
	var bgp_peer_connectivity_profile models.BgpPeerConnectivityProfile
	description := "peer_connectivity_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBgpPeerConnectivityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBgpPeerConnectivityProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpPeerConnectivityProfileExists("aci_bgp_peer_connectivity_profile.example", &bgp_peer_connectivity_profile),
					testAccCheckAciBgpPeerConnectivityProfileAttributes(description, &bgp_peer_connectivity_profile),
				),
			},
		},
	})
}

func TestAccAciBgpPeerConnectivityProfile_update(t *testing.T) {
	var bgp_peer_connectivity_profile models.BgpPeerConnectivityProfile
	description := "peer_connectivity_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBgpPeerConnectivityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBgpPeerConnectivityProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpPeerConnectivityProfileExists("aci_bgp_peer_connectivity_profile.example", &bgp_peer_connectivity_profile),
					testAccCheckAciBgpPeerConnectivityProfileAttributes(description, &bgp_peer_connectivity_profile),
				),
			},
			{
				Config: testAccCheckAciBgpPeerConnectivityProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpPeerConnectivityProfileExists("aci_bgp_peer_connectivity_profile.example", &bgp_peer_connectivity_profile),
					testAccCheckAciBgpPeerConnectivityProfileAttributes(description, &bgp_peer_connectivity_profile),
				),
			},
		},
	})
}

func testAccCheckAciBgpPeerConnectivityProfileConfig_basic(description string) string {	

	resource "aci_bgp_peer_connectivity_profile" "example" {
		logical_node_profile_dn = aci_logical_node_profile.example.id
		addr                    = "10.0.0.1"
		addr_t_ctrl             = "af-mcast, af-ucast"
		allowed_self_as_cnt     = "3"
		description 			= "%s"
		annotation              = "example"
		ctrl                    = "allow-self-as"
		name_alias              = "example"
		peer_ctrl               = "bfd"
		private_a_sctrl         = "remove-all,remove-exclusive"
		ttl                     = "1"
		weight                  = "1"
		as_number               = "1"
		local_asn               = "15"
		local_asn_propagate     = "dual-as"
	  }
	  
	`, description)
}

func testAccCheckAciBgpPeerConnectivityProfileExists(name string, bgp_peer_connectivity_profile *models.BgpPeerConnectivityProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("BgpPeerConnectivityProfile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No BgpPeerConnectivityProfile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bgp_peer_connectivity_profileFound := models.BgpPeerConnectivityProfileFromContainer(cont)
		if bgp_peer_connectivity_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("BgpPeerConnectivityProfile %s not found", rs.Primary.ID)
		}
		*bgp_peer_connectivity_profile = *bgp_peer_connectivity_profileFound
		return nil
	}
}

func testAccCheckAciBgpPeerConnectivityProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_peer_connectivity_profile" {
			cont, err := client.Get(rs.Primary.ID)
			bgp_peer_connectivity_profile := models.BgpPeerConnectivityProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("BgpPeerConnectivityProfile %s Still exists", bgp_peer_connectivity_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciBgpPeerConnectivityProfileAttributes(description string, bgp_peer_connectivity_profile *models.BgpPeerConnectivityProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != bgp_peer_connectivity_profile.Description {
			return fmt.Errorf("Bad bgp_peer_connectivity_profile Description %s", bgp_peer_connectivity_profile.Description)
		}

		if "10.0.0.1" != bgp_peer_connectivity_profile.Addr {
			return fmt.Errorf("Bad bgp_peer_connectivity_profile addr %s", bgp_peer_connectivity_profile.Addr)
		}

		if "af-mcast,af-ucast" != bgp_peer_connectivity_profile.AddrTCtrl {
			return fmt.Errorf("Bad bgp_peer_connectivity_profile addr_t_ctrl %s", bgp_peer_connectivity_profile.AddrTCtrl)
		}

		if "3" != bgp_peer_connectivity_profile.AllowedSelfAsCnt {
			return fmt.Errorf("Bad bgp_peer_connectivity_profile allowed_self_as_cnt %s", bgp_peer_connectivity_profile.AllowedSelfAsCnt)
		}

		if "example" != bgp_peer_connectivity_profile.Annotation {
			return fmt.Errorf("Bad bgp_peer_connectivity_profile annotation %s", bgp_peer_connectivity_profile.Annotation)
		}

		if "allow-self-as" != bgp_peer_connectivity_profile.Ctrl {
			return fmt.Errorf("Bad bgp_peer_connectivity_profile ctrl %s", bgp_peer_connectivity_profile.Ctrl)
		}

		if "example" != bgp_peer_connectivity_profile.NameAlias {
			return fmt.Errorf("Bad bgp_peer_connectivity_profile name_alias %s", bgp_peer_connectivity_profile.NameAlias)
		}

		if "bfd" != bgp_peer_connectivity_profile.PeerCtrl {
			return fmt.Errorf("Bad bgp_peer_connectivity_profile peer_ctrl %s", bgp_peer_connectivity_profile.PeerCtrl)
		}

		if "remove-all,remove-exclusive" != bgp_peer_connectivity_profile.PrivateASctrl {
			return fmt.Errorf("Bad bgp_peer_connectivity_profile private_a_sctrl %s", bgp_peer_connectivity_profile.PrivateASctrl)
		}

		if "1" != bgp_peer_connectivity_profile.Ttl {
			return fmt.Errorf("Bad bgp_peer_connectivity_profile ttl %s", bgp_peer_connectivity_profile.Ttl)
		}

		if "1" != bgp_peer_connectivity_profile.Weight {
			return fmt.Errorf("Bad bgp_peer_connectivity_profile weight %s", bgp_peer_connectivity_profile.Weight)
		}
		return nil
	}
}
