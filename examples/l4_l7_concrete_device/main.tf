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

resource "aci_l4_l7_device" "example" {
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
  l4_l7_device_dn   = aci_l4_l7_device.example.id
  name              = "virtual-Device"
  vmm_controller_dn = "uni/vmmp-VMware/dom-ACI-vDS/ctrlr-vcenter"
  vm_name           = "tenant1-ASA1"
}

resource "aci_concrete_device" "example2" {
  l4_l7_device_dn   = aci_l4_l7_device.example.id
  name              = "physical-Device"
}
