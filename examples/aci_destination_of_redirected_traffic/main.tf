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

// The mac attribute in "aci_destination_of_redirected_traffic" is required for APIC prior to Version 5.2 release.
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

// The mac attribute in "aci_destination_of_redirected_traffic" can be foregone when relation_vns_rs_ipsla_monitoring_pol is set in the resource "aci_service_redirect_policy" for APIC Version 5.2 release and above.
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

resource "aci_destination_of_redirected_traffic" "example" {
  service_redirect_policy_dn            = aci_service_redirect_policy.service_policy.id
  ip                                    = "1.2.3.4"
  ip2                                   = "10.20.30.40"
  dest_name                             = "last"
  pod_id                                = "5"
  annotation                            = "load_traffic_dest"
  description                           = "From Terraform"
  name_alias                            = "load_traffic_dest"
  relation_vns_rs_redirect_health_group = aci_l4_l7_redirect_health_group.health_group.id
}
