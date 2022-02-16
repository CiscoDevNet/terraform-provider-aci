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

func TestAccAciTACACSDestination_Basic(t *testing.T) {
	var tacacs_accounting_destination_default models.TACACSDestination
	var tacacs_accounting_destination_updated models.TACACSDestination
	resourceName := "aci_tacacs_accounting_destination.test"
	rNameUpdated := acctest.RandString(5)
	host := fmt.Sprintf("%s.com", acctest.RandString(5))
	hostUpdated := fmt.Sprintf("%s.com", acctest.RandString(5))
	tacacsGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateTACACSDestinationWithoutRequired(tacacsGroupName, host, "tacacs_accounting_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateTACACSDestinationWithoutRequired(tacacsGroupName, host, "host"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTACACSDestinationConfig(tacacsGroupName, host),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSDestinationExists(resourceName, &tacacs_accounting_destination_default),
					resource.TestCheckResourceAttr(resourceName, "tacacs_accounting_dn", fmt.Sprintf("uni/fabric/tacacsgroup-%s", tacacsGroupName)),
					resource.TestCheckResourceAttr(resourceName, "host", host),
					resource.TestCheckResourceAttr(resourceName, "port", "49"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "name", ""),
					resource.TestCheckResourceAttr(resourceName, "auth_protocol", "pap"),
				),
			},
			{
				Config: CreateAccTACACSDestinationConfigWithOptionalValues(tacacsGroupName, host),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSDestinationExists(resourceName, &tacacs_accounting_destination_updated),
					resource.TestCheckResourceAttr(resourceName, "tacacs_accounting_dn", fmt.Sprintf("uni/fabric/tacacsgroup-%s", tacacsGroupName)),
					resource.TestCheckResourceAttr(resourceName, "host", host),
					resource.TestCheckResourceAttr(resourceName, "port", "49"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_tacacs_accounting_destination"),
					resource.TestCheckResourceAttr(resourceName, "auth_protocol", "chap"),
					resource.TestCheckResourceAttr(resourceName, "key", "test_key"),
					resource.TestCheckResourceAttr(resourceName, "name", "test_name"),
					testAccCheckAciTACACSDestinationIdEqual(&tacacs_accounting_destination_default, &tacacs_accounting_destination_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key"},
			},
			{
				Config:      CreateAccTACACSDestinationRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccTACACSDestinationConfigWithRequiredParams(rNameUpdated, host, "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccTACACSDestinationConfigWithRequiredParams(rNameUpdated, host, "65536"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config: CreateAccTACACSDestinationConfigWithRequiredParams(rNameUpdated, host, "49"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSDestinationExists(resourceName, &tacacs_accounting_destination_updated),
					resource.TestCheckResourceAttr(resourceName, "tacacs_accounting_dn", fmt.Sprintf("uni/fabric/tacacsgroup-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "host", host),
					resource.TestCheckResourceAttr(resourceName, "port", "49"),
					testAccCheckAciTACACSDestinationIdNotEqual(&tacacs_accounting_destination_default, &tacacs_accounting_destination_updated),
				),
			},
			{
				Config: CreateAccTACACSDestinationConfig(tacacsGroupName, host),
			},
			{
				Config: CreateAccTACACSDestinationConfigWithRequiredParams(tacacsGroupName, hostUpdated, "49"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSDestinationExists(resourceName, &tacacs_accounting_destination_updated),
					resource.TestCheckResourceAttr(resourceName, "tacacs_accounting_dn", fmt.Sprintf("uni/fabric/tacacsgroup-%s", tacacsGroupName)),
					resource.TestCheckResourceAttr(resourceName, "host", hostUpdated),
					resource.TestCheckResourceAttr(resourceName, "port", "49"),
					testAccCheckAciTACACSDestinationIdNotEqual(&tacacs_accounting_destination_default, &tacacs_accounting_destination_updated),
				),
			},
			{
				Config: CreateAccTACACSDestinationConfig(tacacsGroupName, host),
			},
			{
				Config: CreateAccTACACSDestinationConfigWithRequiredParams(tacacsGroupName, host, "50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSDestinationExists(resourceName, &tacacs_accounting_destination_updated),
					resource.TestCheckResourceAttr(resourceName, "tacacs_accounting_dn", fmt.Sprintf("uni/fabric/tacacsgroup-%s", tacacsGroupName)),
					resource.TestCheckResourceAttr(resourceName, "port", "50"),
					resource.TestCheckResourceAttr(resourceName, "host", host),
					testAccCheckAciTACACSDestinationIdNotEqual(&tacacs_accounting_destination_default, &tacacs_accounting_destination_updated),
				),
			},
		},
	})
}

