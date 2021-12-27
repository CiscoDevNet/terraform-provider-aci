package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLogicalNodeProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_logical_node_profile.test"
	dataSourceName := "data.aci_logical_node_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalNodeProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccLogicalNodeProfileDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLogicalNodeProfileDSWithoutParentDn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLogicalNodeProfileDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag", resourceName, "tag"),
					resource.TestCheckResourceAttrPair(dataSourceName, "target_dscp", resourceName, "target_dscp"),
					resource.TestCheckResourceAttrPair(dataSourceName, "l3_outside_dn", resourceName, "l3_outside_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "config_issues", resourceName, "config_issues"),
				),
			},
			{
				Config:      CreateAccLogicalNodeProfileDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccLogicalNodeProfileDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`Object may not exists`),
			},
			{
				Config: CreateAccLogicalNodeProfileDataSourceUpdate(rName, "description", "test_annotation_1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccLogicalNodeProfileDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalNodeProfile data source creation without giving Name")
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
		name = "%s"
	  }
	  data "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
	  }
	`, rName, rName, rName)
	return resource
}

func CreateAccLogicalNodeProfileDataSourceUpdate(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing any data source update for attribute: %s = %s \n", attribute, value)
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
		%s = "%s"
      }
	data "aci_logical_node_profile" "test" {
		l3_outside_dn = aci_logical_node_profile.test.l3_outside_dn
        name          = aci_logical_node_profile.test.name
	}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccLogicalNodeProfileDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: testing any data source reading with required arguments only")
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
	data "aci_logical_node_profile" "test" {
		l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
		depends_on = [aci_logical_node_profile.test]
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccLogicalNodeProfileDataSourceUpdateRandomAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing any data source update for attribute: %s = %s \n", attribute, value)
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
	data "aci_logical_node_profile" "test" {
		l3_outside_dn = aci_logical_node_profile.test.l3_outside_dn
        name          = aci_logical_node_profile.test.name
		%s = "%s"
		}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccLogicalNodeProfileDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalNodeProfile data source reading with invalid name")
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
	data "aci_logical_node_profile" "test" {
		l3_outside_dn = aci_logical_node_profile.test.l3_outside_dn
        name          =  "${aci_logical_node_profile.test.name}xyz"
	} 
	`, rName, rName, rName)
	return resource
}

func CreateAccLogicalNodeProfileDSWithoutParentDn(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalNodeProfile creation without giving parent dn")
	resource := fmt.Sprintf(`
	
	  data "aci_logical_node_profile" "test" {
        name = "%s"
	  }
	`, rName)
	return resource
}
