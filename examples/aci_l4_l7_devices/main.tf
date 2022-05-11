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
  url      = "" # <cloud APIC URL>
  insecure = true
}

resource "aci_tenant" "terraform_tenant" {
  name        = "tf_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_l4_l7_devices" "virtual_device" {
  tenant_dn     = aci_tenant.terraform_tenant.id
  name          = "example1"
  active        = "no"
  context_aware = "single-Context"
  devtype       = "VIRTUAL"
  func_type     = "GoTo"
  is_copy       = "no"
  mode          = "legacy-Mode"
  prom_mode     = "no"
  svc_type      = "OTHERS"
  trunking      = "no"
  relation_vns_rs_al_dev_to_dom_p {
    target_dn      = "uni/vmmp-VMware/dom-ESX0-leaf102-vds"
    switching_mode = "AVE"
  }
}

resource "aci_l4_l7_devices" "physical_device" {
  tenant_dn                            = aci_tenant.terraform_tenant.id
  name                                 = "example2"
  active                               = "no"
  context_aware                        = "single-Context"
  devtype                              = "PHYSICAL"
  func_type                            = "GoTo"
  is_copy                              = "no"
  mode                                 = "legacy-Mode"
  prom_mode                            = "no"
  svc_type                             = "OTHERS"
  relation_vns_rs_al_dev_to_phys_dom_p = "uni/phys-test_dom"
}

