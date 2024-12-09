
resource "aci_netflow_exporter_policy" "full_example_tenant" {
  parent_dn              = aci_tenant.example.id
  annotation             = "annotation"
  description            = "description_1"
  qos_dscp_value         = "AF11"
  destination_ip_address = "2.2.2.1"
  destination_port       = "https"
  name                   = "netfow_exporter"
  name_alias             = "name_alias_1"
  owner_key              = "owner_key_1"
  owner_tag              = "owner_tag_1"
  source_ip_type         = "custom-src-ip"
  source_ip_address      = "1.1.1.1/10"
  version                = "v9"
  relation_to_vrf = {
    annotation = "annotation_1"
    target_dn  = aci_vrf.example.id
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
  relation_to_epg = {
    annotation = "annotation_1"
    target_dn  = aci_application_epg.example.id
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
