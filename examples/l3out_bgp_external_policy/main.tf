provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_l3out_bgp_external_policy" "example" {

  l3_outside_dn = aci_l3_outside.example.id
  annotation    = "example"
  name_alias    = "example"

}
