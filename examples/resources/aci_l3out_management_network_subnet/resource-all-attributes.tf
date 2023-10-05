
resource "aci_l3out_management_network_subnet" "example" {
  parent_dn   = aci_l3out_management_network_instance_profile.example.id
  annotation  = "annotation"
  description = "description"
  ip          = "1.1.1.0/24"
  name        = "name"
  name_alias  = "name_alias"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}

