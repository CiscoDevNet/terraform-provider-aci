resource "aci_attachable_access_entity_profile" "test_ep" {
    name = "tf_ep"
  
}

resource "aci_vlan_encapsulationfor_vxlan_traffic" "test_vxlan_traffic" {
    attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test_ep.id
    name = "test_vxlan"
  
}
