package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciHSRPGroupPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_hsrp_group_policy.test"
	dataSourceName := "data.aci_hsrp_group_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciHSRPGroupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateHSRPGroupPolicyDSWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateHSRPGroupPolicyDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccHSRPGroupPolicyConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl", resourceName, "ctrl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hello_intvl", resourceName, "hello_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hold_intvl", resourceName, "hold_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "preempt_delay_min", resourceName, "preempt_delay_min"),
					resource.TestCheckResourceAttrPair(dataSourceName, "preempt_delay_reload", resourceName, "preempt_delay_reload"),
					resource.TestCheckResourceAttrPair(dataSourceName, "preempt_delay_sync", resourceName, "preempt_delay_sync"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
					resource.TestCheckResourceAttrPair(dataSourceName, "timeout", resourceName, "timeout"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hsrp_group_policy_type", resourceName, "hsrp_group_policy_type"),
				),
			},
			{
				Config:      CreateAccHSRPGroupPolicyDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccHSRPGroupPolicyDSWithInvalidName(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccHSRPGroupPolicyDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccHSRPGroupPolicyConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing hsrp_group_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_hsrp_group_policy.test.name
		depends_on = [ aci_hsrp_group_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateHSRPGroupPolicyDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing hsrp_group_policy Data Source without ", attrName)
	rBlock := `

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_hsrp_group_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
		depends_on = [ aci_hsrp_group_policy.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
		depends_on = [ aci_hsrp_group_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccHSRPGroupPolicyDSWithInvalidName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing hsrp_group_policy Data Source with invalid name")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_hsrp_group_policy.test.name}_invalid"
		depends_on = [ aci_hsrp_group_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccHSRPGroupPolicyDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing hsrp_group_policy Data Source with random attribute")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_hsrp_group_policy.test.name
		%s = "%s"
		depends_on = [ aci_hsrp_group_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccHSRPGroupPolicyDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing hsrp_group_policy Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_hsrp_group_policy.test.name
		depends_on = [ aci_hsrp_group_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
