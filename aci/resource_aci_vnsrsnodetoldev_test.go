package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)





func TestAccAciRelationfromaAbsNodetoanLDev_Basic(t *testing.T) {
	var relationfroma_abs_nodetoan_l_dev models.RelationfromaAbsNodetoanLDev
	description := "relationfroma_abs_nodetoan_l_dev created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:	  func(){ testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRelationfromaAbsNodetoanLDevDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRelationfromaAbsNodetoanLDevConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRelationfromaAbsNodetoanLDevExists("aci_relationfroma_abs_nodetoan_l_dev.foorelationfroma_abs_nodetoan_l_dev", &relationfroma_abs_nodetoan_l_dev),
					testAccCheckAciRelationfromaAbsNodetoanLDevAttributes(description, &relationfroma_abs_nodetoan_l_dev),
				),
			},
			{
				ResourceName:      "aci_relationfroma_abs_nodetoan_l_dev",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciRelationfromaAbsNodetoanLDev_update(t *testing.T) {
	var relationfroma_abs_nodetoan_l_dev models.RelationfromaAbsNodetoanLDev
	description := "relationfroma_abs_nodetoan_l_dev created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t)},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRelationfromaAbsNodetoanLDevDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRelationfromaAbsNodetoanLDevConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRelationfromaAbsNodetoanLDevExists("aci_relationfroma_abs_nodetoan_l_dev.foorelationfroma_abs_nodetoan_l_dev", &relationfroma_abs_nodetoan_l_dev),
					testAccCheckAciRelationfromaAbsNodetoanLDevAttributes(description, &relationfroma_abs_nodetoan_l_dev),
				),
			},
			{
				Config: testAccCheckAciRelationfromaAbsNodetoanLDevConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRelationfromaAbsNodetoanLDevExists("aci_relationfroma_abs_nodetoan_l_dev.foorelationfroma_abs_nodetoan_l_dev", &relationfroma_abs_nodetoan_l_dev),
					testAccCheckAciRelationfromaAbsNodetoanLDevAttributes(description, &relationfroma_abs_nodetoan_l_dev),
				),
			},
		},
	})
}

func testAccCheckAciRelationfromaAbsNodetoanLDevConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_relationfroma_abs_nodetoan_l_dev" "foorelationfroma_abs_nodetoan_l_dev" {
		  function_node_dn  = "${aci_function_node.example.id}"
		description = "%s"
		  annotation  = "example"
		  t_dn  = "example"
		}
	`, description)
}

func testAccCheckAciRelationfromaAbsNodetoanLDevExists(name string, relationfroma_abs_nodetoan_l_dev *models.RelationfromaAbsNodetoanLDev) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Relation from a AbsNode to an LDev %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Relation from a AbsNode to an LDev dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		relationfroma_abs_nodetoan_l_devFound := models.RelationfromaAbsNodetoanLDevFromContainer(cont)
		if relationfroma_abs_nodetoan_l_devFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Relation from a AbsNode to an LDev %s not found", rs.Primary.ID)
		}
		*relationfroma_abs_nodetoan_l_dev = *relationfroma_abs_nodetoan_l_devFound
		return nil
	}
}

func testAccCheckAciRelationfromaAbsNodetoanLDevDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		
		 if rs.Type == "aci_relationfroma_abs_nodetoan_l_dev" {
			cont,err := client.Get(rs.Primary.ID)
			relationfroma_abs_nodetoan_l_dev := models.RelationfromaAbsNodetoanLDevFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Relation from a AbsNode to an LDev %s Still exists",relationfroma_abs_nodetoan_l_dev.DistinguishedName)
			}

		}else{
			continue
		}
	}

	return nil
}

func testAccCheckAciRelationfromaAbsNodetoanLDevAttributes(description string, relationfroma_abs_nodetoan_l_dev  *models.RelationfromaAbsNodetoanLDev) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != relationfroma_abs_nodetoan_l_dev.Description {
			return fmt.Errorf("Bad relationfroma_abs_nodetoan_l_dev Description %s", relationfroma_abs_nodetoan_l_dev.Description)
		}
		

	
    
		if "example" != relationfroma_abs_nodetoan_l_dev.Annotation {
			return fmt.Errorf("Bad relationfroma_abs_nodetoan_l_dev annotation %s", relationfroma_abs_nodetoan_l_dev.Annotation)
		}
	
    
		if "example" != relationfroma_abs_nodetoan_l_dev.TDn {
			return fmt.Errorf("Bad relationfroma_abs_nodetoan_l_dev t_dn %s", relationfroma_abs_nodetoan_l_dev.TDn)
		}
	
    

	return nil
	}
}
