terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_l4_l7_service_graph_template" "example" {
  tenant_dn                         = "${aci_tenant.tenentcheck.id}"
  name                              = "second"
  name_alias                        = "alias"
  description                       = "from terraform"
  l4_l7_service_graph_template_type = "cloud"
  ui_template_type                  = "ONE_NODE_ADC_ONE_ARM"
}