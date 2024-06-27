
data "aci_epg_useg_sub_block_statement" "example_epg_useg_block_statement" {
  parent_dn = aci_epg_useg_block_statement.example.id
  name      = "sub_criterion"
}

data "aci_epg_useg_sub_block_statement" "example_epg_useg_sub_block_statement" {
  parent_dn = aci_epg_useg_sub_block_statement.example.id
  name      = "sub_criterion"
}
