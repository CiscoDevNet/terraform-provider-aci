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

resource "aci_tenant" "example" {
  name        = "example"
  description = "This tenant is created by terraform"
}

resource "aci_span_source_group" "example" {
  tenant_dn   = aci_tenant.example.id
  name        = "example"
  admin_st    = "enabled"
  annotation  = "tag_span"
  description = "from terraform"
  name_alias  = "alias_span"
}