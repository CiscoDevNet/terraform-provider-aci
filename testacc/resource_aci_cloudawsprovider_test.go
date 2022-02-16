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

func TestAccAciCloudAWSProvider_Basic(t *testing.T) {
	var cloud_aws_provider_default models.CloudAWSProvider
	var cloud_aws_provider_updated models.CloudAWSProvider
	resourceName := "aci_cloud_aws_provider.test"
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudAWSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudAWSProviderWithoutRequired(fvTenantName, "147258369147", "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudAWSProviderConfig(fvTenantName, "147258369147"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudAWSProviderExists(resourceName, &cloud_aws_provider_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "access_key_id", "AKIA5EO6KGEA74MUY673"),
					resource.TestCheckResourceAttr(resourceName, "account_id", "147258369147"),
					resource.TestCheckResourceAttr(resourceName, "email", ""),
					resource.TestCheckResourceAttr(resourceName, "http_proxy", ""),
					resource.TestCheckResourceAttr(resourceName, "is_account_in_org", "no"),
					resource.TestCheckResourceAttr(resourceName, "is_trusted", "no"),
					resource.TestCheckResourceAttr(resourceName, "provider_id", ""),
					resource.TestCheckResourceAttr(resourceName, "region", ""),
				),
			},
			{
				Config: CreateAccCloudAWSProviderConfigWithOptionalValues(fvTenantName, "147258369147"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudAWSProviderExists(resourceName, &cloud_aws_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_cloud_aws_provider"),
					resource.TestCheckResourceAttr(resourceName, "access_key_id", "AKIA5EO6KGEA74MUY656"),
					resource.TestCheckResourceAttr(resourceName, "account_id", "147258369147"),
					resource.TestCheckResourceAttr(resourceName, "email", "testmail"),
					resource.TestCheckResourceAttr(resourceName, "http_proxy", "http_proxy_test"),
					resource.TestCheckResourceAttr(resourceName, "is_account_in_org", "yes"),
					resource.TestCheckResourceAttr(resourceName, "is_trusted", "no"),
					resource.TestCheckResourceAttr(resourceName, "provider_id", "provider_test"),
					resource.TestCheckResourceAttr(resourceName, "region", ""),
					testAccCheckAciCloudAWSProviderIdEqual(&cloud_aws_provider_default, &cloud_aws_provider_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_access_key"},
			},
			{
				Config:      CreateAccCloudAWSProviderRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudAWSProviderConfigWithRequiredParams(rNameUpdated, "147258369147"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudAWSProviderExists(resourceName, &cloud_aws_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					testAccCheckAciCloudAWSProviderIdNotEqual(&cloud_aws_provider_default, &cloud_aws_provider_updated),
				),
			},
		},
	})
}

func TestAccAciCloudAWSProvider_Update(t *testing.T) {
	var cloud_aws_provider_default models.CloudAWSProvider
	var cloud_aws_provider_updated models.CloudAWSProvider
	resourceName := "aci_cloud_aws_provider.test"
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudAWSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudAWSProviderConfig(fvTenantName, "147258369258"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudAWSProviderExists(resourceName, &cloud_aws_provider_default),
				),
			},
			{
				Config: CreateAccCloudAWSProviderUpdatedAttr(fvTenantName, "147258369258", "is_trusted", "yes"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudAWSProviderExists(resourceName, &cloud_aws_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "is_trusted", "yes"),
					testAccCheckAciCloudAWSProviderIdEqual(&cloud_aws_provider_default, &cloud_aws_provider_updated),
				),
			},
			{
				Config: CreateAccCloudAWSProviderConfig(fvTenantName, "147258369258"),
			},
		},
	})
}

func TestAccAciCloudAWSProvider_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudAWSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudAWSProviderConfig(fvTenantName, "147258369369"),
			},
			{
				Config:      CreateAccCloudAWSProviderWithInValidParentDn(rName, "147258369369"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccCloudAWSProviderUpdatedAttr(fvTenantName, "147258369369", "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudAWSProviderUpdatedAttr(fvTenantName, "147258369369", "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudAWSProviderUpdatedAttr(fvTenantName, "147258369369", "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudAWSProviderUpdatedAttr(fvTenantName, "147258369369", "email", acctest.RandString(513)),
				ExpectError: regexp.MustCompile(`property email of awsprovider failed validation for value`),
			},
			{
				Config:      CreateAccCloudAWSProviderUpdatedAttr(fvTenantName, "147258369369", "http_proxy", acctest.RandString(513)),
				ExpectError: regexp.MustCompile(`property httpProxy of awsprovider failed validation for value`),
			},
			{
				Config:      CreateAccCloudAWSProviderUpdatedAttr(fvTenantName, "147258369369", "is_account_in_org", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccCloudAWSProviderUpdatedAttr(fvTenantName, "147258369369", "is_trusted", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccCloudAWSProviderUpdatedAttr(fvTenantName, "147258369369", "provider_id", acctest.RandString(513)),
				ExpectError: regexp.MustCompile(`property providerId of awsprovider failed validation for value`),
			},
			{
				Config:      CreateAccCloudAWSProviderUpdatedAttr(fvTenantName, "147258369369", randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCloudAWSProviderConfig(fvTenantName, "147258369369"),
			},
		},
	})
}

func testAccCheckAciCloudAWSProviderExists(name string, cloud_aws_provider *models.CloudAWSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud AWS Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud AWS Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_aws_providerFound := models.CloudAWSProviderFromContainer(cont)
		if cloud_aws_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud AWS Provider %s not found", rs.Primary.ID)
		}
		*cloud_aws_provider = *cloud_aws_providerFound
		return nil
	}
}

func testAccCheckAciCloudAWSProviderDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing cloud_aws_provider destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_aws_provider" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_aws_provider := models.CloudAWSProviderFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud AWS Provider %s Still exists", cloud_aws_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudAWSProviderIdEqual(m1, m2 *models.CloudAWSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("cloud_aws_provider DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciCloudAWSProviderIdNotEqual(m1, m2 *models.CloudAWSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("cloud_aws_provider DNs are equal")
		}
		return nil
	}
}

func CreateCloudAWSProviderWithoutRequired(fvTenantName, accId, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_aws_provider creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_cloud_aws_provider" "test" {
	#	tenant_dn  = aci_tenant.test.id
		account_id = "%s"
		access_key_id = "AKIA5EO6KGEA74MUY673"
		secret_access_key = "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"
	}
		`

	}
	return fmt.Sprintf(rBlock, fvTenantName, accId)
}

func CreateAccCloudAWSProviderConfigWithRequiredParams(fvTenantName, accId string) string {
	fmt.Println("=== STEP  testing cloud_aws_provider creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_cloud_aws_provider" "test" {
		tenant_dn  = aci_tenant.test.id
		account_id = "%s"
		access_key_id = "AKIA5EO6KGEA74MUY673"
		secret_access_key = "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"
	}
	`, fvTenantName, accId)
	return resource
}
func CreateAccCloudAWSProviderConfig(fvTenantName, accId string) string {
	fmt.Println("=== STEP  testing cloud_aws_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_cloud_aws_provider" "test" {
		tenant_dn  = aci_tenant.test.id
		account_id = "%s"
		access_key_id = "AKIA5EO6KGEA74MUY673"
		secret_access_key = "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"
	}
	`, fvTenantName, accId)
	return resource
}

func CreateAccCloudAWSProviderWithInValidParentDn(rName, accId string) string {
	fmt.Println("=== STEP  Negative Case: testing cloud_aws_provider creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_aaa_domain" "test"{
		name = "%s"
	}
	resource "aci_cloud_aws_provider" "test" {
		tenant_dn  = aci_aaa_domain.test.id	
		account_id = "%s"
		access_key_id = "AKIA5EO6KGEA74MUY673"
		secret_access_key = "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"
	}
	`, rName, accId)
	return resource
}

func CreateAccCloudAWSProviderConfigWithOptionalValues(fvTenantName, accId string) string {
	fmt.Println("=== STEP  Basic: testing cloud_aws_provider creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_cloud_aws_provider" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_aws_provider"
		email = "testmail"
		http_proxy = "http_proxy_test"
		is_account_in_org = "yes"
		is_trusted = "no"
		provider_id = "provider_test"
		region = ""
		account_id = "%s"
		access_key_id = "AKIA5EO6KGEA74MUY656"
		secret_access_key = "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"
	}
	`, fvTenantName, accId)

	return resource
}

func CreateAccCloudAWSProviderRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing cloud_aws_provider updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_cloud_aws_provider" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_aws_provider"
		access_key_id = ""
		account_id = ""
		email = ""
		http_proxy = ""
		is_account_in_org = "yes"
		is_trusted = "yes"
		provider_id = ""
		region = ""
	}
	`)

	return resource
}

func CreateAccCloudAWSProviderUpdatedAttr(fvTenantName, accId, attribute, value string) string {
	fmt.Printf("=== STEP  testing cloud_aws_provider attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_cloud_aws_provider" "test" {
		tenant_dn  = aci_tenant.test.id
		account_id = "%s"
		access_key_id = "AKIA5EO6KGEA74MUY673"
		secret_access_key = "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"
		%s = "%s"
	}
	`, fvTenantName, accId, attribute, value)
	return resource
}
