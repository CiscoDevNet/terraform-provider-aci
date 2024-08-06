
resource "aci_external_management_network_subnet" "full_example_external_management_network_instance_profile" {
  parent_dn   = aci_external_management_network_instance_profile.example.id
  annotation  = "annotation"
  description = "description_1"
  ip          = "1.1.1.0/24"
  name        = "name_1"
  name_alias  = "name_alias_1"
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
