
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_l3out_bgp_protocol_profile" "example" {

  logical_node_profile_dn = aci_logical_node_profile.example.id
  annotation              = "example"
  name_alias              = "example"
  description             = "from terraform"

}
