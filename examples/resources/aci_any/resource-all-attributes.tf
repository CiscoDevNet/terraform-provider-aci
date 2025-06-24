
resource "aci_any" "full_example_vrf" {
  parent_dn              = aci_vrf.example.id
  annotation             = "annotation"
  description            = "description_1"
  match_criteria         = "All"
  name                   = "name_1"
  name_alias             = "name_alias_1"
  preferred_group_member = "disabled"
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
  relation_from_any_to_consumer_contracts = [
    {
      annotation    = "annotation_1"
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
  ]
  relation_from_any_to_contract_interfaces = [
    {
      annotation             = "annotation_1"
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
  ]
  relation_from_any_to_provider_contracts = [
    {
      annotation     = "annotation_1"
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
  ]
}
