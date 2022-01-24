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

func TestAccAciPortTracking_Basic(t *testing.T) {
	var port_tracking_default models.PortTracking
	var port_tracking_updated models.PortTracking
	resourceName := "aci_port_tracking.test"
	infraPortTrackPol, err := aci.GetRemotePortTracking(sharedAciClient(), "uni/infra/trackEqptFabP-default")
	if err != nil {
		t.Errorf("reading initial config of infraPortTrackPol")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPortTrackingConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "admin_st"),
					resource.TestCheckResourceAttrSet(resourceName, "delay"),
					resource.TestCheckResourceAttrSet(resourceName, "include_apic_ports"),
					resource.TestCheckResourceAttrSet(resourceName, "minlinks"),
				),
			},
			{
				Config: CreateAccPortTrackingConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_port_tracking"),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "on"),
					resource.TestCheckResourceAttr(resourceName, "delay", "1"),
					resource.TestCheckResourceAttr(resourceName, "include_apic_ports", "yes"),
					resource.TestCheckResourceAttr(resourceName, "minlinks", "0"),
					testAccCheckAciPortTrackingIdEqual(&port_tracking_default, &port_tracking_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: restorePortTrackingConfig(infraPortTrackPol),
			},
		},
	})
}

func TestAccAciPortTracking_Update(t *testing.T) {
	var port_tracking_default models.PortTracking
	var port_tracking_updated models.PortTracking
	resourceName := "aci_port_tracking.test"
	infraPortTrackPol, err := aci.GetRemotePortTracking(sharedAciClient(), "uni/infra/trackEqptFabP-default")
	if err != nil {
		t.Errorf("reading initial config of infraPortTrackPol")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPortTrackingConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_default),
				),
			},
			{
				Config: CreateAccPortTrackingUpdatedAttr("delay", "300"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_updated),
					resource.TestCheckResourceAttr(resourceName, "delay", "300"),
					testAccCheckAciPortTrackingIdEqual(&port_tracking_default, &port_tracking_updated),
				),
			},
			{
				Config: CreateAccPortTrackingUpdatedAttr("delay", "150"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_updated),
					resource.TestCheckResourceAttr(resourceName, "delay", "150"),
					testAccCheckAciPortTrackingIdEqual(&port_tracking_default, &port_tracking_updated),
				),
			},
			{
				Config: CreateAccPortTrackingUpdatedAttr("minlinks", "48"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_updated),
					resource.TestCheckResourceAttr(resourceName, "minlinks", "48"),
					testAccCheckAciPortTrackingIdEqual(&port_tracking_default, &port_tracking_updated),
				),
			},
			{
				Config: CreateAccPortTrackingUpdatedAttr("minlinks", "24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_updated),
					resource.TestCheckResourceAttr(resourceName, "minlinks", "24"),
					testAccCheckAciPortTrackingIdEqual(&port_tracking_default, &port_tracking_updated),
				),
			},
			{
				Config: CreateAccPortTrackingUpdatedAttr("admin_st", "off"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_updated),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "off"),
					testAccCheckAciPortTrackingIdEqual(&port_tracking_default, &port_tracking_updated),
				),
			},
			{
				Config: CreateAccPortTrackingUpdatedAttr("include_apic_ports", "no"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_updated),
					resource.TestCheckResourceAttr(resourceName, "include_apic_ports", "no"),
					testAccCheckAciPortTrackingIdEqual(&port_tracking_default, &port_tracking_updated),
				),
			},
			{
				Config: restorePortTrackingConfig(infraPortTrackPol),
			},
		},
	})
}

func TestAccAciPortTracking_Negative(t *testing.T) {
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	infraPortTrackPol, err := aci.GetRemotePortTracking(sharedAciClient(), "uni/infra/trackEqptFabP-default")
	if err != nil {
		t.Errorf("reading initial config of infraPortTrackPol")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPortTrackingConfig(),
			},

			{
				Config:      CreateAccPortTrackingUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccPortTrackingUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccPortTrackingUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccPortTrackingUpdatedAttr("admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccPortTrackingUpdatedAttr("delay", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccPortTrackingUpdatedAttr("delay", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccPortTrackingUpdatedAttr("delay", "301"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccPortTrackingUpdatedAttr("include_apic_ports", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccPortTrackingUpdatedAttr("minlinks", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccPortTrackingUpdatedAttr("minlinks", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccPortTrackingUpdatedAttr("minlinks", "49"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccPortTrackingUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: restorePortTrackingConfig(infraPortTrackPol),
			},
		},
	})
}

func restorePortTrackingConfig(infraPortTrackPol *models.PortTracking) string {
	resource := fmt.Sprintf(`
	resource "aci_port_tracking" "test" {
		description = "%s"
		annotation = "%s"
		name_alias = "%s"
		admin_st = "%s"
		delay = "%s"
		include_apic_ports = "%s"
		minlinks = "%s"
		
	}
	`, infraPortTrackPol.Description, infraPortTrackPol.Annotation, infraPortTrackPol.NameAlias, infraPortTrackPol.AdminSt, infraPortTrackPol.Delay, infraPortTrackPol.IncludeApicPorts, infraPortTrackPol.Minlinks)
	return resource
}

func testAccCheckAciPortTrackingExists(name string, port_tracking *models.PortTracking) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Port Tracking %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Port Tracking dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		port_trackingFound := models.PortTrackingFromContainer(cont)
		if port_trackingFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Port Tracking %s not found", rs.Primary.ID)
		}
		*port_tracking = *port_trackingFound
		return nil
	}
}

func testAccCheckAciPortTrackingIdEqual(m1, m2 *models.PortTracking) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("port_tracking DNs are not equal")
		}
		return nil
	}
}

func CreateAccPortTrackingConfig() string {
	fmt.Println("=== STEP  testing port_tracking creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {}
	`)
	return resource
}

func CreateAccPortTrackingConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing port_tracking creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_port_tracking"
		admin_st = "on"
		delay = "1"
		include_apic_ports = "yes"
		minlinks = "0"
		
	}
	`)

	return resource
}

func CreateAccPortTrackingUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing port_tracking attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
