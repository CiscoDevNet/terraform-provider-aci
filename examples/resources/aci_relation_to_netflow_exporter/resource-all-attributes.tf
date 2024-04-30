
resource "aci_relation_to_netflow_exporter" "full_example_netflow_monitor_policy" {
  parent_dn                    = aci_netflow_monitor_policy.example.id
  annotation                   = "annotation"
  tn_netflow_exporter_pol_name = aci_.example.name
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
