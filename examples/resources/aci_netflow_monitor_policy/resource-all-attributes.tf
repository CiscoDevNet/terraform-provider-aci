
resource "aci_netflow_monitor_policy" "full_example_tenant" {
  parent_dn   = aci_tenant.example.id
  annotation  = "annotation"
  description = "description_1"
  name        = "netfow_monitor"
  name_alias  = "name_alias_1"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  relation_to_netflow_exporters = [
    {
      annotation                   = "annotation_1"
      netflow_exporter_policy_name = aci_netflow_exporter_policy.example.name
    }
  ]
  relation_to_netflow_record = {
    annotation                 = "annotation_1"
    netflow_record_policy_name = aci_netflow_record_policy.example.name
  }
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
