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

resource "aci_multicast_pool" "test-tf-pool" {
  name        = "test-tf-pool"
  description = "This multicast pool is created by terraform"
}
