
data "aci_epg_useg_mac_attribute" "example_epg_useg_block_statement" {
  parent_dn = aci_epg_useg_block_statement.example.id
  name      = "mac_attr"
}
