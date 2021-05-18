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

resource "aci_tenant" "terraform_ten" {
  name = "terraform_ten"
}

resource "aci_vrf" "vrf1" {
  tenant_dn = aci_tenant.terraform_ten.id
  name      = "vrf1"
}

resource "aci_bridge_domain" "bd1" {
  tenant_dn          = aci_tenant.terraform_ten.id
  relation_fv_rs_ctx = aci_vrf.vrf1.name
  name               = "bd1"
}

resource "aci_subnet" "bd1_subnet" {
  bridge_domain_dn = aci_bridge_domain.bd1.id
  name             = "Subnet"
  ip               = var.bd_subnet
}

resource "aci_application_profile" "app1" {
  tenant_dn = aci_tenant.terraform_ten.id
  name      = "app1"
}

# resource "aci_vmm_domain" "vds" {
#   provider_profile_dn = var.provider_profile_dn
#   name                = "ESX0-leaf102"
# }

data "aci_vmm_domain" "vds" {
  provider_profile_dn = "/uni/vmmp-VMware"
  name                = "ESX0-leaf103"
}

resource "aci_application_epg" "epg1" {
  application_profile_dn = aci_application_profile.app1.id
  name                   = "epg1"
  relation_fv_rs_bd      = aci_bridge_domain.bd1.name
  relation_fv_rs_dom_att = [data.aci_vmm_domain.vds.id]
  relation_fv_rs_cons    = [aci_contract.contract_epg1_epg2.name]
}

resource "aci_application_epg" "epg2" {
  application_profile_dn = aci_application_profile.app1.id
  name                   = "epg2"
  relation_fv_rs_bd      = aci_bridge_domain.bd1.name
  relation_fv_rs_dom_att = [data.aci_vmm_domain.vds.id]
  relation_fv_rs_prov    = [aci_contract.contract_epg1_epg2.name]
}

resource "aci_contract" "contract_epg1_epg2" {
  tenant_dn = aci_tenant.terraform_ten.id
  name      = "Web"
}

resource "aci_contract_subject" "Web_subject1" {
  contract_dn                  = aci_contract.contract_epg1_epg2.id
  name                         = "Subject"
  relation_vz_rs_subj_filt_att = [aci_filter.allow_https.name, aci_filter.allow_icmp.name]
}

resource "aci_filter" "allow_https" {
  tenant_dn = aci_tenant.terraform_ten.id
  name      = "allow_https"
}
resource "aci_filter" "allow_icmp" {
  tenant_dn = aci_tenant.terraform_ten.id
  name      = "allow_icmp"
}

resource "aci_filter_entry" "https" {
  name        = "https"
  filter_dn   = aci_filter.allow_https.id
  ether_t     = "ip"
  prot        = "tcp"
  d_from_port = "https"
  d_to_port   = "https"
  stateful    = "yes"
}

resource "aci_filter_entry" "icmp" {
  name        = "icmp"
  filter_dn   = aci_filter.allow_icmp.id
  ether_t     = "ip"
  prot        = "icmp"
  stateful    = "yes"
}

