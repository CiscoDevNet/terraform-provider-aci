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

func TestAccAciAnnotation_Basic(t *testing.T) {
	var annotation_default models.Annotation
	var annotation_updated models.Annotation
	resourceName := "aci_annotation.test"
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
		CheckDestroy:      testAccCheckAciAnnotationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAnnotationWithoutRequired(fvTenantName, key, value, "parent_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAnnotationWithoutRequired(fvTenantName, key, value, "key"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAnnotationWithoutRequired(fvTenantName, key, value, "value"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAnnotationConfig(fvTenantName, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnnotationExists(resourceName, &annotation_default),
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
				Config:      CreateAccAnnotationRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAnnotationConfigWithRequiredParams(rNameUpdated, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnnotationExists(resourceName, &annotation_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					testAccCheckAciAnnotationIdNotEqual(&annotation_default, &annotation_updated),
				),
			},
			{
				Config: CreateAccAnnotationConfig(fvTenantName, key, value),
			},
			{
				Config: CreateAccAnnotationConfigWithRequiredParams(rName, keyUpdated, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnnotationExists(resourceName, &annotation_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "key", keyUpdated),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					testAccCheckAciAnnotationIdNotEqual(&annotation_default, &annotation_updated),
				),
			},
			{
				Config: CreateAccAnnotationConfigWithRequiredParams(rName, key, valueUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnnotationExists(resourceName, &annotation_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "value", valueUpdated),
					testAccCheckAciAnnotationIdNotEqual(&annotation_default, &annotation_updated),
				),
			},
		},
	})
}

func TestAccAciAnnotation_MultipleCreateDelete(t *testing.T) {

	key := makeTestVariable(acctest.RandString(5))
	value := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAnnotationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAnnotationConfigMultiple(fvTenantName, key, value),
			},
		},
	})
}

func testAccCheckAciAnnotationExists(name string, annotation *models.Annotation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Annotation %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Annotation dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		annotationFound := models.AnnotationFromContainer(cont)
		if annotationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Annotation %s not found", rs.Primary.ID)
		}
		*annotation = *annotationFound
		return nil
	}
}

func testAccCheckAciAnnotationDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing annotation destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_annotation" {
			cont, err := client.Get(rs.Primary.ID)
			annotation := models.AnnotationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Annotation %s Still exists", annotation.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciAnnotationIdEqual(m1, m2 *models.Annotation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("annotation DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciAnnotationIdNotEqual(m1, m2 *models.Annotation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("annotation DNs are equal")
		}
		return nil
	}
}

func CreateAnnotationWithoutRequired(fvTenantName, key, value, attrName string) string {
	fmt.Println("=== STEP  Basic: testing annotation creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "parent_dn":
		rBlock += `
	resource "aci_annotation" "test" {
	#	parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "%s"
	}
		`
	case "key":
		rBlock += `
	resource "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
	#	key  = "%s"
		value = "%s"	
	}
		`
	case "value":
		rBlock += `
	resource "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
	#	value = "%s"	
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, key, value)
}

func CreateAccAnnotationConfigWithRequiredParams(fvTenantName, key, value string) string {
	fmt.Println("=== STEP  testing annotation creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "%s"
	}
	`, fvTenantName, key, value)
	return resource
}

func CreateAccAnnotationConfig(fvTenantName, key, value string) string {
	fmt.Println("=== STEP  testing annotation creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "%s"
	}
	`, fvTenantName, key, value)
	return resource
}

func CreateAccAnnotationConfigMultiple(fvTenantName, key, value string) string {
	fmt.Println("=== STEP  testing multiple annotation creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s_${count.index}"
		value = "%s"
		count = 5
	}
	`, fvTenantName, key, value)
	return resource
}

func CreateAccAnnotationWithInValidParentDn(rName, key, value string) string {
	fmt.Println("=== STEP  Negative Case: testing annotation creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"	
		value = "%s"
	}
	`, rName, key, value)
	return resource
}

func CreateAccAnnotationRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing annotation updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_annotation" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_annotation"
		
	}
	`)

	return resource
}
