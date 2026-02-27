
resource "aci_virtual_ip_address" "full_example_application_epg" {
  parent_dn   = aci_application_epg.example.id
  ip          = "1.1.1.4"
  annotation  = "annotation"
  description = "description_1"
  name        = "name_1"
  name_alias  = "name_alias_1"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
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
