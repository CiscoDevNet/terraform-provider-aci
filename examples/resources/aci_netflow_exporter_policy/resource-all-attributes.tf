
resource "aci_netflow_exporter_policy" "full_example_tenant" {
  parent_dn           = aci_tenant.example.id
  annotation          = "annotation"
  description         = "description_1"
  dscp                = "AF11"
  destination_address = "2.2.2.1"
  destination_port    = "https"
  name                = "netfow_exporter"
  name_alias          = "name_alias_1"
  owner_key           = "owner_key_1"
  owner_tag           = "owner_tag_1"
  source_ip_type      = "custom-src-ip"
  source_address      = "1.1.1.1/10"
  version             = "v9"
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
