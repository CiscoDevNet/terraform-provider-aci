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

// depends on fabricPathEp
func TestAccAciOutofServiceFabricPath_Basic(t *testing.T) {
	var interface_blacklist_default models.OutofServiceFabricPath
	var interface_blacklist_updated models.OutofServiceFabricPath
	resourceName := "aci_interface_blacklist.test"
	podId := "1"
	NodeId := "201"
	Interface := "eth1/1"
	InterfaceUpdated := "eth1/2"
	randomValue := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOutofServiceFabricPathDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateOutofServiceFabricPathWithoutRequired(podId, NodeId, Interface, "pod_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateOutofServiceFabricPathWithoutRequired(podId, NodeId, Interface, "node_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateOutofServiceFabricPathWithoutRequired(podId, NodeId, Interface, "interface"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccOutofServiceFabricPathConfig(podId, NodeId, Interface),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutofServiceFabricPathExists(resourceName, &interface_blacklist_default),

					resource.TestCheckResourceAttr(resourceName, "pod_id", podId),
					resource.TestCheckResourceAttr(resourceName, "node_id", NodeId),
					resource.TestCheckResourceAttr(resourceName, "interface", Interface),
				),
			},
			{
				Config: CreateAccOutofServiceFabricPathConfigWithOptionalValues(podId, NodeId, Interface),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutofServiceFabricPathExists(resourceName, &interface_blacklist_updated),

					resource.TestCheckResourceAttr(resourceName, "pod_id", podId),
					resource.TestCheckResourceAttr(resourceName, "node_id", NodeId),
					resource.TestCheckResourceAttr(resourceName, "interface", Interface),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),

					testAccCheckAciOutofServiceFabricPathIdEqual(&interface_blacklist_default, &interface_blacklist_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      CreateAccOutofServiceFabricPathConfigWithRequiredParams(randomValue, NodeId, Interface),
				ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
			},
			{
				Config:      CreateAccOutofServiceFabricPathConfigWithRequiredParams(podId, randomValue, Interface),
				ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
			},
			{
				Config: CreateAccOutofServiceFabricPathConfigWithRequiredParams(podId, NodeId, InterfaceUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutofServiceFabricPathExists(resourceName, &interface_blacklist_updated),
					resource.TestCheckResourceAttr(resourceName, "pod_id", podId),
					resource.TestCheckResourceAttr(resourceName, "node_id", NodeId),
					resource.TestCheckResourceAttr(resourceName, "interface", InterfaceUpdated),
					testAccCheckAciOutofServiceFabricPathIdNotEqual(&interface_blacklist_default, &interface_blacklist_updated),
				),
			},
		},
	})
}

func TestAccAciOutofServiceFabricPath_Negative(t *testing.T) {

	podId := "1"
	NodeId := "201"
	Interface := "eth1/1"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOutofServiceFabricPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccOutofServiceFabricPathConfig(podId, NodeId, Interface),
			},
			{
				Config:      CreateAccOutofServiceFabricPathUpdatedAttr(podId, NodeId, Interface, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccOutofServiceFabricPathUpdatedAttr(podId, NodeId, Interface, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccOutofServiceFabricPathConfig(podId, NodeId, Interface),
			},
		},
	})
}

func TestAccAciOutofServiceFabricPath_MultipleCreateDelete(t *testing.T) {

	podId := "1"
	NodeId := "201"
	Interface := "eth1/1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOutofServiceFabricPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccOutofServiceFabricPathConfigMultiple(podId, NodeId, Interface),
			},
		},
	})
}

func testAccCheckAciOutofServiceFabricPathExists(name string, interface_blacklist *models.OutofServiceFabricPath) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Interface Blacklist %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Interface Blacklist dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		interface_blacklistFound := models.OutofServiceFabricPathFromContainer(cont)
		if interface_blacklistFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Interface Blacklist %s not found", rs.Primary.ID)
		}
		*interface_blacklist = *interface_blacklistFound
		return nil
	}
}

func testAccCheckAciOutofServiceFabricPathDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing interface_blacklist destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_interface_blacklist" {
			cont, err := client.Get(rs.Primary.ID)
			interface_blacklist := models.OutofServiceFabricPathFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Interface Blacklist %s Still exists", interface_blacklist.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciOutofServiceFabricPathIdEqual(m1, m2 *models.OutofServiceFabricPath) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("interface_blacklist DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciOutofServiceFabricPathIdNotEqual(m1, m2 *models.OutofServiceFabricPath) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("interface_blacklist DNs are equal")
		}
		return nil
	}
}

func CreateOutofServiceFabricPathWithoutRequired(podId, NodeId, Interface, attrName string) string {
	fmt.Println("=== STEP  Basic: testing interface_blacklist creation without ", attrName)
	rBlock := ``
	switch attrName {
	case "pod_id":
		rBlock += `
	resource "aci_interface_blacklist" "test" {
	
	#	pod_id  = %s
  		node_id = %s
 		interface = "%s"
	}
		`
	case "node_id":
		rBlock += `
	resource "aci_interface_blacklist" "test" {
	
		pod_id  = %s
  	#	node_id = %s
 		interface = "%s"
	}
		`
	case "interface":
		rBlock += `
	resource "aci_interface_blacklist" "test" {
	
		pod_id  = %s
  		node_id = %s
 	#	interface = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, podId, NodeId, Interface)
}

func CreateAccOutofServiceFabricPathConfigWithRequiredParams(podId, NodeId, Interface string) string {
	fmt.Printf("=== STEP  testing interface_blacklist creation with pod_id %s, node_id %s and interface %s\n", podId, NodeId, Interface)
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		pod_id  = "%s"
  		node_id = "%s"
 		interface = "%s"
	}
	`, podId, NodeId, Interface)
	return resource
}

func CreateAccOutofServiceFabricPathConfig(podId, NodeId, Interface string) string {
	fmt.Println("=== STEP  testing interface_blacklist creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		pod_id  = %s
  		node_id = %s
 		interface = "%s"
	}
	`, podId, NodeId, Interface)
	return resource
}

func CreateAccOutofServiceFabricPathConfigMultiple(podId, NodeId, Interface string) string {
	fmt.Println("=== STEP  testing multiple interface_blacklist creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		pod_id  = %s
  		node_id = %s
 		interface = "%s${count.index+1}"
		count = 5
	}
	`, podId, NodeId, Interface[:len(Interface)-1])
	return resource
}

func CreateAccOutofServiceFabricPathConfigWithOptionalValues(podId, NodeId, Interface string) string {
	fmt.Println("=== STEP  Basic: testing interface_blacklist creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
		pod_id  = %s
  		node_id = %s
 		interface = "%s"
		annotation = "orchestrator:terraform_testacc"
		
	}
	`, podId, NodeId, Interface)

	return resource
}

func CreateAccOutofServiceFabricPathUpdatedAttr(podId, NodeId, Interface, attribute, value string) string {
	fmt.Printf("=== STEP  testing interface_blacklist attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		pod_id  = %s
  		node_id = %s
 		interface = "%s"
		%s = "%s"
	}
	`, podId, NodeId, Interface, attribute, value)
	return resource
}
