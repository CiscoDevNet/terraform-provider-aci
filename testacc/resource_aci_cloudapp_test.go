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

func TestAccAciCloudApplicationcontainer_Basic(t *testing.T) {
	var cloud_applicationcontainer_default models.CloudApplicationcontainer
	var cloud_applicationcontainer_updated models.CloudApplicationcontainer
	resourceName := "aci_cloud_applicationcontainer.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudApplicationcontainerDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudApplicationcontainerWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudApplicationcontainerWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudApplicationcontainerConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudApplicationcontainerExists(resourceName, &cloud_applicationcontainer_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccCloudApplicationcontainerConfigWithOptionalValues(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudApplicationcontainerExists(resourceName, &cloud_applicationcontainer_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_cloud_applicationcontainer"),

					testAccCheckAciCloudApplicationcontainerIdEqual(&cloud_applicationcontainer_default, &cloud_applicationcontainer_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccCloudApplicationcontainerConfigUpdatedName(rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccCloudApplicationcontainerRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudApplicationcontainerConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudApplicationcontainerExists(resourceName, &cloud_applicationcontainer_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciCloudApplicationcontainerIdNotEqual(&cloud_applicationcontainer_default, &cloud_applicationcontainer_updated),
				),
			},
			{
				Config: CreateAccCloudApplicationcontainerConfig(rName, rName),
			},
			{
				Config: CreateAccCloudApplicationcontainerConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudApplicationcontainerExists(resourceName, &cloud_applicationcontainer_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciCloudApplicationcontainerIdNotEqual(&cloud_applicationcontainer_default, &cloud_applicationcontainer_updated),
				),
			},
		},
	})
}

func TestAccAciCloudApplicationcontainer_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudApplicationcontainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudApplicationcontainerConfig(rName, rName),
			},
			{
				Config:      CreateAccCloudApplicationcontainerWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccCloudApplicationcontainerUpdatedAttr(rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudApplicationcontainerUpdatedAttr(rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudApplicationcontainerUpdatedAttr(rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccCloudApplicationcontainerUpdatedAttr(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCloudApplicationcontainerConfig(rName, rName),
			},
		},
	})
}

func TestAccAciCloudApplicationcontainer_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudApplicationcontainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudApplicationcontainerConfigMultiple(fvTenantName, rName),
			},
		},
	})
}

func testAccCheckAciCloudApplicationcontainerExists(name string, cloud_applicationcontainer *models.CloudApplicationcontainer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Applicationcontainer %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Applicationcontainer dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_applicationcontainerFound := models.CloudApplicationcontainerFromContainer(cont)
		if cloud_applicationcontainerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Applicationcontainer %s not found", rs.Primary.ID)
		}
		*cloud_applicationcontainer = *cloud_applicationcontainerFound
		return nil
	}
}

func testAccCheckAciCloudApplicationcontainerDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing cloud_applicationcontainer destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_applicationcontainer" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_applicationcontainer := models.CloudApplicationcontainerFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Applicationcontainer %s Still exists", cloud_applicationcontainer.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudApplicationcontainerIdEqual(m1, m2 *models.CloudApplicationcontainer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("cloud_applicationcontainer DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciCloudApplicationcontainerIdNotEqual(m1, m2 *models.CloudApplicationcontainer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("cloud_applicationcontainer DNs are equal")
		}
		return nil
	}
}

func CreateCloudApplicationcontainerWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_applicationcontainer creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_cloud_applicationcontainer" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccCloudApplicationcontainerConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing cloud_applicationcontainer with tenant name %s and resource name %s\n", fvTenantName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}
func CreateAccCloudApplicationcontainerConfigUpdatedName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing cloud_applicationcontainer creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccCloudApplicationcontainerConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing cloud_applicationcontainer creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccCloudApplicationcontainerConfigMultiple(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing multiple cloud_applicationcontainer creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccCloudApplicationcontainerWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing cloud_applicationcontainer creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	resource "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccCloudApplicationcontainerConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_applicationcontainer creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_applicationcontainer"
		
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccCloudApplicationcontainerRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing cloud_applicationcontainer updation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_cloud_applicationcontainer" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_cloud_applicationcontainer"
	}
	`)

	return resource
}

func CreateAccCloudApplicationcontainerUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing cloud_applicationcontainer attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
