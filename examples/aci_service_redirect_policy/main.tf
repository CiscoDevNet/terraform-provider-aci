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

resource "aci_ip_sla_monitoring_policy" "sla" {
  tenant_dn             = aci_tenant.terraform_tenant.id
  name                  = "example"
  type_of_service       = "0"
  traffic_class_value   = "0"
  sla_detect_multiplier = "3"
  sla_frequency         = "60"
  sla_port              = "10"
  sla_type              = "tcp"
  threshold             = "900"
  timeout               = "900"
}

resource "aci_service_redirect_policy" "example" {
  tenant_dn                            = aci_tenant.terraform_tenant.id
  name                                 = "first"
  name_alias                           = "name_alias"
  dest_type                            = "L3"
  min_threshold_percent                = "30"
  max_threshold_percent                = "50"
  hashing_algorithm                    = "sip"
  description                          = "hello"
  anycast_enabled                      = "no"
  resilient_hash_enabled               = "no"
  threshold_enable                     = "no"
  program_local_pod_only               = "no"
  threshold_down_action                = "permit"
  relation_vns_rs_ipsla_monitoring_pol = aci_ip_sla_monitoring_policy.sla.id
}
