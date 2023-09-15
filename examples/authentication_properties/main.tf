terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

#configure provider with your cisco aci credentials.
provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

resource "aci_authentication_properties" "example" {
  annotation      = "orchestrator:terraform"
  name_alias      = "example_name_alias"
  description     = "from terraform"
  def_role_policy = "no-login"
  ping_check      = "true"
  retries         = "1"
  timeout         = "5"
}
