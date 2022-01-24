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

func TestAccAciFexBundleGroup_Basic(t *testing.T) {
	var fex_bundle_group_default models.FexBundleGroup
	var fex_bundle_group_updated models.FexBundleGroup
	resourceName := "aci_fex_bundle_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	infraFexPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFexBundleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateFexBundleGroupWithoutRequired(infraFexPName, rName, "fex_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateFexBundleGroupWithoutRequired(infraFexPName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFexBundleGroupConfig(infraFexPName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFexBundleGroupExists(resourceName, &fex_bundle_group_default),
					resource.TestCheckResourceAttr(resourceName, "fex_profile_dn", fmt.Sprintf("uni/infra/fexprof-%s", infraFexPName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccFexBundleGroupConfigWithOptionalValues(infraFexPName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFexBundleGroupExists(resourceName, &fex_bundle_group_updated),
					resource.TestCheckResourceAttr(resourceName, "fex_profile_dn", fmt.Sprintf("uni/infra/fexprof-%s", infraFexPName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_fex_bundle_group"),

					testAccCheckAciFexBundleGroupIdEqual(&fex_bundle_group_default, &fex_bundle_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccFexBundleGroupConfigUpdatedName(infraFexPName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccFexBundleGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFexBundleGroupConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFexBundleGroupExists(resourceName, &fex_bundle_group_updated),
					resource.TestCheckResourceAttr(resourceName, "fex_profile_dn", fmt.Sprintf("uni/infra/fexprof-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciFexBundleGroupIdNotEqual(&fex_bundle_group_default, &fex_bundle_group_updated),
				),
			},
			{
				Config: CreateAccFexBundleGroupConfig(infraFexPName, rName),
			},
			{
				Config: CreateAccFexBundleGroupConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFexBundleGroupExists(resourceName, &fex_bundle_group_updated),
					resource.TestCheckResourceAttr(resourceName, "fex_profile_dn", fmt.Sprintf("uni/infra/fexprof-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciFexBundleGroupIdNotEqual(&fex_bundle_group_default, &fex_bundle_group_updated),
				),
			},
		},
	})
}

func TestAccAciFexBundleGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	infraFexPName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFexBundleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFexBundleGroupConfig(infraFexPName, rName),
			},
			{
				Config:      CreateAccFexBundleGroupWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccFexBundleGroupUpdatedAttr(infraFexPName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFexBundleGroupUpdatedAttr(infraFexPName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFexBundleGroupUpdatedAttr(infraFexPName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccFexBundleGroupUpdatedAttr(infraFexPName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFexBundleGroupConfig(infraFexPName, rName),
			},
		},
	})
}

func TestAccAciFexBundleGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	infraFexPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFexBundleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFexBundleGroupConfigMultiple(infraFexPName, rName),
			},
		},
	})
}

func testAccCheckAciFexBundleGroupExists(name string, fex_bundle_group *models.FexBundleGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fex Bundle Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fex Bundle Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fex_bundle_groupFound := models.FexBundleGroupFromContainer(cont)
		if fex_bundle_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fex Bundle Group %s not found", rs.Primary.ID)
		}
		*fex_bundle_group = *fex_bundle_groupFound
		return nil
	}
}

func testAccCheckAciFexBundleGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing fex_bundle_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_fex_bundle_group" {
			cont, err := client.Get(rs.Primary.ID)
			fex_bundle_group := models.FexBundleGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fex Bundle Group %s Still exists", fex_bundle_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFexBundleGroupIdEqual(m1, m2 *models.FexBundleGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("fex_bundle_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFexBundleGroupIdNotEqual(m1, m2 *models.FexBundleGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("fex_bundle_group DNs are equal")
		}
		return nil
	}
}

func CreateFexBundleGroupWithoutRequired(infraFexPName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fex_bundle_group creation without ", attrName)
	rBlock := `
	
	resource "aci_fex_profile" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "fex_profile_dn":
		rBlock += `
	resource "aci_fex_bundle_group" "test" {
	#	fex_profile_dn  = aci_fex_profile.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, infraFexPName, rName)
}

func CreateAccFexBundleGroupConfigWithRequiredParams(infraFexPName, rName string) string {
	fmt.Println("=== STEP  testing fex_bundle_group creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = "%s"
	}
	`, infraFexPName, rName)
	return resource
}
func CreateAccFexBundleGroupConfigUpdatedName(infraFexPName, rName string) string {
	fmt.Println("=== STEP  testing fex_bundle_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = "%s"
	}
	`, infraFexPName, rName)
	return resource
}

func CreateAccFexBundleGroupConfig(infraFexPName, rName string) string {
	fmt.Println("=== STEP  testing fex_bundle_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = "%s"
	}
	`, infraFexPName, rName)
	return resource
}

func CreateAccFexBundleGroupConfigMultiple(infraFexPName, rName string) string {
	fmt.Println("=== STEP  testing multiple fex_bundle_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, infraFexPName, rName)
	return resource
}

func CreateAccFexBundleGroupWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing fex_bundle_group creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccFexBundleGroupConfigWithOptionalValues(infraFexPName, rName string) string {
	fmt.Println("=== STEP  Basic: testing fex_bundle_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = "${aci_fex_profile.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fex_bundle_group"
		
	}
	`, infraFexPName, rName)

	return resource
}

func CreateAccFexBundleGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing fex_bundle_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_fex_bundle_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fex_bundle_group"
		
	}
	`)

	return resource
}

func CreateAccFexBundleGroupUpdatedAttr(infraFexPName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing fex_bundle_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = "%s"
		%s = "%s"
	}
	`, infraFexPName, rName, attribute, value)
	return resource
}
