package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciMatchRouteDestinationRuleDataSource_Basic(t *testing.T) {
	resourceName := "aci_match_route_destination_rule.test"
	dataSourceName := "data.aci_match_route_destination_rule.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	ip, _ := acctest.RandIpAddress("10.4.0.0/16")
	fvTenantName := makeTestVariable(acctest.RandString(5))
	rtctrlSubjPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMatchRouteDestinationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateMatchRouteDestinationRuleDSWithoutRequired(fvTenantName, rtctrlSubjPName, ip, "match_rule_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMatchRouteDestinationRuleDSWithoutRequired(fvTenantName, rtctrlSubjPName, ip, "ip"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMatchRouteDestinationRuleConfigDataSource(fvTenantName, rtctrlSubjPName, ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "match_rule_dn", resourceName, "match_rule_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ip", resourceName, "ip"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "aggregate", resourceName, "aggregate"),
					resource.TestCheckResourceAttrPair(dataSourceName, "greater_than_mask", resourceName, "greater_than_mask"),
					resource.TestCheckResourceAttrPair(dataSourceName, "less_than_mask", resourceName, "less_than_mask"),
				),
			},
			{
				Config:      CreateAccMatchRouteDestinationRuleDataSourceUpdate(fvTenantName, rtctrlSubjPName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccMatchRouteDestinationRuleDSWithInvalidParentDn(fvTenantName, rtctrlSubjPName, ip),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccMatchRouteDestinationRuleDataSourceUpdatedResource(fvTenantName, rtctrlSubjPName, ip, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccMatchRouteDestinationRuleConfigDataSource(fvTenantName, rtctrlSubjPName, ip string) string {
	fmt.Println("=== STEP  testing match_route_destination_rule Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
		ip  = "%s"
	}

	data "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
		ip  = aci_match_route_destination_rule.test.ip
		depends_on = [ aci_match_route_destination_rule.test ]
	}
	`, fvTenantName, rtctrlSubjPName, ip)
	return resource
}

func CreateMatchRouteDestinationRuleDSWithoutRequired(fvTenantName, rtctrlSubjPName, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing match_route_destination_rule Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
		ip  = "%s"
	}
	`
	switch attrName {
	case "match_rule_dn":
		rBlock += `
	data "aci_match_route_destination_rule" "test" {
	#	match_rule_dn  = aci_match_rule.test.id
		ip  = aci_match_route_destination_rule.test.ip
		depends_on = [ aci_match_route_destination_rule.test ]
	}
		`
	case "ip":
		rBlock += `
	data "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
	#	ip  = aci_match_route_destination_rule.test.ip
		depends_on = [ aci_match_route_destination_rule.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rtctrlSubjPName, ip)
}

func CreateAccMatchRouteDestinationRuleDSWithInvalidParentDn(fvTenantName, rtctrlSubjPName, ip string) string {
	fmt.Println("=== STEP  testing match_route_destination_rule Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
		ip  = "%s"
	}

	data "aci_match_route_destination_rule" "test" {
		match_rule_dn  = "${aci_match_rule.test.id}_invalid"
		ip  = "${aci_match_route_destination_rule.test.ip}"
		depends_on = [ aci_match_route_destination_rule.test ]
	}
	`, fvTenantName, rtctrlSubjPName, ip)
	return resource
}

func CreateAccMatchRouteDestinationRuleDataSourceUpdate(fvTenantName, rtctrlSubjPName, ip, key, value string) string {
	fmt.Println("=== STEP  testing match_route_destination_rule Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
		ip  = "%s"
	}

	data "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
		ip  = aci_match_route_destination_rule.test.ip
		%s = "%s"
		depends_on = [ aci_match_route_destination_rule.test ]
	}
	`, fvTenantName, rtctrlSubjPName, ip, key, value)
	return resource
}

func CreateAccMatchRouteDestinationRuleDataSourceUpdatedResource(fvTenantName, rtctrlSubjPName, ip, key, value string) string {
	fmt.Println("=== STEP  testing match_route_destination_rule Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
		ip  = "%s"
		%s = "%s"
	}

	data "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
		ip  = aci_match_route_destination_rule.test.ip
		depends_on = [ aci_match_route_destination_rule.test ]
	}
	`, fvTenantName, rtctrlSubjPName, ip, key, value)
	return resource
}
