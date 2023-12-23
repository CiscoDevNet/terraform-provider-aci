terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "admin"
  password = "password"
  url      = "https://my-cisco-aci.com"
  insecure = true
}

resource "aci_tenant" "example" {
  name = "example_tenant"
}