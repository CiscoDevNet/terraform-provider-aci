package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCloudAWSProviderDataSource_Basic(t *testing.T) {
	resourceName := "aci_cloud_aws_provider.test"
	dataSourceName := "data.aci_cloud_aws_provider.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudAWSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudAWSProviderDSWithoutRequired(fvTenantName, "147258147369", "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudAWSProviderConfigDataSource(fvTenantName, "147258147369"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "access_key_id", resourceName, "access_key_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "account_id", resourceName, "account_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "email", resourceName, "email"),
					resource.TestCheckResourceAttrPair(dataSourceName, "http_proxy", resourceName, "http_proxy"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_account_in_org", resourceName, "is_account_in_org"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_trusted", resourceName, "is_trusted"),
					resource.TestCheckResourceAttrPair(dataSourceName, "provider_id", resourceName, "provider_id"),
				),
			},
			{
				Config:      CreateAccCloudAWSProviderDataSourceUpdate(fvTenantName, "147258147369", randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccCloudAWSProviderDSWithInvalidParentDn(fvTenantName, "147258147369"),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccCloudAWSProviderDataSourceUpdatedResource(fvTenantName, "147258147369", "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccCloudAWSProviderConfigDataSource(fvTenantName, accId string) string {
	fmt.Println("=== STEP  testing cloud_aws_provider Data Source with required arguments only")
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

	data "aci_cloud_aws_provider" "test" {
		tenant_dn  = aci_tenant.test.id
		depends_on = [ aci_cloud_aws_provider.test ]
	}
	`, fvTenantName, accId)
	return resource
}

func CreateCloudAWSProviderDSWithoutRequired(fvTenantName, accId, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_aws_provider Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_cloud_aws_provider" "test" {
		tenant_dn  = aci_tenant.test.id
		account_id = "%s"
		access_key_id = "AKIA5EO6KGEA74MUY673"
		secret_access_key = "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"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_cloud_aws_provider" "test" {
	#	tenant_dn  = aci_tenant.test.id
	
		depends_on = [ aci_cloud_aws_provider.test ]
	}
		`

	}
	return fmt.Sprintf(rBlock, fvTenantName, accId)
}

func CreateAccCloudAWSProviderDSWithInvalidParentDn(fvTenantName, accId string) string {
	fmt.Println("=== STEP  testing cloud_aws_provider Data Source with Invalid Parent Dn")
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

	data "aci_cloud_aws_provider" "test" {
		tenant_dn  = "${aci_tenant.test.id}_invalid"
		depends_on = [ aci_cloud_aws_provider.test ]
	}
	`, fvTenantName, accId)
	return resource
}

func CreateAccCloudAWSProviderDataSourceUpdate(fvTenantName, accId, key, value string) string {
	fmt.Println("=== STEP  testing cloud_aws_provider Data Source with random attribute")
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

	data "aci_cloud_aws_provider" "test" {
		tenant_dn  = aci_tenant.test.id
		%s = "%s"
		depends_on = [ aci_cloud_aws_provider.test ]
	}
	`, fvTenantName, accId, key, value)
	return resource
}

func CreateAccCloudAWSProviderDataSourceUpdatedResource(fvTenantName, accId, key, value string) string {
	fmt.Println("=== STEP  testing cloud_aws_provider Data Source with updated resource")
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

	data "aci_cloud_aws_provider" "test" {
		tenant_dn  = aci_tenant.test.id
		depends_on = [ aci_cloud_aws_provider.test ]
	}
	`, fvTenantName, accId, key, value)
	return resource
}
