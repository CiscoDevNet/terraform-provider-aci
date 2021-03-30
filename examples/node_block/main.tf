resource "aci_leaf_profile" "checkBLK" {
  name        = "example"
}

resource "aci_switch_association" "checkBLK" {
  leaf_profile_dn  = aci_leaf_profile.checkBLK.id
  name  = "example"
  switch_association_type  = "range"
}

resource "aci_node_block" "check" {
  switch_association_dn   = aci_switch_association.checkBLK.id
  name                    = "block"
  annotation              = "aci_node_block"
  from_                   = "101"
  name_alias              = "node_block"
  to_                     = "106"
}