
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_l3out_static_route_next_hop" "example" {

  static_route_dn      = aci_l3out_static_route.example.id
  nh_addr              = "10.0.0.1"
  annotation           = "example"
  name_alias           = "example"
  pref                 = "unspecified"
  nexthop_profile_type = "prefix"
  description          = "from terraform"

}
