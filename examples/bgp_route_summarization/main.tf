provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_tenant" "tenentcheck" {
  name       = "test"
  annotation = "atag"
  name_alias = "alias_tenant"
}

resource "aci_bgp_route_summarization" "example" {

  tenant_dn  = aci_tenant.example.id
  name       = "example"
  annotation = "example"
  attrmap    = "example"
  ctrl       = "as-set"
  name_alias = "example"

}
