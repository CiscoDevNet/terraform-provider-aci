provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_ospf_route_summarization" "example" {

  tenant_dn          = aci_tenant.example.id
  name               = "example"
  annotation         = "example"
  cost               = "1"
  inter_area_enabled = "no"
  name_alias         = "example"
  tag                = "1"

}
