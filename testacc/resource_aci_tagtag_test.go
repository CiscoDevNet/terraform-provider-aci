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

func TestAccAciTag_Basic(t *testing.T) {
	var tag_default models.Tag
	var tag_updated models.Tag
	resourceName := "aci_tag.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	key := makeTestVariable(acctest.RandString(5))
	keyUpdated := makeTestVariable(acctest.RandString(5))
	value := makeTestVariable(acctest.RandString(5))
	valueUpdated := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTagDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateTagWithoutRequired(fvTenantName, key, value, "parent_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateTagWithoutRequired(fvTenantName, key, value, "key"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateTagWithoutRequired(fvTenantName, key, value, "value"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTagConfig(fvTenantName, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTagExists(resourceName, &tag_default),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "value", value),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      CreateAccTagRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTagConfigWithRequiredParams(rNameUpdated, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTagExists(resourceName, &tag_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					testAccCheckAciTagIdNotEqual(&tag_default, &tag_updated),
				),
			},
			{
				Config: CreateAccTagConfig(fvTenantName, key, value),
			},
			{
				Config: CreateAccTagConfigWithRequiredParams(rName, keyUpdated, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTagExists(resourceName, &tag_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "key", keyUpdated),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					testAccCheckAciTagIdNotEqual(&tag_default, &tag_updated),
				),
			},
			{
				Config: CreateAccTagConfigWithRequiredParams(rName, key, valueUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTagExists(resourceName, &tag_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "value", valueUpdated),
					testAccCheckAciTagIdNotEqual(&tag_default, &tag_updated),
				),
			},
		},
	})
}

func TestAccAciTag_MultipleCreateDelete(t *testing.T) {

	key := makeTestVariable(acctest.RandString(5))
	value := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTagConfigMultiple(fvTenantName, key, value),
			},
		},
	})
}

func testAccCheckAciTagExists(name string, tag *models.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Tag %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Tag dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tagFound := models.TagFromContainer(cont)
		if tagFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Tag %s not found", rs.Primary.ID)
		}
		*tag = *tagFound
		return nil
	}
}

func testAccCheckAciTagDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing tag destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tag" {
			cont, err := client.Get(rs.Primary.ID)
			tag := models.TagFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Tag %s Still exists", tag.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTagIdEqual(m1, m2 *models.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("tag DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciTagIdNotEqual(m1, m2 *models.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("tag DNs are equal")
		}
		return nil
	}
}

func CreateTagWithoutRequired(fvTenantName, key, value, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tag creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "parent_dn":
		rBlock += `
	resource "aci_tag" "test" {
	#	parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "%s"
	}
		`
	case "key":
		rBlock += `
	resource "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
	#	key  = "%s"
		value = "%s"	
	}
		`
	case "value":
		rBlock += `
	resource "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
	#	value = "%s"	
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, key, value)
}

func CreateAccTagConfigWithRequiredParams(fvTenantName, key, value string) string {
	fmt.Println("=== STEP  testing tag creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "%s"
	}
	`, fvTenantName, key, value)
	return resource
}

func CreateAccTagConfig(fvTenantName, key, value string) string {
	fmt.Println("=== STEP  testing tag creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "%s"
	}
	`, fvTenantName, key, value)
	return resource
}

func CreateAccTagConfigMultiple(fvTenantName, key, value string) string {
	fmt.Println("=== STEP  testing multiple tag creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s_${count.index}"
		value = "%s"
		count = 5
	}
	`, fvTenantName, key, value)
	return resource
}

func CreateAccTagWithInValidParentDn(rName, key, value string) string {
	fmt.Println("=== STEP  Negative Case: testing tag creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"	
		value = "%s"
	}
	`, rName, key, value)
	return resource
}

func CreateAccTagRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing tag updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_tag" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tag"
		
	}
	`)

	return resource
}
