resource "aci_l3out_management_network_instance_profile" "example" {
  name = "test_l3out_management_network_instance_profile"
  l3out_management_network_oob_contracts = [
    {
      contract_name = "test_l3out_management_network_contract"
    },
  ]
}
