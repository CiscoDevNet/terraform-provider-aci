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

func TestAccAciMatchRouteDestinationRule_Basic(t *testing.T) {
	var match_route_destination_rule_default models.MatchRouteDestinationRule
	var match_route_destination_rule_updated models.MatchRouteDestinationRule
	resourceName := "aci_match_route_destination_rule.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.0.0.0/16")
	ip = fmt.Sprintf("%s/16", ip)
	ipUpdated, _ := acctest.RandIpAddress("10.1.0.0/16")
	fvTenantName := makeTestVariable(acctest.RandString(5))
	rtctrlSubjPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMatchRouteDestinationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateMatchRouteDestinationRuleWithoutRequired(fvTenantName, rtctrlSubjPName, ip, "match_rule_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMatchRouteDestinationRuleWithoutRequired(fvTenantName, rtctrlSubjPName, ip, "ip"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMatchRouteDestinationRuleConfig(fvTenantName, rtctrlSubjPName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRouteDestinationRuleExists(resourceName, &match_route_destination_rule_default),
					resource.TestCheckResourceAttr(resourceName, "match_rule_dn", fmt.Sprintf("uni/tn-%s/subj-%s", fvTenantName, rtctrlSubjPName)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "aggregate", "no"),
					resource.TestCheckResourceAttr(resourceName, "greater_than_mask", "0"),
					resource.TestCheckResourceAttr(resourceName, "less_than_mask", "0"),
				),
			},
			{
				Config: CreateAccMatchRouteDestinationRuleConfigWithOptionalValues(fvTenantName, rtctrlSubjPName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRouteDestinationRuleExists(resourceName, &match_route_destination_rule_updated),
					resource.TestCheckResourceAttr(resourceName, "match_rule_dn", fmt.Sprintf("uni/tn-%s/subj-%s", fvTenantName, rtctrlSubjPName)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_match_route_destination_rule"),
					resource.TestCheckResourceAttr(resourceName, "aggregate", "yes"),
					resource.TestCheckResourceAttr(resourceName, "greater_than_mask", "32"),
					resource.TestCheckResourceAttr(resourceName, "less_than_mask", "32"),
					testAccCheckAciMatchRouteDestinationRuleIdEqual(&match_route_destination_rule_default, &match_route_destination_rule_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"match_rule_dn"},
			},
			{
				Config:      CreateAccMatchRouteDestinationRuleWithInavalidIP(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccMatchRouteDestinationRuleRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMatchRouteDestinationRuleConfigWithRequiredParams(rNameUpdated, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRouteDestinationRuleExists(resourceName, &match_route_destination_rule_updated),
					resource.TestCheckResourceAttr(resourceName, "match_rule_dn", fmt.Sprintf("uni/tn-%s/subj-%s", rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					testAccCheckAciMatchRouteDestinationRuleIdNotEqual(&match_route_destination_rule_default, &match_route_destination_rule_updated),
				),
			},
			{
				Config: CreateAccMatchRouteDestinationRuleConfig(fvTenantName, rtctrlSubjPName, ip),
			},
			{
				Config: CreateAccMatchRouteDestinationRuleConfigWithRequiredParams(rName, ipUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRouteDestinationRuleExists(resourceName, &match_route_destination_rule_updated),
					resource.TestCheckResourceAttr(resourceName, "match_rule_dn", fmt.Sprintf("uni/tn-%s/subj-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "ip", ipUpdated),
					testAccCheckAciMatchRouteDestinationRuleIdNotEqual(&match_route_destination_rule_default, &match_route_destination_rule_updated),
				),
			},
		},
	})
}

func TestAccAciMatchRouteDestinationRule_Update(t *testing.T) {
	var match_route_destination_rule_default models.MatchRouteDestinationRule
	var match_route_destination_rule_updated models.MatchRouteDestinationRule
	resourceName := "aci_match_route_destination_rule.test"
	ip := "2001:db8:85a3::8a2e:370:7334/24"
	fvTenantName := makeTestVariable(acctest.RandString(5))
	rtctrlSubjPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMatchRouteDestinationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMatchRouteDestinationRuleConfig(fvTenantName, rtctrlSubjPName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRouteDestinationRuleExists(resourceName, &match_route_destination_rule_default),
				),
			},
			{
				Config: CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "greater_than_mask", "64"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRouteDestinationRuleExists(resourceName, &match_route_destination_rule_updated),
					resource.TestCheckResourceAttr(resourceName, "greater_than_mask", "64"),
					testAccCheckAciMatchRouteDestinationRuleIdEqual(&match_route_destination_rule_default, &match_route_destination_rule_updated),
				),
			},
			{
				Config: CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "less_than_mask", "64"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRouteDestinationRuleExists(resourceName, &match_route_destination_rule_updated),
					resource.TestCheckResourceAttr(resourceName, "less_than_mask", "64"),
					testAccCheckAciMatchRouteDestinationRuleIdEqual(&match_route_destination_rule_default, &match_route_destination_rule_updated),
				),
			},
			{
				Config: CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "greater_than_mask", "128"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRouteDestinationRuleExists(resourceName, &match_route_destination_rule_updated),
					resource.TestCheckResourceAttr(resourceName, "greater_than_mask", "128"),
					testAccCheckAciMatchRouteDestinationRuleIdEqual(&match_route_destination_rule_default, &match_route_destination_rule_updated),
				),
			},
			{
				Config: CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "less_than_mask", "128"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRouteDestinationRuleExists(resourceName, &match_route_destination_rule_updated),
					resource.TestCheckResourceAttr(resourceName, "less_than_mask", "128"),
					testAccCheckAciMatchRouteDestinationRuleIdEqual(&match_route_destination_rule_default, &match_route_destination_rule_updated),
				),
			},
			{
				Config: CreateAccMatchRouteDestinationRuleConfig(fvTenantName, rtctrlSubjPName, ip),
			},
		},
	})
}

