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

func TestAccAciCloudEndpointSelectorforExternalEPgs_Basic(t *testing.T) {
	var cloud_endpoint_selectorfor_external_epgs_default models.CloudEndpointSelectorforExternalEPgs
	var cloud_endpoint_selectorfor_external_epgs_updated models.CloudEndpointSelectorforExternalEPgs
	resourceName := "aci_cloud_endpoint_selectorfor_external_epgs.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	subnet, _ := acctest.RandIpAddress("10.1.0.0/16")
	subnet = fmt.Sprintf("%s/16", subnet)
	subnetUpdated, _ := acctest.RandIpAddress("10.2.0.0/17")
	subnetUpdated = fmt.Sprintf("%s/17", subnetUpdated)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudEndpointSelectorforExternalEPgsDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudEndpointSelectorforExternalEPgsWithoutRequired(rName, rName, rName, subnet, "cloud_external_epg_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudEndpointSelectorforExternalEPgsWithoutRequired(rName, rName, rName, subnet, "subnet"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudEndpointSelectorforExternalEPgsWithoutRequired(rName, rName, rName, subnet, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudEndpointSelectorforExternalEPgsConfig(rName, rName, rName, subnet, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorforExternalEPgsExists(resourceName, &cloud_endpoint_selectorfor_external_epgs_default),
					resource.TestCheckResourceAttr(resourceName, "cloud_external_epg_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudextepg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "subnet", subnet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "is_shared", "yes"),
					resource.TestCheckResourceAttr(resourceName, "match_expression", ""),
				),
			},
			{
				Config: CreateAccCloudEndpointSelectorforExternalEPgsConfigWithOptionalValues(rName, rName, rName, subnet, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorforExternalEPgsExists(resourceName, &cloud_endpoint_selectorfor_external_epgs_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_external_epg_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudextepg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "subnet", subnet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_cloud_endpoint_selectorfor_external_epgs"),
					resource.TestCheckResourceAttr(resourceName, "is_shared", "no"),
					resource.TestCheckResourceAttr(resourceName, "match_expression", "custom:tag=='provbaz'"),
					testAccCheckAciCloudEndpointSelectorforExternalEPgsIdEqual(&cloud_endpoint_selectorfor_external_epgs_default, &cloud_endpoint_selectorfor_external_epgs_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccCloudEndpointSelectorforExternalEPgsWithInavalidIP(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+`),
			},
			{
				Config:      CreateAccCloudEndpointSelectorforExternalEPgsWithLongerName(rName, subnet, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudEndpointSelectorforExternalEPgsRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudEndpointSelectorforExternalEPgsConfigWithRequiredParams(rNameUpdated, subnet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorforExternalEPgsExists(resourceName, &cloud_endpoint_selectorfor_external_epgs_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_external_epg_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudextepg-%s", rNameUpdated, rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "subnet", subnet),
					testAccCheckAciCloudEndpointSelectorforExternalEPgsIdNotEqual(&cloud_endpoint_selectorfor_external_epgs_default, &cloud_endpoint_selectorfor_external_epgs_updated),
				),
			},
			{
				Config: CreateAccCloudEndpointSelectorforExternalEPgsConfig(rName, rName, rName, subnet, rName),
			},
			{
				Config: CreateAccCloudEndpointSelectorforExternalEPgsConfigWithRequiredParams(rName, subnetUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEndpointSelectorforExternalEPgsExists(resourceName, &cloud_endpoint_selectorfor_external_epgs_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_external_epg_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudextepg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "subnet", subnetUpdated),
					testAccCheckAciCloudEndpointSelectorforExternalEPgsIdNotEqual(&cloud_endpoint_selectorfor_external_epgs_default, &cloud_endpoint_selectorfor_external_epgs_updated),
				),
			},
		},
	})
}

func TestAccAciCloudEndpointSelectorforExternalEPgs_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	subnet, _ := acctest.RandIpAddress("10.3.0.0/18")
	subnet = fmt.Sprintf("%s/16", subnet)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudEndpointSelectorforExternalEPgsDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudEndpointSelectorforExternalEPgsConfig(rName, rName, rName, subnet, rName),
			},
			{
				Config:      CreateAccCloudEndpointSelectorforExternalEPgsWithInValidParentDn(rName, subnet),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccCloudEndpointSelectorforExternalEPgsUpdatedAttr(rName, rName, rName, subnet, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudEndpointSelectorforExternalEPgsUpdatedAttr(rName, rName, rName, subnet, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudEndpointSelectorforExternalEPgsUpdatedAttr(rName, rName, rName, subnet, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccCloudEndpointSelectorforExternalEPgsUpdatedAttr(rName, rName, rName, subnet, rName, "is_shared", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccCloudEndpointSelectorforExternalEPgsUpdatedAttr(rName, rName, rName, subnet, rName, "match_expression", acctest.RandStringFromCharSet(513, "abcdefghijklmnopqrstuvwxyz")),
				ExpectError: regexp.MustCompile(`failed validation for value ''`),
			},
			{
				Config:      CreateAccCloudEndpointSelectorforExternalEPgsUpdatedAttr(rName, rName, rName, subnet, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCloudEndpointSelectorforExternalEPgsConfig(rName, rName, rName, subnet, rName),
			},
		},
	})
}

func TestAccAciCloudEndpointSelectorforExternalEPgs_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	subnet1, _ := acctest.RandIpAddress("10.5.0.0/20")
	subnet1 = fmt.Sprintf("%s/20", subnet1)
	subnet2, _ := acctest.RandIpAddress("10.6.0.0/21")
	subnet2 = fmt.Sprintf("%s/21", subnet2)
	subnet3, _ := acctest.RandIpAddress("10.7.0.0/22")
	subnet3 = fmt.Sprintf("%s/22", subnet3)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudEndpointSelectorforExternalEPgsDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudEndpointSelectorforExternalEPgsConfigMultiple(rName, rName, rName, subnet1, subnet2, subnet3, rName),
			},
		},
	})
}

func testAccCheckAciCloudEndpointSelectorforExternalEPgsExists(name string, cloud_endpoint_selectorfor_external_epgs *models.CloudEndpointSelectorforExternalEPgs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Endpoint Selectorfor External EPgs %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Endpoint Selectorfor External EPgs dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_endpoint_selectorfor_external_epgsFound := models.CloudEndpointSelectorforExternalEPgsFromContainer(cont)
		if cloud_endpoint_selectorfor_external_epgsFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Endpoint Selectorfor External EPgs %s not found", rs.Primary.ID)
		}
		*cloud_endpoint_selectorfor_external_epgs = *cloud_endpoint_selectorfor_external_epgsFound
		return nil
	}
}

func testAccCheckAciCloudEndpointSelectorforExternalEPgsDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing cloud_endpoint_selectorfor_external_epgs destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_endpoint_selectorfor_external_epgs" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_endpoint_selectorfor_external_epgs := models.CloudEndpointSelectorforExternalEPgsFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Endpoint Selectorfor External EPgs %s Still exists", cloud_endpoint_selectorfor_external_epgs.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudEndpointSelectorforExternalEPgsIdEqual(m1, m2 *models.CloudEndpointSelectorforExternalEPgs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("cloud_endpoint_selectorfor_external_epgs DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciCloudEndpointSelectorforExternalEPgsIdNotEqual(m1, m2 *models.CloudEndpointSelectorforExternalEPgs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("cloud_endpoint_selectorfor_external_epgs DNs are equal")
		}
		return nil
	}
}

func CreateCloudEndpointSelectorforExternalEPgsWithoutRequired(fvTenantName, cloudAppName, cloudExtEPgName, subnet, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_endpoint_selectorfor_external_epgs creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	`
	switch attrName {
	case "cloud_external_epg_dn":
		rBlock += `
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
	#	cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = "%s"
		name = "%s"
	}
		`
	case "subnet":
		rBlock += `
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
	#	subnet  = "%s"
		name = "%s"
	}
	
		`
	case "name":
		rBlock += `
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = "%s"
	#	name = "%s"
	}`
	}
	return fmt.Sprintf(rBlock, fvTenantName, cloudAppName, cloudExtEPgName, subnet)
}

func CreateAccCloudEndpointSelectorforExternalEPgsConfigWithRequiredParams(rName, subnet string) string {
	fmt.Printf("=== STEP  testing cloud_endpoint_selectorfor_external_epgs creation with parent resources name %s and subnet %s\n", rName, subnet)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = "%s"
		name = "%s"
	}
	`, rName, rName, rName, subnet, rName)
	return resource
}
func CreateAccCloudEndpointSelectorforExternalEPgsConfigUpdatedName(fvTenantName, cloudAppName, cloudExtEPgName, subnet, name string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selectorfor_external_epgs creation with invalid ip")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
	subnet  = "%s_invalid"
		name = "%s"
		
	}
	`, fvTenantName, cloudAppName, cloudExtEPgName, subnet, name)
	return resource
}

func CreateAccCloudEndpointSelectorforExternalEPgsWithLongerName(rName, subnet, longName string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selectorfor_external_epgs creation with long name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = "%s"
		name = "%s"
	}
	`, rName, rName, rName, subnet, longName)
	return resource
}

func CreateAccCloudEndpointSelectorforExternalEPgsWithInavalidIP(rName string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selectorfor_external_epgs creation with invalid subnet")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = "%s"
		name = "%s"
	}
	`, rName, rName, rName, rName, rName)
	return resource
}

func CreateAccCloudEndpointSelectorforExternalEPgsConfig(fvTenantName, cloudAppName, cloudExtEPgName, subnet, name string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selectorfor_external_epgs creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = "%s"
		name = "%s"
	}
	`, fvTenantName, cloudAppName, cloudExtEPgName, subnet, name)
	return resource
}

func CreateAccCloudEndpointSelectorforExternalEPgsConfigMultiple(fvTenantName, cloudAppName, cloudExtEPgName, subnet1, subnet2, subnet3, name string) string {
	fmt.Println("=== STEP  testing multiple cloud_endpoint_selectorfor_external_epgs creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet = "%s"
		name = "%s"
	}

	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test1" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet = "%s"
		name = "%s"
	}

	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test2" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet = "%s"
		name = "%s"
	}
	`, fvTenantName, cloudAppName, cloudExtEPgName, subnet1, name, subnet2, name, subnet3, name)
	return resource
}

