package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciAAAAuthenticationDataSource_Basic(t *testing.T) {
	resourceName := "aci_authentication_properties.test"
	dataSourceName := "data.aci_authentication_properties.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	aaaAuthRealm, err := aci.GetRemoteAAAAuthentication(sharedAciClient(), "uni/userext/authrealm")
	if err != nil {
		t.Errorf("reading initial config of authenticationProperties")
	}
	aaaPingEp, err := aci.GetRemoteDefaultRadiusAuthenticationSettings(sharedAciClient(), "uni/userext/pingext")
	if err != nil {
		t.Errorf("reading initial config of authenticationProperties")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAAAAuthenticationConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "def_role_policy", resourceName, "def_role_policy"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ping_check", resourceName, "ping_check"),
					resource.TestCheckResourceAttrPair(dataSourceName, "retries", resourceName, "retries"),
					resource.TestCheckResourceAttrPair(dataSourceName, "timeout", resourceName, "timeout"),
				),
			},
			{
				Config:      CreateAccAAAAuthenticationDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config: CreateAccAAAAuthenticationDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},

			{
				Config: CreateAccAAAAuthenticationInitialConfig(aaaAuthRealm, aaaPingEp),
			},
		},
	})
}

func CreateAccAAAAuthenticationConfigDataSource() string {
	fmt.Println("=== STEP  testing authentication_properties Data Source with required arguments only")
	resource := fmt.Sprintln(`
	
	resource "aci_authentication_properties" "test" {
	
	}

	data "aci_authentication_properties" "test" {
	
		depends_on = [ aci_authentication_properties.test ]
	}
	`)
	return resource
}

func CreateAccAAAAuthenticationDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing authentication_properties Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_authentication_properties" "test" {
	
	}

	data "aci_authentication_properties" "test" {
	
		%s = "%s"
		depends_on = [ aci_authentication_properties.test ]
	}
	`, key, value)
	return resource
}

func CreateAccAAAAuthenticationDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing authentication_properties Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_authentication_properties" "test" {
	
		%s = "%s"
	}

	data "aci_authentication_properties" "test" {
	
		depends_on = [ aci_authentication_properties.test ]
	}
	`, key, value)
	return resource
}

func CreateAccAAAAuthenticationInitialConfig(aaaAuthRealm *models.AAAAuthentication, aaaPingEp *models.DefaultRadiusAuthenticationSettings) string {
	resource := fmt.Sprintf(`
	
	resource "aci_authentication_properties" "test" {
		annotation = "%s"
		name_alias = "%s"
		description = "%s"
		def_role_policy = "%s"
		ping_check = "%s"
		retries = "%s"
		timeout = "%s"
	}
	`, aaaAuthRealm.Annotation, aaaAuthRealm.NameAlias, aaaAuthRealm.Description, aaaAuthRealm.DefRolePolicy, aaaPingEp.PingCheck, aaaPingEp.Retries, aaaPingEp.Timeout)
	return resource
}
