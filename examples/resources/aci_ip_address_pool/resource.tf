
resource "aci_ip_address_pool" "example" {
  gateway_address = "10.0.0.1/24"
  name            = "test_name"
}
