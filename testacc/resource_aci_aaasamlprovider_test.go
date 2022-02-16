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
)

func TestAccAciSAMLProvider_Basic(t *testing.T) {
	var saml_provider_default models.SAMLProvider
	var saml_provider_updated models.SAMLProvider
	resourceName := "aci_saml_provider.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSAMLProviderDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateSAMLProviderWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSAMLProviderConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderExists(resourceName, &saml_provider_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "entity_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gui_banner_message", ""),
					resource.TestCheckResourceAttr(resourceName, "https_proxy", ""),
					resource.TestCheckResourceAttr(resourceName, "id_p", "adfs"),
					resource.TestCheckResourceAttr(resourceName, "metadata_url", ""),
					resource.TestCheckResourceAttr(resourceName, "monitor_server", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_user", "default"),
					resource.TestCheckResourceAttr(resourceName, "retries", "1"),
					resource.TestCheckResourceAttr(resourceName, "sig_alg", "SIG_RSA_SHA256"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "5"),
					resource.TestCheckResourceAttr(resourceName, "tp", ""),
					resource.TestCheckResourceAttr(resourceName, "want_assertions_encrypted", "yes"),
					resource.TestCheckResourceAttr(resourceName, "want_assertions_signed", "yes"),
					resource.TestCheckResourceAttr(resourceName, "want_requests_signed", "yes"),
					resource.TestCheckResourceAttr(resourceName, "want_response_signed", "yes"),
				),
			},
			{
				Config: CreateAccSAMLProviderConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderExists(resourceName, &saml_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_saml_provider"),
					resource.TestCheckResourceAttr(resourceName, "entity_id", "entity_id_test"),
					resource.TestCheckResourceAttr(resourceName, "gui_banner_message", "gui_banner_message_test"),
					resource.TestCheckResourceAttr(resourceName, "https_proxy", "https_proxy_test"),
					resource.TestCheckResourceAttr(resourceName, "id_p", "okta"),
					resource.TestCheckResourceAttr(resourceName, "key", "key_test"),
					resource.TestCheckResourceAttr(resourceName, "metadata_url", "metadata_url_test"),
					resource.TestCheckResourceAttr(resourceName, "monitor_server", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_password", "monitoring_password_test"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_user", "monitoring_user_test"),
					resource.TestCheckResourceAttr(resourceName, "retries", "5"),
					resource.TestCheckResourceAttr(resourceName, "sig_alg", "SIG_RSA_SHA1"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "60"),
					resource.TestCheckResourceAttr(resourceName, "tp", "tp_test"),
					resource.TestCheckResourceAttr(resourceName, "want_assertions_encrypted", "no"),
					resource.TestCheckResourceAttr(resourceName, "want_assertions_signed", "no"),
					resource.TestCheckResourceAttr(resourceName, "want_requests_signed", "no"),
					resource.TestCheckResourceAttr(resourceName, "want_response_signed", "no"),
					testAccCheckAciSAMLProviderIdEqual(&saml_provider_default, &saml_provider_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key", "monitoring_password"},
			},
			{
				Config:      CreateAccSAMLProviderConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccSAMLProviderRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSAMLProviderConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderExists(resourceName, &saml_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciSAMLProviderIdNotEqual(&saml_provider_default, &saml_provider_updated),
				),
			},
		},
	})
}

func TestAccAciSAMLProvider_Update(t *testing.T) {
	var saml_provider_default models.SAMLProvider
	var saml_provider_updated models.SAMLProvider
	resourceName := "aci_saml_provider.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSAMLProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSAMLProviderConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderExists(resourceName, &saml_provider_default),
				),
			},
			{
				Config: CreateAccSAMLProviderUpdatedAttr(rName, "retries", "2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderExists(resourceName, &saml_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "retries", "2"),
					testAccCheckAciSAMLProviderIdEqual(&saml_provider_default, &saml_provider_updated),
				),
			},

			{
				Config: CreateAccSAMLProviderUpdatedAttr(rName, "sig_alg", "SIG_RSA_SHA224"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderExists(resourceName, &saml_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "sig_alg", "SIG_RSA_SHA224"),
					testAccCheckAciSAMLProviderIdEqual(&saml_provider_default, &saml_provider_updated),
				),
			},
			{
				Config: CreateAccSAMLProviderUpdatedAttr(rName, "sig_alg", "SIG_RSA_SHA384"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderExists(resourceName, &saml_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "sig_alg", "SIG_RSA_SHA384"),
					testAccCheckAciSAMLProviderIdEqual(&saml_provider_default, &saml_provider_updated),
				),
			},
			{
				Config: CreateAccSAMLProviderUpdatedAttr(rName, "sig_alg", "SIG_RSA_SHA512"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderExists(resourceName, &saml_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "sig_alg", "SIG_RSA_SHA512"),
					testAccCheckAciSAMLProviderIdEqual(&saml_provider_default, &saml_provider_updated),
				),
			},
			{
				Config: CreateAccSAMLProviderUpdatedAttr(rName, "timeout", "27"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderExists(resourceName, &saml_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "27"),
					testAccCheckAciSAMLProviderIdEqual(&saml_provider_default, &saml_provider_updated),
				),
			},
			{
				Config: CreateAccSAMLProviderUpdatedAttr(rName, "id_p", "ping identity"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderExists(resourceName, &saml_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "id_p", "ping identity"),
					testAccCheckAciSAMLProviderIdEqual(&saml_provider_default, &saml_provider_updated),
				),
			},
			{
				Config: CreateAccSAMLProviderConfig(rName),
			},
		},
	})
}

