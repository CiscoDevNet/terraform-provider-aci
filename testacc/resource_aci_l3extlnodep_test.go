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

func TestAccAciLogicalNodeProfile_Basic(t *testing.T) {
	var LogicalNodeProfile_default models.LogicalNodeProfile
	var LogicalNodeProfile_updated models.LogicalNodeProfile
	resourceName := "aci_logical_node_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	prOther := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLogicalNodeProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccLogicalNodeProfileWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLogicalNodeProfileWithoutParentDn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLogicalNodeProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttrSet(resourceName, "tag"),
					resource.TestCheckResourceAttrSet(resourceName, "config_issues"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "l3_outside_dn", fmt.Sprintf("uni/tn-%s/out-%s", rName, rName)),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag_node"),
					resource.TestCheckResourceAttr(resourceName, "description", "sample logical node profile"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "alias_node"),
					resource.TestCheckResourceAttr(resourceName, "tag", "black"),
					resource.TestCheckResourceAttrSet(resourceName, "config_issues"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS0"),
					resource.TestCheckResourceAttr(resourceName, "l3_outside_dn", fmt.Sprintf("uni/tn-%s/out-%s", rName, rName)),

					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccLogicalNodeProfileConfigWithAnotherName(rName, rOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					testAccCheckAciLogicalNodeProfileIdNotEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileConfig(rName),
			},
			{
				Config: CreateAccLogicalNodeProfileConfigWithAnotherParentDn(prOther, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					testAccCheckAciLogicalNodeProfileIdNotEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config:      CreateAccLogicalNodeProfileConfigUpdateWithoutParentdn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLogicalNodeProfileConfigUpdateWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLogicalNodeProfileConfigUpdateWithInvalidName(rName, longrName),
				ExpectError: regexp.MustCompile("failed validation for value"),
			},
			{
				Config: CreateAccLogicalNodeProfileConfig(rName),
			},
		},
	})
}

func TestAccAciLogicalNodeProfile_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLogicalNodeProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLogicalNodeProfileConfig(rName),
			},
			{
				Config:      CreateAccLogicalNodeProfileUpdatedAttr(rName, "tag", longDescAnnotation),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLogicalNodeProfileUpdatedAttr(rName, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalNodeProfileUpdatedAttr(rName, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalNodeProfileUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalNodeProfileUpdatedAttr(rName, "target_dscp", randomValue),
				ExpectError: regexp.MustCompile(`expected target_dscp to be one of (.)+ got (.)+`),
			},
			{
				Config:      CreateAccLogicalNodeProfileUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLogicalNodeProfileConfig(rName),
			},
		},
	})
}
func TestAccAciLogicalNodeProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLogicalNodeProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLogicalNodeProfileConfigMultiple(rName),
			},
		},
	})
}

func TestAccAciLogicalNodeProfile_Update(t *testing.T) {
	var LogicalNodeProfile_default models.LogicalNodeProfile
	var LogicalNodeProfile_updated models.LogicalNodeProfile
	resourceName := "aci_logical_node_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLogicalNodeProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLogicalNodeProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_default),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "target_dscp", "CS1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS1"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "target_dscp", "AF11"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF11"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "target_dscp", "AF12"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF12"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "target_dscp", "AF13"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF13"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "target_dscp", "CS2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS2"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "target_dscp", "AF21"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF21"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "target_dscp", "AF22"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF22"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "target_dscp", "AF23"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF23"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "target_dscp", "CS3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS3"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "tag", "pink"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "pink"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "tag", "chartreuse"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "chartreuse"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "tag", "plum"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "plum"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "tag", "dark-orchid"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-orchid"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalNodeProfileUpdatedAttr(rName, "tag", "lime-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalNodeProfileExists(resourceName, &LogicalNodeProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "lime-green"),
					testAccCheckAciLogicalNodeProfileIdEqual(&LogicalNodeProfile_default, &LogicalNodeProfile_updated),
				),
			},
		},
	})
}

