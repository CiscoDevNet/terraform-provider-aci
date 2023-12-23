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

# Do not use the `multicast_address_block` argument from the `aci_multicast_pool` resource in combination with the `aci_multicast_pool_block` resource!

# Example of using the `aci_multicast_pool` resource with the `multicast_address_block` argument:

resource "aci_multicast_pool" "test-tf-pool-1-resource" {
  name        = "test-tf-pool-1-resource"
  description = "This multicast pool is created by terraform with multicast_address_block argument"
  multicast_address_block {
    from = "224.0.0.40"
    to   = "224.0.0.44"
    name = "testing-1"
  }
  multicast_address_block {
    from = "224.0.0.50"
    to   = "224.0.0.54"
    name = "testing-2"
  }
}

# Example of using the `aci_multicast_pool` resource with the `aci_multicast_pool_block` resource:

resource "aci_multicast_pool" "test-tf-pool-X-resources" {
  name        = "test-tf-pool-X-resources"
  description = "This multicast pool is created by terraform to be used by aci_multicast_pool_block resource"
}

resource "aci_multicast_pool_block" "test-tf-pool-block-1" {
  multicast_pool_dn = aci_multicast_pool.test-tf-pool-X-resources.id
  name              = "test-tf-pool-block-1"
  description       = "This multicast block is created by terraform using the aci_multicast_pool resource"
  from              = "224.0.0.0"
  to                = "224.0.0.10"
}

resource "aci_multicast_pool_block" "test-tf-pool-block-2" {
  multicast_pool_dn = aci_multicast_pool.test-tf-pool-X-resources.id
  name              = "test-tf-pool-block-2"
  description       = "This multicast block is created by terraform using the aci_multicast_pool resource"
  from              = "224.0.0.11"
  to                = "224.0.0.20"
}
