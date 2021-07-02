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

resource "aci_monitoring_policy" "example" {
  tenant_dn   = aci_tenant.example.id
  description = "From Terraform"
  name        = "example"
  annotation  = "example"
  name_alias  = "example"
}
