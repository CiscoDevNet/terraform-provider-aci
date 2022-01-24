package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciOSPFTimersDataSource_Basic(t *testing.T) {
	resourceName := "aci_ospf_timers.test"
	dataSourceName := "data.aci_ospf_timers.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOSPFTimersDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateOSPFTimersDSWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateOSPFTimersDSWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccOSPFTimersConfigDataSource(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bw_ref", resourceName, "bw_ref"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl", resourceName, "ctrl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "dist", resourceName, "dist"),
					resource.TestCheckResourceAttrPair(dataSourceName, "gr_ctrl", resourceName, "gr_ctrl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "lsa_arrival_intvl", resourceName, "lsa_arrival_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "lsa_gp_pacing_intvl", resourceName, "lsa_gp_pacing_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "lsa_hold_intvl", resourceName, "lsa_hold_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "lsa_max_intvl", resourceName, "lsa_max_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "lsa_start_intvl", resourceName, "lsa_start_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_ecmp", resourceName, "max_ecmp"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_lsa_action", resourceName, "max_lsa_action"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_lsa_num", resourceName, "max_lsa_num"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_lsa_reset_intvl", resourceName, "max_lsa_reset_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_lsa_sleep_cnt", resourceName, "max_lsa_sleep_cnt"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_lsa_sleep_intvl", resourceName, "max_lsa_sleep_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_lsa_thresh", resourceName, "max_lsa_thresh"),
					resource.TestCheckResourceAttrPair(dataSourceName, "spf_hold_intvl", resourceName, "spf_hold_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "spf_init_intvl", resourceName, "spf_init_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "spf_max_intvl", resourceName, "spf_max_intvl"),
				),
			},
			{
				Config:      CreateAccOSPFTimersDataSourceUpdate(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccOSPFTimersDSWithInvalidParentDn(fvTenantName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccOSPFTimersDataSourceUpdatedResource(fvTenantName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccOSPFTimersConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing ospf_timers Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_ospf_timers.test.name
		depends_on = [ aci_ospf_timers.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateOSPFTimersDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ospf_timers Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_ospf_timers" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
		depends_on = [ aci_ospf_timers.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
		depends_on = [ aci_ospf_timers.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccOSPFTimersDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing ospf_timers Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_ospf_timers" "test" {
		tenant_dn  = "${aci_tenant.test.id}_invalid"
		name  = aci_ospf_timers.test.name
		depends_on = [ aci_ospf_timers.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccOSPFTimersDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing ospf_timers Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_ospf_timers.test.name
		%s = "%s"
		depends_on = [ aci_ospf_timers.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccOSPFTimersDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing ospf_timers Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_ospf_timers.test.name
		depends_on = [ aci_ospf_timers.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
