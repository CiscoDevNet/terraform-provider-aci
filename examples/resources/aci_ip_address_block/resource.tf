
resource "aci_ip_address_block" "example_ip_address_pool" {
  parent_dn = aci_ip_address_pool.example.id
  from      = "10.0.0.2"
  to        = "10.0.0.5"
}
