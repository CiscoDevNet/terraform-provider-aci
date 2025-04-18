
resource "aci_taboo_contract_subject" "example_taboo_contract" {
  parent_dn = aci_taboo_contract.example.id
  name      = "test_name"
}
