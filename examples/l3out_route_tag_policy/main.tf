
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_l3out_route_tag_policy" "example" {

  tenant_dn  = "${aci_tenant.example.id}"
  name  = "example"
  annotation  = "example"
  name_alias  = "example"
  tag  = "1"
  description = "from terraform"

}
