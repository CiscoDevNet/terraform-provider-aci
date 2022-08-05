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

func TestAccAciSPANSourceGroup_Basic(t *testing.T) {
	var span_source_group_default models.SPANSourceGroup
	var span_source_group_updated models.SPANSourceGroup
	resourceName := "aci_span_source_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSPANSourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateSPANSourceGroupWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSPANSourceGroupWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSPANSourceGroupConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourceGroupExists(resourceName, &span_source_group_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "enabled"),
				),
			},
			{
				Config: CreateAccSPANSourceGroupConfigWithOptionalValues(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourceGroupExists(resourceName, &span_source_group_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_span_source_group"),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "disabled"),
					testAccCheckAciSPANSourceGroupIdEqual(&span_source_group_default, &span_source_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSPANSourceGroupConfigUpdatedName(fvTenantName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccSPANSourceGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSPANSourceGroupConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourceGroupExists(resourceName, &span_source_group_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciSPANSourceGroupIdNotEqual(&span_source_group_default, &span_source_group_updated),
				),
			},
			{
				Config: CreateAccSPANSourceGroupConfig(fvTenantName, rName),
			},
			{
				Config: CreateAccSPANSourceGroupConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourceGroupExists(resourceName, &span_source_group_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciSPANSourceGroupIdNotEqual(&span_source_group_default, &span_source_group_updated),
				),
			},
		},
	})
}

func TestAccAciSPANSourceGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSPANSourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSPANSourceGroupConfig(fvTenantName, rName),
			},
			{
				Config:      CreateAccSPANSourceGroupWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccSPANSourceGroupUpdatedAttr(fvTenantName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSPANSourceGroupUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSPANSourceGroupUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccSPANSourceGroupUpdatedAttr(fvTenantName, rName, "admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccSPANSourceGroupUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSPANSourceGroupConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciSPANSourceGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSPANSourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSPANSourceGroupConfigMultiple(fvTenantName, rName),
			},
		},
	})
}

func testAccCheckAciSPANSourceGroupExists(name string, span_source_group *models.SPANSourceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("SPAN Source Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SPAN Source Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		span_source_groupFound := models.SPANSourceGroupFromContainer(cont)
		if span_source_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("SPAN Source Group %s not found", rs.Primary.ID)
		}
		*span_source_group = *span_source_groupFound
		return nil
	}
}

func testAccCheckAciSPANSourceGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing span_source_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_span_source_group" {
			cont, err := client.Get(rs.Primary.ID)
			span_source_group := models.SPANSourceGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("SPAN Source Group %s Still exists", span_source_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSPANSourceGroupIdEqual(m1, m2 *models.SPANSourceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("span_source_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSPANSourceGroupIdNotEqual(m1, m2 *models.SPANSourceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("span_source_group DNs are equal")
		}
		return nil
	}
}

func CreateSPANSourceGroupWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing span_source_group creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_span_source_group" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccSPANSourceGroupConfigWithRequiredParams(prName, rName string) string {
	fmt.Printf("=== STEP  testing span_source_group creation with parent resource name %s and resource name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, prName, rName)
	return resource
}
func CreateAccSPANSourceGroupConfigUpdatedName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing span_source_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccSPANSourceGroupConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing span_source_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccSPANSourceGroupConfigMultiple(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing multiple span_source_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccSPANSourceGroupWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing span_source_group creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_span_source_group" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccSPANSourceGroupConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing span_source_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_span_source_group"
		admin_st = "disabled"
		
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccSPANSourceGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing span_source_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_span_source_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_span_source_group"
		admin_st = "disabled"
		
	}
	`)

	return resource
}

func CreateAccSPANSourceGroupUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing span_source_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
