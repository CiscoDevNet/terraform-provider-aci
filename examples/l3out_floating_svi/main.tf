provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}
resource "aci_l3out_floating_svi" "example" {
  logical_interface_profile_dn = aci_logical_interface_profile.example.id
  node_dn                      = "topology/pod-1/node-201"
  encap                        = "vlan-20"
  addr                         = "10.20.30.40/16"
  annotation                   = "example"
  description                  = "from terraform"
  autostate                    = "enabled"
  encap_scope                  = "ctx"
  if_inst_t                    = "ext-svi"
  ipv6_dad                     = "disabled"
  ll_addr                      = "::"
  mac                          = "12:23:34:45:56:67"
  mode                         = "untagged"
  mtu                          = "580"
  target_dscp                  = "CS1"
  userdom                      = ":all:"
}
