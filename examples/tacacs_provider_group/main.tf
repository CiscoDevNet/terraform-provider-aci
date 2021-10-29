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

resource "aci_tacacs_provider_group" "example" {
  name  = "example"
  annotation = "orchestrator:terraform"
  description = "from terraform"
  name_alias = "example_name_alias"
}