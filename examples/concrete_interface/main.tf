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

resource "aci_l4_l7_device" "virtual_device" {
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

resource "aci_concrete_device" "virtual_concrete" {
  l4_l7_devices_dn                 = aci_l4_l7_devices.virtual_device.id
  name                             = "tenant1-ASA1"
  clone_count                      = "0"
  is_clone_operation               = "no"
  is_template                      = "no"
  vcenter_name                     = "vcenter"
  vm_name                          = "tenant1-ASA1"
  relation_vns_rs_c_dev_to_ctrlr_p = "uni/vmmp-VMware/dom-ACI-vDS/ctrlr-vcenter"
}

# Creating an interface for a Virtual Concrete Device
resource "aci_concrete_interface" "example1" {
  concrete_device_dn            = aci_concrete_device.virtual_concrete.id
  name                          = "g0/4"
  encap                         = "unknown"
  vnic_name                     = "Network adapter 5"
  relation_vns_rs_c_if_path_att = "topology/pod-1/paths-101/pathep-[eth1/1]"
}

resource "aci_l4_l7_device" "physical_device" {
  tenant_dn                            = aci_tenant.terraform_tenant.id
  name                                 = "example2"
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
  l4_l7_device_dn   = aci_l4_l7_device.physical_device.id
  name              = "physical-Device"
}

# Creating an interface for a Physical Concrete Device
resource "aci_concrete_interface" "example2" {
  concrete_device_dn            = aci_concrete_device.physical_concrete.id
  name                          = "g0/3"
  relation_vns_rs_c_if_path_att = "topology/pod-1/paths-101/pathep-[eth1/2]"
}
