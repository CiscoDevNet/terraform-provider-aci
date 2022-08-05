package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciBdDHCPLabelDataSource_Basic(t *testing.T) {
	resourceName := "aci_bd_dhcp_label.test"
	dataSourceName := "data.aci_bd_dhcp_label.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	fvBDName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBDDHCPLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBdDHCPLabelDSWithoutRequired(fvTenantName, fvBDName, rName, "bridge_domain_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBdDHCPLabelDSWithoutRequired(fvTenantName, fvBDName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBdDHCPLabelConfigDataSource(fvTenantName, fvBDName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "bridge_domain_dn", resourceName, "bridge_domain_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "owner", resourceName, "owner"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag", resourceName, "tag"),
				),
			},
			{
				Config:      CreateAccBdDHCPLabelDataSourceUpdate(fvTenantName, fvBDName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccBdDHCPLabelDSWithInvalidParentDn(fvTenantName, fvBDName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccBdDHCPLabelDataSourceUpdatedResource(fvTenantName, fvBDName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccBdDHCPLabelConfigDataSource(fvTenantName, fvBDName, rName string) string {
	fmt.Println("=== STEP  testing bd_dhcp_label Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bridge_domain" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = "%s"
	}

	data "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = aci_bd_dhcp_label.test.name
		depends_on = [ aci_bd_dhcp_label.test ]
	}
	`, fvTenantName, fvBDName, rName)
	return resource
}

func CreateBdDHCPLabelDSWithoutRequired(fvTenantName, fvBDName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bd_dhcp_label Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bridge_domain" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "bridge_domain_dn":
		rBlock += `
	data "aci_bd_dhcp_label" "test" {
	#	bridge_domain_dn  = aci_bridge_domain.test.id
		name  = aci_bd_dhcp_label.test.name
		depends_on = [ aci_bd_dhcp_label.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
	#	name  = aci_bd_dhcp_label.test.name
		depends_on = [ aci_bd_dhcp_label.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, fvBDName, rName)
}

func CreateAccBdDHCPLabelDSWithInvalidParentDn(fvTenantName, fvBDName, rName string) string {
	fmt.Println("=== STEP  testing bd_dhcp_label Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bridge_domain" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = "%s"
	}

	data "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = "${aci_bd_dhcp_label.test.name}_invalid"
		depends_on = [ aci_bd_dhcp_label.test ]
	}
	`, fvTenantName, fvBDName, rName)
	return resource
}

func CreateAccBdDHCPLabelDataSourceUpdate(fvTenantName, fvBDName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bd_dhcp_label Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bridge_domain" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = "%s"
	}

	data "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = aci_bd_dhcp_label.test.name
		%s = "%s"
		depends_on = [ aci_bd_dhcp_label.test ]
	}
	`, fvTenantName, fvBDName, rName, key, value)
	return resource
}

func CreateAccBdDHCPLabelDataSourceUpdatedResource(fvTenantName, fvBDName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bd_dhcp_label Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bridge_domain" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_bd_dhcp_label" "test" {
		bridge_domain_dn  = aci_bridge_domain.test.id
		name  = aci_bd_dhcp_label.test.name
		depends_on = [ aci_bd_dhcp_label.test ]
	}
	`, fvTenantName, fvBDName, rName, key, value)
	return resource
}
