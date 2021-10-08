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

resource "aci_qos_instance_policy" "example" {
    name_alias            = "qos_instance_alias"
    description           = "From Terraform"
    etrap_age_timer       = "0" 
    etrap_bw_thresh       = "0"
    etrap_byte_ct         = "0"
    etrap_st              = "no"
    fabric_flush_interval = "500"
    fabric_flush_st       = "false"
    annotation            = "orchestrator:terraform"
    ctrl                  = "none"
    uburst_spine_queues   = "10"
    uburst_tor_queues     = "10"
}