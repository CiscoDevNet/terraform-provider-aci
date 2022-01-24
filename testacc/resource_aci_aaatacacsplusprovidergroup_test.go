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

func TestAccAciTACACSPlusProviderGroup_Basic(t *testing.T) {
	var tacacs_provider_group_default models.TACACSPlusProviderGroup
	var tacacs_provider_group_updated models.TACACSPlusProviderGroup
	resourceName := "aci_tacacs_provider_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSPlusProviderGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateTACACSPlusProviderGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTACACSPlusProviderGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSPlusProviderGroupExists(resourceName, &tacacs_provider_group_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccTACACSPlusProviderGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSPlusProviderGroupExists(resourceName, &tacacs_provider_group_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_tacacs_provider_group"),
					testAccCheckAciTACACSPlusProviderGroupIdEqual(&tacacs_provider_group_default, &tacacs_provider_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccTACACSPlusProviderGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccTACACSPlusProviderGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTACACSPlusProviderGroupConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSPlusProviderGroupExists(resourceName, &tacacs_provider_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciTACACSPlusProviderGroupIdNotEqual(&tacacs_provider_group_default, &tacacs_provider_group_updated),
				),
			},
		},
	})
}

func TestAccAciTACACSPlusProviderGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSPlusProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSPlusProviderGroupConfig(rName),
			},

			{
				Config:      CreateAccTACACSPlusProviderGroupUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSPlusProviderGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSPlusProviderGroupUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccTACACSPlusProviderGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccTACACSPlusProviderGroupConfig(rName),
			},
		},
	})
}

func TestAccAciTACACSPlusProviderGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSPlusProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSPlusProviderGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciTACACSPlusProviderGroupExists(name string, tacacs_provider_group *models.TACACSPlusProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("TACACS Plus Provider Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No TACACS Plus Provider Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tacacs_provider_groupFound := models.TACACSPlusProviderGroupFromContainer(cont)
		if tacacs_provider_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("TACACS Plus Provider Group %s not found", rs.Primary.ID)
		}
		*tacacs_provider_group = *tacacs_provider_groupFound
		return nil
	}
}

func testAccCheckAciTACACSPlusProviderGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing tacacs_provider_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tacacs_provider_group" {
			cont, err := client.Get(rs.Primary.ID)
			tacacs_provider_group := models.TACACSPlusProviderGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("TACACS Plus Provider Group %s Still exists", tacacs_provider_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTACACSPlusProviderGroupIdEqual(m1, m2 *models.TACACSPlusProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("tacacs_provider_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciTACACSPlusProviderGroupIdNotEqual(m1, m2 *models.TACACSPlusProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("tacacs_provider_group DNs are equal")
		}
		return nil
	}
}

func CreateTACACSPlusProviderGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_provider_group creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_tacacs_provider_group" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccTACACSPlusProviderGroupConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing tacacs_provider_group creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccTACACSPlusProviderGroupConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing tacacs_provider_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccTACACSPlusProviderGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing tacacs_provider_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccTACACSPlusProviderGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple tacacs_provider_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider_group" "test" {
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccTACACSPlusProviderGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_provider_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider_group" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tacacs_provider_group"
		
	}
	`, rName)

	return resource
}

func CreateAccTACACSPlusProviderGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing tacacs_provider_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_tacacs_provider_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tacacs_provider_group"
		
	}
	`)

	return resource
}

func CreateAccTACACSPlusProviderGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing tacacs_provider_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}