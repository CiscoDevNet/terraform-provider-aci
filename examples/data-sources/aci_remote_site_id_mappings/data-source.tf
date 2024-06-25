
data "aci_remote_site_id_mappings" "example_associated_site" {
  parent_dn = aci_associated_site.example.id
  site_id   = "102"
}
