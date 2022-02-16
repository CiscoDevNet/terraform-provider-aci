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
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciEndpointControls_Basic(t *testing.T) {
	var endpoint_controls_default models.EndpointControlPolicy
	var endpoint_controls_updated models.EndpointControlPolicy
	resourceName := "aci_endpoint_controls.test"
	epControls, err := aci.GetRemoteEndpointControlPolicy(sharedAciClient(), "uni/infra/epCtrlP-default")
	if err != nil {
		t.Errorf("reading initial config of EP Loop Protection Policy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointControlsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlsExists(resourceName, &endpoint_controls_default),
					resource.TestCheckResourceAttrSet(resourceName, "admin_st"),
					resource.TestCheckResourceAttrSet(resourceName, "hold_intvl"),
					resource.TestCheckResourceAttrSet(resourceName, "rogue_ep_detect_intvl"),
					resource.TestCheckResourceAttrSet(resourceName, "rogue_ep_detect_mult"),
				),
			},
			{
				Config: CreateAccEndpointControlsConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlsExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_endpoint_controls"),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "hold_intvl", "300"),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_intvl", "30"),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_mult", "2"),
					testAccCheckAciEndpointControlsIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccEndpointControlsInitialConfig(epControls),
			},
		},
	})
}

func TestAccAciEndpointControls_Update(t *testing.T) {
	var endpoint_controls_default models.EndpointControlPolicy
	var endpoint_controls_updated models.EndpointControlPolicy
	epControls, err := aci.GetRemoteEndpointControlPolicy(sharedAciClient(), "uni/infra/epCtrlP-default")
	if err != nil {
		t.Errorf("reading initial config of EP Loop Protection Policy")
	}
	resourceName := "aci_endpoint_controls.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointControlsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlsExists(resourceName, &endpoint_controls_default),
				),
			},
			{
				Config: CreateAccEndpointControlsUpdatedAttr("admin_st", "disabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlsExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "disabled"),
					testAccCheckAciEndpointControlsIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				Config: CreateAccEndpointControlsUpdatedAttr("hold_intvl", "3600"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlsExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "hold_intvl", "3600"),
					testAccCheckAciEndpointControlsIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				Config: CreateAccEndpointControlsUpdatedAttr("hold_intvl", "1650"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlsExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "hold_intvl", "1650"),
					testAccCheckAciEndpointControlsIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				Config: CreateAccEndpointControlsUpdatedAttr("rogue_ep_detect_intvl", "3600"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlsExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_intvl", "3600"),
					testAccCheckAciEndpointControlsIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				Config: CreateAccEndpointControlsUpdatedAttr("rogue_ep_detect_intvl", "1785"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlsExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_intvl", "1785"),
					testAccCheckAciEndpointControlsIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				Config: CreateAccEndpointControlsUpdatedAttr("rogue_ep_detect_mult", "65535"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlsExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_mult", "65535"),
					testAccCheckAciEndpointControlsIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				Config: CreateAccEndpointControlsUpdatedAttr("rogue_ep_detect_mult", "32766"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlsExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_mult", "32766"),
					testAccCheckAciEndpointControlsIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				Config: CreateAccEndpointControlsInitialConfig(epControls),
			},
		},
	})
}

func TestAccAciEndpointControls_Negative(t *testing.T) {
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	epControls, err := aci.GetRemoteEndpointControlPolicy(sharedAciClient(), "uni/infra/epCtrlP-default")
	if err != nil {
		t.Errorf("reading initial config of EP Loop Protection Policy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointControlsConfig(),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr("admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr("hold_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr("hold_intvl", "299"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr("hold_intvl", "3601"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr("rogue_ep_detect_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr("rogue_ep_detect_intvl", "29"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr("rogue_ep_detect_intvl", "3601"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr("rogue_ep_detect_mult", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr("rogue_ep_detect_mult", "1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr("rogue_ep_detect_mult", "65536"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndpointControlsUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccEndpointControlsInitialConfig(epControls),
			},
		},
	})
}

func testAccCheckAciEndpointControlsExists(name string, endpoint_controls *models.EndpointControlPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Endpoint Controls %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Endpoint Controls dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		endpoint_controlsFound := models.EndpointControlPolicyFromContainer(cont)
		if endpoint_controlsFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Endpoint Controls %s not found", rs.Primary.ID)
		}
		*endpoint_controls = *endpoint_controlsFound
		return nil
	}
}

func testAccCheckAciEndpointControlsIdEqual(m1, m2 *models.EndpointControlPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("endpoint_controls DNs are not equal")
		}
		return nil
	}
}

func CreateAccEndpointControlsConfig() string {
	fmt.Println("=== STEP  testing endpoint_controls creation")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	}
	`)
	return resource
}

func CreateAccEndpointControlsConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing endpoint_controls creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_endpoint_controls"
		admin_st = "enabled"
		hold_intvl = "300"
		rogue_ep_detect_intvl = "30"
		rogue_ep_detect_mult = "2"
		
	}
	`)
	return resource
}

func CreateAccEndpointControlsInitialConfig(epControls *models.EndpointControlPolicy) string {
	fmt.Println("=== STEP  Basic: testing endpoint_controls creation with Initial Config")
	resource := fmt.Sprintf(`
	resource "aci_endpoint_controls" "test" {
		description = "%s"
		annotation = "%s"
		name_alias = "%s"
		admin_st = "%s"
		hold_intvl = "%s"
		rogue_ep_detect_intvl = "%s"
		rogue_ep_detect_mult = "%s"
	}
	`, epControls.Description, epControls.Annotation, epControls.NameAlias, epControls.AdminSt, epControls.HoldIntvl, epControls.RogueEpDetectIntvl, epControls.RogueEpDetectMult)
	return resource
}

func CreateAccEndpointControlsUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing endpoint_controls attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
