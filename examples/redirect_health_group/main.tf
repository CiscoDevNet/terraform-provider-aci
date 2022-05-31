terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

resource "aci_tenant" "terraform_tenant" {
  name        = "tf_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_l4_l7_redirect_health_group" "example" {
  tenant_dn  = aci_tenant.terraform_tenant.id
  name       = "example"
}