func TestAccAciTACACSDestination_Update(t *testing.T) {
	var tacacs_accounting_destination_default models.TACACSDestination
	var tacacs_accounting_destination_updated models.TACACSDestination
	resourceName := "aci_tacacs_accounting_destination.test"
	host := fmt.Sprintf("%s.com", acctest.RandString(5))
	tacacsGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSDestinationConfig(tacacsGroupName, host),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSDestinationExists(resourceName, &tacacs_accounting_destination_default),
				),
			},
			{
				Config: CreateAccTACACSDestinationUpdatedAttr(tacacsGroupName, host, "auth_protocol", "mschap"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSDestinationExists(resourceName, &tacacs_accounting_destination_updated),
					resource.TestCheckResourceAttr(resourceName, "auth_protocol", "mschap"),
					testAccCheckAciTACACSDestinationIdEqual(&tacacs_accounting_destination_default, &tacacs_accounting_destination_updated),
				),
			},
			{
				Config: CreateAccTACACSDestinationConfigWithRequiredParams(tacacsGroupName, host, "1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSDestinationExists(resourceName, &tacacs_accounting_destination_updated),
					resource.TestCheckResourceAttr(resourceName, "tacacs_accounting_dn", fmt.Sprintf("uni/fabric/tacacsgroup-%s", tacacsGroupName)),
					resource.TestCheckResourceAttr(resourceName, "host", host),
					resource.TestCheckResourceAttr(resourceName, "port", "1"),
					testAccCheckAciTACACSDestinationIdNotEqual(&tacacs_accounting_destination_default, &tacacs_accounting_destination_updated),
				),
			},
			{
				Config: CreateAccTACACSDestinationConfigWithRequiredParams(tacacsGroupName, host, "32000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSDestinationExists(resourceName, &tacacs_accounting_destination_updated),
					resource.TestCheckResourceAttr(resourceName, "tacacs_accounting_dn", fmt.Sprintf("uni/fabric/tacacsgroup-%s", tacacsGroupName)),
					resource.TestCheckResourceAttr(resourceName, "host", host),
					resource.TestCheckResourceAttr(resourceName, "port", "32000"),
					testAccCheckAciTACACSDestinationIdNotEqual(&tacacs_accounting_destination_default, &tacacs_accounting_destination_updated),
				),
			},
		},
	})
}

func TestAccAciTACACSDestination_Negative(t *testing.T) {
	host := fmt.Sprintf("%s.com", acctest.RandString(5))
	tacacsGroupName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSDestinationConfig(tacacsGroupName, host),
			},
			{
				Config:      CreateAccTACACSDestinationWithInValidParentDn(tacacsGroupName, host),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccTACACSDestinationUpdatedAttr(tacacsGroupName, host, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSDestinationUpdatedAttr(tacacsGroupName, host, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSDestinationUpdatedAttr(tacacsGroupName, host, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSDestinationUpdatedAttr(tacacsGroupName, host, "name", acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccTACACSDestinationUpdatedAttr(tacacsGroupName, host, "auth_protocol", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccTACACSDestinationUpdatedAttr(tacacsGroupName, host, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccTACACSDestinationConfig(tacacsGroupName, host),
			},
		},
	})
}

func TestAccAciTACACSDestination_MultipleCreateDelete(t *testing.T) {
	host := acctest.RandString(5)
	tacacsGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSDestinationConfigMultiple(tacacsGroupName, host),
			},
		},
	})
}

