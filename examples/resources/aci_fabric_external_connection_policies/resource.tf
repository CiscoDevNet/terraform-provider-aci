
resource "aci_fabric_external_connection_policies" "example_tenant" {
  parent_dn    = aci_tenant.example.id
  id_attribute = "1"
}
