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

resource "aci_tacacs_accounting" "example" {
  name        = "example"
  annotation  = "orchestrator:terraform"
  name_alias  = "tacacs_accounting_alias"
  description = "From Terraform"
}