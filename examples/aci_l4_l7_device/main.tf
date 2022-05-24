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

resource "aci_l4_l7_devices" "example1" {
  tenant_dn        = aci_tenant.terraform_tenant.id
  name             = "example1"
  active           = "no"
  context_aware    = "single-Context"
  device_type      = "VIRTUAL"
  function_type    = "GoTo"
  is_copy          = "no"
  mode             = "legacy-Mode"
  promiscuous_mode = "no"
  service_type     = "OTHERS"
  trunking         = "no"
  relation_vns_rs_al_dev_to_dom_p {
    target_dn      = "uni/vmmp-VMware/dom-ESX0-leaf102-vds"
    switching_mode = "AVE"
  }
}

resource "aci_l4_l7_devices" "example2" {
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

