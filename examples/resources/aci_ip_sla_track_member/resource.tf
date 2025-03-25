
resource "aci_ip_sla_track_member" "example_tenant" {
  parent_dn              = aci_tenant.example.id
  destination_ip_address = "1.1.1.1"
  name                   = "test_name"
  scope                  = "uni/tn-test_tenant/BD-test_bd"
}
