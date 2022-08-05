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

func TestAccAciRSAProvider_Basic(t *testing.T) {
	var rsa_provider_default models.RSAProvider
	var rsa_provider_updated models.RSAProvider
	resourceName := "aci_rsa_provider.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRSAProviderDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateRSAProviderWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRSAProviderConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRSAProviderExists(resourceName, &rsa_provider_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "auth_port", "1812"),
					resource.TestCheckResourceAttr(resourceName, "auth_protocol", "pap"),
					resource.TestCheckResourceAttr(resourceName, "monitor_server", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_user", "default"),
					resource.TestCheckResourceAttr(resourceName, "retries", "1"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "5"),
				),
			},
			{
				Config: CreateAccRSAProviderConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRSAProviderExists(resourceName, &rsa_provider_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_rsa_provider"),
					resource.TestCheckResourceAttr(resourceName, "auth_port", "1"),
					resource.TestCheckResourceAttr(resourceName, "auth_protocol", "chap"),
					resource.TestCheckResourceAttr(resourceName, "key", "test_key"),
					resource.TestCheckResourceAttr(resourceName, "monitor_server", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_password", "test_password"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_user", "test_user"),
					resource.TestCheckResourceAttr(resourceName, "retries", "5"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "1"),
					testAccCheckAciRSAProviderIdEqual(&rsa_provider_default, &rsa_provider_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key", "monitoring_password"},
			},
			{
				Config:      CreateAccRSAProviderConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccRSAProviderRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccRSAProviderConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRSAProviderExists(resourceName, &rsa_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciRSAProviderIdNotEqual(&rsa_provider_default, &rsa_provider_updated),
				),
			},
		},
	})
}

func TestAccAciRSAProvider_Update(t *testing.T) {
	var rsa_provider_default models.RSAProvider
	var rsa_provider_updated models.RSAProvider
	resourceName := "aci_rsa_provider.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRSAProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRSAProviderConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRSAProviderExists(resourceName, &rsa_provider_default),
				),
			},
			{
				Config: CreateAccRSAProviderUpdatedAttr(rName, "auth_port", "65535"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRSAProviderExists(resourceName, &rsa_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "auth_port", "65535"),
					testAccCheckAciRSAProviderIdEqual(&rsa_provider_default, &rsa_provider_updated),
				),
			},
			{
				Config: CreateAccRSAProviderUpdatedAttr(rName, "auth_port", "32767"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRSAProviderExists(resourceName, &rsa_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "auth_port", "32767"),
					testAccCheckAciRSAProviderIdEqual(&rsa_provider_default, &rsa_provider_updated),
				),
			},
			{
				Config: CreateAccRSAProviderUpdatedAttr(rName, "auth_protocol", "mschap"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRSAProviderExists(resourceName, &rsa_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "auth_protocol", "mschap"),
					testAccCheckAciRSAProviderIdEqual(&rsa_provider_default, &rsa_provider_updated),
				),
			},
			{
				Config: CreateAccRSAProviderUpdatedAttr(rName, "retries", "2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRSAProviderExists(resourceName, &rsa_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "retries", "2"),
					testAccCheckAciRSAProviderIdEqual(&rsa_provider_default, &rsa_provider_updated),
				),
			},
			{
				Config: CreateAccRSAProviderUpdatedAttr(rName, "timeout", "60"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRSAProviderExists(resourceName, &rsa_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "60"),
					testAccCheckAciRSAProviderIdEqual(&rsa_provider_default, &rsa_provider_updated),
				),
			},
			{
				Config: CreateAccRSAProviderUpdatedAttr(rName, "timeout", "30"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRSAProviderExists(resourceName, &rsa_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "30"),
					testAccCheckAciRSAProviderIdEqual(&rsa_provider_default, &rsa_provider_updated),
				),
			},

			{
				Config: CreateAccRSAProviderConfig(rName),
			},
		},
	})
}

func TestAccAciRSAProvider_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRSAProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRSAProviderConfig(rName),
			},

			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "auth_port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "auth_port", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "auth_port", "65536"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "auth_protocol", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "monitor_server", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "monitoring_user", acctest.RandString(33)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "retries", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "retries", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "retries", "6"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "timeout", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "timeout", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, "timeout", "61"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccRSAProviderUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccRSAProviderConfig(rName),
			},
		},
	})
}

func TestAccAciRSAProvider_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRSAProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRSAProviderConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciRSAProviderExists(name string, rsa_provider *models.RSAProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("RSA Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No RSA Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		rsa_providerFound := models.RSAProviderFromContainer(cont)
		if rsa_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("RSA Provider %s not found", rs.Primary.ID)
		}
		*rsa_provider = *rsa_providerFound
		return nil
	}
}

func testAccCheckAciRSAProviderDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing rsa_provider destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_rsa_provider" {
			cont, err := client.Get(rs.Primary.ID)
			rsa_provider := models.RSAProviderFromContainer(cont)
			if err == nil {
				return fmt.Errorf("RSA Provider %s Still exists", rsa_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRSAProviderIdEqual(m1, m2 *models.RSAProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("rsa_provider DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciRSAProviderIdNotEqual(m1, m2 *models.RSAProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("rsa_provider DNs are equal")
		}
		return nil
	}
}

func CreateRSAProviderWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing rsa_provider creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_rsa_provider" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccRSAProviderConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing rsa_provider creation with name =", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_rsa_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccRSAProviderConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing rsa_provider creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_rsa_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccRSAProviderConfig(rName string) string {
	fmt.Println("=== STEP  testing rsa_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_rsa_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccRSAProviderConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple rsa_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_rsa_provider" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccRSAProviderConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing rsa_provider creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_rsa_provider" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_rsa_provider"
		auth_port = "1"
		auth_protocol = "chap"
		key = "test_key"
		monitor_server = "enabled"
		monitoring_password = "test_password"
		monitoring_user = "test_user"
		retries = "5"
		timeout = "1"
		
	}
	`, rName)

	return resource
}

func CreateAccRSAProviderRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing rsa_provider updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_rsa_provider" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_rsa_provider"
		auth_port = "2"
		auth_protocol = "chap"
		key = ""
		monitor_server = "enabled"
		monitoring_password = ""
		monitoring_user = ""
		retries = "1"
		timeout = "1"
		
	}
	`)

	return resource
}

func CreateAccRSAProviderUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing rsa_provider attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_rsa_provider" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccRSAProviderUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing rsa_provider attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_rsa_provider" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
