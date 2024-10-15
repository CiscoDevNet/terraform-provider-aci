
data "aci_rogue_coop_exception" "example_bridge_domain" {
  parent_dn = aci_bridge_domain.example.id
  mac       = "00:00:00:00:00:01"
}
