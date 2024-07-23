
resource "aci_l3out_consumer_label" "full_example_l3_outside" {
  parent_dn = aci_l3_outside.example.id
  annotation = "annotation"
  description = "description"
  name = "test_name"
  name_alias = "name_alias"
  owner = "infra"
  owner_key = "owner_key"
  owner_tag = "owner_tag"
  tag = "lemon-chiffon"
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
