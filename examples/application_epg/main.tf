terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "ansible_github_ci"
  password = "sJ94G92#8dq2hx*K4qh"
  url      = "https://173.36.219.70"
  insecure = true
}

resource "aci_tenant" "terraform_tenant" {
  name        = "static_leaf_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_application_profile" "test_ap" {
  tenant_dn   = aci_tenant.terraform_tenant.id
  name        = "test"
  description = "from terraform"
  name_alias  = "test_ap"
  prio        = "level1"
}

resource "aci_bridge_domain" "terraform_bd" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "test_bd"
}

resource "aci_rest" "rest_qos_custom_pol" {
  path       = "api/node/mo/${aci_tenant.terraform_tenant.id}/qoscustom-testpol.json"
  class_name = "qosCustomPol"

  content = {
    "name" = "testpol"
  }
}

resource "aci_contract" "rs_prov_contract" {
  tenant_dn   = aci_tenant.terraform_tenant.id
  name        = "rs_prov_contract"
  description = "This contract is created by terraform ACI provider"
  scope       = "tenant"
  target_dscp = "VA"
  prio        = "unspecified"
}

resource "aci_imported_contract" "rest_vz_cons_if" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "testcontract"
}

resource "aci_application_epg" "inherit_epg" {
  application_profile_dn = aci_application_profile.test_ap.id
  name                   = "inherit_epg"
  description            = "epg to create relation sec_inherited"
  relation_fv_rs_node_att {
    node_dn              = "topology/pod-1/node-108"
    encap                = "vlan-100"
    description          = "this is desc for static leaf"
    deployment_immediacy = "lazy"
    mode                 = "regular"
  }
}

resource "aci_rest" "rest_qos_dpp_pol" {
  path       = "api/node/mo/${aci_tenant.terraform_tenant.id}/qosdpppol-testqospol.json"
  class_name = "qosDppPol"

  content = {
    "name" = "testqospol"
  }
}

resource "aci_contract" "rs_cons_contract" {
  tenant_dn   = aci_tenant.terraform_tenant.id
  name        = "rs_cons_contract"
  description = "This contract is created by terraform ACI provider"
  scope       = "tenant"
  target_dscp = "VA"
  prio        = "unspecified"
}

resource "aci_rest" "rest_trust_ctrl_pol" {
  path       = "api/node/mo/${aci_tenant.terraform_tenant.id}/trustctrlpol-testtrustpol.json"
  class_name = "fhsTrustCtrlPol"

  content = {
    "name" = "testtrustpol"
  }
}

resource "aci_taboo_contract" "rest_taboo_con" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "testcon"
}

resource "aci_monitoring_policy" "rest_mon_epg_pol" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "testpol"
}

resource "aci_contract" "intra_epg_contract" {
  tenant_dn   = aci_tenant.terraform_tenant.id
  name        = "intra_epg_contract"
  description = "This contract is created by terraform ACI provider"
  scope       = "tenant"
  target_dscp = "VA"
  prio        = "unspecified"
}

# data "aci_fabric_path_ep" "example" {
#   pod_id  = "1"
#   node_id = "101"
#   name    = "eth1/7"
# }

resource "aci_application_epg" "fooapplication_epg" {
  application_profile_dn      = aci_application_profile.test_ap.id
  name                        = "demo_epg"
  description                 = "from terraform"
  exception_tag               = "0"
  flood_on_encap              = "disabled"
  fwd_ctrl                    = "none"
  has_mcast_source            = "no"
  is_attr_based_epg           = "no"
  match_t                     = "AtleastOne"
  name_alias                  = "alias_epg"
  pc_enf_pref                 = "unenforced"
  pref_gr_memb                = "exclude"
  prio                        = "unspecified"
  shutdown                    = "no"
  relation_fv_rs_bd           = aci_bridge_domain.terraform_bd.id
  relation_fv_rs_cust_qos_pol = aci_rest.rest_qos_custom_pol.id
  relation_fv_rs_fc_path_att = [
    # data.aci_fabric_path_ep.example.id,
    # "topology/pod-1/paths-101/pathep-[eth1/22]"
  ]
  relation_fv_rs_prov          = [aci_contract.rs_prov_contract.id]
  relation_fv_rs_cons_if       = [aci_imported_contract.rest_vz_cons_if.id]
  relation_fv_rs_sec_inherited = [aci_application_epg.inherit_epg.id]
  relation_fv_rs_dpp_pol       = aci_rest.rest_qos_dpp_pol.id
  relation_fv_rs_cons          = [aci_contract.rs_cons_contract.id]
  relation_fv_rs_trust_ctrl    = aci_rest.rest_trust_ctrl_pol.id
  relation_fv_rs_prot_by       = [aci_taboo_contract.rest_taboo_con.id]
  relation_fv_rs_aepg_mon_pol  = aci_monitoring_policy.rest_mon_epg_pol.id
  relation_fv_rs_intra_epg     = [aci_contract.intra_epg_contract.id]
  relation_fv_rs_node_att {
    node_dn              = "topology/pod-1/node-108"
    encap                = "vlan-100"
    description          = "this is desc for static leaf"
    deployment_immediacy = "lazy"
    mode                 = "regular"
  }
}

/*
The following depicts an example to create and associate an application EPG with the common Tenant's BD and VRF
*/

data "aci_tenant" "common_tenant" {
  name = "common"
}

data "aci_vrf" "default_vrf" {
  tenant_dn = data.aci_tenant.common_tenant.id
  name      = "default"
}

resource "aci_bridge_domain" "test_bd" {
  tenant_dn          = data.aci_tenant.common_tenant.id
  name               = "common_test_bd"
  relation_fv_rs_ctx = data.aci_vrf.default_vrf.id
}

resource "aci_application_epg" "test_epg_common" {
  application_profile_dn = aci_application_profile.test_ap.id
  name                   = "common_test_epg"
  relation_fv_rs_bd      = aci_bridge_domain.test_bd.id
  relation_fv_rs_node_att {
    node_dn              = "topology/pod-1/node-108"
    encap                = "vlan-100"
    description          = "this is desc for static leaf"
    deployment_immediacy = "lazy"
    mode                 = "regular"
  }
}
