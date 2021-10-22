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


resource "aci_spine_port_policy_group" "foospine_access_port_policy_group" {
  description = "from terraform"
  name        = "spine_port_policy_group_1"
  annotation  = "spine_port_policy_group_tag"
  name_alias  = "example"
}

data "aci_spine_port_policy_group" "example4" {
  name = aci_spine_port_policy_group.foospine_access_port_policy_group.name
}

output "name4" {
  value = data.aci_spine_port_policy_group.example4
}