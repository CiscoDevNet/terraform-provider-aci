
resource "aci_epg_useg_dns_attribute" "example_epg_useg_block_statement" {
  parent_dn = aci_epg_useg_block_statement.example.id
  name      = "dns_attribute"
}
