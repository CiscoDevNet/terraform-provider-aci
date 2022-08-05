package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciMatchRuleDataSource_Basic(t *testing.T) {
	resourceName := "aci_match_rule.test"
	dataSourceName := "data.aci_match_rule.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMatchRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateMatchRuleDSWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMatchRuleDSWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMatchRuleConfigDataSource(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccMatchRuleDataSourceUpdate(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccMatchRuleDSWithInvalidParentDn(fvTenantName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccMatchRuleDataSourceUpdatedResource(fvTenantName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccMatchRuleConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing match_rule Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_match_rule" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_match_rule.test.name
		depends_on = [ aci_match_rule.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateMatchRuleDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing match_rule Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_match_rule" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = aci_match_rule.test.name
		depends_on = [ aci_match_rule.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_match_rule" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = aci_match_rule.test.name
		depends_on = [ aci_match_rule.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccMatchRuleDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing match_rule Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_match_rule" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_match_rule.test.name}_invalid"
		depends_on = [ aci_match_rule.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccMatchRuleDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing match_rule Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_match_rule" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_match_rule.test.name
		%s = "%s"
		depends_on = [ aci_match_rule.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccMatchRuleDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing match_rule Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_match_rule" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_match_rule.test.name
		depends_on = [ aci_match_rule.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
