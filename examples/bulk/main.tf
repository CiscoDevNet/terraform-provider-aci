terraform {
  required_providers {
    aci = {
      #source =  "local/providers/aci"
      source =  "ciscodevnet/aci"
      #version = "2.6.1"
    }
  }
  required_version = ">= 0.13"
}

provider "aci" {
  username = ""
  password = ""
  url      = "http://"
  insecure = true
}
