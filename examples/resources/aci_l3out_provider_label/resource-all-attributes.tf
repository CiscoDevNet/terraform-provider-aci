
resource "aci_l3out_provider_label" "full_example_l3_outside" {
  parent_dn = aci_l3_outside.example.id
  annotation = "annotation"
  description = "description"
  name = "prov_label"
  name_alias = "name_alias"
  owner_key = "owner_key"
  owner_tag = "owner_tag"
  tag = "alice-blue"
  annotations = [
    {
      key = "key_0"
      value = "value_1"
    }
  ]
  tags = [
    {
      key = "key_0"
      value = "value_1"
    }
  ]
}
