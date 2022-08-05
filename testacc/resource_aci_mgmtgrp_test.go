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

func TestAccAciManagedNodeConnectivityGroup_Basic(t *testing.T) {
	var managed_node_connectivity_group_default models.ManagedNodeConnectivityGroup
	var managed_node_connectivity_group_updated models.ManagedNodeConnectivityGroup
	resourceName := "aci_managed_node_connectivity_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciManagedNodeConnectivityGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateManagedNodeConnectivityGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccManagedNodeConnectivityGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciManagedNodeConnectivityGroupExists(resourceName, &managed_node_connectivity_group_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
				),
			},
			{
				Config: CreateAccManagedNodeConnectivityGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciManagedNodeConnectivityGroupExists(resourceName, &managed_node_connectivity_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					testAccCheckAciManagedNodeConnectivityGroupIdEqual(&managed_node_connectivity_group_default, &managed_node_connectivity_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccManagedNodeConnectivityGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccManagedNodeConnectivityGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccManagedNodeConnectivityGroupConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciManagedNodeConnectivityGroupExists(resourceName, &managed_node_connectivity_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciManagedNodeConnectivityGroupIdNotEqual(&managed_node_connectivity_group_default, &managed_node_connectivity_group_updated),
				),
			},
		},
	})
}

func TestAccAciManagedNodeConnectivityGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciManagedNodeConnectivityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccManagedNodeConnectivityGroupConfig(rName),
			},
			{
				Config:      CreateAccManagedNodeConnectivityGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccManagedNodeConnectivityGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccManagedNodeConnectivityGroupConfig(rName),
			},
		},
	})
}

func TestAccAciManagedNodeConnectivityGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciManagedNodeConnectivityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccManagedNodeConnectivityGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciManagedNodeConnectivityGroupExists(name string, managed_node_connectivity_group *models.ManagedNodeConnectivityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Managed Node Connectivity Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Managed Node Connectivity Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		managed_node_connectivity_groupFound := models.ManagedNodeConnectivityGroupFromContainer(cont)
		if managed_node_connectivity_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Managed Node Connectivity Group %s not found", rs.Primary.ID)
		}
		*managed_node_connectivity_group = *managed_node_connectivity_groupFound
		return nil
	}
}

func testAccCheckAciManagedNodeConnectivityGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing managed_node_connectivity_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_managed_node_connectivity_group" {
			cont, err := client.Get(rs.Primary.ID)
			managed_node_connectivity_group := models.ManagedNodeConnectivityGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Managed Node Connectivity Group %s Still exists", managed_node_connectivity_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciManagedNodeConnectivityGroupIdEqual(m1, m2 *models.ManagedNodeConnectivityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("managed_node_connectivity_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciManagedNodeConnectivityGroupIdNotEqual(m1, m2 *models.ManagedNodeConnectivityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("managed_node_connectivity_group DNs are equal")
		}
		return nil
	}
}

func CreateManagedNodeConnectivityGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing managed_node_connectivity_group creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_managed_node_connectivity_group" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccManagedNodeConnectivityGroupConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing managed_node_connectivity_group creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccManagedNodeConnectivityGroupConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing managed_node_connectivity_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccManagedNodeConnectivityGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing managed_node_connectivity_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccManagedNodeConnectivityGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple managed_node_connectivity_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccManagedNodeConnectivityGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing managed_node_connectivity_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
	
		name  = "%s"
		annotation = "orchestrator:terraform_testacc"
		
	}
	`, rName)

	return resource
}

func CreateAccManagedNodeConnectivityGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing managed_node_connectivity_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_managed_node_connectivity_group" "test" {
		annotation = "orchestrator:terraform_testacc"
		
	}
	`)

	return resource
}

func CreateAccManagedNodeConnectivityGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing managed_node_connectivity_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
