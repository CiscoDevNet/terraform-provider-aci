
resource "aci_remote_site_id_mappings" "example_associated_site" {
  parent_dn         = aci_associated_site.example.id
  remote_vrf_pc_tag = "2818057"
  remote_pc_tag     = "16386"
  site_id           = "100"
}
