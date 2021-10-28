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

resource "aci_rsa_provider" "example" {
    name                   = "example"
    name_alias             = "rsa_provider_alias"
    description            = "From Terraform"
    annotation             = "orchestrator:terraform"
    auth_port              = "1812"
    auth_protocol          = "pap"
    key                    = ""
    monitor_server         = "disabled"
    monitoring_password    = ""
    monitoring_user        = "default"
    retries                = "1"
    timeout                = "5"
}