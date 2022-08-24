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

# Do not use the `multicast_address_block` from this resource in combination with the `aci_multicast_pool_block` resource.

resource "aci_multicast_pool" "test-tf-pool" {
  name        = "test-tf-pool"
  description = "This multicast pool is created by terraform"
  multicast_address_block {
    from = "224.0.0.40"
    to = "224.0.0.44"
    name = "testing-1"
  }
  multicast_address_block {
    from = "224.0.0.50"
    to = "224.0.0.54"
    name = "testing-2"
  }
}
