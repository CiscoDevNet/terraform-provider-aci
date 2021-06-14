provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_tenant" "footenant" {
  description = "sample aci_tenant from terraform"
  name        = "tenant_1"
  annotation  = "tenant_1_tag"
  name_alias  = "alias_tenant"
}
	  
resource "aci_l3_outside" "fool3_outside" {
  tenant_dn      = aci_tenant.footenant.id
  description    = "sample aci_l3_outside"
  name           = "l3_outside_1"
  annotation     = "l3_outside_1_tag"
  enforce_rtctrl = "export"
  name_alias     = "alias_out"
  target_dscp    = "unspecified"
}

resource "aci_logical_node_profile" "foological_node_profile" {
  l3_outside_dn = aci_l3_outside.fool3_outside.id
  description   = "sample logical node profile"
  name          = "logical_node_profile_1"
  annotation    = "logical_node_profile_1_tag"
  config_issues = "none"
  name_alias    = "alias_node"
  tag           = "black"
  target_dscp   = "unspecified"
}

resource "aci_logical_interface_profile" "foological_interface_profile" {
  logical_node_profile_dn = aci_logical_node_profile.foological_node_profile.id
  description             = "aci_logical_interface_profile from terraform"
  name                    = "logical_interface_profile_1"
  annotation              = "logical_interface_profile_1_tag"
  name_alias              = "alias_prof"
  prio                    = "unspecified"
  tag                     = "black"
}	  


resource "aci_l3out_path_attachment" "fool3out_path_attachment" {
  logical_interface_profile_dn = aci_logical_interface_profile.foological_interface_profile.id
  target_dn                    = "topology/pod-1/paths-101/pathep-[eth1/1]"
  if_inst_t                    = "ext-svi"
  description                  = "from terraform"
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
