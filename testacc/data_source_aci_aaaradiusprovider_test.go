package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciRadiusProviderDataSource_Basic(t *testing.T) {
	resourceName := "aci_radius_provider.test"
	dataSourceName := "data.aci_radius_provider.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	providerType := "radius"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRadiusProviderDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateRadiusProviderDSWithoutRequired(rName, providerType, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateRadiusProviderDSWithoutRequired(rName, providerType, "type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRadiusProviderConfigDataSource(rName, providerType),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
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
				Config: CreateAccRadiusProviderConfigDataSource(rName, "duo"),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_port", resourceName, "auth_port"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_protocol", resourceName, "auth_protocol"),
					resource.TestCheckResourceAttrPair(dataSourceName, "key", resourceName, "key"),
					resource.TestCheckResourceAttrPair(dataSourceName, "monitor_server", resourceName, "monitor_server"),
					resource.TestCheckResourceAttrPair(dataSourceName, "monitoring_password", resourceName, "monitoring_password"),
					resource.TestCheckResourceAttrPair(dataSourceName, "monitoring_user", resourceName, "monitoring_user"),
					resource.TestCheckResourceAttrPair(dataSourceName, "retries", resourceName, "retries"),
					resource.TestCheckResourceAttrPair(dataSourceName, "timeout", resourceName, "timeout"),
				),
			},
			{
				Config:      CreateAccRadiusProviderDataSourceUpdate(rName, providerType, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccRadiusProviderDSWithInvalidName(rName, providerType),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config:      CreateAccRadiusProviderConfigDataSourceWithInvalidType(rName, acctest.RandString(5)),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config: CreateAccRadiusProviderDataSourceUpdatedResource(rName, providerType, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccRadiusProviderConfigDataSourceWithInvalidType(rName, providerType string) string {
	fmt.Println("=== STEP  testing radius_provider Data Source with invalid type")
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider" "test" {
	
		name  = "%s"
		type  = "duo"
		timeout = 60
	}

	data "aci_radius_provider" "test" {
	
		name  = aci_radius_provider.test.name
		type  = "%s"
		depends_on = [ aci_radius_provider.test ]
	}
	`, rName, providerType)
	return resource
}

func CreateAccRadiusProviderConfigDataSource(rName, providerType string) string {
	fmt.Println("=== STEP  testing radius_provider Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider" "test" {
	
		name  = "%s"
		type  = "%s"
		timeout = 60
	}

	data "aci_radius_provider" "test" {
	
		name  = aci_radius_provider.test.name
		type  = aci_radius_provider.test.type
		depends_on = [ aci_radius_provider.test ]
	}
	`, rName, providerType)
	return resource
}

func CreateRadiusProviderDSWithoutRequired(rName, providerType, attrName string) string {
	fmt.Println("=== STEP  Basic: testing radius_provider Data Source without ", attrName, providerType)
	rBlock := `
	
	resource "aci_radius_provider" "test" {
	
		name  = "%s"
		type  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_radius_provider" "test" {
	
	#	name  = aci_radius_provider.test.name
		type  = aci_radius_provider.test.type
		depends_on = [ aci_radius_provider.test ]
	}
		`
	case "type":
		rBlock += `
	data "aci_radius_provider" "test" {
	
		name  = aci_radius_provider.test.name
	#	type  = aci_radius_provider.test.type
		depends_on = [ aci_radius_provider.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, providerType)
}

func CreateAccRadiusProviderDSWithInvalidName(rName, providerType string) string {
	fmt.Println("=== STEP  testing radius_provider Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider" "test" {
	
		name  = "%s"
		type  = "%s"
	}

	data "aci_radius_provider" "test" {
	
		name  = "${aci_radius_provider.test.name}_invalid"
		type  = aci_radius_provider.test.type
		depends_on = [ aci_radius_provider.test ]
	}
	`, rName, providerType)
	return resource
}

func CreateAccRadiusProviderDataSourceUpdate(rName, providerType, key, value string) string {
	fmt.Println("=== STEP  testing radius_provider Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider" "test" {
	
		name  = "%s"
		type  = "%s"
	}

	data "aci_radius_provider" "test" {
	
		name  = aci_radius_provider.test.name
		type  = aci_radius_provider.test.type
		%s = "%s"
		depends_on = [ aci_radius_provider.test ]
	}
	`, rName, providerType, key, value)
	return resource
}

func CreateAccRadiusProviderDataSourceUpdatedResource(rName, providerType, key, value string) string {
	fmt.Println("=== STEP  testing radius_provider Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider" "test" {
	
		name  = "%s"
		type  = "%s"
		%s = "%s"
	}

	data "aci_radius_provider" "test" {
	
		name  = aci_radius_provider.test.name
		type  = aci_radius_provider.test.type
		depends_on = [ aci_radius_provider.test ]
	}
	`, rName, providerType, key, value)
	return resource
}
