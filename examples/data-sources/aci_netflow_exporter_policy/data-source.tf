
data "aci_netflow_exporter_policy" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "netfow_exporter"
}
