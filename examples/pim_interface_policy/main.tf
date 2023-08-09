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

resource "aci_tenant" "example" {
  description = "sample aci_tenant from terraform"
  name        = "tf_tenant"
}

resource "aci_pim_interface_policy" "exampleip" {
  tenant_dn                         = aci_tenant.example.id
  name                              = "example_ip"
  auth_type                         = "ah-md5"
  designated_router_delay           = "3"
  designated_router_priority        = "1"
  hello_interval                    = "30000"
  join_prune_interval               = "60"
  control_state                     = ["border"]
  inbound_join_prune_filter_policy  = "uni/tn-common/rtmap-test-1"
  outbound_join_prune_filter_policy = "uni/tn-common/rtmap-test-2"
  neighbor_filter_policy            = "uni/tn-common/rtmap-test-3"
}
