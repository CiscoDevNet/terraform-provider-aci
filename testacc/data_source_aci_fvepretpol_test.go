package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciEndPointRetentionPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_end_point_retention_policy.test"
	dataSourceName := "data.aci_end_point_retention_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEndPointRetentionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateEndPointRetentionPolicyDSWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateEndPointRetentionPolicyDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEndPointRetentionPolicyConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bounce_age_intvl", resourceName, "bounce_age_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bounce_trig.#", resourceName, "bounce_trig.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bounce_trig.0", resourceName, "bounce_trig.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hold_intvl", resourceName, "hold_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "local_ep_age_intvl", resourceName, "local_ep_age_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "move_freq", resourceName, "move_freq"),
					resource.TestCheckResourceAttrPair(dataSourceName, "remote_ep_age_intvl", resourceName, "remote_ep_age_intvl"),
				),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccEndPointRetentionPolicyDSWithInvalidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccEndPointRetentionPolicyDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccEndPointRetentionPolicyConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing end_point_retention_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_end_point_retention_policy.test.name
		depends_on = [ aci_end_point_retention_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateEndPointRetentionPolicyDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing end_point_retention_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_end_point_retention_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = aci_end_point_retention_policy.test.name
		depends_on = [ aci_end_point_retention_policy.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = aci_end_point_retention_policy.test.name
		depends_on = [ aci_end_point_retention_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccEndPointRetentionPolicyDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing end_point_retention_policy Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_end_point_retention_policy.test.name}_invalid"
		depends_on = [ aci_end_point_retention_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccEndPointRetentionPolicyDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing end_point_retention_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_end_point_retention_policy.test.name
		%s = "%s"
		depends_on = [ aci_end_point_retention_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccEndPointRetentionPolicyDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing end_point_retention_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_end_point_retention_policy.test.name
		depends_on = [ aci_end_point_retention_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
