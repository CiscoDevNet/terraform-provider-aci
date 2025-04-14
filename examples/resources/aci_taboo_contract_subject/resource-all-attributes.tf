
resource "aci_taboo_contract_subject" "full_example_taboo_contract" {
  parent_dn   = aci_taboo_contract.example.id
  annotation  = "annotation"
  description = "description_1"
  name        = "test_name"
  name_alias  = "name_alias_1"
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
  relation_to_filters = [
    {
      annotation  = "annotation_1"
      directives  = ["log", "no_stats"]
      filter_name = aci_filter.example.name
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
