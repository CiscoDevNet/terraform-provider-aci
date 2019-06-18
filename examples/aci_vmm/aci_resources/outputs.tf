output "aci_tenant_name" {
  value = "${aci_tenant.terraform_ten.name}"
}

output "aci_epg_name" {
  value = "${aci_application_epg.epg1.name}"
}

output "aci_application_profile_name" {
  value = "${aci_application_profile.app1.name}"
}


