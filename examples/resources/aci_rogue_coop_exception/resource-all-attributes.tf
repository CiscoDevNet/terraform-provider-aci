
resource "aci_rogue_coop_exception" "full_example_bridge_domain" {
  parent_dn   = aci_bridge_domain.example.id
  annotation  = "annotation"
  description = "description_1"
  mac         = "00:00:00:00:00:01"
  name        = "name_1"
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
}
