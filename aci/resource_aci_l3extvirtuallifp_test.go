package aci

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/ciscoecosystem/aci-go-client/client"
// 	"github.com/ciscoecosystem/aci-go-client/models"
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/terraform"
// )

// func TestAccAciLogicalInterfaceProfile_Basic(t *testing.T) {
// 	var logical_interface_profile models.LogicalInterfaceProfile
// 	description := "logical_interface_profile created while acceptance testing"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckAciLogicalInterfaceProfileDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccCheckAciLogicalInterfaceProfileConfig_basic(description),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciLogicalInterfaceProfileExists("aci_logical_interface_profile.foological_interface_profile", &logical_interface_profile),
// 					testAccCheckAciLogicalInterfaceProfileAttributes(description, &logical_interface_profile),
// 				),
// 			},
// 			{
// 				ResourceName:      "aci_logical_interface_profile",
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

// func TestAccAciLogicalInterfaceProfile_update(t *testing.T) {
// 	var logical_interface_profile models.LogicalInterfaceProfile
// 	description := "logical_interface_profile created while acceptance testing"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckAciLogicalInterfaceProfileDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccCheckAciLogicalInterfaceProfileConfig_basic(description),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciLogicalInterfaceProfileExists("aci_logical_interface_profile.foological_interface_profile", &logical_interface_profile),
// 					testAccCheckAciLogicalInterfaceProfileAttributes(description, &logical_interface_profile),
// 				),
// 			},
// 			{
// 				Config: testAccCheckAciLogicalInterfaceProfileConfig_basic(description),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciLogicalInterfaceProfileExists("aci_logical_interface_profile.foological_interface_profile", &logical_interface_profile),
// 					testAccCheckAciLogicalInterfaceProfileAttributes(description, &logical_interface_profile),
// 				),
// 			},
// 		},
// 	})
// }

// func testAccCheckAciLogicalInterfaceProfileConfig_basic(description string) string {
// 	return fmt.Sprintf(`

// 	resource "aci_logical_interface_profile" "foological_interface_profile" {
// 		  logical_interface_profile_dn  = "${aci_logical_interface_profile.example.id}"
// 		description = "%s"

// 		nodeDn  = "example"

// 		encap  = "example"
// 		  addr  = "example"
// 		  annotation  = "example"
// 		  autostate  = "example"
// 		  encap_scope  = "example"
// 		  if_inst_t  = "example"
// 		  ipv6_dad  = "example"
// 		  ll_addr  = "example"
// 		  mac  = "example"
// 		  mode  = "example"
// 		  mtu  = "example"
// 		  node_dn  = "example"
// 		  target_dscp  = "example"
// 		  userdom  = "example"
// 		}
// 	`, description)
// }

// func testAccCheckAciLogicalInterfaceProfileExists(name string, logical_interface_profile *models.LogicalInterfaceProfile) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[name]

// 		if !ok {
// 			return fmt.Errorf("Logical Interface Profile %s not found", name)
// 		}

// 		if rs.Primary.ID == "" {
// 			return fmt.Errorf("No Logical Interface Profile dn was set")
// 		}

// 		client := testAccProvider.Meta().(*client.Client)

// 		cont, err := client.Get(rs.Primary.ID)
// 		if err != nil {
// 			return err
// 		}

// 		logical_interface_profileFound := models.LogicalInterfaceProfileFromContainer(cont)
// 		if logical_interface_profileFound.DistinguishedName != rs.Primary.ID {
// 			return fmt.Errorf("Logical Interface Profile %s not found", rs.Primary.ID)
// 		}
// 		*logical_interface_profile = *logical_interface_profileFound
// 		return nil
// 	}
// }

// func testAccCheckAciLogicalInterfaceProfileDestroy(s *terraform.State) error {
// 	client := testAccProvider.Meta().(*client.Client)

// 	for _, rs := range s.RootModule().Resources {

// 		if rs.Type == "aci_logical_interface_profile" {
// 			cont, err := client.Get(rs.Primary.ID)
// 			logical_interface_profile := models.LogicalInterfaceProfileFromContainer(cont)
// 			if err == nil {
// 				return fmt.Errorf("Logical Interface Profile %s Still exists", logical_interface_profile.DistinguishedName)
// 			}

// 		} else {
// 			continue
// 		}
// 	}

// 	return nil
// }

// func testAccCheckAciLogicalInterfaceProfileAttributes(description string, logical_interface_profile *models.LogicalInterfaceProfile) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {

// 		if description != logical_interface_profile.Description {
// 			return fmt.Errorf("Bad logical_interface_profile Description %s", logical_interface_profile.Description)
// 		}

// 		if "example" != logical_interface_profile.NodeDn {
// 			return fmt.Errorf("Bad logical_interface_profile node_dn %s", logical_interface_profile.NodeDn)
// 		}

// 		if "example" != logical_interface_profile.Encap {
// 			return fmt.Errorf("Bad logical_interface_profile encap %s", logical_interface_profile.Encap)
// 		}

// 		if "example" != logical_interface_profile.Addr {
// 			return fmt.Errorf("Bad logical_interface_profile addr %s", logical_interface_profile.Addr)
// 		}

// 		if "example" != logical_interface_profile.Annotation {
// 			return fmt.Errorf("Bad logical_interface_profile annotation %s", logical_interface_profile.Annotation)
// 		}

// 		if "example" != logical_interface_profile.Autostate {
// 			return fmt.Errorf("Bad logical_interface_profile autostate %s", logical_interface_profile.Autostate)
// 		}

// 		if "example" != logical_interface_profile.EncapScope {
// 			return fmt.Errorf("Bad logical_interface_profile encap_scope %s", logical_interface_profile.EncapScope)
// 		}

// 		if "example" != logical_interface_profile.IfInstT {
// 			return fmt.Errorf("Bad logical_interface_profile if_inst_t %s", logical_interface_profile.IfInstT)
// 		}

// 		if "example" != logical_interface_profile.Ipv6Dad {
// 			return fmt.Errorf("Bad logical_interface_profile ipv6_dad %s", logical_interface_profile.Ipv6Dad)
// 		}

// 		if "example" != logical_interface_profile.LlAddr {
// 			return fmt.Errorf("Bad logical_interface_profile ll_addr %s", logical_interface_profile.LlAddr)
// 		}

// 		if "example" != logical_interface_profile.Mac {
// 			return fmt.Errorf("Bad logical_interface_profile mac %s", logical_interface_profile.Mac)
// 		}

// 		if "example" != logical_interface_profile.Mode {
// 			return fmt.Errorf("Bad logical_interface_profile mode %s", logical_interface_profile.Mode)
// 		}

// 		if "example" != logical_interface_profile.Mtu {
// 			return fmt.Errorf("Bad logical_interface_profile mtu %s", logical_interface_profile.Mtu)
// 		}

// 		if "example" != logical_interface_profile.NodeDn {
// 			return fmt.Errorf("Bad logical_interface_profile node_dn %s", logical_interface_profile.NodeDn)
// 		}

// 		if "example" != logical_interface_profile.TargetDscp {
// 			return fmt.Errorf("Bad logical_interface_profile target_dscp %s", logical_interface_profile.TargetDscp)
// 		}

// 		if "example" != logical_interface_profile.Userdom {
// 			return fmt.Errorf("Bad logical_interface_profile userdom %s", logical_interface_profile.Userdom)
// 		}

// 		return nil
// 	}
// }
