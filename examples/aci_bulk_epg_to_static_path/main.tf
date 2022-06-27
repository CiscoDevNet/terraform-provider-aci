terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "admin"
  password = "ins3965!"
  url      = "https://10.23.248.120"
  insecure = true
}

resource "aci_tenant" "terraform_tenant" {
  name        = "static_path_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_application_profile" "terraform_ap" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "ap_bulk_static_path"
}

resource "aci_application_epg" "terraform_epg" {
  application_profile_dn = aci_application_profile.terraform_ap.id
  name                   = "epg_bulk_static_path"
}

resource "aci_bulk_epg_to_static_path" "example" {
  application_epg_dn = aci_application_epg.terraform_epg.id
  static_path {
    interface_dn                = "topology/pod-1/paths-129/pathep-[eth1/5]"
    encap = "vlan-1000"
    description = "this is desc for bulk static path"
    deployment_immediacy = "lazy"
    mode = "native"
    primary_encap = "vlan-700"
  }
  static_path {
    interface_dn                = "topology/pod-1/paths-129/pathep-[eth1/6]"
    encap = "vlan-1001"
    description = "this is desc for another bulk static path"
    deployment_immediacy = "immediate"
    mode = "regular"
    primary_encap = "vlan-800"
  }
}