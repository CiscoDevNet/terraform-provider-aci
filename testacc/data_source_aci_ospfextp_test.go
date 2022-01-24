package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outOspfExternalPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_ospf_external_policy.test"
	dataSourceName := "data.aci_l3out_ospf_external_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outOspfExternalPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outOspfExternalPolicyDSWithoutRequired(rName, rName, "l3_outside_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "l3_outside_dn", resourceName, "l3_outside_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "area_cost", resourceName, "area_cost"),
					resource.TestCheckResourceAttrPair(dataSourceName, "area_ctrl.#", resourceName, "area_ctrl.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "area_ctrl.0", resourceName, "area_ctrl.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "area_ctrl.1", resourceName, "area_ctrl.1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "area_id", resourceName, "area_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "area_type", resourceName, "area_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "multipod_internal", resourceName, "multipod_internal"),
				),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL3outOspfExternalPolicyDSWithInvalidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccL3outOspfExternalPolicyDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL3outOspfExternalPolicyConfigDataSource(fvTenantName, l3extOutName string) string {
	fmt.Println("=== STEP  testing l3out_ospf_external_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
	}

	data "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		depends_on = [ aci_l3out_ospf_external_policy.test ]
	}
	`, fvTenantName, l3extOutName)
	return resource
}

func CreateL3outOspfExternalPolicyDSWithoutRequired(fvTenantName, l3extOutName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_ospf_external_policy Data Source without ", attrName)
	rBlock := `

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
	}
	`
	switch attrName {
	case "l3_outside_dn":
		rBlock += `
	data "aci_l3out_ospf_external_policy" "test" {
	#	l3_outside_dn  = aci_l3_outside.test.id

		depends_on = [ aci_l3out_ospf_external_policy.test ]
	}
		`

	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName)
}

func CreateAccL3outOspfExternalPolicyDSWithInvalidParentDn(fvTenantName, l3extOutName string) string {
	fmt.Println("=== STEP  testing l3out_ospf_external_policy Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
	}

	data "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = "${aci_l3_outside.test.id}_invalid"
		depends_on = [ aci_l3out_ospf_external_policy.test ]
	}
	`, fvTenantName, l3extOutName)
	return resource
}

func CreateAccL3outOspfExternalPolicyDataSourceUpdate(fvTenantName, l3extOutName, key, value string) string {
	fmt.Println("=== STEP  testing l3out_ospf_external_policy Data Source with random attribute")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
	}

	data "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		%s = "%s"
		depends_on = [ aci_l3out_ospf_external_policy.test ]
	}
	`, fvTenantName, l3extOutName, key, value)
	return resource
}

func CreateAccL3outOspfExternalPolicyDataSourceUpdatedResource(fvTenantName, l3extOutName, key, value string) string {
	fmt.Println("=== STEP  testing l3out_ospf_external_policy Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		%s = "%s"
	}

	data "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		depends_on = [ aci_l3out_ospf_external_policy.test ]
	}
	`, fvTenantName, l3extOutName, key, value)
	return resource
}
