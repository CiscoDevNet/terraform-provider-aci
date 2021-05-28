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
resource "aci_leaf_access_bundle_policy_group" "fooleaf_access_bundle_policy_group" {
		description = "From Terraform"
		name        = "demo_if_pol_grp"
		annotation  = "tag_if_pol"
		lag_t       = "link"
		name_alias  = "alias_if_pol"
	} 