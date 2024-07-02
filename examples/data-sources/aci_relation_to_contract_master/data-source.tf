
data "aci_relation_to_contract_master" "example_application_epg" {
  parent_dn = aci_application_epg.example.id
  target_dn = aci_application_epg.example_2.id
}

data "aci_relation_to_contract_master" "example_endpoint_security_group" {
  parent_dn = aci_endpoint_security_group.example.id
  target_dn = aci_endpoint_security_group.example_2.id
}