func TestAccAciSAMLProvider_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSAMLProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSAMLProviderConfig(rName),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "entity_id", acctest.RandString(513)),
				ExpectError: regexp.MustCompile(`failed validation for value`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "https_proxy", acctest.RandString(257)),
				ExpectError: regexp.MustCompile(`failed validation for value`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "id_p", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "monitor_server", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "monitoring_user", acctest.RandString(33)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "retries", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "retries", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "retries", "6"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "sig_alg", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "timeout", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "timeout", "4"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "timeout", "61"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "tp", acctest.RandString(27)),
				ExpectError: regexp.MustCompile(`failed validation for value`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "want_assertions_encrypted", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "want_assertions_signed", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "want_requests_signed", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, "want_response_signed", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccSAMLProviderUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSAMLProviderConfig(rName),
			},
		},
	})
}

func TestAccAciSAMLProvider_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSAMLProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSAMLProviderConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciSAMLProviderExists(name string, saml_provider *models.SAMLProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("SAML Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SAML Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		saml_providerFound := models.SAMLProviderFromContainer(cont)
		if saml_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("SAML Provider %s not found", rs.Primary.ID)
		}
		*saml_provider = *saml_providerFound
		return nil
	}
}

func testAccCheckAciSAMLProviderDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing saml_provider destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_saml_provider" {
			cont, err := client.Get(rs.Primary.ID)
			saml_provider := models.SAMLProviderFromContainer(cont)
			if err == nil {
				return fmt.Errorf("SAML Provider %s Still exists", saml_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSAMLProviderIdEqual(m1, m2 *models.SAMLProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("saml_provider DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSAMLProviderIdNotEqual(m1, m2 *models.SAMLProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("saml_provider DNs are equal")
		}
		return nil
	}
}

func CreateSAMLProviderWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing saml_provider creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_saml_provider" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccSAMLProviderConfigWithRequiredParams(rName string) string {
	fmt.Printf("=== STEP  testing saml_provider creation with name %s\n", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccSAMLProviderConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing saml_provider creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSAMLProviderConfig(rName string) string {
	fmt.Println("=== STEP  testing saml_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSAMLProviderConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple saml_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccSAMLProviderConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing saml_provider creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider" "test" {
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_saml_provider"
		entity_id = "entity_id_test"
		gui_banner_message = "gui_banner_message_test"
		https_proxy = "https_proxy_test"
		id_p = "okta"
		key = "key_test"
		metadata_url = "metadata_url_test"
		monitor_server = "disabled"
		monitoring_password = "monitoring_password_test"
		monitoring_user = "monitoring_user_test"
		retries = "5"
		sig_alg = "SIG_RSA_SHA1"
		timeout = "60"
		tp = "tp_test"
		want_assertions_encrypted = "no"
		want_assertions_signed = "no"
		want_requests_signed = "no"
		want_response_signed = "no"
		
	}
	`, rName)

	return resource
}

func CreateAccSAMLProviderRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing saml_provider updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_saml_provider" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_saml_provider"
		entity_id = ""
		gui_banner_message = ""
		https_proxy = ""
		id_p = "okta"
		key = ""
		metadata_url = ""
		monitor_server = "enabled"
		monitoring_password = ""
		monitoring_user = ""
		retries = "1"
		sig_alg = "SIG_RSA_SHA1"
		timeout = "6"
		tp = ""
		want_assertions_encrypted = "no"
		want_assertions_signed = "no"
		want_requests_signed = "no"
		want_response_signed = "no"
		
	}
	`)

	return resource
}

func CreateAccSAMLProviderUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing saml_provider attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
