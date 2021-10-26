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
resource "aci_leaf_access_port_policy_group" "test_port_grp" {
  name        = "tf_port_grp"
  description = "From Terraform"
  annotation  = "tag_ports"
  name_alias  = "name_alias"

}
