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

resource "aci_radius_provider" "example" {
  name  = "example"
  type = "radius"
  annotation = "orchestrator:terraform"
  auth_port = "1812"
  auth_protocol = "pap"
  key = "example_key_value"
  monitor_server = "disabled"
  monitoring_password = "example_monitoring_password"
  monitoring_user = "default"
  retries = "1"
  timeout = "5"
  description = "from terraform"
  name_alias = "example_name_alias_value"
}