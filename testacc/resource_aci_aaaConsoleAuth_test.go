package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/terraform-providers/terraform-provider-aci/aci"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciConsoleAuthentication_Basic(t *testing.T) {
	var console_authentication_default models.ConsoleAuthenticationMethod
	var console_authentication_updated models.ConsoleAuthenticationMethod
	resourceName := "aci_console_authentication.test"
	aaaConsoleAuth, err := aci.GetRemoteConsoleAuthenticationMethod(sharedAciClient(), "uni/userext/authrealm/consoleauth")
	if err != nil {
		t.Errorf("reading initial config of aaaConsoleAuth")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAciConsoleAuthenticationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccAciConsoleAuthenticationExists(resourceName, &console_authentication_default),
					resource.TestCheckResourceAttrSet(resourceName, "realm"),
					resource.TestCheckResourceAttrSet(resourceName, "realm_sub_type"),
				),
			},
			{
				Config: CreateAccAciConsoleAuthenticationConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_console_authentication"),
					resource.TestCheckResourceAttr(resourceName, "provider_group", "60"),
					resource.TestCheckResourceAttr(resourceName, "realm", "ldap"),
					resource.TestCheckResourceAttr(resourceName, "realm_sub_type", "default"),
					testAccCheckAccAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: restoreConsoleAuthentication(aaaConsoleAuth),
			},
		},
	})
}

func TestAccAciConsoleAuthentication_Update(t *testing.T) {
	var console_authentication_default models.ConsoleAuthenticationMethod
	var console_authentication_updated models.ConsoleAuthenticationMethod
	resourceName := "aci_console_authentication.test"
	aaaConsoleAuth, err := aci.GetRemoteConsoleAuthenticationMethod(sharedAciClient(), "uni/userext/authrealm/consoleauth")
	if err != nil {
		t.Errorf("reading initial config of aaaConsoleAuth")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAciConsoleAuthenticationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccAciConsoleAuthenticationExists(resourceName, &console_authentication_default),
				),
			},
			{
				Config: CreateAccAciConsoleAuthenticationUpdatedAttr("realm", "local"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "local"),
					testAccCheckAccAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				Config: CreateAccAciConsoleAuthenticationUpdatedAttr("realm", "radius"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "radius"),
					testAccCheckAccAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				Config: CreateAccAciConsoleAuthenticationUpdatedAttr("realm", "rsa"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "rsa"),
					testAccCheckAccAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				Config: CreateAccAciConsoleAuthenticationUpdatedAttr("realm", "saml"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "saml"),
					testAccCheckAccAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				Config: CreateAccAciConsoleAuthenticationUpdatedAttr("realm", "tacacs"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "tacacs"),
					testAccCheckAccAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				Config: CreateAccAciConsoleAuthenticationUpdatedAttr("realm_sub_type", "duo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm_sub_type", "duo"),
					testAccCheckAccAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				Config: CreateAccAciConsoleAuthenticationUpdatedAttr("provider_group", "0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "provider_group", "0"),
					testAccCheckAccAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				Config: CreateAccAciConsoleAuthenticationUpdatedAttr("provider_group", "63"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "provider_group", "63"),
					testAccCheckAccAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				Config: restoreConsoleAuthentication(aaaConsoleAuth),
			},
		},
	})
}

func TestAccAciConsoleAuthentication_Negative(t *testing.T) {
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	aaaConsoleAuth, err := aci.GetRemoteConsoleAuthenticationMethod(sharedAciClient(), "uni/userext/authrealm/consoleauth")
	if err != nil {
		t.Errorf("reading initial config of aaaConsoleAuth")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: CreateAccAciConsoleAuthenticationConfig(),
			},
			{
				Config:      CreateAccAciConsoleAuthenticationUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAciConsoleAuthenticationUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAciConsoleAuthenticationUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAciConsoleAuthenticationUpdatedAttr("realm", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccAciConsoleAuthenticationUpdatedAttr("realm_sub_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccAciConsoleAuthenticationUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: restoreConsoleAuthentication(aaaConsoleAuth),
			},
		},
	})
}

func restoreConsoleAuthentication(aaaConsoleAuth *models.ConsoleAuthenticationMethod) string {
	resource := fmt.Sprintf(`
	resource "aci_console_authentication" "test" {
		annotation = "%s"
		description = "%s"
		name_alias = "%s"
		provider_group = "%s"
		realm          = "%s"
		realm_sub_type = "%s"
	}
	`, aaaConsoleAuth.Annotation, aaaConsoleAuth.Description, aaaConsoleAuth.NameAlias, aaaConsoleAuth.ProviderGroup, aaaConsoleAuth.Realm, aaaConsoleAuth.RealmSubType)
	return resource
}

func testAccCheckAccAciConsoleAuthenticationExists(name string, console_authentication *models.ConsoleAuthenticationMethod) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Console Autentication %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Console Autentication Data dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		console_authenticationFound := models.ConsoleAuthenticationMethodFromContainer(cont)
		if console_authenticationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Console Autentication %s not found", rs.Primary.ID)
		}
		*console_authentication = *console_authenticationFound
		return nil
	}
}

func testAccCheckAccAciConsoleAuthenticationIdEqual(m1, m2 *models.ConsoleAuthenticationMethod) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("Console Autentication DNs are not equal")
		}
		return nil
	}
}

func CreateAccAciConsoleAuthenticationConfig() string {
	fmt.Println("=== STEP  Testing console_authentication creation")
	resource := fmt.Sprintf(`
	
	resource "aci_console_authentication" "test" {}
	`)
	return resource
}

func CreateAccAciConsoleAuthenticationConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Testing console_authentication creation with optional parameters")
	resource := fmt.Sprint(`
	resource "aci_console_authentication" "test" {
		annotation     = "orchestrator:terraform_testacc"
		provider_group = "60"
		realm          = "ldap"
		realm_sub_type = "default"
		name_alias     = "test_console_authentication"
		description    = "created while acceptance testing"
	}
	`)

	return resource
}

func CreateAccAciConsoleAuthenticationUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  Testing console_authentication attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_console_authentication" "test" {
	
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
