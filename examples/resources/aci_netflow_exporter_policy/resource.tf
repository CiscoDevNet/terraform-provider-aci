
resource "aci_netflow_exporter_policy" "example_tenant" {
  parent_dn           = aci_tenant.example.id
  destination_address = "2.2.2.1"
  destination_port    = "https"
  name                = "netfow_exporter"
}
