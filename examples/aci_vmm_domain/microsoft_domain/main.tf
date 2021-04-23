terraform {
  required_providers {
    aci = {
      source = "CiscoDevNet/aci"
    }
  }
}

provider "aci" {
  username = "admin"
  password = "ins3965!"
  url      = "https://10.23.248.103"
  insecure = true
}

resource "aci_vlan_pool" "vmm_vlan_pool" {
  name       = "vlan_pool_P12"
  alloc_mode = "dynamic"
}

// Microsoft vmm domain resources
resource "aci_vmm_domain" "microsoft_domain" {
  provider_profile_dn = var.microsoft_domain
  relation_infra_rs_vlan_ns = aci_vlan_pool.vmm_vlan_pool.id
  name                = var.vmm_domain
}

resource "aci_vmm_controller" "microsoft_controller" {
  vmm_domain_dn = aci_vmm_domain.microsoft_domain.id
  name = var.aci_vmm_controller
  host_or_ip = "10.10.10.10"
  root_cont_name = "vmmdc"
  scope = "MicrosoftSCVMM"
}

resource "aci_vmm_controller" "microsoft_controller_2" {
  vmm_domain_dn = aci_vmm_domain.microsoft_domain.id
  name = "microsoft_vmm_controller_2"
  host_or_ip = "10.10.10.1"
  root_cont_name = "vmmdc"
  scope = "MicrosoftSCVMM"
}

resource "aci_lldp_interface_policy" "LLDP_policy" {
  name       = "vmm_lldp"
}
resource "aci_lacp_policy" "port_channel_policy" {
  name       = "vmm_lacp"
}
resource "aci_cdp_interface_policy" "foocdp_interface_policy" {
  name = "cdpIfPol1"
}
resource "aci_vswitch_policy" "microsoft_switch_policy" {
  vmm_domain_dn = aci_vmm_domain.microsoft_domain.id
  relation_vmm_rs_vswitch_override_cdp_if_pol = aci_cdp_interface_policy.foocdp_interface_policy.id
  relation_vmm_rs_vswitch_override_lacp_pol = aci_lacp_policy.port_channel_policy.id
  relation_vmm_rs_vswitch_override_lldp_if_pol = aci_lldp_interface_policy.LLDP_policy.id
  // relation_vmm_rs_vswitch_override_stp_pol = "uni/infra/ifPol-stpPolicy"
}
