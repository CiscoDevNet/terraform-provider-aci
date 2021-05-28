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

	resource "aci_tenant" "footenant" {
		description = "sample aci_tenant from terraform"
		name        = "demo_tenant"
		annotation  = "tag_tenant"
		name_alias  = "alias_tenant"
	  }
	  
	resource "aci_l3_outside" "fool3_outside" {
        tenant_dn      = "${aci_tenant.footenant.id}"
        description    = "sample aci_l3_outside"
        name           = "demo_l3out"
        annotation     = "tag_l3out"
        enforce_rtctrl = ["export", "import"]
        name_alias     = "alias_out"
        target_dscp    = "unspecified"
    }

	resource "aci_logical_node_profile" "foological_node_profile" {
        l3_outside_dn = "${aci_l3_outside.fool3_outside.id}"
        description   = "sample logical node profile"
        name          = "demo_node"
        annotation    = "tag_node"
        config_issues = "none"
        name_alias    = "alias_node"
        tag           = "black"
        target_dscp   = "unspecified"
      }
	
	resource "aci_logical_interface_profile" "foological_interface_profile" {
		logical_node_profile_dn = "${aci_logical_node_profile.foological_node_profile.id}"
		description             = "aci_logical_interface_profile from terraform"
		name                    = "demo_int_prof"
		annotation              = "tag_prof"
		name_alias              = "alias_prof"
		prio                    = "unspecified"
		tag                     = "black"
	  }	  
