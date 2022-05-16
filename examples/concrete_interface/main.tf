terraform {
  required_providers {
    aci = {
      source = "CiscoDevNet/aci"
    }
  }
}

provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

resource "aci_tenant" "terraform_tenant" {
  name        = "tf_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_l4_l7_devices" "device" {
  tenant_dn     = aci_tenant.terraform_tenant.id
  name          = "tenant1-ASAv"
  active        = "no"
  context_aware = "single-Context"
  devtype       = "VIRTUAL"
  func_type     = "GoTo"
  is_copy       = "no"
  mode          = "legacy-Mode"
  prom_mode     = "no"
  svc_type      = "FW"
  trunking      = "no"
  relation_vns_rs_al_dev_to_dom_p {
    domain_dn = "uni/vmmp-VMware/dom-ACI-vDS"
  }
}

resource "aci_concrete_device" "concrete" {
  l4_l7_devices_dn                 = aci_l4_l7_devices.device.id
  name                             = "tenant1-ASA1"
  clone_count                      = "0"
  is_clone_operation               = "no"
  is_template                      = "no"
  vcenter_name                     = "vcenter"
  vm_name                          = "tenant1-ASA1"
  relation_vns_rs_c_dev_to_ctrlr_p = "uni/vmmp-VMware/dom-ACI-vDS/ctrlr-vcenter"
}

resource "aci_concrete_interface" "example" {
  concrete_device_dn            = aci_concrete_device.concrete.id
  name                          = "g0/4"
  encap                         = "unknown"
  vnic_name                     = "Network adapter 5"
  relation_vns_rs_c_if_path_att = "topology/pod-1/paths-101/pathep-[eth1/1]"
}