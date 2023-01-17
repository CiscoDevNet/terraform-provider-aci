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

# Tenant Setup
resource "aci_tenant" "terraform_tenant" {
  name = "terraform_tenant"
}

# Concrete Interface Setup
resource "aci_l4_l7_device" "l4_l7_device" {
  tenant_dn   = aci_tenant.terraform_tenant.id
  name        = "tf_l4_l7_device"
  device_type = "CLOUD"
}

resource "aci_concrete_device" "concrete_device" {
  l4_l7_device_dn = aci_l4_l7_device.l4_l7_device.id
  name            = "tf_concrete_device"
}

resource "aci_concrete_interface" "concrete_interface" {
  concrete_device_dn = aci_concrete_device.concrete_device.id
  name               = "tf_concrete_interface"
}

# Redirect Health Group Setup
resource "aci_l4_l7_redirect_health_group" "l4_l7_health_group" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "l4_l7_health_group"
}

# IP SLA Monitoring Policy Setup
resource "aci_ip_sla_monitoring_policy" "ip_sla" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "ip_sla"
}

# Service Redirect Policy(L4-L7 PBR) for L3 Destination with IP SLA
resource "aci_service_redirect_policy" "l4_l7_pbr_with_ip_sla_l3_dst" {
  tenant_dn                            = aci_tenant.terraform_tenant.id
  name                                 = "l4_l7_pbr_with_ip_sla_l3_dst"
  relation_vns_rs_ipsla_monitoring_pol = aci_ip_sla_monitoring_policy.ip_sla.id # Required when we associate L3 Destinations without MAC
}

# Service Redirect Policy(L4-L7 PBR) for L3 Destination without MAC address
resource "aci_destination_of_redirected_traffic" "l3_destinations_without_mac" {
  service_redirect_policy_dn            = aci_service_redirect_policy.l4_l7_pbr_with_ip_sla_l3_dst.id
  ip                                    = "1.1.1.1"
  dest_name                             = "l3_destinations_without_mac"
  relation_vns_rs_redirect_health_group = aci_l4_l7_redirect_health_group.l4_l7_health_group.id
}

# Service Redirect Policy(L4-L7 PBR) for L3 Destination without IP SLA
resource "aci_service_redirect_policy" "l4_l7_pbr_without_ip_sla_l3_dst" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "l4_l7_pbr_without_ip_sla_l3_dst"
}

# Service Redirect Policy(L4-L7 PBR) for L3 Destination with MAC Address
resource "aci_destination_of_redirected_traffic" "l3_destinations_with_mac" {
  service_redirect_policy_dn            = aci_service_redirect_policy.l4_l7_pbr_without_ip_sla_l3_dst.id
  ip                                    = "2.2.2.2"
  mac                                   = "12:25:56:98:45:74"
  dest_name                             = "l3_destinations_with_mac"
  relation_vns_rs_redirect_health_group = aci_l4_l7_redirect_health_group.l4_l7_health_group.id
}

# Service Redirect Policy(L4-L7 PBR) for L1/L2 Destination
resource "aci_service_redirect_policy" "l4_l7_pbr_without_ip_sla_l1_l2_dst" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "l4_l7_pbr_without_ip_sla_l1_l2_dst"
  dest_type = "L1"
}

resource "aci_pbr_l1_l2_destination" "l1_l2_destination" {
  policy_based_redirect_dn                    = aci_service_redirect_policy.l4_l7_pbr_without_ip_sla_l1_l2_dst.id
  destination_name                            = "l1_l2_destination"
  relation_vns_rs_to_c_if                     = aci_concrete_interface.concrete_interface.id
  relation_vns_rs_l1_l2_redirect_health_group = aci_l4_l7_redirect_health_group.l4_l7_health_group.id
}

# Service Redirect Backup Policy(L4-L7 PBR-Backup)
resource "aci_service_redirect_backup_policy" "pbr_backup" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "tf_pbr_backup"
}

# Service Redirect Backup Policy(L4-L7 PBR-Backup) - L3 Destination
resource "aci_destination_of_redirected_traffic" "pbr_backup_l3_destinations" {
  service_redirect_policy_dn = aci_service_redirect_backup_policy.pbr_backup.id
  ip                         = "1.1.1.1"
}

# Service Redirect Backup Policy(L4-L7 PBR-Backup) - L1/L2 Destination
resource "aci_pbr_l1_l2_destination" "pbr_backup_l1_l2_destinations" {
  policy_based_redirect_dn                    = aci_service_redirect_backup_policy.pbr_backup.id
  destination_name                            = "tf_l1_l2_dest"
  name                                        = "tf_test"
  relation_vns_rs_to_c_if                     = aci_concrete_interface.concrete_interface.id
  relation_vns_rs_l1_l2_redirect_health_group = aci_l4_l7_redirect_health_group.l4_l7_health_group.id
}
