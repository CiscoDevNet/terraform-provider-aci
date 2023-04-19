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

data "aci_client_end_point" "this" {
  ip                 = "10.1.1.33"
  allow_empty_result = true
}

output "mac_endpoints" {
  value = data.aci_client_end_point.this.fvcep_objects
}