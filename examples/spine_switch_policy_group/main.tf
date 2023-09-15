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

resource "aci_spine_switch_policy_group" "example" {
  name        = "example"
  annotation  = "orchestrator:terraform"
  name_alias  = "example"
  description = "from terraform"

}