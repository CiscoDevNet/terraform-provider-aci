
data "aci_relation_from_bridge_domain_to_netflow_monitor_policy" "example_bridge_domain" {
  parent_dn                   = aci_bridge_domain.example.id
  filter_type                 = "ipv4"
  netflow_monitor_policy_name = aci_netflow_monitor_policy.example.name
}
