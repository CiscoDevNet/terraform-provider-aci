
resource "aci_ip_sla_track_member" "full_example_tenant" {
  parent_dn              = aci_tenant.example.id
  annotation             = "annotation"
  description            = "description_1"
  destination_ip_address = "1.1.1.1"
  name                   = "test_name"
  name_alias             = "name_alias_1"
  owner_key              = "owner_key_1"
  owner_tag              = "owner_tag_1"
  scope                  = "uni/tn-test_tenant/BD-test_bd"
  relation_to_monitoring_policy = {
    annotation = "annotation_1"
    target_dn  = aci_ip_sla_monitoring_policy.example.id
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
