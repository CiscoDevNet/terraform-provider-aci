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


resource "aci_service_redirect_policy" "example" {
  tenant_dn              = aci_tenant.tenentcheck.id
  name                   = "first"
  name_alias             = "name_alias"
  dest_type              = "L3"
  min_threshold_percent  = "30"
  max_threshold_percent  = "50"
  hashing_algorithm      = "sip"
  description            = "hello"
  anycast_enabled        = "no"
  resilient_hash_enabled = "no"
  threshold_enable       = "no"
  program_local_pod_only = "no"
  threshold_down_action  = "permit"
}