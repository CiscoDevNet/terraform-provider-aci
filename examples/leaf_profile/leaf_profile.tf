resource "aci_leaf_profile" "tf_leaf_prof" {
  name        = "tf_leaf_prof"
  description = "From Terraform"
  annotation  = "example"
  name_alias  = "example"
  leaf_selector {
    name                    = "one"
    switch_association_type = "range"
    node_block {
      name  = "blk1"
      from_ = "105"
      to_   = "106"
    }
    node_block {
      name  = "blk2"
      from_ = "102"
      to_   = "104"
    }
  }
  leaf_selector {
    name                    = "two"
    switch_association_type = "range"
    node_block {
      name  = "blk3"
      from_ = "105"
      to_   = "106"
    }
  }

}