func CreateAccCloudEndpointSelectorforExternalEPgsWithInValidParentDn(rName, subnet string) string {
	fmt.Println("=== STEP  Negative Case: testing cloud_endpoint_selectorfor_external_epgs creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_tenant.test.id
		subnet  = "%s"	
		name = "%s"
	}
	`, rName, subnet, rName)
	return resource
}

func CreateAccCloudEndpointSelectorforExternalEPgsConfigWithOptionalValues(fvTenantName, cloudAppName, cloudExtEPgName, subnet, name string) string {
	fmt.Println("=== STEP  Basic: testing cloud_endpoint_selectorfor_external_epgs creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = "${aci_cloud_external_epg.test.id}"
		subnet  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_endpoint_selectorfor_external_epgs"
		is_shared = "no"
		match_expression = "custom:tag=='provbaz'"
		name = "%s"
	}
	`, fvTenantName, cloudAppName, cloudExtEPgName, subnet, name)

	return resource
}

func CreateAccCloudEndpointSelectorforExternalEPgsRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing cloud_endpoint_selectorfor_external_epgs updation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_endpoint_selectorfor_external_epgs"
		is_shared = "no"
		match_expression = ""
	}
	`)

	return resource
}

func CreateAccCloudEndpointSelectorforExternalEPgsUpdatedAttr(fvTenantName, cloudAppName, cloudExtEPgName, subnet, name, attribute, value string) string {
	fmt.Printf("=== STEP  testing cloud_endpoint_selectorfor_external_epgs attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = "%s"
		name = "%s"
		%s = "%s"
	}
	`, fvTenantName, cloudAppName, cloudExtEPgName, subnet, name, attribute, value)
	return resource
}
