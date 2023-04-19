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

# Import Brownfield VPC in AWS cloud APIC
resource "aci_tenant" "terraform_tenant" {
  name = "tenant1"
}

resource "aci_vrf" "vrf" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "aws_vrf"
}

# AWS cloud
resource "aci_cloud_context_profile" "ctx1" {
  name                     = "cloud_context_profile"
  description              = "import brownfield vpc in aws"
  tenant_dn                = aci_tenant.terraform_tenant.id
  primary_cidr             = "10.2.0.0/24"
  region                   = "us-east-1"
  cloud_vendor             = "aws"
  relation_cloud_rs_to_ctx = aci_vrf.vrf.id
  cloud_brownfield         = "vpc-00a844d6354c53502"
  access_policy_type       = "read-only"
}

resource "aci_cloud_context_profile" "ctx2" {
  name                     = "cloud_context_profile_2"
  description              = "normal vpc in aws"
  tenant_dn                = aci_tenant.terraform_tenant.id
  primary_cidr             = "10.2.1.0/24"
  region                   = "us-west-1"
  cloud_vendor             = "aws"
  relation_cloud_rs_to_ctx = aci_vrf.vrf.id
}