package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciDefaultAuthenticationDataSource_Basic(t *testing.T) {
	resourceName := "aci_default_authentication.test"
	dataSourceName := "data.aci_default_authentication.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	defaultAuthentication, err := aci.GetRemoteDefaultAuthenticationMethodforallLogins(sharedAciClient(), "uni/userext/authrealm/defaultauth")
	if err != nil {
		t.Errorf("reading initial config of Default Authentication")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDefaultAuthenticationConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fallback_check", resourceName, "fallback_check"),
					resource.TestCheckResourceAttrPair(dataSourceName, "provider_group", resourceName, "provider_group"),
					resource.TestCheckResourceAttrPair(dataSourceName, "realm", resourceName, "realm"),
					resource.TestCheckResourceAttrPair(dataSourceName, "realm_sub_type", resourceName, "realm_sub_type"),
				),
			},
			{
				Config:      CreateAccDefaultAuthenticationDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccDefaultAuthenticationDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config: CreateAccDefaultAuthenticationInitialConfig(defaultAuthentication),
			},
		},
	})
}

func CreateAccDefaultAuthenticationConfigDataSource() string {
	fmt.Println("=== STEP  testing default_authentication Data Source")
	resource := fmt.Sprint(`
	
	resource "aci_default_authentication" "test" {
	}

	data "aci_default_authentication" "test" {
		depends_on = [ aci_default_authentication.test ]
	}
	`)
	return resource
}

func CreateAccDefaultAuthenticationDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing default_authentication Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {	
	}

	data "aci_default_authentication" "test" {
		%s = "%s"
		depends_on = [ aci_default_authentication.test ]
	}
	`, key, value)
	return resource
}

func CreateAccDefaultAuthenticationDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing default_authentication Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
		%s = "%s"
	}

	data "aci_default_authentication" "test" {
	
		depends_on = [ aci_default_authentication.test ]
	}
	`, key, value)
	return resource
}
