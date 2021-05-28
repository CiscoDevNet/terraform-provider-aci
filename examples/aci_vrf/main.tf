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

resource "aci_tenant" "tenant_for_vrf" {
	name        = "tenant_for_vrf"
	description = "This tenant is created by terraform ACI provider"
}

resource "aci_vrf" "foovrf" {
	tenant_dn   		   = aci_tenant.tenant_for_vrf.id
	description 		   = "%s"
	name                   = "demo_vrf"
	annotation             = "tag_vrf"
	bd_enforced_enable     = "no"
	ip_data_plane_learning = "enabled"
	knw_mcast_act          = "permit"
	name_alias             = "alias_vrf"
	pc_enf_dir             = "egress"
	pc_enf_pref            = "unenforced"
	}