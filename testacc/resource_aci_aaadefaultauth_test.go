package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciDefaultAuthentication_Basic(t *testing.T) {
	var default_authentication_default models.DefaultAuthenticationMethodforallLogins
	var default_authentication_updated models.DefaultAuthenticationMethodforallLogins
	resourceName := "aci_default_authentication.test"
	defaultAuthentication, err := aci.GetRemoteDefaultAuthenticationMethodforallLogins(sharedAciClient(), "uni/userext/authrealm/defaultauth")
	if err != nil {
		t.Errorf("reading initial config of Default Authentication")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDefaultAuthenticationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_default),
					resource.TestCheckResourceAttrSet(resourceName, "fallback_check"),
					resource.TestCheckResourceAttrSet(resourceName, "realm"),
					resource.TestCheckResourceAttrSet(resourceName, "realm_sub_type"),
				),
			},
			{
				Config: CreateAccDefaultAuthenticationConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_default_authentication"),
					resource.TestCheckResourceAttr(resourceName, "fallback_check", "false"),
					testAccCheckAciDefaultAuthenticationIdEqual(&default_authentication_default, &default_authentication_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccDefaultAuthenticationInitialConfig(defaultAuthentication),
			},
		},
	})
}

func TestAccAciDefaultAuthentication_Update(t *testing.T) {
	var default_authentication_default models.DefaultAuthenticationMethodforallLogins
	var default_authentication_updated models.DefaultAuthenticationMethodforallLogins
	resourceName := "aci_default_authentication.test"
	defaultAuthentication, err := aci.GetRemoteDefaultAuthenticationMethodforallLogins(sharedAciClient(), "uni/userext/authrealm/defaultauth")
	if err != nil {
		t.Errorf("reading initial config of Default Authentication")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDefaultAuthenticationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_default),
				),
			},
			{
				Config: CreateAccDefaultAuthenticationUpdatedAttr("realm", "tacacs"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "tacacs"),
					testAccCheckAciDefaultAuthenticationIdEqual(&default_authentication_default, &default_authentication_updated),
				),
			},
			{
				Config: CreateAccDefaultAuthenticationUpdatedAttr("realm", "rsa"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "rsa"),
					testAccCheckAciDefaultAuthenticationIdEqual(&default_authentication_default, &default_authentication_updated),
				),
			},
			{
				Config: CreateAccDefaultAuthenticationUpdatedAttr("realm", "saml"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "saml"),
					testAccCheckAciDefaultAuthenticationIdEqual(&default_authentication_default, &default_authentication_updated),
				),
			},
			{
				Config: CreateAccDefaultAuthenticationUpdatedAttr("realm", "radius"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "radius"),
					testAccCheckAciDefaultAuthenticationIdEqual(&default_authentication_default, &default_authentication_updated),
				),
			},
			{
				Config: CreateAccDefaultAuthenticationUpdatedAttr("realm", "local"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "local"),
					testAccCheckAciDefaultAuthenticationIdEqual(&default_authentication_default, &default_authentication_updated),
				),
			},
			{
				Config: CreateAccDefaultAuthenticationInitialConfig(defaultAuthentication),
			},
		},
	})
}

func TestAccAciDefaultAuthentication_Negative(t *testing.T) {

	defaultAuthentication, err := aci.GetRemoteDefaultAuthenticationMethodforallLogins(sharedAciClient(), "uni/userext/authrealm/defaultauth")
	if err != nil {
		t.Errorf("reading initial config of Default Authentication")
	}
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDefaultAuthenticationConfig(),
			},
			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr("fallback_check", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr("realm", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr("realm_sub_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccDefaultAuthenticationInitialConfig(defaultAuthentication),
			},
		},
	})
}

func testAccCheckAciDefaultAuthenticationExists(name string, default_authentication *models.DefaultAuthenticationMethodforallLogins) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Default Authentication %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Default Authentication dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		default_authenticationFound := models.DefaultAuthenticationMethodforallLoginsFromContainer(cont)
		if default_authenticationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Default Authentication %s not found", rs.Primary.ID)
		}
		*default_authentication = *default_authenticationFound
		return nil
	}
}

func testAccCheckAciDefaultAuthenticationDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing default_authentication destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_default_authentication" {
			cont, err := client.Get(rs.Primary.ID)
			default_authentication := models.DefaultAuthenticationMethodforallLoginsFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Default Authentication %s Still exists", default_authentication.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciDefaultAuthenticationIdEqual(m1, m2 *models.DefaultAuthenticationMethodforallLogins) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("default_authentication DNs are not equal")
		}
		return nil
	}
}

func CreateAccDefaultAuthenticationConfig() string {
	fmt.Println("=== STEP  testing default_authentication creation")
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
	}
	`)
	return resource
}

func CreateAccDefaultAuthenticationInitialConfig(defaultAuthentication *models.DefaultAuthenticationMethodforallLogins) string {
	fmt.Println("=== STEP  Basic: testing default_authentication creation with Initial Config")
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
	
		description = "%s"
		annotation = "%s"
		name_alias = "%s"
		fallback_check = "%s"
		provider_group = "%s"
		realm = "%s"
		realm_sub_type = "%s"
	}
	`, defaultAuthentication.Description, defaultAuthentication.Annotation, defaultAuthentication.NameAlias, defaultAuthentication.FallbackCheck, defaultAuthentication.ProviderGroup, defaultAuthentication.Realm, defaultAuthentication.RealmSubType)

	return resource
}

func CreateAccDefaultAuthenticationConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing default_authentication creation with optional parameters")
	resource := fmt.Sprint(`
	
	resource "aci_default_authentication" "test" {
	
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_default_authentication"
		fallback_check = "false"		
	}
	`)

	return resource
}

func CreateAccDefaultAuthenticationUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing default_authentication attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
	
		%s = "%s"
	}
	`, attribute, value)
	return resource
}

func CreateAccDefaultAuthenticationUpdatedAttrList(attribute, value string) string {
	fmt.Printf("=== STEP  testing default_authentication attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
	
		%s = %s
	}
	`, attribute, value)
	return resource
}
