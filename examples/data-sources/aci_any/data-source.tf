
data "aci_any" "example_vrf" {
  parent_dn = aci_vrf.example.id
}
