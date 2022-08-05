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

func TestAccAciRouteControlContext_Basic(t *testing.T) {
	var route_control_context_default models.RouteControlContext
	var route_control_context_updated models.RouteControlContext
	resourceName := "aci_route_control_context.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	rtctrlProfileName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRouteControlContextDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateRouteControlContextWithoutRequired(fvTenantName, rtctrlProfileName, rName, "route_control_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateRouteControlContextWithoutRequired(fvTenantName, rtctrlProfileName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRouteControlContextConfig(fvTenantName, rtctrlProfileName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRouteControlContextExists(resourceName, &route_control_context_default),
					resource.TestCheckResourceAttr(resourceName, "route_control_profile_dn", fmt.Sprintf("uni/tn-%s/prof-%s", fvTenantName, rtctrlProfileName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "action", "permit"),
					// resource.TestCheckResourceAttr(resourceName, "set_rule", ""),
					resource.TestCheckResourceAttr(resourceName, "order", "0"),
				),
			},
			{
				Config: CreateAccRouteControlContextConfigWithOptionalValues(fvTenantName, rtctrlProfileName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRouteControlContextExists(resourceName, &route_control_context_updated),
					resource.TestCheckResourceAttr(resourceName, "route_control_profile_dn", fmt.Sprintf("uni/tn-%s/prof-%s", fvTenantName, rtctrlProfileName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_route_control_context"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
					resource.TestCheckResourceAttr(resourceName, "order", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_rule", fmt.Sprintf("uni/tn-%s/attr-%s", fvTenantName, rtctrlProfileName)),

					testAccCheckAciRouteControlContextIdEqual(&route_control_context_default, &route_control_context_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccRouteControlContextConfigUpdatedName(fvTenantName, rtctrlProfileName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccRouteControlContextRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRouteControlContextConfigWithRequiredParams(rName, rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRouteControlContextExists(resourceName, &route_control_context_updated),
					resource.TestCheckResourceAttr(resourceName, "route_control_profile_dn", fmt.Sprintf("uni/tn-%s/prof-%s", rName, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciRouteControlContextIdNotEqual(&route_control_context_default, &route_control_context_updated),
				),
			},
			{
				Config: CreateAccRouteControlContextConfig(fvTenantName, rtctrlProfileName, rName),
			},
			{
				Config: CreateAccRouteControlContextConfigWithRequiredParams(rName, rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRouteControlContextExists(resourceName, &route_control_context_updated),
					resource.TestCheckResourceAttr(resourceName, "route_control_profile_dn", fmt.Sprintf("uni/tn-%s/prof-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciRouteControlContextIdNotEqual(&route_control_context_default, &route_control_context_updated),
				),
			},
		},
	})
}

func TestAccAciRouteControlContext_Update(t *testing.T) {
	var route_control_context_default models.RouteControlContext
	var route_control_context_updated models.RouteControlContext
	resourceName := "aci_route_control_context.test"
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	rtctrlProfileName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRouteControlContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRouteControlContextConfig(fvTenantName, rtctrlProfileName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRouteControlContextExists(resourceName, &route_control_context_default),
				),
			},
			{
				Config: CreateAccRouteControlContextUpdatedAttr(fvTenantName, rtctrlProfileName, rName, "order", "9"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRouteControlContextExists(resourceName, &route_control_context_updated),
					resource.TestCheckResourceAttr(resourceName, "order", "9"),
					testAccCheckAciRouteControlContextIdEqual(&route_control_context_default, &route_control_context_updated),
				),
			},
			{
				Config: CreateAccRouteControlContextUpdatedAttr(fvTenantName, rtctrlProfileName, rName, "order", "4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRouteControlContextExists(resourceName, &route_control_context_updated),
					resource.TestCheckResourceAttr(resourceName, "order", "4"),
					testAccCheckAciRouteControlContextIdEqual(&route_control_context_default, &route_control_context_updated),
				),
			},

			{
				Config: CreateAccRouteControlContextConfig(fvTenantName, rtctrlProfileName, rName),
			},
		},
	})
}

func TestAccAciRouteControlContext_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	rtctrlProfileName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRouteControlContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRouteControlContextConfig(fvTenantName, rtctrlProfileName, rName),
			},
			{
				Config:      CreateAccRouteControlContextWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`Invalid request`),
			},
			{
				Config:      CreateAccRouteControlContextUpdatedAttr(fvTenantName, rtctrlProfileName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRouteControlContextUpdatedAttr(fvTenantName, rtctrlProfileName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRouteControlContextUpdatedAttr(fvTenantName, rtctrlProfileName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccRouteControlContextUpdatedAttr(fvTenantName, rtctrlProfileName, rName, "action", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRouteControlContextUpdatedAttr(fvTenantName, rtctrlProfileName, rName, "order", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRouteControlContextUpdatedAttr(fvTenantName, rtctrlProfileName, rName, "order", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRouteControlContextUpdatedAttr(fvTenantName, rtctrlProfileName, rName, "order", "10"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccRouteControlContextUpdatedAttr(fvTenantName, rtctrlProfileName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccRouteControlContextConfig(fvTenantName, rtctrlProfileName, rName),
			},
		},
	})
}

func TestAccAciRouteControlContext_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	rtctrlProfileName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRouteControlContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRouteControlContextConfigMultiple(fvTenantName, rtctrlProfileName, rName),
			},
		},
	})
}

func testAccCheckAciRouteControlContextExists(name string, route_control_context *models.RouteControlContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Route Control Context %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route Control Context dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		route_control_contextFound := models.RouteControlContextFromContainer(cont)
		if route_control_contextFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Route Control Context %s not found", rs.Primary.ID)
		}
		*route_control_context = *route_control_contextFound
		return nil
	}
}

func testAccCheckAciRouteControlContextDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing route_control_context destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_route_control_context" {
			cont, err := client.Get(rs.Primary.ID)
			route_control_context := models.RouteControlContextFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Route Control Context %s Still exists", route_control_context.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRouteControlContextIdEqual(m1, m2 *models.RouteControlContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("route_control_context DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciRouteControlContextIdNotEqual(m1, m2 *models.RouteControlContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("route_control_context DNs are equal")
		}
		return nil
	}
}

func CreateRouteControlContextWithoutRequired(fvTenantName, rtctrlProfileName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing route_control_context creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		name 		= "%s"
		parent_dn = aci_tenant.test.id
	}
	
	`
	switch attrName {
	case "route_control_profile_dn":
		rBlock += `
	resource "aci_route_control_context" "test" {
	#	route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rtctrlProfileName, rName)
}

func CreateAccRouteControlContextConfigWithRequiredParams(fvTenantName, rtctrlProfileName, rName string) string {
	fmt.Println("=== STEP  testing route_control_context creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		name 		= "%s"
		parent_dn = aci_tenant.test.id
	}
	
	resource "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = "%s"
	}
	`, fvTenantName, rtctrlProfileName, rName)
	return resource
}
func CreateAccRouteControlContextConfigUpdatedName(fvTenantName, rtctrlProfileName, rName string) string {
	fmt.Println("=== STEP  testing route_control_context creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		name 		= "%s"
		parent_dn = aci_tenant.test.id
	}
	
	resource "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = "%s"
	}
	`, fvTenantName, rtctrlProfileName, rName)
	return resource
}

func CreateAccRouteControlContextConfig(fvTenantName, rtctrlProfileName, rName string) string {
	fmt.Println("=== STEP  testing route_control_context creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		name 		= "%s"
		parent_dn = aci_tenant.test.id
	}
	
	resource "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = "%s"
	}
	`, fvTenantName, rtctrlProfileName, rName)
	return resource
}

func CreateAccRouteControlContextConfigMultiple(fvTenantName, rtctrlProfileName, rName string) string {
	fmt.Println("=== STEP  testing multiple route_control_context creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		name 		= "%s"
		parent_dn = aci_tenant.test.id
	}
	
	resource "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, rtctrlProfileName, rName)
	return resource
}

func CreateAccRouteControlContextWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing route_control_context creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccRouteControlContextConfigWithOptionalValues(fvTenantName, rtctrlProfileName, rName string) string {
	fmt.Println("=== STEP  Basic: testing route_control_context creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		name 		= "%s"
		parent_dn = aci_tenant.test.id
	}

	resource "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	
	resource "aci_route_control_context" "test" {
		route_control_profile_dn  = "${aci_bgp_route_control_profile.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_route_control_context"
		action = "deny"
		order = "1"
		set_rule = aci_action_rule_profile.test.id
	}
	`, fvTenantName, rtctrlProfileName, rtctrlProfileName, rName)

	return resource
}

func CreateAccRouteControlContextRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing route_control_context updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_route_control_context" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_route_control_context"
		action = "deny"
		order = "1"
		
	}
	`)

	return resource
}

func CreateAccRouteControlContextUpdatedAttr(fvTenantName, rtctrlProfileName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing route_control_context attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		name 		= "%s"
		parent_dn = aci_tenant.test.id
	}
	
	resource "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rtctrlProfileName, rName, attribute, value)
	return resource
}
