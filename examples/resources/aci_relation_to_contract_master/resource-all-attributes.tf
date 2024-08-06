
resource "aci_relation_to_contract_master" "full_example_application_epg" {
  parent_dn  = aci_application_epg.example.id
  annotation = "annotation"
  target_dn  = aci_application_epg.example_application_epg.id
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

resource "aci_relation_to_contract_master" "full_example_endpoint_security_group" {
  parent_dn  = aci_endpoint_security_group.example.id
  annotation = "annotation"
  target_dn  = aci_endpoint_security_group.example_endpoint_security_group_2.id
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
