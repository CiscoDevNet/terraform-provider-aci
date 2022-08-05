package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSAMLProviderDataSource_Basic(t *testing.T) {
	resourceName := "aci_saml_provider.test"
	dataSourceName := "data.aci_saml_provider.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSAMLProviderDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateSAMLProviderDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSAMLProviderConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "entity_id", resourceName, "entity_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "gui_banner_message", resourceName, "gui_banner_message"),
					resource.TestCheckResourceAttrPair(dataSourceName, "https_proxy", resourceName, "https_proxy"),
					resource.TestCheckResourceAttrPair(dataSourceName, "id_p", resourceName, "id_p"),
					resource.TestCheckResourceAttrPair(dataSourceName, "metadata_url", resourceName, "metadata_url"),
					resource.TestCheckResourceAttrPair(dataSourceName, "monitor_server", resourceName, "monitor_server"),
					resource.TestCheckResourceAttrPair(dataSourceName, "monitoring_user", resourceName, "monitoring_user"),
					resource.TestCheckResourceAttrPair(dataSourceName, "retries", resourceName, "retries"),
					resource.TestCheckResourceAttrPair(dataSourceName, "sig_alg", resourceName, "sig_alg"),
					resource.TestCheckResourceAttrPair(dataSourceName, "timeout", resourceName, "timeout"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tp", resourceName, "tp"),
					resource.TestCheckResourceAttrPair(dataSourceName, "want_assertions_encrypted", resourceName, "want_assertions_encrypted"),
					resource.TestCheckResourceAttrPair(dataSourceName, "want_assertions_signed", resourceName, "want_assertions_signed"),
					resource.TestCheckResourceAttrPair(dataSourceName, "want_requests_signed", resourceName, "want_requests_signed"),
					resource.TestCheckResourceAttrPair(dataSourceName, "want_response_signed", resourceName, "want_response_signed"),
				),
			},
			{
				Config:      CreateAccSAMLProviderDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccSAMLProviderDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccSAMLProviderDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccSAMLProviderConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing saml_provider Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider" "test" {
	
		name  = "%s"
	}

	data "aci_saml_provider" "test" {
	
		name  = aci_saml_provider.test.name
		depends_on = [ aci_saml_provider.test ]
	}
	`, rName)
	return resource
}

func CreateSAMLProviderDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing saml_provider Data Source without ", attrName)
	rBlock := `
	
	resource "aci_saml_provider" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_saml_provider" "test" {
	
	#	name  = aci_saml_provider.test.name
		depends_on = [ aci_saml_provider.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccSAMLProviderDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing saml_provider Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider" "test" {
	
		name  = "%s"
	}

	data "aci_saml_provider" "test" {
	
		name  = "${aci_saml_provider.test.name}_invalid"
		depends_on = [ aci_saml_provider.test ]
	}
	`, rName)
	return resource
}

func CreateAccSAMLProviderDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing saml_provider Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider" "test" {
	
		name  = "%s"
	}

	data "aci_saml_provider" "test" {
	
		name  = aci_saml_provider.test.name
		%s = "%s"
		depends_on = [ aci_saml_provider.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccSAMLProviderDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing saml_provider Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_saml_provider" "test" {
	
		name  = aci_saml_provider.test.name
		depends_on = [ aci_saml_provider.test ]
	}
	`, rName, key, value)
	return resource
}
