package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciBfdInterfacePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_bfd_interface_policy.test"
	dataSourceName := "data.aci_bfd_interface_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBfdInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBfdInterfacePolicyDSWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBfdInterfacePolicyDSWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBfdInterfacePolicyConfigDataSource(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl", resourceName, "ctrl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "detect_mult", resourceName, "detect_mult"),
					resource.TestCheckResourceAttrPair(dataSourceName, "echo_admin_st", resourceName, "echo_admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "echo_rx_intvl", resourceName, "echo_rx_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "min_rx_intvl", resourceName, "min_rx_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "min_tx_intvl", resourceName, "min_tx_intvl"),
				),
			},
			{
				Config:      CreateAccBfdInterfacePolicyDataSourceUpdate(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyDSWithInvalidParentDn(fvTenantName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccBfdInterfacePolicyDataSourceUpdatedResource(fvTenantName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccBfdInterfacePolicyConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bfd_interface_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bfd_interface_policy.test.name
		depends_on = [ aci_bfd_interface_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateBfdInterfacePolicyDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bfd_interface_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_bfd_interface_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = aci_bfd_interface_policy.test.name
		depends_on = [ aci_bfd_interface_policy.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = aci_bfd_interface_policy.test.name
		depends_on = [ aci_bfd_interface_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBfdInterfacePolicyDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bfd_interface_policy Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bfd_interface_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}_invalid"
		name  = aci_bfd_interface_policy.test.name
		depends_on = [ aci_bfd_interface_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBfdInterfacePolicyDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bfd_interface_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bfd_interface_policy.test.name
		%s = "%s"
		depends_on = [ aci_bfd_interface_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccBfdInterfacePolicyDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bfd_interface_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bfd_interface_policy.test.name
		depends_on = [ aci_bfd_interface_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
