output "tenant_dn" {
  value = "${data.aci_tenant.tenant_fetch.id}"
}

output "app_dn" {
  value = "${data.aci_application_profile.ap_fetch.id}"
}

output "epg_dn" {
  value = "${data.aci_application_epg.epg_fetch.id}"
}

output "bridge_domain_dn" {
  value = "${data.aci_bridge_domain.bd_fetch.id}"
}

output "subnet_dn" {
  value = "${data.aci_subnet.subnet_fetch.id}"
}
output "contract_dn" {
  value = "${data.aci_contract.contract_fecth.id}"
}

output "subject_dn" {
  value = "${data.aci_contract_subject.subject_fetch.id}"
}

output "filter_dn" {
  value = "${data.aci_filter.fiter_fetch.id}"
}
output "entry_dn" {
  value = "${data.aci_filter_entry.fetch_entry.id}"
}

output "vrf_dn" {
  value = "${data.aci_vrf.vrf_fetch.id}"
}
output "vm_domain_dn" {
  value = "${data.aci_vmm_domain.fetch_domain.id}"
}







