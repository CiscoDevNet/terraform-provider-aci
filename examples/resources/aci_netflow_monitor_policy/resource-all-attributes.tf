
resource "aci_netflow_monitor_policy" "full_example_tenant" {
  parent_dn   = aci_tenant.example.id
  annotation  = "annotation"
  description = "description"
  name        = "netfow_monitor"
  name_alias  = "name_alias"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
  relation_to_netflow_exporters = [
    {
      annotation                   = "annotation_1"
      tn_netflow_exporter_pol_name = aci_.example.name
    }
  ]
  relation_to_netflow_record = [
    {
      annotation                 = "annotation_1"
      tn_netflow_record_pol_name = aci_.example.name
    }
  ]
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
