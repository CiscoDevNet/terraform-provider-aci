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

resource "aci_user_security_domain_role" "example" {
    user_domain_dn  = aci_user_security_domain.example.id
    annotation     = "orchestrator:terraform"
    name            = "example"
    priv_type       = "readPriv"
    name_alias      = "user_role_alias"
    description     = "From Terraform"
}