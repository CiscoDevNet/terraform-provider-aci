resource "aci_l3out_management_network_instance_profile" "example" {
  name = "test_name"
  l3out_management_network_oob_contracts = [
    {
      contract_name = "l3out_management_network_oob_contracts_1"
    }
  ]
}
