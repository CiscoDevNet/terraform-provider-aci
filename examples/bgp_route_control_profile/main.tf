provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}
resource "aci_bgp_route_control_profile" "example" {
  parent_dn                  = aci_tenant.tenentcheck.id
  name                       = "one"
  annotation                 = "example"
  description                = "from terraform"
  name_alias                 = "example"
  route_control_profile_type = "global"
}

resource "aci_bgp_route_control_profile" "example" {
  parent_dn                  = aci_l3_outside.example.id
  name                       = "one"
  annotation                 = "example"
  description                = "from terraform"
  name_alias                 = "example"
  route_control_profile_type = "global"
}
