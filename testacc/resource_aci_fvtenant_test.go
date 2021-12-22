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

func TestAccAciTenant_Basic(t *testing.T) {
	var tenant_default models.Tenant
	var tenant_updated models.Tenant
	resourceName := "aci_tenant.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccTenantWithoutName(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTenantConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTenantExists(resourceName, &tenant_default),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					// resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_tenant_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: CreateAccTenantConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTenantExists(resourceName, &tenant_updated),
					resource.TestCheckResourceAttr(resourceName, "description", "from terraform"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_ap"),
					// resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_tenant_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_tn_deny_rule.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciTenantIdEqual(&tenant_default, &tenant_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccTenantConfigUpdatedName(longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of tn-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccTenantConfigWithName(rOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTenantExists(resourceName, &tenant_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rOther),
					testAccCheckAciTenantIdNotEqual(&tenant_default, &tenant_updated),
				),
			},
			{
				Config:      CreateAccTenantConfigUpdateWithoutName(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTenantConfig(rName),
			},
		},
	})
}

func TestAccAciTenant_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTenantConfig(rName),
			},
			{
				Config:      CreateAccTenantUpdatedAttr(rName, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTenantUpdatedAttr(rName, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTenantUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTenantUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccTenantConfig(rName),
			},
		},
	})
}

func TestAccAciTenant_reltionalParameters(t *testing.T) {
	var tenant_default models.Tenant
	var tenant_rel1 models.Tenant
	var tenant_rel2 models.Tenant
	resourceName := "aci_tenant.test"
	rName := makeTestVariable(acctest.RandString(5))
	rsRelName1 := acctest.RandString(5)
	rsRelName2 := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTenantConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTenantExists(resourceName, &tenant_default),
					// resource.TestCheckResourceAttr(resourceName, "relation_tenant_rs_ap_mon_pol", ""),       // checking value of relation_fv_rs_ap_mon_pol parameter for given configuration
				),
			},
			{
				Config: CreateAccTenanttUpdatedTenantIntial(rName, rsRelName1, rsRelName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTenantExists(resourceName, &tenant_rel1),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_tenant_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rsRelName1, rsRelName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_tn_deny_rule.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_tn_deny_rule.*", fmt.Sprintf("uni/tn-%s/flt-%s", rsRelName1, rsRelName2)),
					testAccCheckAciTenantIdEqual(&tenant_default, &tenant_rel1),
				),
			},
			{
				Config: CreateAccTenantUpdatedTenantFinal(rName, rsRelName2, rsRelName1, rsRelName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTenantExists(resourceName, &tenant_rel2),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_tenant_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rsRelName2, rsRelName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_tn_deny_rule.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_tn_deny_rule.*", fmt.Sprintf("uni/tn-%s/flt-%s", rsRelName2, rsRelName1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_tn_deny_rule.*", fmt.Sprintf("uni/tn-%s/flt-%s", rsRelName2, rsRelName2)),
					testAccCheckAciTenantIdEqual(&tenant_default, &tenant_rel2),
				),
			},
			{
				Config: CreateAccTenantConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_tenant_mon_pol", ""),
				),
			},
		},
	})
}

func TestAccAciTenant_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTenantsConfig(rName),
			},
		},
	})
}

func testAccCheckAciTenantExists(name string, tenant *models.Tenant) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Tenant %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Tenant dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tenantFound := models.TenantFromContainer(cont)
		if tenantFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Tenant %s not found", rs.Primary.ID)
		}
		*tenant = *tenantFound
		return nil
	}
}

func testAccCheckAciTenantDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing tenant destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_tenant" {
			cont, err := client.Get(rs.Primary.ID)
			tenant := models.TenantFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Tenant %s Still exists", tenant.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciTenantIdEqual(tn1, tn2 *models.Tenant) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if tn1.DistinguishedName != tn2.DistinguishedName {
			return fmt.Errorf("Tenant DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciTenantIdNotEqual(tn1, tn2 *models.Tenant) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if tn1.DistinguishedName == tn2.DistinguishedName {
			return fmt.Errorf("Tenant DNs are equal")
		}
		return nil
	}
}

func CreateAccTenantWithoutName() string {
	fmt.Println("=== STEP  Basic: testing tenant creation without giving Name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
	}
	`)
	return resource
}
func CreateAccTenantConfigUpdateWithoutName() string {
	fmt.Println("=== STEP  Basic: testing tenant update without giving Name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		annotation = "tag"
		description = "from terraform"
		name_alias = "test_ap"
	}
	`)
	return resource
}
func CreateAccTenantConfigWithName(rOther string) string {
	fmt.Printf("=== STEP  Basic: testing tenant creation with tenant name %s \n", rOther)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	`, rOther)
	return resource
}

func CreateAccTenantConfig(rName string) string {
	fmt.Println("=== STEP testing tenant creation with name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	`, rName)
	return resource
}

func CreateAccTenantConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing tenant creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
		annotation = "tag"
		description = "from terraform"
		name_alias = "test_ap"
	}
	`, rName)
	return resource
}

func CreateAccTenantConfigUpdatedName(longrName string) string {
	fmt.Println("=== STEP  Basic: testing tenant creation with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	`, longrName)
	return resource
}

func CreateAccTenantUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
func CreateAccTenantsConfig(rName string) string {
	fmt.Println("=== STEP  creating multiple tenants")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_tenant" "test1"{
		name = "%s"
	}

	resource "aci_tenant" "test2"{
		name = "%s"
	}

	`, rName+"1", rName+"2", rName+"3")
	return resource
}
func CreateAccTenanttUpdatedTenantIntial(rName, tenantToMP, tenantToFilter string) string {
	fmt.Println("=== STEP  Relation Parameters: testing tenant creation with initial relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test1" {
		name = "%s"
	}

	resource "aci_monitoring_policy" "test" {
		tenant_dn = aci_tenant.test1.id
		name = "%s"
	}

	resource "aci_filter" "test" {
		tenant_dn   = aci_tenant.test1.id
		name = "%s"
	}

	resource "aci_tenant" "test" {
		name = "%s"
		relation_fv_rs_tenant_mon_pol = aci_monitoring_policy.test.id
		relation_fv_rs_tn_deny_rule = [aci_filter.test.id]
	}

	`, tenantToMP, tenantToMP, tenantToFilter, rName)
	return resource
}
func CreateAccTenantUpdatedTenantFinal(rName, tenantToMP, tenantToFilter1, tenantToFilter2 string) string {
	fmt.Println("=== STEP  Relation Parameters: testing tenant creation with final relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test1" {
		name = "%s"
	}

	resource "aci_monitoring_policy" "test" {
		tenant_dn = aci_tenant.test1.id
		name = "%s"
	}

	resource "aci_filter" "test" {
		tenant_dn   = aci_tenant.test1.id
		name = "%s"
	}

	resource "aci_filter" "test1" {
		tenant_dn   = aci_tenant.test1.id
		name = "%s"
	}

	resource "aci_tenant" "test" {
		name = "%s"
		relation_fv_rs_tenant_mon_pol = aci_monitoring_policy.test.id
		relation_fv_rs_tn_deny_rule = [aci_filter.test.id, aci_filter.test1.id]
	}

	`, tenantToMP, tenantToMP, tenantToFilter1, tenantToFilter2, rName)
	return resource
}
