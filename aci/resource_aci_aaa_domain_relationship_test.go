package aci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciDomainRelationship_Basic(t *testing.T) {
	var aaa_domain_relationship models.AaaDomainRef
	fv_tenant_name := acctest.RandString(5)
	aaa_domain_ref_name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDomainRelationshipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDomainRelationshipConfig_basic(fv_tenant_name, aaa_domain_ref_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDomainRelationshipExists("aci_aaa_domain_relationship.fooaaa_domain_relationship", &aaa_domain_relationship),
					testAccCheckAciDomainRelationshipAttributes(fv_tenant_name, aaa_domain_ref_name, &aaa_domain_relationship),
				),
			},
		},
	})
}

func TestAccAciDomainRelationship_Update(t *testing.T) {
	var aaa_domain_relationship models.AaaDomainRef
	fv_tenant_name := acctest.RandString(5)
	aaa_domain_ref_name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDomainRelationshipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDomainRelationshipConfig_basic(fv_tenant_name, aaa_domain_ref_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDomainRelationshipExists("aci_aaa_domain_relationship.fooaaa_domain_relationship", &aaa_domain_relationship),
					testAccCheckAciDomainRelationshipAttributes(fv_tenant_name, aaa_domain_ref_name, &aaa_domain_relationship),
				),
			},
			{
				Config: testAccCheckAciDomainRelationshipConfig_basic(fv_tenant_name, aaa_domain_ref_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDomainRelationshipExists("aci_aaa_domain_relationship.fooaaa_domain_relationship", &aaa_domain_relationship),
					testAccCheckAciDomainRelationshipAttributes(fv_tenant_name, aaa_domain_ref_name, &aaa_domain_relationship),
				),
			},
		},
	})
}

func testAccCheckAciDomainRelationshipConfig_basic(fv_tenant_name, aaa_domain_ref_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
	}

	resource "aci_aaa_domain" "foosecurity_domain" {
		name        = "%s"
		description = "from terraform"
		annotation  = "aaa_domain_tag"
		name_alias  = "example"
	}

	resource "aci_aaa_domain_relationship" "fooaaa_domain_relationship" {
		aaa_domain_dn 		= aci_aaa_domain.foosecurity_domain.id
		parent_dn = aci_tenant.footenant.id
	}

	`, fv_tenant_name, aaa_domain_ref_name)
}

func testAccCheckAciDomainRelationshipExists(name string, aaa_domain_relationship *models.AaaDomainRef) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("AAA Domain Relationship Object %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AAA Domain Relationship Object dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		aaa_domain_relationshipFound := models.AaaDomainRefFromContainer(cont)
		if aaa_domain_relationshipFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("AAA Domain Relationship Object %s not found", rs.Primary.ID)
		}
		*aaa_domain_relationship = *aaa_domain_relationshipFound
		return nil
	}
}

func testAccCheckAciDomainRelationshipDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_aaa_domain_relationship" {
			cont, err := client.Get(rs.Primary.ID)
			aaa_domain_relationship := models.AaaDomainRefFromContainer(cont)
			if err == nil {
				return fmt.Errorf("AAA Domain Relationship Object %s Still exists", aaa_domain_relationship.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciDomainRelationshipAttributes(fv_tenant_name, aaa_domain_ref_name string, aaa_domain_relationship *models.AaaDomainRef) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aaa_domain_ref_name != GetMOName(aaa_domain_relationship.DistinguishedName) {
			return fmt.Errorf("Bad aaa_domain_ref %s", GetMOName(aaa_domain_relationship.DistinguishedName))
		}
		return nil
	}
}
