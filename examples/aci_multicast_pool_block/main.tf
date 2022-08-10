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

resource "aci_multicast_pool_block" "test-tf-pool-block-1" {
  multicast_pool_dn = aci_multicast_pool.test-tf-pool.id
  name                      = "test-tf-block-1"
  description               = "This multicast block is created by terraform"
  from                      = "224.0.0.0"
  to                        = "224.0.0.10"
}

resource "aci_multicast_pool_block" "test-tf-pool-block-2" {
  multicast_pool_dn         = aci_multicast_pool.test-tf-pool.id
  description               = "This multicast block is created by terraform"
  from                      = "224.0.0.11"
  to                        = "224.0.0.20"
}
