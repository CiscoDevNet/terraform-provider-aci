
data "aci_epg_useg_block_statement" "example_application_epg" {
  parent_dn = aci_application_epg.example.id
}
