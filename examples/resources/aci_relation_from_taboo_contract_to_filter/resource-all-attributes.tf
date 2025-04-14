
resource "aci_relation_from_taboo_contract_to_filter" "full_example_taboo_contract_subject" {
  parent_dn   = aci_taboo_contract_subject.example.id
  annotation  = "annotation"
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
