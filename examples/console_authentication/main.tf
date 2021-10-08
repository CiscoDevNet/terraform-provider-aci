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

resource "aci_console_authentication" "example" {
    annotation     = "orchestrator:terraform"
    provider_group = "60"
    realm          = "ldap"
    realm_sub_type = "default"
    name_alias     = "console_alias"
    description    = "From Terraform"
}