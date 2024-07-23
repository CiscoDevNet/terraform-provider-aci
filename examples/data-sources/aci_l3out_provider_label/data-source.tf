
data "aci_l3out_provider_label" "example_l3_outside" {
  parent_dn = aci_l3_outside.example.id
  name = "prov_label"
}
