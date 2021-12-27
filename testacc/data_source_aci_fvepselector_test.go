package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciEndpointSecurityGroupSelectorDataSource_Basic(t *testing.T) {
	resourceName := "aci_endpoint_security_group_selector.test"
	dataSourceName := "data.aci_endpoint_security_group_selector.test"
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.30.0.0/17")
	ip = fmt.Sprintf("%s/17", ip)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndpointSecurityGroupSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateEndpointSecurityGroupSelectorDSWithoutRequired(rName, ip, "endpoint_security_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateEndpointSecurityGroupSelectorDSWithoutRequired(rName, ip, "matchExpression"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEndpointSecurityGroupSelectorDSConfig(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "endpoint_security_group_dn", resourceName, "endpoint_security_group_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_expression", resourceName, "match_expression"),
				),
			},
			{
				Config:      CreateAccEndpointSecurityGroupSelectorDataSourceUpdateRandomAttr(rName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupDSWithInvalidEndpointSecurityGroup(rName, ip),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccEndpointSecurityGroupDSWithUpdatedResource(rName, ip, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccEndpointSecurityGroupDSWithUpdatedResource(rName, ip, key, value string) string {
	fmt.Println("=== STEP  testing endpoint_security_group_selector data source creation with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_application_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_endpoint_security_group" "test" {
		name 		= "%s"
		application_profile_dn = aci_application_profile.test.id
	}
	
	resource "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_endpoint_security_group.test.id
		match_expression  = "ip=='%s'"
		%s = "%s"
	}

	data "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_endpoint_security_group_selector.test.endpoint_security_group_dn
		match_expression  = aci_endpoint_security_group_selector.test.match_expression
	}
	`, rName, rName, rName, ip, key, value)
	return resource
}

func CreateAccEndpointSecurityGroupDSWithInvalidEndpointSecurityGroup(rName, ip string) string {
	fmt.Println("=== STEP  testing endpoint_security_group_selector data source creation with invalid endpoint_security_group_dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_application_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_endpoint_security_group" "test" {
		name 		= "%s"
		application_profile_dn = aci_application_profile.test.id
	}
	
	resource "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_endpoint_security_group.test.id
		match_expression  = "ip=='%s'"
	}

	data "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = "${aci_endpoint_security_group_selector.test.endpoint_security_group_dn}xyz"
		match_expression  = aci_endpoint_security_group_selector.test.match_expression
	}
	`, rName, rName, rName, ip)
	return resource
}

func CreateAccEndpointSecurityGroupSelectorDataSourceUpdateRandomAttr(rName, ip, key, value string) string {
	fmt.Println("=== STEP  testing endpoint_security_group_selector data source creation with random parameter")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_application_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_endpoint_security_group" "test" {
		name 		= "%s"
		application_profile_dn = aci_application_profile.test.id
	}
	
	resource "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_endpoint_security_group.test.id
		match_expression  = "ip=='%s'"
	}

	data "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_endpoint_security_group_selector.test.endpoint_security_group_dn
		match_expression  = aci_endpoint_security_group_selector.test.match_expression
		%s = "%s"
	}
	`, rName, rName, rName, ip, key, value)
	return resource
}

func CreateAccEndpointSecurityGroupSelectorDSConfig(rName, ip string) string {
	fmt.Println("=== STEP  testing endpoint_security_group_selector data source creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_application_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_endpoint_security_group" "test" {
		name 		= "%s"
		application_profile_dn = aci_application_profile.test.id
	}
	
	resource "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_endpoint_security_group.test.id
		match_expression  = "ip=='%s'"
	}

	data "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_endpoint_security_group_selector.test.endpoint_security_group_dn
		match_expression  = aci_endpoint_security_group_selector.test.match_expression
	}
	`, rName, rName, rName, ip)
	return resource
}

func CreateEndpointSecurityGroupSelectorDSWithoutRequired(rName, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_security_group_selector data source creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_application_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_endpoint_security_group" "test" {
		name 		= "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	resource "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_endpoint_security_group.test.id
		match_expression  = "ip=='%s'"
	}
	
	`
	switch attrName {
	case "endpoint_security_group_dn":
		rBlock += `
	data "aci_endpoint_security_group_selector" "test" {
	#	endpoint_security_group_dn  = aci_endpoint_security_group_selector.test.endpoint_security_group_dn
		match_expression  = aci_endpoint_security_group_selector.test.match_expression
	}
		`
	case "matchExpression":
		rBlock += `
	data "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_endpoint_security_group_selector.test.endpoint_security_group_dn
	#	match_expression  = aci_endpoint_security_group_selector.test.match_expression
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, rName, rName, ip)
}
