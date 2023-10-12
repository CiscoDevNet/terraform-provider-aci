package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceAciRestManaged_tenant(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAciRestManagedConfigTenant,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aci_rest_managed.infra", "class_name", "fvTenant"),
					resource.TestCheckResourceAttr("data.aci_rest_managed.infra", "id", "uni/tn-infra"),
					resource.TestCheckResourceAttr("data.aci_rest_managed.infra", "content.name", "infra"),
					resource.TestCheckResourceAttr("data.aci_rest_managed.infra", "child.0.class_name", "aaaDomainRef"),
					resource.TestCheckResourceAttr("data.aci_rest_managed.infra", "child.0.content.name", "infra"),
				),
			},
		},
	})
}

const testAccDataSourceAciRestManagedConfigTenant = `
data "aci_rest_managed" "infra" {
  dn = "uni/tn-infra"
  class_name = "fvTenant"
  content = {
	name = "infra"
  }
  child {
	class_name = "aaaDomainRef"
	rn = "domain-infra"
	content = {
		name = "infra"
	}
  }
}
`
