
resource "aci_relation_to_consumed_contract_interface" "example_application_epg" {
  parent_dn               = aci_application_epg.example.id
  contract_interface_name = aci_contract_interface.example.name
}

resource "aci_relation_to_consumed_contract_interface" "example_endpoint_security_group" {
  parent_dn               = aci_endpoint_security_group.example.id
  contract_interface_name = aci_contract_interface.example.name
}
