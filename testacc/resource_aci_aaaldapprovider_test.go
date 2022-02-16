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

func TestAccAciLDAPProvider_Basic(t *testing.T) {
	var ldap_provider_default models.LDAPProvider
	var ldap_provider_updated models.LDAPProvider
	resourceName := "aci_ldap_provider.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPProviderDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLDAPProviderWithoutRequired(rName, "duo", "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLDAPProviderWithoutRequired(rName, "duo", "type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLDAPProviderConfig(rName, "duo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "duo"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "attribute", "CiscoAVPair"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "ssl_validation_level", "strict"),
					resource.TestCheckResourceAttr(resourceName, "basedn", ""),
					resource.TestCheckResourceAttr(resourceName, "enable_ssl", "no"),
					resource.TestCheckResourceAttr(resourceName, "filter", "sAMAccountName=$userid"),
					resource.TestCheckResourceAttr(resourceName, "monitor_server", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_user", "default"),
					resource.TestCheckResourceAttr(resourceName, "port", "389"),
					resource.TestCheckResourceAttr(resourceName, "retries", "1"),
					resource.TestCheckResourceAttr(resourceName, "rootdn", ""),
					resource.TestCheckResourceAttr(resourceName, "timeout", "30"),
				),
			},
			{
				Config: CreateAccLDAPProviderConfigWithOptionalValues(rName, "duo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_ldap_provider"),
					resource.TestCheckResourceAttr(resourceName, "ssl_validation_level", "permissive"),
					resource.TestCheckResourceAttr(resourceName, "type", "duo"),
					resource.TestCheckResourceAttr(resourceName, "attribute", "memberOf"),
					resource.TestCheckResourceAttr(resourceName, "basedn", "CN=Users,DC=host,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "enable_ssl", "yes"),
					resource.TestCheckResourceAttr(resourceName, "filter", "cn=$userid"),
					resource.TestCheckResourceAttr(resourceName, "key", "test_key"),
					resource.TestCheckResourceAttr(resourceName, "monitor_server", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_password", "test_password"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_user", "test_user"),
					resource.TestCheckResourceAttr(resourceName, "port", "1"),
					resource.TestCheckResourceAttr(resourceName, "retries", "5"),
					resource.TestCheckResourceAttr(resourceName, "rootdn", "CN=admin,CN=Users,DC=host,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "60"),
					testAccCheckAciLDAPProviderIdEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key", "monitoring_password"},
			},
			{
				Config:      CreateAccLDAPProviderConfigUpdatedName(acctest.RandString(65), "duo"),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccLDAPProviderConfigWithRequiredParams(rName, rName),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccLDAPProviderConfigWithRequiredParams(rNameUpdated, "duo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "type", "duo"),
					testAccCheckAciLDAPProviderIdNotEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
			{
				Config: CreateAccLDAPProviderConfig(rName, "duo"),
			},
			{
				Config: CreateAccLDAPProviderConfigWithRequiredParams(rName, "ldap"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "ldap"),
					testAccCheckAciLDAPProviderIdNotEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
		},
	})
}

func TestAccAciLDAPProvider_Update(t *testing.T) {
	var ldap_provider_default models.LDAPProvider
	var ldap_provider_updated models.LDAPProvider
	resourceName := "aci_ldap_provider.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLDAPProviderConfig(rName, "ldap"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_default),
				),
			},
			{
				Config: CreateAccLDAPProviderUpdatedAttr(rName, "ldap", "port", "65535"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "port", "65535"),
					testAccCheckAciLDAPProviderIdEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
			{
				Config: CreateAccLDAPProviderUpdatedAttr(rName, "ldap", "port", "32767"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "port", "32767"),
					testAccCheckAciLDAPProviderIdEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
			{
				Config: CreateAccLDAPProviderUpdatedAttr(rName, "ldap", "retries", "2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "retries", "2"),
					testAccCheckAciLDAPProviderIdEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
			{
				Config: CreateAccLDAPProviderUpdatedAttr(rName, "ldap", "timeout", "5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "5"),
					testAccCheckAciLDAPProviderIdEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
			{
				Config: CreateAccLDAPProviderUpdatedAttr(rName, "ldap", "timeout", "18"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "18"),
					testAccCheckAciLDAPProviderIdEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
			{
				Config: CreateAccLDAPProviderConfig(rName, "ldap"),
			},
		},
	})
}

func TestAccAciLDAPProvider_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLDAPProviderConfig(rName, "duo"),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "ssl_validation_level", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "attribute", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "basedn", acctest.RandString(128)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "enable_ssl", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "filter", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "monitor_server", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "monitoring_user", acctest.RandString(33)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "port", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "port", "65536"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "retries", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "retries", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "retries", "6"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "rootdn", acctest.RandString(128)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "timeout", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "timeout", "4"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", "timeout", "61"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "duo", randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLDAPProviderConfig(rName, "duo"),
			},
		},
	})
}

