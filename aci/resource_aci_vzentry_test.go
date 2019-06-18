package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAciFilterEntry_Basic(t *testing.T) {
	var filter_entry models.FilterEntry
	description := "filter_entry created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFilterEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFilterEntryConfig_basic(description, "http"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists("aci_filter_entry.foofilter_entry", &filter_entry),
					testAccCheckAciFilterEntryAttributes(description, "http", &filter_entry),
				),
			},
			{
				ResourceName:      "aci_filter_entry",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciFilterEntry_update(t *testing.T) {
	var filter_entry models.FilterEntry
	description := "filter_entry created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFilterEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFilterEntryConfig_basic(description, "http"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists("aci_filter_entry.foofilter_entry", &filter_entry),
					testAccCheckAciFilterEntryAttributes(description, "http", &filter_entry),
				),
			},
			{
				Config: testAccCheckAciFilterEntryConfig_basic(description, "https"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists("aci_filter_entry.foofilter_entry", &filter_entry),
					testAccCheckAciFilterEntryAttributes(description, "https", &filter_entry),
				),
			},
		},
	})
}

func testAccCheckAciFilterEntryConfig_basic(description, d_from_port string) string {
	return fmt.Sprintf(`

	resource "aci_filter_entry" "foofilter_entry" {
		filter_dn     = "${aci_filter.example.id}"
		description   = "%s"
		name          = "demo_entry"
		annotation    = "tag_entry"
		apply_to_frag = "no"
		arp_opc       = "unspecified"
		d_from_port   = "%s"
		d_to_port     = "unspecified"
		ether_t       = "ipv4"
		icmpv4_t      = "unspecified"
		icmpv6_t      = "unspecified"
		match_dscp    = "CS0"
		name_alias    = "alias_entry"
		prot          = "icmp"
		s_from_port   = "0"
		s_to_port     = "0"
		stateful      = "no"
		tcp_rules     = "ack"
	}
	  
	`, description, d_from_port)
}

func testAccCheckAciFilterEntryExists(name string, filter_entry *models.FilterEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Filter Entry %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Filter Entry dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		filter_entryFound := models.FilterEntryFromContainer(cont)
		if filter_entryFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Filter Entry %s not found", rs.Primary.ID)
		}
		*filter_entry = *filter_entryFound
		return nil
	}
}

func testAccCheckAciFilterEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_filter_entry" {
			cont, err := client.Get(rs.Primary.ID)
			filter_entry := models.FilterEntryFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Filter Entry %s Still exists", filter_entry.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciFilterEntryAttributes(description, d_from_port string, filter_entry *models.FilterEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != filter_entry.Description {
			return fmt.Errorf("Bad filter_entry Description %s", filter_entry.Description)
		}

		if "demo_entry" != filter_entry.Name {
			return fmt.Errorf("Bad filter_entry name %s", filter_entry.Name)
		}

		if "tag_entry" != filter_entry.Annotation {
			return fmt.Errorf("Bad filter_entry annotation %s", filter_entry.Annotation)
		}

		if "no" != filter_entry.ApplyToFrag {
			return fmt.Errorf("Bad filter_entry apply_to_frag %s", filter_entry.ApplyToFrag)
		}

		if "unspecified" != filter_entry.ArpOpc {
			return fmt.Errorf("Bad filter_entry arp_opc %s", filter_entry.ArpOpc)
		}

		if d_from_port != filter_entry.DFromPort {
			return fmt.Errorf("Bad filter_entry d_from_port %s", filter_entry.DFromPort)
		}

		if "unspecified" != filter_entry.DToPort {
			return fmt.Errorf("Bad filter_entry d_to_port %s", filter_entry.DToPort)
		}

		if "ipv4" != filter_entry.EtherT {
			return fmt.Errorf("Bad filter_entry ether_t %s", filter_entry.EtherT)
		}

		if "unspecified" != filter_entry.Icmpv4T {
			return fmt.Errorf("Bad filter_entry icmpv4_t %s", filter_entry.Icmpv4T)
		}

		if "unspecified" != filter_entry.Icmpv6T {
			return fmt.Errorf("Bad filter_entry icmpv6_t %s", filter_entry.Icmpv6T)
		}

		if "CS0" != filter_entry.MatchDscp {
			return fmt.Errorf("Bad filter_entry match_dscp %s", filter_entry.MatchDscp)
		}

		if "alias_entry" != filter_entry.NameAlias {
			return fmt.Errorf("Bad filter_entry name_alias %s", filter_entry.NameAlias)
		}

		if "icmp" != filter_entry.Prot {
			return fmt.Errorf("Bad filter_entry prot %s", filter_entry.Prot)
		}

		if "0" != filter_entry.SFromPort {
			return fmt.Errorf("Bad filter_entry s_from_port %s", filter_entry.SFromPort)
		}

		if "0" != filter_entry.SToPort {
			return fmt.Errorf("Bad filter_entry s_to_port %s", filter_entry.SToPort)
		}

		if "no" != filter_entry.Stateful {
			return fmt.Errorf("Bad filter_entry stateful %s", filter_entry.Stateful)
		}

		if "ack" != filter_entry.TcpRules {
			return fmt.Errorf("Bad filter_entry tcp_rules %s", filter_entry.TcpRules)
		}

		return nil
	}
}
