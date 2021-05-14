package aci

import (
	"fmt"
	"os"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciCloudAWSProvider_Basic(t *testing.T) {
	var cloud_aws_provider models.CloudAWSProvider
	description := "cloud_aws_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudAWSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudAWSProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudAWSProviderExists("aci_cloud_aws_provider.foocloud_aws_provider", &cloud_aws_provider),
					testAccCheckAciCloudAWSProviderAttributes(description, &cloud_aws_provider),
				),
			},
			{
				ResourceName:      "aci_cloud_aws_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAciCloudAWSProviderConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "footenant" {
		description = "Tenant created while acceptance testing"
		name        = "demo_tenant"
	}

	resource "aci_cloud_aws_provider" "foocloud_aws_provider" {
		tenant_dn         = "${aci_tenant.footenant.id}"
		description       = "%s"
		access_key_id     = "%s"
		account_id        = "310368696476"
		annotation        = "tag_aws"
		region            = "us-west-2"
		secret_access_key = "%s"
	}
	  
	`, description, os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_KEY"))
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

func testAccCheckAciCloudAWSProviderAttributes(description string, cloud_aws_provider *models.CloudAWSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_aws_provider.Description {
			return fmt.Errorf("Bad cloud_aws_provider Description %s", cloud_aws_provider.Description)
		}

		if os.Getenv("AWS_ACCESS_KEY_ID") != cloud_aws_provider.AccessKeyId {
			return fmt.Errorf("Bad cloud_aws_provider access_key_id %s", cloud_aws_provider.AccessKeyId)
		}

		if "310368696476" != cloud_aws_provider.AccountId {
			return fmt.Errorf("Bad cloud_aws_provider account_id %s", cloud_aws_provider.AccountId)
		}

		if "tag_aws" != cloud_aws_provider.Annotation {
			return fmt.Errorf("Bad cloud_aws_provider annotation %s", cloud_aws_provider.Annotation)
		}

		if "us-west-2" != cloud_aws_provider.Region {
			return fmt.Errorf("Bad cloud_aws_provider region %s", cloud_aws_provider.Region)
		}

		if os.Getenv("AWS_SECRET_KEY") != cloud_aws_provider.SecretAccessKey {
			return fmt.Errorf("Bad cloud_aws_provider secret_access_key %s", cloud_aws_provider.SecretAccessKey)
		}

		return nil
	}
}
