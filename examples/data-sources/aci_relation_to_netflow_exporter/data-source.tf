
data "aci_relation_to_netflow_exporter" "example_netflow_monitor_policy" {
  parent_dn                    = aci_netflow_monitor_policy.example.id
  tn_netflow_exporter_pol_name = aci_.example.name
}
