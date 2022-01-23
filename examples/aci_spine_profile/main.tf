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

resource "aci_spine_profile" "foospine_profile" {
  name        = "spine_profile_1"
  description = "from terraform"
  annotation  = "spine_profile_tag"
  name_alias  = "example"
  spine_selector {
    name                    = "one"
    switch_association_type = "range"
    node_block {
      name  = "blk1"
      from_ = "101"
      to_   = "102"
    }
    node_block {
      name  = "blk2"
      from_ = "103"
      to_   = "104"
    }
  }
  spine_selector {
    name                    = "two"
    switch_association_type = "range"
    node_block {
      name  = "blk3"
      from_ = "105"
      to_   = "106"
    }
  }
  relation_infra_rs_sp_acc_port_p = [
    aci_spine_interface_profile.example.id
  ]
}

resource "aci_spine_interface_profile" "example" {
  name        = "example"
  description = "from terraform"
  annotation  = "example"
  name_alias  = "example"
}