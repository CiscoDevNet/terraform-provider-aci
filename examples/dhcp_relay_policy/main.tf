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

resource "aci_dhcp_relay_policy" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "name_example"
  annotation  = "annotation_example"
  mode  = "visible"
  name_alias  = "alias_example"
  owner  = "infra"
}