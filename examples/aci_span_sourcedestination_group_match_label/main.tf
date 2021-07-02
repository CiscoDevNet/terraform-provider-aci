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

resource "aci_span_source_group" "foospan_source_group" {
  tenant_dn = aci_tenant.test_tenant1.id
  name = "spanSrcGrp1"
}

resource "aci_span_sourcedestination_group_match_label" "example" {
  span_source_group_dn  = aci_span_source_group.foospan_source_group.id
  description           = "From Terraform"
  name                  = "example"
  annotation            = "tag_label"
  name_alias            = "alias_label"
  tag                   = "yellow-green"
}