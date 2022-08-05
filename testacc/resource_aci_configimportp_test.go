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

func TestAccAciConfigurationImportPolicy_Basic(t *testing.T) {
	var configuration_import_policy_default models.ConfigurationImportPolicy
	var configuration_import_policy_updated models.ConfigurationImportPolicy
	resourceName := "aci_configuration_import_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConfigurationImportPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateConfigurationImportPolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccConfigurationImportPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationImportPolicyExists(resourceName, &configuration_import_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "fail_on_decrypt_errors", "yes"),
					resource.TestCheckResourceAttr(resourceName, "file_name", "file.tar.gz"),
					resource.TestCheckResourceAttr(resourceName, "import_mode", "atomic"),
					resource.TestCheckResourceAttr(resourceName, "import_type", "merge"),
					resource.TestCheckResourceAttr(resourceName, "snapshot", "no"),
				),
			},
			{
				Config: CreateAccConfigurationImportPolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationImportPolicyExists(resourceName, &configuration_import_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_configuration_import_policy"),

					resource.TestCheckResourceAttr(resourceName, "admin_st", "triggered"),

					resource.TestCheckResourceAttr(resourceName, "fail_on_decrypt_errors", "no"),

					resource.TestCheckResourceAttr(resourceName, "file_name", "file2.tar.gz"),

					resource.TestCheckResourceAttr(resourceName, "import_mode", "best-effort"),

					resource.TestCheckResourceAttr(resourceName, "import_type", "merge"),

					resource.TestCheckResourceAttr(resourceName, "snapshot", "yes"),

					testAccCheckAciConfigurationImportPolicyIdEqual(&configuration_import_policy_default, &configuration_import_policy_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_st"},
			},
			{
				Config:      CreateAccConfigurationImportPolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccConfigurationImportPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccConfigurationImportPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationImportPolicyExists(resourceName, &configuration_import_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciConfigurationImportPolicyIdNotEqual(&configuration_import_policy_default, &configuration_import_policy_updated),
				),
			},
		},
	})
}

func TestAccAciConfigurationImportPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConfigurationImportPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccConfigurationImportPolicyConfig(rName),
			},

			{
				Config:      CreateAccConfigurationImportPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccConfigurationImportPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccConfigurationImportPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccConfigurationImportPolicyUpdatedAttr(rName, "admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConfigurationImportPolicyUpdatedAttr(rName, "fail_on_decrypt_errors", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConfigurationImportPolicyUpdatedAttr(rName, "import_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConfigurationImportPolicyUpdatedAttr(rName, "import_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConfigurationImportPolicyUpdatedAttr(rName, "snapshot", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConfigurationImportPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccConfigurationImportPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciConfigurationImportPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConfigurationImportPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccConfigurationImportPolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciConfigurationImportPolicyExists(name string, configuration_import_policy *models.ConfigurationImportPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Configuration Import Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Configuration Import Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		configuration_import_policyFound := models.ConfigurationImportPolicyFromContainer(cont)
		if configuration_import_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Configuration Import Policy %s not found", rs.Primary.ID)
		}
		*configuration_import_policy = *configuration_import_policyFound
		return nil
	}
}

func testAccCheckAciConfigurationImportPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing configuration_import_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_configuration_import_policy" {
			cont, err := client.Get(rs.Primary.ID)
			configuration_import_policy := models.ConfigurationImportPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Configuration Import Policy %s Still exists", configuration_import_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciConfigurationImportPolicyIdEqual(m1, m2 *models.ConfigurationImportPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("configuration_import_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciConfigurationImportPolicyIdNotEqual(m1, m2 *models.ConfigurationImportPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("configuration_import_policy DNs are equal")
		}
		return nil
	}
}

func CreateConfigurationImportPolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing configuration_import_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_configuration_import_policy" "test" {
	
	#	name  = "%s"
		file_name  = "file.tar.gz"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccConfigurationImportPolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing configuration_import_policy creation with name =", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_import_policy" "test" {
	
		name  = "%s"
		file_name  = "file.tar.gz"
	}
	`, rName)
	return resource
}
func CreateAccConfigurationImportPolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing configuration_import_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_import_policy" "test" {
	
		name  = "%s"
		file_name  = "file.tar.gz"
	}
	`, rName)
	return resource
}

func CreateAccConfigurationImportPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing configuration_import_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_import_policy" "test" {
	
		name  = "%s"
		file_name  = "file.tar.gz"
	}
	`, rName)
	return resource
}

func CreateAccConfigurationImportPolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple configuration_import_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_import_policy" "test" {
	
		name  = "%s_${count.index}"
		file_name  = "file.tar.gz"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccConfigurationImportPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing configuration_import_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_import_policy" "test" {
	
		name  = "%s"
		file_name  = "file2.tar.gz"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_configuration_import_policy"
		admin_st = "triggered"
		fail_on_decrypt_errors = "no"
		import_mode = "best-effort"
		import_type = "merge"
		snapshot = "yes"
		
	}
	`, rName)

	return resource
}

func CreateAccConfigurationImportPolicyConfigSecondWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing configuration_import_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_import_policy" "test" {
	
		name  = "%s"
		file_name  = "file2.tar.gz"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_configuration_import_policy"
		admin_st = "triggered"
		fail_on_decrypt_errors = "no"
		import_mode = "atomic"
		import_type = "replace"
		snapshot = "yes"
		
	}
	`, rName)

	return resource
}

func CreateAccConfigurationImportPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing configuration_import_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_configuration_import_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_configuration_import_policy"
		admin_st = "triggered"
		fail_on_decrypt_errors = "no"
		file_name = ""
		import_mode = "best-effort"
		import_type = "replace"
		snapshot = "yes"
		
	}
	`)

	return resource
}

func CreateAccConfigurationImportPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing configuration_import_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_import_policy" "test" {
	
		name  = "%s"
		file_name  = "file.tar.gz"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
