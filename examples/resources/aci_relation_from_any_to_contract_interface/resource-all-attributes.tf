
resource "aci_relation_from_any_to_contract_interface" "full_example_any" {
  parent_dn              = aci_any.example.id
  annotation             = "annotation"
  priority               = "level1"
  imported_contract_name = aci_imported_contract.example.name
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
