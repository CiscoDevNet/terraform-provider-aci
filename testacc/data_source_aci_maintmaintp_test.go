package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciMaintenancePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_maintenance_policy.test"
	dataSourceName := "data.aci_maintenance_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMaintenancePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateMaintenancePolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMaintenancePolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "graceful", resourceName, "graceful"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ignore_compat", resourceName, "ignore_compat"),
					resource.TestCheckResourceAttrPair(dataSourceName, "internal_label", resourceName, "internal_label"),
					resource.TestCheckResourceAttrPair(dataSourceName, "notif_cond", resourceName, "notif_cond"),
					resource.TestCheckResourceAttrPair(dataSourceName, "run_mode", resourceName, "run_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "version", resourceName, "version"),
					resource.TestCheckResourceAttrPair(dataSourceName, "version_check_override", resourceName, "version_check_override"),
				),
			},
			{
				Config:      CreateAccMaintenancePolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccMaintenancePolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccMaintenancePolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccMaintenancePolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing maintenance_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_maintenance_policy" "test" {
	
		name  = "%s"
	}

	data "aci_maintenance_policy" "test" {
	
		name  = aci_maintenance_policy.test.name
		depends_on = [ aci_maintenance_policy.test ]
	}
	`, rName)
	return resource
}

func CreateMaintenancePolicyDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing maintenance_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_maintenance_policy" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_maintenance_policy" "test" {
	
	#	name  = aci_maintenance_policy.test.name
		depends_on = [ aci_maintenance_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccMaintenancePolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing maintenance_policy Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_maintenance_policy" "test" {
	
		name  = "%s"
	}

	data "aci_maintenance_policy" "test" {
	
		name  = "${aci_maintenance_policy.test.name}_invalid"
		depends_on = [ aci_maintenance_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccMaintenancePolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing maintenance_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_maintenance_policy" "test" {
	
		name  = "%s"
	}

	data "aci_maintenance_policy" "test" {
	
		name  = aci_maintenance_policy.test.name
		%s = "%s"
		depends_on = [ aci_maintenance_policy.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccMaintenancePolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing maintenance_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_maintenance_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_maintenance_policy" "test" {
	
		name  = aci_maintenance_policy.test.name
		depends_on = [ aci_maintenance_policy.test ]
	}
	`, rName, key, value)
	return resource
}
