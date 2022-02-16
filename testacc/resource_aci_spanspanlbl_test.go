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

func TestAccAciSPANSourcedestinationGroupMatchLabel_Basic(t *testing.T) {
	var span_sourcedestination_group_match_label_default models.SPANSourcedestinationGroupMatchLabel
	var span_sourcedestination_group_match_label_updated models.SPANSourcedestinationGroupMatchLabel
	resourceName := "aci_span_sourcedestination_group_match_label.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	spanSrcGrpName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSPANSourcedestinationGroupMatchLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateSPANSourcedestinationGroupMatchLabelWithoutRequired(fvTenantName, spanSrcGrpName, rName, "span_source_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSPANSourcedestinationGroupMatchLabelWithoutRequired(fvTenantName, spanSrcGrpName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelConfig(fvTenantName, spanSrcGrpName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists(resourceName, &span_sourcedestination_group_match_label_default),
					resource.TestCheckResourceAttr(resourceName, "span_source_group_dn", fmt.Sprintf("uni/tn-%s/srcgrp-%s", fvTenantName, spanSrcGrpName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttrSet(resourceName, "tag"),
				),
			},
			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelConfigWithOptionalValues(fvTenantName, spanSrcGrpName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists(resourceName, &span_sourcedestination_group_match_label_updated),
					resource.TestCheckResourceAttr(resourceName, "span_source_group_dn", fmt.Sprintf("uni/tn-%s/srcgrp-%s", fvTenantName, spanSrcGrpName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_span_sourcedestination_group_match_label"),
					resource.TestCheckResourceAttr(resourceName, "tag", "blue"),
					testAccCheckAciSPANSourcedestinationGroupMatchLabelIdEqual(&span_sourcedestination_group_match_label_default, &span_sourcedestination_group_match_label_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSPANSourcedestinationGroupMatchLabelConfigUpdatedName(fvTenantName, spanSrcGrpName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccSPANSourcedestinationGroupMatchLabelRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists(resourceName, &span_sourcedestination_group_match_label_updated),
					resource.TestCheckResourceAttr(resourceName, "span_source_group_dn", fmt.Sprintf("uni/tn-%s/srcgrp-%s", rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciSPANSourcedestinationGroupMatchLabelIdNotEqual(&span_sourcedestination_group_match_label_default, &span_sourcedestination_group_match_label_updated),
				),
			},
			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelConfig(fvTenantName, spanSrcGrpName, rName),
			},
			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists(resourceName, &span_sourcedestination_group_match_label_updated),
					resource.TestCheckResourceAttr(resourceName, "span_source_group_dn", fmt.Sprintf("uni/tn-%s/srcgrp-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciSPANSourcedestinationGroupMatchLabelIdNotEqual(&span_sourcedestination_group_match_label_default, &span_sourcedestination_group_match_label_updated),
				),
			},
		},
	})
}

func TestAccAciSPANSourcedestinationGroupMatchLabel_Update(t *testing.T) {
	var span_sourcedestination_group_match_label_default models.SPANSourcedestinationGroupMatchLabel
	var span_sourcedestination_group_match_label_updated models.SPANSourcedestinationGroupMatchLabel
	resourceName := "aci_span_sourcedestination_group_match_label.test"
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	spanSrcGrpName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSPANSourcedestinationGroupMatchLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelConfig(fvTenantName, spanSrcGrpName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists(resourceName, &span_sourcedestination_group_match_label_default),
				),
			},

			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelUpdatedAttr(fvTenantName, spanSrcGrpName, rName, "tag", "black"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists(resourceName, &span_sourcedestination_group_match_label_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "black"),
					testAccCheckAciSPANSourcedestinationGroupMatchLabelIdEqual(&span_sourcedestination_group_match_label_default, &span_sourcedestination_group_match_label_updated),
				),
			},

			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelUpdatedAttr(fvTenantName, spanSrcGrpName, rName, "tag", "gray"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists(resourceName, &span_sourcedestination_group_match_label_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "gray"),
					testAccCheckAciSPANSourcedestinationGroupMatchLabelIdEqual(&span_sourcedestination_group_match_label_default, &span_sourcedestination_group_match_label_updated),
				),
			},
			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelUpdatedAttr(fvTenantName, spanSrcGrpName, rName, "tag", "green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists(resourceName, &span_sourcedestination_group_match_label_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "green"),
					testAccCheckAciSPANSourcedestinationGroupMatchLabelIdEqual(&span_sourcedestination_group_match_label_default, &span_sourcedestination_group_match_label_updated),
				),
			},

			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelUpdatedAttr(fvTenantName, spanSrcGrpName, rName, "tag", "red"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists(resourceName, &span_sourcedestination_group_match_label_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "red"),
					testAccCheckAciSPANSourcedestinationGroupMatchLabelIdEqual(&span_sourcedestination_group_match_label_default, &span_sourcedestination_group_match_label_updated),
				),
			},
			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelUpdatedAttr(fvTenantName, spanSrcGrpName, rName, "tag", "white"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists(resourceName, &span_sourcedestination_group_match_label_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "white"),
					testAccCheckAciSPANSourcedestinationGroupMatchLabelIdEqual(&span_sourcedestination_group_match_label_default, &span_sourcedestination_group_match_label_updated),
				),
			},

			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelUpdatedAttr(fvTenantName, spanSrcGrpName, rName, "tag", "yellow"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists(resourceName, &span_sourcedestination_group_match_label_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "yellow"),
					testAccCheckAciSPANSourcedestinationGroupMatchLabelIdEqual(&span_sourcedestination_group_match_label_default, &span_sourcedestination_group_match_label_updated),
				),
			},

			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelConfig(fvTenantName, spanSrcGrpName, rName),
			},
		},
	})
}

