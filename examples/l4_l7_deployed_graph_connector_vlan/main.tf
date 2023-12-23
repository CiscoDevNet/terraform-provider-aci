terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <APIC URL>
  insecure = true
}

resource "aci_tenant" "terraform_tenant" {
  name        = "tf_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_application_profile" "terraform_ap" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "tf_ap"
}

resource "aci_vrf" "vrf" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "vrf1"
}

resource "aci_bridge_domain" "provider_bd" {
  tenant_dn          = aci_tenant.terraform_tenant.id
  relation_fv_rs_ctx = aci_vrf.vrf.id
  name               = "bd_prov"
}

resource "aci_bridge_domain" "consumer_bd" {
  tenant_dn          = aci_tenant.terraform_tenant.id
  relation_fv_rs_ctx = aci_vrf.vrf.id
  name               = "bd_cons"
}

resource "aci_application_epg" "consumer_epg" {
  application_profile_dn = aci_application_profile.terraform_ap.id
  name                   = "consumer_epg"
  relation_fv_rs_bd      = aci_bridge_domain.consumer_bd.id
}

resource "aci_application_epg" "provider_epg" {
  application_profile_dn = aci_application_profile.terraform_ap.id
  name                   = "provider_epg"
  relation_fv_rs_bd      = aci_bridge_domain.provider_bd.id
}

data "aci_vmm_domain" "vmm_dom" {
  provider_profile_dn = "uni/vmmp-VMware"
  name                = "VMware-VMM"
}

resource "aci_epg_to_domain" "consumer_epg" {
  application_epg_dn = aci_application_epg.consumer_epg.id
  tdn                = data.aci_vmm_domain.vmm_dom.id
  res_imedcy         = "immediate"
}

resource "aci_epg_to_domain" "provider_epg" {
  application_epg_dn = aci_application_epg.provider_epg.id
  tdn                = data.aci_vmm_domain.vmm_dom.id
  res_imedcy         = "immediate"
}

resource "aci_filter" "filter" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "any"
}

resource "aci_filter_entry" "entry" {
  filter_dn = aci_filter.filter.id
  name      = "any"
}

resource "aci_contract" "cntrct" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "contract1"
}

resource "aci_contract_subject" "subj" {
  contract_dn                   = aci_contract.cntrct.id
  name                          = "subject1"
  relation_vz_rs_subj_filt_att  = [aci_filter.filter.id]
  relation_vz_rs_subj_graph_att = aci_l4_l7_service_graph_template.template.id
}

resource "aci_epg_to_contract" "consumer" {
  application_epg_dn = aci_application_epg.consumer_epg.id
  contract_dn        = aci_contract.cntrct.id
  contract_type      = "consumer"
}

resource "aci_epg_to_contract" "provider" {
  application_epg_dn = aci_application_epg.provider_epg.id
  contract_dn        = aci_contract.cntrct.id
  contract_type      = "provider"
}

resource "aci_bridge_domain" "provider_service_bd" {
  tenant_dn          = aci_tenant.terraform_tenant.id
  relation_fv_rs_ctx = aci_vrf.vrf.id
  name               = "terraform-lb-provider-service-bd"
}

resource "aci_bridge_domain" "consumer_service_bd" {
  tenant_dn          = aci_tenant.terraform_tenant.id
  relation_fv_rs_ctx = aci_vrf.vrf.id
  name               = "terraform-lb-consumer-service-bd"
}

resource "aci_subnet" "provider_bd_subnet" {
  parent_dn = aci_bridge_domain.provider_bd.id
  ip        = "10.10.101.1/24"
}

resource "aci_subnet" "consumer_bd_subnet" {
  parent_dn = aci_bridge_domain.consumer_bd.id
  ip        = "10.10.103.1/24"
}

resource "aci_subnet" "provider_service_bd_subnet" {
  parent_dn = aci_bridge_domain.provider_service_bd.id
  ip        = "10.10.105.1/24"
}

resource "aci_subnet" "consumer_service_bd_subnet" {
  parent_dn = aci_bridge_domain.consumer_service_bd.id
  ip        = "10.10.106.1/24"
}

resource "aci_l4_l7_device" "virtual_device" {
  tenant_dn        = aci_tenant.terraform_tenant.id
  name             = "tenant1-ASAv"
  context_aware    = "single-Context"
  device_type      = "VIRTUAL"
  function_type    = "GoTo"
  is_copy          = "no"
  mode             = "legacy-Mode"
  promiscuous_mode = "no"
  service_type     = "ADC"
  trunking         = "no"
  relation_vns_rs_al_dev_to_dom_p {
    domain_dn = data.aci_vmm_domain.vmm_dom.id
  }
}

resource "aci_l4_l7_service_graph_template" "template" {
  tenant_dn                         = aci_tenant.terraform_tenant.id
  name                              = "SG1"
  description                       = "SG TF"
  l4_l7_service_graph_template_type = "legacy"
  ui_template_type                  = "UNSPECIFIED"
}

