package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAciRest_tenant(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAciRestConfigTenant,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aci_rest.infra", "class_name", "fvTenant"),
					resource.TestCheckResourceAttr("data.aci_rest.infra", "id", "uni/tn-infra"),
					resource.TestCheckResourceAttr("data.aci_rest.infra", "content.name", "infra"),
					resource.TestCheckResourceAttr("data.aci_rest.infra", "child.0.class_name", "aaaDomainRef"),
					resource.TestCheckResourceAttr("data.aci_rest.infra", "child.0.content.name", "infra"),
				),
			},
		},
	})
}

const testAccDataSourceAciRestConfigTenant = `
data "aci_rest" "infra" {
  dn = "uni/tn-infra"
}
`
