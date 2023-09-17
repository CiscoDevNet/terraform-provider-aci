
data "aci_l3out_management_network_oob_contract" "example" {
  parent_dn     = aci_l3out_management_network_instance_profile.example.id
  contract_name = "test_tn_vz_oob_br_cp_name"
}

