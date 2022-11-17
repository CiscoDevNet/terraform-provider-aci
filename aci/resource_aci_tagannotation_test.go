package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciAnnotation_Basic(t *testing.T) {
	tag_annotation_key := acctest.RandString(5)
	tag_annotation_value := acctest.RandString(5)
	tenant_name := acctest.RandString(5)
	description := "annotation created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAnnotationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAnnotationConfig_basic(tenant_name, tag_annotation_key, tag_annotation_value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnnotationExists("aci_annotation.fooannotation", &annotation),
					testAccCheckAciAnnotationAttributes(tenant_name, tag_annotation_key, tag_annotation_value),
				),
			},
		},
	})
}

func TestAccAciAnnotation_Update(t *testing.T) {
	tag_annotation_key := acctest.RandString(5)
	tag_annotation_value := acctest.RandString(5)
	tag_annotation_parent_dn := "uni/tn-common"
	description := "annotation created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAnnotationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAnnotationConfig_basic(tenant_name, tag_annotation_key, tag_annotation_value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnnotationExists("aci_annotation.fooannotation", &annotation),
					testAccCheckAciAnnotationAttributes(tenant_name, tag_annotation_key, tag_annotation_value),
				),
			},
			{
				Config: testAccCheckAciAnnotationConfig_basic(tenant_name, tag_annotation_key, tag_annotation_value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnnotationExists("aci_annotation.fooannotation", &annotation),
					testAccCheckAciAnnotationAttributes(tenant_name, tag_annotation_key, tag_annotation_value),
				),
			},
		},
	})
}

func testAccCheckAciAnnotationConfig_basic(tenant_name, key, value string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.foological_node_profile.id
	}

	resource "aci_annotation" "fooannotation" {
		key 		= "%s"
		value       = "%s"
		description = "annotation created while acceptance testing"
		parent_dn = aci_tenant.footenant.id
	}

	`, tenant_name, key, value)
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

func testAccCheckAciAnnotationAttributes(tenant_name, key, value string, annotation *models.Annotation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if tag_annotation_key != GetMOName(annotation.DistinguishedName) {
			return fmt.Errorf("Bad tag_annotation %s", GetMOName(annotation.DistinguishedName))
		}

		if tenant_name != GetMOName(GetParentDn(annotation.DistinguishedName)) {
			return fmt.Errorf(" Bad tenant_name %s", GetMOName(GetParentDn(annotation.DistinguishedName)))
		}
		if value != annotation.Value {
			return fmt.Errorf("Bad annotation value %s", annotation.Value)
		}
		return nil
	}
}
