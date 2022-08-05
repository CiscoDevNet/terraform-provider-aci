package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciOSPFInterfacePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_ospf_interface_policy.test"
	dataSourceName := "data.aci_ospf_interface_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOSPFInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateOSPFInterfacePolicyDSWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateOSPFInterfacePolicyDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccOSPFInterfacePolicyConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cost", resourceName, "cost"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.#", resourceName, "ctrl.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.0", resourceName, "ctrl.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "dead_intvl", resourceName, "dead_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hello_intvl", resourceName, "hello_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "nw_t", resourceName, "nw_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pfx_suppress", resourceName, "pfx_suppress"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
					resource.TestCheckResourceAttrPair(dataSourceName, "rexmit_intvl", resourceName, "rexmit_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "xmit_delay", resourceName, "xmit_delay"),
				),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccOSPFInterfacePolicyDSWithInvalidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccOSPFInterfacePolicyDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccOSPFInterfacePolicyConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing ospf_interface_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_ospf_interface_policy.test.name
		depends_on = [ aci_ospf_interface_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateOSPFInterfacePolicyDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ospf_interface_policy creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_ospf_interface_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
		depends_on = [ aci_ospf_interface_policy.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
		depends_on = [ aci_ospf_interface_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccOSPFInterfacePolicyDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing ospf_interface_policy Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_ospf_interface_policy.test.name}_invalid"
		depends_on = [ aci_ospf_interface_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccOSPFInterfacePolicyDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing ospf_interface_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_ospf_interface_policy.test.name
		%s = "%s"
		depends_on = [ aci_ospf_interface_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccOSPFInterfacePolicyDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing ospf_interface_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_ospf_interface_policy.test.name
		depends_on = [ aci_ospf_interface_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
