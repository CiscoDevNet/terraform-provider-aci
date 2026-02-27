
data "aci_virtual_ip_address" "example_application_epg" {
  parent_dn = aci_application_epg.example.id
  ip        = "1.1.1.4"
}
