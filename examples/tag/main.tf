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

data "aci_tag" "tag" {
  parent_dn = "uni/tn-common"
  key       = "terraform-test"
}

resource "aci_tenant" "tenant_tag" {
  name        = "tenant_tag"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_tag" "tag" {
  parent_dn = aci_tenant.tenant_tag.id
  key       = "example"
  value     = "example-value"
}