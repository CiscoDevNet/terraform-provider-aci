
resource "aci_relation_to_consumed_out_of_band_contract" "full_example_external_management_network_instance_profile" {
  parent_dn                 = aci_external_management_network_instance_profile.example.id
  annotation                = "annotation"
  priority                  = "level1"
  out_of_band_contract_name = aci_out_of_band_contract.example.name
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
