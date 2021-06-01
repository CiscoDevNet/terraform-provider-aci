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

resource "aci_l3_domain_profile" "fool3_domain_profile" {
	  name  = "l3_domain_profile"
		annotation  = "l3_domain_profile_tag"
		name_alias  = "alias_name"
}