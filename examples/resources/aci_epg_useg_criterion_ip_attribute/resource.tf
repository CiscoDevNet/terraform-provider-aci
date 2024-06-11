
resource "aci_epg_useg_criterion_ip_attribute" "example_epg_useg_criterion" {
  parent_dn = aci_epg_useg_criterion.example.id
  ip        = "131.107.1.200"
  name      = "131"
}
