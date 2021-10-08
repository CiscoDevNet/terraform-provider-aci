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

resource "aci_port_tracking" "example" {
    admin_st           = "off"
    annotation         = "orchestrator:terraform"
    delay              = "120"
    include_apic_ports = "no"
    minlinks           = "0"
    name_alias         = "port_tracking_alias"
    description        = "From Terraform"
}