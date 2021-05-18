
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}
resource "aci_bgp_best_path_policy" "foobgp_best_path_policy" {
  tenant_dn  = aci_tenant.example.id
  name       = "example"
  annotation = "example"
  ctrl       = "asPathMultipathRelax"
  name_alias = "example"
}
