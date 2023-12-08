
resource "aci_external_management_network_subnet" "example" {
  parent_dn = aci_external_management_network_instance_profile.example.id
  ip        = "1.1.1.0/24"
}
  