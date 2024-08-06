
resource "aci_relation_to_contract_master" "full_example_endpoint_security_group" {
  parent_dn  = aci_endpoint_security_group.example.id
  annotation = "annotation"
  target_dn  = aci_endpoint_security_group.example_endpoint_security_group.id
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
