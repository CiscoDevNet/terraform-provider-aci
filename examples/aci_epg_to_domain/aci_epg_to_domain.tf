resource "aci_tenant" "terraform_tenant" {
    name        = "tf_tenant"
    description = "This tenant is created by terraform"
}

resource "aci_application_profile" "terraform_ap" {
    tenant_dn  = aci_tenant.terraform_tenant.id
    name       = "tf_ap"
}

resource "aci_application_epg" "terraform_epg" {
    application_profile_dn  = aci_application_profile.terraform_ap.id
    name                    = "tf_epg"
}

#ENSURE DOMAIN IS BOUND TO EPG
resource "aci_epg_to_domain" "terraform_epg_domain" {
    application_epg_dn    = aci_application_epg.terraform_epg.id
    tdn                   = "uni/vmmp-VMware/dom-aci_terraform_lab"
    vmm_allow_promiscuous = "accept"
    vmm_forged_transmits  = "reject"
    vmm_mac_changes       = "accept"
}

