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

resource "aci_l4_l7_device" "virtual" {
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

resource "aci_concrete_device" "example1" {
  l4_l7_device_dn   = aci_l4_l7_device.virtual.id
  name              = "virtual-Device"
  vmm_controller_dn = "uni/vmmp-VMware/dom-ACI-vDS/ctrlr-vcenter"
  vm_name           = "tenant1-ASA1"
}

resource "aci_l4_l7_device" "physical" {
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

resource "aci_concrete_device" "example2" {
  l4_l7_device_dn   = aci_l4_l7_device.physical.id
  name              = "physical-Device"
}
