
resource "aci_relation_to_provided_contract" "full_example_application_epg" {
  parent_dn      = aci_application_epg.example.id
  annotation     = "annotation"
  match_criteria = "All"
  priority       = "level1"
  contract_name  = aci_contract.example.name
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

resource "aci_relation_to_provided_contract" "full_example_endpoint_security_group" {
  parent_dn      = aci_endpoint_security_group.example.id
  annotation     = "annotation"
  match_criteria = "All"
  priority       = "level1"
  contract_name  = aci_contract.example.name
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
