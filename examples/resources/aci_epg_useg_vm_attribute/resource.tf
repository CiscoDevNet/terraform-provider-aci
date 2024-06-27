
resource "aci_epg_useg_vm_attribute" "example_epg_useg_block_statement" {
  parent_dn = aci_epg_useg_block_statement.example.id
  name      = "vm_attribute"
  value     = "default_value"
}

resource "aci_epg_useg_vm_attribute" "example_epg_useg_sub_block_statement" {
  parent_dn = aci_epg_useg_sub_block_statement.example.id
  name      = "vm_attribute"
  value     = "default_value"
}
