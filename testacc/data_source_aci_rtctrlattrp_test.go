package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciActionRuleProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_action_rule_profile.test"
	dataSourceName := "data.aci_action_rule_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciActionRuleProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateActionRuleProfileDSWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateActionRuleProfileDSWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccActionRuleProfileConfigDataSource(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccActionRuleProfileDataSourceUpdate(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccActionRuleProfileDSWithInvalidParentDn(fvTenantName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccActionRuleProfileDataSourceUpdatedResource(fvTenantName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccActionRuleProfileConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing action_rule_profile Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_action_rule_profile.test.name
		depends_on = [ aci_action_rule_profile.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateActionRuleProfileDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing action_rule_profile Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_action_rule_profile" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = aci_action_rule_profile.test.name
		depends_on = [ aci_action_rule_profile.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = aci_action_rule_profile.test.name
		depends_on = [ aci_action_rule_profile.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccActionRuleProfileDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing action_rule_profile Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_action_rule_profile.test.name}_invalid"
		depends_on = [ aci_action_rule_profile.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccActionRuleProfileDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing action_rule_profile Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_action_rule_profile.test.name
		%s = "%s"
		depends_on = [ aci_action_rule_profile.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccActionRuleProfileDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing action_rule_profile Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_action_rule_profile.test.name
		depends_on = [ aci_action_rule_profile.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
