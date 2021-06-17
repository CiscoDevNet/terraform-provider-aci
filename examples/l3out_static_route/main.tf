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

resource "aci_l3out_static_route" "example" {

  fabric_node_dn = aci_logical_node_to_fabric_node.example.id
  ip             = "10.0.0.1"
  aggregate      = "no"
  annotation     = "example"
  name_alias     = "example"
  pref           = "example"
  rt_ctrl        = "bfd"
  description    = "from terraform"

}
