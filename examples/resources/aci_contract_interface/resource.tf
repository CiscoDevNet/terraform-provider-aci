
resource "aci_contract_interface" "example" {
  parent_dn               = aci_application_epg.example.id
  contract_interface_name = "test_contract_interface"
  annotations = [
    {
      key = "test_annotation"
    },
  ]
}

