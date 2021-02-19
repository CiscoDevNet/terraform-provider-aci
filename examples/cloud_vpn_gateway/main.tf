
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_cloud_vpn_gateway" "example" {
  cloud_context_profile_dn  = aci_cloud_context_profile.cloudCtxProf.id
  name                      = "VPN_GW"
  name_alias                = "primaryVPNgw"
  num_instances             = "1"
  cloud_router_profile_type = "vpn-gw"
}
