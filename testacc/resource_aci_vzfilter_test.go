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

func TestAccAciFilter_Basic(t *testing.T) {
	var filter_default models.Filter
	var filter_updated models.Filter
	resourceName := "aci_filter.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)
	prOther := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccFilterWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccFilterWithoutTenant(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{

				Config: CreateAccFilterConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists(resourceName, &filter_default),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_filt_graph_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_fwd_r_flt_p_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_rev_r_flt_p_att", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
				),
			},
			{
				Config: CreateAccFilterConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists(resourceName, &filter_updated),
					resource.TestCheckResourceAttr(resourceName, "description", "From Terraform"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "alias_filter"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag_filter"),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_filt_graph_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_fwd_r_flt_p_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_rev_r_flt_p_att", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAciFilterIdEqual(&filter_default, &filter_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccFilterRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccFilterConfigUpdatedName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of flt-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccFilterConfigWithParentAndName(rName, rOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists(resourceName, &filter_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rOther),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAciFilterIdNotEqual(&filter_default, &filter_updated),
				),
			},
			{
				Config: CreateAccFilterConfig(rName),
			},
			{
				Config: CreateAccFilterConfigWithParentAndName(prOther, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists(resourceName, &filter_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", prOther)),
					testAccCheckAciFilterIdNotEqual(&filter_default, &filter_updated),
				),
			},
		},
	},
	)
}

func TestAccAciFilter_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFilterConfig(rName),
			},
			{
				Config:      CreateAccFilterWithInValidTenantDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class vzFilter (.)+`),
			},
			{
				Config:      CreateAccFilterUpdatedAttr(rName, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFilterUpdatedAttr(rName, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFilterUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccFilterUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFilterConfig(rName),
			},
		},
	})
}

func TestAccAciFilter_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFiltersConfig(rName),
			},
		},
	})
}

func CreateAccFiltersConfig(rName string) string {
	fmt.Println("=== STEP  creating multiple filter")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_filter" "test1"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_filter" "test2"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_filter" "test3"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	`, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccFilterUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing filter attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		%s = "%s"
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccFilterWithInValidTenantDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing filter creation with invalid tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_vrf" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter" "test" {
		tenant_dn = aci_vrf.test.id
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func testAccCheckAciFilterIdNotEqual(f1, f2 *models.Filter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if f1.DistinguishedName == f2.DistinguishedName {
			return fmt.Errorf("Filter DNs are equal")
		}
		return nil
	}
}

func CreateAccFilterConfigWithParentAndName(prName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing filter creation with tenant name %s name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, prName, rName)
	return resource
}

func CreateAccFilterConfigUpdatedName(rName, longrName string) string {
	fmt.Println("=== STEP  Basic: testing filter creation with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, longrName)
	return resource
}

func CreateAccFilterWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter creation without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
	}
	`, rName)
	return resource
}

func CreateAccFilterWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter creation without giving tenant_dn")
	resource := fmt.Sprintf(`

	resource "aci_filter" "test" {
		name="%s"
	}
	`, rName)
	return resource
}

func CreateAccFilterWithoutFilter(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter creation without creating tenant")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_filter" "test" {
		tenant_dn = 
	}
	`, rName)
	return resource
}

func CreateAccFilterConfig(rName string) string {
	fmt.Println("=== STEP  testing filter creation with required arguments")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccFilterConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test" {
        tenant_dn   = aci_tenant.test.id
        description = "From Terraform"
        name        = "%s"
        annotation  = "tag_filter"
        name_alias  = "alias_filter"
    }
	`, rName, rName)
	return resource
}

func CreateAccFilterRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing filter updation without optional parameters")
	resource := fmt.Sprintln(`

	resource "aci_filter" "test" {
        description = "From Terraform"
        annotation  = "tag"
        name_alias  = "alias_filter"
    }
	`)
	return resource
}

func testAccCheckAciFilterIdEqual(f1, f2 *models.Filter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if f1.DistinguishedName != f2.DistinguishedName {
			return fmt.Errorf("filter DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFilterExists(name string, filter *models.Filter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Filter %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Filter dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		filterFound := models.FilterFromContainer(cont)
		if filterFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Filter %s not found", rs.Primary.ID)
		}
		*filter = *filterFound
		return nil
	}
}

func testAccCheckAciFilterDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing filter destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_filter" {
			cont, err := client.Get(rs.Primary.ID)
			filter := models.FilterFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Filter %s Still exists", filter.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}
