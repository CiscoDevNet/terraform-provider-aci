
resource "aci_netflow_record_policy" "full_example_tenant" {
  parent_dn          = aci_tenant.example.id
  annotation         = "annotation"
  collect_parameters = ["count-bytes", "src-intf"]
  description        = "description_1"
  match_parameters   = ["dst-ip", "src-ip"]
  name               = "netfow_record"
  name_alias         = "name_alias_1"
  owner_key          = "owner_key_1"
  owner_tag          = "owner_tag_1"
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
