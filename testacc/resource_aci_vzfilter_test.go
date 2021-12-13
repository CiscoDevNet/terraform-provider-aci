package acctest

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

	resource.Test(t, resource.TestCase{
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
				// step terraform will create application profile with only required arguements i.e. name and tenant_dn
				Config: CreateAccFilterConfig(rName), // configuration to create application profile with required fields only
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists(resourceName, &filter_default), // this function will check whether any resource is exist or not in state file with given resource name
					// now will compare value of all attributes with default for given resource
					resource.TestCheckResourceAttr(resourceName, "description", ""), // no default value for description so comparing with ""
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_filt_graph_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_fwd_r_flt_p_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_rev_r_flt_p_att", ""),   // no default value for name_alias so comparing with ""
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"), // comparing with default value of annotation
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
				),
			},
			{
				Config: CreateAccFilterConfigWithOptionalValues(rName), // configuration to update optional filelds
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists(resourceName, &filter_updated),
					resource.TestCheckResourceAttr(resourceName, "description", "From Terraform"), // comparing description with value which is given in configuration
					resource.TestCheckResourceAttr(resourceName, "name_alias", "alias_filter"),    // comparing name_alias with value which is given in configuration
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag_filter"),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_filt_graph_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_fwd_r_flt_p_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_rev_r_flt_p_att", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)), // comparing prio with value which is given in configuration
					testAccCheckAciFilterIdEqual(&filter_default, &filter_updated),                             // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccFilterRemovingRequiredField(), // configuration to update optional filelds
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccFilterConfigUpdatedName(rName, longrName), // passing invalid name for application profile
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of flt-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccFilterConfigWithParentAndName(rName, rOther), // creating resource with same parent name and different resource name
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists(resourceName, &filter_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rOther),                               // comparing name attribute of applicaiton profile
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)), // comparing tenant_dn attribute of application profile
					testAccCheckAciFilterIdNotEqual(&filter_default, &filter_updated),                          // checking whether id or dn of both resource are different because name changed and terraform need to create another resource
				),
			},
			{
				Config: CreateAccFilterConfig(rName), // creating resource with required parameters only
			},
			{
				Config: CreateAccFilterConfigWithParentAndName(prOther, rName), // creating resource with same name but different parent resource name
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists(resourceName, &filter_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", prOther)),
					testAccCheckAciFilterIdNotEqual(&filter_default, &filter_updated), // checking whether id or dn of both resource are different because tenant_dn changed and terraform need to create another resource
				),
			},
		},
	},
	)
}

func TestAccFilter_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longDescAnnotation := acctest.RandString(129)                                     // creating random string of 129 characters
	longNameAlias := acctest.RandString(64)                                           // creating random string of 64 characters                                              // creating random string of 6 characters
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz") // creating random string of 5 characters (to give as random parameter)
	randomValue := acctest.RandString(5)                                              // creating random string of 5 characters (to give as random value of random parameter)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFilterConfig(rName), // creating application profile with required arguements only
			},
			{
				Config:      CreateAccFilterWithInValidTenantDn(rName),                                       // checking application profile creation with invalid tenant_dn value
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class vzFilter (.)+`), // test step expect error which should be match with defined regex
			},
			{
				Config:      CreateAccFilterUpdatedAttr(rName, "description", longDescAnnotation), // checking application profile creation with invalid description value
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFilterUpdatedAttr(rName, "annotation", longDescAnnotation), // checking application profile creation with invalid annotation value
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFilterUpdatedAttr(rName, "name_alias", longNameAlias), // checking application profile creation with invalid name_alias value
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccFilterUpdatedAttr(rName, randomParameter, randomValue), // checking application profile creation with randomly created parameter and value
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFilterConfig(rName), // creating application profile with required arguements only
			},
		},
	})
}

func TestAccFilter_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
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
	fmt.Println("=== STEP  Basic: testing filter creation without giving name")
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
	resource "aci_filter" "test" {
		name = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFilterConfig(rName string) string {
	fmt.Println("=== STEP  testing filter creation with required arguements")
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
	fmt.Println("=== STEP  Basic: testing filter creation with optional parameters")
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
