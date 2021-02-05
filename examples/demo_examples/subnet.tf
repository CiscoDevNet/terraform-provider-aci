resource "aci_tenant" "tenant_for_subnet" {
  name        = "tenant_for_subnet"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_bridge_domain" "bd_for_subnet" {
  tenant_dn   = aci_tenant.tenant_for_subnet.id
  name        = "bd_for_subnet"
  description = "This bridge domain is created by terraform ACI provider"
  mac         = "00:22:BD:F8:19:FF"
}

resource "aci_subnet" "demosubnet" {
  name             = "10.0.3.28/27"
  bridge_domain_dn = aci_bridge_domain.bd_for_subnet.id
  ip_address       = "10.0.3.28/27"
  scope            = ["private"]
  description      = "This subject is created by terraform"
}
