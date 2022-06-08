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

resource "aci_tenant" "terraform_tenant" {
  name        = "tf_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_l4_l7_device" "virtual_device" {
  tenant_dn        = aci_tenant.terraform_tenant.id
  name             = "tenant1-ASAv"
  active           = "no"
  context_aware    = "single-Context"
  device_type      = "VIRTUAL"
  function_type    = "GoTo"
  is_copy          = "no"
  mode             = "legacy-Mode"
  promiscuous_mode = "no"
  service_type     = "FW"
  trunking         = "no"
  relation_vns_rs_al_dev_to_dom_p {
    domain_dn = "uni/vmmp-VMware/dom-ACI-vDS"
  }
}

resource "aci_concrete_device" "virtual_concrete" {
  l4_l7_device_dn   = aci_l4_l7_device.virtual_device.id
  name              = "tenant1-ASA1"
  vmm_controller_dn = "uni/vmmp-VMware/dom-ACI-vDS/ctrlr-vcenter"
  vm_name           = "tenant1-ASA1"
}

resource "aci_concrete_interface" "virtual_interface" {
  concrete_device_dn            = aci_concrete_device.virtual_concrete.id
  name                          = "g0/4"
  encap                         = "unknown"
  vnic_name                     = "Network adapter 5"
  relation_vns_rs_c_if_path_att = "topology/pod-1/paths-101/pathep-[eth1/1]"
}

resource "aci_l4_l7_logical_interface" "example1" {
  l4_l7_device_dn            = aci_l4_l7_device.virtual_device.id
  name                       = "example1"
  enhanced_lag_policy_name   = "Lacp"
  relation_vns_rs_c_if_att_n = [aci_concrete_interface.virtual_interface.id]
}

resource "aci_l4_l7_device" "physical_device" {
  tenant_dn                            = aci_tenant.terraform_tenant.id
  name                                 = "tenant1-ASAv2"
  active                               = "no"
  context_aware                        = "single-Context"
  device_type                          = "PHYSICAL"
  function_type                        = "GoTo"
  is_copy                              = "no"
  mode                                 = "legacy-Mode"
  promiscuous_mode                     = "no"
  service_type                         = "OTHERS"
  relation_vns_rs_al_dev_to_phys_dom_p = "uni/phys-test_dom"
}

resource "aci_concrete_device" "physical_concrete" {
  l4_l7_device_dn = aci_l4_l7_device.physical_device.id
  name            = "physical-Device"
}

resource "aci_concrete_interface" "physical_interface" {
  concrete_device_dn            = aci_concrete_device.physical_concrete.id
  name                          = "g0/3"
  relation_vns_rs_c_if_path_att = "topology/pod-1/paths-101/pathep-[eth1/2]"
}

resource "aci_l4_l7_logical_interface" "example2" {
  l4_l7_device_dn            = aci_l4_l7_device.physical_device.id
  name                       = "example2"
  encap                      = "vlan-1"
  relation_vns_rs_c_if_att_n = [aci_concrete_interface.physical_interface.id]
}
