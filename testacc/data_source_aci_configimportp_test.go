package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciConfigurationImportPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_configuration_import_policy.test"
	dataSourceName := "data.aci_configuration_import_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConfigurationImportPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateConfigurationImportPolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccConfigurationImportPolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fail_on_decrypt_errors", resourceName, "fail_on_decrypt_errors"),
					resource.TestCheckResourceAttrPair(dataSourceName, "file_name", resourceName, "file_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "import_mode", resourceName, "import_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "import_type", resourceName, "import_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "snapshot", resourceName, "snapshot"),
				),
			},
			{
				Config:      CreateAccConfigurationImportPolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccConfigurationImportPolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccConfigurationImportPolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccConfigurationImportPolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing configuration_import_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_import_policy" "test" {
	
		name  = "%s"
		file_name  = "file.tar.gz"
	}

	data "aci_configuration_import_policy" "test" {
	
		name  = aci_configuration_import_policy.test.name
		depends_on = [ aci_configuration_import_policy.test ]
	}
	`, rName)
	return resource
}

func CreateConfigurationImportPolicyDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing configuration_import_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_configuration_import_policy" "test" {
	
		name  = "%s"
		file_name  = "file.tar.gz"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_configuration_import_policy" "test" {
	
	#	name  = aci_configuration_import_policy.test.name
		depends_on = [ aci_configuration_import_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccConfigurationImportPolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing configuration_import_policy Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_import_policy" "test" {
	
		name  = "%s"
		file_name  = "file.tar.gz"
	}

	data "aci_configuration_import_policy" "test" {
	
		name  = "${aci_configuration_import_policy.test.name}_invalid"
		depends_on = [ aci_configuration_import_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccConfigurationImportPolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing configuration_import_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_import_policy" "test" {
	
		name  = "%s"
		file_name  = "file.tar.gz"
	}

	data "aci_configuration_import_policy" "test" {
	
		name  = aci_configuration_import_policy.test.name
		%s = "%s"
		depends_on = [ aci_configuration_import_policy.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccConfigurationImportPolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing configuration_import_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_import_policy" "test" {
	
		name  = "%s"
		file_name  = "file.tar.gz"
		%s = "%s"
	}

	data "aci_configuration_import_policy" "test" {
	
		name  = aci_configuration_import_policy.test.name
		depends_on = [ aci_configuration_import_policy.test ]
	}
	`, rName, key, value)
	return resource
}
