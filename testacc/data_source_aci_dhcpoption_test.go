package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciDHCPOptionDataSource_Basic(t *testing.T) {
	dataSourceName := "data.aci_dhcp_option.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      CreateDHCPOptionDSWithoutRequired(rName, rName, rName, "dhcp_option_policy_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateDHCPOptionDSWithoutRequired(rName, rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccDHCPOptionConfigDataSource(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_option_policy_dn", fmt.Sprintf("uni/tn-%s/dhcpoptpol-%s", rName, rName)),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "annotation", "test_annotation"),
					resource.TestCheckResourceAttr(dataSourceName, "data", "test_data"),
					resource.TestCheckResourceAttr(dataSourceName, "name_alias", "test_name_alias"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_option_id", "1"),
				),
			},
			{
				Config:      CreateAccDHCPOptionDataSourceUpdate(rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccDHCPOptionDSWithInvalidName(rName, rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccDHCPOptionDataSourceUpdatedResource(rName, rName, rName, "name_alias", "updated_name_alias"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name_alias", "updated_name_alias"),
				),
			},
		},
	})
}

func CreateAccDHCPOptionConfigDataSource(fvTenantName, dhcpOptionPolName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_option Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "name" {
		name = "%s"
	  }
	  
	  resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.name.id
		name  = "%s"
	  
		dhcp_option {
		  name  = "%s"
		  annotation  = "test_annotation"
		  data  = "test_data"
		  dhcp_option_id  = "1"
		  name_alias  = "test_name_alias"
		}
	  }
	  
	  data "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = aci_dhcp_option_policy.test.dhcp_option.0.name
	  }
	`, fvTenantName, dhcpOptionPolName, rName)
	return resource
}

func CreateDHCPOptionDSWithoutRequired(fvTenantName, dhcpOptionPolName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing dhcp_option Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
		dhcp_option {
			name  = "%s"
		}
	}
	
	`
	switch attrName {
	case "dhcp_option_policy_dn":
		rBlock += `
	data "aci_dhcp_option" "test" {
	#	dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = aci_dhcp_option_policy.test.dhcp_option.0.name
	}
		`
	case "name":
		rBlock += `
	data "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
	#	name  = aci_dhcp_option_policy.test.dhcp_option.0.name
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, dhcpOptionPolName, rName)
}

func CreateAccDHCPOptionDSWithInvalidName(fvTenantName, dhcpOptionPolName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_option Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "name" {
		name = "%s"
	  }
	  
	  resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.name.id
		name  = "%s"
	  
		dhcp_option {
		  name  = "%s"
		}
	  }

	data "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = "${aci_dhcp_option_policy.test.dhcp_option.0.name}_invalid"
	}
	`, fvTenantName, dhcpOptionPolName, rName)
	return resource
}

func CreateAccDHCPOptionDataSourceUpdate(fvTenantName, dhcpOptionPolName, rName, key, value string) string {
	fmt.Println("=== STEP  testing dhcp_option Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "name" {
		name = "%s"
	  }
	  
	  resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.name.id
		name  = "%s"
	  
		dhcp_option {
		  name  = "%s"
		}
	  }
	  
	  data "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = aci_dhcp_option_policy.test.dhcp_option.0.name
		%s = "%s"
	  }
	`, fvTenantName, dhcpOptionPolName, rName, key, value)
	return resource
}

func CreateAccDHCPOptionDataSourceUpdatedResource(fvTenantName, dhcpOptionPolName, rName, key, value string) string {
	fmt.Println("=== STEP  testing dhcp_option Data Source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "name" {
		name = "%s"
	  }
	  
	  resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.name.id
		name  = "%s"
	  
		dhcp_option {
		  name  = "%s"
		  %s = "%s"
		}
	  }
	  
	  data "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = aci_dhcp_option_policy.test.dhcp_option.0.name
	  }
	`, fvTenantName, dhcpOptionPolName, rName, key, value)
	return resource
}
