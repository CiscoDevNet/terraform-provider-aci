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
resource "aci_pod_maintenance_group" "example" {
  name                       = "mgmt"
  fwtype                     = "controller"
  description                = "from terraform"
  name_alias                 = "aliasing"
  pod_maintenance_group_type = "ALL"
}

