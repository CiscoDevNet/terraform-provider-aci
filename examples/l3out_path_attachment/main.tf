provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_l3out_path_attachment" "example" {

  logical_interface_profile_dn = aci_logical_interface_profile.example.id
  target_dn                    = "topology/pod-1/paths-101/pathep-[eth1/1]"
  if_inst_t                    = "ext-svi"
  addr                         = "0.0.0.0"
  annotation                   = "example"
  autostate                    = "disabled"
  encap                        = "vlan-1"
  encap_scope                  = "ctx"
  ipv6_dad                     = "disabled"
  ll_addr                      = "::"
  mac                          = "0F:0F:0F:0F:FF:FF"
  mode                         = "native"
  mtu                          = "inherit"
  target_dscp                  = "AF11"

}