func testAccCheckLogicalNodeProfileExists(name string, LogicalNodeProfile *models.LogicalNodeProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LogicalNodeProfile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LogicalNodeProfile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		LogicalNodeProfileFound := models.LogicalNodeProfileFromContainer(cont)
		if LogicalNodeProfileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LogicalNodeProfile %s not found", rs.Primary.ID)
		}
		*LogicalNodeProfile = *LogicalNodeProfileFound
		return nil
	}
}
func CreateAccLogicalNodeProfileConfigMultiple(rName string) string {
	fmt.Println("=== STEP  creating multiple LogicalNodeProfiles")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        name           = "%s"
    }
   resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
	resource "aci_logical_node_profile" "test1" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
	resource "aci_logical_node_profile" "test2" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
	`, rName, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func testAccCheckAciLogicalNodeProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing LogicalNodeProfile destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_logical_node_profile" {
			cont, err := client.Get(rs.Primary.ID)
			aci := models.LogicalNodeProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LogicalNodeProfile %s Still exists", aci.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLogicalNodeProfileIdEqual(LogicalNodeProfile1, LogicalNodeProfile2 *models.LogicalNodeProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if LogicalNodeProfile1.DistinguishedName != LogicalNodeProfile2.DistinguishedName {
			return fmt.Errorf("LogicalNodeProfile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLogicalNodeProfileIdNotEqual(LogicalNodeProfile1, LogicalNodeProfile2 *models.LogicalNodeProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if LogicalNodeProfile1.DistinguishedName == LogicalNodeProfile2.DistinguishedName {
			return fmt.Errorf("LogicalNodeProfile DNs are equal")
		}
		return nil
	}
}

func CreateAccLogicalNodeProfileWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalNodeProfile creation without giving Name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
			tenant_dn      = aci_tenant.test.id
			name 		= "%s"
	  }
	  resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
	  }
	`, rName, rName)
	return resource
}

func CreateAccLogicalNodeProfileWithoutParentDn(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalNodeProfile creation without giving parent dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		name = "%s"
	  }
	  resource "aci_logical_node_profile" "test" {
        name = "%s"
	  }
	`, rName, rName, rName)
	return resource
}

func CreateAccLogicalNodeProfileConfigUpdateWithoutParentdn(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalNodeProfile update without giving Parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		name    = "%s"
	  }
	  resource "aci_logical_interface_profile" "test" {
		name    = "%s"
		annotation    = "tag"
        name_alias    = "alias_node"
        tag           = "black"
        target_dscp   = "CS0"
	  }
	`, rName, rName, rName)
	return resource
}

func CreateAccLogicalNodeProfileConfigUpdateWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalNodeProfile update without giving Name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		tenant_dn  = aci_tenant.test.id
		name ="%s"
	  }
	  resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
		annotation    = "tag"
        name_alias    = "alias_node"
        tag           = "black"
        target_dscp   = "CS0"
	  }
	`, rName, rName)
	return resource
}

func CreateAccLogicalNodeProfileConfigUpdateWithInvalidName(parentName, rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalNodeProfile update with invalid Name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		tenant_dn  = aci_tenant.test.id
		name= "%s"
	  }
	  resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
		name = "%s"
	  }
	`, parentName, parentName, rName)
	return resource
}

func CreateAccLogicalNodeProfileConfigWithAnotherName(parentName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing LogicalNodeProfile creation with different LogicalNodeProfile name %s \n", rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	  }
	  resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
		name = "%s"
	  }
	`, parentName, parentName, rName)
	return resource
}

func CreateAccLogicalNodeProfileConfigWithAnotherParentDn(parentName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing LogicalNodeProfile creation with different parent %s \n", parentName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	  }
	  resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
	  }
	`, parentName, parentName, rName)
	return resource
}

func CreateAccLogicalNodeProfileConfig(rName string) string {
	fmt.Println("=== STEP  testing LogicalNodeProfile creation with required attributes")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		name           = "%s"
	  }
	  resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
	    name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccLogicalNodeProfileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalNodeProfile creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		tenant_dn      = aci_tenant.test.id
	    name           = "%s"
	  }
	  resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        description   = "sample logical node profile"
        name          = "%s"
        annotation    = "tag_node"
        name_alias    = "alias_node"
        tag           = "black"
        target_dscp   = "CS0"
      }
	`, rName, rName, rName)
	return resource
}

func CreateAccLogicalNodeProfileUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		tenant_dn      = aci_tenant.test.id
		name           = "%s"
	
	}
	resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
		name = "%s"
		%s = "%s"
	}
	`, rName, rName, rName, attribute, value)
	return resource
}
