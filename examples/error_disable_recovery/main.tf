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

resource "aci_error_disable_recovery" "example" {
    annotation          = "orchestrator:terraform"
    err_dis_recov_intvl = "300"
    name_alias          = "error_disable_recovery_alias"
    description         = "From Terraform"
    edr_event {
    event               = "event-mcp-loop"
    recover             = "yes"
    description         = "From Terraform"
    name_alias          = "event_alias"
    name                = "example"
    annotation          = "orchestrator:terraform"
    }
    edr_event_ids       = []
}