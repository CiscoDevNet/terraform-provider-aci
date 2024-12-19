
resource "aci_dscp_to_priority_map" "example_custom_qos_policy" {
  parent_dn = aci_custom_qos_policy.example.id
  from      = "AF11"
  to        = "AF22"
}
