package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciDestinationOfRedirectedTraffic_Basic(t *testing.T) {
	var destination_of_redirected_traffic_default models.Destinationofredirectedtraffic
	var destination_of_redirected_traffic_updated models.Destinationofredirectedtraffic
	resourceName := "aci_destination_of_redirected_traffic.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	ip, _ := acctest.RandIpAddress("10.0.0.0/16")
	ipUpdated, _ := acctest.RandIpAddress("10.0.0.0/16")
	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsSvcRedirectPolName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDestinationOfRedirectedTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateDestinationOfRedirectedTrafficWithoutRequired(fvTenantName, vnsSvcRedirectPolName, ip, "service_redirect_policy_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateDestinationOfRedirectedTrafficWithoutRequired(fvTenantName, vnsSvcRedirectPolName, ip, "ip"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccDestinationOfRedirectedTrafficConfig(fvTenantName, vnsSvcRedirectPolName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDestinationOfRedirectedTrafficExists(resourceName, &destination_of_redirected_traffic_default),
					resource.TestCheckResourceAttr(resourceName, "service_redirect_policy_dn", fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", fvTenantName, vnsSvcRedirectPolName)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "dest_name", ""),
					resource.TestCheckResourceAttr(resourceName, "ip2", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "mac", "12:25:56:98:45:74"),
					resource.TestCheckResourceAttr(resourceName, "pod_id", "1"),
				),
			},
			{
				Config: CreateAccDestinationOfRedirectedTrafficConfigWithOptionalValues(fvTenantName, vnsSvcRedirectPolName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDestinationOfRedirectedTrafficExists(resourceName, &destination_of_redirected_traffic_updated),
					resource.TestCheckResourceAttr(resourceName, "service_redirect_policy_dn", fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", fvTenantName, vnsSvcRedirectPolName)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_destination_of_redirected_traffic"),

					resource.TestCheckResourceAttr(resourceName, "dest_name", "dest"),

					resource.TestCheckResourceAttr(resourceName, "ip2", "1.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "mac", "12:25:56:98:45:74"),
					resource.TestCheckResourceAttr(resourceName, "pod_id", "2"),

					testAccCheckAciDestinationOfRedirectedTrafficIdEqual(&destination_of_redirected_traffic_default, &destination_of_redirected_traffic_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccDestinationOfRedirectedTrafficWithInavalidIP(rName, rName, ip),
				ExpectError: regexp.MustCompile(`unknown property value (.)+`),
			},

			{
				Config:      CreateAccDestinationOfRedirectedTrafficRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccDestinationOfRedirectedTrafficConfigWithRequiredParams(rName, rNameUpdated, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDestinationOfRedirectedTrafficExists(resourceName, &destination_of_redirected_traffic_updated),
					resource.TestCheckResourceAttr(resourceName, "service_redirect_policy_dn", fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", rName, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					testAccCheckAciDestinationOfRedirectedTrafficIdNotEqual(&destination_of_redirected_traffic_default, &destination_of_redirected_traffic_updated),
				),
			},
			{
				Config: CreateAccDestinationOfRedirectedTrafficConfig(fvTenantName, vnsSvcRedirectPolName, ip),
			},
			{
				Config: CreateAccDestinationOfRedirectedTrafficConfigWithRequiredParams(rName, rName, ipUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDestinationOfRedirectedTrafficExists(resourceName, &destination_of_redirected_traffic_updated),
					resource.TestCheckResourceAttr(resourceName, "service_redirect_policy_dn", fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "ip", ipUpdated),
					testAccCheckAciDestinationOfRedirectedTrafficIdNotEqual(&destination_of_redirected_traffic_default, &destination_of_redirected_traffic_updated),
				),
			},
		},
	})
}

func TestAccAciDestinationOfRedirectedTraffic_Update(t *testing.T) {
	var destination_of_redirected_traffic_default models.Destinationofredirectedtraffic
	var destination_of_redirected_traffic_updated models.Destinationofredirectedtraffic
	resourceName := "aci_destination_of_redirected_traffic.test"

	ip, _ := acctest.RandIpAddress("10.0.0.0/16")
	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsSvcRedirectPolName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDestinationOfRedirectedTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDestinationOfRedirectedTrafficConfig(fvTenantName, vnsSvcRedirectPolName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDestinationOfRedirectedTrafficExists(resourceName, &destination_of_redirected_traffic_default),
				),
			},
			{
				Config: CreateAccDestinationOfRedirectedTrafficUpdatedAttr(fvTenantName, vnsSvcRedirectPolName, ip, "pod_id", "255"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDestinationOfRedirectedTrafficExists(resourceName, &destination_of_redirected_traffic_updated),
					resource.TestCheckResourceAttr(resourceName, "pod_id", "255"),
					testAccCheckAciDestinationOfRedirectedTrafficIdEqual(&destination_of_redirected_traffic_default, &destination_of_redirected_traffic_updated),
				),
			},
			{
				Config: CreateAccDestinationOfRedirectedTrafficUpdatedAttr(fvTenantName, vnsSvcRedirectPolName, ip, "pod_id", "127"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDestinationOfRedirectedTrafficExists(resourceName, &destination_of_redirected_traffic_updated),
					resource.TestCheckResourceAttr(resourceName, "pod_id", "127"),
					testAccCheckAciDestinationOfRedirectedTrafficIdEqual(&destination_of_redirected_traffic_default, &destination_of_redirected_traffic_updated),
				),
			},

			{
				Config: CreateAccDestinationOfRedirectedTrafficConfig(fvTenantName, vnsSvcRedirectPolName, ip),
			},
		},
	})
}

