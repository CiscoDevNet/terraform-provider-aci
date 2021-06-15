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

resource "aci_maintenance_group_node" "example" {
  pod_maintenance_group_dn = aci_pod_maintenance_group.example.id
  description              = "from terraform"
  name                     = "First"
  annotation               = "example"
  from_                    = "1"
  name_alias               = "aliasing"
  to_                      = "5"
}