func TestAccAciSPANSourcedestinationGroupMatchLabel_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	spanSrcGrpName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSPANSourcedestinationGroupMatchLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelConfig(fvTenantName, spanSrcGrpName, rName),
			},
			{
				Config:      CreateAccSPANSourcedestinationGroupMatchLabelWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccSPANSourcedestinationGroupMatchLabelUpdatedAttr(fvTenantName, spanSrcGrpName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSPANSourcedestinationGroupMatchLabelUpdatedAttr(fvTenantName, spanSrcGrpName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSPANSourcedestinationGroupMatchLabelUpdatedAttr(fvTenantName, spanSrcGrpName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccSPANSourcedestinationGroupMatchLabelUpdatedAttr(fvTenantName, spanSrcGrpName, rName, "tag", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccSPANSourcedestinationGroupMatchLabelUpdatedAttr(fvTenantName, spanSrcGrpName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelConfig(fvTenantName, spanSrcGrpName, rName),
			},
		},
	})
}

func testAccCheckAciSPANSourcedestinationGroupMatchLabelExists(name string, span_sourcedestination_group_match_label *models.SPANSourcedestinationGroupMatchLabel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("SPANSourcedestinationGroupMatchLabel %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SPANSourcedestinationGroupMatchLabel dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		span_sourcedestination_group_match_labelFound := models.SPANSourcedestinationGroupMatchLabelFromContainer(cont)
		if span_sourcedestination_group_match_labelFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("SPANSourcedestinationGroupMatchLabel %s not found", rs.Primary.ID)
		}
		*span_sourcedestination_group_match_label = *span_sourcedestination_group_match_labelFound
		return nil
	}
}

func testAccCheckAciSPANSourcedestinationGroupMatchLabelDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing span_sourcedestination_group_match_label destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_span_sourcedestination_group_match_label" {
			cont, err := client.Get(rs.Primary.ID)
			span_sourcedestination_group_match_label := models.SPANSourcedestinationGroupMatchLabelFromContainer(cont)
			if err == nil {
				return fmt.Errorf("SPANSourcedestinationGroupMatchLabel %s Still exists", span_sourcedestination_group_match_label.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSPANSourcedestinationGroupMatchLabelIdEqual(m1, m2 *models.SPANSourcedestinationGroupMatchLabel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("span_sourcedestination_group_match_label DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSPANSourcedestinationGroupMatchLabelIdNotEqual(m1, m2 *models.SPANSourcedestinationGroupMatchLabel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("span_sourcedestination_group_match_label DNs are equal")
		}
		return nil
	}
}

func CreateSPANSourcedestinationGroupMatchLabelWithoutRequired(fvTenantName, spanSrcGrpName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing span_sourcedestination_group_match_label creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_span_source_group" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	`
	switch attrName {
	case "span_source_group_dn":
		rBlock += `
	resource "aci_span_sourcedestination_group_match_label" "test" {
	#	span_source_group_dn  = aci_span_source_group.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, spanSrcGrpName, rName)
}

func CreateAccSPANSourcedestinationGroupMatchLabelConfigWithRequiredParams(prName, rName string) string {
	fmt.Printf("=== STEP  testing span_sourcedestination_group_match_label creation with parent resource name %s and resource name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = "%s"
	}
	`, prName, prName, rName)
	return resource
}
func CreateAccSPANSourcedestinationGroupMatchLabelConfigUpdatedName(fvTenantName, spanSrcGrpName, rName string) string {
	fmt.Println("=== STEP  testing span_sourcedestination_group_match_label creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = "%s"
	}
	`, fvTenantName, spanSrcGrpName, rName)
	return resource
}

func CreateAccSPANSourcedestinationGroupMatchLabelConfig(fvTenantName, spanSrcGrpName, rName string) string {
	fmt.Println("=== STEP  testing span_sourcedestination_group_match_label creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = "%s"
	}
	`, fvTenantName, spanSrcGrpName, rName)
	return resource
}

func CreateAccSPANSourcedestinationGroupMatchLabelWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing span_sourcedestination_group_match_label creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccSPANSourcedestinationGroupMatchLabelConfigWithOptionalValues(fvTenantName, spanSrcGrpName, rName string) string {
	fmt.Println("=== STEP  Basic: testing span_sourcedestination_group_match_label creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = "${aci_span_source_group.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_span_sourcedestination_group_match_label"
		tag = "blue"
		
	}
	`, fvTenantName, spanSrcGrpName, rName)

	return resource
}

func CreateAccSPANSourcedestinationGroupMatchLabelRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing span_sourcedestination_group_match_label updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_span_sourcedestination_group_match_label" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_span_sourcedestination_group_match_label"
		tag = "alice-blue"
		
	}
	`)

	return resource
}

func CreateAccSPANSourcedestinationGroupMatchLabelUpdatedAttr(fvTenantName, spanSrcGrpName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing span_sourcedestination_group_match_label attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, spanSrcGrpName, rName, attribute, value)
	return resource
}

func CreateAccSPANSourcedestinationGroupMatchLabelUpdatedAttrList(fvTenantName, spanSrcGrpName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing span_sourcedestination_group_match_label attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = "%s"
		%s = %s
	}
	`, fvTenantName, spanSrcGrpName, rName, attribute, value)
	return resource
}
