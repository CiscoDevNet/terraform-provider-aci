
terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_tenant" "example" {
  name = "tf_tenant_l3out"
}

resource "aci_vrf" "vrf_example_tf" {
  tenant_dn = aci_tenant.example.id
  name = "tf_vrf_l3out"
}

resource "aci_l3_outside" "example" {
  tenant_dn = aci_tenant.example.id
  name      = "demo_l3out"
  relation_l3ext_rs_ectx = aci_vrf.vrf_example_tf.id
}

resource "aci_logical_node_profile" "node_profile" {
  l3_outside_dn = aci_l3_outside.example.id
  name          = "demo_node"
}

resource "aci_logical_interface_profile" "interface_profile" {
  logical_node_profile_dn = aci_logical_node_profile.node_profile.id
  name                    = "demo_int_prof"
}


resource "aci_l3out_path_attachment" "l3out_path" {
  logical_interface_profile_dn = aci_logical_interface_profile.interface_profile.id
  target_dn                    = "topology/pod-1/protpaths-101-102/pathep-[hxpel-pod3-ucs-a]"
  if_inst_t                    = "ext-svi"
  encap                        = "vlan-2"
}

resource "aci_l3out_vpc_member" "side_A" {
  leaf_port_dn = aci_l3out_path_attachment.l3out_path.id
  side         = "A"
  addr         = "10.0.0.2/16"
  ipv6_dad     = "enabled"
}

resource "aci_l3out_path_attachment_secondary_ip" "side_A" {
  l3out_path_attachment_dn = aci_l3out_vpc_member.side_A.id
  addr                     = "10.0.0.1/16"
  ipv6_dad                 = "enabled"
}

resource "aci_l3out_vpc_member" "side_B" {
  leaf_port_dn = aci_l3out_path_attachment.l3out_path.id
  side         = "B"
  addr         = "10.0.0.3/16"
  ipv6_dad     = "enabled"
}

resource "aci_l3out_path_attachment_secondary_ip" "side_B" {
  l3out_path_attachment_dn = aci_l3out_vpc_member.side_B.id
  addr                     = "10.0.0.1/16"
  ipv6_dad                 = "enabled"
}
