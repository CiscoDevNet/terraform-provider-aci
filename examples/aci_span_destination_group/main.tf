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

resource "aci_tenant" "test_tenant1" {
  name        = "tf_test_rel_tenant2"
  description = "This tenant is created by terraform"
}

resource "aci_span_destination_group" "foospan_destination_group" {
  tenant_dn = aci_tenant.test_tenant1.id
  name      = "spanDestGrp1"
}