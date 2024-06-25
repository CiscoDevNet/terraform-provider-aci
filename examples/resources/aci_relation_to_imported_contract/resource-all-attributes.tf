
resource "aci_relation_to_imported_contract" "full_example_application_epg" {
  parent_dn              = aci_application_epg.example.id
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

resource "aci_relation_to_imported_contract" "full_example_endpoint_security_group" {
  parent_dn              = aci_endpoint_security_group.example.id
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
