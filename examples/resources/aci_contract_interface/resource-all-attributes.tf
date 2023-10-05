
resource "aci_contract_interface" "example" {
  parent_dn               = aci_application_epg.example.id
  annotation              = "annotation"
  priority                = "level1"
  contract_interface_name = "test_tn_vz_cp_if_name"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}

