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

func TestAccAciCloudExternalEPg_Basic(t *testing.T) {
	var cloud_external_epg_default models.CloudExternalEPg
	var cloud_external_epg_updated models.CloudExternalEPg
	resourceName := "aci_cloud_external_epg.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudExternalEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudExternalEPgWithoutRequired(rName, rName, rName, "cloud_applicationcontainer_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudExternalEPgWithoutRequired(rName, rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudExternalEPgConfig(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_default),
					resource.TestCheckResourceAttr(resourceName, "cloud_applicationcontainer_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", ""),
					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtleastOne"),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "exclude"),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "route_reachability", "inter-site"),
				),
			},
			{
				Config: CreateAccCloudExternalEPgConfigWithOptionalValues(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_applicationcontainer_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_cloud_external_epg"),
					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "match_t", "All"),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "include"),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", "0"),
					resource.TestCheckResourceAttr(resourceName, "route_reachability", "inter-site"),
					testAccCheckAciCloudExternalEPgIdEqual(&cloud_external_epg_default, &cloud_external_epg_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccCloudExternalEPgConfigUpdatedName(rName, rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccCloudExternalEPgRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudExternalEPgConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_applicationcontainer_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s", rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciCloudExternalEPgIdNotEqual(&cloud_external_epg_default, &cloud_external_epg_updated),
				),
			},
			{
				Config: CreateAccCloudExternalEPgConfig(rName, rName, rName),
			},
			{
				Config: CreateAccCloudExternalEPgConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "cloud_applicationcontainer_dn", fmt.Sprintf("uni/tn-%s/cloudapp-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciCloudExternalEPgIdNotEqual(&cloud_external_epg_default, &cloud_external_epg_updated),
				),
			},
		},
	})
}

func TestAccAciCloudExternalEPg_Update(t *testing.T) {
	var cloud_external_epg_default models.CloudExternalEPg
	var cloud_external_epg_updated models.CloudExternalEPg
	resourceName := "aci_cloud_external_epg.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudExternalEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudExternalEPgConfig(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_default),
				),
			},
			{
				Config: CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "exception_tag", "512"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", "512"),
					testAccCheckAciCloudExternalEPgIdEqual(&cloud_external_epg_default, &cloud_external_epg_updated),
				),
			},
			{
				Config: CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "exception_tag", "256"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", "256"),
					testAccCheckAciCloudExternalEPgIdEqual(&cloud_external_epg_default, &cloud_external_epg_updated),
				),
			},
			{
				Config: CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "match_t", "AtmostOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtmostOne"),
					testAccCheckAciCloudExternalEPgIdEqual(&cloud_external_epg_default, &cloud_external_epg_updated),
				),
			},
			{
				Config: CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "match_t", "None"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "None"),
					testAccCheckAciCloudExternalEPgIdEqual(&cloud_external_epg_default, &cloud_external_epg_updated),
				),
			},
			{
				Config: CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciCloudExternalEPgIdEqual(&cloud_external_epg_default, &cloud_external_epg_updated),
				),
			},
			{
				Config: CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciCloudExternalEPgIdEqual(&cloud_external_epg_default, &cloud_external_epg_updated),
				),
			},
			{
				Config: CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciCloudExternalEPgIdEqual(&cloud_external_epg_default, &cloud_external_epg_updated),
				),
			},
			{
				Config: CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciCloudExternalEPgIdEqual(&cloud_external_epg_default, &cloud_external_epg_updated),
				),
			},
			{
				Config: CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciCloudExternalEPgIdEqual(&cloud_external_epg_default, &cloud_external_epg_updated),
				),
			},
			{
				Config: CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rOther, "route_reachability", "internet"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "route_reachability", "internet"),
				),
			},
			{
				Config: CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "route_reachability", "site-ext"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "route_reachability", "site-ext"),
				),
			},
			{
				Config: CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rOther, "route_reachability", "unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "route_reachability", "unspecified"),
				),
			},
			{
				Config: CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "route_reachability", "inter-site-ext"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists(resourceName, &cloud_external_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "route_reachability", "inter-site-ext"),
				),
			},
			{
				Config: CreateAccCloudExternalEPgConfig(rName, rName, rName),
			},
		},
	})
}

