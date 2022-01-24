package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL2OutExtepgDataSource_Basic(t *testing.T) {
	resourceName := "aci_l2out_extepg.test"
	dataSourceName := "data.aci_l2out_extepg.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2outExternalEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL2OutExtepgDSWithoutRequired(rName, rName, rName, "l2_outside_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL2OutExtepgDSWithoutRequired(rName, rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL2OutExtepgConfigDataSource(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "l2_outside_dn", resourceName, "l2_outside_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "flood_on_encap", resourceName, "flood_on_encap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_t", resourceName, "match_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pref_gr_memb", resourceName, "pref_gr_memb"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
					resource.TestCheckResourceAttrPair(dataSourceName, "target_dscp", resourceName, "target_dscp"),
					resource.TestCheckResourceAttrPair(dataSourceName, "exception_tag", resourceName, "exception_tag"),
				),
			},
			{
				Config:      CreateAccL2OutExtepgDataSourceUpdate(rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL2OutExtepgDSWithInvalidParentDn(rName, rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccL2OutExtepgDataSourceUpdatedResource(rName, rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL2OutExtepgConfigDataSource(fvTenantName, l2extOutName, rName string) string {
	fmt.Println("=== STEP  testing l2out_extepg Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
		name  = "%s"
	}

	data "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
		name  = aci_l2out_extepg.test.name
		depends_on = [ aci_l2out_extepg.test ]
	}
	`, fvTenantName, l2extOutName, rName)
	return resource
}

func CreateL2OutExtepgDSWithoutRequired(fvTenantName, l2extOutName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l2out_extepg Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "l2_outside_dn":
		rBlock += `
	data "aci_l2out_extepg" "test" {
	#	l2_outside_dn  = aci_l2_outside.test.id
		name  = aci_l2out_extepg.test.name
		depends_on = [ aci_l2out_extepg.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
	#	name  = aci_l2out_extepg.test.name
		depends_on = [ aci_l2out_extepg.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l2extOutName, rName)
}

func CreateAccL2OutExtepgDSWithInvalidParentDn(fvTenantName, l2extOutName, rName string) string {
	fmt.Println("=== STEP  testing l2out_extepg Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
		name  = "%s"
	}

	data "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
		name  = "${aci_l2out_extepg.test.name}_invalid"
		depends_on = [ aci_l2out_extepg.test ]
	}
	`, fvTenantName, l2extOutName, rName)
	return resource
}

func CreateAccL2OutExtepgDataSourceUpdate(fvTenantName, l2extOutName, rName, key, value string) string {
	fmt.Println("=== STEP  testing l2out_extepg Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
		name  = "%s"
	}

	data "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
		name  = aci_l2out_extepg.test.name
		%s = "%s"
		depends_on = [ aci_l2out_extepg.test ]
	}
	`, fvTenantName, l2extOutName, rName, key, value)
	return resource
}

func CreateAccL2OutExtepgDataSourceUpdatedResource(fvTenantName, l2extOutName, rName, key, value string) string {
	fmt.Println("=== STEP  testing l2out_extepg Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
		name  = aci_l2out_extepg.test.name
		depends_on = [ aci_l2out_extepg.test ]
	}
	`, fvTenantName, l2extOutName, rName, key, value)
	return resource
}
