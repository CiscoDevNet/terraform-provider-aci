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

resource "aci_endpoint_loop_protection" "example" {
    action            = ["port-disable"]
    admin_st          = "disabled"
    annotation        = "orchestrator:terraform"
    loop_detect_intvl = "60"
    loop_detect_mult  = "4"
    name_alias        = "endpoint_loop_protection_alias"
    description       = "From Terraform"
}