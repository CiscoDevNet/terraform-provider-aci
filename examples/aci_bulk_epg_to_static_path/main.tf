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

# Retrive the switch Path Attribute (used for EPG Static Path Binding) from the Leaf Interface Policy Group
data "aci_fabric_path_ep" "vpc_example" {
  vpc     = true
  pod_id  = "1"
  node_id = "101-102"
  name    = "demo_if_pol_grp"
}

resource "aci_bulk_epg_to_static_path" "example" {
  application_epg_dn = aci_application_epg.terraform_epg.id
  static_path {
    interface_dn = "topology/pod-1/paths-129/pathep-[eth1/5]"
    encap        = "vlan-1000"
  }
  static_path {
    interface_dn         = "topology/pod-1/paths-129/pathep-[eth1/6]"
    encap                = "vlan-1001"
    description          = "this is updated desc for another bulk static path"
    deployment_immediacy = "immediate"
    mode                 = "regular"
  }
  static_path {
    interface_dn         = "topology/pod-1/paths-129/pathep-[eth1/7]"
    encap                = "vlan-1002"
    description          = "this is desc for third bulk static path"
    deployment_immediacy = "lazy"
    mode                 = "untagged"
    primary_encap        = "vlan-900"
  }
  static_path {
    interface_dn         = "topology/pod-1/paths-129/pathep-[eth1/8]"
    encap                = "vlan-1003"
    description          = "this is desc for fourth bulk static path"
    deployment_immediacy = "lazy"
    mode                 = "native"
  }
  static_path {
    interface_dn         = data.aci_fabric_path_ep.vpc_example.id
    encap                = "vlan-1003"
    description          = "this is desc for fourth bulk static path"
    deployment_immediacy = "lazy"
    mode                 = "native"
  }
}