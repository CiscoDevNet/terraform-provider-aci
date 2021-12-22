package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLogicalInterfaceProfile_Basic(t *testing.T) {
	var logicalInterfaceProfile_default models.LogicalInterfaceProfile
	var logicalInterfaceProfile_updated models.LogicalInterfaceProfile
	resourceName := "aci_logical_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	prOther := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccLogicalInterfaceProfileWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileWithoutLogicalNodeProfileDn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLogicalInterfaceProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "logical_node_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttrSet(resourceName, "tag"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_l_if_p_to_netflow_monitor_pol.#", "0"),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag_prof"),
					resource.TestCheckResourceAttr(resourceName, "description", "Sample logical interface profile"),
					resource.TestCheckResourceAttr(resourceName, "logical_node_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "alias_prof"),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					resource.TestCheckResourceAttr(resourceName, "tag", "navy"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_l_if_p_to_netflow_monitor_pol.#", "0"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccLogicalInterfaceProfileConfigWithAnotherName(rName, rOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					testAccCheckAciLogicalInterfaceIdNotEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileConfig(rName),
			},
			{
				Config: CreateAccl3outsideConfigWithAnotherLogicalNodeProfileDn(prOther, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					testAccCheckAciLogicalInterfaceIdNotEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileConfigUpdateWithoutRequiredAttri(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileConfigUpdateWithInvalidName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of lifp-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccLogicalInterfaceProfileConfig(rName),
			},
		},
	})
}

func TestAccAciLogicalInterfaceProfile_Update(t *testing.T) {
	var logicalInterfaceProfile_default models.LogicalInterfaceProfile
	var logicalInterfaceProfile_updated models.LogicalInterfaceProfile
	resourceName := "aci_logical_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLogicalInterfaceProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_default),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "lime"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "lime"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "indigo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "indigo"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "purple"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "purple"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "gray"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "gray"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "yellow"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "yellow"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "white"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "white"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
		},
	})
}
func TestAccAciLogicalInterfaceProfile_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLogicalInterfaceProfileConfig(rName),
			},
			{
				Config:      CreateAccl3outsideConfigWithInvalidLogicalNodeProfiledn(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class l3extLIfP (.)+`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "prio", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name prio, class l3extLIfP (.)+`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name tag, class l3extLIfP (.)+`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLogicalInterfaceProfileConfig(rName),
			},
		},
	})
}

func TestAccAciLogicalInterfaceProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLogicalInterfaceProfileConfigMultiple(rName),
			},
		},
	})
}

func CreateAccLogicalInterfaceProfileConfigUpdateWithInvalidName(parentName, rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalInterfaceProfile update with invalid Name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
  resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

   resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
	   resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
  }
	`, parentName, parentName, parentName, rName)
	return resource
}

func testAccCheckLogicalInterfaceProfileExists(name string, logicalInterfaceProfile *models.LogicalInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LogicalInterfaceProfile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LogicalInterfaceProfile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		logicalInterfaceProfileFound := models.LogicalInterfaceProfileFromContainer(cont)
		if logicalInterfaceProfileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LogicalInterfaceProfile %s not found", rs.Primary.ID)
		}
		*logicalInterfaceProfile = *logicalInterfaceProfileFound
		return nil
	}
}

func testAccCheckAciLogicalInterfaceProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing LogicalInterfaceProfile destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_logical_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			aci := models.LogicalInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LogicalInterfaceProfile %s Still exists", aci.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLogicalInterfaceProfileIdEqual(logicalInterfaceProfile1, logicalInterfaceProfile2 *models.LogicalInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if logicalInterfaceProfile1.DistinguishedName != logicalInterfaceProfile2.DistinguishedName {
			return fmt.Errorf("LogicalInterfaceProfile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLogicalInterfaceIdNotEqual(logicalInterfaceProfile1, logicalInterfaceProfile2 *models.LogicalInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if logicalInterfaceProfile1.DistinguishedName == logicalInterfaceProfile2.DistinguishedName {
			return fmt.Errorf("LogicalInterfaceProfile DNs are equal")
		}
		return nil
	}
}

func CreateAccLogicalInterfaceProfileWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalInterfaceProfile creation without giving Name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
  resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

   resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
	   resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccLogicalInterfaceProfileWithoutLogicalNodeProfileDn(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalInterfaceProfile creation without giving LogicalNodeProfile dn")
	resource := fmt.Sprintf(`
	resource "aci_logical_interface_profile" "test" {
        name                    = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLogicalInterfaceProfileConfigUpdateWithoutRequiredAttri() string {
	fmt.Println("=== STEP  Basic: testing LogicalInterfaceProfile update without giving required Attributes")
	resource := fmt.Sprintf(`
    resource "aci_logical_interface_profile" "test" {
		description             = "Sample logical interface profile"
        annotation              = "tag_prof"
        name_alias              = "alias_prof"
        prio                    = "level1"
        tag                     = "navy"
  }
	`)
	return resource
}

func CreateAccLogicalInterfaceProfileConfigWithAnotherName(parentName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing LogicalInterfaceProfile creation with different LogicalInterfaceProfile name %s \n", rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
  resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

   resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
	   resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
  }
	`, parentName, parentName, parentName, rName)
	return resource
}

func CreateAccl3outsideConfigWithAnotherLogicalNodeProfileDn(parentName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing LogicalInterfaceProfile creation with different parent %s \n", parentName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
  resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

   resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
	   resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
  }
	`, parentName, parentName, parentName, rName)
	return resource
}

func CreateAccl3outsideConfigWithInvalidLogicalNodeProfiledn(rName string) string {
	fmt.Printf("=== STEP  Basic: testing LogicalInterfaceProfile creation with invalid Logical Node Profile dn \n")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
    resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_tenant.test.id
        name                    = "%s"
		description             = "Sample logical interface profile"
        annotation              = "tag_prof"
        name_alias              = "alias_prof"
        prio                    = "level1"
        tag                     = "navy"
  }
	`, rName, rName)
	return resource
}

func CreateAccLogicalInterfaceProfileConfig(rName string) string {
	fmt.Println("=== STEP  testing LogicalInterfaceProfile creation with required attributes")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
    resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

    resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
    resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
  }
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccLogicalInterfaceProfileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing logicalInterfaceProfile creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
    resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

    resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
    resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
		description             = "Sample logical interface profile"
        annotation              = "tag_prof"
        name_alias              = "alias_prof"
        prio                    = "level1"
        tag                     = "navy"
  }
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccLogicalInterfaceProfileUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
    resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

    resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
    resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name           = "%s"
		%s = "%s"
  }
	`, rName, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccLogicalInterfaceProfileConfigMultiple(rName string) string {
	fmt.Println("=== STEP  Creating Multiple LogicalInterfaceProfile")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
    resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

    resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
    resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
    }
	resource "aci_logical_interface_profile" "test1" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
    }
	resource "aci_logical_interface_profile" "test2" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
    }
	resource "aci_logical_interface_profile" "test3" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
    }
	`, rName, rName, rName, rName, rName+"1", rName+"2", rName+"3")
	return resource
}
