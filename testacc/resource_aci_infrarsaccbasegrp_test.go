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

func TestAccAciAccessGroup_Basic(t *testing.T) {
	var access_group_default models.AccessAccessGroup
	var access_group_updated models.AccessAccessGroup
	resourceName := "aci_access_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	infraAccPortPName := makeTestVariable(acctest.RandString(5))
	infraHPortSName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccessGroupWithoutRequired(infraAccPortPName, infraHPortSName, "access_port_selector_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAccessGroupConfig(infraAccPortPName, infraHPortSName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessGroupExists(resourceName, &access_group_default),
					resource.TestCheckResourceAttr(resourceName, "access_port_selector_dn", fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-ALL", infraAccPortPName, infraHPortSName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "fex_id", "101"),
					resource.TestCheckResourceAttr(resourceName, "tdn", ""),
				),
			},
			{
				Config: CreateAccAccessGroupConfigWithOptionalValues(infraAccPortPName, infraHPortSName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessGroupExists(resourceName, &access_group_updated),
					resource.TestCheckResourceAttr(resourceName, "access_port_selector_dn", fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-ALL", infraAccPortPName, infraHPortSName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "fex_id", "199"),

					resource.TestCheckResourceAttr(resourceName, "tdn", "uni/infra/fexprof-acctest_fex/fexbundle-acctest_fex"),

					testAccCheckAciAccessGroupIdEqual(&access_group_default, &access_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccAccessGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAccessGroupConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessGroupExists(resourceName, &access_group_updated),
					resource.TestCheckResourceAttr(resourceName, "access_port_selector_dn", fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-ALL", rName, rNameUpdated)),
					testAccCheckAciAccessGroupIdNotEqual(&access_group_default, &access_group_updated),
				),
			},
			{
				Config: CreateAccAccessGroupConfig(infraAccPortPName, infraHPortSName),
			},
		},
	})
}

func TestAccAciAccessGroup_Update(t *testing.T) {
	var access_group_default models.AccessAccessGroup
	var access_group_updated models.AccessAccessGroup
	resourceName := "aci_access_group.test"
	infraAccPortPName := makeTestVariable(acctest.RandString(5))
	infraHPortSName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAccessGroupConfig(infraAccPortPName, infraHPortSName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessGroupExists(resourceName, &access_group_default),
				),
			},
			{
				Config: CreateAccAccessGroupUpdatedAttr(infraAccPortPName, infraHPortSName, "fex_id", "150"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessGroupExists(resourceName, &access_group_updated),
					resource.TestCheckResourceAttr(resourceName, "fex_id", "150"),
					testAccCheckAciAccessGroupIdEqual(&access_group_default, &access_group_updated),
				),
			},
		},
	})
}

func TestAccAciAccessGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	infraAccPortPName := makeTestVariable(acctest.RandString(5))
	infraHPortSName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAccessGroupConfig(infraAccPortPName, infraHPortSName),
			},
			{
				Config:      CreateAccAccessGroupWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccAccessGroupWithInValidTdn(rName),
				ExpectError: regexp.MustCompile(`Invalid target DN`),
			},
			{
				Config:      CreateAccAccessGroupUpdatedAttr(infraAccPortPName, infraHPortSName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccAccessGroupUpdatedAttr(infraAccPortPName, infraHPortSName, "fex_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccAccessGroupUpdatedAttr(infraAccPortPName, infraHPortSName, "fex_id", "100"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccAccessGroupUpdatedAttr(infraAccPortPName, infraHPortSName, "fex_id", "200"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccAccessGroupUpdatedAttr(infraAccPortPName, infraHPortSName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccAccessGroupConfig(infraAccPortPName, infraHPortSName),
			},
		},
	})
}

func TestAccAciAccessGroup_MultipleCreateDelete(t *testing.T) {
	infraAccPortPName := makeTestVariable(acctest.RandString(5))
	infraHPortSName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAccessGroupConfigMultiple(infraAccPortPName, infraHPortSName),
			},
		},
	})
}

