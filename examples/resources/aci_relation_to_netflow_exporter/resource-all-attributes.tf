
resource "aci_relation_to_netflow_exporter" "full_example_netflow_monitor_policy" {
  parent_dn                    = aci_netflow_monitor_policy.example.id
  annotation                   = "annotation"
  netflow_exporter_policy_name = aci_netflow_exporter_policy.example.name
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
}
