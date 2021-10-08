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

resource "aci_coop_policy" "example" {
    annotation  = "orchestrator:terraform"
    type        = "compatible"
    name_alias  = "alias_coop_policy"
    description = "From Terraform"
}