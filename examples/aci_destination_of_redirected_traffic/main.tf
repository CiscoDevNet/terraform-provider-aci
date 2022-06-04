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

resource "aci_tenant" "terraform_tenant" {
  name        = "tf_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_l4_l7_redirect_health_group" "health_group" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "example"
}

resource "aci_service_redirect_policy" "service_policy" {
  tenant_dn              = aci_tenant.terraform_tenant.id
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

resource "aci_destination_of_redirected_traffic" "example" {
  service_redirect_policy_dn            = aci_service_redirect_policy.service_policy.id
  ip                                    = "1.2.3.4"
  mac                                   = "12:25:56:98:45:74"
  ip2                                   = "10.20.30.40"
  dest_name                             = "last"
  pod_id                                = "5"
  annotation                            = "load_traffic_dest"
  description                           = "From Terraform"
  name_alias                            = "load_traffic_dest"
  relation_vns_rs_redirect_health_group = aci_l4_l7_redirect_health_group.health_group.id
}