func TestAccAciMatchRouteDestinationRule_Negative(t *testing.T) {
	ip, _ := acctest.RandIpAddress("10.3.0.0/16")
	rName := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	rtctrlSubjPName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMatchRouteDestinationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMatchRouteDestinationRuleConfig(fvTenantName, rtctrlSubjPName, ip),
			},
			{
				Config:      CreateAccMatchRouteDestinationRuleWithInValidParentDn(rName, ip),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "aggregate", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "greater_than_mask", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "greater_than_mask", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "greater_than_mask", "129"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "less_than_mask", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "less_than_mask", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, "less_than_mask", "129"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccMatchRouteDestinationRuleConfig(fvTenantName, rtctrlSubjPName, ip),
			},
		},
	})
}

func TestAccAciMatchRouteDestinationRule_MultipleCreateDelete(t *testing.T) {
	fvTenantName := makeTestVariable(acctest.RandString(5))
	rtctrlSubjPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMatchRouteDestinationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMatchRouteDestinationRuleConfigMultiple(fvTenantName, rtctrlSubjPName),
			},
		},
	})
}

func testAccCheckAciMatchRouteDestinationRuleExists(name string, match_route_destination_rule *models.MatchRouteDestinationRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Match Route Destination Rule %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Match Route Destination Rule dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		match_route_destination_ruleFound := models.MatchRouteDestinationRuleFromContainer(cont)
		if match_route_destination_ruleFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Match Route Destination Rule %s not found", rs.Primary.ID)
		}
		*match_route_destination_rule = *match_route_destination_ruleFound
		return nil
	}
}

func testAccCheckAciMatchRouteDestinationRuleDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing match_route_destination_rule destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_match_route_destination_rule" {
			cont, err := client.Get(rs.Primary.ID)
			match_route_destination_rule := models.MatchRouteDestinationRuleFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Match Route Destination Rule %s Still exists", match_route_destination_rule.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciMatchRouteDestinationRuleIdEqual(m1, m2 *models.MatchRouteDestinationRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("match_route_destination_rule DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciMatchRouteDestinationRuleIdNotEqual(m1, m2 *models.MatchRouteDestinationRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("match_route_destination_rule DNs are equal")
		}
		return nil
	}
}

func CreateMatchRouteDestinationRuleWithoutRequired(fvTenantName, rtctrlSubjPName, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing match_route_destination_rule creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_match_rule" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	`
	switch attrName {
	case "match_rule_dn":
		rBlock += `
	resource "aci_match_route_destination_rule" "test" {
	#	match_rule_dn  = aci_match_rule.test.id
		ip  = "%s"
	}
		`
	case "ip":
		rBlock += `
	resource "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
	#	ip  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rtctrlSubjPName, ip)
}

func CreateAccMatchRouteDestinationRuleConfigWithRequiredParams(prName, ip string) string {
	fmt.Printf("=== STEP  testing match_route_destination_rule creation with parent resource name %s and ip %s\n", prName, ip)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
		ip  = "%s"
	}
	`, prName, prName, ip)
	return resource
}

func CreateAccMatchRouteDestinationRuleConfig(fvTenantName, rtctrlSubjPName, ip string) string {
	fmt.Println("=== STEP  testing match_route_destination_rule creation with required arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
		ip  = "%s"
	}
	`, fvTenantName, rtctrlSubjPName, ip)
	return resource
}

func CreateAccMatchRouteDestinationRuleWithInavalidIP(rName string) string {
	fmt.Println("=== STEP  testing match_route_destination_rule creation with required invalid ip")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
		ip  = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccMatchRouteDestinationRuleConfigMultiple(fvTenantName, rtctrlSubjPName string) string {
	fmt.Println("=== STEP  testing multiple match_route_destination_rule creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
		ip  = "10.4.0.${count.index}"
		count = 5
	}
	`, fvTenantName, rtctrlSubjPName)
	return resource
}

func CreateAccMatchRouteDestinationRuleWithInValidParentDn(rName, ip string) string {
	fmt.Println("=== STEP  Negative Case: testing match_route_destination_rule creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_tenant.test.id
		ip  = "%s"	
	}
	`, rName, ip)
	return resource
}

func CreateAccMatchRouteDestinationRuleConfigWithOptionalValues(fvTenantName, rtctrlSubjPName, ip string) string {
	fmt.Println("=== STEP  Basic: testing match_route_destination_rule creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_match_route_destination_rule" "test" {
		match_rule_dn  = "${aci_match_rule.test.id}"
		ip  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_match_route_destination_rule"
		aggregate = "yes"
		greater_than_mask = "32"
		less_than_mask = "32"
		
	}
	`, fvTenantName, rtctrlSubjPName, ip)

	return resource
}

func CreateAccMatchRouteDestinationRuleRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing match_route_destination_rule updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_match_route_destination_rule" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_match_route_destination_rule"
		aggregate = "yes"
		from_pfx_len = "1"
		to_pfx_len = "1"
		
	}
	`)

	return resource
}

func CreateAccMatchRouteDestinationRuleUpdatedAttr(fvTenantName, rtctrlSubjPName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing match_route_destination_rule attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_match_rule" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_match_route_destination_rule" "test" {
		match_rule_dn  = aci_match_rule.test.id
		ip  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rtctrlSubjPName, ip, attribute, value)
	return resource
}
