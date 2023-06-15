resource "aci_tenant" "terraform_tenant" {
  name        = "tf_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_application_profile" "terraform_ap" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "tf_ap"
}

resource "aci_application_epg" "terraform_epg" {
  application_profile_dn = aci_application_profile.terraform_ap.id
  name                   = "tf_epg"
}

# Port example
resource "aci_epg_to_static_path" "port" {
  application_epg_dn = aci_application_epg.terraform_epg.id
  tdn                = "topology/pod-1/paths-129/pathep-[eth1/5]"
  encap              = "vlan-100"
  mode               = "regular"
}

# PC example
resource "aci_epg_to_static_path" "pc" {
  application_epg_dn = aci_application_epg.terraform_epg.id
  tdn                = "topology/pod-1/paths-101/pathep-[pc_PolGrp]"
  encap              = "vlan-100"
  mode               = "regular"
}

# VPC example
resource "aci_epg_to_static_path" "vpc" {
  application_epg_dn = aci_application_epg.terraform_epg.id
  tdn                = "topology/pod-1/protpaths-101-104/pathep-[vpc_PolGrp]"
  encap              = "vlan-100"
  mode               = "regular"
}
