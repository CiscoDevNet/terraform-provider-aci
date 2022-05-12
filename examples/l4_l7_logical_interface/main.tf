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

resource "aci_l4_l7_devices" "devices" {
  tenant_dn                            = aci_tenant.terraform_tenant.id
  name                                 = "example"
  active                               = "no"
  context_aware                        = "single-Context"
  devtype                              = "VIRTUAL"
  func_type                            = "GoTo"
  is_copy                              = "no"
  mode                                 = "legacy-Mode"
  prom_mode                            = "no"
  svc_type                             = "OTHERS"
  trunking                             = "no"
  relation_vns_rs_al_dev_to_phys_dom_p = "uni/phys-test_dom"
}

resource "aci_logical_interface" "example" {
  l4_l7_devices_dn           = aci_l4_l7_devices.devices.id
  name                       = "example"
  relation_vns_rs_c_if_att_n = ["uni/tn-tenant1/lDevVip-example/cDev-test/cIf-[g0/0]", "uni/tn-tenant1/lDevVip-example/cDev-test/cIf-[g0/1]"]
}

