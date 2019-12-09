package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciTenant_Basic(t *testing.T) {
	var tenant models.Tenant
	description := "tenant created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTenantConfig_basic(description, "tag_tenant"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTenantExists("aci_tenant.footenant", &tenant),
					testAccCheckAciTenantAttributes(description, "tag_tenant", &tenant),
				),
			},
			{
				ResourceName:      "aci_tenant",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciTenant_update(t *testing.T) {
	var tenant models.Tenant
	description := "tenant created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTenantConfig_basic(description, "tag_tenant"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTenantExists("aci_tenant.footenant", &tenant),
					testAccCheckAciTenantAttributes(description, "tag_tenant", &tenant),
				),
			},
			{
				Config: testAccCheckAciTenantConfig_basic(description, "tag_change"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTenantExists("aci_tenant.footenant", &tenant),
					testAccCheckAciTenantAttributes(description, "tag_change", &tenant),
				),
			},
		},
	})
}

func testAccCheckAciTenantConfig_basic(description, annotation string) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "footenant" {
		description = "%s"
		name        = "demo_tenant"
		annotation  = "%s"
		name_alias  = "alias_tenant"
	} 
	`, description, annotation)
}

func testAccCheckAciTenantExists(name string, tenant *models.Tenant) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Tenant %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Tenant dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tenantFound := models.TenantFromContainer(cont)
		if tenantFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Tenant %s not found", rs.Primary.ID)
		}
		*tenant = *tenantFound
		return nil
	}
}

func testAccCheckAciTenantDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_tenant" {
			cont, err := client.Get(rs.Primary.ID)
			tenant := models.TenantFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Tenant %s Still exists", tenant.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciTenantAttributes(description, annotation string, tenant *models.Tenant) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != tenant.Description {
			return fmt.Errorf("Bad tenant Description %s", tenant.Description)
		}

		if "demo_tenant" != tenant.Name {
			return fmt.Errorf("Bad tenant name %s", tenant.Name)
		}

		if annotation != tenant.Annotation {
			return fmt.Errorf("Bad tenant annotation %s", tenant.Annotation)
		}

		if "alias_tenant" != tenant.NameAlias {
			return fmt.Errorf("Bad tenant name_alias %s", tenant.NameAlias)
		}

		return nil
	}
}