func testAccCheckAciLDAPProviderExists(name string, ldap_provider *models.LDAPProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LDAP Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LDAP Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ldap_providerFound := models.LDAPProviderFromContainer(cont)
		if ldap_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LDAP Provider %s not found", rs.Primary.ID)
		}
		*ldap_provider = *ldap_providerFound
		return nil
	}
}

func testAccCheckAciLDAPProviderDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing ldap_provider destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ldap_provider" {
			cont, err := client.Get(rs.Primary.ID)
			ldap_provider := models.LDAPProviderFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LDAP Provider %s Still exists", ldap_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLDAPProviderIdEqual(m1, m2 *models.LDAPProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("ldap_provider DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLDAPProviderIdNotEqual(m1, m2 *models.LDAPProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("ldap_provider DNs are equal")
		}
		return nil
	}
}

func CreateLDAPProviderWithoutRequired(rName, proType, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ldap_provider creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_ldap_provider" "test" {
	
	#	name  = "%s"
		type = "%s"
	}
		`
	case "type":
		rBlock += `
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
	#	type = "%s"
	}
	`
	}
	return fmt.Sprintf(rBlock, rName, proType)
}

func CreateAccLDAPProviderConfigWithRequiredParams(rName, proType string) string {
	fmt.Printf("=== STEP  testing ldap_provider creation with name %s and type %s\n", rName, proType)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
		type = "%s"
	}
	`, rName, proType)
	return resource
}
func CreateAccLDAPProviderConfigUpdatedName(rName, proType string) string {
	fmt.Println("=== STEP  testing ldap_provider creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
		type = "%s"
	}
	`, rName, proType)
	return resource
}

func CreateAccLDAPProviderConfig(rName, proType string) string {
	fmt.Println("=== STEP  testing ldap_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
		type = "%s"
	}
	`, rName, proType)
	return resource
}

func CreateAccLDAPProviderConfigWithOptionalValues(rName, proType string) string {
	fmt.Println("=== STEP  Basic: testing ldap_provider creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ldap_provider"
		ssl_validation_level = "permissive"
		attribute = "memberOf"
		basedn = "CN=Users,DC=host,DC=com"
		enable_ssl = "yes"
		filter = "cn=$userid"
		key = "test_key"
		monitor_server = "enabled"
		monitoring_password = "test_password"
		monitoring_user = "test_user"
		port = "1"
		retries = "5"
		rootdn = "CN=admin,CN=Users,DC=host,DC=com"
		timeout = "60"
		type = "%s"
	}
	`, rName, proType)

	return resource
}

func CreateAccLDAPProviderRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing ldap_provider updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_ldap_provider" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ldap_provider"
		ssl_validation_level = "permissive"
		attribute = ""
		basedn = ""
		enable_ssl = "yes"
		filter = ""
		key = ""
		monitor_server = "enabled"
		monitoring_password = ""
		monitoring_user = ""
		port = "2"
		retries = "1"
		rootdn = ""
		timeout = "6"
		
	}
	`)

	return resource
}

func CreateAccLDAPProviderUpdatedAttr(rName, proType, attribute, value string) string {
	fmt.Printf("=== STEP  testing ldap_provider attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
		type = "%s"
		%s = "%s"
	}
	`, rName, proType, attribute, value)
	return resource
}
