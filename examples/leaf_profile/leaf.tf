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
resource "aci_leaf_interface_profile" "test_leaf_profile" {
  description = "from terraform"
  name        = "demo_leaf_profile"
  annotation  = "tag_leaf"
  name_alias  = "name_alias"
}
