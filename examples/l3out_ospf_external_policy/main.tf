terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

#configure provider with your cisco aci credentials.
provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

resource "aci_tenant" "foo_tenant" {
  name = "foo_tenant"
}

resource "aci_vrf" "vrf1" {
  tenant_dn          = aci_tenant.foo_tenant.id
  bd_enforced_enable = "no"
  knw_mcast_act      = "permit"
  name               = "vrf1"
  pc_enf_dir         = "ingress"
  pc_enf_pref        = "enforced"
}

resource "aci_l3_outside" "foo_l3_outside" {
  tenant_dn              = aci_tenant.foo_tenant.id
  description            = "aci_l3_outside"
  name                   = "demo_l3out"
  annotation             = "tag_l3out"
  enforce_rtctrl         = ["export", "import"]
  name_alias             = "alias_out"
  target_dscp            = "unspecified"
  relation_l3ext_rs_ectx = aci_vrf.vrf1.id

}

resource "aci_l3out_ospf_external_policy" "example" {
  l3_outside_dn     = aci_l3_outside.example.id
  annotation        = "example"
  area_cost         = "1"
  area_ctrl         = ["redistribute", "summary"]
  area_id           = "0.0.0.1"
  area_type         = "nssa"
  multipod_internal = "no"
  name_alias        = "example"
}
