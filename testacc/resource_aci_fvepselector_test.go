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

func TestAccAciEndpointSecurityGroupSelector_Basic(t *testing.T) {
	var endpoint_security_group_selector_default models.EndpointSecurityGroupSelector
	var endpoint_security_group_selector_updated models.EndpointSecurityGroupSelector
	resourceName := "aci_endpoint_security_group_selector.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.20.0.0/17")
	ip = fmt.Sprintf("%s/17", ip)
	ipother, _ := acctest.RandIpAddress("10.21.0.0/17")
	ipother = fmt.Sprintf("%s/17", ipother)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndpointSecurityGroupSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateEndpointSecurityGroupSelectorWithoutRequired(rName, ip, "endpoint_security_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateEndpointSecurityGroupSelectorWithoutRequired(rName, ip, "matchExpression"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEndpointSecurityGroupSelectorConfig(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupSelectorExists(resourceName, &endpoint_security_group_selector_default),
					resource.TestCheckResourceAttr(resourceName, "endpoint_security_group_dn", fmt.Sprintf("uni/tn-%s/ap-%s/esg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "match_expression", fmt.Sprintf("ip=='%s'", ip)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "name", ""),
				),
			},
			{
				Config: CreateAccEndpointSecurityGroupSelectorConfigWithOptionalValues(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupSelectorExists(resourceName, &endpoint_security_group_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_endpoint_security_group_selector"),
					resource.TestCheckResourceAttr(resourceName, "endpoint_security_group_dn", fmt.Sprintf("uni/tn-%s/ap-%s/esg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "match_expression", fmt.Sprintf("ip=='%s'", ip)),
					resource.TestCheckResourceAttr(resourceName, "name", "test_endpoint_security_group_selector_name"),
					testAccCheckAciEndpointSecurityGroupSelectorIdEqual(&endpoint_security_group_selector_default, &endpoint_security_group_selector_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccEndpointSecurityGroupSelectorRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupSelectorConfigWithRequiredParams(rName, rName),
				ExpectError: regexp.MustCompile(`Invalid IP Address`),
			},
			{
				Config: CreateAccEndpointSecurityGroupSelectorConfigWithRequiredParams(rOther, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupSelectorExists(resourceName, &endpoint_security_group_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "endpoint_security_group_dn", fmt.Sprintf("uni/tn-%s/ap-%s/esg-%s", rOther, rOther, rOther)),
					resource.TestCheckResourceAttr(resourceName, "match_expression", fmt.Sprintf("ip=='%s'", ip)),
					testAccCheckAciEndpointSecurityGroupSelectorIdNotEqual(&endpoint_security_group_selector_default, &endpoint_security_group_selector_updated),
				),
			},
			{
				Config: CreateAccEndpointSecurityGroupSelectorConfig(rName, ip),
			},
			{
				Config: CreateAccEndpointSecurityGroupSelectorConfigWithRequiredParams(rName, ipother),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupSelectorExists(resourceName, &endpoint_security_group_selector_updated),
					resource.TestCheckResourceAttr(resourceName, "endpoint_security_group_dn", fmt.Sprintf("uni/tn-%s/ap-%s/esg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "match_expression", fmt.Sprintf("ip=='%s'", ipother)),
					testAccCheckAciEndpointSecurityGroupSelectorIdNotEqual(&endpoint_security_group_selector_default, &endpoint_security_group_selector_updated),
				),
			},
		},
	})
}

func TestAccAciEndpointSecurityGroupSelector_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longDescAnnotation := acctest.RandString(129)
	longNameAliasName := acctest.RandString(65)
	ip, _ := acctest.RandIpAddress("10.22.0.0/17")
	ip = fmt.Sprintf("%s/17", ip)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndpointSecurityGroupSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointSecurityGroupSelectorConfig(rName, ip),
			},
			{
				Config:      CreateAccEndpointSecurityGroupSelectorWithInValidSecurityGroup(rName, ip),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class fvEPSelector (.)+`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupSelectorUpdatedAttr(rName, ip, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupSelectorUpdatedAttr(rName, ip, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupSelectorUpdatedAttr(rName, ip, "name_alias", longNameAliasName),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupSelectorUpdatedAttr(rName, ip, "name", longNameAliasName),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupSelectorUpdatedAttr(rName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+ is not expected here.`),
			},
			{
				Config: CreateAccEndpointSecurityGroupSelectorConfig(rName, ip),
			},
		},
	})
}
func TestAccAciEndpointSecurityGroupSelector_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	ip1, _ := acctest.RandIpAddress("10.32.0.0/18")
	ip1 = fmt.Sprintf("%s/18", ip1)
	ip2, _ := acctest.RandIpAddress("10.33.0.0/19")
	ip2 = fmt.Sprintf("%s/19", ip2)
	ip3, _ := acctest.RandIpAddress("10.34.0.0/20")
	ip3 = fmt.Sprintf("%s/20", ip3)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndpointSecurityGroupSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointSecurityGroupSelectorsConfig(rName, ip1, ip2, ip3),
			},
		},
	})
}

