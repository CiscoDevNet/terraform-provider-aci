
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

resource "aci_vpc_explicit_protection_group" "example" {
  name                              = "example"
  annotation                        = "tag_vpc"
  switch1                           = "switch1_id"
  switch2                           = "switch2_id"
  vpc_domain_policy                 = "test"
  vpc_explicit_protection_group_id  = "1"
}