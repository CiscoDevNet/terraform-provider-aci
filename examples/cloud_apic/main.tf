resource "aci_tenant" "terraform_ten" {
  name = "DevNet-lab"
}

resource "aci_vrf" "vrf1" {
  tenant_dn = aci_tenant.terraform_ten.id
  name      = "devnet--vrftf"
}

resource "aci_cloud_aws_provider" "cloud_apic_provider" {
  name              = "aws"
  tenant_dn         = aci_tenant.terraform_ten.id
  access_key_id     = ""
  secret_access_key = ""
  account_id        = ""
  is_trusted        = "no"
}

data "aci_cloud_domain_profile" "cloud_domp" {
  name = "demo_domp"
}

resource "aci_cloud_applicationcontainer" "app1" {
  tenant_dn = aci_tenant.terraform_ten.id
  name      = "app1"
}

resource "aci_bridge_domain" "bd1" {
  tenant_dn          = aci_tenant.terraform_ten.id
  relation_fv_rs_ctx = aci_vrf.vrf1.name
  name               = "bd1"
}

resource "aci_cloud_epg" "cloud_apic_epg" {
  name                             = "epg1"
  cloud_applicationcontainer_dn    = aci_cloud_applicationcontainer.app1.id
  relation_fv_rs_prov              = [aci_contract.contract_epg1_epg2.name]
  relation_fv_rs_cons              = [aci_contract.contract_epg1_epg2.name]
  relation_cloud_rs_cloud_epg_ctx = aci_vrf.vrf1.name
}

resource "aci_cloud_endpoint_selector" "cloud_ep_selector" {
  cloud_epg_dn    = aci_cloud_epg.cloud_apic_epg.id
  name             = "devnet-ep-select"
  match_expression = "custom:Name=='-ep2'"
}

resource "aci_contract" "contract_epg1_epg2" {
  tenant_dn = aci_tenant.terraform_ten.id
  name      = "Web"
}

resource "aci_contract_subject" "Web_subject1" {
  contract_dn                  = aci_contract.contract_epg1_epg2.id
  name                         = "Subject"
  relation_vz_rs_subj_filt_att = [aci_filter.allow_https.id, aci_filter.allow_icmp.id]
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
  name      = "icmp"
  filter_dn = aci_filter.allow_icmp.id
  ether_t   = "ip"
  prot      = "icmp"
  stateful  = "yes"
}

resource "aci_cloud_external_epg" "cloud_epic_ext_epg" {
  cloud_applicationcontainer_dn    = aci_cloud_applicationcontainer.app1.id
  name                             = "devnet--inet"
  relation_fv_rs_prov              = [aci_contract.contract_epg1_epg2.name]
  relation_fv_rs_cons              = [aci_contract.contract_epg1_epg2.name]
  relation_cloud_rs_cloud_epg_ctx = aci_vrf.vrf1.name
}

resource "aci_cloud_endpoint_selectorfor_external_epgs" "ext_ep_selector" {
  cloud_external_epg_dn = aci_cloud_external_epg.cloud_epic_ext_epg.id
  name                   = "devnet-ext"
  subnet                 = "0.0.0.0/0"
}

resource "aci_cloud_context_profile" "context_profile" {
  name                     = "devnet--cloud-ctx-profile"
  description              = "context provider created with terraform"
  tenant_dn                = aci_tenant.terraform_ten.id
  primary_cidr             = "10.230.231.1/16"
  region                   = "us-west-1"
  relation_cloud_rs_to_ctx = aci_vrf.vrf1.id
  depends_on               = ["aci_filter_entry.icmp"]
}

data "aci_cloud_cidr_pool" "prim_cidr" {
  cloud_context_profile_dn = aci_cloud_context_profile.context_profile.id
  addr                     = "10.230.231.1/16"
  name                     = "10.230.231.1/16"
}

resource "aci_cloud_subnet" "cloud_apic_subnet" {
  cloud_cidr_pool_dn            = data.aci_cloud_cidr_pool.prim_cidr.id
  name                          = "10.230.231.1/24"
  ip                            = "10.230.231.1/24"
  relation_cloud_rs_zone_attach = "uni/clouddomp/provp-aws/region-us-west-1/zone-us-west-1a"
}

output "demo_vpc_name" {
  value = "context-[${aci_vrf.vrf1.name}]-addr-[${aci_cloud_context_profile.context_profile.primary_cidr}]"
}
