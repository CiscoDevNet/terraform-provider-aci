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

resource "aci_fabric_if_pol" "foolink_level_policy" {
		description = "from terraform"
		name  = "fabric_if_pol_1"
		annotation  = "fabric_if_pol_tag"
		auto_neg  = "on"
		fec_mode  = "kp-fec"
		link_debounce  = "100"
		name_alias  = "example"
		speed  = "inherit"
	}

data "aci_fabric_if_pol" "example5"{
  name  = aci_fabric_if_pol.foolink_level_policy.name
}

output "name5" {
  value = data.aci_fabric_if_pol.example5
}