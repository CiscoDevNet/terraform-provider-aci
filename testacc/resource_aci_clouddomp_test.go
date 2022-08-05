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

var cloud_domain_profile_default models.CloudDomainProfile

func TestAccAciCloudDomainProfile_Default(t *testing.T) {
	resourceName := "aci_cloud_domain_profile.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudDomainProfileConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudDomainProfileExists(resourceName, &cloud_domain_profile_default),
				),
			},
		},
	})
}

func TestAccAciCloudDomainProfile_Basic(t *testing.T) {
	resourceName := "aci_cloud_domain_profile.test"
	dataSourceName := "data.aci_cloud_domain_profile.test"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudDomainProfileConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "annotation", "test_annotation"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_cloud_domain_profile"),
					resource.TestCheckResourceAttr(resourceName, "site_id", "0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "site_id", resourceName, "site_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccCloudDomainProfileUpdatedAttr("site_id", "500"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_id", "500"),
				),
			},
			{
				Config: CreateAccCloudDomainProfileUpdatedAttr("site_id", "1000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_id", "1000"),
				),
			},
			{
				Config:      CreateAccCloudDomainProfileUpdatedAttr("site_id", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccCloudDomainProfileUpdatedAttr("site_id", "1001"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccCloudDomainProfileUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudDomainProfileUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudDomainProfileUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudDomainProfileUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccCloudDomainProfileUpdatedAttrDS(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCloudDomainProfileFromSavedModel(&cloud_domain_profile_default),
			},
		},
	})
}

func testAccCheckAciCloudDomainProfileExists(name string, cloud_domain_profile *models.CloudDomainProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Domain Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Domain Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_domain_profileFound := models.CloudDomainProfileFromContainer(cont)
		if cloud_domain_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Domain Profile %s not found", rs.Primary.ID)
		}
		*cloud_domain_profile = *cloud_domain_profileFound
		return nil
	}
}

func CreateAccCloudDomainProfileConfig() string {
	fmt.Println("=== STEP  testing cloud_domain_profile creation with required arguments only")
	resource := fmt.Sprintln(`
	
	resource "aci_cloud_domain_profile" "test" {
	
	}
	`)
	return resource
}

func CreateAccCloudDomainProfileFromSavedModel(cloud_domain_profile *models.CloudDomainProfile) string {
	fmt.Println("=== STEP  Basic: setting previously stored value in the server")
	resource := fmt.Sprintf(`
	resource "aci_cloud_domain_profile" "test" {
		description = "%s"
		annotation = "%s"
		name_alias = "%s"
		site_id = "%s"
	}
	`, cloud_domain_profile.Description, cloud_domain_profile.Annotation, cloud_domain_profile.NameAlias, cloud_domain_profile.SiteId)
	return resource
}

func CreateAccCloudDomainProfileConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing cloud_domain_profile resource with all parameters and data source")
	resource := fmt.Sprintln(`
	resource "aci_cloud_domain_profile" "test" {
		description = "created while acceptance testing"
		annotation = "test_annotation"
		name_alias = "test_cloud_domain_profile"
		site_id = "0"
	}

	data "aci_cloud_domain_profile" "test" {
		depends_on = [aci_cloud_domain_profile.test]
	}
	`)
	return resource
}

func CreateAccCloudDomainProfileUpdatedAttrDS(attribute, value string) string {
	fmt.Println("=== STEP  testing cloud_domain_profile Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	data "aci_cloud_domain_profile" "test" {
		%s = "%s"
	}
	`, attribute, value)
	return resource
}

func CreateAccCloudDomainProfileUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing cloud_domain_profile attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_cloud_domain_profile" "test" {
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
