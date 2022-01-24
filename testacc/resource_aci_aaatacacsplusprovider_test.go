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

func TestAccAciTACACSProvider_Basic(t *testing.T) {
	var tacacs_provider_default models.TACACSProvider
	var tacacs_provider_updated models.TACACSProvider
	resourceName := "aci_tacacs_provider.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSProviderDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateTACACSProviderWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTACACSProviderConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSProviderExists(resourceName, &tacacs_provider_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "auth_protocol", "pap"),
					resource.TestCheckResourceAttr(resourceName, "monitor_server", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_user", "default"),
					resource.TestCheckResourceAttr(resourceName, "port", "49"),
					resource.TestCheckResourceAttr(resourceName, "retries", "1"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "5"),
				),
			},
			{
				Config: CreateAccTACACSProviderConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSProviderExists(resourceName, &tacacs_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_tacacs_provider"),
					resource.TestCheckResourceAttr(resourceName, "auth_protocol", "chap"),
					resource.TestCheckResourceAttr(resourceName, "key", "example_key"),
					resource.TestCheckResourceAttr(resourceName, "monitor_server", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_password", "monitoring_password_example"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_user", "monitoring_user_example"),
					resource.TestCheckResourceAttr(resourceName, "port", "1"),
					resource.TestCheckResourceAttr(resourceName, "retries", "5"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "1"),
					testAccCheckAciTACACSProviderIdEqual(&tacacs_provider_default, &tacacs_provider_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key", "monitoring_password"},
			},
			{
				Config:      CreateAccTACACSProviderConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccTACACSProviderRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccTACACSProviderConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSProviderExists(resourceName, &tacacs_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciTACACSProviderIdNotEqual(&tacacs_provider_default, &tacacs_provider_updated),
				),
			},
		},
	})
}

func TestAccAciTACACSProvider_Update(t *testing.T) {
	var tacacs_provider_default models.TACACSProvider
	var tacacs_provider_updated models.TACACSProvider
	resourceName := "aci_tacacs_provider.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSProviderConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSProviderExists(resourceName, &tacacs_provider_default),
				),
			},
			{
				Config: CreateAccTACACSProviderUpdatedAttr(rName, "port", "65535"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSProviderExists(resourceName, &tacacs_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "port", "65535"),
					testAccCheckAciTACACSProviderIdEqual(&tacacs_provider_default, &tacacs_provider_updated),
				),
			},
			{
				Config: CreateAccTACACSProviderUpdatedAttr(rName, "port", "32767"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSProviderExists(resourceName, &tacacs_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "port", "32767"),
					testAccCheckAciTACACSProviderIdEqual(&tacacs_provider_default, &tacacs_provider_updated),
				),
			},
			{
				Config: CreateAccTACACSProviderUpdatedAttr(rName, "retries", "2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSProviderExists(resourceName, &tacacs_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "retries", "2"),
					testAccCheckAciTACACSProviderIdEqual(&tacacs_provider_default, &tacacs_provider_updated),
				),
			},
			{
				Config: CreateAccTACACSProviderUpdatedAttr(rName, "timeout", "60"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSProviderExists(resourceName, &tacacs_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "60"),
					testAccCheckAciTACACSProviderIdEqual(&tacacs_provider_default, &tacacs_provider_updated),
				),
			},
			{
				Config: CreateAccTACACSProviderUpdatedAttr(rName, "timeout", "30"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSProviderExists(resourceName, &tacacs_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "30"),
					testAccCheckAciTACACSProviderIdEqual(&tacacs_provider_default, &tacacs_provider_updated),
				),
			},
			{
				Config: CreateAccTACACSProviderUpdatedAttr(rName, "auth_protocol", "mschap"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSProviderExists(resourceName, &tacacs_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "auth_protocol", "mschap"),
					testAccCheckAciTACACSProviderIdEqual(&tacacs_provider_default, &tacacs_provider_updated),
				),
			},
			{
				Config: CreateAccTACACSProviderConfig(rName),
			},
		},
	})
}

func TestAccAciTACACSProvider_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSProviderConfig(rName),
			},

			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "auth_protocol", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "monitor_server", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "monitoring_user", acctest.RandString(33)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "port", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "port", "65536"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "retries", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "retries", "-1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "retries", "6"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "timeout", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "timeout", "-1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, "timeout", "61"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccTACACSProviderUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccTACACSProviderConfig(rName),
			},
		},
	})
}

func TestAccAciTACACSProvider_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSProviderConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciTACACSProviderExists(name string, tacacs_provider *models.TACACSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("TACACS Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No TACACS Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tacacs_providerFound := models.TACACSProviderFromContainer(cont)
		if tacacs_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("TACACS Provider %s not found", rs.Primary.ID)
		}
		*tacacs_provider = *tacacs_providerFound
		return nil
	}
}

func testAccCheckAciTACACSProviderDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing tacacs_provider destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tacacs_provider" {
			cont, err := client.Get(rs.Primary.ID)
			tacacs_provider := models.TACACSProviderFromContainer(cont)
			if err == nil {
				return fmt.Errorf("TACACS Provider %s Still exists", tacacs_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTACACSProviderIdEqual(m1, m2 *models.TACACSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("tacacs_provider DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciTACACSProviderIdNotEqual(m1, m2 *models.TACACSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("tacacs_provider DNs are equal")
		}
		return nil
	}
}

func CreateTACACSProviderWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_provider creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_tacacs_provider" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccTACACSProviderConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing tacacs_provider creation with name", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccTACACSProviderConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing tacacs_provider creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccTACACSProviderConfig(rName string) string {
	fmt.Println("=== STEP  testing tacacs_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccTACACSProviderConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple tacacs_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccTACACSProviderConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_provider creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tacacs_provider"
		auth_protocol = "chap"
		key = "example_key"
		monitor_server = "enabled"
		monitoring_password = "monitoring_password_example"
		monitoring_user = "monitoring_user_example"
		port = "1"
		retries = "5"
		timeout = "1"
		
	}
	`, rName)

	return resource
}

func CreateAccTACACSProviderRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing tacacs_provider updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_tacacs_provider" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tacacs_provider"
		auth_protocol = "chap"
		key = ""
		monitor_server = "enabled"
		monitoring_password = ""
		monitoring_user = ""
		port = "2"
		retries = "1"
		timeout = "1"
		
	}
	`)

	return resource
}

func CreateAccTACACSProviderUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing tacacs_provider attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccTACACSProviderUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing tacacs_provider attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
