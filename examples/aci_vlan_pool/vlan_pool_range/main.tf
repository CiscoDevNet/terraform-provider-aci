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

resource "aci_vlan_pool" "vlan_pool" {
  name  = "VLAN-POOL"
  alloc_mode  = "static"
}

resource "aci_ranges" "pool_range_1" {
  vlan_pool_dn  = aci_vlan_pool.vlan_pool.id
  from          = "vlan-1"
  to            = "vlan-10"
  alloc_mode    = "inherit"
  role          = "external"
}

resource "aci_ranges" "pool_range_2" {
  vlan_pool_dn  = aci_vlan_pool.vlan_pool.id
  from          = "vlan-100"
  to            = "vlan-150"
  description   = "pool range description"
}