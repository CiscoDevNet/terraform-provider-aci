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

func TestAccAciMgmtconnectivitypreference_Basic(t *testing.T) {
	var mgmt_preference_default models.Mgmtconnectivitypreference
	var mgmt_preference_updated models.Mgmtconnectivitypreference
	resourceName := "aci_mgmt_preference.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMgmtconnectivitypreferenceDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMgmtconnectivitypreferenceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMgmtconnectivitypreferenceExists(resourceName, &mgmt_preference_default),
					// all default values varies based on server
				),
			},
			{
				Config: CreateAccMgmtconnectivitypreferenceConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMgmtconnectivitypreferenceExists(resourceName, &mgmt_preference_updated),

					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_mgmt_preference"),

					resource.TestCheckResourceAttr(resourceName, "interface_pref", "ooband"),

					testAccCheckAciMgmtconnectivitypreferenceIdEqual(&mgmt_preference_default, &mgmt_preference_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccMgmtconnectivitypreferenceUpdatedAttr("interface_pref", "inband"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMgmtconnectivitypreferenceExists(resourceName, &mgmt_preference_updated),
					resource.TestCheckResourceAttr(resourceName, "interface_pref", "inband"),
					testAccCheckAciMgmtconnectivitypreferenceIdEqual(&mgmt_preference_default, &mgmt_preference_updated),
				),
			},
		},
	})
}

func TestAccAciMgmtconnectivitypreference_Negative(t *testing.T) {

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMgmtconnectivitypreferenceDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMgmtconnectivitypreferenceConfig(),
			},

			{
				Config:      CreateAccMgmtconnectivitypreferenceUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMgmtconnectivitypreferenceUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMgmtconnectivitypreferenceUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccMgmtconnectivitypreferenceUpdatedAttr("interface_pref", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccMgmtconnectivitypreferenceUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccMgmtconnectivitypreferenceConfig(),
			},
		},
	})
}

func testAccCheckAciMgmtconnectivitypreferenceExists(name string, mgmt_preference *models.Mgmtconnectivitypreference) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Mgmt Preference %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Mgmt Preference dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		mgmt_preferenceFound := models.MgmtconnectivitypreferenceFromContainer(cont)
		if mgmt_preferenceFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Mgmt Preference %s not found", rs.Primary.ID)
		}
		*mgmt_preference = *mgmt_preferenceFound
		return nil
	}
}

func testAccCheckAciMgmtconnectivitypreferenceDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing mgmt_preference destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_mgmt_preference" {
			cont, err := client.Get(rs.Primary.ID)
			mgmt_preference := models.MgmtconnectivitypreferenceFromContainer(cont)
			if err != nil {
				return fmt.Errorf("Mgmt Preference %s Still exists", mgmt_preference.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciMgmtconnectivitypreferenceIdEqual(m1, m2 *models.Mgmtconnectivitypreference) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("mgmt_preference DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciMgmtconnectivitypreferenceIdNotEqual(m1, m2 *models.Mgmtconnectivitypreference) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("mgmt_preference DNs are equal")
		}
		return nil
	}
}

func CreateAccMgmtconnectivitypreferenceConfig() string {
	fmt.Println("=== STEP  testing mgmt_preference creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_mgmt_preference" "test" {
	
	}
	`)
	return resource
}

func CreateAccMgmtconnectivitypreferenceConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing mgmt_preference creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_mgmt_preference" "test" {
	
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_mgmt_preference"
		interface_pref = "ooband"
		
	}
	`)

	return resource
}

func CreateAccMgmtconnectivitypreferenceUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing mgmt_preference attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_mgmt_preference" "test" {
	
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
