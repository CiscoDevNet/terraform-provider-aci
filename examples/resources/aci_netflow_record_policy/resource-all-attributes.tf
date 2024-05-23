
resource "aci_netflow_record_policy" "full_example_tenant" {
  parent_dn   = aci_tenant.example.id
  annotation  = "annotation"
  collect     = "count-bytes"
  description = "description"
  match       = "dst-ip"
  name        = "netfow_record"
  name_alias  = "name_alias"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
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
