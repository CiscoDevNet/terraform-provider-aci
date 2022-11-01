terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "admin"
  password = "Ins3965!12345"
  url      = "https://20.230.92.129"
  insecure = true
}

resource "aci_tenant" "terraform_tenant" {
  name = "unmanaged-tenant1"
}

resource "aci_vrf" "vrf" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "unmanaged-VRF2"
}

resource "aci_cloud_context_profile" "ctx1" {
  name                     = "cloud_context_profile"
  description              = "update desc"
  tenant_dn                = aci_tenant.terraform_tenant.id
  primary_cidr             = "10.1.0.0/16"
  region                   = "eastus2"
  cloud_vendor             = "azure"
  relation_cloud_rs_to_ctx = aci_vrf.vrf.id
  cloud_brownfield         = "/subscriptions/aafaec5f-c828-4651-8504-3a1a01c5daeb/resourceGroups/Unmanaged-test/providers/Microsoft.Network/virtualNetworks/Unmanaged-VNet3"
  access_policy_type       = "read-only"
}
