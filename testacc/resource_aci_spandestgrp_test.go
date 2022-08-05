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

func TestAccAciSPANDestinationGroup_Basic(t *testing.T) {
	var span_destination_group_default models.SPANDestinationGroup
	var span_destination_group_updated models.SPANDestinationGroup
	resourceName := "aci_span_destination_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSPANDestinationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateSPANDestinationGroupWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSPANDestinationGroupWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSPANDestinationGroupConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANDestinationGroupExists(resourceName, &span_destination_group_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccSPANDestinationGroupConfigWithOptionalValues(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANDestinationGroupExists(resourceName, &span_destination_group_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_span_destination_group"),

					testAccCheckAciSPANDestinationGroupIdEqual(&span_destination_group_default, &span_destination_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSPANDestinationGroupConfigUpdatedName(fvTenantName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccSPANDestinationGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSPANDestinationGroupConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANDestinationGroupExists(resourceName, &span_destination_group_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciSPANDestinationGroupIdNotEqual(&span_destination_group_default, &span_destination_group_updated),
				),
			},
			{
				Config: CreateAccSPANDestinationGroupConfig(fvTenantName, rName),
			},
			{
				Config: CreateAccSPANDestinationGroupConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANDestinationGroupExists(resourceName, &span_destination_group_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciSPANDestinationGroupIdNotEqual(&span_destination_group_default, &span_destination_group_updated),
				),
			},
		},
	})
}

func TestAccAciSPANDestinationGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSPANDestinationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSPANDestinationGroupConfig(fvTenantName, rName),
			},
			{
				Config:      CreateAccSPANDestinationGroupWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccSPANDestinationGroupUpdatedAttr(fvTenantName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSPANDestinationGroupUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSPANDestinationGroupUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccSPANDestinationGroupUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSPANDestinationGroupConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciSPANDestinationGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSPANDestinationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSPANDestinationGroupConfigMultiple(fvTenantName, rName),
			},
		},
	})
}

func testAccCheckAciSPANDestinationGroupExists(name string, span_destination_group *models.SPANDestinationGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("SPAN Destination Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SPAN Destination Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		span_destination_groupFound := models.SPANDestinationGroupFromContainer(cont)
		if span_destination_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("SPAN Destination Group %s not found", rs.Primary.ID)
		}
		*span_destination_group = *span_destination_groupFound
		return nil
	}
}

func testAccCheckAciSPANDestinationGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing span_destination_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_span_destination_group" {
			cont, err := client.Get(rs.Primary.ID)
			span_destination_group := models.SPANDestinationGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("SPAN Destination Group %s Still exists", span_destination_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSPANDestinationGroupIdEqual(m1, m2 *models.SPANDestinationGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("span_destination_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSPANDestinationGroupIdNotEqual(m1, m2 *models.SPANDestinationGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("span_destination_group DNs are equal")
		}
		return nil
	}
}

func CreateSPANDestinationGroupWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing span_destination_group creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_span_destination_group" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_span_destination_group" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccSPANDestinationGroupConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing span_destination_group creation with parent resource name %s and resource name %s\n", fvTenantName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_destination_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}
func CreateAccSPANDestinationGroupConfigUpdatedName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing span_destination_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_destination_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccSPANDestinationGroupConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing span_destination_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_destination_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccSPANDestinationGroupConfigMultiple(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing multiple span_destination_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_destination_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccSPANDestinationGroupWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing span_destination_group creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		tenant_dn  = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_span_destination_group" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccSPANDestinationGroupConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing span_destination_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_destination_group" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_span_destination_group"
		
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccSPANDestinationGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing span_destination_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_span_destination_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_span_destination_group"
		
	}
	`)

	return resource
}

func CreateAccSPANDestinationGroupUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing span_destination_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_destination_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