func testAccCheckAciTACACSDestinationExists(name string, tacacs_accounting_destination *models.TACACSDestination) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("TACACSDestination %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No TACACSDestination dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tacacs_accounting_destinationFound := models.TACACSDestinationFromContainer(cont)
		if tacacs_accounting_destinationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("TACACSDestination %s not found", rs.Primary.ID)
		}
		*tacacs_accounting_destination = *tacacs_accounting_destinationFound
		return nil
	}
}

func testAccCheckAciTACACSDestinationDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing tacacs_accounting_destination destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tacacs_accounting_destination" {
			cont, err := client.Get(rs.Primary.ID)
			tacacs_accounting_destination := models.TACACSDestinationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("TACACSDestination %s Still exists", tacacs_accounting_destination.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTACACSDestinationIdEqual(m1, m2 *models.TACACSDestination) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("tacacs_accounting_destination DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciTACACSDestinationIdNotEqual(m1, m2 *models.TACACSDestination) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("tacacs_accounting_destination DNs are equal")
		}
		return nil
	}
}

func CreateTACACSDestinationWithoutRequired(tacacsGroupName, host, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_accounting_destination creation without ", attrName)
	rBlock := `
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tacacs_accounting_dn":
		rBlock += `
	resource "aci_tacacs_accounting_destination" "test" {
	#	tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = "%s"	
	}
		`
	case "host":
		rBlock += `
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
	#	host  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, tacacsGroupName, host)
}

func CreateAccTACACSDestinationConfigWithRequiredParams(tacacsGroupName, host, port string) string {
	fmt.Printf("=== STEP  testing tacacs_accounting_destination creation with parent resource name %s, host %s and port %s\n", tacacsGroupName, host, port)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = "%s"
		port  = "%s"
	}
	`, tacacsGroupName, host, port)
	return resource
}

func CreateAccTACACSDestinationConfig(tacacsGroupName, host string) string {
	fmt.Println("=== STEP  testing tacacs_accounting_destination creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = "%s"
	}
	`, tacacsGroupName, host)
	return resource
}

func CreateAccTACACSDestinationConfigMultiple(tacacsGroupName, host string) string {
	fmt.Println("=== STEP  testing multiple tacacs_accounting_destination creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = "%s${count.index}.com"
		count = 5
	}
	`, tacacsGroupName, host)
	return resource
}

func CreateAccTACACSDestinationWithInValidParentDn(rName, host string) string {
	fmt.Println("=== STEP  Negative Case: testing tacacs_accounting_destination creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tenant.test.id
		host  = "%s"	
	}
	`, rName, host)
	return resource
}

func CreateAccTACACSDestinationConfigWithOptionalValues(tacacsGroupName, host string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_accounting_destination creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = "${aci_tacacs_accounting.test.id}"
		host  = "%s"
		port  = "49"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tacacs_accounting_destination"
		auth_protocol = "chap"
		key = "test_key"
		name = "test_name"
	}
	`, tacacsGroupName, host)

	return resource
}

func CreateAccTACACSDestinationRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing tacacs_accounting_destination updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_tacacs_accounting_destination" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tacacs_accounting_destination"
		auth_protocol = "chap"
		key = ""
		
	}
	`)

	return resource
}

func CreateAccTACACSDestinationUpdatedAttr(tacacsGroupName, host, attribute, value string) string {
	fmt.Printf("=== STEP  testing tacacs_accounting_destination attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = "%s"
		%s = "%s"
	}
	`, tacacsGroupName, host, attribute, value)
	return resource
}
