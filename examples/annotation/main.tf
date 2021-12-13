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

data "aci_annotation" "annotation" {
  parent_dn = "uni/tn-common"
  key       = "terraform-test"
}

resource "aci_tenant" "tenant_annotation" {
  name        = "tenant_annotation"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_annotation" "annotation" {
  parent_dn = aci_tenant.tenant_annotation.id
  key       = "example"
  value     = "example-value"
}