func CreateAccEndpointSecurityGroupSelectorsConfig(rName, ip1, ip2, ip3 string) string {
	fmt.Println("=== STEP  testing endpoint_security_group_selector multiple creation")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_application_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_endpoint_security_group" "test" {
		name 		= "%s"
		application_profile_dn = aci_application_profile.test.id
	}
	
	resource "aci_endpoint_security_group_selector" "test1" {
		endpoint_security_group_dn  = aci_endpoint_security_group.test.id
		match_expression  = "ip=='%s'"
	}

	resource "aci_endpoint_security_group_selector" "test2" {
		endpoint_security_group_dn  = aci_endpoint_security_group.test.id
		match_expression  = "ip=='%s'"
	}

	resource "aci_endpoint_security_group_selector" "test3" {
		endpoint_security_group_dn  = aci_endpoint_security_group.test.id
		match_expression  = "ip=='%s'"
	}
	`, rName, rName, rName, ip1, ip2, ip3)
	return resource
}

func testAccCheckAciEndpointSecurityGroupSelectorExists(name string, endpoint_security_group_selector *models.EndpointSecurityGroupSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Endpoint Security Group Selector %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Endpoint Security Group Selector dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		endpoint_security_group_selectorFound := models.EndpointSecurityGroupSelectorFromContainer(cont)
		if endpoint_security_group_selectorFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Endpoint Security Group Selector %s not found", rs.Primary.ID)
		}
		*endpoint_security_group_selector = *endpoint_security_group_selectorFound
		return nil
	}
}

func testAccCheckAciEndpointSecurityGroupSelectorDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing endpoint_security_group_selector destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_endpoint_security_group_selector" {
			cont, err := client.Get(rs.Primary.ID)
			endpoint_security_group_selector := models.EndpointSecurityGroupSelectorFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Endpoint Security Group Selector %s Still exists", endpoint_security_group_selector.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciEndpointSecurityGroupSelectorIdEqual(m1, m2 *models.EndpointSecurityGroupSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("endpoint_security_group_selector DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciEndpointSecurityGroupSelectorIdNotEqual(m1, m2 *models.EndpointSecurityGroupSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("endpoint_security_group_selector DNs are equal")
		}
		return nil
	}
}

func CreateEndpointSecurityGroupSelectorWithoutRequired(rName, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_security_group_selector creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_application_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_endpoint_security_group" "test" {
		name 		= "%s"
		application_profile_dn = aci_application_profile.test.id
	}
	
	`
	switch attrName {
	case "endpoint_security_group_dn":
		rBlock += `
	resource "aci_endpoint_security_group_selector" "test" {
	#	endpoint_security_group_dn  = aci_endpoint_security_group.test.id
		match_expression  = "ip=='%s'"
	}
		`
	case "matchExpression":
		rBlock += `
	resource "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_endpoint_security_group.test.id
	#	match_expression  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, rName, rName, ip)
}

func CreateAccEndpointSecurityGroupSelectorConfigWithRequiredParams(rName, ip string) string {
	fmt.Println("=== STEP  testing endpoint_security_group_selector creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_application_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_endpoint_security_group" "test" {
		name 		= "%s"
		application_profile_dn = aci_application_profile.test.id
	}
	
	resource "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_endpoint_security_group.test.id
		match_expression  = "ip=='%s'"
	}
	`, rName, rName, rName, ip)
	return resource
}

func CreateAccEndpointSecurityGroupSelectorConfig(rName, ip string) string {
	fmt.Println("=== STEP  testing endpoint_security_group_selector creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_application_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_endpoint_security_group" "test" {
		name 		= "%s"
		application_profile_dn = aci_application_profile.test.id
	}
	
	resource "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_endpoint_security_group.test.id
		match_expression  = "ip=='%s'"
	}
	`, rName, rName, rName, ip)
	return resource
}

func CreateAccEndpointSecurityGroupSelectorWithInValidSecurityGroup(rName, ip string) string {
	fmt.Println("=== STEP  Negative Case: testing endpoint_security_group_selector creation with invalid endpoint_security_group_dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	
	resource "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_tenant.test.id
		match_expression  = "ip=='%s'"	
	}
	`, rName, ip)
	return resource
}

func CreateAccEndpointSecurityGroupSelectorConfigWithOptionalValues(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_security_group_selector creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_application_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_endpoint_security_group" "test" {
		name 		= "%s"
		application_profile_dn = aci_application_profile.test.id
	}
	
	resource "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = "${aci_endpoint_security_group.test.id}"
		match_expression  = "ip=='%s'"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_endpoint_security_group_selector"
		name = "test_endpoint_security_group_selector_name"
	}
	`, rName, rName, rName, ip)

	return resource
}

func CreateAccEndpointSecurityGroupSelectorRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing endpoint_security_group_selector creation with optional parameters")
	resource := fmt.Sprintln(`
	resource "aci_endpoint_security_group_selector" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_endpoint_security_group_selector"
		name = "test_endpoint_security_group_selector_name"
	}
	`)

	return resource
}

func CreateAccEndpointSecurityGroupSelectorUpdatedAttr(rName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing endpoint_security_group_selector attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_application_profile" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_endpoint_security_group" "test" {
		name 		= "%s"
		application_profile_dn = aci_application_profile.test.id
	}
	
	resource "aci_endpoint_security_group_selector" "test" {
		endpoint_security_group_dn  = aci_endpoint_security_group.test.id
		match_expression  = "ip==''%s"
		%s = "%s"
	}
	`, rName, rName, rName, ip, attribute, value)
	return resource
}
