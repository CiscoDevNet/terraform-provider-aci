
resource "aci_relation_from_any_to_consumer_contract" "full_example_any" {
  parent_dn     = aci_any.example.id
  annotation    = "annotation"
  priority      = "level1"
  contract_name = aci_contract.example.name
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