resource "aci_function_node" "function_node" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.template.id
  name                            = "N1"
  func_template_type              = "ADC_TWO_ARM"
  func_type                       = "GoTo"
  is_copy                         = "no"
  managed                         = "no"
  routing_mode                    = "Redirect"
  sequence_number                 = "0"
  share_encap                     = "no"
  relation_vns_rs_node_to_l_dev   = aci_l4_l7_device.virtual_device.id
}

resource "aci_connection" "t1-n1" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.template.id
  name                            = "C2"
  adj_type                        = "L3"
  conn_dir                        = "provider"
  conn_type                       = "external"
  direct_connect                  = "no"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_l4_l7_service_graph_template.template.term_prov_dn,
    aci_function_node.function_node.conn_provider_dn
  ]
}

resource "aci_connection" "n1-t2" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.template.id
  name                            = "C1"
  adj_type                        = "L3"
  conn_dir                        = "consumer"
  conn_type                       = "external"
  direct_connect                  = "no"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_l4_l7_service_graph_template.template.term_cons_dn,
    aci_function_node.function_node.conn_consumer_dn
  ]
}

resource "aci_logical_device_context" "device_context" {
  tenant_dn                          = aci_tenant.terraform_tenant.id
  ctrct_name_or_lbl                  = aci_contract.cntrct.name
  graph_name_or_lbl                  = "SG1"
  node_name_or_lbl                   = "N1"
  relation_vns_rs_l_dev_ctx_to_l_dev = aci_l4_l7_device.virtual_device.id
}

resource "aci_logical_interface_context" "consumer" {
  logical_device_context_dn        = aci_logical_device_context.device_context.id
  conn_name_or_lbl                 = "consumer"
  l3_dest                          = "yes"
  permit_log                       = "no"
  relation_vns_rs_l_if_ctx_to_l_if = aci_l4_l7_logical_interface.External.id
  relation_vns_rs_l_if_ctx_to_bd   = aci_bridge_domain.provider_bd.id
}

resource "aci_logical_interface_context" "provider" {
  logical_device_context_dn                    = aci_logical_device_context.device_context.id
  conn_name_or_lbl                             = "provider"
  l3_dest                                      = "no"
  permit_log                                   = "no"
  relation_vns_rs_l_if_ctx_to_l_if             = aci_l4_l7_logical_interface.Internal.id
  relation_vns_rs_l_if_ctx_to_bd               = aci_bridge_domain.provider_bd.id
  relation_vns_rs_l_if_ctx_to_svc_redirect_pol = aci_service_redirect_policy.service_policy.id
}

data "aci_vmm_controller" "controller" {
  vmm_domain_dn = data.aci_vmm_domain.vmm_dom.id
  name          = "DMZ-vcenter"
}

resource "aci_concrete_device" "virtual_concrete" {
  l4_l7_device_dn   = aci_l4_l7_device.virtual_device.id
  name              = "virtual-Device"
  vmm_controller_dn = data.aci_vmm_controller.controller.id
  vm_name           = "VMware-VMM"
}

resource "aci_concrete_interface" "external_interface" {
  concrete_device_dn = aci_concrete_device.virtual_concrete.id
  name               = "External"
  vnic_name          = "Network adapter 5"
}

resource "aci_concrete_interface" "internal_interface" {
  concrete_device_dn = aci_concrete_device.virtual_concrete.id
  name               = "Internal"
  vnic_name          = "Network adapter 4"
}

resource "aci_l4_l7_logical_interface" "Internal" {
  l4_l7_device_dn            = aci_l4_l7_device.virtual_device.id
  name                       = "Internal"
  relation_vns_rs_c_if_att_n = [aci_concrete_interface.internal_interface.id]
}

resource "aci_l4_l7_logical_interface" "External" {
  l4_l7_device_dn            = aci_l4_l7_device.virtual_device.id
  name                       = "External"
  relation_vns_rs_c_if_att_n = [aci_concrete_interface.external_interface.id]
}

resource "aci_l4_l7_redirect_health_group" "health_group" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "example"
}

resource "aci_ip_sla_monitoring_policy" "sla" {
  tenant_dn             = aci_tenant.terraform_tenant.id
  name                  = "example"
  sla_detect_multiplier = "3"
  sla_frequency         = "60"
  sla_port              = "10"
  sla_type              = "tcp"
}

resource "aci_service_redirect_policy" "service_policy" {
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

resource "aci_destination_of_redirected_traffic" "traffic" {
  service_redirect_policy_dn            = aci_service_redirect_policy.service_policy.id
  ip                                    = "10.10.105.100"
  mac                                   = "01:02:03:04:05:06"
  dest_name                             = "tenant1-ASAv"
  pod_id                                = "5"
  annotation                            = "load_traffic_dest"
  description                           = "From Terraform"
  name_alias                            = "load_traffic_dest"
  relation_vns_rs_redirect_health_group = aci_l4_l7_redirect_health_group.health_group.id
}

data "aci_l4_l7_deployed_graph_connector_vlan" "example1" {
  logical_context_dn = aci_logical_interface_context.consumer.id
}

data "aci_l4_l7_deployed_graph_connector_vlan" "example2" {
  logical_context_dn = aci_logical_interface_context.provider.id
}
