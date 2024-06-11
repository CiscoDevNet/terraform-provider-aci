
resource "aci_epg_useg_sub_criterion" "example_epg_useg_criterion" {
  parent_dn = aci_epg_useg_criterion.example.id
  name      = "sub_criterion"
}
