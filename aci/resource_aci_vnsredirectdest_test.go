package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciDestinationofredirectedtraffic_Basic(t *testing.T) {
	var destinationofredirectedtraffic models.Destinationofredirectedtraffic
	description := "destinationofredirectedtraffic created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDestinationofredirectedtrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDestinationofredirectedtrafficConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDestinationofredirectedtrafficExists("aci_destination_of_redirected_traffic.example", &destinationofredirectedtraffic),
					testAccCheckAciDestinationofredirectedtrafficAttributes(description, &destinationofredirectedtraffic),
				),
			},
		},
	})
}

func TestAccAciDestinationofredirectedtraffic_update(t *testing.T) {
	var destinationofredirectedtraffic models.Destinationofredirectedtraffic
	description := "destinationofredirectedtraffic created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDestinationofredirectedtrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDestinationofredirectedtrafficConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDestinationofredirectedtrafficExists("aci_destination_of_redirected_traffic.example", &destinationofredirectedtraffic),
					testAccCheckAciDestinationofredirectedtrafficAttributes(description, &destinationofredirectedtraffic),
				),
			},
			{
				Config: testAccCheckAciDestinationofredirectedtrafficConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDestinationofredirectedtrafficExists("aci_destination_of_redirected_traffic.example", &destinationofredirectedtraffic),
					testAccCheckAciDestinationofredirectedtrafficAttributes(description, &destinationofredirectedtraffic),
				),
			},
		},
	})
}

func testAccCheckAciDestinationofredirectedtrafficConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_destination_of_redirected_traffic" "example" {
		service_redirect_policy_dn  = "${aci_service_redirect_policy.example.id}"
		ip  = "1.2.3.4"
		mac = "12:25:56:98:45:74"
		ip2 = "10.20.30.40"
		dest_name = "last"
		pod_id = "5"
		description = "%s"
	}
	`, description)
}

func testAccCheckAciDestinationofredirectedtrafficExists(name string, destinationofredirectedtraffic *models.Destinationofredirectedtraffic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Destination of redirected traffic %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Destination of redirected traffic dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		destinationofredirectedtrafficFound := models.DestinationofredirectedtrafficFromContainer(cont)
		if destinationofredirectedtrafficFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Destination of redirected traffic %s not found", rs.Primary.ID)
		}
		*destinationofredirectedtraffic = *destinationofredirectedtrafficFound
		return nil
	}
}

func testAccCheckAciDestinationofredirectedtrafficDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_destination_of_redirected_traffic" {
			cont, err := client.Get(rs.Primary.ID)
			destinationofredirectedtraffic := models.DestinationofredirectedtrafficFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Destination of redirected traffic %s Still exists", destinationofredirectedtraffic.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciDestinationofredirectedtrafficAttributes(description string, destinationofredirectedtraffic *models.Destinationofredirectedtraffic) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != destinationofredirectedtraffic.Description {
			return fmt.Errorf("Bad destinationofredirectedtraffic Description %s", destinationofredirectedtraffic.Description)
		}

		if "1.2.3.4" != destinationofredirectedtraffic.Ip {
			return fmt.Errorf("Bad destinationofredirectedtraffic ip %s", destinationofredirectedtraffic.Ip)
		}

		if "last" != destinationofredirectedtraffic.DestName {
			return fmt.Errorf("Bad destinationofredirectedtraffic dest_name %s", destinationofredirectedtraffic.DestName)
		}

		if "10.20.30.40" != destinationofredirectedtraffic.Ip2 {
			return fmt.Errorf("Bad destinationofredirectedtraffic ip2 %s", destinationofredirectedtraffic.Ip2)
		}

		if "12:25:56:98:45:74" != destinationofredirectedtraffic.Mac {
			return fmt.Errorf("Bad destinationofredirectedtraffic mac %s", destinationofredirectedtraffic.Mac)
		}

		if "5" != destinationofredirectedtraffic.PodId {
			return fmt.Errorf("Bad destinationofredirectedtraffic pod_id %s", destinationofredirectedtraffic.PodId)
		}

		return nil
	}
}