func TestAccAciDestinationOfRedirectedTraffic_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	ip, _ := acctest.RandIpAddress("10.0.0.0/16")
	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsSvcRedirectPolName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDestinationOfRedirectedTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDestinationOfRedirectedTrafficConfig(fvTenantName, vnsSvcRedirectPolName, ip),
			},
			{
				Config:      CreateAccDestinationOfRedirectedTrafficWithInValidParentDn(rName, ip),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccDestinationOfRedirectedTrafficUpdatedAttr(fvTenantName, vnsSvcRedirectPolName, ip, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDestinationOfRedirectedTrafficUpdatedAttr(fvTenantName, vnsSvcRedirectPolName, ip, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDestinationOfRedirectedTrafficUpdatedAttr(fvTenantName, vnsSvcRedirectPolName, ip, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccDestinationOfRedirectedTrafficUpdatedAttr(fvTenantName, vnsSvcRedirectPolName, ip, "ip2", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccDestinationOfRedirectedTrafficUpdatedAttr(fvTenantName, vnsSvcRedirectPolName, ip, "pod_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccDestinationOfRedirectedTrafficUpdatedAttr(fvTenantName, vnsSvcRedirectPolName, ip, "pod_id", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccDestinationOfRedirectedTrafficUpdatedAttr(fvTenantName, vnsSvcRedirectPolName, ip, "pod_id", "256"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccDestinationOfRedirectedTrafficUpdatedAttr(fvTenantName, vnsSvcRedirectPolName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccDestinationOfRedirectedTrafficConfig(fvTenantName, vnsSvcRedirectPolName, ip),
			},
		},
	})
}

func TestAccAciDestinationOfRedirectedTraffic_MultipleCreateDelete(t *testing.T) {
	ip, _ := acctest.RandIpAddress("10.0.0.0/16")
	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsSvcRedirectPolName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDestinationOfRedirectedTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDestinationOfRedirectedTrafficConfigMultiple(fvTenantName, vnsSvcRedirectPolName, ip[:len(ip)-1]),
			},
		},
	})
}

func testAccCheckAciDestinationOfRedirectedTrafficExists(name string, destination_of_redirected_traffic *models.Destinationofredirectedtraffic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Destination Of Redirected Traffic %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Destination Of Redirected Traffic dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		destination_of_redirected_trafficFound := models.DestinationofredirectedtrafficFromContainer(cont)
		if destination_of_redirected_trafficFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Destination Of Redirected Traffic %s not found", rs.Primary.ID)
		}
		*destination_of_redirected_traffic = *destination_of_redirected_trafficFound
		return nil
	}
}

func testAccCheckAciDestinationOfRedirectedTrafficDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing destination_of_redirected_traffic destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_destination_of_redirected_traffic" {
			cont, err := client.Get(rs.Primary.ID)
			destination_of_redirected_traffic := models.DestinationofredirectedtrafficFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Destination Of Redirected Traffic %s Still exists", destination_of_redirected_traffic.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciDestinationOfRedirectedTrafficIdEqual(m1, m2 *models.Destinationofredirectedtraffic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("destination_of_redirected_traffic DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciDestinationOfRedirectedTrafficIdNotEqual(m1, m2 *models.Destinationofredirectedtraffic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("destination_of_redirected_traffic DNs are equal")
		}
		return nil
	}
}

