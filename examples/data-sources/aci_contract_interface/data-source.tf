
data "aci_contract_interface" "example" {
  parent_dn               = aci_application_epg.example.id
  contract_interface_name = "test_tn_vz_cp_if_name"
}

