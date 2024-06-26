
data "aci_associated_site" "example_application_epg" {
  parent_dn = aci_application_epg.example.id
}

data "aci_associated_site" "example_bridge_domain" {
  parent_dn = aci_bridge_domain.example.id
}
