package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciRSAProviderDataSource_Basic(t *testing.T) {
	resourceName := "aci_rsa_provider.test"
	dataSourceName := "data.aci_rsa_provider.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRSAProviderDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateRSAProviderDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRSAProviderConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_port", resourceName, "auth_port"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_protocol", resourceName, "auth_protocol"),
					resource.TestCheckResourceAttrPair(dataSourceName, "monitor_server", resourceName, "monitor_server"),
					resource.TestCheckResourceAttrPair(dataSourceName, "monitoring_user", resourceName, "monitoring_user"),
					resource.TestCheckResourceAttrPair(dataSourceName, "retries", resourceName, "retries"),
					resource.TestCheckResourceAttrPair(dataSourceName, "timeout", resourceName, "timeout"),
				),
			},
			{
				Config:      CreateAccRSAProviderDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccRSAProviderDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccRSAProviderDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccRSAProviderConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing rsa_provider Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_rsa_provider" "test" {
	
		name  = "%s"
	}

	data "aci_rsa_provider" "test" {
	
		name  = aci_rsa_provider.test.name
		depends_on = [ aci_rsa_provider.test ]
	}
	`, rName)
	return resource
}

func CreateRSAProviderDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing rsa_provider Data Source without ", attrName)
	rBlock := `
	
	resource "aci_rsa_provider" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_rsa_provider" "test" {
	
	#	name  = aci_rsa_provider.test.name
		depends_on = [ aci_rsa_provider.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccRSAProviderDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing rsa_provider Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_rsa_provider" "test" {
	
		name  = "%s"
	}

	data "aci_rsa_provider" "test" {
	
		name  = "${aci_rsa_provider.test.name}_invalid"
		depends_on = [ aci_rsa_provider.test ]
	}
	`, rName)
	return resource
}

func CreateAccRSAProviderDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing rsa_provider Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_rsa_provider" "test" {
	
		name  = "%s"
	}

	data "aci_rsa_provider" "test" {
	
		name  = aci_rsa_provider.test.name
		%s = "%s"
		depends_on = [ aci_rsa_provider.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccRSAProviderDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing rsa_provider Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_rsa_provider" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_rsa_provider" "test" {
	
		name  = aci_rsa_provider.test.name
		depends_on = [ aci_rsa_provider.test ]
	}
	`, rName, key, value)
	return resource
}
