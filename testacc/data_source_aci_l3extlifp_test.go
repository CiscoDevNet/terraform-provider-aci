package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLogicalInterfaceProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_logical_interface_profile.test"
	dataSourceName := "data.aci_logical_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLogicalInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccLogicalInterfaceProfileDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileDSWithoutLogicalNodeProfileDn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAciLogicalInterfaceProfileDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag", resourceName, "tag"),
					resource.TestCheckResourceAttrPair(dataSourceName, "logical_node_profile_dn", resourceName, "logical_node_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
				),
			},
			{
				Config:      CreateAccAciLogicalInterfaceProfileDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileDSWithInvalidNodeProfileDn(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccLogicalInterfaceProfileDataSourceUpdate(rName, "description", "test_annotation_1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccAciLogicalInterfaceProfileDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalInterfaceProfile data source reading with required arguments only")
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
	  data "aci_logical_interface_profile" "test" {
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = aci_logical_node_profile.test.name
		depends_on = [aci_logical_interface_profile.test]
	  }
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccAciLogicalInterfaceProfileDataSourceUpdateRandomAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing LogicalInterfaceProfile data source update for attribute: %s = %s \n", attribute, value)
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
	  data "aci_logical_interface_profile" "test" {
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = aci_logical_node_profile.test.name
		depends_on = [aci_logical_interface_profile.test]
		%s = "%s"
	  }
	`, rName, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccLogicalInterfaceProfileDataSourceUpdate(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing LogicalInterfaceProfile data source update for attribute: %s = %s \n", attribute, value)
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
		%s = "%s"
  	}
	  data "aci_logical_interface_profile" "test" {
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = aci_logical_node_profile.test.name
		depends_on = [aci_logical_interface_profile.test]
	  }
	`, rName, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccLogicalInterfaceProfileDSWithInvalidNodeProfileDn(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalInterfaceProfile data source reading with invalid Logical Node Profile dn")
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
	  data "aci_logical_interface_profile" "test" {
		logical_node_profile_dn = "${aci_logical_node_profile.test.id}xyz"
		name = aci_logical_node_profile.test.name
		depends_on = [aci_logical_interface_profile.test]
	  }
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccLogicalInterfaceProfileDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalInterfaceProfile data source reading with invalid name")
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
	  data "aci_logical_interface_profile" "test" {
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = "{aci_logical_node_profile.test.name}xyz"
		depends_on = [aci_logical_interface_profile.test]
	  }
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccLogicalInterfaceProfileDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalInterfaceProfile data source reading without giving name")
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
	  data "aci_logical_interface_profile" "test" {
		logical_node_profile_dn = aci_logical_node_profile.test.id
		depends_on = [aci_logical_interface_profile.test]
	  }
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccLogicalInterfaceProfileDSWithoutLogicalNodeProfileDn(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalInterfaceProfile data source reading without giving Logical Node Profile dn")
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
	  data "aci_logical_interface_profile" "test" {
		name = aci_logical_node_profile.test.name
		depends_on = [aci_logical_interface_profile.test]
	  }
	`, rName, rName, rName, rName)
	return resource
}
