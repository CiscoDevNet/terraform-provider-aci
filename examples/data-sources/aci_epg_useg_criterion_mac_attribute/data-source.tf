
data "aci_epg_useg_criterion_mac_attribute" "example_epg_useg_criterion" {
  parent_dn = aci_epg_useg_criterion.example.id
  name      = "mac_attr"
}
