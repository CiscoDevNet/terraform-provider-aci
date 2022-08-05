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

func TestAccAciEndpointLoopProtection_Basic(t *testing.T) {
	var endpoint_loop_protection_default models.EPLoopProtectionPolicy
	var endpoint_loop_protection_updated models.EPLoopProtectionPolicy
	resourceName := "aci_endpoint_loop_protection.test"
	epLoopProtectPolicy, err := aci.GetRemoteEPLoopProtectionPolicy(sharedAciClient(), "uni/infra/epLoopProtectP-default")
	if err != nil {
		t.Errorf("reading initial config of EP Loop Protection Policy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointLoopProtectionConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointLoopProtectionExists(resourceName, &endpoint_loop_protection_default),
					resource.TestCheckResourceAttrSet(resourceName, "admin_st"),
					resource.TestCheckResourceAttrSet(resourceName, "loop_detect_intvl"),
					resource.TestCheckResourceAttrSet(resourceName, "loop_detect_mult"),
				),
			},
			{
				Config: CreateAccEndpointLoopProtectionConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointLoopProtectionExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_endpoint_loop_protection"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0", "bd-learn-disable"),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_intvl", "30"),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "1"),
					testAccCheckAciEndpointLoopProtectionIdEqual(&endpoint_loop_protection_default, &endpoint_loop_protection_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccEndpointLoopProtectionInitialConfig(epLoopProtectPolicy),
			},
		},
	})
}

