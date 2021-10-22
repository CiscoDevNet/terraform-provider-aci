
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_application_profile" "test_ap" {
  tenant_dn   = aci_tenant.dev_tenant.id
  name        = "demo_ap"
  annotation  = "tag"
  description = "from terraform"
  name_alias  = "test_ap"
  prio        = "level1"
}