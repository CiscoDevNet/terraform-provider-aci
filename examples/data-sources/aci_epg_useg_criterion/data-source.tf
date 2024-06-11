
data "aci_epg_useg_criterion" "example_application_epg" {
  parent_dn = aci_application_epg.example.id
}
