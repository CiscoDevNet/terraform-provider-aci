
resource "aci_netflow_exporter_policy" "full_example_tenant" {
  parent_dn      = aci_tenant.example.id
  annotation     = "annotation"
  description    = "description"
  dscp           = "AF11"
  dst_addr       = "2.2.2.1"
  dst_port       = "https"
  name           = "netfow_exporter"
  name_alias     = "name_alias"
  owner_key      = "owner_key"
  owner_tag      = "owner_tag"
  source_ip_type = "custom-src-ip"
  src_addr       = "1.1.1.1/10"
  ver            = "v9"
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
