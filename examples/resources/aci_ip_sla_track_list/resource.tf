
resource "aci_ip_sla_track_list" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
