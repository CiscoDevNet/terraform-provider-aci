
resource "aci_netflow_record_policy" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "netfow_record"
}
