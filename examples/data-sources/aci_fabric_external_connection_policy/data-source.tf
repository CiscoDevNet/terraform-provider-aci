
data "aci_fabric_external_connection_policy" "example_tenant" {
  parent_dn = aci_tenant.example.id
  fabric_id = "1"
}
