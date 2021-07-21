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

resource "aci_vlan_pool" "VLAN-POOL" {
  name  = "VLAN-POOL"
  alloc_mode  = "static"
}