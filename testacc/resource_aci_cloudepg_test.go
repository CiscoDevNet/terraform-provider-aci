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

func TestAccAciCloudEPg_Basic(t *testing.T) {
	var cloud_epg_default models.CloudEPg
	var cloud_epg_updated models.CloudEPg
	resourceName := "aci_cloud_epg.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudEPgWithoutRequired(rName, rName, rName, "cloud_applicationcontainer_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudEPgWithoutRequired(rName, rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudEPgConfig(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_default),
					resource.TestCheckResourceAttr(resourceName, "cloud_applicationcontainer_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtleastOne"),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "exclude"),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", ""),
				),
			},
			{
				Config: CreateAccCloudEPgConfigWithOptionalValues(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_applicationcontainer_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_cloud_epg"),
					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "match_t", "All"),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "include"),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", "0"),
					testAccCheckAciCloudEPgIdEqual(&cloud_epg_default, &cloud_epg_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccCloudEPgConfigUpdatedName(rName, rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccCloudEPgRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudEPgConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_applicationcontainer_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s", rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciCloudEPgIdNotEqual(&cloud_epg_default, &cloud_epg_updated),
				),
			},
			{
				Config: CreateAccCloudEPgConfig(rName, rName, rName),
			},
			{
				Config: CreateAccCloudEPgConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_applicationcontainer_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciCloudEPgIdNotEqual(&cloud_epg_default, &cloud_epg_updated),
				),
			},
		},
	})
}

func TestAccAciCloudEPg_Update(t *testing.T) {
	var cloud_epg_default models.CloudEPg
	var cloud_epg_updated models.CloudEPg
	resourceName := "aci_cloud_epg.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudEPgConfig(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_default),
				),
			},

			{
				Config: CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "match_t", "AtmostOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtmostOne"),
					testAccCheckAciCloudEPgIdEqual(&cloud_epg_default, &cloud_epg_updated),
				),
			},
			{
				Config: CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "match_t", "None"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "None"),
					testAccCheckAciCloudEPgIdEqual(&cloud_epg_default, &cloud_epg_updated),
				),
			},
			{
				Config: CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciCloudEPgIdEqual(&cloud_epg_default, &cloud_epg_updated),
				),
			},
			{
				Config: CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciCloudEPgIdEqual(&cloud_epg_default, &cloud_epg_updated),
				),
			},
			{
				Config: CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciCloudEPgIdEqual(&cloud_epg_default, &cloud_epg_updated),
				),
			},
			{
				Config: CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciCloudEPgIdEqual(&cloud_epg_default, &cloud_epg_updated),
				),
			},
			{
				Config: CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciCloudEPgIdEqual(&cloud_epg_default, &cloud_epg_updated),
				),
			},
			{
				Config: CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "exception_tag", "512"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", "512"),
					testAccCheckAciCloudEPgIdEqual(&cloud_epg_default, &cloud_epg_updated),
				),
			},
			{
				Config: CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "exception_tag", "256"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists(resourceName, &cloud_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", "256"),
					testAccCheckAciCloudEPgIdEqual(&cloud_epg_default, &cloud_epg_updated),
				),
			},
			{
				Config: CreateAccCloudEPgConfig(rName, rName, rName),
			},
		},
	})
}

func TestAccAciCloudEPg_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudEPgConfig(rName, rName, rName),
			},
			{
				Config:      CreateAccCloudEPgWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "flood_on_encap", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "match_t", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "prio", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccCloudEPgUpdatedAttr(rName, rName, rName, "pref_gr_memb", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccCloudEPgUpdatedAttr(rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCloudEPgConfig(rName, rName, rName),
			},
		},
	})
}

func TestAccAciCloudEPg_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudEPgConfigMultiple(rName, rName, rName),
			},
		},
	})
}

func testAccCheckAciCloudEPgExists(name string, cloud_epg *models.CloudEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud EPg %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud EPg dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_epgFound := models.CloudEPgFromContainer(cont)
		if cloud_epgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud EPg %s not found", rs.Primary.ID)
		}
		*cloud_epg = *cloud_epgFound
		return nil
	}
}

func testAccCheckAciCloudEPgDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing cloud_epg destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_epg" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_epg := models.CloudEPgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud EPg %s Still exists", cloud_epg.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudEPgIdEqual(m1, m2 *models.CloudEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("cloud_epg DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciCloudEPgIdNotEqual(m1, m2 *models.CloudEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("cloud_epg DNs are equal")
		}
		return nil
	}
}

func CreateCloudEPgWithoutRequired(fvTenantName, cloudAppName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_epg creation without ", attrName)
	rBlock := `

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	`
	switch attrName {
	case "cloud_applicationcontainer_dn":
		rBlock += `
	resource "aci_cloud_epg" "test" {
	#	cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, cloudAppName, rName)
}

func CreateAccCloudEPgConfigWithRequiredParams(prName, rName string) string {
	fmt.Printf("=== STEP  testing cloud_epg creation with parent resource name %s and resource name %s\n", prName, rName)
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}
	`, prName, prName, rName)
	return resource
}
func CreateAccCloudEPgConfigUpdatedName(fvTenantName, cloudAppName, rName string) string {
	fmt.Println("=== STEP  testing cloud_epg creation with invalid name = ", rName)
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}
	`, fvTenantName, cloudAppName, rName)
	return resource
}

func CreateAccCloudEPgConfig(fvTenantName, cloudAppName, rName string) string {
	fmt.Println("=== STEP  testing cloud_epg creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}
	`, fvTenantName, cloudAppName, rName)
	return resource
}

func CreateAccCloudEPgConfigMultiple(fvTenantName, cloudAppName, rName string) string {
	fmt.Println("=== STEP  testing multiple cloud_epg creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, cloudAppName, rName)
	return resource
}

func CreateAccCloudEPgWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing cloud_epg creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccCloudEPgConfigWithOptionalValues(fvTenantName, cloudAppName, rName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_epg creation with optional parameters")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = "${aci_cloud_applicationcontainer.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_epg"
		flood_on_encap = "enabled"
		match_t = "All"
		pref_gr_memb = "include"
		prio = "level1"
		exception_tag = "0"
	}
	`, fvTenantName, cloudAppName, rName)

	return resource
}

func CreateAccCloudEPgRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing cloud_epg updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_cloud_epg" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_epg"
		flood_on_encap = "enabled"
		match_t = "All"
		pref_gr_memb = "include"
		prio = "level1"
		exception_tag = "0"
	}
	`)

	return resource
}

func CreateAccCloudEPgUpdatedAttr(fvTenantName, cloudAppName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing cloud_epg attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, cloudAppName, rName, attribute, value)
	return resource
}
