
resource "aci_relation_to_consumed_out_of_band_contract" "example_external_management_network_instance_profile" {
  parent_dn                 = aci_external_management_network_instance_profile.example.id
  out_of_band_contract_name = aci_out_of_band_contract.example.name
}
