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
  name = "cloudTenant"
}

resource "aci_vrf" "vrf" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "cloudVrf"
}

resource "aci_cloud_context_profile" "ctx1" {
  name                     = "cloud_context_profile"
  description              = "cloud_context_profile created while acceptance testing"
  tenant_dn                = aci_tenant.terraform_tenant.id
  primary_cidr             = "10.0.0.0/16"
  region                   = "us-west-1"
  cloud_vendor             = "aws"
  relation_cloud_rs_to_ctx = aci_vrf.vrf.id
  hub_network              = "uni/tn-infra/gwrouterp-default"
}