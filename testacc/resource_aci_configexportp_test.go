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

func TestAccAciConfigurationExportPolicy_Basic(t *testing.T) {
	var configuration_export_policy_default models.ConfigurationExportPolicy
	var configuration_export_policy_updated models.ConfigurationExportPolicy
	resourceName := "aci_configuration_export_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConfigurationExportPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateConfigurationExportPolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccConfigurationExportPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationExportPolicyExists(resourceName, &configuration_export_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "format", "json"),
					resource.TestCheckResourceAttr(resourceName, "include_secure_fields", "yes"),
					resource.TestCheckResourceAttr(resourceName, "max_snapshot_count", "global-limit"),
					resource.TestCheckResourceAttr(resourceName, "snapshot", "no"),
					resource.TestCheckResourceAttr(resourceName, "target_dn", ""),
				),
			},
			{
				Config: CreateAccConfigurationExportPolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationExportPolicyExists(resourceName, &configuration_export_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_configuration_export_policy"),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "triggered"),
					resource.TestCheckResourceAttr(resourceName, "format", "xml"),
					resource.TestCheckResourceAttr(resourceName, "include_secure_fields", "no"),
					resource.TestCheckResourceAttr(resourceName, "max_snapshot_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "snapshot", "yes"),
					resource.TestCheckResourceAttr(resourceName, "target_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAciConfigurationExportPolicyIdEqual(&configuration_export_policy_default, &configuration_export_policy_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_st"},
			},
			{
				Config:      CreateAccConfigurationExportPolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccConfigurationExportPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccConfigurationExportPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationExportPolicyExists(resourceName, &configuration_export_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciConfigurationExportPolicyIdNotEqual(&configuration_export_policy_default, &configuration_export_policy_updated),
				),
			},
		},
	})
}

func TestAccAciConfigurationExportPolicy_Update(t *testing.T) {
	var configuration_export_policy_default models.ConfigurationExportPolicy
	var configuration_export_policy_updated models.ConfigurationExportPolicy
	resourceName := "aci_configuration_export_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConfigurationExportPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccConfigurationExportPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationExportPolicyExists(resourceName, &configuration_export_policy_default),
				),
			},
			{
				Config: CreateAccConfigurationExportPolicyUpdatedAttr(rName, "max_snapshot_count", "10"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationExportPolicyExists(resourceName, &configuration_export_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "max_snapshot_count", "10"),
					testAccCheckAciConfigurationExportPolicyIdEqual(&configuration_export_policy_default, &configuration_export_policy_updated),
				),
			},
			{
				Config: CreateAccConfigurationExportPolicyUpdatedAttr(rName, "max_snapshot_count", "5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationExportPolicyExists(resourceName, &configuration_export_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "max_snapshot_count", "5"),
					testAccCheckAciConfigurationExportPolicyIdEqual(&configuration_export_policy_default, &configuration_export_policy_updated),
				),
			},

			{
				Config: CreateAccConfigurationExportPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciConfigurationExportPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConfigurationExportPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccConfigurationExportPolicyConfig(rName),
			},

			{
				Config:      CreateAccConfigurationExportPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccConfigurationExportPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccConfigurationExportPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccConfigurationExportPolicyUpdatedAttr(rName, "admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConfigurationExportPolicyUpdatedAttr(rName, "format", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConfigurationExportPolicyUpdatedAttr(rName, "include_secure_fields", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConfigurationExportPolicyUpdatedAttr(rName, "max_snapshot_count", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccConfigurationExportPolicyUpdatedAttr(rName, "max_snapshot_count", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccConfigurationExportPolicyUpdatedAttr(rName, "max_snapshot_count", "11"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccConfigurationExportPolicyUpdatedAttr(rName, "snapshot", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConfigurationExportPolicyUpdatedAttr(rName, "target_dn", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccConfigurationExportPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccConfigurationExportPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciConfigurationExportPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConfigurationExportPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccConfigurationExportPolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciConfigurationExportPolicyExists(name string, configuration_export_policy *models.ConfigurationExportPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("ConfigurationExportPolicy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ConfigurationExportPolicy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		configuration_export_policyFound := models.ConfigurationExportPolicyFromContainer(cont)
		if configuration_export_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("ConfigurationExportPolicy %s not found", rs.Primary.ID)
		}
		*configuration_export_policy = *configuration_export_policyFound
		return nil
	}
}

func testAccCheckAciConfigurationExportPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing configuration_export_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_configuration_export_policy" {
			cont, err := client.Get(rs.Primary.ID)
			configuration_export_policy := models.ConfigurationExportPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("ConfigurationExportPolicy %s Still exists", configuration_export_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciConfigurationExportPolicyIdEqual(m1, m2 *models.ConfigurationExportPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("configuration_export_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciConfigurationExportPolicyIdNotEqual(m1, m2 *models.ConfigurationExportPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("configuration_export_policy DNs are equal")
		}
		return nil
	}
}

func CreateConfigurationExportPolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing configuration_export_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_configuration_export_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccConfigurationExportPolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing configuration_export_policy creation with updated name =", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_export_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccConfigurationExportPolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing configuration_export_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_export_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccConfigurationExportPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing configuration_export_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_export_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccConfigurationExportPolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple configuration_export_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_export_policy" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccConfigurationExportPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing configuration_export_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_export_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_configuration_export_policy"
		admin_st = "triggered"
		format = "xml"
		include_secure_fields = "no"
		max_snapshot_count = "1"
		snapshot = "yes"
		target_dn = aci_tenant.test.id
		
	}

	resource "aci_tenant" "test" {
		name = "%s"
	}
	`, rName, rName)

	return resource
}

func CreateAccConfigurationExportPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing configuration_export_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_configuration_export_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_configuration_export_policy"
		admin_st = "triggered"
		format = "xml"
		include_secure_fields = "no"
		max_snapshot_count = "1"
		snapshot = "yes"
		target_dn = ""
		
	}
	`)

	return resource
}

func CreateAccConfigurationExportPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing configuration_export_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_export_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccConfigurationExportPolicyUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing configuration_export_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_configuration_export_policy" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
