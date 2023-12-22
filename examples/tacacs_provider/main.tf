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

resource "aci_tacacs_provider" "example" {
  name                = "example"
  annotation          = "orchestrator:terraform"
  name_alias          = "tacacs_provider_alias"
  description         = "From Terraform"
  auth_protocol       = "pap"
  key                 = "example"
  monitor_server      = "disabled"
  monitoring_password = "example"
  monitoring_user     = "default"
  port                = "49"
  retries             = "1"
  timeout             = "5"
}