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

# provider "aci" {
#   username    = ""
#   private_key = ""
#   cert_name   = ""
#   url         = ""
#   insecure    = true
# }

