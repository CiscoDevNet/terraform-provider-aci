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

resource "aci_saml_provider_group" "example" {
    name        = "example"
    annotation  = "orchestrator:terraform"
    name_alias  = "saml_provider_group_alias"
    description = "From Terraform"
}