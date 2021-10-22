terraform {
  required_providers {
    aci = {
      source = "CiscoDevNet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_vlan_pool" "vmm_vlan_pool" {
  name       = "vlan_pool_P12"
  alloc_mode = "dynamic"
}


// VMWare ave domain resources
resource "aci_vmm_domain" "ave" {
  provider_profile_dn       = var.ave
  relation_infra_rs_vlan_ns = aci_vlan_pool.vmm_vlan_pool.id
  name                      = var.vmm_domain
  enable_ave                = "yes"
  mcast_addr                = "239.10.10.10"
  // create a multicast address pool and add that to below relationship.
  relation_vmm_rs_dom_mcast_addr_ns = "uni/infra/maddrns-testo"
}

resource "aci_vmm_controller" "vmware_controller" {
  vmm_domain_dn  = aci_vmm_domain.ave.id
  name           = var.aci_vmm_controller
  host_or_ip     = "10.10.10.10"
  root_cont_name = "vmmdc"
}

resource "aci_vmm_credential" "vmware_credential" {
  vmm_domain_dn = aci_vmm_domain.ave.id
  name          = var.aci_vmm_credential
  pwd           = "mySecretPassword"
  usr           = "myUsername"
}

resource "aci_vmm_controller" "vmware_controller_2" {
  vmm_domain_dn       = aci_vmm_domain.ave.id
  relation_vmm_rs_acc = aci_vmm_credential.vmware_credential.id
  name                = "vmware_ave_controlller_2"
  host_or_ip          = "10.10.10.1"
  root_cont_name      = "vmmdc"
}

resource "aci_lldp_interface_policy" "LLDP_policy" {
  name = "vmm_lldp"
}
resource "aci_lacp_policy" "port_channel_policy" {
  name = "vmm_lacp"
}
resource "aci_cdp_interface_policy" "foocdp_interface_policy" {
  name = "cdpIfPol1"
}

resource "aci_vswitch_policy" "vmware_switch_policy" {
  vmm_domain_dn                           = aci_vmm_domain.ave.id
  relation_vmm_rs_vswitch_override_fw_pol = "uni/infra/fwP-firewallPolicy"
  // STP_Policy is available for higher versions of ACI only (>5.1(2e))
  relation_vmm_rs_vswitch_override_stp_pol = "uni/infra/ifPol-stpPolicy"
  relation_vmm_rs_vswitch_override_mtu_pol = "uni/fabric/l2pol-l2InstPolicy"
  relation_vmm_rs_vswitch_exporter_pol {
    target_dn            = "uni/infra/vmmexporterpol-exporter_policy"
    active_flow_time_out = 60
    idle_flow_time_out   = 10
    sampling_rate        = 0
  }
  relation_vmm_rs_vswitch_override_cdp_if_pol  = aci_cdp_interface_policy.foocdp_interface_policy.id
  relation_vmm_rs_vswitch_override_lacp_pol    = aci_lacp_policy.port_channel_policy.id
  relation_vmm_rs_vswitch_override_lldp_if_pol = aci_lldp_interface_policy.LLDP_policy.id
}