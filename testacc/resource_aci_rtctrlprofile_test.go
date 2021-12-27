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

func TestAccAciBgpRouteControlProfile_Basic(t *testing.T) {
	var bgp_route_control_profile_default models.RouteControlProfile
	var bgp_route_control_profile_updated models.RouteControlProfile
	resourceName := "aci_bgp_route_control_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBgpRouteControlProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBgpRouteControlProfileWithoutRequired(rName, rName, "parent_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBgpRouteControlProfileWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBgpRouteControlProfileConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteControlProfileExists(resourceName, &bgp_route_control_profile_default),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "route_control_profile_type", "combinable"),
				),
			},
			{
				Config: CreateAccBgpRouteControlProfileConfigWithOptionalValues(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteControlProfileExists(resourceName, &bgp_route_control_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_route_control_profile"),
					resource.TestCheckResourceAttr(resourceName, "route_control_profile_type", "global"),
					testAccCheckAciBgpRouteControlProfileIdEqual(&bgp_route_control_profile_default, &bgp_route_control_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccBgpRouteControlProfileConfigUpdatedName(rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config: CreateAccBgpRouteControlProfileConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteControlProfileExists(resourceName, &bgp_route_control_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciBgpRouteControlProfileIdNotEqual(&bgp_route_control_profile_default, &bgp_route_control_profile_updated),
				),
			},
			{
				Config: CreateAccBgpRouteControlProfileConfig(rName, rName),
			},
			{
				Config: CreateAccBgpRouteControlProfileConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteControlProfileExists(resourceName, &bgp_route_control_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciBgpRouteControlProfileIdNotEqual(&bgp_route_control_profile_default, &bgp_route_control_profile_updated),
				),
			},
			{
				Config:      CreateAccBgpRouteControlProfileUpdateWithoutRequiredAttr(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBgpRouteControlProfileConfigL3Outside(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteControlProfileExists(resourceName, &bgp_route_control_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s/out-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "route_control_profile_type", "combinable"),
					testAccCheckAciBgpRouteControlProfileIdNotEqual(&bgp_route_control_profile_default, &bgp_route_control_profile_updated),
				),
			},
		},
	})
}

func TestAccAciBgpRouteControlProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBgpRouteControlProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBgpRouteControlProfileConfig(rName, rName),
			},
			{
				Config:      CreateAccBgpRouteControlProfileWithInValidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class rtctrlProfile (.)+`),
			},
			{
				Config:      CreateAccBgpRouteControlProfileUpdatedAttr(rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBgpRouteControlProfileUpdatedAttr(rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBgpRouteControlProfileUpdatedAttr(rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBgpRouteControlProfileUpdatedAttr(rName, rName, "route_control_profile_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+to be one of(.)+, got(.)+`),
			},

			{
				Config:      CreateAccBgpRouteControlProfileUpdatedAttr(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config: CreateAccBgpRouteControlProfileConfig(rName, rName),
			},
		},
	})
}

func TestAccAciBgpRouteControlProfile_Multiple(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBgpRouteControlProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBgpRouteControlProfileConfigs(rName, rName),
			},
		},
	})
}

func CreateAccBgpRouteControlProfileConfigs(fvTenantName, rName string) string {
	fmt.Println("=== STEP Testing Multiple route_control_profile creation")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_bgp_route_control_profile" "test1" {
		parent_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_bgp_route_control_profile" "test2" {
		parent_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName, rName+"1", rName+"2")
	return resource
}

func testAccCheckAciBgpRouteControlProfileExists(name string, bgp_route_control_profile *models.RouteControlProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Bgp Route Control Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Bgp Route Control Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bgp_route_control_profileFound := models.RouteControlProfileFromContainer(cont)
		if bgp_route_control_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Bgp Route Control Profile %s not found", rs.Primary.ID)
		}
		*bgp_route_control_profile = *bgp_route_control_profileFound
		return nil
	}
}

func testAccCheckAciBgpRouteControlProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  Testing bgp_route_control_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_bgp_route_control_profile" {
			cont, err := client.Get(rs.Primary.ID)
			bgp_route_control_profile := models.RouteControlProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Bgp Route Control Profile %s Still exists", bgp_route_control_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBgpRouteControlProfileIdEqual(m1, m2 *models.RouteControlProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("bgp_route_control_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciBgpRouteControlProfileIdNotEqual(m1, m2 *models.RouteControlProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("bgp_route_control_profile DNs are equal")
		}
		return nil
	}
}

func CreateBgpRouteControlProfileWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: Testing bgp_route_control_profile creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "parent_dn":
		rBlock += `
	resource "aci_bgp_route_control_profile" "test" {
	#	parent_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBgpRouteControlProfileConfigWithRequiredParams(ParentName, rName string) string {
	fmt.Printf("=== STEP  Testing bgp_route_control_profile creation with tenant name = %s and bgp_route_control_profile name = %s\n", ParentName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, ParentName, rName)
	return resource
}

func CreateAccBgpRouteControlProfileConfigL3Outside(ParentName, rName string) string {
	fmt.Println("=== STEP  Testing bgp_route_control_profile creation when parent is l3outside")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        name           = "%s"
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_l3_outside.test.id
		name  = "%s"
	}
	`, ParentName, ParentName, rName)
	return resource
}

func CreateAccBgpRouteControlProfileConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Testing bgp_route_control_profile creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBgpRouteControlProfileConfigUpdatedName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Testing bgp_route_control_profile creation with Invalid Long Name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBgpRouteControlProfileWithInValidParentDn(ParentName, rName string) string {
	fmt.Println("=== STEP  Negative Case: Testing bgp_route_control_profile creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	resource "aci_application_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name       = "%s"
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_application_profile.test.id
		name  = "%s" 
	}
	`, ParentName, ParentName, rName)
	return resource
}

func CreateAccBgpRouteControlProfileConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: Testing bgp_route_control_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_route_control_profile"
		route_control_profile_type = "global"
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccBgpRouteControlProfileUpdateWithoutRequiredAttr() string {
	fmt.Println("=== STEP  Basic: Testing bgp_route_control_profile updation without required attributes")
	resource := fmt.Sprintln(`
	resource "aci_bgp_route_control_profile" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_route_control_profile"
		route_control_profile_type = "global"
	}
	`)

	return resource
}

func CreateAccBgpRouteControlProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: Testing bgp_route_control_profile creation with optional parameters")
	resource := fmt.Sprintln(`
	resource "aci_bgp_route_control_profile" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_route_control_profile"
		route_control_profile_type = "global"
	}
	`)

	return resource
}

func CreateAccBgpRouteControlProfileUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  Testing bgp_route_control_profile attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		parent_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
