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

resource "aci_tacacs_accounting_destination" "example" {
  tacacs_monitoring_destination_group_dn  = aci_tacacs_accounting.example.id
  host = "cisco.com"
  port = "49"
  annotation = "orchestrator:terraform"
  auth_protocol = "pap"
  key = "example_key_value"
  description = "from terraform"
}