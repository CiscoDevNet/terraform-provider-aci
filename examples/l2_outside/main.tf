terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_tenant" "foo_tenant" {
  name        = "tenant_1"
  description = "This tenant is created by terraform"
}

resource "aci_l2_outside" "fool2_outside" {
  tenant_dn   = aci_tenant.foo_tenant.id
  description = "from terraform"
  name        = "l2_outside_1"
  annotation  = "l2_outside_tag"
  name_alias  = "example"
  target_dscp = "AF11"
}

data "aci_l2_outside" "example9" {
  tenant_dn = aci_tenant.foo_tenant.id
  name      = aci_l2_outside.fool2_outside.name
}

output "name9" {
  value = data.aci_l2_outside.example9
}