func TestAccAciCloudExternalEPg_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudExternalEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudExternalEPgConfig(rName, rName, rName),
			},
			{
				Config:      CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "route_reachability", "internet"),
				ExpectError: regexp.MustCompile(`Create-only and naming props cannot be modified after creation`),
			},
			{
				Config:      CreateAccCloudExternalEPgWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "flood_on_encap", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "match_t", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "pref_gr_memb", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "prio", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, "route_reachability", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccCloudExternalEPgUpdatedAttr(rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCloudExternalEPgConfig(rName, rName, rName),
			},
		},
	})
}

func TestAccAciCloudExternalEPg_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudExternalEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCloudExternalEPgConfigMultiple(rName, rName, rName),
			},
		},
	})
}

func testAccCheckAciCloudExternalEPgExists(name string, cloud_external_epg *models.CloudExternalEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud External EPg %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud External EPg dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_external_epgFound := models.CloudExternalEPgFromContainer(cont)
		if cloud_external_epgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud External EPg %s not found", rs.Primary.ID)
		}
		*cloud_external_epg = *cloud_external_epgFound
		return nil
	}
}

func testAccCheckAciCloudExternalEPgDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing cloud_external_epg destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_external_epg" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_external_epg := models.CloudExternalEPgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud External EPg %s Still exists", cloud_external_epg.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudExternalEPgIdEqual(m1, m2 *models.CloudExternalEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("cloud_external_epg DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciCloudExternalEPgIdNotEqual(m1, m2 *models.CloudExternalEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("cloud_external_epg DNs are equal")
		}
		return nil
	}
}

func CreateCloudExternalEPgWithoutRequired(fvTenantName, cloudAppName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_external_epg creation without ", attrName)
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
	resource "aci_cloud_external_epg" "test" {
	#	cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, cloudAppName, rName)
}

func CreateAccCloudExternalEPgConfigWithRequiredParams(prName, rName string) string {
	fmt.Printf("=== STEP  testing cloud_external_epg creation with parent resource name %s and resource name %s\n", prName, rName)
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}
	`, prName, prName, rName)
	return resource
}
func CreateAccCloudExternalEPgConfigUpdatedName(fvTenantName, cloudAppName, rName string) string {
	fmt.Println("=== STEP  testing cloud_external_epg creation with invalid name = ", rName)
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}
	`, fvTenantName, cloudAppName, rName)
	return resource
}

func CreateAccCloudExternalEPgConfig(fvTenantName, cloudAppName, rName string) string {
	fmt.Println("=== STEP  testing cloud_external_epg creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}
	`, fvTenantName, cloudAppName, rName)
	return resource
}

func CreateAccCloudExternalEPgConfigMultiple(fvTenantName, cloudAppName, rName string) string {
	fmt.Println("=== STEP  testing multiple cloud_external_epg creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, cloudAppName, rName)
	return resource
}

func CreateAccCloudExternalEPgWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing cloud_external_epg creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccCloudExternalEPgConfigWithOptionalValues(fvTenantName, cloudAppName, rName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_external_epg creation with optional parameters")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = "${aci_cloud_applicationcontainer.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_external_epg"
		flood_on_encap = "enabled"
		match_t = "All"
		pref_gr_memb = "include"
		prio = "level1"
		route_reachability = "inter-site"
		exception_tag = "0"
	}
	`, fvTenantName, cloudAppName, rName)

	return resource
}

func CreateAccCloudExternalEPgRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing cloud_external_epg updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_cloud_external_epg" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cloud_external_epg"
		flood_on_encap = "enabled"
		match_t = "All"
		pref_gr_memb = "include"
		prio = "level1"
		route_reachability = "inter-site-ext"

	}
	`)

	return resource
}

func CreateAccCloudExternalEPgUpdatedAttr(fvTenantName, cloudAppName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing cloud_external_epg attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, cloudAppName, rName, attribute, value)
	return resource
}
