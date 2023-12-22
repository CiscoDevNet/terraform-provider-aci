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

resource "aci_fabric_node_control" "example" {
  name        = "example"
  annotation  = "orchestrator:terraform"
  control     = "Dom"
  feature_sel = "telemetry"
  name_alias  = "example_name_alias"
  description = "from terraform"
}