func TestAccAciEndpointLoopProtection_Update(t *testing.T) {
	var endpoint_loop_protection_default models.EPLoopProtectionPolicy
	var endpoint_loop_protection_updated models.EPLoopProtectionPolicy
	resourceName := "aci_endpoint_loop_protection.test"
	epLoopProtectPolicy, err := aci.GetRemoteEPLoopProtectionPolicy(sharedAciClient(), "uni/infra/epLoopProtectP-default")
	if err != nil {
		t.Errorf("reading initial config of EP Loop Protection Policy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointLoopProtectionConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointLoopProtectionExists(resourceName, &endpoint_loop_protection_default),
				),
			},
			{
				Config: CreateAccEndpointLoopProtectionUpdatedAttrList("action", StringListtoString([]string{"port-disable"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointLoopProtectionExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0", "port-disable"),
				),
			},
			{
				Config: CreateAccEndpointLoopProtectionUpdatedAttrList("action", StringListtoString([]string{"bd-learn-disable", "port-disable"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointLoopProtectionExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "action.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "action.0", "bd-learn-disable"),
					resource.TestCheckResourceAttr(resourceName, "action.1", "port-disable"),
				),
			},
			{
				Config: CreateAccEndpointLoopProtectionUpdatedAttrList("action", StringListtoString([]string{"port-disable", "bd-learn-disable"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointLoopProtectionExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "action.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "action.0", "port-disable"),
					resource.TestCheckResourceAttr(resourceName, "action.1", "bd-learn-disable"),
				),
			},
			{
				Config: CreateAccEndpointLoopProtectionUpdatedAttr("admin_st", "disabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointLoopProtectionExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "disabled"),
					testAccCheckAciEndpointLoopProtectionIdEqual(&endpoint_loop_protection_default, &endpoint_loop_protection_updated),
				),
			},
			{
				Config: CreateAccEndpointLoopProtectionUpdatedAttr("loop_detect_intvl", "300"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointLoopProtectionExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_intvl", "300"),
					testAccCheckAciEndpointLoopProtectionIdEqual(&endpoint_loop_protection_default, &endpoint_loop_protection_updated),
				),
			},
			{
				Config: CreateAccEndpointLoopProtectionUpdatedAttr("loop_detect_intvl", "135"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointLoopProtectionExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_intvl", "135"),
					testAccCheckAciEndpointLoopProtectionIdEqual(&endpoint_loop_protection_default, &endpoint_loop_protection_updated),
				),
			},
			{
				Config: CreateAccEndpointLoopProtectionUpdatedAttr("loop_detect_mult", "255"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointLoopProtectionExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "255"),
					testAccCheckAciEndpointLoopProtectionIdEqual(&endpoint_loop_protection_default, &endpoint_loop_protection_updated),
				),
			},
			{
				Config: CreateAccEndpointLoopProtectionUpdatedAttr("loop_detect_mult", "127"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointLoopProtectionExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "127"),
					testAccCheckAciEndpointLoopProtectionIdEqual(&endpoint_loop_protection_default, &endpoint_loop_protection_updated),
				),
			},
			{
				Config: CreateAccEndpointLoopProtectionInitialConfig(epLoopProtectPolicy),
			},
		},
	})
}

func TestAccAciEndpointLoopProtection_Negative(t *testing.T) {
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	epLoopProtectPolicy, err := aci.GetRemoteEPLoopProtectionPolicy(sharedAciClient(), "uni/infra/epLoopProtectP-default")
	if err != nil {
		t.Errorf("reading initial config of EP Loop Protection Policy")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointLoopProtectionConfig(),
			},
			{
				Config:      CreateAccEndpointLoopProtectionUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointLoopProtectionUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointLoopProtectionUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointLoopProtectionUpdatedAttrList("action", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccEndpointLoopProtectionUpdatedAttrList("action", StringListtoString([]string{"bd-learn-disable", "bd-learn-disable"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccEndpointLoopProtectionUpdatedAttr("admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccEndpointLoopProtectionUpdatedAttr("loop_detect_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndpointLoopProtectionUpdatedAttr("loop_detect_intvl", "29"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEndpointLoopProtectionUpdatedAttr("loop_detect_intvl", "301"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEndpointLoopProtectionUpdatedAttr("loop_detect_mult", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndpointLoopProtectionUpdatedAttr("loop_detect_mult", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEndpointLoopProtectionUpdatedAttr("loop_detect_mult", "256"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndpointLoopProtectionUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccEndpointLoopProtectionInitialConfig(epLoopProtectPolicy),
			},
		},
	})
}

func testAccCheckAciEndpointLoopProtectionExists(name string, endpoint_loop_protection *models.EPLoopProtectionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Endpoint Loop Protection %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Endpoint Loop Protection dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		endpoint_loop_protectionFound := models.EPLoopProtectionPolicyFromContainer(cont)
		if endpoint_loop_protectionFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Endpoint Loop Protection %s not found", rs.Primary.ID)
		}
		*endpoint_loop_protection = *endpoint_loop_protectionFound
		return nil
	}
}

func testAccCheckAciEndpointLoopProtectionIdEqual(m1, m2 *models.EPLoopProtectionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("endpoint_loop_protection DNs are not equal")
		}
		return nil
	}
}

func CreateAccEndpointLoopProtectionConfig() string {
	fmt.Println("=== STEP  testing endpoint_loop_protection creation")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	}
	`)
	return resource
}

func CreateAccEndpointLoopProtectionConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing endpoint_loop_protection creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_endpoint_loop_protection" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_endpoint_loop_protection"
		action = ["bd-learn-disable"]
		admin_st = "enabled"
		loop_detect_intvl = "30"
		loop_detect_mult = "1"		
	}
	`)

	return resource
}

func CreateAccEndpointLoopProtectionInitialConfig(epLoopProtection *models.EPLoopProtectionPolicy) string {
	fmt.Println("=== STEP  Basic: testing endpoint_loop_protection creation with initial config")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
		description = "%s"
		annotation = "%s"
		name_alias = "%s"
		action = %s
		admin_st = "%s"
		loop_detect_intvl = "%s"
		loop_detect_mult = "%s"		
	}

	`, epLoopProtection.Description, epLoopProtection.Annotation, epLoopProtection.NameAlias, StringListtoString(convertToStringArray(epLoopProtection.Action)), epLoopProtection.AdminSt, epLoopProtection.LoopDetectIntvl, epLoopProtection.LoopDetectMult)
	return resource
}
func CreateAccEndpointLoopProtectionUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing endpoint_loop_protection attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {	
		%s = "%s"
	}
	`, attribute, value)
	return resource
}

func CreateAccEndpointLoopProtectionUpdatedAttrList(attribute, value string) string {
	fmt.Printf("=== STEP  testing endpoint_loop_protection attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
		%s = %s
	}
	`, attribute, value)
	return resource
}
