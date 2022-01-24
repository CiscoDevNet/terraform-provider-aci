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

func TestAccAciCloudEndpointSelector_Basic(t *testing.T) {
	var cloud_endpoint_selector_default models.CloudEndpointSelector
	var cloud_endpoint_selector_updated models.CloudEndpointSelector
	resourceName := "aci_cloud_endpoint_selector.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudEndpointSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudEndpointSelectorWithoutRequired(rName, rName, rName, rName, "cloud_epg_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudEndpointSelectorWithoutRequired(rName, rName, rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudEndpointSelectorConfig(rName, rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorExists(resourceName, &cloud_endpoint_selector_default),
					resource.TestCheckResourceAttr(resourceName, "cloud_epg_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudepg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "match_expression", ""),
				),
			},
			{
				Config: CreateAccCloudEndpointSelectorConfigWithOptionalValues(rName, rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorExists(resourceName, &cloud_endpoint_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_epg_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudepg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_cloud_endpoint_selector"),
					resource.TestCheckResourceAttr(resourceName, "match_expression", "custom:Name=='admin-ep2'"),
					testAccCheckAciCloudEndpointSelectorIdEqual(&cloud_endpoint_selector_default, &cloud_endpoint_selector_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccCloudEndpointSelectorConfigUpdatedName(rName, rName, rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccCloudEndpointSelectorRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudEndpointSelectorConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorExists(resourceName, &cloud_endpoint_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_epg_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudepg-%s", rNameUpdated, rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciCloudEndpointSelectorIdNotEqual(&cloud_endpoint_selector_default, &cloud_endpoint_selector_updated),
				),
			},
			{
				Config: CreateAccCloudEndpointSelectorConfig(rName, rName, rName, rName),
			},
			{
				Config: CreateAccCloudEndpointSelectorConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorExists(resourceName, &cloud_endpoint_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_epg_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudepg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciCloudEndpointSelectorIdNotEqual(&cloud_endpoint_selector_default, &cloud_endpoint_selector_updated),
				),
			},
		},
	})
}

func TestAccAciCloudEndpointSelector_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudEndpointSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudEndpointSelectorConfig(rName, rName, rName, rName),
			},
			{
				Config:      CreateAccCloudEndpointSelectorWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccCloudEndpointSelectorUpdatedAttr(rName, rName, rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudEndpointSelectorUpdatedAttr(rName, rName, rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudEndpointSelectorUpdatedAttr(rName, rName, rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudEndpointSelectorUpdatedAttr(rName, rName, rName, rName, "match_expression", acctest.RandStringFromCharSet(513, "abcdefghijklmnopqrstuvwxyz")),
				ExpectError: regexp.MustCompile(`failed validation for value ''`),
			},
			{
				Config:      CreateAccCloudEndpointSelectorUpdatedAttr(rName, rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCloudEndpointSelectorConfig(rName, rName, rName, rName),
			},
		},
	})
}

func TestAccAciCloudEndpointSelector_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudEndpointSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudEndpointSelectorConfigMultiple(rName, rName, rName, rName),
			},
		},
	})
}

func testAccCheckAciCloudEndpointSelectorExists(name string, cloud_endpoint_selector *models.CloudEndpointSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Endpoint Selector %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Endpoint Selector dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_endpoint_selectorFound := models.CloudEndpointSelectorFromContainer(cont)
		if cloud_endpoint_selectorFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Endpoint Selector %s not found", rs.Primary.ID)
		}
		*cloud_endpoint_selector = *cloud_endpoint_selectorFound
		return nil
	}
}

func testAccCheckAciCloudEndpointSelectorDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing cloud_endpoint_selector destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_endpoint_selector" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_endpoint_selector := models.CloudEndpointSelectorFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Endpoint Selector %s Still exists", cloud_endpoint_selector.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudEndpointSelectorIdEqual(m1, m2 *models.CloudEndpointSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("cloud_endpoint_selector DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciCloudEndpointSelectorIdNotEqual(m1, m2 *models.CloudEndpointSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("cloud_endpoint_selector DNs are equal")
		}
		return nil
	}
}

func CreateCloudEndpointSelectorWithoutRequired(fvTenantName, cloudAppName, cloudEPgName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_endpoint_selector creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	`
	switch attrName {
	case "cloud_epg_dn":
		rBlock += `
	resource "aci_cloud_endpoint_selector" "test" {
	#	cloud_epg_dn  = aci_cloud_epg.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, cloudAppName, cloudEPgName, rName)
}

func CreateAccCloudEndpointSelectorConfigWithRequiredParams(prName, rName string) string {
	fmt.Printf("=== STEP  testing cloud_endpoint_selector with parent resource name %s and name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = "%s"
	}
	`, prName, prName, prName, rName)
	return resource
}
func CreateAccCloudEndpointSelectorConfigUpdatedName(fvTenantName, cloudAppName, cloudEPgName, rName string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selector creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = "%s"
	}
	`, fvTenantName, cloudAppName, cloudEPgName, rName)
	return resource
}

func CreateAccCloudEndpointSelectorConfig(fvTenantName, cloudAppName, cloudEPgName, rName string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selector creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = "%s"
	}
	`, fvTenantName, cloudAppName, cloudEPgName, rName)
	return resource
}

func CreateAccCloudEndpointSelectorConfigMultiple(fvTenantName, cloudAppName, cloudEPgName, rName string) string {
	fmt.Println("=== STEP  testing multiple cloud_endpoint_selector creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, cloudAppName, cloudEPgName, rName)
	return resource
}

func CreateAccCloudEndpointSelectorWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing cloud_endpoint_selector creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccCloudEndpointSelectorConfigWithOptionalValues(fvTenantName, cloudAppName, cloudEPgName, rName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_endpoint_selector creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = "${aci_cloud_epg.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_endpoint_selector"
		match_expression = "custom:Name=='admin-ep2'"
		
	}
	`, fvTenantName, cloudAppName, cloudEPgName, rName)

	return resource
}

func CreateAccCloudEndpointSelectorRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing cloud_endpoint_selector updation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_cloud_endpoint_selector" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_cloud_endpoint_selector"
	}
	`)

	return resource
}

func CreateAccCloudEndpointSelectorUpdatedAttr(fvTenantName, cloudAppName, cloudEPgName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing cloud_endpoint_selector attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, cloudAppName, cloudEPgName, rName, attribute, value)
	return resource
}
