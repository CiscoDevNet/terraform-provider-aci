
resource "aci_l3out_management_network_oob_contract" "example" {
  parent_dn                 = aci_l3out_management_network_instance_profile.example.id
  out_of_band_contract_name = "test_tn_vz_oob_br_cp_name"
}

