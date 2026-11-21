
resource "aci_macsec_key" "example_macsec_key_chain" {
  parent_dn = aci_macsec_key_chain.example.id
  key_name  = "aa"
}
