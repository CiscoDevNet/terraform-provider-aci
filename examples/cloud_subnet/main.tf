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

resource "aci_tenant" "tenentcheck" {
  description = "cloud resources"
  name        = "cloudTenant"
  annotation  = "cloudTenant"
  name_alias  = "cloud_tenant"
}

resource "aci_vrf" "vrf1" {
  tenant_dn = aci_tenant.tenentcheck.id
  name      = "vrf-1"
}

resource "aci_cloud_context_profile" "ctx1" {
  name                     = "check"
  description              = "cloud_context_profile created while acceptance testing"
  tenant_dn                = aci_tenant.tenentcheck.id
  primary_cidr             = "10.0.0.0/16"
  region                   = "us-west-1"
  cloud_vendor             = "aws"
  relation_cloud_rs_to_ctx = aci_vrf.vrf1.id
  hub_network              = "uni/tn-infra/gwrouterp-default"
}

resource "aci_cloud_cidr_pool" "cloud_cidr_pool" {
  cloud_context_profile_dn = aci_cloud_context_profile.ctx1.id
  addr                     = "10.0.0.0/16"
}

data "aci_cloud_provider_profile" "aws_prof" {
  vendor = "aws"
}

data "aci_cloud_providers_region" "aws_region" {
  cloud_provider_profile_dn = data.aci_cloud_provider_profile.aws_prof.id
  name                      = "us-west-1"
}

data "aci_cloud_availability_zone" "aws_region_availability_zone" {
  cloud_providers_region_dn = data.aci_cloud_providers_region.aws_region.id
  name                      = "default"
}

resource "aci_cloud_subnet" "cloud_subnet" {
  cloud_cidr_pool_dn = aci_cloud_cidr_pool.cloud_cidr_pool.id
  ip                 = "10.0.1.0/24"
  usage              = "gateway"
  zone               = data.aci_cloud_availability_zone.aws_region_availability_zone.id
  scope              = ["shared", "private", "public"]
}

resource "aci_cloud_subnet" "gcp_cloud_subnet" {
  cloud_cidr_pool_dn = aci_cloud_cidr_pool.cloud_cidr_pool.id
  ip                 = "10.0.2.0/24"
  usage              = "gateway"
  subnet_group_label = "subnet_group_label" # Only applicable to the GCP vendor
  scope              = ["shared", "private", "public"]
}
