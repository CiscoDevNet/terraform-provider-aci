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

resource "aci_radius_provider_group" "example" {
    name        = "example"
    annotation  = "orchestrator:terraform"
    name_alias  = "radius_provider_group_alias"
    description = "From Terraform"
}