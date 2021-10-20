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

resource "aci_ldap_group_map_rule" "example" {
    name        = "example"
    type        = "ldap"
    description = "From Terraform"
    name_alias  = "ldap_group_map_rule_alias"
}