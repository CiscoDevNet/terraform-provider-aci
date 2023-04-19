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

# Tenant Setup
resource "aci_tenant" "terraform_tenant" {
  name = "terraform_tenant"
}

data "aci_tenant" "common" {
  name = "common"
}

# Route Control Profile Setup
data "aci_route_control_profile" "shared_route_control_profile" {
  parent_dn = data.aci_tenant.common.id
  name      = "ok"
}

# VRF Setup
data "aci_vrf" "default_vrf" {
  tenant_dn = data.aci_tenant.common.id
  name      = "default"
}

# L3 Domain Setup
resource "aci_l3_domain_profile" "l3_domain_profile" {
  name = "l3_domain_profile"
}

# L3Outside Domain
resource "aci_l3_outside" "l3_outside" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "l3_outside"
}

# Taboo Contract Setup
resource "aci_taboo_contract" "taboo_contract" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "taboo_contract"
}

# Contract Setup
resource "aci_contract" "web_contract" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "web_contract"
}

# Imported Contract Setup
resource "aci_imported_contract" "contract_interface" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "contract_interface"
}

# Route Control Profile Setup
resource "aci_route_control_profile" "route_control_profile_export" {
  parent_dn = aci_tenant.terraform_tenant.id
  name      = "route_control_profile_export"
}

resource "aci_route_control_profile" "route_control_profile_import" {
  parent_dn = aci_tenant.terraform_tenant.id
  name      = "route_control_profile_import"
}

# External EPG Setup
resource "aci_external_network_instance_profile" "external_epg_1" {
  l3_outside_dn = aci_l3_outside.l3_outside.id
  name          = "external_epg_1"
}

resource "aci_external_network_instance_profile" "external_epg_2" {
  l3_outside_dn = aci_l3_outside.l3_outside.id
  name          = "external_epg_2"

  # Route Control Profile - Every direction(export/import) allows only one object
  relation_l3ext_rs_inst_p_to_profile {
    tn_rtctrl_profile_dn = aci_route_control_profile.route_control_profile_export.id
    direction            = "export"
  }
  relation_l3ext_rs_inst_p_to_profile {
    tn_rtctrl_profile_dn = aci_route_control_profile.route_control_profile_import.id
    direction            = "import"
  }

  relation_fv_rs_sec_inherited = [aci_external_network_instance_profile.external_epg_1.id]
  relation_fv_rs_cons_if       = [aci_imported_contract.contract_interface.id]
  relation_fv_rs_prov          = [aci_contract.web_contract.id]
  relation_fv_rs_cons          = [aci_contract.web_contract.id]
  relation_fv_rs_prot_by       = [aci_taboo_contract.taboo_contract.id]
}
