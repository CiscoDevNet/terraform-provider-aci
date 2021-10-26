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

resource "aci_l3out_path_attachment_secondary_ip" "example" {
  l3out_path_attachment_dn = aci_l3out_path_attachment.example.id
  addr                     = "10.0.0.1/24"
  annotation               = "example"
  ipv6_dad                 = "disabled"
  name_alias               = "example"
}