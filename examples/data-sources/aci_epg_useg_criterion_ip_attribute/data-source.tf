
data "aci_epg_useg_criterion_ip_attribute" "example_epg_useg_criterion" {
  parent_dn = aci_epg_useg_criterion.example.id
  name      = "131"
}
