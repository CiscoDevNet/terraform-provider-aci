
data "aci_relation_to_netflow_exporter" "example_netflow_monitor_policy" {
  parent_dn                    = aci_netflow_monitor_policy.example.id
  netflow_exporter_policy_name = aci_netflow_exporter_policy.example.name
}
