
resource "aci_associated_site" "example_application_epg" {
  parent_dn = aci_application_epg.example.id
  site_id   = "0"
}

resource "aci_associated_site" "example_vrf" {
  parent_dn = aci_vrf.example.id
  site_id   = "0"
}
