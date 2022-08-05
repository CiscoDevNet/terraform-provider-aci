package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciConfigurationExportPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_configuration_export_policy.test"
	dataSourceName := "data.aci_configuration_export_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:	  func(){ testAccPreCheck(t) },
		ProviderFactories:    testAccProviders,
		CheckDestroy: testAccCheckAciConfigurationExportPolicyDestroy,
		Steps: []resource.TestStep{
			
			{
				Config:      CreateConfigurationExportPolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccConfigurationExportPolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "format", resourceName, "format"),
					resource.TestCheckResourceAttrPair(dataSourceName, "include_secure_fields", resourceName, "include_secure_fields"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_snapshot_count", resourceName, "max_snapshot_count"),
					resource.TestCheckResourceAttrPair(dataSourceName, "snapshot", resourceName, "snapshot"),
					resource.TestCheckResourceAttrPair(dataSourceName, "target_dn", resourceName, "target_dn"),
					
				),
			},
			{
				Config:      CreateAccConfigurationExportPolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			
			{
				Config:      CreateAccConfigurationExportPolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccConfigurationExportPolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}


func CreateAccConfigurationExportPolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing configuration_export_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_export_policy" "test" {
	
		name  = "%s"
	}

	data "aci_configuration_export_policy" "test" {
	
		name  = aci_configuration_export_policy.test.name
		depends_on = [ aci_configuration_export_policy.test ]
	}
	`, rName)
	return resource
}

func CreateConfigurationExportPolicyDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing configuration_export_policy Data Source without ",attrName)
	rBlock := `
	
	resource "aci_configuration_export_policy" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_configuration_export_policy" "test" {
	
	#	name  = aci_configuration_export_policy.test.name
		depends_on = [ aci_configuration_export_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock,rName)
}


func CreateAccConfigurationExportPolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing configuration_export_policy Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_export_policy" "test" {
	
		name  = "%s"
	}

	data "aci_configuration_export_policy" "test" {
	
		name  = "${aci_configuration_export_policy.test.name}_invalid"
		depends_on = [ aci_configuration_export_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccConfigurationExportPolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing configuration_export_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_export_policy" "test" {
	
		name  = "%s"
	}

	data "aci_configuration_export_policy" "test" {
	
		name  = aci_configuration_export_policy.test.name
		%s = "%s"
		depends_on = [ aci_configuration_export_policy.test ]
	}
	`, rName,key,value)
	return resource
}

func CreateAccConfigurationExportPolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing configuration_export_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_export_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_configuration_export_policy" "test" {
	
		name  = aci_configuration_export_policy.test.name
		depends_on = [ aci_configuration_export_policy.test ]
	}
	`, rName,key,value)
	return resource
}