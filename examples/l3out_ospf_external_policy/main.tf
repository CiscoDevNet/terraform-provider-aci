provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_l3out_ospf_external_policy" "example" {

  l3_outside_dn  = "${aci_l3_outside.example.id}"
  annotation  = "example"
  area_cost  = "1"
  area_ctrl = "redistribute,summary"
  area_id  = "0.0.0.1"
  area_type = "nssa"
  multipod_internal = "no"
  name_alias  = "example"

}

