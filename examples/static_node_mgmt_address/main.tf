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

resource "aci_static_node_mgmt_address" "example" {
  management_epg_dn = aci_node_mgmt_epg.example.id
  t_dn              = "topology/pod-1/node-1"
  type              = "out_of_band"
  addr              = "10.20.30.40/20"
  annotation        = "example"
  description       = "from terraform"
  gw                = "10.20.30.41"
  v6_addr           = "1::40/64"
  v6_gw             = "1::21"
}