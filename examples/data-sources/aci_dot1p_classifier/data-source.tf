
data "aci_dot1p_classifier" "example_custom_qos_policy" {
  parent_dn = aci_custom_qos_policy.example.id
  from      = "1"
  to        = "2"
}
