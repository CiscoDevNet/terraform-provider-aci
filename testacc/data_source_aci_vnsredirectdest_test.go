package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciDestinationOfRedirectedTrafficDataSource_Basic(t *testing.T) {
	resourceName := "aci_destination_of_redirected_traffic.test"
	dataSourceName := "data.aci_destination_of_redirected_traffic.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	ip, _ := acctest.RandIpAddress("10.0.0.0/16")
	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsSvcRedirectPolName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDestinationOfRedirectedTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateDestinationOfRedirectedTrafficDSWithoutRequired(fvTenantName, vnsSvcRedirectPolName, ip, "service_redirect_policy_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateDestinationOfRedirectedTrafficDSWithoutRequired(fvTenantName, vnsSvcRedirectPolName, ip, "ip"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccDestinationOfRedirectedTrafficConfigDataSource(fvTenantName, vnsSvcRedirectPolName, ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "service_redirect_policy_dn", resourceName, "service_redirect_policy_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ip", resourceName, "ip"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "dest_name", resourceName, "dest_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ip2", resourceName, "ip2"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mac", resourceName, "mac"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pod_id", resourceName, "pod_id"),
				),
			},
			{
				Config:      CreateAccDestinationOfRedirectedTrafficDataSourceUpdate(fvTenantName, vnsSvcRedirectPolName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccDestinationOfRedirectedTrafficDSWithInvalidParentDn(fvTenantName, vnsSvcRedirectPolName, ip),
				ExpectError: regexp.MustCompile(`Invalid RN`),
			},

			{
				Config: CreateAccDestinationOfRedirectedTrafficDataSourceUpdatedResource(fvTenantName, vnsSvcRedirectPolName, ip, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccDestinationOfRedirectedTrafficConfigDataSource(fvTenantName, vnsSvcRedirectPolName, ip string) string {
	fmt.Println("=== STEP  testing destination_of_redirected_traffic Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = aci_service_redirect_policy.test.id
		mac = "12:25:56:98:45:74"
		ip  = "%s"
	}

	data "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = aci_service_redirect_policy.test.id
		ip  = aci_destination_of_redirected_traffic.test.ip
		depends_on = [ aci_destination_of_redirected_traffic.test ]
	}
	`, fvTenantName, vnsSvcRedirectPolName, ip)
	return resource
}

func CreateDestinationOfRedirectedTrafficDSWithoutRequired(fvTenantName, vnsSvcRedirectPolName, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing destination_of_redirected_traffic Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = aci_service_redirect_policy.test.id
		mac = "12:25:56:98:45:74"
		ip  = "%s"
	}
	`
	switch attrName {
	case "service_redirect_policy_dn":
		rBlock += `
	data "aci_destination_of_redirected_traffic" "test" {
	#	service_redirect_policy_dn  = aci_service_redirect_policy.test.id
		ip  = aci_destination_of_redirected_traffic.test.ip
		depends_on = [ aci_destination_of_redirected_traffic.test ]
	}
		`
	case "ip":
		rBlock += `
	data "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = aci_service_redirect_policy.test.id
	#	ip  = aci_destination_of_redirected_traffic.test.ip
		depends_on = [ aci_destination_of_redirected_traffic.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, vnsSvcRedirectPolName, ip)
}

func CreateAccDestinationOfRedirectedTrafficDSWithInvalidParentDn(fvTenantName, vnsSvcRedirectPolName, ip string) string {
	fmt.Println("=== STEP  testing destination_of_redirected_traffic Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = aci_service_redirect_policy.test.id
		mac = "12:25:56:98:45:74"
		ip  = "%s"
	}

	data "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = aci_service_redirect_policy.test.id
		ip  = "${aci_destination_of_redirected_traffic.test.ip}_invalid"
		depends_on = [ aci_destination_of_redirected_traffic.test ]
	}
	`, fvTenantName, vnsSvcRedirectPolName, ip)
	return resource
}

func CreateAccDestinationOfRedirectedTrafficDataSourceUpdate(fvTenantName, vnsSvcRedirectPolName, ip, key, value string) string {
	fmt.Println("=== STEP  testing destination_of_redirected_traffic Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = aci_service_redirect_policy.test.id
		ip  = "%s"
		mac = "12:25:56:98:45:74"
	}

	data "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = aci_service_redirect_policy.test.id
		ip  = aci_destination_of_redirected_traffic.test.ip
		%s = "%s"
		depends_on = [ aci_destination_of_redirected_traffic.test ]
	}
	`, fvTenantName, vnsSvcRedirectPolName, ip, key, value)
	return resource
}

func CreateAccDestinationOfRedirectedTrafficDataSourceUpdatedResource(fvTenantName, vnsSvcRedirectPolName, ip, key, value string) string {
	fmt.Println("=== STEP  testing destination_of_redirected_traffic Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = aci_service_redirect_policy.test.id
		ip  = "%s"
		mac = "12:25:56:98:45:74"
		%s = "%s"
	}

	data "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = aci_service_redirect_policy.test.id
		ip  = aci_destination_of_redirected_traffic.test.ip
		depends_on = [ aci_destination_of_redirected_traffic.test ]
	}
	`, fvTenantName, vnsSvcRedirectPolName, ip, key, value)
	return resource
}
