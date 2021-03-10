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
  name       = "vlan_pool_P1"
  alloc_mode = "dynamic"
}

resource "aci_vmm_domain" "vds" {
  provider_profile_dn = var.provider_profile_dn
  relation_infra_rs_vlan_ns = aci_vlan_pool.vmm_vlan_pool.id
  name                = var.vmm_domain
}

resource "aci_vmm_controller" "vmware_controller" {
  vmm_domain_dn = aci_vmm_domain.vds.id
  name = var.aci_vmm_controller
  host_or_ip = "10.10.10.10"
  root_cont_name = "vmmdc"
}

resource "aci_vmm_credential" "vmware_credential" {
  vmm_domain_dn = aci_vmm_domain.vds.id
  name = var.aci_vmm_credential
  pwd = "mySecretPassword"
  usr = "myUsername"
}

resource "aci_vmm_controller" "vmware_controller_2" {
  vmm_domain_dn = aci_vmm_domain.vds.id
  relation_vmm_rs_acc = aci_vmm_credential.vmware_credential.id
  name = "vmware_controller_2"
  host_or_ip = "10.10.10.1"
  root_cont_name = "vmmdc"
}

resource "aci_lldp_interface_policy" "LLDP_policy" {
  name       = "vmm_lldp"
}
resource "aci_lacp_policy" "port_channel_policy" {
  name       = "vmm_lacp"
}

resource "aci_vswitch_policy" "vmware_switch_policy" {
  vmm_domain_dn = aci_vmm_domain.vds.id
  // relation_vmm_rs_vswitch_override_lacp_pol = aci_lldp_interface_policy.port_channel_policy.id
  // relation_vmm_rs_vswitch_override_lldp_if_pol = aci_lacp_policy.LLDP_policy.id
}
