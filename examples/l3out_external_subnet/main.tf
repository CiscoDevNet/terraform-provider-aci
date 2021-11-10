terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

#configure provider with your cisco aci credentials.
provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

resource "aci_tenant" "tenant01" {
  name        = "tenant01"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_l3_domain_profile" "l3_domain_profile" {
  name = "l3extDomp1"
}

resource "aci_vrf" "vrf1" {
  tenant_dn          = aci_tenant.tenant01.id
  bd_enforced_enable = "no"
  knw_mcast_act      = "permit"
  name               = "vrf1"
  pc_enf_dir         = "ingress"
  pc_enf_pref        = "enforced"
}

resource "aci_l3_outside" "l3_outside" {
  tenant_dn                    = aci_tenant.tenant01.id
  name                         = "demo_l3out"
  annotation                   = "tag_l3out"
  name_alias                   = "alias_out"
  target_dscp                  = "unspecified"
  relation_l3ext_rs_ectx       = aci_vrf.vrf1.id
  relation_l3ext_rs_l3_dom_att = aci_l3_domain_profile.l3_domain_profile.id
}


resource "aci_external_network_instance_profile" "l3out_extepg" {
  l3_outside_dn = aci_l3_outside.l3_outside.id
  description   = "ExtEpg for Terraform l3out"
  name          = "l3out_extepg"
}

resource "aci_bgp_route_control_profile" "bgp_route_control_profile" {
  parent_dn = aci_tenant.tenant01.id
  name      = "bgp_route_control_profile"
}

resource "aci_l3_ext_subnet" "l3out_extepg_subnet" {
  external_network_instance_profile_dn = aci_external_network_instance_profile.l3out_extepg.id
  ip                                   = "10.10.10.10/24"
  aggregate                            = "shared-rtctrl"
  annotation                           = "tag_ext_subnet"
  name_alias                           = "alias_ext_subnet"
  scope                                = ["import-rtctrl", "export-rtctrl", "import-security"]
  relation_l3ext_rs_subnet_to_profile {
    tn_rtctrl_profile_dn = aci_bgp_route_control_profile.bgp_route_control_profile.id
    direction            = "import"
  }
}

// The below example uses tn_rtctrl_profile_name parameter which is being deprecated in favor of the above example using tn_rtctrl_profile_dn.
resource "aci_l3_ext_subnet" "l3out_extepg_subnet02" {
  external_network_instance_profile_dn = aci_external_network_instance_profile.l3out_extepg.id
  ip                                   = "10.10.10.10/16"
  scope                                = ["import-rtctrl"]
  relation_l3ext_rs_subnet_to_profile {
    tn_rtctrl_profile_name = "test"
    direction              = "import"
  }
}