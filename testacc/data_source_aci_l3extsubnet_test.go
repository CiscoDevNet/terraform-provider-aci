package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3ExtSubnetDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3_ext_subnet.test"
	dataSourceName := "data.aci_l3_ext_subnet.test"
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.0.5.0/20")
	ip = fmt.Sprintf("%s/20", ip)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3ExtSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3ExtSubnetDSWithoutRequired(rName, rName, rName, ip, "external_network_instance_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3ExtSubnetDSWithoutRequired(rName, rName, rName, ip, "ip"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3ExtSubnetDSConfig(rName, rName, rName, ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "external_network_instance_profile_dn", resourceName, "external_network_instance_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ip", resourceName, "ip"),
					resource.TestCheckResourceAttrPair(dataSourceName, "aggregate", resourceName, "aggregate"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "scope.#", resourceName, "scope.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "scope.0", resourceName, "scope.0"),
				),
			},
			{
				Config:      CreateAccL3ExtSubnetDSUpdateRandomAttr(rName, rName, rName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccL3ExtSubnetDSWithInvalidParentDn(rName, rName, rName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`Object may not exists`),
			},
			{
				Config: CreateAccL3ExtSubnetDSUpdate(rName, rName, rName, ip, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateL3ExtSubnetDSWithoutRequired(fvTenantName, l3extOutName, l3extInstPName, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3_ext_subnet data source creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_external_network_instance_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}

	resource "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_external_network_instance_profile.test.id
		ip  = "%s"
	}
	
	`
	switch attrName {
	case "external_network_instance_profile_dn":
		rBlock += `
	data "aci_l3_ext_subnet" "test" {
	#	external_network_instance_profile_dn  = aci_l3_ext_subnet.test.external_network_instance_profile_dn
		ip  = aci_l3_ext_subnet.test.ip
	}
		`
	case "ip":
		rBlock += `
	data "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_l3_ext_subnet.test.external_network_instance_profile_dn
	#	ip  = aci_l3_ext_subnet.test.ip
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extInstPName, ip)
}

func CreateAccL3ExtSubnetDSConfig(fvTenantName, l3extOutName, l3extInstPName, ip string) string {
	fmt.Println("=== STEP  testing l3_ext_subnet data source creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_external_network_instance_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_external_network_instance_profile.test.id
		ip  = "%s"
	}

	data "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_l3_ext_subnet.test.external_network_instance_profile_dn
		ip  = aci_l3_ext_subnet.test.ip
	}
	`, fvTenantName, l3extOutName, l3extInstPName, ip)
	return resource
}

func CreateAccL3ExtSubnetDSUpdateRandomAttr(fvTenantName, l3extOutName, l3extInstPName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3_ext_subnet data source creation with %s = %s\n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_external_network_instance_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_external_network_instance_profile.test.id
		ip  = "%s"
	}

	data "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_l3_ext_subnet.test.external_network_instance_profile_dn
		ip  = aci_l3_ext_subnet.test.ip
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extInstPName, ip, attribute, value)
	return resource
}

func CreateAccL3ExtSubnetDSWithInvalidParentDn(fvTenantName, l3extOutName, l3extInstPName, ip, attribute, value string) string {
	fmt.Println("=== STEP  testing l3_ext_subnet data source creation with invalid external_network_instance_profile_dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_external_network_instance_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_external_network_instance_profile.test.id
		ip  = "%s"
	}

	data "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = "${aci_l3_ext_subnet.test.external_network_instance_profile_dn}abc" 
		ip  = aci_l3_ext_subnet.test.ip
	}
	`, fvTenantName, l3extOutName, l3extInstPName, ip)
	return resource
}

func CreateAccL3ExtSubnetDSUpdate(fvTenantName, l3extOutName, l3extInstPName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3_ext_subnet data source creation with %s = %s\n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_external_network_instance_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_external_network_instance_profile.test.id
		ip  = "%s"
		%s = "%s"
	}

	data "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_l3_ext_subnet.test.external_network_instance_profile_dn 
		ip  = aci_l3_ext_subnet.test.ip
	}
	`, fvTenantName, l3extOutName, l3extInstPName, ip, attribute, value)
	return resource
}
