package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciConnection_Basic(t *testing.T) {
	var connection models.Connection
	description := "connection created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciConnectionConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConnectionExists("aci_connection.test", &connection),
					testAccCheckAciConnectionAttributes(description, &connection),
				),
			},
		},
	})
}

func TestAccAciConnection_update(t *testing.T) {
	var connection models.Connection
	description := "connection created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciConnectionConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConnectionExists("aci_connection.test", &connection),
					testAccCheckAciConnectionAttributes(description, &connection),
				),
			},
			{
				Config: testAccCheckAciConnectionConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConnectionExists("aci_connection.test", &connection),
					testAccCheckAciConnectionAttributes(description, &connection),
				),
			},
		},
	})
}

func testAccCheckAciConnectionConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = "uni/tn-phase2/AbsGraph-second"
		name  = "conn2"
		adj_type  = "L3"
		description = "%s"
		annotation  = "example"
		conn_dir  = "provider"
		conn_type  = "internal"
		direct_connect  = "yes"
		name_alias  = "example"
		unicast_route  = "yes"
	}
	`, description)
}

func testAccCheckAciConnectionExists(name string, connection *models.Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Connection %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Connection dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		connectionFound := models.ConnectionFromContainer(cont)
		if connectionFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Connection %s not found", rs.Primary.ID)
		}
		*connection = *connectionFound
		return nil
	}
}

func testAccCheckAciConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_connection" {
			cont, err := client.Get(rs.Primary.ID)
			connection := models.ConnectionFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Connection %s Still exists", connection.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciConnectionAttributes(description string, connection *models.Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != connection.Description {
			return fmt.Errorf("Bad connection Description %s", connection.Description)
		}

		if "conn2" != connection.Name {
			return fmt.Errorf("Bad connection name %s", connection.Name)
		}

		if "L3" != connection.AdjType {
			return fmt.Errorf("Bad connection adj_type %s", connection.AdjType)
		}

		if "example" != connection.Annotation {
			return fmt.Errorf("Bad connection annotation %s", connection.Annotation)
		}

		if "provider" != connection.ConnDir {
			return fmt.Errorf("Bad connection conn_dir %s", connection.ConnDir)
		}

		if "internal" != connection.ConnType {
			return fmt.Errorf("Bad connection conn_type %s", connection.ConnType)
		}

		if "yes" != connection.DirectConnect {
			return fmt.Errorf("Bad connection direct_connect %s", connection.DirectConnect)
		}

		if "example" != connection.NameAlias {
			return fmt.Errorf("Bad connection name_alias %s", connection.NameAlias)
		}

		if "yes" != connection.UnicastRoute {
			return fmt.Errorf("Bad connection unicast_route %s", connection.UnicastRoute)
		}

		return nil
	}
}