func CreateDestinationOfRedirectedTrafficWithoutRequired(fvTenantName, vnsSvcRedirectPolName, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing destination_of_redirected_traffic creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_service_redirect_policy" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	`
	switch attrName {
	case "service_redirect_policy_dn":
		rBlock += `
	resource "aci_destination_of_redirected_traffic" "test" {
	#	service_redirect_policy_dn  = aci_service_redirect_policy.test.id
		ip  = "%s"
		mac = "12:25:56:98:45:74"
	}
		`
	case "ip":
		rBlock += `
	resource "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = aci_service_redirect_policy.test.id
	#	ip  = "%s"
		mac = "12:25:56:98:45:74"
	}
		`
	case "mac":
		rBlock += `
	resource "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = aci_service_redirect_policy.test.id
		ip  = "%s"
	#	mac = "12:25:56:98:45:74"
	}
		`

	}
	return fmt.Sprintf(rBlock, fvTenantName, vnsSvcRedirectPolName, ip)
}

func CreateAccDestinationOfRedirectedTrafficConfigWithRequiredParams(fvTenantName, vnsSvcRedirectPolName, ip string) string {
	fmt.Println("=== STEP  testing destination_of_redirected_traffic creation with updated naming arguments")
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
	`, fvTenantName, vnsSvcRedirectPolName, ip)
	return resource
}

func CreateAccDestinationOfRedirectedTrafficConfig(fvTenantName, vnsSvcRedirectPolName, ip string) string {
	fmt.Println("=== STEP  testing destination_of_redirected_traffic creation with required arguments only")
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
	`, fvTenantName, vnsSvcRedirectPolName, ip)
	return resource
}

func CreateAccDestinationOfRedirectedTrafficWithInavalidIP(fvTenantName, vnsSvcRedirectPolName, ip string) string {
	fmt.Println("=== STEP  testing destination_of_redirected_traffic creation with Invalid IP")
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
		ip  = "%s000"
		mac = "12:25:56:98:45:74"
	}
	`, fvTenantName, vnsSvcRedirectPolName, ip)
	return resource
}

func CreateAccDestinationOfRedirectedTrafficConfigMultiple(fvTenantName, vnsSvcRedirectPolName, ip string) string {
	fmt.Println("=== STEP  testing multiple destination_of_redirected_traffic creation with required arguments only")
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
		ip  = "%s${count.index}"
		mac = "12:25:56:98:45:7${count.index}"
		count = 5
	}
	`, fvTenantName, vnsSvcRedirectPolName, ip)
	return resource
}

func CreateAccDestinationOfRedirectedTrafficWithInValidParentDn(rName, ip string) string {
	fmt.Println("=== STEP  Negative Case: testing destination_of_redirected_traffic creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = aci_tenant.test.id
		ip  = "%s"
		mac = "12:25:56:98:45:74"	
	}
	`, rName, ip)
	return resource
}

func CreateAccDestinationOfRedirectedTrafficConfigWithOptionalValues(fvTenantName, vnsSvcRedirectPolName, ip string) string {
	fmt.Println("=== STEP  Basic: testing destination_of_redirected_traffic creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_destination_of_redirected_traffic" "test" {
		service_redirect_policy_dn  = "${aci_service_redirect_policy.test.id}"
		ip  = "%s"
		mac = "12:25:56:98:45:74"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_destination_of_redirected_traffic"
		dest_name = "dest"
		ip2 = "1.0.0.1"
		pod_id = "2"
		
	}
	`, fvTenantName, vnsSvcRedirectPolName, ip)

	return resource
}

func CreateAccDestinationOfRedirectedTrafficRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing destination_of_redirected_traffic updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_destination_of_redirected_traffic" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_destination_of_redirected_traffic"
		dest_name = ""
		ip2 = ""
		mac = ""
		pod_id = "2"
		
	}
	`)

	return resource
}

func CreateAccDestinationOfRedirectedTrafficUpdatedAttr(fvTenantName, vnsSvcRedirectPolName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing destination_of_redirected_traffic attribute: %s = %s \n", attribute, value)
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
	`, fvTenantName, vnsSvcRedirectPolName, ip, attribute, value)
	return resource
}
