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

resource "aci_taboo_contract" "example" {
  tenant_dn   = aci_tenant.example.id
  name        = "example_contract"
  description = "from terraform"
  annotation  = "orchestrator:terraform"
  name_alias  = "example"
}
