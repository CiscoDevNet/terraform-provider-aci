
resource "aci_l3out_management_network_subnet" "example" {
  parent_dn = aci_l3out_management_network_instance_profile.example.id
  ip        = "1.1.1.0/24"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}

