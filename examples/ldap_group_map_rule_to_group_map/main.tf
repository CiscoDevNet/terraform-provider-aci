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

resource "aci_ldap_group_map_rule_to_group_map" "example" {
  ldap_group_map_dn  = aci_ldap_group_map.example.id
  name  = "example"
  annotation = "orchestrator:terraform"
  name_alias = "example_name_alias_value"
  description = "from terraform"
}