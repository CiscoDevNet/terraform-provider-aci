
resource "aci_relation_from_bridge_domain_to_netflow_monitor_policy" "full_example_bridge_domain" {
  parent_dn                   = aci_bridge_domain.example.id
  annotation                  = "annotation"
  filter_type                 = "ipv4"
  netflow_monitor_policy_name = aci_netflow_monitor_policy.example.name
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