func testAccCheckAciAccessGroupExists(name string, access_group *models.AccessAccessGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Access Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Access Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		access_groupFound := models.AccessAccessGroupFromContainer(cont)
		if access_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Access Group %s not found", rs.Primary.ID)
		}
		*access_group = *access_groupFound
		return nil
	}
}

func testAccCheckAciAccessGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing access_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_access_group" {
			cont, err := client.Get(rs.Primary.ID)
			access_group := models.AccessAccessGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Access Group %s Still exists", access_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciAccessGroupIdEqual(m1, m2 *models.AccessAccessGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("access_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciAccessGroupIdNotEqual(m1, m2 *models.AccessAccessGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("access_group DNs are equal")
		}
		return nil
	}
}

func CreateAccessGroupWithoutRequired(infraAccPortPName, infraHPortSName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing access_group creation without ", attrName)
	rBlock := `
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
		access_port_selector_type = "ALL"
	}
	
	`
	switch attrName {
	case "access_port_selector_dn":
		rBlock += `
	resource "aci_access_group" "test" {
	#	access_port_selector_dn  = aci_access_port_selector.test.id
	
	}
		`

	}
	return fmt.Sprintf(rBlock, infraAccPortPName, infraHPortSName)
}

func CreateAccAccessGroupConfigWithRequiredParams(infraAccPortPName, infraHPortSName string) string {
	fmt.Println("=== STEP  testing access_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
		access_port_selector_type = "ALL"
	}
	
	resource "aci_access_group" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
	}
	`, infraAccPortPName, infraHPortSName)
	return resource
}
func CreateAccAccessGroupConfig(infraAccPortPName, infraHPortSName string) string {
	fmt.Println("=== STEP  testing access_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
		access_port_selector_type = "ALL"
	}
	
	resource "aci_access_group" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
	}
	`, infraAccPortPName, infraHPortSName)
	return resource
}

func CreateAccAccessGroupConfigMultiple(infraAccPortPName, infraHPortSName string) string {
	fmt.Println("=== STEP  testing multiple access_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	resource "aci_fex_profile" "test" {
		name        = "acctest_fex"
	}
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name            = "acctest_fex"
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
		access_port_selector_type = "ALL"
	}
	
	resource "aci_access_group" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
		count = 5
	}
	`, infraAccPortPName, infraHPortSName)
	return resource
}

func CreateAccAccessGroupWithInValidTdn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing access_group creation with invalid tdn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	  
	 resource "aci_leaf_interface_profile" "test" {
		 name 		= "acctest_test"
		  
	 }
	  
	resource "aci_access_port_selector" "test" {
		name 		= "acctest_test"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
		access_port_selector_type = "ALL"
	}
	  
	 resource "aci_access_group" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
		tdn = aci_tenant.test.id
	}
	`, rName)
	return resource
}

func CreateAccAccessGroupWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing access_group creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_access_group" "test" {
		access_port_selector_dn  = aci_tenant.test.id	
	}
	`, rName)
	return resource
}

func CreateAccAccessGroupConfigWithOptionalValues(infraAccPortPName, infraHPortSName string) string {
	fmt.Println("=== STEP  Basic: testing access_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
		access_port_selector_type = "ALL"
	}

	resource "aci_fex_profile" "test" {
		name        = "acctest_fex"
	}
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name            = "acctest_fex"
	}
	
	resource "aci_access_group" "test" {
		access_port_selector_dn  = "${aci_access_port_selector.test.id}"
		annotation = "orchestrator:terraform_testacc"
		fex_id = "199"
		tdn = aci_fex_bundle_group.test.id
		
	}
	`, infraAccPortPName, infraHPortSName)

	return resource
}

func CreateAccAccessGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing access_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_access_group" "test" {
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_access_group"
		fex_id = "102"
		tdn = ""
	}
	`)

	return resource
}

func CreateAccAccessGroupUpdatedAttr(infraAccPortPName, infraHPortSName, attribute, value string) string {
	fmt.Printf("=== STEP  testing access_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}

	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
		access_port_selector_type = "ALL"
	}

	resource "aci_fex_profile" "test" {
		name        = "acctest_fex"
	}

	resource "aci_access_group" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
		%s = "%s"
	}
	`, infraAccPortPName, infraHPortSName, attribute, value)
	return resource
}
