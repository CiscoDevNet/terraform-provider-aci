package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciConsoleAuthenticationDataSource_Basic(t *testing.T) {
	resourceName := "aci_console_authentication.test"
	dataSourceName := "data.aci_console_authentication.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	aaaConsoleAuth, err := aci.GetRemoteConsoleAuthenticationMethod(sharedAciClient(), "uni/userext/authrealm/consoleauth")
	if err != nil {
		t.Errorf("reading initial config of aaaConsoleAuth")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAciConsoleAuthenticationConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "provider_group", resourceName, "provider_group"),
					resource.TestCheckResourceAttrPair(dataSourceName, "realm", resourceName, "realm"),
					resource.TestCheckResourceAttrPair(dataSourceName, "realm_sub_type", resourceName, "realm_sub_type"),
				),
			},
			{
				Config:      CreateAccAciConsoleAuthenticationDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccAciConsoleAuthenticationDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config: restoreConsoleAuthentication(aaaConsoleAuth),
			},
		},
	})
}

func CreateAccAciConsoleAuthenticationConfigDataSource() string {
	fmt.Println("=== STEP  Testing console_authentication Data Source")
	resource := fmt.Sprintf(`
	resource "aci_console_authentication" "test" {}

	data "aci_console_authentication" "test" {
		depends_on = [ aci_console_authentication.test ]
	}
	`)
	return resource
}

func CreateAccAciConsoleAuthenticationDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  Testing console_authentication Data Source with random attribute")
	resource := fmt.Sprintf(`

	resource "aci_console_authentication" "test" {}

	data "aci_console_authentication" "test" {
		%s = "%s"
		depends_on = [ aci_console_authentication.test ]
	}
	`, key, value)
	return resource
}

func CreateAccAciConsoleAuthenticationDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  Testing console_authentication Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_console_authentication" "test" {
		%s = "%s"
	}

	data "aci_console_authentication" "test" {
		depends_on = [ aci_console_authentication.test ]
	}
	`, key, value)
	return resource
}
