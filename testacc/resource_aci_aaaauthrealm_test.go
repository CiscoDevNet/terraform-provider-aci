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

func TestAccAciAAAAuthentication_Basic(t *testing.T) {
	var authentication_properties_default models.AAAAuthentication
	var authentication_properties_updated models.AAAAuthentication
	resourceName := "aci_authentication_properties.test"
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
				Config: CreateAccAAAAuthenticationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAAAAuthenticationExists(resourceName, &authentication_properties_default),
					resource.TestCheckResourceAttrSet(resourceName, "def_role_policy"),
					resource.TestCheckResourceAttrSet(resourceName, "ping_check"),
					resource.TestCheckResourceAttrSet(resourceName, "retries"),
					resource.TestCheckResourceAttrSet(resourceName, "timeout"),
				),
			},
			{
				Config: CreateAccAAAAuthenticationConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAAAAuthenticationExists(resourceName, &authentication_properties_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_authentication_properties"),
					resource.TestCheckResourceAttr(resourceName, "def_role_policy", "assign-default-role"),
					resource.TestCheckResourceAttr(resourceName, "ping_check", "false"),
					resource.TestCheckResourceAttr(resourceName, "retries", "0"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "1"),
					testAccCheckAciAAAAuthenticationIdEqual(&authentication_properties_default, &authentication_properties_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccAAAAuthenticationInitialConfig(aaaAuthRealm, aaaPingEp),
			},
		},
	})
}

func TestAccAciAAAAuthentication_Update(t *testing.T) {
	var authentication_properties_default models.AAAAuthentication
	var authentication_properties_updated models.AAAAuthentication
	resourceName := "aci_authentication_properties.test"
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
				Config: CreateAccAAAAuthenticationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAAAAuthenticationExists(resourceName, &authentication_properties_default),
				),
			},
			{
				Config: CreateAccAAAAuthenticationUpdatedAttr("def_role_policy", "no-login"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAAAAuthenticationExists(resourceName, &authentication_properties_updated),
					resource.TestCheckResourceAttr(resourceName, "def_role_policy", "no-login"),
					testAccCheckAciAAAAuthenticationIdEqual(&authentication_properties_default, &authentication_properties_updated),
				),
			},
			{
				Config: CreateAccAAAAuthenticationUpdatedAttr("ping_check", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAAAAuthenticationExists(resourceName, &authentication_properties_updated),
					resource.TestCheckResourceAttr(resourceName, "ping_check", "true"),
					testAccCheckAciAAAAuthenticationIdEqual(&authentication_properties_default, &authentication_properties_updated),
				),
			},
			{
				Config: CreateAccAAAAuthenticationUpdatedAttr("retries", "5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAAAAuthenticationExists(resourceName, &authentication_properties_updated),
					resource.TestCheckResourceAttr(resourceName, "retries", "5"),
					testAccCheckAciAAAAuthenticationIdEqual(&authentication_properties_default, &authentication_properties_updated),
				),
			},
			{
				Config: CreateAccAAAAuthenticationUpdatedAttr("retries", "2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAAAAuthenticationExists(resourceName, &authentication_properties_updated),
					resource.TestCheckResourceAttr(resourceName, "retries", "2"),
					testAccCheckAciAAAAuthenticationIdEqual(&authentication_properties_default, &authentication_properties_updated),
				),
			},
			{
				Config: CreateAccAAAAuthenticationUpdatedAttr("timeout", "60"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAAAAuthenticationExists(resourceName, &authentication_properties_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "60"),
					testAccCheckAciAAAAuthenticationIdEqual(&authentication_properties_default, &authentication_properties_updated),
				),
			},
			{
				Config: CreateAccAAAAuthenticationUpdatedAttr("timeout", "30"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAAAAuthenticationExists(resourceName, &authentication_properties_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "30"),
					testAccCheckAciAAAAuthenticationIdEqual(&authentication_properties_default, &authentication_properties_updated),
				),
			},
			{
				Config: CreateAccAAAAuthenticationInitialConfig(aaaAuthRealm, aaaPingEp),
			},
		},
	})
}

func TestAccAciAAAAuthentication_Negative(t *testing.T) {
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
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
				Config: CreateAccAAAAuthenticationConfig(),
			},

			{
				Config:      CreateAccAAAAuthenticationUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAAAAuthenticationUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAAAAuthenticationUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccAAAAuthenticationUpdatedAttr("def_role_policy", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccAAAAuthenticationUpdatedAttr("retries", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccAAAAuthenticationUpdatedAttr("retries", "6"),
				ExpectError: regexp.MustCompile(`is out of range`),
			},
			{
				Config:      CreateAccAAAAuthenticationUpdatedAttr("retries", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccAAAAuthenticationUpdatedAttr("timeout", "0"),
				ExpectError: regexp.MustCompile(`is out of range`),
			},
			{
				Config:      CreateAccAAAAuthenticationUpdatedAttr("timeout", "61"),
				ExpectError: regexp.MustCompile(`is out of range`),
			},
			{
				Config:      CreateAccAAAAuthenticationUpdatedAttr("timeout", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccAAAAuthenticationUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccAAAAuthenticationInitialConfig(aaaAuthRealm, aaaPingEp),
			},
		},
	})
}

func testAccCheckAciAAAAuthenticationExists(name string, authentication_properties *models.AAAAuthentication) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("AAA Authentication %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AAA Authentication dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		authentication_propertiesFound := models.AAAAuthenticationFromContainer(cont)
		if authentication_propertiesFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("AAA Authentication %s not found", rs.Primary.ID)
		}
		*authentication_properties = *authentication_propertiesFound
		return nil
	}
}

func testAccCheckAciAAAAuthenticationIdEqual(m1, m2 *models.AAAAuthentication) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("authentication_properties DNs are not equal")
		}
		return nil
	}
}

func CreateAccAAAAuthenticationConfig() string {
	fmt.Println("=== STEP  testing authentication_properties creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_authentication_properties" "test" {
	
	}
	`)
	return resource
}

func CreateAccAAAAuthenticationConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing authentication_properties creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_authentication_properties" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_authentication_properties"
		def_role_policy = "assign-default-role"
		ping_check = "false"
		retries = "0"
		timeout = "1"
	}
	`)

	return resource
}

func CreateAccAAAAuthenticationUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing authentication_properties attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_authentication_properties" "test" {
	
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
