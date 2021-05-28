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
resource "aci_l2_interface_policy" "test_l2" {
        name = "tf_l2"
        description = "From Terraform"
		name        = "demo_l2_pol"
		annotation  = "tag_l2_pol"
		name_alias  = "alias_l2_pol"
		qinq        = "disabled"
		vepa        = "disabled"
		vlan_scope  = "global"
  
}
